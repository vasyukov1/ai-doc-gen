package documents

import (
	"path"
	"regexp"
	"strings"

	"github.com/vasyukov1/ai-doc-gen/internal/domain"
)

const maxFactsPerFile = 80
const maxDeclarationLines = 200

var (
	scriptImport  = regexp.MustCompile(`^\s*import\s+(?:[^'"]*?\s+from\s+)?["']([^"']+)["']`)
	scriptRequire = regexp.MustCompile(`require\(\s*["']([^"']+)["']\s*\)`)
	scriptFunc    = regexp.MustCompile(`^\s*(?:export\s+)?(?:default\s+)?(?:async\s+)?function\s+(\w+)\s*\(`)
	scriptClass   = regexp.MustCompile(`^\s*(?:export\s+)?(?:default\s+)?(?:abstract\s+)?class\s+(\w+)`)
	scriptArrow   = regexp.MustCompile(`^\s*(?:export\s+)?(?:const|let|var)\s+(\w+)\s*(?::[^=]+)?=\s*(?:async\s*)?(?:\([^)]*\)|\w+)\s*=>`)
	scriptRoute   = regexp.MustCompile(`\b(?:app|router|server|api)\s*\.\s*(get|post|put|patch|delete|use|route)\s*\(\s*["']([^"']+)["']`)
	scriptCall    = regexp.MustCompile(`\b(?:fetch|got)\s*\(\s*["']([^"']+)["']`)

	pythonImport = regexp.MustCompile(`^\s*(?:import|from)\s+([\w.]+)`)
	pythonDef    = regexp.MustCompile(`^(\s*)(?:async\s+)?def\s+(\w+)\s*\(`)
	pythonClass  = regexp.MustCompile(`^(\s*)class\s+(\w+)`)
	pythonRoute  = regexp.MustCompile(`^\s*@[\w.]+\.(get|post|put|patch|delete|route)\(\s*["']([^"']+)["']`)

	openAPIPaths  = regexp.MustCompile(`^(\s*)paths:\s*$`)
	openAPIRoute  = regexp.MustCompile(`^(\s*)(/[^\s:][^:]*):\s*$`)
	openAPIMethod = regexp.MustCompile(`^(\s*)(get|post|put|patch|delete|head|options|trace):\s*$`)
	openAPISchema = regexp.MustCompile(`^(\s*)(schemas|definitions):\s*$`)
	openAPIName   = regexp.MustCompile(`^(\s*)([A-Za-z_][\w.-]*):\s*$`)
)

var scriptExtensions = map[string]bool{
	".ts": true, ".tsx": true, ".js": true, ".jsx": true,
	".mts": true, ".cts": true, ".mjs": true, ".cjs": true,
}

// scriptFacts keeps one fact per import, declaration and route. Declarations
// are bounded by brace depth, so a long body keeps its signature.
func scriptFacts(filePath string, data []byte) ([]domain.Fact, bool) {
	lines := strings.Split(string(data), "\n")
	facts := make([]domain.Fact, 0, 16)
	add := func(kind, value string, start, end int) {
		if len(facts) >= maxFactsPerFile || strings.TrimSpace(value) == "" {
			return
		}
		facts = append(facts, domain.Fact{Kind: kind, Value: clipText(value, 6000), Path: filePath, StartLine: start, EndLine: end, Provenance: "source"})
	}
	for index, line := range lines {
		if match := scriptImport.FindStringSubmatch(line); match != nil {
			add("import", match[1], index+1, index+1)
		} else if match := scriptRequire.FindStringSubmatch(line); match != nil {
			add("import", match[1], index+1, index+1)
		}
		if match := scriptRoute.FindStringSubmatch(line); match != nil {
			add("route_declaration", strings.ToUpper(match[1])+" "+match[2], index+1, index+1)
		}
		if match := scriptCall.FindStringSubmatch(line); match != nil {
			add("call_reference", "HTTP call to "+match[1], index+1, index+1)
		}
		kind, name := "", ""
		switch {
		case scriptFunc.MatchString(line):
			kind, name = "code_function", scriptFunc.FindStringSubmatch(line)[1]
		case scriptClass.MatchString(line):
			kind, name = "code_type", scriptClass.FindStringSubmatch(line)[1]
		case scriptArrow.MatchString(line):
			kind, name = "code_function", scriptArrow.FindStringSubmatch(line)[1]
		}
		if kind == "" {
			continue
		}
		end := declarationEnd(lines, index)
		add(kind, name+"\n"+strings.Join(lines[index:end+1], "\n"), index+1, end+1)
	}
	return facts, len(facts) > 0
}

// declarationEnd follows balanced braces from a declaration line.
func declarationEnd(lines []string, start int) int {
	depth, opened := 0, false
	for index := start; index < len(lines) && index < start+maxDeclarationLines; index++ {
		depth += strings.Count(lines[index], "{") - strings.Count(lines[index], "}")
		if strings.Contains(lines[index], "{") {
			opened = true
		}
		if opened && depth <= 0 {
			return index
		}
		if !opened && index > start && strings.TrimSpace(lines[index]) == "" {
			return index - 1
		}
	}
	return min(start+maxDeclarationLines-1, len(lines)-1)
}

// pythonFacts follows indentation instead of braces.
func pythonFacts(filePath string, data []byte) ([]domain.Fact, bool) {
	lines := strings.Split(string(data), "\n")
	facts := make([]domain.Fact, 0, 16)
	add := func(kind, value string, start, end int) {
		if len(facts) >= maxFactsPerFile || strings.TrimSpace(value) == "" {
			return
		}
		facts = append(facts, domain.Fact{Kind: kind, Value: clipText(value, 6000), Path: filePath, StartLine: start, EndLine: end, Provenance: "source"})
	}
	for index, line := range lines {
		if match := pythonImport.FindStringSubmatch(line); match != nil {
			add("import", match[1], index+1, index+1)
		}
		if match := pythonRoute.FindStringSubmatch(line); match != nil {
			add("route_declaration", strings.ToUpper(match[1])+" "+match[2], index+1, index+1)
		}
		if match := pythonClass.FindStringSubmatch(line); match != nil {
			end := indentedEnd(lines, index, len(match[1]))
			add("code_type", match[2]+"\n"+strings.Join(lines[index:end+1], "\n"), index+1, end+1)
			continue
		}
		if match := pythonDef.FindStringSubmatch(line); match != nil {
			end := indentedEnd(lines, index, len(match[1]))
			add("code_function", match[2]+"\n"+strings.Join(lines[index:end+1], "\n"), index+1, end+1)
		}
	}
	return facts, len(facts) > 0
}

// indentedEnd returns the last line belonging to an indented block.
func indentedEnd(lines []string, start, indent int) int {
	end := start
	for index := start + 1; index < len(lines) && index < start+maxDeclarationLines; index++ {
		line := lines[index]
		if strings.TrimSpace(line) == "" {
			end = index
			continue
		}
		current := len(line) - len(strings.TrimLeft(line, " \t"))
		if current <= indent {
			break
		}
		end = index
	}
	return end
}

// yamlOpenAPIFacts reads the operations and component names of a YAML schema
// without a general YAML dependency; only indentation and key names are used.
func yamlOpenAPIFacts(filePath string, data []byte) ([]domain.Fact, bool) {
	head := string(data[:min(len(data), 4000)])
	if !strings.Contains(head, "openapi:") && !strings.Contains(head, "swagger:") {
		return nil, false
	}
	lines := strings.Split(string(data), "\n")
	facts := make([]domain.Fact, 0, 16)
	add := func(kind, value string, line int) {
		if len(facts) >= maxFactsPerFile {
			return
		}
		facts = append(facts, domain.Fact{Kind: kind, Value: value, Path: filePath, StartLine: line, EndLine: line, Provenance: "source"})
	}
	inPaths, routeIndent, route, inSchemas, schemaIndent, schemaChild := false, -1, "", false, -1, -1
	for index, line := range lines {
		if match := openAPIPaths.FindStringSubmatch(line); match != nil {
			inPaths, routeIndent, route = true, -1, ""
			continue
		}
		if match := openAPISchema.FindStringSubmatch(line); match != nil {
			inSchemas, schemaIndent, schemaChild = true, len(match[1]), -1
			continue
		}
		indent := len(line) - len(strings.TrimLeft(line, " "))
		if inSchemas {
			if strings.TrimSpace(line) != "" && indent <= schemaIndent {
				inSchemas = false
			} else if match := openAPIName.FindStringSubmatch(line); match != nil && len(match[1]) > schemaIndent {
				// Only direct children of schemas are component names.
				if schemaChild < 0 {
					schemaChild = len(match[1])
				}
				if len(match[1]) == schemaChild {
					add("api_entity", match[2], index+1)
					continue
				}
			}
		}
		if inPaths {
			if strings.TrimSpace(line) != "" && indent <= 0 && !strings.HasPrefix(strings.TrimSpace(line), "#") {
				inPaths, route = false, ""
				continue
			}
			if match := openAPIRoute.FindStringSubmatch(line); match != nil {
				route, routeIndent = match[2], len(match[1])
				continue
			}
			if match := openAPIMethod.FindStringSubmatch(line); match != nil && route != "" && len(match[1]) > routeIndent {
				add("api_operation", strings.ToUpper(match[2])+" "+route, index+1)
			}
		}
	}
	return facts, len(facts) > 0
}

func sourceLanguage(filePath string) string {
	lower := strings.ToLower(filePath)
	switch {
	case scriptExtensions[path.Ext(lower)]:
		return "script"
	case strings.HasSuffix(lower, ".py") || strings.HasSuffix(lower, ".pyi"):
		return "python"
	case strings.HasSuffix(lower, ".yaml") || strings.HasSuffix(lower, ".yml"):
		return "yaml"
	}
	return ""
}
