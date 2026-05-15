FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o task-manager ./cmd

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/task-manager .
COPY --from=builder /app/.env.example .env
COPY --from=builder /app/migration ./migration

EXPOSE 8080

CMD ["./task-manager"]