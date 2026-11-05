# Incremental implementation plan

One numbered item equals one Codex implementation request. A future agent should receive the text of one item, read `../../AGENTS.md` and this directory, implement all its subitems, and report build plus manual smoke results. Do not combine adjacent items. The initial site is local; there is no deployment request.

> Updated 2026-11-04: the website now collects title-page metadata, the specification, the model and the provider key in the initial form. File selection, evidence acceptance, drafting and Typst rendering run automatically; the evidence page is read-only. See [decisions.md](decisions.md).

## 1. Foundation and minimal website — implemented in this request

1.1. Create a Go module with `../../cmd/ai-doc-gen/main.go`, `internal/config`, a small `net/http` server, graceful shutdown, and a health endpoint.

1.2. Add a responsive landing page showing repository/ZIP, profile, language, and optional description controls as a **disabled preview**. The page must say submission is not yet available.

1.3. Add `../../docker-compose.yaml` for local PostgreSQL, reserve the requested package boundaries, and document them. Keep the web server independent of the database until jobs are implemented.

**Acceptance:** `go build ./...` succeeds; `/` serves the site and `/healthz` returns `ok`; unfinished controls cannot submit; PostgreSQL starts when Docker is available.

## 2. Durable asynchronous jobs and working submission — implemented in this request

2.1. Add versioned goose migrations for jobs, inputs, progress events, status history, and artifact metadata in `../../migrations`. Add a `DATABASE_URL` configuration with startup validation.

2.2. Implement `POST /jobs` for one public GitHub URL or ZIP upload, optional description, selected document ID, and `ru`/`en`. Validate size, type, and required fields before persisting. Store uploads outside the web root and record a digest.

2.3. Implement a separate `worker` mode in `../../cmd/ai-doc-gen`. Put job claims, leases, heartbeats, and bounded retries in `internal/pgqueue`. Claim jobs with `FOR UPDATE SKIP LOCKED`. Do not process generation in the HTTP request.

2.4. Implement domain types in `../../internal/domain`, legal status transitions in `internal/statemachine`, and persistence in `internal/storage`. Follow [job_lifecycle.md](job_lifecycle.md): status and stage are separate, and terminal or actionable states are visible promptly. Avoid a state-machine library for this initial flow.

2.5. Enable the form and add a status page that polls `GET /jobs/{id}`. Show stage, last update, needed input, retry time, and safe failure messages. Mark stale worker leases as recoverable rather than leaving the page indefinitely in `running`.

**Acceptance:** submitting a valid input returns a job ID quickly; the worker advances status to `input_ready`; a restart does not lose a queued job; invalid input never enters the queue; a forced worker failure becomes visible. `input_ready` means the input was stored and checked. No document is generated yet.

## 2.5. Typed database access and local commands — implemented

2.5.1. Generate Jet table types from the goose-migrated PostgreSQL schema. Pin the generator version, keep generated code in the repository, and document regeneration.

2.5.2. Move job and queue queries to Jet while preserving transactions, status history, progress events, leases, and `FOR UPDATE SKIP LOCKED`. Keep goose SQL migrations and PostgreSQL-specific migration locking as SQL.

2.5.3. Keep migration execution separate from web and worker startup. Make queue state and lease updates atomic within their existing transactions.

2.5.4. Add `Makefile` targets `run`, `build`, `test`, `down`, and `mocks`. `run` starts the complete local stack; `down` stops it. `mocks` is a documented placeholder until interfaces need generated mocks.

**Acceptance:** `make build` succeeds; `make run` starts PostgreSQL, web and worker; `make down` stops all three; URL and ZIP jobs still reach `input_ready`; an expired lease can be reclaimed. Do not add automated test files before item 8.

## 3. Safe source intake and profile-aware file selection — implemented

3.1. Fetch a public GitHub repository at a pinned commit or unpack a ZIP in an isolated work directory. Apply compressed/expanded size, file-count, path, symlink, and nested-archive limits. Never run source code.

3.2. Build a file inventory with reason codes for include/exclude. Ignore binaries, dependencies, generated output, and likely secrets. Keep test files available for test-related document sections.

3.3. Show the selected file list and allow the user to adjust it before generation. Use `needs_input` while waiting; persist the final selection with the job.

> Updated 2026-11-04: the worker now confirms the default selection automatically; the list remains readable and manual selection stays available as a fallback. See [decisions.md](decisions.md).

**Acceptance:** a small public repository and ZIP produce a visible inventory; unsafe archives fail safely; included files and exclusions are explainable.

## 4. Structured project evidence — implemented

4.1. Extract deterministic facts from repository structure, manifests, README, API schemas, CI files, and selected code. Attach source path and line span when possible.

4.2. Store a versioned evidence record with provenance (`source`, `user`, or `unknown`). Incorporate the optional description as user-supplied context, without treating it as code evidence.

4.3. Add a compact evidence review screen so the user can see and correct obvious gaps before drafting. Persist the review decision and resume the job without losing earlier work.

**Acceptance:** the same pinned source and selection produce a reproducible evidence record; every displayed technical fact has a source or user attribution.

## 5. Section generation through an LLM API — implemented; live-provider acceptance pending

5.1. Add a replaceable model provider interface with server-side encrypted per-job OpenRouter credentials supplied through the website, per-call timeouts, usage recording, and bounded retries. Classify rate limits, exhausted quota, invalid credentials, malformed output, and network timeouts.

5.2. Plan required sections from a versioned profile. Retrieve only relevant evidence for each section; request structured output, preserve required headings, and mark missing information instead of inventing it.

5.3. Cache section results by source digest, selected files, profile/version, language, evidence version, prompt version, and model. Validate schema and citations before accepting a section.

5.4. Save validated sections as they finish and update `completed_sections` and `total_sections`. Show completed drafts immediately. On `429`, honor a reasonable `Retry-After` and expose `retry_wait`; on exhausted quota, fail promptly with an actionable message.

5.5. Allow one section to be regenerated without replacing user edits elsewhere.

**Acceptance:** a real small repository yields reviewable, source-linked sections in Russian or English; a repeated identical request can reuse cached sections; unsupported claims are flagged; rate limits and partial results are visible without waiting for the whole job.

## 6. Academic and technical Typst profiles — implemented; live model quality remains unverified

6.1. Port the existing five academic document variants as versioned assets in this repository. Keep their section order and title-page structure. Accept the technical specification from upload or a structured web form.

6.2. Provide complete Russian and English text for fixed sections, headings, title pages, prompts, and placeholders. Translate the English title-page text while preserving its structure.

6.3. Add one generic technical project overview profile. Keep its requirements distinct from the academic set.

6.4. Render user-reviewed sections into Typst without allowing uploaded text to execute arbitrary Typst functions. Compile in an isolated worker with pinned Typst and bundled fonts/assets.

**Acceptance:** every selected document compiles in both languages; mandatory headings and title-page fields are present; missing source material remains visible; exported Typst can reproduce the PDF.

Implemented with explicit variant/leaf plans and preserved academic heading hierarchy. Minimal capability-protected retrieval is included solely to inspect rendered milestone 6 artifacts. Milestone 7 below adds full browser ownership, retention, cancellation and lifecycle work. Manual verification and provider limitations are recorded in [verification.md](verification.md).

## 7. Review, artifacts, and lifecycle — implemented; provider quality remains unverified

7.1. Add a section editor and final review page. Preserve user changes and show missing-evidence warnings.

7.2. Make the Typst source bundle and PDF downloadable only by the authorized job owner or a temporary capability token. Add retention and cleanup of uploaded code and generated artifacts.

7.3. Add cancellation, retry for transient failures, and clear job progress. Verify that duplicate requests and worker restarts do not create inconsistent artifacts. Show errors and intermediate sections as soon as the corresponding progress event is stored.

**Acceptance:** a user can submit a project, review and edit sections, download reproducible artifacts, and see a clear failure or cancellation state.

Milestone 7 adds browser ownership, private downloads, completed-document editing with explicit rebuild, cancellation and durable retry, creation-time expiry, tombstoned deletion and bounded cleanup. Legacy ownerless records are preserved and excluded from execution. See [verification.md](verification.md) for manual checks and remaining limits. Automated tests remain deferred to milestone 8.

## 8. Final verification and release preparation

8.1. Add table-driven automated tests now: transition rules, archive validation, evidence extraction, template completeness, provider failure handling, and caching.

8.2. Add integration tests with PostgreSQL and a fake LLM provider, plus end-to-end checks from public repo/ZIP to Typst and PDF in both languages.

8.3. Review upload security, privacy and retention behavior, observability, deployment configuration, and documentation. Benchmark latency and cost on a consented evaluation set.

**Acceptance:** the test suite and build pass; critical failure paths have coverage; a reviewer can reproduce the deployment and the end-to-end document flow.

## 9. Free, bounded document drafts — implemented; live provider benchmark pending

9.1. Save source-linked baseline sections before any external API call. Preserve every profile heading and expose baseline/model/user provenance.

9.2. Improve compact groups of at most four sections using optional Gemini, Groq and OpenRouter connections. Keep independent section validation, caching, editing and regeneration. Limit output to 2,400 tokens per group.

9.3. Persist a 240-second model deadline and a maximum of 12 API attempts, including catalog checks. Coordinate concurrency, conservative RPM/TPM reservations and cooldowns through PostgreSQL. Fall back to saved baseline text when providers cannot help.

9.4. Compile the saved draft into PDF and reproducible Typst, marking baseline/mixed content and improvement limitations. Completed means artifacts were compiled, not that claims were verified.

**Acceptance:** a small project yields a useful baseline and PDF within 300 seconds, including when all model providers fail. Record request counts, time, token usage, failure behavior and visual PDF inspection. Live provider quality and speed require a real free account; fixtures do not establish them.

This request explicitly revises milestone 9 without starting milestone 8 or implementing local inference. See [local_model/plan.md](local_model/plan.md) for the current design and [verification.md](verification.md) for measured results.

### Milestone 9 corrective acceptance — generation quality recovery

The previous baseline-PDF target is superseded by [generation_quality_plan.md](generation_quality_plan.md). A successful compiler run is not successful documentation generation. Preserve the bounded deadline and free-only providers, but generate logical chapters from coherent sources, retain actionable diagnostics and distinguish outline/partial/content-ready outcomes. Implementation is available; live-model semantic quality and the five-minute useful-document target remain unverified. Milestone 8 automated tests and local inference remain excluded.

### Academic template correction — FKN layout and visible manual fields

Use the supplied FKN layout for all five academic variants. Keep the technical overview separate, preserve document IDs and user prose, and render manual-completion fields as red TODOs. Add explicit resolution of section-gap warnings in the editor. Acceptance: build, compile the five profiles in RU/EN, visually inspect front matter/body/change register, verify literal TODO/table rendering and reproducible bundles. No new numbered milestone or automated suite is started.

### Additional document profiles — SRS and arc42

Add independently selectable `srs` and `arc42` profiles with bilingual structure, profile/section-specific generation instructions, neutral PDF layout, private reproducible bundles and red manual TODOs. SRS uses an independently authored requirements structure informed by ISO/IEC/IEEE 29148:2018, not copied normative text or a conformity certificate. arc42 preserves its twelve sections and is not renamed ISO/IEC/IEEE 42010. Use seven and six semantic request groups respectively to leave room for discovery/repair within the existing persisted budget. Preserve the FKN profiles and existing jobs. Acceptance: build, manual bilingual render and website generation checks, grouped provider-fixture inspection, visual TODO/profile verification. No final automated testing milestone is started.

### Milestone 9 corrective acceptance — response formatting and form feedback

Normalize Markdown section bodies into trusted Typst structures without executing generated code. Keep internal evidence separate from external bibliography, exclude assistant metadata from source analysis, and show required-input/submit errors within the form. Preserve saved text and user edits. Validate Russian/English PDF rendering, unsafe markup, source filtering and browser form retention manually; do not start the deferred automated testing milestone.

### Milestone 9 corrective acceptance — harness wave 1

See [the wave index](harness-wave-1/plan.md) and [triage](harness-wave-1/triage.md). Five accepted changes cover correctable submission, preserved edits/tables, honest compilation/content status, visible connection failures, and profile/privacy guidance. Keep simulated reviews separate from measured checks; full localization and regeneration history remain deferred. This does not start milestone 8.
