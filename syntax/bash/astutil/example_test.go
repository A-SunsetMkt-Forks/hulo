// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package astutil

import (
	"fmt"

	"github.com/hulo-lang/hulo/syntax/bash/ast"
	"github.com/hulo-lang/hulo/syntax/bash/token"
)

func ExampleNewFileBuilder() {
	// 创建一个简单的 bash 文件
	file := NewFileBuilder().
		Echo(Word("Hello, World!")).
		Assign(Ident("name"), Word("Alice")).
		Echo(Word("Hello,"), Ident("name")).
		Build()

	fmt.Printf("File has %d statements\n", len(file.Stmts))
	// Output: File has 3 statements
}

func ExampleIfBuilder() {
	// 创建一个 if-elif-else 语句
	file := NewFileBuilder().
		If(Eq(Ident("x"), Word("1"))).
		Then(
			CmdStmt("echo", Word("x is 1")),
		).
		ElseIf(Eq(Ident("x"), Word("2"))).
		Then(
			CmdStmt("echo", Word("x is 2")),
		).
		Else(
			CmdStmt("echo", Word("x is something else")),
		).
		EndIf().
		Build()

	fmt.Printf("File has %d statements\n", len(file.Stmts))
	// Output: File has 1 statements
}

func ExampleWhileBuilder() {
	// 创建一个 while 循环
	file := NewFileBuilder().
		Assign(Ident("count"), Word("0")).
		While(Lt(Ident("count"), Word("5"))).
		Do(
			CmdStmt("echo", Word("Count:"), Ident("count")),
			AssignStmt(Ident("count"), Binary(Ident("count"), token.Plus, Word("1"))),
		).
		EndWhile().
		Build()

	fmt.Printf("File has %d statements\n", len(file.Stmts))
	// Output: File has 3 statements
}

func ExampleForInBuilder() {
	// 创建一个 for-in 循环
	file := NewFileBuilder().
		ForIn(Ident("item"), Word("apple banana cherry")).
		Do(
			CmdStmt("echo", Word("Processing:"), Ident("item")),
		).
		EndFor().
		Build()

	fmt.Printf("File has %d statements\n", len(file.Stmts))
	// Output: File has 1 statements
}

func ExampleFunctionBuilder() {
	// 创建一个函数定义
	file := NewFileBuilder().
		Function("greet").
		Body(
			CmdStmt("echo", Word("Hello,"), Ident("1")),
		).
		EndFunction().
		Echo(Word("Function defined")).
		Build()

	fmt.Printf("File has %d statements\n", len(file.Stmts))
	// Output: File has 2 statements
}

func ExampleFileExists() {
	// 创建文件测试条件
	file := NewFileBuilder().
		If(FileExists(Word("config.txt"))).
		Then(
			CmdStmt("echo", Word("Config file exists")),
		).
		Else(
			CmdStmt("echo", Word("Config file not found")),
		).
		EndIf().
		Build()

	fmt.Printf("File has %d statements\n", len(file.Stmts))
	// Output: File has 1 statements
}

func ExampleIsNotEmpty() {
	// 创建字符串测试条件
	file := NewFileBuilder().
		If(IsNotEmpty(Ident("name"))).
		Then(
			CmdStmt("echo", Word("Hello,"), Ident("name")),
		).
		Else(
			CmdStmt("echo", Word("Please provide a name")),
		).
		EndIf().
		Build()

	fmt.Printf("File has %d statements\n", len(file.Stmts))
	// Output: File has 1 statements
}

func ExampleAnd() {
	// 使用逻辑操作符
	file := NewFileBuilder().
		If(And(
			IsFile(Word("data.txt")),
			IsReadable(Word("data.txt")),
		)).
		Then(
			CmdStmt("cat", Word("data.txt")),
		).
		Else(
			CmdStmt("echo", Word("File is not readable or not a file")),
		).
		EndIf().
		Build()

	fmt.Printf("File has %d statements\n", len(file.Stmts))
	// Output: File has 1 statements
}

func ExampleNewFileBuilder_complex() {
	// 创建一个复杂的脚本
	file := NewFileBuilder().
		// 定义变量
		Assign(Ident("script_name"), Word("my_script.sh")).
		Assign(Ident("log_file"), Word("output.log")).

		// 检查脚本是否存在且可执行
		If(And(
			FileExists(Ident("script_name")),
			IsExecutable(Ident("script_name")),
		)).
		Then(
			// 执行脚本并重定向输出
			CmdStmt("bash", Ident("script_name"), Word(">"), Ident("log_file")),
			CmdStmt("echo", Word("Script executed successfully")),
		).
		Else(
			// 错误处理
			CmdStmt("echo", Word("Error: Script not found or not executable")),
			ExprStmt(Exit(Word("1"))),
		).
		EndIf().

		// 检查日志文件是否创建
		If(FileExists(Ident("log_file"))).
		Then(
			CmdStmt("echo", Word("Log file created:"), Ident("log_file")),
			CmdStmt("wc", Word("-l"), Ident("log_file")),
		).
		EndIf().
		Build()

	fmt.Printf("Complex script has %d statements\n", len(file.Stmts))
	// Output: Complex script has 4 statements
}

func ExampleBlockStmt() {
	// 直接使用辅助函数构建 AST
	file := &ast.File{
		Stmts: []ast.Stmt{
			// 赋值语句
			AssignStmt(Ident("message"), Word("Hello from direct construction")),

			// 条件语句
			IfStmt(
				Eq(Ident("message"), Word("Hello")),
				Block(
					CmdStmt("echo", Word("Message is Hello")),
				),
				Block(
					CmdStmt("echo", Word("Message is not Hello")),
				),
			),

			// 循环语句
			WhileStmt(
				Lt(Ident("counter"), Word("3")),
				[]ast.Stmt{
					CmdStmt("echo", Word("Counter:"), Ident("counter")),
					AssignStmt(Ident("counter"), Binary(Ident("counter"), token.Plus, Word("1"))),
				},
			),
		},
	}

	fmt.Printf("Directly constructed file has %d statements\n", len(file.Stmts))
	// Output: Directly constructed file has 3 statements
}
