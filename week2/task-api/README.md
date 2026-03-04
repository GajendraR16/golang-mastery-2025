# Task API

A RESTful API for managing tasks built with Go, PostgreSQL, and Docker.

## Features

- Create, read, update, and delete tasks
- Mark tasks as complete
- Search tasks by description
- PostgreSQL database with Docker support
- CORS and logging middleware
- Input validation

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
└── dockerfile
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
docker-compose up --build
```

The API will be available at `http://localhost:8080`

### Running Locally

1. Start PostgreSQL:
```bash
# Using Docker
docker run --name taskdb -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=taskdb -p 5432:5432 -d postgres:15
```
If port `5432` is already used on your machine, map to another host port (example `5434:5432`) and update `DATABASE_URL` accordingly.

2. Initialize the database:
```bash
docker exec taskdb pg_isready -U postgres
docker exec -i taskdb psql -U postgres -d taskdb < schema.sql
```

3. Set the database connection string (optional):
```bash
export DATABASE_URL="postgres://postgres:postgres@localhost:5432/taskdb?sslmode=disable"
```

4. Install dependencies:
```bash
go mod download
```

5. Run the application:
```bash
go run main.go
```

## Running Tests
1. Start PostgreSQL:
```bash
# Using Docker
docker run --name task-db-test  -e POSTGRES_USER=postgres  -e POSTGRES_PASSWORD=postgres  -e POSTGRES_DB=taskdb_test  -p 5433:5432 -d postgres
```

2. Initialize the database:
```bash
docker exec task-db-test pg_isready -U postgres
docker exec -i task-db-test psql -U postgres -d taskdb_test < schema.sql
```

3. Run all tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific test file
go test -v ./handler/
```

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `DATABASE_URL` | PostgreSQL connection string | `postgres://postgres:postgres@localhost:5432/taskdb?sslmode=disable` |

## API Documentation

See [API.md](API.md) for detailed endpoint documentation.

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
