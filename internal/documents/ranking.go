package documents

import (
	"math"
	"strings"
	"unicode"
)

// bm25 ranks source facts by term relevance instead of substring luck.
type bm25 struct {
	docs  [][]string
	df    map[string]int
	avgdl float64
}

func newBM25(texts []string) *bm25 {
	b := &bm25{docs: make([][]string, len(texts)), df: make(map[string]int, len(texts))}
	total := 0
	for index, text := range texts {
		terms := tokenize(text)
		b.docs[index] = terms
		total += len(terms)
		seen := make(map[string]bool, len(terms))
		for _, term := range terms {
			if !seen[term] {
				seen[term] = true
				b.df[term]++
			}
		}
	}
	if len(texts) > 0 {
		b.avgdl = float64(total) / float64(len(texts))
	}

	return b
}

func (b *bm25) score(query []string, index int) float64 {
	if index < 0 || index >= len(b.docs) || b.avgdl == 0 || len(query) == 0 {
		return 0
	}
	const k1, saturation = 1.2, 0.75
	counts := make(map[string]int, len(b.docs[index]))
	for _, term := range b.docs[index] {
		counts[term]++
	}
	documents := float64(len(b.docs))
	length := float64(len(b.docs[index]))
	result := 0.0
	seen := make(map[string]bool, len(query))
	for _, term := range query {
		if seen[term] {
			continue
		}
		seen[term] = true
		frequency := float64(counts[term])
		if frequency == 0 {
			continue
		}
		appearances := float64(b.df[term])
		idf := math.Log(1 + (documents-appearances+0.5)/(appearances+0.5))
		result += idf * frequency * (k1 + 1) / (frequency + k1*(1-saturation+saturation*length/b.avgdl))
	}

	return result
}

func tokenize(text string) []string {
	fields := strings.FieldsFunc(strings.ToLower(text), func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsDigit(r)
	})
	terms := fields[:0]
	for _, field := range fields {
		if len(field) >= 3 {
			terms = append(terms, field)
		}
	}

	return terms
}

// queryWords collects unique topic words from the section plans.
func queryWords(plans []SectionPlan) []string {
	seen := make(map[string]bool, len(plans)*4)
	words := []string{}
	for _, plan := range plans {
		for _, word := range topicWords(plan) {
			if !seen[word] {
				seen[word] = true
				words = append(words, word)
			}
		}
	}

	return words
}
