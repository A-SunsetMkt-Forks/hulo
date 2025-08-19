// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package module

import "fmt"

type Symbol interface {
	GetName() string
	GetMangledName() string
	SetMangledName(name string)
	GetKind() SymbolKind
	GetScopeLevel() int
	GetModuleID() int
	IsExported() bool
	IsImported() bool
	IsMangled() bool
	SetMangled(mangled bool)
	GetScope() *Scope
	SetScope(scope *Scope)
}

type SymbolKind int

const (
	SymbolVar SymbolKind = iota
	SymbolFunc
	SymbolConst
	SymbolClass
	SymbolEnum
	SymbolNamespace
)

// BaseSymbol 基础符号实现
type BaseSymbol struct {
	Name        string
	MangledName string
	Kind        SymbolKind
	ScopeLevel  int
	ModuleID    int
	Exported    bool
	Imported    bool
	Mangled     bool
	Scope       *Scope
}

func (s *BaseSymbol) GetName() string            { return s.Name }
func (s *BaseSymbol) GetMangledName() string     { return s.MangledName }
func (s *BaseSymbol) SetMangledName(name string) { s.MangledName = name }
func (s *BaseSymbol) GetKind() SymbolKind        { return s.Kind }
func (s *BaseSymbol) GetScopeLevel() int         { return s.ScopeLevel }
func (s *BaseSymbol) GetModuleID() int           { return s.ModuleID }
func (s *BaseSymbol) IsExported() bool           { return s.Exported }
func (s *BaseSymbol) IsImported() bool           { return s.Imported }
func (s *BaseSymbol) IsMangled() bool            { return s.Mangled }
func (s *BaseSymbol) SetMangled(mangled bool)    { s.Mangled = mangled }
func (s *BaseSymbol) GetScope() *Scope           { return s.Scope }
func (s *BaseSymbol) SetScope(scope *Scope)      { s.Scope = scope }

// VariableSymbol 变量符号
type VariableSymbol struct {
	BaseSymbol
	Value    any
	DataType string
}

func NewVariableSymbol(name string, value any, dataType string, moduleID int) *VariableSymbol {
	return &VariableSymbol{
		BaseSymbol: BaseSymbol{
			Name:     name,
			Kind:     SymbolVar,
			ModuleID: moduleID,
			Mangled:  false,
		},
		Value:    value,
		DataType: dataType,
	}
}

// FunctionSymbol 函数符号
type FunctionSymbol struct {
	BaseSymbol
	Parameters []*Parameter
	ReturnType string
	Body       any    // 函数体AST
	IsMethod   bool   // 是否是类方法
	Receiver   string // 接收者类型（如果是方法）
}

type Parameter struct {
	Name string
	Type string
}

func NewFunctionSymbol(name string, returnType string, moduleID int) *FunctionSymbol {
	return &FunctionSymbol{
		BaseSymbol: BaseSymbol{
			Name:     name,
			Kind:     SymbolFunc,
			ModuleID: moduleID,
			Mangled:  false,
		},
		ReturnType: returnType,
		Parameters: make([]*Parameter, 0),
	}
}

// ConstantSymbol 常量符号
type ConstantSymbol struct {
	BaseSymbol
	Value    any
	DataType string
}

func NewConstantSymbol(name string, value any, dataType string, moduleID int) *ConstantSymbol {
	return &ConstantSymbol{
		BaseSymbol: BaseSymbol{
			Name:     name,
			Kind:     SymbolConst,
			ModuleID: moduleID,
			Mangled:  false,
		},
		Value:    value,
		DataType: dataType,
	}
}

// ClassSymbol 类符号
type ClassSymbol struct {
	BaseSymbol
	fields     map[string]*Field
	methods    map[string]*FunctionSymbol
	superClass string
	traits     []string
}

func (c *ClassSymbol) GetField(name string) *Field {
	return c.fields[name]
}

func (c *ClassSymbol) GetMethod(name string) *FunctionSymbol {
	return c.methods[name]
}

func (c *ClassSymbol) Fields() map[string]*Field {
	return c.fields
}

func (c *ClassSymbol) Methods() map[string]*FunctionSymbol {
	return c.methods
}

func (c *ClassSymbol) Traits() []string {
	return c.traits
}

func (c *ClassSymbol) SuperClass() string {
	return c.superClass
}

type Field struct {
	name  string
	typ   string
	isPub bool
	val   any
}

func (f *Field) Name() string {
	return f.name
}

func (f *Field) Type() string {
	return f.typ
}

func (f *Field) IsPublic() bool {
	return f.isPub
}

func (f *Field) Value() any {
	return f.val
}

type EnumValue struct {
	Name  string
	Value string
	Type  string
}

type EnumSymbol struct {
	BaseSymbol
	Values []*EnumValue
}

func NewEnumSymbol(name string, moduleID int) *EnumSymbol {
	return &EnumSymbol{
		BaseSymbol: BaseSymbol{
			Name:     name,
			Kind:     SymbolEnum,
			ModuleID: moduleID,
			Mangled:  false,
		},
	}
}

func NewClassSymbol(name string, moduleID int) *ClassSymbol {
	return &ClassSymbol{
		BaseSymbol: BaseSymbol{
			Name:     name,
			Kind:     SymbolClass,
			ModuleID: moduleID,
			Mangled:  false,
		},
		fields:  make(map[string]*Field),
		methods: make(map[string]*FunctionSymbol),
		traits:  make([]string, 0),
	}
}

// NamespaceSymbol 命名空间符号
type NamespaceSymbol struct {
	BaseSymbol
	Symbols map[string]Symbol
}

func NewNamespaceSymbol(name string, moduleID int) *NamespaceSymbol {
	return &NamespaceSymbol{
		BaseSymbol: BaseSymbol{
			Name:     name,
			Kind:     SymbolNamespace,
			ModuleID: moduleID,
			Mangled:  false,
		},
		Symbols: make(map[string]Symbol),
	}
}

type Scope struct {
	Parent  *Scope
	Symbols map[string]Symbol // 统一的符号表，函数、常量、变量都不能同名

	Level int
	Name  string // 作用域名称，如 "global", "function:add", "block:1"
}

type SymbolTable struct {
	// 全局符号映射
	GlobalSymbols map[string]Symbol
	// 作用域栈
	ScopeStack []*Scope
	// 当前作用域
	CurrentScope *Scope
	// 模块ID
	ModuleID int
	// 混淆器
	Mangler *SymbolMangler
}

func (st *SymbolTable) LookupSymbol(name string) Symbol {
	// 1. 在当前作用域查找
	if st.CurrentScope != nil {
		if symbol, exists := st.CurrentScope.Symbols[name]; exists {
			return symbol
		}
	}

	// 2. 向上查找父作用域
	current := st.CurrentScope
	for current != nil && current.Parent != nil {
		current = current.Parent
		if symbol, exists := current.Symbols[name]; exists {
			return symbol
		}
	}

	// 3. 查找全局作用域
	if symbol, exists := st.GlobalSymbols[name]; exists {
		return symbol
	}

	return nil
}

func (st *SymbolTable) HasSymbol(name string) bool {
	return st.LookupSymbol(name) != nil
}

// EnterScope 进入新作用域
func (st *SymbolTable) EnterScope(name string) {
	newScope := &Scope{
		Parent:  st.CurrentScope,
		Symbols: make(map[string]Symbol),
		Level: func() int {
			if st.CurrentScope != nil {
				return st.CurrentScope.Level + 1
			}
			return 0
		}(),
		Name: name,
	}

	st.ScopeStack = append(st.ScopeStack, newScope)
	st.CurrentScope = newScope
}

// ExitScope 退出当前作用域
func (st *SymbolTable) ExitScope() {
	if len(st.ScopeStack) > 0 {
		st.ScopeStack = st.ScopeStack[:len(st.ScopeStack)-1]
		if len(st.ScopeStack) > 0 {
			st.CurrentScope = st.ScopeStack[len(st.ScopeStack)-1]
		} else {
			st.CurrentScope = nil
		}
	}
}

func (st *SymbolTable) DeclareVariable(name string, value interface{}, dataType string) error {
	return st.DeclareSymbol(name, SymbolVar, value, dataType)
}

func (st *SymbolTable) MangleAllSymbols() error {
	// 混淆全局符号
	for name, symbol := range st.GlobalSymbols {
		if !symbol.IsMangled() {
			symbol.SetMangledName(st.Mangler.MangleSymbol(st.ModuleID, symbol.GetKind(), name))
			symbol.SetMangled(true)
		}
	}

	// 混淆所有作用域中的符号
	for _, scope := range st.ScopeStack {
		for name, symbol := range scope.Symbols {
			if !symbol.IsMangled() {
				symbol.SetMangledName(st.Mangler.MangleSymbolWithScope(st.ModuleID, symbol.GetKind(), name, symbol.GetScopeLevel()))
				symbol.SetMangled(true)
			}
		}
	}

	return nil
}

type SymbolMangler struct {
	moduleCounter int
	typeCounter   map[int]int       // 模块ID -> 类型计数器
	symbolCounter map[string]int    // "模块ID_类型ID" -> 符号计数器
	mangledMap    map[string]string // 原始名称 -> 混淆名称
}

func (sm *SymbolMangler) MangleSymbol(moduleID int, symbolType SymbolKind, originalName string) string {
	return originalName
	// 检查是否已经混淆过
	if mangled, exists := sm.mangledMap[originalName]; exists {
		return mangled
	}

	// 生成新的混淆名称
	typeID := int(symbolType)
	key := fmt.Sprintf("%d_%d", moduleID, typeID)
	symbolID := sm.symbolCounter[key]
	sm.symbolCounter[key]++

	mangledName := fmt.Sprintf("_%d_%d_%d", moduleID, typeID, symbolID)

	// 保存映射关系
	sm.mangledMap[originalName] = mangledName

	return mangledName
}

// MangleSymbolWithScope 带作用域级别的混淆
func (sm *SymbolMangler) MangleSymbolWithScope(moduleID int, symbolType SymbolKind, originalName string, scopeLevel int) string {
	return originalName
	// 检查是否已经混淆过
	if mangled, exists := sm.mangledMap[originalName]; exists {
		return mangled
	}

	// 生成新的混淆名称，包含作用域级别
	typeID := int(symbolType)
	key := fmt.Sprintf("%d_%d_%d", moduleID, typeID, scopeLevel)
	symbolID := sm.symbolCounter[key]
	sm.symbolCounter[key]++

	// 混淆名称中包含作用域级别
	mangledName := fmt.Sprintf("_%d_%d_%d_%d", moduleID, typeID, scopeLevel, symbolID)

	// 保存映射关系
	sm.mangledMap[originalName] = mangledName

	return mangledName
}

// DeclareSymbol 声明符号（统一方法）
func (st *SymbolTable) DeclareSymbol(name string, kind SymbolKind, value interface{}, dataType string) error {
	// 检查是否已存在
	if st.LookupSymbol(name) != nil {
		return fmt.Errorf("symbol %s already declared", name)
	}

	// 创建符号
	var symbol Symbol
	switch kind {
	case SymbolVar:
		symbol = NewVariableSymbol(name, value, dataType, st.ModuleID)
	case SymbolFunc:
		symbol = NewFunctionSymbol(name, dataType, st.ModuleID)
	case SymbolConst:
		symbol = NewConstantSymbol(name, value, dataType, st.ModuleID)
	case SymbolClass:
		symbol = NewClassSymbol(name, st.ModuleID)
	case SymbolNamespace:
		symbol = NewNamespaceSymbol(name, st.ModuleID)
	default:
		return fmt.Errorf("unknown symbol kind: %d", kind)
	}

	// 添加到当前作用域
	if st.CurrentScope != nil {
		st.CurrentScope.Symbols[name] = symbol
	} else {
		st.GlobalSymbols[name] = symbol
	}

	return nil
}

// DeclareFunction 声明函数
func (st *SymbolTable) DeclareFunction(name string, value interface{}, dataType string) error {
	return st.DeclareSymbol(name, SymbolFunc, value, dataType)
}

// DeclareConstant 声明常量
func (st *SymbolTable) DeclareConstant(name string, value interface{}, dataType string) error {
	return st.DeclareSymbol(name, SymbolConst, value, dataType)
}

// DeclareClass 声明类
func (st *SymbolTable) DeclareClass(name string, value interface{}, dataType string) error {
	return st.DeclareSymbol(name, SymbolClass, value, dataType)
}

// GetSymbolsByKind 根据类型获取符号
func (st *SymbolTable) GetSymbolsByKind(kind SymbolKind) []Symbol {
	var symbols []Symbol

	// 从全局符号中查找
	for _, symbol := range st.GlobalSymbols {
		if symbol.GetKind() == kind {
			symbols = append(symbols, symbol)
		}
	}

	// 从当前作用域中查找
	if st.CurrentScope != nil {
		for _, symbol := range st.CurrentScope.Symbols {
			if symbol.GetKind() == kind {
				symbols = append(symbols, symbol)
			}
		}
	}

	return symbols
}

// 类型断言辅助方法
func AsVariableSymbol(symbol Symbol) (*VariableSymbol, bool) {
	if vs, ok := symbol.(*VariableSymbol); ok {
		return vs, true
	}
	return nil, false
}

func AsFunctionSymbol(symbol Symbol) (*FunctionSymbol, bool) {
	if fs, ok := symbol.(*FunctionSymbol); ok {
		return fs, true
	}
	return nil, false
}

func AsConstantSymbol(symbol Symbol) (*ConstantSymbol, bool) {
	if cs, ok := symbol.(*ConstantSymbol); ok {
		return cs, true
	}
	return nil, false
}

func AsEnumSymbol(symbol Symbol) (*EnumSymbol, bool) {
	if es, ok := symbol.(*EnumSymbol); ok {
		return es, true
	}
	return nil, false
}

func AsClassSymbol(symbol Symbol) (*ClassSymbol, bool) {
	if cs, ok := symbol.(*ClassSymbol); ok {
		return cs, true
	}
	return nil, false
}

func AsNamespaceSymbol(symbol Symbol) (*NamespaceSymbol, bool) {
	if ns, ok := symbol.(*NamespaceSymbol); ok {
		return ns, true
	}
	return nil, false
}

// 符号表辅助方法
func (st *SymbolTable) DeclareVariableSymbol(name string, value interface{}, dataType string) error {
	return st.DeclareSymbol(name, SymbolVar, value, dataType)
}

func (st *SymbolTable) DeclareFunctionSymbol(name string, returnType string) error {
	symbol := NewFunctionSymbol(name, returnType, st.ModuleID)

	// 检查是否已存在
	if st.LookupSymbol(name) != nil {
		return fmt.Errorf("symbol %s already declared", name)
	}

	// 添加到当前作用域
	if st.CurrentScope != nil {
		st.CurrentScope.Symbols[name] = symbol
	} else {
		st.GlobalSymbols[name] = symbol
	}

	return nil
}

func (st *SymbolTable) DeclareClassSymbol(name string) error {
	symbol := NewClassSymbol(name, st.ModuleID)

	// 检查是否已存在
	if st.LookupSymbol(name) != nil {
		return fmt.Errorf("symbol %s already declared", name)
	}

	// 添加到当前作用域
	if st.CurrentScope != nil {
		st.CurrentScope.Symbols[name] = symbol
	} else {
		st.GlobalSymbols[name] = symbol
	}

	return nil
}

func (st *SymbolTable) DeclareNamespaceSymbol(name string) error {
	symbol := NewNamespaceSymbol(name, st.ModuleID)

	// 检查是否已存在
	if st.LookupSymbol(name) != nil {
		return fmt.Errorf("symbol %s already declared", name)
	}

	// 添加到当前作用域
	if st.CurrentScope != nil {
		st.CurrentScope.Symbols[name] = symbol
	} else {
		st.GlobalSymbols[name] = symbol
	}

	return nil
}
