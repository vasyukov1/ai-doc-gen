# Requirements index

Each requirement has a stable id, the user action it serves and the check that shows it works. Earlier phases are traced in `improvements_plan.md`; Phase C is detailed here because its checks are package-level.

| ID | Requirement | User action | Check |
|---|---|---|---|
| A1 | One submission source at a time | Choose a ZIP, then type a repository (or the reverse) | The other source clears and the form submits |
| A2 | Archive intake names the exact problem | Upload a ZIP with unsafe paths or read-only entries | Unsafe paths fail with the entry name; symlinks and duplicates are skipped with a reason |
| A3 | Fewer decisions before generation | Open the create form on a fresh browser and on a repeat visit | Language is chosen by the user, profile/language are restored only from the user's own choice, year defaults to the current one and stays editable |
| A4 | Paper signature fields | Render an academic profile | Signature and date render as blank lines, never as red TODO |
| B1 | The budget fits the document | Inspect generation details for a long profile | Attempts and time scale with chapter count; the shown limit matches the worker |
| B2 | Later chapters always get a call | Run a multi-chapter generation | Each unfinished chapter receives one call per pass before repair calls repeat |
| B3 | Generation continues after the deadline | Let a run hit its deadline with baseline sections left | Up to two bounded passes run; accepted sections stay intact |
| E1 | Model output follows the Markdown contract | Inspect a generated section with bold, lists and a table | The prompt states the subset, Typst directives are neutralized, cached sections from the old contract are not reused |
| E2 | Markdown renders correctly in PDF | Render RU and EN documents with rich Markdown and literal characters | Bold, emphasis, lists, quote, code, table and `#`, `$`, `[`, `@` all render |
| C1 | TypeScript and JavaScript declarations become evidence | Analyze a repository with `.ts`/`.js` sources | Imports, functions, classes and routes appear as facts with paths and line ranges |
| C2 | Python declarations become evidence | Analyze a repository with `.py` sources | Imports, functions, classes and route decorators appear as facts |
| C3 | YAML OpenAPI schemas become evidence | Analyze a repository whose only schema is `openapi.yaml` | Operations and component schemas appear as facts without a new dependency |
| C4 | Evidence ranking uses meaning, not substring luck | Compare two sections with different vocabulary | BM25 plus kind weights select relevant facts, duplicates are removed and services stay balanced |
| C5 | Long declarations keep their signatures | Analyze a file with a long function body | The fact starts with the declaration signature and is bounded |
| D1 | Generation works without a provider key | Enable the local model on the create form and run a document | The worker calls the local OpenAI-compatible endpoint with no key and the PDF is produced |
| F1 | Free-tier limits shape the draft | Connect one or several keys per provider and run a long profile | Published per-minute limits override larger form values, a context rejection retries with half the pack, and each key has its own rate window |
