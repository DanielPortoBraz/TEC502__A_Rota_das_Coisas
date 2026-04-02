# device/Dockerfile
FROM golang:1.22-alpine

WORKDIR /app

COPY . .

RUN go mod init sensores || true
RUN go mod tidy
RUN go build -o sensores

CMD ["./sensores"]