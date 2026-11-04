package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/vasyukov1/ai-doc-gen/internal/domain"
	"github.com/vasyukov1/ai-doc-gen/internal/statemachine"
	"github.com/vasyukov1/ai-doc-gen/internal/storage/generated/ai_doc_gen/public/table"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Ping(ctx context.Context) error {
	return r.db.PingContext(ctx)
}

func (r *Repository) CreateJob(ctx context.Context, job domain.NewJob) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("start job: %w", err)
	}
	defer tx.Rollback()

	jobs := table.Jobs
	stmt := jobs.INSERT(jobs.ID, jobs.Status, jobs.ProfileID, jobs.Language, jobs.Description).
		VALUES(job.ID, domain.StatusQueued, job.ProfileID, job.Language, job.Description)
	query, args := stmt.Sql()
	if _, err := tx.ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("save job: %w", err)
	}

	inputs := table.JobInputs
	stmt = inputs.INSERT(inputs.JobID, inputs.SourceType, inputs.GithubURL, inputs.UploadPath,
		inputs.OriginalName, inputs.Sha256, inputs.SizeBytes).
		VALUES(job.ID, job.Input.SourceType, nullable(job.Input.GitHubURL), nullable(job.Input.UploadPath),
			job.Input.OriginalName, job.Input.SHA256, job.Input.SizeBytes)
	query, args = stmt.Sql()
	if _, err := tx.ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("save input: %w", err)
	}

	if err := addStatus(ctx, tx, job.ID, "", domain.StatusQueued, "", "queued", "Job queued."); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit job: %w", err)
	}

	return nil
}

func (r *Repository) GetJob(ctx context.Context, id string) (domain.Job, error) {
	jobs := table.Jobs
	stmt := jobs.SELECT(jobs.ID, jobs.Status, jobs.Stage, jobs.ProfileID, jobs.Language,
		jobs.CreatedAt, jobs.UpdatedAt, jobs.LastHeartbeatAt, jobs.LeaseUntil, jobs.NextRetryAt,
		jobs.ErrorCode, jobs.ErrorMessage, jobs.RequiredAction, jobs.CompletedSections, jobs.TotalSections).
		WHERE(jobs.ID.EQ(postgres.CAST(postgres.Text(id)).AS_UUID()))
	query, args := stmt.Sql()

	var job domain.Job
	var heartbeat, leaseUntil, retryAt sql.NullTime
	err := r.db.QueryRowContext(ctx, query, args...).Scan(&job.ID, &job.Status, &job.Stage,
		&job.ProfileID, &job.Language, &job.CreatedAt, &job.UpdatedAt, &heartbeat, &leaseUntil,
		&retryAt, &job.ErrorCode, &job.ErrorMessage, &job.RequiredAction, &job.CompletedSections,
		&job.TotalSections)
	if err != nil {
		return domain.Job{}, fmt.Errorf("get job: %w", err)
	}

	events := table.JobProgressEvents
	stmt = events.SELECT(events.Message).
		WHERE(events.JobID.EQ(postgres.CAST(postgres.Text(id)).AS_UUID())).
		ORDER_BY(events.ID.DESC()).LIMIT(1)
	query, args = stmt.Sql()
	err = r.db.QueryRowContext(ctx, query, args...).Scan(&job.LatestMessage)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return domain.Job{}, fmt.Errorf("get job progress: %w", err)
	}

	if heartbeat.Valid {
		job.LastHeartbeatAt = &heartbeat.Time
	}
	if retryAt.Valid {
		job.NextRetryAt = &retryAt.Time
	}
	if job.Status == domain.StatusRunning && leaseUntil.Valid && leaseUntil.Time.Before(time.Now()) {
		job.RecoveryPending = true
		job.LatestMessage = "Worker stopped. Waiting for recovery."
	}

	return job, nil
}

func (r *Repository) GetInput(ctx context.Context, id string) (domain.Input, error) {
	inputs := table.JobInputs
	stmt := inputs.SELECT(inputs.SourceType, inputs.GithubURL, inputs.UploadPath,
		inputs.OriginalName, inputs.Sha256, inputs.SizeBytes).
		WHERE(inputs.JobID.EQ(postgres.CAST(postgres.Text(id)).AS_UUID()))
	query, args := stmt.Sql()

	var input domain.Input
	var githubURL, uploadPath sql.NullString
	err := r.db.QueryRowContext(ctx, query, args...).Scan(&input.SourceType, &githubURL,
		&uploadPath, &input.OriginalName, &input.SHA256, &input.SizeBytes)
	if err != nil {
		return domain.Input{}, fmt.Errorf("get input: %w", err)
	}

	input.GitHubURL = githubURL.String
	input.UploadPath = uploadPath.String

	return input, nil
}

func TransitionTx(ctx context.Context, tx *sql.Tx, id string, from, to domain.Status, stage, code, message string) error {
	if !statemachine.CanTransition(from, to) {
		return fmt.Errorf("invalid job transition: %s to %s", from, to)
	}

	jobs := table.Jobs
	stmt := jobs.UPDATE(jobs.Status, jobs.Stage, jobs.UpdatedAt).
		SET(to, stage, postgres.NOW()).
		WHERE(jobs.ID.EQ(postgres.CAST(postgres.Text(id)).AS_UUID()).AND(jobs.Status.EQ(postgres.String(string(from)))))
	query, args := stmt.Sql()
	result, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("update job: %w", err)
	}

	count, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("check job update: %w", err)
	}
	if count != 1 {
		return errors.New("job state changed")
	}

	return addStatus(ctx, tx, id, from, to, stage, code, message)
}

func addStatus(ctx context.Context, tx *sql.Tx, id string, from, to domain.Status, stage, code, message string) error {
	history := table.JobStatusHistory
	stmt := history.INSERT(history.JobID, history.OldStatus, history.NewStatus, history.Stage).
		VALUES(id, nullable(string(from)), to, stage)
	query, args := stmt.Sql()
	if _, err := tx.ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("save job status: %w", err)
	}

	return AddProgressTx(ctx, tx, id, code, message)
}

func AddProgressTx(ctx context.Context, tx *sql.Tx, id, code, message string) error {
	events := table.JobProgressEvents
	stmt := events.INSERT(events.JobID, events.Code, events.Message).VALUES(id, code, message)
	query, args := stmt.Sql()
	if _, err := tx.ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("save job event: %w", err)
	}

	return nil
}

func nullable(value string) any {
	if value == "" {
		return nil
	}

	return value
}
