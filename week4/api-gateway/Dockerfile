FROM golang:1.26-alpine AS builder

WORKDIR /app
COPY api-gateway/go.mod api-gateway/go.sum ./
COPY api-gateway/ .
COPY grpc-basics/ /grpc-basics/
COPY user-service/ /user-service/

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o api-gateway .

FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root
COPY --from=builder /app/api-gateway .

EXPOSE 8080
CMD ["./api-gateway"]