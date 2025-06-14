#!/bin/bash
# Copyright 2025 The Hulo Authors. All rights reserved.
# Use of this source code is governed by a MIT-style
# license that can be found in the LICENSE file.

set -e

INFO='\033[1;38;5;12m'
ERROR='\033[0;31m'
NC='\033[0m' # none color

info() {
  echo -e "  ${INFO}•${NC} $1"
}

error_exit() {
  echo -e "${ERROR}Error: $1${NC}"
  exit 1
}

command_exists() {
  command -v "$1" >/dev/null 2>&1
}

info "Checking Scoop installation..."
if ! command_exists scoop; then
  info "Scoop not found. Please install Scoop manually on Linux or adapt script for Linux package manager."
else
  info "Scoop already installed"
fi

info "Checking Go installation..."
if ! command_exists go; then
  info "Installing Go..."
  # 这里示例用apt-get安装，实际根据你环境换成对应命令
  if command_exists apt-get; then
    sudo apt-get update
    sudo apt-get install -y golang
  else
    error_exit "No supported package manager found to install Go. Please install Go manually."
  fi
else
  info "Go already installed"
fi

info "Setting Go proxy..."
go env -w GOPROXY=https://goproxy.cn,direct

info "Running go mod tidy..."
go mod tidy

info "Checking Java installation..."
if ! command_exists java; then
  info "Installing Java..."
  if command_exists apt-get; then
    sudo apt-get install -y default-jdk
  else
    error_exit "No supported package manager found to install Java. Please install Java manually."
  fi
else
  info "Java already installed"
fi

info "Checking Mage installation..."
if ! command_exists mage; then
  info "Installing Mage..."
  GO111MODULE=on go install github.com/magefile/mage@latest
else
  info "Mage already installed"
fi

info "Running mage setup..."
mage setup
