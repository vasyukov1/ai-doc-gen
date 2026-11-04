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

func (h *Handler) fileJSON(w http.ResponseWriter, r *http.Request) {
	_, inventory, ok := h.loadInventory(w, r)
	if !ok {
		return
	}

	writeJSON(w, http.StatusOK, inventory)
}

func (h *Handler) filePageView(w http.ResponseWriter, r *http.Request) {
	job, inventory, ok := h.loadInventory(w, r)
	if !ok {
		return
	}

	hasEligible := false
	for _, file := range inventory.Files {
		if file.Eligible {
			hasEligible = true
			break
		}
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.filePage.Execute(w, struct {
		Job        domain.Job
		Inventory  domain.Inventory
		Editable   bool
		NoEligible bool
	}{Job: job, Inventory: inventory, Editable: job.Status == domain.StatusNeedsInput && job.Stage == "selecting_files" && hasEligible,
		NoEligible: !hasEligible}); err != nil {
		log.Printf("render file list %s: %v", job.ID, err)
	}
}

func (h *Handler) loadInventory(w http.ResponseWriter, r *http.Request) (domain.Job, domain.Inventory, bool) {
	job, err := h.findJob(r)
	if err != nil {
		writeJobError(w, err)
		return domain.Job{}, domain.Inventory{}, false
	}

	inventory, err := h.jobs.GetInventory(r.Context(), job.ID)
	if errors.Is(err, sql.ErrNoRows) {
		http.Error(w, "The file inventory is still being prepared.", http.StatusConflict)
		return domain.Job{}, domain.Inventory{}, false
	}
	if err != nil {
		log.Printf("get file inventory %s: %v", job.ID, err)
		http.Error(w, "Could not load the file inventory.", http.StatusServiceUnavailable)
		return domain.Job{}, domain.Inventory{}, false
	}

	return job, inventory, true
}

func (h *Handler) saveFileSelection(w http.ResponseWriter, r *http.Request) {
	job, err := h.findJob(r)
	if err != nil {
		writeJobError(w, err)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	if err := r.ParseForm(); err != nil {
		writeFormError(w, r, http.StatusBadRequest, "Could not read the file selection.")
		return
	}

	paths := r.PostForm["file"]
	if err := h.jobs.ConfirmSelection(r.Context(), job.ID, paths); err != nil {
		switch {
		case errors.Is(err, storage.ErrInvalidSelection):
			writeFormError(w, r, http.StatusBadRequest, "Select 1 to 500 available files.")
		case errors.Is(err, storage.ErrSelectionUnavailable):
			writeFormError(w, r, http.StatusConflict, "File selection is no longer available.")
		default:
			log.Printf("save file selection %s: %v", job.ID, err)
			writeFormError(w, r, http.StatusServiceUnavailable, "Could not save the file selection.")
		}
		return
	}

	location := "/jobs/" + job.ID
	if strings.Contains(r.Header.Get("Accept"), "application/json") {
		writeJSON(w, http.StatusOK, map[string]string{"status": string(domain.StatusQueued), "url": location})
		return
	}

	http.Redirect(w, r, location, http.StatusSeeOther)
}
