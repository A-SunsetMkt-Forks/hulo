// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package build_test

import (
	"testing"

	"github.com/caarlos0/log"
	build "github.com/hulo-lang/hulo/internal/build/bash"
	"github.com/hulo-lang/hulo/internal/config"
	bast "github.com/hulo-lang/hulo/syntax/bash/ast"
	hast "github.com/hulo-lang/hulo/syntax/hulo/ast"
	"github.com/hulo-lang/hulo/syntax/hulo/parser"
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

func TestSorceScriptBuild(t *testing.T) {
	log.SetLevel(log.ErrorLevel)
	script := `echo "Hello, World!" 3.14 true`
	node, err := parser.ParseSourceScript(script, parser.ParseOptions{})
	if err != nil {
		t.Fatal(err)
	}
	// hast.Print(node)
	bnode, err := build.Translate(&config.BashOptions{}, node)
	if err != nil {
		t.Fatal(err)
	}
	bast.Print(bnode)
}
