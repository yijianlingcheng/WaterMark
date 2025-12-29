<#
tools.ps1 - PowerShell 版本的任务分发脚本（check / dev / win-build）
功能：
  - 在 PowerShell 环境下运行（也可由 .\build.bat 委派调用）
  - 检查依赖工具：golangci-lint, swag, wails
  - 支持 version 文件多格式解析：
      APP_VERSION = v1.0.3.Releases
      APP_VERSION:v1.0.3.Releases
    会忽略空行并自动 Trim 前后空白与可选引号
用法:
  PS> .\tools.ps1 check
  PS> .\tools.ps1 dev
  PS> .\tools.ps1 win-build
#>

[CmdletBinding()]
param(
  [Parameter(Mandatory=$true, Position=0)]
  [ValidateSet("check","dev","win-build")]
  [string]$Action
)

function Check-Tool {
  param([string]$Tool)
  $cmd = Get-Command $Tool -ErrorAction SilentlyContinue
  if (-not $cmd) {
    Write-Host "[缺失] 未找到可执行程序: $Tool" -ForegroundColor Yellow
    Write-Host "建议安装（示例）:" -ForegroundColor Yellow
    switch ($Tool.ToLower()) {
      "golangci-lint" {
        Write-Host "  go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"
        Write-Host "  或: choco install golangci-lint"
      }
      "swag" {
        Write-Host "  go install github.com/swaggo/swag/cmd/swag@latest"
        Write-Host "  或: choco install swag"
      }
      "wails" {
        Write-Host "  go install github.com/wailsapp/wails/v2/cmd/wails@latest"
        Write-Host "  或: choco install wails"
      }
      default {
        Write-Host "  请自行安装 $Tool 并确保其在 PATH 中"
      }
    }
    Write-Host "安装后请确保将 Go 的 bin 目录加入 PATH（例如: $env:USERPROFILE\go\bin）" -ForegroundColor Yellow
    return $false
  } else {
    Write-Host "[检查] $Tool -> 已找到"
    return $true
  }
}

function Ensure-Tools {
  param([string[]]$Tools)
  $missing = @()
  foreach ($t in $Tools) {
    if (-not (Check-Tool $t)) { $missing += $t }
  }
  if ($missing.Count -gt 0) {
    Write-Host "" 
    Write-Host "检测到 $($missing.Count) 个缺失的依赖: $($missing -join ', ')" -ForegroundColor Red
    throw "缺少依赖"
  }
}

function Parse-VersionFile {
  param([string]$Path)
  if (-not (Test-Path $Path)) {
    throw "找不到 version 文件: $Path"
  }
  $lines = Get-Content -Raw -Path $Path -ErrorAction Stop -Encoding UTF8
  # 使用换行分割，按行处理，过滤空行，Trim 两端空白
  foreach ($rawLine in $lines -split "`r?`n") {
    $line = $rawLine.Trim()
    if ([string]::IsNullOrWhiteSpace($line)) { continue }
    # 支持 "key = value" 或 "key:value"
    if ($line -match "^\s*APP_VERSION\s*[:=]\s*(.+)$") {
      $val = $Matches[1].Trim()
      # 去掉前后双引号（如果有）
      if ($val.StartsWith('"') -and $val.EndsWith('"')) {
        $val = $val.Substring(1, $val.Length - 2)
      }
      if ($val.Length -eq 0) {
        throw "解析到的 APP_VERSION 为空"
      }
      return $val
    }
  }
  throw "未能在 $Path 中找到 APP_VERSION 的值或无法解析（请确保包含类似 'APP_VERSION = v1.0.3.Releases' 或 'APP_VERSION:v1.0.3.Releases' 的行）"
}

try {
  switch ($Action) {
    'check' {
      Ensure-Tools -Tools @('golangci-lint')
      Write-Host "==================================================" -ForegroundColor Cyan
      Write-Host "执行: golangci-lint run" -ForegroundColor Cyan
      Write-Host "=================================================="
      & golangci-lint run
      exit $LASTEXITCODE
    }
    'dev' {
      Ensure-Tools -Tools @('swag','wails')
      Write-Host "==================================================" -ForegroundColor Cyan
      Write-Host "步骤1: 执行 swag init" -ForegroundColor Cyan
      Write-Host "=================================================="
      & swag init
      if ($LASTEXITCODE -ne 0) { throw "swag init 失败，错误码 $LASTEXITCODE" }

      Write-Host "==================================================" -ForegroundColor Cyan
      Write-Host "步骤2: 启动 wails dev" -ForegroundColor Cyan
      Write-Host "=================================================="
      Write-Host "注意: wails dev 会阻塞当前控制台。若需在新窗口运行可使用: Start-Process wails -ArgumentList 'dev' " -ForegroundColor Yellow
      & wails dev
      exit $LASTEXITCODE
    }
    'win-build' {
      Ensure-Tools -Tools @('swag','golangci-lint','wails')
      Write-Host "==================================================" -ForegroundColor Cyan
      Write-Host "步骤1: 执行 swag init" -ForegroundColor Cyan
      Write-Host "=================================================="
      & swag init
      if ($LASTEXITCODE -ne 0) { throw "swag init 失败，错误码 $LASTEXITCODE" }

      Write-Host "==================================================" -ForegroundColor Cyan
      Write-Host "步骤2: 执行 golangci-lint run" -ForegroundColor Cyan
      Write-Host "=================================================="
      & golangci-lint run
      if ($LASTEXITCODE -ne 0) { throw "golangci-lint run 返回错误，错误码 $LASTEXITCODE" }

      Write-Host "==================================================" -ForegroundColor Cyan
      Write-Host "步骤3: 解析同目录下的 version 文件以获取 APP_VERSION" -ForegroundColor Cyan
      Write-Host "支持格式: APP_VERSION = v1.0.3.Releases 或 APP_VERSION:v1.0.3.Releases（会忽略空行并 Trim）"
      Write-Host "=================================================="
      $versionFile = Join-Path -Path $PSScriptRoot -ChildPath 'version'
      $appVer = Parse-VersionFile -Path $versionFile
      Write-Host "解析到 APP_VERSION = '$appVer'"

      if ([string]::IsNullOrWhiteSpace($appVer)) { throw "APP_VERSION 为空，停止构建" }

      $outName = "win10_amd64_WaterMark_{0}.exe" -f $appVer
      Write-Host "==================================================" -ForegroundColor Cyan
      Write-Host "步骤4: 执行 wails build -clean -o `"$outName`"" -ForegroundColor Cyan
      Write-Host "=================================================="
      & wails build -clean -o $outName
      exit $LASTEXITCODE
    }
  }
} catch {
  Write-Host "错误: $_" -ForegroundColor Red
  exit 1
}
