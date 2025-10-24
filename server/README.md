# Server Building

## Requirements

- Go 1.25+
- [Docker](https://docs.docker.com/get-started/get-docker/) (Docker Engine and Docker Compose)
- *(Optional)* [Taskfile](https://taskfile.dev) for task automation

---

### 0. Create Local Configuration

Before running the application, you need to create a `local.yaml` configuration file in the `config` directory with the following contents:

```yaml
env: "local" # local, dev, prod

http_server:
  port: ":8080"
  read_timeout: 30s
  write_timeout: 30s

postgres:
  dsn: "postgres://user:password123@localhost:5437/dbname?sslmode=disable"
  max_retries: 10
  retry_interval: 5s
  connect_timeout: 30s
```

This file is used for local development.
For production, `prod.yaml` is created with the actual secrets and production settings.

### 1. Start PostgreSQL Database

Before running the application, you need to start the PostgreSQL database using Docker Compose. Navigate to the project root directory and run:

```bash
docker compose up -d
```

This command will:
- Download the `postgres:15-alpine` Docker image (if not already present).
- Create and start a PostgreSQL container named `db` (or similar, based on your project name).
- Map the container's port `5437` to your host's port `5437`.
- Initialize the database with the specified `POSTGRES_DB`, `POSTGRES_USER`, and `POSTGRES_PASSWORD` from `docker-compose.yml`.

**To check the status of your Docker containers:**

```bash
docker-compose ps
```

**To stop the PostgreSQL container:**

```bash
docker-compose -down
```

**To stop and remove the PostgreSQL container and its data volume (start fresh):**

```bash
docker-compose down -v
```

### 2. Database Migrations

After starting the PostgreSQL container, you need to apply the database migrations to set up the schema.

Install dependencies:

```bash
go mod tidy
```

Then, run the migrations using Taskfile:

```bash
# Apply all UP migrations
task migrate:up

# Rollback the last DOWN migration
task migrate:down
```

**Manual Migration Command (if not using Taskfile):**

If you prefer to run migrations manually without Taskfile, you can use the following command (replace with your actual DSN):

```bash
go run ./cmd/migrator --database-dsn "postgres://user:password123@localhost:5437/dbname?sslmode=disable" --migrations-path ./migrations --command up
```

**To go into psql run:**

```bash
docker exec -it postgres_database psql -U user -d dbname
```

### 3. Run the Server

```bash
# Run the server using Taskfile
task server

# Or manually:
go run ./cmd/app/main.go --config=./config/local.yaml
```
