// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package transpiler

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/hulo-lang/hulo/internal/config"
	"github.com/hulo-lang/hulo/internal/container"
	"github.com/hulo-lang/hulo/internal/interpreter"
	"github.com/hulo-lang/hulo/internal/util"
	"github.com/hulo-lang/hulo/internal/vfs"
	hast "github.com/hulo-lang/hulo/syntax/hulo/ast"
	"github.com/hulo-lang/hulo/syntax/hulo/parser"
	past "github.com/hulo-lang/hulo/syntax/powershell/ast"
)

type Module struct {
	Name         string
	Path         string
	AST          *hast.File
	Exports      map[string]*ExportInfo
	Imports      []*ImportInfo
	Dependencies []string
	Transpiled   *past.File
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

type SymbolType int

const (
	SymbolVariable SymbolType = iota
	SymbolFunction
	SymbolClass
	SymbolConstant
	SymbolArray
	SymbolObject
)

type Symbol struct {
	Name        string     // 原始名称
	MangledName string     // 混淆后的名称
	Type        SymbolType // 符号类型
	Value       hast.Node  // 原始AST节点
	Scope       *Scope     // 所属作用域
	Module      *Module    // 所属模块
	IsExported  bool       // 是否导出
}

type Scope struct {
	ID        string             // 作用域唯一标识
	Type      ScopeType          // 作用域类型
	Variables map[string]*Symbol // 变量名 -> 符号
	Parent    *Scope             // 父作用域
	Module    *Module            // 所属模块
}

type ScopeStack struct {
	scopes []*Scope
}

func (ss *ScopeStack) Push(scope *Scope) {
	ss.scopes = append(ss.scopes, scope)
}

func (ss *ScopeStack) Pop() {
	if len(ss.scopes) > 0 {
		ss.scopes = ss.scopes[:len(ss.scopes)-1]
	}
}

func (ss *ScopeStack) Current() *Scope {
	if len(ss.scopes) == 0 {
		return nil
	}
	return ss.scopes[len(ss.scopes)-1]
}

type SymbolTable struct {
	// 保留现有的模块级别符号
	functions map[string]hast.Node
	classes   map[string]hast.Node
	constants map[string]hast.Node
	variables map[string]hast.Node

	// 作用域管理
	scopes      *ScopeStack
	globalScope *Scope

	// 名称分配器
	moduleAllocator *util.Allocator            // 模块级别的分配器
	scopeAllocators map[string]*util.Allocator // 作用域级别的分配器
	enableMangle    bool
}

// NewSymbolTable 创建新的符号表
func NewSymbolTable(moduleName string, enableMangle bool) *SymbolTable {
	return &SymbolTable{
		functions: make(map[string]hast.Node),
		classes:   make(map[string]hast.Node),
		constants: make(map[string]hast.Node),
		variables: make(map[string]hast.Node),

		scopes: &ScopeStack{},
		globalScope: &Scope{
			ID:        "global",
			Type:      GlobalScope,
			Variables: make(map[string]*Symbol),
			Module:    nil, // 全局作用域不属于特定模块
		},

		moduleAllocator: util.NewAllocator(
			util.WithPrefix(fmt.Sprintf("_%s_", moduleName)),
			util.WithInitialCounter(0),
		),
		scopeAllocators: make(map[string]*util.Allocator),
		enableMangle:    enableMangle,
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

// 作用域管理方法
func (st *SymbolTable) PushScope(scopeType ScopeType, module *Module) *Scope {
	parent := st.CurrentScope()
	scopeID := fmt.Sprintf("scope_%d", len(st.scopes.scopes))

	newScope := &Scope{
		ID:        scopeID,
		Type:      scopeType,
		Variables: make(map[string]*Symbol),
		Parent:    parent,
		Module:    module,
	}

	// 为每个作用域创建独立的分配器
	st.scopeAllocators[scopeID] = util.NewAllocator(
		util.WithPrefix(fmt.Sprintf("_%s_", scopeID)),
		util.WithInitialCounter(0),
	)

	st.scopes.Push(newScope)
	return newScope
}

func (st *SymbolTable) PopScope() {
	current := st.CurrentScope()
	if current != nil {
		// 清理作用域分配器
		delete(st.scopeAllocators, current.ID)
	}
	st.scopes.Pop()
}

func (st *SymbolTable) CurrentScope() *Scope {
	return st.scopes.Current()
}

// 变量分配策略
func (st *SymbolTable) AllocVariableName(originalName string, scope *Scope) string {
	if !st.enableMangle {
		return originalName
	}
	// 检查当前作用域是否已有同名变量
	if _, exists := scope.Variables[originalName]; exists {
		// 如果存在，使用作用域级别的分配器
		allocator := st.scopeAllocators[scope.ID]
		return allocator.AllocName(originalName)
	}

	// 检查模块级别的符号
	if st.HasVariable(originalName) || st.HasFunction(originalName) ||
		st.HasClass(originalName) || st.HasConstant(originalName) {
		// 模块级别符号，使用模块分配器
		return st.moduleAllocator.AllocName(originalName)
	}

	// 新变量，使用当前作用域的分配器
	allocator := st.scopeAllocators[scope.ID]
	return allocator.AllocName(originalName)
}

// 符号查找
func (st *SymbolTable) LookupSymbol(name string) *Symbol {
	current := st.CurrentScope()
	for current != nil {
		if symbol, exists := current.Variables[name]; exists {
			return symbol
		}
		current = current.Parent
	}
	return nil
}

// 添加符号到当前作用域
func (st *SymbolTable) AddSymbol(symbol *Symbol) {
	current := st.CurrentScope()
	if current != nil {
		current.Variables[symbol.Name] = symbol
	}
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

	options *config.Huloc

	interp *interpreter.Interpreter
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
		// 解析相对路径
		resolvedDep := mm.resolveRelativePath(filePath, dep)
		if err := mm.resolveDependencies(resolvedDep); err != nil {
			return err
		}
	}

	// 添加到解析顺序
	mm.resolver.order = append(mm.resolver.order, filePath)

	return nil
}

// resolveRelativePath 解析相对于当前模块的导入路径
func (mm *ModuleManager) resolveRelativePath(currentModulePath, importPath string) string {
	if filepath.IsAbs(importPath) {
		return importPath
	}
	// 如果导入路径是绝对路径或以 ./ 或 ../ 开头，则相对于当前模块解析
	if strings.HasPrefix(importPath, "./") || strings.HasPrefix(importPath, "../") {
		currentDir := filepath.Dir(currentModulePath)
		return filepath.Join(currentDir, importPath)
	}

	// 否则，尝试在当前模块的目录中查找
	currentDir := filepath.Dir(currentModulePath)
	relativePath := filepath.Join(currentDir, importPath)

	// 检查相对路径是否存在
	if mm.vfs.Exists(relativePath) || mm.vfs.Exists(relativePath+".hl") {
		return relativePath
	}

	// 如果相对路径不存在，返回原始路径（让上层处理）
	return importPath
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

	opts := []parser.ParserOptions{}
	if len(mm.options.Parser.ShowASTTree) > 0 {
		switch mm.options.Parser.ShowASTTree {
		case "stdout":
			opts = append(opts, parser.OptionTracerASTTree(os.Stdout))
		case "stderr":
			opts = append(opts, parser.OptionTracerASTTree(os.Stderr))
		case "file":
			file, err := mm.vfs.Open(mm.options.Parser.ShowASTTree)
			if err != nil {
				return nil, fmt.Errorf("failed to open file %s: %w", mm.options.Parser.ShowASTTree, err)
			}
			defer file.Close()
			opts = append(opts, parser.OptionTracerASTTree(file))
		}
	}
	if !mm.options.Parser.EnableTracer {
		opts = append(opts, parser.OptionDisableTracer())
	}
	if mm.options.Parser.DisableTiming {
		opts = append(opts, parser.OptionTracerDisableTiming())
	}

	// 解析AST
	ast, err := parser.ParseSourceScript(string(content), opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to parse file %s: %w", filePath, err)
	}

	ast = mm.interp.Eval(ast).(*hast.File)

	// 提取导入信息
	imports, dependencies := mm.extractImports(ast, filePath)

	// 创建模块
	moduleName := mm.getModuleName(filePath)
	module := &Module{
		Name:         moduleName,
		Path:         filePath,
		AST:          ast,
		Exports:      make(map[string]*ExportInfo),
		Imports:      imports,
		Dependencies: dependencies,
		Symbols:      NewSymbolTable(moduleName, mm.options.EnableMangle),
	}

	// 缓存模块
	mm.modules[filePath] = module

	return module, nil
}

func (mm *ModuleManager) extractImports(ast *hast.File, currentModulePath string) ([]*ImportInfo, []string) {
	var imports []*ImportInfo
	var dependencies []string
	// 使用map来去重依赖
	dependencySet := make(map[string]bool)

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
				dep := importInfo.ModulePath

				// 非 vbs 文件才补全 .hl
				depPath := dep
				if !filepath.IsAbs(dep) && !strings.HasSuffix(dep, ".hl") && !strings.HasPrefix(dep, "./") && !strings.HasPrefix(dep, "../") {
					currentDir := filepath.Dir(currentModulePath)
					depPath = filepath.Join(currentDir, dep+".hl")
				}

				// 去重：只有未添加过的依赖才添加
				if !dependencySet[depPath] {
					dependencies = append(dependencies, depPath)
					dependencySet[depPath] = true
				}
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
