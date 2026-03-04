# Inventory API

[![Go](https://img.shields.io/badge/Go-1.24-00ADD8?logo=go&logoColor=white)](https://go.dev/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-17-4169E1?logo=postgresql&logoColor=white)](https://www.postgresql.org/)
[![Redis](https://img.shields.io/badge/Redis-7-DC382D?logo=redis&logoColor=white)](https://redis.io/)
[![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?logo=docker&logoColor=white)](https://www.docker.com/)

A robust Inventory Management REST API built with Go, focused on clean architecture, production-grade engineering, and high-performance stock management. It implements comprehensive features for managing products, suppliers, warehouses, and complex inventory movement scenarios.

## Overview

The primary goal of this project is to provide a solid foundation for industrial-scale inventory management with a layered architecture and clear separation of concerns.

Key features include:

- **Authentication & RBAC** — JWT-based authentication using RSA key pairs and Role-Based Access Control (RBAC) powered by Casbin.
- **Inventory Management** — Advanced stock control including additions, removals, reservations, and releases, with automatic transaction logging.
- **Product & Supplier Catalog** — Complete management of products, categories, and suppliers with relational integrity.
- **Warehousing** — Multi-warehouse support for distributed inventory tracking.
- **Caching** — Redis-backed caching for high-frequency data (like roles and categories) with automatic invalidation.
- **Rate Limiting** — Built-in protection against API abuse using a Redis-backed window rate limiter.
- **Observability** — Distributed tracing with OpenTelemetry and structured logging via `slog`, with Jaeger support.
- **Database Migrations** — Version-controlled schema migrations using `golang-migrate`.
- **API Documentation** — Fully interactive Swagger/OpenAPI documentation.

## Tech Stack

| Layer | Technology |
|---|---|
| Language | Go 1.24 |
| HTTP Framework | Gin |
| ORM | GORM |
| Authorization | Casbin (RBAC) |
| Database | PostgreSQL 17 |
| Cache | Redis 7 |
| Authentication | JWT (RSA) + Argon2id |
| Rate Limiting | ulule/limiter |
| Tracing | OpenTelemetry + Jaeger |
| Documentation | Swagger (swag) |
| Testing | testify + testcontainers-go |
| Containerization | Docker + Docker Compose |

## Project Structure

```
.
├── cmd/
│   └── api/                # Application entrypoint
├── deployments/            
│   ├── compose.yaml        # Docker Compose infrastructure
│   └── seed.sql            # Manual data seeding script
├── docs/                   # Auto-generated Swagger documentation
├── internal/
│   ├── apperrors/          # Error mapping and handling
│   ├── config/             # Configuration management (Env/File)
│   ├── database/           # Database initialization and seeding logic
│   ├── dto/                # Request/Response Data Transfer Objects
│   ├── handlers/           # HTTP controllers
│   ├── middleware/         # Auth, RBAC, Tracing, and Rate Limit middlewares
│   ├── models/             # Domain entities (Product, Inventory, etc.)
│   ├── pkg/                # Shared utilities (Logger, Cache, JWT)
│   ├── repository/         # Data access layer
│   ├── routes/             # API routing definitions
│   └── service/            # Business logic layer
├── migrations/             # SQL migration files
├── tests/
│   ├── integration/        # Integration tests with testcontainers
│   └── unit/               # Mocked unit tests
├── Makefile                # Development automation
└── policy.csv              # Casbin RBAC policies
```

## Prerequisites

| Tool | Version | Purpose |
|---|---|---|
| [Go](https://go.dev/dl/) | 1.24+ | Development environment |
| [Docker](https://www.docker.com/) | 20.10+ | Infrastructure containers |
| [Makefile](https://www.gnu.org/software/make/) | latest | Running automation commands |
| [OpenSSL](https://www.openssl.org/) | any | Generating RSA keys |

## Environment Variables

The application is configured via environment variables. Create a `.env` file in the root directory:

```bash
# App
ENV=development
APP_PORT=8080

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=inventory_user
DB_PASSWORD=inventory_password
DB_NAME=inventory

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379

# Tracing
OTLP_ENDPOINT=localhost:4317
```

## Getting Started

### 1. Initialize Infrastructure
Start the database and cache containers:
```bash
make docker-up
```

### 2. Prepare Security
Generate the RSA keys for JWT signing:
```bash
make generate-keys
```

### 3. Run Migrations
Apply the database schema:
```bash
make migration-up
```

### 4. Seed Data (Optional)
Populate the database with initial test data:
```bash
make seed
```

### 5. Start the API
```bash
make start
```
The API will be available at `http://localhost:8080` and Swagger docs at `http://localhost:8080/docs/index.html`.

## Available Makefile Commands

| Command | Description |
|---|---|
| `make start` | Run the application in development mode |
| `make seed` | Populate the database from `deployments/seed.sql` |
| `make build` | Compile the binary to `bin/api` |
| `make test` | Run all tests |
| `make lint` | Run code quality checks (golangci-lint) |
| `make migration-up` | Apply pending SQL migrations |
| `make migration-down` | Rollback the last migration |
| `make db-dump` | Create a data-only dump of the current DB |
| `make swagger` | Update API documentation |
| `make docker-up-all` | Run the full stack in Docker (API included) |

## Database Migrations

We use a two-step approach for database management:
1. **Schema**: Managed via `golang-migrate` scripts in the `migrations/` folder.
2. **Data**: Managed via `deployments/seed.sql` for manual developer seeding.

To create a new migration:
```bash
make migration-create name=my_new_table
```

## Testing

The project implements a testing pyramid approach:
- **Unit Tests**: Fast tests using mocks (located in subpackages).
- **Integration Tests**: Comprehensive tests using `testcontainers-go` to spin up real Postgres and Redis instances.

Run tests:
```bash
make test-unit
make test-integration
```

## Observability

- **Tracing**: Fully instrumented with OpenTelemetry. You can view spans in **Jaeger** by accessing `http://localhost:16686`.
- **Logging**: Structured JSON logs using `slog` for production environments.

## Security Notice

- Use `make generate-keys` to create unique signing keys.
- **Casbin** is used for RBAC. Policy definitions can be found in `policy.csv`.
- Rate limiting is enabled by default to prevent brute-force attacks on auth endpoints.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
