@echo off
REM Test runner for WaterMark project

setlocal enabledelayedexpansion

echo ========================================
echo WaterMark Project Test Runner
echo ========================================
echo.

if "%1"=="help" goto :help
if "%1"=="all" goto :all
if "%1"=="pkg" goto :pkg
if "%1"=="internal" goto :internal
if "%1"=="scripts" goto :scripts
if "%1"=="cover" goto :cover
if "%1"=="race" goto :race
if "%1"=="-v" goto :custom
if "%1"=="-cover" goto :custom
if "%1"=="-race" goto :custom

:default
echo Running all tests...
go test ./... -v
goto :end

:all
echo Running all tests...
go test ./... -v
goto :end

:pkg
echo Running pkg package tests...
go test ./pkg/... -v
goto :end

:internal
echo Running internal package tests...
go test ./internal/... -v
goto :end

:scripts
echo Running scripts package tests...
go test ./scripts/... -v
goto :end

:cover
echo Running tests with coverage...
go test ./... -cover -coverprofile=%CD%\coverage.out
echo.
echo Coverage report generated: %CD%\coverage.out
echo To view coverage report, run: go tool cover '-html=%CD%\coverage.out'
goto :end

:race
echo Running tests with race detector...
go test ./... -race -v
goto :end

:custom
echo Running tests with custom arguments...

go test ./... %*
echo.
echo Test run completed

REM Check if coverage file was generated and rename it if needed
if exist coverage (
    REM Delete existing coverage.out file if it exists
    if exist coverage.out (
        del coverage.out
    )
    echo Renaming coverage file from "coverage" to "coverage.out"...
    ren coverage coverage.out
    echo Coverage file renamed to: coverage.out
)

goto :end

:help
echo Usage: run_tests.bat [command] [options]
echo.
echo Commands:
echo   (none)    - Run all tests
echo   all       - Run all tests
echo   pkg       - Run pkg package tests
echo   internal  - Run internal package tests
echo   scripts   - Run scripts package tests
echo   cover     - Run tests with coverage
echo   race      - Run tests with race detector
echo   help      - Show this help message
echo.
echo Custom options (passed directly to go test):
echo   -v                - Verbose output
echo   -cover            - Enable coverage
echo   -coverprofile     - Specify coverage profile file
echo   -race             - Enable race detector
echo   -run              - Run tests matching pattern
echo.
echo Examples:
echo   run_tests.bat
echo   run_tests.bat -v
echo   run_tests.bat -cover -coverprofile=coverage.out
echo   run_tests.bat -cover -coverprofile=coverage.out -v
echo   run_tests.bat -run TestLoadImageWithDecode
echo.
goto :end

:end
echo.
echo ========================================
echo Test run completed
echo ========================================
