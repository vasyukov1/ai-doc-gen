package pgqueue

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/vasyukov1/ai-doc-gen/internal/domain"
	"github.com/vasyukov1/ai-doc-gen/internal/storage"
	"github.com/vasyukov1/ai-doc-gen/internal/storage/generated/ai_doc_gen/public/table"
)

var ErrDraftBudget = errors.New("draft budget exhausted")

func draftID(id string) postgres.StringExpression { return postgres.CAST(postgres.Text(id)).AS_UUID() }

// ReserveDraftCall serializes admission across workers and records a crash-safe attempt.
// The caller supplies the attempt limit because it scales with the document profile.
func (q *Queue) ReserveDraftCall(ctx context.Context, id, owner, scope string, c domain.Connection, tokens int, sections []string, attemptLimit int) (int64, time.Duration, error) {
	tx, err := q.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, 0, err
	}
	defer tx.Rollback()
	if _, err = tx.ExecContext(ctx, "SELECT pg_advisory_xact_lock(721849)"); err != nil {
		return 0, 0, err
	}
	if err = checkOwner(ctx, tx, id, owner); err != nil {
		return 0, 0, err
	}
	runs := table.JobDraftRuns
	query, args := runs.SELECT(runs.Deadline, runs.Attempts, runs.SectionAttempts).WHERE(runs.JobID.EQ(draftID(id))).FOR(postgres.UPDATE()).Sql()
	var deadline time.Time
	var attempts int
	var raw []byte
	if err = tx.QueryRowContext(ctx, query, args...).Scan(&deadline, &attempts, &raw); err != nil {
		return 0, 0, err
	}
	now := time.Now()
	if attemptLimit < 1 {
		attemptLimit = 1
	}
	if !now.Before(deadline) || attempts >= attemptLimit {
		return 0, 0, ErrDraftBudget
	}
	counts := map[string]int{}
	if err = json.Unmarshal(raw, &counts); err != nil {
		return 0, 0, err
	}
	for _, section := range sections {
		if counts[section] >= 2 {
			return 0, 0, ErrDraftBudget
		}
		counts[section]++
	}
	limits := table.ProviderLimits
	query, args = limits.INSERT(limits.Scope, limits.Provider, limits.WindowStart, limits.NextAt, limits.ActiveUntil).VALUES(scope, c.Provider, now, now.Add(-time.Second), now.Add(-time.Second)).ON_CONFLICT(limits.Scope).DO_NOTHING().Sql()
	if _, err = tx.ExecContext(ctx, query, args...); err != nil {
		return 0, 0, err
	}
	query, args = limits.SELECT(limits.Scope, limits.Provider, limits.WindowStart, limits.Requests, limits.Tokens, limits.NextAt, limits.ActiveUntil).FOR(postgres.UPDATE()).Sql()
	rows, err := tx.QueryContext(ctx, query, args...)
	if err != nil {
		return 0, 0, err
	}
	active := 0
	busy := false
	wait := time.Duration(0)
	window := now
	requests, spent := 0, 0
	for rows.Next() {
		var key, provider string
		var start, next, until time.Time
		var req, tok int
		if err = rows.Scan(&key, &provider, &start, &req, &tok, &next, &until); err != nil {
			rows.Close()
			return 0, 0, err
		}
		if until.After(now) {
			active++
			if provider == c.Provider {
				busy = true
			}
		}
		if key == scope {
			window, requests, spent = start, req, tok
			if !now.Before(start.Add(time.Minute)) {
				window, requests, spent = now, 0, 0
			}
			if next.After(now) {
				wait = next.Sub(now)
			}
			if requests >= c.RPM || spent+tokens > c.TPM {
				if d := window.Add(time.Minute).Sub(now); d > wait {
					wait = d
				}
			}
		}
	}
	err = rows.Err()
	rows.Close()
	if err != nil {
		return 0, 0, err
	}
	if busy || active >= 2 {
		if wait < time.Second {
			wait = time.Second
		}
	}
	if tokens > c.TPM {
		return 0, 0, ErrDraftBudget
	}
	if wait > 0 {
		return 0, wait, nil
	}
	raw, _ = json.Marshal(counts)
	query, args = runs.UPDATE(runs.Attempts, runs.SectionAttempts).SET(attempts+1, string(raw)).WHERE(runs.JobID.EQ(draftID(id))).Sql()
	if _, err = tx.ExecContext(ctx, query, args...); err != nil {
		return 0, 0, err
	}
	calls := table.ProviderCalls
	sectionJSON, _ := json.Marshal(sections)
	query, args = calls.INSERT(calls.JobID, calls.SectionID, calls.Provider, calls.Model, calls.ErrorCode, calls.InputTokens, calls.OutputTokens).VALUES(id, string(sectionJSON), c.Provider, c.Model, "in_flight", 0, 0).RETURNING(calls.ID).Sql()
	var callID int64
	if err = tx.QueryRowContext(ctx, query, args...).Scan(&callID); err != nil {
		return 0, 0, err
	}
	query, args = limits.UPDATE(limits.WindowStart, limits.Requests, limits.Tokens, limits.ActiveUntil, limits.Holder).SET(window, requests+1, spent+tokens, now.Add(60*time.Second), owner).WHERE(limits.Scope.EQ(postgres.String(scope))).Sql()
	if _, err = tx.ExecContext(ctx, query, args...); err != nil {
		return 0, 0, err
	}
	return callID, 0, tx.Commit()
}

// StartDraftBudget applies the profile budget before the first model call. A
// run with recorded attempts keeps its stored deadline, so a worker restart
// cannot extend the budget indefinitely.
func (q *Queue) StartDraftBudget(ctx context.Context, id, owner string, duration time.Duration) error {
	return q.setDraftDeadline(ctx, id, owner, duration, true)
}

// ExtendDraftDeadline grants one automatic continuation pass after the deadline.
func (q *Queue) ExtendDraftDeadline(ctx context.Context, id, owner string, duration time.Duration) error {
	return q.setDraftDeadline(ctx, id, owner, duration, false)
}

func (q *Queue) setDraftDeadline(ctx context.Context, id, owner string, duration time.Duration, onlyFresh bool) error {
	tx, err := q.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if err = checkOwner(ctx, tx, id, owner); err != nil {
		return err
	}
	runs := table.JobDraftRuns
	query, args := runs.SELECT(runs.Attempts).WHERE(runs.JobID.EQ(draftID(id))).FOR(postgres.UPDATE()).Sql()
	var attempts int
	if err = tx.QueryRowContext(ctx, query, args...).Scan(&attempts); err != nil {
		return err
	}
	if onlyFresh && attempts > 0 {
		return tx.Commit()
	}
	query, args = runs.UPDATE(runs.Deadline).SET(postgres.TimestampzT(time.Now().Add(duration))).WHERE(runs.JobID.EQ(draftID(id))).Sql()
	if _, err = tx.ExecContext(ctx, query, args...); err != nil {
		return err
	}

	return tx.Commit()
}

func (q *Queue) FinishDraftCall(ctx context.Context, id, owner, scope string, callID int64, usage domain.Usage, code string, duration, cooldown time.Duration) error {
	tx, err := q.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	calls := table.ProviderCalls
	query, args := calls.UPDATE(calls.ErrorCode, calls.InputTokens, calls.OutputTokens, calls.DurationMs).SET(code, usage.InputTokens, usage.OutputTokens, duration.Milliseconds()).WHERE(calls.ID.EQ(postgres.Int64(callID)).AND(calls.JobID.EQ(draftID(id)))).Sql()
	if _, err = tx.ExecContext(ctx, query, args...); err != nil {
		return err
	}
	if code != "" && code != "provider_invalid_output" && code != "catalog_ok" {
		query, args = calls.SELECT(calls.SectionID).WHERE(calls.ID.EQ(postgres.Int64(callID))).Sql()
		var sectionJSON []byte
		if err = tx.QueryRowContext(ctx, query, args...).Scan(&sectionJSON); err != nil {
			return err
		}
		var ids []string
		if err = json.Unmarshal(sectionJSON, &ids); err != nil {
			return err
		}
		runs := table.JobDraftRuns
		query, args = runs.SELECT(runs.SectionAttempts).WHERE(runs.JobID.EQ(draftID(id))).FOR(postgres.UPDATE()).Sql()
		var raw []byte
		if err = tx.QueryRowContext(ctx, query, args...).Scan(&raw); err != nil {
			return err
		}
		counts := map[string]int{}
		if err = json.Unmarshal(raw, &counts); err != nil {
			return err
		}
		for _, section := range ids {
			if counts[section] > 0 {
				counts[section]--
			}
		}
		raw, _ = json.Marshal(counts)
		query, args = runs.UPDATE(runs.SectionAttempts).SET(string(raw)).WHERE(runs.JobID.EQ(draftID(id))).Sql()
		if _, err = tx.ExecContext(ctx, query, args...); err != nil {
			return err
		}
	}
	if usage.InputTokens+usage.OutputTokens > 0 {
		query, args = calls.SELECT(calls.Diagnostics, calls.CreatedAt).WHERE(calls.ID.EQ(postgres.Int64(callID))).Sql()
		var diagnostic []byte
		var created time.Time
		if err = tx.QueryRowContext(ctx, query, args...).Scan(&diagnostic, &created); err != nil {
			return err
		}
		var saved struct {
			Reserved int `json:"reserved_tokens"`
		}
		if err = json.Unmarshal(diagnostic, &saved); err != nil {
			return err
		}
		limits := table.ProviderLimits
		query, args = limits.SELECT(limits.Tokens, limits.WindowStart).WHERE(limits.Scope.EQ(postgres.String(scope))).FOR(postgres.UPDATE()).Sql()
		var tokens int
		var window time.Time
		if err = tx.QueryRowContext(ctx, query, args...).Scan(&tokens, &window); err != nil {
			return err
		}
		if !created.Before(window) && saved.Reserved > 0 {
			tokens = max(0, tokens-saved.Reserved+usage.InputTokens+usage.OutputTokens)
			query, args = limits.UPDATE(limits.Tokens).SET(tokens).WHERE(limits.Scope.EQ(postgres.String(scope))).Sql()
			if _, err = tx.ExecContext(ctx, query, args...); err != nil {
				return err
			}
		}
	}
	limits := table.ProviderLimits
	query, args = limits.UPDATE(limits.ActiveUntil, limits.NextAt, limits.Holder).SET(time.Now(), time.Now().Add(cooldown), "").WHERE(limits.Scope.EQ(postgres.String(scope)).AND(limits.Holder.EQ(postgres.String(owner)))).Sql()
	if _, err = tx.ExecContext(ctx, query, args...); err != nil {
		return err
	}
	return tx.Commit()
}

func (q *Queue) DraftNote(ctx context.Context, id, owner, reason string, disabled map[string]string) error {
	tx, err := q.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if err = checkOwner(ctx, tx, id, owner); err != nil {
		return err
	}
	t := table.JobDraftRuns
	raw, _ := json.Marshal(disabled)
	query, args := t.UPDATE(t.Reason, t.Disabled).SET(reason, string(raw)).WHERE(t.JobID.EQ(draftID(id))).Sql()
	if _, err = tx.ExecContext(ctx, query, args...); err != nil {
		return err
	}
	return tx.Commit()
}

// DraftReadyForRender closes generation and queues rendering in one transaction.
func (q *Queue) DraftReadyForRender(ctx context.Context, id, owner string) error {
	tx, err := q.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if err = checkOwner(ctx, tx, id, owner); err != nil {
		return err
	}
	jobs := table.Jobs
	query, args := jobs.SELECT(jobs.CompletedSections, jobs.TotalSections).WHERE(jobs.ID.EQ(draftID(id))).Sql()
	var completed, total int
	if err = tx.QueryRowContext(ctx, query, args...).Scan(&completed, &total); err != nil {
		return err
	}
	if total == 0 || completed != total {
		return errors.New("draft sections incomplete")
	}
	inputs := table.JobDocumentInputs
	query, args = inputs.UPDATE(inputs.Reviewed, inputs.CapabilityHash).SET(true, "").WHERE(inputs.JobID.EQ(draftID(id))).Sql()
	result, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	if count, err := result.RowsAffected(); err != nil || count != 1 {
		return errors.New("document input missing")
	}
	if err = storage.TransitionTx(ctx, tx, id, domain.StatusRunning, domain.StatusReadyForReview, "draft_ready", "draft_ready", "Saved draft is ready for compilation; content still needs review."); err != nil {
		return err
	}
	if err = storage.TransitionTx(ctx, tx, id, domain.StatusReadyForReview, domain.StatusQueued, "rendering", "render_queued", "Compiling the saved baseline and model improvements."); err != nil {
		return err
	}
	query, args = jobs.UPDATE(jobs.Attempts, jobs.LeaseOwner, jobs.LeaseUntil, jobs.ErrorCode, jobs.ErrorMessage).SET(0, postgres.NULL, postgres.NULL, "", "").WHERE(jobs.ID.EQ(draftID(id))).Sql()
	if _, err = tx.ExecContext(ctx, query, args...); err != nil {
		return err
	}
	return tx.Commit()
}
