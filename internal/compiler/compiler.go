// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package compiler

import (
	"fmt"
	"runtime"

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
	moduleMgr := module.NewDependecyResolver()
	if err := module.ResolveAllDependencies(moduleMgr, cfg.Main); err != nil {
		return err
	}

	c := &Compiler{
		transpilers: make(map[string]transpiler.Transpiler[any]),
	}

	// 注册映射关系

	for _, target := range cfg.Targets {
		env := interpreter.NewEnvironment()
		env.SetWithScope("TARGET", &object.StringValue{Value: target}, token.CONST, true)
		env.SetWithScope("OS", &object.StringValue{Value: runtime.GOOS}, token.CONST, true)
		env.SetWithScope("ARCH", &object.StringValue{Value: runtime.GOARCH}, token.CONST, true)
		interp := interpreter.NewInterpreter(env)

		err := moduleMgr.VisitModules(func(mod *module.Module) error {
			// Eval 也应该按照 Module 的
			mod.AST = interp.Eval(mod.AST).(*ast.File)
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
			targetNode, err := transpiler.Convert(mod.AST)
			if err != nil {
				return err
			}

			results[mod.Path] = targetNode
			return nil
		})
		if err != nil {
			return err
		}

		// 未知符号怎么拿？
		// 转译器给符号表吧 因为Unsafe会链接
		// unresolvedSymbols := transpiler.UnresolvedSymbols()
		// 未知符号是一个数据结构？然后有待链接的文件名以及符号信息？？？按照Module遍历吧？每个Module独立符号
		// for _, symbol := range unresolvedSymbols {
		// 	fmt.Println(symbol)
		// }


	}

	return nil
}

func (c *Compiler) BindTarget(transpiler transpiler.Transpiler[any], targets ...string) {
	for _, target := range targets {
		c.transpilers[target] = transpiler
	}
}
