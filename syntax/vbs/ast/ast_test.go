// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package ast_test

import (
	"testing"

	"github.com/hulo-lang/hulo/syntax/vbs/ast"
	"github.com/hulo-lang/hulo/syntax/vbs/token"
	"github.com/stretchr/testify/assert"
)

func TestAST(t *testing.T) {
	testset := []struct {
		node     ast.Node
		expected string
	}{
		{
			&ast.File{
				Stmts: []ast.Stmt{
					&ast.DimDecl{List: []ast.Expr{&ast.Ident{Name: "A"}}},
					&ast.AssignStmt{
						Lhs: &ast.Ident{Name: "A"},
						Rhs: &ast.CallExpr{
							Func: &ast.Ident{Name: "Array"},
							Recv: []ast.Expr{&ast.Ident{Name: "10"}, &ast.Ident{Name: "20"}, &ast.Ident{Name: "30"}},
						},
					},
				},
			}, `Dim A
A = Array(10,20,30)`},
		{&ast.File{Stmts: []ast.Stmt{
			&ast.DimDecl{List: []ast.Expr{&ast.IndexExpr{X: &ast.Ident{Name: "Names"}, Index: &ast.Ident{Name: "9"}}}},
			&ast.DimDecl{List: []ast.Expr{&ast.IndexListExpr{X: &ast.Ident{Name: "Names"}, Indices: []ast.Expr{&ast.Ident{Name: "10"}, &ast.Ident{Name: "10"}, &ast.Ident{Name: "10"}}}}},
			&ast.DimDecl{List: []ast.Expr{&ast.Ident{Name: "MyVar"}, &ast.Ident{Name: "MyNum"}}}}}, `Dim Names(9)
Dim Names(10, 10, 10)
Dim MyVar, MyNum`},
		{&ast.BlockStmt{List: []ast.Stmt{
			&ast.OnErrorStmt{OnErrorResume: &ast.OnErrorResume{}},
			&ast.ExprStmt{
				X: &ast.CallExpr{
					Func: &ast.SelectorExpr{X: &ast.Ident{Name: "Err"}, Sel: &ast.Ident{Name: "Raise"}},
					Recv: []ast.Expr{&ast.BasicLit{Kind: token.INTEGER, Value: "6"}},
				}},
			&ast.ExprStmt{
				Doc: &ast.CommentGroup{List: []*ast.Comment{{Tok: token.SGL_QUOTE, Text: "Clear the error"}}},
				X: &ast.CallExpr{
					Func: &ast.Ident{Name: "MsgBox"},
					Recv: []ast.Expr{&ast.BinaryExpr{
						X:  &ast.Ident{Name: `"Error # "`},
						Op: token.AND,
						Y: &ast.BinaryExpr{
							X: &ast.CallExpr{
								Func: &ast.Ident{Name: "CStr"},
								Recv: []ast.Expr{&ast.SelectorExpr{X: &ast.Ident{Name: "Err"}, Sel: &ast.Ident{Name: "Number"}}},
							},
							Op: token.AND,
							Y: &ast.BinaryExpr{
								X:  &ast.BasicLit{Kind: token.STRING, Value: " "},
								Op: token.AND,
								Y: &ast.SelectorExpr{
									X:   &ast.Ident{Name: "Err"},
									Sel: &ast.Ident{Name: "Description"},
								},
							},
						}}}}},
		}}, `On Error Resume Next
Err.Raise 6   ' Raise an overflow error.
MsgBox ("Error # " & CStr(Err.Number) & " " & Err.Description)
Err.Clear      ' Clear the error`},
	}
	for _, tt := range testset {
		assert.Equal(t, tt.expected, tt.node)
	}
}

func TestPrint(t *testing.T) {
	ast.Print(&ast.File{
		Stmts: []ast.Stmt{
			&ast.DimDecl{
				List:  []ast.Expr{&ast.Ident{Name: "x"}},
				Colon: token.DynPos,
				Set: &ast.SetStmt{
					Lhs: &ast.Ident{Name: "x"},
					Rhs: &ast.NewExpr{
						X: &ast.CallExpr{
							Func: &ast.Ident{Name: "CreateObject"},
							Recv: []ast.Expr{&ast.BasicLit{Kind: token.STRING, Value: "Scripting.Dictionary"}},
						},
					},
				},
			},
			&ast.ReDimDecl{
				Preserve: token.DynPos,
				List: []ast.Expr{
					&ast.IndexListExpr{
						X:       &ast.Ident{Name: "X"},
						Indices: []ast.Expr{&ast.Ident{Name: "10"}, &ast.Ident{Name: "10"}, &ast.Ident{Name: "15"}},
					},
				},
			},
			&ast.ClassDecl{
				Mod:  token.PUBLIC,
				Name: &ast.Ident{Name: "RGB"},
				Stmts: []ast.Stmt{
					&ast.PrivateStmt{
						List: []ast.Expr{&ast.Ident{Name: "m_value"}},
					},
					&ast.PropertyGetStmt{
						Mod:  token.PUBLIC,
						Name: &ast.Ident{Name: "Value"},
						Recv: []*ast.Field{
							{Tok: token.BYREF, TokPos: token.DynPos, Name: &ast.Ident{Name: "val"}},
						},
						Body: &ast.BlockStmt{
							List: []ast.Stmt{
								&ast.IfStmt{
									Cond: &ast.BinaryExpr{
										X:  &ast.Ident{Name: "val"},
										Op: token.LTQ,
										Y:  &ast.BasicLit{Kind: token.INTEGER, Value: "0"},
									},
									Body: &ast.BlockStmt{
										List: []ast.Stmt{
											&ast.AssignStmt{
												Lhs: &ast.Ident{Name: "m_value"},
												Rhs: &ast.Ident{Name: "val"},
											},
										},
									},
									Else: &ast.BlockStmt{
										List: []ast.Stmt{
											&ast.ExprStmt{
												X: &ast.CallExpr{
													Func: &ast.SelectorExpr{
														X:   &ast.Ident{Name: "Err"},
														Sel: &ast.Ident{Name: "Raise"},
													},
													Recv: []ast.Expr{
														&ast.BasicLit{Kind: token.INTEGER, Value: "1001"},
														&ast.BasicLit{Kind: token.STRING, Value: "Color"},
														&ast.BasicLit{Kind: token.STRING, Value: "Invalid color value"},
													},
												},
											},
										},
									},
								},
							},
						},
					},
					&ast.FuncDecl{
						Mod:  token.PUBLIC,
						Name: &ast.Ident{Name: "color"},
						Body: &ast.BlockStmt{
							List: []ast.Stmt{
								&ast.IfStmt{
									Cond: &ast.BinaryExpr{X: &ast.Ident{Name: "X"}, Op: token.EQ, Y: &ast.Ident{Name: "10"}},
									Body: &ast.BlockStmt{
										List: []ast.Stmt{
											&ast.ExprStmt{X: &ast.CallExpr{Func: &ast.Ident{Name: "MsgBox"}, Recv: []ast.Expr{&ast.BasicLit{Kind: token.STRING, Value: "right"}}}},
											&ast.IfStmt{
												Cond: &ast.BinaryExpr{
													X: &ast.BinaryExpr{
														X:  &ast.Ident{Name: "a"},
														Op: token.LT,
														Y:  &ast.BasicLit{Kind: token.INTEGER, Value: "10"},
													},
													Op: token.AND,
													Y: &ast.BinaryExpr{
														X:  &ast.Ident{Name: "b"},
														Op: token.LT,
														Y:  &ast.BasicLit{Kind: token.INTEGER, Value: "10"},
													},
												},
												Body: &ast.BlockStmt{
													List: []ast.Stmt{
														&ast.ExprStmt{X: &ast.CallExpr{Func: &ast.Ident{Name: "MsgBox"}, Recv: []ast.Expr{&ast.BasicLit{Kind: token.STRING, Value: "right"}}}},
													},
												},
												Else: &ast.IfStmt{Cond: &ast.BinaryExpr{X: &ast.Ident{Name: "X"}, Op: token.GT, Y: &ast.BasicLit{Kind: token.INTEGER, Value: "17"}},
													Body: &ast.BlockStmt{List: []ast.Stmt{
														&ast.ExprStmt{X: &ast.CallExpr{Func: &ast.Ident{Name: "MsgBox"}, Recv: []ast.Expr{&ast.BasicLit{Kind: token.STRING, Value: "right"}}}},
													}},
													Else: &ast.IfStmt{
														Cond: &ast.BinaryExpr{X: &ast.Ident{Name: "X"}, Op: token.GT, Y: &ast.BasicLit{Kind: token.INTEGER, Value: "17"}},
														Body: &ast.BlockStmt{
															List: []ast.Stmt{
																&ast.ExprStmt{X: &ast.CallExpr{Func: &ast.Ident{Name: "MsgBox"}, Recv: []ast.Expr{&ast.BasicLit{Kind: token.STRING, Value: "right"}}}},
															}},
														Else: &ast.BlockStmt{List: []ast.Stmt{
															&ast.ExprStmt{X: &ast.CallExpr{Func: &ast.Ident{Name: "MsgBox"}, Recv: []ast.Expr{&ast.BasicLit{Kind: token.STRING, Value: "right"}}}},
														}},
													},
												},
											},
											&ast.ExprStmt{X: &ast.CallExpr{Func: &ast.Ident{Name: "MsgBox"}, Recv: []ast.Expr{&ast.BasicLit{Kind: token.STRING, Value: "right"}}}},
										},
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
