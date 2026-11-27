.PHONY: help build run clean test fmt deps db db-setup db-reset

APP_NAME=GrandExchange
BINARY_DIR=bin
MAIN_PATH=./cmd/main.go
ELM_DIR=ui
ELM_SRC=$(ELM_DIR)/src/Main.elm
ELM_OUT=$(ELM_DIR)/dist/elm.js
ELM_PUBLIC=$(ELM_DIR)/public

include .env
export

help: ## Show available commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-12s\033[0m %s\n", $$1, $$2}'

build: ## Build the application
	@mkdir -p $(BINARY_DIR)
	go build -o $(BINARY_DIR)/$(APP_NAME) $(MAIN_PATH)

run: ## Run the application
	go run $(MAIN_PATH)

deps: ## Download dependencies
	go mod download && go mod tidy

fmt: ## Format code
	go fmt ./...

test: ## Run tests
	go test -v ./...

clean: ## Remove build artifacts
	rm -rf $(BINARY_DIR)

db: ## Connect to PostgreSQL database
	psql $(DB_URL)

db-setup: ## Create database tables
	psql $(DB_URL) -f schema.sql
	@echo "Database schema created successfully!"

db-reset: ## Drop and recreate tables
	psql $(DB_URL) -c "DROP TABLE IF EXISTS search_items CASCADE;"
	@$(MAKE) db-setup

elm-build: ## Build Elm frontend
	@mkdir -p $(ELM_DIR)/dist
	cd $(ELM_DIR) && elm make src/Main.elm --output=dist/elm.js
 
elm-watch: ## Build Elm in watch mode (requires elm-live)
	cd $(ELM_DIR) && elm-live src/Main.elm --dir=public -- --output=public/elm.js
 
elm-dev: ## Run backend + frontend in dev mode
	@echo "Starting backend on :8080 and frontend on :8000..."
	@make -j2 run elm-watch
 
elm-open: ## Open the frontend in browser
	@which xdg-open > /dev/null && xdg-open $(ELM_PUBLIC)/index.html || open $(ELM_PUBLIC)/index.html
 
elm-clean: ## Remove Elm build artifacts
	rm -rf $(ELM_DIR)/dist $(ELM_DIR)/elm-stuff
 
dev: ## Run both backend and frontend
	@echo "Start backend: make run"
	@echo "Start frontend: make elm-watch"
	@echo "Or run both: make elm-dev"
