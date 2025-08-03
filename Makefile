
ENV_FILE := .env
MIGRATIONS_DIR=./migrations

help:
	@echo "Usage:"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

## run: run golang app
run:
	go run ./cmd/api

# --- golang-migrate Commands ---

## check-migrate: Ensure migrate tool is installed
check-migrate:
	@command -v migrate >/dev/null 2>&1 || \
	{ echo >&2 "ERROR: golang-migrate is not installed. Install it: https://github.com/golang-migrate/migrate/tree/master/cmd/migrate"; exit 1; }


## migrate-add: add new migration files
migrate-add:
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir $(MIGRATIONS_DIR) -seq $$name

## migrate-up: apply all pending migrations
migrate-up:
	@if [ -z "$(URL)" ]; then \
		echo "❌ Error: You must provide the database URL as URL=..."; \
		exit 1; \
	fi
	migrate -path $(MIGRATIONS_DIR) -database "$(URL)" up

## migrate-down: revert one migration
migrate-down:
	@if [ -z "$(URL)" ]; then \
		echo "❌ Error: You must provide the database URL as URL=..."; \
		exit 1; \
	fi
	migrate -path $(MIGRATIONS_DIR) -database "$(URL)" down 1

## migrate-force: force migration
migrate-force:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" force

## migrate-drop: reset db
migrate-drop:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" drop -f

## migrate-status: show current migration version
migrate-status:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" version

.PHONY: \
  help \
  migrate-add \
  migrate-up \
  migrate-down \
  migrate-force \
  migrate-drop \
  migrate-status \
  check-migrate \
  run