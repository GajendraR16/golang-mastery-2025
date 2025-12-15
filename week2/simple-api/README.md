# Simple API

A RESTful HTTP API server built with Go's standard library as part of the golang-mastery-2025 learning journey.

## Features

- ✅ **Health Check Endpoint**: Monitor API status
- ✅ **Hello Endpoint**: Personalized greeting with path parameters
- ✅ **Echo Endpoint**: JSON request/response handling
- ✅ **JSON Responses**: All endpoints return structured JSON
- ✅ **Error Handling**: Proper HTTP status codes and error messages
- ✅ **Method Validation**: Enforces correct HTTP methods

## API Endpoints

### GET /health
Returns the API health status.

**Response:**
```json
{
  "status": "ok"
}
```

### GET /hello/{name}
Returns a personalized greeting. If no name is provided, defaults to "World!".

**Examples:**
- `GET /hello/` → `{"message": "Hello, World!!"}`
- `GET /hello/John` → `{"message": "Hello, John!"}`

### POST /echo
Echoes back the JSON payload sent in the request body.

**Request:**
```json
{
  "name": "John",
  "age": 30
}
```

**Response:**
```json
{
  "name": "John",
  "age": 30
}
```

## Installation & Usage

### Prerequisites
- Go 1.24.5 or later

### Running the Server

```bash
# Navigate to the project directory
cd week2/simple-api

# Start the server
go run .
```

The server will start on `http://localhost:8080`

### Testing the API

```bash
# Health check
curl http://localhost:8080/health

# Hello endpoint
curl http://localhost:8080/hello/
curl http://localhost:8080/hello/Alice

# Echo endpoint
curl -X POST http://localhost:8080/echo \
  -H "Content-Type: application/json" \
  -d '{"message": "Hello API", "timestamp": "2025-01-01"}'
```

## Project Structure

```
simple-api/
├── main.go         # Server setup and routing
├── handlers.go     # HTTP handlers and response types
├── go.mod          # Go module definition
└── README.md       # This file
```

## Architecture

### Core Components

1. **HTTP Server**: Built with Go's standard `net/http` package
2. **JSON Handlers**: Structured response types and encoding
3. **Error Handling**: Consistent error responses with proper status codes
4. **Route Handlers**: Separate functions for each endpoint

### Response Types

```go
type ResponseStatus struct {
    Status string `json:"status"`
}

type MessageResponse struct {
    Message string `json:"message"`
}

type ErrorResponse struct {
    Error string `json:"error"`
}
```

## Key Go Concepts Demonstrated

- **HTTP Server**: Using `http.ListenAndServe` and `http.HandleFunc`
- **JSON Encoding/Decoding**: Marshal and unmarshal JSON data
- **HTTP Methods**: Method validation and routing
- **Request Handling**: Reading request bodies and URL paths
- **Response Writing**: Setting headers and status codes
- **Error Handling**: Proper error responses and status codes
- **Struct Tags**: JSON serialization tags
- **Interface Usage**: Using `any` interface for flexible JSON handling

## Error Handling

The API returns consistent error responses:

```json
{
  "error": "Method not allowed"
}
```

Common HTTP status codes used:
- `200 OK`: Successful requests
- `405 Method Not Allowed`: Wrong HTTP method
- `500 Internal Server Error`: JSON parsing or encoding errors

## Learning Objectives Achieved

- ✅ HTTP server development with Go standard library
- ✅ RESTful API design principles
- ✅ JSON request/response handling
- ✅ HTTP method validation
- ✅ Error handling in web applications
- ✅ URL path parameter extraction
- ✅ Request body parsing
- ✅ Response header management
- ✅ Clean code organization and separation of concerns