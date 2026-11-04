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
