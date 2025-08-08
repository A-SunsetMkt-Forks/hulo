// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package astutil provides utilities for building and manipulating bash AST nodes.
//
// This package offers two main approaches for creating bash AST:
//
// 1. Builder Pattern: Use FileBuilder and related builders for fluent API
// 2. Direct Construction: Use helper functions to create individual AST nodes
//
// # Builder Pattern Example
//
//	file := NewFileBuilder().
//		Echo(Word("Hello, World!")).
//		If(Eq(Ident("x"), Word("1"))).
//		Then(
//			CmdStmt("echo", Word("x is 1")),
//		).
//		Else(
//			CmdStmt("echo", Word("x is not 1")),
//		).
//		EndIf().
//		Build()
//
// # Direct Construction Example
//
//	file := &ast.File{
//		Stmts: []ast.Stmt{
//			CmdStmt("echo", Word("Hello")),
//			IfStmt(
//				Eq(Ident("x"), Word("1")),
//				BlockStmt(CmdStmt("echo", Word("x is 1"))),
//				BlockStmt(CmdStmt("echo", Word("x is not 1"))),
//			),
//		},
//	}
//
// # Available Builders
//
// - FileBuilder: Main builder for creating bash files
// - IfBuilder: For if-elif-else statements
// - WhileBuilder: For while loops
// - UntilBuilder: For until loops
// - ForInBuilder: For for-in loops
// - FunctionBuilder: For function declarations
//
// # Helper Functions
//
// The package provides many helper functions for creating AST nodes:
//
// - Word(), Ident(): Create basic expressions
// - CmdExpr(), CmdStmt(): Create command expressions and statements
// - Eq(), Ne(), Lt(), Gt(): Create comparison expressions
// - FileExists(), IsFile(), IsDir(): Create file test expressions
// - Echo(), Exit(), Break(): Create common commands
// - And(), Or(), Not(): Create logical expressions
//
// # Usage Guidelines
//
// 1. Use builders for complex, nested structures
// 2. Use helper functions for simple, individual nodes
// 3. Combine both approaches as needed
// 4. Always set appropriate token positions for proper AST structure
package astutil
