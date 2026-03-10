# golang-mastery-2025
A structured roadmap and collection of projects, notes, and exercises to achieve Go (Golang) mastery in 2025 — from fundamentals to advanced backend engineering.

## Project Structure

### Week 1
- **Daily Exercises**: Fundamental Go concepts and syntax practice
- **LeetCode Solutions**: Algorithm and data structure problems solved in Go
- **Task Manager CLI**: A complete command-line task management application

### Week 2
- **Simple API**: Basic HTTP server implementation with Go's standard library
- **Task API**: RESTful API with PostgreSQL, Docker, and advanced routing

### Week 3
- **Concurrency Basics**: Goroutines, channels, select, pipelines, and worker patterns
- **Connection Pool Test**: Database connection pooling behavior and concurrency testing

## Projects

### Task Manager CLI (`week1/task-manager/`)
A fully functional command-line task management application built with Go.

**Features:**
- Add new tasks
- List all tasks
- Mark tasks as complete
- Delete tasks
- Search tasks by description
- Persistent storage with JSON

**Usage:**
```bash
cd week1/task-manager
go run . add "Buy groceries"
go run . list
go run . complete 1
go run . delete 2
go run . search "buy"
```

### Task API (`week2/task-api/`)
A production-ready RESTful API for task management with PostgreSQL database.

**Features:**
- Full CRUD operations for tasks
- Search functionality
- PostgreSQL database with Docker support
- CORS and logging middleware
- Input validation
- Comprehensive test suite
- Complete API documentation

**Tech Stack:**
- Go 1.24.5
- PostgreSQL 15
- Gorilla Mux router
- Docker & Docker Compose

**Quick Start:**
```bash
cd week2/task-api
docker-compose up --build
```

**API Endpoints:**
- `GET /tasks` - List all tasks
- `POST /tasks` - Create new task
- `GET /tasks/{id}` - Get task by ID
- `PUT /tasks/{id}` - Mark task as complete
- `DELETE /tasks/{id}` - Delete task
- `GET /tasks?q=search` - Search tasks

**Documentation:**
- [README.md](week2/task-api/README.md) - Setup and usage guide
- [API.md](week2/task-api/API.md) - Complete API documentation
- [ARCHITECTURE.md](week2/task-api/ARCHITECTURE.md) - Architecture overview

### Concurrency Basics (`week3/concurrency-basics/`)
Hands-on Go concurrency exercises and patterns.

**Includes:**
- Basic goroutine usage
- Channel communication
- `select` pattern
- Pipeline pattern
- Worker pattern

**Run:**
```bash
cd week3/concurrency-basics
go run .
```

### Connection Pool Test (`week3/connection-pool-test/`)
Experiments for understanding DB connection pool behavior under concurrent load.

**Run:**
```bash
cd week3/connection-pool-test
go run .
```

## Learning Progress

### Week 1
- ✅ Go basics and syntax
- ✅ Structs and methods
- ✅ Error handling
- ✅ JSON marshaling/unmarshaling
- ✅ File I/O operations
- ✅ Command-line argument parsing
- ✅ Unit testing

### Week 2
- ✅ HTTP server fundamentals
- ✅ RESTful API design
- ✅ Third-party routing (Gorilla Mux)
- ✅ PostgreSQL integration
- ✅ Database connection management
- ✅ SQL queries and parameterization
- ✅ HTTP middleware patterns
- ✅ CORS handling
- ✅ Request/response logging
- ✅ Docker containerization
- ✅ Docker Compose orchestration
- ✅ Environment configuration
- ✅ API documentation
- ✅ Integration testing

### Week 3
- ✅ Goroutines and concurrency fundamentals
- ✅ Channels (buffered and unbuffered)
- ✅ `select` for multiplexing and cancellation
- ✅ Fan-in / fan-out and pipeline patterns
- ✅ Worker pool pattern
- ✅ Semaphore-style concurrency limiting
- ✅ Context propagation and timeout handling
- ✅ Graceful shutdown concepts
- ✅ Structured logging basics
- ✅ Connection pooling experiments
