// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package ast

import (
	"strings"
	"testing"

	"github.com/hulo-lang/hulo/syntax/bash/token"
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
						X: &Word{Val: "echo"},
					},
					&AssignStmt{
						Lhs: &Ident{Name: "var"},
						Rhs: &Word{Val: "value"},
					},
				},
			},
			expected: "*ast.File",
		},
		{
			name: "Function declaration",
			node: &FuncDecl{
				Name: &Ident{Name: "test"},
				Body: &BlockStmt{
					List: []Stmt{
						&ExprStmt{X: &Word{Val: "echo"}},
					},
				},
			},
			expected: "*ast.FuncDecl",
		},
		{
			name: "If statement",
			node: &IfStmt{
				Cond: &BinaryExpr{
					X:  &Word{Val: "a"},
					Op: token.Equal,
					Y:  &Word{Val: "b"},
				},
				Body: &BlockStmt{
					List: []Stmt{
						&ExprStmt{X: &Word{Val: "echo"}},
					},
				},
			},
			expected: "*ast.IfStmt",
		},
		{
			name: "While statement",
			node: &WhileStmt{
				Cond: &Word{Val: "condition"},
				Body: &BlockStmt{
					List: []Stmt{
						&ExprStmt{X: &Word{Val: "echo"}},
					},
				},
			},
			expected: "*ast.WhileStmt",
		},
		{
			name: "Until statement",
			node: &UntilStmt{
				Cond: &Word{Val: "condition"},
				Body: &BlockStmt{
					List: []Stmt{
						&ExprStmt{X: &Word{Val: "echo"}},
					},
				},
			},
			expected: "*ast.UntilStmt",
		},
		{
			name: "For statement",
			node: &ForStmt{
				Init: &AssignStmt{
					Lhs: &Ident{Name: "i"},
					Rhs: &Word{Val: "0"},
				},
				Cond: &BinaryExpr{
					X:  &Ident{Name: "i"},
					Op: token.TsLss,
					Y:  &Word{Val: "10"},
				},
				Post: &AssignStmt{
					Lhs: &Ident{Name: "i"},
					Rhs: &BinaryExpr{
						X:  &Ident{Name: "i"},
						Op: token.Plus,
						Y:  &Word{Val: "1"},
					},
				},
				Body: &BlockStmt{
					List: []Stmt{
						&ExprStmt{X: &Word{Val: "echo"}},
					},
				},
			},
			expected: "*ast.ForStmt",
		},
		{
			name: "ForIn statement",
			node: &ForInStmt{
				Var:  &Ident{Name: "item"},
				List: &Word{Val: "1 2 3"},
				Body: &BlockStmt{
					List: []Stmt{
						&ExprStmt{X: &Word{Val: "echo"}},
					},
				},
			},
			expected: "*ast.ForInStmt",
		},
		{
			name: "Select statement",
			node: &SelectStmt{
				Var:  &Ident{Name: "choice"},
				List: &Word{Val: "option1 option2 option3"},
				Body: &BlockStmt{
					List: []Stmt{
						&ExprStmt{X: &Word{Val: "echo"}},
					},
				},
			},
			expected: "*ast.SelectStmt",
		},
		{
			name: "Case statement",
			node: &CaseStmt{
				X: &Ident{Name: "var"},
				Patterns: []*CaseClause{
					{
						Conds: []Expr{&Word{Val: "pattern1"}},
						Body: &BlockStmt{
							List: []Stmt{
								&ExprStmt{X: &Word{Val: "echo"}},
							},
						},
					},
				},
			},
			expected: "*ast.CaseStmt",
		},
		{
			name: "Assignment statement",
			node: &AssignStmt{
				Lhs: &Ident{Name: "var"},
				Rhs: &Word{Val: "value"},
			},
			expected: "*ast.AssignStmt",
		},
		{
			name: "Pipeline expression",
			node: &PipelineExpr{
				CtrOp: token.Or,
				Cmds: []Expr{
					&CmdExpr{Name: &Ident{Name: "cmd1"}},
					&CmdExpr{Name: &Ident{Name: "cmd2"}},
				},
			},
			expected: "*ast.PipelineExpr",
		},
		{
			name: "Redirect",
			node: &Redirect{
				N:      &Word{Val: "2"},
				CtrOp:  token.RdrOut,
				Word:   &Word{Val: "file.txt"},
			},
			expected: "*ast.Redirect",
		},
		{
			name: "Word",
			node: &Word{Val: "hello world"},
			expected: "*ast.Word",
		},
		{
			name: "Identifier",
			node: &Ident{Name: "variable"},
			expected: "*ast.Ident",
		},
		{
			name: "Binary expression",
			node: &BinaryExpr{
				X:  &Word{Val: "a"},
				Op: token.Equal,
				Y:  &Word{Val: "b"},
			},
			expected: "*ast.BinaryExpr",
		},
		{
			name: "Command expression",
			node: &CmdExpr{
				Name: &Ident{Name: "echo"},
				Recv: []Expr{
					&Word{Val: "hello"},
					&Word{Val: "world"},
				},
			},
			expected: "*ast.CmdExpr",
		},
		{
			name: "Test expression",
			node: &TestExpr{
				X: &BinaryExpr{
					X:  &Word{Val: "a"},
					Op: token.Equal,
					Y:  &Word{Val: "b"},
				},
			},
			expected: "*ast.TestExpr",
		},
		{
			name: "Extended test expression",
			node: &ExtendedTestExpr{
				X: &BinaryExpr{
					X:  &Word{Val: "a"},
					Op: token.Equal,
					Y:  &Word{Val: "b"},
				},
			},
			expected: "*ast.ExtendedTestExpr",
		},
		{
			name: "Arithmetic evaluation expression",
			node: &ArithEvalExpr{
				X: &BinaryExpr{
					X:  &Word{Val: "a"},
					Op: token.Plus,
					Y:  &Word{Val: "b"},
				},
			},
			expected: "*ast.ArithEvalExpr",
		},
		{
			name: "Command substitution",
			node: &CmdSubst{
				Tok: token.DollParen,
				X:   &CmdExpr{Name: &Ident{Name: "echo"}},
			},
			expected: "*ast.CmdSubst",
		},
		{
			name: "Process substitution",
			node: &ProcSubst{
				CtrOp: token.CmdIn,
				X:     &CmdExpr{Name: &Ident{Name: "echo"}},
			},
			expected: "*ast.ProcSubst",
		},
		{
			name: "Arithmetic expression",
			node: &ArithExpr{
				X: &BinaryExpr{
					X:  &Word{Val: "a"},
					Op: token.Plus,
					Y:  &Word{Val: "b"},
				},
			},
			expected: "*ast.ArithExpr",
		},
		{
			name: "Variable expansion expression",
			node: &VarExpExpr{
				X: &Ident{Name: "variable"},
			},
			expected: "*ast.VarExpExpr",
		},
		{
			name: "Parameter expansion expression",
			node: &ParamExpExpr{
				Var: &Ident{Name: "variable"},
			},
			expected: "*ast.ParamExpExpr",
		},
		{
			name: "Index expression",
			node: &IndexExpr{
				X: &Ident{Name: "array"},
				Y: &Word{Val: "0"},
			},
			expected: "*ast.IndexExpr",
		},
		{
			name: "Array expression",
			node: &ArrExpr{
				Vars: []Expr{
					&Word{Val: "item1"},
					&Word{Val: "item2"},
				},
			},
			expected: "*ast.ArrExpr",
		},
		{
			name: "Unary expression",
			node: &UnaryExpr{
				Op: token.ExclMark,
				X:  &Word{Val: "condition"},
			},
			expected: "*ast.UnaryExpr",
		},
		{
			name: "Command list expression",
			node: &CmdListExpr{
				CtrOp: token.AndAnd,
				Cmds: []Expr{
					&CmdExpr{Name: &Ident{Name: "cmd1"}},
					&CmdExpr{Name: &Ident{Name: "cmd2"}},
				},
			},
			expected: "*ast.CmdListExpr",
		},
		{
			name: "Command group",
			node: &CmdGroup{
				Op: token.LeftBrace,
				List: []Expr{
					&CmdExpr{Name: &Ident{Name: "cmd1"}},
					&CmdExpr{Name: &Ident{Name: "cmd2"}},
				},
			},
			expected: "*ast.CmdGroup",
		},
		{
			name: "Comment",
			node: &Comment{
				Text: " This is a comment",
			},
			expected: "*ast.Comment",
		},
		{
			name: "Comment group",
			node: &CommentGroup{
				List: []*Comment{
					{Text: " Comment 1"},
					{Text: " Comment 2"},
				},
			},
			expected: "*ast.CommentGroup",
		},
		{
			name: "Return statement",
			node: &ReturnStmt{
				X: &Word{Val: "0"},
			},
			expected: "*ast.ReturnStmt",
		},
		{
			name: "Block statement",
			node: &BlockStmt{
				List: []Stmt{
					&ExprStmt{X: &Word{Val: "echo"}},
					&ExprStmt{X: &Word{Val: "hello"}},
				},
			},
			expected: "*ast.BlockStmt",
		},
		{
			name: "Expression statement",
			node: &ExprStmt{
				X: &CmdExpr{Name: &Ident{Name: "echo"}},
			},
			expected: "*ast.ExprStmt",
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

func TestInspectComplexStructure(t *testing.T) {
	// Test a more complex AST structure
	node := &File{
		Stmts: []Stmt{
			&FuncDecl{
				Name: &Ident{Name: "main"},
				Body: &BlockStmt{
					List: []Stmt{
						&AssignStmt{
							Lhs: &Ident{Name: "var"},
							Rhs: &Word{Val: "value"},
						},
						&IfStmt{
							Cond: &BinaryExpr{
								X:  &Ident{Name: "var"},
								Op: token.Equal,
								Y:  &Word{Val: "test"},
							},
							Body: &BlockStmt{
								List: []Stmt{
									&ExprStmt{
										X: &CmdExpr{
											Name: &Ident{Name: "echo"},
											Recv: []Expr{
												&Word{Val: "Hello"},
												&Word{Val: "World"},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	var buf strings.Builder
	Inspect(node, &buf)
	output := buf.String()

	// Check that the output contains expected node types
	expectedTypes := []string{
		"*ast.File",
		"*ast.FuncDecl",
		"*ast.BlockStmt",
		"*ast.AssignStmt",
		"*ast.IfStmt",
		"*ast.BinaryExpr",
		"*ast.ExprStmt",
		"*ast.CmdExpr",
	}

	for _, expectedType := range expectedTypes {
		if !strings.Contains(output, expectedType) {
			t.Errorf("Inspect() output missing expected type %s", expectedType)
		}
	}
}
