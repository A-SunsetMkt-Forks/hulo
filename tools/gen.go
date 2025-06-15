// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package tools

import (
	"time"

	"github.com/caarlos0/log"
	"github.com/magefile/mage/mg"
)

func Generate() {
	log.Info("generating code...")
	start := time.Now()

	mg.Deps(generateParser, generateStringer, generateUnsafe)

	elapsed := time.Since(start)
	log.Infof("generator completed in %.2fs", elapsed.Seconds())
}

// generateParser generates the Go parser code from ANTLR grammar files.
func generateParser() error {
	log.Info("Generating parser")
	return runCmdInDir("syntax/hulo/parser/grammar",
		"java", "-jar", "../antlr.jar",
		"-Dlanguage=Go", "-visitor", "-no-listener",
		"-o", "../generated", "-package", "generated", "*.g4")
}

func generateUnsafe() error {
	return nil
}

func generateStringer() error {
	return nil
}
