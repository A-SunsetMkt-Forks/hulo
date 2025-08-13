// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package astutil

import (
	"github.com/hulo-lang/hulo/syntax/batch/ast"
	"github.com/hulo-lang/hulo/syntax/batch/token"
)

func Goto(label string) *ast.GotoStmt {
	return &ast.GotoStmt{Label: label}
}

func Label(label string) *ast.LabelStmt {
	return &ast.LabelStmt{Name: label}
}

func Block(stmts ...ast.Stmt) *ast.BlockStmt {
	return &ast.BlockStmt{List: stmts}
}

func Word(parts ...ast.Expr) *ast.Word {
	return &ast.Word{Parts: parts}
}

func Unary(op token.Token, x ast.Expr) *ast.UnaryExpr {
	return &ast.UnaryExpr{Op: op, X: x}
}

func Binary(x ast.Expr, op token.Token, y ast.Expr) *ast.BinaryExpr {
	return &ast.BinaryExpr{X: x, Op: op, Y: y}
}

func Call(fun ast.Expr, recv ...ast.Expr) *ast.CallExpr {
	return &ast.CallExpr{Fun: fun, Recv: recv}
}

func Lit(val string) *ast.Lit {
	return &ast.Lit{Val: val}
}

func If(cond ast.Expr, body ast.Stmt, elseStmt ast.Stmt) *ast.IfStmt {
	return &ast.IfStmt{Cond: cond, Body: body, Else: elseStmt}
}

func For(x ast.Expr, list ast.Expr, body *ast.BlockStmt) *ast.ForStmt {
	return &ast.ForStmt{X: x, List: list, Body: body}
}

func FuncDecl(name string, body *ast.BlockStmt) *ast.FuncDecl {
	return &ast.FuncDecl{Name: name, Body: body}
}

func Assign(lhs, rhs ast.Expr) *ast.AssignStmt {
	return &ast.AssignStmt{Lhs: lhs, Rhs: rhs}
}

func CommentNode(cmt string, others ...string) ast.Node {
	if len(others) > 0 {
		comments := make([]*ast.Comment, len(others)+1)
		comments[0] = &ast.Comment{Text: cmt}
		for i, text := range others {
			comments[i+1] = &ast.Comment{Text: text}
		}
		return &ast.CommentGroup{Comments: comments}
	}
	return &ast.Comment{Text: cmt}
}

func Comment(cmt string) *ast.Comment {
	return &ast.Comment{Text: cmt}
}

func CommentGroup(cmts ...*ast.Comment) *ast.CommentGroup {
	return &ast.CommentGroup{Comments: cmts}
}

func ExprStmt(x ast.Expr) *ast.ExprStmt {
	return &ast.ExprStmt{X: x}
}

func CallExpr(name string, recv ...ast.Expr) *ast.CallExpr {
	return &ast.CallExpr{Fun: &ast.Lit{Val: name}, Recv: recv}
}

func SglQuote(val ast.Expr) *ast.SglQuote {
	return &ast.SglQuote{Val: val}
}

func DblQuote(val ast.Expr) *ast.DblQuote {
	return &ast.DblQuote{Val: val}
}

func CmdExpr(name string, recv ...ast.Expr) *ast.CmdExpr {
	return &ast.CmdExpr{Name: &ast.Lit{Val: name}, Recv: recv}
}

func File(docs []*ast.CommentGroup, stmts ...ast.Stmt) *ast.File {
	return &ast.File{Docs: docs, Stmts: stmts}
}
