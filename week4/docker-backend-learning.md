# Docker Basics — Backend Project Notes

## 1. Docker workflow

```text
Dockerfile
   ↓
docker build
   ↓
Image
   ↓
docker run
   ↓
Container
```

- **Dockerfile**: instructions for building an image.
- **Image**: packaged application/runtime.
- **Container**: a running instance of an image.

---

## 2. Dockerfile basics

A typical Go Dockerfile uses a multi-stage build:

```dockerfile
FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o app .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root
COPY --from=builder /app/app .
EXPOSE 8080
CMD ["./app"]
```

| Instruction | Meaning |
|---|---|
| `FROM` | Selects the base image |
| `WORKDIR` | Sets the working directory |
| `COPY` | Copies files from the build context |
| `RUN` | Runs a command while building |
| `EXPOSE` | Documents the container port |
| `CMD` | Default command when the container starts |

---

## 3. Multi-stage builds

The first stage compiles the Go application:

```dockerfile
FROM golang:1.26-alpine AS builder
```

The second stage runs only the compiled binary:

```dockerfile
FROM alpine:latest
```

```dockerfile
COPY --from=builder /app/app .
```

This keeps the final image smaller because the Go toolchain is not included in the runtime image.

---

## 4. Docker build context

The **build context** is the directory Docker is allowed to access during `docker build`.

Project structure:

```text
week4/
├── api-gateway/
├── grpc-basics/
└── user-service/
```

The Gateway needs `week4/` as its build context because it depends on the sibling modules.

From `week4/`:

```bash
sudo docker build -f api-gateway/Dockerfile .
```

- `-f api-gateway/Dockerfile` = Dockerfile to use.
- `.` = `week4/` is the build context.

---

## 5. Go local module replacements

The Gateway has:

```go
replace grpc-basics => ../grpc-basics
replace user-service => ../user-service
```

Therefore the Docker build must make those modules available at the paths expected by Go.

Gateway Dockerfile setup:

```dockerfile
WORKDIR /app

COPY api-gateway/go.mod api-gateway/go.sum ./
COPY api-gateway/ .
COPY grpc-basics/ /grpc-basics/
COPY user-service/ /user-service/

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o api-gateway .
```

The important idea is that the replacement paths must exist where Go expects them.

---

## 6. Go version compatibility

The Gateway's `go.mod` specifies:

```go
go 1.26.4
```

Therefore the Docker builder should use Go 1.26.x or newer.

Example:

```dockerfile
FROM golang:1.26-alpine AS builder
```

Using an older version such as `golang:1.24.5-alpine` can cause a version mismatch.

---

## 7. Building an image

Give the image a name with `-t`:

```bash
sudo docker build -t api-gateway -f api-gateway/Dockerfile .
```

Check images:

```bash
sudo docker images
```

---

## 8. Running a container

Basic form:

```bash
sudo docker run --name api-gateway api-gateway
```

Useful options:

```text
--name     Give the container a name
-p         Publish/map a port
-d         Run in the background
--network  Attach the container to a network
-e         Set an environment variable
```

Example:

```bash
sudo docker run -d   --name api-gateway   -p 8080:8080   api-gateway
```

---

## 9. Port mapping

```bash
-p 8080:8080
```

means:

```text
Host port 8080 → Container port 8080
```

`EXPOSE` does **not** publish a port to the host. You still need `-p` when running the container.

---

## 10. Checking containers

Running containers:

```bash
sudo docker ps
```

All containers:

```bash
sudo docker ps -a
```

Images:

```bash
sudo docker images
```

---

## 11. Stopping and removing containers

```bash
sudo docker stop api-gateway
sudo docker rm api-gateway
```

Force stop:

```bash
sudo docker kill api-gateway
```

For normal development, prefer `docker stop` followed by `docker rm` when recreating a container.

---

## 12. Docker logs

```bash
sudo docker logs api-gateway
```

Follow logs:

```bash
sudo docker logs -f api-gateway
```

---

## 13. Executing commands inside a container

Open a shell:

```bash
sudo docker exec -it api-gateway sh
```

Example:

```bash
ps
```

Exit:

```bash
exit
```

---

## 14. Docker networks

Manual Docker networking:

```bash
sudo docker network create backend-network
```

List networks:

```bash
sudo docker network ls
```

Inspect:

```bash
sudo docker network inspect backend-network
```

Search the output:

```bash
sudo docker network inspect backend-network | grep -E '"Name"|"IPv4Address"'
```

---

## 15. Container-to-container communication

Architecture:

```text
API Gateway :8080
      |
      ├── Task Service :50051
      |
      └── User Service :50052
```

Containers on the same Docker network can communicate using container/service names:

```text
task-service:50051
user-service:50052
```

Do not use:

```text
localhost:50051
localhost:50052
```

from inside the Gateway container.

Inside a container, `localhost` means **that same container**.

---

## 16. Environment variables for service addresses

The Gateway reads:

```text
TASK_SERVICE_ADDR
USER_SERVICE_ADDR
```

Local defaults can be:

```text
localhost:50051
localhost:50052
```

Inside Docker, use service names:

```text
TASK_SERVICE_ADDR=task-service:50051
USER_SERVICE_ADDR=user-service:50052
```

Manual Docker example:

```bash
sudo docker run -d   --name api-gateway   --network backend-network   -p 8080:8080   -e TASK_SERVICE_ADDR=task-service:50051   -e USER_SERVICE_ADDR=user-service:50052   api-gateway
```

---

## 17. Running the three services manually

Create network:

```bash
sudo docker network create backend-network
```

Task Service:

```bash
sudo docker run -d   --name task-service   --network backend-network   task-service
```

User Service:

```bash
sudo docker run -d   --name user-service   --network backend-network   user-service
```

API Gateway:

```bash
sudo docker run -d   --name api-gateway   --network backend-network   -p 8080:8080   -e TASK_SERVICE_ADDR=task-service:50051   -e USER_SERVICE_ADDR=user-service:50052   api-gateway
```

Check:

```bash
sudo docker ps
```

---

## 18. Testing the API

Gateway:

```bash
curl http://localhost:8080/tasks
```

The host connects to:

```text
localhost:8080
```

The Gateway communicates internally with:

```text
task-service:50051
user-service:50052
```

---

## 19. Docker Compose

Docker Compose lets you define multiple services in one `compose.yaml`.

Instead of manually creating a network and running three containers, Compose handles the setup.

Project structure:

```text
week4/
├── compose.yaml
├── api-gateway/
├── grpc-basics/
└── user-service/
```

Example:

```yaml
services:
  task-service:
    build:
      context: ./grpc-basics
      dockerfile: Dockerfile

  user-service:
    build:
      context: ./user-service
      dockerfile: Dockerfile

  api-gateway:
    build:
      context: .
      dockerfile: ./api-gateway/Dockerfile
    ports:
      - "8080:8080"
    environment:
      TASK_SERVICE_ADDR: task-service:50051
      USER_SERVICE_ADDR: user-service:50052
```

### Compose build contexts

Task Service:

```yaml
context: ./grpc-basics
```

User Service:

```yaml
context: ./user-service
```

Gateway:

```yaml
context: .
```

The Gateway needs the root `week4/` context because its Dockerfile accesses the other modules.

### Validate

```bash
sudo docker compose config
```

### Build and start

```bash
sudo docker compose up --build
```

### Stop and remove Compose containers/network

```bash
sudo docker compose down
```

`docker compose down` does not remove the images by default.

---

## 20. Compose networking

Compose automatically creates a default network for the services.

For example:

```text
week4_default
```

The services are automatically connected to it.

They can communicate using their service names:

```text
task-service:50051
user-service:50052
```

Therefore no manual:

```bash
docker network create
```

is required for the basic Compose setup.

---

## 21. Compose vs manual Docker

### Manual Docker

```text
docker network create
        ↓
docker run task-service
        ↓
docker run user-service
        ↓
docker run api-gateway
        ↓
configure networking
        ↓
configure environment variables
```

### Docker Compose

```text
compose.yaml
      ↓
docker compose up --build
      ↓
network + containers + configuration
```

Compose automates much of the setup we performed manually.

---

## 22. Useful troubleshooting commands

Check containers:

```bash
sudo docker ps
```

Check all containers:

```bash
sudo docker ps -a
```

Check images:

```bash
sudo docker images
```

Check logs:

```bash
sudo docker logs api-gateway
```

Enter a container:

```bash
sudo docker exec -it api-gateway sh
```

Inspect a network:

```bash
sudo docker network inspect backend-network
```

Search network output:

```bash
sudo docker network inspect backend-network | grep -E 'api-gateway|task-service|user-service'
```

Validate Compose:

```bash
sudo docker compose config
```

Start Compose:

```bash
sudo docker compose up --build
```

Stop Compose:

```bash
sudo docker compose down
```

---

## 23. Main concepts learned

### Dockerfile
Defines how an image is built.

### Image
The packaged application.

### Container
A running instance of an image.

### Build context
The directory Docker can access during `docker build`.

### Port mapping
Connects a host port to a container port:

```text
-p HOST:CONTAINER
```

### Docker network
Allows containers to communicate.

### Container DNS
Containers on the same Docker network can use container/service names:

```text
task-service:50051
user-service:50052
```

### Environment variables
Allow configuration to change between environments without changing application code.

### Multi-stage builds
Build with a Go image, then run the resulting binary in a smaller runtime image.

### Docker Compose
Defines and runs multiple related containers from one `compose.yaml`.

Compose also automatically handles the default network for the services.

---

## Project architecture after Dockerization

```text
                    Host
                     |
                  :8080
                     |
                     ▼
              ┌─────────────┐
              │ API Gateway │
              │   :8080     │
              └──────┬──────┘
                     │
               week4_default
                ┌────┴────┐
                │         │
                ▼         ▼
        ┌────────────┐ ┌────────────┐
        │Task Service│ │User Service│
        │   :50051   │ │   :50052   │
        └────────────┘ └────────────┘
```

## Current Docker learning progress

Covered:

- Dockerfile
- `FROM`
- `WORKDIR`
- `COPY`
- `RUN`
- `EXPOSE`
- `CMD`
- Multi-stage builds
- Docker images
- Docker containers
- Port mapping
- Build context
- Docker networks
- Container-to-container communication
- Container DNS
- Environment variables
- Docker Compose
- `compose.yaml`
- Compose build contexts
- `docker compose config`
- `docker compose up --build`
- `docker compose down`

## Next Docker topics

- Volumes
- PostgreSQL with Docker Compose
- Persistent data
- Health checks
- `depends_on`
- More Compose configuration
