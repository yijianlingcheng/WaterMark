@echo off
REM tools.bat - minimal ASCII-safe dispatcher: check / api / dev / build
REM Usage: tools.bat <check|api|dev|build>
setlocal enabledelayedexpansion

for %%I in ("%CD%\..") do set "WORKSPACE=%%~fI"
if "%WORKSPACE:~-1%"=="\" set "WORKSPACE=%WORKSPACE:~0,-1%"
cd %WORKSPACE%
:: If no argument, show usage and exit
if "%~1"=="" goto :usage

set "ACTION=%~1"
shift

if /I "%ACTION%"=="check" goto :check
if /I "%ACTION%"=="api" goto :api
if /I "%ACTION%"=="dev" goto :dev
if /I "%ACTION%"=="build" goto :build

echo Unknown action: "%ACTION%"
goto :usage

:maybe_delegate_check
call :maybe_delegate
goto :check

:maybe_delegate_api
call :maybe_delegate
goto :api

:maybe_delegate_dev
call :maybe_delegate
goto :dev

:maybe_delegate_build
call :maybe_delegate
goto :build

:: -- batch implementations (ASCII messages) --

:check
set "REQUIRED=golangci-lint"
call :ensure_tools
echo ==================================================
echo Running: golangci-lint run
echo ==================================================
golangci-lint run
exit /B %ERRORLEVEL%

:api
set "REQUIRED=swag"
call :ensure_tools
echo ==================================================
echo Step 1: swag init
echo ==================================================
swag init
if ERRORLEVEL 1 (
  echo swag init failed, code %ERRORLEVEL%
  exit /B %ERRORLEVEL%
)
echo ==================================================
echo Step 2: change mode
echo ==================================================
go run .\scripts\tool.go -appMode=api-dev
if ERRORLEVEL 1 (
  echo change mode failed, code %ERRORLEVEL%
  exit /B %ERRORLEVEL%
)
echo ==================================================
echo Step 3: go run main.go
echo ==================================================
go run .\main.go
exit /B %ERRORLEVEL%

:dev
set "REQUIRED=swag wails"
call :ensure_tools
echo ==================================================
echo Step 1: swag init
echo ==================================================
swag init
if ERRORLEVEL 1 (
  echo swag init failed, code %ERRORLEVEL%
  exit /B %ERRORLEVEL%
)
echo ==================================================
echo Step 2: change mode
echo ==================================================
go run .\scripts\tool.go -appMode=dev
if ERRORLEVEL 1 (
  echo change mode failed, code %ERRORLEVEL%
  exit /B %ERRORLEVEL%
)
echo ==================================================
echo Step 3: wails dev
echo ==================================================
echo Note: wails dev will block this console. Use: start "" wails dev to run in new window.
wails dev
exit /B %ERRORLEVEL%

:build
set "REQUIRED=swag golangci-lint wails"
call :ensure_tools
echo ==================================================
echo Step 1: swag init
echo ==================================================
swag init
if ERRORLEVEL 1 (
  echo swag init failed, code %ERRORLEVEL%
  exit /B %ERRORLEVEL%
)
echo ==================================================
echo Step 2: golangci-lint run
echo ==================================================
golangci-lint run
if ERRORLEVEL 1 (
  echo golangci-lint run failed, code %ERRORLEVEL%
  exit /B %ERRORLEVEL%
)
echo ==================================================
echo Step 3: change mode
echo ==================================================
go run .\scripts\tool.go -appMode=release
if ERRORLEVEL 1 (
  echo change mode failed, code %ERRORLEVEL%
  exit /B %ERRORLEVEL%
)
echo ==================================================
echo Step 4: parse version file
echo ==================================================
set "VERSION_FILE=%~dp0\version"
if not exist "%VERSION_FILE%" (
  echo [Error] version file not found: "%VERSION_FILE%"
  exit /B 2
)
set "APP_VERSION="
set "FOUND=0"
for /f "usebackq delims=" %%L in ("%VERSION_FILE%") do (
  for /f "tokens=* delims= " %%A in ("%%L") do (
    set "LINE=%%A"
    if not "!LINE!"=="" (
      set "KEY="
      set "VAL="
      rem try split by '=' first (tokens=1* -> %%B is key, %%C is rest)
      for /f "tokens=1* delims==" %%B in ("!LINE!") do (
        set "KEY=%%B"
        set "VAL=%%C"
      )
      rem if VAL empty, try ':' as delimiter
      if "!VAL!"=="" (
        for /f "tokens=1* delims=:" %%B in ("!LINE!") do (
          set "KEY=%%B"
          set "VAL=%%C"
        )
      )
      rem Trim KEY and VAL (use unique loop vars)
      if defined KEY for /f "tokens=* delims= " %%D in ("!KEY!") do set "KEY=%%D"
      if defined VAL for /f "tokens=* delims= " %%E in ("!VAL!") do set "VAL=%%E"
      rem remove surrounding double quotes from VAL if any
      if defined VAL set "VAL=!VAL:"=!"
      if /I "!KEY!"=="APP_VERSION" (
        set "APP_VERSION=!VAL!"
        set "FOUND=1"
        goto :version_parsed
      )
    )
  )
)
:version_parsed
if "!FOUND!"=="0" (
  echo [Error] APP_VERSION not found or cannot be parsed.
  exit /B 3
)
for /f "tokens=* delims= " %%X in ("!APP_VERSION!") do set "APP_VERSION=%%X"
echo Parsed APP_VERSION = "!APP_VERSION!"
if "!APP_VERSION!"=="" (
  echo [Error] APP_VERSION empty.
  exit /B 4
)
set "OUTNAME=win10_amd64_WaterMark_!APP_VERSION!.exe"
echo ==================================================
echo Step 5: wails build -clean -o "%OUTNAME%"
echo ==================================================
wails build -clean -o "%OUTNAME%"
exit /B %ERRORLEVEL%

:check_tool
set "TOOL=%~1"
where "%TOOL%" >nul 2>&1
if errorlevel 1 (
  echo.
  echo [Missing] %TOOL% not found in PATH.
  echo Install suggestion:
  if /I "%TOOL%"=="golangci-lint" (
    echo   go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
  ) else if /I "%TOOL%"=="swag" (
    echo   go install github.com/swaggo/swag/cmd/swag@latest
  ) else if /I "%TOOL%"=="wails" (
    echo   go install github.com/wailsapp/wails/v2/cmd/wails@latest
  )
  set /A MISSING+=1
) else (
  echo [Check] %TOOL% found.
)
exit /B 0

:ensure_tools
set /A MISSING=0
for %%T in (%REQUIRED%) do (
  call :check_tool %%T
)
if %MISSING% GTR 0 (
  echo.
  echo %MISSING% missing dependencies detected. Please install them and retry.
  exit /B 10
)
exit /B 0

:usage
echo Usage: %~nx0 ^<check^|api^|dev^|build^>
echo.
echo   check      - run: golangci-lint run
echo   api        - run: swag init then ^change appMode to api then exec go run main.go
echo   dev        - run: swag init then ^change appMode to dev then wails dev
echo   build      - run: swag init -^> golangci-lint run -^> ^change appMode to release -^> parse version -^> wails build -clean -o win10_amd64_WaterMark_{APP_VERSION}.exe
exit /B 1