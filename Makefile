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

.PHONY: mocks
mocks:
	@echo "Generating mocks..."
	@mkdir -p mocks/repository 
	@mockgen -source=internal/port/repository/user-repository.go -destination=mocks/repository/user_repository_mock.go -package=repository
	@mockgen -source=internal/port/repository/user-log-repository.go -destination=mocks/repository/user_log_repository_mock.go -package=repository
	@echo "Mocks generated successfully."

.PHONY: clean-mocks
clean-mocks:
	@echo "Cleaning mocks..."
	@rm -rf mocks/
	@echo "Mocks cleaned."

.PHONY: tests
tests:
	@echo "Running tests..."
	@go test -v ./...
	@echo "Tests completed. Coverage report generated at coverage.html."