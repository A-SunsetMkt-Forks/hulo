// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package astutil

import (
	"testing"

	"github.com/hulo-lang/hulo/syntax/bash/ast"
	"github.com/hulo-lang/hulo/syntax/bash/token"
	"github.com/stretchr/testify/assert"
)

func TestFileBuilder_Basic(t *testing.T) {
	// 测试基本的文件构建
	file := NewFileBuilder().
		Echo(Word("Hello")).
		Assign(Ident("x"), Word("1")).
		Echo(Ident("x")).
		Build()

	assert.NotNil(t, file)
	assert.Len(t, file.Stmts, 3)

	// 验证第一个语句是 echo
	echoStmt, ok := file.Stmts[0].(*ast.ExprStmt)
	assert.True(t, ok)
	cmdExpr, ok := echoStmt.X.(*ast.CmdExpr)
	assert.True(t, ok)
	assert.Equal(t, "echo", cmdExpr.Name.Name)

	// 验证第二个语句是赋值
	assignStmt, ok := file.Stmts[1].(*ast.AssignStmt)
	assert.True(t, ok)
	assert.Equal(t, "x", assignStmt.Lhs.(*ast.Ident).Name)
}

func TestFileBuilder_Assign(t *testing.T) {
	// 测试赋值语句
	file := NewFileBuilder().
		Assign(Ident("var1"), Word("value1")).
		LocalAssign(Ident("var2"), Word("value2")).
		Build()

	assert.Len(t, file.Stmts, 2)

	// 验证普通赋值
	assignStmt, ok := file.Stmts[0].(*ast.AssignStmt)
	assert.True(t, ok)
	assert.Equal(t, "var1", assignStmt.Lhs.(*ast.Ident).Name)
	assert.Equal(t, "value1", assignStmt.Rhs.(*ast.Word).Val)
	assert.False(t, assignStmt.Local.IsValid())

	// 验证 local 赋值
	localStmt, ok := file.Stmts[1].(*ast.AssignStmt)
	assert.True(t, ok)
	assert.Equal(t, "var2", localStmt.Lhs.(*ast.Ident).Name)
	assert.Equal(t, "value2", localStmt.Rhs.(*ast.Word).Val)
	assert.True(t, localStmt.Local.IsValid())
}

func TestFileBuilder_Return(t *testing.T) {
	// 测试 return 语句
	file := NewFileBuilder().
		Return(Word("0")).
		Return(Ident("exit_code")).
		Build()

	assert.Len(t, file.Stmts, 2)

	// 验证第一个 return
	returnStmt1, ok := file.Stmts[0].(*ast.ReturnStmt)
	assert.True(t, ok)
	assert.Equal(t, "0", returnStmt1.X.(*ast.Word).Val)

	// 验证第二个 return
	returnStmt2, ok := file.Stmts[1].(*ast.ReturnStmt)
	assert.True(t, ok)
	assert.Equal(t, "exit_code", returnStmt2.X.(*ast.Ident).Name)
}

func TestIfBuilder_Basic(t *testing.T) {
	// 测试基本的 if-then-else
	cond := Eq(Ident("x"), Word("1"))
	thenStmt := CmdStmt("echo", Word("true"))
	elseStmt := CmdStmt("echo", Word("false"))

	file := NewFileBuilder().
		If(cond).
		Then(thenStmt).
		Else(elseStmt).
		EndIf().
		Build()

	assert.NotNil(t, file)
	assert.Len(t, file.Stmts, 1)

	ifStmt, ok := file.Stmts[0].(*ast.IfStmt)
	assert.True(t, ok)
	assert.Equal(t, cond, ifStmt.Cond)
	assert.NotNil(t, ifStmt.Body)
	assert.NotNil(t, ifStmt.Else)

	// 验证 body 是 BlockStmt
	assert.Len(t, ifStmt.Body.List, 1)
	assert.Equal(t, thenStmt, ifStmt.Body.List[0])

	// 验证 else 是 BlockStmt
	elseBlock, ok := ifStmt.Else.(*ast.BlockStmt)
	assert.True(t, ok)
	assert.Len(t, elseBlock.List, 1)
	assert.Equal(t, elseStmt, elseBlock.List[0])
}

func TestIfBuilder_ElseIf(t *testing.T) {
	// 测试 if-elif-else 链
	cond1 := Eq(Ident("x"), Word("1"))
	cond2 := Eq(Ident("x"), Word("2"))
	stmt1 := CmdStmt("echo", Word("one"))
	stmt2 := CmdStmt("echo", Word("two"))
	elseStmt := CmdStmt("echo", Word("other"))

	file := NewFileBuilder().
		If(cond1).
		Then(stmt1).
		ElseIf(cond2).
		Then(stmt2).
		Else(elseStmt).
		EndIf().
		Build()

	assert.NotNil(t, file)
	assert.Len(t, file.Stmts, 1)

	ifStmt, ok := file.Stmts[0].(*ast.IfStmt)
	assert.True(t, ok)
	assert.Equal(t, cond1, ifStmt.Cond)

	// 验证根节点的 body
	assert.Len(t, ifStmt.Body.List, 1)
	assert.Equal(t, stmt1, ifStmt.Body.List[0])

	// 验证 else 是 IfStmt (elif)
	elif, ok := ifStmt.Else.(*ast.IfStmt)
	assert.True(t, ok)
	assert.Equal(t, cond2, elif.Cond)

	// 验证 elif 的 body
	assert.Len(t, elif.Body.List, 1)
	assert.Equal(t, stmt2, elif.Body.List[0])

	// 验证 elif 的 else
	elifElse, ok := elif.Else.(*ast.BlockStmt)
	assert.True(t, ok)
	assert.Len(t, elifElse.List, 1)
	assert.Equal(t, elseStmt, elifElse.List[0])
}

func TestIfBuilder_MultipleElseIf(t *testing.T) {
	// 测试多个 elif 的情况
	cond1 := Eq(Ident("x"), Word("1"))
	cond2 := Eq(Ident("x"), Word("2"))
	cond3 := Eq(Ident("x"), Word("3"))
	stmt1 := CmdStmt("echo", Word("one"))
	stmt2 := CmdStmt("echo", Word("two"))
	stmt3 := CmdStmt("echo", Word("three"))
	elseStmt := CmdStmt("echo", Word("other"))

	file := NewFileBuilder().
		If(cond1).
		Then(stmt1).
		ElseIf(cond2).
		Then(stmt2).
		ElseIf(cond3).
		Then(stmt3).
		Else(elseStmt).
		EndIf().
		Build()

	assert.NotNil(t, file)
	assert.Len(t, file.Stmts, 1)

	ifStmt, ok := file.Stmts[0].(*ast.IfStmt)
	assert.True(t, ok)
	assert.Equal(t, cond1, ifStmt.Cond)

	// 验证第一个 elif
	elif1, ok := ifStmt.Else.(*ast.IfStmt)
	assert.True(t, ok)
	assert.Equal(t, cond2, elif1.Cond)

	// 验证第二个 elif
	elif2, ok := elif1.Else.(*ast.IfStmt)
	assert.True(t, ok)
	assert.Equal(t, cond3, elif2.Cond)

	// 验证最终的 else
	elseBlock, ok := elif2.Else.(*ast.BlockStmt)
	assert.True(t, ok)
	assert.Len(t, elseBlock.List, 1)
	assert.Equal(t, elseStmt, elseBlock.List[0])
}

func TestIfBuilder_IfOnly(t *testing.T) {
	// 测试只有 if 没有 else 的情况
	cond := Eq(Ident("x"), Word("1"))
	stmt := CmdStmt("echo", Word("true"))

	file := NewFileBuilder().
		If(cond).
		Then(stmt).
		EndIf().
		Build()

	assert.NotNil(t, file)
	assert.Len(t, file.Stmts, 1)

	ifStmt, ok := file.Stmts[0].(*ast.IfStmt)
	assert.True(t, ok)
	assert.Equal(t, cond, ifStmt.Cond)
	assert.NotNil(t, ifStmt.Body)
	assert.Nil(t, ifStmt.Else)

	assert.Len(t, ifStmt.Body.List, 1)
	assert.Equal(t, stmt, ifStmt.Body.List[0])
}

func TestWhileBuilder(t *testing.T) {
	// 测试 while 循环
	cond := Lt(Ident("count"), Word("5"))
	body := []ast.Stmt{
		CmdStmt("echo", Word("Count:"), Ident("count")),
		AssignStmt(Ident("count"), Binary(Ident("count"), token.Plus, Word("1"))),
	}

	file := NewFileBuilder().
		While(cond).
		Do(body...).
		EndWhile().
		Build()

	assert.Len(t, file.Stmts, 1)

	whileStmt, ok := file.Stmts[0].(*ast.WhileStmt)
	assert.True(t, ok)
	assert.Equal(t, cond, whileStmt.Cond)
	assert.Len(t, whileStmt.Body.List, 2)
	assert.Equal(t, body[0], whileStmt.Body.List[0])
	assert.Equal(t, body[1], whileStmt.Body.List[1])
}

func TestUntilBuilder(t *testing.T) {
	// 测试 until 循环
	cond := Eq(Ident("flag"), Word("true"))
	body := []ast.Stmt{
		CmdStmt("echo", Word("Waiting...")),
		CmdStmt("sleep", Word("1")),
	}

	file := NewFileBuilder().
		Until(cond).
		Do(body...).
		EndUntil().
		Build()

	assert.Len(t, file.Stmts, 1)

	untilStmt, ok := file.Stmts[0].(*ast.UntilStmt)
	assert.True(t, ok)
	assert.Equal(t, cond, untilStmt.Cond)
	assert.Len(t, untilStmt.Body.List, 2)
	assert.Equal(t, body[0], untilStmt.Body.List[0])
	assert.Equal(t, body[1], untilStmt.Body.List[1])
}

func TestForInBuilder(t *testing.T) {
	// 测试 for-in 循环
	varName := Ident("item")
	list := Word("apple banana cherry")
	body := []ast.Stmt{
		CmdStmt("echo", Word("Processing:"), Ident("item")),
	}

	file := NewFileBuilder().
		ForIn(varName, list).
		Do(body...).
		EndFor().
		Build()

	assert.Len(t, file.Stmts, 1)

	forStmt, ok := file.Stmts[0].(*ast.ForInStmt)
	assert.True(t, ok)
	assert.Equal(t, varName, forStmt.Var)
	assert.Equal(t, list, forStmt.List)
	assert.Len(t, forStmt.Body.List, 1)
	assert.Equal(t, body[0], forStmt.Body.List[0])
}

func TestFunctionBuilder(t *testing.T) {
	// 测试函数定义
	name := "greet"
	body := []ast.Stmt{
		CmdStmt("echo", Word("Hello,"), Ident("1")),
		ReturnStmt(Word("0")),
	}

	file := NewFileBuilder().
		Function(name).
		Body(body...).
		EndFunction().
		Build()

	assert.Len(t, file.Stmts, 1)

	funcDecl, ok := file.Stmts[0].(*ast.FuncDecl)
	assert.True(t, ok)
	assert.Equal(t, name, funcDecl.Name.Name)
	assert.Len(t, funcDecl.Body.List, 2)
	assert.Equal(t, body[0], funcDecl.Body.List[0])
	assert.Equal(t, body[1], funcDecl.Body.List[1])
}

func TestComplexNestedStructure(t *testing.T) {
	// 测试复杂的嵌套结构
	file := NewFileBuilder().
		// 定义变量
		Assign(Ident("max_count"), Word("10")).
		Assign(Ident("current"), Word("0")).

		// 函数定义
		Function("increment").
		Body(
			AssignStmt(Ident("current"), Binary(Ident("current"), token.Plus, Word("1"))),
			CmdStmt("echo", Word("Current:"), Ident("current")),
		).
		EndFunction().

		// while 循环
		While(Lt(Ident("current"), Ident("max_count"))).
		Do(
			// if 语句
			IfStmt(
				Eq(Ident("current"), Word("5")),
				Block(
					CmdStmt("echo", Word("Halfway point!")),
				),
				Block(
					CmdStmt("echo", Word("Continuing...")),
				),
			),
			CmdStmt("increment"),
		).
		EndWhile().

		// 最终输出
		Echo(Word("Loop completed")).
		Build()

	assert.NotNil(t, file)
	assert.Len(t, file.Stmts, 4) // 2 assignments + 1 function + 1 while + 1 echo

	// 验证函数定义
	funcDecl, ok := file.Stmts[2].(*ast.FuncDecl)
	assert.True(t, ok)
	assert.Equal(t, "increment", funcDecl.Name.Name)

	// 验证 while 循环
	whileStmt, ok := file.Stmts[3].(*ast.WhileStmt)
	assert.True(t, ok)
	assert.Len(t, whileStmt.Body.List, 2) // if statement + increment call
}

func TestBuilderChaining(t *testing.T) {
	// 测试构建器链式调用
	file := NewFileBuilder().
		Echo(Word("Start")).
		If(Eq(Ident("debug"), Word("true"))).
		Then(
			CmdStmt("echo", Word("Debug mode enabled")),
		).
		EndIf().
		Echo(Word("End")).
		Build()

	assert.Len(t, file.Stmts, 3) // echo + if + echo

	// 验证第一个 echo
	echo1, ok := file.Stmts[0].(*ast.ExprStmt)
	assert.True(t, ok)
	cmd1, ok := echo1.X.(*ast.CmdExpr)
	assert.True(t, ok)
	assert.Equal(t, "echo", cmd1.Name.Name)

	// 验证 if 语句
	ifStmt, ok := file.Stmts[1].(*ast.IfStmt)
	assert.True(t, ok)
	assert.NotNil(t, ifStmt.Cond)

	// 验证最后一个 echo
	echo2, ok := file.Stmts[2].(*ast.ExprStmt)
	assert.True(t, ok)
	cmd2, ok := echo2.X.(*ast.CmdExpr)
	assert.True(t, ok)
	assert.Equal(t, "echo", cmd2.Name.Name)
}

func TestNestedIfChaining(t *testing.T) {
	// 测试嵌套的 if 语句链式调用
	file := NewFileBuilder().
		Echo(Word("Starting nested test")).
		If(Eq(Ident("outer"), Word("true"))).
		Then(
			CmdStmt("echo", Word("Outer condition is true")),
			IfStmt(
				Eq(Ident("inner"), Word("true")),
				Block(
					CmdStmt("echo", Word("Inner condition is also true")),
				),
				Block(
					CmdStmt("echo", Word("Inner condition is false")),
				),
			),
		).
		Else(
			CmdStmt("echo", Word("Outer condition is false")),
		).
		EndIf().
		Echo(Word("Nested test completed")).
		Build()

	assert.Len(t, file.Stmts, 3) // echo + if + echo

	// 验证 if 语句
	ifStmt, ok := file.Stmts[1].(*ast.IfStmt)
	assert.True(t, ok)
	assert.Len(t, ifStmt.Body.List, 2) // echo + nested if
}

func TestMultipleControlStructures(t *testing.T) {
	// 测试多个控制结构的链式调用
	file := NewFileBuilder().
		Assign(Ident("i"), Word("0")).
		While(Lt(Ident("i"), Word("3"))).
		Do(
			CmdStmt("echo", Word("While iteration:"), Ident("i")),
			AssignStmt(Ident("i"), Binary(Ident("i"), token.Plus, Word("1"))),
		).
		EndWhile().
		ForIn(Ident("item"), Word("a b c")).
		Do(
			CmdStmt("echo", Word("For item:"), Ident("item")),
		).
		EndFor().
		Function("test_func").
		Body(
			CmdStmt("echo", Word("Function called")),
		).
		EndFunction().
		Build()

	assert.Len(t, file.Stmts, 4) // assign + while + for + function

	// 验证 while 循环
	whileStmt, ok := file.Stmts[1].(*ast.WhileStmt)
	assert.True(t, ok)
	assert.Equal(t, "i", whileStmt.Cond.(*ast.BinaryExpr).X.(*ast.Ident).Name)

	// 验证 for-in 循环
	forStmt, ok := file.Stmts[2].(*ast.ForInStmt)
	assert.True(t, ok)
	assert.Equal(t, "item", forStmt.Var.(*ast.Ident).Name)

	// 验证函数定义
	funcDecl, ok := file.Stmts[3].(*ast.FuncDecl)
	assert.True(t, ok)
	assert.Equal(t, "test_func", funcDecl.Name.Name)
}
