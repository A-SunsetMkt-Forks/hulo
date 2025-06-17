// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package build_test

import (
	"testing"

	"github.com/caarlos0/log"
	build "github.com/hulo-lang/hulo/internal/build/vbs"
	"github.com/hulo-lang/hulo/internal/config"
	hast "github.com/hulo-lang/hulo/syntax/hulo/ast"
	"github.com/hulo-lang/hulo/syntax/hulo/parser"
	vast "github.com/hulo-lang/hulo/syntax/vbs/ast"
)

func TestCommandStmt(t *testing.T) {
	log.SetLevel(log.ErrorLevel)
	script := `echo "Hello, World!" 3.14 true`
	node, err := parser.ParseSourceScript(script, parser.ParseOptions{})
	if err != nil {
		t.Fatal(err)
	}
	// hast.Print(node)
	bnode, err := build.Translate(&config.VBScriptOptions{}, node)
	if err != nil {
		t.Fatal(err)
	}
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
	node, err := parser.ParseSourceScript(script, parser.ParseOptions{})
	if err != nil {
		t.Fatal(err)
	}
	// hast.Print(node)
	bnode, err := build.Translate(&config.VBScriptOptions{
		CommentStyle: "'",
	}, node)
	if err != nil {
		t.Fatal(err)
	}
	vast.Print(bnode)
}

func TestAssign(t *testing.T) {
	log.SetLevel(log.ErrorLevel)
	script := `let a = 10
	var b = 3.14
	const c = "Hello, World!"
	$d := true`
	node, err := parser.ParseSourceScript(script, parser.ParseOptions{})
	if err != nil {
		t.Fatal(err)
	}
	// hast.Print(node)
	bnode, err := build.Translate(&config.VBScriptOptions{}, node)
	if err != nil {
		t.Fatal(err)
	}
	vast.Print(bnode)
}

func TestIf(t *testing.T) {
	log.SetLevel(log.ErrorLevel)
	script := `if $a > 10 {
		echo "a is greater than 10"
	} else {
		echo "a is less than or equal to 10"
	}`
	node, err := parser.ParseSourceScript(script, parser.ParseOptions{})
	if err != nil {
		t.Fatal(err)
	}
	// hast.Print(node)
	bnode, err := build.Translate(&config.VBScriptOptions{}, node)
	if err != nil {
		t.Fatal(err)
	}
	vast.Print(bnode)
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
	node, err := parser.ParseSourceScript(script, parser.ParseOptions{})
	if err != nil {
		t.Fatal(err)
	}
	hast.Print(node)
	bnode, err := build.Translate(&config.VBScriptOptions{}, node)
	if err != nil {
		t.Fatal(err)
	}
	vast.Print(bnode)
}

func TestMatch(t *testing.T) {
}
