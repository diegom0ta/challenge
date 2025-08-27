.PHONY: up down migrate clean logs run-api run-cli build test

# Start all services
up:
	docker compose up -d

# Stop all services
down:
	docker compose down

# Stop and remove all containers, networks, and volumes
clean:
	docker compose down -v --remove-orphans

# Run migrations only
migrate:
	docker compose up flyway

# View logs
logs:
	docker compose logs -f

# View postgres logs
logs-db:
	docker compose logs -f postgres

# View flyway logs
logs-flyway:
	docker compose logs -f flyway

# Connect to postgres
psql:
	docker compose exec postgres psql -U postgres -d challenge

# Restart services
restart:
	docker compose restart

# Run API server
run-api:
	go run cmd/api/main.go

# Run CLI with sample data
run-cli:
	go run cmd/cli/main.go sample_data.csv

# Build all binaries
build:
	go build -o bin/api cmd/api/main.go
	go build -o bin/cli cmd/cli/main.go

# Run tests
test:
	go test ./...

# Install dependencies
deps:
	go mod tidy
