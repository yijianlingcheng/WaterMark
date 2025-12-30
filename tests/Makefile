.PHONY: help test test-pkg test-internal test-scripts test-cover test-race clean

help:
	@echo "WaterMark Project Test Commands"
	@echo ""
	@echo "Usage:"
	@echo "  make help          - Show this help message"
	@echo "  make test          - Run all tests"
	@echo "  make test-pkg      - Run pkg package tests"
	@echo "  make test-internal - Run internal package tests"
	@echo "  make test-scripts  - Run scripts package tests"
	@echo "  make test-cover    - Run tests with coverage"
	@echo "  make test-race     - Run tests with race detector"
	@echo "  make clean         - Clean test artifacts"
	@echo ""

test:
	@echo "Running all tests..."
	@go test ./... -v

test-pkg:
	@echo "Running pkg package tests..."
	@go test ./pkg/... -v

test-internal:
	@echo "Running internal package tests..."
	@go test ./internal/... -v

test-scripts:
	@echo "Running scripts package tests..."
	@go test ./scripts/... -v

test-cover:
	@echo "Running tests with coverage..."
	@go test ./... -cover -coverprofile=coverage.out
	@echo "Coverage report generated: coverage.out"
	@echo "To view coverage report, run: go tool cover -html=coverage.out"

test-race:
	@echo "Running tests with race detector..."
	@go test ./... -race -v

clean:
	@echo "Cleaning test artifacts..."
	@rm -f coverage.out
	@rm -f *.test
	@echo "Clean completed"
