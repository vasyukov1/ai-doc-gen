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

## 2026-11-04 — Render explicit academic variants with literal data

Milestone 6 introduces versioned `tz`, `pz`, `ro`, `pmi`, and `pmi-team` profiles and a distinct technical overview. Preserve original heading order and hierarchy, the approval sheet and separate title sheet, annotation, contents and change register. Every leaf heading has a stable independent section ID; structural parents do not borrow unrelated generic paragraphs. Legacy `academic` jobs are not silently mapped to a document variant: submit a new job with a specific variant. Profile and prompt versions are bumped to `sections-2` and `grounded-2`, with localized prompts and evidence-relevance hints.

Collect title metadata plus an academic technical specification before any model call. Accept one UTF-8 `.txt`/`.md` upload, literal text, or four complete structured fields (purpose, functions, constraints, acceptance). Specifications are bounded to 64,000 bytes and explicitly attributed to the user in model evidence. They are requirements, not verified implementation facts. Missing title fields remain localized visible placeholders; no standards references, economic indicators or test results are invented. Generic overview requires a title but no academic approvals or specification.

Store document input in `job_document_inputs` (goose migration 7 and generated Jet types). Freeze it when drafting begins. After the existing section review, an explicit review confirmation atomically queues rendering and generates a 256-bit capability; only its SHA-256 is persisted. A seven-day HttpOnly, SameSite Strict cookie gives the initiating browser artifact access. HTTPS additionally sets Secure. This minimal retrieval supports milestone 6 acceptance; accounts, retention, cancellation and complete lifecycle work remain milestone 7. Completed artifacts are immutable in this milestone. Failed-render edits/regeneration clear render approval and the capability atomically; generation retries reject active render intent.

Only the embedded trusted template is executable. Metadata, specifications and section text are JSON strings rendered by `text()`, never `eval` or string interpolation into Typst code. The worker compiles in a new private root with a clean environment, no source tree or credentials, isolated package locations, a 45-second deadline, exact Typst 0.15.1 and embedded Liberation Serif fonts under the SIL Open Font License. Docker verifies official compiler archive hashes for Linux amd64/arm64. No user-uploaded Typst is accepted. JSON, template, versioned profiles, fonts/license, reviewed evidence and README form the editable source bundle. A fixed creation timestamp and bundled fonts reproduce PDF bytes on the same pinned compiler/platform.

Academic leaf counts are 35 (TZ), 17 (PZ), 9 (RO), and 16 each (PMI variants); technical overview has four. The worker's overall bound is raised to 60 minutes to accommodate these sequential 45-second calls, while each provider call and render retains its own short timeout. Per-section saves and heartbeats preserve resumability.

## 2026-11-04 — Anonymous browser ownership and bounded retention

Milestone 7 uses a random 256-bit bearer cookie (`job_owner`) with HttpOnly, SameSite Strict, a one-year browser lifetime, and Secure on HTTPS. HTTP remains supported for the local Docker site. PostgreSQL stores only its SHA-256 digest. Central middleware checks ownership, expiry and deletion tombstones before every HTML/API job read, mutation and artifact download. Cross-site mutation requests are rejected, including submission. Downloads use owner authorization; no sharing capabilities are issued. Clearing cookies permanently loses access: there is no account or ID-based recovery. The jobs page shows up to 100 retained jobs belonging to this session.

Migration 8 preserves legacy ownerless records, files and credentials, makes them inaccessible by guessed job ID, and excludes them from worker claims and automatic expiry. They are not claimable through the website. An explicit operator-assisted ownership/data migration is future work; do not assign ownership merely because someone knows an ID.

New jobs have an immutable creation-time `expires_at`: `RETENTION_DAYS` defaults to 7 (allowed 1–365). Expiry and user deletion revoke leases and first persist a tombstone, then wait 90 seconds before purging. A per-job PostgreSQL advisory process lock additionally prevents filesystem disposal while a worker still uses that job; acquisition and heartbeat database calls are bounded to three seconds. Workers poll lease validity every second and cancel provider requests and Typst subprocesses on lease loss. Each claim receives a fresh lease-owner token so a prior attempt cannot publish into a later claim. Database writes already lock and check the live lease. Cancellation preserves saved sections, evidence and review decisions; retry resumes saved review checkpoints using the pre-cancellation status, otherwise the durable rendering/evidence/inventory checkpoint. Cancellation does not revoke a provider charge already accepted remotely.

A worker-owned cleaner runs immediately and every minute; `cleanup` runs one bounded pass. Each pass expires at most 20 jobs, purges at most 10, removes at most 100 expired cache entries and scavenges at most 50 old filesystem entries. Database cascades remove evidence, sections, calls and encrypted credentials. Purging evicts associated section-cache keys, even when shared with another job. Cache entries also expire after seven days from insertion, independent of provider account. Superseded renders, abandoned render directories, intake directories and uploads are scavenged after 24 hours only when unreferenced and not associated with an active, tombstoned or legacy job. Known legacy references are always protected. Crashed or ambiguous commits therefore leave recoverable data until safe cleanup rather than deleting a potentially committed artifact.

`DATA_DIR` is the application-owned root (default `./artifacts`, Docker `/data`); source and upload directories must be strictly inside it. Cleanup constructs paths from validated server-generated UUIDs, never persisted arbitrary paths, checks containment and rejects symlink parents before removal. Tombstones and purge hold job locks, making concurrent cleaners and publication coherent. Filesystem removal is repeatable; failed removal retains the tombstone for a later pass. Concurrent downloads may fail if their job expires or is deleted during transfer.

Completed sections and document metadata remain editable. Saving a change preserves other drafts and missing-evidence warnings, clears reviewed state, and moves completed jobs back to final review. Downloads are unavailable until explicit review and rebuilding succeeds. Typst publications replace both artifact metadata records atomically; superseded files remain private until scavenging. Repeated mutations may return conflict after the first succeeds; they do not create multiple live publications. Separate submission POSTs intentionally create separate jobs.

## 2026-11-04 — Confirm the default file selection automatically

Milestone 3 paused every job at `needs_input` / `selecting_files` for a manual confirmation even though the inventory already preselected every eligible file. The product owner found this step redundant. The worker now confirms the default selection immediately after saving the inventory and queues `analyzing`; the file page stays readable with per-file reason codes, and manual selection remains available as a fallback when automatic confirmation fails. Lock files and repository metadata (`go.sum` and other lock files, `.gitignore`, `.gitattributes`, `.dockerignore`, `.editorconfig`, `.gitmodules`) are excluded as `tooling_metadata`. Build and deployment files (`Dockerfile*`, `docker-compose*`, `compose.*`, `Makefile`) stay eligible but are not selected by default as `build_config`. Manifests such as `go.mod` and `package.json` remain selected because evidence extraction derives module, dependency and script facts from them.

## 2026-11-04 — Trust browser fetch metadata for same-origin checks

The previous guard compared the `Origin` header with the literal `http(s)://` + `Host` value. That rejected legitimate browser submissions whenever a proxy or tunnel rewrote `Host`, the browser reported a null origin, or the host case differed. `Sec-Fetch-Site` is set by the browser and cannot be altered by page scripts, so `same-origin` and `none` are accepted outright, `cross-site` is rejected, and clients without fetch metadata fall back to a case-insensitive `Origin` host comparison. Rejected mutations log origin, host and fetch-site so the exact cause stays visible. `SameSite=Strict` owner cookies remain the second line of defense.

## 2026-11-04 — One-form submission with automatic drafting and rendering

The product goal is to submit a project once and receive a finished document. The website form now collects the document variant, language, title-page metadata, one technical specification source, the model and the per-job OpenRouter key. `POST /jobs` validates the whole input before anything is stored. The job row, document input and encrypted credentials are created in one transaction, so the worker never starts a job whose input or key cannot be used.

File selection, evidence extraction and drafting no longer pause for confirmation. The worker confirms the default selection, stores evidence as already accepted, and queues generation. Evidence remains the grounding layer: the model receives only extracted facts with source paths, but the page is read-only and requires no action. When drafting finishes, the worker invokes `QueueRender` automatically, so a successful run ends with a compiled Typst/PDF download on the job page. Editing a completed section still requires the explicit rebuild path, so user edits are never overwritten silently.

Per-job keys remain until a local provider replaces OpenRouter. The master key is resolved at startup: when `LLM_ENCRYPTION_KEY` is unset, the process generates and reuses `DATA_DIR/state/llm-encryption.key`; a shared `state` volume keeps web and worker consistent. Jobs created before this change keep their stored records and need a fresh submission if they have no saved credentials.

## 2026-11-04 — Baseline-first free drafts (milestone 9)

The approved revision replaces the proposed local-model rollout with a bounded free-cloud improvement pass over a deterministic baseline. New jobs need no API key. Preserve all headings, cite source excerpts and user requirements, and mark missing evidence. Persist `baseline` / `model` / `user` provenance; completed means PDF rendering succeeded.

Connections are encrypted per job and ordered Gemini, Groq, OpenRouter. The form requires free-account declarations for Gemini/Groq; OpenRouter additionally enforces zero catalog pricing and zero-price routing. No paid fallback or billing activation is implemented. Runtime model discovery replaces unverified model presets.

A durable 240-second deadline and 12 API admissions, including discovery, survive retries/restarts. Groups contain at most four sections with a 2,400-token response cap. PostgreSQL coordinates provider/global admission and conservative token/request windows. One job submits sequentially; separate workers may call different providers concurrently. Invalid output has one repair opportunity; quota/network failures spend the global allowance without using that repair allowance. No API result may overwrite a manual edit or an already accepted model section. Explicit single-section regeneration renews the budget for that section only.

The previous unsupported performance estimates are superseded by measured smoke results and explicit live-provider limitations. No local inference adapter or milestone 8 automated suite is included. Existing dirty changes from earlier work were preserved on the new feature branch.

## 2026-11-04 — Restore chapter context and separate content from execution

The supplied failed run showed that line harvesting plus four tiny fact excerpts lost the program's meaning. This correction supersedes milestone 9's arbitrary four-leaf/2,400-token grouping. Profile contracts now use logical chapters, up to eight leaves with independent IDs; per-call output budgets scale with chapter size and account TPM (up to 6,000 tokens). Fit the entire prompt including shared context and repair material. Input-token estimates are labeled estimates; reported provider usage reconciles reserved TPM after successful calls. A temporary 429 mentioning billing in a generic message is no longer classified as daily quota exhaustion.

Use standard-library Go AST and JSON parsing, bounded SQL declaration recognition, paragraph/source chunks and canonical Swagger JSON. Keep source bodies because a declaration index alone cannot support workflows. Store a versioned static model in the existing evidence JSONB. Route/call references are observed syntax, not proven dynamic edges. The shared brief is deterministic and adds no API call. Deep YAML and dynamic resolution require separate evidence before expansion.

A model response controls only body/citations/gaps and stable IDs. Headings and provenance are server-owned. Ordinary square brackets in prose are permitted. Local schema validation reports exact member failures, saves valid siblings, and retains private rejected output for targeted repair across restart. Migration 10 stores diagnostics and rejected prose on existing cascading provider-call rows; rejected content is not exposed through diagnostic endpoints or public logs. This is private job data with existing deletion/expiry rules.

The content report (`outline`, `partial`, `draft_ready`) is derived from persisted sections and does not certify factual correctness. Rendering status remains independent. Explicit continuation renews the 240-second/12-admission budget only for remaining baseline sections and preserves all accepted/model/user text. No automatic continuation runs after the deadline.

The recovery was implemented in the existing dirty feature branch to preserve prior user work. Automated suites remain deferred. Controlled responses establish transport behavior only; the live model quality gate remains open. Source examples guide prose style, never executable templates or unsupported project requirements.

## 2026-11-04 — Align academic output with the supplied FKN layout

The user requested the FKN template family and conspicuous manual-completion markers. Academic profile assets now identify layout version `fkn-2`. Adapt the supplied MIT-licensed shared layout into embedded `fkn.typ`, retaining the approval sheet, separate title sheet, vertical inventory tables, page counter reset after the approval sheet, body page margins and bottom inscription, centered chapter headings, and ten-column change-registration sheet. The generic technical overview keeps its independent layout. The operator-message tail is unnumbered as in the supplied reference; stable section IDs and saved prose are unchanged.

Data still enters through literal JSON/text, never executable model or uploaded Typst. Personal names, dates, signatures, project code and missing requirements are not fabricated. Empty fields and section gaps render with `TODO` in dark red; literal TODO text in prose and table cells uses the same style. Institution/faculty/program labels default to the FKN Software Engineering profile and remain editable. The bundled Liberation Serif fonts remain the deterministic, distributable substitute for the reference's Times New Roman; there is no dependency on system fonts or remote packages. The source copyright/license is retained in `FKN-LICENSE` and exported with the bundle. Profile annotations describe actual section headings instead of copying unverified compliance assertions or unrelated project claims.

The editor can clear a section's missing-information warnings only through an explicit checked “missing information is filled in” action on save. Ordinary saves retain warnings; other sections and stored artifacts are unchanged until an explicit rebuild. Signature/date TODOs are deliberate paper-signing placeholders. Blank archival/change-registration cells remain form cells rather than invented change records. Previously downloaded bundles retain their own layout and remain reproducible.

## 2026-11-05 — Live provider measurements and request economy

Measured with the owner's keys: Gemini `gemini-2.5-flash` answers a JSON request in about 1.7 s; Groq `openai/gpt-oss-120b` in about 0.3 s but publishes a hard `x-ratelimit-limit-tokens: 8000` per minute; OpenRouter free models are unstable (Gemma returned an upstream 429), while `nvidia/nemotron-3.5-lightning:free` answered in about 2 s. The create form now prefills these IDs and the matching per-minute values (Gemini 10 RPM, Groq 5 RPM / 8000 TPM, OpenRouter 20 RPM).

Model identity is compared after lowercasing, dropping a leading `models/` prefix (Gemini's OpenAI catalog uses it while chat does not) and a trailing `:latest`. Published per-minute limits from `x-ratelimit-limit-tokens` override a larger form value for later calls, so the worker fits the next chapter instead of repeating 429s. A context-window rejection no longer disables the provider: it stays temporary and the next attempt halves the evidence pack. The key field accepts several keys separated by commas or new lines; each key is a separate connection with its own rate window and quota, which rotates load without extra services.

Ensemble generation and parallel chapter calls are deliberately not implemented: free tiers are the constraint, and merging two weak drafts costs twice the quota without a measured quality gain.

## 2026-11-05 — Local model provider

A `local` provider joins the free cloud providers. It uses an OpenAI-compatible endpoint from `LOCAL_LLM_ENDPOINT` (default `http://127.0.0.1:11434/v1/chat/completions`), needs no API key, and omits the `Authorization` header. Free-tier confirmation and price checks do not apply; RPM stays bounded to 1–60 and TPM to 1,000–100,000, and a request gets 180 seconds instead of 45 because CPU-only generation is slow. The create form exposes an explicit opt-in with `qwen2.5:7b-instruct` as the default model, and Docker Compose runs Ollama with a persisted model volume; `make models` pulls the model.

The worker still verifies the model list at the endpoint before the first call, with local-specific messages when Ollama is not running. Expectations remain honest: the local path is slower and produces weaker prose than hosted free tiers, and its throughput and JSON-validity were not measured because Docker is unavailable in the development sandbox.

## 2026-11-05 — Multi-language extraction and BM25 ranking

Source extraction now covers more than Go without adding a runtime dependency. TypeScript/JavaScript files provide imports, functions, classes, arrow functions, route registrations and outgoing HTTP calls; Python files provide imports, functions, classes and route decorators; YAML OpenAPI files provide operations and component schema names through an indentation-based reader. Each declaration is one fact bounded by brace or indentation depth (at most 200 lines) and clipped to 6,000 runes, so a long body keeps its signature. The analysis version moves to `source-context-4`, which invalidates section caches built from older evidence.

Section and chapter packs rank facts with BM25 over kind, path and value, fused with the existing kind bonuses; substring counting is gone. Identical fact values are removed before ranking, and chapter packs keep the service round-robin so one service cannot consume the budget. This is a deterministic relevance score, not a semantic model: it cannot infer meaning that the source does not state, and runtime behavior remains unverified.

## 2026-11-05 — Markdown contract for model output

The generation prompt now names the supported Markdown subset — paragraphs separated by blank lines, `**bold**`, `*emphasis*`, lists with `-` or `1.` indented by two spaces, `>` quotes, inline code, `[text](https://…)` links and GFM pipe tables — and forbids Typst commands, math and raw HTML. The chapter prompt version moves to `chapters-4`, so previously cached model sections are not reused with the old contract.

Model bodies pass through `normalizeMarkdown` before validation: a Typst directive at the start of a line (`#set`, `#figure(`, `#text(`, `#link(` and similar) is escaped to literal text, and an unterminated code fence is closed so it cannot swallow later paragraphs. Unclosed inline emphasis or brackets stay literal text because the parser never drops content. Rendering stays allowlisted: `rich.typ` receives leaf text as Typst strings, so no model output can become executable Typst even when the text contains `#`, `$`, `[` or `@`.

## 2026-11-05 — Profile-scaled generation budget

The fixed 12-admission and 240-second model budget could not cover long profiles: late chapters never received a call. The budget now scales with the number of chapters a profile produces: `12 + 2 × chapters` admissions and `240s + 45s × chapters` of model time, overridable with `DRAFT_BASE_ATTEMPTS`, `DRAFT_ATTEMPTS_PER_CHAPTER`, `DRAFT_BASE_SECONDS` and `DRAFT_SECONDS_PER_CHAPTER`. The worker writes the profile deadline once before the first call; a run with recorded attempts keeps its stored deadline, so a restart cannot renew the budget.

Scheduling is coverage first. Each pass gives every unfinished chapter at most one call, so a repair on the first chapter cannot starve the rest of the document; at most three passes per budget are made. After the deadline the worker grants up to two bounded continuation passes of 90 seconds for chapters that still have baseline text, preserves every accepted section, and extends the admission limit only for those remaining chapters. This supersedes the earlier statement that no automatic continuation runs after the deadline; the explicit "Continue" action still renews the full profile budget.

Generation details now report the covered chapter count, the used and allowed model calls, and the providers disabled during the run with their error codes.

## 2026-11-05 — Submission friction, archive tolerance and paper signatures

The first submission form asked for more decisions than the task needs. One source now excludes the other: choosing a ZIP clears the repository field and typing a repository clears the chosen ZIP, so the request is never blocked by two sources at once. The user selects the document language; only an explicit earlier choice is restored locally, and it is never inferred from the browser. The previous profile is remembered the same way, the year defaults to the current year and stays editable in the project fields, and the project title is derived from the repository or archive name when the field is empty. The technical specification remains required for academic profiles; the remaining metadata stays optional with visible red TODOs.

Archive intake no longer rejects a whole project because of one entry. Symlinks, special files and duplicate names are excluded from the inventory with explicit reason codes. Only unsafe paths and file/directory conflicts fail, and the message names the offending entry. The upload request inspects the ZIP central directory before creating a job, so the user receives that reason immediately instead of after a failed job. GitHub archives follow the same rules in the worker.

Signature and signing-date fields are paper placeholders: they render as blank lines, not as missing data, because they are filled by hand after the draft is approved.

## 2026-11-05 — Provider rejections and outline documents

A technical-overview run finished as `completed` while every section stayed a baseline placeholder. Persisted diagnostics showed Groq returning HTTP 400 `json_validate_failed` ("Failed to generate JSON") for the chapter call. The worker classified every unknown rejection as a permanent configuration error, disabled the only connected provider, and never retried, although the 12-admission budget was almost untouched.

Provider rejections now follow the model output they describe. A malformed-JSON rejection becomes `provider_invalid_output`: the provider stays enabled, the rejected generation is stored as private job data, and the next attempt receives it as repair material with an explicit instruction. Context-window rejections become `provider_context_budget`. A free account may publish per-minute input/output token windows smaller than the account TPM entered on the form; the worker reads those limits from the rejection and fits the next chapter and response to them instead of repeating an oversized request until the budget is gone.

Content outcome stays independent of execution status. Completed jobs that are not `draft_ready` now show a visible warning on the job and sections pages with the recorded reason and the next step; the generation details block remains the full diagnostic view. Baseline sections still render with `TODO` markers, so an outline is never presented as generated prose.

## 2026-11-05 — Distinct SRS and arc42 document profiles

The user requested an SRS profile and arc42. The supplied architecture URL identifies ISO/IEC/IEEE 42010:2022, not arc42; implement the explicitly named arc42 format using its official twelve-part structure. SRS is an independently authored requirements specification informed by ISO/IEC/IEEE 29148:2018's public scope. Neither a complete licensed standard reproduction nor formal conformity is claimed.

Profile assets own optional generation guidance and logical chapter labels. Stable leaf IDs, profile guidance and section guidance participate in section cache keys through the section plan. SRS has 22 editable leaves in seven groups; arc42 has 12 numbered sections in six related groups. Both retain the existing 12-admission and 240-second model budget. Source selection includes build/deployment and test files for these profiles without executing them. Model instructions distinguish requirements from observed implementation, proposed decisions and inferred scenarios; unknown thresholds/approvals are TODO. Traceability must not invent requirement IDs from unseen generated chapters.

Both profiles use neutral document rendering without FKN university defaults, approval sheets, inventory frames or change registers. Display the actual profile identifier/version and localized annotation. User requirements/context are optional input so a source-based architecture draft or SRS outline can be produced when no complete specification exists. Missing stakeholder decisions remain visible. Existing jobs and profile IDs are unchanged; no database migration is required because profile IDs are text and evidence/sections are JSONB.

The adapted arc42 asset credits Gernot Starke and Peter Hruschka and uses CC BY-SA 4.0 for the template, separately from project-specific content. Attribution and adaptation notes appear in the document annotation and the exported ARC42-NOTICE.md. No external runtime template dependency is added.

References checked on 2026-11-05:
- https://www.iso.org/standard/72089.html — ISO/IEC/IEEE 29148:2018 public scope.
- https://www.iso.org/standard/74393.html — ISO/IEC/IEEE 42010:2022 architecture description.
- https://arc42.org/overview/ — twelve sections.
- https://arc42.org/license/ — template attribution and adaptation terms.

## 2026-11-05 — Build the browser interface from `web`

The interface lives in `web/` so presentation assets stay outside the Go delivery package. `web/templates` holds the Go page templates, `web/src` the TypeScript and CSS sources, and `web/static` the Vite output. The `web` Go package embeds the templates and built assets; `internal/web` only routes requests and renders pages. This keeps a single binary, allows `go build ./...` without Node installed, and still gives the frontend normal tooling. Rebuild with `make ui` and commit the refreshed `static/app.css` and `static/app.js` with their sources.

Light and dark themes are resolved from one token set with CSS `light-dark()`: the system preference is the default and the header menu stores an explicit choice. Pages work without JavaScript; scripts add theme selection, file-name previews, auto-growing textareas, and submit feedback.

## 2026-11-05 — Safe Markdown rendering and bibliography separation

Model output remains structured JSON. Section bodies may use CommonMark with tables and strikethrough; they are never executable Typst. Goldmark v1.7.13 parses body text into a small server-owned node tree; trusted `rich.typ` renders those nodes. This dependency replaces an unsafe/incomplete regex-based Markdown conversion. HTML and Typst-looking text remain literal, images do not load files or networks, and only HTTP(S) links are active. Existing editable section bodies and cache records are retained. Rebuilding a PDF applies the renderer without an LLM request.

Internal evidence IDs and file locations remain in evidence JSON and the website, not the PDF bibliography. Non-user bibliography sections use deterministic external documentation suggestions based on selected source technologies. They do not claim that the author consulted these resources. The initial catalog covers Go, TypeScript, JavaScript, Python, Typst and PostgreSQL (identified through Go drivers); unknown technologies require manual bibliography completion. User-authored bibliography sections are preserved. No articles, access dates or standard compliance are invented.

Assistant metadata (`.claude`, `.codex`, `.agents`, `.cursor`, `.codegraph`, `AGENTS.md`, `CLAUDE.md`) is excluded from new analysis, including stale selected inventories. These files dominated the supplied sample PDF and described development helpers instead of the uploaded application. Analysis and chapter prompt versions are bumped. Existing generated prose is not silently regenerated; correcting its subject requires a new analysis/generation.

Forms progressively enhance submission with inline errors, retain fields/uploads, and preserve their original encoding. Academic specification validation accepts text, a file, or complete structured fields. Optional title metadata omissions require an explicit continue-with-TODO action on creation. Required fields open collapsed panels before reporting validation. Server validation remains authoritative; JavaScript does not bypass it.

## 2026-11-05 — Harness wave 1: recoverable document journeys

Apply the supplied D/U/A rules through the existing milestone 9 scope, with the requirements and manual checks indexed in [harness-wave-1/plan.md](harness-wave-1/plan.md). The imported teaching demo is not product evidence. Three fictional personas reviewed the existing implementation independently twice; their responses are hypotheses, distinct from browser/HTTP/SQL observations. No application provider keys were used. Root milestone instructions continue to defer automated suites.

Saving section prose now retains structured tables. Section-save requests accepting JSON return `{saved:true,url:...}`; ordinary browser POSTs retain their redirect. The frontend saves in place, tracks dirty forms and in-flight revisions, and warns on navigation. A save response timeout preserves local text and does not assert that the server rolled back. Regeneration still deliberately replaces a section and now requires explicit confirmation; version history is deferred. Other pending inputs are not stored in browser persistent storage.

Content reports expose an additional `baseline` count. UI compilation status is separate from content quality, and continuation is offered only where baseline sections exist. Browser polling reports stale/disconnected state, times out, supports manual refresh and preserves lease-recovery labeling. Submission validates existing file/UTF-8 limits and offers explicit removal/correction; bilingual profile guidance and accurate server/provider-processing disclosure replace ambiguous local-only wording.

## 2026-11-05 — Harness wave 2: repeatable document setup and visible usage

The second bounded milestone 9 UX correction covers five changes indexed in [harness-wave-2/plan.md](harness-wave-2/plan.md). Purpose shortcuts select existing profiles and technical overview becomes the first-visit default; explicit saved profile/language preferences still win. Language support remains RU/EN. Unsupported PDF/Typst specification input now has text-copy guidance.

An explicit opt-in stores only allowlisted creation text, repository address and model IDs in browser local storage. Keys, files, provider confirmations and enabled states are excluded. Forget removes the creation draft; existing profile/language preferences remain independent. This revises wave 1's no-persistent-pending-input default only when the user explicitly opts in. New jobs still need credentials/uploads; existing jobs reuse their server-side connection/source. No credential vault, ownership transfer or new retention policy is introduced.

Progress pages expose saved sections and recorded provider usage, with bounded manual refresh. Diagnostic calls include discovery and failures; the endpoint caps at 100 entries, so the UI labels possible truncation. Token totals are recorded usage, not a billing statement or complete cost estimate. No-generation outcomes and PDF compilation remain distinct. Sections link directly to review/rebuild and describe editing availability. No generation policy or model-quality claim changed.

## 2026-11-06 — Live generation quality and reusable delivery

The five-change follow-up is indexed in [live-quality/plan.md](live-quality/plan.md). Actual model runs use isolated storage and explicit user authorization. Persona simulations remain separate and use no provider keys. A controlled Gemini 2.5 Flash baseline demonstrated truncated structured output and an empty model result; source review of the first complete iteration additionally exposed a false error-fallback claim.

Generation now reserves more output, splits chapters into at most three leaves and bounds thinking only for the supported Gemini 2.5 Flash endpoint/model. Source selection favors implementation with a soft per-file penalty, preserves compact manifest facts, excludes license/import/call-name noise, and retains both ends of long declarations. Prompts distinguish bundled samples, source behavior and unverified execution, and constrain each section's scope. Known inline source markers move to citation metadata; unknown fact markers and empty normalized prose are rejected. Evidence provenance is retained; no runtime correctness or universal model-quality guarantee follows.

References: https://ai.google.dev/gemini-api/docs/openai documents the nested `extra_body.google.thinking_config` REST field, checked on 2026-11-06. Private provider keys, cookies and rejected model responses are never committed.

Retained-job reuse creates an independent workspace/ZIP and a fresh job, with atomic owner/expiry checks and digest verification. An unchecked explicit choice re-encrypts existing connections for that new job; no browser key vault is introduced. The URL-encoded form limit is 1 MiB to accommodate legal UTF-8 fields within the existing per-field limits. Failed credential decryption does not block source-only reuse.

Catalan expands neutral document content only (overview, SRS, arc42 and handover); academic FKN profiles stay RU/EN. Migration 00011 persists this distinction without relabeling old jobs. A dedicated handover profile becomes the first-visit purpose default; stored explicit profile/language choices still win. This revises the wave-2 overview default and RU/EN-only scope explicitly.

Section reading uses explicit refresh rather than periodic page replacement. A history-entry checkpoint stores only section ID/coordinates; dirty and in-flight forms block refresh. Native anchors and query-bearing refresh links work without JavaScript. Version history, live prose replacement and cross-device identity remain deferred.

Later source audits in the same wave exposed unsupported hardware/economic claims and a visible TODO omitted from model gap metadata. PZ section guidance now explicitly separates runtime/build dependencies, historical decisions, proposed rationale, intended benefits and measurements. Visible TODO lines are conservatively included in unresolved metadata. Keyword relevance is capped so it cannot overwhelm typed priorities; test functions do not receive production entrypoint bonuses. Shared context is bounded and trimmed before useful implementation excerpts.

Actual provider failures also required compatibility repairs: unwrap one-item error arrays before quota classification, and decode selected-model scalar prices without assuming every model's pricing metadata is flat. Unknown/nonzero charges remain rejected. These are evidenced Q1 reliability repairs, not a new provider, paid fallback or additional improvement wave.

Gemini 3.5 Flash-Lite returned unquoted object keys in JSON-object mode during the real PZ run. Enable the existing strict chapter schema for explicitly supported Gemini 2.5 Flash, 3.5 Flash-Lite and 3.8 Flash models, and preserve the actual schema flag in generation diagnostics. Official compatibility documentation supports this mode; a manual local transport check verifies the payload/diagnostic contract, and a real Catalan handover used the negotiated schema. Factual source review remains separate from schema conformance.

## 2026-11-06 — Consolidated main and architecture handoff

The user requested integration of main, harness wave 2 and live generation quality. After fetching the remote, both feature branches and remote main formed one ancestry chain. Main fast-forwards to the complete live-quality implementation without conflict resolution or discarded changes. Preserve unrelated tmp_example work. Remove only the two requested feature branch references after main is successfully pushed; their commits remain reachable from main.

[architecture-handoff.md](architecture-handoff.md) records the current architecture and proposed evidence-first evolution, bounded future waves, research metrics and a portable next-agent assignment. Its H1 and later waves are proposals, not implemented capabilities or authorization to start milestone 8. A fresh builder review found no important unanswered planning question. The merged tree passed Go build and frontend typecheck/build; automated suites remain deferred by the existing user policy.

## 2026-11-06 — Railway deployment increment

The user requested actual deployment and selected Railway: separate web and worker processes in one application container sharing /data, with PostgreSQL as a separate private service. Google sign-in is explicitly deferred. The prior local-only plan statement is superseded only for this deployment increment; milestone 8 and the architecture handoff proposals remain separate.

Use a dedicated Railway image and supervised startup without changing the local Docker entry point. Initialize volume ownership, drop privileges, run migrations, prepare one persistent encryption key, and stop both processes when either fails. An explicit COOKIE_SECURE option protects ownership cookies behind managed HTTPS without trusting arbitrary forwarded headers. Keep one replica and document volume downtime and backup requirements. No uploaded project code is executed, no local secrets are copied and no LLM keys are needed for the deployment smoke check. See deployment-plan.md and deployment.md.
