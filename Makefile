.PHONY: up down migrate clean logs

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
