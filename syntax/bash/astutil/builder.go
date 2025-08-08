// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package astutil

import (
	"github.com/hulo-lang/hulo/syntax/bash/ast"
	"github.com/hulo-lang/hulo/syntax/bash/token"
)

// FileBuilder 用于构建 bash 文件的 AST
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
func (b *FileBuilder) AddStmt(stmt ast.Stmt) *FileBuilder {
	b.stmts = append(b.stmts, stmt)
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
	return b.AddStmt(&ast.ExprStmt{
		X: CmdExpr("echo", args...),
	})
}

// Assign 添加赋值语句
func (b *FileBuilder) Assign(lhs, rhs ast.Expr) *FileBuilder {
	return b.AddStmt(&ast.AssignStmt{
		Lhs:    lhs,
		Assign: token.Pos(1),
		Rhs:    rhs,
	})
}

// LocalAssign 添加 local 赋值语句
func (b *FileBuilder) LocalAssign(lhs, rhs ast.Expr) *FileBuilder {
	return b.AddStmt(&ast.AssignStmt{
		Local:  token.Pos(1),
		Lhs:    lhs,
		Assign: token.Pos(1),
		Rhs:    rhs,
	})
}

// Return 添加 return 语句
func (b *FileBuilder) Return(x ast.Expr) *FileBuilder {
	return b.AddStmt(&ast.ReturnStmt{
		Return: token.Pos(1),
		X:      x,
	})
}

// If 开始构建 if 语句
func (b *FileBuilder) If(cond ast.Expr) *IfBuilder {
	return &IfBuilder{
		fileBuilder: b,
		root: &ast.IfStmt{
			If:   token.Pos(1),
			Cond: cond,
			Semi: token.Pos(1),
			Then: token.Pos(1),
		},
		current: nil,
	}
}

// While 开始构建 while 循环
func (b *FileBuilder) While(cond ast.Expr) *WhileBuilder {
	return &WhileBuilder{
		fileBuilder: b,
		cond:        cond,
	}
}

// Until 开始构建 until 循环
func (b *FileBuilder) Until(cond ast.Expr) *UntilBuilder {
	return &UntilBuilder{
		fileBuilder: b,
		cond:        cond,
	}
}

// ForIn 开始构建 for-in 循环
func (b *FileBuilder) ForIn(varName ast.Expr, list ast.Expr) *ForInBuilder {
	return &ForInBuilder{
		fileBuilder: b,
		varName:     varName,
		list:        list,
	}
}

// Function 开始构建函数声明
func (b *FileBuilder) Function(name string) *FunctionBuilder {
	return &FunctionBuilder{
		fileBuilder: b,
		name:        name,
	}
}

// ===== IfBuilder =====

type IfBuilder struct {
	fileBuilder *FileBuilder
	root        *ast.IfStmt // 根 if 语句
	current     *ast.IfStmt // 当前正在构建的 if 语句
}

func (b *IfBuilder) Then(stmts ...ast.Stmt) *IfBuilder {
	// 如果是第一次调用 Then，设置根节点的 body
	if b.current == nil {
		b.root.Body = &ast.BlockStmt{
			List: stmts,
		}
		b.current = b.root
	} else {
		// 否则设置当前节点的 body
		b.current.Body = &ast.BlockStmt{
			List: stmts,
		}
	}
	return b
}

func (b *IfBuilder) Else(stmts ...ast.Stmt) *IfBuilder {
	// 设置当前节点的 else 分支
	if b.current == nil {
		b.root.Else = &ast.BlockStmt{
			List: stmts,
		}
	} else {
		b.current.Else = &ast.BlockStmt{
			List: stmts,
		}
	}
	return b
}

func (b *IfBuilder) ElseIf(cond ast.Expr) *IfBuilder {
	// 创建新的 elif 节点
	elifStmt := &ast.IfStmt{
		If:   token.Pos(1),
		Cond: cond,
		Semi: token.Pos(1),
		Then: token.Pos(1),
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

// TODO 这个代码很有问题 if 嵌套的适合也会加到 file 里面
// EndIf 结束 if 语句构建，返回 FileBuilder 以支持链式调用
func (b *IfBuilder) EndIf() *FileBuilder {
	b.root.Fi = token.Pos(1)
	return b.fileBuilder.AddStmt(b.root)
}

// ===== WhileBuilder =====

type WhileBuilder struct {
	fileBuilder *FileBuilder
	cond        ast.Expr
	body        []ast.Stmt
}

func (b *WhileBuilder) Do(stmts ...ast.Stmt) *WhileBuilder {
	b.body = stmts
	return b
}

func (b *WhileBuilder) EndWhile() *FileBuilder {
	return b.fileBuilder.AddStmt(&ast.WhileStmt{
		While: token.Pos(1),
		Cond:  b.cond,
		Semi:  token.Pos(1),
		Do:    token.Pos(1),
		Body: &ast.BlockStmt{
			List: b.body,
		},
		Done: token.Pos(1),
	})
}

// ===== UntilBuilder =====

type UntilBuilder struct {
	fileBuilder *FileBuilder
	cond        ast.Expr
	body        []ast.Stmt
}

func (b *UntilBuilder) Do(stmts ...ast.Stmt) *UntilBuilder {
	b.body = stmts
	return b
}

func (b *UntilBuilder) EndUntil() *FileBuilder {
	return b.fileBuilder.AddStmt(&ast.UntilStmt{
		Until: token.Pos(1),
		Cond:  b.cond,
		Semi:  token.Pos(1),
		Do:    token.Pos(1),
		Body: &ast.BlockStmt{
			List: b.body,
		},
		Done: token.Pos(1),
	})
}

// ===== ForInBuilder =====

type ForInBuilder struct {
	fileBuilder *FileBuilder
	varName     ast.Expr
	list        ast.Expr
	body        []ast.Stmt
}

func (b *ForInBuilder) Do(stmts ...ast.Stmt) *ForInBuilder {
	b.body = stmts
	return b
}

func (b *ForInBuilder) EndFor() *FileBuilder {
	return b.fileBuilder.AddStmt(&ast.ForInStmt{
		For:  token.Pos(1),
		Var:  b.varName,
		In:   token.Pos(1),
		List: b.list,
		Semi: token.Pos(1),
		Do:   token.Pos(1),
		Body: &ast.BlockStmt{
			List: b.body,
		},
		Done: token.Pos(1),
	})
}

// ===== FunctionBuilder =====

type FunctionBuilder struct {
	fileBuilder *FileBuilder
	name        string
	body        []ast.Stmt
}

func (b *FunctionBuilder) Body(stmts ...ast.Stmt) *FunctionBuilder {
	b.body = stmts
	return b
}

func (b *FunctionBuilder) EndFunction() *FileBuilder {
	return b.fileBuilder.AddStmt(&ast.FuncDecl{
		Function: token.Pos(1),
		Name:     Ident(b.name),
		Lparen:   token.Pos(1),
		Rparen:   token.Pos(1),
		Body: &ast.BlockStmt{
			List: b.body,
		},
	})
}
