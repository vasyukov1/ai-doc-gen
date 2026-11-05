# Job lifecycle and user feedback

The web request creates a durable job and returns its ID immediately. A Go worker does repository analysis and generation outside the HTTP request. The interface reads status and shows the latest useful result.

## State model

Keep **status** separate from **stage**. Status controls what the user can do; stage describes the work currently running.

| Status | Meaning | User-facing behavior |
|---|---|---|
| `queued` | Waiting for a worker | Show queue position only if it can be measured |
| `running` | Worker owns a lease | Show current stage and completed sections |
| `input_ready` | Evidence review is confirmed | Legacy boundary before generation is queued |
| `needs_input` | File selection, specification or required metadata is missing | Show a specific action and preserve completed work |
| `retry_wait` | Transient dependency failure with a scheduled retry | Show the reason and next retry time |
| `ready_for_review` | Draft sections are available | Let the user inspect and edit them |
| `completed` | Reviewed source and PDF are ready | Show downloads |
| `failed` | No automatic recovery remains | Show a short, actionable error and saved partial results |
| `cancelled` | User cancelled | Stop pending work and retain only what policy allows |

Stages include `validating`, `fetching`, `extracting`, `saving_inventory`, `selecting_files`, `source_selected`, `analyzing`, `generating`, and `rendering`. Store `stage`, `updated_at`, `last_heartbeat_at`, `completed_sections`, `total_sections`, and an append-only progress event for meaningful changes. Do not invent a percentage when the work size is unknown.

Implement legal transitions in `../../internal/statemachine` and enforce them in the storage transaction. A worker lease prevents simultaneous processing; expired leases are recoverable. Repeated work must be idempotent by job, stage and artifact version. Milestone 3 requeues old `input_ready` jobs once during migration. New jobs confirm their default file selection automatically and queue `analyzing`; `needs_input` / `selecting_files` remains the fallback when automatic confirmation fails. Confirmed selections reset the attempt budget. Extraction stores the evidence record as already accepted and queues `generating`; `needs_input` / `reviewing_evidence` remains a fallback when automatic acceptance cannot be completed. Migration 5 requeues legacy `input_ready` / `evidence_confirmed` jobs, which fail actionably if credentials are missing. When drafting completes, the worker queues `rendering` automatically, so a successful run ends `completed` / `rendered` without a manual step. The worker dispatches from persisted inventory during retries, preserving intake and selection.

## Prompt user attention early

- Reject invalid URLs, archive type/size, unsupported profile, and missing required fields in the upload request.
- Detect missing specification or file-selection decisions before model calls. Switch to `needs_input` and show exactly what the user must provide.
- Set bounded timeouts for repository fetch, model calls and Typst compilation. After a stage timeout, report `retry_wait` or `failed` promptly; do not leave the UI in `running` indefinitely.
- For provider rate limits, honor `Retry-After` when reasonable, store `next_retry_at`, and show the pause immediately. For exhausted quota or invalid credentials, fail without automatic retries and show a safe message.
- Save each validated section as it completes. If a later section fails, show completed drafts and the failed section. Never claim the whole document is ready without a successful render.
- Poll the status endpoint at a short interval in the temporary UI. An event stream can replace polling later if needed; it is not required for the first working flow.

## Error contract

Storage keeps an internal error code and short diagnostic. The API returns a stable code, safe English message, and optional `retry_at` or `required_action`. Logs include job ID and stage, but never uploaded source or secrets. Separate temporary failures from permanent input or quota failures. A user can retry a failed stage after fixing its cause; retries must not overwrite prior manual edits.

Document metadata, the specification, the model and the provider key are validated and stored in the same transaction as the job, so no worker claim can start without usable input or credentials. The document page remains the place to correct them later. Academic variants have explicit leaf-section plans; a legacy `academic` profile requires resubmission with a variant. `ready_for_review` now exposes the document page: confirming all saved sections and input queues `rendering` with a persisted reviewed flag. The worker dispatches reviewed input to the compiler before generation, then atomically persists both artifact records and transitions to `completed` / `rendered`. Artifact publication validates the active worker lease. Failed renders remain actionable and can be requeued from the document page. Editing or regenerating after a failed render invalidates its review/capability so generation is dispatched correctly. Milestone 7 extends this flow with completed-document editing, owner access and lifecycle controls below.

## Milestone 7 ownership, review and disposal

Every new job belongs to the anonymous `job_owner` browser session. All `/jobs/{id}` and `/api/jobs/{id}` endpoints require the owner, including read-only file/evidence/section views and downloads. Expired/deleting/foreign/legacy jobs return 404. `/jobs` lists the browser's retained jobs. Lost cookies have no ID-based recovery.

Completed drafts and title/specification input can be edited. A saved change moves `completed` to `ready_for_review`, invalidates review, preserves missing-information warnings and other sections, and requires explicit compilation. Metadata/specification changes preserve existing section prose; review it for consistency before rebuilding. Only a successful atomic artifact publication returns `completed`.

`POST /jobs/{id}/lifecycle` accepts `cancel`, `retry`, or `delete`. Cancellation atomically revokes the current lease; the worker cancels its context within the next one-second heartbeat. A per-job advisory process lock prevents disposal of an in-flight workspace even after the grace period. The previous status is saved: cancelled file/evidence reviews return to `needs_input`, and cancelled final reviews return to `ready_for_review`. Other retries queue work from durable input, inventory, evidence or reviewed rendering state and preserve completed drafts. Failed credentials may be replaced through the section page. Duplicate queued retries/renders return conflict; duplicate cancellation is harmless. A fresh claim token fences prior workers even within the same process.

New jobs expose `expires_at` and default to seven days from creation. Worker cleanup tombstones expired jobs, revokes leases, and waits 90 seconds before filesystem/database disposal. User deletion uses the same grace period and hides the job immediately. Credentials, private evidence, uploads, workspaces, artifacts and job-associated cache entries are included. Periodic cache expiry and 24-hour orphan/superseded-render scavenging are bounded. Restarted cleanup repeats an incomplete purge safely. Legacy ownerless jobs have no automatic expiry and are never executed by the new worker; preserve them until an explicit operator migration.

## Milestone 9 free draft runs

New submissions store an encrypted connection list (possibly empty) and a 240-second deadline atomically with the job. Existing jobs without that record keep their legacy behavior. All baseline sections are saved during `running/generating` before any model request. `completed_sections` counts available sections, while `model_sections` counts model-improved sections. Both are visible before compilation.

Provider unavailability, invalid output or exhausted model budgets are soft limitations: keep the baseline, persist the reason and continue to `ready_for_review`, automatic render queue and `completed/rendered`. Source, persistence and compiler errors remain failures. Waiting for provider admission is an observable `provider_wait` event within the bounded running stage. Cancellation still revokes the lease and stops external requests; automatic restarts never renew the draft deadline or request budget.

A section's explicit regeneration starts a fresh budget and targets only that section. Other baselines, model drafts and edits remain unchanged. Model metadata is server-owned. A final PDF includes baseline/mixed provenance and the reason improvement stopped.

### Recovery: content outcome and continuation

`completed/rendered` continues to mean successful artifact compilation. The separate `content.outcome` is `outline`, `partial`, or `draft_ready`; model-section count alone is not a quality score. Missing-information entries keep the result partial. Owners can read sanitized call diagnostics at `/api/jobs/{id}/generation-report`. Rejected prose stays private in database retention records. `/jobs/{id}/generation/continue` is a same-origin owner mutation available only after a run stops and only when baseline sections remain; it renews the durable model budget, invalidates rendered review, and queues unfinished work without deleting accepted prose or edits.

Academic rendering now uses `fkn-2`. Existing artifacts are never rewritten on startup or template update; a new job or explicit rebuild uses the new embedded layout. Section saves retain missing-information warnings unless the owner explicitly checks `resolve_missing=yes`. That resolution affects only the edited section, marks it user-origin, invalidates reviewed artifacts, and requires the normal explicit rebuild. Red TODO text remaining in the body is not silently removed.
