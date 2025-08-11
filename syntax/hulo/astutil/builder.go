// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package astutil

import (
	"github.com/hulo-lang/hulo/syntax/hulo/ast"
	"github.com/hulo-lang/hulo/syntax/hulo/token"
)

type controlFlowBuilder struct{}

type declarationBuilder struct{}

type expressionBuilder struct{}

type ifBuilder struct {
	controlFlowBuilder
	parent any
	node   *ast.IfStmt
}

func NewIfBuilder(cond ast.Expr) *ifBuilder {
	return &ifBuilder{node: &ast.IfStmt{Cond: cond}}
}

func (b *ifBuilder) Add(node any) *ifBuilder {
	if b.node.Body == nil {
		b.node.Body = &ast.BlockStmt{}
	}
	if stmt, ok := node.(ast.Stmt); ok {
		b.node.Body.List = append(b.node.Body.List, stmt)
	}
	return b
}

func (b *ifBuilder) End() any {
	return b.parent
}

// StatementBuilder 是一个接口，用于支持嵌套的语句构建
type StatementBuilder interface {
	// AddStmt 添加一个语句到当前构建器
	AddStmt(stmt ast.Stmt) StatementBuilder
	// End 结束当前构建，返回父构建器
	End() StatementBuilder
}

// FileBuilder 用于构建 hulo 文件的 AST
type FileBuilder struct {
	stmts []ast.Stmt
}

// NewFileBuilder 创建一个新的 FileBuilder
func NewFileBuilder() *FileBuilder {
	return &FileBuilder{
		stmts: make([]ast.Stmt, 0),
	}
}

// AddStmt 添加一个语句到文件中
func (b *FileBuilder) AddStmt(stmt ast.Stmt) StatementBuilder {
	b.stmts = append(b.stmts, stmt)
	return b
}

func (b *FileBuilder) End() StatementBuilder {
	return b
}

// Build 构建并返回完整的文件 AST
func (b *FileBuilder) Build() *ast.File {
	return &ast.File{
		Stmts: b.stmts,
	}
}

// Echo 添加 echo 命令
func (b *FileBuilder) Echo(args ...ast.Expr) *FileBuilder {
	b.AddStmt(&ast.ExprStmt{
		X: CmdExpr("echo", args...),
	})
	return b
}

// Assign 添加赋值语句
func (b *FileBuilder) Assign(lhs, rhs ast.Expr) *FileBuilder {
	b.AddStmt(&ast.AssignStmt{
		Lhs: lhs,
		Tok: token.ASSIGN,
		Rhs: rhs,
	})
	return b
}

// VarAssign 添加 var 赋值语句
func (b *FileBuilder) VarAssign(lhs, rhs ast.Expr) *FileBuilder {
	b.AddStmt(&ast.AssignStmt{
		Scope: token.VAR,
		Lhs:   lhs,
		Tok:   token.ASSIGN,
		Rhs:   rhs,
	})
	return b
}

// ConstAssign 添加 const 赋值语句
func (b *FileBuilder) ConstAssign(lhs, rhs ast.Expr) *FileBuilder {
	b.AddStmt(&ast.AssignStmt{
		Scope: token.CONST,
		Lhs:   lhs,
		Tok:   token.ASSIGN,
		Rhs:   rhs,
	})
	return b
}

// Return 添加 return 语句
func (b *FileBuilder) Return(x ast.Expr) *FileBuilder {
	b.AddStmt(&ast.ReturnStmt{
		Return: token.Pos(1),
		X:      x,
	})
	return b
}

// If 开始构建 if 语句
func (b *FileBuilder) If(cond ast.Expr) *IfBuilder {
	return &IfBuilder{
		parent: b,
		root: &ast.IfStmt{
			If:   token.Pos(1),
			Cond: cond,
		},
		current: nil,
	}
}

// While 开始构建 while 循环
func (b *FileBuilder) While(cond ast.Expr) *WhileBuilder {
	return &WhileBuilder{
		parent: b,
		cond:   cond,
	}
}

// Loop 开始构建 loop 循环
func (b *FileBuilder) Loop() *LoopBuilder {
	return &LoopBuilder{
		parent: b,
	}
}

// ForIn 开始构建 for-in 循环
func (b *FileBuilder) ForIn(index *ast.Ident, var_ ast.Expr) *ForInBuilder {
	return &ForInBuilder{
		parent: b,
		index:  index,
		var_:   var_,
	}
}

// Function 开始构建函数声明
func (b *FileBuilder) Function(name string) *FunctionBuilder {
	return &FunctionBuilder{
		parent: b,
		name:   name,
	}
}

// Class 开始构建类声明
func (b *FileBuilder) Class(name string) *ClassBuilder {
	return &ClassBuilder{
		parent: b,
		name:   name,
	}
}

// ===== IfBuilder =====

type IfBuilder struct {
	parent  StatementBuilder
	root    *ast.IfStmt // 根 if 语句
	current *ast.IfStmt // 当前正在构建的 if 语句
	body    []ast.Stmt  // 当前分支的语句列表
}

func (b *IfBuilder) AddStmt(stmt ast.Stmt) StatementBuilder {
	b.body = append(b.body, stmt)
	return b
}

func (b *IfBuilder) Then(stmts ...ast.Stmt) *IfBuilder {
	// 如果是第一次调用 Then，设置根节点的 body
	if b.current == nil {
		b.root.Body = &ast.BlockStmt{
			Lbrace: token.Pos(1),
			List:   stmts,
			Rbrace: token.Pos(1),
		}
		b.current = b.root
	} else {
		// 否则设置当前节点的 body
		b.current.Body = &ast.BlockStmt{
			Lbrace: token.Pos(1),
			List:   stmts,
			Rbrace: token.Pos(1),
		}
	}
	return b
}

func (b *IfBuilder) Else(stmts ...ast.Stmt) *IfBuilder {
	// 设置当前节点的 else 分支
	if b.current == nil {
		b.root.Else = &ast.BlockStmt{
			Lbrace: token.Pos(1),
			List:   stmts,
			Rbrace: token.Pos(1),
		}
	} else {
		b.current.Else = &ast.BlockStmt{
			Lbrace: token.Pos(1),
			List:   stmts,
			Rbrace: token.Pos(1),
		}
	}
	return b
}

func (b *IfBuilder) ElseIf(cond ast.Expr) *IfBuilder {
	// 创建新的 elif 节点
	elifStmt := &ast.IfStmt{
		If:   token.Pos(1),
		Cond: cond,
	}

	// 如果是第一次调用 ElseIf，设置根节点的 else
	if b.current == nil {
		b.root.Else = elifStmt
	} else {
		// 否则设置当前节点的 else
		b.current.Else = elifStmt
	}

	// 更新当前节点为新的 elif 节点
	b.current = elifStmt
	return b
}

// EndIf 结束 if 语句构建，返回父构建器
func (b *IfBuilder) End() StatementBuilder {
	// 将完整的 if 语句添加到父构建器
	if fileBuilder, ok := b.parent.(*FileBuilder); ok {
		fileBuilder.AddStmt(b.root)
		return fileBuilder
	}
	// 如果是嵌套在其他构建器中，添加到当前构建器
	b.parent.AddStmt(b.root)
	return b.parent
}

// ===== WhileBuilder =====

type WhileBuilder struct {
	parent StatementBuilder
	cond   ast.Expr
	body   []ast.Stmt
}

func (b *WhileBuilder) AddStmt(stmt ast.Stmt) StatementBuilder {
	b.body = append(b.body, stmt)
	return b
}

func (b *WhileBuilder) Do(stmts ...ast.Stmt) *WhileBuilder {
	b.body = stmts
	return b
}

func (b *WhileBuilder) End() StatementBuilder {
	whileStmt := &ast.WhileStmt{
		Loop: token.Pos(1),
		Cond: b.cond,
		Body: &ast.BlockStmt{
			Lbrace: token.Pos(1),
			List:   b.body,
			Rbrace: token.Pos(1),
		},
	}

	if fileBuilder, ok := b.parent.(*FileBuilder); ok {
		fileBuilder.AddStmt(whileStmt)
		return fileBuilder
	}
	b.parent.AddStmt(whileStmt)
	return b.parent
}

// ===== LoopBuilder =====

type LoopBuilder struct {
	parent StatementBuilder
	body   []ast.Stmt
}

func (b *LoopBuilder) AddStmt(stmt ast.Stmt) StatementBuilder {
	b.body = append(b.body, stmt)
	return b
}

func (b *LoopBuilder) Do(stmts ...ast.Stmt) *LoopBuilder {
	b.body = stmts
	return b
}

func (b *LoopBuilder) End() StatementBuilder {
	// 这里需要根据具体的loop类型来决定
	// 暂时使用WhileStmt作为示例
	loopStmt := &ast.WhileStmt{
		Loop: token.Pos(1),
		Cond: &ast.TrueLiteral{ValuePos: token.Pos(1)}, // 无限循环
		Body: &ast.BlockStmt{
			Lbrace: token.Pos(1),
			List:   b.body,
			Rbrace: token.Pos(1),
		},
	}

	if fileBuilder, ok := b.parent.(*FileBuilder); ok {
		fileBuilder.AddStmt(loopStmt)
		return fileBuilder
	}
	b.parent.AddStmt(loopStmt)
	return b.parent
}

// ===== ForInBuilder =====

type ForInBuilder struct {
	parent StatementBuilder
	index  *ast.Ident
	var_   ast.Expr
	body   []ast.Stmt
}

func (b *ForInBuilder) AddStmt(stmt ast.Stmt) StatementBuilder {
	b.body = append(b.body, stmt)
	return b
}

func (b *ForInBuilder) Do(stmts ...ast.Stmt) *ForInBuilder {
	b.body = stmts
	return b
}

func (b *ForInBuilder) End() StatementBuilder {
	forInStmt := &ast.ForInStmt{
		Loop:  token.Pos(1),
		Index: b.index,
		In:    token.Pos(1),
		RangeExpr: ast.RangeExpr{
			Start: b.var_,
		},
		Body: &ast.BlockStmt{
			Lbrace: token.Pos(1),
			List:   b.body,
			Rbrace: token.Pos(1),
		},
	}

	if fileBuilder, ok := b.parent.(*FileBuilder); ok {
		fileBuilder.AddStmt(forInStmt)
		return fileBuilder
	}
	b.parent.AddStmt(forInStmt)
	return b.parent
}

// ===== FunctionBuilder =====

type FunctionBuilder struct {
	parent StatementBuilder
	name   string
	body   []ast.Stmt
	params []ast.Expr
}

func (b *FunctionBuilder) AddStmt(stmt ast.Stmt) StatementBuilder {
	b.body = append(b.body, stmt)
	return b
}

func (b *FunctionBuilder) Params(params ...ast.Expr) *FunctionBuilder {
	b.params = params
	return b
}

func (b *FunctionBuilder) Body(stmts ...ast.Stmt) *FunctionBuilder {
	b.body = stmts
	return b
}

func (b *FunctionBuilder) End() StatementBuilder {
	funcDecl := &ast.FuncDecl{
		Fn:   token.Pos(1),
		Name: &ast.Ident{NamePos: token.Pos(1), Name: b.name},
		Recv: b.params,
		Body: &ast.BlockStmt{
			Lbrace: token.Pos(1),
			List:   b.body,
			Rbrace: token.Pos(1),
		},
	}

	if fileBuilder, ok := b.parent.(*FileBuilder); ok {
		fileBuilder.AddStmt(funcDecl)
		return fileBuilder
	}
	b.parent.AddStmt(funcDecl)
	return b.parent
}

// ===== ClassBuilder =====

type ClassBuilder struct {
	parent  StatementBuilder
	name    string
	fields  []*ast.Field
	methods []ast.Stmt
}

func (b *ClassBuilder) AddStmt(stmt ast.Stmt) StatementBuilder {
	b.methods = append(b.methods, stmt)
	return b
}

func (b *ClassBuilder) Field(name string, fieldType ast.Expr) *ClassBuilder {
	field := &ast.Field{
		Name:  &ast.Ident{NamePos: token.Pos(1), Name: name},
		Colon: token.Pos(1),
		Type:  fieldType,
	}
	b.fields = append(b.fields, field)
	return b
}

func (b *ClassBuilder) Method(name string) *FunctionBuilder {
	return &FunctionBuilder{
		parent: b,
		name:   name,
	}
}

func (b *ClassBuilder) End() StatementBuilder {
	classDecl := &ast.ClassDecl{
		Class:  token.Pos(1),
		Name:   &ast.Ident{NamePos: token.Pos(1), Name: b.name},
		Lbrace: token.Pos(1),
		Fields: &ast.FieldList{
			List: b.fields,
		},
		Methods: make([]*ast.FuncDecl, 0),
		Rbrace:  token.Pos(1),
	}

	// 将方法转换为FuncDecl
	for _, stmt := range b.methods {
		if funcDecl, ok := stmt.(*ast.FuncDecl); ok {
			classDecl.Methods = append(classDecl.Methods, funcDecl)
		}
	}

	if fileBuilder, ok := b.parent.(*FileBuilder); ok {
		fileBuilder.AddStmt(classDecl)
		return fileBuilder
	}
	b.parent.AddStmt(classDecl)
	return b.parent
}

// ===== 辅助函数 =====

// BinaryExpr 创建二元表达式
func BinaryExpr(x ast.Expr, op token.Token, y ast.Expr) *ast.BinaryExpr {
	return &ast.BinaryExpr{
		X:  x,
		Op: op,
		Y:  y,
	}
}
