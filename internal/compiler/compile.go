// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package compiler

import (
	"github.com/hulo-lang/hulo/syntax/hulo/ast"

	// "github.com/hulo-lang/hulo/internal/build/bash"
	"github.com/hulo-lang/hulo/internal/config"
	"github.com/hulo-lang/hulo/internal/interpreter"
	"github.com/hulo-lang/hulo/internal/optimizer"
	build "github.com/hulo-lang/hulo/internal/transpiler/bash"
	"github.com/hulo-lang/hulo/syntax/hulo/parser"
)

type Compiler struct {
	analyzer    *parser.Analyzer
	optimizer   *optimizer.Optimizer
	interpreter *interpreter.Interpreter
	// transpiler
}

// TODO transpiler 输出未知符号信息，让 ld 链接上

func Compile(cfg *config.Huloc) error {
	file, err := parser.ParseSourceFile(cfg.Main)
	if err != nil {
		return err
	}
	ctx := interpreter.DefaultContext()
	for canEval(file) {
		interpreter.Evaluate(ctx, file)
	}

	for _, target := range cfg.Targets {
		switch target {
		case config.L_BASH:
			_, err := build.Transpile(cfg, nil, ".")
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func canEval(file *ast.File) bool {
	return true
}
