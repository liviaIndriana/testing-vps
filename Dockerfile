# syntax=docker/dockerfile:1
FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o server cmd/api/main.go

FROM alpine:3.18
WORKDIR /app
COPY --from=builder /app/server .

EXPOSE 8081
CMD ["./server"]