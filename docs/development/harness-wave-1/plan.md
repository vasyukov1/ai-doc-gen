# Harness wave 1 — recoverable document journeys

Date: 2026-11-05. Scope: milestone 9 correction, not the deferred automated testing milestone.
Rules applied: D1–D4, D6, D8–D10, D12–D15; U3–U13; A6, A10, A13.

## Fixed scope and defaults

This is an agent's working interpretation of the owner's existing requests, not a replacement for a human-owned product directive. The imported directive template and demo are untouched. The reader is the project owner deciding whether the small website makes document creation and correction reliable. Synthetic critiques are hypotheses, not user measurements.

- I1: A usable baseline and PDF remain possible without model credentials; no new provider calls or billing features.
- I2: User edits, private ownership and saved evidence remain under the user's control; no silent data loss.
- I3: Compiled PDF does not imply verified content; gaps and baseline origin remain visible.
- I4: Mistakes must be correctable on the current page, including on a phone.
- I5: Existing Go/PostgreSQL/Typst boundaries and previously uncommitted PDF viewer work are preserved.

D10 defaults: prioritize observed journey failures, not a visual rewrite; use existing Codex fresh-context reviewers instead of the paid API simulation script; retain manual smoke checks because root AGENTS explicitly defers automated suites. Do not publish or restart shared services as part of evaluation; run a separate local preview. Do not rewrite the supplied teaching examples or claim they describe this product.

## Requirement / situation / task / manual check index

| Requirement | Situation | Task (estimated 20–30 minutes each) | Check and expected result |
|---|---|---|---|
| R1: correctable input | US1: choose ZIP/specification incorrectly | T1: validate limits, exclusive sources and readable field labels | C1: wrong/oversized file and text+file blocked; correcting clears error; field preserved |
| R2: safe editing | US2: edit one of several sections, then navigate | T2: retain tables, inline saves, unsaved-change guard | C2: server retains table JSON; other dirty forms prevent accidental redirect; failed save retains text |
| R3: honest result | US3: download an outline believing it is final | T3: explicit available/model counts and actionable review/continue links | C3: outline/partial/ready text differs; no impossible retry advice |
| R4: visible connection state | US4: lose network while generation continues | T4: show polling failure and manual retry, bound request time | C4: failing API shows stale state; recovery clears it; no overlap |
| R5: informed submission | US5: choose SRS vs FKN; assume files remain local | T5: profile guidance and accurate privacy/key-free expectations | C5: guide updates with profile; phone layout works; no unsupported guarantee |

## Execution

- [x] Baseline: fresh persona critiques twice each; capture current phone/laptop screens; triage at most five accepted changes.
- [x] T1 + T5: creation guidance and validation.
- [x] T2: saved tables and recoverable edit flow.
- [x] T3 + T4: honest result and connection feedback.
- [x] Walk again, inspect screenshots, build Go and frontend; review own diff and preserve unrelated changes.

## Coordination

Claim: primary agent owns forms.ts, job.ts, section save, index/section/job templates and this report. Reviewer agents only read source/screens and write uniquely named critique responses. Existing dirty files are snapshotted under ignored tmp/harness-wave-1/before; no database records or shared services are changed by this wave's checks.

## Resume

Status: implementation and manual checks complete; verification evidence is in README.md and checks/. Full localization and regeneration history are deferred. Next command: `go build ./...`.

Verification adjustment: isolated native services exercised the wave before local Compose refresh. Existing PDF-viewer work remains unstaged; shared-file deltas and generated assets are staged separately against HEAD. Provided rule documents are imported unchanged as task references, not installed over the root product instructions.
