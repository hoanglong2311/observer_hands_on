DATABASE_URL ?= postgres://observer:observer@localhost:5432/observer?sslmode=disable

.PHONY: run build test lint up down

run:
	go run ./cmd/server

build:
	go build -o bin/server ./cmd/server

test:
	go test ./...

lint:
	go vet ./...

up:
	docker compose up --build -d

down:
	docker compose down

logs:
	docker compose logs -f app

# Quick smoke test against running server
smoke:
	@echo "==> health"
	@curl -sf http://localhost:8080/healthz | jq .
	@echo "==> create task"
	@curl -sf -X POST http://localhost:8080/api/tasks \
		-H "Content-Type: application/json" \
		-d '{"title":"learn observability"}' | jq .
	@echo "==> list tasks"
	@curl -sf http://localhost:8080/api/tasks | jq .
