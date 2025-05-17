// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package token

type Pos int

// IsValid reports whether the position is valid.
func (p Pos) IsValid() bool {
	return p != NoPos
}

// The zero value for Pos is NoPos; there is no file and line information
// associated with it, and NoPos.IsValid() is false. NoPos is always
// smaller than any other Pos value. The corresponding Position value
// for NoPos is the zero value for Position.
const NoPos Pos = 0

const DynPos Pos = -1

type Token string

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// Identifiers
	IDENT = "IDENT" // add, foobar, x, y, ...
	NUM   = "num"   // 1343456
	STR   = "str"   // "foobar"
	TRUE  = "true"
	FALSE = "false"

	// Operators
	ASSIGN = "="
	PLUS   = "+"
	MINUS  = "-"

	INC = "++"
	DEC = "--"

	ASTERISK = "*"
	SLASH    = "/"
	MOD      = "%"
	POWER    = "**"

	HASH      = "#"
	REFERENCE = "$"
	OR        = "||"
	AND       = "&&"

	CONCAT = "&"

	LT = "<"
	GT = ">"

	EQ  = "=="
	NEQ = "!="

	LEQ = "<=" // <=
	GEQ = ">=" // >=

	DOCS = "'"

	// Delimiters
	COMMA     = ","
	SEMICOLON = ";"
	COLON     = ":"
	DOT       = "."
	NS        = "::"

	LPAREN   = "("
	RPAREN   = ")"
	LBRACE   = "{"
	RBRACE   = "}"
	LBRACKET = "["
	RBRACKET = "]"

	/// keywords

	// Package

	MODLIT = "mod"
	IMPORT = "import"
	FROM   = "from"
	AS     = "as"

	// Scopes

	LET   = "let"
	VAR   = "var"
	CONST = "const"

	// Function

	FN     = "fn"
	RETURN = "return"

	// Controll Flow

	IF    = "if"
	ELSE  = "else"
	MATCH = "match"
	LOOP  = "loop"
	DO    = "do"

	// Modifiers

	PUB      = "pub"
	FINAL    = "final"
	STATIC   = "static"
	COMPTIME = "comptime"
	UNSAFE   = "unsafe"

	// Class and Traits

	CLASS   = "class"
	EXTENDS = "extends"
	IMPL    = "impl"
	FOR     = "for"

	NEW    = "new"
	DELETE = "delete"

	EXTENSIONS = "extensions"
)
