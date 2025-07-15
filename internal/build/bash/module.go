// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package build

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/hulo-lang/hulo/internal/container"
	"github.com/hulo-lang/hulo/internal/vfs"
	bast "github.com/hulo-lang/hulo/syntax/bash/ast"
	hast "github.com/hulo-lang/hulo/syntax/hulo/ast"
	"github.com/hulo-lang/hulo/syntax/hulo/parser"
)

type Module struct {
	Name         string
	Path         string
	AST          *hast.File
	Exports      map[string]*ExportInfo
	Imports      []*ImportInfo
	Dependencies []string
	Transpiled   *bast.File
	Symbols      *SymbolTable
}

type ExportInfo struct {
	Name   string
	Value  hast.Node
	Kind   ExportKind
	Public bool
}

// Function, Class, Constant, etc.
type ExportKind int

const (
	ExportFunction ExportKind = iota
	ExportClass
	ExportConstant
	ExportVariable
)

type ImportInfo struct {
	ModulePath string
	SymbolName []string
	Alias      string
	Kind       ImportKind
}

// Named, All, Default
type ImportKind int

const (
	ImportSingle ImportKind = iota
	ImportMulti
	ImportAll
)

type SymbolTable struct {
	functions map[string]hast.Node
	classes   map[string]hast.Node
	constants map[string]hast.Node
	variables map[string]hast.Node
}

// NewSymbolTable 创建新的符号表
func NewSymbolTable() *SymbolTable {
	return &SymbolTable{
		functions: make(map[string]hast.Node),
		classes:   make(map[string]hast.Node),
		constants: make(map[string]hast.Node),
		variables: make(map[string]hast.Node),
	}
}

// HasVariable 检查是否包含变量
func (st *SymbolTable) HasVariable(name string) bool {
	_, exists := st.variables[name]
	return exists
}

// HasFunction 检查是否包含函数
func (st *SymbolTable) HasFunction(name string) bool {
	_, exists := st.functions[name]
	return exists
}

// HasClass 检查是否包含类
func (st *SymbolTable) HasClass(name string) bool {
	_, exists := st.classes[name]
	return exists
}

// HasConstant 检查是否包含常量
func (st *SymbolTable) HasConstant(name string) bool {
	_, exists := st.constants[name]
	return exists
}

// AddVariable 添加变量
func (st *SymbolTable) AddVariable(name string, node hast.Node) {
	st.variables[name] = node
}

// AddFunction 添加函数
func (st *SymbolTable) AddFunction(name string, node hast.Node) {
	st.functions[name] = node
}

// AddClass 添加类
func (st *SymbolTable) AddClass(name string, node hast.Node) {
	st.classes[name] = node
}

// AddConstant 添加常量
func (st *SymbolTable) AddConstant(name string, node hast.Node) {
	st.constants[name] = node
}

type ScopeType int

const (
	GlobalScope ScopeType = iota
	FunctionScope
	ClassScope
	BlockScope
)

type ModuleManager struct {
	vfs      vfs.VFS
	basePath string
	modules  map[string]*Module
	resolver *DependencyResolver

	isMain bool // 首次要收集 builtin 包
}

type DependencyResolver struct {
	visited container.Set[string]
	stack   container.Set[string]
	order   []string
}

func (mm *ModuleManager) ResolveAllDependencies(mainFile string) ([]*Module, error) {
	// 重置解析器状态
	mm.resolver.visited.Clear()
	mm.resolver.stack.Clear()
	mm.resolver.order = make([]string, 0)

	// 从主文件开始解析
	if err := mm.resolveDependencies(mainFile); err != nil {
		return nil, err
	}

	// 按依赖顺序返回模块
	var modules []*Module
	for _, path := range mm.resolver.order {
		if module, exists := mm.modules[path]; exists {
			modules = append(modules, module)
		}
	}

	return modules, nil
}

func (mm *ModuleManager) resolvePath(filePath string) string {
	exts := []string{".hl", ".bash", ".vbs"}
	for _, ext := range exts {
		if strings.HasSuffix(filePath, ext) {
			return filePath
		}
	}
	if filepath.IsAbs(filePath) {
		return filePath
	}
	return filePath + ".hl"
}

func (mm *ModuleManager) resolveDependencies(filePath string) error {
	filePath = mm.resolvePath(filePath)

	// 检查是否已经访问过
	if mm.resolver.visited.Contains(filePath) {
		return nil
	}

	// 检查循环依赖
	if mm.resolver.stack.Contains(filePath) {
		return fmt.Errorf("circular dependency detected: %s", filePath)
	}

	mm.resolver.visited.Add(filePath)
	mm.resolver.stack.Add(filePath)
	defer func() { mm.resolver.stack.Remove(filePath) }()

	// 加载模块
	module, err := mm.loadModule(filePath)
	if err != nil {
		return err
	}

	// 递归解析依赖
	for _, dep := range module.Dependencies {
		if err := mm.resolveDependencies(dep); err != nil {
			return err
		}
	}

	// 添加到解析顺序
	mm.resolver.order = append(mm.resolver.order, filePath)

	return nil
}

func (mm *ModuleManager) loadModule(filePath string) (*Module, error) {
	// 检查缓存
	if module, exists := mm.modules[filePath]; exists {
		return module, nil
	}

	// 读取文件内容
	content, err := mm.vfs.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	// 解析AST
	ast, err := parser.ParseSourceScript(string(content))
	if err != nil {
		return nil, fmt.Errorf("failed to parse file %s: %w", filePath, err)
	}

	// 提取导入信息
	imports, dependencies := mm.extractImports(ast)

	// 创建模块
	module := &Module{
		Name:         mm.getModuleName(filePath),
		Path:         filePath,
		AST:          ast,
		Exports:      make(map[string]*ExportInfo),
		Imports:      imports,
		Dependencies: dependencies,
		Symbols:      NewSymbolTable(),
	}

	// 缓存模块
	mm.modules[filePath] = module

	return module, nil
}

func (mm *ModuleManager) extractImports(ast *hast.File) ([]*ImportInfo, []string) {
	var imports []*ImportInfo
	var dependencies []string

	for _, stmt := range ast.Stmts {
		if importStmt, ok := stmt.(*hast.Import); ok {
			var importInfo *ImportInfo

			switch {
			case importStmt.ImportSingle != nil:
				// import "module" as alias
				importInfo = &ImportInfo{
					ModulePath: importStmt.ImportSingle.Path,
					Alias:      importStmt.ImportSingle.Alias,
					Kind:       ImportSingle,
				}

			case importStmt.ImportMulti != nil:
				// import { symbol1, symbol2 } from "module"
				var symbols []string
				for _, field := range importStmt.ImportMulti.List {
					symbols = append(symbols, field.Field)
				}
				importInfo = &ImportInfo{
					ModulePath: importStmt.ImportMulti.Path,
					SymbolName: symbols,
					Kind:       ImportMulti,
				}

			case importStmt.ImportAll != nil:
				// import * from "module"
				importInfo = &ImportInfo{
					ModulePath: importStmt.ImportAll.Path,
					Alias:      importStmt.ImportAll.Alias,
					Kind:       ImportAll,
				}
			}

			if importInfo != nil {
				imports = append(imports, importInfo)
				dependencies = append(dependencies, importInfo.ModulePath)
			}
		}
	}

	return imports, dependencies
}

func (mm *ModuleManager) getModuleName(filePath string) string {
	baseName := filepath.Base(filePath)
	ext := filepath.Ext(baseName)
	if ext != "" {
		baseName = baseName[:len(baseName)-len(ext)]
	}
	return baseName
}

func (mm *ModuleManager) GetModule(path string) *Module {
	return mm.modules[path]
}
