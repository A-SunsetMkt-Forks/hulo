// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package build_test

import (
	"os"
	"testing"

	"github.com/caarlos0/log"
	build "github.com/hulo-lang/hulo/internal/build/vbs"
	"github.com/hulo-lang/hulo/internal/config"
	"github.com/hulo-lang/hulo/syntax/hulo/parser"
	vast "github.com/hulo-lang/hulo/syntax/vbs/ast"
	"github.com/stretchr/testify/assert"
)

func TestCommandStmt(t *testing.T) {
	log.SetLevel(log.ErrorLevel)
	script := `echo "Hello, World!" 3.14 true`
	node, err := parser.ParseSourceScript(script)
	assert.NoError(t, err)
	// hast.Print(node)
	bnode, err := build.TranspileToVBScript(&config.VBScriptOptions{}, node)
	assert.NoError(t, err)
	vast.Print(bnode)
}

func TestComment(t *testing.T) {
	log.SetLevel(log.ErrorLevel)
	script := `// this is a comment

	/*
	 * this
	 *    is a multi
	 * comment
	 */

	/**
	 * this
	 *    is a multi
	 * comment
	 */`
	node, err := parser.ParseSourceScript(script)
	assert.NoError(t, err)
	// hast.Print(node)
	bnode, err := build.TranspileToVBScript(&config.VBScriptOptions{
		CommentSyntax: "'",
	}, node)
	assert.NoError(t, err)
	vast.Print(bnode)
}

func TestAssign(t *testing.T) {
	log.SetLevel(log.ErrorLevel)
	script := `let a = 10
	var b = 3.14
	const c = "Hello, World!"
	$d := true`
	node, err := parser.ParseSourceScript(script)
	assert.NoError(t, err)
	// hast.Print(node)
	bnode, err := build.TranspileToVBScript(&config.VBScriptOptions{}, node)
	assert.NoError(t, err)
	vast.Print(bnode)
}

func TestIf(t *testing.T) {
	log.SetLevel(log.ErrorLevel)
	script := `if $a > 10 {
		echo "a is greater than 10"
	} else {
		echo "a is less than or equal to 10"
	}`
	node, err := parser.ParseSourceScript(script)
	assert.NoError(t, err)
	// hast.Print(node)
	bnode, err := build.TranspileToVBScript(&config.VBScriptOptions{}, node)
	assert.NoError(t, err)
	vast.Print(bnode)
}

func TestWhile(t *testing.T) {
	log.SetLevel(log.ErrorLevel)
	script := `loop {
		echo "Hello, World!"
	}

	loop $a == true {
		echo "Hello, World!"
	}

	do {
		echo "Hello, World!"
	} loop ($a > 10)`
	node, err := parser.ParseSourceScript(script)
	assert.NoError(t, err)
	// hast.Print(node)
	bnode, err := build.TranspileToVBScript(&config.VBScriptOptions{}, node)
	assert.NoError(t, err)
	vast.Print(bnode)
}

// TODO j-- error
func TestFor(t *testing.T) {
	// log.SetLevel(log.ErrorLevel)
	script := `
	loop $i := 0; $i < 10; $i++ {
		echo "Hello, World!"
		loop $j := 10; $j > 0; $j-- {
			echo "Hello, World!"
		}
	}`
	node, err := parser.ParseSourceScript(script, parser.OptionTracerASTTree(os.Stdout))
	assert.NoError(t, err)
	// hast.Print(node)
	bnode, err := build.TranspileToVBScript(&config.VBScriptOptions{}, node)
	assert.NoError(t, err)
	vast.Print(bnode)
}

func TestMatch(t *testing.T) {
}
