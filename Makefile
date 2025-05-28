include .env

.PHONY: setup
setup:
	@echo "Setting up environment..."
	@go mod tidy
	@docker compose up -d
	@echo "Environment setup complete."

.PHONY: build
build:
	@rm -rf ./bin
	@go build -o ./bin/app ./cmd/main.go

.PHONY: run
run: build
	@./bin/app

# Usage: make migration-create name=your_migration_name
.PHONY: migration-create
migration-create:
	@if [ -z "$(name)" ]; then \
		echo "Error: Migration name is required."; \
		exit 1; \
	fi
	goose create $(name) --dir database/migrations sql

.PHONY: migration-up
migration-up:
	goose postgres $(POSTGRES_DB_URL) --dir database/migrations up

.PHONY: migration-down
migration-down:
	goose postgres $(POSTGRES_DB_URL) --dir database/migrations down

.PHONY: migration-status
migration-status:
		goose postgres $(POSTGRES_DB_URL) --dir database/migrations status

.PHONY: docs
docs:
	@swag init -g ./cmd/main.go 

.PHONY: seed
seed:
	@echo "Seeding users..."
	@go run ./database/seeders/main.go