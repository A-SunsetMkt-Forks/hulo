// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package build

import (
	"fmt"

	"github.com/hulo-lang/hulo/internal/config"
	bast "github.com/hulo-lang/hulo/syntax/bash/ast"
	hast "github.com/hulo-lang/hulo/syntax/hulo/ast"
)

func Translate(opts *config.BashOptions, node hast.Node) (bast.Node, error) {
	bnode := translate2Bash(opts, node)
	return bnode, nil
}

func translate2Bash(opts *config.BashOptions, node hast.Node) bast.Node {
	switch node := node.(type) {
	case *hast.File:
		stmts := []bast.Stmt{
			&bast.Comment{
				Text: "!/bin/bash",
			},
		}
		for _, s := range node.Stmts {
			stmts = append(stmts, translate2Bash(opts, s).(bast.Stmt))
		}

		return &bast.File{
			Stmts: stmts,
		}
	case *hast.IfStmt:
		// opts.Boolean = "number & >= 1.2.3"
		parse, err := hcrDispatcher.Get(opts.BooleanFormat)
		if err != nil {
			return nil
		}
		cond, err := parse.Apply(&hast.IfStmt{}, node.Cond)
		if err != nil {
			return nil
		}
		return &bast.IfStmt{
			Cond: cond.(bast.Expr),
			Body: translate2Bash(opts, node.Body).(*bast.BlockStmt),
		}
	case *hast.BlockStmt:
		return &bast.BlockStmt{}
	case *hast.BinaryExpr:
		return &bast.BinaryExpr{
			X:  translate2Bash(opts, node.X).(bast.Expr),
			Op: Token(node.Op),
			Y:  translate2Bash(opts, node.Y).(bast.Expr),
		}
	case *hast.RefExpr:
		return &bast.VarExpExpr{
			X: translate2Bash(opts, node.X).(*bast.Ident),
		}
	case *hast.Ident:
		return &bast.Ident{
			Name: node.Name,
		}
	case *hast.BasicLit:
		return &bast.Word{
			Val: node.Value,
		}
	case *hast.ExprStmt:
		return &bast.ExprStmt{
			X: translate2Bash(opts, node.X).(bast.Expr),
		}
	case *hast.CallExpr:
		fun := translate2Bash(opts, node.Fun).(*bast.Ident)

		recv := []bast.Expr{}
		for _, r := range node.Recv {
			recv = append(recv, translate2Bash(opts, r).(bast.Expr))
		}

		return &bast.CmdExpr{
			Name: fun,
			Recv: recv,
		}
	default:
		fmt.Printf("%T\n", node)
	}

	return nil
}
