package web

import (
	"database/sql"
	"errors"
	"html/template"
	"io"
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
		fields = append(fields, metadataField{key, labels[n], input.Metadata[key]})
	}
	page, err := template.ParseFS(assets, "static/document.html")
	if err != nil {
		http.Error(w, "Page unavailable.", 500)
		return
	}
	_ = page.Execute(w, struct {
		Job                                   domain.Job
		Input                                 domain.DocumentInput
		Title                                 string
		Fields                                []metadataField
		Editable, CanRender, Refresh, Russian bool
	}{job, input, profile.Title(job.Language), fields, (job.Status == domain.StatusNeedsInput && job.CompletedSections == 0) || job.Status == domain.StatusReadyForReview || job.Status == domain.StatusCompleted || job.Status == domain.StatusFailed, job.Status == domain.StatusReadyForReview || job.Status == domain.StatusFailed, job.Status == domain.StatusRunning || job.Status == domain.StatusQueued || job.Status == domain.StatusRetryWait, job.Language == "ru"})
}

func (h *Handler) saveDocument(w http.ResponseWriter, r *http.Request) {
	if !sameOrigin(r) {
		http.Error(w, "Cross-origin request rejected.", 403)
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
	input := domain.DocumentInput{Metadata: map[string]string{}}
	for _, key := range documents.MetadataFields {
		input.Metadata[key] = r.FormValue("metadata." + key)
	}
	input.Specification = strings.TrimSpace(r.FormValue("specification"))
	structured := []struct{ key, en, ru string }{{"purpose", "Purpose", "Назначение"}, {"functions", "Required functions", "Требуемые функции"}, {"constraints", "Constraints", "Ограничения"}, {"acceptance", "Acceptance criteria", "Критерии приемки"}}
	parts := []string{}
	for _, field := range structured {
		value := strings.TrimSpace(r.FormValue("spec." + field.key))
		if value != "" {
			label := field.en
			if job.Language == "ru" {
				label = field.ru
			}
			parts = append(parts, label+":\n"+value)
		}
	}
	uploaded := r.MultipartForm.File["specification_file"]
	if len(uploaded) > 1 {
		http.Error(w, "Upload one specification.", 400)
		return
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
		http.Error(w, "Use one specification source: text, structured fields, or upload.", 400)
		return
	}
	if len(parts) > 0 {
		if len(parts) != len(structured) {
			http.Error(w, "Complete all four structured specification fields. Unknown requirements may be stated explicitly.", 400)
			return
		}
		input.Specification = strings.Join(parts, "\n\n")
		input.SpecificationName = "structured form"
	}
	if len(uploaded) == 1 {
		file := uploaded[0]
		extension := strings.ToLower(filepath.Ext(file.Filename))
		if (extension != ".txt" && extension != ".md") || file.Size > 64000 {
			http.Error(w, "Upload a UTF-8 .txt or .md specification below 64000 bytes. PDF, DOCX and Typst are not accepted.", 400)
			return
		}
		opened, err := file.Open()
		if err != nil {
			http.Error(w, "Could not read specification.", 400)
			return
		}
		raw, err := io.ReadAll(io.LimitReader(opened, 64001))
		_ = opened.Close()
		if err != nil || len(raw) > 64000 || !utf8.Valid(raw) {
			http.Error(w, "Specification must be UTF-8 text below 64000 bytes.", 400)
			return
		}
		input.Specification = string(raw)
		input.SpecificationName = filepath.Base(strings.ReplaceAll(file.Filename, "\\", "/"))
	}
	input, err = documents.ValidateDocumentInput(input, job.ProfileID, job.Language)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	if err = h.jobs.SaveDocumentInput(r.Context(), job.ID, input); err != nil {
		http.Error(w, "Document input is editable while source or evidence review is waiting, before drafting begins.", 409)
		return
	}
	http.Redirect(w, r, "/jobs/"+job.ID+"/document", 303)
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
	w.Header().Set("Cache-Control", "private, no-store")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Disposition", `attachment; filename="`+name+`"`)
	http.ServeFile(w, r, artifact.Path)
}
