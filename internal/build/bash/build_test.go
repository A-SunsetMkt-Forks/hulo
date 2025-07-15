// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package build_test

import (
	"fmt"
	"testing"

	"github.com/caarlos0/log"

	build "github.com/hulo-lang/hulo/internal/build/bash"
	"github.com/hulo-lang/hulo/internal/config"
	"github.com/hulo-lang/hulo/internal/vfs/memvfs"
	bast "github.com/hulo-lang/hulo/syntax/bash/ast"
	hast "github.com/hulo-lang/hulo/syntax/hulo/ast"
	"github.com/hulo-lang/hulo/syntax/hulo/parser"
	htok "github.com/hulo-lang/hulo/syntax/hulo/token"
	"github.com/stretchr/testify/assert"
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
	assert.NoError(t, err)
	bast.Print(node)
}

func TestSorceScriptBuild(t *testing.T) {
	log.SetLevel(log.ErrorLevel)
	script := `echo "Hello, World!" 3.14 true`
	node, err := parser.ParseSourceScript(script)
	assert.NoError(t, err)
	// hast.Print(node)
	bnode, err := build.Translate(&config.BashOptions{}, node)
	assert.NoError(t, err)
	bast.Print(bnode)
}

func TestBuildModule(t *testing.T) {
	fs := memvfs.New()
	testFiles := map[string]string{
		"std/unsafe/bash/constant.hl": `
declare fn echo(message: str);
		`,
		"my_script.hl": `
			pub fn greet() {
				echo("Hello, World!")
			}
		`,
		"main.hl": `
			import "my_script"

			my_script.greet()
		`,
	}

	for path, content := range testFiles {
		fs.WriteFile(path, []byte(content), 0644)
	}

	results, err := build.Transpile(&config.BashOptions{}, nil, fs, ".", ".", "main.hl")
	assert.NoError(t, err)

	for file, code := range results {
		fmt.Printf("=== %s ===\n", file)
		fmt.Println(code)
		fmt.Println()
	}
}

func TestConvertSelectExpr(t *testing.T) {
	fs := memvfs.New()
	testFiles := map[string]string{
		"math.hl": `
// Math module for Bash

const PI=3.14159

pub fn sqrt(x: num) -> num {
	return $x * $x
}`,
		"main.hl": `
			import "math"

			echo($PI)
			echo(math.sqrt(16))

			$p := "hello"
			echo($p)
			echo($p.length())
		`,
	}

	for path, content := range testFiles {
		fs.WriteFile(path, []byte(content), 0644)
	}

	results, err := build.Transpile(&config.BashOptions{}, nil, fs, ".", ".", "main.hl")
	assert.NoError(t, err)

	for file, code := range results {
		fmt.Printf("=== %s ===\n", file)
		fmt.Println(code)
		fmt.Println()
	}
}
