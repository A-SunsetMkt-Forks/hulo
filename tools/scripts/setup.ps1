# REM Copyright 2025 The Hulo Authors. All rights reserved.
# REM Use of this source code is governed by a MIT-style
# REM license that can be found in the LICENSE file.

$ErrorActionPreference = "Stop"

function Write-Info($msg) {
    Write-Output $msg
}

function Write-ErrorAndExit($msg) {
    Write-Host $msg -ForegroundColor Red
    exit 1
}

function Test-CommandExists($cmd) {
    $null -ne (Get-Command $cmd -ErrorAction SilentlyContinue)
}

Write-Info "Checking Scoop installation..."
if (-not (Test-CommandExists "scoop")) {
    Write-Info "Installing Scoop..."
    Set-ExecutionPolicy RemoteSigned -Scope CurrentUser -Force
    Invoke-Expression (Invoke-WebRequest -UseBasicParsing get.scoop.sh).Content
} else {
    Write-Info "Scoop already installed"
}

Write-Info "Checking Go installation..."
if (-not (Test-CommandExists "go")) {
    Write-Info "Installing Go via Scoop..."
    scoop install go
} else {
    Write-Info "Go already installed"
}

Write-Info "Setting Go proxy..."
go env -w GOPROXY=https://goproxy.cn,direct

Write-Info "Running go mod tidy..."
go mod tidy

Write-Info "Checking Java installation..."
if (-not (Test-CommandExists "java")) {
    Write-Info "Installing Java via Scoop..."
    scoop bucket add java
    scoop install temurin22-jdk
} else {
    Write-Info "Java already installed"
}

Write-Info "Checking GitVersion installation..."
if (-not (Test-CommandExists "gitversion")) {
    Write-Info "Installing GitVersion..."
    scoop install gitversion
} else {
    Write-Info "GitVersion already installed"
}

Write-Info "Checking Mage installation..."
if (-not (Test-CommandExists "mage")) {
    Write-Info "Installing Mage..."
    go install github.com/magefile/mage@latest
} else {
    Write-Info "Mage already installed"
}

Write-Info "Running mage setup..."
mage setup
