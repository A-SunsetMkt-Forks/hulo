// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package transpiler

import (
	"fmt"

	"github.com/hulo-lang/hulo/internal/container"
	bast "github.com/hulo-lang/hulo/syntax/batch/ast"
	hast "github.com/hulo-lang/hulo/syntax/hulo/ast"
)

type Module struct {
	Name         string // 模块名，要处理别名的情况 import "time" as t
	Path         string // 模块绝对路径
	AST          *hast.File
	Exports      map[string]*ExportInfo
	Imports      []*ImportInfo
	Dependencies []string
	Transpiled   *bast.File
	Symbols      *SymbolTable
	State        ModuleState
}

type ModuleState int

const (
	ModuleStateUnresolved ModuleState = iota
	ModuleStateResolved
	ModuleStateMangled
)

func (m *Module) IsMangled() bool {
	return m.State == ModuleStateMangled
}

type ImportInfo struct {
	ModulePath string
	SymbolName []string
	Alias      string
	Kind       ImportKind
}

type ImportKind int

type SymbolTable struct{}

type ExportInfo struct{}

// 从 main 入口开始解析
func ResolveAllDependencies(mainFile string) error {
	resolver := NewDependecyResolver()

	// 递归解析所有依赖，会收集依赖建立 Imports
	if err := resolver.resolveRecursive(mainFile); err != nil {
		return err
	}

	if cycles, err := resolver.detectCycles(); err != nil {
		return fmt.Errorf("detect cycles: %v", cycles)
	}

	// 拓扑确认引入顺序
	if err := resolver.topologicalSort(); err != nil {
		return fmt.Errorf("topological sort: %w", err)
	}

	// 混淆叶节点
	for _, module := range resolver.modules {
		if len(module.Dependencies) == 0 {
			if err := MangleModule(module); err != nil {
				return fmt.Errorf("mangle module: %w", err)
			}
		}
	}

	// 按拓扑顺序混淆其他模块
	for _, module := range resolver.modules {
		if len(module.Dependencies) > 0 {
			// 确保所有依赖都已混淆
			for _, dep := range module.Dependencies {
				if !resolver.modules[dep].IsMangled() {
					return fmt.Errorf("dependency %s not mangled", dep)
				}
			}

			// 混淆当前模块
			if err := MangleModule(module); err != nil {
				return err
			}
		}
	}

	// 混淆完可以建立符号表了

	// math.add() -> _math_add()
	// m.add() -> _math_add()
	// add() -> _math_add()
	// a() -> _math_add()
	// 所有导出 -> 直接映射

	// 对于不同的输出方式，需要建立引用

	// 输出
	return nil
}

func MangleModule(module *Module) error {
	return nil
}

type DependecyResolver struct {
	modules map[string]*Module
	visited container.Set[string]
	order   []string
}

func NewDependecyResolver() *DependecyResolver {
	return &DependecyResolver{}
}

func (r *DependecyResolver) resolveRecursive(file string) error {
	return nil
}

// 检查循环依赖
func (r *DependecyResolver) detectCycles() (*Cycles, error) {
	return nil, nil
}

func (r *DependecyResolver) topologicalSort() error {
	return nil
}

type Cycles struct{}
