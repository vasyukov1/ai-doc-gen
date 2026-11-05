package web

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/vasyukov1/ai-doc-gen/internal/documents"
	"github.com/vasyukov1/ai-doc-gen/internal/domain"
)

func connectionsFromForm(r *http.Request) ([]domain.Connection, error) {
	connections := []domain.Connection{}
	for _, provider := range []string{"gemini", "groq", "openrouter"} {
		keys := splitKeys(r.FormValue(provider + "_key"))
		if len(keys) == 0 {
			continue
		}
		if r.FormValue(provider+"_free") != "yes" {
			return nil, fmt.Errorf("confirm the free account conditions for %s", provider)
		}
		rpm, tpm := 5, 15000
		var err error
		if value := r.FormValue(provider + "_rpm"); value != "" {
			rpm, err = strconv.Atoi(value)
			if err != nil {
				return nil, fmt.Errorf("invalid RPM for %s", provider)
			}
		}
		if value := r.FormValue(provider + "_tpm"); value != "" {
			tpm, err = strconv.Atoi(value)
			if err != nil {
				return nil, fmt.Errorf("invalid TPM for %s", provider)
			}
		}
		for _, key := range keys {
			c := domain.Connection{Provider: provider, Key: key, Model: strings.TrimSpace(r.FormValue(provider + "_model")), RPM: rpm, TPM: tpm}
			if err := documents.ValidateConnection(c); err != nil {
				return nil, fmt.Errorf("%s: %w", provider, err)
			}
			connections = append(connections, c)
		}
	}
	if r.FormValue("local_enabled") == "yes" {
		model := strings.TrimSpace(r.FormValue("local_model"))
		if model == "" {
			model = "qwen2.5:7b-instruct"
		}
		rpm, tpm := 1, 24000
		var err error
		if value := r.FormValue("local_rpm"); value != "" {
			rpm, err = strconv.Atoi(value)
			if err != nil {
				return nil, fmt.Errorf("invalid RPM for the local model")
			}
		}
		if value := r.FormValue("local_tpm"); value != "" {
			tpm, err = strconv.Atoi(value)
			if err != nil {
				return nil, fmt.Errorf("invalid TPM for the local model")
			}
		}
		c := domain.Connection{Provider: "local", Model: model, RPM: rpm, TPM: tpm}
		if err := documents.ValidateConnection(c); err != nil {
			return nil, fmt.Errorf("local: %w", err)
		}
		connections = append(connections, c)
	}

	return connections, nil
}

// splitKeys accepts one key or several separated by commas, semicolons or
// newlines. Each key gets its own rate window, so free-tier limits stack.
func splitKeys(value string) []string {
	seen := map[string]bool{}
	keys := []string{}
	for _, part := range strings.FieldsFunc(value, func(r rune) bool {
		return r == ',' || r == ';' || r == '\n' || r == '\r'
	}) {
		part = strings.TrimSpace(part)
		if part == "" || seen[part] {
			continue
		}
		seen[part] = true
		keys = append(keys, part)
		if len(keys) == 5 {
			break
		}
	}
	return keys
}
