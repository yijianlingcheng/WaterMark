#!/usr/bin/env bash
# tools.sh - cross-platform POSIX-style script for Linux/macOS
# Usage: ./tools.sh <check|api|dev|build|check-env>
# - check      : run golangci-lint run
# - api        : swag init && change appMode to api && go run main.go
# - dev        : swag init && change appMode to dev && wails dev
# - build      : swag init -> golangci-lint run -> change appMode to release -> parse version file -> wails build -clean -o waterMark_{APP_VERSION}
# - check-env  : check environment tools: go, wails, exiftool, ImageMagick (magick or convert)
#
# Version file:
#   - placed in the same directory as this script, named: version
#   - supported formats:
#       APP_VERSION = v1.0.3.Releases
#       APP_VERSION:v1.0.3.Releases
#       APP_VERSION = "v1.0.3.Releases"
#   - empty lines and surrounding whitespace are ignored

set -o errexit
set -o nounset
set -o pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" >/dev/null 2>&1 && pwd)"
WORKSPACE_DIR="$(cd "$(dirname "${SCRIPT_DIR}")" >/dev/null 2>&1 && pwd)"
VERSION_FILE="${WORKSPACE_DIR}/version"
cd ${WORKSPACE_DIR}

err() { printf '%s\n' "$*" >&2; }
info() { printf '%s\n' "$*"; }

usage() {
  cat <<'USAGE'
Usage: tools.sh <check|api|dev|build|check-env>

  check      - run: golangci-lint run
  api        - run: swag init && change appMode to api && go run main.go
  dev        - run: swag init && change appMode to dev && wails dev (will block the terminal)
  build      - swag init -> golangci-lint run -> parse version ->
               wails build -clean -o waterMark_{APP_VERSION}
  check-env  - check environment tools: go, wails, exiftool, ImageMagick (magick or convert)

Notes:
  - Place a file named "version" next to this script containing a line like:
      APP_VERSION = v1.0.3.Releases
    or
      APP_VERSION:v1.0.3.Releases
    Quotes are allowed around the value and will be stripped.
  - The script checks for required tools and prints install hints when missing.
USAGE
}

# Check executable presence; collect missing and show hints
check_tools() {
  local missing=()
  for tool in "$@"; do
    if ! command -v "$tool" >/dev/null 2>&1; then
      missing+=("$tool")
      err "[Missing] $tool not found in PATH."
    else
      info "[Check] $tool found."
    fi
  done

  if [ "${#missing[@]}" -ne 0 ]; then
    err
    err "Suggested install commands (examples):"
    for t in "${missing[@]}"; do
      case "$t" in
        golangci-lint)
          err "  go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"
          err "  or use your package manager (e.g. brew install golangci-lint / apt)"
          ;;
        swag)
          err "  go install github.com/swaggo/swag/cmd/swag@latest"
          err "  or: brew install swag (if available) / apt / package manager"
          ;;
        wails)
          err "  go install github.com/wailsapp/wails/v2/cmd/wails@latest"
          err "  or: refer to https://wails.io for installation"
          ;;
        exiftool)
          err "  brew install exiftool"
          err "  or apt: apt install libimage-exiftool-perl"
          ;;
        magick|convert)
          err "  brew install imagemagick"
          err "  or apt: apt install imagemagick"
          ;;
        go)
          err "  Install Go from https://golang.org/dl/ or use your package manager (brew/apt/etc.)"
          ;;
        *)
          err "  Please install $t and ensure it's in PATH."
          ;;
      esac
    done
    return 1
  fi
  return 0
}

# Parse version file and extract APP_VERSION; supports separators '=' or ':'
# Trims whitespace and strips surrounding double quotes.
parse_version_file() {
  if [ ! -f "$VERSION_FILE" ]; then
    err "[Error] version file not found: $VERSION_FILE"
    return 2
  fi

  local line
  line="$(grep -iE '^[[:space:]]*APP_VERSION[[:space:]]*[:=]' "$VERSION_FILE" | head -n 1 || true)"
  if [ -z "$line" ]; then
    err "[Error] APP_VERSION not found in $VERSION_FILE"
    return 3
  fi

  local val
  # Extract RHS after '=' or ':', case-insensitive for APP_VERSION
  val="$(printf '%s' "$line" | sed -E 's/^[[:space:]]*APP_VERSION[[:space:]]*[:=][[:space:]]*(.*)$/\1/I')"
  val="$(printf '%s' "$val" | sed -E 's/^[[:space:]]+//; s/[[:space:]]+$//')"
  if printf '%s' "$val" | grep -qE '^".*"$'; then
    val="$(printf '%s' "$val" | sed -E 's/^"(.*)"$/\1/')"
  fi

  if [ -z "$val" ]; then
    err "[Error] Parsed APP_VERSION is empty."
    return 4
  fi

  printf '%s' "$val"
  return 0
}

run_or_exit() {
  # run command and exit on non-zero with message
  local desc="$1"; shift
  info ">>> $desc"
  if ! "$@"; then
    err "[Error] Command failed: $*"
    return 1
  fi
  return 0
}

check_env() {
  local ok=0

  # 1) Go
  if command -v go >/dev/null 2>&1; then
    info "[Check] go -> $(go version 2>/dev/null || echo 'version unknown')"
  else
    err "[Missing] go not found in PATH."
    err "  Install: https://golang.org/dl/ or brew install go / apt install golang-go"
    ok=1
  fi

  # 2) wails
  if command -v wails >/dev/null 2>&1; then
    info "[Check] wails -> $(wails version 2>/dev/null || echo 'installed')"
  else
    err "[Missing] wails not found in PATH."
    err "  Install: go install github.com/wailsapp/wails/v2/cmd/wails@latest"
    ok=1
  fi

  # 3) exiftool
  if command -v exiftool >/dev/null 2>&1; then
    info "[Check] exiftool -> $(exiftool -ver 2>/dev/null || echo 'installed')"
  else
    err "[Missing] exiftool not found in PATH."
    err "  Install: brew install exiftool  or apt install libimage-exiftool-perl"
    ok=1
  fi

  # 4) ImageMagick: prefer `magick` (v7), fallback to `convert` (v6)
  if command -v magick >/dev/null 2>&1; then
    info "[Check] ImageMagick -> magick present ($(magick -version | head -n1 2>/dev/null || echo 'version unknown'))"
  elif command -v convert >/dev/null 2>&1; then
    info "[Check] ImageMagick -> convert present ($(convert -version | head -n1 2>/dev/null || echo 'version unknown'))"
  else
    err "[Missing] ImageMagick (magick or convert) not found in PATH."
    err "  Install: brew install imagemagick  or apt install imagemagick"
    ok=1
  fi

  if [ "$ok" -ne 0 ]; then
    err
    err "Environment check failed: some tools are missing."
    return 10
  fi

  info
  info "Environment check passed: all required tools are present."
  return 0
}

main() {
  if [ $# -lt 1 ]; then
    usage
    exit 1
  fi

  local action="$1"; shift

  case "$action" in
    check)
      check_tools golangci-lint || exit 10
      run_or_exit "Running: golangci-lint run" golangci-lint run
      ;;
    
    api)
      check_tools swag || exit 10
      run_or_exit "Step 1: swag init" swag init || exit $?
      info "Step 2: change appMode to api"
      go run ./scripts/tool.go -appMode=api-dev
      info "Step 3: go run main.go"
      go run ./main.go
      ;;

    dev)
      check_tools swag wails || exit 10
      run_or_exit "Step 1: swag init" swag init || exit $?
      info "Step 2: change appMode to api"
      go run ./scripts/tool.go -appMode=dev
      info "Step 3: starting wails dev (will block this terminal)"
      # wails dev usually blocks; run directly
      wails dev
      ;;

    build)
      check_tools swag golangci-lint wails || exit 10
      run_or_exit "Step 1: swag init" swag init || exit $?
      run_or_exit "Step 2: golangci-lint run" golangci-lint run || exit $?
      info "Step 3: change appMode to api"
      go run ./scripts/tool.go -appMode=release
      info "Step 4: parsing version file: $VERSION_FILE"
      local app_ver
      app_ver="$(parse_version_file)" || exit $?
      info "Parsed APP_VERSION = ${app_ver}"
      local outname="waterMark_${app_ver}"
      run_or_exit "Step 5: wails build -clean -o ${outname}" wails build -clean -o "${outname}"
      ;;

    check-env)
      check_env || exit $?
      ;;

    *)
      err "Unknown action: $action"
      usage
      exit 2
      ;;
  esac
}

main "$@"