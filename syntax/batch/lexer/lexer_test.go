// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package lexer

import (
	"testing"
	"github.com/hulo-lang/hulo/syntax/batch/token"
	"github.com/stretchr/testify/assert"
)

func TestNextToken(t *testing.T) {
	input := `@echo off
echo "Hello World"
if exist file.txt (
    echo File exists
) else (
    echo File not found
)
for %%i in (*.txt) do (
    echo %%i
)
set var=value
call :function
:function
echo In function
goto :eof
rem This is a comment
:: This is also a comment`

	tests := []struct {
		expectedType    token.Token
		expectedLiteral string
	}{
		{token.AT, "@"},
		{token.IDENT, "echo"},
		{token.IDENT, "off"},
		{token.IDENT, "echo"},
		{token.LITERAL, "\"Hello World\""},
		{token.IF, "if"},
		{token.IDENT, "exist"},
		{token.IDENT, "file.txt"},
		{token.LPAREN, "("},
		{token.IDENT, "echo"},
		{token.IDENT, "File"},
		{token.IDENT, "exists"},
		{token.RPAREN, ")"},
		{token.ELSE, "else"},
		{token.LPAREN, "("},
		{token.IDENT, "echo"},
		{token.IDENT, "File"},
		{token.IDENT, "not"},
		{token.IDENT, "found"},
		{token.RPAREN, ")"},
		{token.FOR, "for"},
		{token.IDENT, "%%i"},
		{token.IN, "in"},
		{token.LPAREN, "("},
		{token.IDENT, "*.txt"},
		{token.RPAREN, ")"},
		{token.DO, "do"},
		{token.LPAREN, "("},
		{token.IDENT, "echo"},
		{token.IDENT, "%%i"},
		{token.RPAREN, ")"},
		{token.IDENT, "set"},
		{token.IDENT, "var"},
		{token.ASSIGN, "="},
		{token.IDENT, "value"},
		{token.CALL, "call"},
		{token.COLON, ":"},
		{token.IDENT, "function"},
		{token.COLON, ":"},
		{token.IDENT, "function"},
		{token.IDENT, "echo"},
		{token.IDENT, "In"},
		{token.IDENT, "function"},
		{token.GOTO, "goto"},
		{token.COLON, ":"},
		{token.IDENT, "eof"},
		{token.REM, "rem"},
		{token.IDENT, "This"},
		{token.IDENT, "is"},
		{token.IDENT, "a"},
		{token.IDENT, "comment"},
		{token.DOUBLE_COLON, "::"},
		{token.IDENT, "This"},
		{token.IDENT, "is"},
		{token.IDENT, "also"},
		{token.IDENT, "a"},
		{token.IDENT, "comment"},
		{token.EOF, ""},
	}

	l := NewLexer(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestOperators(t *testing.T) {
	input := `+ - * / % = == != < <= > >= && || ! & | ^ ~ += -= *= /= %=`

	tests := []struct {
		expectedType    token.Token
		expectedLiteral string
	}{
		{token.ADD, "+"},
		{token.SUB, "-"},
		{token.MUL, "*"},
		{token.DIV, "/"},
		{token.MOD, "%"},
		{token.ASSIGN, "="},
		{token.DOUBLE_ASSIGN, "=="},
		{token.NEQ, "!="},
		{token.LSS, "<"},
		{token.LEQ, "<="},
		{token.GTR, ">"},
		{token.GEQ, ">="},
		{token.ANDAND, "&&"},
		{token.OROR, "||"},
		{token.NOT, "!"},
		{token.BITAND, "&"},
		{token.BITOR, "|"},
		{token.BITXOR, "^"},
		{token.TILDE, "~"},
		{token.ADD_ASSIGN, "+="},
		{token.SUB_ASSIGN, "-="},
		{token.MUL_ASSIGN, "*="},
		{token.DIV_ASSIGN, "/="},
		{token.MOD_ASSIGN, "%="},
		{token.EOF, ""},
	}

	l := NewLexer(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestKeywords(t *testing.T) {
	input := `if else for do in call goto rem`

	tests := []struct {
		expectedType    token.Token
		expectedLiteral string
	}{
		{token.IF, "if"},
		{token.ELSE, "else"},
		{token.FOR, "for"},
		{token.DO, "do"},
		{token.IN, "in"},
		{token.CALL, "call"},
		{token.GOTO, "goto"},
		{token.REM, "rem"},
		{token.EOF, ""},
	}

	l := NewLexer(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestIdentifiers(t *testing.T) {
	input := `foobar _bar foo_bar`

	tests := []struct {
		expectedType    token.Token
		expectedLiteral string
	}{
		{token.IDENT, "foobar"},
		{token.IDENT, "_bar"},
		{token.IDENT, "foo_bar"},
		{token.EOF, ""},
	}

	l := NewLexer(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestNumbers(t *testing.T) {
	input := `123 456.789`

	tests := []struct {
		expectedType    token.Token
		expectedLiteral string
	}{
		{token.LITERAL, "123"},
		{token.LITERAL, "456.789"},
		{token.EOF, ""},
	}

	l := NewLexer(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestWhitespace(t *testing.T) {
	input := `  \t\n\r  if`

	l := NewLexer(input)
	tok := l.NextToken()

	assert.Equal(t, token.IF, tok.Type)
	assert.Equal(t, "if", tok.Literal)
}

func TestIllegalToken(t *testing.T) {
	input := `$`

	l := NewLexer(input)
	tok := l.NextToken()

	assert.Equal(t, token.ILLEGAL, tok.Type)
	assert.Equal(t, "$", tok.Literal)
}
