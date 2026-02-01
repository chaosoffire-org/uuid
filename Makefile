# UUID Project Makefile
# Temporary files are stored in ./bin

BIN_DIR := ./bin

# OS detection for executable extension
ifeq ($(OS),Windows_NT)
    EXE_EXT := .exe
else
    EXE_EXT :=
endif

GOLANGCI_LINT := $(BIN_DIR)/golangci-lint$(EXE_EXT)
FUZZ_RUNNER := $(BIN_DIR)/fuzz-runner$(EXE_EXT)

.PHONY: all setup format lint build test benchmark benchmark-compare fuzz fuzz-list fuzz-single clean help

# Default target
all: format lint build test

# Create bin directory
setup:
	@mkdir -p $(BIN_DIR)

# Format code
format:
	@echo "==> Formatting code..."
	go fmt ./...

# Install golangci-lint to bin directory
$(GOLANGCI_LINT): setup
	@echo "==> Installing golangci-lint..."
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(BIN_DIR) latest

# Lint code
lint: $(GOLANGCI_LINT)
	@echo "==> Running linter..."
	$(GOLANGCI_LINT) run ./...

# Build project
build:
	@echo "==> Building..."
	go build ./...

# Run tests
test:
	@echo "==> Running tests..."
	go test ./...

# Run tests with coverage
test-cover: setup
	@echo "==> Running tests with coverage..."
	go test -coverprofile=$(BIN_DIR)/coverage.out ./...
	go tool cover -html=$(BIN_DIR)/coverage.out -o $(BIN_DIR)/coverage.html
	@echo "Coverage report: $(BIN_DIR)/coverage.html"

# Run benchmarks
benchmark:
	@echo "==> Running benchmarks..."
	go test -bench=. -benchmem ./...

# Run comparison benchmarks
benchmark-compare:
	@echo "==> Running comparison benchmarks..."
	cd benchmarks && go test -bench=. -benchmem

# Build fuzz runner
$(FUZZ_RUNNER): setup
	@echo "==> Building fuzz runner..."
	cd scripts && go build -o ../$(FUZZ_RUNNER) .

# Run all fuzz tests
fuzz: $(FUZZ_RUNNER)
	@echo "==> Running fuzz tests..."
	$(FUZZ_RUNNER) run

# List available fuzz tests
fuzz-list: $(FUZZ_RUNNER)
	@echo "==> Listing fuzz tests..."
	$(FUZZ_RUNNER) list

# Run single fuzz test (usage: make fuzz-single NAME=FuzzParse)
# Run single fuzz test (usage: make fuzz-single NAME=FuzzParse)
fuzz-single: $(FUZZ_RUNNER)
	@echo "==> Running fuzz test: $(NAME)..."
	$(FUZZ_RUNNER) run --filter $(NAME)

# Run fuzz with custom time (usage: make fuzz-time TIME=30s)
fuzz-time: $(FUZZ_RUNNER)
	@echo "==> Running fuzz tests for $(TIME)..."
	$(FUZZ_RUNNER) run --fuzztime $(TIME)

# Clean temporary files
clean:
	@echo "==> Cleaning..."
	rm -rf $(BIN_DIR)

# Help
help:
	@echo "Available targets:"
	@echo "  setup           - Create bin directory"
	@echo "  format          - Format Go code"
	@echo "  lint            - Run golangci-lint"
	@echo "  build           - Build the project"
	@echo "  test            - Run unit tests"
	@echo "  test-cover      - Run tests with coverage report"
	@echo "  benchmark       - Run all benchmarks"
	@echo "  benchmark-compare - Run comparison benchmarks"
	@echo "  fuzz            - Run all fuzz tests"
	@echo "  fuzz-list       - List available fuzz tests"
	@echo "  fuzz-single     - Run single fuzz test (NAME=FuzzName)"
	@echo "  fuzz-time       - Run fuzz with custom time (TIME=30s)"
	@echo "  clean           - Remove bin directory"
	@echo "  help            - Show this help"
