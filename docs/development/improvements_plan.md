# Usability and generation improvements

Status: approved on 2026-11-05 after the user answered the open questions. Phases A, B and E are implemented on 2026-11-05; Phases C and D remain. Work proceeds in phases; each phase ends with `go build ./...`, the template check and a manual smoke note in `verification.md`. Automated tests stay deferred to the final milestone.

## User decisions

- Local model spike uses `qwen2.5:7b-instruct` in Ollama.
- Form simplification is soft: metadata and model fields move out of the way, but the technical specification stays required for academic profiles.
- Extraction stays deterministic: Go AST plus new language parsers and BM25-style ranking. No tree-sitter, no embedding service, no runtime dependency on Serena or CodeGraph.
- Signature and signing-date placeholders are paper fields, not missing data. They must render as blank lines, never as red `TODO`.

## Phase A — frictionless submission

### A1. Mutually exclusive sources

- Choosing a ZIP clears the GitHub field; typing a GitHub URL clears the chosen ZIP. The dropzone also exposes an explicit clear action.
- Fixes the reported dead end where a chosen archive blocked submission after a repository URL was entered.

### A2. Archive preflight and tolerant intake

- The upload request inspects the ZIP central directory before creating a job and returns the exact reason and path.
- Symlinks, special files and duplicate names are skipped with an inventory reason instead of rejecting the whole archive. Only unsafe paths (traversal, absolute, control characters) remain hard failures, and their message names the offending entry.
- The worker uses the same rules for GitHub archives.

### A3. Fewer required fields

- The user chooses the document language; only their own earlier choice is remembered in `localStorage`. Nothing is inferred from the browser.
- Year defaults to the current year. The project title is derived from the repository name or archive name when empty.
- The technical specification remains required for academic profiles, as decided.
- Metadata, specification and provider keys stay in collapsed advanced blocks.

### A4. Paper signature fields

- `signature_todo` and `date_todo` become plain blanks in both languages; the red TODO style applies only to missing data.

## Phase B — generation reaches the end

- Budget per profile: `12 + 2 × chapter count` admissions and `240s + 45s × chapter count` model time, configurable by environment.
- Coverage first: one call per chapter before any repair call, so repairs cannot starve later chapters.
- Automatic continuation: up to two extra bounded passes after the deadline, preserving accepted sections; the manual button remains.
- Generation details report covered chapters, disabled providers with reasons and remaining budget.

## Phase C — extraction quality

- Deterministic project graph: service, route, entity, dependency and call sites with source spans. Go AST stays; TypeScript/JavaScript and Python get import/export/function/route extraction, and YAML OpenAPI joins the JSON parser.
- BM25-style ranking and score fusion replace substring counting in section and chapter packs, with deduplication and service balancing.
- AST-based chunking: functions, types and endpoints are the unit; long bodies keep their signatures.

## Phase D — local model

- Ollama service in Compose with a persisted model volume and `make models`.
- Provider `local`: OpenAI-compatible endpoint from `LOCAL_LLM_ENDPOINT`, no key, no free-tier check, smaller context fitting.
- Spike measures throughput and JSON validity for `qwen2.5:7b-instruct`; the form gains an explicit "Use local model" option.
- Expectations stay honest: CPU-only generation is slow and produces weaker prose than free cloud tiers.

## Phase E — text contract and Typst

- The prompt states the Markdown subset explicitly and forbids Typst syntax.
- Model output that contains Typst-like commands is normalized to literal text; unbalanced Markdown is repaired or rejected.
- Acceptance renders RU/EN PDFs with bold, emphasis, nested lists, a GFM table and a link, plus literal `#`, `$`, `[` and `@`.

## Order and out of scope

Order: A → B → E → C → D. Tree-sitter, embedding retrieval, extra providers and paid fallback are out of scope until the deterministic path is measured.
