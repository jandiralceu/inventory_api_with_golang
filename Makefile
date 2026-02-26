.PHONY: start help build test test-unit test-integration test-cover clean migration-create migration-up migration-down docker-up docker-up-all docker-stop docker-down generate-keys

include .env
export

BINARY_NAME=api
MIGRATIONS_PATH = migrations
DB_URL = postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable

start: ## Start the application
	@echo "Starting the application..."
	go run cmd/api/main.go

build: ## Build the application
	@echo "Building the application..."
	@go build -o bin/$(BINARY_NAME) cmd/api/main.go

test: ## Run all tests
	@make test-unit
	@make test-integration

test-unit: ## Run unit tests
	go test ./... -short -v

test-integration: ## Run integration tests
	go test -tags=integration ./tests/integration/... -v

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
	docker compose up -d 

docker-up-all: generate-keys ## Start ALL docker containers including the API app in detached mode
	docker compose --profile all up -d --build

docker-stop: ## Stop docker containers without removing them
	docker compose --profile all stop

docker-down: ## Stop and remove docker containers
	docker compose --profile all down

generate-keys: ## Generate RSA keys for JWT
	@if [ ! -f private.pem ]; then \
		echo "Generating RSA keys..."; \
		openssl genrsa -out private.pem 2048; \
		openssl rsa -in private.pem -pubout -out public.pem; \
		echo "Keys generated: private.pem, public.pem"; \
	else \
		echo "Keys already exist, skipping generation."; \
	fi

help: ## Display all available commands
	@echo "Available commands:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-18s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)
