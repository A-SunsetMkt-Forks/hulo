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

const (
	Illegal Token = iota // ILLEGAL

	// Literals

	IDENT  // identifier
	INT    // integer literal
	FLOAT  // float literal
	STRING // string literal

	// Operators

	ADD // +
	SUB // -
	MUL // *
	DIV // /
	MOD // %

	// Punctuation

	LPAREN       // (
	RPAREN       // )
	LBRACK       // [
	RBRACK       // ]
	LBRACE       // {
	RBRACE       // }
	SEMICOLON    // ;
	COMMA        // ,
	DOT          // .
	COLON        // :
	DOUBLE_COLON // ::

	// Quotes and Comments

	SGLQUOTE            // '
	DBLQUOTE            // "
	HASH                // #
	DELIMITED_CMT_BEGIN // <#
	DELIMITED_CMT_END   // #>

	// PowerShell Keywords

	BEGIN            // begin
	BREAK            // break
	CATCH            // catch
	CLASS            // class
	CONTINUE         // continue
	DATA             // data
	DEFINE           // define
	DO               // do
	DYNAMICPARAM     // dynamicparam
	ELSE             // else
	ELSEIF           // elseif
	END              // end
	EXIT             // exit
	FILTER           // filter
	FINALLY          // finally
	FOR              // for
	FOREACH          // foreach
	FROM             // from
	FUNCTION         // function
	IF               // if
	IN_KEYWORD       // in (keyword)
	INLINESCRIPT     // inlinescript
	PARALLEL_KEYWORD // parallel (keyword)
	PARAM            // param
	PROCESS          // process
	RETURN           // return
	SWITCH           // switch
	THROW            // throw
	TRAP             // trap
	TRY              // try
	UNTIL            // until
	USING            // using
	VAR              // var
	WHILE            // while
	WORKFLOW         // workflow

	// PowerShell Special Keywords

	CONSTANT         // constant
	ENUM             // enum
	HIDDEN           // hidden
	INLINE           // inline
	PARALLEL_SPECIAL // parallel (special)
	PRIVATE          // private
	REQUIRES         // requires
	STATIC           // static
	THROWS           // throws
	VALIDATE         // validate
	VALIDATESET      // validateset

	// Variable and Scope

	DOLLAR          // $
	AT              // @
	QUESTION        // ?
	DOUBLE_QUESTION // ??

	// Assignment Operators

	ASSIGN      // =
	ADD_ASSIGN  // +=
	SUB_ASSIGN  // -=
	MUL_ASSIGN  // *=
	DIV_ASSIGN  // /=
	MOD_ASSIGN  // %=
	BAND_ASSIGN // -band=
	BOR_ASSIGN  // -bor=
	BXOR_ASSIGN // -bxor=
	SHL_ASSIGN  // -shl=
	SHR_ASSIGN  // -shr=

	// Increment/Decrement

	INC // ++
	DEC // --

	// Arithmetic Operators

	BAND // -band
	BNOT // -bnot
	BOR  // -bor
	BXOR // -bxor
	SHL  // -shl
	SHR  // -shr

	// Comparison Operators

	EQ  // -eq
	IEQ // -ieq
	CEQ // -ceq
	NE  // -ne
	INE // -ine
	CNE // -cne
	GT  // -gt
	IGT // -igt
	CGT // -cgt
	GE  // -ge
	IGE // -ige
	CGE // -cge
	LT  // -lt
	ILT // -ilt
	CLT // -clt
	LE  // -le
	ILE // -ile
	CLE // -cle

	// String Comparison Operators

	// Matches a string against a wildcard pattern
	LIKE // -like
	// Matches a string against a wildcard pattern (case-insensitive)
	ILIKE // -ilike
	// Matches a string against a wildcard pattern (case-sensitive)
	CLIKE // -clike
	// Does not match a string against a wildcard pattern
	NOTLIKE // -notlike
	// Does not match a string against a wildcard pattern (case-insensitive)
	INOTLIKE // -inotlike
	// Does not match a string against a wildcard pattern (case-sensitive)
	CNOTLIKE // -cnotlike

	// Regular Expression Operators

	// Matches a string against a regular expression pattern
	MATCH // -match
	// Matches a string against a regular expression pattern (case-insensitive)
	IMATCH // -imatch
	// Matches a string against a regular expression pattern (case-sensitive)
	CMATCH // -cmatch
	// Does not match a string against a regular expression pattern
	NOTMATCH // -notmatch
	// Does not match a string against a regular expression pattern (case-insensitive)
	INOTMATCH // -inotmatch
	// Does not match a string against a regular expression pattern (case-sensitive)
	CNOTMATCH // -cnotmatch
	// Replaces parts of a string that match a regular expression pattern
	REPLACE // -replace
	// Replaces parts of a string that match a regular expression pattern (case-insensitive)
	IREPLACE // -ireplace
	// Replaces parts of a string that match a regular expression pattern (case-sensitive)
	CREPLACE // -creplace

	// Collection Operators

	// Checks if a collection contains a specific value
	CONTAINS // -contains
	// Checks if a collection contains a specific value (case-insensitive)
	ICONTAINS // -icontains
	// Checks if a collection contains a specific value (case-sensitive)
	CCONTAINS // -ccontains
	// Checks if a collection does not contain a specific value
	NOTCONTAINS // -notcontains
	// Checks if a collection does not contain a specific value (case-insensitive)
	INOTCONTAINS // -inotcontains
	// Checks if a collection does not contain a specific value (case-sensitive)
	CNOTCONTAINS // -cnotcontains
	// Checks if a value is within a collection
	IN_OPERATOR // -in
	// Checks if a value is not within a collection
	NOTIN // -notin

	// Logical Operators

	AND // -and
	OR  // -or
	NOT // -not
	XOR // -xor

	// Bitwise Operators (symbolic)

	BITAND // &
	BITOR  // |
	BITNOT // !

	// Pipeline Operators

	PIPE    // |
	PIPEAND // &&
	PIPEOR  // ||

	// Redirection Operators

	// Sends the specified stream to a file
	REWRITE // >
	// Appends the specified stream to the file
	APPEND // >>
	// Redirects the specified stream to the Success stream
	REDIRECT // >&
	// Redirects input from a file
	REDIRECT_IN // <

	// String Operators

	SPLIT  // -split
	ISPLIT // -isplit
	JOIN   // -join

	// Type Operators

	// Checks if an object is an instance of a specified .NET type.
	// Returns TRUE if the type matches; otherwise, FALSE.
	IS // -is
	// Checks if an object is NOT an instance of a specified .NET type.
	// Returns FALSE if the type matches; otherwise, TRUE.
	ISNOT // -isnot
	// Attempts to convert the input object to a specified .NET type.
	// Returns the converted object if successful; otherwise, $null.
	// Does not throw an error if conversion fails.
	AS // -as

	// Format Operator
	FORMAT // -f

	// Range Operator
	RANGE // ..

	// Special Tokens

	BACKTICK   // `
	NEWLINE    // \n
	WHITESPACE // whitespace

	EOF // EOF
)
