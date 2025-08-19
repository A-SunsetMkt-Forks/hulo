// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package parser

import (
	"testing"

	"github.com/hulo-lang/hulo/syntax/batch/lexer"
	"github.com/stretchr/testify/assert"
)

func TestParseProgram(t *testing.T) {
	input := `@echo off
echo Hello World
if exist file.txt (
    echo File exists
) else (
    echo File not found
)`

	l := lexer.NewLexer(input)
	p := NewParser(l)
	program := p.ParseProgram()

	// 打印错误信息以便调试
	if len(p.Errors()) > 0 {
		for _, err := range p.Errors() {
			t.Logf("Parser error: %s", err)
		}
	}

	assert.NotNil(t, program)
	assert.Greater(t, len(program.Stmts), 0)
}

func TestParseIfStatement(t *testing.T) {
	input := `if exist file.txt (
    echo File exists
) else (
    echo File not found
)`

	l := lexer.NewLexer(input)
	p := NewParser(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		for _, err := range p.Errors() {
			t.Logf("Parser error: %s", err)
		}
	}

	assert.NotNil(t, program)
	assert.Greater(t, len(program.Stmts), 0)
}

func TestParseForStatement(t *testing.T) {
	input := `for %%i in (*.txt) do (
    echo %%i
)`

	l := lexer.NewLexer(input)
	p := NewParser(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		for _, err := range p.Errors() {
			t.Logf("Parser error: %s", err)
		}
	}

	assert.NotNil(t, program)
	assert.Greater(t, len(program.Stmts), 0)
}

func TestParseCallStatement(t *testing.T) {
	input := `call :function arg1 arg2`

	l := lexer.NewLexer(input)
	p := NewParser(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		for _, err := range p.Errors() {
			t.Logf("Parser error: %s", err)
		}
	}

	assert.NotNil(t, program)
	assert.Greater(t, len(program.Stmts), 0)
}

func TestParseGotoStatement(t *testing.T) {
	input := `goto :label`

	l := lexer.NewLexer(input)
	p := NewParser(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		for _, err := range p.Errors() {
			t.Logf("Parser error: %s", err)
		}
	}

	assert.NotNil(t, program)
	assert.Greater(t, len(program.Stmts), 0)
}

func TestParseLabelStatement(t *testing.T) {
	input := `:label`

	l := lexer.NewLexer(input)
	p := NewParser(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		for _, err := range p.Errors() {
			t.Logf("Parser error: %s", err)
		}
	}

	assert.NotNil(t, program)
	assert.Greater(t, len(program.Stmts), 0)
}

func TestParseComment(t *testing.T) {
	input := `rem This is a comment
:: This is also a comment`

	l := lexer.NewLexer(input)
	p := NewParser(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		for _, err := range p.Errors() {
			t.Logf("Parser error: %s", err)
		}
	}

	assert.NotNil(t, program)
	assert.Greater(t, len(program.Stmts), 0)
}

func TestParseAtCommand(t *testing.T) {
	input := `@echo off`

	l := lexer.NewLexer(input)
	p := NewParser(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		for _, err := range p.Errors() {
			t.Logf("Parser error: %s", err)
		}
	}

	assert.NotNil(t, program)
	assert.Greater(t, len(program.Stmts), 0)
}

func TestParseEmptyInput(t *testing.T) {
	input := ""

	l := lexer.NewLexer(input)
	p := NewParser(l)
	program := p.ParseProgram()

	assert.Empty(t, p.Errors())
	assert.NotNil(t, program)
	assert.Empty(t, program.Stmts)
}

func TestParseOnlyComments(t *testing.T) {
	input := `rem This is a comment
:: This is another comment`

	l := lexer.NewLexer(input)
	p := NewParser(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		for _, err := range p.Errors() {
			t.Logf("Parser error: %s", err)
		}
	}

	assert.NotNil(t, program)
	assert.Greater(t, len(program.Stmts), 0)
}
