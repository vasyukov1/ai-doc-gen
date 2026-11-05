package documents

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"unicode"

	"github.com/vasyukov1/ai-doc-gen/internal/domain"
)

// sectionOutput excludes server-owned cache, provenance and edit metadata.
type sectionOutput struct {
	Tables    []domain.DocumentTable `json:"tables,omitempty"`
	ID        string                 `json:"id"`
	Heading   string                 `json:"heading"`
	Body      string                 `json:"body"`
	Citations []string               `json:"citations"`
	Missing   []string               `json:"missing_information"`
}

const DraftVersion = "chapters-3"

func clipText(value string, limit int) string {
	value = strings.Map(func(r rune) rune {
		if unicode.IsControl(r) && r != '\n' && r != '\t' && r != '\r' {
			return -1
		}
		return r
	}, value)
	runes := []rune(value)
	if len(runes) <= limit {
		return value
	}
	return string(runes[:limit]) + "…"
}

func topicWords(plan SectionPlan) []string {
	words := strings.FieldsFunc(strings.ToLower(plan.Topic+" "+plan.Heading), func(r rune) bool { return !unicode.IsLetter(r) })
	result := []string{}
	for _, w := range words {
		if len([]rune(w)) >= 4 {
			result = append(result, w)
		}
	}
	return result
}

// CompactFacts ranks specific evidence before general README context.
func CompactFacts(evidence domain.Evidence, plan SectionPlan) []domain.Fact {
	type ranked struct {
		fact  domain.Fact
		score int
	}
	candidates := []ranked{}
	for _, fact := range evidence.Facts {
		if fact.Provenance == "unknown" || fact.Kind == "selected_file" {
			continue
		}
		if correction, ok := evidence.Corrections[fact.ID]; ok {
			fact.Value = correction
			fact.Provenance = "user"
		}
		score := 0
		topic := strings.ToLower(plan.Topic)
		if strings.Contains(fact.Kind, "readme") && (strings.Contains(topic, "overview") || strings.Contains(topic, "purpose") || strings.Contains(topic, "application")) {
			score += 4
		}
		if strings.Contains(fact.Kind, "api") && (strings.Contains(topic, "interface") || strings.Contains(topic, "input") || strings.Contains(topic, "output")) {
			score += 4
		}
		text := strings.ToLower(fact.Kind + " " + fact.Path + " " + fact.Value)
		for _, word := range topicWords(plan) {
			if strings.Contains(text, word) {
				score += 3
			}
		}
		for _, kind := range plan.Kinds {
			if kind != "readme" && kind != "description" && kind != "unknown" && strings.Contains(strings.ToLower(fact.Kind+" "+fact.Path), kind) {
				score += 2
			}
		}
		if strings.Contains(strings.ToLower(plan.Topic), "name") && fact.Kind == "project_title" {
			score += 20
		}
		if fact.Kind == "specification" {
			// Select a relevant excerpt instead of repeating the entire specification.
			best := ""
			bestScore := 0
			for _, line := range strings.FieldsFunc(fact.Value, func(r rune) bool { return r == '\n' || r == '.' || r == ';' }) {
				n := 0
				for _, word := range topicWords(plan) {
					if strings.Contains(strings.ToLower(line), word) {
						n++
					}
				}
				if n > bestScore {
					bestScore = n
					best = line
				}
			}
			if bestScore == 0 {
				continue
			}
			fact.Value = strings.TrimSpace(best)
			score = bestScore * 3
		}
		if score == 0 {
			continue
		}
		fact.Value = clipText(fact.Value, 450)
		candidates = append(candidates, ranked{fact, score})
	}
	sort.SliceStable(candidates, func(i, j int) bool { return candidates[i].score > candidates[j].score })
	facts := []domain.Fact{}
	for _, entry := range candidates {
		facts = append(facts, entry.fact)
		if len(facts) == 4 {
			break
		}
	}
	// General project descriptions belong in overview/purpose sections only.
	if len(facts) == 0 && (strings.Contains(strings.ToLower(plan.Topic), "overview") || strings.Contains(strings.ToLower(plan.Topic), "purpose") || strings.Contains(strings.ToLower(plan.Topic), "application")) {
		for _, fact := range evidence.Facts {
			if strings.Contains(fact.Kind, "readme") || fact.Kind == "project_context" {
				fact.Value = clipText(fact.Value, 350)
				facts = append(facts, fact)
				if len(facts) == 3 {
					break
				}
			}
		}
	}
	return facts
}

func BaselineSection(plan SectionPlan, facts []domain.Fact) domain.Section {
	section := domain.Section{ID: plan.ID, Heading: plan.Heading, Origin: "baseline", Citations: []string{}, Missing: []string{}}
	paragraphs := []string{}
	for _, fact := range facts {
		section.Citations = append(section.Citations, fact.ID)
		if fact.Provenance == "user" && (fact.Kind == "project_title" || fact.Kind == "project_context") {
			paragraphs = append(paragraphs, clipText(fact.Value, 500))
		}
	}
	gap := "Section “" + plan.Heading + "” needs review: the excerpts do not establish complete requirements or verified runtime behavior."
	if len(facts) == 0 {
		gap = "No specific evidence for “" + plan.Heading + "” was found. Provide the relevant requirements or source references."
	}
	if plan.Language == "ru" {
		gap = "Раздел «" + plan.Heading + "» требует проверки: выдержки не подтверждают полноту требований или поведение при выполнении."
		if len(facts) == 0 {
			gap = "Для раздела «" + plan.Heading + "» не найдены конкретные сведения. Добавьте требования или ссылки на исходные материалы."
		}
	}
	section.Missing = append(section.Missing, gap)
	if len(paragraphs) == 0 {
		paragraphs = append(paragraphs, gap)
	}
	section.Body = strings.Join(paragraphs, "\n\n")
	return section
}

type draftGroup struct {
	Plans    []SectionPlan
	Facts    map[string][]domain.Fact
	Title    string
	Profile  string
	Brief    domain.ProjectBrief
	Sources  []domain.Fact
	Issues   []domain.ValidationIssue
	Previous string
}

func GroupPrompt(group draftGroup, language string) string { return chapterPrompt(group, language) }

func ParseGroup(raw string, group draftGroup) map[string]domain.Section {
	result, _ := parseChapter(raw, group)
	return result
}

func parseChapter(raw string, group draftGroup) (map[string]domain.Section, []domain.ValidationIssue) {
	var envelope struct {
		Sections []json.RawMessage `json:"sections"`
	}
	result := map[string]domain.Section{}
	issues := []domain.ValidationIssue{}
	raw = strings.TrimSpace(raw)
	if strings.HasPrefix(raw, "```json") && strings.HasSuffix(raw, "```") {
		raw = strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(raw, "```json"), "```"))
	}
	if json.Unmarshal([]byte(raw), &envelope) != nil {
		envelope.Sections = completeMembers(raw)
		if len(envelope.Sections) == 0 {
			return result, []domain.ValidationIssue{{Field: "response", Reason: "invalid or truncated JSON; return only requested unfinished sections"}}
		}
	}
	counts := map[string]int{}
	for _, member := range envelope.Sections {
		var s sectionOutput
		if json.Unmarshal(member, &s) == nil {
			counts[s.ID]++
		}
	}
	for _, p := range group.Plans {
		if counts[p.ID] != 1 {
			issues = append(issues, domain.ValidationIssue{SectionID: p.ID, Field: "id", Reason: "expected exactly one section"})
			continue
		}
		for _, member := range envelope.Sections {
			var output sectionOutput
			if json.Unmarshal(member, &output) != nil || output.ID != p.ID {
				continue
			}
			section, err := validateChapterSection(string(member), p, group.Facts[p.ID])
			if err != nil {
				issues = append(issues, domain.ValidationIssue{SectionID: p.ID, Field: "section", Reason: err.Error()})
				continue
			}
			section.Origin = "model"
			result[p.ID] = section
		}
	}
	return result, issues
}

func validateChapterSection(raw string, plan SectionPlan, facts []domain.Fact) (domain.Section, error) {
	var output sectionOutput
	decoder := json.NewDecoder(strings.NewReader(raw))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&output); err != nil {
		return domain.Section{}, fmt.Errorf("invalid section fields or types")
	}
	section := domain.Section{Tables: output.Tables, ID: output.ID, Heading: plan.Heading, Body: strings.TrimSpace(output.Body), Citations: output.Citations, Missing: output.Missing}
	if err := validateTables(section.Tables); err != nil {
		return section, err
	}
	if section.ID != plan.ID {
		return section, fmt.Errorf("unknown section ID")
	}
	if section.Body == "" || !safeText(section.Body) || len(section.Body) > 40000 {
		return section, fmt.Errorf("body must be nonempty valid text below 40000 bytes")
	}
	if section.Citations == nil {
		section.Citations = []string{}
	}
	if section.Missing == nil {
		section.Missing = []string{}
	}
	allowed := map[string]bool{}
	for _, f := range facts {
		allowed[f.ID] = true
	}
	for _, c := range section.Citations {
		if !allowed[c] {
			return section, fmt.Errorf("citation %q is not in the supplied source pack", clipText(c, 80))
		}
	}
	for _, gap := range section.Missing {
		if !safeText(gap) || len(gap) > 4000 {
			return section, fmt.Errorf("invalid missing_information entry")
		}
	}
	if len(section.Citations) == 0 && len(section.Missing) == 0 {
		return section, fmt.Errorf("provide supporting citations or a specific evidence gap")
	}
	return section, nil
}

// DocumentEvidence exposes the same user-attributed facts to generation, review and export.
func DocumentEvidence(evidence domain.Evidence, input domain.DocumentInput) domain.Evidence {
	evidence.Facts = append(append([]domain.Fact(nil), evidence.Facts...), domain.Fact{ID: "project_title", Kind: "project_title", Value: input.Metadata["project"], Provenance: "user"})
	if input.Specification != "" {
		evidence.Facts = append(evidence.Facts, domain.Fact{ID: "technical_specification", Kind: "specification", Value: input.Specification, Provenance: "user"})
	}
	if value := evidence.Corrections["additional_context"]; value != "" {
		evidence.Facts = append(evidence.Facts, domain.Fact{ID: "additional_context", Kind: "project_context", Value: value, Provenance: "user"})
	}
	return evidence
}

// completeMembers recovers only fully decoded objects from a truncated sections array.
func completeMembers(raw string) []json.RawMessage {
	decoder := json.NewDecoder(strings.NewReader(raw))
	token, err := decoder.Token()
	if err != nil || token != json.Delim('{') {
		return nil
	}
	for decoder.More() {
		key, err := decoder.Token()
		if err != nil {
			return nil
		}
		if key != "sections" {
			var ignored json.RawMessage
			if decoder.Decode(&ignored) != nil {
				return nil
			}
			continue
		}
		token, err = decoder.Token()
		if err != nil || token != json.Delim('[') {
			return nil
		}
		members := []json.RawMessage{}
		for decoder.More() {
			var member json.RawMessage
			if decoder.Decode(&member) != nil {
				break
			}
			members = append(members, member)
		}
		return members
	}
	return nil
}

func validateTables(tables []domain.DocumentTable) error {
	if len(tables) > 5 {
		return fmt.Errorf("at most five tables per section")
	}
	for _, t := range tables {
		if len(t.Columns) < 1 || len(t.Columns) > 6 || len(t.Rows) > 40 || !safeText(t.Caption) || len(t.Caption) > 500 {
			return fmt.Errorf("invalid table dimensions or caption")
		}
		cells := append([]string{}, t.Columns...)
		for _, row := range t.Rows {
			if len(row) != len(t.Columns) {
				return fmt.Errorf("table row width differs from columns")
			}
			cells = append(cells, row...)
		}
		for _, cell := range cells {
			if !safeText(cell) || len(cell) > 2000 {
				return fmt.Errorf("invalid table cell")
			}
		}
	}
	return nil
}
