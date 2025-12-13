# Task Manager CLI

A command-line task management application built with Go as part of the golang-mastery-2025 learning journey.

## Features

- ✅ **Add Tasks**: Create new tasks with descriptions
- ✅ **List Tasks**: View all tasks with their status
- ✅ **Complete Tasks**: Mark tasks as completed
- ✅ **Delete Tasks**: Remove tasks from the list
- ✅ **Search Tasks**: Find tasks by description keywords
- ✅ **Persistent Storage**: Tasks are saved to `tasks.json`

## Installation & Usage

### Prerequisites
- Go 1.24.5 or later

### Running the Application

```bash
# Navigate to the project directory
cd week1/task-manager

# Add a new task
go run . add "Buy groceries"

# List all tasks
go run . list

# Complete a task (by ID)
go run . complete 1

# Delete a task (by ID)
go run . delete 2

# Search for tasks
go run . search "buy"
```

### Example Output

```bash
$ go run . add "Learn Go programming"
$ go run . add "Build a CLI app"
$ go run . list
-----Task List-----
1. [ ] Learn Go programming
2. [ ] Build a CLI app

$ go run . complete 1
Task completed.

$ go run . list
-----Task List-----
1. [✓] Learn Go programming
2. [ ] Build a CLI app
```

## Project Structure

```
task-manager/
├── main.go          # Entry point and CLI argument parsing
├── task.go          # Task struct and methods
├── manager.go       # TaskManager struct and business logic
├── storage.go       # JSON persistence layer
├── helper.go        # CLI command handlers
├── main_test.go     # Unit tests
├── go.mod           # Go module definition
└── tasks.json       # Persistent task storage (created at runtime)
```

## Architecture

### Core Components

1. **Task**: Represents a single task with ID, description, completion status, and timestamps
2. **TaskManager**: Manages collections of tasks and provides CRUD operations
3. **Storage**: Handles JSON serialization and file persistence
4. **CLI Interface**: Command-line argument parsing and user interaction

### Key Go Concepts Demonstrated

- **Structs and Methods**: Task and TaskManager types with associated methods
- **Error Handling**: Custom error types and proper error propagation
- **JSON Marshaling**: Struct tags and JSON serialization
- **File I/O**: Reading and writing files with proper error handling
- **Command-line Args**: Using `os.Args` for CLI interface
- **Unit Testing**: Test functions with table-driven tests
- **Memory Management**: Efficient slice operations for task deletion

## Testing

Run the test suite:

```bash
go test
```

Run tests with coverage:

```bash
go test -cover
```

## Learning Objectives Achieved

- ✅ Go project structure and organization
- ✅ Working with structs and methods
- ✅ Error handling patterns
- ✅ JSON data persistence
- ✅ Command-line application development
- ✅ Unit testing in Go
- ✅ Memory-efficient operations
- ✅ Clean code practices