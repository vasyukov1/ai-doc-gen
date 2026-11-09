package documents

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/vasyukov1/ai-doc-gen/internal/domain"
	"github.com/vasyukov1/ai-doc-gen/internal/pgqueue"
	"github.com/vasyukov1/ai-doc-gen/internal/storage"
)

func (w *Worker) generateDraft(ctx, parent context.Context, claim *pgqueue.Claim, evidence domain.Evidence, run domain.DraftRun) {
	if err := w.buildDraft(ctx, claim, evidence, run); err != nil {
		if ctx.Err() == nil {
			w.generationFailure(parent, claim, err)
		}
		return
	}
	if err := w.queue.DraftReadyForRender(ctx, claim.ID, w.owner); err != nil && ctx.Err() == nil {
		w.generationFailure(parent, claim, err)
	}
}

func (w *Worker) buildDraft(ctx context.Context, claim *pgqueue.Claim, evidence domain.Evidence, run domain.DraftRun) error {
	job, err := w.jobs.GetJob(ctx, claim.ID)
	if err != nil {
		return err
	}
	input, err := w.jobs.DocumentInput(ctx, claim.ID)
	if err != nil {
		return err
	}
	if _, err = ValidateDocumentInput(input, job.ProfileID, job.Language); err != nil {
		return err
	}
	evidence = DocumentEvidence(evidence, input)
	plan := Plan(job.ProfileID, job.Language)
	if len(plan) == 0 {
		return errors.New("profile unavailable")
	}
	if err = w.queue.SectionTotal(ctx, claim.ID, w.owner, len(plan)); err != nil {
		return err
	}
	if err = w.queue.Stage(ctx, claim.ID, w.owner, "generating", "baseline", "Building a source-linked baseline before model requests."); err != nil {
		return err
	}
	saved, err := w.jobs.Sections(ctx, claim.ID)
	if err != nil {
		return err
	}
	existing := map[string]domain.Section{}
	for _, s := range saved {
		existing[s.ID] = s
	}
	facts := map[string][]domain.Fact{}
	positions := map[string]int{}
	pending := []SectionPlan{}
	for position, item := range plan {
		positions[item.ID] = position
		facts[item.ID] = CompactFacts(evidence, item)
		section, ok := existing[item.ID]
		if !ok {
			section = BaselineSection(item, facts[item.ID])
			if err = w.queue.SaveSection(ctx, claim.ID, w.owner, position, len(plan), section); err != nil {
				return err
			}
		}
		if section.Origin == "baseline" && !section.Edited && run.SectionAttempts[item.ID] < 2 {
			pending = append(pending, item)
		}
	}
	if len(run.Connections) == 0 {
		return w.queue.DraftNote(ctx, claim.ID, w.owner, "No providers connected; baseline draft generated from source evidence.", run.Disabled)
	}
	budget := DraftBudgetFor(job.ProfileID, job.Language)
	if err = w.queue.StartDraftBudget(ctx, claim.ID, w.owner, budget.Duration); err != nil {
		return err
	}
	// A fresh run receives the profile budget; a restarted run keeps its stored
	// deadline so a restart cannot renew model time.
	deadline := run.Deadline
	if run.Attempts == 0 {
		deadline = time.Now().Add(budget.Duration)
		if run.Deadline.After(deadline) {
			deadline = run.Deadline
		}
	}
	attemptLimit := budget.Attempts
	modelExpired := false
	coveredChapters := map[string]bool{}
	verified := map[string]bool{}
	schemas := map[string]bool{}
	// Free accounts publish input/output token windows that may be smaller than
	// the account TPM entered on the form. Learn them from provider rejections.
	inputCaps := map[string]int{}
	outputCaps := map[string]int{}
	if w.provider == nil {
		for _, c := range run.Connections {
			ok, schema, cacheErr := w.jobs.DraftCatalog(ctx, claim.ID, c.Provider, c.Model)
			if cacheErr != nil {
				return cacheErr
			}
			verified[c.Provider] = ok
			schemas[c.Provider] = schema
		}
	}
	reason := ""
	groups := chapterGroups(job.ProfileID, pending, evidence)
	chapterTotal := len(groups)
	previous, issues, err := w.jobs.DraftRepair(ctx, claim.ID)
	if err != nil {
		return err
	}
	// Coverage first: every pass gives each unfinished chapter at most one call,
	// so repairs on the first chapters cannot starve the rest of the document.
	for continuation := 0; ; continuation++ {
		modelCtx, cancel := context.WithDeadline(ctx, deadline)
		runPasses := func(modelCtx context.Context) error {
			for pass := 0; pass < 3; pass++ {
				calls := 0
				for groupIndex := 0; groupIndex < len(groups) && modelCtx.Err() == nil; groupIndex++ {
					group := groups[groupIndex]
					if len(group.Plans) == 0 {
						continue
					}
					callsLeft := 1
					for _, issue := range issues {
						for _, p := range group.Plans {
							if issue.SectionID == p.ID {
								group.Previous = clipText(previous, 12000)
								group.Issues = issues
							}
						}
					}
					nextProvider := 0
					waitReported := false
					for len(group.Plans) > 0 && modelCtx.Err() == nil && callsLeft > 0 {
						progress := false
						earliest := time.Duration(0)
						available := false
						for offset := range run.Connections {
							providerIndex := (nextProvider + offset) % len(run.Connections)
							connection := run.Connections[providerIndex]
							limitKey := connection.Provider + "/" + connection.Model
							if modelCtx.Err() != nil || run.Attempts >= attemptLimit {
								break
							}
							if run.Disabled[connection.Provider] != "" {
								continue
							}
							available = true
							local := connection.Provider == "local"
							timeout := time.Duration(0)
							if local {
								// A CPU-only model needs a longer window than a hosted endpoint.
								timeout = 180 * time.Second
							}
							adapter := Provider(&ChatProvider{Endpoint: ProviderEndpoint(connection.Provider), Model: connection.Model, Key: connection.Key, MaxTokens: 6000, FreeOnly: connection.Provider == "openrouter", Keyless: local, Timeout: timeout})
							if w.provider != nil {
								adapter = w.provider
								if p, ok := adapter.(*ChatProvider); ok {
									copy := *p
									copy.MaxTokens = 6000
									copy.Model = connection.Model
									adapter = &copy
								}
							}
							scope := connection.Provider + "/" + storage.CredentialScope(connection.Key)
							cacheProvider := ProviderEndpoint(connection.Provider) + "/account/" + storage.CredentialScope(connection.Key) + "/" + DraftVersion
							if w.provider != nil {
								endpoint, model := w.provider.Identity()
								cacheProvider += "/fixture/" + endpoint + "/" + model
							}
							// Reuse validated sections before making any network request.
							remaining := []SectionPlan{}
							for _, item := range group.Plans {
								cacheKey := SectionKey(evidence, job, item, cacheProvider, connection.Model)
								cached, cacheErr := w.jobs.CachedSection(ctx, cacheKey)
								if cacheErr != nil && !errors.Is(cacheErr, sql.ErrNoRows) && !errors.Is(cacheErr, storage.ErrInvalidCache) {
									return cacheErr
								}
								if cacheErr == nil && (cached.Origin != "model" || cached.Edited) {
									cacheErr = storage.ErrInvalidCache
								}
								if cacheErr == nil {
									cached, cacheErr = validateChapterSection(sectionJSON(cached), item, group.Facts[item.ID])
								}
								if cacheErr == nil {
									cached.Origin = "model"
									cached.Provider = connection.Provider
									cached.Model = connection.Model
									cached.CacheKey = cacheKey
									if err = w.queue.SaveSection(ctx, claim.ID, w.owner, positions[item.ID], len(plan), cached); err != nil {
										return err
									}
									progress = true
								} else {
									remaining = append(remaining, item)
								}
							}
							group.Plans = remaining
							if len(remaining) == 0 {
								break
							}
							if !verified[connection.Provider] {
								if w.provider == nil {
									wait, attempted, verifyErr := w.verifyDraftProvider(modelCtx, ctx, claim, scope, connection, adapter.(*ChatProvider), attemptLimit)
									if attempted {
										run.Attempts++
									}
									if wait > 0 {
										if earliest == 0 || wait < earliest {
											earliest = wait
										}
										continue
									}
									if errors.Is(verifyErr, pgqueue.ErrDraftBudget) {
										reason = "Generation budget reached; remaining sections retain their baseline."
										continue
									}
									if verifyErr != nil {
										if ctx.Err() != nil {
											return ctx.Err()
										}
										var failure *ProviderError
										if !errors.As(verifyErr, &failure) {
											if modelCtx.Err() != nil {
												break
											}
											return verifyErr
										}
										run.Disabled[connection.Provider] = failure.Code
										reason = connection.Provider + ": " + failure.Message
										if err = w.queue.DraftNote(ctx, claim.ID, w.owner, reason, run.Disabled); err != nil {
											return err
										}
										continue
									}
								}
								verified[connection.Provider] = true
								if p, ok := adapter.(*ChatProvider); ok {
									schemas[connection.Provider] = p.Schema
								}
							}
							if p, ok := adapter.(*ChatProvider); ok {
								p.Schema = schemas[connection.Provider]
							}
							ids := []string{}
							for _, item := range group.Plans {
								ids = append(ids, item.ID)
							}
							tpm := connection.TPM
							if cap := inputCaps[limitKey]; cap > 0 && cap < tpm {
								tpm = cap
							}
							outputTokens := min(6000, max(3000, len(group.Plans)*1800))
							outputTokens = min(outputTokens, tpm/2)
							if cap := outputCaps[limitKey]; cap > 0 {
								outputTokens = min(outputTokens, cap)
							}
							if outputTokens < 400 {
								reason = connection.Provider + ": the account token limits are too small for a chapter response."
								run.Disabled[connection.Provider] = "provider_context_budget"
								if err = w.queue.DraftNote(ctx, claim.ID, w.owner, reason, run.Disabled); err != nil {
									return err
								}
								continue
							}
							fitted := fitChapter(group, job.Language, tpm-outputTokens-256)
							prompt := GroupPrompt(fitted, job.Language)
							if estimateTokens(prompt)+outputTokens+128 > tpm || len(fitted.Sources)+len(fitted.Brief.Context) == 0 {
								reason = connection.Provider + ": account TPM is too small for this chapter and useful source context. Increase the configured limit only if your account supports it."
								run.Disabled[connection.Provider] = "provider_context_budget"
								if err = w.queue.DraftNote(ctx, claim.ID, w.owner, reason, run.Disabled); err != nil {
									return err
								}
								continue
							}
							if p, ok := adapter.(*ChatProvider); ok {
								p.MaxTokens = outputTokens
								p.SectionIDs = ids
								p.CitationIDs = nil
								seen := map[string]bool{}
								for _, fact := range append(append([]domain.Fact{}, fitted.Sources...), fitted.Brief.Context...) {
									if !seen[fact.ID] {
										p.CitationIDs = append(p.CitationIDs, fact.ID)
										seen[fact.ID] = true
									}
								}
							}
							reservedTokens := estimateTokens(prompt) + outputTokens + 128
							limited := connection
							limited.TPM = tpm
							callID, wait, reserveErr := w.queue.ReserveDraftCall(modelCtx, claim.ID, w.owner, scope, limited, reservedTokens, ids, attemptLimit)
							if errors.Is(reserveErr, pgqueue.ErrDraftBudget) {
								reason = "Generation budget reached; remaining sections retain their baseline."
								continue
							}
							if reserveErr != nil {
								if modelCtx.Err() != nil {
									break
								}
								return reserveErr
							}
							if wait > 0 {
								if earliest == 0 || wait < earliest {
									earliest = wait
								}
								continue
							}
							run.Attempts++
							callsLeft--
							calls++
							for _, id := range ids {
								run.SectionAttempts[id]++
							}
							if err = w.queue.Stage(ctx, claim.ID, w.owner, "generating", "model_group", fmt.Sprintf("Improving %d sections with %s; attempt %d of %d.", len(ids), connection.Provider, run.Attempts, attemptLimit)); err != nil {
								return err
							}
							began := time.Now()
							raw, usage, callErr := adapter.Generate(modelCtx, prompt)
							cooldown := time.Duration(0)
							if p, ok := adapter.(*ChatProvider); ok {
								cooldown = p.Cooldown
							}
							code := ""
							result := map[string]domain.Section{}
							if callErr == nil {
								result, group.Issues = parseChapter(raw, fitted)
								group.Previous = clipText(raw, 12000)
								if len(result) != len(group.Plans) {
									code = "provider_invalid_output"
									if len(group.Issues) > 0 {
										reason = connection.Provider + ": " + group.Issues[0].Reason
									}
								}
							}
							var failure *ProviderError
							if callErr != nil {
								failure = providerFailure(callErr)
								if failure.Code != "provider_invalid_output" {
									for _, id := range ids {
										run.SectionAttempts[id]--
									}
								}
								code = failure.Code
								reason = connection.Provider + ": " + failure.Message
								if failure.RetryAfter > cooldown {
									cooldown = failure.RetryAfter
								}
								if failure.InputLimit > 0 {
									inputCaps[limitKey] = failure.InputLimit
								}
								if failure.OutputLimit > 0 {
									outputCaps[limitKey] = failure.OutputLimit
								}
								if failure.Code == "provider_context_budget" {
									// The real window is smaller than the form TPM; halve the pack
									// so the next attempt has a chance to fit.
									bound := max(2000, estimateTokens(prompt)/2)
									if inputCaps[limitKey] == 0 || bound < inputCaps[limitKey] {
										inputCaps[limitKey] = bound
									}
								}
								if failure.Rejected != "" {
									group.Previous = clipText(failure.Rejected, 12000)
									group.Issues = []domain.ValidationIssue{{
										SectionID: ids[0],
										Reason:    "The previous response was not valid JSON. Return one complete JSON object with every requested section.",
									}}
								}
								if !failure.Temporary {
									run.Disabled[connection.Provider] = code
								}
							}
							diagnostic := CallDiagnostics{Stage: "generation", PromptVersion: DraftVersion, EstimatedTokens: estimateTokens(prompt)}
							if p, ok := adapter.(*ChatProvider); ok {
								diagnostic = p.Diagnostics
							}
							// A published per-minute window smaller than the form value wins;
							// fit later calls to it instead of retrying 429 responses.
							if diagnostic.LimitTokens > 0 && (inputCaps[limitKey] == 0 || diagnostic.LimitTokens < inputCaps[limitKey]) {
								inputCaps[limitKey] = diagnostic.LimitTokens
							}
							diagnostic.Issues = group.Issues
							diagnostic.ReservedTokens = reservedTokens
							rejected := ""
							if failure != nil && failure.Rejected != "" {
								rejected = clipText(failure.Rejected, 100000)
							}
							if len(group.Issues) > 0 && raw != "" {
								rejected = clipText(raw, 100000)
							}
							if ctx.Err() == nil {
								if err = w.queue.SaveCallDiagnostics(ctx, claim.ID, w.owner, callID, diagnostic, rejected); err != nil {
									return err
								}
							}
							finishCtx, finishCancel := context.WithTimeout(context.WithoutCancel(ctx), 3*time.Second)
							err = w.queue.FinishDraftCall(finishCtx, claim.ID, w.owner, scope, callID, usage, code, time.Since(began), cooldown)
							finishCancel()
							if err != nil {
								return err
							}
							if ctx.Err() != nil {
								return ctx.Err()
							}
							if len(result) > 0 {
								coveredChapters[group.Title] = true
							}
							remaining = nil
							for _, item := range group.Plans {
								if section, ok := result[item.ID]; ok {
									section.Provider = connection.Provider
									section.Model = connection.Model
									section.CacheKey = SectionKey(evidence, job, item, cacheProvider, connection.Model)
									if err = w.queue.SaveSection(ctx, claim.ID, w.owner, positions[item.ID], len(plan), section); err != nil {
										return err
									}
								} else if run.SectionAttempts[item.ID] < 2 {
									remaining = append(remaining, item)
								} else {
									reason = "Some sections could not be improved after two attempts; their baseline is retained."
								}
							}
							if diagnostic.FinishReason == "length" && len(remaining) > 1 {
								cut := (len(remaining) + 1) / 2
								tail := group
								tail.Plans = append([]SectionPlan{}, remaining[cut:]...)
								groups = append(groups, tail)
								remaining = remaining[:cut]
							}
							group.Plans = remaining
							nextProvider = providerIndex
							if callErr != nil {
								nextProvider = (providerIndex + 1) % len(run.Connections)
							}
							progress = true
							if err = w.queue.DraftNote(ctx, claim.ID, w.owner, reason, run.Disabled); err != nil {
								return err
							}
							break
						}
						if len(group.Plans) == 0 {
							break
						}
						if !progress {
							if !available || earliest == 0 || run.Attempts >= attemptLimit {
								break
							}
							if !waitReported {
								if err = w.queue.Stage(ctx, claim.ID, w.owner, "generating", "provider_wait", "Waiting for a free provider slot; baseline sections are already available."); err != nil {
									return err
								}
								waitReported = true
							}
							timer := time.NewTimer(earliest)
							select {
							case <-modelCtx.Done():
								timer.Stop()
							case <-timer.C:
							}
						}
					}
					groups[groupIndex] = group
				}
				if calls == 0 {
					break
				}
			}
			return nil
		}
		err = runPasses(modelCtx)
		modelExpired = modelCtx.Err() != nil
		cancel()
		if err != nil {
			return err
		}
		if ctx.Err() != nil {
			return ctx.Err()
		}
		if !modelExpired {
			break
		}
		remaining := 0
		for _, group := range groups {
			remaining += len(group.Plans)
		}
		if remaining == 0 || continuation >= MaxContinuationPasses {
			break
		}
		if err = w.queue.ExtendDraftDeadline(ctx, claim.ID, w.owner, ContinuationPassDuration); err != nil {
			return err
		}
		attemptLimit += min(remaining, maxContinuationAttempts)
		deadline = time.Now().Add(ContinuationPassDuration)
		reason = "Model time budget reached; continuing with a bounded extra pass for unfinished sections."
	}
	if ctx.Err() != nil {
		return ctx.Err()
	}
	final, err := w.jobs.Sections(ctx, claim.ID)
	if err != nil {
		return err
	}
	baseline := 0
	for _, s := range final {
		if s.Origin == "baseline" {
			baseline++
		}
	}
	switch {
	case baseline == 0:
		reason = "All sections have model improvements or preserved edits."
	case modelExpired:
		reason = "Model time budget reached; PDF uses the saved baseline and validated improvements."
	case run.Attempts >= attemptLimit:
		reason = "Generation budget reached; remaining sections retain their baseline."
	case strings.TrimSpace(reason) == "":
		reason = "Some sections remain a baseline draft; review the marked gaps."
	}
	if run.Attempts > 0 {
		reason += fmt.Sprintf(" Model coverage: %d of %d chapters; %d of %d model calls used.", len(coveredChapters), max(chapterTotal, len(coveredChapters)), run.Attempts, attemptLimit)
	}
	if len(run.Disabled) > 0 {
		providers := make([]string, 0, len(run.Disabled))
		for provider, code := range run.Disabled {
			providers = append(providers, provider+" ("+code+")")
		}
		sort.Strings(providers)
		reason += " Disabled providers: " + strings.Join(providers, ", ") + "."
	}
	return w.queue.DraftNote(ctx, claim.ID, w.owner, reason, run.Disabled)
}

func (w *Worker) verifyDraftProvider(modelCtx, ctx context.Context, claim *pgqueue.Claim, scope string, c domain.Connection, p *ChatProvider, attemptLimit int) (time.Duration, bool, error) {
	callID, wait, err := w.queue.ReserveDraftCall(modelCtx, claim.ID, w.owner, scope, c, 0, nil, attemptLimit)
	if err != nil || wait > 0 {
		return wait, false, err
	}
	began := time.Now()
	verifyErr := p.Verify(modelCtx)
	code := "catalog_ok"
	p.Diagnostics.Stage = "catalog"
	p.Diagnostics.JSONSchema = p.Schema
	if verifyErr != nil {
		code = providerFailure(verifyErr).Code
	}
	if ctx.Err() == nil {
		if err = w.queue.SaveCallDiagnostics(ctx, claim.ID, w.owner, callID, p.Diagnostics, ""); err != nil {
			return 0, true, err
		}
	}
	finishCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), 3*time.Second)
	defer cancel()
	if err = w.queue.FinishDraftCall(finishCtx, claim.ID, w.owner, scope, callID, domain.Usage{}, code, time.Since(began), p.Cooldown); err != nil {
		return 0, true, err
	}
	return 0, true, verifyErr
}
