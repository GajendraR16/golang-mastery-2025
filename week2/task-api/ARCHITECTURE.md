# Task API Architecture

## Overview

This is a RESTful API built with Go following clean architecture principles, separating concerns into distinct layers.

## Architecture Diagram

```
┌─────────────────────────────────────────────────────────────┐
│                         Client Layer                        │
│                  (Browser, cURL, Postman)                   │
└────────────────────────┬────────────────────────────────────┘
                         │ HTTP Requests
                         ▼
┌─────────────────────────────────────────────────────────────┐
│                      Middleware Layer                       │
│  ┌──────────────────┐         ┌──────────────────┐          │
│  │ CORS Middleware  │────────▶│Logging Middleware│          │
│  └──────────────────┘         └──────────────────┘          │
└────────────────────────┬────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│                       Router Layer                          │
│                     (Gorilla Mux)                           │
│  ┌─────────────────────────────────────────────────────┐    │
│  │  Route Matching & Parameter Extraction              │    │
│  │  - /tasks                                           │    │
│  │  - /tasks/{id}                                      │    │
│  │  - /tasks?q={query}                                 │    │
│  └─────────────────────────────────────────────────────┘    │
└────────────────────────┬────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│                      Handler Layer                          │
│                    (handler/task.go)                        │
│  ┌──────────────┐   ┌──────────────┐   ┌──────────────┐     │
│  │ TaskHandler  │   │CreateHandler │   │DeleteHandler │     │
│  └──────────────┘   └──────────────┘   └──────────────┘     │
│  ┌───────────────┐  ┌───────────────┐  ┌──────────────┐     │
│  │TaskHandlerById│  │CompleteHandler│  │SearchHandler │     │
│  └───────────────┘  └───────────────┘  └──────────────┘     │
│                                                             │
│  Responsibilities:                                          │
│  - Request validation                                       │
│  - JSON encoding/decoding                                   │
│  - HTTP response formatting                                 │
│  - Error handling                                           │
└────────────────────────┬────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│                       Models Layer                          │
│                    (models/task.go)                         │
│  ┌─────────────────────────────────────────────────────┐    │
│  │  Task Model                                         │    │
│  │  - ID, Description, Completed                       │    │
│  │  - CreatedAt, CompletedAt                           │    │
│  │  - Business logic methods (Complete, String)        │    │
│  └─────────────────────────────────────────────────────┘    │
│  ┌─────────────────────────────────────────────────────┐    │
│  │  TaskData (DTO for creation)                        │    │
│  └─────────────────────────────────────────────────────┘    │
└────────────────────────┬────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│                      Storage Layer                          │
│                  (storage/postgres.go)                      │
│  ┌─────────────────────────────────────────────────────┐    │
│  │  PostgresStore                                      │    │
│  │  - GetAllTasks()                                    │    │
│  │  - GetTaskById(id)                                  │    │
│  │  - CreateTask(description)                          │    │
│  │  - CompletedTaskById(id)                            │    │
│  │  - DeleteTaskById(id)                               │    │
│  └─────────────────────────────────────────────────────┘    │
│                                                             │
│  Responsibilities:                                          │
│  - Database connection management                           │
│  - SQL query execution                                      │
│  - Data persistence                                         │
└────────────────────────┬────────────────────────────────────┘
                         │
                         ▼
┌────────────────────────────────────────────────────────────┐
│                      Database Layer                        │
│                      PostgreSQL 15                         │
│  ┌─────────────────────────────────────────────────────┐   │
│  │  tasks table                                        │   │
│  │  - id (SERIAL PRIMARY KEY)                          │   │
│  │  - description (TEXT)                               │   │
│  │  - completed (BOOLEAN)                              │   │
│  │  - created_at (TIMESTAMP)                           │   │
│  │  - completed_at (TIMESTAMP)                         │   │
│  └─────────────────────────────────────────────────────┘   │
└────────────────────────────────────────────────────────────┘
```

## Layer Responsibilities

### 1. Middleware Layer
**Location**: `middleware/`

**Purpose**: Cross-cutting concerns that apply to all requests

**Components**:
- **CORS Middleware**: Handles cross-origin resource sharing
  - Allows all origins
  - Supports GET, POST, PUT, DELETE methods
  - Handles preflight OPTIONS requests

- **Logging Middleware**: Request/response logging
  - Logs HTTP method and path
  - Logs response status code
  - Tracks request duration

### 2. Router Layer
**Location**: `main.go`

**Purpose**: Route HTTP requests to appropriate handlers

**Technology**: Gorilla Mux

**Features**:
- Path parameter extraction (`/tasks/{id:[0-9]+}`)
- Query parameter handling (`/tasks?q={query}`)
- HTTP method routing
- Regex validation for path parameters

### 3. Handler Layer
**Location**: `handler/`

**Purpose**: HTTP request/response handling

**Key Functions**:
- `TaskHandler`: Get all tasks
- `TaskHandlerById`: Get single task
- `CreateHandler`: Create new task
- `TaskCompleteHandler`: Mark task complete
- `DeleteHandler`: Delete task
- `SearchHandler`: Search tasks

**Responsibilities**:
- Parse request body/parameters
- Validate input
- Call storage layer
- Format JSON responses
- Handle errors with appropriate status codes

**Helper Functions**:
- `jsonError()`: Standardized error responses
- `jsonHandler()`: Standardized success responses

### 4. Models Layer
**Location**: `models/`

**Purpose**: Data structures and business logic

**Components**:
- **Task**: Main domain model
  - Fields: ID, Description, Completed, CreatedAt, CompletedAt
  - Methods: `Complete()`, `String()`, `NewTask()`
  
- **TaskData**: Data transfer object for task creation
  - Validation tags for input validation

**Design Pattern**: Domain-Driven Design (DDD)

### 5. Storage Layer
**Location**: `storage/`

**Purpose**: Data persistence abstraction

**Components**:
- **PostgresStore**: Database operations
  - Connection management
  - CRUD operations
  - Query execution

**Design Pattern**: Repository Pattern

**Benefits**:
- Database implementation can be swapped
- Easy to mock for testing
- Separates SQL from business logic

### 6. Database Layer
**Technology**: PostgreSQL 15

**Schema**: Single `tasks` table with proper indexing

**Features**:
- Auto-incrementing primary key
- Timestamp tracking
- Boolean completion status

## Data Flow

### Example: Creating a Task

```
1. Client sends POST /tasks with JSON body
   ↓
2. CORS middleware adds headers
   ↓
3. Logging middleware starts timer
   ↓
4. Router matches POST /tasks → CreateHandler
   ↓
5. CreateHandler:
   - Decodes JSON to TaskData
   - Validates description (length, empty check)
   - Trims whitespace
   ↓
6. Storage layer:
   - Executes INSERT query
   - Returns created Task with ID
   ↓
7. CreateHandler:
   - Encodes Task to JSON
   - Returns 201 Created
   ↓
8. Logging middleware logs request
   ↓
9. Client receives response
```

## Design Patterns

### 1. Repository Pattern
The storage layer abstracts database operations, making it easy to:
- Switch databases
- Mock for testing
- Maintain clean separation of concerns

### 2. Dependency Injection
The `App` struct receives its dependencies:
```go
type App struct {
    Store *storage.PostgresStore
}
```

Benefits:
- Testability
- Flexibility
- Loose coupling

### 3. Middleware Chain
Middleware functions wrap handlers:
```go
router.Use(middleware.LoggingMiddleware)
router.Use(middleware.CorsMiddleware)
```

### 4. Factory Pattern
Task creation uses a constructor:
```go
func NewTask(id int, description string) *Task
```

## Error Handling Strategy

### Layered Error Handling

1. **Storage Layer**: Returns Go errors
2. **Handler Layer**: Converts to HTTP status codes
3. **Client**: Receives JSON error responses

### Error Response Format
```json
{
  "error": "Human-readable error message"
}
```

### Status Code Mapping
- 400: Client input errors
- 404: Resource not found
- 500: Server/database errors

## Testing Strategy

### Unit Tests
- Handler tests with mock storage
- Model method tests
- Validation tests

### Integration Tests
- Full request/response cycle
- Database interactions
- Test database isolation

### Test Database
Separate `taskdb_test` database for testing

## Scalability Considerations

### Current Architecture
- Single instance
- Direct database connection
- Synchronous request handling

### Future Improvements
1. **Connection Pooling**: Already supported by `lib/pq`
2. **Caching**: Add Redis for frequently accessed tasks
3. **Rate Limiting**: Protect against abuse
4. **Authentication**: Add JWT middleware
5. **Pagination**: For large task lists
6. **Background Jobs**: For async operations
7. **Metrics**: Prometheus integration
8. **Tracing**: OpenTelemetry support

## Security Considerations

### Current Implementation
- SQL injection protection (parameterized queries)
- Input validation
- CORS configuration

### Recommended Additions
- Authentication/Authorization
- Rate limiting
- Request size limits
- HTTPS enforcement
- API versioning
- Input sanitization

## Deployment

### Docker Compose
- Multi-container setup
- Health checks
- Volume persistence
- Environment configuration

### Container Architecture
```
┌──────────────┐      ┌──────────────┐
│   API        │─────▶│  PostgreSQL  │
│  Container   │      │   Container  │
│  (Port 8080) │      │  (Port 5432) │
└──────────────┘      └──────────────┘
```

## Configuration Management

### Environment Variables
- `DATABASE_URL`: Database connection string
- Fallback to local defaults for development

### Configuration Files
- `docker-compose.yml`: Container orchestration
- `schema.sql`: Database initialization
- `go.mod`: Dependency management

## Monitoring & Observability

### Current Logging
- Request/response logging via middleware
- Console output

### Recommended Additions
- Structured logging (JSON format)
- Log aggregation (ELK stack)
- Application metrics
- Health check endpoint
- Readiness/liveness probes

## Performance Characteristics

### Strengths
- Fast Go runtime
- Efficient JSON encoding
- Connection pooling
- Minimal dependencies

### Bottlenecks
- Database queries (single connection)
- No caching layer
- Synchronous processing

### Optimization Opportunities
- Add database indexes
- Implement caching
- Use prepared statements
- Add query optimization
