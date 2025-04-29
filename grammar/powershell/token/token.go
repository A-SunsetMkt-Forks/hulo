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
	ADD = "+"
	SUB = "-"
	MUL = "*"
	DIV = "/"
	MOD = "%"

	// Arithmetic Operators

	BAND = "-band"
	BNOT = "-bnot"
	BOR  = "-bor"
	BXOR = "-bxor"

	SHL = "-shl"
	SHR = "-shr"

	// Assignment Operators

	ADD_ASSIGN          = "+="
	SUB_ASSIGN          = "-="
	MUL_ASSIGN          = "*="
	MOD_ASSIGN          = "%="
	DOUBLE_QUEST_ASSIGN = "??="

	INC = "++"
	DEC = "--"

	// Comparison Operators

	EQ  = "-eq"
	IEQ = "-ieq"
	CEQ = "-ceq"
	NE  = "-ne"
	INE = "-ine"
	CNE = "-cne"
	GT  = "-gt"
	IGT = "-igt"
	CGT = "-cgt"
	GE  = "-ge"
	IGE = "-ige"
	CGE = "-cge"
	LT  = "-lt"
	ILT = "-ilt"
	CLT = "-clt"
	LE  = "-le"
	ILE = "-ile"
	CLE = "-cle"

	LIKE  = "-like"  // Matches a string against a wildcard pattern
	ILIKE = "-ilike" // Matches a string against a wildcard pattern (case-insensitive)
	CLIKE = "-clike" // Matches a string against a wildcard pattern (case-sensitive)

	NOTLIKE  = "-notlike"  // Does not match a string against a wildcard pattern
	INOTLIKE = "-inotlike" // Does not match a string against a wildcard pattern (case-insensitive)
	CNOTLIKE = "-cnotlike" // Does not match a string against a wildcard pattern (case-sensitive)

	MATCH  = "-match"  // Matches a string against a regular expression pattern
	IMATCH = "-imatch" // Matches a string against a regular expression pattern (case-insensitive)
	CMATCH = "-cmatch" // Matches a string against a regular expression pattern (case-sensitive)

	NOTMATCH  = "-notmatch"  // Does not match a string against a regular expression pattern
	INOTMATCH = "-inotmatch" // Does not match a string against a regular expression pattern (case-insensitive)
	CNOTMATCH = "-cnotmatch" // Does not match a string against a regular expression pattern (case-sensitive)

	REPLACE  = "-replace"  // Replaces parts of a string that match a regular expression pattern
	IREPLACE = "-ireplace" // Replaces parts of a string that match a regular expression pattern (case-insensitive)
	CREPLACE = "-creplace" // Replaces parts of a string that match a regular expression pattern (case-sensitive)

	CONTAINS  = "-contains"  // Checks if a collection contains a specific value
	ICONTAINS = "-icontains" // Checks if a collection contains a specific value (case-insensitive)
	CCONTAINS = "-ccontains" // Checks if a collection contains a specific value (case-sensitive)

	NOTCONTAINS  = "-notcontains"  // Checks if a collection does not contain a specific value
	INOTCONTAINS = "-inotcontains" // Checks if a collection does not contain a specific value (case-insensitive)
	CNOTCONTAINS = "-cnotcontains" // Checks if a collection does not contain a specific value (case-sensitive)

	IN    = "-in"    // Checks if a value is within a collection
	NOTIN = "-notin" // Checks if a value is not within a collection

	// Logical Operators

	AND = "-and"
	OR  = "-or"
	NOT = "-not" // -not or !
	XOR = "-xor"

	// Redirection

	REWRITE  = ">"   // Sends the specified stream to a file
	APPEND   = ">>"  // Appends the specified stream to the file
	REDIRECT = ">&1" // Redirects the specified stream to the Success stream

	SPLIT  = "-split"
	ISPLIT = "-iSplit"

	JOIN = "-Join"

	// Type Operators

	// Checks if an object is an instance of a specified .NET type.
	// Returns TRUE if the type matches; otherwise, FALSE.
	IS = "-is"

	// Checks if an object is NOT an instance of a specified .NET type.
	// Returns FALSE if the type matches; otherwise, TRUE.
	ISNOT = "-isnot"

	// Attempts to convert the input object to a specified .NET type.
	// Returns the converted object if successful; otherwise, $null.
	// Does not throw an error if conversion fails.
	AS = "-as"

	FORMAT = "-f"

	LPAREN = "("
	RPAREN = ")"
	LBRACK = "["
	RBRACK = "]"
	LBRACE = "{"
	RBRACE = "}"

	AT     = "@"
	DOLLAR = "$"
	RANGE  = ".."

	STR = "STR"
	NUM = "NUM"

	DOT   = "."
	COMMA = ","

	DOUBLE_COLON = "::"

	QUEST        = "?"
	DOUBLE_QUEST = "??"

	BITOR  = "|"
	BITAND = "&"
	BITNOT = "!"

	PIPE    = BITOR
	PIPEAND = "&&"
	PIPEOR  = "||"
)
