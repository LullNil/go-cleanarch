# Agent Guidelines

This repository is a minimal Go backend template based on Clean Architecture.
Keep changes small, explicit, and consistent with the existing package boundaries.

## Architecture Rules

- `domain/` contains business entities, domain-specific errors, and repository ports.
- `internal/service/` contains use cases and application commands.
- `internal/delivery/http/` contains HTTP servers, handlers, routes, transport DTOs, and HTTP error mapping.
- `internal/delivery/grpc/` contains gRPC servers, handlers, generated protobuf code, and gRPC error mapping.
- `internal/repository/<driver>/` contains infrastructure adapters that implement domain repository ports or service cache ports.
- `internal/integration/<service>/` contains external service clients that implement service-owned ports.
- `internal/app/` is the composition root and lifecycle orchestrator. Do not put routes or transport-specific logic there.
- `internal/lib/` contains reusable technical helpers that do not know business rules.

Dependencies must point inward:

```text
delivery -> service -> domain
repository -> domain
integration -> domain
app -> delivery/service/repository/integration
```

`domain` must not import `internal`.

## Adding New Business Logic

For a new entity or feature:

1. Add the domain entity under `domain/<entity>/`.
2. Add repository ports under `domain/<entity>/repository.go` when persistence is needed.
3. Add use case commands under `internal/service/<entity>/command.go`.
4. Add use case methods in `internal/service/<entity>/service.go`.
5. Add cache ports under `internal/service/<entity>/cache.go` when cache is needed.
6. Add external service ports under `internal/service/<entity>/` when outside calls are needed.
7. Add repository/cache implementation under `internal/repository/<driver>/`.
8. Add external client implementations under `internal/integration/<service>/`.
9. Wire technical resources in `internal/app/modules.go`.
10. Wire external clients in `internal/app/integrations.go`.
11. Wire repositories, caches, integrations, and services in `internal/app/services.go`.
12. Add delivery handlers:
   - HTTP DTOs and handlers under `internal/delivery/http/<entity>/`.
   - gRPC contract under `docs/proto/v1/` and handler under `internal/delivery/grpc/<entity>/`.
13. Register routes/services in the relevant delivery server package.

Handlers should only do transport-level work: decode input, read path/query params or metadata, call services, and map application errors to protocol responses.

## DTOs And Commands

- HTTP request/response DTOs stay in `internal/delivery/http/<entity>/dto.go`.
- gRPC request/response types come from generated protobuf code in `internal/delivery/grpc/pb`.
- Service input types are commands, not DTOs, and live in `internal/service/<entity>/command.go`.
- Do not import HTTP DTOs into services or gRPC handlers.
- Do not import generated protobuf or concrete gRPC clients into services; services depend on small interfaces they own.

## Errors

- Domain-level categories live in `domain/errors.go`.
- Application error classification lives in `internal/apperr`.
- HTTP error mapping lives in `internal/delivery/http/httperror`.
- gRPC error mapping lives in `internal/delivery/grpc/grpcerror`.
- Repositories should translate known infrastructure errors into domain errors and return unknown errors wrapped with context.
- Log errors once at the delivery boundary. Avoid duplicate logging in repository and service layers.
- Best-effort cache failures may be logged in services only when the error is intentionally swallowed.

## Logging

Use `log/slog`.

HTTP error logs should include:

- `request_id`
- `trace_id`

gRPC error logs should include:

- `method`
- `code`
- `request_id`
- `trace_id`

## Comments And Naming

- Exported types, functions, methods, vars, and consts must have Go-style comments.
- Exported comments start with the identifier name and end with a period.
- Unexported comments are optional; when added, keep them short and omit the final period.
- Prefer Go names over framework names: `Command` for service input, `DTO` only for transport shapes.
- Avoid package names like `common`, `utils`, or `contracts`.

## gRPC Generation

Proto contracts live under `docs/proto/v1/`.

Generated Go code goes under `internal/delivery/grpc/pb`.

Use:

```bash
PATH="$HOME/go/bin:$PATH" protoc \
  --go_out=. --go_opt=module=github.com/LullNil/go-cleanarch \
  --go-grpc_out=. --go-grpc_opt=module=github.com/LullNil/go-cleanarch \
  docs/proto/v1/entity1.proto
```

## Verification

Before finishing code changes, run:

```bash
gofmt -w <changed-go-paths>
go test ./...
```

If dependencies changed, run:

```bash
go mod tidy
```

## Tests

- Prefer unit tests for service/use case behavior.
- Keep small hand-written fakes in `fake_test.go` next to the tested package.
- Keep test cases in `service_test.go` or another focused `*_test.go` file.
- Use generated mocks only when interfaces become too large for readable fakes.
