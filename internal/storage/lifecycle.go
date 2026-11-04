package storage

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/vasyukov1/ai-doc-gen/internal/domain"
	"github.com/vasyukov1/ai-doc-gen/internal/storage/generated/ai_doc_gen/public/table"
)

func (r *Repository) OwnsJob(ctx context.Context, id, owner string) bool {
	t := table.Jobs
	query, args := t.SELECT(t.ID).WHERE(t.ID.EQ(postgres.CAST(postgres.Text(id)).AS_UUID()).AND(t.OwnerHash.EQ(postgres.String(owner))).AND(t.ExpiresAt.GT(postgres.NOW())).AND(t.PurgeAfter.IS_NULL())).Sql()
	var found string
	return owner != "" && r.db.QueryRowContext(ctx, query, args...).Scan(&found) == nil
}

func (r *Repository) OwnJobs(ctx context.Context, owner string) ([]domain.Job, error) {
	t := table.Jobs
	query, args := t.SELECT(t.ID, t.Status, t.ProfileID, t.ExpiresAt).WHERE(t.OwnerHash.EQ(postgres.String(owner)).AND(t.ExpiresAt.GT(postgres.NOW())).AND(t.PurgeAfter.IS_NULL())).ORDER_BY(t.CreatedAt.DESC()).LIMIT(100).Sql()
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	jobs := []domain.Job{}
	for rows.Next() {
		var job domain.Job
		if err := rows.Scan(&job.ID, &job.Status, &job.ProfileID, &job.ExpiresAt); err != nil {
			return nil, err
		}
		jobs = append(jobs, job)
	}
	return jobs, rows.Err()
}

func (r *Repository) Lifecycle(ctx context.Context, id, action string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	t := table.Jobs
	query, args := t.SELECT(t.Status, t.Stage, t.PurgeAfter, t.ExpiresAt, t.ResumeStatus).WHERE(t.ID.EQ(postgres.CAST(postgres.Text(id)).AS_UUID())).FOR(postgres.UPDATE()).Sql()
	var status domain.Status
	var stage string
	var purge, expires sql.NullTime
	var prior sql.NullString
	if err = tx.QueryRowContext(ctx, query, args...).Scan(&status, &stage, &purge, &expires, &prior); err != nil {
		return err
	}
	if purge.Valid || (expires.Valid && time.Now().After(expires.Time)) {
		return errors.New("job expired")
	}
	switch action {
	case "cancel", "delete":
		if action == "cancel" && (status == domain.StatusCompleted || status == domain.StatusFailed || status == domain.StatusCancelled) {
			return nil
		}
		if status != domain.StatusCancelled {
			query, args = t.UPDATE(t.ResumeStatus).SET(string(status)).WHERE(t.ID.EQ(postgres.CAST(postgres.Text(id)).AS_UUID())).Sql()
			if _, err = tx.ExecContext(ctx, query, args...); err != nil {
				return err
			}
			if err = TransitionTx(ctx, tx, id, status, domain.StatusCancelled, stage, "cancelled", "Cancelled. Saved work is retained until expiry or deletion."); err != nil {
				return err
			}
		}
		query, args = t.UPDATE(t.LeaseOwner, t.LeaseUntil, t.NextRetryAt, t.RequiredAction).SET(postgres.NULL, postgres.NULL, postgres.NULL, "").WHERE(t.ID.EQ(postgres.CAST(postgres.Text(id)).AS_UUID())).Sql()
		if _, err = tx.ExecContext(ctx, query, args...); err != nil {
			return err
		}
		if action == "delete" {
			query, args = t.UPDATE(t.PurgeAfter).SET(postgres.NOW().ADD(postgres.INTERVAL(90, postgres.SECOND))).WHERE(t.ID.EQ(postgres.CAST(postgres.Text(id)).AS_UUID())).Sql()
			if _, err = tx.ExecContext(ctx, query, args...); err != nil {
				return err
			}
		}
	case "retry":
		if status != domain.StatusFailed && status != domain.StatusCancelled && status != domain.StatusRetryWait {
			return errors.New("retry unavailable")
		}
		resume := domain.StatusQueued
		if status == domain.StatusCancelled && (prior.String == string(domain.StatusNeedsInput) || prior.String == string(domain.StatusReadyForReview)) {
			resume = domain.Status(prior.String)
		}
		if stage == "selecting_files" || stage == "reviewing_evidence" {
			resume = domain.StatusNeedsInput
		}
		if stage == "draft_ready" || stage == "reviewing_document" {
			resume = domain.StatusReadyForReview
		}
		if err = TransitionTx(ctx, tx, id, status, resume, stage, "retry_queued", "Retry queued. Completed drafts and review decisions are preserved."); err != nil {
			return err
		}
		query, args = t.UPDATE(t.Attempts, t.LeaseOwner, t.LeaseUntil, t.NextRetryAt, t.ErrorCode, t.ErrorMessage).SET(0, postgres.NULL, postgres.NULL, postgres.NULL, "", "").WHERE(t.ID.EQ(postgres.CAST(postgres.Text(id)).AS_UUID())).Sql()
		if _, err = tx.ExecContext(ctx, query, args...); err != nil {
			return err
		}
	default:
		return errors.New("invalid action")
	}
	return tx.Commit()
}
