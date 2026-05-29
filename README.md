# Go Clean Architecture Template

A minimal and extensible Go backend template following the principles of Clean Architecture.
It separates domain models, use cases, delivery adapters, and infrastructure adapters while keeping the project small enough to extend without ceremony.

## Features

- Clean dependency direction: delivery and infrastructure depend inward, domain stays independent.
- Domain-first repository ports per entity.
- Concrete application services with consumer-side interfaces in adapters.
- HTTP DTOs isolated in delivery packages.
- gRPC example with versioned protobuf contracts.
- PostgreSQL repository example with migrations.
- Redis cache example wired into the entity service.
- Replaceable external integration client example under `internal/integration`.
- Versioned REST routes under `/v1`.
- Taskfile for local development.
- GitHub Actions CI for formatting, tests, and module tidy checks.

## Stack

- Go 1.25+
- HTTP: standard `net/http`, Chi router, CORS middleware
- gRPC: `google.golang.org/grpc`, protobuf contracts under `docs/proto`
- Database: PostgreSQL with SQL migrations
- Cache: Redis via `github.com/redis/go-redis/v9`
- Config: YAML via `cleanenv`, optional `.env` for `CONFIG_PATH`
- Logging: standard `log/slog`
- API docs: static OpenAPI spec and Swagger UI

## Structure

```text
.
├── cmd
│   ├── app                 # HTTP application entry point
│   └── migrator            # Database migration CLI
├── config                  # Configuration loading and local config
├── docs                    # OpenAPI and protobuf contracts
├── .github/workflows       # CI workflows
├── domain
│   ├── errors.go           # Shared domain errors
│   └── entity1
│       ├── entity.go       # Business entity
│       └── repository.go   # Repository port for this domain
├── internal
│   ├── app                 # Composition root and lifecycle orchestration
│   ├── delivery
│   │   ├── grpc            # gRPC server, handlers, generated pb code
│   │   └── http            # HTTP server, handlers, routes, transport DTOs
│   ├── integration         # External service clients
│   ├── lib                 # Shared internal helpers
│   ├── repository          # Infrastructure adapters
│   └── service             # Use cases / application services
├── migrations
├── Dockerfile
├── docker-compose.yml      # Local PostgreSQL and Redis
├── go.mod
└── Taskfile.yaml
```

## Requirements

- Go 1.25+
- Docker Engine and Docker Compose
- golangci-lint for local linting
- Optional: Taskfile

Install local development tools:

```bash
task tools:install
```

By default this installs the latest `golangci-lint` into `$(go env GOPATH)/bin`. Make sure that directory is in your `PATH`.

Pin a version when needed:

```bash
task tools:install GOLANGCI_LINT_VERSION=v2.12.2
```

## Configuration

Create a local config from the example:

```bash
cp config/local.example.yaml config/local.yaml
```

Local config lives in `config/local.yaml`

You can pass a custom config path with:

```bash
go run ./cmd/app --config=./config/local.yaml
```

or set:

```bash
CONFIG_PATH=./config/local.yaml go run ./cmd/app
```

## Run Locally

```bash
docker compose up -d
task migrate:up
task server
```

Without Taskfile:

```bash
go run ./cmd/migrator --database-dsn "postgres://user:password123@localhost:5437/dbname?sslmode=disable" --migrations-path ./migrations --command up
go run ./cmd/app --config=./config/local.yaml
```

## Docker

Build the application image:

```bash
task docker:build
```

Without Taskfile:

```bash
docker build -t go-cleanarch:local .
```

Run PostgreSQL and Redis, then start the application container with the Docker example config:

```bash
docker compose up -d
docker run --rm \
  --network go-cleanarch_default \
  -p 8080:8080 \
  -p 9090:9090 \
  -e CONFIG_PATH=/app/config/docker.yaml \
  -v "$PWD/config/docker.example.yaml:/app/config/docker.yaml:ro" \
  go-cleanarch:local
```

The image contains the application binary and OpenAPI assets. Runtime configuration is supplied from outside the image.

## Example Routes

```text
POST   /v1/entity1
GET    /v1/entity1/{id}
PUT    /v1/entity1/{id}
DELETE /v1/entity1/{id}
```

Swagger UI is available locally when `http_server.enable_swagger` is `true`:

```text
GET /swagger
GET /swagger/openapi.yaml
```

## gRPC

The versioned protobuf contract lives in `docs/proto/v1/entity1.proto`.

Generated Go bindings are placed under `internal/delivery/grpc/pb` and are used only by the gRPC delivery adapter.

Generate protobuf and gRPC bindings with Taskfile:

```bash
task proto:gen
```

Or run `protoc` directly:

```bash
PATH="$HOME/go/bin:$PATH" protoc \
  --go_out=. --go_opt=module=github.com/LullNil/go-cleanarch \
  --go-grpc_out=. --go-grpc_opt=module=github.com/LullNil/go-cleanarch \
  docs/proto/v1/entity1.proto
```

The app starts both HTTP and gRPC servers from the same composition root.

## External Integrations

External service clients live under `internal/integration/<service>`.

The example auth gRPC client is in `internal/integration/auth/grpc_client.go`. Application services depend on the service-owned port, not on protobuf or gRPC types:

```text
internal/service/entity1/auth.go        # AuthClient interface
internal/integration/auth/grpc_client.go # gRPC adapter
internal/app/integrations.go             # external client lifecycle
internal/app/services.go                 # dependency injection
```

Leave `integrations.auth.grpc_target` empty to disable the example client. When adding a real auth contract, keep generated protobuf types inside the integration adapter and map them to the service interface.

## CI

CI provider files are intentionally thin. The reusable checks live in `Taskfile.yaml` so GitHub Actions, GitLab CI, and local runs use the same commands.

GitHub Actions is enabled by default in `.github/workflows/ci.yml`. It runs on pushes to `main` and on pull requests:

```text
task ci
task docker:build
```

The workflow uses `go.mod` as the source of truth for the Go version.

A GitLab CI template is available at `docs/ci/gitlab-ci.yml`. To enable GitLab, copy that template to `.gitlab-ci.yml` in the repository root. GitLab requires the file at the root by default, while GitHub requires active workflows under `.github/workflows`.

Run the same checks locally with Taskfile:

```bash
task tools:install
task ci
task docker:build
```

Or run the individual checks:

```bash
task fmt:check
task lint
task test
task mod:tidy:check
```

Go linting is configured in `.golangci.yml`. The default lint gate enables `govet`, `staticcheck`, `errcheck`, `revive`, and a small set of bug-focused checks. Add future checks as separate Taskfile tasks or explicit linters so each failure remains easy to diagnose.

## Notes

- Keep domain packages focused on entities, domain-specific errors, and repository ports.
- Keep handlers focused on transport-level work: decode input, read path/query params, call services, and map domain errors to protocol responses.
- Keep transport DTOs in delivery adapters such as `internal/delivery/http`; map responses through explicit DTO structs.
- Keep protocol-specific error mapping in delivery packages and shared application error categories in `internal/apperr`.
- Prefer concrete services and define small consumer-side interfaces where adapters need them.
- Keep `internal/app/modules.go` limited to low-level technical resources such as databases, caches, brokers, and metrics.
- Keep `internal/app/integrations.go` limited to external service client lifecycle and wiring.
- Keep `internal/app/services.go` limited to use case service wiring.
- Keep cache ports in the service package and cache adapters under `internal/repository/<driver>`.
- Keep external service ports in the service package and concrete clients under `internal/integration/<service>`.
- Add new infrastructure implementations under `internal/repository/<driver>`.
- The `entity1` CRUD example is intentionally abstract. Replace it with real business entities instead of building new services around the placeholder name.
- Service unit tests use small hand-written fakes to show how business logic can be tested without infrastructure.

## License

This project is licensed under the [MIT License](./LICENSE).
