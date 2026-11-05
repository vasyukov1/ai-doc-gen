package web

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/vasyukov1/ai-doc-gen/internal/platform"
)

// The bearer secret stays in an HttpOnly cookie; only its digest is persisted.
func ownerCookie(w http.ResponseWriter, r *http.Request, create bool) (string, error) {
	cookie, err := r.Cookie("job_owner")
	value := ""
	if err == nil {
		value = cookie.Value
	}
	raw, err := hex.DecodeString(value)
	if err != nil || len(raw) != 32 {
		if !create {
			return "", errors.New("owner session missing")
		}
		raw = make([]byte, 32)
		if _, err = rand.Read(raw); err != nil {
			return "", err
		}
		value = hex.EncodeToString(raw)
		http.SetCookie(w, &http.Cookie{Name: "job_owner", Value: value, Path: "/", MaxAge: 365 * 24 * 3600, HttpOnly: true, Secure: r.TLS != nil, SameSite: http.SameSiteStrictMode})
	}
	digest := sha256.Sum256([]byte(value))
	return hex.EncodeToString(digest[:]), nil
}

// sameOrigin rejects state-changing requests that a cross-site page could have
// started. Sec-Fetch-Site is set by the browser and cannot be altered by page
// scripts, so same-origin and user-initiated values stay authoritative even when
// a proxy or tunnel rewrites the Host header. Older clients fall back to
// comparing Origin with the host observed by the server.
func sameOrigin(r *http.Request) bool {
	switch r.Header.Get("Sec-Fetch-Site") {
	case "same-origin", "none":
		return true
	case "cross-site":
		return false
	}

	origin := r.Header.Get("Origin")
	if origin == "" {
		return true
	}

	parsed, err := url.Parse(origin)
	if err != nil || parsed.Host == "" {
		return false
	}

	return strings.EqualFold(parsed.Host, r.Host)
}

func (h *Handler) rejectCrossOrigin(w http.ResponseWriter, r *http.Request) {
	log.Printf("cross-origin %s %s rejected: origin=%q host=%q sec_fetch_site=%q",
		r.Method, r.URL.Path, r.Header.Get("Origin"), r.Host, r.Header.Get("Sec-Fetch-Site"))
	http.Error(w, "Cross-origin request rejected.", http.StatusForbidden)
}

func (h *Handler) authorize(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet && r.Method != http.MethodHead && !sameOrigin(r) {
			h.rejectCrossOrigin(w, r)
			return
		}
		path := strings.TrimPrefix(r.URL.Path, "/api")
		if strings.HasPrefix(path, "/jobs/") {
			w.Header().Set("Cache-Control", "private, no-store")
			w.Header().Set("Referrer-Policy", "no-referrer")
			id := strings.Split(strings.TrimPrefix(path, "/jobs/"), "/")[0]
			owner, err := ownerCookie(w, r, false)
			if err != nil || !platform.ValidID(id) || !h.jobs.OwnsJob(r.Context(), id, owner) {
				http.NotFound(w, r)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

func (h *Handler) lifecycle(w http.ResponseWriter, r *http.Request) {
	job, err := h.findJob(r)
	if err != nil {
		writeJobError(w, err)
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, 1000)
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid action.", 400)
		return
	}
	action := r.FormValue("action")
	if err := h.jobs.Lifecycle(r.Context(), job.ID, action); err != nil {
		http.Error(w, "Action unavailable in the current state.", 409)
		return
	}
	if action == "delete" {
		http.Redirect(w, r, "/jobs", 303)
		return
	}
	http.Redirect(w, r, "/jobs/"+job.ID, 303)
}

func (h *Handler) listJobs(w http.ResponseWriter, r *http.Request) {
	owner, err := ownerCookie(w, r, false)
	if err != nil {
		owner = ""
	}
	jobs, err := h.jobs.OwnJobs(r.Context(), owner)
	if err != nil {
		http.Error(w, "Jobs unavailable.", 503)
		return
	}
	w.Header().Set("Cache-Control", "private, no-store")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_ = h.jobsPage.Execute(w, jobs)
}
