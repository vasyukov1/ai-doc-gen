# Production student journey and PZ quality

Date: 2026-11-06. Authorized scope: walk the deployed new-document form as a fictional FKN student, generate a complete Russian explanatory note using a real model, diagnose observed issues, fix them and repeat. This is a corrective increment, not milestone 8. Use native execution with independent review before delivery. Preserve unrelated tmp_example work.

## Evidence and design

The reader is the project owner deciding whether the deployed student journey produces a useful document. The persona supplies fictional cover data and a clearly labelled demonstration assignment about the public CLI fixture. It is not an approved university specification. Actual generation is explicitly authorized and kept separate from simulated UX hypotheses. No uploaded code is executed. Credentials and raw artifacts remain ignored and private. Public production requires a separate explicit permission to transfer the existing connection; received in this session.

Observed baseline: the new form is on production. The real run uses Gemini 3.5 Flash-Lite. The model picker offers only a different preset and a custom ID field. At 3 model sections, the progress card displays 17/17 and 100% because baseline availability drives its bar. Optional model setup is collapsed, while the enabled local-model control has no matching hosted inference service.

Approach: correct the existing form/status/prompt contracts rather than build a new wizard or inference architecture. Preserve saved inputs, edits, ownership, bounded provider budgets, and honest incomplete outcomes. Confirm technical prose against source; never remove genuine unknowns by fabrication. At most five accepted changes; additional findings remain deferred.

## Requirement/action/check index

| ID | Requirement and action | Files | Manual acceptance |
| --- | --- | --- | --- |
| SP1 | Make prose versus outline intent explicit before submission; keep provider setup understandable and supported model presets selectable | creation template, setup/forms, provider defaults as needed | Missing key cannot silently produce an outline when prose is selected; working model is selectable; custom IDs restore correctly |
| SP2 | Show substantive drafting progress separately from saved outlines | job template, job.ts, presentation helpers | 3 model leaves out of 17 does not show 100% prose progress; edits and compilation remain distinct |
| SP3 | Fix source-grounding or completeness defects evidenced by the real PZ | documents generation/retrieval and profile guidance, only after audit | Fresh raw output has substantive text in every generated leaf; source checks cover commands, architecture, error handling, dependencies and economic claims |
| SP4 | Show local inference only when configured on the server | configuration/web templates if confirmed | Hosted form cannot advertise an unavailable local service; configured local deployment retains support |

## Execution

- [x] Read requirements, existing deployment and quality records; start isolated branch.
- [x] Submit a new production PZ from the browser with fictional input and the authorized real connection.
- [ ] Capture terminal report, source-pinned generated text and PDF; audit raw output before corrections.
- [ ] Implement accepted fixes and build Go/frontend; perform manual pre/post checks.
- [ ] Repeat real generation and audit raw output; inspect the PDF and downloadable Typst.
- [ ] Independent review, documentation, explicit-path commit/push, integrate deployed branch and verify production.

## Guardrails and limits

UX findings are hypotheses, not user population measurements. Model sections and words measure coverage, not truth. Provider-reported token usage is not total billing. Separate raw output from editorial delivery; preserve failed attempts in the evidence report. Unknown approvals, runtime measurements and benchmarks remain explicit limitations. Do not claim general superiority to coding agents.

Resume: inspect production job e1126a09-85f3-432a-859c-e29eed88b7bf in the owned browser, then audit its saved sections. Command for source status: `git status --short`.

## Resume verification — 2026-11-09

The interrupted baseline completed with 14 model leaves, 3 baselines and 17 recorded gap entries. The report recorded 35,032 input and 8,417 output tokens across 10 calls including discovery/failures. These are observed provider reports, not billed totals. Baseline PDF and raw sections/report were retained privately. Source audit found repeated introduction, unsupported document-type/standards wording, a fallback bundling path described as the whole algorithm, unnecessary completion requests and duplicated table representations.

The independent reviewer found an accidental diagnostics JSON-tag regression in the new schema delta. It was corrected before deployment; a disposable HTTP transport smoke confirmed `json_schema` remains a boolean, scoped enums reach the payload and table limits match the validator. Go build, TypeScript checking/build and Docker build passed. Local browser checks confirmed missing-key prose blocking, explicit outline disabling keys, no unconfigured local option, and model preset restoration after reload. No automated suites were added.

The stricter schema binds citations to the actual fitted source pack and section IDs to the requested leaves. This prevents malformed identifiers, not factual hallucinations. Fresh model acceptance remains to be checked on production.

### Live schema compatibility correction — 2026-11-09

Production attempt 241b7239-2733-4a7c-87b8-e3f0cc7e56ae was rejected before drafting (HTTP 400). A minimal real Gemini probe reproduced the rejection with nested table array bounds and succeeded with those bounds removed, preserving section/citation enums. Remove nested array bounds from the provider schema; keep exact IDs, prompt limits and existing server-side table validation. This corrects SP3 rather than adding scope. Mock transport acceptance alone was insufficient evidence of provider compatibility. Private probe reports preserve both results.

### Single table representation

The corrected-schema raw run still duplicated pz-22 tables and rejected pz-15 for row width mismatch. SP3 now requests Markdown-only tables in body and omits the structured tables property from new model response schemas; legacy stored tables remain supported. Validate with real section regeneration and inspect the resulting PDF. This supersedes the earlier schema table-bound approach.
