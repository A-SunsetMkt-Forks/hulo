package module

import "fmt"

type Symbol struct {
	Name        string // 原始名称
	MangledName string // 混淆名称
	Kind        SymbolKind

	Value    any
	DataType string // 数据类型，可能在对象的方法调用的时候混淆用

	ScopeLevel int  // 作用域层级
	IsExported bool // 是否导出 (pub)
	IsImported bool // 是否导入 (import)
	ModuleID   int  // 模块ID
	Scope      *Scope

	IsMangled bool // 是否已混淆
}

type SymbolKind int

const (
	SymbolVar SymbolKind = iota
	SymbolFunc
	SymbolConst
	SymbolClass
	SymbolNamespace
)

type Scope struct {
	Parent  *Scope
	Symbols map[string]*Symbol

	Variables map[string]*Symbol
	Functions map[string]*Symbol
	Constants map[string]*Symbol

	Level int
	Name  string // 作用域名称，如 "global", "function:add", "block:1"
}

type SymbolTable struct {
	// 全局符号映射
	GlobalSymbols map[string]*Symbol
	// 作用域栈
	Scopes []*Scope
	// 当前作用域
	CurrentScope *Scope
	// 模块ID
	ModuleID int
	// 混淆器
	Mangler *SymbolMangler
}

func (st *SymbolTable) LookupSymbol(name string) *Symbol {
	// 1. 在当前作用域查找
	if st.CurrentScope != nil {
		if symbol, exists := st.CurrentScope.Variables[name]; exists {
			return symbol
		}
		if symbol, exists := st.CurrentScope.Functions[name]; exists {
			return symbol
		}
		if symbol, exists := st.CurrentScope.Constants[name]; exists {
			return symbol
		}
	}

	// 2. 向上查找父作用域
	current := st.CurrentScope
	for current != nil && current.Parent != nil {
		current = current.Parent
		if symbol, exists := current.Variables[name]; exists {
			return symbol
		}
		if symbol, exists := current.Functions[name]; exists {
			return symbol
		}
		if symbol, exists := current.Constants[name]; exists {
			return symbol
		}
	}

	// 3. 查找全局作用域
	if symbol, exists := st.GlobalSymbols[name]; exists {
		return symbol
	}

	return nil
}

func (st *SymbolTable) DeclareVariable(name string, value interface{}, dataType string) error {
	// 检查是否已存在
	if st.LookupSymbol(name) != nil {
		return fmt.Errorf("symbol %s already declared", name)
	}

	// 创建符号
	symbol := &Symbol{
		Name:       name,
		Kind:       SymbolVar,
		Value:      value,
		DataType:   dataType,
		ScopeLevel: st.CurrentScope.Level,
		ModuleID:   st.ModuleID,
		IsMangled:  false,
	}

	// 添加到当前作用域
	if st.CurrentScope != nil {
		st.CurrentScope.Variables[name] = symbol
	} else {
		st.GlobalSymbols[name] = symbol
	}

	return nil
}

func (st *SymbolTable) MangleAllSymbols() error {
	// 混淆全局符号
	for name, symbol := range st.GlobalSymbols {
		if !symbol.IsMangled {
			symbol.MangledName = st.Mangler.MangleSymbol(st.ModuleID, symbol.Kind, name)
			symbol.IsMangled = true
		}
	}

	// 混淆所有作用域中的符号
	for _, scope := range st.Scopes {
		for name, symbol := range scope.Variables {
			if !symbol.IsMangled {
				symbol.MangledName = st.Mangler.MangleSymbol(st.ModuleID, symbol.Kind, name)
				symbol.IsMangled = true
			}
		}
		for name, symbol := range scope.Functions {
			if !symbol.IsMangled {
				symbol.MangledName = st.Mangler.MangleSymbol(st.ModuleID, symbol.Kind, name)
				symbol.IsMangled = true
			}
		}
		for name, symbol := range scope.Constants {
			if !symbol.IsMangled {
				symbol.MangledName = st.Mangler.MangleSymbol(st.ModuleID, symbol.Kind, name)
				symbol.IsMangled = true
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
