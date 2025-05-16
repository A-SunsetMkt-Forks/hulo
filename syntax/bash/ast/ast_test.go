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
		Lhs: &ast.Ident{Name: "count"},
		Rhs: &ast.BasicLit{Kind: token.NUMBER, Value: "0"},
	})
	assert.Equal(t, "count=0", actual)

	actual = ast.String(&ast.AssignStmt{
		Local: NotNullPos,
		Lhs:   &ast.Ident{Name: "count"},
		Rhs:   &ast.BasicLit{Kind: token.NUMBER, Value: "0"},
	})
	assert.Equal(t, "local count=0", actual)
}

func TestIfStmt(t *testing.T) {
	nestedIf := &ast.IfStmt{
		Cond: &ast.TestExpr{
			X: &ast.BinaryExpr{
				X:  &ast.VarExpExpr{X: &ast.Ident{Name: "count"}},
				Op: token.ASSIGN,
				Y:  &ast.BasicLit{Kind: token.STRING, Value: ""},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.AssignStmt{
					Lhs: &ast.VarExpExpr{X: &ast.Ident{Name: "count"}},
					Rhs: &ast.BasicLit{Kind: token.NUMBER, Value: "0"},
				},
			},
		},
	}

	rootIf := &ast.IfStmt{
		Cond: &ast.ExtendedTestExpr{
			X: &ast.BinaryExpr{
				X:  &ast.VarExpExpr{X: &ast.Ident{Name: "count"}},
				Op: token.ASSIGN,
				Y:  &ast.BasicLit{Kind: token.STRING, Value: ""},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.AssignStmt{
					Lhs: &ast.VarExpExpr{X: &ast.Ident{Name: "count"}},
					Rhs: &ast.BasicLit{Kind: token.NUMBER, Value: "0"},
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
		Name: &ast.Ident{Name: "myecho"},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.ExprStmt{
					X: &ast.CallExpr{
						Func: &ast.Ident{Name: "echo"},
						Recv: []ast.Expr{&ast.BasicLit{Kind: token.STRING, Value: "Hello, World!"}},
					},
				},
			},
		},
	})
}

func TestForStmt(t *testing.T) {
	ast.Print(&ast.ForStmt{
		Init: &ast.AssignStmt{
			Lhs: &ast.Ident{Name: "i"},
			Rhs: &ast.BasicLit{Kind: token.NUMBER, Value: "0"},
		},
		Cond: &ast.BinaryExpr{
			X:  &ast.Ident{Name: "i"},
			Op: token.LT,
			Y:  &ast.BasicLit{Kind: token.NUMBER, Value: "10"},
		},
		Post: &ast.AssignStmt{
			Lhs: &ast.Ident{Name: "i"},
			Rhs: &ast.BinaryExpr{X: &ast.Ident{Name: "i"}, Op: token.ADD, Y: &ast.BasicLit{Kind: token.NUMBER, Value: "1"}},
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
							Op: token.ASSIGN,
							Y:  &ast.BasicLit{Kind: token.STRING, Value: "."},
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
				Op: token.LT_LIT,
				Y:  &ast.BasicLit{Kind: token.NUMBER, Value: "10"},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.AssignStmt{
					Lhs: &ast.Ident{Name: "number"},
					Rhs: &ast.ArithExpr{
						X: &ast.BinaryExpr{
							X:  &ast.VarExpExpr{X: &ast.Ident{Name: "number"}},
							Op: token.ADD,
							Y:  &ast.BasicLit{Kind: token.NUMBER, Value: "1"},
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
					X: &ast.BinaryExpr{
						X:  &ast.Ident{Name: "-d"},
						Op: token.NONE,
						Y:  &ast.Ident{Name: "file.txt"},
					},
				},
				Body: &ast.BlockStmt{
					Tok: token.NONE,
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
							Recv: []ast.Expr{&ast.Ident{Name: "-e"}, &ast.BasicLit{Kind: token.STRING, Value: "${string}"}},
						},
						Op: token.BITOR,
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
		Decls: []ast.Decl{
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
													Recv: []ast.Expr{&ast.BasicLit{Kind: token.STRING, Value: "string"}},
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
											Recv: []ast.Expr{&ast.BasicLit{Kind: token.STRING, Value: "string"}},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		Stmts: []ast.Stmt{
			&ast.IfStmt{
				Cond: &ast.ExtendedTestExpr{
					X: &ast.BinaryExpr{
						X:  &ast.Ident{Name: "-d"},
						Op: token.NONE,
						Y:  &ast.Ident{Name: "test.txt"},
					},
				},
				Body: &ast.BlockStmt{
					Tok: token.NONE,
					List: []ast.Stmt{
						&ast.ExprStmt{
							X: &ast.CallExpr{
								Func: &ast.Ident{Name: "echo"},
								Recv: []ast.Expr{},
							},
						},
					},
				},
				Else: &ast.IfStmt{Cond: &ast.BinaryExpr{
					X:  &ast.Ident{Name: "!"},
					Op: token.NONE,
					Y: &ast.ExtendedTestExpr{
						X: &ast.BinaryExpr{
							X:  &ast.BasicLit{Kind: token.STRING, Value: "$number"},
							Op: token.ASSIGN_BITNEG,
							Y:  &ast.Ident{Name: "^[0-9]+$"},
						},
					},
				}, Body: &ast.BlockStmt{
					List: []ast.Stmt{
						&ast.ExprStmt{
							X: &ast.CallExpr{
								Func: &ast.Ident{Name: "echo"},
								Recv: []ast.Expr{&ast.BasicLit{Kind: token.STRING, Value: "input is invalid"}},
							},
						},
					},
				}, Else: &ast.IfStmt{Cond: &ast.BinaryExpr{
					X:  &ast.Ident{Name: "!"},
					Op: token.NONE,
					Y: &ast.ExtendedTestExpr{
						X: &ast.BinaryExpr{
							X:  &ast.BasicLit{Kind: token.STRING, Value: "$number"},
							Op: token.ASSIGN_BITNEG,
							Y:  &ast.Ident{Name: "^[0-9]+$"},
						},
					},
				}, Body: &ast.BlockStmt{
					List: []ast.Stmt{
						&ast.ExprStmt{
							X: &ast.CallExpr{
								Func: &ast.Ident{Name: "echo"},
								Recv: []ast.Expr{&ast.BasicLit{Kind: token.STRING, Value: "input is invalid"}},
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
