# Production student walkthrough and generation audit

Cut-off: 2026-11-09. Reader: project owner/supervisor deciding whether the student journey produces reviewable explanatory notes. Scope: SP1–SP4 in [plan.md](plan.md); no next milestone or automated test suite started. Fictional FKN cover data and an explicitly unapproved demonstration assignment were used. This is an agent walkthrough, not evidence from recruited students (U1, U9, A1, A6).

## Evidence ledger

Raw JSON, PDFs, screenshots and probe responses are retained privately under ignored `tmp/production-pz`. No keys or owner cookies are published. Source fixture snapshot: `99b05b1ca25310546207361bbdca2c53f98a8063`, 59 retained text files; uploaded code was never executed.

| Evidence | Observation / reported value | Date, kind, grade | Covers / misses |
| --- | --- | --- | --- |
| `before-report.json`, job `e1126a09-85f3-432a-859c-e29eed88b7bf` | 14 model leaves, 3 baseline leaves; 17 recorded gap entries | Run Nov 6, inspected Nov 9; observed, A | Saved origin and completeness, not factual correctness |
| Same report | 35,032 input / 8,417 output tokens; 10 calls including discovery/failures | Provider-reported, D | Reported usage only, not billing or unreported usage |
| `before.pdf` | 18 compiled pages | Inspected Nov 9; observed, A | Compilation/page count, not a completed assignment |
| `schema-rejection-report.json`, job `241b7239-2733-4a7c-87b8-e3f0cc7e56ae` | 0/17 model leaves; catalog success then HTTP 400 generation rejection | Nov 9; observed, A | Failed first attempt on stricter schema; must remain in comparison |
| `probe-bounded.json`, `probe-unbounded.json` | Minimal bounded schema HTTP 400; otherwise identical unbounded schema HTTP 200 | Nov 9; observed, A | Gemini 3.5 Flash-Lite schema acceptance; not long-document quality |
| `mobile-form.png` and DOM dimensions | 390 px viewport and document width; full form screenshot height 6,784 px | Nov 9; observed, A | One synthetic mobile viewport, not physical phone usability |

Values above are single-run diagnostics (A3, A5, A8). A valid schema or citation identifier is not proof of a claim. No causal speed/cost claim is made between runs several days apart.

## Baseline content audit

- Introduction leaves repeated broad product descriptions instead of their distinct purposes.
- The prose conflated the presence of technical-specification templates with supported AI generation. The actual generation dispatcher accepts PZ, PMI, RO and team PMI; PDF compilation has a different selection contract.
- Bundling prose described lexicographic concatenation as the whole algorithm. Pinned `internal/typst/bundle.go` first handles a file, an existing requested entry, then `main.typ` / `body.typ`; recursive include expansion applies there. Sorted concatenation is the fallback.
- A comparison table appeared both in body Markdown and structured tables. Unsupported standards/comparison claims and unnecessary completion demands reduced usefulness.
- Missing formal approval and benchmarks were repeatedly treated as required input despite the demonstration assignment expressly excluding them.
- Raw output remained partial after invalid source identifiers, oversized tables and provider availability errors. A compiled PDF was not a success criterion for content.

Source-check anchors: pinned `cmd/term-paper/main.go` dispatch/flags/error exit, `internal/app/app.go` generation/compile/improve orchestration, `internal/typst/bundle.go` entry/fallback behavior. The baseline review is static; it does not assert that the CLI has passed runtime tests.

## Implemented corrections

SP1 makes prose versus outline an explicit choice, validates a missing prose connection on both client and server, offers selectable model presets and restores the selected preset consistently. SP2 counts drafted sections rather than stored placeholders and timestamps usage snapshots. SP3 narrows section guidance, requires exact scoped source/leaf IDs, reduces unsupported claims and duplicate tables, and preserves post-response table validation. SP4 hides and rejects local inference when no local endpoint is configured.

The first live SP3 schema revision exposed an additional provider compatibility failure. Nested array bounds were removed only after a minimal real probe reproduced HTTP 400 and the same schema without those bounds succeeded. This is a correction within SP3. Server-side limits remain enforced; syntax validation does not become semantic certification.

## UX hypotheses and deferred design

The laptop walkthrough confirmed missing-key blocking, explicit outline selection, preset persistence and saved-source/connection reuse. Reuse copied the source snapshot and input without requesting another upload or key. A new full draft does not copy prior prose. The phone-width walkthrough confirmed no horizontal overflow; long bilingual labels can be truncated inside selects, although the helper explains the choice.

Hypothesis, medium confidence: a single-provider selector followed by that provider's key/help would reduce scrolling. The form currently exposes all providers and is long even after conditional profile fields. Test this with actual students before claiming fewer errors or faster completion. A full interface locale selector would reduce repeated bilingual copy. Reuse guidance still uses the older English form and should share the new-document field components in a future bounded UX wave.

Opposite reading (A9): detailed visible help may benefit first-time students more than compact controls. Measure completion errors, requests for help and later edits before hiding guidance. This walkthrough cannot establish superiority over a general-purpose coding agent.

## Verification status

Go build, frontend typecheck/build, Docker build, missing-key HTTP smoke, browser preset/outline checks and a disposable transport diagnostic check passed. The real provider probe additionally checked schema compilation. No automated suites were introduced before the selected testing milestone. Post-correction results and the raw-versus-reviewed distinction are recorded below.

## Final observed results — 2026-11-09

The resumed production run reached 17/17 model-origin leaves and 0 baseline leaves, with 2,646 whitespace-separated body words and a 19-page PDF. The raw report still classified it as partial because 16 gap entries remained. Its 9 recorded calls (including the original schema rejection and catalog call) reported 39,637 input and 8,026 output tokens. These are observed saved counts / provider-reported usage, not quality or billing scores. Evidence: `raw-complete-sections.json`, `raw-complete-report.json`, `raw-complete.pdf` (A3, A5).

The full raw result was preserved before targeted regeneration. It still contained errors: a runtime requirement contradicted the later standalone Go description; input TZ and generated document types were mixed; bundling was conflated with PDF compilation; fallback concatenation was described as the whole bundling algorithm; unsupported bibliography/comparison claims remained; and pz-22 repeated the same table in two representations. The pz-15 structured table initially failed width validation and was repaired on the bounded retry. **The raw factual-quality gate did not pass.** All leaves containing prose establishes completeness only.

SP3 was therefore simplified again: `chapters-10` asks for Markdown tables in body only, with no separate structured tables field in new response schemas. A real minimal probe returned HTTP 200. Targeted pz-22 regeneration on production returned HTTP 200 with the new prompt version, a single Markdown table and no structured tables. Its semantic comparison claims still require review. Legacy stored tables remain readable. Count/dimension guidance for new Markdown tables is prompt-level; body size remains bounded. This limitation is explicit rather than claimed to be an enforced table contract.

After that targeted regeneration the cumulative saved report contained 11 calls, 46,307 input / 8,543 output tokens and 17 recorded gaps. These totals must not replace the pre-regeneration full-run totals. Three additional minimal diagnostic requests were outside the job report: one rejected bounded-schema request, one unbounded-schema success (23 input / 33 output tokens reported), and one Markdown-only success (29 input / 74 output tokens reported). Their purpose was compatibility diagnosis, not document-quality evaluation.

The progress card visibly showed 3/17, 8/17, 13/17, 16/17 and finally 17/17 during the walkthrough. A remaining partial-result banner incorrectly said sections were not generated even when every section had prose but some gaps remained; SP2 wording was corrected to say the draft contains unresolved gaps and needs review. Usage has an explicit snapshot timestamp, and compilation remains distinct from content verification.

## Editorial delivery and reproducibility

A separate reviewed version was prepared from the preserved raw full run, not presented as unattended model output. All 17 leaves were editorially changed or normalized; the reviewed body has 2,596 whitespace-separated words. Corrections cover runtime versus build dependencies, exact supported document types, command flow, the ignored generate-doc apply flag, the entry/include/fallback bundling branches, unsupported savings/competitor claims, and fabricated bibliography attribution. Unknown approvals, historical decisions, measurements, adoption and comparative results are stated as limitations rather than invented values. Specific official Typst, Pandoc and Doxygen pages were consulted only for this editorial version and are linked in its bibliography.

The reviewed document was compiled locally through this application's trusted `documents.Render` and pinned Typst 0.15.1, using the retained input and evidence. It was **not uploaded over the production raw draft**. A local artifact can therefore be reviewed without confusing it with the model-only production result. Private reproducibility files are `edit-delivery.py`, `reviewed-sections.json`, `render-helper/main.go`, `reviewed-artifacts.txt`, and the exported Typst bundle. No uploaded project code was executed. The editorial text contains no TODO markers; this does not turn missing measurements into known facts.

Raw and reviewed PDF pages were rendered to PNGs for inspection. The raw PDF exposes duplicate tables and gap boxes; the 17-page reviewed version removes these content defects and retains the academic cover, contents, numbering, source assignment and change sheet. The signature/date lines are intentionally unfilled for this fictional demonstration, not fabricated approvals.

## Checks and remaining work

- `go build ./...`, frontend typecheck/build and local Docker build passed for the relevant changes. Final schema changes passed Go build and real provider probes; no automated suites were added.
- Independent static review of the original SP1–SP4 diff and final Markdown-only delta reported no actionable regressions after the diagnostics-tag correction.
- Railway displayed the Markdown-only commit as Active with deployment successful; public `/readyz` returned `ready`. One brief failed-response page occurred during rollout and recovered on reload; no availability guarantee is inferred.
- Phone and laptop checks are simulations. They do not establish ease of use for actual students. The page remains long, reuse input is not yet visually unified, and generation details can retain a stale snapshot until refresh.
- Do not claim automatic submission readiness. The next quality work should add bounded related-definition retrieval, evidence-backed claim checking and cross-chapter contradiction checks, as described in [architecture handoff](../architecture-handoff.md). This wave does not implement that roadmap automatically.

Confidence is high for the observed structural/compatibility failures and fixes on this fixture, medium for the proposed UX benefit, and low for generalization to other projects. A different provider/model or a new fixture can overturn general quality expectations. The alternative explanation is that some gains come from stochastic model variation; this was not a controlled benchmark (A9, A10).
