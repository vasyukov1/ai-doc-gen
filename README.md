# AI Doc Gen

AI Doc Gen is a local web application that turns a GitHub repository or a source archive into a compiled PDF with editable Typst source. The user submits the project once, chooses a document profile and language, fills in the title-page metadata, and receives a generated document.

The [product and niche analysis](docs/analysis/README.md) records evidence, assumptions, and product decisions. The [development guide](docs/development/README.md) contains the incremental implementation plan. The service accepts a public GitHub URL or ZIP, inspects files in a separate worker, and selects the documentation source files automatically; the list stays readable with per-file exclusion reasons. Evidence extraction, drafting and rendering then run without further user input.

Run `make run` to start PostgreSQL, migrations, the website, and the worker. Open `http://localhost:8080`. Run `make down` to stop the stack without deleting its data. `make build` builds the binary, `make test` runs the Go test suite when tests are added, and `make mocks` reports that mock generation is not configured yet. Use `make db-generate` after changing migrations to regenerate the checked-in Jet schema. The [development guide](docs/development/README.md) has local commands and architecture details.


### Evidence (internal grounding)

After the automatic file selection, the worker extracts deterministic project evidence: manifests, README excerpts, API schemas, CI files and code declarations with source paths. Evidence grounds the model so sections cite real facts instead of inventing them. The extraction is accepted automatically and the job continues; the **Extracted evidence** page is read-only, and the same record is available at `/api/jobs/{id}/evidence`. The optional project description is treated as user context, not as verified code evidence.

### Section generation (milestone 5)

The submission form asks for the model and your own OpenRouter API key, so no server-wide provider configuration is required. The server encrypts that key with a master key: when `LLM_ENCRYPTION_KEY` is unset, the web process generates one under `DATA_DIR/state/llm-encryption.key` (Docker shares the `state` volume with the worker). Set `LLM_ENCRYPTION_KEY` explicitly to a base64-encoded 32-byte value if you prefer to manage it yourself; both processes must use the same value, and losing it makes saved credentials unreadable. Do not commit it.

Generation starts automatically after evidence extraction. Drafts appear incrementally at the job's **Review draft sections** link. You can edit completed drafts and regenerate a single section after processing stops; other edits remain intact. Failed generation can retry remaining sections with a replacement key. Keys are encrypted in PostgreSQL and retained for regeneration until their job is deleted. This MVP is intended for local use.


### Document profiles and rendering (milestone 6)

Choose one academic variant: technical specification (TZ), explanatory note (PZ), operator’s manual (RO), individual test plan/methodology (PMI), or team test plan/methodology (PMI team). The technical project overview has a separate simpler title page and section plan. Both Russian and English are supported. Old `academic` jobs require resubmission with a specific variant.

The first form also collects the document title, title-page metadata and one technical specification: paste text, upload a UTF-8 `.txt`/`.md` file up to 64,000 bytes, or complete all four structured fields (purpose, functions, constraints, acceptance). PDF, DOCX and Typst specification uploads are unsupported. State unknown requirements explicitly. Empty title fields appear as localized missing-information placeholders. Input freezes once drafting starts.

When drafting finishes, the worker compiles automatically and the job page exposes the PDF and the editable Typst ZIP. The worker uses exactly Typst 0.15.1 with bundled Liberation Serif fonts. Docker installs it automatically; local worker runs require that compiler version on PATH. Compilation has a 45-second deadline, a private project root and a clean environment. User text is literal JSON data, never executable Typst.

To change the result, open the document or section page: completed documents and metadata stay editable, but saving a change invalidates the review, so confirm final review on the document page to rebuild. Downloads are unavailable until that rebuild succeeds.

On completion, the owner browser can download the PDF and editable Typst ZIP. The ZIP contains the trusted template, text/metadata JSON, profile assets, reviewed evidence, fonts/license and a reproduction command. Keep the owner cookie to retain access until the job expires. These templates preserve academic structure but do not certify standards compliance or verify runtime behavior.

## Private jobs and retention

Use [My private jobs](http://localhost:8080/jobs) in the browser that submitted the project. The private HttpOnly owner cookie authorizes every job view, edit and download. Clearing or losing it permanently loses access; a job ID does not restore ownership. HTTPS adds the Secure flag. Legacy jobs created before ownership was introduced remain preserved but inaccessible and paused pending an explicit operator migration.

Review sections and missing-evidence warnings on the section page, then confirm final review on the document page to compile. Completed documents can be edited; saving an edit requires explicit review and rebuilding before downloads resume. Changes to the specification preserve existing prose, which should be checked again. Cancellation stops pending work promptly and retains saved progress. Retry resumes a review checkpoint or the saved processing stage; update invalid provider credentials through the section page. Delete hides a job immediately and cancels work before disposal.

`RETENTION_DAYS` defaults to `7` (1–365). New jobs show their immutable `expires_at`; expiry covers source, drafts, credentials, evidence and downloads. The worker runs bounded cleanup every minute. Deletion/expiry uses a 90-second worker shutdown grace period. Associated cache keys are evicted on purge; cache records also expire seven days after insertion. Unreferenced uploads, workspaces and superseded render directories are scavenged after 24 hours. Ownerless legacy records and referenced files are preserved. Run one pass manually with the same database and storage configuration: `go run ./cmd/ai-doc-gen cleanup`.

`DATA_DIR` defaults to `./artifacts` (`/data` in Docker). `UPLOAD_DIR` and `SOURCE_DIR` must be inside that application-owned root; do not point it at unrelated files. Cleanup rejects symlink parents and removes only controlled job paths or old unreferenced entries under these directories. The worker needs writable storage; the web process only reads compiled source artifacts.

### Free baseline drafts

New jobs generate a source-linked baseline even without an API key. Optionally connect Gemini, Groq and/or OpenRouter in the initial form to improve the text. Enter the exact model ID from the provider console and its account RPM/TPM limits. Use a Gemini project without billing, a Groq Free Plan, or a verified zero-price OpenRouter model. The application never activates billing or switches to a paid model; it cannot inspect Gemini/Groq account billing settings.

The model stage ends 240 seconds after job acceptance and admits at most 12 API calls, including model discovery. Remaining baseline text is compiled into PDF with an explicit quality label. The five-minute target is for small projects; external availability, long downloads and queued work are not guaranteed. Drafts remain editable. Regenerating one section starts a new bounded pass only for that section.

Run `make run` to build and migrate the local stack (migration 9). Existing jobs retain their original generation path. See [the current design](docs/development/local_model/plan.md) and [verification results](docs/development/verification.md). Live model quality and speed must be measured using a free account; fixture runs are not model benchmarks.

## Deployment

See [Railway operations](docs/development/deployment.md) for the supervised web/worker image, PostgreSQL, persistent storage and HTTPS configuration. Google sign-in remains deferred; browser ownership still applies.
