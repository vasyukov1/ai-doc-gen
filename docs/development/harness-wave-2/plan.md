# Harness wave 2 implementation plan

Date: 2026-11-05. Reader: the owner choosing how to make document delivery easier than assembling prompts and exports manually. Scope: milestone 9 corrective UX, not milestone 8.

## Design and defaults

Working interpretation of the user's request, not a new human-owned directive: students, freelancers and employees should turn project evidence into editable documentation with clear remaining work. I1: preserve evidence and user text. I2: distinguish outline, model prose and verified content. I3: no uploaded code execution or application keys for simulation. I4: reversible, explicit reuse and existing owner isolation. I5: existing Go/PostgreSQL/Typst boundaries, no new dependency.

Alternatives: a cosmetic redesign leaves repeated setup and hidden outcomes; a project/account/credential vault would reduce more repetition but needs retention and ownership design; this bounded wave improves the current workflow first. The latter is the agent's reversible default under D10 and the user's authorization to improve. No deployment or broad new milestone.

Use writing-plans and executing-plans guidance; root instructions replace automated test steps with manual browser/HTTP checks. Execute inline, fresh code review at the end. Keep the unrelated tmp_example directory untouched. Attachments are reference prompts, not authorization to publish a prototype, rename the product or buy domains.

## Requirement/action/check index

| ID | Action and implementation | Manual acceptance |
|---|---|---|
| W2-R1 | Choose compact handover, FKN explanatory note or architecture via purpose shortcuts; general overview first; specification recovery and supported-language disclosure in index template | C1 phone/laptop shortcuts update profile guidance; explicit prior selection preserved; unsupported PDF explains conversion |
| W2-R2 | Choose a provider model from presets with custom-ID escape; key setup links/instructions; identical outline wording | C2 known/custom model submitted correctly; no provider chosen implicitly; no-key outcome clear |
| W2-R3 | Explicitly remember nonsecret creation text and model IDs on this browser; restore and forget; source files and keys excluded | C3 reload restores opt-in fields, key absent from storage, forget removes draft, denied storage doesn't break form |
| W2-R4 | Visible recorded input/output tokens and call count from owned generation report, with refresh and truthful missing-usage/cost labels | C4 recorded data displayed, failure recoverable, capped diagnostic list not misrepresented as full totals |
| W2-R5 | Direct intermediate sections and rebuild links; distinguish saved text from PDF version | C5 no-key generation for personas, inspect partial text, edit and rebuild; phone/laptop layout |

## Execution

- [x] Inspect current implementation and capture before phone/laptop screenshots.
- [x] Preserve independent persona source reviews and browser observations; triage maximum five accepted changes.
- [x] Implement R1–R3 in templates plus focused browser module.
- [x] Implement R4–R5 without changing generation policy.
- [x] Build Go/frontend, run isolated manual scenarios, inspect screenshots and PDF, obtain fresh code review.
- [x] Record limitations, deferred work, decisions and resume point.

Review focus: browser storage unavailable or corrupt; keys accidentally persisted; stale asynchronous usage response; custom model value; browser refresh during editing. No live-model quality or competitive superiority claim without an actual measured comparison.

## Resume

2026-11-05: five changes implemented and manual checks complete. Three no-provider generation paths and edit/rebuild passed; live prose quality, Catalan, project/credential reuse and cross-device recovery remain deferred. Fresh review found no important issues; wording about retained language/profile preferences was corrected. No tests were added. Next exact command: `go build ./...`. Preview uses isolated port 8083; shared Compose was not restarted.
