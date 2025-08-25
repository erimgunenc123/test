# --- Stage 1: Build ---
FROM golang:1.25-alpine AS builder

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o ob_test ./cmd

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/ob_test .
COPY --from=builder /app/ .

EXPOSE 8080

CMD ["./ob_test"]
