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
type Token uint8

func (t Token) IsValid() bool {
	return t != ILLEGAL
}

const (
	ILLEGAL Token = iota // ILLEGAL

	/// Operators

	// Arithmetic

	// Addition Operator
	ADD // +
	// Subtraction or Negation Operator
	SUB // -
	// Multiplication Operator
	MUL // *
	// Division Operator
	DIV // /
	// Integer Division Operator
	IDIV // \
	// Mod Operator
	MOD // Mod
	// Exponentiation Operator
	EXP // ^
	// Concatenation Operator or Logical
	CONCAT // &

	// Comparison

	// Equality or Assignment Operator
	EQ // =
	// Inequality
	IEQ // <>
	// Less than
	LT // <
	// Greater than
	GT // >
	// Less than or equal to
	LTQ // <=
	// Greater than or equal to
	GTQ // >=
	// Is Operator
	IS // Is

	// Logical

	// Not Operator
	NOT // Not
	// And Operator
	AND // And
	// Or Operator
	OR // Or
	// Xor Operator
	XOR // Xor
	// Eqv Operator
	EQV // Eqv
	// Imp Operator
	IMP // Imp

	COLON // :
	COMMA // ,

	// Single Quote
	SGL_QUOTE // '
	// Double Quote
	DBL_QUOTE // "

	FALSE   // False
	TRUE    // True
	NOTHING // Nothing

	// By Value
	BYVAL // ByVal
	// By Reference
	BYREF // ByRef

	/// Property

	GET // Get
	LET // Let
	SET // Set

	CONST // Const

	/// Data Types

	EMPTY    // Empty
	NULL     // Null
	BOOLEAN  // Boolean
	BYTE     // Byte
	INTEGER  // Integer
	CURRENCY // Currency
	LONG     // Long
	SINGLE   // Single
	DOUBLE   // Double
	DATE     // Date
	STRING   // String
	OBJECT   // Object
	ERROR    // Error

	DIM       // Dim
	REDIM     // ReDim
	PRESERVE  // Preserve
	FOR       // For
	EACH      // Each
	IN        // In
	TO        // To
	STEP      // Step
	NEXT      // Next
	EXIT      // Exit
	SELECT    // Select
	CASE      // Case
	THEN      // Then
	IF        // If
	ELSEIF    // ElseIf
	ELSE      // Else
	WITH      // With
	WHILE     // While
	WEND      // Wend
	END       // End
	SUB_LIT   // Sub
	PROPERTY  // Property
	FUNCTION  // Function
	DEFAULT   // Default
	CLASS     // Class
	PUBLIC    // Public
	PRIVATE   // Private
	CALL      // Call
	ON        // On
	GOTO      // GoTo
	RESUME    // Resume
	STOP      // Stop
	RANDOMIZE // Randomize
	OPTION    // Option
	EXPLICIT  // Explicit
	ERASE     // Erase
	EXECUTE   // Excute
	NEW       // New
	REM       // Rem

	EOF // EOF
)
