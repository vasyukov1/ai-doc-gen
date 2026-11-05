package domain

// ProjectModel contains static observations, not verified runtime behavior.
type ProjectModel struct {
	Version   string `json:"version"`
	Entities  []Fact `json:"entities"`
	Relations []Fact `json:"relations"`
}

type ProjectBrief struct {
	RequirementsStatus string   `json:"requirements_status"`
	Version            string   `json:"version"`
	Context            []Fact   `json:"context"`
	Services           []string `json:"services"`
	Unknowns           []string `json:"unknowns"`
}

type ChapterContract struct {
	ID       string   `json:"id"`
	Purpose  string   `json:"purpose"`
	Sections []string `json:"sections"`
}

type ChapterEvidencePack struct {
	Brief   ProjectBrief `json:"brief"`
	Sources []Fact       `json:"sources"`
	Omitted int          `json:"omitted"`
}

type ValidationIssue struct {
	SectionID string `json:"section_id"`
	Field     string `json:"field"`
	Reason    string `json:"reason"`
}

type GenerationReport struct {
	Baseline   int    `json:"baseline"`
	Outcome    string `json:"outcome"`
	Available  int    `json:"available"`
	Model      int    `json:"model"`
	Unresolved int    `json:"unresolved"`
}

// ContentReport deliberately does not equate model provenance with quality.
func ContentReport(sections []Section, total int) GenerationReport {
	r := GenerationReport{Outcome: "outline", Available: len(sections)}
	substantive := 0
	for _, s := range sections {
		if s.Origin == "baseline" {
			r.Baseline++
		}
		if s.Origin == "model" {
			r.Model++
		}
		if s.Origin != "baseline" && len(s.Body) > 0 {
			substantive++
		}
		r.Unresolved += len(s.Missing)
	}
	if substantive > 0 {
		r.Outcome = "partial"
	}
	if total > 0 && substantive == total && r.Unresolved == 0 {
		r.Outcome = "draft_ready"
	}
	return r
}
