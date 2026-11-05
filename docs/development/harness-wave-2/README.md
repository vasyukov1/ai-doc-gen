# Harness wave 2 — document delivery journeys

Cut-off: 2026-11-05. Reader: the product owner deciding what makes this workflow easier than manual interaction with a general-purpose coding assistant. These are simulated hypotheses, not interviews or measured conversion improvements. [Requirements and checks](plan.md).

## Method and evidence

Three fictional personas received independent source reviews twice, with fresh contexts and no other reviews visible: Mira, a Russian-speaking FKN student on a phone; Nil, a Catalan-speaking freelancer on a laptop; Lena, an employee architect using laptop and phone. Each had a further fresh-context rewalk against changed source, screenshots and available manual logs. Nine persona reviews used the inherited Codex session model; no application provider keys or provider calls were used for simulation. Prompts and exact responses are in `prompts/` and `responses/`. A separate final code reviewer found no important issues; its minor wording concern about independently retained profile/language preferences was corrected.

Recompute normalized screen/heuristic overlap with `python3 docs/development/harness-wave-2/summarize.py`; [agreement.json](agreement.json) retains the inputs, union and intersection. Student and freelancer reviews have identical normalized pairs; employee overlap is 3/7. This reflects heuristic-review stability, not truth or a fraction of real users. Normalization is the analyst’s and collapses details.

Before reviews are source-based heuristic inspections, not independent browser interactions. The primary agent operated Chrome for the manual scenarios. The three after reviews are one targeted pass each, not two independent repeated experiments. Model identity/version was inherited, not independently recorded. These limits prevent claims of fully reproducible model behavior.

The student scenario used the public sample repository supplied in the request as an uploaded ZIP, never a runtime dependency and never executed. Other end-to-end scenarios used the same small synthetic order-service archive retained from wave 1, with technical-overview and arc42 profiles. The freelancer's booking-API persona is hypothetical; runtime evidence concerns the order-service fixture, not a separate booking product. Metadata and requirements were invented and labeled synthetic. No real company code was uploaded.

Isolated local database `harness_wave2`, separate data directory under ignored `tmp/harness-wave-2`, web on port 8083, native worker. Existing jobs and services were not modified. Existing `tmp_example/` remains untouched. Root milestone instructions override generic automated-suite requirements.

## Five accepted changes and deferred findings

| Change | Findings addressed | Decision and invariant |
|---|---|---|
| R1: purpose shortcuts, general handover default, PDF-text recovery, language and privacy disclosure | M1, M2-1, NIL-1, NIL-3, freelancer round-2 rows 1/3, L1, L2-1 | ACCEPT: I2/I3, correct template choice and informed inputs |
| R2: provider model presets/custom ID, key links/instructions, consistent no-model wording | M2, M2-2, provider onboarding in owner's request | ACCEPT: I2, recognize configuration without memorizing IDs |
| R3: opt-in remembered text/source address/model IDs, restore and forget | M3, M2-3, NIL-2, freelancer round-2 row 2 | ACCEPT: I1/I4, reduce repeated typing without storing keys/files in browser |
| R4: visible recorded tokens and diagnostic-call count with explicit limitations | L2, L2-2 | ACCEPT: I2, consumption observable without fabricated cost |
| R5: available-section preview/rebuild links and honest editing/availability labels | M4, M2-4, NIL-4, freelancer round-2 row 4, L3/L4/L5, L2-3/L2-4 | ACCEPT: I1/I2, preserve text and expose next action; full-page refresh and section navigator remain deferred |
| Cross-device ownership transfer | L2-5 | DEFER: I4; authenticated transfer/recovery needs explicit ownership and retention design |

Remaining parts of accepted findings are retained here, not treated as solved: no automatic creation-page leave warning; remembrance is opt-in and low in the form; new documents still need credentials and file reselection; no shared project workspace; intermediate sections still use full-page refresh; long documents lack a section navigator. Existing document continuation/rebuild reuses its stored source and connection. Catalan remains unsupported, now disclosed.

## Manual checks

- `go build ./...`, `npm --prefix web run typecheck`, `npm --prefix web run build`, `git diff --check`: recorded separately in `checks/build.log`.
- `node tmp/harness-wave-2/before.cjs`: baseline Chrome screenshots at 1366×900 and 390×844; content width matched viewport, model controls were text inputs, languages were EN/RU.
- `node tmp/harness-wave-2/walk.cjs`: [log](checks/walk.log). Purpose shortcuts, RU guidance, custom model restoration, key exclusion, forgetting, no horizontal overflow. Three no-provider jobs produced compiled PDFs with outline outcomes. Observed section counts: student 17, freelancer 4, employee 12, all model counts zero. These counts describe fixtures, not quality or user statistics.
- `node tmp/harness-wave-2/recovery.cjs`: [log](checks/recovery.log). Controlled synthetic usage response, failing response/retry, saved edit/direct rebuild/preserved prose, disabled browser storage, corrupt saved draft recovery. Synthetic token values are UI fixtures, not measured model use.
- `node tmp/harness-wave-2/screens.cjs`: [log](checks/screens.log), all three status screens at phone and laptop widths without horizontal overflow; PDF-input error corrected by removing the file and pasting text.
- PDF visual inspection used `pdftoppm`: technical-overview title page, FKN contents page, arc42 contents page. Titles, TODO styling and layout were readable. This is a small layout check, not exhaustive PDF or semantic validation.

One-off manual scripts, private owner cookies and generated sample PDFs remain in ignored `tmp`; screenshots/logs are retained here. Scripts initially needed two harness corrections: open the collapsed model group before selecting, and expect the existing compilation redirect to `/document`. Neither was a product failure; final evidence comes from corrected runs.

## What should differentiate the product

Analysis, medium confidence: the defensible value is a persistent document workspace: explicit structure, traceable sources, editable sections, missing-information markers, recoverable jobs and repeatable PDF/Typst delivery. Presets alone are not a moat. Strong contrary case: an experienced user with an agent, repository access and a known template may produce a better document faster. This wave does not establish superiority, better prose, time saved or standards compliance.

Next priorities, deferred under the five-change cap:

1. **Quality evaluation:** same consented projects/specifications/model budgets through this product and a generic agent; blind section review for correctness, source support, completeness and editing effort. Record token/time/cost measurements and failures. No key should be spent on persona simulation.
2. **Reusable project workspace:** source snapshot + requirements + metadata, multiple document types/versions, explicit encrypted connection reuse, expiry/deletion and owner isolation. Avoid silent source drift or exposing old keys to the browser.
3. **Catalan and other languages:** translate profile structure, prompts, baseline text, bibliography/labels, metadata, PDF fonts and validation; bilingual review against actual examples. A language dropdown alone is insufficient.
4. **Task templates:** client handover with deployment/operation/maintenance, team onboarding, API integration guide; keep observed facts, proposed acceptance criteria and missing stakeholder decisions separate. Improve existing profiles before multiplying near-duplicates.
5. **Preview and review:** chapter navigation, incremental reading without page reload, saved-versus-rendered versions, edit history, live measured usage, clear connection repair.
6. **Optional external research/SERP:** useful for external bibliography and comparison, not evidence of repository behavior. Needs explicit external sources, citations, retrieval dates and a bounded opt-in budget. Do not add another key before validating this use case.

Confidence is high only in executed paths, medium in the proposed usability benefit, and unmeasured for prose quality. Real users still failing to recover, or an equal-budget agent needing less correction, would overturn the benefit hypothesis.

## Sources and claim boundaries

Provider setup links were checked against primary instructions: [Groq quickstart](https://console.groq.com/docs/quickstart), [OpenRouter quickstart](https://openrouter.ai/docs/quickstart), [Gemini API keys](https://ai.google.dev/gemini-api/docs/api-key). Gemini's Markdown response was fetched with curl after the browsing tool rejected its content type. Existing configured model IDs were retained as presets; they are not a new model ranking or guarantee of current availability. Worker catalog checks remain authoritative.

Observed claims may cite the manual logs; simulated concerns must be labeled hypotheses. Never claim real-user percentages, measured competitive advantage, live-model quality, monetary savings, or Catalan support from this wave. No external publication or production deployment was performed.
