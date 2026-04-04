# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY . .

RUN go mod init client 2>/dev/null || true
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o client

# Runtime stage
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/client .

CMD ["./client"]