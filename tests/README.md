# WaterMark Project Test Suite

This directory contains test entry files and utilities for running tests across the entire WaterMark project.

## Test Coverage

The project currently includes tests for the following packages:

- `pkg/` - Core utility packages
  - `pkg/any_test.go` - Type conversion and array membership tests
  - `pkg/compress_test.go` - ZLIB compression/decompression tests
  - `pkg/csv_test.go` - CSV generation with BOM handling tests
  - `pkg/eerror_test.go` - Custom error type tests
  - `pkg/image_test.go` - Image processing and manipulation tests

- `internal/` - Internal packages
  - `internal/zip_test.go` - ZIP file handling tests
  - `internal/zip_7z_test.go` - 7z file handling tests

- `scripts/` - Script utilities
  - `scripts/tool_test.go` - Tool function tests

## Running Tests

### Using Batch Script (Windows)

```bash
# Run all tests
tests\run_tests.bat

# Run specific package tests
tests\run_tests.bat pkg
tests\run_tests.bat internal
tests\run_tests.bat scripts

# Run tests with coverage
tests\run_tests.bat cover

# Run tests with race detector
tests\run_tests.bat race

# Show help
tests\run_tests.bat help
```

### Using Shell Script (Linux/Mac)

```bash
# Make script executable first
chmod +x tests/run_tests.sh

# Run all tests
./tests/run_tests.sh

# Run specific package tests
./tests/run_tests.sh pkg
./tests/run_tests.sh internal
./tests/run_tests.sh scripts

# Run tests with coverage
./tests/run_tests.sh cover

# Run tests with race detector
./tests/run_tests.sh race

# Show help
./tests/run_tests.sh help
```

### Using Make

```bash
# Run all tests
make -C tests test

# Run specific package tests
make -C tests test-pkg
make -C tests test-internal
make -C tests test-scripts

# Run tests with coverage
make -C tests test-cover

# Run tests with race detector
make -C tests test-race

# Clean test artifacts
make -C tests clean

# Show help
make -C tests help
```

### Using Go Directly

```bash
# Run all tests
go test ./... -v

# Run specific package tests
go test ./pkg/... -v
go test ./internal/... -v
go test ./scripts/... -v

# Run tests with coverage
go test ./... -cover -coverprofile=coverage.out

# View coverage report
go tool cover -html=coverage.out

# Run tests with race detector
go test ./... -race -v
```

### Using Go Program

```bash
# Run all tests
go run tests/run_tests.go

# Run with verbose output
go run tests/run_tests.go -v

# Run with coverage
go run tests/run_tests.go -cover -coverprofile=coverage.out

# Run with race detector
go run tests/run_tests.go -race

# Run specific packages
go run tests/run_tests.go -pkgs ./pkg/...
```

## Test Structure

### Test Entry Files

- `all_test.go` - Main test entry point for the test suite
- `run_tests.go` - Go program for running tests with various options
- `run_tests.bat` - Windows batch script for running tests
- `run_tests.sh` - Unix shell script for running tests
- `Makefile` - Makefile for running tests with make commands

### Test Organization

Tests are organized by package and follow Go testing conventions:

```
pkg/
├── any.go
├── any_test.go
├── compress.go
├── compress_test.go
├── csv.go
├── csv_test.go
├── ecode.go
├── eerror.go
├── eerror_test.go
├── image.go
└── image_test.go

internal/
├── zip.go
├── zip_test.go
├── zip_7z.go
└── zip_7z_test.go

scripts/
├── tool.go
└── tool_test.go
```

## Adding New Tests

When adding new tests to the project:

1. Create test files alongside the source files using the `_test.go` suffix
2. Follow Go testing conventions and naming patterns
3. Use table-driven tests for multiple test cases
4. Include proper error handling and cleanup
5. Run the full test suite to ensure no regressions

Example test file structure:

```go
package pkg

import (
    "testing"
)

func TestNewFunction(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected string
    }{
        {
            name:     "test case 1",
            input:    "input",
            expected: "output",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := NewFunction(tt.input)
            if result != tt.expected {
                t.Errorf("Expected %s, got %s", tt.expected, result)
            }
        })
    }
}
```

## Continuous Integration

These test files are designed to work with CI/CD pipelines. The test commands can be integrated into:

- GitHub Actions
- GitLab CI
- Jenkins
- Other CI/CD systems

Example GitHub Actions workflow:

```yaml
name: Tests

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - uses: actions/setup-go@v2
        with:
          go-version: '1.24'
      - name: Run tests
        run: go test ./... -v -race -cover
```

## Troubleshooting

### Tests Fail to Run

If tests fail to run, ensure:

1. Go is properly installed and configured
2. All dependencies are installed: `go mod download`
3. Test files are in the correct packages
4. Test functions are properly named (TestXxx)

### Coverage Report Issues

If coverage reports are not generated:

1. Ensure tests are passing first
2. Check that the coverage output path is writable
3. Use `go tool cover -html=coverage.out` to view the report

### Race Detector Issues

If the race detector reports issues:

1. Review the reported race conditions
2. Add proper synchronization (mutexes, channels)
3. Fix the race conditions before proceeding
