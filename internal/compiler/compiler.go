// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package compiler

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/hulo-lang/hulo/syntax/hulo/ast"

	// "github.com/hulo-lang/hulo/internal/build/bash"
	"github.com/hulo-lang/hulo/internal/config"
	"github.com/hulo-lang/hulo/internal/interpreter"
	"github.com/hulo-lang/hulo/internal/linker"
	"github.com/hulo-lang/hulo/internal/module"
	"github.com/hulo-lang/hulo/internal/object"
	"github.com/hulo-lang/hulo/internal/optimizer"
	"github.com/hulo-lang/hulo/internal/transpiler"
	"github.com/hulo-lang/hulo/internal/vfs"
	"github.com/hulo-lang/hulo/syntax/hulo/token"
)

type Compiler struct {
	optimizer   *optimizer.Optimizer
	interpreter *interpreter.Interpreter
	transpilers map[string]transpiler.Transpiler[any]
	moduleMgr   *module.DependecyResolver
	linker      *linker.Linker
	fs          vfs.VFS
}

func Compile(cfg *config.Huloc, fs vfs.VFS) error {
	// 直接调用moduleMgr来执行
	moduleMgr := module.NewDependecyResolver(cfg, fs)
	if err := module.ResolveAllDependencies(moduleMgr, cfg.Main); err != nil {
		return err
	}

	c := &Compiler{
		transpilers: make(map[string]transpiler.Transpiler[any]),
	}

	// 注册转译器映射关系

	ld := linker.NewLinker(fs)
	ld.Listen(".bat", linker.BeginEnd{Begin: "REM HULO_LINK_BEGIN", End: "REM HULO_LINK_END"}, linker.BeginEnd{Begin: ":: HULO_LINK_BEGIN", End: ":: HULO_LINK_END"})
	ld.Listen(".vbs", linker.BeginEnd{Begin: "' HULO_LINK_BEGIN", End: "' HULO_LINK_END"})
	ld.Listen(".cmd", linker.BeginEnd{Begin: "REM HULO_LINK_BEGIN", End: "REM HULO_LINK_END"})
	ld.Listen(".ps1", linker.BeginEnd{Begin: "# HULO_LINK_BEGIN", End: "# HULO_LINK_END"})
	ld.Listen(".sh", linker.BeginEnd{Begin: "# HULO_LINK_BEGIN", End: "# HULO_LINK_END"})

	for _, target := range cfg.Targets {
		env := interpreter.NewEnvironment()
		env.SetWithScope("TARGET", &object.StringValue{Value: target}, token.CONST, true)
		env.SetWithScope("OS", &object.StringValue{Value: runtime.GOOS}, token.CONST, true)
		env.SetWithScope("ARCH", &object.StringValue{Value: runtime.GOARCH}, token.CONST, true)
		interp := interpreter.NewInterpreter(env)

		err := moduleMgr.VisitModules(func(mod *module.Module) error {
			// Eval 也应该按照 Module 的
			// dap 调试器启动
			mod.AST = interp.Eval(mod.AST).(*ast.File)
			mod.AST = c.optimizer.Optimize(mod.AST).(*ast.File)
			return nil
		})
		if err != nil {
			return err
		}

		transpiler, ok := c.transpilers[target]
		if !ok {
			return fmt.Errorf("transpiler for target %s not found", target)
		}

		var results map[string]any

		err = moduleMgr.VisitModules(func(mod *module.Module) error {
			targetAST := transpiler.Convert(mod.AST)
			results[strings.Replace(mod.Path, ".hl", transpiler.GetTargetExt(), 1)] = targetAST
			return nil
		})
		if err != nil {
			return err
		}

		for _, nodes := range transpiler.UnresolvedSymbols() {
			for _, node := range nodes {
				ld.Read(node.AST.Path)
				linkable := ld.Load(node.AST.Path)
				symbol := linkable.Lookup(node.AST.Symbol)
				if symbol == nil {
					return fmt.Errorf("symbol %s not found in %s", node.AST.Symbol, node.AST.Path)
				}
				// node.Node.Val = symbol.Text()
			}
		}

		// 写文件
	}

	return nil
}

func (c *Compiler) BindTarget(transpiler transpiler.Transpiler[any], targets ...string) {
	for _, target := range targets {
		c.transpilers[target] = transpiler
	}
}
