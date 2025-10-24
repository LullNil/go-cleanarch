
---

# Go Clean Architecture Template

A minimal and extensible Go backend template following the principles of **Clean Architecture**.
It separates application layers to improve testability, scalability, and maintainability.

## Overview

The project is organized into a `server` directory, which contains all backend-related code.
Each domain entity (e.g. `entity1`) has its own folder with separate `entity`, `repository`, and `service` layers.

## Getting Started

```bash
# Start dependencies (PostgreSQL)
docker-compose up -d

# Run migrations
task migrate:up

# Start application
task server
```

## Features

* Layered structure with clear boundaries (domain, service, repository, delivery)
* Config-driven application setup
* PostgreSQL integration via `docker-compose`
* Taskfile for simplified local development
* Migration support via `cmd/migrator`

## Structure

```bash
.
├── server
    ├── README.md              # Server-specific documentation
    ├── Taskfile.yaml          # Task definitions for development
    ├── cmd                    # Application entry points
    │   ├── app
    │   │   └── main.go        # Main application entry
    │   └── migrator
    │       └── main.go        # Database migration tool
    ├── config
    │   ├── config.go          # Configuration loader
    │   └── local.yaml         # Local environment configuration
    ├── docker-compose.yml     # Local development services (DB, etc.)
    ├── domain
    │   └── entity1
    │       ├── entity.go      # Domain model
    │       ├── repository.go  # Repository interface
    │       └── service.go     # Domain service logic
    ├── internal
    │   ├── app
    │   │   ├── app.go         # Application initialization
    │   │   ├── modules.go     # Dependency wiring
    │   │   ├── router.go      # HTTP routing setup
    │   │   └── services.go    # Service initialization
    │   ├── delivery
    │   │   └── http
    │   │       └── entity1
    │   │           └── handler.go  # HTTP handler for entity1
    │   ├── lib
    │   │   └── logger
    │   │       └── pretty.go   # Pretty-printed structured logging
    │   ├── repository
    │   │   ├── postgres
    │   │   │   ├── entity1_repository.go  # PostgreSQL implementation
    │   │   │   └── postgres.go            # DB connection setup
    │   │   └── repository.go  # Common repository interfaces
    │   └── service
    │       └── entity1
    │           └── service.go  # Service implementation
    └── migrations
        ├── 1_create_entity1_table.up.sql
        └── 1_drop_entity1_table.down.sql
```

## TODO

- [ ] Add gRPC examples under `internal/delivery/grpc/`
- [ ] Add Redis repository example (caching, sessions, etc.)
- [ ] Add Kafka producer/consumer examples in `internal/repository/`
- [ ] Add Prometheus metrics and health endpoints
- [ ] Add unit and integration tests for core modules
- [ ] Add CI/CD workflow example using GitHub Actions


## License

This project is licensed under the [MIT License](./LICENSE).

---
