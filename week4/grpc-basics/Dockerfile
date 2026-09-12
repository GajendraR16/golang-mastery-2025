FROM golang:1.26-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY server/ ./server/
COPY proto/ ./proto/

RUN CGO_ENABLED=0 GOOS=linux go build -o task-server ./server/

FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root
COPY --from=builder /app/task-server .

EXPOSE 50051
CMD ["./task-server"]