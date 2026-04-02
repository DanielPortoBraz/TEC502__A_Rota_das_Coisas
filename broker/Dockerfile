# broker/Dockerfile
FROM golang:1.22-alpine

WORKDIR /app

COPY . .

RUN go mod init broker || true
RUN go mod tidy
RUN go build -o broker

CMD ["./broker"]