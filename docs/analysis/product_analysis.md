# Product and implementation analysis

Status: proposed design, not an implemented feature list.

## User outcome and core workflow

1. User submits a public GitHub URL or a ZIP archive and optional description. Private repositories require a separate authorization decision.
2. User selects a document profile, document language, and project metadata required by that profile.
3. The service validates the input and displays the files and evidence it intends to use. Users may include or exclude files.
4. The service extracts a structured project inventory, drafts sections, and shows uncertainties and missing information.
5. User reviews and edits sections, then downloads the Typst project and PDF.

The optional description is user-provided context; it must not silently override repository facts. Output should distinguish `observed in source`, `provided by user`, and `unknown`.

## Processing pipeline

```text
URL / ZIP + user description
  -> intake validation and isolated unpacking
  -> file inventory and profile-aware filtering
  -> deterministic extraction + source references
  -> section plan and evidence allocation
  -> section drafts with explicit unknowns
  -> factual and structural checks
  -> Typst template rendering
  -> isolated Typst compilation and PDF checks
  -> review, edit, download, retention cleanup
```

### Intake and filtering

- MVP inputs: public GitHub repositories at a pinned commit and ZIP archives. Record the commit hash or archive digest for repeatability.
- Enforce compressed and expanded size limits, file count and path-length limits, archive depth, and per-file read limits before processing. Reject path traversal, symlinks, special files, nested archive expansion, and encrypted archives. Never run uploaded project code.
- Skip dependency caches, generated outputs, binaries, secrets, and very large files by default. Detect secrets before sending content to any model. A filename alone is not sufficient for secret detection.
- Do not blindly exclude tests: include test metadata and relevant test files for testing or quality profiles. Keep the selection visible and user-correctable.
- Treat README, manifests, API schemas, CI configuration, and source structure as different evidence types. Preserve source paths and line spans where available.

### Extraction and generation

- Build a typed intermediate representation: project summary, languages, modules, entry points, interfaces, dependencies, deployment hints, tests, and cited evidence. Give each fact a confidence or `unverified` state.
- Use deterministic parsers for manifests, OpenAPI, tree structure, and CI configuration; use model inference only where deterministic extraction is insufficient.
- Plan required sections from the chosen profile, then retrieve relevant evidence per section. A section cannot present a technical claim without a source reference or a user-provided attribution.
- Validate the section schema, mandatory headings, unresolved placeholders, citations, and contradictions before rendering. Use `Information needed` rather than invented requirements or performance figures.
- Render data into versioned Typst templates. Keep generated content separate from template code, escape user text, compile in a restricted environment, and inspect compilation diagnostics. Offer a source bundle with its assets so the user can edit it.

## Selected MVP direction and remaining scope choices

The product owner selected students preparing coursework or thesis projects, an existing academic document set adapted to the website, a generic technical documentation family, and hosted LLM generation with checks and caching (2026-11-02). The first concrete profiles are the existing academic set and a **technical project overview**. The academic set includes a technical specification, explanatory report, test programme and method, team test programme and method, and operator guide. The technical specification will be uploaded or entered through a structured form, then used as source material for the other documents. See [template_adaptation.md](template_adaptation.md) for the migration contract. The first release accepts public GitHub repositories and ZIP archives, optional description, file-selection preview, section preview/edit, Typst download, and PDF download. Each job selects Russian or English output; the English academic template retains the existing title-page structure with translated fixed text. Account, billing, private-repository OAuth, and collaboration are later decisions. This is analysis scope only: no website or template migration is authorized by the current decision.

## Architecture recommendation to test

| Component | First implementation | Why |
|---|---|---|
| Web | Go HTTP server and server-rendered pages with modest JavaScript | Small deployment and clear upload/progress flow |
| Metadata | PostgreSQL | Jobs, documents, template versions, audit metadata |
| Background work | PostgreSQL job table with Go workers | Asynchronous processing without a message broker |
| Artifact storage | Local filesystem in development; object storage before multi-instance deployment | ZIP, editable Typst bundle, PDF |
| Generation | Replaceable provider interface, initially hosted LLM API | Fast baseline for measuring quality and cost |
| Renderer | Pinned Typst version in an isolated worker | Reproducible builds and controlled file access |

### Asynchronous job contract

1. `POST /jobs` validates basic input, document profile and `language` (`ru` or `en`), stores the upload or pinned repository reference, inserts a `queued` job in PostgreSQL, and responds quickly with `202 Accepted` and a job ID. It does not wait for LLM generation or PDF compilation.
2. A Go worker process claims a queued job in a short transaction using `FOR UPDATE SKIP LOCKED`, sets a lease and attempt count, then commits before doing expensive work. Job status is separate from work stage; the [lifecycle contract](../development/job_lifecycle.md) defines both.
3. The browser polls `GET /jobs/{id}` for stage, completed sections, required action, retry time, or failure; it can later use server-sent events if polling proves inadequate. `GET /jobs/{id}/artifacts` returns authorized links to the Typst bundle and PDF only after completion.
4. Workers renew leases during long tasks. Expired leases are retried with bounded attempts and idempotent artifact writes; permanent validation errors fail immediately. Store a short, safe error message for the user and a detailed internal diagnostic.
5. A single Go binary can expose `serve` and `worker` commands, allowing both processes in development or separate replicas in production. PostgreSQL is the durable queue and status store; work is not kept only in a Go goroutine, because a server restart would lose it.

Kafka is unnecessary for this request/worker/result workflow. Revisit a broker only if measured queue throughput, multiple independent consumers, or event replay requirements justify it. Redis is also unnecessary for the initial polling and job lifecycle.

## Generation improvement sequence

1. Establish an evaluation set: consented repositories, target documents, reviewer rubric, and section-level failure labels. Track factual support, missing important facts, template completeness, compile success, edit time, latency, and cost.
2. Improve extraction and prompts: evidence retrieval by section, structured outputs, validation, retry only failed sections, caching by commit/profile/template/model version, and model routing by task.
3. Compare hosted models and an open-weight local model on the same evaluation set, including privacy and operating cost.
4. Consider fine-tuning only if recurring errors remain after better context and validation, and enough licensed, high-quality training examples exist. A custom model from scratch is outside the initial plan.

## Privacy, trust, and product boundaries

Uploaded code may contain personal data or secrets. Set short retention by default, delete temporary extraction data, encrypt stored artifacts, restrict job access, and obtain explicit consent for use of uploaded content in training. For EU users, map data roles and processing against the [GDPR](https://eur-lex.europa.eu/legal-content/ENG/ALL/?uri=celex%3A32016R0679), especially minimization and processor obligations. These are design requirements to review with counsel before a public launch, not a claim of legal compliance.

Academic output should support attribution and user review. It cannot guarantee acceptance under a university's rules or claim that the code proves requirements, results, or test coverage that it does not prove.

## Release stages and gates

| Stage | Deliverable | Exit evidence |
|---|---|---|
| Discovery | Interviews, real artifacts, chosen niche and profiles | Repeated job and review constraints documented |
| Vertical slice | Public repo/ZIP to editable Typst and PDF for one profile | End-to-end artifact compiles and can be edited |
| Quality baseline | Blinded review on held-out projects | Failure modes and correction time measured |
| Usable beta | File preview, job state, section editing, retention controls | Users complete a real task without assistance |
| Optimization | Improved retrieval, prompts, caching, model comparison | Measured quality/cost gain against baseline |
