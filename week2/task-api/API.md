# Task API Documentation

Base URL: `http://localhost:8080`

## Endpoints

### 1. Get All Tasks

Retrieve all tasks from the database.

**Endpoint**: `GET /tasks`

**Response**: `200 OK`

```json
[
  {
    "id": 1,
    "description": "Buy groceries",
    "complete": false,
    "created_at": "2026-01-15T10:30:00Z",
    "completed_at": null
  },
  {
    "id": 2,
    "description": "Finish project documentation",
    "complete": true,
    "created_at": "2026-01-14T09:00:00Z",
    "completed_at": "2026-01-15T11:00:00Z"
  }
]
```

**Example**:
```bash
curl http://localhost:8080/tasks
```

---

### 2. Get Task by ID

Retrieve a specific task by its ID.

**Endpoint**: `GET /tasks/{id}`

**Parameters**:
- `id` (path parameter): Task ID (integer)

**Response**: `302 Found`

```json
{
  "id": 1,
  "description": "Buy groceries",
  "complete": false,
  "created_at": "2026-01-15T10:30:00Z",
  "completed_at": null
}
```

**Error Response**: `404 Not Found`

```json
{
  "error": "Task not found"
}
```

**Example**:
```bash
curl http://localhost:8080/tasks/1
```

---

### 3. Create Task

Create a new task.

**Endpoint**: `POST /tasks`

**Request Body**:
```json
{
  "description": "Buy groceries"
}
```

**Validation Rules**:
- `description` is required
- Minimum length: 3 characters
- Whitespace is trimmed

**Response**: `201 Created`

```json
{
  "id": 3,
  "description": "Buy groceries",
  "complete": false,
  "created_at": "2026-01-15T12:00:00Z",
  "completed_at": null
}
```

**Error Responses**:

`400 Bad Request` - Invalid JSON:
```json
{
  "error": "Invalid Json"
}
```

`400 Bad Request` - Empty description:
```json
{
  "error": "Description cannot be empty"
}
```

`400 Bad Request` - Description too short:
```json
{
  "error": "Description is too short (min 3 chars)"
}
```

**Example**:
```bash
curl -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{"description": "Buy groceries"}'
```

---

### 4. Mark Task as Complete

Mark a task as completed.

**Endpoint**: `PUT /tasks/{id}`

**Parameters**:
- `id` (path parameter): Task ID (integer)

**Response**: `200 OK`

```json
{
  "id": 1,
  "description": "Buy groceries",
  "complete": true,
  "created_at": "2026-01-15T10:30:00Z",
  "completed_at": "2026-01-15T14:00:00Z"
}
```

**Error Response**: `404 Not Found`

```json
{
  "error": "Task Not Found"
}
```

**Example**:
```bash
curl -X PUT http://localhost:8080/tasks/1
```

---

### 5. Delete Task

Delete a task by ID.

**Endpoint**: `DELETE /tasks/{id}`

**Parameters**:
- `id` (path parameter): Task ID (integer)

**Response**: `204 No Content`

(Empty response body)

**Error Response**: `404 Not Found`

```json
{
  "error": "Incorrect Id"
}
```

**Example**:
```bash
curl -X DELETE http://localhost:8080/tasks/1
```

---

### 6. Search Tasks

Search for tasks by description.

**Endpoint**: `GET /tasks?q={query}`

**Query Parameters**:
- `q`: Search query string (case-insensitive)

**Response**: `200 OK`

```json
[
  {
    "id": 1,
    "description": "Buy groceries",
    "complete": false,
    "created_at": "2026-01-15T10:30:00Z",
    "completed_at": null
  }
]
```

**Example**:
```bash
curl "http://localhost:8080/tasks?q=groceries"
```

---

## Error Handling

All error responses follow this format:

```json
{
  "error": "Error message description"
}
```

### Common HTTP Status Codes

| Code | Description |
|------|-------------|
| 200 | OK - Request succeeded |
| 201 | Created - Resource created successfully |
| 204 | No Content - Request succeeded with no response body |
| 302 | Found - Resource found |
| 400 | Bad Request - Invalid input |
| 404 | Not Found - Resource not found |
| 500 | Internal Server Error - Server error |

---

## Testing with cURL

### Complete Workflow Example

```bash
# 1. Create a task
curl -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{"description": "Learn Go programming"}'

# 2. Get all tasks
curl http://localhost:8080/tasks

# 3. Get specific task (replace {id} with actual ID)
curl http://localhost:8080/tasks/1

# 4. Search for tasks
curl "http://localhost:8080/tasks?q=programming"

# 5. Mark task as complete
curl -X PUT http://localhost:8080/tasks/1

# 6. Delete task
curl -X DELETE http://localhost:8080/tasks/1
```

---

## Testing with HTTPie

```bash
# Create a task
http POST localhost:8080/tasks description="Learn Go programming"

# Get all tasks
http GET localhost:8080/tasks

# Get specific task
http GET localhost:8080/tasks/1

# Search tasks
http GET localhost:8080/tasks q==programming

# Mark as complete
http PUT localhost:8080/tasks/1

# Delete task
http DELETE localhost:8080/tasks/1
```

---

## Postman Collection

You can import these endpoints into Postman using the following structure:

1. Create a new collection named "Task API"
2. Set base URL variable: `{{baseUrl}}` = `http://localhost:8080`
3. Add each endpoint as described above

---

## CORS

The API includes CORS middleware that allows:
- All origins (`*`)
- Methods: GET, POST, PUT, DELETE, OPTIONS
- Headers: Content-Type

---

## Middleware

### Logging Middleware
All requests are logged with:
- HTTP method
- Request path
- Response status code
- Response time

### CORS Middleware
Handles cross-origin requests for browser-based clients.
