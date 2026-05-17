# Load .env file if exists
ifneq (,$(wildcard .env))
	include .env
	export
endif

DB_URL=$DB_MS://$DB_USER:$DB_PASSWORD@$DB_HOST:$DB_PORT/$DB_NAME?sslmode=$DB_SSLMODE
MIGRATIONS_PATH=migration
MAIN_FILE=cmd/main.go

.PHONY: help install migrate-up migrate-down swagger run docker docker-down test coverage bench metrics pprof clean

help:
	@echo "Available commands:"
	@echo " make deps					- Install dependencies"
	@echo " make migrate-up				- Run database migrations up"
	@echo " make migrate-down			- Rollback last migration"
	@echo " make swagger        		- Generate swagger docs"
	@echo " make run            		- Run the application"
	@echo " make docker-compose-up		- Run containers"
	@echo " make docker-compose-down	- Stop containers"
	@echo " make docker-compose-buildup	- Run containers with --build"
	@echo " make test					- Run tests"
	@echo " make coverage			    - Generate coverage"
	@echo " make bench					- Run benchmarks"

deps:
	go mod tidy
	go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	go install github.com/swaggo/swag/cmd/swag@latest

migrate-up:
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" up

migrate-down:
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" down -all

swagger:
	swag init -g $(MAIN_FILE)

run:
	go run $(MAIN_FILE)

docker-compose-buildup:
	docker compose up --build

docker-compose-up:
	docker compose up

docker-compose-down:
	docker compose down

test:
	go test ./...

coverage:
	go test ./... -coverprofile=coverage.out
	go tool cover -func=coverage.out