package pgqueue

import (
	"context"
	"github.com/go-jet/jet/v2/postgres"
	"github.com/vasyukov1/ai-doc-gen/internal/domain"
	"github.com/vasyukov1/ai-doc-gen/internal/storage"
	"github.com/vasyukov1/ai-doc-gen/internal/storage/generated/ai_doc_gen/public/table"
)

func (q *Queue) Rendered(ctx context.Context, id, owner string, artifacts []domain.Artifact) error {
	tx, err := q.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if err = checkOwner(ctx, tx, id, owner); err != nil {
		return err
	}
	t := table.JobArtifacts
	query, args := t.DELETE().WHERE(t.JobID.EQ(postgres.CAST(postgres.Text(id)).AS_UUID())).Sql()
	if _, err = tx.ExecContext(ctx, query, args...); err != nil {
		return err
	}
	for _, artifact := range artifacts {
		query, args = t.INSERT(t.JobID, t.Kind, t.StoragePath, t.Sha256).VALUES(id, artifact.Kind, artifact.Path, artifact.SHA256).Sql()
		if _, err = tx.ExecContext(ctx, query, args...); err != nil {
			return err
		}
	}
	if err = storage.TransitionTx(ctx, tx, id, domain.StatusRunning, domain.StatusCompleted, "rendered", "rendered", "PDF and reproducible Typst source are ready."); err != nil {
		return err
	}
	jobs := table.Jobs
	query, args = jobs.UPDATE(jobs.LeaseOwner, jobs.LeaseUntil, jobs.RequiredAction, jobs.ErrorCode, jobs.ErrorMessage).SET(postgres.NULL, postgres.NULL, "", "", "").WHERE(jobs.ID.EQ(postgres.CAST(postgres.Text(id)).AS_UUID())).Sql()
	if _, err = tx.ExecContext(ctx, query, args...); err != nil {
		return err
	}
	return tx.Commit()
}
