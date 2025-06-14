@echo off
@REM Copyright 2025 The Hulo Authors. All rights reserved.
@REM Use of this source code is governed by a MIT-style
@REM license that can be found in the LICENSE file.

echo Installing scoop
Set-ExecutionPolicy RemoteSigned -Scope CurrentUser -Force
irm get.scoop.sh | iex

@REM resolve golang

echo Installing golang
scoop install go
go env -w GOPROXY=https://goproxy.cn,direct
go mod tidy

@REM resolve java
scoop bucket add java
scoop install temurin22-jdk

echo Installing mage...
go install github.com/magefile/mage@latest
mage setup
