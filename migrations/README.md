# Database migrations

Goose SQL migrations are embedded into the Go binary. Run `ai-doc-gen migrate` before starting web or worker; `make run` does this through a one-shot Compose service.
Add one numbered migration for each schema change.
Keep each migration reversible where possible and document data migrations that cannot be rolled back safely.
Run `make db-generate` after applying a schema change and commit the regenerated Jet code.

`00011_catalan_language.sql` expands the existing text-language constraint to `en`, `ru` and `ca`; Go/Jet column types are unchanged. Its down migration deliberately fails if Catalan jobs remain, rather than relabeling or deleting their documents. Apply the up migration before accepting Catalan submissions.
