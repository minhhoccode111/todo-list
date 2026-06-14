![Logo](docs/img/logo.svg)

# Todo List API

A REST API for a todo list, built with Go and Clean Architecture.

[![Release](https://img.shields.io/github/v/release/minhhoccode111/todo-list.svg)](https://github.com/minhhoccode111/todo-list/releases/)
[![License](https://img.shields.io/badge/License-MIT-success)](https://github.com/minhhoccode111/todo-list/blob/master/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/minhhoccode111/todo-list)](https://goreportcard.com/report/github.com/minhhoccode111/todo-list)
[![codecov](https://codecov.io/gh/minhhoccode111/todo-list/branch/master/graph/badge.svg?token=XE3E0X3EVQ)](https://codecov.io/gh/minhhoccode111/todo-list)

[![Web Framework](https://img.shields.io/badge/Gin-Web%20Framework-blue)](https://github.com/gin-gonic/gin)
[![API Documentation](https://img.shields.io/badge/Swagger-API%20Documentation-blue)](https://github.com/swaggo/swag)
[![Validation](https://img.shields.io/badge/Validator-Data%20Integrity-blue)](https://github.com/go-playground/validator)
[![JSON Handling](https://img.shields.io/badge/Go--JSON-Fast%20Serialization-blue)](https://github.com/goccy/go-json)
[![Query Builder](https://img.shields.io/badge/sqlc-SQL%20Compiler-blue)](https://sqlc.dev/)
[![Database Migrations](https://img.shields.io/badge/Migrations-Seamless%20Schema%20Updates-blue)](https://github.com/golang-migrate/migrate)
[![Logging](https://img.shields.io/badge/ZeroLog-Structured%20Logging-blue)](https://github.com/rs/zerolog)
[![Metrics](https://img.shields.io/badge/Prometheus-Metrics%20Integration-blue)](https://github.com/zsais/go-gin-prometheus)
[![Testing](https://img.shields.io/badge/Testify-Testing%20Framework-blue)](https://github.com/stretchr/testify)
[![Mocking](https://img.shields.io/badge/Mock-Mocking%20Library-blue)](https://go.uber.org/mock)

## Overview

This service exposes a JSON REST API for managing todos behind user
authentication. It started as a fork of
[go-clean-template](https://github.com/evrone/go-clean-template) but was
adapted into a single REST application:

- Gin instead of Fiber
- sqlc instead of Squirrel
- A `validatorx` wrapper in `pkg`
- An [Otter](https://github.com/maypok86/otter) in-memory cache in `pkg`
- RPC transports (gRPC, NATS, RabbitMQ) removed — REST only

The goal is to keep business logic independent, clean, and extensible by
following Clean Architecture (Robert C. Martin / "Uncle Bob"). A SvelteKit
single-page app under `frontend/` consumes the API; see the top-level
[README](README.md) for the product-facing notes and screenshots.

## Features

- **Auth** — register, login, refresh, and logout flows using short-lived JWT
  access tokens plus hashed refresh tokens stored in the database and delivered
  in http-only cookies (so they can be revoked).
- **Sessions** — list active device sessions, log out a single session by id,
  or log out everywhere.
- **Todos** — create, read (paginated, cached), update, and delete todos scoped
  to the authenticated user.
- **Cross-cutting** — request logging, panic recovery, CORS, per-IP rate
  limiting, Prometheus metrics, and Swagger documentation with bearer auth.

## Tech stack

- [Gin](https://github.com/gin-gonic/gin) — HTTP framework
- [sqlc](https://sqlc.dev/) — type-safe Go from SQL, with
  [pgx](https://github.com/jackc/pgx) for PostgreSQL
- [golang-migrate](https://github.com/golang-migrate/migrate) — schema migrations
- [swag](https://github.com/swaggo/swag) — Swagger/OpenAPI generation
- [validator](https://github.com/go-playground/validator) — request validation
- [zerolog](https://github.com/rs/zerolog) — structured logging
- [Otter](https://github.com/maypok86/otter) — in-memory cache
- [testify](https://github.com/stretchr/testify) + [go.uber.org/mock](https://go.uber.org/mock) — tests and mocks

## Quick start

Requires Go, Docker (for PostgreSQL), and the tools installed by `make bin-deps`
(`migrate`, `sqlc`, `swag`, `air`, `swagger-typescript-api`).

```sh
# Start PostgreSQL via docker compose
make db-up

# Generate code, run migrations, and start the API (http://127.0.0.1:8080)
make run

# Or run with live reload
air
```

Configuration is read from `.env` (falling back to `.env.example`). Once
running:

- Health check: http://127.0.0.1:8080/healthz
- Swagger UI: http://127.0.0.1:8080/swagger/index.html
- Metrics (if enabled): http://127.0.0.1:8080/metrics

To run the full stack (API + reverse proxy) in Docker, use `make compose-up-all`.

## API

All endpoints are served under `/api/v1`. Routes marked 🔒 require a
`Bearer <access token>` header.

| Method | Path                 | Description                          |
| ------ | -------------------- | ------------------------------------ |
| POST   | `/register`          | Create an account                    |
| POST   | `/login`             | Log in, receive access + refresh     |
| POST   | `/refresh`           | Exchange a refresh token             |
| POST   | `/logout`         🔒 | Log out the current session          |
| POST   | `/logout/all`     🔒 | Log out all sessions                 |
| GET    | `/sessions`       🔒 | List active sessions                 |
| DELETE | `/sessions/:id`   🔒 | Revoke a session by id               |
| GET    | `/todos`          🔒 | List todos (paginated)               |
| POST   | `/todos`          🔒 | Create a todo                        |
| PUT    | `/todos/:id`      🔒 | Update a todo                        |
| DELETE | `/todos/:id`      🔒 | Delete a todo                        |

The Swagger spec (`docs/swagger.yaml`) is the source of truth and is
auto-generated from handler annotations; see [Make commands](#make-commands).

## Project structure

### `cmd/app/main.go`

Loads configuration and the logger, then hands off to `internal/app`.

### `config`

Twelve-factor configuration from environment variables. See
[config.go](config/config.go) and [.env.example](.env.example).

### `docs`

Auto-generated Swagger documentation (`docs.go`, `swagger.json`, `swagger.yaml`)
plus the database schema dump (`schema.sql`) and images. Generated files should
not be edited by hand.

### `internal/app`

Holds the single `Run` function that "continues" `main`. This is where every
object is created and wired together via `New...` constructors (see
[Dependency Injection](#dependency-injection)), the HTTP server is started, and
graceful shutdown is handled. With the `migrate` build tag, migrations run
automatically on start:

```sh
go run -tags migrate ./cmd/app
```

### `internal/controller/restapi`

The HTTP delivery layer (Gin). `router.go` wires global middleware, metrics,
Swagger, and the health check; `middleware/` holds cross-cutting middleware
(auth, logging, recovery, CORS, rate limiting); `v1/` holds the versioned
handlers grouped by area (`user`, `todo`) with their request/response types.
Swagger annotations live above the handler methods. Adding a `v2/` package
alongside `v1/` is all it takes to version the API.

### `internal/entity`

Core business models, usable in any layer and independent of storage or
transport. They may carry methods such as validation.

### `internal/usecase`

Business logic, grouped by area (`todo`, `user`) — one structure per group.
Repositories and caches are injected as interfaces (`contracts.go`), so the use
cases stay independent of concrete implementations.

### `internal/repo/persistent`

The repository: an abstract data store the business logic talks to. SQL lives in
`queries/`, sqlc generates the type-safe code in `sqlc/`, and the repository
maps generated database models to business entities.

### `internal/repo/cache`

An abstract cache the business logic talks to. It uses the adapter pattern to
sit in front of a concrete cache implementation (Otter, from `pkg/cache`).

### `pkg`

Reusable, project-agnostic packages: `cache` (Otter wrapper), `httpserver`,
`jwt`, `logger` (zerolog), `postgres` (pgx pool), `validatorx`, and small
`util`/`utils` helpers.

### `migrations`

golang-migrate SQL migration pairs (`*.up.sql` / `*.down.sql`).

### `frontend`

SvelteKit single-page app (static adapter) that consumes the API. Its TypeScript
client (`src/lib/types/api.ts`) is generated from `docs/swagger.yaml`.

## Dependency Injection

To keep business logic free of external dependencies, dependencies are injected
through `New` constructors as interfaces. The implementation behind an interface
can be swapped without touching the `usecase` package:

```go
package usecase

type Repository interface {
	Get()
}

type UseCase struct {
	repo Repository
}

func New(r Repository) *UseCase {
	return &UseCase{repo: r}
}

func (uc *UseCase) Do() {
	uc.repo.Get()
}
```

This also lets us auto-generate mocks (`make mock`, via
[go.uber.org/mock](https://go.uber.org/mock)) and write straightforward unit
tests.

> We aren't tied to specific implementations, so any component can be replaced
> with another. If the new component implements the interface, the business
> logic doesn't change.

## Clean Architecture

### The main principle

Dependency Inversion (from SOLID): dependencies point from the outer layer
inward, so business logic and entities stay independent of the rest of the
system. The application is split into two layers:

1. **Business logic** — entities and use cases (standard library only).
2. **Tools** — databases, HTTP servers, caches, and other frameworks.

![Clean Architecture](docs/img/layers-1.png)

**The inner layer** must stay clean: no imports from the outer layer, and all
calls outward happen through interfaces. Business logic knows nothing about
PostgreSQL specifically — only about an _abstract_ repository.

**The outer layer** components are unaware of each other. They communicate only
through the inner layer, always via interfaces, exchanging data in the shape the
business logic expects (`internal/entity`).

For example, an HTTP handler reaching the database goes through a use case:

```
    HTTP > usecase
           usecase > repository (Postgres)
           usecase < repository (Postgres)
    HTTP < usecase
```

The `>` and `<` symbols mark layer boundaries crossed through interfaces.

![Layers](docs/img/layers-2.png)

### Terminology

- **Entities** (`internal/entity`) — core business objects, agnostic of storage
  or transport, used throughout the app.
- **Database models** (`internal/repo/persistent/sqlc`) — generated code mirroring
  the schema, private to the repository layer.
- **Mapping** — the repository converts database models to business entities, so
  a schema change doesn't ripple into the business logic.
- **Use cases** (`internal/usecase`) — the application's business logic.
- **Controllers** (`internal/controller`) — the entry points (also called
  delivery, transport, or gateways elsewhere).

### Alternative approaches

Onion architecture and Hexagonal (Ports and Adapters) are close relatives, all
built on Dependency Inversion; the differences are mostly terminology.

## Make commands

Run `make help` for the full list. Common targets:

| Command            | Description                                       |
| ------------------ | ------------------------------------------------- |
| `make run`         | Generate code, migrate, and run the API           |
| `make build`       | Build the binary into `./main`                     |
| `make test`        | Run unit tests with race detector and coverage    |
| `make sqlc`        | Generate Go from SQL                              |
| `make swag-v1`     | Regenerate Swagger docs                           |
| `make mock`        | Regenerate test mocks                             |
| `make migrate-up`  | Apply migrations (`migrate-create name=...` etc.) |
| `make db-up`       | Start PostgreSQL via docker compose              |
| `make format`      | Format the codebase                              |

## Useful links

- [Project requirements](https://roadmap.sh/projects/todo-list-api)
- [The Clean Architecture article](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [The Twelve-Factor App](https://12factor.net/)
- Upstream template: [evrone/go-clean-template](https://github.com/evrone/go-clean-template)
