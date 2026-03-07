# Task API Implementation Notes

This file contains implementation-heavy notes from the Task API project (HTTP, routing, middleware, PostgreSQL, Docker). Core reusable Go concepts remain in `learning.md`.

## net/http Package

The `net/http` package provides HTTP client and server implementations for building web applications and APIs.

### HTTP Server Basics

`ListenAndServe` accepts a port and a type that implements the `ServeHTTP` method. We can use `r.URL.Path` (where `r` is of type `*http.Request`) to segregate endpoints and functionality.

The http package also provides `NewServeMux` which has a `.Handle` function that accepts `HandlerFunc`. `HandlerFunc` implements the `ServeHTTP` interface so we don't have to implement it for each type or endpoint. `HandlerFunc` is just a wrapper or adapter function with the signature `func(w http.ResponseWriter, r *http.Request)`.

To make it convenient, http provides us with a global `ServeMux` instance, which reduces the complexity of implementation:

```go
http.HandleFunc("/", handler)
http.ListenAndServe(":8080", nil)
```

### JSON Encoding and Decoding

When building APIs, JSON encoding and decoding are essential for handling request/response data.

**JSON Encoder (`json.NewEncoder`):**
- Used to encode Go values directly to an `io.Writer` (like `http.ResponseWriter`)
- More efficient for streaming data as it writes directly to the output
- Commonly used for HTTP responses

```go
func jsonHandler(w http.ResponseWriter, data any) {
    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(data); err != nil {
        http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
    }
}
```

**JSON Decoder (`json.NewDecoder`):**
- Used to decode JSON directly from an `io.Reader` (like `http.Request.Body`)
- More efficient for streaming data as it reads directly from the input
- Commonly used for parsing HTTP request bodies

```go
func parseJSON(r *http.Request, v any) error {
    defer r.Body.Close()
    return json.NewDecoder(r.Body).Decode(v)
}
```

**Encoder vs Marshal / Decoder vs Unmarshal:**
- `json.Marshal/Unmarshal`: Work with byte slices, require loading entire data into memory
- `json.NewEncoder/NewDecoder`: Work with streams (`io.Writer`/`io.Reader`), more memory efficient
- For HTTP handlers, Encoder/Decoder are preferred as they work directly with request/response streams

**Common HTTP Response Patterns:**
```go
// Success response
type Response struct {
    Status string `json:"status"`
    Data   any    `json:"data,omitempty"`
}

// Error response
type ErrorResponse struct {
    Error string `json:"error"`
}

func jsonError(w http.ResponseWriter, message string, code int) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
    json.NewEncoder(w).Encode(ErrorResponse{Error: message})
}
```

## Third-Party Routing with Gorilla Mux

While Go's standard `net/http` package provides basic routing capabilities, complex APIs often require more sophisticated routing features. The Gorilla Mux package is a popular third-party router that extends Go's routing capabilities.

### Why Use Gorilla Mux?

**Standard Library Limitations:**
```go
// Standard library - basic pattern matching
http.HandleFunc("/tasks/", taskHandler)  // Matches /tasks/anything
```

**Gorilla Mux Advantages:**
```go
// Gorilla Mux - precise routing with constraints
router.HandleFunc("/tasks/{id:[0-9]+}", taskHandler).Methods("GET")
```

### Key Features

**1. Path Variables with Regex Constraints:**
```go
// Extract ID from URL path with validation
router.HandleFunc("/tasks/{id:[0-9]+}", taskHandler)

// In handler:
vars := mux.Vars(r)
id, _ := strconv.Atoi(vars["id"])  // Safe conversion - regex ensures numeric
```

**2. HTTP Method Routing:**
```go
router.HandleFunc("/tasks", getAllTasks).Methods("GET")
router.HandleFunc("/tasks", createTask).Methods("POST")
router.HandleFunc("/tasks/{id}", updateTask).Methods("PUT")
router.HandleFunc("/tasks/{id}", deleteTask).Methods("DELETE")
```

**3. Query Parameter Routing:**
```go
// Route based on query parameters
router.HandleFunc("/tasks", searchTasks).Methods("GET").Queries("q", "{q}")
router.HandleFunc("/tasks", getAllTasks).Methods("GET")  // Fallback route

// Access query parameters:
query := r.URL.Query().Get("q")
```

**4. Route Priority and Specificity:**
```go
// More specific routes should be registered first
router.HandleFunc("/tasks", searchHandler).Methods("GET").Queries("q", "{q}")
router.HandleFunc("/tasks", listHandler).Methods("GET")  // General fallback
```

### RESTful API Design Patterns

The task-api demonstrates RESTful API design principles:

**Resource-Based URLs:**
- `GET /tasks` - List all tasks
- `POST /tasks` - Create new task
- `GET /tasks/{id}` - Get specific task
- `PUT /tasks/{id}` - Update/complete task
- `DELETE /tasks/{id}` - Delete task
- `GET /tasks?q=search` - Search tasks

**HTTP Status Codes:**
```go
// Success responses
jsonHandler(w, http.StatusOK, data)        // 200 - Success
jsonHandler(w, http.StatusCreated, task)   // 201 - Created
w.WriteHeader(http.StatusNoContent)        // 204 - No Content (DELETE)

// Error responses
jsonError(w, "Not Found", http.StatusNotFound)           // 404
jsonError(w, "Bad Request", http.StatusBadRequest)       // 400
jsonError(w, "Internal Error", http.StatusInternalServerError) // 500
```

### Input Validation Patterns

**Struct-Based Validation:**
```go
type TaskData struct {
    Description string `json:"description" validate:"required,min=3"`
}

// Manual validation in handler
task.Description = strings.TrimSpace(task.Description)
if task.Description == "" {
    jsonError(w, "Description cannot be empty", http.StatusBadRequest)
    return
}
if len(task.Description) < 3 {
    jsonError(w, "Description too short", http.StatusBadRequest)
    return
}
```

**Separation of Concerns:**
- **Models** (`models.go`): Data structures and business logic
- **Handlers** (`handlers.go`): HTTP request/response handling
- **Storage** (`storage.go`): Data persistence layer
- **Main** (`main.go`): Application setup and routing

### Error Handling in APIs

**Consistent Error Response Format:**
```go
type ErrorResponse struct {
    Error string `json:"error"`
}

func jsonError(w http.ResponseWriter, message string, code int) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
    json.NewEncoder(w).Encode(ErrorResponse{Error: message})
}
```

**Graceful Error Handling:**
```go
// Handle missing resources
if task == nil {
    jsonError(w, "Task Not Found", http.StatusNotFound)
    return
}

// Handle invalid input
if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
    jsonError(w, "Invalid JSON", http.StatusBadRequest)
    return
}
```

### API Response Patterns

**Consistent JSON Responses:**
```go
func jsonHandler(w http.ResponseWriter, code int, data any) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
    if err := json.NewEncoder(w).Encode(data); err != nil {
        http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
    }
}
```

**Resource Creation Pattern:**
```go
// POST /tasks - Return created resource with 201 status
createdTask := tm.Add(task.Description)
SaveTasks(tm.Tasks, filename)
jsonHandler(w, http.StatusCreated, createdTask)
```

**Resource Deletion Pattern:**
```go
// DELETE /tasks/{id} - Return 204 No Content on success
if tm.Delete(id) {
    SaveTasks(tm.Tasks, filename)
    w.WriteHeader(http.StatusNoContent)  // No response body needed
    return
}
```

### Key Learnings from Task API Implementation

1. **Third-Party Routing**: Gorilla Mux provides powerful routing features beyond standard library
2. **RESTful Design**: Consistent URL patterns and HTTP methods for resource operations
3. **Input Validation**: Always validate and sanitize user input before processing
4. **Error Handling**: Provide consistent, informative error responses with proper status codes
5. **Separation of Concerns**: Organize code into logical layers (models, handlers, storage)
6. **Resource Management**: Proper handling of request bodies with `defer r.Body.Close()`
7. **Query Parameters**: Handle both path variables and query parameters for flexible APIs
8. **Status Codes**: Use appropriate HTTP status codes to communicate operation results

## HTTP Middleware

Middleware is a powerful pattern in web development that allows you to wrap HTTP handlers with additional functionality. Middleware functions execute before and/or after the main handler, enabling cross-cutting concerns like logging, authentication, CORS, and request validation.

### Understanding Middleware Pattern

**Basic Middleware Concept:**
```go
type Middleware func(http.Handler) http.Handler
```

A middleware is a function that takes an `http.Handler` and returns an `http.Handler`. This allows you to "wrap" handlers with additional behavior.

**Simple Middleware Example:**
```go
func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        
        // Call the next handler
        next.ServeHTTP(w, r)
        
        // Log after handler completes
        log.Printf("%s %s %v", r.Method, r.URL.Path, time.Since(start))
    })
}
```

### Common Middleware Patterns

**1. Logging Middleware:**
```go
func LoggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        
        // Create a response writer wrapper to capture status code
        wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
        
        next.ServeHTTP(wrapped, r)
        
        log.Printf("[%s] %s %s %d %v",
            time.Now().Format("2006-01-02 15:04:05"),
            r.Method,
            r.URL.Path,
            wrapped.statusCode,
            time.Since(start),
        )
    })
}

type responseWriter struct {
    http.ResponseWriter
    statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
    rw.statusCode = code
    rw.ResponseWriter.WriteHeader(code)
}
```

**2. CORS Middleware:**
```go
func CORSMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
        
        // Handle preflight requests
        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }
        
        next.ServeHTTP(w, r)
    })
}
```

**3. Authentication Middleware:**
```go
func AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        token := r.Header.Get("Authorization")
        
        if token == "" {
            http.Error(w, "Missing authorization header", http.StatusUnauthorized)
            return
        }
        
        // Validate token (simplified)
        if !isValidToken(token) {
            http.Error(w, "Invalid token", http.StatusUnauthorized)
            return
        }
        
        // Add user info to context
        ctx := context.WithValue(r.Context(), "userID", getUserID(token))
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
```

**4. Rate Limiting Middleware:**
```go
func RateLimitMiddleware(requestsPerMinute int) func(http.Handler) http.Handler {
    limiter := rate.NewLimiter(rate.Limit(requestsPerMinute), requestsPerMinute)
    
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            if !limiter.Allow() {
                http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
                return
            }
            next.ServeHTTP(w, r)
        })
    }
}
```

### Middleware Chaining

**Manual Chaining:**
```go
func main() {
    handler := http.HandlerFunc(homeHandler)
    
    // Wrap with middleware (innermost first)
    handler = loggingMiddleware(handler)
    handler = corsMiddleware(handler)
    handler = authMiddleware(handler)
    
    http.Handle("/", handler)
}
```

**Chain Helper Function:**
```go
func Chain(h http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
    for i := len(middlewares) - 1; i >= 0; i-- {
        h = middlewares[i](h)
    }
    return h
}

// Usage
handler := Chain(
    http.HandlerFunc(homeHandler),
    loggingMiddleware,
    corsMiddleware,
    authMiddleware,
)
```

### Middleware with Gorilla Mux

**Global Middleware:**
```go
router := mux.NewRouter()
router.Use(loggingMiddleware)
router.Use(corsMiddleware)

router.HandleFunc("/api/tasks", tasksHandler)
```

**Route-Specific Middleware:**
```go
// Protected routes
protected := router.PathPrefix("/api").Subrouter()
protected.Use(authMiddleware)
protected.HandleFunc("/tasks", tasksHandler)

// Public routes
router.HandleFunc("/health", healthHandler)
```

### Context Usage in Middleware

**Passing Data Between Middleware:**
```go
type contextKey string

const UserIDKey contextKey = "userID"

func AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        userID := extractUserID(r)
        ctx := context.WithValue(r.Context(), UserIDKey, userID)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

func protectedHandler(w http.ResponseWriter, r *http.Request) {
    userID := r.Context().Value(UserIDKey).(string)
    // Use userID in handler
}
```

### Error Handling in Middleware

**Panic Recovery Middleware:**
```go
func RecoveryMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            if err := recover(); err != nil {
                log.Printf("Panic recovered: %v", err)
                http.Error(w, "Internal Server Error", http.StatusInternalServerError)
            }
        }()
        next.ServeHTTP(w, r)
    })
}
```

### Middleware Best Practices

**1. Order Matters:**
```go
// Correct order (outer to inner)
handler = recoveryMiddleware(     // Catch panics from all inner middleware
    loggingMiddleware(            // Log all requests
        corsMiddleware(           // Handle CORS for all requests
            authMiddleware(       // Authenticate specific routes
                rateLimitMiddleware(  // Rate limit authenticated users
                    actualHandler,
                ),
            ),
        ),
    ),
)
```

**2. Early Returns:**
```go
func ValidationMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if r.Header.Get("Content-Type") != "application/json" {
            http.Error(w, "Content-Type must be application/json", http.StatusBadRequest)
            return // Don't call next.ServeHTTP
        }
        next.ServeHTTP(w, r)
    })
}
```

**3. Configurable Middleware:**
```go
func TimeoutMiddleware(timeout time.Duration) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            ctx, cancel := context.WithTimeout(r.Context(), timeout)
            defer cancel()
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}

// Usage
handler = TimeoutMiddleware(30 * time.Second)(handler)
```

### Testing Middleware

**Unit Testing Middleware:**
```go
func TestLoggingMiddleware(t *testing.T) {
    // Create a test handler
    testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("test response"))
    })
    
    // Wrap with middleware
    handler := LoggingMiddleware(testHandler)
    
    // Create test request
    req := httptest.NewRequest("GET", "/test", nil)
    rr := httptest.NewRecorder()
    
    // Execute
    handler.ServeHTTP(rr, req)
    
    // Assert
    assert.Equal(t, http.StatusOK, rr.Code)
    assert.Equal(t, "test response", rr.Body.String())
}
```

### Key Middleware Concepts

1. **Wrapper Pattern**: Middleware wraps handlers to add functionality
2. **Chain of Responsibility**: Multiple middleware can be chained together
3. **Context Propagation**: Use `context.Context` to pass data between middleware
4. **Early Termination**: Middleware can stop the chain by not calling `next.ServeHTTP`
5. **Order Dependency**: The order of middleware application affects behavior
6. **Reusability**: Well-designed middleware can be reused across different routes
7. **Separation of Concerns**: Each middleware should handle one specific concern
8. **Performance Impact**: Middleware adds overhead, so use judiciously

Middleware is essential for building robust web applications, providing a clean way to handle cross-cutting concerns without cluttering your main business logic.


## PostgreSQL Integration in Go

PostgreSQL is a powerful, open-source relational database. Go provides excellent support for PostgreSQL through the `database/sql` package and the `lib/pq` driver.

### Database Connection

**Basic Connection Setup:**
```go
import (
    "database/sql"
    _ "github.com/lib/pq"  // PostgreSQL driver
)

func NewPostgresStore(connStr string) (*PostgresStore, error) {
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        return nil, err
    }
    
    // Verify connection
    if err := db.Ping(); err != nil {
        return nil, err
    }
    
    return &PostgresStore{db: db}, nil
}
```

**Connection String Format:**
```
postgres://username:password@host:port/database?sslmode=disable
```

### Repository Pattern

The repository pattern abstracts database operations, providing a clean interface for data access:

```go
type PostgresStore struct {
    db *sql.DB
}

// CRUD operations
func (s *PostgresStore) CreateTask(description string) (*Task, error)
func (s *PostgresStore) GetAllTasks() ([]*Task, error)
func (s *PostgresStore) GetTaskById(id int) (*Task, error)
func (s *PostgresStore) UpdateTask(id int) (*Task, error)
func (s *PostgresStore) DeleteTaskById(id int) error
```

### SQL Query Patterns

**1. INSERT with RETURNING:**
```go
func (s *PostgresStore) CreateTask(description string) (*Task, error) {
    query := `
        INSERT INTO tasks (description)
        VALUES ($1)
        RETURNING id, description, completed, created_at, completed_at
    `
    
    var task Task
    var completedAt sql.NullTime
    
    err := s.db.QueryRow(query, description).Scan(
        &task.ID,
        &task.Description,
        &task.Completed,
        &task.CreatedAt,
        &completedAt,
    )
    
    if completedAt.Valid {
        task.CompletedAt = &completedAt.Time
    }
    
    return &task, err
}
```

**2. SELECT Multiple Rows:**
```go
func (s *PostgresStore) GetAllTasks() ([]*Task, error) {
    query := `SELECT id, description, completed, created_at, completed_at FROM tasks`
    
    rows, err := s.db.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var tasks []*Task
    for rows.Next() {
        var task Task
        var completedAt sql.NullTime
        
        err := rows.Scan(
            &task.ID,
            &task.Description,
            &task.Completed,
            &task.CreatedAt,
            &completedAt,
        )
        if err != nil {
            return nil, err
        }
        
        if completedAt.Valid {
            task.CompletedAt = &completedAt.Time
        }
        
        tasks = append(tasks, &task)
    }
    
    return tasks, rows.Err()
}
```

**3. UPDATE with RETURNING:**
```go
func (s *PostgresStore) CompletedTaskById(id int) (*Task, error) {
    query := `
        UPDATE tasks 
        SET completed = true, completed_at = $1 
        WHERE id = $2
        RETURNING id, description, completed, created_at, completed_at
    `
    
    var task Task
    var completedAt sql.NullTime
    now := time.Now()
    
    err := s.db.QueryRow(query, now, id).Scan(
        &task.ID,
        &task.Description,
        &task.Completed,
        &task.CreatedAt,
        &completedAt,
    )
    
    if err == sql.ErrNoRows {
        return nil, fmt.Errorf("task not found")
    }
    
    if completedAt.Valid {
        task.CompletedAt = &completedAt.Time
    }
    
    return &task, err
}
```

**4. DELETE with Row Count Check:**
```go
func (s *PostgresStore) DeleteTaskById(id int) error {
    query := `DELETE FROM tasks WHERE id = $1`
    
    res, err := s.db.Exec(query, id)
    if err != nil {
        return err
    }
    
    rowsAffected, _ := res.RowsAffected()
    if rowsAffected == 0 {
        return sql.ErrNoRows
    }
    
    return nil
}
```

### Handling NULL Values

PostgreSQL allows NULL values, which don't map directly to Go types. Use `sql.Null*` types:

```go
var completedAt sql.NullTime

err := rows.Scan(&task.ID, &task.Description, &completedAt)

// Check if value is NULL
if completedAt.Valid {
    task.CompletedAt = &completedAt.Time  // Convert to *time.Time
} else {
    task.CompletedAt = nil
}
```

**Common NULL Types:**
- `sql.NullString` - for nullable strings
- `sql.NullInt64` - for nullable integers
- `sql.NullFloat64` - for nullable floats
- `sql.NullBool` - for nullable booleans
- `sql.NullTime` - for nullable timestamps

### SQL Injection Prevention

**Always use parameterized queries:**
```go
// ✓ Safe - parameterized query
query := `SELECT * FROM tasks WHERE id = $1`
rows, err := db.Query(query, userInput)

// ✗ Dangerous - SQL injection vulnerability
query := fmt.Sprintf("SELECT * FROM tasks WHERE id = %s", userInput)
rows, err := db.Query(query)
```

PostgreSQL uses `$1, $2, $3` for parameter placeholders (not `?` like MySQL).

### Error Handling

**Common Database Errors:**
```go
import "database/sql"

// No rows found
if err == sql.ErrNoRows {
    return nil, fmt.Errorf("task not found")
}

// Connection errors
if err != nil {
    log.Printf("Database error: %v", err)
    return nil, err
}

// Check rows affected
rowsAffected, err := result.RowsAffected()
if rowsAffected == 0 {
    return fmt.Errorf("no rows affected")
}
```

### Database Schema Management

**Schema Definition (schema.sql):**
```sql
CREATE TABLE IF NOT EXISTS tasks (
    id SERIAL PRIMARY KEY,
    description TEXT NOT NULL,
    completed BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMP
);
```

**Key SQL Features:**
- `SERIAL` - Auto-incrementing integer
- `PRIMARY KEY` - Unique identifier
- `DEFAULT` - Default values
- `NOT NULL` - Required fields
- `TIMESTAMP` - Date and time storage

### Connection Pooling

The `database/sql` package automatically manages a connection pool:

```go
// Configure connection pool
db.SetMaxOpenConns(25)           // Maximum open connections
db.SetMaxIdleConns(5)            // Maximum idle connections
db.SetConnMaxLifetime(5 * time.Minute)  // Connection lifetime
```

**Best Practices:**
- Don't create a new connection for each request
- Reuse the `*sql.DB` instance across your application
- The pool handles concurrent access automatically

### Testing with PostgreSQL

**Test Database Setup:**
```go
func setupTestDB(t *testing.T) *PostgresStore {
    connStr := "postgres://postgres:postgres@localhost:5432/taskdb_test?sslmode=disable"
    store, err := NewPostgresStore(connStr)
    if err != nil {
        t.Fatal(err)
    }
    return store
}

func TestCreateTask(t *testing.T) {
    store := setupTestDB(t)
    defer store.TruncateTasks()  // Clean up after test
    
    task, err := store.CreateTask("Test task")
    if err != nil {
        t.Fatalf("Failed to create task: %v", err)
    }
    
    if task.ID == 0 {
        t.Error("Expected task ID to be set")
    }
}
```

**Test Isolation:**
- Use a separate test database
- Truncate tables between tests
- Use transactions that rollback for isolation

### Docker Integration

**Docker Compose for PostgreSQL:**
```yaml
services:
  postgres:
    image: postgres:15
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
      POSTGRES_DB: taskdb
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./schema.sql:/docker-entrypoint-initdb.d/init.sql
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres -d taskdb"]
      interval: 3s
      timeout: 3s
      retries: 5
```

**Benefits:**
- Automatic schema initialization
- Health checks for service readiness
- Data persistence with volumes
- Easy local development setup

### Key PostgreSQL Concepts

1. **Parameterized Queries**: Always use `$1, $2` placeholders to prevent SQL injection
2. **NULL Handling**: Use `sql.Null*` types for nullable database columns
3. **RETURNING Clause**: Get inserted/updated data in one query
4. **Connection Pooling**: Reuse `*sql.DB` instance, don't create new connections
5. **Error Handling**: Check for `sql.ErrNoRows` and handle connection errors
6. **Repository Pattern**: Abstract database operations behind clean interfaces
7. **Test Isolation**: Use separate test databases and clean up between tests
8. **Docker Integration**: Containerize database for consistent development environment

## Docker and Containerization

Docker enables packaging applications with their dependencies into containers, ensuring consistent behavior across different environments.

### Dockerfile Basics

**Multi-Stage Build for Go:**
```dockerfile
# Build stage
FROM golang:1.24.5-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main .

# Runtime stage
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/main .
EXPOSE 8080
CMD ["./main"]
```

**Benefits of Multi-Stage Builds:**
- Smaller final image (only runtime dependencies)
- Build tools not included in production image
- Faster deployment and reduced attack surface

### Docker Compose

**Multi-Container Application:**
```yaml
services:
  postgres:
    image: postgres:15
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
      POSTGRES_DB: taskdb
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./schema.sql:/docker-entrypoint-initdb.d/init.sql
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres -d taskdb"]
      interval: 3s
      timeout: 3s
      retries: 5

  api:
    build: .
    ports:
      - "8080:8080"
    environment:
      DATABASE_URL: postgres://postgres:postgres@postgres:5432/taskdb?sslmode=disable
    depends_on:
      postgres:
        condition: service_healthy

volumes:
  postgres_data:
```

**Key Concepts:**
- **Services**: Define containers (postgres, api)
- **Volumes**: Persist data across container restarts
- **Networks**: Containers can communicate by service name
- **Health Checks**: Ensure services are ready before starting dependents
- **depends_on**: Control startup order

**Common Commands:**
```bash
docker-compose up --build    # Build and start all services
docker-compose down          # Stop and remove containers
docker-compose logs -f       # Follow logs
docker-compose ps            # List running services
```

### Environment Configuration

**Environment Variables in Docker:**
```go
connStr := os.Getenv("DATABASE_URL")
if connStr == "" {
    connStr = "postgres://localhost:5432/taskdb?sslmode=disable"  // Fallback
}
```

**Benefits:**
- Different configurations for dev/staging/production
- Secrets management (passwords, API keys)
- Easy configuration changes without rebuilding

### Container Networking

**Service Discovery:**
```yaml
api:
  environment:
    DATABASE_URL: postgres://postgres:postgres@postgres:5432/taskdb
    # 'postgres' resolves to the postgres service
```

Containers in the same Docker Compose network can communicate using service names as hostnames.

### Volume Management

**Types of Volumes:**
1. **Named Volumes**: Managed by Docker, persist data
   ```yaml
   volumes:
     - postgres_data:/var/lib/postgresql/data
   ```

2. **Bind Mounts**: Map host directory to container
   ```yaml
   volumes:
     - ./schema.sql:/docker-entrypoint-initdb.d/init.sql
   ```

**Use Cases:**
- Named volumes: Database data persistence
- Bind mounts: Configuration files, development code

### Health Checks

**PostgreSQL Health Check:**
```yaml
healthcheck:
  test: ["CMD-SHELL", "pg_isready -U postgres -d taskdb"]
  interval: 3s
  timeout: 3s
  retries: 5
```

**Benefits:**
- Ensures database is ready before starting API
- Automatic restart on failure
- Better orchestration with `depends_on`

### Docker Best Practices

1. **Multi-Stage Builds**: Reduce image size
2. **Layer Caching**: Order Dockerfile commands for efficient caching
3. **Health Checks**: Ensure service readiness
4. **Named Volumes**: Persist important data
5. **Environment Variables**: Externalize configuration
6. **Service Dependencies**: Use `depends_on` with health checks
7. **.dockerignore**: Exclude unnecessary files from build context

### Development Workflow

**Local Development:**
```bash
# Start all services
docker-compose up --build

# View logs
docker-compose logs -f api

# Execute commands in container
docker-compose exec api sh

# Stop services
docker-compose down

# Remove volumes (clean slate)
docker-compose down -v
```

**Benefits:**
- Consistent environment across team
- Easy onboarding for new developers
- No need to install PostgreSQL locally
- Isolated development environment
