package documents

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/vasyukov1/ai-doc-gen/internal/domain"
)

const AnalysisVersion = "source-context-4"

// sourceFacts keeps whole declarations and paragraphs within explicit bounds.
func sourceFacts(path string, data []byte) []domain.Fact {
	facts := []domain.Fact{}
	add := func(kind, value string, start, end int) {
		if len(facts) >= 80 {
			return
		}
		facts = append(facts, domain.Fact{Kind: kind, Value: clipText(value, 6000), Path: path, StartLine: start, EndLine: end, Provenance: "source"})
	}
	lower := strings.ToLower(path)
	if strings.HasSuffix(lower, ".json") && (strings.Contains(lower, "swagger") || strings.Contains(lower, "openapi")) {
		var schema struct {
			Paths       map[string]map[string]json.RawMessage `json:"paths"`
			Definitions map[string]json.RawMessage            `json:"definitions"`
			Components  struct {
				Schemas map[string]json.RawMessage `json:"schemas"`
			} `json:"components"`
		}
		if json.Unmarshal(data, &schema) == nil && len(schema.Paths) > 0 {
			for _, route := range sortedKeys(schema.Paths) {
				for _, method := range sortedKeys(schema.Paths[route]) {
					if !strings.Contains(" get post put patch delete head options trace ", " "+method+" ") {
						continue
					}
					var op any
					if json.Unmarshal(schema.Paths[route][method], &op) != nil {
						continue
					}
					raw, _ := json.Marshal(op)
					add("api_operation", strings.ToUpper(method)+" "+route+" "+string(raw), 0, 0)
				}
			}
			for _, defs := range []map[string]json.RawMessage{schema.Definitions, schema.Components.Schemas} {
				for _, name := range sortedKeys(defs) {
					add("api_entity", name+": "+string(defs[name]), 0, 0)
				}
			}
			return facts
		}
	}
	switch sourceLanguage(path) {
	case "script":
		if facts, ok := scriptFacts(path, data); ok {
			return facts
		}
	case "python":
		if facts, ok := pythonFacts(path, data); ok {
			return facts
		}
	case "yaml":
		if facts, ok := yamlOpenAPIFacts(path, data); ok {
			return facts
		}
	}
	if strings.HasSuffix(lower, ".go") {
		// Generated Swagger Go duplicates the canonical schema and obscures source coverage.
		if strings.Contains(string(data[:min(len(data), 1000)]), "Code generated") {
			return facts
		}
		fs := token.NewFileSet()
		f, err := parser.ParseFile(fs, path, data, parser.ParseComments)
		if err == nil {
			for _, imp := range f.Imports {
				add("import", imp.Path.Value, fs.Position(imp.Pos()).Line, fs.Position(imp.End()).Line)
			}
			for _, decl := range f.Decls {
				a, b := fs.Position(decl.Pos()), fs.Position(decl.End())
				kind := "code_type"
				caller := "declaration"
				if fn, ok := decl.(*ast.FuncDecl); ok {
					caller = fn.Name.Name
					kind = "code_function"
				}
				add(kind, string(data[a.Offset:b.Offset]), a.Line, b.Line)
				references := []string{}
				ast.Inspect(decl, func(n ast.Node) bool {
					call, ok := n.(*ast.CallExpr)
					if !ok {
						return true
					}
					selector, ok := call.Fun.(*ast.SelectorExpr)
					if !ok {
						return true
					}
					name := selector.Sel.Name
					pos, end := fs.Position(call.Fun.Pos()), fs.Position(call.Fun.End())
					if len(references) < 30 {
						references = append(references, string(data[pos.Offset:end.Offset]))
					}
					if strings.Contains(" GET POST PUT PATCH DELETE Handle HandleFunc Group ", " "+name+" ") {
						a, b := fs.Position(call.Pos()), fs.Position(call.End())
						add("route_declaration", string(data[a.Offset:b.Offset]), a.Line, b.Line)
					}
					return true
				})
				if len(references) > 0 {
					add("call_reference", caller+" contains static call expressions: "+strings.Join(references, ", "), a.Line, b.Line)
				}
			}
			return facts
		}
	}
	if strings.HasSuffix(lower, ".sql") {
		tablePattern := regexp.MustCompile(`(?is)CREATE\s+TABLE\s+(?:IF\s+NOT\s+EXISTS\s+)?([\w."]+)\s*\((.*?)\)\s*;`)
		matches := tablePattern.FindAllStringSubmatchIndex(string(data), -1)
		for _, m := range matches {
			value := string(data[m[0]:m[1]])
			start := strings.Count(string(data[:m[0]]), "\n") + 1
			add("sql_table", value, start, start+strings.Count(value, "\n"))
		}
		if len(matches) > 0 {
			return facts
		}
	}
	// Preserve paragraphs and bounded consecutive source chunks; never skip long lines.
	lines := strings.Split(string(data), "\n")
	start := 0
	for start < len(lines) && len(facts) < 80 {
		end := start
		size := 0
		for end < len(lines) && size < 5000 {
			size += len(lines[end])
			end++
			if end > start+1 && strings.TrimSpace(lines[end-1]) == "" {
				break
			}
		}
		value := strings.TrimSpace(strings.Join(lines[start:end], "\n"))
		if value != "" {
			kind := "source_excerpt"
			if strings.HasPrefix(strings.ToLower(filepath.Base(path)), "readme") {
				kind = "readme_section"
			}
			if strings.HasSuffix(lower, ".sql") {
				kind = "sql_declaration"
			}
			if strings.Contains(lower, "config") || strings.Contains(lower, "compose") {
				kind = "configuration"
			}
			add(kind, value, start+1, end)
		}
		start = end
	}
	if start < len(lines) {
		facts = append(facts, domain.Fact{Kind: "extraction_limit", Value: "Remaining source exceeds the bounded extraction budget.", Path: path, Provenance: "unknown"})
	}
	return facts
}

func sortedKeys[V any](m map[string]V) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func projectBrief(evidence domain.Evidence) domain.ProjectBrief {
	b := domain.ProjectBrief{Version: AnalysisVersion, Context: []domain.Fact{}, Services: []string{}, Unknowns: []string{"Static source analysis does not establish successful execution or test results."}}
	seen := map[string]bool{}
	ordered := append([]domain.Fact{}, evidence.Facts...)
	sort.SliceStable(ordered, func(i, j int) bool { return ordered[i].Provenance == "user" && ordered[j].Provenance != "user" })
	for _, f := range ordered {
		if v, ok := evidence.Corrections[f.ID]; ok {
			f.Value = v
			f.Provenance = "user"
		}
		if f.Kind == "specification" {
			b.RequirementsStatus = "provided_text_not_verified"
			if len([]rune(f.Value)) < 1000 {
				b.RequirementsStatus = "short_description_not_complete_requirements"
			}
		}
		if f.Kind == "go_module" {
			b.Services = append(b.Services, f.Value)
		}
		if f.Provenance == "user" || strings.Contains(f.Kind, "readme") {
			// Preserve the short project description regardless of chapter keywords.
			if !seen[f.Value] && len(b.Context) < 8 {
				f.Value = clipText(f.Value, 2200)
				b.Context = append(b.Context, f)
				seen[f.Value] = true
			}
		}
	}
	if len(b.Context) == 0 {
		b.Unknowns = append(b.Unknowns, "Project purpose was not supplied.")
	}
	return b
}

func buildProjectModel(evidence domain.Evidence) domain.ProjectModel {
	model := domain.ProjectModel{Version: AnalysisVersion, Entities: []domain.Fact{}, Relations: []domain.Fact{}}
	for _, f := range evidence.Facts {
		switch f.Kind {
		case "go_module", "api_entity", "code_type", "sql_declaration", "sql_table":
			model.Entities = append(model.Entities, f)
		case "route_declaration", "api_operation", "call_reference", "import":
			model.Relations = append(model.Relations, f)
		}
	}
	return model
}

func sourceLabel(f domain.Fact) string {
	return fmt.Sprintf("%s:%d-%d", f.Path, f.StartLine, f.EndLine)
}
