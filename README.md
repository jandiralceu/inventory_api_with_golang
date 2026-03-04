# Inventory API

[![Go](https://img.shields.io/badge/Go-1.24-00ADD8?logo=go&logoColor=white)](https://go.dev/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-17-4169E1?logo=postgresql&logoColor=white)](https://www.postgresql.org/)
[![Redis](https://img.shields.io/badge/Redis-7-DC382D?logo=redis&logoColor=white)](https://redis.io/)
[![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?logo=docker&logoColor=white)](https://www.docker.com/)
[![Kubernetes](https://img.shields.io/badge/Kubernetes-Ready-326CE5?logo=kubernetes&logoColor=white)](https://kubernetes.io/)

A production-grade Inventory Management REST API built with Go, designed around clean architecture principles and cloud-native best practices. It provides comprehensive stock control — including product cataloging, multi-warehouse tracking, supplier management, and advanced inventory operations (additions, removals, reservations, and releases) — all backed by JWT authentication, role-based access control, distributed tracing, and Kubernetes-ready deployment manifests.

## Overview

This project goes beyond a simple CRUD API. It demonstrates idiomatic Go development with a layered architecture, clear separation of concerns, and production-ready tooling that can be extended and deployed to any cloud environment:

- **Authentication & RBAC** — JWT-based authentication using RSA key pairs with secure password hashing via Argon2id. Fine-grained Role-Based Access Control powered by [Casbin](https://casbin.org/) with three predefined roles: `admin`, `manager`, and `operator`.
- **Inventory Operations** — Full stock lifecycle management: add stock, remove stock, reserve for orders, and release reservations. Every operation is recorded as an immutable transaction for complete audit trails.
- **Product & Supplier Catalog** — Complete CRUD for products, categories (with parent-child hierarchy), and suppliers with structured address data (JSONB).
- **Multi-Warehouse Support** — Track inventory across multiple warehouses with location codes, min/max quantity thresholds, and per-warehouse stock levels.
- **Caching** — Redis-backed caching layer with automatic prefix-based invalidation on write operations to ensure data consistency.
- **Rate Limiting** — Redis-backed sliding window rate limiter with configurable limits per route group (global, auth, password changes).
- **Observability** — Distributed tracing with OpenTelemetry (OTLP/gRPC) and structured logging via `slog`, with Jaeger as the default trace backend.
- **Database Migrations** — Version-controlled schema migrations using `golang-migrate`, keeping schema management separate from application code.
- **API Documentation** — Auto-generated Swagger/OpenAPI docs served at `/docs/index.html`.
- **Graceful Shutdown** — The server handles OS signals (SIGINT/SIGTERM) and drains in-flight requests before stopping.
- **CI/CD Pipelines** — GitHub Actions workflows for linting, testing, building, and pushing Docker images to Docker Hub.
- **Kubernetes Manifests** — Ready-to-deploy K8s configurations with Deployments, Services, ConfigMaps, Secrets, health probes, and resource limits.

## Tech Stack

| Layer | Technology |
|---|---|
| Language | Go 1.24 |
| HTTP Framework | Gin |
| ORM | GORM |
| Database | PostgreSQL 17 |
| Cache | Redis 7 |
| Authentication | JWT (RSA) + Argon2id |
| Authorization | Casbin (RBAC) |
| Rate Limiting | ulule/limiter |
| Migrations | golang-migrate |
| Tracing | OpenTelemetry + Jaeger |
| Documentation | Swagger (swag) |
| Testing | testify + testcontainers-go |
| Containerization | Docker + Docker Compose |
| Orchestration | Kubernetes (manifests included) |
| CI/CD | GitHub Actions |

## Architecture

The project follows a **layered architecture** with strict dependency rules — each layer only depends on the layer directly below it:

```
┌─────────────────────────────────────────────────────┐
│                   HTTP Layer                        │
│   Middleware (Auth, RBAC, Tracing, Rate Limit)      │
│   Handlers (request parsing, validation, response)  │
├─────────────────────────────────────────────────────┤
│                 Service Layer                       │
│   Business logic, caching strategies, validation    │
├─────────────────────────────────────────────────────┤
│                Repository Layer                     │
│   Data access (GORM), query building, error mapping │
├─────────────────────────────────────────────────────┤
│               Infrastructure Layer                  │
│   PostgreSQL, Redis, OpenTelemetry, JWT, Argon2id   │
└─────────────────────────────────────────────────────┘
```

### RBAC Model

Access control is defined in `policy.csv` using the Casbin ACL model:

| Role | Permissions |
|---|---|
| **admin** | Full access to all endpoints |
| **manager** | CRUD on categories, suppliers, warehouses, products. Read-only on roles. Can change own password. |
| **operator** | Read-only on categories, suppliers, warehouses, products. Can change own password. |

## Project Structure

```
.
├── cmd/api/                   # Application entrypoint and server bootstrap
├── deployments/
│   ├── compose.yaml           # Docker Compose for local development
│   ├── seed.sql               # Data-only seed script (schema is managed by migrations)
│   └── k8s/
│       ├── base/              # K8s manifests: API Deployment, Service, ConfigMap, Secret, Namespace
│       └── local/             # K8s manifests: Postgres, Redis, Jaeger (for local cluster testing)
├── docs/                      # Auto-generated Swagger documentation
├── internal/
│   ├── apperrors/             # Sentinel errors mapped to HTTP status codes
│   ├── config/                # Environment-based configuration loading
│   ├── database/              # PostgreSQL connection initialization
│   ├── dto/                   # Request/Response Data Transfer Objects
│   ├── handlers/              # Gin HTTP handlers with Swagger annotations
│   ├── middleware/            # Auth (JWT), RBAC (Casbin), Tracing, Rate Limiting
│   ├── models/                # GORM domain models
│   ├── pkg/                   # Shared utilities: Logger, Cache, JWT, Hasher, Telemetry
│   ├── repository/            # Data access layer (GORM implementations)
│   ├── routes/                # Route registration and middleware wiring
│   └── service/               # Business logic with caching strategies
├── migrations/                # SQL migration files (golang-migrate)
├── tests/integration/         # Integration tests using testcontainers
├── .github/workflows/         # CI (lint, test, build) and CD (push to Docker Hub)
├── Dockerfile                 # Multi-stage production build
├── Makefile                   # Developer workflow automation
├── model.conf                 # Casbin RBAC model configuration
└── policy.csv                 # Casbin RBAC policy definitions
```

## Prerequisites

| Tool | Version | Purpose |
|---|---|---|
| [Go](https://go.dev/dl/) | 1.24+ | Build and run the application |
| [Docker](https://docs.docker.com/get-docker/) | 20.10+ | Run PostgreSQL, Redis, and Jaeger containers |
| [Docker Compose](https://docs.docker.com/compose/install/) | 2.0+ | Orchestrate local infrastructure |
| [OpenSSL](https://www.openssl.org/) | any | Generate RSA key pairs for JWT signing |

## Environment Variables

The application loads settings from environment variables, with an optional `.env` file for local development. Copy the example file to get started:

```bash
cp .env.example .env
```

| Variable | Required | Default | Description |
|---|---|---|---|
| `DB_HOST` | ✅ | — | PostgreSQL host address |
| `DB_PORT` | ✅ | — | PostgreSQL port |
| `DB_USER` | ✅ | — | PostgreSQL username |
| `DB_PASSWORD` | ✅ | — | PostgreSQL password |
| `DB_NAME` | ✅ | — | PostgreSQL database name |
| `APP_PORT` | ❌ | `8080` | HTTP server listening port |
| `APP_NAME` | ❌ | `inventory-api` | Application name used in tracing spans |
| `ENV` | ❌ | `development` | Environment (`development` or `production`) |
| `PRIVATE_KEY_PATH` | ❌ | `private.pem` | Path to RSA private key for JWT signing |
| `PUBLIC_KEY_PATH` | ❌ | `public.pem` | Path to RSA public key for JWT validation |
| `OTLP_ENDPOINT` | ❌ | `localhost:4317` | OpenTelemetry Collector gRPC endpoint |
| `REDIS_HOST` | ❌ | `localhost` | Redis server host address |
| `REDIS_PORT` | ❌ | `6379` | Redis server port |
| `REDIS_PASSWORD` | ❌ | *(empty)* | Redis authentication password |

> **Note:** In `development` mode the logger outputs human-readable text at DEBUG level. In `production` it switches to structured JSON at INFO level.

## Getting Started

### Option 1 — Run everything with Docker

The fastest way to get the full stack running (API + PostgreSQL + Redis + Jaeger + migrations):

```bash
# 1. Clone the repository
git clone https://github.com/jandiralceu/inventory_api_with_golang.git
cd inventory_api_with_golang

# 2. Generate RSA keys for JWT authentication
make generate-keys

# 3. Start all containers (builds the API image and runs migrations automatically)
make docker-up-all
```

The API will be available at `http://localhost:8080` and Swagger docs at `http://localhost:8080/docs/index.html`.

### Option 2 — Local development

Run the Go application natively while using Docker only for infrastructure:

```bash
# 1. Clone the repository
git clone https://github.com/jandiralceu/inventory_api_with_golang.git
cd inventory_api_with_golang

# 2. Start infrastructure containers (PostgreSQL, Redis, Jaeger)
#    The migrator container applies database migrations automatically.
make docker-up

# 3. Generate RSA keys for JWT authentication
make generate-keys

# 4. Copy and configure environment variables
cp .env.example .env
# Edit .env with your database credentials (match compose.yaml values)

# 5. Install Go dependencies
go mod download

# 6. (Optional) Seed the database with sample data
make seed

# 7. Run the application
make start
```

### Available Makefile Commands

| Command | Description |
|---|---|
| `make start` | Run the application with `go run` |
| `make build` | Compile the application into `bin/api` |
| `make seed` | Populate the database from `deployments/seed.sql` |
| `make test` | Run all tests (unit + integration) |
| `make test-unit` | Run unit tests only |
| `make test-integration` | Run integration tests only |
| `make test-bench` | Run benchmarks |
| `make test-cover` | Generate HTML coverage report |
| `make lint` | Run code quality checks (`golangci-lint`) |
| `make swagger` | Regenerate Swagger documentation |
| `make generate-keys` | Generate RSA key pair for JWT |
| `make migration-create name=xxx` | Create a new migration file |
| `make migration-up` | Apply all pending migrations |
| `make migration-down` | Rollback the last migration |
| `make db-dump` | Create a data-only dump of the current database |
| `make db-restore file=xxx` | Restore database from a SQL dump file |
| `make docker-up` | Start infrastructure containers only |
| `make docker-up-all` | Start all containers including the API |
| `make docker-stop` | Stop all containers |
| `make docker-down` | Stop and remove all containers |
| `make clean` | Remove binaries, coverage files, and keys |

## API Documentation

Full interactive documentation is available via Swagger UI at `/docs/index.html` when the server is running.

The API is organized under `/api/v1` with the following resource groups:

| Group | Prefix | Auth | Description |
|---|---|---|---|
| Authentication | `/auth` | No | Sign in, register, refresh tokens, sign out |
| Roles | `/roles` | Partial | Public listing; CRUD requires auth |
| Users | `/users` | Yes | User management, password & role changes |
| Categories | `/categories` | Yes | Product categorization with hierarchy support |
| Suppliers | `/suppliers` | Yes | Supplier management with structured addresses |
| Warehouses | `/warehouses` | Yes | Multi-warehouse location management |
| Products | `/products` | Yes | Product catalog with SKU, pricing, and metadata |
| Inventories | `/inventories` | Yes | Stock tracking, operations (add/remove/reserve/release), transaction history |

## Database Migrations

Schema is managed exclusively by `golang-migrate`. The `deployments/seed.sql` file contains **data only** — table creation is the responsibility of the migration files in `migrations/`.

### Workflow

```bash
# 1. Create a new migration
make migration-create name=add_phone_to_suppliers

# 2. Edit the generated SQL files in migrations/

# 3. Apply migrations
make migration-up

# 4. (Optional) Seed with test data
make seed
```

> **Note:** When using `make docker-up` or `make docker-up-all`, migrations are applied automatically by the `migrate` container before the API starts.

## Testing

The project has two layers of tests: **unit tests** that run in isolation with mocks, and **integration tests** that spin up real PostgreSQL and Redis containers via [testcontainers-go](https://golang.testcontainers.org/).

### Run all tests

```bash
make test
```

### Unit tests only

```bash
make test-unit
```

Unit tests cover the following packages:

| Package | What is tested |
|---|---|
| `handlers` | HTTP status codes, request binding, error responses |
| `service` | Business logic, caching behavior, stock operations |
| `repository` | Query building, error mapping (using sqlmock) |
| `middleware` | JWT validation, rate limiting, trace ID propagation |
| `pkg` | Cache operations (miniredis), JWT signing, Argon2id hashing |
| `dto` | Pagination defaults and calculations |

### Integration tests only

```bash
make test-integration
```

Integration tests use the `integration` build tag and require Docker to be running. They start real PostgreSQL and Redis containers, apply migrations, and exercise the full HTTP request lifecycle (handler → service → repository → database).

### Coverage report

```bash
make test-cover
# Opens coverage.html in the project root
```

## Observability

The application is instrumented with [OpenTelemetry](https://opentelemetry.io/) for distributed tracing. Traces are exported via OTLP/gRPC to a collector — by default, [Jaeger](https://www.jaegertracing.io/) running as a Docker container.

**What is traced:**

- Every incoming HTTP request (via `otelgin` middleware)
- All database queries (via `otelgorm` plugin)
- Redis cache operations (via `redisotel` instrumentation)

**Accessing the Jaeger UI:**

Once the containers are running (`make docker-up`), open [http://localhost:16686](http://localhost:16686) in your browser to explore traces, inspect latency, and debug request flows across the entire stack.

> **Tip:** In production, you can replace Jaeger by pointing `OTLP_ENDPOINT` to any OpenTelemetry-compatible backend (e.g., Google Cloud Trace, Datadog, Grafana Tempo) without changing application code.

## Kubernetes Deployment

The project includes production-ready Kubernetes manifests organized following the base/overlay pattern:

```
deployments/k8s/
├── base/                                # Core application manifests
│   ├── inventory-namespace.yaml         # Dedicated namespace
│   ├── inventory-api-deployment.yaml    # 2 replicas, rolling updates, health probes, resource limits
│   ├── inventory-api-service.yaml       # ClusterIP service
│   ├── inventory-api-configmap.yaml     # Non-sensitive configuration
│   └── inventory-api-secret.yaml        # DB credentials, RSA keys
└── local/                               # Infrastructure for local cluster (minikube, kind)
    ├── inventory-postgres-*.yaml        # PostgreSQL Deployment + Service
    ├── inventory-redis-*.yaml           # Redis Deployment + Service
    └── inventory-jaeger-*.yaml          # Jaeger Deployment + Service
```

### Highlights

- **Rolling Updates** with `maxSurge: 1` and `maxUnavailable: 0` for zero-downtime deployments.
- **Health Probes** — `livenessProbe` and `readinessProbe` configured against the `/health` endpoint.
- **Resource Limits** — CPU and memory requests/limits defined for predictable scheduling.
- **Secrets Management** — RSA keys and database credentials stored as Kubernetes Secrets, mounted as read-only volumes.
- **Non-root User** — The Docker image runs as a non-root user (`appuser`) for security.

## CI/CD

GitHub Actions workflows automate the entire delivery pipeline:

| Workflow | Trigger | Steps |
|---|---|---|
| **CI** (`ci.yml`) | Pull Request → `main` | Lint → Unit Tests → Integration Tests → Go Build → Docker Build (dry run) |
| **CD** (`cd.yml`) | Push → `main` | Unit Tests + Coverage → Integration Tests + Coverage → Upload to Codecov → Build & Push to Docker Hub |

## ⚠️ Security Notice

The `compose.yaml` file contains **hardcoded database credentials** for local development convenience. RSA private keys (`*.pem`) and `.env` files are already **excluded from version control** via `.gitignore`, but are mounted into the Docker container as read-only volumes for local use.

In a production environment you **must**:

- Store database credentials, Redis passwords, and API keys in a secrets manager (e.g., AWS Secrets Manager, GCP Secret Manager, HashiCorp Vault) or inject them via environment variables at deploy time.
- Provision RSA key pairs through your CI/CD pipeline or secrets manager instead of generating them locally.

## License

This project is licensed under the MIT License — see the [LICENSE](LICENSE) file for details.
