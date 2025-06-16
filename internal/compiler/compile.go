// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package compiler

import (
	"github.com/antlr4-go/antlr/v4"
	"github.com/hulo-lang/hulo/syntax/hulo/ast"

	// "github.com/hulo-lang/hulo/internal/build/bash"
	build "github.com/hulo-lang/hulo/internal/build/bash"
	"github.com/hulo-lang/hulo/internal/comptime"
	"github.com/hulo-lang/hulo/internal/config"
	"github.com/hulo-lang/hulo/syntax/hulo/parser"
)

func Compile(cfg *config.Huloc) error {
	file, err := parser.ParseSourceFile(cfg.Main, parser.ParseOptions{
		Channel: antlr.TokenDefaultChannel,
	})
	if err != nil {
		return err
	}
	ctx := comptime.DefaultContext()
	for canEval(file) {
		comptime.Evaluate(ctx, file)
	}

	switch cfg.Language {
	case config.L_BASH:
		_, err := build.Translate(cfg.CompilerOptions.Bash, file)
		if err != nil {
			return err
		}
	}

	return nil
}

func canEval(file *ast.File) bool {
	return true
}
