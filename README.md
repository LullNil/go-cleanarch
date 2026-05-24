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
- Versioned REST routes under `/v1`.
- Taskfile for local development.

## Structure

```text
.
├── cmd
│   ├── app                 # HTTP application entry point
│   └── migrator            # Database migration CLI
├── config                  # Configuration loading and local config
├── docs                    # OpenAPI and protobuf contracts
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
│   ├── lib                 # Shared internal helpers
│   ├── repository          # Infrastructure adapters
│   └── service             # Use cases / application services
├── migrations
├── docker-compose.yml
├── go.mod
└── Taskfile.yaml
```

## Requirements

- Go 1.25+
- Docker Engine and Docker Compose
- Optional: Taskfile

## Configuration

Create a local config from the example:

```bash
cp config/local.example.yaml config/local.yaml
```

Local config lives in `config/local.yaml`:

```yaml
env: "local"

http_server:
  port: ":8080"
  read_timeout: 30s
  write_timeout: 30s
  enable_swagger: true

grpc_server:
  port: ":9090"

postgres:
  dsn: "postgres://user:password123@localhost:5437/dbname?sslmode=disable"
  max_retries: 10
  retry_interval: 5s
  connect_timeout: 30s

redis:
  addr: "localhost:6379"
  password: ""
  db: 0
  ttl: 5m
  max_retries: 10
  retry_interval: 5s
  connect_timeout: 30s
```

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

## Notes

- Keep domain packages focused on entities, domain-specific errors, and repository ports.
- Keep handlers focused on transport-level work: decode input, read path/query params, call services, and map domain errors to protocol responses.
- Keep transport DTOs in delivery adapters such as `internal/delivery/http`; map responses through explicit DTO structs.
- Keep protocol-specific error mapping in delivery packages and shared application error categories in `internal/apperr`.
- Prefer concrete services and define small consumer-side interfaces where adapters need them.
- Keep `internal/app/services.go` limited to use case service wiring; external clients belong in `modules.go`, and technical providers can be split out when they appear.
- Keep cache ports in the service package and cache adapters under `internal/repository/<driver>`.
- Add new infrastructure implementations under `internal/repository/<driver>`.

## License

This project is licensed under the [MIT License](./LICENSE).
