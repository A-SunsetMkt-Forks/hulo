// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package ast_test

import (
	"testing"

	"github.com/hulo-lang/hulo/grammar/fish/ast"
	"github.com/hulo-lang/hulo/grammar/fish/token"
)

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
						&ast.SwitchStmt{
							Var: &ast.Ident{Name: "$ANIMAL"},
							Cases: []*ast.CaseClause{
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
				Elif: []*ast.IfStmt{
					{Cond: &ast.BinaryExpr{
						X:  &ast.Ident{Name: "!"},
						Op: token.NONE,
						Y: &ast.ExtendedTestExpr{
							X: &ast.BinaryExpr{
								X:  &ast.BasicLit{Kind: token.STRING, Value: "$number"},
								Op: "=~",
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
					}},
					{Cond: &ast.BinaryExpr{
						X:  &ast.Ident{Name: "!"},
						Op: token.NONE,
						Y: &ast.ExtendedTestExpr{
							X: &ast.BinaryExpr{
								X:  &ast.BasicLit{Kind: token.STRING, Value: "$number"},
								Op: "=~",
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
					}},
				},
				Else: &ast.BlockStmt{
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
	})
}