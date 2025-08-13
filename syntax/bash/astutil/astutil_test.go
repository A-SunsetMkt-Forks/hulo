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

// ===== 基础表达式测试 =====

func TestWord(t *testing.T) {
	word := Word("hello")
	assert.NotNil(t, word)
	assert.Equal(t, "hello", word.Val)
	assert.True(t, word.ValPos.IsValid())
}

func TestIdent(t *testing.T) {
	ident := Ident("variable")
	assert.NotNil(t, ident)
	assert.Equal(t, "variable", ident.Name)
	assert.True(t, ident.NamePos.IsValid())
}

func TestBinary(t *testing.T) {
	x := Ident("a")
	y := Ident("b")
	binary := Binary(x, token.Equal, y)

	assert.NotNil(t, binary)
	assert.Equal(t, x, binary.X)
	assert.Equal(t, token.Equal, binary.Op)
	assert.Equal(t, y, binary.Y)
	assert.True(t, binary.OpPos.IsValid())
}

func TestUnary(t *testing.T) {
	x := Ident("flag")
	unary := Unary(token.ExclMark, x)

	assert.NotNil(t, unary)
	assert.Equal(t, token.ExclMark, unary.Op)
	assert.Equal(t, x, unary.X)
	assert.True(t, unary.OpPos.IsValid())
}

// ===== 语句测试 =====

func TestExprStmt(t *testing.T) {
	expr := CmdExpr("echo", Word("hello"))
	stmt := ExprStmt(expr)

	assert.NotNil(t, stmt)
	assert.Equal(t, expr, stmt.X)
}

func TestCmdExpr(t *testing.T) {
	cmd := CmdExpr("echo", Word("hello"), Word("world"))

	assert.NotNil(t, cmd)
	assert.Equal(t, "echo", cmd.Name.Name)
	assert.Len(t, cmd.Recv, 2)
	assert.Equal(t, "hello", cmd.Recv[0].(*ast.Word).Val)
	assert.Equal(t, "world", cmd.Recv[1].(*ast.Word).Val)
}

func TestCmdStmt(t *testing.T) {
	stmt := CmdStmt("echo", Word("hello"))

	assert.NotNil(t, stmt)
	cmdExpr, ok := stmt.X.(*ast.CmdExpr)
	assert.True(t, ok)
	assert.Equal(t, "echo", cmdExpr.Name.Name)
}

func TestAssignStmt(t *testing.T) {
	lhs := Ident("var")
	rhs := Word("value")
	stmt := AssignStmt(lhs, rhs)

	assert.NotNil(t, stmt)
	assert.Equal(t, lhs, stmt.Lhs)
	assert.Equal(t, rhs, stmt.Rhs)
	assert.True(t, stmt.Assign.IsValid())
}

func TestLocalAssignStmt(t *testing.T) {
	lhs := Ident("var")
	rhs := Word("value")
	stmt := LocalAssignStmt(lhs, rhs)

	assert.NotNil(t, stmt)
	assert.Equal(t, lhs, stmt.Lhs)
	assert.Equal(t, rhs, stmt.Rhs)
	assert.True(t, stmt.Local.IsValid())
	assert.True(t, stmt.Assign.IsValid())
}

func TestReturnStmt(t *testing.T) {
	expr := Word("0")
	stmt := ReturnStmt(expr)

	assert.NotNil(t, stmt)
	assert.Equal(t, expr, stmt.X)
	assert.True(t, stmt.Return.IsValid())
}

func TestBlockStmt(t *testing.T) {
	stmt1 := CmdStmt("echo", Word("hello"))
	stmt2 := CmdStmt("echo", Word("world"))
	block := Block(stmt1, stmt2)

	assert.NotNil(t, block)
	assert.Len(t, block.List, 2)
	assert.Equal(t, stmt1, block.List[0])
	assert.Equal(t, stmt2, block.List[1])
}

// ===== 测试表达式测试 =====

func TestTestExpr(t *testing.T) {
	cond := Eq(Ident("x"), Word("1"))
	test := TestExpr(cond)

	assert.NotNil(t, test)
	assert.Equal(t, cond, test.X)
	assert.True(t, test.Lbrack.IsValid())
	assert.True(t, test.Rbrack.IsValid())
}

func TestExtendedTestExpr(t *testing.T) {
	cond := Eq(Ident("x"), Word("1"))
	test := ExtendedTestExpr(cond)

	assert.NotNil(t, test)
	assert.Equal(t, cond, test.X)
	assert.True(t, test.Lbrack.IsValid())
	assert.True(t, test.Rbrack.IsValid())
}

// ===== 变量扩展测试 =====

func TestVarExp(t *testing.T) {
	name := Ident("var")
	exp := VarExp(name)

	assert.NotNil(t, exp)
	assert.Equal(t, name, exp.X)
	assert.True(t, exp.Dollar.IsValid())
}

// ===== 比较操作符测试 =====

func TestEq(t *testing.T) {
	x := Ident("a")
	y := Ident("b")
	eq := Eq(x, y)

	assert.NotNil(t, eq)
	assert.Equal(t, x, eq.X)
	assert.Equal(t, token.Equal, eq.Op)
	assert.Equal(t, y, eq.Y)
}

func TestNe(t *testing.T) {
	x := Ident("a")
	y := Ident("b")
	ne := Ne(x, y)

	assert.NotNil(t, ne)
	assert.Equal(t, x, ne.X)
	assert.Equal(t, token.NotEqual, ne.Op)
	assert.Equal(t, y, ne.Y)
}

func TestLt(t *testing.T) {
	x := Ident("a")
	y := Ident("b")
	lt := Lt(x, y)

	assert.NotNil(t, lt)
	assert.Equal(t, x, lt.X)
	assert.Equal(t, token.TsLss, lt.Op)
	assert.Equal(t, y, lt.Y)
}

func TestGt(t *testing.T) {
	x := Ident("a")
	y := Ident("b")
	gt := Gt(x, y)

	assert.NotNil(t, gt)
	assert.Equal(t, x, gt.X)
	assert.Equal(t, token.TsGtr, gt.Op)
	assert.Equal(t, y, gt.Y)
}

func TestLe(t *testing.T) {
	x := Ident("a")
	y := Ident("b")
	le := Le(x, y)

	assert.NotNil(t, le)
	assert.Equal(t, x, le.X)
	assert.Equal(t, token.TsLeq, le.Op)
	assert.Equal(t, y, le.Y)
}

func TestGe(t *testing.T) {
	x := Ident("a")
	y := Ident("b")
	ge := Ge(x, y)

	assert.NotNil(t, ge)
	assert.Equal(t, x, ge.X)
	assert.Equal(t, token.TsGeq, ge.Op)
	assert.Equal(t, y, ge.Y)
}

// ===== 文件测试操作符测试 =====

func TestFileExists(t *testing.T) {
	file := Word("test.txt")
	exists := FileExists(file)

	assert.NotNil(t, exists)
	assert.Equal(t, "-e", exists.X.(*ast.Word).Val)
	assert.Equal(t, token.TsExists, exists.Op)
	assert.Equal(t, file, exists.Y)
}

func TestIsFile(t *testing.T) {
	file := Word("test.txt")
	isFile := IsFile(file)

	assert.NotNil(t, isFile)
	assert.Equal(t, "-f", isFile.X.(*ast.Word).Val)
	assert.Equal(t, token.TsRegFile, isFile.Op)
	assert.Equal(t, file, isFile.Y)
}

func TestIsDir(t *testing.T) {
	file := Word("testdir")
	isDir := IsDir(file)

	assert.NotNil(t, isDir)
	assert.Equal(t, "-d", isDir.X.(*ast.Word).Val)
	assert.Equal(t, token.TsDirect, isDir.Op)
	assert.Equal(t, file, isDir.Y)
}

func TestIsReadable(t *testing.T) {
	file := Word("test.txt")
	readable := IsReadable(file)

	assert.NotNil(t, readable)
	assert.Equal(t, "-r", readable.X.(*ast.Word).Val)
	assert.Equal(t, token.TsRead, readable.Op)
	assert.Equal(t, file, readable.Y)
}

func TestIsWritable(t *testing.T) {
	file := Word("test.txt")
	writable := IsWritable(file)

	assert.NotNil(t, writable)
	assert.Equal(t, "-w", writable.X.(*ast.Word).Val)
	assert.Equal(t, token.TsWrite, writable.Op)
	assert.Equal(t, file, writable.Y)
}

func TestIsExecutable(t *testing.T) {
	file := Word("test.sh")
	executable := IsExecutable(file)

	assert.NotNil(t, executable)
	assert.Equal(t, "-x", executable.X.(*ast.Word).Val)
	assert.Equal(t, token.TsExec, executable.Op)
	assert.Equal(t, file, executable.Y)
}

// ===== 字符串测试操作符测试 =====

func TestIsEmpty(t *testing.T) {
	str := Ident("var")
	empty := IsEmpty(str)

	assert.NotNil(t, empty)
	assert.Equal(t, "-z", empty.X.(*ast.Word).Val)
	assert.Equal(t, token.TsEmpStr, empty.Op)
	assert.Equal(t, str, empty.Y)
}

func TestIsNotEmpty(t *testing.T) {
	str := Ident("var")
	notEmpty := IsNotEmpty(str)

	assert.NotNil(t, notEmpty)
	assert.Equal(t, "-n", notEmpty.X.(*ast.Word).Val)
	assert.Equal(t, token.TsNempStr, notEmpty.Op)
	assert.Equal(t, str, notEmpty.Y)
}

// ===== 常用命令测试 =====

func TestEcho(t *testing.T) {
	cmd := Echo(Word("hello"), Word("world"))

	assert.NotNil(t, cmd)
	assert.Equal(t, "echo", cmd.Name.Name)
	assert.Len(t, cmd.Recv, 2)
	assert.Equal(t, "hello", cmd.Recv[0].(*ast.Word).Val)
	assert.Equal(t, "world", cmd.Recv[1].(*ast.Word).Val)
}

func TestExit(t *testing.T) {
	// 测试无参数的 exit
	cmd1 := Exit()
	assert.NotNil(t, cmd1)
	assert.Equal(t, "exit", cmd1.Name.Name)
	assert.Len(t, cmd1.Recv, 0)

	// 测试带参数的 exit
	cmd2 := Exit(Word("1"))
	assert.NotNil(t, cmd2)
	assert.Equal(t, "exit", cmd2.Name.Name)
	assert.Len(t, cmd2.Recv, 1)
	assert.Equal(t, "1", cmd2.Recv[0].(*ast.Word).Val)
}

func TestBreak(t *testing.T) {
	// 测试无参数的 break
	cmd1 := Break()
	assert.NotNil(t, cmd1)
	assert.Equal(t, "break", cmd1.Name.Name)
	assert.Len(t, cmd1.Recv, 0)

	// 测试带参数的 break
	cmd2 := Break(2)
	assert.NotNil(t, cmd2)
	assert.Equal(t, "break", cmd2.Name.Name)
	assert.Len(t, cmd2.Recv, 1)
	assert.Equal(t, "2", cmd2.Recv[0].(*ast.Word).Val)
}

func TestContinue(t *testing.T) {
	// 测试无参数的 continue
	cmd1 := Continue()
	assert.NotNil(t, cmd1)
	assert.Equal(t, "continue", cmd1.Name.Name)
	assert.Len(t, cmd1.Recv, 0)

	// 测试带参数的 continue
	cmd2 := Continue(1)
	assert.NotNil(t, cmd2)
	assert.Equal(t, "continue", cmd2.Name.Name)
	assert.Len(t, cmd2.Recv, 1)
	assert.Equal(t, "1", cmd2.Recv[0].(*ast.Word).Val)
}

func TestTrue(t *testing.T) {
	cmd := True()
	assert.NotNil(t, cmd)
	assert.Equal(t, "true", cmd.Name.Name)
	assert.Len(t, cmd.Recv, 0)
}

func TestFalse(t *testing.T) {
	cmd := False()
	assert.NotNil(t, cmd)
	assert.Equal(t, "false", cmd.Name.Name)
	assert.Len(t, cmd.Recv, 0)
}

// ===== 逻辑操作符测试 =====

func TestAnd(t *testing.T) {
	x := Eq(Ident("a"), Word("1"))
	y := Eq(Ident("b"), Word("2"))
	and := And(x, y)

	assert.NotNil(t, and)
	assert.Equal(t, x, and.X)
	assert.Equal(t, token.AndAnd, and.Op)
	assert.Equal(t, y, and.Y)
}

func TestOr(t *testing.T) {
	x := Eq(Ident("a"), Word("1"))
	y := Eq(Ident("b"), Word("2"))
	or := Or(x, y)

	assert.NotNil(t, or)
	assert.Equal(t, x, or.X)
	assert.Equal(t, token.OrOr, or.Op)
	assert.Equal(t, y, or.Y)
}

func TestNot(t *testing.T) {
	x := Eq(Ident("a"), Word("1"))
	not := Not(x)

	assert.NotNil(t, not)
	assert.Equal(t, token.ExclMark, not.Op)
	assert.Equal(t, x, not.X)
}

// ===== 控制流语句测试 =====

func TestIfStmt(t *testing.T) {
	cond := Eq(Ident("x"), Word("1"))
	body := Block(CmdStmt("echo", Word("true")))
	elseBody := Block(CmdStmt("echo", Word("false")))

	ifStmt := IfStmt(cond, body, elseBody)

	assert.NotNil(t, ifStmt)
	assert.Equal(t, cond, ifStmt.Cond)
	assert.Equal(t, body, ifStmt.Body)
	assert.Equal(t, elseBody, ifStmt.Else)
	assert.True(t, ifStmt.If.IsValid())
	assert.True(t, ifStmt.Semi.IsValid())
	assert.True(t, ifStmt.Then.IsValid())
	assert.True(t, ifStmt.Fi.IsValid())
}

func TestWhileStmt(t *testing.T) {
	cond := Lt(Ident("i"), Word("10"))
	body := []ast.Stmt{
		CmdStmt("echo", Ident("i")),
		AssignStmt(Ident("i"), Binary(Ident("i"), token.Plus, Word("1"))),
	}

	whileStmt := WhileStmt(cond, body)

	assert.NotNil(t, whileStmt)
	assert.Equal(t, cond, whileStmt.Cond)
	assert.Len(t, whileStmt.Body.List, 2)
	assert.True(t, whileStmt.While.IsValid())
	assert.True(t, whileStmt.Semi.IsValid())
	assert.True(t, whileStmt.Do.IsValid())
	assert.True(t, whileStmt.Done.IsValid())
}

func TestUntilStmt(t *testing.T) {
	cond := Eq(Ident("flag"), Word("true"))
	body := []ast.Stmt{
		CmdStmt("echo", Word("waiting")),
		CmdStmt("sleep", Word("1")),
	}

	untilStmt := UntilStmt(cond, body)

	assert.NotNil(t, untilStmt)
	assert.Equal(t, cond, untilStmt.Cond)
	assert.Len(t, untilStmt.Body.List, 2)
	assert.True(t, untilStmt.Until.IsValid())
	assert.True(t, untilStmt.Semi.IsValid())
	assert.True(t, untilStmt.Do.IsValid())
	assert.True(t, untilStmt.Done.IsValid())
}

func TestForInStmt(t *testing.T) {
	varName := Ident("item")
	list := Word("apple banana cherry")
	body := []ast.Stmt{
		CmdStmt("echo", Ident("item")),
	}

	forStmt := ForInStmt(varName, list, body)

	assert.NotNil(t, forStmt)
	assert.Equal(t, varName, forStmt.Var)
	assert.Equal(t, list, forStmt.List)
	assert.Len(t, forStmt.Body.List, 1)
	assert.True(t, forStmt.For.IsValid())
	assert.True(t, forStmt.In.IsValid())
	assert.True(t, forStmt.Semi.IsValid())
	assert.True(t, forStmt.Do.IsValid())
	assert.True(t, forStmt.Done.IsValid())
}

func TestFuncDecl(t *testing.T) {
	name := "test_func"
	body := []ast.Stmt{
		CmdStmt("echo", Word("Hello")),
		ReturnStmt(Word("0")),
	}

	funcDecl := FuncDecl(name, body)

	assert.NotNil(t, funcDecl)
	assert.Equal(t, name, funcDecl.Name.Name)
	assert.Len(t, funcDecl.Body.List, 2)
	assert.True(t, funcDecl.Function.IsValid())
	assert.True(t, funcDecl.Lparen.IsValid())
	assert.True(t, funcDecl.Rparen.IsValid())
}

// ===== 复杂组合测试 =====

func TestComplexCondition(t *testing.T) {
	// 测试复杂的条件组合
	fileExists := FileExists(Word("config.txt"))
	isReadable := IsReadable(Word("config.txt"))
	notEmpty := IsNotEmpty(Ident("content"))

	complexCond := And(
		And(fileExists, isReadable),
		notEmpty,
	)

	assert.NotNil(t, complexCond)
	assert.Equal(t, token.AndAnd, complexCond.Op)

	// 验证第一个 And 操作
	firstAnd, ok := complexCond.X.(*ast.BinaryExpr)
	assert.True(t, ok)
	assert.Equal(t, token.AndAnd, firstAnd.Op)
	assert.Equal(t, fileExists, firstAnd.X)
	assert.Equal(t, isReadable, firstAnd.Y)

	// 验证第二个操作数
	assert.Equal(t, notEmpty, complexCond.Y)
}

func TestNestedCommands(t *testing.T) {
	// 测试嵌套的命令结构
	innerCmd := CmdExpr("grep", Word("pattern"), Ident("file"))
	outerCmd := CmdExpr("echo", innerCmd)

	assert.NotNil(t, outerCmd)
	assert.Equal(t, "echo", outerCmd.Name.Name)
	assert.Len(t, outerCmd.Recv, 1)
	assert.Equal(t, innerCmd, outerCmd.Recv[0])
}

func TestMultipleAssignments(t *testing.T) {
	// 测试多个赋值语句
	assign1 := AssignStmt(Ident("var1"), Word("value1"))
	assign2 := AssignStmt(Ident("var2"), Word("value2"))
	assign3 := LocalAssignStmt(Ident("local_var"), Word("local_value"))

	assert.NotNil(t, assign1)
	assert.NotNil(t, assign2)
	assert.NotNil(t, assign3)

	// 验证 local 赋值有 Local 字段
	assert.True(t, assign3.Local.IsValid())
	assert.False(t, assign1.Local.IsValid())
	assert.False(t, assign2.Local.IsValid())
}
