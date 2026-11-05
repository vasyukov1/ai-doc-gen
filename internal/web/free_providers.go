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
		key := strings.TrimSpace(r.FormValue(provider + "_key"))
		if key == "" {
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
		c := domain.Connection{Provider: provider, Key: key, Model: strings.TrimSpace(r.FormValue(provider + "_model")), RPM: rpm, TPM: tpm}
		if err := documents.ValidateConnection(c); err != nil {
			return nil, fmt.Errorf("%s: %w", provider, err)
		}
		connections = append(connections, c)
	}
	return connections, nil
}
