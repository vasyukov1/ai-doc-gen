# Guided document creation

Date: 2026-11-06. User request: choose the document first, explain inputs next to the controls, hide irrelevant fields, and remember information with an option to fill it again.

## Design

Keep a single adaptive form, rather than a multi-page wizard that adds navigation and could lose files. Start with document type and output language. Use the actual Typst metadata contract: neutral profiles need project, author, organization and year; academic fields appear only for academic profiles. Hidden metadata controls are disabled for submission but retain their values in memory and explicitly opted-in browser storage. Specification is required for academic profiles and optional source requirements/context for neutral profiles, with profile-specific guidance. Optional is not the same as irrelevant.

Keep explicit local-storage opt-in and never persist keys or uploaded files. Move the existing save controls near the beginning, add explicit Restore saved form and Start over actions. Start over clears current text, credentials and files while retaining the saved snapshot for recovery; disable autosave until the user opts in again. Forget removes the saved snapshot. No account storage or Google login is implied.

## Accepted changes and manual acceptance (D1–D4)

| ID | Change | Check |
| --- | --- | --- |
| GF1 | Type/language first; profile-specific visibility | Technical overview hides academic fields and excludes them from FormData; switching back retains typed values; academic specification remains required |
| GF2 | Inline instructions/examples and upfront title-page requirements | Each visible metadata field has a linked hint; neutral and academic forms explain what to enter; no newly mandatory backend fields |
| GF3 | Visible save, restore, start-over and forget | Opt-in persists text across reload; restore and reset work; keys/files never stored; storage failure stays actionable |
| GF4 | Stable repository-derived project title | Complete repository URL yields complete project name; explicit user title is not overwritten |

## Execution

- [x] Inspect existing frontend, persistence and renderer contracts.
- [x] Implement adaptive form and reuse controls.
- [x] Build Go/frontend and manually walk laptop/mobile flows with fictional input, without model keys.
- [x] Independent code review, record limitations, commit and push.
- [ ] Integrate the verified change into the deployed branch and verify the published form.

This increment does not start milestone 8 or alter generation semantics. Preserve unrelated tmp_example work. Findings from walkthroughs are simulated hypotheses, not real-user measurements.

## Deferred

Per-project libraries of metadata, account sync and Google login remain outside this four-change wave. Specifications are not fabricated from examples, and hidden controls cannot create undisclosed requirements.

## Verification — 2026-11-06

Frontend typecheck/build, Go build, Docker image build and diff whitespace checks passed. Independent static review found no important issues. Browser smoke on isolated localhost:8086 confirmed hidden academic fields are disabled, values survive profile changes, academic requirements open upfront, progressive URL typing produces the complete title, and explicit titles survive source edits. Opt-in save survived reload; Start over cleared current fields without deleting the snapshot; Restore recovered author/title/group; Forget followed by reload left text blank. At 390×844 the document scroll width equalled the viewport width, with readable controls. No provider keys or generation calls were used.

Credential/file exclusion was source-reviewed; file reattachment and browser-storage failure were not exercised interactively in this increment. This is a simulated walkthrough, not measured user research.
