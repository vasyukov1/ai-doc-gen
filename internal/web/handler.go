package web

import (
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"

	"github.com/vasyukov1/ai-doc-gen/internal/storage"
)

//go:embed static/*
var assets embed.FS

type Handler struct {
	jobs          *storage.Repository
	uploadDir     string
	jobPage       *template.Template
	filePage      *template.Template
	evidencePage  *template.Template
	encryptionKey string
}

func NewHandler(jobs *storage.Repository, uploadDir, encryptionKey string) (http.Handler, error) {
	static, err := fs.Sub(assets, "static")
	if err != nil {
		return nil, fmt.Errorf("load static assets: %w", err)
	}

	jobPage, err := template.ParseFS(assets, "static/job.html")
	if err != nil {
		return nil, fmt.Errorf("load job page: %w", err)
	}
	filePage, err := template.ParseFS(assets, "static/files.html")
	if err != nil {
		return nil, fmt.Errorf("load file page: %w", err)
	}

	evidencePage, err := template.ParseFS(assets, "static/evidence.html")
	if err != nil {
		return nil, fmt.Errorf("load evidence page: %w", err)
	}

	handler := &Handler{jobs: jobs, uploadDir: uploadDir, jobPage: jobPage, filePage: filePage, evidencePage: evidencePage, encryptionKey: encryptionKey}
	mux := http.NewServeMux()

	mux.HandleFunc("GET /healthz", handler.health)
	mux.HandleFunc("GET /readyz", handler.ready)
	mux.Handle("GET /assets/", http.StripPrefix("/assets/", http.FileServer(http.FS(static))))
	mux.HandleFunc("POST /jobs", handler.createJob)
	mux.HandleFunc("GET /jobs/{id}", handler.jobPageView)
	mux.HandleFunc("GET /jobs/{id}/files", handler.filePageView)
	mux.HandleFunc("POST /jobs/{id}/files", handler.saveFileSelection)
	mux.HandleFunc("GET /jobs/{id}/evidence", handler.evidenceView)
	mux.HandleFunc("POST /jobs/{id}/evidence", handler.confirmEvidence)
	mux.HandleFunc("GET /api/jobs/{id}/evidence", handler.evidenceView)
	mux.HandleFunc("POST /jobs/{id}/generation/retry", handler.retryGeneration)
	mux.HandleFunc("GET /jobs/{id}/sections", handler.sectionsView)
	mux.HandleFunc("GET /api/jobs/{id}/sections", handler.sectionsView)
	mux.HandleFunc("POST /jobs/{id}/sections/{section}", handler.changeSection)
	mux.HandleFunc("GET /api/jobs/{id}", handler.jobJSON)
	mux.HandleFunc("GET /api/jobs/{id}/files", handler.fileJSON)
	mux.HandleFunc("GET /", handler.home)

	return mux, nil
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

	page, err := assets.ReadFile("static/index.html")
	if err != nil {
		http.Error(w, "page unavailable", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write(page)
}
