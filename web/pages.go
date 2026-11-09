// Package webassets holds the browser assets and page templates.
package webassets

import (
	"embed"
	"fmt"
	"html/template"

	"github.com/vasyukov1/ai-doc-gen/internal/domain"
)

//go:embed templates/*.html
//go:embed static/app.css static/app.js
var FS embed.FS

// Page parses one page together with the shared partials.
func Page(name string) (*template.Template, error) {
	set, err := template.New("pages").Funcs(funcs()).ParseFS(FS, "templates/partials.html", "templates/"+name)
	if err != nil {
		return nil, fmt.Errorf("parse %s: %w", name, err)
	}

	page := set.Lookup(name)
	if page == nil {
		return nil, fmt.Errorf("page %s is not defined", name)
	}

	return page, nil
}

func funcs() template.FuncMap {
	return template.FuncMap{
		"pct":             percent,
		"writtenSections": writtenSections,
		"profileName":     profileName,
		"statusClass":     statusClass,
		"statusLabel":     statusLabel,
		"stageIndex":      stageIndex,
	}
}

// Written sections exclude initial outlines but include the owner's edits.
func writtenSections(job domain.Job) int {
	if job.FreeDraft {
		return max(0, job.Content.Available-job.Content.Baseline)
	}
	return job.CompletedSections
}

func percent(done, total int) int {
	if total <= 0 || done <= 0 {
		return 0
	}
	if done >= total {
		return 100
	}

	return done * 100 / total
}

var profileNames = map[string]string{
	"client-handover":    "Client project handover",
	"arc42":              "Architecture documentation — arc42",
	"pmi":                "Test plan and methodology",
	"pmi-team":           "Test plan and methodology (team)",
	"pz":                 "Explanatory note",
	"ro":                 "Operator’s manual",
	"srs":                "Software requirements specification — SRS",
	"technical-overview": "Technical project overview",
	"tz":                 "Technical specification",
}

func profileName(v any) string {
	id := fmt.Sprint(v)
	if name, ok := profileNames[id]; ok {
		return name
	}

	return id
}

func statusLabel(v any) string {
	switch fmt.Sprint(v) {
	case "queued":
		return "Queued"
	case "running":
		return "Running"
	case "input_ready":
		return "Input ready"
	case "needs_input":
		return "Needs input"
	case "retry_wait":
		return "Waiting to retry"
	case "ready_for_review":
		return "Ready for review"
	case "completed":
		return "PDF compiled"
	case "failed":
		return "Failed"
	case "cancelled":
		return "Cancelled"
	default:
		return fmt.Sprint(v)
	}
}

func statusClass(v any) string {
	switch fmt.Sprint(v) {
	case "completed":
		return "ok"
	case "failed", "cancelled":
		return "error"
	case "input_ready", "needs_input", "ready_for_review":
		return "warn"
	case "queued", "running", "retry_wait":
		return "running"
	default:
		return ""
	}
}

// stageIndex maps the stored stage to the five visible steps.
func stageIndex(stage, status any) int {
	if fmt.Sprint(status) == "completed" {
		return 5
	}

	switch fmt.Sprint(stage) {
	case "selecting_files", "reviewing_evidence", "analyzing":
		return 2
	case "generating":
		return 3
	case "rendering":
		return 4
	default:
		return 1
	}
}
