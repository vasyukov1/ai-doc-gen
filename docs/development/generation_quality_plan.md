# Generation quality recovery plan

Status: implementation authorized on 2026-11-04; corrective changes delivered for milestone 9, with live-model content acceptance still pending. This does not begin the final testing milestone. Repository documentation remains in English. Existing work, jobs, edits, ownership and artifact safety must be preserved.

## Goal and priority

Restore the original product purpose: produce coherent, project-specific academic/technical documentation through a small website. A five-minute target means useful generated content, not merely a compilable PDF. Quality and honest completion take precedence over silently replacing prose with source excerpts. Keep free-provider support; never enable paid billing or paid fallback without an explicit user choice.

No runtime dependency on an earlier implementation is allowed. Reference behavior is a specification, not executable application code to import.

## Evidence from the failed run

Read-only investigation on 2026-11-04 inspected the supplied 20-page PDF, the current implementation and persisted records for job c705d56f-b168-49d1-8d23-9087383d95bf.

- The document contains 34 baseline sections and one model section. The model section itself only reports missing information about the development topic. Provenance is not a quality score.
- There are three recorded Groq calls: catalog discovery; a response with 642 input and 1,074 output tokens, partially rejected as invalid; and an error classified as quota exhaustion. The exact provider error body/status details and per-member validation failures were not retained. It is not established that this job exhausted the account's quota or that the classification was correct.
- Intake selected 61 files and extracted 475 records: 175 code declaration lines, 120 API-schema lines, 85 dependency entries, 61 inventory entries, 20 README excerpts and a few metadata/limit records. This is an excerpt inventory, not semantic analysis.
- Input/output requirements in the PDF contain opening Swagger syntax instead of endpoint/DTO descriptions. Functional-purpose sections repeat README headings. The supplied specification is a short project description, not detailed acceptance requirements.
- Extraction discards code bodies, skips lines longer than 300 bytes and limits most files to 20 lines. Retrieval then picks at most four keyword-matched facts, clips each to 450 characters and further prunes group context. No shared semantic project model or component relationship analysis exists.
- The current prompt asks for 100–180 words per leaf, with a 2,400-token cap shared by four leaves. Validation requires exact model-echoed headings and inline citation syntax; every failure collapses to the same error code. The specific reasons three members of the first response were rejected cannot be reconstructed.
- Historical behavior supplied complete bounded file excerpts, full specification, chapter templates, project metadata, notes, chapter-specific instructions and optional style examples. PZ used three generative chapters; PMI and operator guidance used four each. TZ was input/static template material, not 35 generative leaf calls. Preserve that distinction when comparing results.

## Alternatives and chosen direction

Switching free providers alone does not restore information already discarded before the prompt. Adding many summarization calls or a graph database first would introduce cost before proving a quality improvement. The recommended approach is to restore chapter-level context and prompts as a control, then introduce deterministic extraction and one reusable project brief. Add analysis only when an equal-input comparison shows a benefit.

## Q1. Make failures explainable and stop reporting skeletons as success

- Persist HTTP status, provider request ID, typed error code, sanitized message, Retry-After/reset information, reported usage, finish reason, stage, prompt version and request/response sizes. Never log credentials, raw uploaded source or full prompts.
- Return structured validation failures with section ID and field/reason: invalid JSON, missing section, wrong ID, invalid source ID, truncation, empty content or unsupported output. Separate structural validity, evidence coverage and content quality.
- Preserve rejected generated content as a private job artifact under existing retention/ownership rules so repair does not require regenerating accepted text. Raw diagnostics stay out of public logs and final PDFs.
- Classify failures per provider using status and structured fields. Do not infer exhausted daily quota from arbitrary occurrences of words such as billing or credit. Unknown 429 is an observable bounded cooldown unless specific evidence identifies a permanent quota failure.
- Add a content outcome independent of execution status: `outline`, `partial`, `draft_ready`. Baseline artifacts are optional outline downloads, never the primary successful documentation result. At a time limit, show saved prose, coverage and the actual remaining blocker.

Acceptance: every rejected member and provider failure has an actionable reason; a repeat of the supplied failure cannot be presented as a completed substantive document. No model call is needed to verify these diagnostics with controlled responses.

## Q2. Restore chapter-level generation as the quality control

- Give each profile an explicit chapter contract: mandatory headings, section purpose, required topics, source priority, fixed content, generated content and expected output shape. Populate headings and fixed metadata on the server; do not spend model output on echoing them exactly.
- For PZ generate purpose, technical characteristics and economic discussion as three logical chapters; for PMI generate goals, requirements, procedure and methods; for operator guidance generate its four practical chapters. Preserve independent leaf IDs for edits and progress even when one call writes a chapter.
- Keep supplied TZ authoritative for requirements. TZ creation from a short description is a separate drafting behavior: mark proposed requirements, missing decisions and code/requirement conflicts, never invent approved dates, performance targets or approvals.
- Build a bounded source pack containing project overview, full relevant README sections, full relevant specification passages, entry points, router/handler bodies, domain types and selected repository/integration code. Select by feature/service coverage rather than alphabetical order; deduplicate generated API docs and exclude dependencies/secrets.
- Use chapter-specific Russian/English prompts that request explanations, workflows, inputs, outputs, failure cases and tables where appropriate. Remove the global 100–180-word restriction and arbitrary four-leaf grouping. Fit input/output to the selected model's context and account budget, splitting only oversized logical chapters.
- Keep output as safe document data, not executable Typst. The renderer owns the template, numbering and headings. Use stable section IDs plus paragraphs, lists, tables and source references; reject unknown executable/content constructs.

Acceptance: on the same source snapshot and supplied requirements, the control produces a PZ with project-specific service roles, workflows and API descriptions. No braces, import declarations or raw README headings stand in for prose. It becomes the comparison baseline for later analysis work.

## Q3. Replace line harvesting with structured project understanding

- Parse README/requirements by sections and paragraphs. Preserve long meaningful paragraphs; report omissions instead of silently dropping every long line. Validate whether the submitted material is a project description, full requirements document or irrelevant text, without pretending a short description is a complete TZ.
- Parse OpenAPI/Swagger JSON as operations, parameters, request/response schemas and errors. Prefer one canonical representation over repeated generated Go/JSON/YAML documentation. YAML support may use a dedicated parser only if necessary and justified.
- For Go use `go/parser` and `go/ast`: imports, types/fields, function bodies, route registrations, constructor wiring and call sites. Resolve local paths statically; do not execute uploaded code, install its dependencies or run its build.
- Extract supported SQL table/column definitions and known messaging/configuration declarations. Unsupported constructs remain explicit unknowns. Initial deep support is Go plus OpenAPI, matching the failing project; other languages retain coherent bounded source excerpts instead of fake semantic support.
- Store typed entities and relations with source spans: service, endpoint, entity, storage, external dependency, producer/consumer and workflow step. Distinguish explicit source declarations, derived relationships, user requirements and unresolved assumptions. Dynamic dispatch and ambiguous calls must not become asserted facts.
- Keep this model in versioned JSONB; no graph database, vector store or embedding service is needed for this iteration.

Acceptance: the failing project yields an evidence-backed map of the gateway, order and payment services, their actual endpoints and persisted entities. Kafka and inbox/outbox relationships are included only where the source establishes them. Extracted relationships link to real file spans; none claims tests were executed.

## Q4. Build one reusable project brief, then evidence packs for chapters

- Define `ProjectModel`, `ProjectBrief`, `ChapterContract`, `ChapterEvidencePack` and `GenerationReport` in the domain layer. Business decisions stay in internal/documents; persistence stays in internal/storage and internal/pgqueue.
- Construct a deterministic project brief from typed facts first. Use at most one model synthesis call when coherent source packs cannot be assembled deterministically. Cache the result by source snapshot, selections, user requirements, analysis/prompt versions, provider/model and ownership boundary.
- Brief fields cover purpose, actors, services, entities, key workflows, integrations, constraints, terminology, requirement/code discrepancies and unanswered questions. Every technical claim retains supporting source IDs. A model-authored summary is not new evidence.
- Each chapter receives the same compact shared brief plus its requirements and relevant raw/structured evidence. Retrieve related facts across services so an order workflow does not lose its payment/Kafka steps. Expand source coverage for a genuine gap before asking a model to invent more prose.
- Set budgets in tokens appropriate to the provider; record estimated and actual usage separately. Remove the unconditional 6,000-byte pruning rule. Show a preflight request/token estimate and whether the account can plausibly finish within the selected budget.

Acceptance: the same project is explained consistently across chapters without repeating full source packs in every leaf call. The enriched route must match or improve the Q2 control on content before becoming default; fewer tokens alone do not justify worse output.

## Q5. Repair and scheduling that preserve useful text

- Use schema-constrained output when the chosen model supports it. Map headings from IDs on the server. Normalize only harmless representation issues locally; never fabricate missing evidence, claims or citations to make validation pass.
- Keep citations as structured metadata rather than requiring every pair of square brackets in prose to be a valid fact ID. Produce reader-friendly source notes or an appendix; retain source navigation in the web editor.
- Repair only rejected sections with the exact validation feedback and their previous response. Keep accepted siblings. On truncated output, split the unfinished chapter at section boundaries and request only the missing content; do not replay a complete paid/budgeted chapter blindly.
- Persist progress and budgets across worker restarts. Account for request/minute, token/minute and explicit daily limits independently. Cache model discovery, keep its calls visible but distinguish them from generation.
- Target one shared synthesis plus the profile's substantive chapter calls, not 35 leaf calls. Preserve the finite call budget; when a model/account cannot fit the job, report the tradeoff instead of silently substituting an outline. Keep provider fallback optional and never silently select paid models.
- Five minutes is the target for a useful semantic draft. At the deadline preserve partial work and offer explicit continuation with an estimate; never advertise a structural PDF as meeting the content target.

Acceptance: one schema failure does not destroy a chapter, retries do not overwrite edits, quota failures identify their class and reset condition where available, and continued generation resumes only unfinished work.

## Q6. Quality gate and evidence-based acceptance

Use the failing three-service project, a small single-service project and a project with genuinely missing requirements. Use the same pinned source, language, profile and model for before/after comparisons. Compare corresponding PZ/PMI/RO behavior; do not compare a formerly static TZ to a generated PZ. Keep provider capacity and pipeline quality as separate measurements.

Required manual acceptance before declaring recovery:

- Core service roles, domain entities, actual interfaces and at least two source-supported end-to-end scenarios are explained in coherent prose.
- Every mandatory heading is present. All profile-critical content is covered or explicitly identified as unavailable; no unresolved critical gap is hidden behind a success label.
- Claims do not contradict the source/requirements. Unsupported performance, security, testing and compliance claims are absent. Source references are resolvable; matching IDs alone is not proof of semantic support.
- There are no raw schema fragments masquerading as paragraphs, repeated generic gap boxes, irrelevant warnings on explanatory prose or duplicated content across chapters.
- Record time to first useful chapter and full draft, input/output tokens, calls by purpose, repair calls, model/account limits and reason for interruption. Human comparison must find the result at least as useful as the Q2 control.
- Exercise malformed output, genuine RPM/TPM/daily exhaustion, unsupported JSON mode, long source files, cancellation/restart, cache invalidation and isolated section regeneration.
- Compile and inspect RU/EN PDF output with literal untrusted text. Preserve private downloads and reproducible source bundles.

Follow the user's test sequence: build and small manual checks per corrective step; add the table-driven/integration suite only when milestone 8 is explicitly started. Synthetic provider success is infrastructure evidence, never content-quality acceptance.

## Execution boundaries

The subsequent user request authorizes the complete recovery. Implement the deterministic chapter/source control and bounded structured analysis together; keep Q6 content acceptance open until an actual equal-input model comparison. Do not infer successful quality acceptance from synthetic responses. No separate model synthesis call is enabled without evidence that it improves this control.

Before implementation, update the affected milestone 9 acceptance criteria and decision record to replace the earlier baseline-PDF success definition. New analysis and prompt versions invalidate only incompatible caches; existing saved documents and user edits remain intact. A/B comparison operates on new runs, not destructive changes to prior artifacts. No deployment architecture rewrite is required.

## Implementation and remaining acceptance

The implementation uses profile chapter groups (PZ 3, PMI 4, RO 4, TZ 9), preserving leaf IDs. Generation uses a shared deterministic brief and service/category-balanced source excerpts; budgets include the entire serialized prompt and repair context. The source analyzer extracts Go declarations with bodies, static selector-call and route references, canonical JSON OpenAPI operations/schemas, recognized SQL table declarations and coherent text chunks. These are static observations with source spans, not a fully resolved inter-service workflow graph. Deep YAML parsing and dynamic call resolution are not claimed. There is no additional model synthesis or embedding call.

Responses use server-owned headings and structured citation metadata. Schema mode is enabled when OpenRouter discovery advertises structured outputs; other connections retain JSON-object mode and local validation. Complete siblings can be recovered from a truncated JSON envelope; unfinished truncated groups are split for the remaining attempt. Recent successful model discovery is reused within the same job across worker restarts (five-minute freshness). Private rejected prose and exact member issues survive worker restarts in PostgreSQL under job retention. Migration 10 adds these diagnostics; Jet types were regenerated. The owner can inspect call diagnostics and explicitly continue baseline sections. Saved model/user sections are preserved.

Content outcome is separate from execution status. `outline` means no substantive generated/edited sections; `partial` includes unresolved gaps even when all sections have model provenance; `draft_ready` requires complete non-baseline sections without reported gaps and still requires human review. This is a coverage gate, not automated verification of truth. Outline downloads are labeled accordingly, with a source-location appendix instead of raw schema passages.

The supplied local Typst example was read as a content/style reference: connected explanations, concrete requirements and structured chapter hierarchy. Its claims, dates, names and executable Typst are not imported into generated documents. Generated tables use a bounded data DTO (caption, columns, rows) rendered literally by the trusted template. Plain-text section edits explicitly replace/remove its generated tables; full layout parity with the reference is not claimed.

Open acceptance: real-provider RU/EN comparison, semantic support of two end-to-end workflows, model/context capability measurements for each connected provider, and source-pack quality versus the chapter control. No five-minute content-quality claim may be made from the manual transport fixture measurements. See verification.md for measured checks.
