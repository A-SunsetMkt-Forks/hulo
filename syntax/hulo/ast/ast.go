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

// All statement nodes implement the Stmt interface.
type Stmt interface {
	Node
	stmtNode()
}

// All expression nodes implement the Expr interface.
type Expr interface {
	Node
	exprNode()
	String() string
}

// ----------------------------------------------------------------------------
// Declarations

type (
	// A TraitDecl node represents a trait declaration.
	TraitDecl struct {
		Docs      *CommentGroup // associated documentation
		Modifiers []Modifier
		Trait     token.Pos  // position of "trait" keyword
		Name      *Ident     // trait name
		Lbrace    token.Pos  // position of "{"
		Fields    *FieldList // list of method signatures
		Rbrace    token.Pos  // position of "}"
	}

	// An ImplDecl node represents an implementation of a trait for a type.
	ImplDecl struct {
		Docs  *CommentGroup // associated documentation
		Impl  token.Pos     // position of "impl" keyword
		Trait *Ident        // trait being implemented
		For   token.Pos     // position of "for" keyword
		*ImplDeclBody
		*ImplDeclBinding
	}

	// An ImplDeclBody represents the implementation for a single class.
	ImplDeclBody struct {
		Class *Ident     // class implementing the trait
		Body  *BlockStmt // implementation body
	}

	// An ImplDeclBinding represents the implementation for multiple classes.
	ImplDeclBinding struct {
		Classes []*Ident // classes implementing the trait
	}

	// An EnumDecl node represents an enum declaration.
	EnumDecl struct {
		Docs       *CommentGroup // associated documentation
		Decs       []*Decorator  // decorators
		Modifiers  []Modifier    // modifiers array (final, const, pub, etc.)
		Enum       token.Pos     // position of "enum" keyword
		Name       *Ident        // enum name
		Lt         token.Pos     // position of "<", if any
		TypeParams []Expr        // type parameters for a generic enum
		Gt         token.Pos     // position of ">", if any
		Body       EnumBody      // enum body
	}

	// EnumBody represents the body of an enum declaration.
	EnumBody interface {
		Node
		enumBodyNode()
	}

	// BasicEnumBody represents a simple enum with basic values.
	// enum Status { Pending, Approved, Rejected }
	// enum HttpCode { OK = 200, NotFound = 404 }
	BasicEnumBody struct {
		Lbrace token.Pos    // position of "{"
		Values []*EnumValue // enum values
		Rbrace token.Pos    // position of "}"
	}

	// AssociatedEnumBody represents an enum with associated values.
	// enum Protocol { port: num; tcp(6), udp(17) }
	AssociatedEnumBody struct {
		Lbrace token.Pos    // position of "{"
		Fields *FieldList   // associated fields declaration
		Values []*EnumValue // enum values with data
		Rbrace token.Pos    // position of "}"
	}

	// ADTEnumBody represents an algebraic data type enum.
	// enum NetworkPacket { TCP { src_port: num, dst_port: num }, UDP { port: num } }
	ADTEnumBody struct {
		Lbrace   token.Pos      // position of "{"
		Variants []*EnumVariant // enum variants
		Methods  []Stmt         // methods
		Rbrace   token.Pos      // position of "}"
	}

	// EnumValue represents a single enum value.
	EnumValue struct {
		Docs   *CommentGroup // associated documentation
		Name   *Ident        // value name
		Assign token.Pos     // position of "=", if any
		Value  Expr          // value expression, if any
		Lparen token.Pos     // position of "(", for associated values
		Data   []Expr        // associated data, if any
		Rparen token.Pos     // position of ")"
		Comma  token.Pos     // position of ","
	}

	// EnumVariant represents a variant in an ADT enum.
	EnumVariant struct {
		Docs   *CommentGroup // associated documentation
		Name   *Ident        // variant name
		Lbrace token.Pos     // position of "{"
		Fields *FieldList    // variant fields
		Rbrace token.Pos     // position of "}"
		Comma  token.Pos     // position of ","
	}

	// An EnumMember represents a member of an enum.
	EnumMember struct{}

	// An EnumMethod represents a method within an enum.
	EnumMethod struct {
		Mod Modifier
		Fn  token.Pos // position of "fn"
	}

	// A ClassDecl node represents a class declaration.
	ClassDecl struct {
		Docs       *CommentGroup      // associated documentation
		Decs       []*Decorator       // decorators
		Pub        token.Pos          // position of "pub" keyword, if any
		Class      token.Pos          // position of "class" keyword
		Name       *Ident             // class name
		Lt         token.Pos          // position of "<", if any
		TypeParams []Expr             // type parameters for a generic class
		Gt         token.Pos          // position of ">", if any
		Extends    token.Pos          // position of "extends" keyword, if any
		Parent     *Ident             // parent classes
		Lbrace     token.Pos          // position of "{"
		Fields     *FieldList         // list of fields
		Ctors      []*ConstructorDecl // list of constructors
		Methods    []*FuncDecl        // list of methods
		Rbrace     token.Pos          // position of "}"
	}

	// A FuncDecl node represents a function declaration.
	FuncDecl struct {
		Docs       *CommentGroup // associated documentation
		Decs       []*Decorator  // decorators
		Modifiers  []Modifier    // modifiers array (const, comptime, pub, etc.)
		Fn         token.Pos     // position of "fn" keyword
		Name       *Ident        // function name
		Lt         token.Pos     // position of "<", if any
		TypeParams []Expr        // type parameters for a generic function
		Gt         token.Pos     // position of ">", if any
		Recv       []Expr        // function parameters
		Throw      bool          // if the function can throw an exception
		Type       Expr          // result type
		Body       *BlockStmt    // function body
	}

	// A ConstructorDecl node represents a constructor declaration.
	// e.g., const Constants(pi: num, e: num) $this.PI = $pi, $this.E = $e {}
	ConstructorDecl struct {
		Docs       *CommentGroup // associated documentation
		Decs       []*Decorator  // decorators
		Modifiers  []Modifier    // modifiers array (const, pub, etc.)
		ClsName    *Ident        // class name (same as class name)
		Name       *Ident        // constructor name
		Lparen     token.Pos     // position of "("
		Recv       []Expr        // constructor parameters
		Rparen     token.Pos     // position of ")"
		InitFields []Expr        // init fields
		Body       *BlockStmt    // constructor body
	}

	// An OperatorDecl node represents an operator declaration.
	OperatorDecl struct {
		Docs     *CommentGroup // associated documentation
		Decs     []*Decorator  // decorators
		Operator token.Pos     // position of "operator" keyword
		Name     token.Token   // operator token
		Lparen   token.Pos     // position of "("
		Params   []Expr        // parameters
		Rparen   token.Pos     // position of ")"
		Arrow    token.Pos     // position of "->"
		Type     Expr          // result type
		Body     *BlockStmt    // function body
	}

	// A ModDecl node represents a module declaration.
	ModDecl struct {
		Docs   *CommentGroup // associated documentation
		Pub    token.Pos     // position of "pub" keyword, if any
		Mod    token.Pos     // position of "mod" keyword
		Name   *Ident        // module name
		Lbrace token.Pos     // position of "{"
		Body   *BlockStmt    // module body
		Rbrace token.Pos     // position of "}"
	}

	// A TypeDecl node represents a type declaration (type alias).
	TypeDecl struct {
		Docs   *CommentGroup // associated documentation
		Type   token.Pos     // position of "type" keyword
		Name   *Ident        // type name
		Assign token.Pos     // position of "="
		Value  Expr          // type value
	}

	// An ExtensionDecl node represents an extension declaration.
	ExtensionDecl struct {
		Docs      *CommentGroup // associated documentation
		Extension token.Pos     // position of "extension" keyword
		Name      *Ident
		Body      ExtensionBody // the extension body
	}

	// ExtensionBody represents the body of an extension declaration.
	ExtensionBody interface {
		Node
		extensionBodyNode()
	}

	// An ExtensionEnum node represents an enum extension.
	ExtensionEnum struct {
		Enum token.Pos // position of "enum"
		Body *BlockStmt
	}

	// An ExtensionClass node represents a class extension.
	ExtensionClass struct {
		Class token.Pos // position of "class"
		Body  *BlockStmt
	}

	// An ExtensionTrait node represents a trait extension.
	ExtensionTrait struct {
		Trait token.Pos // position of "trait"
		Body  *BlockStmt
	}

	// An ExtensionType node represents a type extension.
	ExtensionType struct {
		Type token.Pos // position of "type"
		Body *BlockStmt
	}

	// An ExtensionMod node represents a module extension.
	ExtensionMod struct {
		Mod  token.Pos // position of "mod"
		Body *BlockStmt
	}

	// A DeclareDecl node represents a declare declaration.
	DeclareDecl struct {
		Docs    *CommentGroup
		Declare token.Pos // position of "declare"
		X       Node
	}

	// A UseDecl node represents a use declaration.
	UseDecl struct {
		Docs   *CommentGroup
		Use    token.Pos // position of "use"
		Lhs    Expr
		Assign token.Pos // position of "="
		Rhs    Expr
	}

	// An ExternDecl represents a foreign function interface declaration.
	ExternDecl struct {
		Extern token.Pos // position of "extern"
		List   []Expr
	}

	// A GetAccessor represents a getter in a property.
	GetAccessor struct {
		Get  token.Pos // position of "get"
		Name *Ident
		Body *BlockStmt
	}

	// A SetAccessor represents a setter in a property.
	SetAccessor struct {
		Set  token.Pos // position of "set"
		Name *Ident
		Body *BlockStmt
	}
)

func (d *TraitDecl) Pos() token.Pos {
	if len(d.Modifiers) > 0 {
		return d.Modifiers[0].Pos()
	}
	return d.Trait
}
func (d *ImplDecl) Pos() token.Pos { return d.Impl }
func (d *EnumDecl) Pos() token.Pos {
	if len(d.Modifiers) > 0 {
		return d.Modifiers[0].Pos()
	}
	return d.Enum
}
func (d *ClassDecl) Pos() token.Pos {
	if d.Pub.IsValid() {
		return d.Pub
	}
	return d.Class
}
func (d *ConstructorDecl) Pos() token.Pos {
	if len(d.Modifiers) > 0 {
		return d.Modifiers[0].Pos()
	}
	return d.ClsName.Pos()
}
func (d *FuncDecl) Pos() token.Pos {
	if len(d.Modifiers) > 0 {
		return d.Modifiers[0].Pos()
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
func (d *UseDecl) Pos() token.Pos       { return d.Use }
func (d *OperatorDecl) Pos() token.Pos  { return d.Operator }
func (d *ExternDecl) Pos() token.Pos    { return d.Extern }
func (d *GetAccessor) Pos() token.Pos   { return d.Get }
func (d *SetAccessor) Pos() token.Pos   { return d.Set }

func (d *TraitDecl) End() token.Pos { return d.Rbrace }
func (d *ImplDecl) End() token.Pos  { return d.Body.Rbrace }
func (d *EnumDecl) End() token.Pos {
	if d.Body != nil {
		return d.Body.End()
	}
	return d.Enum
}
func (d *ClassDecl) End() token.Pos { return d.Rbrace }
func (d *FuncDecl) End() token.Pos  { return d.Body.Rbrace }
func (d *ConstructorDecl) End() token.Pos {
	if d.Body != nil {
		return d.Body.End()
	}
	return d.Rparen
}
func (d *ModDecl) End() token.Pos  { return d.Rbrace }
func (d *TypeDecl) End() token.Pos { return d.Value.End() }
func (d *ExtensionDecl) End() token.Pos {
	return d.Body.End()
}
func (d *DeclareDecl) End() token.Pos  { return d.X.End() }
func (d *UseDecl) End() token.Pos      { return d.Rhs.End() }
func (d *OperatorDecl) End() token.Pos { return d.Body.Rbrace }
func (d *ExternDecl) End() token.Pos   { return d.List[len(d.List)-1].End() }
func (d *GetAccessor) End() token.Pos  { return d.Body.Rbrace }
func (d *SetAccessor) End() token.Pos  { return d.Body.Rbrace }

func (*TraitDecl) stmtNode()       {}
func (*ImplDecl) stmtNode()        {}
func (*EnumDecl) stmtNode()        {}
func (*ClassDecl) stmtNode()       {}
func (*FuncDecl) stmtNode()        {}
func (*ConstructorDecl) stmtNode() {}
func (*ModDecl) stmtNode()         {}
func (*TypeDecl) stmtNode()        {}
func (*ExtensionDecl) stmtNode()   {}
func (*DeclareDecl) stmtNode()     {}
func (*UseDecl) stmtNode()         {}
func (*OperatorDecl) stmtNode()    {}
func (*ExternDecl) stmtNode()      {}
func (*GetAccessor) stmtNode()     {}
func (*SetAccessor) stmtNode()     {}

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
	// An AssignStmt node represents an assignment or a short variable declaration.
	AssignStmt struct {
		Docs     *CommentGroup
		Scope    token.Token // const | var | let
		ScopePos token.Pos
		Lhs      Expr
		Colon    token.Pos // position of ":"
		Type     Expr
		Tok      token.Token // assignment token, e.g., := or =
		Rhs      Expr
	}

	// A CmdStmt node represents a command statement.
	CmdStmt struct {
		Docs *CommentGroup
		Name Expr
		Recv []Expr
	}

	// A BreakStmt node represents a break statement.
	BreakStmt struct {
		Docs  *CommentGroup
		Break token.Pos // position of "break"
		Label *Ident
	}

	// A ContinueStmt node represents a continue statement.
	ContinueStmt struct {
		Docs     *CommentGroup
		Continue token.Pos // position of "continue"
	}

	// A ComptimeStmt node represents a compile-time execution block.
	ComptimeStmt struct {
		Comptime token.Pos // position of "comptime"
		X        Stmt
	}

	// An UnsafeStmt node represents an unsafe block.
	// unsafe { ... } or ${ ... }
	UnsafeStmt struct {
		Docs   *CommentGroup
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
		Docs *CommentGroup
		If   token.Pos // position of "if"
		Cond Expr
		Body *BlockStmt
		Else Stmt
	}

	// ---------------------------------------------------
	// loop statement

	// A LabelStmt node represents a labeled statement.
	LabelStmt struct {
		Name  *Ident
		Colon token.Pos
	}

	// A ForeachStmt node represents a for-each statement.
	ForeachStmt struct {
		Loop   token.Pos // position of "loop"
		Lparen token.Pos // position of "("
		Index  Expr
		Value  Expr
		Rparen token.Pos   // position of ")"
		Tok    token.Token // Token.IN or Token.OF
		TokPos token.Pos   // position of "in" or "of"
		Var    Expr
		Body   *BlockStmt
	}

	// A ForInStmt node represents a for-range statement.
	ForInStmt struct {
		Loop  token.Pos // position of "loop"
		Index *Ident
		In    token.Pos // position of "in"
		RangeExpr
		Body *BlockStmt
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
		Docs    *CommentGroup
		Match   token.Pos // position of "match"
		Lbrace  token.Pos // position of "{"
		Cases   []*CaseClause
		Default *CaseClause
		Rbrace  token.Pos // position of "}"
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
		Docs   *CommentGroup
		Return token.Pos // position of "return"
		X      Expr
	}

	// A DeferStmt node represents a defer statement.
	// defer {...} or defer fn()
	DeferStmt struct {
		Docs  *CommentGroup
		Defer token.Pos // position of "defer"
		X     Node
	}

	// An ExprStmt node represents a (stand-alone) expression
	// in a statement list.
	ExprStmt struct {
		Docs *CommentGroup
		Decs []*Decorator
		X    Expr // expression
	}

	// A TryStmt node represents a try-catch-finally statement.
	TryStmt struct {
		Docs    *CommentGroup
		Try     token.Pos // position of "try"
		Body    *BlockStmt
		Catches []*CatchClause
		Finally *FinallyStmt
	}

	// A CatchClause node represents a catch block.
	CatchClause struct {
		Docs   *CommentGroup
		Catch  token.Pos // position of "catch"
		Lparen token.Pos // position of "("
		Cond   Expr
		Rparen token.Pos // position of ")"
		Body   *BlockStmt
	}

	// A FinallyStmt node represents a finally block.
	FinallyStmt struct {
		Docs    *CommentGroup
		Finally token.Pos // position of "finally"
		Body    *BlockStmt
	}

	// A ThrowStmt node represents a throw statement.
	ThrowStmt struct {
		Docs  *CommentGroup
		Throw token.Pos // position of "throw"
		X     Expr
	}

	// A Decorator node represents a decorator.
	Decorator struct {
		At   token.Pos // position of "@"
		Name *Ident
		Recv []Expr
	}
)

func (s *AssignStmt) Pos() token.Pos   { return s.ScopePos }
func (s *CmdStmt) Pos() token.Pos      { return s.Name.Pos() }
func (s *ComptimeStmt) Pos() token.Pos { return s.Comptime }
func (s *UnsafeStmt) Pos() token.Pos   { return s.Unsafe }
func (s *TryStmt) Pos() token.Pos      { return s.Try }
func (s *CatchClause) Pos() token.Pos  { return s.Catch }
func (s *FinallyStmt) Pos() token.Pos  { return s.Finally }
func (s *ExprStmt) Pos() token.Pos     { return s.X.Pos() }
func (s *BlockStmt) Pos() token.Pos    { return s.Lbrace }
func (s *ThrowStmt) Pos() token.Pos    { return s.Throw }
func (s *Decorator) Pos() token.Pos    { return s.At }
func (s *ReturnStmt) Pos() token.Pos   { return s.Return }
func (s *IfStmt) Pos() token.Pos       { return s.If }
func (s *MatchStmt) Pos() token.Pos    { return s.Match }
func (s *DeferStmt) Pos() token.Pos    { return s.Defer }
func (s *WhileStmt) Pos() token.Pos    { return s.Loop }
func (s *ForeachStmt) Pos() token.Pos  { return s.Loop }
func (s *ForInStmt) Pos() token.Pos    { return s.Loop }
func (s *DoWhileStmt) Pos() token.Pos  { return s.Do }
func (s *ForStmt) Pos() token.Pos      { return s.Loop }
func (s *BreakStmt) Pos() token.Pos    { return s.Break }
func (s *ContinueStmt) Pos() token.Pos { return s.Continue }
func (s *LabelStmt) Pos() token.Pos    { return s.Name.Pos() }

func (s *AssignStmt) End() token.Pos   { return s.Rhs.End() }
func (s *CmdStmt) End() token.Pos      { return s.Recv[len(s.Recv)-1].End() }
func (s *ComptimeStmt) End() token.Pos { return s.X.End() }
func (s *UnsafeStmt) End() token.Pos   { return s.EndPos }
func (s *TryStmt) End() token.Pos      { return s.Body.Rbrace }
func (s *CatchClause) End() token.Pos  { return s.Body.Rbrace }
func (s *FinallyStmt) End() token.Pos  { return s.Body.Rbrace }
func (s *ExprStmt) End() token.Pos     { return s.X.End() }
func (s *BlockStmt) End() token.Pos    { return s.Rbrace }
func (s *ThrowStmt) End() token.Pos    { return s.X.End() }
func (s *Decorator) End() token.Pos {
	if len(s.Recv) > 0 {
		return s.Recv[len(s.Recv)-1].End()
	}
	return s.Name.End()
}
func (s *ReturnStmt) End() token.Pos   { return s.X.End() }
func (s *IfStmt) End() token.Pos       { return s.Body.Rbrace }
func (s *MatchStmt) End() token.Pos    { return s.Rbrace }
func (s *DeferStmt) End() token.Pos    { return s.X.End() }
func (s *WhileStmt) End() token.Pos    { return s.Body.Rbrace }
func (s *ForeachStmt) End() token.Pos  { return s.Body.Rbrace }
func (s *ForInStmt) End() token.Pos    { return s.Body.Rbrace }
func (s *DoWhileStmt) End() token.Pos  { return s.Rparen }
func (s *ForStmt) End() token.Pos      { return s.Body.Rbrace }
func (s *BreakStmt) End() token.Pos    { return s.Break }
func (s *ContinueStmt) End() token.Pos { return s.Continue }
func (s *LabelStmt) End() token.Pos    { return s.Colon }

func (*AssignStmt) stmtNode()   {}
func (*CmdStmt) stmtNode()      {}
func (*ComptimeStmt) stmtNode() {}
func (*UnsafeStmt) stmtNode()   {}
func (*TryStmt) stmtNode()      {}
func (*CatchClause) stmtNode()  {}
func (*FinallyStmt) stmtNode()  {}
func (*ExprStmt) stmtNode()     {}
func (*BlockStmt) stmtNode()    {}
func (*ThrowStmt) stmtNode()    {}
func (*Decorator) stmtNode()    {}
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
func (*ForInStmt) stmtNode()    {}
func (*LabelStmt) stmtNode()    {}

type (
	// A RangeExpr node represents a range expression.
	// e.g., 1..5..0.1 or 1..0.5..-0.1
	RangeExpr struct {
		Start     Expr
		DblColon1 token.Pos // position of ".."
		End_      Expr
		DblColon2 token.Pos // position of ".."
		Step      Expr
	}

	// A CallExpr node represents an expression followed by an argument list.
	CallExpr struct {
		Fun        Expr      // function expression
		Lt         token.Pos // position of "<"
		TypeParams []Expr    // type arguments for a generic function call
		Gt         token.Pos // position of ">"
		Lparen     token.Pos // position of "("
		Recv       []Expr    // function arguments; or nil
		Rparen     token.Pos // position of ")"
	}

	// A CmdSubstExpr node represents a command substitution. e.g. $(...)
	CmdSubstExpr struct {
		Dollar token.Pos
		Lparen token.Pos
		Fun    Expr
		Recv   []Expr
		Rparen token.Pos
	}

	// An Ident node represents an identifier.
	Ident struct {
		NamePos token.Pos // identifier position
		Name    string    // identifier name
	}

	// A BasicLit node represents a literal of basic type.
	BasicLit struct {
		Kind     token.Token // token.NUM, token.STR, token.IDENT
		Value    string      // literal string; e.g. 42, 0x7f, 3.14, "foo"
		ValuePos token.Pos   // literal position
	}

	// An IndexExpr node represents an expression followed by an index.
	IndexExpr struct {
		X      Expr      // expression
		Quest  bool      // true if optional chaining is used, e.g., a?[i]
		Lbrack token.Pos // position of "["
		Index  Expr      // index expression
		Rbrack token.Pos // position of "]"
	}

	// An IndexListExpr node represents an expression followed by multiple
	// indices. e.g., a[i, j]
	IndexListExpr struct {
		X       Expr   // expression
		Indices []Expr // index expressions
	}

	// A SliceExpr node represents an expression followed by slice indices.
	SliceExpr struct {
		X       Expr      // expression
		Lbrack  token.Pos // position of "["
		Low     Expr      // begin of slice range; or nil
		DblCol  token.Pos // position of ".."
		High    Expr      // end of slice range; or nil
		DblCol2 token.Pos // position of ".."
		Max     Expr      // maximum capacity of slice; or nil
		Rbrack  token.Pos // position of "]"
	}

	// A BinaryExpr node represents a binary expression.
	BinaryExpr struct {
		X     Expr        // left operand
		OpPos token.Pos   // position of Op
		Op    token.Token // operator
		Y     Expr        // right operand
	}

	// An IncDecExpr node represents an increment or decrement expression.
	IncDecExpr struct {
		Pre    bool // true if it's a prefix operator, false for postfix
		X      Expr
		Tok    token.Token // INC or DEC
		TokPos token.Pos   // position of "++" or "--"
	}

	// A NewDelExpr node represents a new or delete expression.
	NewDelExpr struct {
		TokPos token.Pos   // position of "new" or "delete"
		Tok    token.Token // Token.NEW or Token.DELETE
		X      Expr
	}

	// A RefExpr node represents a reference expression. e.g., ${x} or $x
	RefExpr struct {
		Dollar token.Pos // position of "$"
		Lbrace token.Pos // position of "{"
		X      Expr
		Rbrace token.Pos // position of "}"
	}

	// A SelectExpr node represents a selector expression. e.g., X.Y
	SelectExpr struct {
		X     Expr
		Quest bool      // true if optional chaining is used, e.g., X?.Y
		Dot   token.Pos // position of "."
		Y     Expr
	}

	// A ModAccessExpr node represents a module access expression. e.g., X::Y
	ModAccessExpr struct {
		X        Expr
		DblColon token.Pos // position of "::"
		Y        Expr
	}

	// A UnaryExpr node represents a unary expression.
	UnaryExpr struct {
		OpPos token.Pos   // position of Op
		Op    token.Token // operator
		X     Expr        // operand
	}

	// A ComptimeExpr node represents a compile-time expression.
	ComptimeExpr struct {
		Comptime token.Pos // position of comptime
		X        Expr
	}

	// A LiteralType represents a basic literal type.
	LiteralType BasicLit

	// A UnionType node represents a union type. e.g., T | U
	UnionType struct {
		Types []Expr
	}

	// An IntersectionType node represents an intersection type. e.g., T & U
	IntersectionType struct {
		Types []Expr
	}

	// A TypeLiteral node represents a type literal. e.g., { a: int, b: string }
	TypeLiteral struct {
		Memebers []Expr
	}

	// A TypeReference node represents a reference to a type.
	// e.g., MyMap<T, U>
	TypeReference struct {
		Name       Expr
		TypeParams []Expr
	}

	// An ArrayType node represents an array type. e.g., T[]
	ArrayType struct {
		Name Expr
	}

	// A FunctionType node represents a function type. e.g., (int, string) -> bool
	FunctionType struct {
		Lparen  token.Pos
		Recv    []Expr
		Rparent token.Pos
		RetVal  Expr
	}

	// A TupleType node represents a tuple type. e.g., [int, string]
	TupleType struct {
		Types []Expr
	}

	// A SetType node represents a set type. e.g., {int}
	SetType struct {
		Types []Expr
	}

	// A ConditionalType node represents a conditional type. e.g., T extends U ? X : Y
	ConditionalType struct {
		CheckType   Expr
		Extends     token.Pos // position of "extends"
		ExtendsType Expr
		Quest       token.Pos // position of "?"
		TrueType    Expr
		Colon       token.Pos // position of ":"
		FalseType   Expr
	}

	// An InferType node represents a type inference variable in a conditional type. e.g. infer T
	InferType struct {
		Infer token.Pos // position of "infer"
		X     Expr
	}

	// A KeyValueExpr node represents (key : value) pairs
	// in composite literals.
	KeyValueExpr struct {
		Docs  *CommentGroup
		Decs  []*Decorator
		Mod   Modifier
		Key   Expr
		Quest bool      // true if the key is optional
		Colon token.Pos // position of ":"
		Value Expr
	}

	// A TypeParameter node represents a type parameter in a generic declaration.
	// e.g. RW extends Readable + Writable
	TypeParameter struct {
		Name        Expr
		Extends     token.Pos // position of "extends"
		Constraints []Expr
	}

	// A ConditionalExpr node represents a conditional expression. e.g., cond ? a : b
	ConditionalExpr struct {
		Cond      Expr
		Quest     token.Pos // position of "?"
		WhenTrue  Expr
		Colon     token.Pos // position of ":"
		WhneFalse Expr
	}

	// A CascadeExpr node represents a cascade expression. e.g., x..a()..b()
	CascadeExpr struct {
		X      Expr
		DblDot token.Pos // position of ".."
		Y      Expr
	}

	// A NullableType node represents a nullable type. e.g., T?
	NullableType struct {
		X Expr
	}

	// A NumericLiteral node represents a numeric literal.
	NumericLiteral struct {
		ValuePos token.Pos
		Value    string
	}

	// A StringLiteral node represents a string literal.
	StringLiteral struct {
		ValuePos token.Pos
		Value    string
	}

	// A TrueLiteral node represents the boolean literal 'true'.
	TrueLiteral struct {
		ValuePos token.Pos
	}

	// A FalseLiteral node represents the boolean literal 'false'.
	FalseLiteral struct {
		ValuePos token.Pos
	}

	// An AnyLiteral node represents the 'any' type literal.
	AnyLiteral struct {
		ValuePos token.Pos
	}

	// A NullLiteral node represents the 'null' literal.
	NullLiteral struct {
		ValuePos token.Pos
	}

	// A SetLiteralExpr node represents a set literal.
	SetLiteralExpr struct {
		Lbrace token.Pos // position of "{"
		Elems  []Expr
		Rbrace token.Pos // position of "}"
	}

	// A TupleLiteralExpr node represents a tuple literal.
	TupleLiteralExpr struct {
		Lbrack token.Pos // position of "["
		Elems  []Expr
		Rbrack token.Pos // position of "]"
	}

	// A NamedTupleLiteralExpr node represents a named tuple literal.
	NamedTupleLiteralExpr struct {
		Lparen token.Pos // position of "("
		Elems  []Expr
		Rparen token.Pos // position of ")"
	}

	// An ArrayLiteralExpr node represents an array literal.
	ArrayLiteralExpr struct {
		Lbrack token.Pos // position of "["
		Elems  []Expr
		Rbrack token.Pos // position of "]"
	}

	// An ObjectLiteralExpr node represents an object literal.
	ObjectLiteralExpr struct {
		Lbrace token.Pos // position of "{"
		Props  []Expr
		Rbrace token.Pos // position of "}"
	}

	// A NamedObjectLiteralExpr node represents a composite literal expression,
	// used for instantiating a struct or class.
	// e.g. User{name: "root", pwd: "123456"}
	NamedObjectLiteralExpr struct {
		Name   Expr      // Type name, e.g., 'User'.
		Lbrace token.Pos // Position of "{"
		Props  []Expr    // List of elements, typically KeyValueExpr.
		Rbrace token.Pos // Position of "}"
	}

	// A TypeofExpr node represents a typeof expression.
	// e.g. typeof X
	TypeofExpr struct {
		Typeof token.Pos // position of "typeof"
		X      Expr
	}

	// An AsExpr node represents a type assertion or conversion.
	// e.g. X as Y
	AsExpr struct {
		X  Expr
		As token.Pos
		Y  Expr
	}

	// A Parameter node represents a parameter in a function declaration.
	// e.g. required x: type = value
	Parameter struct {
		TokPos   token.Pos // position of "required" or "..."
		Modifier Modifier
		Name     *Ident
		Colon    token.Pos // position of ":"
		Type     Expr
		Assign   token.Pos // position of "="
		Value    Expr      // default value
	}

	// ANamedParameters node represents a set of named parameters.
	NamedParameters struct {
		Lbrace token.Pos // position of "{"
		Params []Expr
		Rbrace token.Pos // position of "}"
	}

	// A LambdaExpr node represents a lambda expression.
	// TODO: This might be converted to a regular anonymous function during analysis.
	LambdaExpr struct {
		Lparen token.Pos // position of "("
		Recv   []Expr
		Rparen token.Pos // position of ")"
		DArrow token.Pos // position of "=>"
		Body   *BlockStmt
	}
)

func (x *IndexExpr) Pos() token.Pos             { return x.X.Pos() }
func (x *BinaryExpr) Pos() token.Pos            { return x.X.Pos() }
func (x *RangeExpr) Pos() token.Pos             { return x.Start.Pos() }
func (x *CallExpr) Pos() token.Pos              { return x.Fun.Pos() }
func (x *CmdSubstExpr) Pos() token.Pos          { return x.Dollar }
func (x *Ident) Pos() token.Pos                 { return x.NamePos }
func (x *BasicLit) Pos() token.Pos              { return x.ValuePos }
func (x *NewDelExpr) Pos() token.Pos            { return x.TokPos }
func (x *RefExpr) Pos() token.Pos               { return x.Dollar }
func (x *IncDecExpr) Pos() token.Pos            { return x.X.Pos() }
func (x *SelectExpr) Pos() token.Pos            { return x.X.Pos() }
func (x *ModAccessExpr) Pos() token.Pos         { return x.X.Pos() }
func (x *UnaryExpr) Pos() token.Pos             { return x.OpPos }
func (x *ComptimeExpr) Pos() token.Pos          { return x.Comptime }
func (x *TypeParameter) Pos() token.Pos         { return x.Name.Pos() }
func (x *UnionType) Pos() token.Pos             { return x.Types[0].Pos() }
func (x *IntersectionType) Pos() token.Pos      { return x.Types[0].Pos() }
func (x *TypeLiteral) Pos() token.Pos           { return x.Memebers[0].Pos() }
func (x *TypeReference) Pos() token.Pos         { return x.Name.Pos() }
func (x *ArrayType) Pos() token.Pos             { return x.Name.Pos() }
func (x *FunctionType) Pos() token.Pos          { return x.Lparen }
func (x *TupleType) Pos() token.Pos             { return x.Types[0].Pos() }
func (x *SetType) Pos() token.Pos               { return x.Types[0].Pos() }
func (x *ConditionalExpr) Pos() token.Pos       { return x.Cond.Pos() }
func (x *CascadeExpr) Pos() token.Pos           { return x.X.Pos() }
func (x *ArrayLiteralExpr) Pos() token.Pos      { return x.Lbrack }
func (x *SetLiteralExpr) Pos() token.Pos        { return x.Lbrace }
func (x *TupleLiteralExpr) Pos() token.Pos      { return x.Lbrack }
func (x *NamedTupleLiteralExpr) Pos() token.Pos { return x.Lparen }
func (x *NumericLiteral) Pos() token.Pos        { return x.ValuePos }
func (x *StringLiteral) Pos() token.Pos         { return x.ValuePos }
func (x *AnyLiteral) Pos() token.Pos            { return x.ValuePos }
func (x *NullLiteral) Pos() token.Pos           { return x.ValuePos }
func (x *TrueLiteral) Pos() token.Pos           { return x.ValuePos }
func (x *FalseLiteral) Pos() token.Pos          { return x.ValuePos }
func (x *ObjectLiteralExpr) Pos() token.Pos     { return x.Lbrace }
func (x *SliceExpr) Pos() token.Pos             { return x.X.Pos() }
func (x *TypeofExpr) Pos() token.Pos            { return x.Typeof }
func (x *AsExpr) Pos() token.Pos                { return x.X.Pos() }
func (x *Parameter) Pos() token.Pos             { return x.Name.Pos() }
func (x *NamedParameters) Pos() token.Pos       { return x.Lbrace }
func (x *KeyValueExpr) Pos() token.Pos          { return x.Key.Pos() }
func (x *NullableType) Pos() token.Pos          { return x.X.Pos() }
func (x *NamedObjectLiteralExpr) Pos() token.Pos {
	if x.Name != nil {
		return x.Name.Pos()
	}
	return x.Lbrace
}
func (x *LambdaExpr) Pos() token.Pos      { return x.Lparen }
func (x *ConditionalType) Pos() token.Pos { return x.CheckType.Pos() }
func (x *InferType) Pos() token.Pos       { return x.Infer }

func (x *IndexExpr) End() token.Pos  { return x.Rbrack }
func (x *BinaryExpr) End() token.Pos { return x.Y.End() }
func (x *RangeExpr) End() token.Pos {
	switch {
	case x.Step != nil:
		return x.Step.End()
	case x.End_ != nil:
		return x.End_.End()
	}
	return x.Step.End()
}
func (x *CallExpr) End() token.Pos     { return x.Rparen }
func (x *CmdSubstExpr) End() token.Pos { return x.Rparen }
func (x *Ident) End() token.Pos        { return token.Pos(int(x.NamePos) + len(x.Name)) }
func (x *BasicLit) End() token.Pos     { return token.Pos(int(x.ValuePos) + len(x.Value)) }
func (x *NewDelExpr) End() token.Pos   { return x.X.End() }
func (x *RefExpr) End() token.Pos {
	if x.Rbrace.IsValid() {
		return x.Rbrace
	}
	return x.X.End()
}
func (x *IncDecExpr) End() token.Pos {
	return x.TokPos + 2
}
func (x *SelectExpr) End() token.Pos    { return x.Y.End() }
func (x *ModAccessExpr) End() token.Pos { return x.Y.End() }
func (x *UnaryExpr) End() token.Pos     { return x.X.End() }
func (x *ComptimeExpr) End() token.Pos  { return x.X.End() }
func (x *TypeParameter) End() token.Pos {
	if len(x.Constraints) != 0 {
		return x.Constraints[len(x.Constraints)-1].End()
	}
	return x.Name.End()
}
func (x *UnionType) End() token.Pos             { return x.Types[0].End() }
func (x *IntersectionType) End() token.Pos      { return x.Types[0].End() }
func (x *TypeLiteral) End() token.Pos           { return x.Memebers[0].End() }
func (x *TypeReference) End() token.Pos         { return x.Name.End() }
func (x *ArrayType) End() token.Pos             { return x.Name.End() }
func (x *FunctionType) End() token.Pos          { return x.Lparen }
func (x *TupleType) End() token.Pos             { return x.Types[0].End() }
func (x *SetType) End() token.Pos               { return x.Types[0].End() }
func (x *ConditionalExpr) End() token.Pos       { return x.WhneFalse.End() }
func (x *CascadeExpr) End() token.Pos           { return x.Y.End() }
func (x *ArrayLiteralExpr) End() token.Pos      { return x.Rbrack }
func (x *SetLiteralExpr) End() token.Pos        { return x.Rbrace }
func (x *TupleLiteralExpr) End() token.Pos      { return x.Rbrack }
func (x *NamedTupleLiteralExpr) End() token.Pos { return x.Rparen }
func (x *NumericLiteral) End() token.Pos        { return token.Pos(int(x.ValuePos) + len(x.Value)) }
func (x *StringLiteral) End() token.Pos         { return token.Pos(int(x.ValuePos) + len(x.Value)) }
func (x *AnyLiteral) End() token.Pos            { return token.Pos(int(x.ValuePos) + 3) }
func (x *NullLiteral) End() token.Pos           { return token.Pos(int(x.ValuePos) + 4) }
func (x *TrueLiteral) End() token.Pos           { return token.Pos(int(x.ValuePos) + 4) }
func (x *FalseLiteral) End() token.Pos          { return token.Pos(int(x.ValuePos) + 4) }
func (x *ObjectLiteralExpr) End() token.Pos     { return x.Rbrace }
func (x *SliceExpr) End() token.Pos             { return x.Rbrack }
func (x *TypeofExpr) End() token.Pos            { return x.X.End() }
func (x *AsExpr) End() token.Pos                { return x.Y.End() }
func (x *Parameter) End() token.Pos {
	if x.Type != nil {
		return x.Type.End()
	}
	return x.Name.End()
}
func (x *NamedParameters) End() token.Pos        { return x.Rbrace }
func (x *KeyValueExpr) End() token.Pos           { return x.Value.End() }
func (x *NullableType) End() token.Pos           { return x.X.End() }
func (x *NamedObjectLiteralExpr) End() token.Pos { return x.Rbrace }
func (x *LambdaExpr) End() token.Pos             { return x.Body.Rbrace }
func (x *ConditionalType) End() token.Pos        { return x.FalseType.End() }
func (x *InferType) End() token.Pos              { return x.X.End() }

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
	switch {
	case x.Low != nil && x.High != nil && x.Max != nil:
		return fmt.Sprintf("%s[%s..%s..%s]", x.X, x.Low, x.High, x.Max)
	case x.Low != nil && x.High != nil:
		return fmt.Sprintf("%s[%s..%s]", x.X, x.Low, x.High)
	case x.Low != nil:
		return fmt.Sprintf("%s[%s..]", x.X, x.Low)
	case x.High != nil:
		return fmt.Sprintf("%s[..%s]", x.X, x.High)
	}
	panic("unreached")
}

func (x *BinaryExpr) String() string {
	return fmt.Sprintf("%s %s %s", x.X, x.Op, x.Y)
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

func (x *RefExpr) String() string {
	return fmt.Sprintf("$%s", x.X)
}

func (x *SelectExpr) String() string {
	if x.Quest {
		return fmt.Sprintf("%s?.%s", x.X, x.Y)
	}
	return fmt.Sprintf("%s.%s", x.X, x.Y)
}

func (x *ModAccessExpr) String() string {
	return fmt.Sprintf("%s::%s", x.X, x.Y)
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

func (x *CmdSubstExpr) String() string {
	args := []string{}
	for _, arg := range x.Recv {
		args = append(args, arg.String())
	}
	return fmt.Sprintf("$(%s %s)", x.Fun, strings.Join(args, " "))
}

func (x *ComptimeExpr) String() string {
	return fmt.Sprintf("comptime { %s }", x.X)
}

func (x *TypeParameter) String() string {
	if x.Constraints != nil {
		constraints := []string{}
		for _, c := range x.Constraints {
			constraints = append(constraints, c.String())
		}
		return fmt.Sprintf("%s extends %s", x.Name, strings.Join(constraints, " + "))
	}
	return x.Name.String()
}

func (x *UnionType) String() string {
	expr := []string{}
	for _, t := range x.Types {
		expr = append(expr, t.String())
	}
	return strings.Join(expr, " | ")
}

func (x *IntersectionType) String() string {
	expr := []string{}
	for _, t := range x.Types {
		expr = append(expr, t.String())
	}
	return strings.Join(expr, " & ")
}

func (x *TypeLiteral) String() string {
	mems := []string{}
	for _, t := range x.Memebers {
		mems = append(mems, t.String())
	}
	return strings.Join(mems, "\n")
}

func (x *TypeReference) String() string {
	exprs := []string{}
	for _, e := range x.TypeParams {
		exprs = append(exprs, e.String())
	}
	if x.TypeParams != nil {
		return fmt.Sprintf("%s<%s>", x.Name, strings.Join(exprs, ", "))
	}
	return x.Name.String()
}

func (x *ArrayType) String() string {
	return fmt.Sprintf("%s[]", x.Name)
}

func (x *FunctionType) String() string {
	recv := []string{}
	for _, t := range x.Recv {
		recv = append(recv, t.String())
	}
	return fmt.Sprintf("(%s) -> %s", strings.Join(recv, ", "), x.RetVal)
}

func (x *TupleType) String() string {
	exprs := []string{}
	for _, e := range x.Types {
		exprs = append(exprs, e.String())
	}
	return fmt.Sprintf("[%s]", strings.Join(exprs, ", "))
}

func (x *SetType) String() string {
	exprs := []string{}
	for _, e := range x.Types {
		exprs = append(exprs, e.String())
	}
	return fmt.Sprintf("(%s)", strings.Join(exprs, ", "))
}

func (x *ConditionalExpr) String() string {
	return fmt.Sprintf("%s ? %s : %s", x.Cond, x.WhenTrue, x.WhneFalse)
}

func (x *CascadeExpr) String() string {
	return fmt.Sprintf("%s..%s", x.X, x.Y)
}

func (x *ArrayLiteralExpr) String() string {
	expr := []string{}
	for _, e := range x.Elems {
		expr = append(expr, e.String())
	}
	return fmt.Sprintf("[%s]", strings.Join(expr, ", "))
}

func (x *SetLiteralExpr) String() string {
	expr := []string{}
	for _, e := range x.Elems {
		expr = append(expr, e.String())
	}
	return fmt.Sprintf("{%s}", strings.Join(expr, ", "))
}

func (x *TupleLiteralExpr) String() string {
	expr := []string{}
	for _, e := range x.Elems {
		expr = append(expr, e.String())
	}
	return fmt.Sprintf("[%s]", strings.Join(expr, ", "))
}

func (x *NamedTupleLiteralExpr) String() string {
	expr := []string{}
	for _, e := range x.Elems {
		expr = append(expr, e.String())
	}
	return fmt.Sprintf("(%s)", strings.Join(expr, ", "))
}

func (x *NumericLiteral) String() string {
	return x.Value
}

func (x *StringLiteral) String() string {
	return fmt.Sprintf(`"%s"`, x.Value)
}

func (x *AnyLiteral) String() string   { return "any" }
func (x *NullLiteral) String() string  { return "null" }
func (x *TrueLiteral) String() string  { return "true" }
func (x *FalseLiteral) String() string { return "false" }

func (x *TypeofExpr) String() string {
	return fmt.Sprintf("typeof %s", x.X)
}

func (x *AsExpr) String() string {
	return fmt.Sprintf("%s as %s", x.X, x.Y)
}

func (x *Parameter) String() string {
	ret := x.Name.String()
	// if x.Variadic {
	// 	ret = fmt.Sprintf("...%s", ret)
	// }
	// if x.Required {
	// 	ret = fmt.Sprintf("required %s", ret)
	// }
	if x.Type != nil {
		ret = fmt.Sprintf("%s: %s", ret, x.Type)
	}
	if x.Value != nil {
		ret = fmt.Sprintf("%s = %s", ret, x.Value)
	}
	return ret
}

func (x *NamedParameters) String() string {
	params := []string{}
	for _, p := range x.Params {
		params = append(params, p.String())
	}
	return fmt.Sprintf("{%s}", strings.Join(params, ", "))
}

func (x *ObjectLiteralExpr) String() string {
	props := []string{}
	for _, prop := range x.Props {
		props = append(props, prop.String())
	}
	if len(props) == 0 {
		return "{}"
	}
	return fmt.Sprintf("{ %s }", strings.Join(props, ", "))
}

func (x *KeyValueExpr) String() string {
	return fmt.Sprintf("%s: %s", x.Key, x.Value)
}

func (x *NullableType) String() string {
	return fmt.Sprintf("%s?", x.X)
}

func (x *NamedObjectLiteralExpr) String() string {
	elems := []string{}
	for _, el := range x.Props {
		elems = append(elems, el.String())
	}
	if x.Name != nil {
		return fmt.Sprintf("%s{%s}", x.Name.String(), strings.Join(elems, ", "))
	}
	return fmt.Sprintf("{%s}", strings.Join(elems, ", "))
}

func (*LambdaExpr) String() string { panic("method no support") }

func (x *ConditionalType) String() string {
	return fmt.Sprintf("%s extends %s ? %s : %s", x.CheckType, x.ExtendsType, x.TrueType, x.FalseType)
}

func (x *InferType) String() string {
	return fmt.Sprintf("infer %s", x.X)
}

func (*BinaryExpr) exprNode()             {}
func (*CallExpr) exprNode()               {}
func (*CmdSubstExpr) exprNode()           {}
func (*Ident) exprNode()                  {}
func (*BasicLit) exprNode()               {}
func (*NewDelExpr) exprNode()             {}
func (*RefExpr) exprNode()                {}
func (*IncDecExpr) exprNode()             {}
func (*SelectExpr) exprNode()             {}
func (*ModAccessExpr) exprNode()          {}
func (*IndexExpr) exprNode()              {}
func (*UnaryExpr) exprNode()              {}
func (*ComptimeExpr) exprNode()           {}
func (*UnionType) exprNode()              {}
func (*IntersectionType) exprNode()       {}
func (*TypeLiteral) exprNode()            {}
func (*TypeReference) exprNode()          {}
func (*ArrayType) exprNode()              {}
func (*FunctionType) exprNode()           {}
func (*TupleType) exprNode()              {}
func (*SetType) exprNode()                {}
func (*ConditionalExpr) exprNode()        {}
func (*CascadeExpr) exprNode()            {}
func (*ArrayLiteralExpr) exprNode()       {}
func (*TypeParameter) exprNode()          {}
func (*NumericLiteral) exprNode()         {}
func (*StringLiteral) exprNode()          {}
func (*AnyLiteral) exprNode()             {}
func (*NullLiteral) exprNode()            {}
func (*TrueLiteral) exprNode()            {}
func (*FalseLiteral) exprNode()           {}
func (*ObjectLiteralExpr) exprNode()      {}
func (*SliceExpr) exprNode()              {}
func (*TypeofExpr) exprNode()             {}
func (*AsExpr) exprNode()                 {}
func (*Parameter) exprNode()              {}
func (*NamedParameters) exprNode()        {}
func (*KeyValueExpr) exprNode()           {}
func (*NullableType) exprNode()           {}
func (*NamedObjectLiteralExpr) exprNode() {}
func (*LambdaExpr) exprNode()             {}
func (*ConditionalType) exprNode()        {}
func (*InferType) exprNode()              {}

// TODO add position infromation
// Modifier represents a set of modifiers that can be applied to declarations
type Modifier interface {
	Node
	Kind() ModifierKind
}

var (
	_ Modifier = (*FinalModifier)(nil)
	_ Modifier = (*ConstModifier)(nil)
	_ Modifier = (*PubModifier)(nil)
	_ Modifier = (*StaticModifier)(nil)
	_ Modifier = (*RequiredModifier)(nil)
	_ Modifier = (*ReadonlyModifier)(nil)
	_ Modifier = (*EllipsisModifier)(nil)
)

// ModifierKind represents the type of modifier.
type ModifierKind uint8

// Modifier kind constants
const (
	ModKindNone ModifierKind = iota
	ModKindFinal
	ModKindConst
	ModKindPub
	ModKindStatic
	ModKindRequired
	ModKindReadonly
	ModKindEllipsis
)

// FinalModifier represents a "final" modifier.
type FinalModifier struct {
	Final token.Pos // position of "final" keyword
}

// ConstModifier represents a "const" modifier.
type ConstModifier struct {
	Const token.Pos // position of "const" keyword
}

// PubModifier represents a "pub" modifier.
type PubModifier struct {
	Pub token.Pos // position of "pub" keyword
}

// StaticModifier represents a "static" modifier.
type StaticModifier struct {
	Static token.Pos // position of "static" keyword
}

// RequiredModifier represents a "required" modifier.
type RequiredModifier struct {
	Required token.Pos // position of "required" keyword
}

// ReadonlyModifier represents a "readonly" modifier.
type ReadonlyModifier struct {
	Readonly token.Pos // position of "readonly" keyword
}

// EllipsisModifier represents a "..." modifier.
type EllipsisModifier struct {
	Ellipsis token.Pos // position of "..."
}

// Implementation methods for modifier types
func (m *FinalModifier) Pos() token.Pos     { return m.Final }
func (m *FinalModifier) End() token.Pos     { return token.Pos(int(m.Final) + 5) } // "final"
func (m *FinalModifier) Kind() ModifierKind { return ModKindFinal }

func (m *ConstModifier) Pos() token.Pos     { return m.Const }
func (m *ConstModifier) End() token.Pos     { return token.Pos(int(m.Const) + 5) } // "const"
func (m *ConstModifier) Kind() ModifierKind { return ModKindConst }

func (m *PubModifier) Pos() token.Pos     { return m.Pub }
func (m *PubModifier) End() token.Pos     { return token.Pos(int(m.Pub) + 3) } // "pub"
func (m *PubModifier) Kind() ModifierKind { return ModKindPub }

func (m *StaticModifier) Pos() token.Pos     { return m.Static }
func (m *StaticModifier) End() token.Pos     { return token.Pos(int(m.Static) + 6) } // "static"
func (m *StaticModifier) Kind() ModifierKind { return ModKindStatic }

func (m *RequiredModifier) Pos() token.Pos     { return m.Required }
func (m *RequiredModifier) End() token.Pos     { return token.Pos(int(m.Required) + 8) } // "required"
func (m *RequiredModifier) Kind() ModifierKind { return ModKindRequired }

func (m *ReadonlyModifier) Pos() token.Pos     { return m.Readonly }
func (m *ReadonlyModifier) End() token.Pos     { return token.Pos(int(m.Readonly) + 8) } // "readonly"
func (m *ReadonlyModifier) Kind() ModifierKind { return ModKindReadonly }

func (m *EllipsisModifier) Pos() token.Pos     { return m.Ellipsis }
func (m *EllipsisModifier) End() token.Pos     { return token.Pos(int(m.Ellipsis) + 3) } // "..."
func (m *EllipsisModifier) Kind() ModifierKind { return ModKindEllipsis }

// A Field represents a field in a struct or a parameter in a function.
// e.g., (pub | const | final) x: str | num = defaultValue
type Field struct {
	Docs      *CommentGroup
	Decs      []*Decorator
	Modifiers []Modifier
	Name      *Ident
	Colon     token.Pos // position of ":"
	Type      Expr
	Assign    token.Pos // position of '='
	Value     Expr      // default value
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
	// An Import node represents an import statement.
	Import struct {
		Docs      *CommentGroup
		ImportPos token.Pos // position of "import"
		*ImportAll
		*ImportSingle
		*ImportMulti
	}

	// An ImportAll node represents an 'import * as ...' declaration.
	// import * as <alias> from <path>
	ImportAll struct {
		Mul   token.Pos // position of "*"
		As    token.Pos // position of "as"
		Alias string
		From  token.Pos // position of "from"
		Path  string
	}

	// An ImportSingle node represents an 'import ... as ...' declaration.
	// import <path> as <alias>
	ImportSingle struct {
		Path  string
		As    token.Pos // position of "as"
		Alias string
	}

	// An ImportMulti node represents an 'import { ... } from ...' declaration.
	// import { <p1> as <alias1>, <p2> } from <path>
	ImportMulti struct {
		Lbrace token.Pos // position of "{"
		List   []*ImportField
		Rbrace token.Pos // position of "}"
		From   token.Pos // position of "from"
		Path   string
	}

	// An ImportField node represents a single imported entity in a multi-import declaration.
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

// A File node represents a Hulo source file.
type File struct {
	Docs    []*CommentGroup
	Name    *Ident // filename
	Imports map[string]*Import
	Stmts   []Stmt
	Decls   []Stmt
}

func (f *File) Pos() token.Pos { return token.NoPos }
func (f *File) End() token.Pos { return token.NoPos }

func (*File) stmtNode() {}

// A Package node represents a collection of source files.
type Package struct {
	Name *Ident

	Files map[string]*File
}

func (p *Package) Pos() token.Pos { return token.NoPos }
func (p *Package) End() token.Pos { return token.NoPos }

func (e *ExtensionEnum) Pos() token.Pos  { return e.Enum }
func (e *ExtensionEnum) End() token.Pos  { return e.Body.Rbrace }
func (e *ExtensionClass) Pos() token.Pos { return e.Class }
func (e *ExtensionClass) End() token.Pos { return e.Body.Rbrace }
func (e *ExtensionTrait) Pos() token.Pos { return e.Trait }
func (e *ExtensionTrait) End() token.Pos { return e.Body.Rbrace }
func (e *ExtensionType) Pos() token.Pos  { return e.Type }
func (e *ExtensionType) End() token.Pos  { return e.Body.Rbrace }
func (e *ExtensionMod) Pos() token.Pos   { return e.Mod }
func (e *ExtensionMod) End() token.Pos   { return e.Body.Rbrace }

func (*ExtensionEnum) extensionBodyNode()  {}
func (*ExtensionClass) extensionBodyNode() {}
func (*ExtensionTrait) extensionBodyNode() {}
func (*ExtensionType) extensionBodyNode()  {}
func (*ExtensionMod) extensionBodyNode()   {}

func (b *BasicEnumBody) Pos() token.Pos      { return b.Lbrace }
func (b *BasicEnumBody) End() token.Pos      { return b.Rbrace }
func (b *AssociatedEnumBody) Pos() token.Pos { return b.Lbrace }
func (b *AssociatedEnumBody) End() token.Pos { return b.Rbrace }
func (b *ADTEnumBody) Pos() token.Pos        { return b.Lbrace }
func (b *ADTEnumBody) End() token.Pos        { return b.Rbrace }

func (v *EnumValue) Pos() token.Pos { return v.Name.Pos() }
func (v *EnumValue) End() token.Pos {
	if v.Rparen.IsValid() {
		return v.Rparen
	}
	if v.Value != nil {
		return v.Value.End()
	}
	return v.Name.End()
}

func (v *EnumVariant) Pos() token.Pos { return v.Name.Pos() }
func (v *EnumVariant) End() token.Pos { return v.Rbrace }

func (*BasicEnumBody) enumBodyNode()      {}
func (*AssociatedEnumBody) enumBodyNode() {}
func (*ADTEnumBody) enumBodyNode()        {}
