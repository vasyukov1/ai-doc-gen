package domain

type DocumentInput struct {
	Metadata          map[string]string `json:"metadata"`
	Specification     string            `json:"specification"`
	SpecificationName string            `json:"specification_name"`
	Reviewed          bool              `json:"reviewed"`
	CapabilityHash    string            `json:"-"`
}
type Artifact struct {
	Kind   string
	Path   string
	SHA256 string
}
