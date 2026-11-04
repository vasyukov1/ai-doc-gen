# Versioned document profiles

`v1/` embeds five academic variants (`tz`, `pz`, `ro`, `pmi`, `pmi-team`) and the separate `technical-overview` profile. JSON assets preserve the academic heading hierarchy and section order and contain Russian and English labels. Leaf headings become independent grounded drafting requests; structural parent headings are rendered without invented paragraphs. Approval sheet, title sheet, annotation, contents, body and change registration remain separate academic blocks. Title fields and specification content are supplied by the user. Fixed annotations describe the draft without asserting standards compliance.

`main.typ` is trusted code. All project names, metadata, specifications, draft bodies and warnings enter through `document.json` and `text()`. Never evaluate user strings or compile uploaded Typst. Compilation uses a private root, a clean environment, isolated package directories, a 45-second deadline and exactly Typst 0.15.1. System and embedded fonts are disabled. Liberation Serif font files and their SIL Open Font License are bundled in `fonts/`. The Docker image verifies official compiler archive SHA-256 values for amd64/arm64 before installation.

The editable ZIP contains `main.typ`, `document.json`, profile assets, fonts, their license and reproduction instructions. Worker exports also retain `evidence.json` so citation identifiers can be resolved after download. A fixed creation timestamp makes unchanged inputs reproducible with the pinned compiler. The application compiles its own template only; users editing downloaded trusted source are responsible for the code they subsequently compile locally.

Official references checked during implementation:

- [JSON data loading](https://typst.app/docs/reference/data-loading/json/)
- [Literal text](https://typst.app/docs/reference/text/text/)
- [Compiler CLI and root](https://github.com/typst/typst/tree/v0.15.1/crates/typst-cli)
- [Pinned compiler release](https://github.com/typst/typst/releases/tag/v0.15.1)
