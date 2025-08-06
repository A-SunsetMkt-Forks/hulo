// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package compiler

import (
	"fmt"

	"github.com/hulo-lang/hulo/syntax/hulo/ast"

	// "github.com/hulo-lang/hulo/internal/build/bash"
	"github.com/hulo-lang/hulo/internal/config"
	"github.com/hulo-lang/hulo/internal/interpreter"
	"github.com/hulo-lang/hulo/internal/linker"
	"github.com/hulo-lang/hulo/internal/optimizer"
	"github.com/hulo-lang/hulo/internal/transpiler"
	"github.com/hulo-lang/hulo/syntax/hulo/parser"
)

type Compiler struct {
	analyzer    *parser.Analyzer
	optimizer   *optimizer.Optimizer
	interpreter *interpreter.Interpreter
	transpilers map[string]transpiler.Transpiler[any]
	// moduleMgr
	linker *linker.Linker
}

func Compile(cfg *config.Huloc) error {
	file, err := parser.ParseSourceFile(cfg.Main)
	if err != nil {
		return err
	}
	ctx := interpreter.DefaultContext()
	for canEval(file) {
		interpreter.Evaluate(ctx, file)
	}

	c := &Compiler{
		transpilers: make(map[string]transpiler.Transpiler[any]),
	}

	// 注册映射关系

	for _, target := range cfg.Targets {
		transpiler, ok := c.transpilers[target]
		if !ok {
			return fmt.Errorf("transpiler for target %s not found", target)
		}

		files, unresolvedSymbols, err := transpiler.Transpile(cfg.Main)
		if err != nil {
			return err
		}

		fmt.Println(files)
		fmt.Println(unresolvedSymbols)
	}

	return nil
}

func (c *Compiler) BindTarget(transpiler transpiler.Transpiler[any], targets ...string) {
	for _, target := range targets {
		c.transpilers[target] = transpiler
	}
}

func canEval(file *ast.File) bool {
	return true
}
