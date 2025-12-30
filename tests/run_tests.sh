#!/bin/bash

# Test runner for WaterMark project

echo "========================================"
echo "WaterMark Project Test Runner"
echo "========================================"
echo ""

case "$1" in
    help)
        echo "Usage: run_tests.sh [command] [options]"
        echo ""
        echo "Commands:"
        echo "  (none)    - Run all tests"
        echo "  all       - Run all tests"
        echo "  pkg       - Run pkg package tests"
        echo "  internal  - Run internal package tests"
        echo "  scripts   - Run scripts package tests"
        echo "  cover     - Run tests with coverage"
        echo "  race      - Run tests with race detector"
        echo "  help      - Show this help message"
        echo ""
        echo "Custom options (passed directly to go test):"
        echo "  -v                - Verbose output"
        echo "  -cover            - Enable coverage"
        echo "  -coverprofile     - Specify coverage profile file"
        echo "  -race             - Enable race detector"
        echo "  -run              - Run tests matching pattern"
        echo ""
        echo "Examples:"
        echo "  ./run_tests.sh"
        echo "  ./run_tests.sh -v"
        echo "  ./run_tests.sh -cover -coverprofile=coverage.out"
        echo "  ./run_tests.sh -cover -coverprofile=coverage.out -v"
        echo "  ./run_tests.sh -run TestLoadImageWithDecode"
        echo ""
        ;;
    all)
        echo "Running all tests..."
        go test ./... -v
        ;;
    pkg)
        echo "Running pkg package tests..."
        go test ./pkg/... -v
        ;;
    internal)
        echo "Running internal package tests..."
        go test ./internal/... -v
        ;;
    scripts)
        echo "Running scripts package tests..."
        go test ./scripts/... -v
        ;;
    cover)
        echo "Running tests with coverage..."
        go test ./... -cover -coverprofile=coverage.out
        echo ""
        echo "Coverage report generated: coverage.out"
        echo "To view coverage report, run: go tool cover -html=coverage.out"
        ;;
    race)
        echo "Running tests with race detector..."
        go test ./... -race -v
        ;;
    -v|-cover|-race|-run|-coverprofile=*)
        echo "Running tests with custom arguments..."
        go test ./... "$@"
        echo ""
        echo "Test run completed"
        
        # Check if coverage file was generated and rename it if needed
        if [ -f "coverage" ]; then
            # Delete existing coverage.out file if it exists
            if [ -f "coverage.out" ]; then
                rm coverage.out
            fi
            echo "Renaming coverage file from 'coverage' to 'coverage.out'..."
            mv coverage coverage.out
            echo "Coverage file renamed to: coverage.out"
        fi
        ;;
    *)
        echo "Running all tests..."
        go test ./... -v
        ;;
esac

echo ""
echo "========================================"
echo "Test run completed"
echo "========================================"
