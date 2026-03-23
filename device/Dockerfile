# device/Dockerfile
FROM golang:1.22-alpine

WORKDIR /app

COPY . .

RUN go mod init device || true
RUN go mod tidy
RUN go build -o device

CMD ["./device"]