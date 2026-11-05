JET_VERSION := v2.16.0
PG_PORT ?= 5434
HTTP_PORT ?= 8080
DATABASE_URL ?= postgres://ai_doc_gen:local_dev_only@127.0.0.1:$(PG_PORT)/ai_doc_gen?sslmode=disable

export PG_PORT HTTP_PORT DATABASE_URL

.PHONY: run build ui test down mocks db-generate models

run:
	@docker info >/dev/null 2>&1 || { echo "Docker is not running. Start it first, for example: colima start"; exit 1; }
	docker compose build migrate
	docker compose up -d
	@echo "Open http://localhost:$(HTTP_PORT)"

build:
	@mkdir -p bin
	go build -o bin/ai-doc-gen ./cmd/ai-doc-gen

ui:
	cd web && npm install --no-audit --no-fund && npm run build

test:
	go test ./...

down:
	docker compose down --remove-orphans

mocks:
	@echo "Mock generation is not configured yet."

models:
	@docker info >/dev/null 2>&1 || { echo "Docker is not running. Start it first, for example: colima start"; exit 1; }
	docker compose up -d ollama
	docker compose exec ollama ollama pull $${LOCAL_MODEL:-qwen2.5:7b-instruct}

db-generate:
	docker compose up -d postgres
	@go run ./cmd/ai-doc-gen migrate
	@mkdir -p bin
	@GOBIN='$(CURDIR)/bin' CGO_ENABLED=0 go install github.com/go-jet/jet/v2/cmd/jet@$(JET_VERSION)
	@bin/jet -dsn="$$DATABASE_URL" -schema=public -path=./internal/storage/generated -ignore-tables=goose_db_version
