package web

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/vasyukov1/ai-doc-gen/internal/documents"
	"github.com/vasyukov1/ai-doc-gen/internal/domain"
	"github.com/vasyukov1/ai-doc-gen/internal/platform"
	"github.com/vasyukov1/ai-doc-gen/internal/storage"
)

const maxArchiveBytes = documents.MaxArchiveBytes
const maxRequestBytes = maxArchiveBytes + (1 << 20)

var errArchiveTooLarge = errors.New("ZIP archive is too large")
var errArchiveInvalid = errors.New("file is not a ZIP archive")

func (h *Handler) createJob(w http.ResponseWriter, r *http.Request) {
	owner, err := h.ownerCookie(w, r, true)
	if err != nil {
		http.Error(w, "Could not create private session.", 503)
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, maxRequestBytes)
	if err := r.ParseMultipartForm(1 << 20); err != nil {
		var sizeError *http.MaxBytesError
		if errors.As(err, &sizeError) {
			writeFormError(w, r, http.StatusRequestEntityTooLarge, "Upload is too large.")
			return
		}

		writeFormError(w, r, http.StatusBadRequest, "Submit the form as multipart data.")
		return
	}
	defer r.MultipartForm.RemoveAll()

	archives := r.MultipartForm.File["archive"]
	if len(archives) > 1 {
		writeFormError(w, r, http.StatusBadRequest, "Upload one ZIP archive.")
		return
	}

	archiveName := ""
	if len(archives) == 1 {
		archiveName = path.Base(strings.ReplaceAll(archives[0].Filename, "\\", "/"))
	}

	request, err := documents.ValidateRequest(documents.Request{
		GitHubURL:   r.FormValue("github_url"),
		HasArchive:  len(archives) == 1,
		ArchiveName: archiveName,
		ProfileID:   r.FormValue("profile_id"),
		Language:    r.FormValue("language"),
		Description: r.FormValue("description"),
	})
	if err != nil {
		writeFormError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	documentInput, err := h.documentInputFromRequest(r, request.Language)
	if err != nil {
		writeFormError(w, r, http.StatusBadRequest, err.Error())
		return
	}
	documentInput, err = documents.ValidateDocumentInput(documentInput, request.ProfileID, request.Language)
	if err != nil {
		writeFormError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	var connections []domain.Connection
	if r.FormValue("generation_intent") != "outline" {
		connections, err = connectionsFromForm(r)
		if err != nil {
			writeFormError(w, r, 400, err.Error())
			return
		}
	}
	if r.FormValue("local_enabled") == "yes" && !h.localModelEnabled {
		writeFormError(w, r, 400, "Local inference is not configured on this server. Choose a hosted provider.")
		return
	}
	switch r.FormValue("generation_intent") {
	case "prose":
		if len(connections) == 0 {
			writeFormError(w, r, 400, "Connect a model to generate text, or explicitly choose Outline only. / Подключите модель или выберите только структуру.")
			return
		}
	case "outline":
		connections = nil
	case "": // Preserve existing API clients that predate the explicit intent control.
	default:
		writeFormError(w, r, 400, "Choose document text or outline.")
		return
	}
	// The current form always uses the free draft pipeline. Legacy stored jobs keep their adapter.
	if r.FormValue("api_key") != "" {
		writeFormError(w, r, 400, "Use the free provider fields on the current form.")
		return
	}

	id, err := platform.NewID()
	if err != nil {
		writeFormError(w, r, http.StatusInternalServerError, "Could not create a job.")
		return
	}

	input := domain.Input{SourceType: domain.SourceGitHub, GitHubURL: request.GitHubURL}
	if request.HasArchive {
		input, err = h.saveArchive(r.Context(), id, archiveName, archives[0])
		if err != nil {
			if errors.Is(err, errArchiveTooLarge) || errors.Is(err, errArchiveInvalid) {
				writeFormError(w, r, http.StatusBadRequest, err.Error())
				return
			}

			log.Printf("save upload: %v", err)
			writeFormError(w, r, http.StatusServiceUnavailable, "Could not store the upload. Try again.")
			return
		}
		// Reject unusable archives in the request with the exact path instead
		// of failing the job later in the worker.
		if err = documents.ValidateArchive(input.UploadPath); err != nil {
			_ = os.Remove(input.UploadPath)
			writeFormError(w, r, http.StatusBadRequest, err.Error())
			return
		}
	}

	job := domain.NewJob{
		OwnerHash: owner, ExpiresAt: time.Now().Add(h.retention),
		ID: id, ProfileID: request.ProfileID, Language: request.Language,
		Description: request.Description, Input: input,
		DocumentInput: documentInput, FreeDraft: true, Connections: connections,
	}
	if err := h.jobs.CreateJob(r.Context(), job, h.encryptionKey); err != nil {
		if input.UploadPath != "" {
			_ = os.Remove(input.UploadPath)
		}
		if errors.Is(err, storage.ErrEncryptionKey) {
			writeFormError(w, r, http.StatusServiceUnavailable, "Server encryption key is unavailable. Restart the web and worker processes.")
			return
		}
		log.Printf("save job %s: %v", id, err)
		writeFormError(w, r, http.StatusServiceUnavailable, "Could not save the job. Try again.")
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

func (h *Handler) saveArchive(ctx context.Context, id, name string, header *multipart.FileHeader) (domain.Input, error) {
	if header.Size > maxArchiveBytes {
		return domain.Input{}, errArchiveTooLarge
	}

	source, err := header.Open()
	if err != nil {
		return domain.Input{}, fmt.Errorf("open ZIP archive: %w", err)
	}
	defer source.Close()

	var signature [4]byte
	if _, err := io.ReadFull(source, signature[:]); err != nil || string(signature[:]) != "PK\x03\x04" {
		return domain.Input{}, errArchiveInvalid
	}
	if _, err := source.Seek(0, io.SeekStart); err != nil {
		return domain.Input{}, fmt.Errorf("read ZIP archive: %w", err)
	}

	if err := os.MkdirAll(h.uploadDir, 0o700); err != nil {
		return domain.Input{}, fmt.Errorf("create upload directory: %w", err)
	}

	file, err := os.CreateTemp(h.uploadDir, id+"-*.zip")
	if err != nil {
		return domain.Input{}, fmt.Errorf("save ZIP archive: %w", err)
	}
	saved := false
	defer func() {
		_ = file.Close()
		if !saved {
			_ = os.Remove(file.Name())
		}
	}()

	hash := sha256.New()
	size, err := io.Copy(io.MultiWriter(file, hash), io.LimitReader(source, maxArchiveBytes+1))
	if err != nil {
		return domain.Input{}, fmt.Errorf("save ZIP archive: %w", err)
	}
	if err := ctx.Err(); err != nil {
		return domain.Input{}, err
	}
	if size > maxArchiveBytes {
		return domain.Input{}, errArchiveTooLarge
	}
	if err := file.Close(); err != nil {
		return domain.Input{}, fmt.Errorf("close ZIP archive: %w", err)
	}
	saved = true

	return domain.Input{
		SourceType: domain.SourceZIP, UploadPath: file.Name(),
		OriginalName: name, SHA256: hex.EncodeToString(hash.Sum(nil)), SizeBytes: size,
	}, nil
}

func (h *Handler) jobJSON(w http.ResponseWriter, r *http.Request) {
	job, err := h.findJob(r)
	if err != nil {
		writeJobError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, job)
}

func (h *Handler) jobPageView(w http.ResponseWriter, r *http.Request) {
	job, err := h.findJob(r)
	if err != nil {
		writeJobError(w, err)
		return
	}

	refresh := job.Status == domain.StatusQueued || job.Status == domain.StatusRunning || job.Status == domain.StatusRetryWait
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_ = h.jobPage.Execute(w, struct {
		Job     domain.Job
		Refresh bool
	}{Job: job, Refresh: refresh})
}

func (h *Handler) findJob(r *http.Request) (domain.Job, error) {
	id := r.PathValue("id")
	if !platform.ValidID(id) {
		return domain.Job{}, errInvalidID
	}

	job, err := h.jobs.GetJob(r.Context(), id)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.Printf("get job %s: %v", id, err)
	}
	if err == nil {
		job.GenerationAttemptLimit = documents.DraftBudgetFor(job.ProfileID, job.Language).Attempts
	}

	return job, err
}

var errInvalidID = errors.New("invalid job ID")

func writeJobError(w http.ResponseWriter, err error) {
	if errors.Is(err, errInvalidID) {
		http.Error(w, "Invalid job ID.", http.StatusBadRequest)
		return
	}
	if errors.Is(err, sql.ErrNoRows) {
		http.Error(w, "Job not found.", http.StatusNotFound)
		return
	}

	http.Error(w, "Could not load the job.", http.StatusServiceUnavailable)
}

func writeFormError(w http.ResponseWriter, r *http.Request, status int, message string) {
	if strings.Contains(r.Header.Get("Accept"), "application/json") {
		writeJSON(w, status, map[string]string{"error": message})
		return
	}

	http.Error(w, message, status)
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
