# Incremental implementation plan

One numbered item equals one Codex implementation request. A future agent should receive the text of one item, read `../../AGENTS.md` and this directory, implement all its subitems, and report build plus manual smoke results. Do not combine adjacent items. The initial site is local; there is no deployment request.

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

## 6. Academic and technical Typst profiles

6.1. Port the existing five academic document variants as versioned assets in this repository. Keep their section order and title-page structure. Accept the technical specification from upload or a structured web form.

6.2. Provide complete Russian and English text for fixed sections, headings, title pages, prompts, and placeholders. Translate the English title-page text while preserving its structure.

6.3. Add one generic technical project overview profile. Keep its requirements distinct from the academic set.

6.4. Render user-reviewed sections into Typst without allowing uploaded text to execute arbitrary Typst functions. Compile in an isolated worker with pinned Typst and bundled fonts/assets.

**Acceptance:** every selected document compiles in both languages; mandatory headings and title-page fields are present; missing source material remains visible; exported Typst can reproduce the PDF.

## 7. Review, artifacts, and lifecycle

7.1. Add a section editor and final review page. Preserve user changes and show missing-evidence warnings.

7.2. Make the Typst source bundle and PDF downloadable only by the authorized job owner or a temporary capability token. Add retention and cleanup of uploaded code and generated artifacts.

7.3. Add cancellation, retry for transient failures, and clear job progress. Verify that duplicate requests and worker restarts do not create inconsistent artifacts. Show errors and intermediate sections as soon as the corresponding progress event is stored.

**Acceptance:** a user can submit a project, review and edit sections, download reproducible artifacts, and see a clear failure or cancellation state.

## 8. Final verification and release preparation

8.1. Add table-driven automated tests now: transition rules, archive validation, evidence extraction, template completeness, provider failure handling, and caching.

8.2. Add integration tests with PostgreSQL and a fake LLM provider, plus end-to-end checks from public repo/ZIP to Typst and PDF in both languages.

8.3. Review upload security, privacy and retention behavior, observability, deployment configuration, and documentation. Benchmark latency and cost on a consented evaluation set.

**Acceptance:** the test suite and build pass; critical failure paths have coverage; a reviewer can reproduce the deployment and the end-to-end document flow.
