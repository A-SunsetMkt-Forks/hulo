// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package transpiler

import (
	"fmt"

	"github.com/hulo-lang/hulo/internal/config"
	bast "github.com/hulo-lang/hulo/syntax/batch/ast"
	btok "github.com/hulo-lang/hulo/syntax/batch/token"
	hast "github.com/hulo-lang/hulo/syntax/hulo/ast"
)

func Translate(opts *config.BatchOptions, node hast.Node) (bast.Node, error) {
	bnode := translate2Batch(opts, node)
	return bnode, nil
}

func translate2Batch(opts *config.BatchOptions, node hast.Node) bast.Node {
	switch node := node.(type) {
	case *hast.File:
		docs := make([]*bast.CommentGroup, len(node.Docs))
		for i, d := range node.Docs {
			docs[i] = translate2Batch(opts, d).(*bast.CommentGroup)
		}
		stmts := []bast.Stmt{&bast.ExprStmt{
			X: &bast.Word{
				Parts: []bast.Expr{
					&bast.Lit{
						Val: "@echo",
					},
					&bast.Lit{
						Val: "off",
					},
				},
			},
		}}
		for _, s := range node.Stmts {
			stmts = append(stmts, translate2Batch(opts, s).(bast.Stmt))
		}
		return &bast.File{
			Docs:  docs,
			Stmts: stmts,
		}
	case *hast.CommentGroup:
		var tok btok.Token
		if opts.CommentSyntax == "::" {
			tok = btok.DOUBLE_COLON
		} else {
			tok = btok.REM
		}
		docs := make([]*bast.Comment, len(node.List))
		for i, d := range node.List {
			docs[i] = &bast.Comment{
				Tok:  tok,
				Text: d.Text,
			}
		}
		return &bast.CommentGroup{
			Comments: docs,
		}
	case *hast.Comment:
		return &bast.Comment{
			Text: node.Text,
		}
	case *hast.ExprStmt:
		return &bast.ExprStmt{
			X: translate2Batch(opts, node.X).(bast.Expr),
		}
	case *hast.AssignStmt:
		lhs := translate2Batch(opts, node.Lhs).(bast.Expr)
		rhs := translate2Batch(opts, node.Rhs).(bast.Expr)

		// TODO 需要判断是不是可以计算的类型，如果是可以计算的类型的话，
		// 需要加上 /a 参数。
		return &bast.AssignStmt{
			Lhs: lhs,
			Rhs: rhs,
		}
	case *hast.RefExpr:
		return &bast.Lit{
			Val: node.X.(*hast.Ident).Name,
		}
	case *hast.CallExpr:
		recv := make([]bast.Expr, len(node.Recv))
		for i, r := range node.Recv {
			recv[i] = translate2Batch(opts, r).(bast.Expr)
		}
		return &bast.CallExpr{
			Fun:  translate2Batch(opts, node.Fun).(bast.Expr),
			Recv: recv,
		}
	case *hast.Ident:
		return &bast.Lit{
			Val: node.Name,
		}
	case *hast.BasicLit:
		return &bast.Lit{
			Val: node.Value,
		}
	case *hast.IfStmt:
		cond := translate2Batch(opts, node.Cond).(bast.Expr)
		body := translate2Batch(opts, node.Body).(bast.Stmt)
		var else_ bast.Stmt
		if node.Else != nil {
			else_ = translate2Batch(opts, node.Else).(bast.Stmt)
		}
		return &bast.IfStmt{
			Cond: cond,
			Body: body,
			Else: else_,
		}
	case *hast.BinaryExpr:
		lhs := translate2Batch(opts, node.X).(bast.Expr)
		rhs := translate2Batch(opts, node.Y).(bast.Expr)
		return &bast.BinaryExpr{
			X:  lhs,
			Op: Token(node.Op),
			Y:  rhs,
		}
	case *hast.BlockStmt:
		stmts := make([]bast.Stmt, len(node.List))
		for i, s := range node.List {
			stmts[i] = translate2Batch(opts, s).(bast.Stmt)
		}
		return &bast.BlockStmt{
			List: stmts,
		}
	default:
		fmt.Printf("unsupported node type: %T\n", node)
	}
	return nil
}
