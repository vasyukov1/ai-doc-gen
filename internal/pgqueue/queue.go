package pgqueue

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/vasyukov1/ai-doc-gen/internal/domain"
	"github.com/vasyukov1/ai-doc-gen/internal/storage"
	"github.com/vasyukov1/ai-doc-gen/internal/storage/generated/ai_doc_gen/public/table"
)

type Queue struct {
	db *sql.DB
}

type Claim struct {
	ID        string
	ProfileID string
	Attempts  int
}

func New(db *sql.DB) *Queue {
	return &Queue{db: db}
}

func (q *Queue) Claim(ctx context.Context, owner string) (*Claim, error) {
	tx, err := q.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("start claim: %w", err)
	}
	defer tx.Rollback()

	jobs := table.Jobs
	available := jobs.Status.EQ(postgres.String(string(domain.StatusQueued))).
		OR(jobs.Status.EQ(postgres.String(string(domain.StatusRetryWait))).
			AND(jobs.NextRetryAt.LT_EQ(postgres.NOW()))).
		OR(jobs.Status.EQ(postgres.String(string(domain.StatusRunning))).
			AND(jobs.LeaseUntil.LT(postgres.NOW())))
	stmt := jobs.SELECT(jobs.ID, jobs.Status, jobs.ProfileID, jobs.Attempts, jobs.MaxAttempts).
		WHERE(available).
		ORDER_BY(jobs.CreatedAt).
		LIMIT(1).
		FOR(postgres.UPDATE().SKIP_LOCKED())
	query, args := stmt.Sql()

	var id string
	var status domain.Status
	var profileID string
	var attempts, maxAttempts int
	err = tx.QueryRowContext(ctx, query, args...).Scan(&id, &status, &profileID, &attempts, &maxAttempts)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("find job: %w", err)
	}

	if attempts >= maxAttempts {
		if err := storage.TransitionTx(ctx, tx, id, status, domain.StatusFailed, "", "attempts_exhausted", "Processing attempts exhausted."); err != nil {
			return nil, err
		}

		stmt := jobs.UPDATE(jobs.LeaseOwner, jobs.LeaseUntil, jobs.ErrorCode, jobs.ErrorMessage).
			SET(postgres.NULL, postgres.NULL, "attempts_exhausted", "Processing attempts exhausted.").
			WHERE(jobs.ID.EQ(postgres.CAST(postgres.Text(id)).AS_UUID()))
		query, args := stmt.Sql()
		if _, err := tx.ExecContext(ctx, query, args...); err != nil {
			return nil, fmt.Errorf("fail exhausted job: %w", err)
		}

		if err := tx.Commit(); err != nil {
			return nil, fmt.Errorf("commit exhausted job: %w", err)
		}

		return nil, nil
	}

	if err = storage.TransitionTx(ctx, tx, id, status, domain.StatusRunning, "validating", "validating", "Checking project input."); err != nil {
		return nil, err
	}

	leaseStmt := jobs.UPDATE(jobs.Attempts, jobs.LeaseOwner, jobs.LeaseUntil, jobs.LastHeartbeatAt,
		jobs.NextRetryAt, jobs.ErrorCode, jobs.ErrorMessage).
		SET(jobs.Attempts.ADD(postgres.Int(1)), owner,
			postgres.NOW().ADD(postgres.INTERVAL(30, postgres.SECOND)), postgres.NOW(),
			postgres.NULL, "", "").
		WHERE(jobs.ID.EQ(postgres.CAST(postgres.Text(id)).AS_UUID()))
	query, args = leaseStmt.Sql()
	if _, err = tx.ExecContext(ctx, query, args...); err != nil {
		return nil, fmt.Errorf("lease job: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit claim: %w", err)
	}

	return &Claim{ID: id, ProfileID: profileID, Attempts: attempts + 1}, nil
}

func (q *Queue) Heartbeat(ctx context.Context, id, owner string) error {
	jobs := table.Jobs
	stmt := jobs.UPDATE(jobs.LeaseUntil, jobs.LastHeartbeatAt, jobs.UpdatedAt).
		SET(postgres.NOW().ADD(postgres.INTERVAL(30, postgres.SECOND)), postgres.NOW(), postgres.NOW()).
		WHERE(jobs.ID.EQ(postgres.CAST(postgres.Text(id)).AS_UUID()).
			AND(jobs.LeaseOwner.EQ(postgres.String(owner))).
			AND(jobs.Status.EQ(postgres.String(string(domain.StatusRunning)))).
			AND(jobs.LeaseUntil.GT(postgres.NOW())))
	query, args := stmt.Sql()
	result, err := q.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("renew job lease: %w", err)
	}

	count, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("check job lease: %w", err)
	}
	if count != 1 {
		return errors.New("job lease lost")
	}

	return nil
}

func (q *Queue) InputReady(ctx context.Context, id, owner string) error {
	return q.finish(ctx, id, owner, domain.StatusInputReady, "input_ready", "Input saved. Source inspection is the next step.", "", "")
}

func (q *Queue) Fail(ctx context.Context, id, owner, code, message string) error {
	return q.finish(ctx, id, owner, domain.StatusFailed, code, message, code, message)
}

func (q *Queue) Retry(ctx context.Context, id, owner, code, message string, delay time.Duration) error {
	tx, err := q.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("start retry: %w", err)
	}
	defer tx.Rollback()

	if err := checkOwner(ctx, tx, id, owner); err != nil {
		return err
	}
	if err := storage.TransitionTx(ctx, tx, id, domain.StatusRunning, domain.StatusRetryWait, "", code, message); err != nil {
		return err
	}

	jobs := table.Jobs
	stmt := jobs.UPDATE(jobs.LeaseOwner, jobs.LeaseUntil, jobs.NextRetryAt,
		jobs.ErrorCode, jobs.ErrorMessage).
		SET(postgres.NULL, postgres.NULL, postgres.NOW().ADD(postgres.INTERVALd(delay)), code, message).
		WHERE(jobs.ID.EQ(postgres.CAST(postgres.Text(id)).AS_UUID()))
	query, args := stmt.Sql()
	if _, err := tx.ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("schedule retry: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit retry: %w", err)
	}

	return nil
}

func (q *Queue) finish(ctx context.Context, id, owner string, status domain.Status, code, message, errorCode, errorMessage string) error {
	tx, err := q.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("start job update: %w", err)
	}
	defer tx.Rollback()

	if err := checkOwner(ctx, tx, id, owner); err != nil {
		return err
	}
	if err := storage.TransitionTx(ctx, tx, id, domain.StatusRunning, status, "", code, message); err != nil {
		return err
	}

	jobs := table.Jobs
	stmt := jobs.UPDATE(jobs.LeaseOwner, jobs.LeaseUntil, jobs.ErrorCode, jobs.ErrorMessage).
		SET(postgres.NULL, postgres.NULL, errorCode, errorMessage).
		WHERE(jobs.ID.EQ(postgres.CAST(postgres.Text(id)).AS_UUID()))
	query, args := stmt.Sql()
	if _, err := tx.ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("finish job: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit job update: %w", err)
	}

	return nil
}

func checkOwner(ctx context.Context, tx *sql.Tx, id, owner string) error {
	jobs := table.Jobs
	stmt := jobs.SELECT(jobs.LeaseOwner).
		WHERE(jobs.ID.EQ(postgres.CAST(postgres.Text(id)).AS_UUID()).
			AND(jobs.Status.EQ(postgres.String(string(domain.StatusRunning)))).
			AND(jobs.LeaseUntil.GT(postgres.NOW()))).
		FOR(postgres.UPDATE())
	query, args := stmt.Sql()

	var currentOwner string
	if err := tx.QueryRowContext(ctx, query, args...).Scan(&currentOwner); err != nil {
		return fmt.Errorf("check job owner: %w", err)
	}
	if currentOwner != owner {
		return errors.New("job lease lost")
	}

	return nil
}
