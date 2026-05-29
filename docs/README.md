# Documentation

Project documentation lives here.

- `openapi/` contains Swagger/OpenAPI documentation for HTTP APIs.
- `proto/` contains versioned gRPC protobuf contracts.

CI is defined in `.github/workflows/ci.yml` and checks Go formatting, `go test ./...`, and `go mod tidy`.
