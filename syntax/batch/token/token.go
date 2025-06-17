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

//go:generate stringer -type=Token -linecomment
type Token uint32

func (t Token) IsValid() bool {
	return t != ILLEGAL
}

const (
	ILLEGAL Token = iota // ILLEGAL

	/// Operators

	// Arithmetic Operators

	ADD // +
	SUB // -
	MUL // *
	DIV // /
	MOD // %

	// Relational Operators

	EQU // EQU
	NEQ // NEQ
	LSS // LSS
	LEQ // LEQ
	GTR // GTR
	GEQ // GEQ

	// Logical Operators

	AND // AND
	OR  // OR
	NOT // NOT

	// Assignment Operators

	ADD_ASSIGN // +=
	SUB_ASSIGN // -=
	MUL_ASSIGN // *=
	DIV_ASSIGN // /=
	MOD_ASSIGN // %=

	// This is the bitwise and operator
	BITAND // &
	// This is the bitwise or operator
	BITOR // |
	// This is the bitwise xor or Exclusive or operator
	BITXOR // ^

	ASSIGN // =

	SEMI // ;

	EOF // EOF
)
