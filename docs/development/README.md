# Development guide

This directory turns the product analysis into bounded implementation requests. Each top-level item in [plan.md](plan.md) is one Codex request. Finish its subitems and acceptance criteria, then review the result before starting the next item. Items 1, 2, and 2.5 are implemented: the site accepts input, the worker checks it, and database queries use Jet. Source inspection and document generation remain future items.

Read [structure.md](structure.md) for package ownership, [code_style.md](code_style.md) for Go style, [job_lifecycle.md](job_lifecycle.md) for progress and failures, and [decisions.md](decisions.md) for trade-offs.

## Local start

```bash
make run
make down
```

Open `http://localhost:8080`. `GET /healthz` returns `ok`; `GET /readyz` checks PostgreSQL. A one-shot migration container applies embedded goose migrations before the web and worker start. PostgreSQL is exposed on `127.0.0.1:5434` by default for local development; set `PG_PORT` or `HTTP_PORT` to change host ports. The credentials in `../../docker-compose.yaml` are for local use only. `make down` preserves database and upload volumes.

For a host-only Go workflow, start PostgreSQL with `docker compose up -d postgres`, export `DATABASE_URL` from `../../.env.example`, run `go run ./cmd/ai-doc-gen migrate`, then run `serve` and `worker` in separate terminals. After changing a migration, run `make db-generate` and commit the generated Jet files with the migration.

## Working rules

- Use [plan.md](plan.md) to select the next numbered request.
- Use [decisions.md](decisions.md) for architecture choices and changes.
- Keep the UI and server small and understandable. Prefer explicit data flow to framework setup.
- Automated tests are intentionally scheduled for the final item, per the product owner's instruction. Build and manually smoke-check each intermediate item.
- The installed Superpowers skills are available to Codex on the next turn. This repository's direct instructions determine the architecture and the timing of tests.
