package storage

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/vasyukov1/ai-doc-gen/internal/domain"
	"github.com/vasyukov1/ai-doc-gen/internal/storage/generated/ai_doc_gen/public/table"
)

var ErrSectionConflict = errors.New("section is not editable")

func (r *Repository) Sections(ctx context.Context, id string) ([]domain.Section, error) {
	t := table.JobSections
	query, args := t.SELECT(t.Record).WHERE(t.JobID.EQ(postgres.CAST(postgres.Text(id)).AS_UUID())).ORDER_BY(t.Position).Sql()
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := []domain.Section{}
	for rows.Next() {
		var raw []byte
		var section domain.Section
		if err := rows.Scan(&raw); err != nil {
			return nil, err
		}
		if err := json.Unmarshal(raw, &section); err != nil {
			return nil, err
		}
		result = append(result, section)
	}
	return result, rows.Err()
}

func (r *Repository) CachedSection(ctx context.Context, key string) (domain.Section, error) {
	t := table.SectionCache
	query, args := t.SELECT(t.Record).WHERE(t.CacheKey.EQ(postgres.String(key))).Sql()
	var raw []byte
	var section domain.Section
	if err := r.db.QueryRowContext(ctx, query, args...).Scan(&raw); err != nil {
		return section, err
	}
	err := json.Unmarshal(raw, &section)
	return section, err
}

func SaveSectionTx(ctx context.Context, tx *sql.Tx, id string, position int, section domain.Section) error {
	raw, err := json.Marshal(section)
	if err != nil {
		return err
	}
	t := table.JobSections
	query, args := t.INSERT(t.JobID, t.SectionID, t.Position, t.Record).VALUES(id, section.ID, position, string(raw)).ON_CONFLICT(t.JobID, t.SectionID).DO_NOTHING().Sql()
	if _, err = tx.ExecContext(ctx, query, args...); err != nil {
		return err
	}
	cache := table.SectionCache
	query, args = cache.INSERT(cache.CacheKey, cache.Record).VALUES(section.CacheKey, string(raw)).ON_CONFLICT(cache.CacheKey).DO_NOTHING().Sql()
	_, err = tx.ExecContext(ctx, query, args...)
	return err
}

func (r *Repository) RecordUsage(ctx context.Context, id, section, provider, model, code string, usage domain.Usage) error {
	t := table.ProviderCalls
	query, args := t.INSERT(t.JobID, t.SectionID, t.Provider, t.Model, t.ErrorCode, t.InputTokens, t.OutputTokens).VALUES(id, section, provider, model, code, usage.InputTokens, usage.OutputTokens).Sql()
	_, err := r.db.ExecContext(ctx, query, args...)
	return err
}

func (r *Repository) ChangeSection(ctx context.Context, id, sectionID, body string, regenerate bool) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	jobs := table.Jobs
	query, args := jobs.SELECT(jobs.Status).WHERE(jobs.ID.EQ(postgres.CAST(postgres.Text(id)).AS_UUID())).FOR(postgres.UPDATE()).Sql()
	var status domain.Status
	if err := tx.QueryRowContext(ctx, query, args...).Scan(&status); err != nil {
		return err
	}
	if status != domain.StatusReadyForReview && status != domain.StatusFailed {
		return ErrSectionConflict
	}
	t := table.JobSections
	condition := t.JobID.EQ(postgres.CAST(postgres.Text(id)).AS_UUID()).AND(t.SectionID.EQ(postgres.String(sectionID)))
	query, args = t.SELECT(t.Record).WHERE(condition).Sql()
	var raw []byte
	var section domain.Section
	if err := tx.QueryRowContext(ctx, query, args...).Scan(&raw); err != nil {
		return err
	}
	if err := json.Unmarshal(raw, &section); err != nil {
		return err
	}
	if regenerate {
		cache := table.SectionCache
		query, args = cache.DELETE().WHERE(cache.CacheKey.EQ(postgres.String(section.CacheKey))).Sql()
		if _, err := tx.ExecContext(ctx, query, args...); err != nil {
			return err
		}
		query, args = t.DELETE().WHERE(condition).Sql()
		if _, err := tx.ExecContext(ctx, query, args...); err != nil {
			return err
		}
		if err := TransitionTx(ctx, tx, id, status, domain.StatusQueued, "generating", "regenerate", "Section regeneration queued. Other drafts are preserved."); err != nil {
			return err
		}
		query, args = jobs.UPDATE(jobs.Attempts, jobs.ErrorCode, jobs.ErrorMessage, jobs.CompletedSections).SET(0, "", "", jobs.CompletedSections.SUB(postgres.Int(1))).WHERE(jobs.ID.EQ(postgres.CAST(postgres.Text(id)).AS_UUID())).Sql()
	} else {
		if len(body) > 40000 || body == "" {
			return ErrSectionConflict
		}
		section.Body = body
		section.Edited = true
		raw, _ = json.Marshal(section)
		query, args = t.UPDATE(t.Record).SET(string(raw)).WHERE(condition).Sql()
	}
	if _, err := tx.ExecContext(ctx, query, args...); err != nil {
		return err
	}
	documentInput := table.JobDocumentInputs
	query, args = documentInput.UPDATE(documentInput.Reviewed, documentInput.CapabilityHash).SET(false, "").WHERE(documentInput.JobID.EQ(postgres.CAST(postgres.Text(id)).AS_UUID())).Sql()
	if _, err := tx.ExecContext(ctx, query, args...); err != nil {
		return err
	}

	return tx.Commit()
}
