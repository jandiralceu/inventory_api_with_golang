.PHONY: start seed help build test lint test-unit test-integration test-bench test-cover clean migration-create migration-up migration-down docker-up docker-up-all docker-stop docker-down generate-keys swagger db-dump db-restore

include .env
export

BINARY_NAME=api
MIGRATIONS_PATH = migrations
DB_URL = postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable

start: ## Start the application
	@echo "Starting the application..."
	go run cmd/api/main.go

seed: ## Seed individual roles into the database (idempotent)
	@echo "Seeding database..."
	go run cmd/seed/main.go

build: swagger ## Build the application
	@echo "Building the application..."
	@go build -o bin/$(BINARY_NAME) cmd/api/main.go

lint: ## Run linter
	@echo "Running linter..."
	@golangci-lint run --build-tags integration

test: ## Run all tests
	@make test-unit
	@make test-integration

test-unit: ## Run unit tests
	go test ./... -short -v

test-integration: ## Run integration tests
	go test -tags=integration ./tests/integration/... -v

test-bench: ## Run benchmarks
	go test -bench=. -benchmem ./internal/...

test-cover: ## Generate test coverage report
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated at coverage.html"

clean: ## Remove binary, test artifacts and keys
	@echo "Cleaning up..."
	@rm -rf bin/ coverage.out coverage.html *.pem

migration-create: ## Create a new migration. Usage: make migration-create name=<name>
	@if [ -z "$(name)" ]; then echo "Error: 'name' is required. Usage: make migration-create name=my_migration_name"; exit 1; fi
	migrate create -ext sql -dir $(MIGRATIONS_PATH) -seq $(name)

migration-up: ## Run all pending migrations
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" up

migration-down: ## Rollback migrations. Usage: make migration-down [step=1]
	@migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" down $(if $(step),$(step),1)

docker-up: ## Start only background docker containers (DB, Redis, Jaeger) in detached mode
	docker compose -f deployments/compose.yaml up -d 

docker-up-all: generate-keys ## Start ALL docker containers including the API app in detached mode
	docker compose -f deployments/compose.yaml --profile all up -d --build

docker-stop: ## Stop docker containers without removing them
	docker compose -f deployments/compose.yaml --profile all stop

docker-down: ## Stop and remove docker containers
	docker compose -f deployments/compose.yaml --profile all down

db-dump: ## Dump database to deployments/seed.sql (override with file=path/to/file.sql)
	$(eval DUMP_FILE := $(if $(file),$(file),deployments/seed_$(shell date +%Y%m%d_%H%M%S).sql))
	@echo "Creating database dump -> $(DUMP_FILE)"
	@docker exec -t inventory-postgres pg_dump -U $(DB_USER) -d $(DB_NAME) --no-owner --no-acl > $(DUMP_FILE)
	@echo "Dump saved to $(DUMP_FILE)"

db-restore: ## Restore database from a dump. Usage: make db-restore file=deployments/seed.sql
	@if [ -z "$(file)" ]; then echo "Error: 'file' is required. Usage: make db-restore file=deployments/seed.sql"; exit 1; fi
	@echo "Restoring database from $(file)..."
	@docker exec -i inventory-postgres psql -U $(DB_USER) -d $(DB_NAME) < $(file)
	@echo "Database restored from $(file)"

generate-keys: ## Generate RSA keys for JWT
	@if [ ! -f private.pem ]; then \
		echo "Generating RSA keys..."; \
		openssl genrsa -out private.pem 2048; \
		openssl rsa -in private.pem -pubout -out public.pem; \
		echo "Keys generated: private.pem, public.pem"; \
	else \
		echo "Keys already exist, skipping generation."; \
	fi

swagger: ## Generate swagger documentation
	@echo "Generating swagger documentation..."
	swag init -g cmd/api/main.go

help: ## Display all available commands
	@echo "Available commands:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-18s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)
