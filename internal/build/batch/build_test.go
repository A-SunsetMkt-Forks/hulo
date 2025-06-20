// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package build_test

import (
	"testing"

	"github.com/caarlos0/log"
	build "github.com/hulo-lang/hulo/internal/build/batch"
	"github.com/hulo-lang/hulo/internal/config"
	bast "github.com/hulo-lang/hulo/syntax/batch/ast"
	"github.com/hulo-lang/hulo/syntax/hulo/parser"
	"github.com/stretchr/testify/assert"
)

func TestCommandStmt(t *testing.T) {
	log.SetLevel(log.ErrorLevel)
	script := `echo "Hello, World!" 3.14 true`
	node, err := parser.ParseSourceScript(script) // parser.OptionDebuggerDisableTiming(),
	// parser.OptionDebuggerIgnore("MulDivExpression", "AddSubExpression", "ShiftExpression", "LogicalExpression", "ConditionalExpression"),
	// parser.OptionDebuggerWatchNode("CommandExpression",
	// 	func(node hast.Node, pos parser.Position) {
	// 		fmt.Printf("Entering CommandExpression at %d:%d\n", pos.Line, pos.Column)
	// 	},
	// 	func(node hast.Node, result any, err error) {
	// 		fmt.Printf("Exiting CommandExpression with result: %v\n", result)
	// 	},
	// ),
	// parser.OptionDisplayASTTree(os.Stdout),

	assert.NoError(t, err)
	// hast.Print(node)
	bnode, err := build.Translate(&config.BatchOptions{}, node)
	assert.NoError(t, err)
	bast.Print(bnode)
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
	bnode, err := build.Translate(&config.BatchOptions{
		CommentSyntax: "::",
	}, node)
	assert.NoError(t, err)
	bast.Print(bnode)
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
	bnode, err := build.Translate(&config.BatchOptions{}, node)
	assert.NoError(t, err)
	bast.Print(bnode)
}

func TestIf(t *testing.T) {
	log.SetLevel(log.ErrorLevel)
	script := `
	$a := 20
	if $a > 10 {
		echo "a is greater than 10"
	} else {
		echo "a is less than or equal to 10"
	}`
	node, err := parser.ParseSourceScript(script)
	assert.NoError(t, err)
	// hast.Print(node)
	bnode, err := build.Translate(&config.BatchOptions{}, node)
	assert.NoError(t, err)
	bast.Print(bnode)
}

func TestLoop(t *testing.T) {
	// log.SetLevel(log.ErrorLevel)
	script := `loop {
		echo "Hello, World!"
	}

	do {
		echo "Hello, World!"
	} loop ($a > 10)

	loop $i := 0; $i < 10; $i++ {
		echo "Hello, World!"
	}`
	node, err := parser.ParseSourceScript(script)
	assert.NoError(t, err)
	// hast.Print(node)
	bnode, err := build.Translate(&config.BatchOptions{}, node)
	assert.NoError(t, err)
	bast.Print(bnode)
}

func TestMatch(t *testing.T) {
}
