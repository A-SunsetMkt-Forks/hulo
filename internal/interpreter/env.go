package interpreter

import (
	"fmt"

	"github.com/hulo-lang/hulo/internal/object"
	"github.com/hulo-lang/hulo/syntax/hulo/token"
)

// VariableInfo 存储变量的详细信息
type VariableInfo struct {
	Value    object.Value
	Scope    token.Token // CONST, VAR, LET
	ReadOnly bool        // 是否为只读（const 变量）
}

type Environment struct {
	types map[string]object.Type
	// 分别管理不同类型的变量
	consts map[string]*VariableInfo // const 变量
	vars   map[string]*VariableInfo // var 变量
	lets   map[string]*VariableInfo // let 变量
	outer  *Environment
}

func NewEnvironment() *Environment {
	return &Environment{
		types:  make(map[string]object.Type),
		consts: make(map[string]*VariableInfo),
		vars:   make(map[string]*VariableInfo),
		lets:   make(map[string]*VariableInfo),
		outer:  nil,
	}
}

func (e *Environment) Fork() *Environment {
	env := NewEnvironment()
	env.outer = e
	return env
}

// Get 获取变量值，按优先级查找：const -> let -> var
func (e *Environment) Get(name string) (object.Value, bool) {
	// 先查找 const
	if info, ok := e.consts[name]; ok {
		return info.Value, true
	}
	// 再查找 let
	if info, ok := e.lets[name]; ok {
		return info.Value, true
	}
	// 最后查找 var
	if info, ok := e.vars[name]; ok {
		return info.Value, true
	}
	// 如果当前环境没有，查找外层环境
	if e.outer != nil {
		return e.outer.Get(name)
	}
	return nil, false
}

// GetVariableInfo 获取变量的完整信息
func (e *Environment) GetVariableInfo(name string) (*VariableInfo, bool) {
	if info, ok := e.consts[name]; ok {
		return info, true
	}
	if info, ok := e.lets[name]; ok {
		return info, true
	}
	if info, ok := e.vars[name]; ok {
		return info, true
	}
	if e.outer != nil {
		return e.outer.GetVariableInfo(name)
	}
	return nil, false
}

// Set 设置变量值（兼容旧接口）
func (e *Environment) Set(name string, obj object.Value) object.Value {
	// 默认使用 var 作用域
	return e.SetWithScope(name, obj, token.LET, false)
}

// SetWithScope 根据作用域设置变量
func (e *Environment) SetWithScope(name string, obj object.Value, scope token.Token, readOnly bool) object.Value {
	info := &VariableInfo{
		Value:    obj,
		Scope:    scope,
		ReadOnly: readOnly,
	}

	switch scope {
	case token.CONST:
		// 检查是否已存在
		if existing, exists := e.consts[name]; exists {
			if existing.ReadOnly {
				panic(fmt.Sprintf("cannot reassign to const variable '%s'", name))
			}
		}
		e.consts[name] = info
	case token.LET:
		// 检查是否已存在
		if existing, exists := e.lets[name]; exists {
			if existing.ReadOnly {
				panic(fmt.Sprintf("cannot reassign to let variable '%s'", name))
			}
		}
		e.lets[name] = info
	case token.VAR:
		// var 变量可以重新赋值
		e.vars[name] = info
	default:
		panic(fmt.Sprintf("unknown scope token: %v", scope))
	}

	return obj
}

// Assign 重新赋值变量
func (e *Environment) Assign(name string, obj object.Value) error {
	// 查找变量信息
	if info, ok := e.consts[name]; ok {
		if info.ReadOnly {
			return fmt.Errorf("cannot reassign to const variable '%s'", name)
		}
		info.Value = obj
		return nil
	}

	if info, ok := e.lets[name]; ok {
		if info.ReadOnly {
			return fmt.Errorf("cannot reassign to let variable '%s'", name)
		}
		info.Value = obj
		return nil
	}

	if info, ok := e.vars[name]; ok {
		info.Value = obj
		return nil
	}

	// 如果当前环境没有，尝试在外层环境赋值
	if e.outer != nil {
		return e.outer.Assign(name, obj)
	}

	return fmt.Errorf("variable '%s' not found", name)
}

// Declare 声明新变量
func (e *Environment) Declare(name string, obj object.Value, scope token.Token) error {
	// 检查是否已存在
	if _, exists := e.GetVariableInfo(name); exists {
		return fmt.Errorf("variable '%s' already declared", name)
	}

	readOnly := (scope == token.CONST)
	e.SetWithScope(name, obj, scope, readOnly)
	return nil
}

// IsDeclared 检查变量是否已声明
func (e *Environment) IsDeclared(name string) bool {
	_, ok := e.GetVariableInfo(name)
	return ok
}

// GetScope 获取变量的作用域
func (e *Environment) GetScope(name string) (token.Token, bool) {
	if info, ok := e.GetVariableInfo(name); ok {
		return info.Scope, true
	}
	return token.ILLEGAL, false
}

// IsReadOnly 检查变量是否为只读
func (e *Environment) IsReadOnly(name string) bool {
	if info, ok := e.GetVariableInfo(name); ok {
		return info.ReadOnly
	}
	return false
}
