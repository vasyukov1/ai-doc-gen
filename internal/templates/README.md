# Versioned document profiles

`v1/` embeds five academic variants (`tz`, `pz`, `ro`, `pmi`, `pmi-team`) and the separate `technical-overview` profile. JSON assets preserve the academic heading hierarchy and section order and contain Russian and English labels. Leaf headings keep independent edit IDs within logical chapter requests; structural parent headings are rendered without invented paragraphs. Approval sheet, title sheet, annotation, contents, body and change registration remain separate academic blocks. Title fields and specification content are supplied by the user. Fixed annotations describe the draft without asserting standards compliance.

`main.typ` is trusted code. All project names, metadata, specifications, draft bodies and warnings enter through `document.json` and `text()`. Never evaluate user strings or compile uploaded Typst. Compilation uses a private root, a clean environment, isolated package directories, a 45-second deadline and exactly Typst 0.15.1. System and embedded fonts are disabled. Liberation Serif font files and their SIL Open Font License are bundled in `fonts/`. The Docker image verifies official compiler archive SHA-256 values for amd64/arm64 before installation.

The editable ZIP contains `main.typ`, `document.json`, profile assets, fonts, their license and reproduction instructions. Worker exports also retain `evidence.json` so citation identifiers can be resolved after download. A fixed creation timestamp makes unchanged inputs reproducible with the pinned compiler. The application compiles its own template only; users editing downloaded trusted source are responsible for the code they subsequently compile locally.

Official references checked during implementation:

- [JSON data loading](https://typst.app/docs/reference/data-loading/json/)
- [Literal text](https://typst.app/docs/reference/text/text/)
- [Compiler CLI and root](https://github.com/typst/typst/tree/v0.15.1/crates/typst-cli)
- [Pinned compiler release](https://github.com/typst/typst/releases/tag/v0.15.1)

## FKN layout version 2

Academic profiles use the embedded `fkn.typ` adaptation of the user-supplied MIT-licensed FKN layout (copyright Alexander Vasyukov; see `v1/FKN-LICENSE`). The layout preserves the two front sheets, vertical inventory cells, approval-sheet counter reset, main-page footer, heading hierarchy and ten-column change register. Localized labels support RU/EN. `main.typ` renders all data literally and highlights TODO text in dark red. The generic overview remains independent.

`fkn-2` identifies the layout revision. The bundled Liberation Serif family is used instead of the reference's system-dependent Times New Roman. Annotations are generated from the actual profile headings; copying a template does not assert verified standards compliance. User fields and project-specific facts are never copied from the reference. The downloaded source bundle includes the adapted layout and license; no external runtime import is needed.
