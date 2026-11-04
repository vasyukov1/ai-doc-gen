# AI Doc Gen

AI Doc Gen is a planned web application that turns a GitHub repository or a source archive into editable Typst and a compiled PDF. Users choose a document profile and may add a short project description.

The [product and niche analysis](docs/analysis/README.md) records evidence, assumptions, and product decisions. The [development guide](docs/development/README.md) contains the incremental implementation plan. The service accepts a public GitHub URL or ZIP, inspects files in a separate worker, and lets users confirm the file selection. Evidence extraction and document generation are not enabled yet.

Run `make run` to start PostgreSQL, migrations, the website, and the worker. Open `http://localhost:8080`. Run `make down` to stop the stack without deleting its data. `make build` builds the binary, `make test` runs the Go test suite when tests are added, and `make mocks` reports that mock generation is not configured yet. Use `make db-generate` after changing migrations to regenerate the checked-in Jet schema. The [development guide](docs/development/README.md) has local commands and architecture details.


### Evidence review

After confirming files, the worker extracts deterministic project evidence. Open the evidence link on the job page, review source citations and user context, add optional corrections, and confirm. Corrections retain the original facts and are attributed to the user. Save document metadata and a specification before confirming evidence. Confirmation queues grounded drafting. The read-only record is available at `/api/jobs/{id}/evidence`.

### Section generation (milestone 5)

Before starting the stack, generate a server encryption key in your shell:

```sh
export LLM_ENCRYPTION_KEY="$(openssl rand -base64 32)"
make run
```

Store this value securely and reuse it when restarting. Both web and worker need the same value. Losing it makes saved credentials unreadable; failed jobs can accept a replacement OpenRouter key after configuration is restored. Do not commit the value. Without it, evidence confirmation fails with an actionable setup message.

After file selection and evidence review, choose a model and enter your own OpenRouter API key in the website. Confirming starts background section generation. Drafts appear incrementally at the job's **Review draft sections** link. You can edit completed drafts and regenerate a single section after processing stops; other edits remain intact. Failed generation can retry remaining sections with a replacement key. Keys are encrypted in PostgreSQL and retained for regeneration until their job is deleted. This MVP is intended for local use. Reviewed sections can now be compiled with the document page.


### Document profiles and rendering (milestone 6)

Choose one academic variant: technical specification (TZ), explanatory note (PZ), operator’s manual (RO), individual test plan/methodology (PMI), or team test plan/methodology (PMI team). The technical project overview has a separate simpler title page and section plan. Both Russian and English are supported. Old `academic` jobs require resubmission with a specific variant.

After source intake pauses for review, open **Document input and PDF / Typst rendering**. Save the project title and title-page metadata. Academic profiles require one technical specification: paste text, upload a UTF-8 `.txt`/`.md` file up to 64,000 bytes, or complete all four structured fields. PDF, DOCX and Typst specification uploads are unsupported. State unknown requirements explicitly. Empty title fields appear as localized missing-information placeholders. Input freezes once drafting starts.

Review and edit the generated leaf sections, then confirm the review on the document page to queue compilation. The worker uses exactly Typst 0.15.1 with bundled Liberation Serif fonts. Docker installs it automatically; local worker runs require that compiler version on PATH. Compilation has a 45-second deadline, a private project root and a clean environment. User text is literal JSON data, never executable Typst.

On completion, the initiating browser can download the PDF and editable Typst ZIP using a seven-day capability cookie. The ZIP contains the trusted template, text/metadata JSON, profile assets, reviewed evidence, fonts/license and a reproduction command. Preserve the cookie until download; account ownership and full artifact lifecycle/retention are deferred to milestone 7. These templates preserve academic structure but do not certify standards compliance or verify runtime behavior.
