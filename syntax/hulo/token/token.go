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

//go:generate stringer -type=Token -linecomment
type Token uint

const (
	ILLEGAL Token = iota
	EOF

	// Identifiers
	// add, foobar, x, y, ...
	IDENT // IDENT
	// 1343456
	NUM // NUM
	// "foobar"
	STR   // STR
	TRUE  // true
	FALSE // false

	// Operators
	ASSIGN // =
	PLUS   // +
	MINUS  // -

	COLON_ASSIGN // :=

	PLUS_ASSIGN // +=
	MINUS_ASSIGN // -=
	ASTERISK_ASSIGN // *=
	SLASH_ASSIGN // /=
	MOD_ASSIGN // %=
	AND_ASSIGN // &&=
	OR_ASSIGN // ||=
	POWER_ASSIGN // **=

	INC // ++
	DEC // --

	AT // @

	ASTERISK // *
	SLASH    // /
	MOD      // %
	POWER    // **

	HASH      // #
	REFERENCE // $
	OR        // ||
	AND       // &&

	CONCAT // &

	LT // <
	GT // >

	NOT   // !
	QUEST // ?

	EQ  // ==
	NEQ // !=

	LE // <=
	GE // >=

	DOCS // '

	// Delimiters
	COMMA     // ,
	SEMICOLON // ;
	COLON     // :
	DOT       // .
	NS        // ::

	LPAREN   // (
	RPAREN   // )
	LBRACE   // {
	RBRACE   // }
	LBRACKET // [
	RBRACKET // ]

	/// keywords

	// Package

	MODLIT // mod
	IMPORT // import
	FROM   // from
	AS     // as

	// Scopes

	LET   // let
	VAR   // var
	CONST // const

	// Function

	FN     // fn
	RETURN // return

	// Controll Flow

	IF    // if
	ELSE  // else
	MATCH // match
	LOOP  // loop
	DO    // do

	// Modifiers

	PUB      // pub
	FINAL    // final
	STATIC   // static
	COMPTIME // comptime
	UNSAFE   // unsafe

	// Class and Traits

	CLASS   // class
	EXTENDS // extends
	IMPL    // impl
	FOR     // for

	NEW    // new
	DELETE // delete

	EXTENSIONS // extensions

	// Error Handling
	TRY     // try
	CATCH   // catch
	FINALLY // finally

	THROW  // throw
	THROWS // throws
)
