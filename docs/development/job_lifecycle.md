# Job lifecycle and user feedback

The web request creates a durable job and returns its ID immediately. A Go worker does repository analysis and generation outside the HTTP request. The interface reads status and shows the latest useful result.

## State model

Keep **status** separate from **stage**. Status controls what the user can do; stage describes the work currently running.

| Status | Meaning | User-facing behavior |
|---|---|---|
| `queued` | Waiting for a worker | Show queue position only if it can be measured |
| `running` | Worker owns a lease | Show current stage and completed sections |
| `input_ready` | Evidence review is confirmed | State that document drafting is not available yet |
| `needs_input` | File selection, specification or required metadata is missing | Show a specific action and preserve completed work |
| `retry_wait` | Transient dependency failure with a scheduled retry | Show the reason and next retry time |
| `ready_for_review` | Draft sections are available | Let the user inspect and edit them |
| `completed` | Reviewed source and PDF are ready | Show downloads |
| `failed` | No automatic recovery remains | Show a short, actionable error and saved partial results |
| `cancelled` | User cancelled | Stop pending work and retain only what policy allows |

Stages include `validating`, `fetching`, `extracting`, `saving_inventory`, `selecting_files`, `source_selected`, `analyzing`, `generating`, and `rendering`. Store `stage`, `updated_at`, `last_heartbeat_at`, `completed_sections`, `total_sections`, and an append-only progress event for meaningful changes. Do not invent a percentage when the work size is unknown.

Implement legal transitions in `../../internal/statemachine` and enforce them in the storage transaction. A worker lease prevents simultaneous processing; expired leases are recoverable. Repeated work must be idempotent by job, stage and artifact version. Milestone 3 requeues old `input_ready` jobs once during migration and moves new jobs to `needs_input` for file selection. Confirmed selections queue `analyzing`, resetting the attempt budget. Extraction pauses at `needs_input` / `reviewing_evidence`. Evidence confirmation reaches `input_ready` / `evidence_confirmed`, the boundary before milestone 5. The worker dispatches from persisted inventory during retries, preserving intake and selection.

## Prompt user attention early

- Reject invalid URLs, archive type/size, unsupported profile, and missing required fields in the upload request.
- Detect missing specification or file-selection decisions before model calls. Switch to `needs_input` and show exactly what the user must provide.
- Set bounded timeouts for repository fetch, model calls and Typst compilation. After a stage timeout, report `retry_wait` or `failed` promptly; do not leave the UI in `running` indefinitely.
- For provider rate limits, honor `Retry-After` when reasonable, store `next_retry_at`, and show the pause immediately. For exhausted quota or invalid credentials, fail without automatic retries and show a safe message.
- Save each validated section as it completes. If a later section fails, show completed drafts and the failed section. Never claim the whole document is ready without a successful render.
- Poll the status endpoint at a short interval in the temporary UI. An event stream can replace polling later if needed; it is not required for the first working flow.

## Error contract

Storage keeps an internal error code and short diagnostic. The API returns a stable code, safe English message, and optional `retry_at` or `required_action`. Logs include job ID and stage, but never uploaded source or secrets. Separate temporary failures from permanent input or quota failures. A user can retry a failed stage after fixing its cause; retries must not overwrite prior manual edits.
