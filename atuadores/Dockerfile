# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY . .

RUN go mod init atuador 2>/dev/null || true
RUN go mod tidy
RUN go build -o atuador

# Runtime stage
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/atuador .

CMD ["./atuador"]