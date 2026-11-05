# Q3 — Catalan neutral documents

Date: 2026-11-06. Owner: localization agent. Requirement Q3 in [plan.md](plan.md). No automated suites, provider requests, service restarts or commits.

## Requirement/action/check index

| ID | Requirement and action | Manual check |
|---|---|---|
| Q3.1 | Add `ca` profile titles, every heading, notices and neutral render text for technical overview, SRS and arc42 | Render each neutral profile in Catalan, inspect JSON/PDF text and visible labels |
| Q3.2 | Use Catalan baseline gaps and bibliography labels; preserve source/user prose | Render no-provider sections and technology bibliography |
| Q3.3 | Reject unsupported academic Catalan in request and document-input validation; show scope before submission | Manual validation helper, frontend language/profile constraint |
| Q3.4 | Preserve Russian/English behavior and explicit saved language preference | Render overview in RU/EN, Go build and frontend type check |

## Decisions

- Catalan is `ca`. `Profile.SupportsLanguage` accepts it only for nonacademic profiles with a Catalan title, every heading and any required notice translated. This makes support explicit and lets the separately owned handover profile use the same contract.
- Browser chrome remains English/Russian; the language selector controls document output. Catalan profile guidance states the actual scope. Unsupported academic selections stay visible and receive a correction message instead of a silent language switch.
- Uploaded/source prose, user metadata and provider diagnostics are not automatically translated. Fixed draft notices receive native Catalan equivalents; unknown diagnostic detail is summarized without pretending it is translated.

## Verification

- `go build ./...`: exit 0, no output (repeated after the final Go changes).
- `npm --prefix web run typecheck`: exit 0, `tsc --noEmit`.
- Temporary manual helper `go run ./.q3-smoke`: exit 0. Rendered technical overview (4 leaves), SRS (22 leaves), arc42 (12 leaves) in CA/RU/EN through `documents.Render`; each produced PDF and source ZIP using pinned Typst 0.15.1. Both request and document-input validators rejected `ca` for all five academic IDs, and accepted the three neutral IDs. Inspected native CA profile/headings and fixed render/bibliography labels in output JSON. Helper retained outside the repository at `/tmp/ai-doc-gen-q3-smoke.go`; no test files or suites were added.
- Generated output root: `/var/folders/f3/dxghgy_523jb28fmh3znks180000gn/T/ai-doc-gen-q3-67967065`. Poppler text extraction confirmed native CA title metadata, annotations, labels and notices in all three CA PDFs. Visually inspected SRS annotation/body, arc42 contents and technical-overview contents PNGs: accents render, headings are legible, red TODOs remain visible without clipping. Existing baseline rendering repeats a gap in the body and warning box; this is unchanged across languages and deferred as a separate content/layout issue.
- Inspected prompt construction: existing chapter/section prompts already pass `ca` via the selected language. No external model call was made, so actual model-language quality is unverified.
- Added `00011_catalan_language.sql` after identifying the original PostgreSQL language check. Dedicated-database migration/application and browser phone/laptop checks belong to primary integration; not claimed by this isolated render check. Structured specification labels now have a CA branch, but its HTTP path has not been manually exercised here.

## Changed files and integration

- `internal/templates/profiles.go`, the three neutral JSON profiles and `ARC42-NOTICE.md`: CA metadata and completeness contract; primary later adds client-handover through the same contract.
- `internal/documents/request.go`, `render.go`, `bibliography.go`, `draft_sections.go` (`BaselineSection` branch only), `sections.go` (legacy uncited warning branch only): language validation, native fixed text, and export labels.
- `internal/web/documents.go`: native headings for optional structured specification text.
- `web/templates/index.html`, `job.html`, `web/src/profile.ts`, `forms.ts`: selectable scoped document language, visible correctable unsupported pairing and explicit saved CA preference. Browser chrome is not presented as translated. Primary owns compiled static assets and handover additions.
- `migrations/00011_catalan_language.sql`, migration/template README files and this report: persistence constraint, rollback limits and support contract.

Resume point (2026-11-06): Q3 implementation and isolated rendering complete; primary integrates persistence/browser checks and compiles static assets. Next command: `go run ./cmd/ai-doc-gen migrate` with the dedicated integration `DATABASE_URL` already set. Apply before creating a CA job. No worker/service restart was performed by this agent.

## Integrated real-model handover

The primary generated all six client-handover sections in Catalan using the normal saved-project and encrypted-connection reuse flow. Source/language review corrected a meaningful typo, command order, prerequisites and credential wording; the ordinary editor and worker rebuilt the final 17-page PDF. This proves persistence/rendering and a substantive Catalan example, not native-speaker certification. Original supplied specification/metadata retain their original language. Raw and reviewed counts are in the [main report](README.md).
