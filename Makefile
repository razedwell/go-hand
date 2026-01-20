include .env
export

DB_URL=postgresql://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)

.PHONY: migrate-up migrate-down migrate-create migrate-force db-reset

# Run all pending migrations
migrate-up:
	migrate -path ./migrations -database "$(DB_URL)" up

# Rollback last migration
migrate-down:
	migrate -path ./migrations -database "$(DB_URL)" down 1

# Create a new migration file
migrate-create:
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir ./migrations -seq $$name

# Force version (use carefully!)
migrate-force:
	@read -p "Enter version to force: " version; \
	migrate -path ./migrations -database "$(DB_URL)" force $$version

# Drop everything and re-migrate (DESTRUCTIVE - dev only!)
db-reset:
	migrate -path ./migrations -database "$(DB_URL)" drop -f
	migrate -path ./migrations -database "$(DB_URL)" up

# Check migration status
migrate-version:
	migrate -path ./migrations -database "$(DB_URL)" version