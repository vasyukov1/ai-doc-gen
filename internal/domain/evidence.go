package domain

// Evidence keeps extracted facts separate from later user corrections.
type Evidence struct {
	Model        ProjectModel      `json:"project_model"`
	Version      string            `json:"version"`
	Digest       string            `json:"digest"`
	SourceDigest string            `json:"source_digest"`
	Selection    []string          `json:"selection"`
	Facts        []Fact            `json:"facts"`
	Corrections  map[string]string `json:"corrections"`
	Confirmed    bool              `json:"confirmed"`
}

type Fact struct {
	ID         string `json:"id"`
	Kind       string `json:"kind"`
	Value      string `json:"value"`
	Provenance string `json:"provenance"`
	Path       string `json:"path,omitempty"`
	StartLine  int    `json:"start_line,omitempty"`
	EndLine    int    `json:"end_line,omitempty"`
}
