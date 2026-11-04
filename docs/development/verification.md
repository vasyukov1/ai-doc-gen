# Manual verification record

## Milestone 3 — source intake and file selection

Checked locally on 2026-11-03. No automated test files were added.

- `go build ./...`, `make run`, `docker compose config --quiet`, and `git diff --check` passed. The Jet generator produced models for `job_sources` and `job_files` after goose migration `00002_source_inventory.sql`.
- A public GitHub repository resolved to a commit SHA and produced a visible file inventory. A ZIP upload reached `needs_input` at `selecting_files`; the HTML page and JSON endpoint showed selected files and exclusion reasons.
- Academic selection left test files available but unselected. Technical overview selected them by default. A 600-file source selected at most 500 files by default.
- Saving a valid selection moved the job to `input_ready` at `source_selected`. An excluded file was rejected, and a second confirmation returned a conflict.
- Unsafe ZIP paths, symlinks, conflicting paths, more than 2,000 entries, and more than 100 MiB declared expanded size failed with clear error codes. Known credential files, dependencies, binaries, and nested archives were excluded. A source with only excluded files kept its inventory visible and asked for another project.
- A queued ZIP job completed after a worker restart. The inventory API did not expose the private workspace path.

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
