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

type Token string

const (
	NONE  = ""
	SPACE = " "

	ADD = "+"
	SUB = "-"
	MUL = "*"
	DIV = "/"
	MOD = "%"
	EXP = "**"

	HASH   = "#"
	QUEST  = "?"
	AT     = "@"
	DOLLAR = "$"

	INC = "++"
	DEC = "--"
	EQ  = "=="
	NEQ = "!="

	LT_AND = ">&"
	AND_LT = "&>"

	DOLLAR_MUL   = "$*"
	DOLLAR_AT    = "$@"
	DOLLAR_HASH  = "$#"
	DOLLAR_QUEST = "$?"
	DOLLAR_SUB   = "$-"
	DOLLAR_TWO   = "$$"
	DOLLAR_NOT   = "$!"
	DOLLAR_ZERO  = "$0"

	SINGLE_QUOTE = "'"
	DOUBLE_QUOTE = "\""

	BACK_QUOTE = "`"

	LT = "<"
	GT = ">"

	LT_ASSIGN        = "<="
	GT_ASSIGN        = ">="
	MUL_ASSIGN       = "*="
	DIV_ASSIGN       = "/="
	ADD_ASSIGN       = "+="
	SUB_ASSIGN       = "-="
	DOUBLE_LT_ASSIGN = "<<="
	DOUBLE_GT_ASSIGN = ">>="
	AND_ASSIGN       = "&="
	XOR_ASSIGN       = "^="
	OR_ASSIGN        = "|="

	DOUBLE_LT = "<<"
	TRIPLE_LT = "<<<"
	DOUBLE_GT = ">>"

	LPAREN   = "("
	RPAREN   = ")"
	LBRACE   = "{"
	RBRACE   = "}"
	LBRACKET = "["
	RBRACKET = "]"

	DOUBLE_LPAREN = "(("
	DOUBLE_RPAREN = "))"

	BITOR  = "|"
	BITAND = "&"
	BITNOT = "!"
	BITNEG = "~"

	ASSIGN = "="

	COMMA = ","
	COLON = ":"
	SEMI  = ";"

	DOUBLE_SEMI = ";;"

	AND = "&&"
	OR  = "||"
	XOR = "^"

	IF       = "if"
	ELSE     = "else"
	FOR      = "for"
	IN       = "in"
	UNTIL    = "until"
	WHILE    = "while"
	SWITCH   = "switch"
	CASE     = "case"
	SELECT   = "select"
	FUNCTION = "function"
	LOCAL    = "local"
	RETURN   = "return"
	BREAK    = "break"
	CONTINUE = "continue"
	END      = "end"

	STRING = "STRING"
	NUMBER = "NUMBER"
	WORD   = "WORD"

	EOF = "EOF"
)