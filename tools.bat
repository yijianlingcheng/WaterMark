@echo off
REM tools.bat - minimal ASCII-safe dispatcher: check / dev / win-build
REM Usage: tools.bat <check|dev|win-build>
setlocal enabledelayedexpansion

:: If no argument, show usage and exit
if "%~1"=="" goto :usage

set "ACTION=%~1"
shift

if /I "%ACTION%"=="check" goto :check
if /I "%ACTION%"=="dev" goto :dev
if /I "%ACTION%"=="win-build" goto :winbuild

echo Unknown action: "%ACTION%"
goto :usage

:: Try to delegate to PowerShell script only for valid actions and only if tools.ps1 exists
:maybe_delegate
where powershell >nul 2>&1
if errorlevel 1 exit /B 0
if exist "%~dp0build.ps1" (
  echo Detected powershell.exe, delegating to tools.ps1...
  powershell -NoProfile -ExecutionPolicy Bypass -File "%~dp0build.ps1" %*
  exit /B %ERRORLEVEL%
)
exit /B 0

:maybe_delegate_check
call :maybe_delegate
goto :check

:maybe_delegate_dev
call :maybe_delegate
goto :dev

:maybe_delegate_winbuild
call :maybe_delegate
goto :winbuild

:: -- batch implementations (ASCII messages) --

:check
set "REQUIRED=golangci-lint"
call :ensure_tools
echo ==================================================
echo Running: golangci-lint run
echo ==================================================
golangci-lint run
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
echo Step 2: wails dev
echo ==================================================
echo Note: wails dev will block this console. Use: start "" wails dev to run in new window.
wails dev
exit /B %ERRORLEVEL%

:winbuild
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
echo Step 3: parse version file
echo ==================================================
set "VERSION_FILE=%~dp0version"
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
echo Step 4: wails build -clean -o "%OUTNAME%"
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
echo Usage: %~nx0 ^<check^|dev^|win-build^>
echo.
echo   check      - run: golangci-lint run
echo   dev        - run: swag init then wails dev
echo   win-build  - swag init -^> golangci-lint run -^> parse version -^> wails build -clean -o win10_amd64_WaterMark_{APP_VERSION}.exe
exit /B 1