// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package astutil

import (
	"github.com/hulo-lang/hulo/syntax/vbs/ast"
	"github.com/hulo-lang/hulo/syntax/vbs/token"
)

func Comment(text string, tok ...token.Token) *ast.Comment {
	defaultTok := token.REM
	if len(tok) > 0 {
		defaultTok = tok[0]
	}
	return &ast.Comment{Tok: defaultTok, Text: text}
}

func CommentGroup(texts ...string) *ast.CommentGroup {
	comments := make([]*ast.Comment, len(texts))
	for i, text := range texts {
		comments[i] = Comment(text)
	}
	return &ast.CommentGroup{List: comments}
}

func OptionStmt() *ast.OptionStmt {
	return &ast.OptionStmt{}
}

func AssignStmt(lhs ast.Expr, rhs ast.Expr) *ast.AssignStmt {
	return &ast.AssignStmt{Lhs: lhs, Rhs: rhs}
}

func BasicLit(s string) *ast.BasicLit {
	return &ast.BasicLit{Value: s}
}

func Ident(name string) *ast.Ident {
	return &ast.Ident{Name: name}
}

func IndexExpr(x ast.Expr, index ast.Expr) *ast.IndexExpr {
	return &ast.IndexExpr{X: x, Index: index}
}

func IndexListExpr(x ast.Expr, indices ...ast.Expr) *ast.IndexListExpr {
	return &ast.IndexListExpr{X: x, Indices: indices}
}

func NewExpr(x ast.Expr) *ast.NewExpr {
	return &ast.NewExpr{X: x}
}

func CallExpr(fun ast.Expr, recv ...ast.Expr) *ast.CallExpr {
	return &ast.CallExpr{Func: fun, Recv: recv}
}

func CmdExpr(cmd ast.Expr, recv ...ast.Expr) *ast.CmdExpr {
	return &ast.CmdExpr{Cmd: cmd, Recv: recv}
}

func SelectorExpr(x ast.Expr, sel ast.Expr) *ast.SelectorExpr {
	return &ast.SelectorExpr{X: x, Sel: sel}
}

func BinaryExpr(x ast.Expr, op token.Token, y ast.Expr) *ast.BinaryExpr {
	return &ast.BinaryExpr{X: x, Op: op, Y: y}
}
