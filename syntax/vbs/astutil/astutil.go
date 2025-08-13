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
	return &ast.AssignStmt{
		Lhs: lhs,
		Rhs: rhs,
	}
}
