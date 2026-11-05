package web

import (
	"database/sql"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
	"unicode/utf8"

	"github.com/vasyukov1/ai-doc-gen/internal/documents"
	"github.com/vasyukov1/ai-doc-gen/internal/domain"
	"github.com/vasyukov1/ai-doc-gen/internal/templates"
)

type metadataField struct{ Key, Label, Value string }

func (h *Handler) documentView(w http.ResponseWriter, r *http.Request) {
	job, err := h.findJob(r)
	if err != nil {
		writeJobError(w, err)
		return
	}
	profile, err := templates.Load(job.ProfileID)
	if err != nil {
		http.Error(w, "Legacy profile: submit a new job using a specific document variant.", 409)
		return
	}
	input, err := h.jobs.DocumentInput(r.Context(), job.ID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		http.Error(w, "Document input unavailable.", 503)
		return
	}
	fields := []metadataField{}
	english := []string{"Project title", "Organization", "Faculty", "Program", "Author or team members", "Student group", "Supervisor name", "Supervisor position", "Approver name", "Approver position", "Year", "Project code"}
	russian := []string{"Название проекта", "Организация", "Факультет", "Образовательная программа", "Автор или участники команды", "Учебная группа", "Руководитель", "Должность руководителя", "Утверждающее лицо", "Должность утверждающего лица", "Год", "Код проекта"}
	labels := english
	if job.Language == "ru" {
		labels = russian
	}
	for n, key := range documents.MetadataFields {
		value := input.Metadata[key]
		if value == "" && profile.Academic {
			value = documents.FKNMetadataDefaults(job.Language)[key]
		}
		fields = append(fields, metadataField{key, labels[n], value})
	}
	_ = h.documentPage.Execute(w, struct {
		Job                                   domain.Job
		Input                                 domain.DocumentInput
		Title                                 string
		Fields                                []metadataField
		Editable, CanRender, Refresh, Russian bool
	}{job, input, profile.Title(job.Language), fields, (job.Status == domain.StatusNeedsInput && job.CompletedSections == 0) || job.Status == domain.StatusReadyForReview || job.Status == domain.StatusCompleted || job.Status == domain.StatusFailed, job.Status == domain.StatusReadyForReview || job.Status == domain.StatusFailed, job.Status == domain.StatusRunning || job.Status == domain.StatusQueued || job.Status == domain.StatusRetryWait, job.Language == "ru"})
}

func (h *Handler) saveDocument(w http.ResponseWriter, r *http.Request) {
	if !sameOrigin(r) {
		h.rejectCrossOrigin(w, r)
		return
	}
	job, err := h.findJob(r)
	if err != nil {
		writeJobError(w, err)
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, 200000)
	if err = r.ParseMultipartForm(200000); err != nil {
		http.Error(w, "Provide text fields and at most one UTF-8 specification file below 64000 bytes.", 400)
		return
	}
	defer r.MultipartForm.RemoveAll()
	if r.FormValue("action") == "render" {
		if r.FormValue("reviewed") != "yes" {
			http.Error(w, "Confirm that you reviewed all sections and document input.", 400)
			return
		}
		input, err := h.jobs.DocumentInput(r.Context(), job.ID)
		if err != nil {
			http.Error(w, "Save document input before rendering.", 409)
			return
		}
		if _, err = documents.ValidateDocumentInput(input, job.ProfileID, job.Language); err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		sections, err := h.jobs.Sections(r.Context(), job.ID)
		if err != nil {
			http.Error(w, "Sections unavailable.", 503)
			return
		}
		expected := documents.Plan(job.ProfileID, job.Language)
		existing := map[string]string{}
		for _, s := range sections {
			existing[s.ID] = s.Heading
		}
		for _, s := range expected {
			if existing[s.ID] != s.Heading {
				http.Error(w, "The draft uses an obsolete section plan. Submit a new job using a specific variant.", 409)
				return
			}
		}
		if err = h.jobs.QueueRender(r.Context(), job.ID, ""); err != nil {
			http.Error(w, "Rendering is unavailable until all sections are ready for review.", 409)
			return
		}

		http.Redirect(w, r, "/jobs/"+job.ID+"/document", 303)
		return
	}
	input, err := h.documentInputFromRequest(r, job.Language)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	input, err = documents.ValidateDocumentInput(input, job.ProfileID, job.Language)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err = h.jobs.SaveDocumentInput(r.Context(), job.ID, input); err != nil {
		http.Error(w, "Document input is editable while source or evidence review is waiting, before drafting begins.", 409)
		return
	}
	http.Redirect(w, r, "/jobs/"+job.ID+"/document", 303)
}

// documentInputFromRequest reads title-page metadata and one specification
// source from a parsed multipart form.
func (h *Handler) documentInputFromRequest(r *http.Request, language string) (domain.DocumentInput, error) {
	input := domain.DocumentInput{Metadata: map[string]string{}}
	for _, key := range documents.MetadataFields {
		input.Metadata[key] = r.FormValue("metadata." + key)
	}
	input.Specification = strings.TrimSpace(r.FormValue("specification"))

	structured := []struct{ key, en, ru, ca string }{{"purpose", "Purpose", "Назначение", "Propòsit"}, {"functions", "Required functions", "Требуемые функции", "Funcions requerides"}, {"constraints", "Constraints", "Ограничения", "Restriccions"}, {"acceptance", "Acceptance criteria", "Критерии приемки", "Criteris d’acceptació"}}
	parts := []string{}
	for _, field := range structured {
		value := strings.TrimSpace(r.FormValue("spec." + field.key))
		if value != "" {
			label := field.en
			if language == "ru" {
				label = field.ru
			}
			if language == "ca" {
				label = field.ca
			}
			parts = append(parts, label+":\n"+value)
		}
	}

	uploaded := []*multipart.FileHeader{}
	if r.MultipartForm != nil {
		uploaded = r.MultipartForm.File["specification_file"]
	}
	if len(uploaded) > 1 {
		return input, errors.New("upload one specification")
	}
	sources := 0
	if input.Specification != "" {
		sources++
	}
	if len(parts) > 0 {
		sources++
	}
	if len(uploaded) > 0 {
		sources++
	}
	if sources > 1 {
		return input, errors.New("use one specification source: text, structured fields, or upload")
	}
	if len(parts) > 0 {
		if len(parts) != len(structured) {
			return input, errors.New("complete all four structured specification fields; unknown requirements may be stated explicitly")
		}
		input.Specification = strings.Join(parts, "\n\n")
		input.SpecificationName = "structured form"
	}
	if len(uploaded) == 1 {
		file := uploaded[0]
		extension := strings.ToLower(filepath.Ext(file.Filename))
		if (extension != ".txt" && extension != ".md") || file.Size > 64000 {
			return input, errors.New("upload a UTF-8 .txt or .md specification below 64000 bytes; PDF, DOCX and Typst are not accepted")
		}
		opened, err := file.Open()
		if err != nil {
			return input, errors.New("could not read the specification")
		}
		raw, err := io.ReadAll(io.LimitReader(opened, 64001))
		_ = opened.Close()
		if err != nil || len(raw) > 64000 || !utf8.Valid(raw) {
			return input, errors.New("specification must be UTF-8 text below 64000 bytes")
		}
		input.Specification = string(raw)
		input.SpecificationName = filepath.Base(strings.ReplaceAll(file.Filename, "\\", "/"))
	}

	return input, nil
}

func (h *Handler) downloadArtifact(w http.ResponseWriter, r *http.Request) {
	job, err := h.findJob(r)
	if err != nil {
		writeJobError(w, err)
		return
	}
	kind := r.PathValue("kind")
	if kind != "pdf" && kind != "typst" {
		http.NotFound(w, r)
		return
	}
	artifact, err := h.jobs.Artifact(r.Context(), job.ID, kind, "")
	if err != nil || job.Status != domain.StatusCompleted {
		http.Error(w, "Artifact unavailable. Review and rebuild saved changes.", 409)
		return
	}
	name, contentType := "document.pdf", "application/pdf"
	if kind == "typst" {
		name, contentType = "typst-source.zip", "application/zip"
	}
	// The in-site preview frames the PDF, so that request asks for inline
	// rendering. Downloads keep the attachment disposition.
	disposition := "attachment"
	if kind == "pdf" && r.URL.Query().Get("view") == "1" {
		disposition = "inline"
	}
	w.Header().Set("Cache-Control", "private, no-store")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Disposition", disposition+`; filename="`+name+`"`)
	http.ServeFile(w, r, artifact.Path)
}
