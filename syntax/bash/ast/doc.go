// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package ast provides the abstract syntax tree (AST) representation for bash scripts.
//
// This package defines the AST nodes that represent bash script syntax,
// including statements, expressions, and other language constructs.
//
// # Overview
//
// The AST package provides a complete representation of bash script syntax,
// allowing for code generation, analysis, and transformation of bash code. The AST
// nodes are designed to capture all the semantic information while maintaining
// the hierarchical structure. This package focuses on AST manipulation and
// does not include parsing functionality.
//
// # Key Features
//
//   - Complete bash syntax support including control structures, functions, and expressions
//   - Visitor pattern for traversing and analyzing AST nodes
//   - Pretty printing and detailed inspection capabilities
//   - Support for all bash language constructs including pipelines, redirections, and substitutions
//
// # Usage
//
// The AST package is used for creating and manipulating AST representations
// of bash scripts programmatically:
//
//	// Create AST nodes programmatically
//	file := &ast.File{
//		Stmts: []ast.Stmt{
//			&ast.ExprStmt{
//				X: &ast.CmdExpr{
//					Name: &ast.Ident{Name: "echo"},
//					Recv: []ast.Expr{&ast.Word{Val: "Hello, World!"}},
//				},
//			},
//		},
//	}
//
//	// Inspect the AST structure
//	ast.Inspect(file, os.Stdout)
//
//	// Pretty print the AST back to bash code
//	ast.Print(file)
//
// # AST Node Types
//
// The package defines several categories of AST nodes:
//
// ## Statements
//   - File: represents a complete bash script file
//   - FuncDecl: function declarations
//   - IfStmt, WhileStmt, ForStmt: control flow statements
//   - AssignStmt: variable assignments
//   - ExprStmt: expression statements
//   - BlockStmt: statement blocks
//
// ## Expressions
//   - BinaryExpr: binary operations (+, -, ==, etc.)
//   - UnaryExpr: unary operations (!, -, etc.)
//   - CmdExpr: command executions
//   - Word, Ident: literals and identifiers
//   - VarExpExpr: variable expansions ($var)
//   - ParamExpExpr: parameter expansions (${var})
//
// ## Special Constructs
//   - PipelineExpr: command pipelines (cmd1 | cmd2)
//   - Redirect: input/output redirections
//   - CmdSubst: command substitutions ($(cmd))
//   - ProcSubst: process substitutions (<(cmd))
//   - ArithExpr: arithmetic expressions ($((expr)))
//
// # Visitor Pattern
//
// The AST supports the visitor pattern for traversing nodes:
//
//	type MyVisitor struct{}
//
//	func (v *MyVisitor) Visit(node ast.Node) ast.Visitor {
//	    // Process the node
//	    return v // Return v to continue visiting children
//	}
//
//	ast.Walk(&MyVisitor{}, file)
//
// # Printing and Inspection
//
// The package provides several ways to examine AST structures:
//
//   - Print(node): Pretty print to stdout
//   - Write(node, writer): Pretty print to a writer
//   - String(node): Return pretty printed string
//   - Inspect(node, writer): Detailed structural inspection
//
// # Examples
//
// See the example_test.go file for comprehensive usage examples
// demonstrating how to create, traverse, and inspect AST nodes.
//
// # Design Philosophy
//
// The AST design follows these principles:
//
//   - Complete representation: All bash syntax constructs are represented
//   - Position-aware: Each node includes token position information
//   - Visitor-friendly: Easy traversal and analysis
//   - Debug-friendly: Rich inspection and printing capabilities
//   - Code generation ready: Optimized for generating bash code from AST
//
// # Performance Considerations
//
// AST nodes are designed for correctness and clarity over performance.
// For large scripts, consider using the visitor pattern to process
// nodes incrementally rather than building the entire AST in memory.
package ast
