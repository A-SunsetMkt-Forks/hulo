// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package ast

import "github.com/hulo-lang/hulo/syntax/vbs/token"

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

type CommentGroup struct {
	List []*Comment
}

func (g *CommentGroup) Pos() token.Pos { return g.List[0].Pos() }
func (g *CommentGroup) End() token.Pos { return g.List[len(g.List)-1].End() }

type Comment struct {
	TokPos token.Pos
	Tok    token.Token // ' or Rem
	Text   string
}

func (c *Comment) Pos() token.Pos { return c.TokPos }
func (c *Comment) End() token.Pos { return token.Pos(int(c.TokPos) + len(c.Text)) }

type Modifier int

func (m Modifier) IsNone() bool {
	return m == M_NONE
}

func (m Modifier) HasPublic() bool {
	return m&M_PUBLIC != 0
}

func (m Modifier) HasPrivate() bool {
	return m&M_PRIVATE != 0
}

func (m Modifier) IsAll() bool {
	return m == M_ALL
}

const (
	M_NONE   = 0
	M_PUBLIC = 1 << iota
	M_PRIVATE
	M_ALL = M_PUBLIC | M_PRIVATE
)

// ----------------------------------------------------------------------------
// Declarations

type (
	// A SubDecl node represents a sub declaration.
	SubDecl struct {
		Mod    Modifier
		ModPos token.Pos
		Sub    token.Pos // position of "Sub"
		Name   *Ident
		Recv   []*Field
		Body   *BlockStmt
		EndSub token.Pos // position of "End Sub"
	}

	// A FuncDecl node represents a function declaration.
	FuncDecl struct {
		Mod      Modifier
		ModPos   token.Pos
		Function token.Pos // position of "Function"
		Name     *Ident
		Recv     []*Field
		Body     *BlockStmt
		EndFunc  token.Pos // position of "End Function"
	}

	// A PropertyDecl node represents a property declaration.
	PropertyDecl struct {
		Mod         Modifier
		ModPos      token.Pos
		Property    token.Pos   // position of "Property"
		Tok         token.Token // Token.LET | Token.GET | Token.SET
		TokPos      token.Pos
		Name        *Ident
		Recv        []*Field
		Body        *BlockStmt
		EndProverty token.Pos // position of "End Property"
	}

	// A ClassDecl node represents a class declaration.
	ClassDecl struct {
		Mod    Modifier
		ModPos token.Pos
		Class  token.Pos // position of "Class"
		Name   *Ident
		// dim, func, member, proverpty, assign
		Stmts    []Stmt
		Decls    []Decl
		EndClass token.Pos // position of "End Class"
	}

	// A DimDecl node represents an dim declaration.
	DimDecl struct {
		Dim   token.Pos // position of "Dim"
		List  []Expr
		Colon token.Pos // position of ":"
		Set   *AssignStmt
	}

	// A ReDimDecl node represents a redim declaration.
	ReDimDecl struct {
		ReDim    token.Pos // position of "ReDim"
		Preserve token.Pos // position of "Preserve"
		List     []Expr
	}
)

func (d *SubDecl) Pos() token.Pos {
	if d.Mod.HasPublic() {
		return d.ModPos
	}
	if d.Mod.HasPrivate() {
		return d.ModPos
	}
	return d.Sub
}
func (d *PropertyDecl) Pos() token.Pos {
	if d.Mod.HasPublic() {
		return d.ModPos
	}
	if d.Mod.HasPrivate() {
		return d.ModPos
	}
	return d.Property
}
func (d *FuncDecl) Pos() token.Pos {
	if d.Mod.HasPublic() {
		return d.ModPos
	}
	if d.Mod.HasPrivate() {
		return d.ModPos
	}
	return d.Function
}
func (d *ClassDecl) Pos() token.Pos {
	if d.Mod.HasPublic() {
		return d.ModPos
	}
	if d.Mod.HasPrivate() {
		return d.ModPos
	}
	return d.Class
}
func (s *DimDecl) Pos() token.Pos   { return s.Dim }
func (s *ReDimDecl) Pos() token.Pos { return s.ReDim }

func (d *SubDecl) End() token.Pos      { return d.EndSub }
func (d *PropertyDecl) End() token.Pos { return d.EndProverty }
func (d *FuncDecl) End() token.Pos     { return d.EndFunc }
func (d *ClassDecl) End() token.Pos    { return d.EndClass }
func (d *DimDecl) End() token.Pos      { return d.List[len(d.List)-1].End() }
func (d *ReDimDecl) End() token.Pos    { return d.List[len(d.List)-1].End() }

func (*SubDecl) declNode()      {}
func (*PropertyDecl) declNode() {}
func (*FuncDecl) declNode()     {}
func (*ClassDecl) declNode()    {}
func (*DimDecl) declNode()      {}
func (*ReDimDecl) declNode()    {}

type Field struct {
	TokPos token.Pos
	Tok    token.Token // Token.BYVAL | Token.BYREF
	Name   *Ident
}

// ----------------------------------------------------------------------------
// Statement

type (
	// An OptionStmt node represents an option statement.
	OptionStmt struct {
		Option   token.Pos // position of "Option"
		Explicit token.Pos // position of "Explicit"
	}

	// A RandomizeStmt node represents a randomize statement.
	RandomizeStmt struct {
		Randomize token.Pos // position of "Randomize"
	}

	// A WithStmt node represents a with statement.
	WithStmt struct {
		With    token.Pos // position of "With"
		Cond    Expr
		Body    *BlockStmt
		EndWith token.Pos // position of "End With"
	}

	// An AssignStmt node represents an assign statement.
	AssignStmt struct {
		Tok    token.Token // Token.SET | Token.CONST
		TokPos token.Pos   // position of Tok
		Lhs    Expr
		Assign token.Pos // position of '='
		Rhs    Expr
	}

	// A StopStmt node represents a stop statement.
	StopStmt struct {
		Stop token.Pos // position of "Stop"
	}

	// A SelectStmt node represents a select statement.
	SelectStmt struct {
		Select    token.Pos // position of "Select"
		Var       Expr
		Cases     []*CaseStmt
		Else      *CaseStmt
		EndSelect token.Pos // position of "End Select"
	}

	// A CaseStmt node represents a case statement.
	CaseStmt struct {
		Case token.Pos // position of "Case"
		Cond Expr
		Body *BlockStmt
	}

	// An IfStmt node represents an if statement.
	IfStmt struct {
		If     token.Pos // position of "If"
		Cond   Expr
		Then   token.Pos // position of "Then"
		Body   *BlockStmt
		ElseIf []*IfStmt
		Else   *BlockStmt
		EndIf  token.Pos // position of "End If"
	}

	// A BlockStmt node represents a block statement.
	BlockStmt struct {
		List []Stmt
	}

	// A CallStmt node represents a call statement.
	CallStmt struct {
		Call token.Pos // position of "Call"
		Name *Ident
		Recv []Expr
	}

	// An ExitStmt node represents an exit statement.
	ExitStmt struct {
		Exit token.Pos   // position of "Exit"
		X    token.Token // Token.Do | For | Function | Property | Sub
	}

	// A ForNextStmt node represents a For..Next statement.
	ForNextStmt struct {
		For     token.Pos // position of "For"
		Start   Expr
		To      token.Pos // position of "To"
		End_    Expr
		StepPos token.Pos // position of "Step"
		Step    Expr
		Body    *BlockStmt
		Next    token.Pos // position of "Next"
	}

	// A ForEachStmt node represents a For..Each statement.
	ForEachStmt struct {
		For   token.Pos // position of "For"
		Each  token.Pos // position of "Each"
		Elem  Expr
		In    token.Pos // position of "In"
		Group Expr
		Body  *BlockStmt
		Next  token.Pos // position of "Next"
		Stmt  Stmt
	}

	// A WhileWendStmt node represents a While..Wend statement.
	WhileWendStmt struct {
		While token.Pos // position of "While"
		Cond  Expr
		Body  *BlockStmt
		Wend  token.Pos // position of "Wend"
	}

	// A DoLoopStmt node represents a Doop..Loop statement.
	DoLoopStmt struct {
		Do     token.Pos // position of "Do"
		Pre    bool
		Tok    token.Token // Token.WHILE | Token.UNTIL
		TokPos token.Pos
		Cond   Expr
		Body   *BlockStmt
		Loop   token.Pos // position of "Loop"
	}

	// A OnErrorStmt node represents a on..error statement.
	OnErrorStmt struct {
		On    token.Pos
		Error token.Pos
		*OnErrorResume
		*OnErrorGoto
	}

	OnErrorResume struct {
		Resume token.Pos
		Next   token.Pos
	}

	OnErrorGoto struct {
		GoTo token.Pos
		Zero token.Pos
	}

	MemberStmt struct {
		Mod    Modifier // public or private
		ModPos token.Pos
		Name   *Ident
	}

	// An ExprStmt node represents a (stand-alone) expression
	// in a statement list.
	ExprStmt struct {
		Doc *CommentGroup
		X   Expr // expression
	}
)

func (s *OptionStmt) Pos() token.Pos    { return s.Option }
func (s *RandomizeStmt) Pos() token.Pos { return s.Randomize }
func (s *WithStmt) Pos() token.Pos      { return s.With }
func (s *AssignStmt) Pos() token.Pos {
	if s.TokPos.IsValid() {
		return s.TokPos
	}
	return s.Lhs.Pos()
}
func (s *StopStmt) Pos() token.Pos   { return s.Stop }
func (s *SelectStmt) Pos() token.Pos { return s.Select }
func (s *IfStmt) Pos() token.Pos     { return s.If }
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
func (s *MemberStmt) Pos() token.Pos {
	if !s.Mod.IsNone() {
		return s.ModPos
	}
	return s.Name.End()
}
func (s *ExprStmt) Pos() token.Pos { return s.X.Pos() }

func (s *OptionStmt) End() token.Pos    { return s.Explicit }
func (s *RandomizeStmt) End() token.Pos { return s.Randomize }
func (s *WithStmt) End() token.Pos      { return s.EndWith }
func (s *AssignStmt) End() token.Pos    { return s.Rhs.End() }
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
func (s *ExitStmt) End() token.Pos    { return token.Pos(int(s.Exit) + len(s.X)) }
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
func (s *MemberStmt) End() token.Pos { return s.Name.Pos() }
func (s *ExprStmt) End() token.Pos   { return s.X.End() }

func (*OptionStmt) stmtNode()    {}
func (*RandomizeStmt) stmtNode() {}
func (*WithStmt) stmtNode()      {}
func (*AssignStmt) stmtNode()    {}
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
func (*MemberStmt) stmtNode()    {}
func (*ExprStmt) stmtNode()      {}

// ----------------------------------------------------------------------------
// Expression

type (
	// A BasicLit node represents a literal of basic type.
	BasicLit struct {
		Kind     token.Token // Token.Empty | Token.Null | Token.Boolean | Token.Byte | Token.Integer | Token.Currency | Token.Long | Token.Single | Token.Double | Token.Date | Token.String | Token.Object | Token.Error
		Value    string
		ValuePos token.Pos // literal position
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

	// A SelectorExpr node represents an expression followed by a selector.
	SelectorExpr struct {
		X   Expr   // expression
		Sel *Ident // field selector
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
func (x *IndexExpr) End() token.Pos     { return x.Rparen }
func (x *IndexListExpr) End() token.Pos { return x.Rparen }
func (x *NewExpr) End() token.Pos       { return x.X.End() }
func (x *SelectorExpr) End() token.Pos  { return x.Sel.End() }
func (x *BinaryExpr) End() token.Pos    { return x.Y.End() }
func (x *BasicLit) End() token.Pos      { return token.Pos(int(x.ValuePos) + len(x.Value)) }

func (*Ident) exprNode()         {}
func (*CallExpr) exprNode()      {}
func (*IndexExpr) exprNode()     {}
func (*IndexListExpr) exprNode() {}
func (*NewExpr) exprNode()       {}
func (*SelectorExpr) exprNode()  {}
func (*BinaryExpr) exprNode()    {}
func (*BasicLit) exprNode()      {}

type File struct {
	Doc *CommentGroup

	Stmts []Stmt
	Decls []Decl
}

func (*File) Pos() token.Pos { return token.NoPos }
func (*File) End() token.Pos { return token.NoPos }
