package storage

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"time"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/vasyukov1/ai-doc-gen/internal/domain"
	"github.com/vasyukov1/ai-doc-gen/internal/storage/generated/ai_doc_gen/public/table"
)

func jobUUID(id string) postgres.StringExpression { return postgres.CAST(postgres.Text(id)).AS_UUID() }

func saveDraftRunTx(ctx context.Context, tx *sql.Tx, id string, connections []domain.Connection, master string) error {
	aead, err := credentialCipher(master)
	if err != nil {
		return err
	}
	nonce := make([]byte, aead.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return err
	}
	raw, err := json.Marshal(connections)
	if err != nil {
		return err
	}
	encrypted := base64.StdEncoding.EncodeToString(aead.Seal(nonce, nonce, raw, []byte(id+"/draft")))
	t := table.JobDraftRuns
	query, args := t.INSERT(t.JobID, t.EncryptedConnections).VALUES(id, encrypted).Sql()
	_, err = tx.ExecContext(ctx, query, args...)
	return err
}

func (r *Repository) DraftRun(ctx context.Context, id, master string) (domain.DraftRun, error) {
	var run domain.DraftRun
	var encrypted string
	var disabled, attempts []byte
	t := table.JobDraftRuns
	query, args := t.SELECT(t.EncryptedConnections, t.Deadline, t.Attempts, t.Reason, t.Disabled, t.SectionAttempts).WHERE(t.JobID.EQ(jobUUID(id))).Sql()
	if err := r.db.QueryRowContext(ctx, query, args...).Scan(&encrypted, &run.Deadline, &run.Attempts, &run.Reason, &disabled, &attempts); err != nil {
		return run, err
	}
	if err := json.Unmarshal(disabled, &run.Disabled); err != nil {
		return run, err
	}
	if err := json.Unmarshal(attempts, &run.SectionAttempts); err != nil {
		return run, err
	}
	aead, err := credentialCipher(master)
	if err != nil {
		return run, err
	}
	raw, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil || len(raw) < aead.NonceSize() {
		return run, errors.New("invalid saved connections")
	}
	raw, err = aead.Open(nil, raw[:aead.NonceSize()], raw[aead.NonceSize():], []byte(id+"/draft"))
	if err != nil {
		return run, err
	}
	err = json.Unmarshal(raw, &run.Connections)
	return run, err
}

func (r *Repository) draftSummary(ctx context.Context, job *domain.Job) error {
	t := table.JobDraftRuns
	query, args := t.SELECT(t.Deadline, t.Attempts, t.Reason).WHERE(t.JobID.EQ(jobUUID(job.ID))).Sql()
	var deadline time.Time
	err := r.db.QueryRowContext(ctx, query, args...).Scan(&deadline, &job.GenerationAttempts, &job.GenerationReason)
	if errors.Is(err, sql.ErrNoRows) {
		return nil
	}
	if err != nil {
		return err
	}
	job.FreeDraft = true
	job.GenerationDeadline = &deadline
	sections, err := r.Sections(ctx, job.ID)
	job.Content = domain.ContentReport(sections, job.TotalSections)
	for _, s := range sections {
		if s.Origin == "model" {
			job.ModelSections++
		}
	}
	return err
}

// ResetDraftBudgetTx is called only for an explicit section regeneration.
func ResetDraftBudgetTx(ctx context.Context, tx *sql.Tx, id, sectionID string) error {
	sections := table.JobSections
	query, args := sections.SELECT(sections.SectionID).WHERE(sections.JobID.EQ(jobUUID(id))).Sql()
	rows, err := tx.QueryContext(ctx, query, args...)
	if err != nil {
		return err
	}
	counts := map[string]int{}
	for rows.Next() {
		var existingID string
		if err = rows.Scan(&existingID); err != nil {
			rows.Close()
			return err
		}
		if existingID != sectionID {
			counts[existingID] = 2
		}
	}
	err = rows.Err()
	rows.Close()
	if err != nil {
		return err
	}
	raw, _ := json.Marshal(counts)
	t := table.JobDraftRuns
	query, args = t.UPDATE(t.Deadline, t.Attempts, t.Reason, t.Disabled, t.SectionAttempts).SET(postgres.NOW().ADD(postgres.INTERVAL(240, postgres.SECOND)), 0, "", "{}", string(raw)).WHERE(t.JobID.EQ(jobUUID(id))).Sql()
	_, err = tx.ExecContext(ctx, query, args...)
	return err
}
