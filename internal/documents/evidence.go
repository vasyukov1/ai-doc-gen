package documents

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/vasyukov1/ai-doc-gen/internal/domain"
)

// ExtractEvidence reports declarations as text, without inferring runtime behavior.
func ExtractEvidence(ctx context.Context, inventory domain.Inventory, description string) (domain.Evidence, error) {
	evidence := domain.Evidence{Version: "1", SourceDigest: inventory.Source.SourceDigest, Selection: []string{}, Facts: []domain.Fact{}, Corrections: map[string]string{}}
	files := append([]domain.SourceFile(nil), inventory.Files...)
	sort.Slice(files, func(i, j int) bool { return files[i].Path < files[j].Path })

	truncated := false
	add := func(kind, value, provenance, path string, line int) {
		if len(evidence.Facts) >= 1000 && kind != "project_context" && kind != "analysis_limit" && kind != "record_limit" {
			truncated = true
			return
		}
		evidence.Facts = append(evidence.Facts, domain.Fact{ID: fmt.Sprintf("fact-%04d", len(evidence.Facts)+1), Kind: kind, Value: value, Provenance: provenance, Path: path, StartLine: line, EndLine: line})
	}
	for _, file := range files {
		if !file.Selected {
			continue
		}
		if err := ctx.Err(); err != nil {
			return domain.Evidence{}, err
		}
		evidence.Selection = append(evidence.Selection, file.Path)
		data, err := os.ReadFile(filepath.Join(inventory.Source.WorkspacePath, filepath.FromSlash(file.Path)))
		if err != nil {
			return domain.Evidence{}, fmt.Errorf("read evidence file: %w", err)
		}
		digest := sha256.Sum256(data)
		if hex.EncodeToString(digest[:]) != file.SHA256 {
			return domain.Evidence{}, fmt.Errorf("source file changed")
		}
		add("selected_file", file.Path, "source", file.Path, 0)
		kind := evidenceKind(file.Path)
		if facts, parsed := manifestFacts(file.Path, data); parsed {
			for index, fact := range facts {
				if index >= 500 {
					add("extraction_limit", "Manifest facts truncated after 500 entries.", "unknown", file.Path, 0)
					break
				}
				if len(fact.Value) > 4000 {
					add("extraction_limit", "Manifest value omitted because it exceeds 4000 bytes.", "unknown", file.Path, fact.Line)
					continue
				}
				add(fact.Kind, fact.Value, "source", file.Path, fact.Line)
			}
			continue
		}
		count := 0
		for index, line := range strings.Split(string(data), "\n") {
			value := strings.TrimSpace(line)
			if value == "" || len(value) > 300 {
				continue
			}
			declaration := strings.HasPrefix(value, "func ") || strings.HasPrefix(value, "type ") || strings.HasPrefix(value, "class ") || strings.HasPrefix(value, "def ") || strings.HasPrefix(value, "export ")
			if kind == "code" && !declaration {
				continue
			}
			if count >= 20 {
				add("extraction_limit", "Further excerpts or declarations omitted after 20 entries.", "unknown", file.Path, 0)
				break
			}
			factKind := kind
			if kind == "code" {
				factKind = "code_declaration"
			}
			add(factKind, value, "source", file.Path, index+1)
			count++
		}
	}
	if truncated {
		add("record_limit", "Evidence truncated after 1000 extracted facts. Review omitted information using user context.", "unknown", "", 0)
	}
	if strings.TrimSpace(description) != "" {
		add("project_context", description, "user", "", 0)
	}
	add("analysis_limit", "Runtime behavior and undeclared requirements have not been verified.", "unknown", "", 0)
	payload, err := json.Marshal(evidence)
	if err != nil {
		return domain.Evidence{}, err
	}
	digest := sha256.Sum256(payload)
	evidence.Digest = hex.EncodeToString(digest[:])

	return evidence, nil
}

func evidenceKind(path string) string {
	name := strings.ToLower(filepath.Base(path))
	switch {
	case strings.HasPrefix(name, "readme"):
		return "readme_excerpt"
	case name == "go.mod" || name == "package.json" || name == "cargo.toml" || name == "pyproject.toml" || name == "requirements.txt" || name == "pom.xml":
		return "manifest_excerpt"
	case strings.Contains(path, ".github/workflows/") || name == ".gitlab-ci.yml" || name == "jenkinsfile":
		return "ci_excerpt"
	case strings.Contains(name, "openapi") || strings.Contains(name, "swagger") || strings.HasSuffix(name, ".proto") || strings.HasSuffix(name, ".graphql"):
		return "api_schema_excerpt"
	default:
		return "code"
	}
}
