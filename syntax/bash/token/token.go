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

type Token uint

const (
	NONE = iota

	ADD // +
	SUB // -
	MUL // *
	DIV // /
	MOD // %
	EXP // **

	HASH   // #
	QUEST  // ?
	AT     // @
	DOLLAR // $

	INC // ++
	DEC // --
	EQ  // ==
	NEQ // !=

	LT_AND // >&
	AND_LT // &>

	DOLLAR_MUL   // $*
	DOLLAR_AT    // $@
	DOLLAR_HASH  // $#
	DOLLAR_QUEST // $?
	DOLLAR_SUB   // $-
	DOLLAR_TWO   // $$
	DOLLAR_NOT   // $!
	DOLLAR_ZERO  // $0

	SINGLE_QUOTE // '
	DOUBLE_QUOTE // "

	BACK_QUOTE // `

	LT // <
	GT // >

	LT_ASSIGN        // <=
	GT_ASSIGN        // >=
	MUL_ASSIGN       // *=
	DIV_ASSIGN       // /=
	ADD_ASSIGN       // +=
	SUB_ASSIGN       // -=
	DOUBLE_LT_ASSIGN // <<=
	DOUBLE_GT_ASSIGN // >>=
	AND_ASSIGN       // &=
	XOR_ASSIGN       // ^=
	OR_ASSIGN        // |=

	DOUBLE_LT // <<
	TRIPLE_LT // <<<
	DOUBLE_GT // >>

	LPAREN   // (
	RPAREN   // )
	LBRACE   // {
	RBRACE   // }
	LBRACKET // [
	RBRACKET // ]

	DOUBLE_LPAREN // ((
	DOUBLE_RPAREN // ))

	BITOR  // |
	BITAND // &
	BITNOT // !
	BITNEG // ~

	ASSIGN        // =
	ASSIGN_BITNEG // =~

	COMMA // ,
	COLON // :
	SEMI  // ;

	DOUBLE_SEMI // ;;

	AND // &&
	OR  // ||
	XOR // ^

	IF       // if
	THEN     // then
	ELIF     // elif
	ELSE     // else
	FI       // fi
	FOR      // for
	IN       // in
	UNTIL    // until
	WHILE    // while
	DO       // do
	DONE     // done
	CASE     // case
	ESAC     // esac
	SELECT   // select
	FUNCTION // function
	LOCAL    // local
	RETURN   // return
	BREAK    // break
	CONTINUE // continue

	STRING // string
	NUMBER // number
	WORD   // word

	EOF // eof
)

var ToString = map[Token]string{
	NONE: "",

	// Arithmetic operators
	ADD: "+",
	SUB: "-",
	MUL: "*",
	DIV: "/",
	MOD: "%",
	EXP: "**",

	// Special symbols
	HASH:   "#",
	QUEST:  "?",
	AT:     "@",
	DOLLAR: "$",

	// Increment/decrement
	INC: "++",
	DEC: "--",

	// Comparison operators
	EQ:  "==",
	NEQ: "!=",

	// Redirection operators
	LT_AND: ">&",
	AND_LT: "&>",

	// Special variables
	DOLLAR_MUL:   "$*",
	DOLLAR_AT:    "$@",
	DOLLAR_HASH:  "$#",
	DOLLAR_QUEST: "$?",
	DOLLAR_SUB:   "$-",
	DOLLAR_TWO:   "$$",
	DOLLAR_NOT:   "$!",
	DOLLAR_ZERO:  "$0",

	// Quotes
	SINGLE_QUOTE: "'",
	DOUBLE_QUOTE: "\"",
	BACK_QUOTE:   "`",

	// Comparison
	LT: "<",
	GT: ">",

	// Compound assignments
	LT_ASSIGN:        "<=",
	GT_ASSIGN:        ">=",
	MUL_ASSIGN:       "*=",
	DIV_ASSIGN:       "/=",
	ADD_ASSIGN:       "+=",
	SUB_ASSIGN:       "-=",
	DOUBLE_LT_ASSIGN: "<<=",
	DOUBLE_GT_ASSIGN: ">>=",
	AND_ASSIGN:       "&=",
	XOR_ASSIGN:       "^=",
	OR_ASSIGN:        "|=",

	// Bit shifting
	DOUBLE_LT: "<<",
	TRIPLE_LT: "<<<",
	DOUBLE_GT: ">>",

	// Brackets
	LPAREN:   "(",
	RPAREN:   ")",
	LBRACE:   "{",
	RBRACE:   "}",
	LBRACKET: "[",
	RBRACKET: "]",

	// Arithmetic expansion
	DOUBLE_LPAREN: "((",
	DOUBLE_RPAREN: "))",

	// Bitwise operators
	BITOR:  "|",
	BITAND: "&",
	BITNOT: "!",
	BITNEG: "~",

	// Assignment and pattern matching
	ASSIGN:        "=",
	ASSIGN_BITNEG: "=~",

	// Punctuation
	COMMA: ",",
	COLON: ":",
	SEMI:  ";",

	DOUBLE_SEMI: ";;",

	// Logical operators
	AND: "&&",
	OR:  "||",
	XOR: "^",

	// Keywords
	IF:       "if",
	THEN:     "then",
	ELIF:     "elif",
	ELSE:     "else",
	FI:       "fi",
	FOR:      "for",
	IN:       "in",
	UNTIL:    "until",
	WHILE:    "while",
	DO:       "do",
	DONE:     "done",
	CASE:     "case",
	ESAC:     "esac",
	SELECT:   "select",
	FUNCTION: "function",
	LOCAL:    "local",
	RETURN:   "return",
	BREAK:    "break",
	CONTINUE: "continue",

	// Literals
	STRING: "string",
	NUMBER: "number",
	WORD:   "word",

	// End of file
	EOF: "",
}
