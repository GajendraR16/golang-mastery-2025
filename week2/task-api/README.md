# Task API

A RESTful API for managing tasks built with Go, PostgreSQL, and Docker.

## Features

- Create, read, update, and delete tasks
- Mark tasks as complete
- Search tasks by description
- Batch create operation for multiple tasks
- Context-aware database operations
- PostgreSQL database with Docker support
- CORS and logging middleware
- Input validation
- Comprehensive test suite
- Makefile for development workflow

## Tech Stack

- **Language**: Go 1.24.5
- **Database**: PostgreSQL 15
- **Router**: Gorilla Mux
- **Containerization**: Docker & Docker Compose

## Project Structure

```
task-api/
├── config/          # Configuration files
├── handler/         # HTTP request handlers
├── middleware/      # HTTP middleware (CORS, logging)
├── models/          # Data models
├── storage/         # Database layer
├── main.go          # Application entry point
├── schema.sql       # Database schema
├── docker-compose.yml
└── Dockerfile
```

## Prerequisites

- Go 1.24.5 or higher
- Docker and Docker Compose
- PostgreSQL 15 (if running locally without Docker)

## Quick Start

### Using Docker Compose (Recommended)

1. Clone the repository and navigate to the project directory:
```bash
cd week2/task-api
```

2. Start the application and database:
```bash
docker compose up --build
```

The API will be available at `http://localhost:8080`

### Using Makefile (Recommended)

```bash
cd week2/task-api
make up
```

Useful targets:

```bash
make down
make logs
make local-db-up
make local-db-init
make run
make test-db-up
make test-db-init
make test
make clean
```

### Running Locally (With Makefile)

1. Start local Postgres for the app:
```bash
make local-db-up
make local-db-init
```

2. Set connection string and run API:
```bash
export DATABASE_URL="postgres://postgres:postgres@localhost:5434/taskdb?sslmode=disable"
export PORT=8080
make run
```

3. Cleanup:
```bash
make local-db-down
```

### Running Locally (Manual Sequence)

1. Check what is using Postgres default port:
```bash
docker ps --format 'table {{.Names}}\t{{.Ports}}' | rg 5432
ss -ltn 'sport = :5432'
```

2. Start Postgres on an available host port (example `5434`):
```bash
docker rm -f taskdb 2>/dev/null || true
docker run --name taskdb -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=taskdb -p 5434:5432 -d postgres:15
```

3. Wait for DB and initialize schema:
```bash
docker exec taskdb pg_isready -U postgres -d taskdb
docker exec -i taskdb psql -U postgres -d taskdb < schema.sql
```

4. Set connection string and run API:
```bash
export DATABASE_URL="postgres://postgres:postgres@localhost:5434/taskdb?sslmode=disable"
export PORT=8080
go run main.go
```

5. Cleanup:
```bash
docker rm -f taskdb
```

## Running Tests
1. Start PostgreSQL for tests:
```bash
make test-db-up
```

2. Initialize the database:
```bash
make test-db-init
```

3. Run all tests

```bash
make test

# Optional
make test-cover
make test-handler
```

4. Cleanup:
```bash
make test-db-down
```

One-command flow (setup + test + cleanup):
```bash
make test-with-db
```

If Docker requires sudo on your system, avoid `sudo make` and run:
```bash
make test-with-db DOCKER="sudo docker"
```

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `DATABASE_URL` | PostgreSQL connection string | `postgres://postgres:postgres@localhost:5432/taskdb?sslmode=disable` |

## API Documentation

See [API.md](API.md) for detailed endpoint documentation.

### Batch Endpoint Scope

- `POST /tasks/batch` currently supports batch **create only**
- Batch **update** (`PUT`) and batch **delete** (`DELETE`) are not implemented in this version

## Database Schema

```sql
CREATE TABLE IF NOT EXISTS tasks (
    id SERIAL PRIMARY KEY,
    description TEXT NOT NULL,
    completed BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMP
);
```

## Development

### Adding New Endpoints

1. Define the handler in `handler/task.go`
2. Register the route in `server/router.go`
3. Add tests in `handler/handlers_test.go`

### Code Organization

- **handlers**: HTTP request/response logic
- **storage**: Database operations
- **models**: Data structures and business logic
- **middleware**: Cross-cutting concerns (logging, CORS)

## Troubleshooting

- `container name "/taskdb" is already in use`:
```bash
docker rm -f taskdb
```

- `failed to bind host port 0.0.0.0:5432 ... address already in use`:
```bash
docker run --name taskdb -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=taskdb -p 5434:5432 -d postgres:15
export DATABASE_URL="postgres://postgres:postgres@localhost:5434/taskdb?sslmode=disable"
```

- `dial tcp 127.0.0.1:5434: connect: connection refused`:
```bash
docker ps --format 'table {{.Names}}\t{{.Status}}\t{{.Ports}}' | rg 5434
docker start taskdb
docker exec taskdb pg_isready -U postgres -d taskdb
```

- `container ... is not running`:
```bash
docker start task-db-test
```

- `psql: ... No such file or directory` right after starting container:
```bash
docker exec task-db-test pg_isready -U postgres
docker exec -i task-db-test psql -U postgres -d taskdb_test < schema.sql
```

## License

MIT
