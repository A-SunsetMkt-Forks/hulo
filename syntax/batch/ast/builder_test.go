// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package ast

import (
	"testing"

	"github.com/hulo-lang/hulo/syntax/batch/token"
	"github.com/stretchr/testify/assert"
)

type FileBuilder struct {
	stmts []Stmt
}

func NewFile() *FileBuilder {
	return &FileBuilder{
		stmts: make([]Stmt, 0),
	}
}

func (b *FileBuilder) AddStmt(stmt Stmt) *FileBuilder {
	b.stmts = append(b.stmts, stmt)
	return b
}

func (b *FileBuilder) Build() (*File, error) {
	return &File{
		Stmts: b.stmts,
	}, nil
}

// ===== 基础操作 =====

// EchoOff 添加 @echo off
func (b *FileBuilder) EchoOff() *FileBuilder {
	return b.AddStmt(&ExprStmt{
		X: Words(Literal("@echo"), Literal("off")),
	})
}

// SetLocal 添加 SETLOCAL
func (b *FileBuilder) SetLocal() *FileBuilder {
	return b.AddStmt(&ExprStmt{
		X: CmdExpression("SETLOCAL"),
	})
}

// Echo 添加 ECHO 命令
func (b *FileBuilder) Echo(args ...Expr) *FileBuilder {
	return b.AddStmt(&ExprStmt{
		X: &CmdExpr{
			Name: &Lit{Val: "echo"},
			Recv: args,
		},
	})
}

// Exit 添加 EXIT 命令
func (b *FileBuilder) Exit(args ...Expr) *FileBuilder {
	return b.AddStmt(&ExprStmt{
		X: &CallExpr{
			Fun:  Literal("EXIT"),
			Recv: args,
		},
	})
}

// Call 添加函数调用
func (b *FileBuilder) Call(name string, args ...Expr) *FileBuilder {
	return b.AddStmt(&CallStmt{
		Name: name,
		Recv: args,
	})
}

// Goto 添加 GOTO 语句
func (b *FileBuilder) Goto(label string) *FileBuilder {
	return b.AddStmt(&GotoStmt{Label: label})
}

// Label 添加标签
func (b *FileBuilder) Label(name string) *FileBuilder {
	return b.AddStmt(&LabelStmt{Name: name})
}

// Assign 添加赋值语句
func (b *FileBuilder) Assign(lhs, rhs Expr) *FileBuilder {
	return b.AddStmt(&AssignStmt{
		Lhs: lhs,
		Rhs: rhs,
	})
}

// ===== 控制流操作 =====

// If 添加 IF 语句
func (b *FileBuilder) If(cond Expr) *IfBuilder {
	return &IfBuilder{
		fileBuilder: b,
		cond:        cond,
	}
}

// For 添加 FOR 循环
func (b *FileBuilder) For(varName string, list Expr) *ForBuilder {
	return &ForBuilder{
		fileBuilder: b,
		varName:     varName,
		list:        list,
	}
}

// Function 添加函数定义
func (b *FileBuilder) Function(name string) *FunctionBuilder {
	return &FunctionBuilder{
		fileBuilder: b,
		name:        name,
	}
}

// ===== 便捷表达式函数 =====

// Str 创建双引号字符串
func Str(val string) *DblQuote {
	return &DblQuote{Val: &Lit{Val: val}}
}

// Quote 创建单引号字符串
func Quote(val string) *SglQuote {
	return &SglQuote{Val: &Lit{Val: val}}
}

// Var 创建变量引用
func Var(name string) *DblQuote {
	return &DblQuote{Val: &Lit{Val: name}}
}

// Cmd 创建命令表达式
func Cmd(name string, args ...Expr) *CmdExpr {
	return &CmdExpr{
		Name: &Lit{Val: name},
		Recv: args,
	}
}

// Eq 创建相等比较
func Eq(x, y Expr) *BinaryExpr {
	return &BinaryExpr{
		X:  x,
		Op: token.DOUBLE_ASSIGN,
		Y:  y,
	}
}

// Exit 创建 EXIT 调用
func Exit(args ...Expr) *CallExpr {
	return &CallExpr{
		Fun:  Literal("EXIT"),
		Recv: args,
	}
}

// ===== 辅助函数 =====

// EchoStmt 创建 ECHO 语句
func EchoStmt(args ...Expr) *ExprStmt {
	return &ExprStmt{
		X: &CmdExpr{
			Name: &Lit{Val: "ECHO"},
			Recv: args,
		},
	}
}

// ExitStmt 创建 EXIT 语句
func ExitStmt(args ...Expr) *ExprStmt {
	return &ExprStmt{
		X: &CallExpr{
			Fun:  Literal("EXIT"),
			Recv: args,
		},
	}
}

// GotoStmt 创建 GOTO 语句
func Goto(label string) *GotoStmt {
	return &GotoStmt{Label: label}
}

// ===== IfBuilder =====

type IfBuilder struct {
	fileBuilder *FileBuilder
	cond        Expr
	body        []Stmt
	elseBody    []Stmt
}

// Then 设置 IF 的 then 分支
func (ib *IfBuilder) Then(stmts ...Stmt) *IfBuilder {
	ib.body = stmts
	return ib
}

// ThenIf 在 then 分支中添加嵌套的 IF
func (ib *IfBuilder) ThenIf(cond Expr) *IfBuilder {
	nestedIf := &IfStmt{
		Cond: cond,
		Body: &BlockStmt{List: []Stmt{}},
	}
	ib.body = append(ib.body, nestedIf)
	return ib
}

// Else 设置 IF 的 else 分支
func (ib *IfBuilder) Else(stmts ...Stmt) *IfBuilder {
	ib.elseBody = stmts
	return ib
}

// EndIf 结束 IF 语句
func (ib *IfBuilder) EndIf() *FileBuilder {
	stmt := &IfStmt{
		Cond: ib.cond,
		Body: &BlockStmt{List: ib.body},
	}

	if len(ib.elseBody) > 0 {
		stmt.Else = &BlockStmt{List: ib.elseBody}
	}

	return ib.fileBuilder.AddStmt(stmt)
}

// ===== ForBuilder =====

type ForBuilder struct {
	fileBuilder *FileBuilder
	varName     string
	list        Expr
	body        []Stmt
}

// Do 设置 FOR 循环体
func (fb *ForBuilder) Do(stmts ...Stmt) *ForBuilder {
	fb.body = stmts
	return fb
}

// EndFor 结束 FOR 循环
func (fb *ForBuilder) EndFor() *FileBuilder {
	return fb.fileBuilder.AddStmt(&ForStmt{
		X:    &SglQuote{Val: &SglQuote{Val: &Lit{Val: fb.varName}}},
		List: fb.list,
		Body: &BlockStmt{List: fb.body},
	})
}

// ===== FunctionBuilder =====

type FunctionBuilder struct {
	fileBuilder *FileBuilder
	name        string
	body        []Stmt
}

// Body 设置函数体
func (fb *FunctionBuilder) Body(stmts ...Stmt) *FunctionBuilder {
	fb.body = stmts
	return fb
}

// EndFunction 结束函数定义
func (fb *FunctionBuilder) EndFunction() *FileBuilder {
	return fb.fileBuilder.AddStmt(&FuncDecl{
		Name: fb.name,
		Body: &BlockStmt{List: fb.body},
	})
}

// ===== 使用示例 =====

func TestBuilderUsage(t *testing.T) {
	// 基础使用
	file, err := NewFile().
		EchoOff().
		SetLocal().
		Echo(Literal("Hello World")).
		Exit(Literal("0")).
		Build()

	if err != nil {
		t.Fatal(err)
	}

	Print(file)
}

func TestIfStatement(t *testing.T) {
	file, err := NewFile().
		EchoOff().
		If(Eq(Str("VAR"), Literal("value"))).
		Then(
			&ExprStmt{X: Cmd("ECHO", Literal("Match"))},
		).
		Else(
			&ExprStmt{X: Cmd("ECHO", Literal("No Match"))},
		).
		EndIf().
		Exit(Literal("0")).
		Build()

	if err != nil {
		t.Fatal(err)
	}

	Print(file)
}

func TestForLoop(t *testing.T) {
	file, err := NewFile().
		EchoOff().
		For("file", Literal("(*.txt)")).
		Do(
			&ExprStmt{X: Cmd("ECHO", Quote("file"))},
		).
		EndFor().
		Exit(Literal("0")).
		Build()

	if err != nil {
		t.Fatal(err)
	}

	Print(file)
}

func TestFunction(t *testing.T) {
	file, err := NewFile().
		EchoOff().
		Call("myFunction", Literal("arg1"), Literal("arg2")).
		Goto("EOF").
		Function("myFunction").
		Body(
			&ExprStmt{X: Cmd("ECHO", Literal("Function called"))},
			&GotoStmt{Label: "EOF"},
		).
		EndFunction().
		Build()

	if err != nil {
		t.Fatal(err)
	}

	Print(file)
}

func TestComplexExample(t *testing.T) {
	file, err := NewFile().
		EchoOff().
		SetLocal().
		If(Eq(Str("DEBUG"), Literal("1"))).
		Then(
			&ExprStmt{X: Cmd("ECHO", Literal("Debug mode enabled"))},
		).
		EndIf().
		For("file", Literal("(*.txt *.log)")).
		Do(
			&ExprStmt{X: Cmd("ECHO", Literal("Processing"), Quote("file"))},
		).
		EndFor().
		Echo(Literal("Done")).
		Exit(Literal("0")).
		Build()

	if err != nil {
		t.Fatal(err)
	}

	Print(file)
}

// TestNestedIfExample 测试嵌套 IF 的正确用法
func TestNestedIfExample(t *testing.T) {
	file, err := NewFile().
		EchoOff().
		If(Eq(Str("VAR1"), Literal("value1"))).
		Then(
			&IfStmt{
				Cond: Eq(Str("VAR2"), Literal("value2")),
				Body: &BlockStmt{
					List: []Stmt{
						&ExprStmt{X: Cmd("ECHO", Literal("Both conditions match"))},
					},
				},
				Else: &BlockStmt{
					List: []Stmt{
						&ExprStmt{X: Cmd("ECHO", Literal("Only VAR1 matches"))},
					},
				},
			},
		).
		Else(
			&ExprStmt{X: Cmd("ECHO", Literal("No conditions match"))},
		).
		EndIf().
		Exit(Literal("0")).
		Build()

	assert.NoError(t, err)

	Print(file)
}

// testAST 测试辅助函数
func testAST(t *testing.T, name string, builder func() (*File, error)) {
	t.Run(name, func(t *testing.T) {
		file, err := builder()
		if err != nil {
			t.Fatalf("Failed to build AST: %v", err)
		}

		// 验证和打印
		Print(file)
	})
}

func TestWithHelper(t *testing.T) {
	testAST(t, "Basic IF", func() (*File, error) {
		return NewFile().
			EchoOff().
			If(Eq(Str("VAR"), Literal("value"))).
			Then(&ExprStmt{X: Cmd("ECHO", Literal("Match"))}).
			Else(&ExprStmt{X: Cmd("ECHO", Literal("No Match"))}).
			EndIf().
			Build()
	})
}
