// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package ast declares the types used to represent syntax trees for batch scripts.
// All node types implement the Node interface.
package ast

import (
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
}

// Stmt represents a statement node in the AST.
// All statement types implement the Stmt interface.
type Stmt interface {
	Node
	stmtNode() // dummy method to distinguish Stmt from other Node types
}

// ----------------------------------------------------------------------------
// Statements

type (
	// An IfStmt node represents an if statement.
	// if condition ( body ) else body
	IfStmt struct {
		Docs   *CommentGroup // associated documentation; or nil
		If     token.Pos     // position of "if" keyword
		Cond   Expr          // condition expression
		Lparen token.Pos     // position of "("
		Body   Stmt          // body of the if statement
		Rparen token.Pos     // position of ")"
		Else   Stmt          // else branch; or nil
	}

	// A ForStmt node represents a for statement.
	// for variable in list do ( body )
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

	// An ExprStmt node represents a (stand-alone) expression
	// in a statement list.
	ExprStmt struct {
		Docs *CommentGroup // associated documentation; or nil
		X    Expr          // expression
	}

	// An AssignStmt node represents an assignment statement.
	// set variable = value
	AssignStmt struct {
		Docs   *CommentGroup // associated documentation; or nil
		Set    token.Pos     // position of "set" keyword
		Opt    Expr          // optional expression
		Lhs    Expr          // left-hand side of assignment
		Assign token.Pos     // position of "="
		Rhs    Expr          // right-hand side of assignment
	}

	// A CallStmt node represents a call statement.
	// call :label arguments
	CallStmt struct {
		Docs   *CommentGroup // associated documentation; or nil
		Call   token.Pos     // position of "call" keyword
		Colon  token.Pos     // position of ":"
		Name   string        // name of the function to call
		Recv   []Expr        // arguments to the function
		IsFile bool          // whether the function is a file
	}

	// A FuncDecl node represents a function declaration.
	// :label body
	FuncDecl struct {
		Docs  *CommentGroup // associated documentation; or nil
		Colon token.Pos     // position of ":"
		Name  string        // function name
		Body  *BlockStmt    // function body
	}

	// A BlockStmt node represents a block of statements.
	BlockStmt struct {
		List []Stmt // list of statements
	}

	// A GotoStmt node represents a goto statement.
	// goto :label
	GotoStmt struct {
		Docs  *CommentGroup // associated documentation; or nil
		GoTo  token.Pos     // position of "goto" keyword
		Colon token.Pos     // position of ":"
		Label string        // label to jump to
	}

	// A LabelStmt node represents a label statement.
	// :label
	LabelStmt struct {
		Docs  *CommentGroup // associated documentation; or nil
		Colon token.Pos     // position of ":"
		Name  string        // label name
	}
)

// Position methods for statements
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
func (s *CmdExpr) Pos() token.Pos {
	return s.Name.Pos()
}

// End position methods for statements
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
func (s *CmdExpr) End() token.Pos {
	if len(s.Recv) == 0 {
		return s.Name.End()
	}
	return s.Recv[len(s.Recv)-1].End()
}
func (s *AssignStmt) End() token.Pos {
	return s.Rhs.End()
}

// Statement type assertions
func (*IfStmt) stmtNode()     {}
func (*ForStmt) stmtNode()    {}
func (*ExprStmt) stmtNode()   {}
func (*AssignStmt) stmtNode() {}
func (*FuncDecl) stmtNode()   {}
func (*BlockStmt) stmtNode()  {}
func (*CallStmt) stmtNode()   {}
func (*GotoStmt) stmtNode()   {}
func (*LabelStmt) stmtNode()  {}

// ----------------------------------------------------------------------------
// Expressions

type (
	// A Word node represents a word expression, which is a sequence of expressions.
	// Words are the basic building blocks of batch commands and arguments.
	Word struct {
		Parts []Expr // parts of the word
	}

	// A UnaryExpr node represents a unary expression.
	// Unary expressions include operators like ! (logical NOT)
	UnaryExpr struct {
		Op token.Token // operator
		X  Expr        // operand
	}

	// A BinaryExpr node represents a binary expression.
	// Binary expressions include operators like ==, !=, <, >, etc.
	BinaryExpr struct {
		X  Expr        // left operand
		Op token.Token // operator
		Y  Expr        // right operand
	}

	// A CallExpr node represents a function call expression.
	// Function calls in batch scripts
	CallExpr struct {
		Fun  Expr   // function expression
		Recv []Expr // arguments
	}

	// A Lit node represents a literal expression.
	// Literals include strings, numbers, and other constant values
	Lit struct {
		ValPos token.Pos // position of value
		Val    string    // literal value
	}

	// A SglQuote node represents a single quote expression.
	// Single quotes are used for literal strings in batch scripts
	SglQuote struct {
		ModPos token.Pos // position of "%"
		Val    Expr      // quoted expression
	}

	// A DblQuote node represents a double quote expression.
	// Double quotes are used for variable expansion in batch scripts
	DblQuote struct {
		DelayedExpansion bool      // whether to use delayed expansion (!var!)
		Left             token.Pos // position of "%"
		Val              Expr      // quoted expression
		Right            token.Pos // position of "%"
	}

	// A CmdExpr node represents a command statement.
	// Command expressions represent executable commands in batch scripts
	CmdExpr struct {
		Docs *CommentGroup // associated documentation; or nil
		Name Expr          // command name
		Recv []Expr        // command arguments
	}
)

// Position methods for expressions
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

// End position methods for expressions
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

// Expression type assertions
func (*UnaryExpr) exprNode()  {}
func (*BinaryExpr) exprNode() {}
func (*CallExpr) exprNode()   {}
func (*Word) exprNode()       {}
func (*Lit) exprNode()        {}
func (*SglQuote) exprNode()   {}
func (*DblQuote) exprNode()   {}
func (*CmdExpr) exprNode()    {}

// ----------------------------------------------------------------------------
// Comments

// A Comment node represents a single comment.
// Comments in batch scripts start with :: or REM
type Comment struct {
	TokPos token.Pos   // position of comment token
	Tok    token.Token // token.DOUBLE_COLON | token.REM
	Text   string      // comment text
}

func (c *Comment) Pos() token.Pos { return c.TokPos }
func (c *Comment) End() token.Pos { return c.TokPos + token.Pos(len(c.Text)) }
func (c *Comment) stmtNode()      {}

// A CommentGroup represents a sequence of comments.
// Comment groups are used to associate documentation with AST nodes
type CommentGroup struct {
	Comments []*Comment // list of comments
}

func (c *CommentGroup) Pos() token.Pos { return c.Comments[0].Pos() }
func (c *CommentGroup) End() token.Pos { return c.Comments[len(c.Comments)-1].End() }
func (c *CommentGroup) stmtNode()      {}

// ----------------------------------------------------------------------------
// File

// A File node represents a batch script file.
// Files contain the complete AST for a batch script
type File struct {
	Docs  []*CommentGroup // list of documentation comments
	Stmts []Stmt          // list of statements
}

func (f *File) Pos() token.Pos { return f.Stmts[0].Pos() }
func (f *File) End() token.Pos { return f.Stmts[len(f.Stmts)-1].End() }
