package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/vasyukov1/ai-doc-gen/internal/domain"
	"github.com/vasyukov1/ai-doc-gen/internal/storage/generated/ai_doc_gen/public/table"
)

// ErrReuseCommitUnknown preserves copied files when commit acknowledgement fails.
var ErrReuseCommitUnknown = errors.New("reused job commit outcome unknown")

// CreateReusedJob checks the parent again after source copying and atomically
// publishes the new input, inventory, credentials and queued job.
func (r *Repository) CreateReusedJob(ctx context.Context, parentID string, job domain.NewJob, inventory domain.Inventory, master string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("start reuse: %w", err)
	}
	defer tx.Rollback()

	t := table.Jobs
	query, args := t.SELECT(t.ID).WHERE(t.ID.EQ(jobUUID(parentID)).
		AND(t.OwnerHash.EQ(postgres.String(job.OwnerHash))).
		AND(t.ExpiresAt.GT(postgres.NOW())).AND(t.PurgeAfter.IS_NULL())).FOR(postgres.UPDATE()).Sql()
	var found string
	if job.OwnerHash == "" {
		return sql.ErrNoRows
	}
	if err := tx.QueryRowContext(ctx, query, args...).Scan(&found); err != nil {
		return err
	}
	if inventory.Source.JobID != job.ID {
		return errors.New("reused inventory job mismatch")
	}

	if err := createJobTx(ctx, tx, job, master); err != nil {
		return err
	}
	if err := SaveInventoryTx(ctx, tx, inventory); err != nil {
		return err
	}
	query, args = t.UPDATE(t.Stage).SET("analyzing").WHERE(t.ID.EQ(jobUUID(job.ID))).Sql()
	if _, err := tx.ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("set reuse stage: %w", err)
	}
	if err := AddProgressTx(ctx, tx, job.ID, "source_reused", "Retained project copied. New evidence and document generation are queued."); err != nil {
		return err
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("%w: %v", ErrReuseCommitUnknown, err)
	}

	return nil
}
