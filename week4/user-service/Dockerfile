FROM golang:1.26-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY server/ ./server/
COPY proto/ ./proto/

RUN CGO_ENABLED=0 GOOS=linux go build -o user-server ./server

FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root
COPY --from=builder /app/user-server .

EXPOSE 50052
CMD [ "./user-server" ]