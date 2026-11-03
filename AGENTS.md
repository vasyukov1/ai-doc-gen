# Project instructions

Read `docs/development/plan.md`, `structure.md`, `code_style.md`, `job_lifecycle.md`, and `decisions.md` before implementing a milestone. Complete only the requested numbered milestone and its acceptance criteria. Do not begin the next milestone automatically.

- Keep repository documentation, code comments, commit messages, and agent instructions in English. The product may generate Russian or English documents as selected by the user.
- Keep the application in this repository. Do not import the earlier prototype at runtime or mention its repository name in project files.
- Follow the package boundaries in `docs/development`. Keep business logic in `internal/documents` and database operations in `internal/storage` or `internal/pgqueue`.
- Prefer Go standard library and explicit constructors. Add a dependency only when a concrete requirement makes it simpler overall. Do not add `uber/fx`, Kafka, Redis, or a state-machine package without measured need and a short decision record.
- Use short English comments, errors, and logs. Separate logical phases inside functions with blank lines. Follow `docs/development`.
- Model document progress with explicit states and legal transitions. Persist job state and progress events in PostgreSQL when the jobs milestone is implemented. Use a separate Go worker process for long-running work. Surface actionable errors and partial results promptly.
- Never execute uploaded repository code. Treat uploaded files and generated Typst as untrusted input.
- Keep user-facing controls honest: unfinished actions must be disabled or clearly marked as unavailable.
- Follow the user's sequence: defer automated test files and test suites to the final testing milestone. At each earlier milestone, run `go build ./...` and perform a small manual smoke check of the behavior delivered. Record any unverified behavior. Use table-driven tests when the final milestone starts.
- Superpowers skills are installed for Codex. Use relevant planning, debugging, and verification guidance, while preserving the user's milestone sequence and test timing.
- Avoid changing existing analysis decisions silently. When a choice changes, update `docs/development` and the affected plan section.
