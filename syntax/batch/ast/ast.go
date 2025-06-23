// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package ast

import (
	"fmt"
	"strings"

	"github.com/hulo-lang/hulo/syntax/batch/token"
)

// Node represents a node in the AST.
// All node types implement the Node interface.
type Node interface {
	Pos() token.Pos // position of first character belonging to the node
	End() token.Pos // position of first character immediately after the node
}

// Expr represents an expression node in the AST.
// All expression types implement the Expr interface.
type Expr interface {
	Node
	exprNode() // dummy method to distinguish Expr from other Node types
	String() string
}

// Stmt represents a statement node in the AST.
// All statement types implement the Stmt interface.
type Stmt interface {
	Node
	stmtNode() // dummy method to distinguish Stmt from other Node types
}

type (
	// IfStmt represents an if statement.
	IfStmt struct {
		Docs   *CommentGroup // associated documentation; or nil
		If     token.Pos     // position of "if" keyword
		Cond   Expr          // condition
		Lparen token.Pos     // position of "("
		Body   Stmt          // body of the if statement
		Rparen token.Pos     // position of ")"
		Else   Stmt          // else branch; or nil
	}

	// ForStmt represents a for statement.
	ForStmt struct {
		Docs   *CommentGroup // associated documentation; or nil
		For    token.Pos     // position of "for" keyword
		X      Expr          // variable to iterate over
		In     token.Pos     // position of "in" keyword
		List   Expr          // list to iterate over
		Do     token.Pos     // position of "do" keyword
		Lparen token.Pos     // position of "("
		Body   *BlockStmt    // body of the for statement
		Rparen token.Pos     // position of ")"
	}

	// ExprStmt represents an expression statement.
	ExprStmt struct {
		Docs *CommentGroup // associated documentation; or nil
		X    Expr          // expression
	}

	// AssignStmt represents an assignment statement.
	AssignStmt struct {
		Docs   *CommentGroup // associated documentation; or nil
		Set    token.Pos     // position of "set" keyword
		Lhs    Expr          // left-hand side of assignment
		Assign token.Pos     // position of "="
		Rhs    Expr          // right-hand side of assignment
	}

	// CallStmt represents a call statement.
	CallStmt struct {
		Docs  *CommentGroup // associated documentation; or nil
		Call  token.Pos     // position of "call" keyword
		Colon token.Pos     // position of ":"
		Name  string        // name of the function to call
		Recv  []Expr        // arguments to the function
	}

	// FuncDecl represents a function declaration.
	FuncDecl struct {
		Docs  *CommentGroup // associated documentation; or nil
		Colon token.Pos     // position of ":"
		Name  string        // function name
		Body  *BlockStmt    // function body
	}

	// BlockStmt represents a block of statements.
	BlockStmt struct {
		List []Stmt // list of statements
	}

	// GotoStmt represents a goto statement.
	GotoStmt struct {
		Docs  *CommentGroup // associated documentation; or nil
		GoTo  token.Pos     // position of "goto" keyword
		Colon token.Pos     // position of ":"
		Label string        // label to jump to
	}

	// LabelStmt represents a label statement.
	LabelStmt struct {
		Docs  *CommentGroup // associated documentation; or nil
		Colon token.Pos     // position of ":"
		Name  string        // label name
	}

	// Command represents a command statement.
	Command struct {
		Docs *CommentGroup // associated documentation; or nil
		Name Expr          // command name
		Recv []Expr        // command arguments
	}
)

func (s *FuncDecl) Pos() token.Pos   { return s.Colon }
func (s *IfStmt) Pos() token.Pos     { return s.If }
func (s *ForStmt) Pos() token.Pos    { return s.For }
func (s *ExprStmt) Pos() token.Pos   { return s.X.Pos() }
func (s *AssignStmt) Pos() token.Pos { return s.Set }
func (s *GotoStmt) Pos() token.Pos   { return s.GoTo }
func (s *LabelStmt) Pos() token.Pos  { return s.Colon }
func (s *BlockStmt) Pos() token.Pos {
	if len(s.List) == 0 {
		return token.NoPos
	}
	return s.List[0].Pos()
}
func (s *CallStmt) Pos() token.Pos {
	return s.Call
}
func (s *Command) Pos() token.Pos {
	return s.Name.Pos()
}

func (s *FuncDecl) End() token.Pos {
	if s.Body == nil {
		return token.NoPos
	}
	return s.Body.List[len(s.Body.List)-1].End()
}
func (s *IfStmt) End() token.Pos {
	return s.Else.End()
}
func (s *ForStmt) End() token.Pos {
	return s.Rparen
}
func (s *BlockStmt) End() token.Pos {
	if len(s.List) == 0 {
		return token.NoPos
	}
	return s.List[len(s.List)-1].End()
}
func (s *GotoStmt) End() token.Pos {
	return token.Pos(int(s.GoTo) + len(s.Label))
}
func (s *LabelStmt) End() token.Pos {
	return token.Pos(int(s.Colon) + len(s.Name))
}
func (s *ExprStmt) End() token.Pos {
	return s.X.End()
}
func (s *CallStmt) End() token.Pos {
	return token.Pos(int(s.Call) + len(s.Name))
}
func (s *Command) End() token.Pos {
	if len(s.Recv) == 0 {
		return s.Name.End()
	}
	return s.Recv[len(s.Recv)-1].End()
}
func (s *AssignStmt) End() token.Pos {
	return s.Rhs.End()
}

func (*IfStmt) stmtNode()      {}
func (*ForStmt) stmtNode()     {}
func (*ExprStmt) stmtNode()    {}
func (*AssignStmt) stmtNode()  {}
func (*FuncDecl) stmtNode()    {}
func (*BlockStmt) stmtNode()   {}
func (*CallStmt) stmtNode()    {}
func (*Command) stmtNode()     {}
func (*GotoStmt) stmtNode()  {}
func (*LabelStmt) stmtNode() {}

type (
	// Word represents a word expression, which is a sequence of expressions.
	Word struct {
		Parts []Expr // parts of the word
	}

	// UnaryExpr represents a unary expression.
	UnaryExpr struct {
		Op token.Token // operator
		X  Expr        // operand
	}

	// BinaryExpr represents a binary expression.
	BinaryExpr struct {
		X  Expr        // left operand
		Op token.Token // operator
		Y  Expr        // right operand
	}

	// CallExpr represents a function call expression.
	CallExpr struct {
		Fun  Expr   // function expression
		Recv []Expr // arguments
	}

	// Lit represents a literal expression.
	Lit struct {
		ValPos token.Pos // position of value
		Val    string    // literal value
	}

	// SglQuote represents a single quote expression.
	SglQuote struct {
		ModPos token.Pos // position of "%"
		Val    Expr      // quoted expression
	}

	// DblQuote represents a double quote expression.
	DblQuote struct {
		DelayedExpansion bool      // whether to use delayed expansion (!var!)
		Left             token.Pos // position of "%"
		Val              Expr      // quoted expression
		Right            token.Pos // position of "%"
	}
)

func (e *UnaryExpr) Pos() token.Pos {
	return e.X.Pos()
}
func (e *Word) Pos() token.Pos {
	if len(e.Parts) == 0 {
		return token.NoPos
	}
	return e.Parts[0].Pos()
}
func (e *Lit) Pos() token.Pos {
	return e.ValPos
}
func (e *CallExpr) Pos() token.Pos {
	return e.Fun.Pos()
}
func (e *SglQuote) Pos() token.Pos {
	return e.ModPos
}
func (e *DblQuote) Pos() token.Pos {
	return e.Left
}
func (b *BinaryExpr) Pos() token.Pos {
	return b.X.Pos()
}

func (e *UnaryExpr) End() token.Pos {
	return e.X.End()
}
func (e *Word) End() token.Pos {
	if len(e.Parts) == 0 {
		return token.NoPos
	}
	return e.Parts[len(e.Parts)-1].End()
}
func (e *Lit) End() token.Pos {
	return token.Pos(int(e.ValPos) + len(e.Val))
}
func (e *CallExpr) End() token.Pos {
	if len(e.Recv) == 0 {
		return e.Fun.End()
	}
	return e.Recv[len(e.Recv)-1].End()
}
func (e *SglQuote) End() token.Pos {
	return e.Val.End()
}
func (e *DblQuote) End() token.Pos {
	return e.Right
}
func (b *BinaryExpr) End() token.Pos {
	return b.Y.End()
}

func (*UnaryExpr) exprNode()  {}
func (*BinaryExpr) exprNode() {}
func (*CallExpr) exprNode()   {}
func (*Word) exprNode()       {}
func (*Lit) exprNode()        {}
func (*SglQuote) exprNode()   {}
func (*DblQuote) exprNode()   {}

func (w *Word) String() string {
	parts := make([]string, len(w.Parts))
	for i, part := range w.Parts {
		parts[i] = part.String()
	}
	return strings.Join(parts, " ")
}
func (u *UnaryExpr) String() string {
	return fmt.Sprintf("%s%s", u.Op, u.X)
}
func (b *BinaryExpr) String() string {
	switch b.Op {
	case token.EQU, token.NEQ, token.LSS,
		token.LEQ, token.GTR, token.GEQ,
		token.AND, token.OR, token.NOT:
		return fmt.Sprintf("%s %s %s", b.X, b.Op, b.Y)
	}
	return fmt.Sprintf("%s%s%s", b.X, b.Op, b.Y)
}
func (c *CallExpr) String() string {
	recv := make([]string, len(c.Recv))
	for i, r := range c.Recv {
		recv[i] = r.String()
	}
	return fmt.Sprintf("%s %s", c.Fun, strings.Join(recv, " "))
}
func (l *Lit) String() string {
	return l.Val
}
func (s *SglQuote) String() string {
	return fmt.Sprintf("%%%s", s.Val)
}
func (d *DblQuote) String() string {
	if d.DelayedExpansion {
		return fmt.Sprintf("!%s!", d.Val)
	}
	return fmt.Sprintf("%%%s%%", d.Val)
}

// Comment represents a single comment.
type Comment struct {
	TokPos token.Pos   // position of comment token
	Tok    token.Token // token.DOUBLE_COLON | token.REM
	Text   string      // comment text
}

func (c *Comment) Pos() token.Pos { return c.TokPos }
func (c *Comment) End() token.Pos { return c.TokPos + token.Pos(len(c.Text)) }

// CommentGroup represents a sequence of comments.
type CommentGroup struct {
	Comments []*Comment // list of comments
}

func (c *CommentGroup) Pos() token.Pos { return c.Comments[0].Pos() }
func (c *CommentGroup) End() token.Pos { return c.Comments[len(c.Comments)-1].End() }

// File represents a batch script file.
type File struct {
	Docs  []*CommentGroup // list of documentation comments
	Stmts []Stmt          // list of statements
}

func (f *File) Pos() token.Pos { return f.Stmts[0].Pos() }
func (f *File) End() token.Pos { return f.Stmts[len(f.Stmts)-1].End() }
