package web

import (
	"errors"
	"html/template"
	"net/http"
	"strings"

	"github.com/vasyukov1/ai-doc-gen/internal/domain"
	"github.com/vasyukov1/ai-doc-gen/internal/storage"
)

func (h *Handler) sectionsView(w http.ResponseWriter, r *http.Request) {
	job, err := h.findJob(r)
	if err != nil {
		writeJobError(w, err)
		return
	}
	sections, err := h.jobs.Sections(r.Context(), job.ID)
	if err != nil {
		http.Error(w, "Could not load draft sections.", 503)
		return
	}
	if strings.HasPrefix(r.URL.Path, "/api/") {
		writeJSON(w, 200, sections)
		return
	}
	page, err := template.ParseFS(assets, "static/sections.html")
	if err != nil {
		http.Error(w, "Page unavailable.", 500)
		return
	}
	_ = page.Execute(w, struct {
		Job      domain.Job
		Sections []domain.Section
		Editable bool
		Refresh  bool
	}{job, sections, job.Status == domain.StatusReadyForReview || job.Status == domain.StatusFailed || job.Status == domain.StatusCompleted, job.Status == domain.StatusRunning || job.Status == domain.StatusQueued || job.Status == domain.StatusRetryWait})
}

func (h *Handler) changeSection(w http.ResponseWriter, r *http.Request) {
	if !sameOrigin(r) {
		h.rejectCrossOrigin(w, r)
		return
	}
	job, err := h.findJob(r)
	if err != nil {
		writeJobError(w, err)
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, 60000)
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Could not read section.", 400)
		return
	}
	action := r.FormValue("action")
	if action != "save" && action != "regenerate" {
		http.Error(w, "Invalid action.", 400)
		return
	}
	if err := h.jobs.ChangeSection(r.Context(), job.ID, r.PathValue("section"), r.FormValue("body"), action == "regenerate", r.FormValue("resolve_missing") == "yes"); err != nil {
		http.Error(w, "Section is unavailable or generation is still running.", 409)
		return
	}
	http.Redirect(w, r, "/jobs/"+job.ID+"/sections", http.StatusSeeOther)
}

func validModel(model string) bool {
	return model == "openai/gpt-4.1-mini" || model == "deepseek/deepseek-chat-v3-0324" || model == "google/gemini-2.5-flash"
}

func (h *Handler) retryGeneration(w http.ResponseWriter, r *http.Request) {
	if !sameOrigin(r) {
		h.rejectCrossOrigin(w, r)
		return
	}
	job, err := h.findJob(r)
	if err != nil {
		writeJobError(w, err)
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, 4000)
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Could not read credentials.", 400)
		return
	}
	key := r.FormValue("api_key")
	if len(key) < 8 || len(key) > 1000 {
		http.Error(w, "Enter your OpenRouter key.", 400)
		return
	}
	if err := h.jobs.RetryGeneration(r.Context(), job.ID, key, h.encryptionKey); err != nil {
		if errors.Is(err, storage.ErrEncryptionKey) {
			http.Error(w, "Server encryption key is invalid. Check LLM_ENCRYPTION_KEY or the state directory.", 503)
			return
		}
		http.Error(w, "Generation retry unavailable.", 409)
		return
	}
	http.Redirect(w, r, "/jobs/"+job.ID+"/sections", 303)
}

func (h *Handler) generationReport(w http.ResponseWriter, r *http.Request) {
	job, err := h.findJob(r)
	if err != nil {
		writeJobError(w, err)
		return
	}
	calls, err := h.jobs.GenerationCalls(r.Context(), job.ID)
	if err != nil {
		http.Error(w, "Could not load generation diagnostics.", 503)
		return
	}
	writeJSON(w, 200, map[string]any{"content": job.Content, "calls": calls})
}

func (h *Handler) continueDraft(w http.ResponseWriter, r *http.Request) {
	if !sameOrigin(r) {
		h.rejectCrossOrigin(w, r)
		return
	}
	job, err := h.findJob(r)
	if err != nil {
		writeJobError(w, err)
		return
	}
	if !job.FreeDraft || job.Content.Available <= job.Content.Model || h.jobs.ContinueDraft(r.Context(), job.ID) != nil {
		http.Error(w, "No unfinished sections available to continue.", 409)
		return
	}
	http.Redirect(w, r, "/jobs/"+job.ID, 303)
}
