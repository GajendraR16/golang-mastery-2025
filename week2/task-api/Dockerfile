FROM golang:1.24.5-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o task-api

FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root
COPY --from=builder /app/task-api .

EXPOSE 8080
CMD ["./task-api"]
