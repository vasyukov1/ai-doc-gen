package storage

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/vasyukov1/ai-doc-gen/internal/domain"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/vasyukov1/ai-doc-gen/internal/storage/generated/ai_doc_gen/public/table"
)

// GenerationCalls excludes rejected prose; it remains private repair data in storage.
func (r *Repository) GenerationCalls(ctx context.Context, id string) ([]map[string]any, error) {
	t := table.ProviderCalls
	query, args := t.SELECT(t.Provider, t.Model, t.ErrorCode, t.InputTokens, t.OutputTokens, t.DurationMs, t.Diagnostics).WHERE(t.JobID.EQ(jobUUID(id))).ORDER_BY(t.ID.ASC()).LIMIT(100).Sql()
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := []map[string]any{}
	for rows.Next() {
		var provider, model, code string
		var input, output, duration int64
		var raw []byte
		if err = rows.Scan(&provider, &model, &code, &input, &output, &duration, &raw); err != nil {
			return nil, err
		}
		var diagnostic any
		if err = json.Unmarshal(raw, &diagnostic); err != nil {
			return nil, err
		}
		result = append(result, map[string]any{"provider": provider, "model": model, "error_code": code, "input_tokens": input, "output_tokens": output, "duration_ms": duration, "diagnostics": diagnostic})
	}
	return result, rows.Err()
}

// ContinueDraft renews only unfinished baseline work after an explicit owner action.
func (r *Repository) ContinueDraft(ctx context.Context, id string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	jobs := table.Jobs
	query, args := jobs.SELECT(jobs.Status).WHERE(jobs.ID.EQ(jobUUID(id))).FOR(postgres.UPDATE()).Sql()
	var status string
	if err = tx.QueryRowContext(ctx, query, args...).Scan(&status); err != nil {
		return err
	}
	if status != "completed" && status != "ready_for_review" && status != "failed" {
		return ErrSectionConflict
	}
	sections := table.JobSections
	query, args = sections.SELECT(sections.Record).WHERE(sections.JobID.EQ(jobUUID(id))).Sql()
	rows, err := tx.QueryContext(ctx, query, args...)
	if err != nil {
		return err
	}
	unfinished := false
	for rows.Next() {
		var raw []byte
		var section domain.Section
		if err = rows.Scan(&raw); err != nil {
			rows.Close()
			return err
		}
		if err = json.Unmarshal(raw, &section); err != nil {
			rows.Close()
			return err
		}
		if section.Origin == "baseline" && !section.Edited {
			unfinished = true
		}
	}
	err = rows.Err()
	rows.Close()
	if err != nil {
		return err
	}
	if !unfinished {
		return ErrSectionConflict
	}
	// Existing model/user sections remain immutable to generation; only the budget resets.
	runs := table.JobDraftRuns
	query, args = runs.UPDATE(runs.Deadline, runs.Attempts, runs.Disabled, runs.SectionAttempts, runs.Reason).SET(postgres.NOW().ADD(postgres.INTERVAL(240, postgres.SECOND)), 0, "{}", "{}", "").WHERE(runs.JobID.EQ(jobUUID(id))).Sql()
	result, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	if count, _ := result.RowsAffected(); count != 1 {
		return ErrSectionConflict
	}
	if err = TransitionTx(ctx, tx, id, domain.Status(status), "queued", "generating", "continued", "Explicitly continuing unfinished sections; saved prose is preserved."); err != nil {
		return err
	}
	inputs := table.JobDocumentInputs
	query, args = inputs.UPDATE(inputs.Reviewed, inputs.CapabilityHash).SET(false, "").WHERE(inputs.JobID.EQ(jobUUID(id))).Sql()
	if _, err = tx.ExecContext(ctx, query, args...); err != nil {
		return err
	}
	query, args = jobs.UPDATE(jobs.Attempts, jobs.ErrorCode, jobs.ErrorMessage).SET(0, "", "").WHERE(jobs.ID.EQ(jobUUID(id))).Sql()
	if _, err = tx.ExecContext(ctx, query, args...); err != nil {
		return err
	}
	return tx.Commit()
}

func (r *Repository) DraftRepair(ctx context.Context, id string) (string, []domain.ValidationIssue, error) {
	t := table.ProviderCalls
	query, args := t.SELECT(t.RejectedContent, t.Diagnostics).WHERE(t.JobID.EQ(jobUUID(id)).AND(t.RejectedContent.NOT_EQ(postgres.String("")))).ORDER_BY(t.ID.DESC()).LIMIT(1).Sql()
	var content string
	var raw []byte
	err := r.db.QueryRowContext(ctx, query, args...).Scan(&content, &raw)
	if errors.Is(err, sql.ErrNoRows) {
		return "", nil, nil
	}
	if err != nil {
		return "", nil, err
	}
	var diagnostic struct {
		Issues []domain.ValidationIssue `json:"issues"`
	}
	err = json.Unmarshal(raw, &diagnostic)
	return content, diagnostic.Issues, err
}

// DraftCatalog reuses recent successful discovery across worker restarts of this job.
func (r *Repository) DraftCatalog(ctx context.Context, id, provider, model string) (bool, bool, error) {
	t := table.ProviderCalls
	query, args := t.SELECT(t.Diagnostics).WHERE(t.JobID.EQ(jobUUID(id)).AND(t.Provider.EQ(postgres.String(provider))).AND(t.Model.EQ(postgres.String(model))).AND(t.ErrorCode.EQ(postgres.String("catalog_ok"))).AND(t.CreatedAt.GT(postgres.NOW().SUB(postgres.INTERVAL(5, postgres.MINUTE))))).ORDER_BY(t.ID.DESC()).LIMIT(1).Sql()
	var raw []byte
	err := r.db.QueryRowContext(ctx, query, args...).Scan(&raw)
	if errors.Is(err, sql.ErrNoRows) {
		return false, false, nil
	}
	if err != nil {
		return false, false, err
	}
	var diagnostic struct {
		Schema bool `json:"json_schema"`
	}
	if json.Unmarshal(raw, &diagnostic) != nil {
		return false, false, nil
	}
	return true, diagnostic.Schema, nil
}
