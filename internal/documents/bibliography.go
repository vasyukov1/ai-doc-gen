package documents

import (
	"path/filepath"
	"sort"
	"strings"

	"github.com/vasyukov1/ai-doc-gen/internal/domain"
)

// References are suggested documentation, not a claim that an author consulted it.
func bibliography(evidence []domain.Evidence, language string) string {
	refs := map[string]string{}
	for _, e := range evidence {
		for _, f := range e.Facts {
			if f.Path == "" || agentMetadata(f.Path) {
				continue
			}
			switch strings.ToLower(filepath.Ext(f.Path)) {
			case ".go":
				refs["Go documentation"] = "https://go.dev/doc/"
			case ".ts", ".tsx":
				refs["TypeScript documentation"] = "https://www.typescriptlang.org/docs/"
			case ".js", ".jsx", ".mjs":
				refs["JavaScript documentation (MDN)"] = "https://developer.mozilla.org/en-US/docs/Web/JavaScript"
			case ".py":
				refs["Python documentation"] = "https://docs.python.org/3/"
			case ".typ":
				refs["Typst documentation"] = "https://typst.app/docs/"
			}
			if filepath.Base(f.Path) == "go.mod" {
				refs["Go documentation"] = "https://go.dev/doc/"
				if strings.Contains(f.Value, "github.com/jackc/pgx") || strings.Contains(f.Value, "github.com/lib/pq") {
					refs["PostgreSQL documentation"] = "https://www.postgresql.org/docs/"
				}
			}
		}
	}
	names := make([]string, 0, len(refs))
	for name := range refs {
		names = append(names, name)
	}
	sort.Strings(names)
	intro := "Suggested reference documentation for technologies identified in the selected source files. Confirm the applicable versions and add articles actually used."
	if language == "ru" {
		intro = "Справочная документация технологий, обнаруженных в выбранных исходниках. Проверьте применимые версии и добавьте действительно использованные статьи."
	}
	if language == "ca" {
		intro = "Documentació de referència suggerida per a les tecnologies identificades als fitxers font seleccionats. Confirmeu les versions aplicables i afegiu els articles realment consultats."
	}
	if len(names) == 0 {
		if language == "ru" {
			return "TODO: указать внешнюю документацию и использованную литературу; исходные файлы проекта не являются библиографией."
		}
		if language == "ca" {
			return "TODO: indicar la documentació externa i la bibliografia consultada; els fitxers font del projecte no són una bibliografia."
		}
		return "TODO: provide external documentation and consulted literature; project source files are not a bibliography."
	}
	lines := []string{intro, ""}
	for _, name := range names {
		label := name
		if language == "ca" {
			label = map[string]string{"Go documentation": "Documentació de Go", "TypeScript documentation": "Documentació de TypeScript", "JavaScript documentation (MDN)": "Documentació de JavaScript (MDN)", "Python documentation": "Documentació de Python", "Typst documentation": "Documentació de Typst", "PostgreSQL documentation": "Documentació de PostgreSQL"}[name]
		}
		lines = append(lines, "- "+label+": "+refs[name])
	}
	return strings.Join(lines, "\n")
}

func agentMetadata(name string) bool {
	parts := strings.Split(strings.ToLower(filepath.ToSlash(name)), "/")
	for _, part := range parts {
		switch part {
		case ".claude", ".codex", ".agents", ".cursor", ".codegraph", "agents.md", "claude.md":
			return true
		}
	}
	return false
}
