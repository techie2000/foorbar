.PHONY: help build run test clean migrate-up migrate-down docker-up docker-down

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the backend application
	cd backend && go build -o bin/api cmd/api/main.go

run: ## Run the backend application
	cd backend && go run cmd/api/main.go

test: ## Run tests
	cd backend && go test -v ./...

test-coverage: ## Run tests with coverage
	cd backend && go test -v -coverprofile=coverage.out ./...
	cd backend && go tool cover -html=coverage.out -o coverage.html

clean: ## Clean build artifacts
	rm -rf backend/bin
	rm -rf backend/tmp
	rm -f backend/coverage.out backend/coverage.html

docker-up: ## Start all services with Docker Compose
	docker-compose up -d

docker-down: ## Stop all services
	docker-compose down

docker-logs: ## Show logs from all services
	docker-compose logs -f

migrate-create: ## Create a new migration (usage: make migrate-create name=create_users_table)
	migrate create -ext sql -dir backend/migrations -seq $(name)

migrate-up: ## Run database migrations
	migrate -path backend/migrations -database "postgresql://axiom:axiom@localhost:5432/axiom?sslmode=disable" up

migrate-down: ## Rollback database migrations
	migrate -path backend/migrations -database "postgresql://axiom:axiom@localhost:5432/axiom?sslmode=disable" down

migrate-force: ## Force migration version (usage: make migrate-force version=1)
	migrate -path backend/migrations -database "postgresql://axiom:axiom@localhost:5432/axiom?sslmode=disable" force $(version)

swagger: ## Generate Swagger documentation
	cd backend && swag init -g cmd/api/main.go -o docs

lint: ## Run linter
	cd backend && golangci-lint run

fmt: ## Format code
	cd backend && go fmt ./...

install-tools: ## Install development tools
	go install github.com/swaggo/swag/cmd/swag@latest
	go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
