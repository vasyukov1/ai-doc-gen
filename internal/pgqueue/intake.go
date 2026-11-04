package pgqueue

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/vasyukov1/ai-doc-gen/internal/domain"
	"github.com/vasyukov1/ai-doc-gen/internal/storage"
	"github.com/vasyukov1/ai-doc-gen/internal/storage/generated/ai_doc_gen/public/table"
)

var ErrCommitUnknown = errors.New("inventory commit outcome unknown")

func (q *Queue) Stage(ctx context.Context, id, owner, stage, code, message string) error {
	tx, err := q.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("start stage update: %w", err)
	}
	defer tx.Rollback()

	if err := checkOwner(ctx, tx, id, owner); err != nil {
		return err
	}

	jobs := table.Jobs
	stmt := jobs.UPDATE(jobs.Stage, jobs.UpdatedAt).SET(stage, postgres.NOW()).
		WHERE(jobs.ID.EQ(postgres.CAST(postgres.Text(id)).AS_UUID()))
	query, args := stmt.Sql()
	if _, err := tx.ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("update job stage: %w", err)
	}
	if err := storage.AddProgressTx(ctx, tx, id, code, message); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit stage update: %w", err)
	}

	return nil
}

func (q *Queue) InventoryReady(ctx context.Context, id, owner string, inventory domain.Inventory) error {
	tx, err := q.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("start inventory update: %w", err)
	}
	defer tx.Rollback()

	if err := checkOwner(ctx, tx, id, owner); err != nil {
		return err
	}
	if err := storage.SaveInventoryTx(ctx, tx, inventory); err != nil {
		return err
	}

	message := "Review the file list and confirm the selection."
	action := "Confirm the file selection."
	hasEligible := false
	for _, file := range inventory.Files {
		if file.Eligible {
			hasEligible = true
			break
		}
	}
	if !hasEligible {
		message = "No selectable text files were found. Review exclusions and submit another project."
		action = "Submit a project with usable text files."
	}
	if err := storage.TransitionTx(ctx, tx, id, domain.StatusRunning, domain.StatusNeedsInput,
		"selecting_files", "files_ready", message); err != nil {
		return err
	}

	jobs := table.Jobs
	stmt := jobs.UPDATE(jobs.LeaseOwner, jobs.LeaseUntil, jobs.RequiredAction).
		SET(postgres.NULL, postgres.NULL, action).
		WHERE(jobs.ID.EQ(postgres.CAST(postgres.Text(id)).AS_UUID()))
	query, args := stmt.Sql()
	if _, err := tx.ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("finish inventory: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit inventory: %w", errors.Join(ErrCommitUnknown, err))
	}

	return nil
}
