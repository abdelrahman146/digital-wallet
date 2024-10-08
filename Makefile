.PONY: dev build run migrate-up migrate-down create-migration swag test coverage gen-mocks

include .env
export $(shell sed 's/=.*//' .env)

dev:
	@echo "start digital wallet engine..."
	@go run main.go

build:
	@echo "building digital wallet engine..."
	@go build -o bin/digital-wallet main.go
	@echo "build completed"

run:
	@echo "running digital wallet engine..."
	@./bin/digital-wallet

create-migration:
	@echo "creating migration..."
	@migrate create -ext sql -dir db/migrations -seq $(name)

migrate-up:
	@echo "migrating up..."
	@migrate -path db/migrations -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" up

migrate-down:
	@echo "migrating down..."
	@migrate -path db/migrations -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" down

seed:
	@echo "seeding data..."
	@go run db/seed.go

swag:
	@echo "generating swagger docs..."
	@swag init --pdl 3 --parseDependency --parseInternal
	@echo "swagger docs generated"

test:
	@echo "running tests..."
	@go test ./...
	@echo "tests completed"

coverage:
	@echo "running tests with coverage..."
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out
	@echo "coverage completed"

gen-mocks:
	@echo "generating mocks..."
	@chmod +x ./scripts/generate-mocks.sh
	@./scripts/generate-mocks.sh
	@echo "mocks generated"