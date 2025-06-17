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
	Illegal Token = iota

	SglQuote // '
	DblQuote // "
	BckQuote // `

	And    // &
	AndAnd // &&
	Or     // |
	OrOr   // ||
	OrAnd  // |&

	Dollar       // $
	DollSglQuote // $'
	DollDblQuote // $"
	DollBrace    // ${
	DollBrack    // $[
	DollParen    // $(
	DollDblParen // $((

	LeftBrack    // [
	DblLeftBrack // [[
	LeftParen    // (
	DblLeftParen // ((

	RightBrace    // }
	RightBrack    // ]
	RightParen    // )
	DblRightParen // ))

	Semicolon    // ;
	DblSemicolon // ;;
	SemiAnd      // ;&
	DblSemiAnd   // ;;&
	SemiOr       // ;|

	ExclMark   // !
	Tilde      // ~
	AddAdd     // ++
	SubSub     // --
	Star       // *
	Power      // **
	Equal      // ==
	NotEqual   // !=
	LessEqual  // <=
	GreatEqual // >=

	AddAssgn // +=
	SubAssgn // -=
	MulAssgn // *=
	QuoAssgn // /=
	RemAssgn // %=
	AndAssgn // &=
	OrAssgn  // |=
	XorAssgn // ^=
	ShlAssgn // <<=
	ShrAssgn // >>=

	RdrOut   // >
	AppOut   // >>
	RdrIn    // <
	RdrInOut // <>
	DplIn    // <&
	DplOut   // >&
	ClbOut   // >|
	Hdoc     // <<
	DashHdoc // <<-
	WordHdoc // <<<
	RdrAll   // &>
	AppAll   // &>>

	CmdIn  // <(
	CmdOut // >(

	Plus     // +
	ColPlus  // :+
	Minus    // -
	ColMinus // :-
	Quest    // ?
	ColQuest // :?
	Assgn    // =
	ColAssgn // :=
	Perc     // %
	DblPerc  // %%
	Hash     // #
	DblHash  // ##
	Caret    // ^
	DblCaret // ^^
	Comma    // ,
	DblComma // ,,
	At       // @
	Slash    // /
	DblSlash // //
	Colon    // :

	TsExists  // -e
	TsRegFile // -f
	TsDirect  // -d
	TsCharSp  // -c
	TsBlckSp  // -b
	TsNmPipe  // -p
	TsSocket  // -S
	TsSmbLink // -L
	TsSticky  // -k
	TsGIDSet  // -g
	TsUIDSet  // -u
	TsGrpOwn  // -G
	TsUsrOwn  // -O
	TsModif   // -N
	TsRead    // -r
	TsWrite   // -w
	TsExec    // -x
	TsNoEmpty // -s
	TsFdTerm  // -t
	TsEmpStr  // -z
	TsNempStr // -n
	TsOptSet  // -o
	TsVarSet  // -v
	TsRefVar  // -R

	TsReMatch // =~
	TsNewer   // -nt
	TsOlder   // -ot
	TsDevIno  // -ef
	TsEql     // -eq
	TsNeq     // -ne
	TsLeq     // -le
	TsGeq     // -ge
	TsLss     // -lt
	TsGtr     // -gt

	GlobQuest // ?(
	GlobStar  // *(
	GlobPlus  // +(
	GlobAt    // @(
	GlobExcl  // !(

	EOF // EOF

	// compitable with old version
	STRING
	NUMBER
)
