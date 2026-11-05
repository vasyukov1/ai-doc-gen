package documents

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/vasyukov1/ai-doc-gen/internal/domain"
)

func ProviderEndpoint(provider string) string {
	switch provider {
	case "gemini":
		return "https://generativelanguage.googleapis.com/v1beta/openai/chat/completions"
	case "groq":
		return "https://api.groq.com/openai/v1/chat/completions"
	case "openrouter":
		return "https://openrouter.ai/api/v1/chat/completions"
	}
	return ""
}

func ValidateConnection(c domain.Connection) error {
	if ProviderEndpoint(c.Provider) == "" || len(c.Key) < 8 || len(c.Key) > 1000 || len(c.Model) == 0 || len(c.Model) > 150 || strings.ContainsAny(c.Model, "\r\n\t ") {
		return errors.New("enter a provider model ID and API key")
	}
	if c.RPM < 1 || c.RPM > 60 || c.TPM < 4000 || c.TPM > 100000 {
		return errors.New("use RPM 1-60 and TPM 4000-100000 from your free account limits")
	}
	return nil
}

// Verify checks model existence and OpenRouter's published zero-price constraint.
func (p *ChatProvider) Verify(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 8*time.Second)
	defer cancel()
	endpoint := strings.TrimSuffix(p.Endpoint, "chat/completions") + "models"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+p.Key)
	client := p.Client
	if client == nil {
		client = &http.Client{Timeout: 8 * time.Second, CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}
	}
	resp, err := client.Do(req)
	if err != nil {
		return &ProviderError{Code: "provider_unavailable", Message: "Could not verify model availability."}
	}
	defer resp.Body.Close()
	p.Cooldown = responseCooldown(resp.Header)
	p.Diagnostics = CallDiagnostics{Stage: "catalog", HTTPStatus: resp.StatusCode, RequestID: clipText(resp.Header.Get("X-Request-ID"), 200), RetryAfter: clipText(resp.Header.Get("Retry-After"), 100)}
	if resp.StatusCode != 200 {
		return &ProviderError{Code: "provider_catalog", Message: "Model list unavailable. Check the provider key and free account."}
	}
	raw, err := io.ReadAll(io.LimitReader(resp.Body, 4<<20))
	if err != nil {
		return err
	}
	var catalog struct {
		Data []struct {
			ID         string            `json:"id"`
			Pricing    map[string]string `json:"pricing"`
			Parameters []string          `json:"supported_parameters"`
		} `json:"data"`
	}
	if json.Unmarshal(raw, &catalog) != nil {
		return invalidOutput()
	}
	for _, model := range catalog.Data {
		if model.ID != p.Model {
			continue
		}
		if p.FreeOnly {
			for _, field := range []string{"prompt", "completion"} {
				value, err := strconv.ParseFloat(model.Pricing[field], 64)
				if err != nil || value != 0 {
					return &ProviderError{Code: "provider_not_free", Message: "Selected OpenRouter model is not verified as free."}
				}
			}
			for field, value := range model.Pricing {
				if field == "discount" {
					continue
				}
				price, err := strconv.ParseFloat(value, 64)
				if err != nil || price != 0 {
					return &ProviderError{Code: "provider_not_free", Message: "Selected model has nonzero or unknown charges."}
				}
			}
			supported := false
			for _, parameter := range model.Parameters {
				if parameter == "structured_outputs" {
					p.Schema = true
				}
				if parameter == "response_format" || parameter == "structured_outputs" {
					supported = true
				}
			}
			if !supported {
				return &ProviderError{Code: "provider_format", Message: "Selected model does not advertise structured JSON output."}
			}
		}
		return nil
	}
	return &ProviderError{Code: "provider_model_missing", Message: "Selected model is not available in the provider model list."}
}

func responseCooldown(header http.Header) time.Duration {
	delay := time.Duration(0)
	if value := header.Get("Retry-After"); value != "" {
		delay = retryAfter(value)
	}
	for _, kind := range []string{"requests", "tokens"} {
		if header.Get("x-ratelimit-remaining-"+kind) != "0" {
			continue
		}
		value := header.Get("x-ratelimit-reset-" + kind)
		d, err := time.ParseDuration(value)
		if err != nil {
			if seconds, e := strconv.ParseFloat(value, 64); e == nil {
				d = time.Duration(seconds * float64(time.Second))
				if seconds > 1e9 {
					d = time.Until(time.Unix(int64(seconds), 0))
				}
			}
		}
		if d <= 0 {
			d = time.Minute
		}
		if d > delay {
			delay = d
		}
	}
	return delay
}
