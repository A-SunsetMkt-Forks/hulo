// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package ast

import "github.com/hulo-lang/hulo/grammar/powershell/token"

type Node interface {
	Pos() token.Pos
	End() token.Pos
}

type Decl interface {
	Node
	declNode()
}

type Stmt interface {
	Node
	stmtNode()
}

type Expr interface {
	Node
	exprNode()
}

type (
	UsingDecl struct {
		Using token.Pos
		*UsingModule
		*UsingNamespace
		*UsingAssembly
	}

	UsingNamespace struct {
		Namespace token.Pos
		Name      *Ident
	}

	UsingModule struct {
		Module token.Pos
		Name   *Ident
	}

	UsingAssembly struct {
		Assembly token.Pos
		Path     string
	}
)

type (
	// A BinaryExpr node represents a binary expression.
	BinaryExpr struct {
		X     Expr        // left operand
		OpPos token.Pos   // position of Op
		Op    token.Token // operator
		Y     Expr        // right operand
	}

	CallExpr struct {
		Func  Expr
		Args  []Expr
		Brace bool
	}

	// An Ident node represents an identifier.
	Ident struct {
		NamePos token.Pos // identifier position
		Name    string    // identifier name
	}

	List []Expr

	BasicLit struct {
		Kind  token.Token
		Value string
	}

	Fields struct {
		Sep string // " " "," ";"
	}

	VarExpr struct {
		Tok     token.Token // lbrace | lbrack
		Opening token.Pos
		X       Expr
		Closing token.Pos
	}

	BlockExpr struct {
		List []Expr
	}

	SubExpr struct {
		X Expr
	}

	Attribute struct {
		Key   Expr
		Value Expr
	}

	HashTable map[any]any

	IndexExpr struct {
		X      Expr
		Lbrack token.Token
		Index  Expr
		Rbrack token.Token
	}

	ExplictExpr struct {
		X Expr
		Y Expr
	}

	ModExpr struct {
		X Expr
		Y Expr
	}

	// ++X --X
	IncDecExpr struct {
		X      Expr
		Tok    token.Token // INC or DEC
		TokPos token.Pos   // position of "++" or "--"
	}

	MemberAccess struct {
		Type    Expr
		Colon   token.Pos
		Method  Expr
		Generic Expr
		Recv    []Expr
	}
)

func (x *BinaryExpr) Pos() token.Pos { return x.X.Pos() }
func (x *Ident) Pos() token.Pos      { return x.NamePos }
func (x *IncDecExpr) Pos() token.Pos { return x.X.Pos() }

func (x *BinaryExpr) End() token.Pos { return x.Y.End() }
func (x *Ident) End() token.Pos      { return token.Pos(int(x.NamePos) + len(x.Name)) }
func (x *IncDecExpr) End() token.Pos {
	return x.TokPos + 2
}

func (*BinaryExpr) exprNode()  {}
func (*CallExpr) exprNode()    {}
func (*Ident) exprNode()       {}
func (*List) exprNode()        {}
func (*BasicLit) exprNode()    {}
func (*VarExpr) exprNode()     {}
func (*BlockExpr) exprNode()   {}
func (*SubExpr) exprNode()     {}
func (*Attribute) exprNode()   {}
func (*HashTable) exprNode()   {}
func (*IndexExpr) exprNode()   {}
func (*ExplictExpr) exprNode() {}
func (*ModExpr) exprNode()     {}
func (*IncDecExpr) exprNode()  {}

type (
	File struct {
		Name string

		Stmts []Stmt
	}

	BlcokStmt struct {
		List []Stmt
	}

	ExprStmt struct {
		Expr
	}

	AssignStmt struct {
		Lhs Expr
		Rhs Expr
	}

	TryStmt struct {
		Body    *BlcokStmt
		Catches []*CatchStmt
		Finally *BlcokStmt
	}

	CatchStmt struct {
		Cond Expr
		Body *BlcokStmt
	}

	FuncStmt struct {
		Func Expr
		Recv []Expr
		Body *BlcokStmt
	}

	ForeachStmt struct {
		Lhs  Expr
		Rhs  Expr
		Body *BlcokStmt
	}

	ForStmt struct {
		Lhs  Expr
		Cond Expr
		Rhs  Expr
		Body *BlcokStmt
	}

	IfStmt struct {
		Cond Expr
		Body *BlcokStmt
		Else Stmt
	}

	ReturnStmt struct {
		X Expr
	}

	ThrowStmt struct {
		X Expr
	}

	LabelStmt struct {
		X Expr
	}

	WhileStmt struct {
		Do   bool
		Cond Expr
		Body *BlcokStmt
	}

	ContinueStmt struct{}

	SwitchStmt struct {
		Pattern       SwitchPattern
		Casesensitive bool
		Cond          Expr
		Cases         []*CaseStmt
		Default       *CaseStmt
	}

	CaseStmt struct {
		Cond Expr
		Body *BlcokStmt
	}

	DataStmt struct {
		Param []Expr
		Body  *BlcokStmt
	}

	BreakStmt struct {
		Break token.Pos
		Label *Ident
	}
)

type SwitchPattern int

func (p SwitchPattern) String() string {
	switch p {
	case SwitchPatternRegex:
		return "-regex"
	case SwitchPatternWildcard:
		return "-wildcard"
	case SwitchPatternExact:
		return "-exact"
	}
	return ""
}

const (
	SwitchPatternNone SwitchPattern = iota
	SwitchPatternRegex
	SwitchPatternWildcard
	SwitchPatternExact
)

func (*File) stmtNode()         {}
func (*BlcokStmt) stmtNode()    {}
func (*ExprStmt) stmtNode()     {}
func (*AssignStmt) stmtNode()   {}
func (*FuncStmt) stmtNode()     {}
func (*IfStmt) stmtNode()       {}
func (*ForStmt) stmtNode()      {}
func (*ForeachStmt) stmtNode()  {}
func (*WhileStmt) stmtNode()    {}
func (*ReturnStmt) stmtNode()   {}
func (*ThrowStmt) stmtNode()    {}
func (*LabelStmt) stmtNode()    {}
func (*ContinueStmt) stmtNode() {}
func (*SwitchStmt) stmtNode()   {}
func (*DataStmt) stmtNode()     {}
