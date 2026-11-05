# Live quality and reusable delivery plan

Date: 2026-11-05. Authorized follow-up to harness wave 2, within milestone 9. No automated suites or uploaded-code execution. Preserve existing jobs, private ownership and tmp_example. English repository documentation.

## Design

The user explicitly requests real-model quality on the supplied public repository and the remaining UX tasks. Prefer measured source/prompt repairs over adding another provider or disguising gaps. Use existing configured free connections for actual document generation only; never for persona simulation. Preserve a pinned input and baseline before changes. Unknown approvals, measurements and personal data must not be invented to remove TODOs.

Five accepted changes (D1–D4/U11):

| ID | Deliverable | Owner/files | Manual acceptance |
|---|---|---|---|
| Q1 | Measure real baseline, fix evidenced generation/context defects, rerun and inspect complete prose/PDF | primary, internal/documents generation/context files | Same source and model before/after, record time/calls/tokens, inspect claim support and gaps |
| Q2 | Create another document from an owned saved project without reupload/retyping, explicitly reuse encrypted connection | reuse agent, internal/storage reuse, internal/web reuse routes, new reuse page | Same owner succeeds, foreign/expired job denied, source independence after parent deletion, no keys returned |
| Q3 | Catalan for neutral technical profiles, complete headings/render labels/baseline and validation | localization agent, templates, request/render/bibliography, language UI | CA overview/arc42/SRS render, unsupported academic CA rejected honestly, RU/EN preserved |
| Q4 | Dedicated client handover profile with operation/deployment/maintenance guidance | primary after Q3, new template | Useful project-specific deliverable, no fabricated production promises |
| Q5 | Section navigation, stable intermediate reading and saved-vs-rendered guidance | preview agent, sections template and focused frontend module | Read running sections without forced reload, jump to chapter, preserve unsaved editing, phone/laptop |

Deferred to a later wave: cross-device identity, arbitrary languages, version-history storage, optional SERP and comparative agent benchmark. These are separate capabilities, not silently accepted extras. Reuse is explicit per owned job, not an account-wide key vault. Catalan support is limited to neutral documents; FKN translation is not implied.

## Execution ledger

- [x] Real-provider baseline and source inspection.
- [x] Q1 evidenced fixes and live rerun; raw versus source-reviewed results are separated in README.md.
- [x] Q2 reusable source/connection flow.
- [x] Q3 Catalan neutral profiles.
- [x] Q4 client handover profile; real Catalan delivery and source review completed.
- [x] Q5 reading/review flow.
- [x] Manual integration, independent code/quality reviews, documentation and build. See README.md and checks/.

Use subagent-driven-development for isolated Q2/Q3/Q5 tasks; primary owns Q1 and integration. No shared-file changes without coordination: localization owns index.html and job.html; preview owns sections.html; reuse adds reuse.html and supplies insertion instructions to primary. Primary integrates main.ts imports and builds static assets. No worker/service restarts by subagents. Dedicated database/port for live checks. Current feature branch provides task isolation; no move of user files.
