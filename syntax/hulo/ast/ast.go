// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package ast

import (
	"fmt"
	"strings"

	"github.com/hulo-lang/hulo/syntax/hulo/token"
)

type Node interface {
	Pos() token.Pos
	End() token.Pos
}

type Stmt interface {
	Node
	stmtNode()
}

type Decl interface {
	Node
	declNode()
}

type Expr interface {
	Node
	exprNode()
	String() string
}

// ----------------------------------------------------------------------------
// Declarations

type (
	// An TraitDecl node represents a trait declaration.
	TraitDecl struct {
		Doc    *CommentGroup
		Pub    token.Pos // position of "pub"
		Trait  token.Pos // position of "trait"
		Name   *Ident
		Lbrace token.Pos // position of "{"
		Fields *FieldList
		Rbrace token.Pos // position of "}"
	}

	// An ImplDecl node represents an impl declaration.
	ImplDecl struct {
		Doc   *CommentGroup
		Impl  token.Pos // position of "impl"
		Trait *Ident
		For   token.Pos // position of "for"
		*ImplDeclBody
		*ImplDeclBinding
	}

	ImplDeclBody struct {
		Class *Ident
		Body  *BlockStmt
	}

	ImplDeclBinding struct {
		Classes []*Ident
	}

	// An EnumDecl node represents an enum declaration.
	EnumDecl struct {
		Doc     *CommentGroup
		Macros  []*Annotation
		Pub     token.Pos // position of "pub"
		Enum    token.Pos // position of "enum"
		Name    *Ident
		Generic *GenericExpr
		*EnumBody
		*EnumBodySimple
	}

	EnumMember struct {
	}

	EnumMethod struct {
		Mod Modifier
		Fn  token.Pos // position of "fn"
	}

	EnumBody struct {
		Lbrace  token.Pos // position of "{"
		Fields  *FieldList
		Methods []Node
		/*
			blue("blue"),
			red("red");
		*/
		Gen    []CallExpr
		Rbrace token.Pos // position of "}"
	}

	EnumBodySimple struct {
		Lbrace token.Pos // position of "{"
		Fields []*Ident
		Rbrace token.Pos // position of "}"
	}

	// A ClassDecl node represents a class declaration.
	ClassDecl struct {
		Doc         *CommentGroup
		Annotations []*Annotation
		Pub         token.Pos // position of "pub"
		Class       token.Pos // position of "class"
		Name        *Ident
		Extends     token.Pos // position of "extends" or nil
		Parent      []*Ident
		Lbrace      token.Pos // position of "{"
		Fields      *FieldList
		Methods     []Node
		Rbrace      token.Pos // position of "}"
	}

	// A FuncDecl node represents a function declaration.
	FuncDecl struct {
		Doc    *CommentGroup
		Macros []*Annotation
		Mod    *FuncModifier
		Fn     token.Pos // position of "fn"
		Name   *Ident
		Recv   *FieldList
		Throw  bool
		Body   *BlockStmt
	}

	FuncModifier struct {
		Pub    token.Pos // position of "pub"
		Const  token.Pos // position of "const"
		Static token.Pos // position of "static"
	}

	// A ModDecl node represents a mod declaration.
	ModDecl struct {
		Doc    *CommentGroup
		Pub    token.Pos // position of "pub"
		Mod    token.Pos // position of "mod"
		Name   *Ident
		Lbrace token.Pos // position of "{"
		Decls  []Decl
		Rbrace token.Pos // position of "}"
	}

	// A TypeDecl node represents a type declaration.
	TypeDecl struct {
		Doc    *CommentGroup
		Type   token.Pos // position of "type"
		Name   *Ident
		Assign token.Pos // position of "="
		Value  Expr
	}

	// An ExtensionDecl node represents an extension declaration.
	ExtensionDecl struct {
		Doc       *CommentGroup
		Pub       token.Pos // position of "pub"
		Extension token.Pos // position of "extension"
		Name      *Ident
		*ExtensionEnum
		*ExtensionClass
		*ExtensionTrait
		*ExtensionType
		*ExtensionMod
	}

	ExtensionEnum struct {
		Enum token.Pos // position of "enum"
		Body *BlockStmt
	}

	ExtensionClass struct {
		Class token.Pos // position of "class"
		Body  *BlockStmt
	}

	ExtensionTrait struct {
		Trait token.Pos // position of "trait"
		Body  *BlockStmt
	}

	ExtensionType struct {
		Type token.Pos // position of "type"
		Body *BlockStmt
	}

	ExtensionMod struct {
		Mod  token.Pos // position of "mod"
		Body *BlockStmt
	}

	// An DeclareDecl node represents a declare declaration.
	DeclareDecl struct {
		Doc     *CommentGroup
		Declare token.Pos // position of "declare"
		X       Node
	}
)

func (d *TraitDecl) Pos() token.Pos {
	if d.Pub.IsValid() {
		return d.Pub
	}
	return d.Trait
}
func (d *ImplDecl) Pos() token.Pos { return d.Impl }
func (d *EnumDecl) Pos() token.Pos {
	if d.Pub.IsValid() {
		return d.Pub
	}
	return d.Enum
}
func (d *ClassDecl) Pos() token.Pos {
	if d.Pub.IsValid() {
		return d.Pub
	}
	return d.Class
}
func (d *FuncDecl) Pos() token.Pos {
	if d.Mod != nil {
		if d.Mod.Pub.IsValid() {
			return d.Mod.Pub
		}
		if d.Mod.Const.IsValid() {
			return d.Mod.Const
		}
		if d.Mod.Static.IsValid() {
			return d.Mod.Static
		}
	}
	return d.Fn
}
func (d *ModDecl) Pos() token.Pos {
	if d.Pub.IsValid() {
		return d.Pub
	}
	return d.Mod
}
func (d *TypeDecl) Pos() token.Pos      { return d.Type }
func (d *ExtensionDecl) Pos() token.Pos { return d.Extension }
func (d *DeclareDecl) Pos() token.Pos   { return d.Declare }

func (d *TraitDecl) End() token.Pos { return d.Rbrace }
func (d *ImplDecl) End() token.Pos  { return d.Body.Rbrace }
func (d *EnumDecl) End() token.Pos {
	if d.EnumBody != nil {
		return d.EnumBody.Rbrace
	}
	return d.EnumBodySimple.Rbrace
}
func (d *ClassDecl) End() token.Pos { return d.Rbrace }
func (d *FuncDecl) End() token.Pos  { return d.Body.Rbrace }
func (d *ModDecl) End() token.Pos   { return d.Rbrace }
func (d *TypeDecl) End() token.Pos  { return d.Value.End() }
func (d *ExtensionDecl) End() token.Pos {
	switch {
	case d.ExtensionClass != nil:
		return d.ExtensionClass.Body.Rbrace
	case d.ExtensionEnum != nil:
		return d.ExtensionEnum.Body.Rbrace
	case d.ExtensionTrait != nil:
		return d.ExtensionTrait.Body.Rbrace
	case d.ExtensionType != nil:
		return d.ExtensionType.Body.Rbrace
	case d.ExtensionMod != nil:
		return d.ExtensionMod.Body.Rbrace
	}
	return token.NoPos
}
func (d *DeclareDecl) End() token.Pos { return d.X.End() }

func (*TraitDecl) declNode()     {}
func (*ImplDecl) declNode()      {}
func (*EnumDecl) declNode()      {}
func (*ClassDecl) declNode()     {}
func (*FuncDecl) declNode()      {}
func (*ModDecl) declNode()       {}
func (*TypeDecl) declNode()      {}
func (*ExtensionDecl) declNode() {}
func (*DeclareDecl) declNode()   {}

// ----------------------------------------------------------------------------
// Comments

// A Comment node represents a single //-style or /*-style comment.
//
// The Text field contains the comment text without carriage returns (\r) that
// may have been present in the source. Because a comment's end position is
// computed using len(Text), the position reported by End() does not match the
// true source end position for comments containing carriage returns.
type Comment struct {
	Slash token.Pos // position of "/" starting the comment
	Text  string    // comment text (excluding '\n' for //-style comments)
}

func (c *Comment) Pos() token.Pos { return c.Slash }
func (c *Comment) End() token.Pos { return token.Pos(int(c.Slash) + len(c.Text)) }

// A CommentGroup represents a sequence of comments
// with no other tokens and no empty lines between.
type CommentGroup struct {
	List []*Comment // len(List) > 0
}

func (g *CommentGroup) Pos() token.Pos { return g.List[0].Pos() }
func (g *CommentGroup) End() token.Pos { return g.List[len(g.List)-1].End() }

// ----------------------------------------------------------------------------
// Statement

type (
	AssignStmt struct {
		Doc      *CommentGroup
		Scope    token.Token // const | var | let
		ScopePos token.Pos
		Lhs      Expr
		Tok      token.Token // := or =
		Rhs      Expr
	}

	CmdStmt struct {
		Doc  *CommentGroup
		Name Expr
		Recv []Expr
	}

	// An BreakStmt node represents a break statement.
	BreakStmt struct {
		Doc   *CommentGroup
		Break token.Pos // position of "break"
	}

	ContinueStmt struct {
		Doc      *CommentGroup
		Continue token.Pos // position of "continue"
	}

	// A ComptimeStmt node represents a comptime statement.
	ComptimeStmt struct {
		Comptime token.Pos // position of "comptime"
		X        Node
	}

	// unsafe { ... }
	// ${ ... }
	UnsafeStmt struct {
		Doc    *CommentGroup
		Unsafe token.Pos // position of "unsafe" or "$"
		Start  token.Pos // position of `{`
		Text   string
		EndPos token.Pos // position of `}`
	}

	// A BlockStmt node represents a braced statement list.
	BlockStmt struct {
		Lbrace token.Pos // position of "{"
		List   []Stmt
		Rbrace token.Pos // position of "}"
	}

	// An IfStmt node represents an if statement.
	IfStmt struct {
		Doc  *CommentGroup
		If   token.Pos // position of "if"
		Cond Expr
		Body *BlockStmt
		Else Stmt
	}

	// ---------------------------------------------------
	// loop statement

	// A ForeachStmt node represents a foreach statement.
	ForeachStmt struct {
		Loop   token.Pos // position of "loop"
		Lparen token.Pos // position of "("
		Index  Expr
		Value  Expr
		Rparen token.Pos // position of ")"
		In     token.Pos // position of "in"
		Var    Expr
		Body   *BlockStmt
	}

	// A RangeStmt node represents a range statement.
	RangeStmt struct {
		Loop   token.Pos // position of "loop"
		Index  *Ident
		In     token.Pos // position of "in"
		Range  token.Pos // position of "range"
		Lparen token.Pos // position of "("
		RangeClauseExpr
		Rparen token.Pos // position of ")"
		Body   *BlockStmt
	}

	// A WhileStmt node represents a while statement.
	WhileStmt struct {
		Loop token.Pos // position of "loop"
		Cond Expr
		Body *BlockStmt
	}

	// A ForStmt node represents a for statement.
	ForStmt struct {
		Loop   token.Pos // position of "loop"
		Lparen token.Pos // position of "("
		Init   Stmt
		Comma1 token.Pos // position of ";"
		Cond   Expr
		Comma2 token.Pos // position of ";"
		Post   Expr
		Rparen token.Pos // position of ")"
		Body   *BlockStmt
	}

	// A DoWhileStmt node represents a do..while statement.
	DoWhileStmt struct {
		Do     token.Pos // position of "do"
		Body   *BlockStmt
		Loop   token.Pos // position of "loop"
		Lparen token.Pos // position of "("
		Cond   Expr
		Rparen token.Pos // position of ")"
	}

	// A MatchStmt node represents a match statement.
	MatchStmt struct {
		Doc    *CommentGroup
		Match  token.Pos // position of "match"
		Lbrace token.Pos // position of "{"
		Cases  []*CaseClause
		Rbrace token.Pos // position of "}"
	}

	// A CaseClause represents a case of an expression or type match statement.
	CaseClause struct {
		Cond   Expr
		DArrow token.Pos // position of "=>"
		Body   *BlockStmt
		Comma  token.Pos // position of ","
	}

	// A ReturnStmt node represents a return statement.
	ReturnStmt struct {
		Doc    *CommentGroup
		Return token.Pos // position of "return"
		X      Expr
	}

	// A DeferStmt node represents a defer statement.
	DeferStmt struct {
		Doc   *CommentGroup
		Defer token.Pos // position of "defer"
		Call  *CallExpr
	}

	// An ExprStmt node represents a (stand-alone) expression
	// in a statement list.
	ExprStmt struct {
		Doc    *CommentGroup
		Macros []*Annotation
		X      Expr // expression
	}

	// A TryStmt node represents a try statement.
	TryStmt struct {
		Doc     *CommentGroup
		Try     token.Pos // position of "try"
		Body    *BlockStmt
		Catches []*CatchStmt
		Finally *FinallyStmt
	}

	// A CatchStmt node represents a catch statement.
	CatchStmt struct {
		Doc   *CommentGroup
		Catch token.Pos // position of "catch"
		Cond  Expr
		Body  *BlockStmt
	}

	// A FinallyStmt node represents a finally statement.
	FinallyStmt struct {
		Doc     *CommentGroup
		Finally token.Pos // position of "finally"
		Body    *BlockStmt
	}

	// A ThrowStmt node represents a throw statement.
	ThrowStmt struct {
		Doc   *CommentGroup
		Throw token.Pos // position of "throw"
		X     Expr
	}

	// @builtin()
	Annotation struct {
		At     token.Pos // position of "@"
		Name   *Ident
		Lparen token.Pos // position of "("
		Recv   *FieldList
		Rparen token.Pos // position of ")"
	}
)

func (s *AssignStmt) Pos() token.Pos   { return s.ScopePos }
func (s *CmdStmt) Pos() token.Pos      { return s.Name.Pos() }
func (s *ComptimeStmt) Pos() token.Pos { return s.Comptime }
func (s *UnsafeStmt) Pos() token.Pos   { return s.Unsafe }
func (s *TryStmt) Pos() token.Pos      { return s.Try }
func (s *CatchStmt) Pos() token.Pos    { return s.Catch }
func (s *FinallyStmt) Pos() token.Pos  { return s.Finally }
func (s *ExprStmt) Pos() token.Pos     { return s.X.Pos() }
func (s *BlockStmt) Pos() token.Pos    { return s.Lbrace }
func (s *ThrowStmt) Pos() token.Pos    { return s.Throw }
func (s *Annotation) Pos() token.Pos   { return s.At }
func (s *ReturnStmt) Pos() token.Pos   { return s.Return }
func (s *IfStmt) Pos() token.Pos       { return s.If }
func (s *MatchStmt) Pos() token.Pos    { return s.Match }
func (s *DeferStmt) Pos() token.Pos    { return s.Defer }
func (s *WhileStmt) Pos() token.Pos    { return s.Loop }
func (s *ForeachStmt) Pos() token.Pos  { return s.Loop }
func (s *RangeStmt) Pos() token.Pos    { return s.Loop }
func (s *DoWhileStmt) Pos() token.Pos  { return s.Do }
func (s *ForStmt) Pos() token.Pos      { return s.Loop }
func (s *BreakStmt) Pos() token.Pos    { return s.Break }
func (s *ContinueStmt) Pos() token.Pos { return s.Continue }

func (s *AssignStmt) End() token.Pos   { return s.Rhs.End() }
func (s *CmdStmt) End() token.Pos      { return s.Recv[len(s.Recv)-1].End() }
func (s *ComptimeStmt) End() token.Pos { return s.X.End() }
func (s *UnsafeStmt) End() token.Pos   { return s.EndPos }
func (s *TryStmt) End() token.Pos      { return s.Body.Rbrace }
func (s *CatchStmt) End() token.Pos    { return s.Body.Rbrace }
func (s *FinallyStmt) End() token.Pos  { return s.Body.Rbrace }
func (s *ExprStmt) End() token.Pos     { return s.X.End() }
func (s *BlockStmt) End() token.Pos    { return s.Rbrace }
func (s *ThrowStmt) End() token.Pos    { return s.X.End() }
func (s *Annotation) End() token.Pos {
	if s.Rparen.IsValid() {
		return s.Rparen
	}
	return s.Name.End()
}
func (s *ReturnStmt) End() token.Pos   { return s.X.End() }
func (s *IfStmt) End() token.Pos       { return s.Body.Rbrace }
func (s *MatchStmt) End() token.Pos    { return s.Rbrace }
func (s *DeferStmt) End() token.Pos    { return s.Call.Rparen }
func (s *WhileStmt) End() token.Pos    { return s.Body.Rbrace }
func (s *ForeachStmt) End() token.Pos  { return s.Body.Rbrace }
func (s *RangeStmt) End() token.Pos    { return s.Body.Rbrace }
func (s *DoWhileStmt) End() token.Pos  { return s.Rparen }
func (s *ForStmt) End() token.Pos      { return s.Body.Rbrace }
func (s *BreakStmt) End() token.Pos    { return s.Break }
func (s *ContinueStmt) End() token.Pos { return s.Continue }

func (*AssignStmt) stmtNode()   {}
func (*CmdStmt) stmtNode()      {}
func (*ComptimeStmt) stmtNode() {}
func (*UnsafeStmt) stmtNode()   {}
func (*TryStmt) stmtNode()      {}
func (*CatchStmt) stmtNode()    {}
func (*FinallyStmt) stmtNode()  {}
func (*ExprStmt) stmtNode()     {}
func (*BlockStmt) stmtNode()    {}
func (*ThrowStmt) stmtNode()    {}
func (*Annotation) stmtNode()   {}
func (*ReturnStmt) stmtNode()   {}
func (*IfStmt) stmtNode()       {}
func (*MatchStmt) stmtNode()    {}
func (*DeferStmt) stmtNode()    {}
func (*ForStmt) stmtNode()      {}
func (*ForeachStmt) stmtNode()  {}
func (*WhileStmt) stmtNode()    {}
func (*DoWhileStmt) stmtNode()  {}
func (*BreakStmt) stmtNode()    {}
func (*ContinueStmt) stmtNode() {}
func (*RangeStmt) stmtNode()    {}

type (
	// range(1, 5, 0.1)
	RangeClauseExpr struct {
		Start Expr
		End   Expr
		Step  Expr
	}

	// A CallExpr node represents an expression followed by an argument list.
	CallExpr struct {
		Fun     Expr      // function expression
		Lparen  token.Pos // position of "("
		Generic *GenericExpr
		Recv    []Expr    // function arguments; or nil
		Rparen  token.Pos // position of ")"
	}

	// An Ident node represents an identifier.
	Ident struct {
		NamePos token.Pos // identifier position
		Name    string    // identifier name
	}

	// An Ellipsis node stands for the "..." type in a
	// parameter list or the "..." length in an array type.
	Ellipsis struct {
		Ellipsis token.Pos // position of "..."
		Elt      Expr      // ellipsis element type (parameter lists only); or nil
	}

	// A BasicLit node represents a literal of basic type.
	BasicLit struct {
		Kind     token.Token // token.NUM, token.STR, token.IDENT
		Value    string      // TODO literal string; e.g. 42, 0x7f, 3.14, 1e-9, 2.4i, 'a', '\x7f', "foo" or `\m\n\o`
		ValuePos token.Pos   // literal position
	}

	// An IndexExpr node represents an expression followed by an index.
	IndexExpr struct {
		X      Expr      // expression
		Lbrack token.Pos // position of "["
		Index  Expr      // index expression
		Rbrack token.Pos // position of "]"
	}

	// [1, 2]
	// An IndexListExpr node represents an expression followed by multiple
	// indices.
	IndexListExpr struct {
		X       Expr   // expression
		Indices []Expr // index expressions
	}

	// A SliceExpr node represents an expression followed by slice indices.
	SliceExpr struct {
		X      Expr // expression
		Low    Expr // begin of slice range; or nil
		High   Expr // end of slice range; or nil
		Max    Expr // maximum capacity of slice; or nil
		Slice3 bool // true if 3-index slice (2 colons present)
	}

	// A BinaryExpr node represents a binary expression.
	BinaryExpr struct {
		X     Expr        // left operand
		OpPos token.Pos   // position of Op
		Op    token.Token // operator
		Y     Expr        // right operand
	}

	// A KeyValueExpr node represents (key : value) pairs
	// in composite literals.
	KeyValueExpr struct {
		Key   Expr
		Value Expr
	}

	// An IncDecExpr node represents an increase or decrease expression.
	IncDecExpr struct {
		Pre    bool
		X      Expr
		Tok    token.Token // INC or DEC
		TokPos token.Pos   // position of "++" or "--"
	}

	// A NewDelExpr node represents a new or delete expression.
	NewDelExpr struct {
		XPos token.Pos // position of "new" or "delete"
		X    Expr
	}

	// str | num | bool & user
	TypeAnnotation struct {
		List []Expr
	}

	// ${x} or $x
	RefExpr struct {
		Dollar token.Pos // position of "$"
		Lbrace token.Pos // position of "{"
		X      Expr
		Rbrace token.Pos // position of "}"
	}

	// $ raw $
	UnsafeExpr struct {
		Dollar token.Pos // position of "$"
		Div    token.Pos // position of "/"
		Text   string
		EndPos token.Pos // position of "/"
	}

	// <T, U>
	GenericExpr struct {
		Lt    token.Pos // position of "<"
		Types []*Ident
		Gt    token.Pos // position of ">"
	}

	// X.Y
	SelectExpr struct {
		X   Expr
		Dot token.Pos // position of "."
		Y   Expr
	}

	// A UnaryExpr node represents a unary expression.
	UnaryExpr struct {
		OpPos token.Pos   // position of Op
		Op    token.Token // operator
		X     Expr        // operand
	}

	ComptimeExpr struct {
		Comptime token.Pos // position of comptime
		X        Expr
	}
)

func (x *IndexExpr) Pos() token.Pos      { return x.X.Pos() }
func (x *BinaryExpr) Pos() token.Pos     { return x.X.Pos() }
func (x *CallExpr) Pos() token.Pos       { return x.Fun.Pos() }
func (x *Ident) Pos() token.Pos          { return x.NamePos }
func (x *BasicLit) Pos() token.Pos       { return x.ValuePos }
func (x *NewDelExpr) Pos() token.Pos     { return x.XPos }
func (x *TypeAnnotation) Pos() token.Pos { return x.List[0].Pos() }
func (x *UnsafeExpr) Pos() token.Pos     { return x.Dollar }
func (x *RefExpr) Pos() token.Pos        { return x.Dollar }
func (x *IncDecExpr) Pos() token.Pos     { return x.X.Pos() }
func (x *SelectExpr) Pos() token.Pos     { return x.X.Pos() }
func (x *UnaryExpr) Pos() token.Pos      { return x.OpPos }
func (x *ComptimeExpr) Pos() token.Pos   { return x.Comptime }

func (x *IndexExpr) End() token.Pos      { return x.Rbrack }
func (x *BinaryExpr) End() token.Pos     { return x.Y.End() }
func (x *CallExpr) End() token.Pos       { return x.Rparen }
func (x *Ident) End() token.Pos          { return token.Pos(int(x.NamePos) + len(x.Name)) }
func (x *BasicLit) End() token.Pos       { return token.Pos(int(x.ValuePos) + len(x.Value)) }
func (x *NewDelExpr) End() token.Pos     { return x.X.End() }
func (x *TypeAnnotation) End() token.Pos { return x.List[len(x.List)-1].End() }
func (x *UnsafeExpr) End() token.Pos     { return x.EndPos }
func (x *RefExpr) End() token.Pos {
	if x.Rbrace.IsValid() {
		return x.Rbrace
	}
	return x.X.End()
}
func (x *IncDecExpr) End() token.Pos {
	return x.TokPos + 2
}
func (x *SelectExpr) End() token.Pos   { return x.Y.End() }
func (x *UnaryExpr) End() token.Pos    { return x.X.End() }
func (x *ComptimeExpr) End() token.Pos { return x.X.End() }

func (x *Ident) String() string {
	return x.Name
}

func (x *BasicLit) String() string {
	if x.Kind == token.STR {
		return fmt.Sprintf("\"%s\"", x.Value)
	}
	return x.Value
}

func (x *IndexExpr) String() string {
	return fmt.Sprintf("%s[%s]", x.X, x.Index)
}

func (x *IndexListExpr) String() string {
	indices := []string{}
	for _, index := range x.Indices {
		indices = append(indices, index.String())
	}
	return fmt.Sprintf("%s[%s]", x.X, strings.Join(indices, ", "))
}

func (x *SliceExpr) String() string {
	return fmt.Sprintf("%s[%s..%s]", x.X, x.Low, x.High)
}

func (x *BinaryExpr) String() string {
	return fmt.Sprintf("%s %s %s", x.X, x.Op, x.Y)
}

func (x *KeyValueExpr) String() string {
	return fmt.Sprintf("%s: %s", x.Key, x.Value)
}

func (x *IncDecExpr) String() string {
	if x.Pre {
		return fmt.Sprintf("%s%s", x.Tok, x.X)
	}
	return fmt.Sprintf("%s%s", x.X, x.Tok)
}

func (x *NewDelExpr) String() string {
	return fmt.Sprintf("%s %s", token.NEW, x.X)
}

func (x *TypeAnnotation) String() string {
	types := []string{}
	for _, t := range x.List {
		types = append(types, t.String())
	}
	return fmt.Sprintf("<%s>", strings.Join(types, ", "))
}

func (x *UnsafeExpr) String() string {
	return fmt.Sprintf("%s %s", token.UNSAFE, x.Text)
}

func (x *RefExpr) String() string {
	return fmt.Sprintf("$%s", x.X)
}

func (x *GenericExpr) String() string {
	types := []string{}
	for _, t := range x.Types {
		types = append(types, t.String())
	}
	return fmt.Sprintf("<%s>", strings.Join(types, ", "))
}

func (x *SelectExpr) String() string {
	return fmt.Sprintf("%s.%s", x.X, x.Y)
}

func (x *UnaryExpr) String() string {
	return fmt.Sprintf("%s%s", x.Op, x.X)
}

func (x *CallExpr) String() string {
	args := []string{}
	for _, arg := range x.Recv {
		args = append(args, arg.String())
	}
	return fmt.Sprintf("%s(%s)", x.Fun, strings.Join(args, ", "))
}

func (x *ComptimeExpr) String() string {
	return fmt.Sprintf("comptime { %s }", x.X)
}

func (*BinaryExpr) exprNode()     {}
func (*CallExpr) exprNode()       {}
func (*Ident) exprNode()          {}
func (*BasicLit) exprNode()       {}
func (*NewDelExpr) exprNode()     {}
func (*TypeAnnotation) exprNode() {}
func (*UnsafeExpr) exprNode()     {}
func (*RefExpr) exprNode()        {}
func (*IncDecExpr) exprNode()     {}
func (*SelectExpr) exprNode()     {}
func (*IndexExpr) exprNode()      {}
func (*UnaryExpr) exprNode()      {}
func (*ComptimeExpr) exprNode()   {}

type Modifier int

const (
	M_NONE = 0
	M_PUB  = 1 << iota
	M_FINAL
	M_CONST
	M_STATIC
	M_REQUIRED
	M_ALL = M_PUB | M_FINAL | M_CONST | M_STATIC | M_REQUIRED
)

func (m Modifier) IsNone() bool {
	return m == M_NONE
}

func (m Modifier) HasPub() bool {
	return m&M_PUB != 0
}

func (m Modifier) HasFinal() bool {
	return m&M_FINAL != 0
}

func (m Modifier) HasConst() bool {
	return m&M_CONST != 0
}

func (m Modifier) HasStatic() bool {
	return m&M_STATIC != 0
}

func (m Modifier) HasRequired() bool {
	return m&M_REQUIRED != 0
}

func (m Modifier) IsAll() bool {
	return m == M_ALL
}

// (pub | const | final) x: str | num = defaultValue
type Field struct {
	Doc    *CommentGroup
	Macros []*Annotation
	Mod    Modifier
	Name   *Ident
	Type   Expr
	Assign token.Pos // position of '='
	Value  Expr      // default value
}

type FieldList struct {
	Opening token.Pos // position of opening parenthesis/brace/bracket, if any
	List    []*Field  // field list; or nil
	Closing token.Pos // position of closing parenthesis/brace/bracket, if any
}

func (f *FieldList) NumFields() int { return len(f.List) }
func (f *FieldList) Pos() token.Pos { return f.Opening }
func (f *FieldList) End() token.Pos { return f.Closing }

type (
	// A Import node represents an import statement.
	Import struct {
		Doc       *CommentGroup
		ImportPos token.Pos // position of "import"
		*ImportAll
		*ImportSingle
		*ImportMulti
	}

	// import * as <alias> from <path>
	ImportAll struct {
		Mul   token.Pos // position of "*"
		As    token.Pos // position of "as"
		Alias string
		From  token.Pos // position of "from"
		Path  string
	}

	// import <path> as <alias>
	ImportSingle struct {
		Path  string
		As    token.Pos // position of "as"
		Alias string
	}

	// import { <p1> as <alias1>, <p2> } from <path>
	ImportMulti struct {
		Lbrace token.Pos // position of "{"
		List   []*ImportField
		Rbrace token.Pos // position of "}"
		From   token.Pos // position of "from"
		Path   string
	}

	// <field> as <alias>
	ImportField struct {
		Field string
		As    token.Pos // position of "as"
		Alias string
	}
)

func (s *Import) Pos() token.Pos { return s.ImportPos }
func (s *Import) End() token.Pos {
	if s.ImportAll != nil {
		return token.Pos(int(s.ImportAll.From) + len(s.ImportAll.Path))
	}
	if s.ImportSingle != nil {
		if len(s.ImportSingle.Path) == 0 {
			return token.Pos(int(s.ImportPos) + len(s.ImportSingle.Path))
		}
		return token.Pos(int(s.ImportSingle.As) + len(s.ImportSingle.Path))
	}
	return token.Pos(int(s.ImportMulti.From) + len(s.ImportMulti.Path))
}

type File struct {
	Docs []*CommentGroup
	Name *Ident // filename

	Imports map[string]*Import
	Stmts   []Stmt
	Decls   []Decl
}

func (f *File) Pos() token.Pos { return token.NoPos }
func (f *File) End() token.Pos { return token.NoPos }

func (*File) declNode() {}
func (*File) stmtNode() {}

type Package struct {
	Name *Ident

	Files map[string]*File
}

func (p *Package) Pos() token.Pos { return token.NoPos }
func (p *Package) End() token.Pos { return token.NoPos }
