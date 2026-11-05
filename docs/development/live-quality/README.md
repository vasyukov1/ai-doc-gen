# Live generation quality and reusable delivery

Date: 2026-11-06. Corrective milestone 9 work, bounded to Q1–Q5 in [the plan](plan.md). The uploaded public Go CLI was inspected as text; its code was never executed. Synthetic metadata identifies this as a demonstration. These are engineering observations on one fixture, not user-study results or a general model benchmark.

## Observed problem and changes

The old pipeline's controlled Gemini 2.5 Flash run returned no accepted model sections after one initial four-section pass. Three replies ended with `finish_reason=length` and incomplete JSON, despite short visible output; another call timed out. The run was cancelled after that pass (241.7 seconds), not allowed to exhaust every possible repair. A preliminary run using the previously configured model timed out and returned HTTP 503, so it is excluded from the content comparison.

Q1 increases response space, bounds each generation group to three leaves, and explicitly limits Gemini 2.5 Flash thinking through its supported OpenAI-compatible `extra_body.google.thinking_config`. It narrows section responsibilities, separates the documented software from its bundled examples, and requests concrete workflows rather than repeated introductions. The source pack ranks implementation facts with a soft repeated-file penalty; license text, import lists and call-name fragments no longer have equal priority. Compact manifest values remain eligible. Long declaration excerpts preserve beginning and ending with an explicit omission marker. These changes do not execute the project or establish runtime correctness.

Validation moves model-emitted, supplied evidence markers into citation metadata, including tables. Unknown fact markers and marker-only empty prose are rejected. It does not invent citations. Source/version cache keys change so old prose is not reused silently. HTTP timeout follows the selected provider timeout.

The first complete iteration had 4/4 model sections and 1,731 whitespace-separated words, but source review found an incorrect claim about a TODO fallback after an API error. The source actually returns the error. The next prompt revision explicitly distinguishes initialized templates from error recovery. An early context-budget edge case also caused a worker panic; a bounded-loop guard fixed it and lease recovery resumed the job. That iteration's 322.4-second wall time includes recovery and is not a clean speed measurement.

## Other accepted improvements

- [Q2 source and connection reuse](reuse.md): create a new independent document from retained owned input, review prefilled fields, and explicitly opt into encrypted connection reuse. No key is sent to the browser. Ownership, expiry, source integrity and parent-deletion independence were checked with synthetic data.
- [Q3 Catalan](catalan.md): neutral overview, SRS, arc42 and client handover support Catalan. Academic FKN profiles reject it explicitly. Migration 00011 updates persisted language validation; rollback refuses to relabel Catalan jobs.
- Q4 client handover: six sections cover scope, components, setup, operation, acceptance and maintenance. Instructions distinguish a CLI/library from a deployed web service and forbid invented credentials, production URLs, SLAs or support commitments.
- [Q5 stable reading](reading.md): native section navigation, explicit refresh, retained reading position and dirty/in-flight guards. Saved prose is distinguished from rendered PDF and reviewed content.

The product's practical advantage is a repeatable source-to-document workflow: one retained project can produce multiple structured deliverables, keep editable intermediate prose, expose recorded usage and rebuild PDF separately. This work does not establish superiority to a competent general-purpose agent; that requires a separate controlled comparison.

## Checks and limits

The isolated preview uses port 8084 and a dedicated database. The shared service on port 8080 was not restarted or modified. Private provider material and raw generated artifacts stay outside tracked files. No automated test files/suites were added; milestone 8 remains deferred.

[Phone/laptop browser log](checks/browser.log) covers 1366×900 and 390×844 layouts, language/profile persistence, explicit refresh position and unchecked connection reuse. [Citation validation smoke](checks/normalization.log) covers empty-only, supported, unknown and ordinary bracket/link text plus table markers. Additional integration evidence is recorded below once terminal results are available.

Deferred: arbitrary languages, cross-device account recovery, document version history, optional external search/SERP, and a direct agent benchmark. External search is not required to explain the supplied implementation and would introduce separate provenance and privacy decisions. Unknown approvals, real author information, measurements and acceptance results remain explicit gaps; removing those by invention would lower quality.

## Additional real-provider findings

A later raw overview produced 4/4 model sections (1,666 words, 189.3 seconds) but source review rejected several exact path/function claims and a visible TODO omitted from gap metadata. The first PZ produced 14/17 model leaves (2,722 words, 632.5 seconds including 157.3 seconds queued); it contained unsupported hardware/economic claims and three baseline leaves. These are **not accepted quality outcomes**. [The audit](audit.md) records the source checks and corrective behavior.

Gemini 2.5 Flash then exhausted its actual free daily request quota. Google returned a one-element JSON error array; the old parser lost the message and treated it as a transient rate limit. The parser now recognizes it and records `provider_quota`; [the manual check](checks/provider-error.log) verifies this classification. Raw account limits are observations of this account/run, not published universal limits.

OpenRouter's live catalog exposed unrelated nested pricing metadata that broke the old all-string catalog decoder. Pricing is now decoded only as scalars for the chosen model, retaining strict zero-price checks and rejection of unknown charges. The chosen free Gemma model subsequently passed discovery but its shared upstream pool returned HTTP 429. No paid route was enabled. Gemini 2.5 Flash-Lite appeared in the catalog but returned HTTP 404 for new-user access, recommending 3.5 Flash-Lite. This establishes why catalog existence alone is not a successful generation check.

The final continuation uses Gemini 3.5 Flash-Lite; it is a different model and cannot support a controlled causal speed/quality comparison with the 2.5 baseline. Official references checked: [thinking configurations](https://ai.google.dev/gemini-api/docs/thinking), [OpenAI-compatible parameters](https://ai.google.dev/gemini-api/docs/openai), and [pricing](https://ai.google.dev/gemini-api/docs/pricing). Actual account responses remain authoritative for availability.

The 3.5 continuation also returned JavaScript-style unquoted field names despite JSON-object mode. The existing chapter JSON Schema is now enabled for the explicitly supported Gemini 2.5 Flash, 3.5 Flash-Lite and 3.8 Flash models, following the provider's [OpenAI compatibility documentation](https://ai.google.dev/gemini-api/docs/openai). Unknown models keep their existing mode. Validation and bounded repair remain necessary; schema conformance does not establish factual correctness.

## Final measured runs and delivery

Pinned ZIP SHA-256: `6a01ec9e8fe4e65354b9c1197fe6f949f3f4958bbf3a0117157c592aeb46fae7`. The source was not executed. Raw metrics are in [live-results.json](checks/live-results.json); words mean whitespace-separated body words, excluding front matter, supplied specification and structured tables. Calls include discovery and unsuccessful requests. Tokens are recorded provider usage, not complete billing or failed-request consumption.

| Raw run | Model | Model sections / total | Baseline leaves | Body words | Recorded input / output tokens | Elapsed seconds |
|---|---|---:|---:|---:|---:|---:|
| Controlled baseline, stopped after one initial coverage pass | Gemini 2.5 Flash | 0 / 4 | 4 | 59 | 23,062 / 394 | 241.7 |
| Same-model intermediate, rejected factual quality | Gemini 2.5 Flash | 4 / 4 | 0 | 1,666 | 28,060 / 4,756 | 189.3 |
| Technical overview RU | Gemini 3.5 Flash-Lite | 4 / 4 | 0 | 1,185 | 41,572 / 5,317 | 382.7 |
| Explanatory note RU, reused source/connection | Gemini 3.5 Flash-Lite | 17 / 17 | 0 | 2,732 | 36,922 / 8,698 | 507.1 |
| Client handover CA, reused source/connection | Gemini 3.5 Flash-Lite | 6 / 6 | 0 | 1,372 | 20,309 / 3,916 | 185.2 |

The final three raw runs have substantive model text in every leaf, but still required editorial review. They initially recorded 2, 7 and 1 missing-information entries respectively. Source review corrected network access, command order, ignored flags, build/runtime dependencies, unsupported comparisons and wording; it also fixed a Catalan typo. Genuine unknown authorizations, tested environment versions and measurements are stated as limitations or proposed evaluation methods, never invented. Editorial changes also remove unnecessary completion instructions.

All three delivery documents were saved through the ordinary section-save endpoint and rebuilt by the worker, with explicit `user`/`edited` provenance and no further generation calls. Structured tables were checked and preserved. The raw files remain separate. Delivery body sizes are 1,594 / 3,183 / 1,568 words; all 4 / 17 / 6 sections carry editorial provenance, no `TODO:` body instruction or inline fact marker remains, and no baseline remains. This is **reviewed delivery**, not a claim of unattended factual correctness. Formal university approval and real runtime testing are not established.

Local delivery folder: `/Users/alevasyukovv/Downloads/ai-doc-gen-quality-2026-11-06/`, with three PDFs and matching editable Typst ZIPs. PDFs have 16 / 26 / 17 pages including front matter and the original supplied specification; these are not body-page counts. Representative body/table pages were rendered and visually inspected, including academic framing and Catalan diacritics. The original-language supplied specification remains unchanged as source material.

## Final verification

- Go build: [build.log](checks/build.log). Frontend typecheck/build: [frontend.log](checks/frontend.log). No automated test files or suites; milestone 8 remains deferred.
- Reuse ownership, tamper/purge isolation and explicit connection checks: [reuse.md](reuse.md); actual new-profile generations above verify the saved connection path.
- Browser laptop/phone, refresh/dirty guards and native no-JS navigation: [reading.md](reading.md).
- Compact manifests, declaration tails and visible TODO accounting: [retrieval.log](checks/retrieval.log). Latest source selection coverage: [pack-coverage.log](checks/pack-coverage.log).
- Provider errors and transport schema: [provider-error.log](checks/provider-error.log), [schema.log](checks/schema.log). An instrumentation bug initially recorded generation schema mode as false even when catalog negotiation enabled it; the final code preserves the transmitted mode in diagnostics. Historical reports are not rewritten.
- Independent reviews covered code security/correctness and source claims. A reviewer build raced removal of a disposable helper; the primary reran the full build after all helpers were removed. This was not hidden as a product build pass.
- No application code, tests or commands from the uploaded project were executed. No paid route was enabled. Temporary plaintext connection export was deleted; credentials/cookies/raw source artifacts are excluded from Git. Shared port 8080 services were not restarted.

The principal remaining limitation is factual reliability without source review. A complete JSON response, valid citations and full section coverage still do not guarantee correct technical claims. Automated semantic review, editable structured-table controls and a controlled comparison with general-purpose agents remain future work rather than claimed capabilities.
