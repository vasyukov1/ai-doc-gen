# Development guide

This directory turns the product analysis into bounded implementation requests. Each top-level item in [plan.md](plan.md) is one Codex request. Finish its subitems and acceptance criteria, then review the result before starting the next item. The current repository includes item 1 only: the minimal Go website and local PostgreSQL service. The form is a preview and cannot submit a job yet.

Read [structure.md](structure.md) for package ownership, [code_style.md](code_style.md) for Go style, [job_lifecycle.md](job_lifecycle.md) for progress and failures, and [decisions.md](decisions.md) for trade-offs.

## Local start

```bash
docker compose up -d postgres
go run ./cmd/ai-doc-gen
```

Open `http://localhost:8080`. `GET /healthz` returns `ok`. The site itself has no database dependency until item 2. PostgreSQL is exposed on `127.0.0.1:5433` for local development; its credentials in `docker-compose.yaml` are for local use only.

## Working rules

- Use [plan.md](plan.md) to select the next numbered request.
- Use [decisions.md](decisions.md) for architecture choices and changes.
- Keep the UI and server small and understandable. Prefer explicit data flow to framework setup.
- Automated tests are intentionally scheduled for the final item, per the product owner's instruction. Build and manually smoke-check each intermediate item.
- The installed Superpowers skills are available to Codex on the next turn. This repository's direct instructions determine the architecture and the timing of tests.
