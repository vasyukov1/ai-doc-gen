package web

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/vasyukov1/ai-doc-gen/internal/domain"
	"github.com/vasyukov1/ai-doc-gen/internal/storage"
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
	if strings.HasPrefix(r.URL.Path, "/api/") {
		writeJSON(w, http.StatusOK, evidence)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.evidencePage.Execute(w, struct {
		Job      domain.Job
		Evidence domain.Evidence
		Editable bool
	}{job, evidence, job.Status == domain.StatusNeedsInput && job.Stage == "reviewing_evidence"}); err != nil {
		log.Printf("render evidence %s: %v", job.ID, err)
	}
}

func (h *Handler) confirmEvidence(w http.ResponseWriter, r *http.Request) {
	job, err := h.findJob(r)
	if err != nil {
		writeJobError(w, err)
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, 2<<20)
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Could not read evidence review.", http.StatusBadRequest)
		return
	}
	corrections := map[string]string{}
	for key, values := range r.PostForm {
		if strings.HasPrefix(key, "correction.") && len(values) == 1 && strings.TrimSpace(values[0]) != "" {
			corrections[strings.TrimPrefix(key, "correction.")] = strings.TrimSpace(values[0])
		}
	}
	if err := h.jobs.ConfirmEvidence(r.Context(), job.ID, r.PostForm.Get("digest"), corrections); err != nil {
		if errors.Is(err, storage.ErrInvalidEvidence) {
			http.Error(w, "Keep corrections below 4000 bytes.", http.StatusBadRequest)
			return
		}
		if errors.Is(err, storage.ErrEvidenceUnavailable) {
			http.Error(w, "Evidence review is no longer available.", http.StatusConflict)
			return
		}
		log.Printf("confirm evidence %s: %v", job.ID, err)
		http.Error(w, "Could not save evidence review.", http.StatusServiceUnavailable)
		return
	}
	http.Redirect(w, r, "/jobs/"+job.ID+"/evidence", http.StatusSeeOther)
}
