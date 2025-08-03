// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package tools

import (
	"time"

	"github.com/caarlos0/log"
	"github.com/magefile/mage/mg"
	"github.com/opencommand/tinge"
)

// generates code (Go bindings, DSL grammar, etc).
func Generate() {
	log.Info(tinge.Styled().Bold("generating code...").String())
	start := time.Now()

	mg.Deps(generateHuloParser, generateUnsafeParser, generateStringer)

	elapsed := time.Since(start)
	log.Infof("generator completed in %.2fs", elapsed.Seconds())
}

// generateHuloParser generates the Go parser code from ANTLR grammar files.
func generateHuloParser() error {
	log.Info("generating hulo parser")
	return runCmdInDir("syntax/hulo/parser/grammar",
		"java", "-jar", "../antlr.jar",
		"-Dlanguage=Go", "-visitor", "-no-listener",
		"-o", "../generated", "-package", "generated", "*.g4")
}

// generateUnsafeParser generates the Go parser code from ANTLR grammar files.
func generateUnsafeParser() error {
	log.Info("gnerating unsafe parser")
	return runCmdInDir("internal/unsafe/grammar",
		"java", "-jar", "../antlr.jar",
		"-Dlanguage=Go", "-visitor", "-no-listener",
		"-o", "../generated", "-package", "generated", "*.g4")
}

func generateStringer() error {
	return runCmd("go", "generate", "./...")
}
