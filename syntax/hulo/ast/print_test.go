// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package ast

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrintNodeType(t *testing.T) {
	// Test with a simple identifier
	ident := &Ident{Name: "hello"}

	var buf strings.Builder
	Inspect(ident, &buf)

	expected := `*ast.Ident (Name: "hello")`
	assert.Equal(t, expected, strings.TrimSpace(buf.String()))
}

func TestPrintNodeTypeComplex(t *testing.T) {
	// Test with a more complex structure
	file := &File{
		Stmts: []Stmt{
			&AssignStmt{
				Lhs: &Ident{Name: "x"},
				Rhs: &NumericLiteral{Value: "42"},
			},
			&FuncDecl{
				Name: &Ident{Name: "main"},
				Body: &BlockStmt{
					List: []Stmt{
						&ExprStmt{
							X: &CallExpr{
								Fun: &Ident{Name: "println"},
								Recv: []Expr{
									&StringLiteral{Value: "Hello, World!"},
								},
							},
						},
					},
				},
			},
		},
	}

	var buf strings.Builder
	Inspect(file, &buf)

	result := buf.String()

	// Check that it contains expected node types
	assert.Contains(t, result, "*ast.File")
	assert.Contains(t, result, "*ast.AssignStmt")
	assert.Contains(t, result, "*ast.FuncDecl")
	assert.Contains(t, result, "*ast.BlockStmt")
	assert.Contains(t, result, "*ast.ExprStmt")
	assert.Contains(t, result, "*ast.CallExpr")
	assert.Contains(t, result, "*ast.Ident")
	assert.Contains(t, result, "*ast.NumericLiteral")
	assert.Contains(t, result, "*ast.StringLiteral")

	// Check that it contains the actual values
	assert.Contains(t, result, `Name: "x"`)
	assert.Contains(t, result, `Value: "42"`)
	assert.Contains(t, result, `Name: "main"`)
	assert.Contains(t, result, `Name: "println"`)
	assert.Contains(t, result, `Value: "Hello, World!"`)
}

func TestInspect(t *testing.T) {
	file := &File{
		Stmts: []Stmt{
			&AssignStmt{
				Lhs: &Ident{Name: "x"},
				Rhs: &NumericLiteral{Value: "42"},
			},
			&FuncDecl{
				Name: &Ident{Name: "main"},
				Body: &BlockStmt{
					List: []Stmt{
						&ExprStmt{
							X: &CallExpr{
								Fun: &Ident{Name: "println"},
								Recv: []Expr{
									&StringLiteral{Value: "Hello, World!"},
								},
							},
						},
					},
				},
			},
		},
	}
	Inspect(file, os.Stdout)
}
