package storage

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/vasyukov1/ai-doc-gen/internal/domain"
	"github.com/vasyukov1/ai-doc-gen/internal/storage/generated/ai_doc_gen/public/table"
)

// ErrEncryptionKey marks a missing or malformed server master key.
var ErrEncryptionKey = errors.New("server encryption key is not configured")

func credentialCipher(master string) (cipher.AEAD, error) {
	key, err := base64.StdEncoding.DecodeString(master)
	if err != nil || len(key) != 32 {
		return nil, fmt.Errorf("%w: LLM_ENCRYPTION_KEY must be a base64 encoded 32 byte key", ErrEncryptionKey)
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	return cipher.NewGCM(block)
}

func saveCredentialsTx(ctx context.Context, tx *sql.Tx, id, model, key, master string) error {
	aead, err := credentialCipher(master)
	if err != nil {
		return err
	}
	nonce := make([]byte, aead.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return err
	}
	encrypted := base64.StdEncoding.EncodeToString(aead.Seal(nonce, nonce, []byte(key), []byte(id)))
	t := table.JobGenerationCredentials
	query, args := t.INSERT(t.JobID, t.Model, t.EncryptedKey).VALUES(id, model, encrypted).ON_CONFLICT(t.JobID).DO_NOTHING().Sql()
	_, err = tx.ExecContext(ctx, query, args...)
	return err
}

func (r *Repository) Credentials(ctx context.Context, id, master string) (string, string, error) {
	t := table.JobGenerationCredentials
	query, args := t.SELECT(t.Model, t.EncryptedKey).WHERE(t.JobID.EQ(postgres.CAST(postgres.Text(id)).AS_UUID())).Sql()
	var model, encoded string
	if err := r.db.QueryRowContext(ctx, query, args...).Scan(&model, &encoded); err != nil {
		return "", "", err
	}
	aead, err := credentialCipher(master)
	if err != nil {
		return "", "", err
	}
	raw, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil || len(raw) < aead.NonceSize() {
		return "", "", errors.New("invalid saved credential")
	}
	key, err := aead.Open(nil, raw[:aead.NonceSize()], raw[aead.NonceSize():], []byte(id))
	return model, string(key), err
}

func CredentialScope(key string) string {
	digest := sha256.Sum256([]byte(key))
	return hex.EncodeToString(digest[:])
}

func (r *Repository) RetryGeneration(ctx context.Context, id, key, master string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	jobs := table.Jobs
	query, args := jobs.SELECT(jobs.Status).WHERE(jobs.ID.EQ(postgres.CAST(postgres.Text(id)).AS_UUID())).FOR(postgres.UPDATE()).Sql()
	var status domain.Status
	if err := tx.QueryRowContext(ctx, query, args...).Scan(&status); err != nil {
		return err
	}
	if status != domain.StatusFailed {
		return ErrSectionConflict
	}
	inputs := table.JobDocumentInputs
	query, args = inputs.SELECT(inputs.Reviewed).WHERE(inputs.JobID.EQ(postgres.CAST(postgres.Text(id)).AS_UUID())).Sql()
	var rendering bool
	if err := tx.QueryRowContext(ctx, query, args...).Scan(&rendering); err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}
	if rendering {
		return ErrSectionConflict
	}
	credentials := table.JobGenerationCredentials
	query, args = credentials.SELECT(credentials.Model).WHERE(credentials.JobID.EQ(postgres.CAST(postgres.Text(id)).AS_UUID())).Sql()
	var model string
	if err := tx.QueryRowContext(ctx, query, args...).Scan(&model); err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return err
		}
		model = "deepseek/deepseek-chat-v3-0324"
	}
	query, args = credentials.DELETE().WHERE(credentials.JobID.EQ(postgres.CAST(postgres.Text(id)).AS_UUID())).Sql()
	if _, err := tx.ExecContext(ctx, query, args...); err != nil {
		return err
	}
	if err := saveCredentialsTx(ctx, tx, id, model, key, master); err != nil {
		return err
	}
	if err := TransitionTx(ctx, tx, id, status, domain.StatusQueued, "generating", "generation_retry", "Generation retry queued. Saved sections are preserved."); err != nil {
		return err
	}
	query, args = jobs.UPDATE(jobs.Attempts, jobs.ErrorCode, jobs.ErrorMessage, jobs.NextRetryAt).SET(0, "", "", postgres.NULL).WHERE(jobs.ID.EQ(postgres.CAST(postgres.Text(id)).AS_UUID())).Sql()
	if _, err := tx.ExecContext(ctx, query, args...); err != nil {
		return err
	}
	return tx.Commit()
}
