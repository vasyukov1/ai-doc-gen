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

## 2026-11-04 — Version deterministic evidence and preserve source facts

Milestone 4 extracts evidence only from confirmed selected files. Extractor version `1` sorts paths and manifest keys, verifies each stored file digest, and hashes the canonical unreviewed record. Go and package.json manifests yield structured declarations; README, CI and API schema text remains explicitly labeled excerpts. Source facts cite file paths and lines where accurately known; JSON manifest facts cite the complete file. Extraction never executes project code or infers runtime behavior. Per-file bounds and unknown runtime behavior remain visible in review.

Store one versioned JSON record per job with immutable extracted facts and separate user corrections. Optional submission descriptions and later corrections carry user attribution. Confirming review atomically persists corrections and moves the job to `input_ready` / `evidence_confirmed`; drafting remains unavailable until milestone 5. Worker recovery uses the persisted source inventory to dispatch analysis, avoiding repeated intake or selection loss. An evidence insert never overwrites an existing record.

## 2026-11-04 — Generate grounded sections with per-job OpenRouter credentials

Milestone 5 uses a replaceable provider interface and the OpenRouter chat completions protocol. The evidence confirmation form selects DeepSeek Chat V3, GPT-4.1 mini, or Gemini 2.5 Flash and accepts the user's own OpenRouter key. The worker sends that key only to the fixed OpenRouter HTTPS endpoint; redirects are rejected. No model is inferred from a missing selection. A future local provider can implement the same interface. Protocol reference: https://openrouter.ai/docs/api/api-reference/chat/send-chat-completion-request . Initial model identifiers were checked against https://openrouter.ai/api/v1/models .

A server-generated 32-byte base64 `LLM_ENCRYPTION_KEY` encrypts job credentials with AES-GCM and the job ID as associated data. Keep the same master key for web and worker and across restarts. Keys remain encrypted until job deletion so section regeneration can reuse them. Keys never appear in API responses, logs, prompts or cache keys. This is a local, unauthenticated MVP; access control and retention cleanup remain milestone 7. Key-spending POSTs reject cross-origin browser requests. A failed-generation retry accepts a replacement key, keeps the model and preserves saved sections.

Profiles and prompts are versioned independently. Each section receives a bounded relevant evidence subset with separately attributed user corrections. Cache hashes include the complete reviewed evidence, source digest, selected paths, profile/version, language, prompt/version, section, model and account-scoped provider identity. A cache hit is validated again. Source citations are structurally checked, not a proof of semantic entailment. Uncited paragraphs are explicitly flagged for review. User edits are marked and never overwritten when other sections regenerate.

Each call has a 45-second timeout and a bounded response. Persistent usage records include successful and failed attempts. Rate limits and transient output/network errors retry at most three job claims, retaining completed sections. Retry-After is bounded to one second through five minutes. Quota and authentication errors stop promptly. Per-section writes hold the job lock and validate its active owner and lease. Generation resumes from confirmed evidence after restarts.

`LLM_FIXTURE_MODE=1` enables a development-only adapter override using `LLM_ENDPOINT`, `LLM_MODEL`, and a separate server `LLM_API_KEY`. That override never receives the submitted key. Leave fixture mode unset for real OpenRouter use.
