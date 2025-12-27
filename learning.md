# Go Learning Notes

## Variables, Functions, Constants, and Type Overflows

Variables in Go are declared using `var` keyword or short declaration `:=`. Functions are declared with `func` keyword. Constants use `const` keyword and are immutable. Type overflows occur when a value exceeds the range of its data type.

## Looping Constructs

For loop is the only looping construct in Go.

**Syntax:** `for init; condition; increment/decrement { //code }`

- `init` runs just once at the beginning of the loop
- `condition` is evaluated every single time from the beginning. The loop won't begin if condition is false
- At the end of each iteration, `increment/decrement` is executed

## Conditions

**if-else statements:** Go supports standard if-else conditional statements.

Go also has a short statement which can be written in the if clause to compute anything and check the condition. The variable initialized in this short statement has a lifetime limited to just the if-else block.

Example: `if x := getValue(); x > 10 { // use x }`

## Types

### Pointers

- `p = &x` → returns the address of variable x
- `*p` → dereferences the pointer and returns the value itself

### Structs

Structs can be described as a grouping of other generic or aggregate types to form a new user-defined type.

Struct fields can be accessed using dot notation: `structName.fieldName`

### Arrays

The length and capacity of arrays are defined in the declaration. Using `len()` and `cap()` functions have no computational complexity overhead while accessing array properties.

### Slices

Slices are a dynamic way of declaring arrays. They provide a more flexible alternative to arrays with variable length and built-in methods for manipulation.

**Important Memory Detail:**
When we check the slice address versus the underlying array address, they are different:
- The **array address** points to the first element of the underlying array
- The **slice address** points to the slice struct, which contains:
  - Pointer to the first element of the underlying array
  - `len` field (current length)
  - `cap` field (capacity)

This means a slice is actually a struct that references an underlying array, not the array itself.

**Slices are Non-Comparable:**
Slices cannot be compared using the `==` operator (except for comparing to `nil`). This is because:
- Slices contain a pointer to the underlying array
- Comparing slices would require deep comparison of elements
- Go doesn't provide built-in deep comparison for slices

To compare slices, you must iterate through elements manually or use helper functions.

**Append Functionality:**
The `append()` function is used to add elements to a slice. Behind the scenes:
- If the slice has sufficient capacity, `append()` adds elements to the existing underlying array
- If capacity is exceeded, `append()` allocates a new larger array and **copies** all existing elements to it
- The new slice is returned with updated length and capacity
- This copy operation can be expensive for large slices

Example:
```go
s := []int{1, 2, 3}
s = append(s, 4, 5)  // May trigger copy if capacity exceeded
```

### Maps

Maps are Go's built-in associative data type (hash tables or dictionaries in other languages). They map keys to values and provide fast lookups.

**Syntax:** `map[KeyType]ValueType`

**Creating Maps:**
```go
// Using make
m := make(map[string]int)

// Using map literal
m := map[string]int{"key1": 1, "key2": 2}
```

**Checking if a Key Exists:**
Maps return two values when accessing a key: the value and a boolean indicating if the key exists.

```go
value, ok := m["key"]
if !ok {
    // Key doesn't exist
}

// Or compare in one line
if val, ok := m["key"]; !ok || val != expectedValue {
    // Key doesn't exist OR value doesn't match expected
}
```

**Important Map Properties:**
- Maps are **reference types** - passing a map to a function passes a reference, not a copy
- Maps are **not safe for concurrent use** - require synchronization (mutex) for concurrent access
- The **zero value** of a map is `nil` - a nil map behaves like an empty map for reads but causes panic on writes
- Maps are **not comparable** using `==` (except for comparing to `nil`)
- Iteration order is **not guaranteed** - maps iterate in random order

**Deleting from Maps:**
```go
delete(m, "key")  // Removes key from map, safe even if key doesn't exist
```

### Closures

Closures are anonymous functions that can access variables from their outer scope. They "close over" variables from the enclosing function, maintaining access to them even after the outer function returns.

**Basic Closure Example:**
```go
func counter() func() int {
    count := 0
    return func() int {
        count++  // Accesses variable from outer scope
        return count
    }
}

// Usage
c := counter()
fmt.Println(c()) // 1
fmt.Println(c()) // 2
```

**Key Characteristics:**
- **Variable Capture:** Closures capture variables by reference, not by value
- **Lifetime Extension:** Variables captured by closures remain alive even after the outer function returns
- **State Preservation:** Each closure maintains its own copy of captured variables

**Common Use Cases:**
- **Event Handlers:** Capturing context for callback functions
- **Factory Functions:** Creating specialized functions with pre-configured behavior
- **Iterators:** Maintaining state between function calls
- **Decorators:** Wrapping functions with additional behavior

**Variable Capture Gotcha:**
```go
// Common mistake - all closures capture the same variable
funcs := make([]func(), 3)
for i := 0; i < 3; i++ {
    funcs[i] = func() {
        fmt.Println(i) // All print 3 (final value of i)
    }
}

// Correct approach - capture by value
for i := 0; i < 3; i++ {
    i := i  // Create new variable in loop scope
    funcs[i] = func() {
        fmt.Println(i) // Each prints its own value
    }
}
```

**Memory Considerations:**
- Closures keep references to captured variables, preventing garbage collection
- Be mindful of memory leaks when closures capture large objects or long-lived references

### Methods

Methods are functions with a special receiver argument. They allow you to define functions on types, enabling object-oriented programming patterns in Go.

**Syntax:** `func (receiver Type) methodName() returnType { }`

**Value vs Pointer Receivers:**

**Value Receiver:**
```go
func (p Person) getName() string {
    return p.name  // Receives a copy of the struct
}
```
- Method receives a **copy** of the value
- Cannot modify the original struct
- Use when you don't need to modify the receiver
- More memory efficient for small structs

**Pointer Receiver:**
```go
func (p *Person) setName(name string) {
    p.name = name  // Modifies the original struct
}
```
- Method receives a **pointer** to the original value
- Can modify the original struct
- Use when you need to modify the receiver or for large structs (avoids copying)
- Required for methods that modify the receiver

**Method Values and Method Expressions:**

**Method Value:**
A method bound to a specific receiver instance.
```go
p := Person{name: "John"}
methodValue := p.getName  // Method bound to instance p
result := methodValue()   // Calls p.getName()
```

**Method Expression:**
A function that takes the receiver as its first argument.
```go
methodExpr := Person.getName     // Method expression
result := methodExpr(p)          // Pass receiver as first argument
// Equivalent to: result := p.getName()
```

Method expressions are useful for:
- Passing methods as function arguments
- Creating generic functions that work with different receiver types

### Interfaces

Interfaces define a set of method signatures. A type implements an interface by implementing all the methods in the interface. Go uses implicit interface satisfaction - no explicit declaration needed.

**Basic Interface Example:**
```go
type Writer interface {
    Write([]byte) (int, error)
}

type File struct { /* ... */ }

func (f *File) Write(data []byte) (int, error) {
    // File now implements Writer interface
}
```

**Smaller Interfaces are Better:**
- **Interface Segregation Principle:** Prefer many small, focused interfaces over large ones
- Small interfaces are easier to implement and compose
- Better abstraction - clients depend only on methods they need
- Example: `io.Reader` has just one method `Read()`, making it highly reusable

**Interface Internal Structure:**
An interface value has two components:
1. **Dynamic Type (Concrete Type):** The actual type of the value stored
2. **Dynamic Value:** The actual value of the concrete type

```go
var w io.Writer
w = os.Stdout  // Dynamic type: *os.File, Dynamic value: os.Stdout

// Interface is nil only when both type and value are nil
var w io.Writer  // w is nil (type=nil, value=nil)
```

**Pointer Receivers and Interface Satisfaction:**

**Why Pointer Receiver Methods Need Pointers:**
```go
type Counter struct { count int }

func (c *Counter) Increment() {  // Pointer receiver
    c.count++
}

var c Counter
var inc interface{ Increment() }

inc = &c  // ✓ Works - &c has type *Counter
inc = c   // ✗ Fails - c has type Counter, doesn't have Increment method
```
- Only `*Counter` satisfies the interface, not `Counter`
- The method is defined on the pointer type
- You cannot call a pointer receiver method on a non-addressable value

**Why Value Receiver Methods Work with Both:**
```go
type Counter struct { count int }

func (c Counter) Get() int {  // Value receiver
    return c.count
}

var c Counter
var getter interface{ Get() int }

getter = c   // ✓ Works - c has Get method
getter = &c  // ✓ Also works - Go automatically dereferences
```
- Both `Counter` and `*Counter` satisfy the interface
- Go can automatically dereference `&c` to call the value receiver method
- Value receiver methods are in the method set of both the type and its pointer

**Type Assertions:**
Type assertions extract the concrete value from an interface.

```go
var w io.Writer = os.Stdout

// Type assertion
f := w.(*os.File)  // Extracts *os.File from interface

// Safe type assertion with ok check
f, ok := w.(*os.File)
if !ok {
    // w doesn't hold *os.File
}
```

**Why Type Assertions are Necessary:**
```go
var w io.Writer = os.Stdout

// Cannot call File-specific methods directly
w.Sync()  // ✗ Error: io.Writer has no Sync method

// Must use type assertion first
if f, ok := w.(*os.File); ok {
    f.Sync()  // ✓ Works - f is *os.File
}
```
- Interface only exposes methods in its definition
- Concrete type may have additional methods
- Type assertion is required to access methods not in the interface
- Without type assertion, you're limited to the interface's method set

**Type Switches:**
```go
func describe(i interface{}) {
    switch v := i.(type) {
    case int:
        fmt.Printf("Integer: %d\n", v)
    case string:
        fmt.Printf("String: %s\n", v)
    case *os.File:
        fmt.Printf("File: %s\n", v.Name())
    default:
        fmt.Printf("Unknown type\n")
    }
}
```

**Empty Interface:**
```go
var any interface{}  // or any (Go 1.18+)
```
- Can hold values of any type
- Provides no methods
- Useful for generic containers or functions that accept any type
- Requires type assertion to use the underlying value

**Interface Comparison:**
- Interfaces are comparable if their dynamic types are comparable
- Two interface values are equal if they have identical dynamic types and equal dynamic values
- Comparing interfaces with non-comparable dynamic types (like slices) causes panic

**Best Practice: "Accept Interfaces, Return Concrete Types"**

**Why Return Concrete Types (Not Interfaces):**
```go
// ✗ Bad - Returns interface
func NewWriter() io.Writer {
    return &bytes.Buffer{}
}

// ✓ Good - Returns concrete type
func NewWriter() *bytes.Buffer {
    return &bytes.Buffer{}
}
```

**Reasons to return concrete types:**
1. **Flexibility for Callers:** Callers can access all methods of the concrete type, not just interface methods
2. **No Premature Abstraction:** Don't create abstractions until you need them
3. **Easier to Extend:** Adding methods to concrete types doesn't break existing code
4. **Better Documentation:** Concrete types show exactly what's being returned
5. **Avoid Interface Pollution:** Prevents creating unnecessary interfaces

**Example of the Problem:**
```go
// Returns interface - limits caller
func GetConfig() io.Reader {
    return &bytes.Buffer{}
}

config := GetConfig()
config.Reset()  // ✗ Error: io.Reader has no Reset method

// Returns concrete type - caller has full access
func GetConfig() *bytes.Buffer {
    return &bytes.Buffer{}
}

config := GetConfig()
config.Reset()  // ✓ Works: *bytes.Buffer has Reset method
```

**Why Accept Interfaces as Parameters:**
```go
// ✓ Good - Accepts interface
func WriteData(w io.Writer, data []byte) error {
    _, err := w.Write(data)
    return err
}

// Can be called with any type that implements io.Writer
WriteData(os.Stdout, data)
WriteData(&bytes.Buffer{}, data)
WriteData(httpResponseWriter, data)
```

**Reasons to accept interfaces:**
1. **Maximum Flexibility:** Function works with any type implementing the interface
2. **Testability:** Easy to pass mock implementations for testing
3. **Decoupling:** Function doesn't depend on concrete implementations
4. **Composability:** Promotes code reuse across different types
5. **Dependency Inversion:** Depend on abstractions, not concretions

**Summary:**
- **Parameters:** Use interfaces for flexibility and testability
- **Return Values:** Use concrete types to avoid limiting callers
- **Proverb:** "Be conservative in what you send, be liberal in what you accept"


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

## Bit Operations

Bit operations are fundamental operations that work directly on binary representations of numbers:

- `x & 1` → checks the last bit (determines if number is odd/even)
- `x >> 1` → right shifts by 1 bit, removes the last bit. Equivalent to `x / 2^1`
- `x << 1` → left shifts by 1 bit. Equivalent to `x * 2^1`
- `x & (x - 1)` → clears the rightmost 1-bit

**Example:**
```
x = 5 (binary: 101)
x - 1 = 4 (binary: 100)
x & (x - 1) = 101 & 100 = 100 (decimal: 4)
```

This operation is commonly used in algorithms to count set bits or check if a number is a power of 2.

