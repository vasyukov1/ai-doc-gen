package web

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/vasyukov1/ai-doc-gen/internal/documents"
	"github.com/vasyukov1/ai-doc-gen/internal/domain"
)

func (h *Handler) evidenceView(w http.ResponseWriter, r *http.Request) {
	job, err := h.findJob(r)
	if err != nil {
		writeJobError(w, err)
		return
	}
	evidence, err := h.jobs.GetEvidence(r.Context(), job.ID)
	if errors.Is(err, sql.ErrNoRows) {
		http.Error(w, "Evidence is still being prepared.", http.StatusConflict)
		return
	}
	if err != nil {
		http.Error(w, "Could not load evidence.", http.StatusServiceUnavailable)
		return
	}
	if job.FreeDraft {
		input, err := h.jobs.DocumentInput(r.Context(), job.ID)
		if err != nil {
			http.Error(w, "Could not load document evidence.", 503)
			return
		}
		evidence = documents.DocumentEvidence(evidence, input)
	}
	if strings.HasPrefix(r.URL.Path, "/api/") {
		writeJSON(w, http.StatusOK, evidence)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.evidencePage.Execute(w, struct {
		Job      domain.Job
		Evidence domain.Evidence
	}{job, evidence}); err != nil {
		log.Printf("render evidence %s: %v", job.ID, err)
	}
}
