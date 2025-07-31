// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package module

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"maps"

	"github.com/caarlos0/log"
	"github.com/hulo-lang/hulo/internal/config"
	"github.com/hulo-lang/hulo/internal/container"
	"github.com/hulo-lang/hulo/internal/vfs"
	hast "github.com/hulo-lang/hulo/syntax/hulo/ast"
	"github.com/hulo-lang/hulo/syntax/hulo/parser"
	"github.com/hulo-lang/hulo/syntax/hulo/token"
	"gopkg.in/yaml.v3"
)

type Module struct {
	ModuleID     int
	Name         string // 模块名，要处理别名的情况 import "time" as t
	Path         string // 模块绝对路径
	AST          *hast.File
	Exports      map[string]*ExportInfo
	Imports      []*ImportInfo
	Dependencies []string
	Symbols      *SymbolTable
	state        ModuleState
	Pkg          *config.HuloPkg
}

// ResolveImportedSymbols 解析导入的符号
func (m *Module) ResolveImportedSymbols(key string) (Symbol, error) {
	// 1. 先查当前模块的符号表
	if symbol := m.Symbols.LookupSymbol(key); symbol != nil {
		return symbol, nil
	}

	// 2. 查当前模块的导入
	for _, importInfo := range m.Imports {
		switch importInfo.Kind {
		case ImportSingle:
			if importInfo.Alias == key {
				// 通过别名查找
				return m.resolveSymbolFromModule(importInfo.ModulePath, importInfo.Alias)
			}
		case ImportMulti:
			for _, symbolName := range importInfo.SymbolName {
				if symbolName == key {
					return m.resolveSymbolFromModule(importInfo.ModulePath, symbolName)
				}
			}
		case ImportAll:
			// 通配符导入，需要查找所有导出的符号
			return m.resolveSymbolFromModule(importInfo.ModulePath, key)
		}
	}

	return nil, fmt.Errorf("symbol %s not found in module %s", key, m.Name)
}

// resolveSymbolFromModule 从指定模块解析符号
func (m *Module) resolveSymbolFromModule(modulePath, symbolName string) (Symbol, error) {
	// 这里需要访问模块解析器来获取模块信息
	// 暂时返回nil，实际实现需要访问模块解析器的modules map
	return nil, fmt.Errorf("module resolution not implemented for %s::%s", modulePath, symbolName)
}

// GetExportedSymbols 获取模块导出的符号
func (m *Module) GetExportedSymbols() []Symbol {
	var exported []Symbol

	for _, symbol := range m.Symbols.GlobalSymbols {
		if symbol.IsExported() {
			exported = append(exported, symbol)
		}
	}

	return exported
}

// GetSymbolByName 根据名称获取符号
func (m *Module) GetSymbolByName(name string) Symbol {
	return m.Symbols.LookupSymbol(name)
}

// GetFunctionSymbols 获取所有函数符号
func (m *Module) GetFunctionSymbols() []*FunctionSymbol {
	var functions []*FunctionSymbol

	for _, symbol := range m.Symbols.GlobalSymbols {
		if funcSymbol, ok := AsFunctionSymbol(symbol); ok {
			functions = append(functions, funcSymbol)
		}
	}

	return functions
}

func (m *Module) LookupFunctionSymbol(name string) *FunctionSymbol {
	for _, symbol := range m.Symbols.GlobalSymbols {
		if funcSymbol, ok := AsFunctionSymbol(symbol); ok && funcSymbol.GetName() == name {
			return funcSymbol
		}
	}
	return nil
}

func (m *Module) LookupNamespaceSymbol(name string) *NamespaceSymbol {
	for _, symbol := range m.Symbols.GlobalSymbols {
		if namespaceSymbol, ok := AsNamespaceSymbol(symbol); ok && namespaceSymbol.GetName() == name {
			return namespaceSymbol
		}
	}
	return nil
}

// GetClassSymbols 获取所有类符号
func (m *Module) GetClassSymbols() []*ClassSymbol {
	var classes []*ClassSymbol

	for _, symbol := range m.Symbols.GlobalSymbols {
		if classSymbol, ok := AsClassSymbol(symbol); ok {
			classes = append(classes, classSymbol)
		}
	}

	return classes
}

func (m *Module) LookupClassSymbol(name string) *ClassSymbol {
	for _, symbol := range m.Symbols.GlobalSymbols {
		if classSymbol, ok := AsClassSymbol(symbol); ok && classSymbol.GetName() == name {
			return classSymbol
		}
	}
	return nil
}

// GetVariableSymbols 获取所有变量符号
func (m *Module) GetVariableSymbols() []*VariableSymbol {
	var variables []*VariableSymbol

	for _, symbol := range m.Symbols.GlobalSymbols {
		if varSymbol, ok := AsVariableSymbol(symbol); ok {
			variables = append(variables, varSymbol)
		}
	}

	return variables
}

func (m *Module) LookupVariableSymbol(name string) *VariableSymbol {
	for _, symbol := range m.Symbols.GlobalSymbols {
		if varSymbol, ok := AsVariableSymbol(symbol); ok && varSymbol.GetName() == name {
			return varSymbol
		}
	}
	return nil
}

// GetSymbolStats 获取符号统计信息
func (m *Module) GetSymbolStats() map[string]int {
	stats := make(map[string]int)

	for _, symbol := range m.Symbols.GlobalSymbols {
		switch symbol.GetKind() {
		case SymbolVar:
			stats["variables"]++
		case SymbolFunc:
			stats["functions"]++
		case SymbolConst:
			stats["constants"]++
		case SymbolClass:
			stats["classes"]++
		case SymbolNamespace:
			stats["namespaces"]++
		}
	}

	return stats
}

// PrintSymbolTable 打印符号表（调试用）
func (m *Module) PrintSymbolTable() {
	fmt.Printf("=== Symbol Table for Module: %s ===\n", m.Name)

	// 打印全局符号
	fmt.Printf("Global Symbols:\n")
	for name, symbol := range m.Symbols.GlobalSymbols {
		fmt.Printf("  %s (%s)", name, symbolKindToString(symbol.GetKind()))
		if symbol.IsExported() {
			fmt.Printf(" [exported]")
		}
		if symbol.IsMangled() {
			fmt.Printf(" [mangled: %s]", symbol.GetMangledName())
		}
		fmt.Printf("\n")

		// 打印详细信息
		switch s := symbol.(type) {
		case *FunctionSymbol:
			fmt.Printf("    Return: %s, Params: %d\n", s.ReturnType, len(s.Parameters))
		case *ClassSymbol:
			fmt.Printf("    Fields: %d, Methods: %d\n", len(s.Fields), len(s.Methods))
		case *VariableSymbol:
			fmt.Printf("    Type: %s\n", s.DataType)
		}
	}

	// 打印所有作用域的符号
	fmt.Printf("\nScope Symbols:\n")
	for i, scope := range m.Symbols.ScopeStack {
		fmt.Printf("  Scope %d (%s, Level: %d):\n", i, scope.Name, scope.Level)
		for name, symbol := range scope.Symbols {
			fmt.Printf("    %s (%s)", name, symbolKindToString(symbol.GetKind()))
			if symbol.IsMangled() {
				fmt.Printf(" [mangled: %s]", symbol.GetMangledName())
			}
			fmt.Printf("\n")
		}
	}
}

// symbolKindToString 将符号类型转换为字符串
func symbolKindToString(kind SymbolKind) string {
	switch kind {
	case SymbolVar:
		return "variable"
	case SymbolFunc:
		return "function"
	case SymbolConst:
		return "constant"
	case SymbolClass:
		return "class"
	case SymbolNamespace:
		return "namespace"
	default:
		return "unknown"
	}
}

type ModuleState int

const (
	ModuleStateUnresolved ModuleState = iota
	ModuleStateResolved
	ModuleStateMangled
)

func (m *Module) IsMangled() bool {
	return m.state == ModuleStateMangled
}

func (m *Module) SetState(state ModuleState) {
	m.state = state
}

func (m *Module) State() ModuleState {
	return m.state
}

type ImportInfo struct {
	ModulePath string
	SymbolName []string
	Alias      string
	Kind       ImportKind
}

type ImportKind int

const (
	ImportSingle ImportKind = iota
	ImportMulti
	ImportAll
)

type ExportInfo struct {
	Symbol     Symbol
	ModulePath string
	IsDefault  bool
}

// ExportSymbol 导出符号信息
func (m *Module) ExportSymbol(symbolName string) *ExportInfo {
	symbol := m.Symbols.LookupSymbol(symbolName)
	if symbol == nil {
		return nil
	}

	return &ExportInfo{
		Symbol:     symbol,
		ModulePath: m.Path,
		IsDefault:  false,
	}
}

// ExportAllSymbols 导出所有公共符号
func (m *Module) ExportAllSymbols() {
	for name, symbol := range m.Symbols.GlobalSymbols {
		if symbol.IsExported() {
			m.Exports[name] = &ExportInfo{
				Symbol:     symbol,
				ModulePath: m.Path,
				IsDefault:  false,
			}
		}
	}
}

// 从 main 入口开始解析
func ResolveAllDependencies(resolver *DependecyResolver, mainFile string) error {

	// 递归解析所有依赖，会收集依赖建立 Imports, 默认情况下 mainFile 是 ./main.hl
	if err := resolver.resolveRecursive(nil, ".", mainFile); err != nil {
		return err
	}

	// 按依赖序混淆
	for i, path := range resolver.order {
		module := resolver.modules[path]
		module.ModuleID = i

		// 构建符号表
		if err := module.BuildSymbolTable(); err != nil {
			return fmt.Errorf("build symbol table for %s: %w", path, err)
		}

		// 导出所有公共符号
		module.ExportAllSymbols()

		// if len(module.Dependencies) > 0 {
		// 	for _, dep := range module.Dependencies {
		// 		if !resolver.modules[dep].IsMangled() {
		// 			return fmt.Errorf("dependency %s not mangled", dep)
		// 		}
		// 	}
		// }

		if err := MangleModule(module); err != nil {
			return fmt.Errorf("mangle module: %w", err)
		}
	}

	// 混淆完可以建立符号表了

	// math.add() -> _math_add()
	// m.add() -> _math_add()
	// add() -> _math_add()
	// a() -> _math_add()
	// 所有导出 -> 直接映射

	// 对于不同的输出方式，需要建立引用

	return nil
}

func (m *Module) BuildSymbolTable() error {
	st := &SymbolTable{
		GlobalSymbols: make(map[string]Symbol),
		ScopeStack:    []*Scope{},
		ModuleID:      m.ModuleID,
		Mangler: &SymbolMangler{
			typeCounter:   make(map[int]int),
			symbolCounter: make(map[string]int),
			mangledMap:    make(map[string]string),
		},
	}

	st.EnterScope("global")

	// 提取AST中的符号信息
	if err := m.extractSymbolsFromAST(st); err != nil {
		return fmt.Errorf("extract symbols from AST: %w", err)
	}

	// 混淆所有符号
	if err := st.MangleAllSymbols(); err != nil {
		return fmt.Errorf("mangle symbols: %w", err)
	}

	m.Symbols = st
	return nil
}

// extractSymbolsFromAST 从AST中提取符号信息
func (m *Module) extractSymbolsFromAST(st *SymbolTable) error {
	if m.AST == nil {
		return nil
	}

	// 处理所有语句
	for _, stmt := range m.AST.Stmts {
		if err := m.extractSymbolFromStmt(stmt, st); err != nil {
			return err
		}
	}

	return nil
}

// extractSymbolFromStmt 从语句中提取符号
func (m *Module) extractSymbolFromStmt(stmt hast.Stmt, st *SymbolTable) error {
	switch s := stmt.(type) {
	case *hast.FuncDecl:
		return m.extractFunctionSymbol(s, st)
	case *hast.ClassDecl:
		return m.extractClassSymbol(s, st)
	case *hast.AssignStmt:
		return m.extractVariableSymbol(s, st)
	case *hast.EnumDecl:
		return m.extractEnumSymbol(s, st)
	case *hast.TraitDecl:
		return m.extractTraitSymbol(s, st)
	case *hast.ModDecl:
		return m.extractModuleSymbol(s, st)
	case *hast.TypeDecl:
		return m.extractTypeAliasSymbol(s, st)
	case *hast.BlockStmt:
		return m.extractSymbolsFromBlock(s, st)
	case *hast.ForStmt:
		return m.extractForStmtSymbols(s, st)
	}
	return nil
}

// extractForStmtSymbols 提取 for 语句中的符号
func (m *Module) extractForStmtSymbols(forStmt *hast.ForStmt, st *SymbolTable) error {
	// 为 for 语句创建新的作用域
	st.EnterScope(fmt.Sprintf("for:%d", st.CurrentScope.Level+1))
	defer st.ExitScope()

	// 处理 Init 语句中的变量声明
	if forStmt.Init != nil {
		if assignStmt, ok := forStmt.Init.(*hast.AssignStmt); ok {
			// 将 for 循环中的变量声明当作局部变量处理
			if assignStmt.Tok == token.COLON_ASSIGN {
				assignStmt.Scope = token.LET // 强制设置为局部变量
			}
			if err := m.extractVariableSymbol(assignStmt, st); err != nil {
				return err
			}
		}
	}

	// 处理循环体中的符号
	if forStmt.Body != nil {
		if err := m.extractSymbolsFromBlock(forStmt.Body, st); err != nil {
			return err
		}
	}

	return nil
}

// extractFunctionSymbol 提取函数符号
func (m *Module) extractFunctionSymbol(fn *hast.FuncDecl, st *SymbolTable) error {
	funcSymbol := NewFunctionSymbol(fn.Name.Name, "", m.ModuleID)

	// 设置修饰符
	for _, modifier := range fn.Modifiers {
		switch modifier.Kind() {
		case hast.ModKindPub:
			funcSymbol.Exported = true
		case hast.ModKindConst:
			// 函数不能是const
		case hast.ModKindStatic:
			// 静态函数
		}
	}

	// 设置返回类型
	if fn.Type != nil {
		funcSymbol.ReturnType = m.typeToString(fn.Type)
	}

	// 设置参数
	for _, param := range fn.Recv {
		if paramExpr, ok := param.(*hast.Parameter); ok {
			param := &Parameter{
				Name: paramExpr.Name.Name,
				Type: m.typeToString(paramExpr.Type),
			}
			funcSymbol.Parameters = append(funcSymbol.Parameters, param)
		}
	}

	// 设置函数体
	funcSymbol.Body = fn.Body

	// 添加到符号表
	st.GlobalSymbols[fn.Name.Name] = funcSymbol
	return nil
}

// extractClassSymbol 提取类符号
func (m *Module) extractClassSymbol(cls *hast.ClassDecl, st *SymbolTable) error {
	classSymbol := NewClassSymbol(cls.Name.Name, m.ModuleID)

	// 设置修饰符
	for _, modifier := range cls.Modifiers {
		switch modifier.Kind() {
		case hast.ModKindPub:
			classSymbol.Exported = true
		}
	}

	// 设置父类
	if cls.Parent != nil {
		classSymbol.SuperClass = cls.Parent.Name
	}

	// 提取字段
	if cls.Fields != nil {
		for _, field := range cls.Fields.List {
			classField := &Field{
				Name:     field.Name.Name,
				Type:     m.typeToString(field.Type),
				IsPublic: false,
				Default:  nil,
			}

			// 检查字段修饰符
			for _, modifier := range field.Modifiers {
				switch modifier.Kind() {
				case hast.ModKindPub:
					classField.IsPublic = true
				}
			}

			// 设置默认值
			if field.Value != nil {
				classField.Default = field.Value
			}

			classSymbol.Fields[field.Name.Name] = classField
		}
	}

	// 提取方法
	for _, method := range cls.Methods {
		if err := m.extractMethodSymbol(method, classSymbol); err != nil {
			return err
		}
	}

	// 添加到符号表
	st.GlobalSymbols[cls.Name.Name] = classSymbol
	return nil
}

// extractMethodSymbol 提取方法符号
func (m *Module) extractMethodSymbol(method *hast.FuncDecl, classSymbol *ClassSymbol) error {
	methodSymbol := NewFunctionSymbol(method.Name.Name, "", m.ModuleID)
	methodSymbol.IsMethod = true
	methodSymbol.Receiver = classSymbol.GetName()

	for _, modifier := range method.Modifiers {
		switch modifier.Kind() {
		case hast.ModKindPub:
			methodSymbol.Exported = true
		}
	}

	// 设置返回类型
	if method.Type != nil {
		methodSymbol.ReturnType = m.typeToString(method.Type)
	}

	// 设置参数
	for _, param := range method.Recv {
		if paramExpr, ok := param.(*hast.Parameter); ok {
			param := &Parameter{
				Name: paramExpr.Name.Name,
				Type: m.typeToString(paramExpr.Type),
			}
			methodSymbol.Parameters = append(methodSymbol.Parameters, param)
		}
	}

	// 设置函数体
	methodSymbol.Body = method.Body

	// 添加到类的方法表
	classSymbol.Methods[method.Name.Name] = methodSymbol
	return nil
}

// extractVariableSymbol 提取变量符号
func (m *Module) extractVariableSymbol(assign *hast.AssignStmt, st *SymbolTable) error {
	// 获取变量名

	var varName string
	if ident, ok := assign.Lhs.(*hast.Ident); ok {
		varName = ident.Name
	} else if ref, ok := assign.Lhs.(*hast.RefExpr); ok {
		varName = ref.X.(*hast.Ident).Name
	}

	// 确定变量类型
	var dataType string
	if assign.Type != nil {
		dataType = m.typeToString(assign.Type)
	}

	// 创建变量符号
	varSymbol := NewVariableSymbol(varName, assign.Rhs, dataType, m.ModuleID)

	// 设置修饰符
	for _, modifier := range assign.Modifiers {
		switch modifier.Kind() {
		case hast.ModKindPub:
			varSymbol.Exported = true
		}
	}

	// 根据作用域类型决定添加到哪里
	switch assign.Scope {
	case token.VAR, token.CONST:
		// 全局变量或常量
		varSymbol.ScopeLevel = 0
		varSymbol.SetScope(nil) // 全局作用域
		st.GlobalSymbols[varName] = varSymbol
	case token.LET:
		// 局部变量
		if st.CurrentScope != nil {
			varSymbol.ScopeLevel = st.CurrentScope.Level
			varSymbol.SetScope(st.CurrentScope)
			st.CurrentScope.Symbols[varName] = varSymbol
		} else {
			// 如果没有当前作用域，当作全局变量处理
			varSymbol.ScopeLevel = 0
			varSymbol.SetScope(nil)
			st.GlobalSymbols[varName] = varSymbol
		}
	default:
		// 检查是否是 `:=` 赋值，如果是则当作局部变量处理
		if assign.Tok == token.COLON_ASSIGN {
			varSymbol.ScopeLevel = st.CurrentScope.Level
			varSymbol.SetScope(st.CurrentScope)
			st.CurrentScope.Symbols[varName] = varSymbol
		} else {
			// 普通赋值，不创建新符号
			return nil
		}
	}

	return nil
}

// extractEnumSymbol 提取枚举符号
func (m *Module) extractEnumSymbol(enum *hast.EnumDecl, st *SymbolTable) error {
	// 枚举可以看作是一种特殊的类
	enumSymbol := NewClassSymbol(enum.Name.Name, m.ModuleID)

	// 设置修饰符
	for _, modifier := range enum.Modifiers {
		switch modifier.Kind() {
		case hast.ModKindPub:
			enumSymbol.Exported = true
		}
	}

	// 添加到符号表
	st.GlobalSymbols[enum.Name.Name] = enumSymbol
	return nil
}

// extractTraitSymbol 提取trait符号
func (m *Module) extractTraitSymbol(trait *hast.TraitDecl, st *SymbolTable) error {
	// trait可以看作是一种特殊的类
	traitSymbol := NewClassSymbol(trait.Name.Name, m.ModuleID)

	// 设置修饰符
	for _, modifier := range trait.Modifiers {
		switch modifier.Kind() {
		case hast.ModKindPub:
			traitSymbol.Exported = true
		}
	}

	// 添加到符号表
	st.GlobalSymbols[trait.Name.Name] = traitSymbol
	return nil
}

// extractModuleSymbol 提取模块符号
func (m *Module) extractModuleSymbol(mod *hast.ModDecl, st *SymbolTable) error {
	modSymbol := NewNamespaceSymbol(mod.Name.Name, m.ModuleID)

	// 设置修饰符
	if mod.Pub.IsValid() {
		modSymbol.Exported = true
	}

	// 递归提取模块内部的符号
	if mod.Body != nil {
		// 为模块创建一个新的作用域
		moduleScope := &Scope{
			Parent:  st.CurrentScope,
			Symbols: make(map[string]Symbol),
		}

		// 临时切换到模块作用域
		originalScope := st.CurrentScope
		st.CurrentScope = moduleScope

		// 提取模块内部的符号
		if err := m.extractSymbolsFromBlock(mod.Body, st); err != nil {
			st.CurrentScope = originalScope
			return err
		}

		// 将模块作用域中的符号添加到模块符号中
		maps.Copy(modSymbol.Symbols, moduleScope.Symbols)

		// 恢复原来的作用域
		st.CurrentScope = originalScope
	}

	// 添加到当前作用域的符号表
	if st.CurrentScope != nil {
		// 如果是嵌套模块，添加到当前作用域
		st.CurrentScope.Symbols[mod.Name.Name] = modSymbol
	} else {
		// 如果是顶级模块，添加到全局符号表
		st.GlobalSymbols[mod.Name.Name] = modSymbol
	}

	return nil
}

// extractTypeAliasSymbol 提取类型别名符号
func (m *Module) extractTypeAliasSymbol(typeDecl *hast.TypeDecl, st *SymbolTable) error {
	// 类型别名可以看作是一种常量
	typeSymbol := NewConstantSymbol(typeDecl.Name.Name, typeDecl.Value, "", m.ModuleID)

	// 添加到符号表
	st.GlobalSymbols[typeDecl.Name.Name] = typeSymbol
	return nil
}

// extractSymbolsFromBlock 从代码块中提取符号
func (m *Module) extractSymbolsFromBlock(block *hast.BlockStmt, st *SymbolTable) error {
	// 进入块作用域
	// st.EnterScope(fmt.Sprintf("block:%d", st.CurrentScope.Level+1))
	// defer st.ExitScope()

	// 提取块中的符号
	for _, stmt := range block.List {
		if err := m.extractSymbolFromStmt(stmt, st); err != nil {
			return err
		}
	}

	return nil
}

// typeToString 将类型表达式转换为字符串
func (m *Module) typeToString(typeExpr hast.Expr) string {
	if typeExpr == nil {
		return "any"
	}

	// 这里需要根据具体的类型表达式来转换
	switch t := typeExpr.(type) {
	case *hast.Ident:
		return t.Name
	case *hast.TypeReference:
		if ident, ok := t.Name.(*hast.Ident); ok {
			return ident.Name
		}
		return m.typeToString(t.Name)
	default:
		return "any" // 默认类型
	}
}

func (m *Module) TransformAST() error {
	return m.transformNode(m.AST)
}

func (m *Module) transformNode(node hast.Node) error {
	switch n := node.(type) {
	case *hast.File:
		for _, stmt := range n.Stmts {
			if err := m.transformNode(stmt); err != nil {
				return err
			}
		}
	case *hast.Ident:
		// 查找符号，替换为混淆名称
		if symbol := m.Symbols.LookupSymbol(n.Name); symbol != nil && symbol.IsMangled() {
			n.Name = symbol.GetMangledName() // 直接修改AST！
		}

	case *hast.AssignStmt:
		// 检查 Lhs 是否为 Ident 类型
		if left, ok := n.Lhs.(*hast.Ident); ok {
			// 变量声明中的标识符
			if symbol := m.Symbols.LookupSymbol(left.Name); symbol != nil && symbol.IsMangled() {
				left.Name = symbol.GetMangledName()
			}
		}

	case *hast.CallExpr:
		// 检查 Fun 是否为 Ident 类型
		if ident, ok := n.Fun.(*hast.Ident); ok {
			// 函数调用中的标识符
			if symbol := m.Symbols.LookupSymbol(ident.Name); symbol != nil && symbol.IsMangled() {
				ident.Name = symbol.GetMangledName()
			}
		}
	case *hast.CmdExpr:
		if ident, ok := n.Cmd.(*hast.Ident); ok {
			if symbol := m.Symbols.LookupSymbol(ident.Name); symbol != nil && symbol.IsMangled() {
				ident.Name = symbol.GetMangledName()
			}
		}
		for _, arg := range n.Args {
			if err := m.transformNode(arg); err != nil {
				return err
			}
		}

	case *hast.BlockStmt:
		// 递归处理块中的语句
		for _, stmt := range n.List {
			if err := m.transformNode(stmt); err != nil {
				return err
			}
		}

	case *hast.FuncDecl:
		// 处理函数声明中的标识符
		if symbol := m.Symbols.LookupSymbol(n.Name.Name); symbol != nil && symbol.IsMangled() {
			n.Name.Name = symbol.GetMangledName()
		}
		// 递归处理函数体
		if n.Body != nil {
			if err := m.transformNode(n.Body); err != nil {
				return err
			}
		}

	case *hast.ClassDecl:
		// 处理类声明中的标识符
		if symbol := m.Symbols.LookupSymbol(n.Name.Name); symbol != nil && symbol.IsMangled() {
			n.Name.Name = symbol.GetMangledName()
		}
		// 递归处理类体
		if n.Fields != nil {
			for _, field := range n.Fields.List {
				// 处理字段名称
				if symbol := m.Symbols.LookupSymbol(field.Name.Name); symbol != nil && symbol.IsMangled() {
					field.Name.Name = symbol.GetMangledName()
				}
				// 处理字段类型
				if field.Type != nil {
					if err := m.transformNode(field.Type); err != nil {
						return err
					}
				}
				// 处理字段默认值
				if field.Value != nil {
					if err := m.transformNode(field.Value); err != nil {
						return err
					}
				}
			}
		}
		// 处理方法
		for _, method := range n.Methods {
			if err := m.transformNode(method); err != nil {
				return err
			}
		}

	case *hast.ReturnStmt:
		// 处理返回语句中的表达式
		if n.X != nil {
			if err := m.transformNode(n.X); err != nil {
				return err
			}
		}

	case *hast.BinaryExpr:
		// 处理二元表达式
		if err := m.transformNode(n.X); err != nil {
			return err
		}
		if err := m.transformNode(n.Y); err != nil {
			return err
		}

	case *hast.RefExpr:
		// 处理引用表达式 $x
		if err := m.transformNode(n.X); err != nil {
			return err
		}

	case *hast.IncDecExpr:
		// 处理自增自减表达式 x++ 或 ++x
		if ident, ok := n.X.(*hast.Ident); ok {
			if symbol := m.Symbols.LookupSymbol(ident.Name); symbol != nil && symbol.IsMangled() {
				ident.Name = symbol.GetMangledName()
			}
		} else {
			// 如果不是标识符，递归处理
			if err := m.transformNode(n.X); err != nil {
				return err
			}
		}

	case *hast.SelectExpr:
		// 处理选择表达式 x.y
		if err := m.transformNode(n.X); err != nil {
			return err
		}
		if err := m.transformNode(n.Y); err != nil {
			return err
		}

	case *hast.ModAccessExpr:
		// 处理模块访问表达式 x::y
		if err := m.transformNode(n.X); err != nil {
			return err
		}
		if err := m.transformNode(n.Y); err != nil {
			return err
		}

	case *hast.IfStmt:
		// 处理 if 语句
		if err := m.transformNode(n.Cond); err != nil {
			return err
		}
		if err := m.transformNode(n.Body); err != nil {
			return err
		}
		if n.Else != nil {
			if err := m.transformNode(n.Else); err != nil {
				return err
			}
		}

	case *hast.WhileStmt:
		// 处理 while 语句
		if err := m.transformNode(n.Cond); err != nil {
			return err
		}
		if err := m.transformNode(n.Body); err != nil {
			return err
		}

	case *hast.ForStmt:
		// 处理 for 语句
		if n.Init != nil {
			// Init 是 Stmt 类型，可能是 AssignStmt
			if assignStmt, ok := n.Init.(*hast.AssignStmt); ok {
				// 处理赋值语句中的变量
				if left, ok := assignStmt.Lhs.(*hast.Ident); ok {
					if symbol := m.Symbols.LookupSymbol(left.Name); symbol != nil && symbol.IsMangled() {
						left.Name = symbol.GetMangledName()
					}
				}
				// 递归处理右侧表达式
				if err := m.transformNode(assignStmt.Rhs); err != nil {
					return err
				}
			} else {
				// 其他类型的语句
				if err := m.transformNode(n.Init); err != nil {
					return err
				}
			}
		}
		if err := m.transformNode(n.Cond); err != nil {
			return err
		}
		if err := m.transformNode(n.Post); err != nil {
			return err
		}
		if err := m.transformNode(n.Body); err != nil {
			return err
		}

	case *hast.ForeachStmt:
		// 处理 foreach 语句
		if err := m.transformNode(n.Index); err != nil {
			return err
		}
		if err := m.transformNode(n.Value); err != nil {
			return err
		}
		if err := m.transformNode(n.Var); err != nil {
			return err
		}
		if err := m.transformNode(n.Body); err != nil {
			return err
		}

	case *hast.ForInStmt:
		// 处理 for-in 语句
		if err := m.transformNode(n.Index); err != nil {
			return err
		}
		// 处理 RangeExpr 的各个部分
		if err := m.transformNode(n.RangeExpr.Start); err != nil {
			return err
		}
		if err := m.transformNode(n.RangeExpr.End_); err != nil {
			return err
		}
		if err := m.transformNode(n.RangeExpr.Step); err != nil {
			return err
		}
		if err := m.transformNode(n.Body); err != nil {
			return err
		}

	case *hast.ModDecl:
		// 处理模块声明
		if symbol := m.Symbols.LookupSymbol(n.Name.Name); symbol != nil && symbol.IsMangled() {
			n.Name.Name = symbol.GetMangledName()
		}
		if n.Body != nil {
			if err := m.transformNode(n.Body); err != nil {
				return err
			}
		}

	case *hast.EnumDecl:
		// 处理枚举声明
		if symbol := m.Symbols.LookupSymbol(n.Name.Name); symbol != nil && symbol.IsMangled() {
			n.Name.Name = symbol.GetMangledName()
		}

	case *hast.TraitDecl:
		// 处理 trait 声明
		if symbol := m.Symbols.LookupSymbol(n.Name.Name); symbol != nil && symbol.IsMangled() {
			n.Name.Name = symbol.GetMangledName()
		}

	case *hast.TypeDecl:
		// 处理类型声明
		if symbol := m.Symbols.LookupSymbol(n.Name.Name); symbol != nil && symbol.IsMangled() {
			n.Name.Name = symbol.GetMangledName()
		}
		if err := m.transformNode(n.Value); err != nil {
			return err
		}
	}

	return nil
}

func MangleModule(module *Module) error {
	if err := module.TransformAST(); err != nil {
		return err
	}
	return nil
}

type DependecyResolver struct {
	modules         map[string]*Module
	visited         container.Set[string]
	stack           container.Set[string]
	order           []string
	current         *Module
	pkgs            map[string]*config.HuloPkg
	options         *config.Huloc
	fs              vfs.VFS
	huloModulesPath string
	huloPath        string
}

func NewDependecyResolver() *DependecyResolver {
	return &DependecyResolver{}
}

// 返回绝对路径
func (r *DependecyResolver) resolvePath(parent string, path string) (string, error) {
	log.WithField("parent", parent).Infof("resolvePath: %s", path)
	if r.isRelativePath(path) {
		ret, err := filepath.Abs(filepath.Join(filepath.Dir(parent), path))
		if err != nil {
			return "", fmt.Errorf("failed to resolve path: %w", err)
		}
		if r.fs.Exists(ret) {
			stat, err := r.fs.Stat(ret)
			if err != nil {
				return "", fmt.Errorf("failed to stat path: %w", err)
			}
			if !stat.IsDir() {
				return ret, nil
			}
		}
		if r.fs.Exists(ret + ".hl") {
			return ret + ".hl", nil
		}
		if r.fs.Exists(filepath.Join(ret, "index.hl")) {
			return filepath.Join(ret, "index.hl"), nil
		}
		return "", fmt.Errorf("no file found for %s", ret)
	}
	if r.isCoreLibPath(path) {
		return filepath.Join(r.huloPath, path, "index.hl"), nil
	}

	// hulo-lang/hello-world hulo-lang/hello-world/abc
	// 要区分是 abc.hl abc/index.hl 以及 hello-world 的入口点在什么地方
	// 根据版本信息，以及suffix 以及 main 表示文件夹入口 确认到底哪个文件
	if ownerRepoVer, suffix, ok := r.isModulePath(path); ok {
		if pkg, ok := r.pkgs[ownerRepoVer]; ok {
			pend := filepath.Join(r.huloModulesPath, ownerRepoVer, pkg.Main, suffix)
			if r.fs.Exists(pend + ".hl") {
				return pend + ".hl", nil
			}
			if r.fs.Exists(filepath.Join(pend, "index.hl")) {
				return filepath.Join(pend, "index.hl"), nil
			}
			return "", fmt.Errorf("no index.hl found for %s", pend)
		}
		pkgPath := filepath.Join(r.huloModulesPath, ownerRepoVer, config.HuloPkgFileName)
		if !r.fs.Exists(pkgPath) {
			return "", fmt.Errorf("no hulo.pkg.yaml found for %s", pkgPath)
		}
		content, err := r.fs.ReadFile(pkgPath)
		if err != nil {
			return "", fmt.Errorf("failed to read pkg file: %w", err)
		}
		pkg := &config.HuloPkg{}
		if err := yaml.Unmarshal(content, pkg); err != nil {
			return "", fmt.Errorf("failed to unmarshal pkg file: %w", err)
		}
		r.pkgs[ownerRepoVer] = pkg
		pend := filepath.Join(r.huloModulesPath, ownerRepoVer, pkg.Main, suffix)
		if r.fs.Exists(filepath.Join(pend, "index.hl")) {
			return filepath.Join(pend, "index.hl"), nil
		}
		return "", fmt.Errorf("no index.hl found for %s", pend)
	}

	return "", fmt.Errorf("invalid path: %s", path)
}

func (r *DependecyResolver) isCoreLibPath(path string) bool {
	return r.fs.Exists(filepath.Join(r.huloPath, path, "index.hl"))
}

func (r *DependecyResolver) isRelativePath(path string) bool {
	return strings.HasPrefix(path, ".")
}

func (r *DependecyResolver) isModulePath(path string) (string, string, bool) {
	if r.current.Pkg == nil {
		return "", "", false
	}

	for name, version := range r.current.Pkg.Dependencies {
		if strings.HasPrefix(path, name) {
			symbolName := strings.TrimPrefix(path, name)
			return filepath.Join(name, version), symbolName, true
		}
	}

	return "", "", false
}

func (r *DependecyResolver) resolveRecursive(parentPkg *Module, parent, filepath string) error {
	absPath, err := r.resolvePath(parent, filepath)
	if err != nil {
		return err
	}
	log.IncreasePadding()
	defer log.DecreasePadding()

	if r.visited.Contains(absPath) {
		return nil
	}

	if r.stack.Contains(absPath) {
		return fmt.Errorf("circular dependency detected: %s", absPath)
	}

	r.visited.Add(absPath)
	r.stack.Add(absPath)
	defer func() { r.stack.Remove(absPath) }()

	log.Infof("loadModule: %s", absPath)
	module, err := r.loadModule(parentPkg, absPath, filepath)
	if err != nil {
		return err
	}

	// 保存当前的 current 模块
	oldCurrent := r.current
	r.current = module
	defer func() { r.current = oldCurrent }()

	// 递归解析依赖
	for _, dep := range module.Dependencies {
		if err := r.resolveRecursive(module, absPath, dep); err != nil {
			return err
		}
	}

	r.order = append(r.order, absPath)

	return nil
}

func (r *DependecyResolver) findNearestPkgYaml(filePath string) (string, error) {
	dir := filepath.Dir(filePath)
	for {
		pkgPath := filepath.Join(dir, config.HuloPkgFileName)
		if r.fs.Exists(pkgPath) {
			return pkgPath, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			// 已经到达根目录
			break
		}
		dir = parent
	}
	return "", fmt.Errorf("no hulo.pkg.yaml found for %s", filePath)
}

// 收集 importInfo 以及 载入配置文件
func (r *DependecyResolver) loadModule(parentPkg *Module, absPath string, symbolName string) (*Module, error) {
	if module, ok := r.modules[absPath]; ok {
		return module, nil
	}

	content, err := r.fs.ReadFile(absPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", absPath, err)
	}

	opts := []parser.ParserOptions{}
	if len(r.options.Parser.ShowASTTree) > 0 {
		switch r.options.Parser.ShowASTTree {
		case "stdout":
			opts = append(opts, parser.OptionTracerASTTree(os.Stdout))
		case "stderr":
			opts = append(opts, parser.OptionTracerASTTree(os.Stderr))
		case "file":
			file, err := r.fs.Open(r.options.Parser.ShowASTTree)
			if err != nil {
				return nil, fmt.Errorf("failed to open file %s: %w", r.options.Parser.ShowASTTree, err)
			}
			defer file.Close()
			opts = append(opts, parser.OptionTracerASTTree(file))
		}
	}
	if !r.options.Parser.EnableTracer {
		opts = append(opts, parser.OptionDisableTracer())
	}
	if r.options.Parser.DisableTiming {
		opts = append(opts, parser.OptionTracerDisableTiming())
	}

	log.WithField("symbolName", symbolName).Infof("parse file: %s", absPath)
	ast, err := parser.ParseSourceScript(string(content), opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to parse file %s: %w", absPath, err)
	}

	imports, dependencies := r.extractImports(ast, absPath)
	log.WithField("dependencies", dependencies).Infof("extractImports: %s", absPath)

	module := &Module{
		Path:         absPath,
		Imports:      imports,
		Dependencies: dependencies,
		AST:          ast,
		Exports:      make(map[string]*ExportInfo),
	}

	if r.isRelativePath(symbolName) {
		pkgPath, err := r.findNearestPkgYaml(absPath)
		if err != nil {
			return nil, fmt.Errorf("failed to find nearest pkg file: %w", err)
		}
		content, err := r.fs.ReadFile(pkgPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read pkg file: %w", err)
		}
		module.Pkg = &config.HuloPkg{}
		if err := yaml.Unmarshal(content, module.Pkg); err != nil {
			return nil, fmt.Errorf("failed to unmarshal pkg file: %w", err)
		}

		r.pkgs[pkgPath] = module.Pkg
	} else if _, _, ok := r.isModulePath(symbolName); ok {
		// 按照 owner/repo/version 的路径去拿包下面的 hulo.pkg.yaml 文件
		pkgPath, err := r.findNearestPkgYaml(absPath)
		if err != nil {
			return nil, fmt.Errorf("failed to find nearest pkg file: %w", err)
		}
		// pkgPath 出来后要读取这个文件
		content, err := r.fs.ReadFile(pkgPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read pkg file: %w", err)
		}

		module.Pkg = &config.HuloPkg{}
		if err := yaml.Unmarshal(content, module.Pkg); err != nil {
			return nil, fmt.Errorf("failed to unmarshal pkg file: %w", err)
		}

		r.pkgs[pkgPath] = module.Pkg
	}

	// 将模块添加到 modules 中
	r.modules[absPath] = module

	return module, nil
}

func (r *DependecyResolver) extractImports(ast *hast.File, currentModulePath string) ([]*ImportInfo, []string) {
	imports := []*ImportInfo{}
	dependencies := []string{}
	dependencySet := container.NewMapSet[string]()

	for _, stmt := range ast.Stmts {
		if importStmt, ok := stmt.(*hast.Import); ok {
			var importInfo *ImportInfo

			switch {
			case importStmt.ImportSingle != nil:
				importInfo = &ImportInfo{
					ModulePath: importStmt.ImportSingle.Path,
					Alias:      importStmt.ImportSingle.Alias,
					Kind:       ImportSingle,
				}
			case importStmt.ImportMulti != nil:
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
				importInfo = &ImportInfo{
					ModulePath: importStmt.ImportAll.Path,
					Alias:      importStmt.ImportAll.Alias,
					Kind:       ImportAll,
				}
			}

			if importInfo != nil {
				imports = append(imports, importInfo)
				dependencySet.Add(importInfo.ModulePath)
			}
		}
	}

	for _, dep := range dependencySet.Items() {
		dependencies = append(dependencies, dep)
	}

	return imports, dependencies
}

func (r *DependecyResolver) parseSymbolName(symbolName string) (string, string) {
	return "", ""
}
