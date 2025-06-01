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
	ADD  = "+"
	SUB  = "-"
	MUL  = "*"
	DIV  = "/"
	IDIV = "\\"
	MOD  = "Mod"
	EXP  = "^"

	IS = "Is"

	BITAND = "&"

	NOT = "Not"
	AND = "And"
	OR  = "Or"
	XOR = "Xor"
	EQV = "Eqv"
	IMP = "Imp"

	EQ  = "="
	NEQ = "<>"

	LT = "<"
	GT = ">"

	COLON = ":"
	COMMA = ","

	LT_ASSIGN = "<="
	GT_ASSIGN = ">="

	APOSTROPHE = "'"

	FALSE   = "False"
	TRUE    = "True"
	NOTHING = "Nothing"

	BYVAL = "ByVal"
	BYREF = "ByRef"

	GET   = "Get"
	LET   = "Let"
	SET   = "Set"
	CONST = "Const"

	// Data Types

	EMPTY    = "Empty"
	NULL     = "Null"
	BOOLEAN  = "Boolean"
	BYTE     = "Byte"
	INTEGER  = "Integer"
	CURRENCY = "Currency"
	LONG     = "Long"
	SINGLE   = "Single"
	DOUBLE   = "Double"
	DATE     = "Date"
	STRING   = "String"
	OBJECT   = "Object"
	ERROR    = "Error"

	DIM       = "Dim"
	REDIM     = "ReDim"
	PRESERVE  = "Preserve"
	FOR       = "For"
	EACH      = "Each"
	IN        = "In"
	TO        = "To"
	STEP      = "Step"
	NEXT      = "Next"
	EXIT      = "Exit"
	SELECT    = "Select"
	CASE      = "Case"
	THEN      = "Then"
	IF        = "If"
	ELSEIF    = "ElseIf"
	ELSE      = "Else"
	WITH      = "With"
	WHILE     = "While"
	WEND      = "Wend"
	END       = "End"
	SUB_LIT   = "Sub"
	PROPERTY  = "Property"
	FUNCTION  = "Function"
	CLASS     = "Class"
	PUBLIC    = "Public"
	PRIVATE   = "Private"
	CALL      = "Call"
	ON        = "On"
	GOTO      = "GoTo"
	RESUME    = "Resume"
	STOP      = "Stop"
	RANDOMIZE = "Randomize"
	OPTION    = "Option"
	EXPLICIT  = "Explicit"
)
