# Manual verification record

## Milestone 3 — source intake and file selection

Checked locally on 2026-11-03. No automated test files were added.

- `go build ./...`, `make run`, `docker compose config --quiet`, and `git diff --check` passed. The Jet generator produced models for `job_sources` and `job_files` after goose migration `00002_source_inventory.sql`.
- A public GitHub repository resolved to a commit SHA and produced a visible file inventory. A ZIP upload reached `needs_input` at `selecting_files`; the HTML page and JSON endpoint showed selected files and exclusion reasons.
- Academic selection left test files available but unselected. Technical overview selected them by default. A 600-file source selected at most 500 files by default.
- Saving a valid selection moved the job to `input_ready` at `source_selected`. An excluded file was rejected, and a second confirmation returned a conflict.
- Unsafe ZIP paths, symlinks, conflicting paths, more than 2,000 entries, and more than 100 MiB declared expanded size failed with clear error codes. Known credential files, dependencies, binaries, and nested archives were excluded. A source with only excluded files kept its inventory visible and asked for another project.
- A queued ZIP job completed after a worker restart. The inventory API did not expose the private workspace path.

Later on 2026-11-05 the intake limits were raised to 100 MiB compressed, 500 MiB declared expanded size and 20,000 entries after oversized projects were rejected too often. `go build ./...` passed after the change. The rejection paths are unchanged but were not re-run with the new thresholds; manual re-checks are pending.

Section generation and document rendering remain in later milestones.

## Milestone 2.5 — Jet and local commands

Checked locally on 2026-11-03. No automated test files were added.

- `make build`, `go build ./...`, `docker compose config --quiet`, and `git diff --check` passed.
- `make db-generate` applied migrations and regenerated Jet table and model code from PostgreSQL.
- `make run` started PostgreSQL, a successful one-shot migration, the web server, and the worker. `GET /readyz` returned `ready`.
- A GitHub URL submitted through HTTP reached `input_ready`; `GET /api/jobs/{id}` returned the last progress message.
- A ZIP job remained `queued` while the worker was stopped, then reached `input_ready` after restart.
- A manually expired worker lease showed `recovery_pending`; the restarted worker reclaimed the job and reached `input_ready`.
- `make mocks` reported its placeholder state. `make -n test` showed the intended `go test ./...` command without running tests before milestone 8.
- `make down` removed all four Compose containers. Manual smoke jobs and their ZIP upload were removed.

The default PostgreSQL host port changed from 5433 to 5434 because another local project already used 5433. Queue index tuning is deferred until real workload measurements exist.

## Milestone 2 — durable input jobs

Checked locally on 2026-11-02. No automated test files were added, as requested.

- `go build ./...` and `docker compose config --quiet` passed.
- The service applied goose migration `00001_jobs.sql`; a second start reported no pending migrations.
- `GET /healthz`, `GET /readyz`, `/`, and the stylesheet returned success.
- A public GitHub URL returned `202 Accepted` with a job ID. The worker moved the job to `input_ready` and the API reported the latest message.
- An invalid URL and a non-ZIP upload returned `400` before creating jobs.
- A ZIP job stayed `queued` while the worker was stopped, then moved to `input_ready` after the worker restarted.
- A ZIP job whose saved file was removed before the worker started moved to `failed` with `source_missing` and an actionable message.
- A job with an expired worker lease reported `recovery_pending` and a recovery message; a new worker reclaimed it and moved it to `input_ready`.
- Invalid browser form input returned a readable `400` response.

Still outside this milestone: fetching GitHub contents, inspecting ZIP entries, generating sections, rendering Typst, and testing provider rate limits. The status model and retry behavior for those stages are specified in [job_lifecycle.md](job_lifecycle.md) and will be implemented with their respective milestones.

## Milestone 4 — 2026-11-04

- `go build ./...` and `git diff --check` passed. No automated test files were added, per milestone sequence.
- Started the existing Colima instance, built the Docker stack, and applied goose migration 4. Generated Jet table/model types for `job_evidence` from PostgreSQL.
- Submitted a small ZIP with README, Go module and code, package.json, OpenAPI and CI files. Confirmed six files; the worker paused at `needs_input` / `reviewing_evidence`. JSON contained structured module/dependency/script facts, source paths and line spans, explicit excerpts, and separately attributed user description.
- Rendered the review HTML, saved a correction plus additional context, and confirmed. Original source facts remained unchanged; saved review reached `input_ready` / `evidence_confirmed`. Duplicate review returned HTTP 409.
- A second identical ZIP with the same selection and description produced the identical complete evidence record and SHA-256 digest. Confirmed evidence and user corrections survived web/worker restart.
- Removed only the two smoke jobs and their uploaded/extracted files. The local stack remains running.
- Limits: structured manifest parsing covers Go and package.json. Other manifests, README, API schemas and CI use bounded literal excerpts. Code extraction recognizes common Go/Python/JavaScript declaration prefixes; it does not verify runtime behavior. JSON facts cite files without fabricated line numbers. Extraction limits are visible in review. Browser visual layout, public GitHub extraction in this milestone, and forced analysis lease recovery were not manually exercised; intake and recovery use existing milestone 3 mechanics.

## Milestone 5 — 2026-11-04

- `go build ./...`, `git diff --check`, and `docker compose config --quiet` passed. No automated test files were added. Goose migrations 5 and 6 applied; Jet regenerated all 13 table definitions. The configured Go proxy was unavailable, so regeneration used `GOPROXY=https://proxy.golang.org`.
- A temporary fake HTTP provider outside the repository produced four persisted sections in English and Russian from an uploaded small Go/README ZIP. Evidence confirmation accepted a selected model and encrypted per-job key. The sections HTML returned HTTP 200 and source citation links point to evidence entries.
- Repeating an identical Russian request produced four sections with zero provider call records, confirming cache reuse. Editing the overview then regenerating interfaces preserved the overview edit.
- A fake 429 showed `retry_wait`, `provider_rate_limit` and a retry timestamp, then resumed. An exhausted-quota response failed immediately. An invented citation was rejected as `provider_invalid_output` and recovered on retry. A zero-section failed job accepted a replacement key and completed on retry. A cross-origin regenerate POST returned HTTP 403.
- These checks exercised the provider interface using explicit fixture mode, not a real OpenRouter request. Model quality, real provider accounting, real timeout behavior and browser visual layout remain unverified. Semantic claim support still requires human review; schema/reference checks alone cannot establish truth. The source-linked review UI flags uncited paragraphs and missing information.
- Keys remain encrypted until their job is deleted. Automatic retention cleanup and ownership controls remain milestone 7. No Typst or PDF work was started.

- Docker image rebuild passed. Removed nine smoke jobs and their local uploaded/extracted files, stopped temporary fake services, and restored the normal Compose web/worker stack without fixture mode.

## Milestone 6 — 2026-11-04

- `go build ./...`, `git diff --check`, `docker compose config --quiet`, and the final Docker image build passed. No automated tests or test files were added. Goose migration 7 applied and Jet generated the new document-input table/model.
- Manually compiled all six profiles in both languages using explicitly synthetic reviewed text and missing-evidence warnings. Final PDF page counts (RU/EN): TZ 18/17, PZ 13/13, RO 10/10, PMI 13/13, PMI team 13/13, technical overview 8/8. Rendered every PDF page with Poppler into contact sheets and inspected the title structures, mandatory headings, Cyrillic glyphs, contents, warning boxes and change-register pages. Reduced displayed contents depth to two while retaining all body subsections. Inspected the final revised Russian specification and overview, plus a 100-DPI live-application approval sheet with missing title fields. No clipping or overlap was observed in the checked fixtures.
- Unpacked each of the 12 source bundles into a separate directory and ran its documented Typst 0.15.1 command. SHA-256 of every rebuilt PDF matched the original byte for byte. Fonts/license and profile assets are included; worker exports also include reviewed evidence JSON for citation lookup.
- A malicious-text manual fixture containing `#read`, `#include`, `#eval`, markup, HTML-like text and Cyrillic compiled as literal visible text. An external private marker did not appear in PDF text. No uploaded Typst or source code was executed.
- Built/restarted the normal Docker stack with pinned compiler hashes and bundled fonts. Actual ZIP intake reached file selection. A synthetic technical overview was manually advanced to section review without model calls; review POST queued the worker and reached `completed` / `rendered`. Capability-cookie PDF download returned 200; no cookie returned 403.
- With the final image, submitted a separate Russian TZ ZIP job. UTF-8 specification upload returned 303; replacing it with four structured fields returned 303; conflicting upload-plus-text sources returned 400. Seeded 35 explicitly synthetic reviewed leaf drafts to isolate rendering from provider availability. Review POST queued and completed background compilation; PDF and Typst downloads returned 200 with its capability cookie, and PDF download without the cookie returned 403. The downloaded PDF reports Typst 0.15.1, Russian metadata and no JavaScript. Missing approval/title fields were visibly localized.
- Review found and fixed render-intent invalidation on failed-render section edits/regeneration, and rejected generation-retry requests while render review is active. The JSONB input upsert was corrected after the first HTTP smoke exposed a typed SQL expression mismatch. Saved input now inserts and replaces successfully through the running web form.
- Live OpenRouter validation was performed separately with tiny synthetic evidence and returned `provider_quota` before any text was generated. Real model-generated document quality remains unverified. The rendering fixtures above are deliberately synthetic and do not claim an end-to-end successful real-provider generation. Accounts, retention/cleanup, cancellation and broader artifact lifecycle remain milestone 7; automated verification remains milestone 8.
- Manually forced the synthetic rendered TZ job to `failed` / `render_failed`: generation retry returned 409 while render intent was active; saving one section returned 303 and atomically cleared both reviewed state and capability hash. A fresh review returned 303 and the worker completed the corrected render. These checks did not use provider credentials or make model calls.

## Milestone 7 — 2026-11-04

- `go build ./...`, `git diff --check`, and Docker image build passed. Jet types were regenerated from the migrated local schema. A fresh disposable PostgreSQL database migrated through all versions to 8, including owner/expiry/tombstone/resume columns; that disposable database was then removed. Existing user records were preserved.
- Actual Docker HTTP smoke used independent owner-cookie clients. The foreign client received 404 for every job HTML/API view, file/evidence/section read, mutation, generation retry and PDF/Typst download. The owner could see its dashboard/job. Cross-origin submission returned 403. Cookies contain a random bearer secret, while the database stores only a digest.
- Submitted a synthetic ZIP through the real website endpoint and observed `needs_input/selecting_files`; saved document metadata. Seeded four clearly synthetic technical-overview drafts in that newly owned job to isolate rendering from unavailable provider quota. Actual Typst 0.15.1 produced private PDF and reproducible source ZIP downloads. No provider call or real API key was used.
- Edited a completed section through HTTP: the job returned to final review; edited text, citations, missing-evidence warnings and other sections were preserved; old PDF download returned 409. Explicit review/rebuild returned completed artifacts containing the saved literal text. Duplicate render requests left exactly one current PDF and one source record.
- Cancel/retry probes restored file-selection and evidence reviews to `needs_input` and an ordinary final review with an empty stage to `ready_for_review`, using persisted pre-cancellation status. Queued cancellation survived worker restart. A temporary compiler sleep probe inside the worker container was terminated 1.04 seconds after the owner cancelled; it was restored immediately. Four saved edited drafts remained. Retry plus worker restart compiled those drafts and left exactly two artifact records.
- Deleted only known synthetic jobs through owner-authorized HTTP: access disappeared immediately, cancellation/tombstone persisted, and the cleanup command subsequently removed artifact directories, associated cache keys and cascading job data. For manual speed, only those known synthetic tombstones had `purge_after` advanced rather than waiting the normal 90 seconds. Another known synthetic job was expired; access returned 404, cleanup tombstoned it, then purged it in the next probe. All three pre-existing legacy ownerless rows were retained throughout.
- Synthetic old upload/workspace/render directories and an expired cache entry were removed. An active job directory and the outside target of an orphan symlink were preserved. Cleanup does not use arbitrary stored filesystem paths for deletion. Per-job advisory guards prevent purge while a worker still uses its directories; lease publication checks remain transactional.

Limits: these are manual smoke checks, not the automated suite reserved for milestone 8. Live LLM output quality, provider-side cancellation billing, production HTTPS/proxy configuration, browser automation and large-dataset cleanup performance remain unverified. Separate submission POSTs create separate jobs; duplicate render/lifecycle mutations are serialized and may return conflict. Lost owner cookies cannot be recovered. Legacy ownerless jobs are deliberately inaccessible and paused pending explicit operator migration. Metadata/specification edits preserve existing prose and require user consistency review before rebuilding. Cache records expire seven days after insertion; job purge also evicts its associated keys. Superseded/orphan files wait at least 24 hours and may need multiple bounded passes.

## Post-milestone fixes — 2026-11-04

- `go build ./...`, `go vet ./...`, and `gofmt` passed. No automated tests were added.
- The cross-origin guard was reproduced as the cause of a blocked file-selection step: an `Origin` of `http://localhost:8080` reached a server that saw `Host: 127.0.0.1:8080`, and the literal comparison returned 403. The guard now trusts `Sec-Fetch-Site: same-origin` and `none`, rejects `cross-site`, and otherwise compares `Origin` hosts case-insensitively.
- A disposable PostgreSQL database and local web/worker processes on port 8090 verified the behavior. A submission with a rewritten `Origin` and `Sec-Fetch-Site: same-origin` returned 303; `cross-site` and an unmatched `Origin` without fetch metadata returned 403; an unmatched `Origin` was accepted when its host matched. Rejections logged origin, host, and fetch-site.
- A GitHub submission over the fixed guard reached `needs_input` / `reviewing_evidence` without any manual file-selection POST. The inventory classified `.gitignore`, `.gitattributes`, and `go.sum` as excluded `tooling_metadata`; `Dockerfile` and `docker-compose.yml` as eligible but unselected `build_config`; and kept `go.mod`, README, and source files selected. Evidence contained module, dependency, README, API-schema, and code-declaration facts; the file page rendered read-only.
- The disposable database, build output, and local smoke processes were created for this check. The agent sandbox denied process signals, so their termination is left to the operator; the cleanup command is recorded in the session report.

## One-form automatic flow — 2026-11-04

- `go build ./...`, `go vet ./...`, and `gofmt -l` passed. No automated tests were added.
- A disposable PostgreSQL database and local web/worker processes on port 8090, with a fake OpenAI-compatible provider outside the repository, exercised the complete new path. One `POST /jobs` carried the source, profile, language, title-page metadata, specification, model and per-job key; the job row, document input and encrypted credentials were written in one transaction.
- No manual file selection or evidence confirmation POST was sent. The worker confirmed the default selection, stored evidence as accepted, drafted four sections and queued rendering automatically. The job reached `completed` / `rendered` with 4/4 sections, and the PDF download returned `application/pdf` (512 bytes) from the real Typst 0.15.1 compiler.
- The evidence page rendered read-only and the evidence API record still exposed attributed facts.
- Invalid submissions (missing key, missing model, conflicting specification sources, malformed upload) returned `400` before any job row was created.
- When `LLM_ENCRYPTION_KEY` was unset, the web process generated `DATA_DIR/state/llm-encryption.key` with `0600` permissions and reused it; the worker decrypted the per-job key with the same file. The explicit environment override still takes precedence.
- Limits: the provider was a local fake, so real OpenRouter output quality, accounting and latency remain unverified. The shared `state` volume between web and worker in Docker was configured but not exercised under Docker in this check. Mutating an already completed job still follows the explicit review/rebuild path.

## Milestone 9 — free baseline drafts — 2026-11-04

Environment: Apple M1, 8 GB unified memory, native Go worker and Typst 0.15.1; isolated PostgreSQL database on the local Docker server. Source fixture: a small task-board ZIP with README, Go manifest/declarations and an API schema. Uploaded source was never executed. No automated test files or final test suite were added.

Measured manual end-to-end runs (HTTP submission, persisted sections, actual Typst PDF and source ZIP):

| Scenario | Sections / model-improved | Model calls | Observed time |
| --- | --- | --- | --- |
| No keys, technical overview EN | 4 / 0 | 0 | 0.9 s |
| No keys, technical overview RU, ordinary Docker stack | 4 / 0 | 0 | 0.9 s |
| No keys, TZ RU | 35 / 0 | 0 | 0.8–1.3 s |
| No keys, TZ EN | 35 / 0 | 0 | 1.3 s |
| Successful fixture provider, overview RU | 4 / 4 | 1 | 1.3 s |
| Successful fixture provider, TZ RU / EN | 35 / 35 | 9 | 1.7 s each |
| All three fixture providers return exhausted daily quota | 4 / 0 | 3 | 1.3 s |
| First fixture provider exhausted, second succeeds | 4 / 4 | 2 | 1.3 s |
| Malformed fixture JSON, two providers | 4 / 0 | 2 | 1.3 s |
| One response member attempts to set server provenance | 4 / 3 | 2 | 1.3 s |
| Repeated temporary 429 with Retry-After 2 seconds | 4 / 0 | 12 | 23.5 s |
| Provider hangs on every request | 4 / 0 | 6 | 240.8 s |

**Fixture times are not real LLM latency or prose quality measurements.** The hanging-provider run exercised real request cancellation: five approximately 45-second calls and one deadline-shortened call, followed by a successful baseline PDF. Database timestamps measured 240.32 seconds from acceptance to completion; polling/download observation measured 240.77 seconds.

Additional checks:

- Identical inputs/account reused all four model sections with zero new calls. Replacing one fixture cache record with a JSON scalar caused exactly one new call, repaired the cache and completed in 1.3 seconds. Fixture endpoint identity is isolated from real-provider caches.
- A provider with a three-second cooldown remained eligible when a later connection could not fit the prompt into its TPM budget: the job waited, then completed with one call in 3.4 seconds.
- Invalid response members were rejected independently; valid members remained saved. Metadata such as `origin` cannot be supplied by the model. The failed member was tried only once more. Invalid-output fallback used Gemini then Groq in the recorded fixture call history.
- A mixed-document edit/regeneration check saved a manual architecture edit, then regenerated only overview. Final origins were model/user/baseline/baseline; the manual text and both unrelated baseline records were unchanged.
- Cancelling an in-flight hanging call produced `cancelled`; the request ended after approximately 965 ms and saved drafts remained. Restarting a worker during a different hanging call retained the original deadline and persisted attempt, then resumed to four model sections with a second call. The restart did not renew the deadline.
- Provider catalog smoke accepted zero-price JSON-capable metadata, rejected paid pricing, rejected missing JSON capability and rejected a missing model. The public OpenRouter catalog was reachable and listed 11 zero-price models advertising `response_format` at inspection time; no authenticated generation was performed.
- Migration 9 and Jet generation succeeded against a fresh disposable schema. The ordinary local Docker stack was rebuilt and migrated. Its no-key Russian overview reached `completed/rendered`; foreign artifact access returned 404; section/document pages returned 200; evidence included the attributed project-title reference. Only that synthetic production-stack job was deleted through its owner route; existing user jobs were preserved.
- Poppler rendered selected annotation/content pages of Russian TZ baseline/model fixtures. Cyrillic text, provenance labels and source references were readable; no clipping was observed on inspected pages. Baseline TZ produced 19 pages. A duplicated provenance paragraph was removed and the final annotation was checked again. English outputs compiled and retained their headings.
- Smoke checks exposed and corrected a JSONB cache-upsert assignment, new scheduler rows looking spuriously busy due to database/host clock differences, and empty warning slices serializing as null. Reruns produced saved model sections and PDFs. Fresh code review found regeneration scope and fallback-selection problems; their fixes were exercised by the edit/regeneration and fallback checks above.

Validation: `go build ./...`, `go vet ./...`, `git diff --check`, Docker build and migration passed. Automated tests remain deferred to milestone 8.

Remaining limits: no real free account key was provided for a live model run, so model prose quality, actual free quotas, provider-side billing settings and actual token accounting remain unverified. A five-minute full model-written document is not claimed. Large downloads, worker backlog, source failures and database/compiler outages can exceed or prevent the small-project target. Gemini/Groq free account status relies on the explicit form declaration; their billing cannot be inferred from the key. Provider limits are shared per configured key fingerprint, not across unrelated keys on the same external account. A single job submits sequentially; multi-worker admission is coordinated by PostgreSQL, but a broad concurrency/load test suite remains for the final testing milestone.

## 2026-11-04 — Milestone 9 quality recovery checks

Implementation and acceptance are distinct. No automated test suite or local LLM was added. Fresh `go build ./...` and `git diff --check` passed. Migration 10 was applied to an isolated smoke database, Jet provider-call types regenerated, then the additive migration was applied to the local Docker site. The earlier user job and artifacts were preserved. The local site serves the new chapter/outcome controls.

Manual evidence on a copy of the failed three-service source snapshot:

- Extraction now retains 13 API operations, 11 schemas, 104 function bodies, 93 static call-reference records, 6 route declarations and 6 SQL table declarations. It also retains source locations, paragraphs and module records. Generated Swagger representations no longer repeat canonical JSON. These are static declarations, not proof of executed workflows.
- Controlled PZ RU response: all 17 leaves saved in 3 chapter calls, PDF in 2.11 seconds. This fixture exercises transport and storage, not document quality.
- Controlled invalid-citation responses: six calls (three chapters plus one targeted repair each), all 17 leaves retained, 1.68 seconds. The diagnostic endpoint records the exact invalid source ID and section; ordinary bracket text and omitted model headings are accepted. All fixture outputs remain `partial` because their explicit quality gaps are retained.
- Controlled truncated responses: each initial chapter ended with `finish_reason=length`; complete earlier members were recovered. Three repairs saved only missing members. All 17 leaves available, six calls, 1.66 seconds. Earlier runs accidentally used a non-truncating fixture and are not evidence for this check.
- Repeated identical PZ/account: zero API attempts, all 17 sections restored from validated cache, 1.67 seconds. A direct isolated database smoke also confirmed persisted catalog reuse and schema capability retention.
- No-provider TZ EN: 35 outline leaves, zero calls, PDF in 1.75 seconds. Explicit continuation after a manual edit preserved the user-origin body exactly, 1.25 seconds.
- Low TPM=4,000 admitted real fixture requests rather than failing preflight because of the shared brief. A deliberately/interleaved truncated response consumed an extra repair; the run stopped at the persisted model deadline and rendered 3/4 model sections plus one outline section in 240.52 seconds. This verifies the bounded partial path, not full low-TPM completion.
- Direct HTTP-adapter check: a 429 `rate_limit_exceeded` with a generic billing sentence remains temporary `provider_rate_limit`, honors a two-second Retry-After, and records status 429.
- Native RU and EN rendering with table data compiled successfully. Visual inspection confirmed Cyrillic, table headings/cells, source notes and literal `#read(...)` content; uploaded/model text was not executed. Separate PZ page inspection checked hierarchy and font readability. This checks layout, not the fixture's intentionally repetitive prose.
- Deployed Docker no-provider overview: `outline`, four leaves, zero model calls, PDF in 0.84 seconds. The synthetic production job was deleted through its owner route; the original user job remains completed and unchanged.

Two important independent-review findings were fixed: fit the entire serialized prompt including shared brief/repair material; salvage complete members and split unfinished truncated groups. Manual low-TPM and truncation checks above cover those failure modes.

Still unverified: live-provider semantic quality in both languages, five-minute useful-document completion with real quotas, model-specific context limits, and equal-input comparison against the chapter control. Cancellation/restart ownership behavior retains prior milestone coverage; the complete fault matrix was not rerun here. No fixture time is advertised as a real model benchmark. A full resolved cross-service workflow graph is not claimed.

## 2026-11-04 — FKN layout and red TODO acceptance

`go build ./...` and `git diff --check` passed after the final edits. The five academic profiles (TZ, PZ, PMI, team PMI, operator guide) and the generic overview compiled in both RU and EN: twelve manual PDFs. Academic sample lengths were TZ 20, PZ 14, PMI 14, team PMI 14 and operator guide 11 physical pages. The title-sheet count excludes the separate approval sheet, matching the supplied layout's page-counter reset.

Visual contact-sheet review covered RU/EN approval sheets, title sheets, body/outline pages and complete ten-column change registers. Cyrillic, vertical inventory cells, footer inscriptions and red TODOs in metadata/body/table cells were legible without clipping. PDF character inspection found red text in every sample. Literal `#read(...)` remained visible text. Recompiling the exported TZ RU and operator-guide EN roots with the pinned compiler produced identical PDF SHA-256 values.

The independent review found one material issue: HTML text inputs would strip line breaks from the default institution banner. Organization, supervisor/approver positions and author/team fields now use textareas. A deployed owner-session smoke verified that the exact multiline institution default survived a save round trip; checking `resolve_missing` cleared only the edited section's warnings and preserved warnings in all other sections. The synthetic no-provider TZ compiled in 1.26 seconds with all 35 sections and no API calls; this measures layout/infrastructure, not model prose. That synthetic production job was deleted through its owner route.

The updated Docker site is running. Prior user jobs and PDFs were not automatically rewritten. Rebuilding existing reviewed content or creating a new job uses the FKN layout. Font rendering intentionally uses the bundled Liberation Serif rather than an unlicensed/system-dependent Times New Roman copy; full pixel identity and formal institutional approval are not claimed.

## 2026-11-05 — SRS and arc42 profiles

Manual checks only; no automated suite or real model key was used. `go build ./...` and `git diff --check` passed. A separate PostgreSQL smoke database ran the existing migrations unchanged; no profile migration is necessary.

- Native rendering compiled SRS and arc42 in RU and EN. SRS sample PDFs contain 22 editable leaves across seven major headings (12 physical pages); arc42 samples contain all twelve numbered sections (17 pages). Visual review of title, annotation, body and source-context pages confirmed neutral layout, localized headings, red TODOs, literal table text and no FKN institution defaults or inventory frames. Annotation carries the SRS scope qualification or arc42 attribution.
- Controlled HTTP-provider runs through the website completed SRS RU with 22 saved sections in seven calls (0.87 seconds) and arc42 EN with 12 sections in six calls (1.26 seconds). The fixture checked that every request contained nonempty profile guidance and guidance for every requested section, and recorded the seven/six logical groups. These are infrastructure times, not real-provider speed or prose-quality evidence.
- No-key runs completed SRS EN (22 outline sections, 0.84 seconds) and arc42 RU (12 outline sections, 1.25 seconds). Separate submissions without any input specification also completed as honest outlines (SRS 1.28 seconds; arc42 1.24 seconds), retaining missing-context TODOs.
- Recompiling exported SRS RU and arc42 EN roots with pinned Typst/fonts/timestamp reproduced identical PDF bytes. Both source bundles contain the profile assets and ARC42-NOTICE.md.
- A scoped independent review found no material defects in the profile, grouping, cache, selection or rendering changes. Profile/section instructions participate in cache identities; empty optional fields preserve existing-profile plan serialization.

Live-provider SRS requirement correctness, cross-chapter requirement-level traceability and arc42 architecture accuracy require human review. Neither ISO conformity nor a live five-minute content benchmark is established by the controlled responses.

Repeated identical SRS and arc42 fixture submissions restored all model sections from cache with zero API calls (1.27 and 1.25 seconds respectively), confirming both new profiles use the existing validated section cache.

The local Docker site was rebuilt and restarted after confirming no active production jobs. Both options appeared in the submission form; deployed no-key SRS RU and arc42 EN jobs rendered successfully (1.26 and 0.83 seconds). Only these synthetic production jobs were deleted through their owner session. Temporary native workers/provider and the isolated smoke database were stopped/removed; prior user jobs were preserved.

## 2026-11-05 — Provider rejection classification and outline warnings

Manual checks only; no automated suite was added. `go build ./...`, `gofmt -l internal/documents` and the template parse/execute check passed.

Read-only diagnosis of the local database for job `135a6e7f-96d5-49e6-9239-dece23f4c99c` (technical-overview, EN) found: two admissions (catalog plus one generation call), `provider_rejected` disabling Groq, and four baseline sections. The recorded Groq response was HTTP 400 `json_validate_failed` with the message "Failed to generate JSON". The same database shows job `f7af9915-75e1-4211-8e8a-2acda4bb7169` completed with no providers connected, and job `0f234fbb-8fc7-4f4a-b7b2-6b00de3a749b` wasting eight of twelve admissions on 429 responses that reported an output limit of 1000 tokens per minute.

- A local HTTP stub returned the recorded Groq 400 body: the adapter now reports `provider_invalid_output`, temporary, with the rejected generation retained for repair. The same stub returned a 429 OTPM rejection: the adapter reports `provider_rate_limit` with the parsed 1000-token output limit. Context-window and unknown 400 bodies map to `provider_context_budget` and `provider_rejected` respectively.
- The next-attempt fitting for learned input/output limits and the repair prompt reuse were not exercised against a live provider. Live Groq retry behavior, repair success rate and the actual effect of the 1000-token output window remain unverified.
- Outline and partial completed jobs now render a visible warning on the job and sections pages; the template check executed all seven pages successfully.
- No existing job, section or artifact was modified. Jobs created before this change keep their stored diagnostics and reasons.
