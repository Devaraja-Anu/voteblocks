



ENV_FILE := .env
MIGRATIONS_DIR=./migrations

.PHONY: help
help:
	@echo "Usage:"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'


# --- golang-migrate Commands ---

## migrate-add: add new migration files
migrate-add:
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir ./migrations -seq $$name

## migrate-up: apply all pending migrations
migrate-up:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" up

## migrate-down: revert one migration
migrate-down:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" down 1

## migrate-force: force migration
migrate-force:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" force

## migrate-drop: reset db
migrate-drop:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" drop -f

## migrate-status: show current migration version
migrate-status:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" version
