# Harness wave 1 — evidence and handoff

As of 2026-11-05. Reader: the project owner deciding whether document creation and correction are safer. This is a milestone 9 correction; it does not measure LLM prose quality or real-user conversion.

## Findings and decisions

[Plan and stable IDs](plan.md) · [triage](triage.md) · [raw simulated responses](responses/) · [findings index](findings.json).

Three fictional personas each had two independent implementation reviews before changes and two independent targeted reviews afterwards (12 fresh-context Codex subagents, inherited session model). No application-provider API calls, keys or paid simulation-script calls were made. Prompts and answers are retained. Before-review inputs were frozen source snapshots, with optional screenshots; reviewers explicitly reported implementation-only inputs. After reviews used actual synthetic-fixture screenshots, code and manual logs. Do not interpret reviewer agreement as a percentage of real users.

Recompute before-review agreement without API access:

```sh
python3 docs/development/harness-wave-1/summarize.py
```

The normalized (screen, heuristic) intersection/union scores are recorded in [agreement.json](agreement.json), including all pairs behind each value. The primary agent's normalization collapses page variants into creation/sections/status; this may hide differences between detailed findings. The two after-review rounds straddle small feedback fixes (success styling, stale labels, table preview); they are targeted rechecks, not a controlled before/after experiment.

## Manual acceptance evidence

These are developer-run synthetic browser/HTTP/SQL checks, not real-person research. Native Go web/worker processes used the isolated `harness_wave1` database and port 8082, leaving existing user jobs untouched. Uploaded example Go code was analyzed but never executed.

| Check | Persona / device / action | Observed result | Evidence |
|---|---|---|---|
| C1 | Beginner, phone: PDF specification, two spec sources, correct the inputs | PASS: actionable validation, file removal, title-field focus; UTF-8 byte limit enforced | [walk](checks/walk.log), [recovery](checks/recovery.log), [phone](screens/after-phone.png) |
| C2 | Maintainer, laptop and phone: edit two sections, save one, leave the page | PASS: other text survives; native navigation warning; table JSON persists in PostgreSQL | [edits](checks/edits.log), [phone check](checks/edits-phone.log), [laptop](screens/after-edits-laptop.png), [phone](screens/after-edits-phone.png) |
| C2 | Maintainer: HTTP conflict, network failure, retry, type while save is pending | PASS: text retained; retry succeeds; later local revision stays dirty; table preview readable | [recovery](checks/recovery.log), [screen](screens/after-recovery-laptop.png) |
| C2 | Maintainer: server never returns save response | PASS: 30-second client deadline releases controls and retains text; browser clock advanced for this scenario | [timeout](checks/timeout.log) |
| C3 | Reviewer: complete a no-key job | PASS: PDF compiled; four baseline sections, zero model sections, four gaps; outline disclosure and actionable links visible | [walk](checks/walk.log), [phone](screens/after-status-phone.png), [laptop](screens/after-status-laptop.png) |
| C4 | Reviewer, phone: status API fails, then recovers | PASS: stale-state warning and manual refresh; warning clears and recovery/model counts update | [poll](checks/poll.log), [screen](screens/after-poll-phone.png) |
| C5 | Beginner/reviewer: choose arc42 in Russian; inspect processing disclosure | PASS: matching guidance, server/provider disclosure; 390px viewport has 390px document width | [walk](checks/walk.log), [phone](screens/after-phone.png) |

Commands used for the manual checks: `node tmp/harness-wave-1/walk.cjs`, `node tmp/harness-wave-1/edits.cjs`, `node tmp/harness-wave-1/recovery.cjs`, `node tmp/harness-wave-1/timeout.cjs`, `node tmp/harness-wave-1/poll.cjs`. These one-off local harnesses and private synthetic owner cookies stay under ignored `tmp`, not in an automated test suite or the published evidence. Logged successful output is retained here. Repeating the browser checks requires provisioning a new isolated database, fixture and browser session; no replay claims are made for expired cookie state.

Build checks: `go build ./...`, `npm --prefix web run typecheck`, `npm --prefix web run build`, `git diff --check` (all exit 0). Initial manual checks caught CSS overriding the `hidden` attribute and named `action` controls shadowing `form.action`; these were corrected rather than weakening the checks. Snapshot `.go` files were renamed `.go.txt` so manual-review inputs are not accidentally compiled as application packages.

Final code review caught a redirect race: edits made while regeneration was pending could be cleared before navigation. The corrected revision check retains the page and supplies an explicit guarded link. [The delayed-response smoke](checks/redirect-race.log) confirms the newer text remains. A fresh reviewer rechecked the fix; no further P1/P2 findings were reported. The isolated staged tree also passed `go build ./...`; its generated assets exclude pre-existing uncommitted PDF-viewer changes.

## Limits and contrary view

Better status text and safeguards do not make generated prose more accurate. Full Russian UI localization, cross-session recovery and regeneration version history remain deferred. English status warnings can still be missed by a beginner with limited English. Regeneration intentionally replaces the saved section after confirmation; it is not reversible. Continuation with a real provider and tablet/Safari navigation behavior were not validated. Confidence is high in the executed paths and medium in general usability improvement; a real user unable to recover from these mistakes would overturn the latter claim.

Only claims backed by the check logs may be stated as observed. Simulated concerns require attribution. Do not claim verified PDF content, ISO certification, a real-user improvement percentage or a production latency improvement from this wave.

Local Compose was refreshed after confirming no active user jobs. `/readyz` returned `ready`; the updated homepage contains the profile guide and server-processing disclosure. The isolated database recorded zero provider calls. Existing unrelated PDF-viewer changes remain in the working tree, outside this wave’s commit.
