// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package ast_test

import (
	"testing"

	"github.com/hulo-lang/hulo/syntax/bash/ast"
	"github.com/hulo-lang/hulo/syntax/bash/token"
	"github.com/stretchr/testify/assert"
)

var NotNullPos = token.Pos(1)

func TestAssignStmt(t *testing.T) {
	actual := ast.String(&ast.AssignStmt{
		Lhs: ast.Identifier("count"),
		Rhs: ast.Literal("0"),
	})
	assert.Equal(t, "count=0\n", actual)

	actual = ast.String(&ast.AssignStmt{
		Local: NotNullPos,
		Lhs:   ast.Identifier("count"),
		Rhs:   ast.Literal("0"),
	})
	assert.Equal(t, "local count=0\n", actual)
}

func TestIfStmt(t *testing.T) {
	nestedIf := &ast.IfStmt{
		Cond: &ast.TestExpr{
			X: &ast.BinaryExpr{
				X:  &ast.VarExpExpr{X: ast.Identifier("count")},
				Op: token.Assgn,
				Y:  ast.Literal(`""`),
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.AssignStmt{
					Lhs: &ast.VarExpExpr{X: &ast.Ident{Name: "count"}},
					Rhs: ast.Literal("0"),
				},
			},
		},
	}

	rootIf := &ast.IfStmt{
		Cond: &ast.ExtendedTestExpr{
			X: &ast.BinaryExpr{
				X:  &ast.VarExpExpr{X: &ast.Ident{Name: "count"}},
				Op: token.Assgn,
				Y:  ast.Literal(`""`),
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.AssignStmt{
					Lhs: &ast.VarExpExpr{X: &ast.Ident{Name: "count"}},
					Rhs: ast.Literal("0"),
				},
				nestedIf,
			},
		},
		Else: nestedIf,
	}

	ast.Print(rootIf)
}

func TestFuncDecl(t *testing.T) {
	ast.Print(&ast.FuncDecl{
		Name: ast.Identifier("myecho"),
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.ExprStmt{
					X: ast.CallExpression("echo", ast.Literal(`"Hello, World!"`)),
				},
			},
		},
	})
}

func TestForStmt(t *testing.T) {
	ast.Print(&ast.ForStmt{
		Init: &ast.AssignStmt{
			Lhs: ast.Identifier("i"),
			Rhs: ast.Literal("0"),
		},
		Cond: &ast.BinaryExpr{
			X:  ast.Identifier("i"),
			Op: token.TsLss,
			Y:  ast.Literal("10"),
		},
		Post: &ast.AssignStmt{
			Lhs: ast.Identifier("i"),
			Rhs: ast.BinaryExpression(ast.Identifier("i"), token.Plus, ast.Literal("1")),
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.ExprStmt{X: &ast.CallExpr{
					Func: &ast.Ident{Name: "echo"},
					Recv: []ast.Expr{&ast.Ident{Name: "i"}},
				}},
			},
		},
	})

	// infinite loop
	ast.Print(&ast.ForStmt{
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.ExprStmt{X: &ast.CallExpr{
					Func: &ast.Ident{Name: "read"},
					Recv: []ast.Expr{&ast.Ident{Name: "var"}},
				}},
				&ast.IfStmt{
					Cond: &ast.TestExpr{
						X: &ast.BinaryExpr{
							// TODO "$var"
							X:  &ast.VarExpExpr{X: &ast.Ident{Name: "var"}},
							Op: token.Assgn,
							Y:  ast.Literal(`"."`),
						},
					},
					Body: &ast.BlockStmt{
						List: []ast.Stmt{&ast.BreakStmt{}},
					},
				},
			},
		},
	})
}

func TestForInStmt(t *testing.T) {
	ast.Print(&ast.ForInStmt{
		Var:  &ast.Ident{Name: "i"},
		List: &ast.Ident{Name: "*.png"},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.ExprStmt{
					X: &ast.CallExpr{
						Func: &ast.Ident{Name: "ls"},
						Recv: []ast.Expr{&ast.Ident{Name: "-l"}, &ast.VarExpExpr{X: &ast.Ident{Name: "i"}}},
					},
				},
			},
		},
	})
}

func TestWhileStmt(t *testing.T) {
	ast.Print(&ast.WhileStmt{
		Cond: &ast.TestExpr{
			X: &ast.BinaryExpr{
				X:  &ast.VarExpExpr{X: &ast.Ident{Name: "number"}},
				Op: token.TsLss,
				Y:  &ast.Lit{Val: "10"},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.AssignStmt{
					Lhs: &ast.Ident{Name: "number"},
					Rhs: &ast.ArithExpr{
						X: &ast.BinaryExpr{
							X:  &ast.VarExpExpr{X: &ast.Ident{Name: "number"}},
							Op: token.Plus,
							Y:  &ast.Lit{Val: "1"},
						},
					},
				},
			},
		},
	})
}

func TestStmt(t *testing.T) {
	ast.Print(&ast.BlockStmt{
		List: []ast.Stmt{
			&ast.IfStmt{
				Cond: &ast.ArithEvalExpr{
					X: &ast.UnaryExpr{
						Op: token.TsDirect,
						X:  &ast.Ident{Name: "file.txt"},
					},
				},
				Body: &ast.BlockStmt{
					List: []ast.Stmt{
						&ast.ExprStmt{
							X: &ast.CallExpr{
								Func: &ast.Ident{Name: "echo"},
								Recv: []ast.Expr{},
							},
						},
					},
				},
			},
			&ast.AssignStmt{
				Lhs: &ast.Ident{Name: "reversed"},
				Rhs: &ast.CmdSubst{
					X: &ast.BinaryExpr{
						X: &ast.CallExpr{
							Func: &ast.Ident{Name: "echo"},
							Recv: []ast.Expr{&ast.Ident{Name: "-e"}, ast.Literal(`"${string}"`)},
						},
						Op: token.Or,
						Y: &ast.CallExpr{
							Func: &ast.Ident{Name: "rev"},
						},
					},
				},
			},
		},
	})
}

func TestPrint(t *testing.T) {
	ast.Print(&ast.File{
		Stmts: []ast.Stmt{
			&ast.FuncDecl{
				Name: &ast.Ident{Name: "scan"},
				Body: &ast.BlockStmt{
					List: []ast.Stmt{
						&ast.CaseStmt{
							X: &ast.Ident{Name: "$ANIMAL"},
							Patterns: []*ast.CaseClause{
								{Conds: []ast.Expr{&ast.Ident{Name: "cat"}, &ast.Ident{Name: "horse"}},
									Body: &ast.BlockStmt{
										List: []ast.Stmt{
											&ast.ExprStmt{
												&ast.CallExpr{
													Func: &ast.Ident{Name: "echo"},
													Recv: []ast.Expr{ast.Literal(`"string"`)},
												},
											},
										},
									}},
							},
							Else: &ast.BlockStmt{
								List: []ast.Stmt{
									&ast.ExprStmt{
										&ast.CallExpr{
											Func: &ast.Ident{Name: "echo"},
											Recv: []ast.Expr{ast.Literal(`"string"`)},
										},
									},
								},
							},
						},
					},
				},
			},
			&ast.IfStmt{
				Cond: &ast.ExtendedTestExpr{
					X: &ast.UnaryExpr{
						Op: token.TsDirect,
						X:  &ast.Ident{Name: "test.txt"},
					},
				},
				Body: &ast.BlockStmt{
					List: []ast.Stmt{
						&ast.ExprStmt{
							X: &ast.CallExpr{
								Func: &ast.Ident{Name: "echo"},
								Recv: []ast.Expr{},
							},
						},
					},
				},
				Else: &ast.IfStmt{Cond: &ast.UnaryExpr{
					Op: token.ExclMark,
					X: &ast.ExtendedTestExpr{
						X: &ast.BinaryExpr{
							X:  ast.Literal(`"$number"`),
							Op: token.NotEqual,
							Y:  &ast.Ident{Name: "^[0-9]+$"},
						},
					},
				}, Body: &ast.BlockStmt{
					List: []ast.Stmt{
						&ast.ExprStmt{
							X: &ast.CallExpr{
								Func: &ast.Ident{Name: "echo"},
								Recv: []ast.Expr{ast.Literal(`"input is invalid"`)},
							},
						},
					},
				}, Else: &ast.IfStmt{Cond: &ast.UnaryExpr{
					Op: token.ExclMark,
					X: &ast.ExtendedTestExpr{
						X: &ast.BinaryExpr{
							X:  ast.Literal(`"$number"`),
							Op: token.NotEqual,
							Y:  &ast.Ident{Name: "^[0-9]+$"},
						},
					},
				}, Body: &ast.BlockStmt{
					List: []ast.Stmt{
						&ast.ExprStmt{
							X: &ast.CallExpr{
								Func: &ast.Ident{Name: "echo"},
								Recv: []ast.Expr{ast.Literal(`"input is invalid"`)},
							},
						},
					},
				}, Else: &ast.BlockStmt{
					List: []ast.Stmt{
						&ast.ExprStmt{
							X: &ast.CallExpr{
								Func: &ast.Ident{Name: "echo"},
							},
						},
					},
				},
				},
				},
			},
		},
	})
}
