// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package astutil

import (
	"testing"

	"github.com/hulo-lang/hulo/syntax/hulo/ast"
	"github.com/hulo-lang/hulo/syntax/hulo/token"
	"github.com/stretchr/testify/assert"
)

func TestFileBuilder_Basic(t *testing.T) {
	// 测试基本的文件构建
	builder := NewFileBuilder()

	// 添加一些基本语句
	builder.Echo(Str("Hello, World!"))
	builder.Assign(Ident("x"), Num("42"))
	builder.Return(Ident("x"))

	file := builder.Build()

	assert.Len(t, file.Stmts, 3)

	// 验证第一个语句是echo
	echoStmt, ok := file.Stmts[0].(*ast.ExprStmt)
	assert.True(t, ok)
	cmdExpr, ok := echoStmt.X.(*ast.CmdExpr)
	assert.True(t, ok)
	assert.Equal(t, "echo", cmdExpr.Cmd.(*ast.Ident).Name)

	// 验证第二个语句是赋值
	assignStmt, ok := file.Stmts[1].(*ast.AssignStmt)
	assert.True(t, ok)
	assert.Equal(t, token.ASSIGN, assignStmt.Tok)

	// 验证第三个语句是return
	returnStmt, ok := file.Stmts[2].(*ast.ReturnStmt)
	assert.True(t, ok)
	assert.NotNil(t, returnStmt.X)
}

func TestFileBuilder_IfStatement(t *testing.T) {
	// 测试if语句构建
	builder := NewFileBuilder()

	// 构建if语句
	ifBuilder := builder.If(BinaryExpr(Ident("x"), token.EQ, Num("10")))
	ifBuilder.Then(
		&ast.ExprStmt{X: CmdExpr("echo", Str("x is 10"))},
		&ast.AssignStmt{Lhs: Ident("result"), Tok: token.ASSIGN, Rhs: Str("found")},
	)
	ifBuilder.Else(
		&ast.ExprStmt{X: CmdExpr("echo", Str("x is not 10"))},
	)

	// 结束if构建并返回文件构建器
	fileBuilder := ifBuilder.End().(*FileBuilder)
	file := fileBuilder.Build()

	assert.Len(t, file.Stmts, 1)

	// 验证if语句结构
	ifStmt, ok := file.Stmts[0].(*ast.IfStmt)
	assert.True(t, ok)
	assert.NotNil(t, ifStmt.Cond)
	assert.NotNil(t, ifStmt.Body)
	assert.NotNil(t, ifStmt.Else)

	// 验证条件表达式
	binaryExpr, ok := ifStmt.Cond.(*ast.BinaryExpr)
	assert.True(t, ok)
	assert.Equal(t, token.EQ, binaryExpr.Op)

	// 验证then分支
	assert.Len(t, ifStmt.Body.List, 2)

	// 验证else分支
	elseBlock, ok := ifStmt.Else.(*ast.BlockStmt)
	assert.True(t, ok)
	assert.Len(t, elseBlock.List, 1)
}

func TestFileBuilder_IfElseIf(t *testing.T) {
	// 测试if-elseif-else语句构建
	builder := NewFileBuilder()

	ifBuilder := builder.If(BinaryExpr(Ident("x"), token.EQ, Num("10")))
	ifBuilder.Then(
		&ast.ExprStmt{X: CmdExpr("echo", Str("x is 10"))},
	)
	ifBuilder.ElseIf(BinaryExpr(Ident("x"), token.GT, Num("10")))
	ifBuilder.Then(
		&ast.ExprStmt{X: CmdExpr("echo", Str("x is greater than 10"))},
	)
	ifBuilder.Else(
		&ast.ExprStmt{X: CmdExpr("echo", Str("x is less than 10"))},
	)

	fileBuilder := ifBuilder.End().(*FileBuilder)
	file := fileBuilder.Build()

	assert.Len(t, file.Stmts, 1)

	// 验证嵌套的if-else结构
	ifStmt, ok := file.Stmts[0].(*ast.IfStmt)
	assert.True(t, ok)

	// 验证else分支是另一个if语句
	elseIfStmt, ok := ifStmt.Else.(*ast.IfStmt)
	assert.True(t, ok)
	assert.NotNil(t, elseIfStmt.Cond)
	assert.NotNil(t, elseIfStmt.Body)
	assert.NotNil(t, elseIfStmt.Else)
}

func TestFileBuilder_WhileLoop(t *testing.T) {
	// 测试while循环构建
	builder := NewFileBuilder()

	whileBuilder := builder.While(BinaryExpr(Ident("i"), token.LT, Num("10")))
	whileBuilder.Do(
		&ast.ExprStmt{X: CmdExpr("echo", Ident("i"))},
		&ast.AssignStmt{Lhs: Ident("i"), Tok: token.PLUS_ASSIGN, Rhs: Num("1")},
	)

	fileBuilder := whileBuilder.End().(*FileBuilder)
	file := fileBuilder.Build()

	assert.Len(t, file.Stmts, 1)

	// 验证while语句结构
	whileStmt, ok := file.Stmts[0].(*ast.WhileStmt)
	assert.True(t, ok)
	assert.NotNil(t, whileStmt.Cond)
	assert.NotNil(t, whileStmt.Body)
	assert.Len(t, whileStmt.Body.List, 2)
}

func TestFileBuilder_ForInLoop(t *testing.T) {
	// 测试for-in循环构建
	builder := NewFileBuilder()

	forInBuilder := builder.ForIn(Ident("i"), Ident("items"))
	forInBuilder.Do(
		&ast.ExprStmt{X: CmdExpr("echo", Ident("i"))},
	)

	fileBuilder := forInBuilder.End().(*FileBuilder)
	file := fileBuilder.Build()

	assert.Len(t, file.Stmts, 1)

	// 验证for-in语句结构
	forInStmt, ok := file.Stmts[0].(*ast.ForInStmt)
	assert.True(t, ok)
	assert.NotNil(t, forInStmt.Index)
	assert.NotNil(t, forInStmt.Body)
	assert.Len(t, forInStmt.Body.List, 1)
}

func TestFileBuilder_Function(t *testing.T) {
	// 测试函数声明构建
	builder := NewFileBuilder()

	funcBuilder := builder.Function("add")
	funcBuilder.Params(
		&ast.Parameter{Name: Ident("a"), Type: Ident("num")},
		&ast.Parameter{Name: Ident("b"), Type: Ident("num")},
	)
	funcBuilder.Body(
		&ast.ReturnStmt{
			Return: token.Pos(1),
			X:      BinaryExpr(Ident("a"), token.PLUS, Ident("b")),
		},
	)

	fileBuilder := funcBuilder.End().(*FileBuilder)
	file := fileBuilder.Build()

	assert.Len(t, file.Stmts, 1)

	// 验证函数声明结构
	funcDecl, ok := file.Stmts[0].(*ast.FuncDecl)
	assert.True(t, ok)
	assert.Equal(t, "add", funcDecl.Name.Name)
	assert.Len(t, funcDecl.Recv, 2)
	assert.NotNil(t, funcDecl.Body)
	assert.Len(t, funcDecl.Body.List, 1)
}

func TestFileBuilder_Class(t *testing.T) {
	// 测试类声明构建
	builder := NewFileBuilder()

	classBuilder := builder.Class("Person")
	classBuilder.Field("name", Ident("string"))
	classBuilder.Field("age", Ident("num"))

	// 添加方法
	methodBuilder := classBuilder.Method("greet")
	methodBuilder.Body(
		&ast.ExprStmt{X: CmdExpr("echo", Str("Hello, I am"), Ident("name"))},
	)
	methodBuilder.End()

	fileBuilder := classBuilder.End().(*FileBuilder)
	file := fileBuilder.Build()

	assert.Len(t, file.Stmts, 1)

	// 验证类声明结构
	classDecl, ok := file.Stmts[0].(*ast.ClassDecl)
	assert.True(t, ok)
	assert.Equal(t, "Person", classDecl.Name.Name)
	assert.Len(t, classDecl.Fields.List, 2)
	assert.Len(t, classDecl.Methods, 1)
}

func TestFileBuilder_NestedStructures(t *testing.T) {
	// 测试嵌套结构构建
	builder := NewFileBuilder()

	// 在while循环中嵌套if语句
	whileBuilder := builder.While(BinaryExpr(Ident("i"), token.LT, Num("10")))

	// 在while循环体中构建if语句
	// ifBuilder := whileBuilder.If(BinaryExpr(Ident("i"), token.EQ, Num("5")))
	// ifBuilder.Then(
	// 	&ast.ExprStmt{X: CmdExpr("echo", Str("Found 5"))},
	// )
	// ifBuilder.Else(
	// 	&ast.ExprStmt{X: CmdExpr("echo", Ident("i"))},
	// )

	// 结束if构建，继续while构建
	// whileBuilder = ifBuilder.End().(*WhileBuilder)

	// 添加更多语句到while循环
	whileBuilder.AddStmt(&ast.AssignStmt{
		Lhs: Ident("i"),
		Tok: token.PLUS_ASSIGN,
		Rhs: Num("1"),
	})

	fileBuilder := whileBuilder.End().(*FileBuilder)
	file := fileBuilder.Build()

	assert.Len(t, file.Stmts, 1)

	// 验证嵌套结构
	whileStmt, ok := file.Stmts[0].(*ast.WhileStmt)
	assert.True(t, ok)
	assert.Len(t, whileStmt.Body.List, 2)

	// 验证第一个语句是if语句
	ifStmt, ok := whileStmt.Body.List[0].(*ast.IfStmt)
	assert.True(t, ok)
	assert.NotNil(t, ifStmt.Cond)
	assert.NotNil(t, ifStmt.Body)
	assert.NotNil(t, ifStmt.Else)
}

func TestFileBuilder_Chaining(t *testing.T) {
	// 测试链式调用
	builder := NewFileBuilder()

	file := builder.
		Echo(Str("First")).
		Assign(Ident("x"), Num("1")).
		Echo(Str("Second")).
		Assign(Ident("y"), Num("2")).
		Build()

	assert.Len(t, file.Stmts, 4)

	// 验证语句顺序
	firstEcho, ok := file.Stmts[0].(*ast.ExprStmt)
	assert.True(t, ok)
	firstCmd, ok := firstEcho.X.(*ast.CmdExpr)
	assert.True(t, ok)
	assert.Equal(t, "echo", firstCmd.Cmd.(*ast.Ident).Name)

	firstAssign, ok := file.Stmts[1].(*ast.AssignStmt)
	assert.True(t, ok)
	assert.Equal(t, "x", firstAssign.Lhs.(*ast.Ident).Name)
}

func TestFileBuilder_ComplexNesting(t *testing.T) {
	// 测试复杂的嵌套结构
	builder := NewFileBuilder()

	// 构建一个复杂的嵌套结构：函数中包含if-while嵌套
	funcBuilder := builder.Function("process")
	funcBuilder.Params(Ident("data"))

	// 在函数体中构建if语句
	// ifBuilder := funcBuilder.If(BinaryExpr(Ident("data"), token.NEQ, NULL))
	// ifBuilder.Then(
		// &ast.AssignStmt{Lhs: Ident("i"), Tok: token.ASSIGN, Rhs: Num("0")},
	// )

	// 在if的then分支中构建while循环
	// whileBuilder := ifBuilder.While(BinaryExpr(Ident("i"), token.LT, Ident("data")))
	// whileBuilder.Do(
		// &ast.ExprStmt{X: CmdExpr("process_item", Ident("i"))},
		// &ast.AssignStmt{Lhs: Ident("i"), Tok: token.PLUS_ASSIGN, Rhs: Num("1")},
	// )

	// 结束while，继续if
	// ifBuilder = whileBuilder.End().(*IfBuilder)
	// ifBuilder.Else(
		// &ast.ExprStmt{X: CmdExpr("echo", Str("No data"))},
	// )

	// 结束if，继续函数
	// funcBuilder = ifBuilder.End().(*FunctionBuilder)
	funcBuilder.Body(
		&ast.ReturnStmt{Return: token.Pos(1), X: Ident("result")},
	)

	fileBuilder := funcBuilder.End().(*FileBuilder)
	file := fileBuilder.Build()

	assert.Len(t, file.Stmts, 1)

	// 验证复杂的嵌套结构
	funcDecl, ok := file.Stmts[0].(*ast.FuncDecl)
	assert.True(t, ok)
	assert.Equal(t, "process", funcDecl.Name.Name)
	assert.Len(t, funcDecl.Body.List, 1)

	// 验证函数体中的if语句
	ifStmt, ok := funcDecl.Body.List[0].(*ast.IfStmt)
	assert.True(t, ok)
	assert.NotNil(t, ifStmt.Cond)
	assert.NotNil(t, ifStmt.Body)
	assert.NotNil(t, ifStmt.Else)

	// 验证if的then分支中的while语句
	assert.Len(t, ifStmt.Body.List, 2)
	whileStmt, ok := ifStmt.Body.List[1].(*ast.WhileStmt)
	assert.True(t, ok)
	assert.Len(t, whileStmt.Body.List, 2)
}

func TestBuilderHelpers(t *testing.T) {
	// 测试辅助函数
	ident := Ident("test")
	assert.Equal(t, "test", ident.Name)

	str := Str("hello")
	assert.Equal(t, "hello", str.Value)

	num := Num("42")
	assert.Equal(t, "42", num.Value)

	binary := BinaryExpr(Ident("a"), token.PLUS, Ident("b"))
	assert.Equal(t, token.PLUS, binary.Op)
	assert.Equal(t, "a", binary.X.(*ast.Ident).Name)
	assert.Equal(t, "b", binary.Y.(*ast.Ident).Name)

	cmd := CmdExpr("echo", Str("hello"))
	assert.Equal(t, "echo", cmd.Cmd.(*ast.Ident).Name)
	assert.Len(t, cmd.Args, 1)
}

func TestStatementBuilderInterface(t *testing.T) {
	// 测试StatementBuilder接口的实现
	builder := NewFileBuilder()

	// 测试FileBuilder实现了StatementBuilder接口
	var sb StatementBuilder = builder
	assert.NotNil(t, sb)

	// 测试IfBuilder实现了StatementBuilder接口
	ifBuilder := builder.If(BinaryExpr(Ident("x"), token.EQ, Num("10")))
	sb = ifBuilder
	assert.NotNil(t, sb)

	// 测试WhileBuilder实现了StatementBuilder接口
	whileBuilder := builder.While(BinaryExpr(Ident("i"), token.LT, Num("10")))
	sb = whileBuilder
	assert.NotNil(t, sb)
}
