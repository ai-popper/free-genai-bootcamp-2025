.PHONY: help setup run clean test reset migrate seed init

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

setup: ## Setup database (init + migrate + seed)
	@echo "Setting up database..."
	@~/go/bin/mage setup

init: ## Initialize database file
	@echo "Initializing database..."
	@~/go/bin/mage init

migrate: ## Run database migrations
	@echo "Running migrations..."
	@~/go/bin/mage migrate

seed: ## Seed database with sample data
	@echo "Seeding database..."
	@~/go/bin/mage seed

run: ## Run the server
	@echo "Starting server..."
	@go run cmd/server/main.go

clean: ## Remove database file
	@echo "Cleaning database..."
	@~/go/bin/mage clean

reset: ## Clean and setup from scratch
	@echo "Resetting database..."
	@~/go/bin/mage reset

test: ## Run API tests
	@echo "Running API tests..."
	@./test_api.sh

deps: ## Download Go dependencies
	@echo "Downloading dependencies..."
	@go mod download

tidy: ## Tidy Go modules
	@echo "Tidying modules..."
	@go mod tidy

build: ## Build the server binary
	@echo "Building server..."
	@go build -o server cmd/server/main.go

install-mage: ## Install Mage task runner
	@echo "Installing Mage..."
	@go install github.com/magefile/mage@latest

dev: setup run ## Setup and run in one command
