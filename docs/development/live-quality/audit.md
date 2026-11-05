# Source audit and acceptance criteria

Date: 2026-11-06. The fixture is the public Go CLI supplied by the user; paths below belong to that uploaded source, not this application. Inspection was static. No uploaded program, tests or compiler instructions were executed.

| ID | Check against supplied source | Observed defect in an earlier raw model result | Required treatment |
|---|---|---|---|
| Q1-C1 | `internal/app/app.go`, `generateDocDrafts` propagates `client.Complete` errors | Invented graceful TODO fallback | Explain propagation; initialized files do not establish fallback behavior |
| Q1-C2 | Draft mode writes unchanged section filenames under `docs/<doc>/drafts/` | Invented `.draft.typ` suffix | Use the actual branch/path or omit an unsupported exact path |
| Q1-C3 | `Init` creates input/configuration; `GenerateDoc` calls template extraction | Template extraction attributed to initialization | Separate initialization from generation |
| Q1-C4 | CLI init supports `--output`; improve supports `--file`, `--prompt`, `--apply` | Mandatory current project directory; unexplained improve TODO | Cover supported flags without inventing restrictions |
| Q1-C5 | `CreatePDF` invokes Typst compilation; generation calls the model | Claimed AI was future work and application did not compile PDF | Describe implemented main flow; optional bundling is separate |
| Q1-C6 | Go toolchain is a build dependency; no measured RAM/disk minima supplied | Invented 2 GB/1 GB prerequisites | Distinguish build/runtime; avoid invented numeric thresholds |
| Q1-C7 | CI source declares tests, build and CLI/Homebrew checks | Verification omitted CI or implied tests passed | Describe declared steps without claiming execution |
| Q1-C8 | User description is not an approved development specification | Public examples treated as formal authorization | Explicitly identify the absent formal basis |
| Q1-C9 | No measured time savings, demand or competitor benchmark supplied | Claimed savings, superiority and support speed | Label intended benefits/hypotheses; propose an evaluation method |
| Q1-C10 | Exported prose and tables need complete content, resolved citation metadata and honest gaps | Visible TODO with empty missing-information metadata | Count visible TODO conservatively; review actual text and output, not provenance counts alone |

## Engineering changes informed by these observations

Output budgets and bounded chapter groups address truncation. Typed source weighting, capped keyword contributions, test-versus-production distinction, section-specific priorities and bounded shared context reduce misleading or missing implementation context. Profile guidance distinguishes implementation from future plans and inferred rationale from documented history. Source citation normalization and validation reject invalid claims of provenance; they cannot establish the semantic truth of every sentence.

The latest retrieval revision followed a manual pack inspection: uncapped BM25 relevance outweighed all kind bonuses, and names of test functions accidentally received production-generation bonuses. Keyword contribution is now bounded, tests receive a low production priority, and verification explicitly boosts test/CI material. The pack remains bounded and can omit useful facts; this is a reason for source review, not a guarantee of completeness.

A visible `TODO:` line is counted conservatively even when a model forgets its gap array. A discussion quoting a template's literal TODO can therefore require editorial review; this conservative false-positive is preferable to reporting an unresolved instruction as gap-free. Genuine author/approval/measurement gaps are not filled by invention.

Final delivery must distinguish raw model output from any source-reviewed editorial corrections. Corrections should use the ordinary section-save and PDF rebuild workflow, leaving edited provenance visible. No automated claim-verification or agent-superiority benchmark is asserted.
