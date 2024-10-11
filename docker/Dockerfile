# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o chisme_server ./cmd/web

# Production image
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/chisme_server /app/chisme_server

COPY example/data.json /app/example/data.json

EXPOSE 4004

CMD ["/app/chisme_server"]
