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
	NULL  // null

	// Operators
	ASSIGN // =
	PLUS   // +
	MINUS  // -

	COLON_ASSIGN // :=

	PLUS_ASSIGN     // +=
	MINUS_ASSIGN    // -=
	ASTERISK_ASSIGN // *=
	SLASH_ASSIGN    // /=
	MOD_ASSIGN      // %=
	AND_ASSIGN      // &&=
	OR_ASSIGN       // ||=
	POWER_ASSIGN    // **=

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

	PIPE // |
	CONCAT // &

	LT // <
	GT // >

	SHL // <<
	SHR // >>

	NOT   // !
	QUEST // ?

	EQ  // ==
	NEQ // !=

	LE // <=
	GE // >=

	SGL_QUOTE // '
	DBL_QUOTE // "

	DEFER   // defer
	DECLARE // declare
	WHEN    // when

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

	// Module

	MODLIT // mod
	IMPORT // import
	FROM   // from
	USE    // use

	// type system

	TYPE   // type
	TYPEOF // typeof
	AS     // as

	// Scopes

	LET   // let
	VAR   // var
	CONST // const

	// Function

	FN       // fn
	OPERATOR // operator
	RETURN   // return

	// Controll Flow

	IF       // if
	ELSE     // else
	MATCH    // match
	LOOP     // loop
	DO       // do
	IN       // in
	OF       // of
	CONTINUE // continue
	BREAK    // break

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
	THROW   // throw
	THROWS  // throws

	// a-z A-Z Operators

	Opa // a
	Opb // b
	Opc // c
	Opd // d
	Ope // e
	Opf // f
	Opg // g
	Oph // h
	Opi // i
	Opj // j
	Opk // k
	Opl // l
	Opm // m
	Opn // n
	Opo // o
	Opp // p
	Opq // q
	Opr // r
	Ops // s
	Opt // t
	Opu // u
	Opv // v
	Opw // w
	Opz // z

	OpA // A
	OpB // B
	OpC // C
	OpD // D
	OpE // E
	OpF // F
	OpG // G
	OpH // H
	OpI // I
	OpJ // J
	OpK // K
	OpL // L
	OpM // M
	OpN // N
	OpO // O
	OpP // P
	OpQ // Q
	OpR // R
	OpS // S
	OpT // T
	OpU // U
	OpV // V
	OpW // W
	OpX // X
	OpY // Y
	OpZ // Z

	ASYNC // async
	AWAIT // await

	EXTERN // extern
)
