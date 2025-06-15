//go:build mage

package main

// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
//
// These tasks are executed via Mage (https://magefile.org) and are used to automate
// formatting, linting, testing, cleaning, and building the project.
//
// Available targets:
//
//   build   - Compiles and packages the project using GoReleaser.
//   clean   - Removes generated files and build artifacts.
//   format  - Automatically formats the source code using goimports.
//   lint    - Runs golangci-lint on the codebase.
//   setup   - Prepares the development environment and downloads dependencies
//             (e.g., Go tools, ANTLR jars, and Go modules).
//   test    - Runs all unit tests with coverage.
//
// Run `mage -l` to see all available targets.
// Run `mage <target>` to execute a specific task.

import (
	//mage:import
	"github.com/hulo-lang/hulo/tools"
	_ "github.com/hulo-lang/hulo/tools"
)

// TODO
// mg.Namespace, argument

var Aliases = map[string]any{
	"s":   tools.Setup,
	"env": tools.Setup,
	"b":   tools.Build,
	"c":   tools.Clean,
	"fmt": tools.Format,
	"l":   tools.Lint,
	"t":   tools.Test,
	"g":   tools.Generate,
	"gen": tools.Generate,
}

// A var named Default indicates which target is the default.  If there is no
// default, running mage will list the targets available.
// var Default = Build
