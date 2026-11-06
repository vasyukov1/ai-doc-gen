package web

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/vasyukov1/ai-doc-gen/internal/documents"
	"github.com/vasyukov1/ai-doc-gen/internal/domain"
	"github.com/vasyukov1/ai-doc-gen/internal/platform"
	"github.com/vasyukov1/ai-doc-gen/internal/storage"
	"github.com/vasyukov1/ai-doc-gen/internal/templates"
)

type reuseProfile struct{ ID, Title string }
type reuseConnection struct {
	Provider, Model string
	RPM, TPM        int
}
type reusePageData struct {
	Job                        domain.Job
	ProfileID, Language        string
	Description, Specification string
	Fields                     []metadataField
	Profiles                   []reuseProfile
	Source                     domain.Source
	RetainedFiles              int
	Connections                []reuseConnection
	ReuseConnections           bool
	ReviewedInput              bool
	Available                  bool
	Error                      string
	UncertainID                string
	ConnectionError            string
}

func (h *Handler) reuseData(r *http.Request, job domain.Job) (reusePageData, error) {
	data := reusePageData{Job: job, ProfileID: job.ProfileID, Language: job.Language}
	input, err := h.jobs.DocumentInput(r.Context(), job.ID)
	if err != nil {
		return data, err
	}
	data.Specification = input.Specification
	data.Description, err = h.jobs.Description(r.Context(), job.ID)
	if err != nil {
		return data, err
	}
	data.Fields = reuseMetadata(input.Metadata)
	entries, err := templates.Assets.ReadDir("v1")
	if err != nil {
		return data, err
	}
	for _, entry := range entries {
		id := strings.TrimSuffix(entry.Name(), ".json")
		profile, err := templates.Load(id)
		if err == nil {
			data.Profiles = append(data.Profiles, reuseProfile{ID: id, Title: profile.Title("en")})
		}
	}
	inventory, err := h.jobs.GetInventory(r.Context(), job.ID)
	if errors.Is(err, sql.ErrNoRows) {
		data.Error = "The source is not retained yet. Wait for source inspection, then open this page again."
		return data, nil
	}
	if err != nil {
		return data, err
	}
	data.Source = inventory.Source
	for _, file := range inventory.Files {
		if file.Eligible {
			data.RetainedFiles++
		}
	}
	data.Available = data.RetainedFiles > 0
	if !data.Available {
		data.Error = "No retained text files are available. Submit the source again."
	}
	if job.FreeDraft {
		run, err := h.jobs.DraftRun(r.Context(), job.ID, h.encryptionKey)
		if err != nil {
			data.ConnectionError = "Saved connections cannot be opened. You can still copy the source without them."
			return data, nil
		}
		for _, connection := range run.Connections {
			data.Connections = append(data.Connections, reuseConnection{connection.Provider, connection.Model, connection.RPM, connection.TPM})
		}
	}

	return data, nil
}

func reuseMetadata(values map[string]string) []metadataField {
	labels := []string{"Project title", "Organization", "Faculty", "Program", "Author or team members", "Student group", "Supervisor name", "Supervisor position", "Approver name", "Approver position", "Year", "Project code"}
	fields := make([]metadataField, 0, len(documents.MetadataFields))
	for index, key := range documents.MetadataFields {
		fields = append(fields, metadataField{key, labels[index], values[key]})
	}

	return fields
}

func (h *Handler) reuseView(w http.ResponseWriter, r *http.Request) {
	job, err := h.findJob(r)
	if err != nil {
		writeJobError(w, err)
		return
	}
	data, err := h.reuseData(r, job)
	if err != nil {
		http.Error(w, "Saved project input is unavailable. Try again later.", http.StatusServiceUnavailable)
		return
	}

	h.renderReuse(w, r, data, http.StatusOK, data.Error)
}

func (h *Handler) renderReuse(w http.ResponseWriter, r *http.Request, data reusePageData, status int, message string) {
	if status >= 400 && strings.Contains(r.Header.Get("Accept"), "application/json") {
		writeFormError(w, r, status, message)
		return
	}
	data.Error = message
	knownProfile := false
	for _, profile := range data.Profiles {
		if profile.ID == data.ProfileID {
			knownProfile = true
		}
	}
	if !knownProfile {
		data.Profiles = append(data.Profiles, reuseProfile{ID: data.ProfileID, Title: "Unavailable profile: " + data.ProfileID})
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	if err := h.reusePage.Execute(w, data); err != nil {
		log.Printf("render reuse %s: %v", data.Job.ID, err)
	}
}

func (h *Handler) createReusedJob(w http.ResponseWriter, r *http.Request) {
	if !sameOrigin(r) {
		h.rejectCrossOrigin(w, r)
		return
	}
	job, err := h.findJob(r)
	if err != nil {
		writeJobError(w, err)
		return
	}
	data, err := h.reuseData(r, job)
	if err != nil {
		http.Error(w, "Saved project input is unavailable. Try again later.", http.StatusServiceUnavailable)
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	if err := r.ParseForm(); err != nil {
		h.renderReuse(w, r, data, http.StatusBadRequest, "The form is too large or invalid. Keep the specification below 64,000 UTF-8 bytes.")
		return
	}
	data.ProfileID, data.Language = r.FormValue("profile_id"), r.FormValue("language")
	data.Description, data.Specification = r.FormValue("description"), r.FormValue("specification")
	data.ReuseConnections = r.FormValue("reuse_connections") == "yes"
	data.ReviewedInput = r.FormValue("reviewed_input") == "yes"
	input, err := h.documentInputFromRequest(r, data.Language)
	data.Fields = reuseMetadata(input.Metadata)
	if err != nil {
		h.renderReuse(w, r, data, http.StatusBadRequest, err.Error())
		return
	}
	request, err := documents.ValidateRequest(documents.Request{HasArchive: true, ArchiveName: "retained-project.zip", ProfileID: data.ProfileID, Language: data.Language, Description: data.Description})
	if err != nil {
		h.renderReuse(w, r, data, http.StatusBadRequest, err.Error())
		return
	}
	input, err = documents.ValidateDocumentInput(input, request.ProfileID, request.Language)
	if err != nil {
		h.renderReuse(w, r, data, http.StatusBadRequest, err.Error())
		return
	}
	if r.FormValue("reviewed_input") != "yes" {
		h.renderReuse(w, r, data, http.StatusBadRequest, "Review the copied input and confirm that missing metadata may appear as TODOs.")
		return
	}
	if !data.Available {
		h.renderReuse(w, r, data, http.StatusConflict, "The source is not available yet. Wait for source inspection or submit it again.")
		return
	}
	connections := []domain.Connection{}
	if data.ReuseConnections {
		if !job.FreeDraft || len(data.Connections) == 0 {
			h.renderReuse(w, r, data, http.StatusBadRequest, "This job has no reusable free-draft connections. Uncheck connection reuse to create a source-based outline.")
			return
		}
		run, err := h.jobs.DraftRun(r.Context(), job.ID, h.encryptionKey)
		if err != nil {
			h.renderReuse(w, r, data, http.StatusServiceUnavailable, "Saved connections are unavailable. Try again later.")
			return
		}
		for _, connection := range run.Connections {
			if err := documents.ValidateConnection(connection); err != nil {
				h.renderReuse(w, r, data, http.StatusBadRequest, "A saved connection no longer meets the current requirements. Uncheck connection reuse to create an outline.")
				return
			}
		}
		connections = run.Connections
	}

	release, err := h.jobs.GuardJob(r.Context(), job.ID)
	if err != nil {
		h.renderReuse(w, r, data, http.StatusConflict, "The source is in use. Your choices are kept; try again after the current work stops.")
		return
	}
	defer release()
	inventory, err := h.jobs.GetInventory(r.Context(), job.ID)
	if err != nil {
		h.renderReuse(w, r, data, http.StatusConflict, "The retained source is unavailable. Submit the source again.")
		return
	}
	id, err := platform.NewID()
	if err != nil {
		h.renderReuse(w, r, data, http.StatusServiceUnavailable, "Could not create a job. Try again.")
		return
	}
	copied, source, err := documents.CopyRetainedSource(r.Context(), inventory, id, request.ProfileID, h.uploadDir)
	if err != nil {
		log.Printf("copy retained source %s: %v", job.ID, err)
		h.renderReuse(w, r, data, http.StatusConflict, "The retained source could not be copied. It may be missing or changed. Try again or submit the source again.")
		return
	}
	owner, err := h.ownerCookie(w, r, false)
	if err != nil {
		_ = os.RemoveAll(copied.Source.WorkspacePath)
		_ = os.Remove(source.UploadPath)
		http.NotFound(w, r)
		return
	}
	newJob := domain.NewJob{ID: id, OwnerHash: owner, ExpiresAt: time.Now().Add(h.retention), ProfileID: request.ProfileID, Language: request.Language, Description: request.Description, Input: source, DocumentInput: input, FreeDraft: true, Connections: connections}
	if err := h.jobs.CreateReusedJob(r.Context(), job.ID, newJob, copied, h.encryptionKey); err != nil {
		if errors.Is(err, storage.ErrReuseCommitUnknown) {
			data.UncertainID = id
			h.renderReuse(w, r, data, http.StatusServiceUnavailable, "The save response was interrupted. Check your documents before submitting again; a new document may already exist.")
			return
		}
		_ = os.RemoveAll(copied.Source.WorkspacePath)
		_ = os.Remove(source.UploadPath)
		if errors.Is(err, sql.ErrNoRows) {
			http.NotFound(w, r)
			return
		}
		log.Printf("create reused job %s: %v", id, err)
		h.renderReuse(w, r, data, http.StatusServiceUnavailable, "Could not save the new job. Your choices are kept; try again.")
		return
	}

	location := "/jobs/" + id
	w.Header().Set("Location", location)
	if strings.Contains(r.Header.Get("Accept"), "application/json") {
		writeJSON(w, http.StatusAccepted, map[string]string{"id": id, "url": location})
		return
	}
	http.Redirect(w, r, location, http.StatusSeeOther)
}
