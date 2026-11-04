package storage

import (
	"context"
	"database/sql"
	"errors"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/vasyukov1/ai-doc-gen/internal/domain"
	"github.com/vasyukov1/ai-doc-gen/internal/storage/generated/ai_doc_gen/public/table"
)

// Expiry first revokes leases. Purging waits beyond the maximum live lease.
func (r *Repository) ExpireJobs(ctx context.Context) error {
	for i := 0; i < 20; i++ {
		tx, err := r.db.BeginTx(ctx, nil)
		if err != nil {
			return err
		}
		t := table.Jobs
		query, args := t.SELECT(t.ID, t.Status, t.Stage).WHERE(t.ExpiresAt.LT_EQ(postgres.NOW()).AND(t.PurgeAfter.IS_NULL()).AND(t.OwnerHash.IS_NOT_NULL())).ORDER_BY(t.ExpiresAt).LIMIT(1).FOR(postgres.UPDATE().SKIP_LOCKED()).Sql()
		var id, stage string
		var status domain.Status
		err = tx.QueryRowContext(ctx, query, args...).Scan(&id, &status, &stage)
		if errors.Is(err, sql.ErrNoRows) {
			tx.Rollback()
			return nil
		}
		if err != nil {
			tx.Rollback()
			return err
		}
		if status != domain.StatusCancelled {
			err = TransitionTx(ctx, tx, id, status, domain.StatusCancelled, stage, "expired", "Retention expired. Data is scheduled for removal.")
		}
		if err == nil {
			query, args = t.UPDATE(t.PurgeAfter, t.LeaseOwner, t.LeaseUntil, t.NextRetryAt).SET(postgres.NOW().ADD(postgres.INTERVAL(90, postgres.SECOND)), postgres.NULL, postgres.NULL, postgres.NULL).WHERE(t.ID.EQ(postgres.CAST(postgres.Text(id)).AS_UUID())).Sql()
			_, err = tx.ExecContext(ctx, query, args...)
		}
		if err != nil {
			tx.Rollback()
			return err
		}
		if err = tx.Commit(); err != nil {
			return err
		}
	}
	return nil
}

// The callback only removes controlled paths for this locked, revoked job.
func (r *Repository) PurgeJobs(ctx context.Context, remove func(string) error) error {
	for i := 0; i < 10; i++ {
		tx, err := r.db.BeginTx(ctx, nil)
		if err != nil {
			return err
		}
		t := table.Jobs
		query, args := t.SELECT(t.ID).WHERE(t.PurgeAfter.LT_EQ(postgres.NOW())).ORDER_BY(t.PurgeAfter).LIMIT(1).FOR(postgres.UPDATE().SKIP_LOCKED()).Sql()
		var id string
		err = tx.QueryRowContext(ctx, query, args...).Scan(&id)
		if errors.Is(err, sql.ErrNoRows) {
			tx.Rollback()
			return nil
		}
		if err != nil {
			tx.Rollback()
			return err
		}
		query, args = postgres.SELECT(postgres.RawBool("pg_try_advisory_xact_lock(hashtextextended(:id, 0))", postgres.RawArgs{":id": id})).Sql()
		var idle bool
		if err = tx.QueryRowContext(ctx, query, args...).Scan(&idle); err != nil {
			tx.Rollback()
			return err
		}
		if !idle {
			tx.Rollback()
			return nil
		}
		if err = remove(id); err != nil {
			tx.Rollback()
			return err
		}
		sections, cache := table.JobSections, table.SectionCache
		query, args = cache.DELETE().WHERE(cache.CacheKey.IN(sections.SELECT(postgres.RawString("record->>'cache_key'")).WHERE(sections.JobID.EQ(postgres.CAST(postgres.Text(id)).AS_UUID())))).Sql()
		if _, err = tx.ExecContext(ctx, query, args...); err != nil {
			tx.Rollback()
			return err
		}
		query, args = t.DELETE().WHERE(t.ID.EQ(postgres.CAST(postgres.Text(id)).AS_UUID())).Sql()
		if _, err = tx.ExecContext(ctx, query, args...); err != nil {
			tx.Rollback()
			return err
		}
		if err = tx.Commit(); err != nil {
			return err
		}
	}
	return nil
}

// Referenced paths include legacy data, which cleanup must preserve.
func (r *Repository) ProtectedPaths(ctx context.Context) (map[string]bool, map[string]bool, error) {
	paths := map[string]bool{}
	active := map[string]bool{}
	for _, spec := range []struct {
		table  postgres.Table
		column postgres.ColumnString
	}{
		{table.JobInputs, table.JobInputs.UploadPath}, {table.JobSources, table.JobSources.WorkspacePath}, {table.JobArtifacts, table.JobArtifacts.StoragePath},
	} {
		query, args := spec.table.SELECT(spec.column).Sql()
		rows, err := r.db.QueryContext(ctx, query, args...)
		if err != nil {
			return nil, nil, err
		}
		for rows.Next() {
			var path sql.NullString
			if err = rows.Scan(&path); err != nil {
				rows.Close()
				return nil, nil, err
			}
			if path.Valid {
				paths[path.String] = true
			}
		}
		err = rows.Err()
		rows.Close()
		if err != nil {
			return nil, nil, err
		}
	}
	t := table.Jobs
	query, args := t.SELECT(t.ID).WHERE(t.Status.IN(postgres.String("queued"), postgres.String("running"), postgres.String("retry_wait")).OR(t.OwnerHash.IS_NULL()).OR(t.PurgeAfter.IS_NOT_NULL())).Sql()
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, nil, err
		}
		active[id] = true
	}
	return paths, active, rows.Err()
}

func (r *Repository) CleanupCache(ctx context.Context) error {
	t := table.SectionCache
	query, args := t.SELECT(t.CacheKey).WHERE(t.CreatedAt.LT(postgres.NOW().SUB(postgres.INTERVAL(7, postgres.DAY)))).ORDER_BY(t.CreatedAt).LIMIT(100).Sql()
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return err
	}
	keys := []string{}
	for rows.Next() {
		var key string
		if err := rows.Scan(&key); err != nil {
			rows.Close()
			return err
		}
		keys = append(keys, key)
	}
	err = rows.Err()
	rows.Close()
	if err != nil {
		return err
	}
	for _, key := range keys {
		query, args = t.DELETE().WHERE(t.CacheKey.EQ(postgres.String(key))).Sql()
		if _, err = r.db.ExecContext(ctx, query, args...); err != nil {
			return err
		}
	}
	return nil
}
