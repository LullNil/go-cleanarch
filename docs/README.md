# Documentation

Project documentation lives here.

- `openapi/` contains Swagger/OpenAPI documentation for HTTP APIs.
- `proto/` contains versioned gRPC protobuf contracts.
- `ci/` contains CI provider templates that are not active by default.

GitHub CI is active in `.github/workflows/ci.yml`; GitLab CI is available as an opt-in template in `ci/gitlab-ci.yml`.

External integration clients live under `internal/integration/<service>` and should implement service-owned interfaces.
