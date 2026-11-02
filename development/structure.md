# Backend package map

The directories below are stable boundaries. A directory with only a README is reserved for a later numbered milestone; it is not an empty Go package.

```text
cmd/ai-doc-gen/main.go       application entry point
internal/config/             environment and runtime configuration
internal/documents/          code intake and documentation business logic
internal/domain/             job, document, source and evidence types
internal/kafka/              optional broker adapter, only if needed
internal/pb/                 generated protobuf code, only if needed
internal/pgqueue/            PostgreSQL claim, lease and retry operations
internal/platform/           logging, metrics, locks, parsing, transactions,
                             identifiers and validation primitives
internal/statemachine/       legal job and document state transitions
internal/storage/            PostgreSQL repositories and artifact storage
internal/templates/          versioned document profile and Typst assets
internal/tests/              integration and end-to-end tests at final milestone
internal/web/                HTTP routes and temporary HTML UI
migrations/                  goose SQL migrations
docker-compose.yaml          local PostgreSQL
go.mod                       module definition
```

Use `internal/documents` for orchestration and domain decisions about source files and sections. `internal/storage` persists data without deciding what a document means. `internal/pgqueue` owns queue mechanics, not document rules. `internal/platform` holds only small reusable infrastructure adapters; avoid putting arbitrary helpers there. `internal/web` remains a delivery adapter and may later serve a TypeScript frontend without moving business logic.

The entry point wires configuration, storage, queue, worker and HTTP components explicitly. A worker mode can run from the same binary. Add Kafka or protobuf only when a concrete integration requires them.
