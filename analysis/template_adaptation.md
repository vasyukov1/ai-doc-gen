# Academic template adaptation for the website

Status: source-template inventory completed; migration is planned, not implemented. The product owner chose analysis only for the current stage (2026-11-02). The existing source templates remain untouched. This analysis describes the intended destination behavior without depending on the source repository at runtime.

## Existing document variants

| Internal code | User-facing name (English) | Current source structure | Website behavior |
|---|---|---|---|
| `tz` | Technical specification | 10 ordered sections | Accept an uploaded specification or structured web-form answers; do not draft intended requirements from code alone |
| `pz` | Explanatory report | 5 ordered sections; purpose, technical characteristics, economics, references | Generate supported sections and retain required fixed sections |
| `pmi` | Test programme and method | 7 ordered sections | Use requirements and test evidence; mark missing test results clearly |
| `pmi-team` | Team test programme and method | 7 ordered sections | Add team-specific source material when supplied |
| `ro` | Operator guide | 4 ordered sections | Use actual setup and usage paths; request missing screenshots or UI details |

The existing Typst layout shares a rendering core, then assembles each document from ordered section files. Some sections are fixed template text; others have model prompts. The website should present human-readable document choices while keeping stable internal IDs and an explicit `profile_version`. Users should be able to choose one document or a bundle where prerequisites are satisfied.

## Changes needed for web generation

1. **Input model:** replace local project paths and CLI flags with a job record containing repository commit or archive digest, chosen document IDs, author/team metadata, optional description, an uploaded specification or structured specification form, language, and template version. Ask for missing title-page fields in the form.
2. **Source rules:** avoid assuming that a technical specification exists inside every repository. Show which sections depend on it, and request missing requirements rather than inferring them from implementation code.
3. **Profile packaging:** copy and adapt the Typst assets into this repository as versioned assets during implementation. Do not require the original CLI or its directory layout at runtime. Check inherited text, fonts, images, and any template reuse rights before publishing.
4. **Language:** provide complete `ru` and `en` variants for fixed text, headings, prompts, title pages, references, and fallback messages. Translate fixed text for English while preserving the existing academic title-page structure. Review the English wording with a human. Never mix languages because a fixed section was left untranslated.
5. **Rendering:** populate a structured data file, render sections without allowing uploaded code to become executable Typst, compile in an isolated worker, and return both the source bundle and PDF.
6. **Review:** expose section drafts and missing-evidence warnings before final export. Retain template headings and order. Regenerate one section without silently replacing user edits in other sections.

## Acceptance checks

- Every chosen variant compiles in both languages with sample metadata and an empty-code fixture.
- Required headings and order match the versioned profile. All user metadata appears in the correct title-page fields.
- Russian output has no unintended English fixed text; English output has no unintended Russian fixed text. Project names and quotations may remain in their original language.
- The downloaded Typst bundle compiles to the same PDF using the pinned compiler version and included assets.
- Missing specification, test evidence, economic figures, or screenshots appear as explicit user actions or placeholders, never fabricated statements.
