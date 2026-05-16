# Task Manager Service

Production-grade Golang task manager service.

## Features

- Gin REST API
- PostgreSQL
- Redis cache-aside
- Swagger/OpenAPI
- Prometheus metrics
- OpenTelemetry tracing
- Docker multi-stage build
- Pagination & filtering
- Unit & integration tests
- Benchmark & pprof

---

## Run Project

```bash
cp .env.example .env
docker compose up --build
```

---

## Run Migrations

```bash
migrate \
-path migrations \
-database "postgres://postgres:postgres@localhost:5432/task_manager?sslmode=disable" \
up
```

---

## Run Tests

```bash
go test ./...
```

---

## Generate Coverage

```bash
go test ./... -coverprofile=coverage.out
go tool cover -func=coverage.out
```

---

## Swagger

http://localhost:8080/swagger/index.html

---

## Metrics

http://localhost:8080/metrics

---

## pprof

http://localhost:8080/debug/pprof/

---

## Architecture

Handler -> Service -> Repository -> PostgreSQL
                               -> Redis Cache

---

## Architecture Diagram (Mermaid)

```mermaid
graph TD;
Client --> Gin;
Gin --> Handler;
Handler --> Service;
Service --> Repository;
Repository --> PostgreSQL;
Repository --> Redis;
Gin --> Prometheus;
Gin --> OpenTelemetry;
```

---

## Example Request

```bash
curl --location 'localhost:8080/tasks' \
--header 'Content-Type: application/json' \
--data '{
    "title": "learn golang",
    "description": "practice gin",
    "status": "todo",
    "assignee": "amir"
}'
```

---

## Tradeoffs

- In-memory tracing exporter for interview simplicity
- Cache invalidation strategy optimized for list endpoint