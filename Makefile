# ==============================================================================
# Variables
# ==============================================================================
# Define the docker compose command to use
DOCKER_COMPOSE := docker compose

# Define service names for easier reference
SERVICE_API := api
SERVICE_DB := postgres
SERVICE_REDIS := redis

# ==============================================================================
# Main Targets
# ==============================================================================
# Default target: Executed when you run `make`
# This will build, start the services, and then run the tests.
all: build up

# A complete reset: stops containers, removes volumes, and starts fresh.
reset: down up

# ==============================================================================
# Docker Lifecycle Management
# ==============================================================================
# Build or rebuild service images
build:
	@echo "Building Docker images..."
	$(DOCKER_COMPOSE) build

# Start all services, wait for the database, and then run migrations
up:
	@echo "Starting all services..."
	$(DOCKER_COMPOSE) up -d
	@echo "Waiting 10 seconds for the database to be ready..."
	@sleep 10
	@echo "Applying database migrations..."
	$(MAKE) migrate-up

# Stop and remove all containers, networks, and volumes
down:
	@echo "Stopping and removing all services and volumes..."
	$(DOCKER_COMPOSE) down -v --remove-orphans

# Restart all services
restart:
	@echo "Restarting services..."
	$(DOCKER_COMPOSE) restart

# View the status of all running services
ps:
	@echo "Current status of services:"
	$(DOCKER_COMPOSE) ps

# Follow the logs of all services in real-time
logs:
	@echo "Following logs..."
	$(DOCKER_COMPOSE) logs -f

# ==============================================================================
# Testing and CI
# ==============================================================================
# Wait for the API service to become healthy before proceeding
wait-ready:
	@echo "Waiting for API service to be ready..."
	@until curl -s -f http://localhost:8080/swagger/index.html > /dev/null; do \
		printf "."; \
		sleep 1; \
	done
	@echo "✅ API is ready!"

# Run all Go tests inside the API container
test: wait-ready
	@echo "🧪 Running tests inside the container..."
	$(DOCKER_COMPOSE) exec $(SERVICE_API) go test -v -cover ./...

# Run the linter to check code quality
lint:
	@echo "Running linter..."
	docker run --rm -v $(shell pwd):/app -w /app golangci/golangci-lint:v1.55 golangci-lint run -v

# ==============================================================================
# Development & Utilities
# ==============================================================================
# Run the application locally (without Docker)
start-local:
	@echo "Starting application locally..."
	go run main.go

# Generate Go code from SQL using sqlc
sqlc:
	@echo "Generating Go code from SQL..."
	docker-compose run --rm $(SERVICE_API) sqlc generate

# Open an interactive shell inside the running API container
shell-api:
	@echo "Connecting to the API container shell..."
	$(DOCKER_COMPOSE) exec $(SERVICE_API) /bin/sh

# Connect to the PostgreSQL database using psql
psql:
	@echo "Connecting to the PostgreSQL database..."
	$(DOCKER_COMPOSE) exec $(SERVICE_DB) psql -U root -d bilitioo

# Connect to the Redis database using redis-cli
redis-cli:
	@echo "Connecting to Redis..."
	$(DOCKER_COMPOSE) exec $(SERVICE_REDIS) redis-cli

# ==============================================================================
# Database Migrations
# ==============================================================================
# Apply all available database migrations (runs on the host machine)
migrate-up:
	@echo "Applying database migrations..."
	migrate -path db/migrate -database "postgresql://root:secret@localhost:5432/bilitioo?sslmode=disable" -verbose up

# Roll back the last database migration (runs on the host machine)
migrate-down:
	@echo "Rolling back the last migration..."
	migrate -path db/migrate -database "postgresql://root:secret@localhost:5432/bilitioo?sslmode=disable" -verbose down 1

# Create a new migration file. Usage: make new-migration name=add_users_table
new-migration:
	@if [ -z "$(name)" ]; then \
        echo "Usage: make new-migration name=<migration_name>"; \
        exit 1; \
    fi
	@echo "Creating new migration file: $(name)..."
	migrate create -ext sql -dir db/migrate -seq $(name)

# Define all targets as phony to prevent conflicts with file names
.PHONY: all reset build up down restart ps logs wait-ready test lint start-local sqlc shell-api psql redis-cli migrate-up migrate-down new-migration
