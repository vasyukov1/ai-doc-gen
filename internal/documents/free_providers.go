package documents

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
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
	case "local":
		if endpoint := strings.TrimSpace(os.Getenv("LOCAL_LLM_ENDPOINT")); endpoint != "" {
			return endpoint
		}
		return "http://127.0.0.1:11434/v1/chat/completions"
	}
	return ""
}

func ValidateConnection(c domain.Connection) error {
	if ProviderEndpoint(c.Provider) == "" {
		return errors.New("choose a supported provider")
	}
	if len(c.Model) == 0 || len(c.Model) > 150 || strings.ContainsAny(c.Model, "\r\n\t ") {
		if c.Provider == "local" {
			return errors.New("enter the local model ID")
		}
		return errors.New("enter a provider model ID and API key")
	}
	if c.Provider == "local" {
		if c.RPM < 1 || c.RPM > 60 || c.TPM < 1000 || c.TPM > 100000 {
			return errors.New("use RPM 1-60 and TPM 1000-100000 for the local model")
		}
		return nil
	}
	if len(c.Key) < 8 || len(c.Key) > 1000 {
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
		message := "Could not verify model availability."
		if p.Keyless {
			message = "The local model endpoint did not answer. Start Ollama and pull the model."
		}
		return &ProviderError{Code: "provider_unavailable", Message: message}
	}
	defer resp.Body.Close()
	p.Cooldown = responseCooldown(resp.Header)
	p.Diagnostics = CallDiagnostics{Stage: "catalog", HTTPStatus: resp.StatusCode, RequestID: clipText(resp.Header.Get("X-Request-ID"), 200), RetryAfter: clipText(resp.Header.Get("Retry-After"), 100)}
	if resp.StatusCode != 200 {
		message := "Model list unavailable. Check the provider key and free account."
		if p.Keyless {
			message = "The local model endpoint rejected the model list request. Check LOCAL_LLM_ENDPOINT."
		}
		return &ProviderError{Code: "provider_catalog", Message: message}
	}
	raw, err := io.ReadAll(io.LimitReader(resp.Body, 4<<20))
	if err != nil {
		return err
	}
	var catalog struct {
		Data []struct {
			ID         string                     `json:"id"`
			Pricing    map[string]json.RawMessage `json:"pricing"`
			Parameters []string                   `json:"supported_parameters"`
		} `json:"data"`
	}
	if json.Unmarshal(raw, &catalog) != nil {
		return &ProviderError{Code: "provider_catalog", Message: "The provider model catalog could not be parsed. Retry later or choose another provider."}
	}
	available := make([]string, 0, len(catalog.Data))
	for _, model := range catalog.Data {
		available = append(available, model.ID)
		if !sameModelID(model.ID, p.Model) {
			continue
		}
		if p.FreeOnly {
			for _, field := range []string{"prompt", "completion"} {
				value, err := catalogPrice(model.Pricing[field])
				if err != nil || value != 0 {
					return &ProviderError{Code: "provider_not_free", Message: "Selected OpenRouter model is not verified as free."}
				}
			}
			for field, value := range model.Pricing {
				if field == "discount" {
					continue
				}
				price, err := catalogPrice(value)
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
		// These Gemini text models support the existing chapter JSON schema.
		if p.Endpoint == ProviderEndpoint("gemini") && (sameModelID(p.Model, "gemini-2.5-flash") || sameModelID(p.Model, "gemini-3.5-flash-lite") || sameModelID(p.Model, "gemini-3.8-flash")) {
			p.Schema = true
		}
		return nil
	}
	if p.Keyless {
		if len(available) == 0 {
			return &ProviderError{Code: "provider_model_missing", Message: "Ollama is reachable but has no models installed. Run \"ollama pull " + p.Model + "\" (or make models with Docker Compose), then retry."}
		}
		return &ProviderError{Code: "provider_model_missing", Message: "The local model is not installed. Available models: " + clipText(strings.Join(available, ", "), 300) + "."}
	}
	return &ProviderError{Code: "provider_model_missing", Message: "Selected model is not available in the provider model list."}
}

// sameModelID treats "model:latest" and "model" as the same local tag.
func sameModelID(available, requested string) bool {
	normalize := func(value string) string {
		value = strings.ToLower(strings.TrimSpace(value))
		// Gemini's OpenAI catalog prefixes IDs with "models/"; chat requests do not.
		value = strings.TrimPrefix(value, "models/")
		return strings.TrimSuffix(value, ":latest")
	}
	return normalize(available) == normalize(requested)
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

// Decode only scalar prices; unrelated catalog pricing metadata may be structured.
func catalogPrice(raw json.RawMessage) (float64, error) {
	var text string
	if json.Unmarshal(raw, &text) == nil {
		return strconv.ParseFloat(text, 64)
	}
	var number json.Number
	if json.Unmarshal(raw, &number) == nil && len(raw) > 0 && string(raw) != "null" {
		return number.Float64()
	}
	return 0, errors.New("unknown model price")
}
