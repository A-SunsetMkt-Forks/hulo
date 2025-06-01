// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package build_test

import (
	"testing"

	build "github.com/hulo-lang/hulo/internal/build/bash"
	"github.com/hulo-lang/hulo/internal/config"
	bast "github.com/hulo-lang/hulo/syntax/bash/ast"
	hast "github.com/hulo-lang/hulo/syntax/hulo/ast"
	htok "github.com/hulo-lang/hulo/syntax/hulo/token"
)

func TestBuild(t *testing.T) {
	node, err := build.Translate(&config.BashOptions{}, &hast.File{
		Stmts: []hast.Stmt{
			&hast.IfStmt{
				Cond: &hast.BinaryExpr{
					X:  &hast.Ident{Name: "x"},
					Op: htok.EQ,
					Y: &hast.BasicLit{
						Kind:  htok.NUM,
						Value: "10",
					},
				},
				Body: &hast.BlockStmt{
					List: []hast.Stmt{
						&hast.ExprStmt{
							X: &hast.CallExpr{
								Fun: &hast.Ident{Name: "echo"},
								Recv: []hast.Expr{
									&hast.BasicLit{Kind: htok.STR, Value: "Hello, World!"},
								},
							},
						},
					},
				},
			},
		}})
	if err != nil {
		t.Fatal(err)
	}
	bast.Print(node)
}
