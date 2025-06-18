package ast

import "github.com/hulo-lang/hulo/syntax/batch/token"

type Node interface {
	Pos() token.Pos
	End() token.Pos
}

type Expr interface {
	exprNode()
}

type Stmt interface {
	stmtNode()
}

type (
	IfStmt struct {
		Cond Expr
		Body Stmt
		Else Stmt
	}

	ForStmt struct {
		Init Stmt
		Cond Expr
		Post Stmt
		Body Stmt
	}

	ExprStmt struct {
		X Expr
	}

	AssignStmt struct {
		Lhs Expr
	}
)

func (*IfStmt) exprNode()     {}
func (*ForStmt) exprNode()    {}
func (*ExprStmt) exprNode()   {}
func (*AssignStmt) exprNode() {}

type (
	Word struct {
		Parts []Expr
	}

	UnaryExpr struct {
		Op token.Token
		X  Expr
	}

	BinaryExpr struct {
		X  Expr
		Op token.Token
		Y  Expr
	}

	CallExpr struct {
		Fun  Expr
		Recv []Expr
	}

	Lit struct {
		Val string
	}

	SglQuote struct {
		ModPos token.Pos // position of "%"
		Val    string
	}

	DblQuote struct {
		Left  token.Pos // position of "%"
		Val   string
		Right token.Pos // position of "%"
	}
)

func (*UnaryExpr) exprNode()  {}
func (*BinaryExpr) exprNode() {}
func (*CallExpr) exprNode()   {}
func (*Word) exprNode()       {}
func (*Lit) exprNode()        {}
func (*SglQuote) exprNode()   {}
func (*DblQuote) exprNode()   {}
