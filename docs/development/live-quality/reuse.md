# Q2 — Reuse a retained project

Date: 2026-11-06. Scope: authorized live-quality Q2 within milestone 9. No automated test files or provider requests.

## Requirement/action/check index (D1–D4)

| ID | Requirement and user action | Implementation | Manual check |
|---|---|---|---|
| Q2.1 | An owner opens another document from an existing retained job and reviews source, profile, language, metadata and specification | `/jobs/{id}/reuse` | GET has saved input; another owner and expired job receive 404 |
| Q2.2 | The owner submits a new profile without reuploading or retyping | Retained-file copy, new queued job and inventory | New independent files and fresh document state; validation errors retain choices |
| Q2.3 | Connection reuse is an explicit unchecked choice; credentials remain private | Server-side decrypt/validate/re-encrypt under new job ID | New run has connection only on opt-in; response contains no key |
| Q2.4 | Deleting the parent does not break the child | New job-prefixed workspace and upload snapshot | Delete/purge parent, inspect child files and queue input |
| Q2.5 | Foreign, expired, deleting and cross-origin mutation is rejected | Owner middleware, transaction owner/expiry recheck, same-origin policy | 404/403 and no created child |

## Decisions (D6, D10, D13, D14)

- Reuse operates on a retained job, within the existing anonymous browser session. It is not a persistent account-wide project or credential vault.
- Copy only retained eligible text files; binaries, dependencies, likely secrets and excluded files are not recovered. Reclassify selection for the selected document profile. Preserve original source digest and pinned commit in the copied inventory. Fresh evidence/sections are produced through the existing worker queue.
- Retained text is copied into a new job-prefixed workspace, with verified file lengths/digests. A new private ZIP snapshot provides an independent queue input, including for GitHub sources. No fetch or uploaded-code execution occurs in the request.
- Hold the parent process guard during copying to prevent source purging. Atomically recheck owner/expiry/deletion while creating the child. Refuse a busy source with a recoverable form error.
- Offer saved free-draft connections only; legacy adapters are not silently converted. Never pass a `domain.Connection` or encrypted connection field to a template or response. The new run receives a fresh budget and no disabled-provider state.
- Shared file ownership: reuse owns new reuse files, `internal/web/handler.go`, and the small `internal/storage/jobs.go` transaction extraction. Primary integrates links in existing pages. Other contributors' work and `tmp_example/` remain untouched.

## Implementation steps

- [x] Copy and validate retained source independently.
- [x] Add atomic repository creation with fresh draft state.
- [x] Add owner-protected review form and routes, retaining invalid submitted values.
- [x] Build and run isolated manual HTTP/storage/security smoke checks.

## Verification / resume

On 2026-11-06, `go build ./...` completed with exit 0 and no output. A later shared-workspace build briefly encountered another task's temporary `ManualInspect` helper; after its removal, the full build again completed with exit 0. `git diff --check` for Q2 files completed with exit 0 and no output. No automated suite was run.

A disposable Go manual helper (`go run ./.manual-reuse`, removed after use) used the actual HTTP handler, repository, file intake, cleaner, evidence extractor and queue against the isolated migrated `reuse_quality_manual` database. It used synthetic source and a nonfunctional dummy credential; no worker or provider request ran. Two runs exited 0. Observed outcomes:

- GET 200 prefilled metadata/specification with HTML escaping, showed model/rate summaries without a saved key, and left reuse unchecked.
- Foreign GET/POST and expired GET/POST returned 404; cross-origin POST returned 403. A direct storage creation attempt with a foreign owner returned `sql.ErrNoRows`.
- Invalid HTML submission returned 400 while preserving changed specification, author and profile choices.
- Owner POST created a new queued job without an upload; default connection list was empty. Explicit reuse produced the same synthetic connection encrypted to a different ciphertext under the child ID, with fresh attempts and unreviewed document input.
- The child had separate upload/workspace paths, retained pinned commit provenance, selected build configuration for SRS, and continued excluding the synthetic `.env`.
- Changed parent source failed digest validation before child creation and removed temporary copies.
- Parent deletion hid it immediately; actual cleaner purge removed its workspace and upload. The child's original text and valid ZIP survived, normal evidence extraction returned facts, and the normal queue claimed the child.

Not verified here: real-provider generation through reused credentials (deliberately no real keys/calls), an interrupted database commit acknowledgement, and phone/laptop visual walkthroughs. Primary integration owns links and the combined UI walkthrough. The feature retains the existing private browser-session ownership and seven-day default retention; no account recovery or long-term vault is added.

Next exact command: `go build ./...` after integration; then open an owned retained job at `/jobs/{id}/reuse` for the combined phone/laptop walkthrough. No service was restarted and no commit was created by the Q2 agent.

## Primary integration

The integrated laptop/phone walkthrough passed prefilled input, unchecked connection reuse, no password fields and no horizontal overflow. A valid 63,999-byte UTF-8 specification plus legal metadata exceeded 200,000 URL-encoded bytes; the 1 MiB transport bound accepted it ([log](checks/max-form.log)). Per-field validation remains authoritative. A failed connection decryption now leaves source-only reuse available with a warning; decryption-failure injection was not repeated in the integrated browser.

Actual PZ and Catalan handover generation are submitted through the reuse route using an explicitly checked saved connection. Their terminal quality/results belong to the main report, not the synthetic security checks above.
