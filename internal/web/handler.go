package web

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"
)

//go:embed static/*
var assets embed.FS

func NewHandler() (http.Handler, error) {
	mux := http.NewServeMux()

	static, err := fs.Sub(assets, "static")
	if err != nil {
		return nil, fmt.Errorf("load static assets: %w", err)
	}

	files := http.FileServer(http.FS(static))

	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, _ = w.Write([]byte("ok\n"))
	})
	mux.Handle("GET /assets/", http.StripPrefix("/assets/", files))
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		data, err := assets.ReadFile("static/index.html")
		if err != nil {
			http.Error(w, "page unavailable", http.StatusInternalServerError)
			return
		}

		_, _ = w.Write(data)
	})

	return mux, nil
}
