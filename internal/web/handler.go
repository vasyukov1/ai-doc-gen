package web

import (
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"time"

	"github.com/vasyukov1/ai-doc-gen/internal/storage"
	webassets "github.com/vasyukov1/ai-doc-gen/web"
)

type Handler struct {
	jobs          *storage.Repository
	uploadDir     string
	indexPage     *template.Template
	jobPage       *template.Template
	jobsPage      *template.Template
	filePage      *template.Template
	evidencePage  *template.Template
	documentPage  *template.Template
	sectionsPage  *template.Template
	reusePage     *template.Template
	encryptionKey string
	retention     time.Duration
	cookieSecure  bool
}

func NewHandler(jobs *storage.Repository, uploadDir, encryptionKey string, retention time.Duration, cookieSecure bool) (http.Handler, error) {
	static, err := fs.Sub(webassets.FS, "static")
	if err != nil {
		return nil, fmt.Errorf("load static assets: %w", err)
	}

	pages := make(map[string]*template.Template, 8)
	for _, name := range []string{"index.html", "job.html", "jobs.html", "files.html", "evidence.html", "document.html", "sections.html", "reuse.html"} {
		page, err := webassets.Page(name)
		if err != nil {
			return nil, fmt.Errorf("load %s: %w", name, err)
		}
		pages[name] = page
	}

	handler := &Handler{
		jobs:          jobs,
		uploadDir:     uploadDir,
		indexPage:     pages["index.html"],
		jobPage:       pages["job.html"],
		jobsPage:      pages["jobs.html"],
		filePage:      pages["files.html"],
		evidencePage:  pages["evidence.html"],
		documentPage:  pages["document.html"],
		sectionsPage:  pages["sections.html"],
		reusePage:     pages["reuse.html"],
		encryptionKey: encryptionKey,
		retention:     retention,
		cookieSecure:  cookieSecure,
	}
	mux := http.NewServeMux()

	mux.HandleFunc("GET /healthz", handler.health)
	mux.HandleFunc("GET /readyz", handler.ready)
	assetServer := http.StripPrefix("/assets/", http.FileServer(http.FS(static)))
	mux.Handle("GET /assets/", noCache(assetServer))
	mux.HandleFunc("POST /jobs", handler.createJob)
	mux.HandleFunc("POST /jobs/{id}/lifecycle", handler.lifecycle)
	mux.HandleFunc("GET /jobs", handler.listJobs)
	mux.HandleFunc("GET /jobs/{id}", handler.jobPageView)
	mux.HandleFunc("GET /jobs/{id}/reuse", handler.reuseView)
	mux.HandleFunc("POST /jobs/{id}/reuse", handler.createReusedJob)
	mux.HandleFunc("GET /jobs/{id}/files", handler.filePageView)
	mux.HandleFunc("POST /jobs/{id}/files", handler.saveFileSelection)
	mux.HandleFunc("GET /jobs/{id}/evidence", handler.evidenceView)
	mux.HandleFunc("GET /api/jobs/{id}/evidence", handler.evidenceView)
	mux.HandleFunc("POST /jobs/{id}/generation/retry", handler.retryGeneration)
	mux.HandleFunc("GET /jobs/{id}/document", handler.documentView)
	mux.HandleFunc("POST /jobs/{id}/document", handler.saveDocument)
	mux.HandleFunc("GET /jobs/{id}/artifacts/{kind}", handler.downloadArtifact)
	mux.HandleFunc("GET /api/jobs/{id}/generation-report", handler.generationReport)
	mux.HandleFunc("POST /jobs/{id}/generation/continue", handler.continueDraft)
	mux.HandleFunc("GET /jobs/{id}/sections", handler.sectionsView)
	mux.HandleFunc("GET /api/jobs/{id}/sections", handler.sectionsView)
	mux.HandleFunc("POST /jobs/{id}/sections/{section}", handler.changeSection)
	mux.HandleFunc("GET /api/jobs/{id}", handler.jobJSON)
	mux.HandleFunc("GET /api/jobs/{id}/files", handler.fileJSON)
	mux.HandleFunc("GET /", handler.home)

	return handler.authorize(mux), nil
}

// noCache revalidates embedded assets, whose fixed names change with every build.
func noCache(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache")
		next.ServeHTTP(w, r)
	})
}

func (h *Handler) health(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, _ = w.Write([]byte("ok\n"))
}

func (h *Handler) ready(w http.ResponseWriter, r *http.Request) {
	if err := h.jobs.Ping(r.Context()); err != nil {
		http.Error(w, "database unavailable", http.StatusServiceUnavailable)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, _ = w.Write([]byte("ready\n"))
}

func (h *Handler) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.indexPage.Execute(w, nil); err != nil {
		log.Printf("render home: %v", err)
	}
}
