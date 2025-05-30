MODULE_NAME := test-task-scout-go

APP_NAME := test-task-scout-go

BUILD_DIR := bin

MAIN_FILE := main.go

SQLITE_DB := quotes.db

.PHONY: all
all: build

.PHONY: build
build:
	@echo "Building $(APP_NAME)..."
	
	go build -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_FILE)
	@echo "Build complete. Executable in $(BUILD_DIR)/"

.PHONY: run
run: build
	@echo "Running $(APP_NAME)..."
	@export $$(cat .env | xargs) && $(BUILD_DIR)/$(APP_NAME)

.PHONY: test
test:
	@echo "Running all tests..."
	go test -v ./...

.PHONY: test-repository
test-repository:
	@echo "Running repository tests..."
	go test -v ./internal/repository

.PHONY: test-service
test-service:
	@echo "Running service tests..."
	go test -v ./internal/service

.PHONY: test-router
test-router:
	@echo "Running router tests..."
	go test -v ./internal/router

.PHONY: clean
clean:
	@echo "Cleaning build artifacts and database..."
	@rm -rf $(BUILD_DIR)
	@rm -f $(SQLITE_DB)
	@echo "Clean complete."


SCRIPTS_DIR := scripts

.PHONY: scripts-executable
scripts-executable:
	@chmod +x $(SCRIPTS_DIR)/*.sh

.PHONY: run-create-quotes
run-create-quotes: scripts-executable
	@echo "Running create_quotes.sh..."
	@$(SCRIPTS_DIR)/create_quotes.sh

.PHONY: run-get-all-quotes
run-get-all-quotes: scripts-executable
	@echo "Running get_all_quotes.sh..."
	@$(SCRIPTS_DIR)/get_all_quotes.sh

.PHONY: run-get-quotes-by-author
run-get-quotes-by-author: scripts-executable
	@echo "Running get_quotes_by_author.sh..."
	@$(SCRIPTS_DIR)/get_quotes_by_author.sh "Steve Jobs"

.PHONY: run-get-random-quote
run-get-random-quote: scripts-executable
	@echo "Running get_random_quote.sh..."
	@$(SCRIPTS_DIR)/get_random_quote.sh

.PHONY: run-get-by-id
run-get-by-id: scripts-executable
	@echo "Running get_by_id.sh..."
	@$(SCRIPTS_DIR)/get_by_id.sh $(ID)

.PHONY: run-delete-quote
run-delete-quote: scripts-executable
	@echo "Running delete_quote.sh..."
	@$(SCRIPTS_DIR)/delete_quote.sh $(ID)

.PHONY: run-all-scripts
run-all-scripts: scripts-executable
	@echo "--- Running all API interaction scripts ---"
	@echo "1. Creating quotes..."
	@make run-create-quotes

	@echo "\n2. Getting all quotes..."
	@make run-get-all-quotes

	@echo "\n3. Getting quotes by author (Steve Jobs)..."
	@make run-get-quotes-by-author

	@echo "\n4. Getting a random quote..."
	@make run-get-random-quote

	@echo "\n--- All API interaction scripts finished ---"

.PHONY: help
help:
	@echo "Available targets:"
	@echo "  all: Build the application (default)"
	@echo "  build: Build the application executable"
	@echo "  run: Build and run the application (in-memory or sqlite based on .env)"
	@echo "  run-sqlite: Build and run the application (explicitly with SQLite DB path)"
	@echo "  test: Run all tests"
	@echo "  test-repository: Run repository tests"
	@echo "  test-service: Run service tests"
	@echo "  test-router: Run router tests"
	@echo "  clean: Remove build artifacts and database file"
	@echo "  run-create-quotes: Run script to create example quotes"
	@echo "  run-get-all-quotes: Run script to get all quotes"
	@echo "  run-get-quotes-by-author: Run script to get quotes by author"
	@echo "  run-get-random-quote: Run script to get a random quote"
	@echo "  run-get-by-id: Run script to get a quote by ID (requires ID=...)"
	@echo "  run-delete-quote: Run script to delete a quote by ID (requires ID=...)"
	@echo "  run-all-scripts: Run all basic API interaction scripts sequentially"
	@echo "  scripts-executable: Make all scripts in scripts/ executable"
	@echo "  help: Display this help message" 