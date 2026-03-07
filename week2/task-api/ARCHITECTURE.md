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

## Concurrency Patterns

### Context-Aware Operations

The API implements context-aware database operations for better resource management and cancellation support:

```go
func (s *PostgresStore) CreateTask(ctx context.Context, description string) (*Task, error) {
    query := `INSERT INTO tasks (description) VALUES ($1) RETURNING ...`
    err := s.db.QueryRowContext(ctx, query, description).Scan(...)
    return &task, err
}
```

**Benefits**:
- Request cancellation propagation
- Timeout handling
- Resource cleanup
- Distributed tracing support

### Batch Create with Goroutines

Batch create uses goroutines in the handler for concurrent processing:

```go
func (app *App) processBatch(ctx context.Context, tasks []models.TaskData) BatchResponse {
    var wg sync.WaitGroup
    sem := make(chan struct{}, 10)

    for _, task := range tasks {
        wg.Add(1)
        go func(t models.TaskData) {
            defer wg.Done()
            sem <- struct{}{}
            defer func() { <-sem }()

            // validate + app.Store.CreateTask(ctx, ...)
        }(task)
    }
    wg.Wait()
    return response
}
```

### Worker Pool Pattern

For high-throughput scenarios, the API can implement worker pools:

```go
type TaskProcessor struct {
    workers   int
    taskQueue chan TaskJob
    results   chan TaskResult
    wg        sync.WaitGroup
}

func NewTaskProcessor(workers int) *TaskProcessor {
    return &TaskProcessor{
        workers:   workers,
        taskQueue: make(chan TaskJob, workers*2),
        results:   make(chan TaskResult, workers*2),
    }
}

func (tp *TaskProcessor) Start(ctx context.Context) {
    for i := 0; i < tp.workers; i++ {
        tp.wg.Add(1)
        go tp.worker(ctx)
    }
}

func (tp *TaskProcessor) worker(ctx context.Context) {
    defer tp.wg.Done()
    for {
        select {
        case job := <-tp.taskQueue:
            result := tp.processTask(ctx, job)
            tp.results <- result
        case <-ctx.Done():
            return
        }
    }
}
```

### Database Connection Pooling

PostgreSQL connections are managed with built-in pooling:

```go
func NewPostgresStore(connStr string) (*PostgresStore, error) {
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        return nil, err
    }
    
    // Configure connection pool
    db.SetMaxOpenConns(25)                 // Maximum open connections
    db.SetMaxIdleConns(5)                  // Maximum idle connections
    db.SetConnMaxLifetime(5 * time.Minute) // Connection lifetime
    
    return &PostgresStore{db: db}, nil
}
```

### Graceful Shutdown

The server implements graceful shutdown with context cancellation:

```go
func main() {
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()
    
    server := &http.Server{
        Addr:    ":8080",
        Handler: router,
    }
    
    // Start server in goroutine
    go func() {
        if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("Server failed: %v", err)
        }
    }()
    
    // Wait for interrupt signal
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
    <-sigChan
    
    // Graceful shutdown with timeout
    shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer shutdownCancel()
    
    if err := server.Shutdown(shutdownCtx); err != nil {
        log.Printf("Server shutdown error: %v", err)
    }
}
```

### Rate Limiting with Channels

Implement rate limiting using channel-based semaphores:

```go
type RateLimiter struct {
    semaphore chan struct{}
    rate      time.Duration
}

func NewRateLimiter(maxConcurrent int, rate time.Duration) *RateLimiter {
    return &RateLimiter{
        semaphore: make(chan struct{}, maxConcurrent),
        rate:      rate,
    }
}

func (rl *RateLimiter) Allow(ctx context.Context) error {
    select {
    case rl.semaphore <- struct{}{}:
        go func() {
            time.Sleep(rl.rate)
            <-rl.semaphore
        }()
        return nil
    case <-ctx.Done():
        return ctx.Err()
    }
}
```

### Concurrent Safety Patterns

**Mutex for Shared State**:
```go
type SafeCounter struct {
    mu    sync.RWMutex
    count map[string]int
}

func (sc *SafeCounter) Increment(key string) {
    sc.mu.Lock()
    defer sc.mu.Unlock()
    sc.count[key]++
}

func (sc *SafeCounter) Get(key string) int {
    sc.mu.RLock()
    defer sc.mu.RUnlock()
    return sc.count[key]
}
```

**Channel-Based Communication**:
```go
type TaskNotifier struct {
    subscribers []chan TaskEvent
    mu          sync.RWMutex
}

func (tn *TaskNotifier) Subscribe() <-chan TaskEvent {
    tn.mu.Lock()
    defer tn.mu.Unlock()
    
    ch := make(chan TaskEvent, 10)
    tn.subscribers = append(tn.subscribers, ch)
    return ch
}

func (tn *TaskNotifier) Notify(event TaskEvent) {
    tn.mu.RLock()
    defer tn.mu.RUnlock()
    
    for _, ch := range tn.subscribers {
        select {
        case ch <- event:
        default: // Non-blocking send
        }
    }
}
```

### Performance Monitoring

**Concurrent Metrics Collection**:
```go
type Metrics struct {
    requests    int64
    errors      int64
    avgDuration int64
}

func (m *Metrics) RecordRequest(duration time.Duration, err error) {
    atomic.AddInt64(&m.requests, 1)
    if err != nil {
        atomic.AddInt64(&m.errors, 1)
    }
    
    // Update average duration using atomic operations
    currentAvg := atomic.LoadInt64(&m.avgDuration)
    newAvg := (currentAvg + duration.Nanoseconds()) / 2
    atomic.StoreInt64(&m.avgDuration, newAvg)
}
```

### Key Concurrency Benefits

1. **Scalability**: Handle multiple requests simultaneously
2. **Responsiveness**: Non-blocking operations with timeouts
3. **Resource Efficiency**: Connection pooling and worker pools
4. **Fault Tolerance**: Graceful degradation and error isolation
5. **Observability**: Concurrent metrics and logging
6. **Cancellation**: Context-based request cancellation
7. **Backpressure**: Rate limiting and queue management

These patterns enable the Task API to handle high-concurrency scenarios while maintaining data consistency and system stability.
