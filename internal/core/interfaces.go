// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package core

import (
	"github.com/hulo-lang/hulo/syntax/hulo/ast"
)

// Type 类型接口
type Type interface {
	Name() string
	Kind() TypeKind
	AssignableTo(other Type) bool
	Equals(other Type) bool
}

// Value 值接口
type Value interface {
	Type() Type
	Interface() any
	Equals(other Value) bool
}

// TypeKind 类型种类
type TypeKind int

const (
	TypePrimitive TypeKind = iota
	TypeArray
	TypeMap
	TypeSet
	TypeTuple
	TypeUnion
	TypeIntersection
	TypeNullable
	TypeFunction
	TypeGeneric
	TypeAny
	TypeVoid
	TypeNull
	TypeError
)

// Module 模块接口
type Module interface {
	Name() string
	Path() string
	GetType(name string) (Type, bool)
	GetValue(name string) (Value, bool)
	Export(name string, value interface{}) error
	Import(name string, alias string, from string) error
}

// Registry 注册表接口
type Registry interface {
	RegisterModule(module Module) error
	GetModule(name string) (Module, bool)
	ResolveSymbol(moduleName, symbolName string) (any, bool)
	ResolveDependencies() error
}

// Converter AST转换器接口
type Converter interface {
	ConvertType(node ast.Node) (Type, error)
	ConvertValue(node ast.Node) (Value, error)
	ConvertModule(node *ast.File) (Module, error)
}

// TypeChecker 类型检查器接口
type TypeChecker interface {
	CheckType(expr ast.Expr) (Type, error)
	InferType(expr ast.Expr) (Type, error)
	CheckAssignment(target Type, value Value) error
}

// Evaluator 求值器接口
type Evaluator interface {
	Evaluate(node ast.Node) (Value, error)
	EvaluateInContext(node ast.Node, context *Context) (Value, error)
}

// FunctionEvaluator 函数求值器接口
type FunctionEvaluator interface {
	// EvaluateFunction 求值函数
	EvaluateFunction(funcDecl *ast.FuncDecl, args []any) (any, error)
}

// Context 执行上下文
type Context struct {
	Module    Module
	Registry  Registry
	Converter Converter
	Checker   TypeChecker
	Evaluator Evaluator
	Variables map[string]Value
	Parent    *Context
}
