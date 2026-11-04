package documents

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"path/filepath"

	"github.com/vasyukov1/ai-doc-gen/internal/domain"
	"github.com/vasyukov1/ai-doc-gen/internal/pgqueue"
)

func (w *Worker) render(ctx, parent context.Context, claim *pgqueue.Claim, input domain.DocumentInput) {
	if err := w.queue.Stage(ctx, claim.ID, w.owner, "rendering", "rendering", "Compiling reviewed sections with the pinned Typst template."); err != nil {
		return
	}
	job, err := w.jobs.GetJob(ctx, claim.ID)
	if err != nil {
		w.reportFailure(parent, claim, err)
		return
	}
	sections, err := w.jobs.Sections(ctx, claim.ID)
	if err != nil {
		w.reportFailure(parent, claim, err)
		return
	}
	evidence, evidenceErr := w.jobs.GetEvidence(ctx, claim.ID)
	if evidenceErr != nil && !errors.Is(evidenceErr, sql.ErrNoRows) {
		w.reportFailure(parent, claim, evidenceErr)
		return
	}
	artifacts, err := Render(ctx, filepath.Join(w.sourceDir, "artifacts", claim.ID), job, input, sections, evidence)
	if err != nil {
		log.Printf("render job %s: %v", claim.ID, err)
		_ = w.queue.Fail(parent, claim.ID, w.owner, "render_failed", "Rendering failed. Check the worker Typst 0.15.1 installation and the document input, then retry rendering from the document page.")
		return
	}
	if err = w.queue.Rendered(ctx, claim.ID, w.owner, artifacts); err != nil {
		log.Printf("save rendered artifacts %s: %v", claim.ID, err)
	}
}
