package documents

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/vasyukov1/ai-doc-gen/internal/domain"
	"github.com/vasyukov1/ai-doc-gen/internal/pgqueue"
	"github.com/vasyukov1/ai-doc-gen/internal/storage"
)

func (w *Worker) generate(ctx, parent context.Context, claim *pgqueue.Claim, evidence domain.Evidence) {
	run, err := w.jobs.DraftRun(ctx, claim.ID, w.encryptionKey)
	if err == nil {
		w.generateDraft(ctx, parent, claim, evidence, run)
		return
	}
	if !errors.Is(err, sql.ErrNoRows) {
		w.generationFailure(parent, claim, err)
		return
	}
	if err := w.queue.Stage(ctx, claim.ID, w.owner, "generating", "generating", "Generating grounded draft sections."); err != nil {
		return
	}
	job, err := w.jobs.GetJob(ctx, claim.ID)
	if err != nil {
		w.generationFailure(parent, claim, err)
		return
	}
	documentInput, inputErr := w.jobs.DocumentInput(ctx, claim.ID)
	if inputErr != nil {
		w.generationFailure(parent, claim, &ProviderError{Code: "document_input", Message: "Provide document metadata and a technical specification before generation."})
		return
	}
	if _, inputErr = ValidateDocumentInput(documentInput, job.ProfileID, job.Language); inputErr != nil {
		w.generationFailure(parent, claim, &ProviderError{Code: "document_input", Message: inputErr.Error()})
		return
	}
	if documentInput.Specification != "" {
		evidence.Facts = append([]domain.Fact{{ID: "technical_specification", Kind: "specification", Value: documentInput.Specification, Provenance: "user"}}, evidence.Facts...)
	}

	sections, err := w.jobs.Sections(ctx, claim.ID)
	if err != nil {
		w.generationFailure(parent, claim, err)
		return
	}
	existing := map[string]bool{}
	for _, section := range sections {
		existing[section.ID] = true
	}
	model, key, credentialErr := w.jobs.Credentials(ctx, claim.ID, w.encryptionKey)
	if credentialErr != nil {
		w.generationFailure(parent, claim, &ProviderError{Code: "provider_credentials", Message: "Enter an OpenRouter key and configure the server encryption key before generation."})
		return
	}
	adapter := Provider(&ChatProvider{Endpoint: "https://openrouter.ai/api/v1/chat/completions", Model: model, Key: key})
	if w.provider != nil {
		if endpoint, _ := w.provider.Identity(); endpoint != "" {
			adapter = w.provider
		}
	}
	provider, model := adapter.Identity()
	provider += "/account/" + storage.CredentialScope(key)
	plan := Plan(job.ProfileID, job.Language)
	if len(plan) == 0 {
		w.generationFailure(parent, claim, &ProviderError{Code: "profile_obsolete", Message: "This legacy profile is unavailable. Submit the project using a specific document variant."})
		return
	}
	if err := w.queue.SectionTotal(ctx, claim.ID, w.owner, len(plan)); err != nil {
		w.generationFailure(parent, claim, err)
		return
	}
	for position, item := range plan {
		if existing[item.ID] {
			continue
		}
		key := SectionKey(evidence, job, item, provider, model)
		facts := Relevant(evidence, item)
		section, err := w.jobs.CachedSection(ctx, key)
		if err == nil {
			raw := section
			raw.CacheKey = ""
			raw.Edited = false
			// Cache records are checked again against the current evidence.
			encoded := sectionJSON(raw)
			section, err = ValidateSection(encoded, item, facts)
		}
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, storage.ErrInvalidCache) || (err != nil && providerFailure(err).Code == "provider_invalid_output") {
			output, usage, callErr := adapter.Generate(ctx, SectionPrompt(item, job.Language, facts))
			code := ""
			if callErr == nil {
				section, callErr = ValidateSection(output, item, facts)
			}
			if callErr != nil {
				code = providerFailure(callErr).Code
			}
			if saveErr := w.jobs.RecordUsage(ctx, claim.ID, item.ID, provider, model, code, usage); saveErr != nil {
				w.generationFailure(parent, claim, saveErr)
				return
			}
			err = callErr
		}
		if err != nil {
			w.generationFailure(parent, claim, err)
			return
		}
		section.CacheKey = key
		if err := w.queue.SaveSection(ctx, claim.ID, w.owner, position, len(plan), section); err != nil {
			w.generationFailure(parent, claim, err)
			return
		}
	}
	if err := w.queue.DraftReady(ctx, claim.ID, w.owner); err != nil {
		log.Printf("finish draft %s: %v", claim.ID, err)
		return
	}

	// Drafts are rendered without a manual review step so a submitted project
	// returns a finished document. Editing a completed draft still asks for an
	// explicit rebuild.
	renderCtx, cancel := context.WithTimeout(parent, 10*time.Second)
	defer cancel()
	if err := w.jobs.QueueRender(renderCtx, claim.ID, ""); err != nil {
		log.Printf("auto-render %s: %v", claim.ID, err)
	}
}

func (w *Worker) generationFailure(ctx context.Context, claim *pgqueue.Claim, err error) {
	failure := providerFailure(err)
	if !failure.Temporary || claim.Attempts >= 3 {
		_ = w.queue.Fail(ctx, claim.ID, w.owner, failure.Code, failure.Message)
		return
	}
	delay := failure.RetryAfter
	if delay == 0 {
		delay = time.Duration(claim.Attempts*5) * time.Second
	}
	_ = w.queue.Retry(ctx, claim.ID, w.owner, failure.Code, failure.Message, delay)
}

func sectionJSON(section domain.Section) string {
	raw, _ := json.Marshal(sectionOutput{ID: section.ID, Heading: section.Heading, Body: section.Body, Citations: section.Citations, Missing: section.Missing, Tables: section.Tables})
	return string(raw)
}
