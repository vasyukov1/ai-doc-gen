# Architecture decisions

## 2026-11-02 — Start with explicit Go wiring

The first server uses `net/http`, `../../internal/config`, and a small handler package. Explicit constructors are sufficient for the current dependency graph. Reconsider `uber/fx` only if wiring becomes difficult to understand after real components exist. This avoids adding lifecycle and DI concepts before they solve a concrete problem.

## 2026-11-02 — Persist jobs in PostgreSQL, without Kafka

The intended flow is `POST /jobs` → durable queued row → Go worker claim → status polling → artifact download. PostgreSQL is already needed for job metadata and can serve the initial queue. A separate broker is useful only if measured load or multiple event consumers require it. Never hold a long generation request open in the web handler.

## 2026-11-02 — Use explicit states, not a state-machine package

Define named statuses and legal transitions in Go, then enforce transitions inside database transactions. Keep status separate from stage as described in [job_lifecycle.md](job_lifecycle.md). The worker reports section progress and actionable errors separately. Reconsider a library if branching, rollback, and parallel paths make the transition table hard to audit.

## 2026-11-02 — Report progress and failure early

The user must receive a job ID immediately. Store stage changes and completed sections as work proceeds. Missing input becomes `needs_input`; temporary provider limits become `retry_wait` with a visible retry time; exhausted quotas and permanent errors become `failed`. Set bounded timeouts and keep completed drafts available. Polling is enough for the temporary HTML interface.

## 2026-11-02 — Use requested package boundaries

The entry point is `../../cmd/ai-doc-gen/main.go`. `internal/documents` owns business logic, `internal/statemachine` owns transition rules, `internal/pgqueue` owns job leases, and `internal/storage` owns persistence. Kafka and protobuf remain reserved until a real interface requires them. Goose migrations live in `migrations/`.

## 2026-11-02 — Defer automated tests to the final milestone

At earlier steps, compile and manually verify the delivered behavior. Add the automated unit, integration, and end-to-end suite after the pipeline and UI are stable, as requested by the product owner. Record manual checks and known gaps in each milestone summary.

## 2026-11-03 — Use Jet for database access

Use Jet v2.16.0 to generate typed table definitions from the goose-migrated PostgreSQL schema. Keep generated files in `internal/storage/generated` and use Jet in `internal/storage` and `internal/pgqueue`. Migrations and their advisory lock remain SQL. Jet improves schema-aware query maintenance; it does not by itself improve query latency. Run migrations once before starting the web and worker processes. Keep the initial queue index until workload measurements justify a change.

## 2026-11-03 — Pin and inspect source archives

Use the GitHub REST API to resolve the public repository's default branch to a commit SHA, then download the ZIP for that SHA. Follow redirects only to GitHub API or codeload hosts. Inspect ZIP paths and file metadata before extraction, enforce size and entry limits, and unpack only eligible text files into a private worker directory. Persist the inventory and pause at `needs_input` for user selection. A server-side GitHub token is optional; API limits and permanent intake errors are visible to the user.
