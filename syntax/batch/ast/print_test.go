// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package ast

import (
	"strings"
	"testing"

	"github.com/hulo-lang/hulo/syntax/batch/token"
)

func TestInspect(t *testing.T) {
	tests := []struct {
		name     string
		node     Node
		expected string
	}{
		{
			name: "File with statements",
			node: &File{
				Stmts: []Stmt{
					&ExprStmt{
						X: &Lit{Val: "echo"},
					},
					&AssignStmt{
						Lhs: &Lit{Val: "var"},
						Rhs: &Lit{Val: "value"},
					},
				},
			},
			expected: "*ast.File",
		},
		{
			name: "If statement",
			node: &IfStmt{
				Cond: &BinaryExpr{
					X:  &Lit{Val: "a"},
					Op: token.EQU,
					Y:  &Lit{Val: "b"},
				},
				Body: &BlockStmt{
					List: []Stmt{
						&ExprStmt{X: &Lit{Val: "echo"}},
					},
				},
			},
			expected: "*ast.IfStmt",
		},
		{
			name: "For statement",
			node: &ForStmt{
				X:    &Lit{Val: "i"},
				List: &Lit{Val: "1 2 3"},
				Body: &BlockStmt{
					List: []Stmt{
						&ExprStmt{X: &Lit{Val: "echo"}},
					},
				},
			},
			expected: "*ast.ForStmt",
		},
		{
			name: "Function declaration",
			node: &FuncDecl{
				Name: "test",
				Body: &BlockStmt{
					List: []Stmt{
						&ExprStmt{X: &Lit{Val: "echo"}},
					},
				},
			},
			expected: "*ast.FuncDecl",
		},
		{
			name: "Call statement",
			node: &CallStmt{
				Name:   "function",
				IsFile: false,
				Recv: []Expr{
					&Lit{Val: "arg1"},
					&Lit{Val: "arg2"},
				},
			},
			expected: "*ast.CallStmt",
		},
		{
			name: "Word expression",
			node: &Word{
				Parts: []Expr{
					&Lit{Val: "echo"},
					&Lit{Val: "hello"},
				},
			},
			expected: "*ast.Word",
		},
		{
			name: "Binary expression",
			node: &BinaryExpr{
				X:  &Lit{Val: "a"},
				Op: token.EQU,
				Y:  &Lit{Val: "b"},
			},
			expected: "*ast.BinaryExpr",
		},
		{
			name: "Unary expression",
			node: &UnaryExpr{
				Op: token.NOT,
				X:  &Lit{Val: "condition"},
			},
			expected: "*ast.UnaryExpr",
		},
		{
			name: "Single quote",
			node: &SglQuote{
				Val: &Lit{Val: "variable"},
			},
			expected: "*ast.SglQuote",
		},
		{
			name: "Double quote",
			node: &DblQuote{
				DelayedExpansion: false,
				Val:              &Lit{Val: "variable"},
			},
			expected: "*ast.DblQuote",
		},
		{
			name: "Command expression",
			node: &CmdExpr{
				Name: &Lit{Val: "echo"},
				Recv: []Expr{
					&Lit{Val: "hello"},
					&Lit{Val: "world"},
				},
			},
			expected: "*ast.CmdExpr",
		},
		{
			name: "Comment",
			node: &Comment{
				Tok:  token.REM,
				Text: " This is a comment",
			},
			expected: "*ast.Comment",
		},
		{
			name: "Comment group",
			node: &CommentGroup{
				Comments: []*Comment{
					{Tok: token.REM, Text: " Comment 1"},
					{Tok: token.DOUBLE_COLON, Text: " Comment 2"},
				},
			},
			expected: "*ast.CommentGroup",
		},
		{
			name: "Goto statement",
			node: &GotoStmt{
				Label: "label",
			},
			expected: "*ast.GotoStmt",
		},
		{
			name: "Label statement",
			node: &LabelStmt{
				Name: "label",
			},
			expected: "*ast.LabelStmt",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf strings.Builder
			Inspect(tt.node, &buf)
			output := buf.String()

			if !strings.Contains(output, tt.expected) {
				t.Errorf("Inspect() output = %v, want to contain %v", output, tt.expected)
			}
		})
	}
}

func TestInspectNil(t *testing.T) {
	var buf strings.Builder
	Inspect(nil, &buf)
	output := buf.String()

	if output != "nil" {
		t.Errorf("Inspect(nil) = %q, want %q", output, "nil")
	}
}
