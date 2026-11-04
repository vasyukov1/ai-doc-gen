package storage

import (
	"context"
	"crypto/sha256"
	"crypto/subtle"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/vasyukov1/ai-doc-gen/internal/domain"
	"github.com/vasyukov1/ai-doc-gen/internal/storage/generated/ai_doc_gen/public/table"
)

var ErrDocumentConflict = errors.New("document input is not editable")

func (r *Repository) DocumentInput(ctx context.Context, id string) (domain.DocumentInput, error) {
	t := table.JobDocumentInputs
	query, args := t.SELECT(t.Record, t.Reviewed, t.CapabilityHash).WHERE(t.JobID.EQ(postgres.CAST(postgres.Text(id)).AS_UUID())).Sql()
	var input domain.DocumentInput
	var raw []byte
	var reviewed bool
	var hash string
	if err := r.db.QueryRowContext(ctx, query, args...).Scan(&raw, &reviewed, &hash); err != nil {
		return input, err
	}
	err := json.Unmarshal(raw, &input)
	input.Reviewed = reviewed
	input.CapabilityHash = hash
	return input, err
}

func (r *Repository) SaveDocumentInput(ctx context.Context, id string, input domain.DocumentInput) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	jobs := table.Jobs
	query, args := jobs.SELECT(jobs.Status, jobs.CompletedSections).WHERE(jobs.ID.EQ(postgres.CAST(postgres.Text(id)).AS_UUID())).FOR(postgres.UPDATE()).Sql()
	var status domain.Status
	var count int
	if err = tx.QueryRowContext(ctx, query, args...).Scan(&status, &count); err != nil {
		return err
	}
	if status != domain.StatusNeedsInput || count != 0 {
		return ErrDocumentConflict
	}
	raw, err := json.Marshal(input)
	if err != nil {
		return err
	}
	t := table.JobDocumentInputs
	query, args = t.INSERT(t.JobID, t.Record, t.Reviewed).VALUES(id, string(raw), false).ON_CONFLICT(t.JobID).DO_UPDATE(postgres.SET(t.Record.SET(t.EXCLUDED.Record), t.Reviewed.SET(postgres.Bool(false)))).Sql()
	if _, err = tx.ExecContext(ctx, query, args...); err != nil {
		return err
	}
	return tx.Commit()
}

func (r *Repository) QueueRender(ctx context.Context, id, hash string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	jobs := table.Jobs
	query, args := jobs.SELECT(jobs.Status, jobs.CompletedSections, jobs.TotalSections).WHERE(jobs.ID.EQ(postgres.CAST(postgres.Text(id)).AS_UUID())).FOR(postgres.UPDATE()).Sql()
	var status domain.Status
	var count, total int
	if err = tx.QueryRowContext(ctx, query, args...).Scan(&status, &count, &total); err != nil {
		return err
	}
	if (status != domain.StatusReadyForReview && status != domain.StatusFailed) || count != total || total == 0 {
		return ErrDocumentConflict
	}
	t := table.JobDocumentInputs
	query, args = t.UPDATE(t.Reviewed, t.CapabilityHash).SET(true, hash).WHERE(t.JobID.EQ(postgres.CAST(postgres.Text(id)).AS_UUID())).Sql()
	result, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	n, _ := result.RowsAffected()
	if n != 1 {
		return ErrDocumentConflict
	}
	if err = TransitionTx(ctx, tx, id, status, domain.StatusQueued, "rendering", "render_queued", "Reviewed sections queued for Typst rendering."); err != nil {
		return err
	}
	query, args = jobs.UPDATE(jobs.Attempts, jobs.ErrorCode, jobs.ErrorMessage).SET(0, "", "").WHERE(jobs.ID.EQ(postgres.CAST(postgres.Text(id)).AS_UUID())).Sql()
	if _, err = tx.ExecContext(ctx, query, args...); err != nil {
		return err
	}
	return tx.Commit()
}

func (r *Repository) Artifact(ctx context.Context, id, kind, capability string) (domain.Artifact, error) {
	input, err := r.DocumentInput(ctx, id)
	if err != nil {
		return domain.Artifact{}, err
	}
	hash := sha256.Sum256([]byte(capability))
	if input.CapabilityHash == "" || subtle.ConstantTimeCompare([]byte(input.CapabilityHash), []byte(hex.EncodeToString(hash[:]))) != 1 {
		return domain.Artifact{}, sql.ErrNoRows
	}
	t := table.JobArtifacts
	query, args := t.SELECT(t.Kind, t.StoragePath, t.Sha256).WHERE(t.JobID.EQ(postgres.CAST(postgres.Text(id)).AS_UUID()).AND(t.Kind.EQ(postgres.String(kind)))).ORDER_BY(t.ID.DESC()).LIMIT(1).Sql()
	var artifact domain.Artifact
	err = r.db.QueryRowContext(ctx, query, args...).Scan(&artifact.Kind, &artifact.Path, &artifact.SHA256)
	return artifact, err
}
