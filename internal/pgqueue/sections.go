package pgqueue

import (
	"context"
	"github.com/go-jet/jet/v2/postgres"
	"github.com/vasyukov1/ai-doc-gen/internal/domain"
	"github.com/vasyukov1/ai-doc-gen/internal/storage"
	"github.com/vasyukov1/ai-doc-gen/internal/storage/generated/ai_doc_gen/public/table"
)

func (q *Queue) SaveSection(ctx context.Context, id, owner string, position, total int, section domain.Section) error {
	tx, err := q.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if err := checkOwner(ctx, tx, id, owner); err != nil {
		return err
	}
	if err := storage.SaveSectionTx(ctx, tx, id, position, section); err != nil {
		return err
	}
	t := table.JobSections
	query, args := t.SELECT(postgres.COUNT(t.SectionID)).WHERE(t.JobID.EQ(postgres.CAST(postgres.Text(id)).AS_UUID())).Sql()
	var completed int
	if err := tx.QueryRowContext(ctx, query, args...).Scan(&completed); err != nil {
		return err
	}
	jobs := table.Jobs
	query, args = jobs.UPDATE(jobs.CompletedSections, jobs.TotalSections).SET(completed, total).WHERE(jobs.ID.EQ(postgres.CAST(postgres.Text(id)).AS_UUID())).Sql()
	if _, err := tx.ExecContext(ctx, query, args...); err != nil {
		return err
	}
	if err := storage.TransitionTx(ctx, tx, id, domain.StatusRunning, domain.StatusRunning, "generating", "section_saved", "A validated draft section is available."); err != nil {
		return err
	}
	return tx.Commit()
}

func (q *Queue) DraftReady(ctx context.Context, id, owner string) error {
	return q.finish(ctx, id, owner, domain.StatusReadyForReview, "draft_ready", "Draft sections are ready. Review them and provide document metadata to render Typst and PDF.", "", "")
}

func (q *Queue) SectionTotal(ctx context.Context, id, owner string, total int) error {
	tx, err := q.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if err := checkOwner(ctx, tx, id, owner); err != nil {
		return err
	}
	jobs := table.Jobs
	query, args := jobs.UPDATE(jobs.TotalSections).SET(total).WHERE(jobs.ID.EQ(postgres.CAST(postgres.Text(id)).AS_UUID())).Sql()
	if _, err := tx.ExecContext(ctx, query, args...); err != nil {
		return err
	}
	return tx.Commit()
}
