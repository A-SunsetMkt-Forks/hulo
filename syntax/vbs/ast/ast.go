// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package ast declares the types used to represent syntax trees for VBScript.
package ast

import (
	"fmt"

	"github.com/hulo-lang/hulo/syntax/vbs/token"
)

// All node types implement the Node interface.
type Node interface {
	Pos() token.Pos // position of first character belonging to the node
	End() token.Pos // position of first character immediately after the node
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
}

// A CommentGroup represents a sequence of comments
// with no other tokens and no empty lines between.
type CommentGroup struct {
	List []*Comment // len(List) > 0
}

func (g *CommentGroup) Pos() token.Pos { return g.List[0].Pos() }
func (g *CommentGroup) End() token.Pos { return g.List[len(g.List)-1].End() }

// A Comment represents a single comment.
type Comment struct {
	TokPos token.Pos   // position of comment token
	Tok    token.Token // ' or Rem
	Text   string      // comment text (excluding '\n' for line comments)
}

func (c *Comment) Pos() token.Pos { return c.TokPos }
func (c *Comment) End() token.Pos { return token.Pos(int(c.TokPos) + len(c.Text)) }

type (
	// A SubDecl node represents a sub declaration.
	SubDecl struct {
		Doc    *CommentGroup // associated documentation; or nil
		Mod    token.Token   // modifier token (PUBLIC, PRIVATE)
		ModPos token.Pos     // position of Mod
		Sub    token.Pos     // position of "Sub"
		Name   *Ident        // sub name
		Recv   []*Field      // receiver parameters; or nil
		Body   *BlockStmt    // sub body
		EndSub token.Pos     // position of "End Sub"
	}

	// A FuncDecl node represents a function declaration.
	FuncDecl struct {
		Doc    *CommentGroup // associated documentation; or nil
		Mod    token.Token   // modifier token (PUBLIC, PRIVATE)
		ModPos token.Pos     // position of Mod
		// Default is used only with the Public keyword in a Class block
		// to indicate that the Function procedure is the default method for the class.
		Default  token.Pos  // position of "Default"
		Function token.Pos  // position of "Function"
		Name     *Ident     // function name
		Recv     []*Field   // receiver parameters; or nil
		Body     *BlockStmt // function body
		EndFunc  token.Pos  // position of "End Function"
	}

	// A PropertyStmt node represents a property declaration.
	PropertyStmt struct {
		Doc         *CommentGroup // associated documentation; or nil
		Mod         token.Token   // modifier token (PUBLIC, PRIVATE)
		ModPos      token.Pos     // position of Mod
		Property    token.Pos     // position of "Property"
		Name        *Ident        // property name
		LParen      token.Pos     // position of "("
		Recv        []*Field      // receiver parameters; or nil
		RParen      token.Pos     // position of ")"
		Body        *BlockStmt    // property body
		EndProverty token.Pos     // position of "End Property"
	}

	// A PropertySetStmt node represents a property set declaration.
	PropertySetStmt struct {
		Doc         *CommentGroup // associated documentation; or nil
		Mod         token.Token   // modifier token (PUBLIC, PRIVATE)
		ModPos      token.Pos     // position of Mod
		Property    token.Pos     // position of "Property"
		Set         token.Pos     // position of "Set"
		Name        *Ident        // property name
		LParen      token.Pos     // position of "("
		Recv        []*Field      // receiver parameters; or nil
		RParen      token.Pos     // position of ")"
		Body        *BlockStmt    // property body
		EndProverty token.Pos     // position of "End Property"
	}

	// A PropertyLetStmt node represents a property let declaration.
	PropertyLetStmt struct {
		Doc         *CommentGroup // associated documentation; or nil
		Mod         token.Token   // modifier token (PUBLIC, PRIVATE)
		ModPos      token.Pos     // position of Mod
		Property    token.Pos     // position of "Property"
		Let         token.Pos     // position of "Let"
		Name        *Ident        // property name
		LParen      token.Pos     // position of "("
		Recv        []*Field      // receiver parameters; or nil
		RParen      token.Pos     // position of ")"
		Body        *BlockStmt    // property body
		EndProverty token.Pos     // position of "End Property"
	}

	// A PropertyGetStmt node represents a property get declaration.
	PropertyGetStmt struct {
		Doc         *CommentGroup // associated documentation; or nil
		Mod         token.Token   // modifier token (PUBLIC, PRIVATE)
		ModPos      token.Pos     // position of Mod
		Property    token.Pos     // position of "Property"
		Get         token.Pos     // position of "Get"
		Name        *Ident        // property name
		LParen      token.Pos     // position of "("
		Recv        []*Field      // receiver parameters; or nil
		RParen      token.Pos     // position of ")"
		Body        *BlockStmt    // property body
		EndProverty token.Pos     // position of "End Property"
	}

	// A ClassDecl node represents a class declaration.
	ClassDecl struct {
		Doc    *CommentGroup // associated documentation; or nil
		Mod    token.Token   // modifier token (PUBLIC, PRIVATE)
		ModPos token.Pos     // position of Mod
		Class  token.Pos     // position of "Class"
		Name   *Ident        // class name
		// Stmts contains all statements in the class body:
		// dim, func, member, property, assign statements
		Stmts    []Stmt    // class body statements
		EndClass token.Pos // position of "End Class"
	}

	// A DimDecl node represents a variable declaration.
	DimDecl struct {
		Doc   *CommentGroup // associated documentation; or nil
		Dim   token.Pos     // position of "Dim"
		List  []Expr        // variable list
		Colon token.Pos     // position of ":"
		Set   *SetStmt      // optional initial value
	}

	// A ReDimDecl node represents a variable redimension declaration.
	ReDimDecl struct {
		Doc      *CommentGroup // associated documentation; or nil
		ReDim    token.Pos     // position of "ReDim"
		Preserve token.Pos     // position of "Preserve"
		List     []Expr        // variable list
	}
)

func (d *SubDecl) Pos() token.Pos {
	if d.ModPos.IsValid() {
		return d.ModPos
	}
	return d.Sub
}
func (d *PropertyGetStmt) Pos() token.Pos {
	if d.ModPos.IsValid() {
		return d.ModPos
	}
	return d.Property
}
func (d *PropertyLetStmt) Pos() token.Pos {
	if d.ModPos.IsValid() {
		return d.ModPos
	}
	return d.Property
}
func (d *PropertySetStmt) Pos() token.Pos {
	if d.ModPos.IsValid() {
		return d.ModPos
	}
	return d.Property
}
func (d *FuncDecl) Pos() token.Pos {
	if d.ModPos.IsValid() {
		return d.ModPos
	}
	return d.Function
}
func (d *ClassDecl) Pos() token.Pos {
	if d.ModPos.IsValid() {
		return d.ModPos
	}
	return d.Class
}
func (s *DimDecl) Pos() token.Pos   { return s.Dim }
func (s *ReDimDecl) Pos() token.Pos { return s.ReDim }

func (d *SubDecl) End() token.Pos         { return d.EndSub }
func (d *PropertyLetStmt) End() token.Pos { return d.EndProverty }
func (d *PropertyGetStmt) End() token.Pos { return d.EndProverty }
func (d *PropertySetStmt) End() token.Pos { return d.EndProverty }
func (d *FuncDecl) End() token.Pos        { return d.EndFunc }
func (d *ClassDecl) End() token.Pos       { return d.EndClass }
func (d *DimDecl) End() token.Pos         { return d.List[len(d.List)-1].End() }
func (d *ReDimDecl) End() token.Pos       { return d.List[len(d.List)-1].End() }

func (*SubDecl) stmtNode()         {}
func (*PropertyGetStmt) stmtNode() {}
func (*PropertyLetStmt) stmtNode() {}
func (*PropertySetStmt) stmtNode() {}
func (*FuncDecl) stmtNode()        {}
func (*ClassDecl) stmtNode()       {}
func (*DimDecl) stmtNode()         {}
func (*ReDimDecl) stmtNode()       {}

type Field struct {
	TokPos token.Pos
	Tok    token.Token // Token.BYVAL | Token.BYREF
	Name   *Ident
}

func (f *Field) Pos() token.Pos { return f.TokPos }
func (f *Field) End() token.Pos { return f.Name.End() }
func (*Field) exprNode()        {}
func (f *Field) String() string {
	if f.Tok.IsValid() {
		return fmt.Sprintf("%s %s", f.Tok.String(), f.Name.Name)
	}
	return f.Name.Name
}

// ----------------------------------------------------------------------------
// Statement

type (
	// An OptionStmt node represents an option statement.
	OptionStmt struct {
		Doc      *CommentGroup // associated documentation; or nil
		Option   token.Pos     // position of "Option"
		Explicit token.Pos     // position of "Explicit"
	}

	// An AssignStmt node represents an assignment statement.
	AssignStmt struct {
		Doc    *CommentGroup // associated documentation; or nil
		Lhs    Expr          // left hand side
		Assign token.Pos     // position of "="
		Rhs    Expr          // right hand side
	}

	// A RandomizeStmt node represents a randomize statement.
	RandomizeStmt struct {
		Doc       *CommentGroup // associated documentation; or nil
		Randomize token.Pos     // position of "Randomize"
	}

	// A WithStmt node represents a with statement.
	WithStmt struct {
		Doc     *CommentGroup // associated documentation; or nil
		With    token.Pos     // position of "With"
		Cond    Expr          // condition expression
		Body    *BlockStmt    // with body
		EndWith token.Pos     // position of "End With"
	}

	// A StopStmt node represents a stop statement.
	StopStmt struct {
		Doc  *CommentGroup // associated documentation; or nil
		Stop token.Pos     // position of "Stop"
	}

	// A SelectStmt node represents a select statement.
	SelectStmt struct {
		Doc       *CommentGroup // associated documentation; or nil
		Select    token.Pos     // position of "Select"
		Var       Expr          // select variable
		Cases     []*CaseStmt   // case statements
		Else      *CaseStmt     // else case; or nil
		EndSelect token.Pos     // position of "End Select"
	}

	// A CaseStmt node represents a case statement.
	CaseStmt struct {
		Doc  *CommentGroup // associated documentation; or nil
		Case token.Pos     // position of "Case"
		Cond Expr          // case condition
		Body *BlockStmt    // case body
	}

	// An IfStmt node represents an if statement.
	IfStmt struct {
		Doc   *CommentGroup // associated documentation; or nil
		If    token.Pos     // position of "If"
		Cond  Expr          // condition expression
		Then  token.Pos     // position of "Then"
		Body  *BlockStmt    // if body
		Else  Stmt          // elseif or else
		EndIf token.Pos     // position of "End If"
	}

	// A BlockStmt node represents a block statement.
	BlockStmt struct {
		List []Stmt // statement list
	}

	// A CallStmt node represents a call statement.
	CallStmt struct {
		Call token.Pos // position of "Call"
		Name *Ident    // procedure name
		Recv []Expr    // arguments
	}

	// A ConstStmt node represents a constant declaration.
	ConstStmt struct {
		Doc         *CommentGroup // associated documentation; or nil
		ModifierPos token.Pos     // position of Modifier
		Modifier    token.Token   // Token.PUBLIC | PRIVATE
		Const       token.Pos     // position of "Const"
		Lhs         Expr          // left hand side
		Assign      token.Pos     // position of "="
		Rhs         Expr          // right hand side
	}

	// An ExitStmt node represents an exit statement.
	ExitStmt struct {
		Doc  *CommentGroup // associated documentation; or nil
		Exit token.Pos     // position of "Exit"
		Tok  token.Token   // Token.Do | For | Function | Property | Sub
	}

	// A ForNextStmt node represents a For..Next statement.
	ForNextStmt struct {
		Doc     *CommentGroup // associated documentation; or nil
		For     token.Pos     // position of "For"
		Start   Expr          // start expression
		To      token.Pos     // position of "To"
		End_    Expr          // end expression
		StepPos token.Pos     // position of "Step"
		Step    Expr          // step expression; or nil
		Body    *BlockStmt    // for body
		Next    token.Pos     // position of "Next"
	}

	// A ForEachStmt node represents a For..Each statement.
	ForEachStmt struct {
		Doc   *CommentGroup // associated documentation; or nil
		For   token.Pos     // position of "For"
		Each  token.Pos     // position of "Each"
		Elem  Expr          // element expression
		In    token.Pos     // position of "In"
		Group Expr          // group expression
		Body  *BlockStmt    // for body
		Next  token.Pos     // position of "Next"
		Stmt  Stmt          // next statement; or nil
	}

	// A WhileWendStmt node represents a While..Wend statement.
	WhileWendStmt struct {
		Doc   *CommentGroup // associated documentation; or nil
		While token.Pos     // position of "While"
		Cond  Expr          // condition expression
		Body  *BlockStmt    // while body
		Wend  token.Pos     // position of "Wend"
	}

	// A DoLoopStmt node represents a Do..Loop statement.
	DoLoopStmt struct {
		Doc    *CommentGroup // associated documentation; or nil
		Do     token.Pos     // position of "Do"
		Pre    bool          // whether condition is before loop
		Tok    token.Token   // Token.WHILE | Token.UNTIL
		TokPos token.Pos     // position of Tok
		Cond   Expr          // condition expression
		Body   *BlockStmt    // do body
		Loop   token.Pos     // position of "Loop"
	}

	// An OnErrorStmt node represents an on error statement.
	OnErrorStmt struct {
		Doc            *CommentGroup // associated documentation; or nil
		On             token.Pos     // position of "On"
		Error          token.Pos     // position of "Error"
		*OnErrorResume               // resume next; or nil
		*OnErrorGoto                 // goto 0; or nil
	}

	// An OnErrorResume node represents a resume next statement.
	OnErrorResume struct {
		Resume token.Pos // position of "Resume"
		Next   token.Pos // position of "Next"
	}

	// An OnErrorGoto node represents a goto 0 statement.
	OnErrorGoto struct {
		GoTo token.Pos // position of "Goto"
		Zero token.Pos // position of "0"
	}

	// A PrivateStmt node represents a private statement.
	PrivateStmt struct {
		Doc     *CommentGroup // associated documentation; or nil
		Private token.Pos     // position of "Private"
		List    []Expr        // variable list
	}

	// A PublicStmt node represents a public statement.
	PublicStmt struct {
		Doc    *CommentGroup // associated documentation; or nil
		Public token.Pos     // position of "Public"
		List   []Expr        // variable list
	}

	// An ExprStmt node represents a (stand-alone) expression
	// in a statement list.
	ExprStmt struct {
		Doc *CommentGroup // associated documentation; or nil
		X   Expr          // expression
	}

	// An EraseStmt node represents an erase statement.
	EraseStmt struct {
		Doc   *CommentGroup // associated documentation; or nil
		Erase token.Pos     // position of "Erase"
		X     Expr          // array expression
	}

	// An ExecuteStmt node represents an execute statement.
	ExecuteStmt struct {
		Doc       *CommentGroup // associated documentation; or nil
		Execute   token.Pos     // position of "Execute"
		Statement string        // statement to execute
	}

	// A SetStmt node represents a set statement.
	SetStmt struct {
		Doc    *CommentGroup // associated documentation; or nil
		Set    token.Pos     // position of "Set"
		Lhs    Expr          // left hand side
		Assign token.Pos     // position of "="
		Rhs    Expr          // right hand side
	}
)

func (s *OptionStmt) Pos() token.Pos    { return s.Option }
func (s *RandomizeStmt) Pos() token.Pos { return s.Randomize }
func (s *WithStmt) Pos() token.Pos      { return s.With }
func (s *StopStmt) Pos() token.Pos      { return s.Stop }
func (s *SelectStmt) Pos() token.Pos    { return s.Select }
func (s *IfStmt) Pos() token.Pos        { return s.If }
func (s *BlockStmt) Pos() token.Pos {
	if len(s.List) > 0 {
		return s.List[0].Pos()
	}
	return token.NoPos
}
func (s *CallStmt) Pos() token.Pos      { return s.Call }
func (s *ExitStmt) Pos() token.Pos      { return s.Exit }
func (s *ForNextStmt) Pos() token.Pos   { return s.For }
func (s *ForEachStmt) Pos() token.Pos   { return s.For }
func (s *WhileWendStmt) Pos() token.Pos { return s.While }
func (s *DoLoopStmt) Pos() token.Pos    { return s.Do }
func (s *OnErrorStmt) Pos() token.Pos   { return s.Error }
func (s *ExprStmt) Pos() token.Pos      { return s.X.Pos() }
func (s *ConstStmt) Pos() token.Pos {
	if s.Modifier == token.ILLEGAL {
		return s.Const
	}
	return s.ModifierPos
}
func (s *EraseStmt) Pos() token.Pos   { return s.Erase }
func (s *ExecuteStmt) Pos() token.Pos { return s.Execute }
func (s *PrivateStmt) Pos() token.Pos { return s.Private }
func (s *PublicStmt) Pos() token.Pos  { return s.Public }
func (s *SetStmt) Pos() token.Pos     { return s.Set }
func (s *AssignStmt) Pos() token.Pos  { return s.Lhs.Pos() }

func (s *OptionStmt) End() token.Pos    { return s.Explicit }
func (s *RandomizeStmt) End() token.Pos { return s.Randomize }
func (s *WithStmt) End() token.Pos      { return s.EndWith }
func (s *StopStmt) End() token.Pos      { return s.Stop }
func (s *SelectStmt) End() token.Pos    { return s.EndSelect }
func (s *IfStmt) End() token.Pos        { return s.EndIf }
func (s *BlockStmt) End() token.Pos {
	if len(s.List) > 0 {
		return s.List[len(s.List)-1].End()
	}
	return token.NoPos
}
func (s *CallStmt) End() token.Pos {
	if len(s.Recv) > 0 {
		return s.Recv[len(s.Recv)-1].End()
	}
	return s.Name.End()
}
func (s *ExitStmt) End() token.Pos    { return token.Pos(int(s.Exit) + len(s.Tok.String())) }
func (s *ForNextStmt) End() token.Pos { return s.Next }
func (s *ForEachStmt) End() token.Pos {
	if s.Stmt != nil {
		return s.Stmt.End()
	}
	return s.Next
}
func (s *WhileWendStmt) End() token.Pos { return s.Wend }
func (s *DoLoopStmt) End() token.Pos    { return s.Loop }
func (s *OnErrorStmt) End() token.Pos {
	if s.OnErrorGoto != nil {
		return s.OnErrorGoto.Zero
	}
	if s.OnErrorResume != nil {
		return s.OnErrorResume.Next
	}
	return token.NoPos
}
func (s *ExprStmt) End() token.Pos    { return s.X.End() }
func (s *ConstStmt) End() token.Pos   { return s.Rhs.End() }
func (s *EraseStmt) End() token.Pos   { return s.X.End() }
func (s *ExecuteStmt) End() token.Pos { return token.Pos(int(s.Execute) + len(s.Statement)) }
func (s *PrivateStmt) End() token.Pos { return s.List[len(s.List)-1].End() }
func (s *PublicStmt) End() token.Pos  { return s.List[len(s.List)-1].End() }
func (s *AssignStmt) End() token.Pos  { return s.Rhs.End() }
func (s *SetStmt) End() token.Pos     { return s.Rhs.End() }

func (*OptionStmt) stmtNode()    {}
func (*RandomizeStmt) stmtNode() {}
func (*WithStmt) stmtNode()      {}
func (*StopStmt) stmtNode()      {}
func (*SelectStmt) stmtNode()    {}
func (*IfStmt) stmtNode()        {}
func (*BlockStmt) stmtNode()     {}
func (*CallStmt) stmtNode()      {}
func (*ExitStmt) stmtNode()      {}
func (*ForNextStmt) stmtNode()   {}
func (*ForEachStmt) stmtNode()   {}
func (*WhileWendStmt) stmtNode() {}
func (*DoLoopStmt) stmtNode()    {}
func (*OnErrorStmt) stmtNode()   {}
func (*ExprStmt) stmtNode()      {}
func (*ConstStmt) stmtNode()     {}
func (*EraseStmt) stmtNode()     {}
func (*ExecuteStmt) stmtNode()   {}
func (*PrivateStmt) stmtNode()   {}
func (*PublicStmt) stmtNode()    {}
func (*AssignStmt) stmtNode()    {}
func (*SetStmt) stmtNode()       {}

// ----------------------------------------------------------------------------
// Expression

type (
	// A BasicLit node represents a literal of basic type.
	BasicLit struct {
		Kind     token.Token // Token.Empty | Token.Null | Token.Boolean | Token.Byte | Token.Integer | Token.Currency | Token.Long | Token.Single | Token.Double | Token.Date | Token.String | Token.Object | Token.Error
		Value    string      // literal value
		ValuePos token.Pos   // literal position
	}

	// An Ident node represents an identifier.
	Ident struct {
		NamePos token.Pos // identifier position
		Name    string    // identifier name
	}

	// An IndexExpr node represents an expression followed by an index.
	IndexExpr struct {
		X      Expr      // expression
		Lparen token.Pos // position of "("
		Index  Expr      // index expression
		Rparen token.Pos // position of ")"
	}

	// An IndexListExpr node represents an expression followed by multiple
	// indices.
	IndexListExpr struct {
		X       Expr      // expression
		Lparen  token.Pos // position of "("
		Indices []Expr    // index expressions
		Rparen  token.Pos // position of ")"
	}

	NewExpr struct {
		New token.Pos // position of "New"
		X   Expr
	}

	// A CallExpr node represents an expression followed by an argument list.
	CallExpr struct {
		Func   Expr
		Lparen token.Pos // position of "("
		Recv   []Expr
		Rparen token.Pos // position of ")"
	}

	// A CmdExpr node represents a command expression.
	CmdExpr struct {
		Cmd  Expr
		Recv []Expr
	}

	// A SelectorExpr node represents an expression followed by a selector.
	SelectorExpr struct {
		X   Expr // expression
		Sel Expr // field selector
	}

	// A BinaryExpr node represents a binary expression.
	BinaryExpr struct {
		X     Expr        // left operand
		OpPos token.Pos   // position of Op
		Op    token.Token // operator
		Y     Expr        // right operand
	}
)

func (x *Ident) Pos() token.Pos         { return x.NamePos }
func (x *CallExpr) Pos() token.Pos      { return x.Func.Pos() }
func (x *CmdExpr) Pos() token.Pos       { return x.Cmd.Pos() }
func (x *IndexExpr) Pos() token.Pos     { return x.X.Pos() }
func (x *IndexListExpr) Pos() token.Pos { return x.X.Pos() }
func (x *NewExpr) Pos() token.Pos       { return x.New }
func (x *SelectorExpr) Pos() token.Pos  { return x.X.Pos() }
func (x *BinaryExpr) Pos() token.Pos    { return x.X.Pos() }
func (x *BasicLit) Pos() token.Pos      { return x.ValuePos }

func (x *Ident) End() token.Pos { return token.Pos(len(x.Name) + int(x.NamePos)) }
func (x *CallExpr) End() token.Pos {
	if x.Rparen.IsValid() {
		return x.Rparen
	}
	if len(x.Recv) > 0 {
		return x.Recv[len(x.Recv)-1].End()
	}
	return x.Func.End()
}
func (x *CmdExpr) End() token.Pos {
	if len(x.Recv) > 0 {
		return x.Recv[len(x.Recv)-1].End()
	}
	return x.Cmd.End()
}
func (x *IndexExpr) End() token.Pos     { return x.Rparen }
func (x *IndexListExpr) End() token.Pos { return x.Rparen }
func (x *NewExpr) End() token.Pos       { return x.X.End() }
func (x *SelectorExpr) End() token.Pos  { return x.Sel.End() }
func (x *BinaryExpr) End() token.Pos    { return x.Y.End() }
func (x *BasicLit) End() token.Pos      { return token.Pos(int(x.ValuePos) + len(x.Value)) }

func (*Ident) exprNode()         {}
func (*CallExpr) exprNode()      {}
func (*CmdExpr) exprNode()       {}
func (*IndexExpr) exprNode()     {}
func (*IndexListExpr) exprNode() {}
func (*NewExpr) exprNode()       {}
func (*SelectorExpr) exprNode()  {}
func (*BinaryExpr) exprNode()    {}
func (*BasicLit) exprNode()      {}

// A File node represents a VBScript source file.
type File struct {
	Doc []*CommentGroup

	Stmts []Stmt
}

func (*File) Pos() token.Pos { return token.NoPos }
func (f *File) End() token.Pos {
	if len(f.Stmts) > 0 {
		return f.Stmts[len(f.Stmts)-1].End()
	}
	return token.NoPos
}
