package documents

import (
	"encoding/json"
	"fmt"
	"regexp"
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

const DraftVersion = "chapters-8"

func clipText(value string, limit int) string {
	if limit <= 0 {
		return ""
	}
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
	words := strings.FieldsFunc(strings.ToLower(plan.Topic+" "+plan.Heading+" "+plan.Guidance), func(r rune) bool { return !unicode.IsLetter(r) })
	result := []string{}
	for _, w := range words {
		if len([]rune(w)) >= 4 {
			result = append(result, w)
		}
	}
	return result
}

// CompactFacts ranks specific evidence before general README context. BM25
// supplies the relevance signal; kind bonuses keep section intent in control.
func CompactFacts(evidence domain.Evidence, plan SectionPlan) []domain.Fact {
	type candidate struct {
		fact  domain.Fact
		score int
	}
	topic := strings.ToLower(plan.Topic)
	entries := []candidate{}
	seen := map[string]bool{}
	for _, fact := range evidence.Facts {
		if fact.Provenance == "unknown" || fact.Kind == "selected_file" {
			continue
		}
		if correction, ok := evidence.Corrections[fact.ID]; ok {
			fact.Value = correction
			fact.Provenance = "user"
		}
		score := 0
		if strings.Contains(fact.Kind, "readme") && (strings.Contains(topic, "overview") || strings.Contains(topic, "purpose") || strings.Contains(topic, "application")) {
			score += 40
		}
		if strings.Contains(fact.Kind, "api") && (strings.Contains(topic, "interface") || strings.Contains(topic, "input") || strings.Contains(topic, "output")) {
			score += 40
		}
		for _, kind := range plan.Kinds {
			if kind != "readme" && kind != "description" && kind != "unknown" && strings.Contains(strings.ToLower(fact.Kind+" "+fact.Path), kind) {
				score += 20
			}
		}
		if strings.Contains(topic, "name") && fact.Kind == "project_title" {
			score += 200
		}
		if fact.Kind == "specification" {
			best := specificationExcerpt(fact.Value, plan)
			if best == "" {
				continue
			}
			fact.Value = best
		}
		fact.Value = clipText(fact.Value, 450)
		if seen[fact.Value] {
			continue
		}
		seen[fact.Value] = true
		entries = append(entries, candidate{fact, score})
	}
	texts := make([]string, len(entries))
	for index, entry := range entries {
		texts[index] = entry.fact.Kind + " " + entry.fact.Path + " " + entry.fact.Value
	}
	bm := newBM25(texts)
	query := topicWords(plan)
	ranked := make([]candidate, 0, len(entries))
	for index, entry := range entries {
		entry.score += int(bm.score(query, index) * 10)
		if entry.score > 0 {
			ranked = append(ranked, entry)
		}
	}
	sort.SliceStable(ranked, func(i, j int) bool { return ranked[i].score > ranked[j].score })
	facts := []domain.Fact{}
	for _, entry := range ranked {
		facts = append(facts, entry.fact)
		if len(facts) == 4 {
			break
		}
	}
	// General project descriptions belong in overview/purpose sections only.
	if len(facts) == 0 && (strings.Contains(topic, "overview") || strings.Contains(topic, "purpose") || strings.Contains(topic, "application")) {
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

// specificationExcerpt selects the most relevant sentence instead of repeating
// the entire technical specification in every section.
func specificationExcerpt(value string, plan SectionPlan) string {
	best := ""
	bestScore := 0
	for _, line := range strings.FieldsFunc(value, func(r rune) bool { return r == '\n' || r == '.' || r == ';' }) {
		score := 0
		for _, word := range topicWords(plan) {
			if strings.Contains(strings.ToLower(line), word) {
				score++
			}
		}
		if score > bestScore {
			bestScore = score
			best = line
		}
	}
	return strings.TrimSpace(best)
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
	if plan.Language == "ca" {
		gap = "Cal revisar la secció «" + plan.Heading + "»: els fragments no estableixen requisits complets ni un comportament verificat en execució."
		if len(facts) == 0 {
			gap = "No s’han trobat proves específiques per a «" + plan.Heading + "». Aporteu els requisits o les referències a les fonts pertinents."
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
	section := domain.Section{Tables: output.Tables, ID: output.ID, Heading: plan.Heading, Body: strings.TrimSpace(normalizeMarkdown(output.Body)), Citations: output.Citations, Missing: output.Missing}
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
	section.Body, section.Citations = normalizeEvidenceMarkers(section.Body, section.Citations, facts)
	section.Body = strings.TrimSpace(section.Body)
	if section.Body == "" {
		return section, fmt.Errorf("body is empty after citation normalization")
	}
	marker := regexp.MustCompile(`\[\^?fact-[^\]]+\]`)
	clean := func(value string) (string, bool) {
		value, section.Citations = normalizeEvidenceMarkers(value, section.Citations, facts)
		return value, !marker.MatchString(value)
	}
	if marker.MatchString(section.Body) {
		return section, fmt.Errorf("body contains unsupported evidence markers")
	}
	for i := range section.Tables {
		table := &section.Tables[i]
		var valid bool
		table.Caption, valid = clean(table.Caption)
		if !valid {
			return section, fmt.Errorf("table contains unsupported evidence markers")
		}
		for j := range table.Columns {
			table.Columns[j], valid = clean(table.Columns[j])
			if !valid {
				return section, fmt.Errorf("table contains unsupported evidence markers")
			}
		}
		for j := range table.Rows {
			for k := range table.Rows[j] {
				table.Rows[j][k], valid = clean(table.Rows[j][k])
				if !valid {
					return section, fmt.Errorf("table contains unsupported evidence markers")
				}
			}
		}
	}
	// A visible placeholder is a gap even if the model omitted its metadata.
	visible := section.Body
	for _, table := range section.Tables {
		visible += "\n" + table.Caption + "\n" + strings.Join(table.Columns, "\n")
		for _, row := range table.Rows {
			visible += "\n" + strings.Join(row, "\n")
		}
	}
	for _, line := range strings.Split(visible, "\n") {
		if index := strings.Index(line, "TODO:"); index >= 0 {
			gap := clipText(strings.TrimSpace(line[index+5:]), 1000)
			if gap != "" {
				found := false
				for _, existing := range section.Missing {
					if strings.TrimSpace(existing) != "" && (strings.Contains(existing, gap) || strings.Contains(gap, existing)) {
						found = true
					}
				}
				if !found {
					section.Missing = append(section.Missing, gap)
				}
			}
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

// Move only supplied evidence markers into citation metadata; preserve ordinary brackets.
func normalizeEvidenceMarkers(body string, citations []string, facts []domain.Fact) (string, []string) {
	allowed, seen := map[string]bool{}, map[string]bool{}
	for _, fact := range facts {
		allowed[fact.ID] = true
	}
	for _, id := range citations {
		seen[id] = true
	}
	marker := regexp.MustCompile(`\[\^?([a-zA-Z0-9_,-]+(?:\s+[a-zA-Z0-9_,-]+)*)\]`)
	var out strings.Builder
	previous := 0
	for _, location := range marker.FindAllStringIndex(body, -1) {
		value := body[location[0]:location[1]]
		if location[1] < len(body) && body[location[1]] == '(' {
			continue
		}
		inside := strings.TrimPrefix(strings.TrimSuffix(strings.TrimPrefix(value, "["), "]"), "^")
		ids := strings.FieldsFunc(inside, func(r rune) bool { return r == ',' || unicode.IsSpace(r) })
		valid := len(ids) > 0
		for _, id := range ids {
			if !allowed[id] {
				valid = false
			}
		}
		if !valid {
			continue
		}
		out.WriteString(body[previous:location[0]])
		previous = location[1]
		for _, id := range ids {
			if !seen[id] {
				citations = append(citations, id)
				seen[id] = true
			}
		}
	}
	out.WriteString(body[previous:])
	return out.String(), citations
}
