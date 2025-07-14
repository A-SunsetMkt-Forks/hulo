// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
//
// PowerShell reference: https://learn.microsoft.com/en-us/powershell/scripting/lang-spec/chapter-01?view=powershell-7.5
package ast

import "github.com/hulo-lang/hulo/syntax/powershell/token"

type Node interface {
	Pos() token.Pos
	End() token.Pos
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
	List []Comment
}

func (g *CommentGroup) Pos() token.Pos {
	if len(g.List) > 0 {
		return g.List[0].Pos()
	}
	return token.NoPos
}

func (g *CommentGroup) End() token.Pos {
	if len(g.List) > 0 {
		return g.List[len(g.List)-1].End()
	}
	return token.NoPos
}

func (*CommentGroup) stmtNode() {}

type Comment interface {
	Stmt
	commentNode()
}

// # <comment>
type SingleLineComment struct {
	Hash token.Pos
	Text string
}

func (c *SingleLineComment) Pos() token.Pos { return c.Hash }
func (c *SingleLineComment) End() token.Pos { return c.Hash + token.Pos(len(c.Text)) }

func (*SingleLineComment) stmtNode()    {}
func (*SingleLineComment) commentNode() {}

// <#
// <help-directive-1>
// <help-content-1>
// ...

// <help-directive-n>
// <help-content-n>
// #>
type DelimitedComment struct {
	Opening token.Pos // position of '<#'
	Text    string
	Closing token.Pos // position of '#>'
}

func (c *DelimitedComment) Pos() token.Pos { return c.Opening }
func (c *DelimitedComment) End() token.Pos { return c.Closing }

func (*DelimitedComment) stmtNode()    {}
func (*DelimitedComment) commentNode() {}

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

	FuncDecl struct {
		Function   token.Pos // position of function keyword
		Name       *Ident
		Attributes []*Attribute
		Params     []Expr
		Body       *BlcokStmt
	}

	WorkflowDecl struct {
		Workflow token.Pos // position of workflow keyword
		Name     *Ident
		Body     *BlcokStmt
	}

	ProcessDecl struct {
		Process token.Pos // position of process keyword
		Body    *BlcokStmt
	}

	Attribute struct {
		Lparen token.Pos // position of '['
		Name   *Ident
		Recv   []Expr
		Rparen token.Pos // position of ']'
	}
)

func (d *FuncDecl) Pos() token.Pos     { return d.Function }
func (d *WorkflowDecl) Pos() token.Pos { return d.Workflow }
func (d *ProcessDecl) Pos() token.Pos  { return d.Process }
func (d *Attribute) Pos() token.Pos    { return d.Lparen }

func (d *FuncDecl) End() token.Pos     { return d.Body.End() }
func (d *WorkflowDecl) End() token.Pos { return d.Body.End() }
func (d *ProcessDecl) End() token.Pos  { return d.Body.End() }
func (d *Attribute) End() token.Pos    { return d.Rparen }

func (*FuncDecl) stmtNode()     {}
func (*WorkflowDecl) stmtNode() {}
func (*ProcessDecl) stmtNode()  {}
func (*Attribute) stmtNode()    {}

type (
	ConstrainedVarExpr struct {
		Lbrack token.Pos // position of '['
		Type   Expr
		Rbrack token.Pos // position of ']'
		X      Expr
	}

	// A BinaryExpr node represents a binary expression.
	BinaryExpr struct {
		X     Expr        // left operand
		OpPos token.Pos   // position of Op
		Op    token.Token // operator
		Y     Expr        // right operand
	}

	CallExpr struct {
		Func Expr
		Recv []Expr
	}

	CmdExpr struct {
		Cmd  Expr
		Args []Expr
	}

	// An Ident node represents an identifier.
	Ident struct {
		NamePos token.Pos // identifier position
		Name    string    // identifier name
	}

	Lit struct {
		Val    string
		ValPos token.Pos
	}

	List []Expr

	BasicLit struct {
		Kind  token.Token
		Value string
	}

	// [attribute-1] [attribute-2] ...
	// [type]$name
	Parameter struct {
		Attributes []*Attribute
		X          Expr
	}

	Fields struct {
		Sep string // " " "," ";"
	}

	VarExpr struct {
		Dollar token.Pos // position of '$'
		X      Expr
	}

	BlockExpr struct {
		List []Expr
	}

	SubExpr struct {
		X Expr
	}

	HashTable struct {
		At      token.Pos // position of '@'
		Lbrace  token.Pos // position of '{'
		Entries []*HashEntry
		Rbrace  token.Pos // position of '}'
	}

	HashEntry struct {
		Key    *Ident
		Assign token.Pos // position of '='
		Value  Expr
	}

	IndexExpr struct {
		X      Expr
		Lbrack token.Pos
		Index  Expr
		Rbrack token.Pos
	}

	IndicesExpr struct {
		X       Expr
		Lbrack  token.Pos // position of '['
		Indices []Expr
		Rbrack  token.Pos // position of ']'
	}

	ExplictExpr struct {
		X Expr
		Y Expr
	}

	// [module]::expr
	ModExpr struct {
		Lbrack   token.Pos // position of '['
		X        Expr
		Rbrack   token.Pos // position of ']'
		DblColon token.Pos // position of '::'
		Y        Expr
	}

	// ++$X --$X
	IncDecExpr struct {
		Pre    bool
		X      Expr
		Tok    token.Token // Token.INC or Token.DEC
		TokPos token.Pos   // position of "++" or "--"
	}

	// $X:Y or ${X:Y}
	MemberAccess struct {
		Dollar token.Pos // position of '$'
		Long   bool
		Lbrace token.Pos // position of '{'
		X      Expr
		Colon  token.Pos // position of ':'
		Y      Expr
		Rbrace token.Pos // position of '}'
	}

	// [X]::Y
	StaticMemberAccess struct {
		Lbrack   token.Pos // position of '['
		X        Expr
		Rbrack   token.Pos // position of ']'
		DblColon token.Pos // position of '::'
		Y        Expr
	}

	// @(10, "blue", 12.54e3, 16.30D)
	ArrayExpr struct {
		At     token.Pos // position of '@'
		Lparen token.Pos // position of '('
		Elems  []Expr
		Rparen token.Pos // position of ')'
	}

	// 10, "blue", 12.54e3, 16.30D
	ArrayConstructor struct {
		Elems []Expr
	}

	SelectExpr struct {
		X   Expr
		Dot token.Pos // position of '.'
		Sel Expr
	}
)

func (x *Parameter) Pos() token.Pos {
	if len(x.Attributes) > 0 {
		return x.Attributes[0].Pos()
	}
	return x.X.Pos()
}
func (x *BinaryExpr) Pos() token.Pos         { return x.X.Pos() }
func (x *Ident) Pos() token.Pos              { return x.NamePos }
func (x *Lit) Pos() token.Pos                { return x.ValPos }
func (x *IncDecExpr) Pos() token.Pos         { return x.X.Pos() }
func (x *HashTable) Pos() token.Pos          { return x.At }
func (x *HashEntry) Pos() token.Pos          { return x.Key.Pos() }
func (x *IndexExpr) Pos() token.Pos          { return x.X.Pos() }
func (x *IndicesExpr) Pos() token.Pos        { return x.X.Pos() }
func (x *VarExpr) Pos() token.Pos            { return x.Dollar }
func (x *CallExpr) Pos() token.Pos           { return x.Func.Pos() }
func (x *CmdExpr) Pos() token.Pos            { return x.Cmd.Pos() }
func (x *ArrayExpr) Pos() token.Pos          { return x.At }
func (x *ArrayConstructor) Pos() token.Pos   { return x.Elems[0].Pos() }
func (x *ConstrainedVarExpr) Pos() token.Pos { return x.Lbrack }
func (x *SelectExpr) Pos() token.Pos         { return x.X.Pos() }
func (x *MemberAccess) Pos() token.Pos       { return x.Dollar }
func (x *StaticMemberAccess) Pos() token.Pos { return x.Lbrack }

func (x *Parameter) End() token.Pos  { return x.X.End() }
func (x *BinaryExpr) End() token.Pos { return x.Y.End() }
func (x *Ident) End() token.Pos      { return token.Pos(int(x.NamePos) + len(x.Name)) }
func (x *Lit) End() token.Pos        { return token.Pos(int(x.ValPos) + len(x.Val)) }
func (x *IncDecExpr) End() token.Pos {
	return x.TokPos + 2
}
func (x *HashTable) End() token.Pos          { return x.Rbrace }
func (x *HashEntry) End() token.Pos          { return x.Value.End() }
func (x *IndexExpr) End() token.Pos          { return x.Rbrack }
func (x *IndicesExpr) End() token.Pos        { return x.Rbrack }
func (x *VarExpr) End() token.Pos            { return x.X.End() }
func (x *CallExpr) End() token.Pos           { return x.Recv[len(x.Recv)-1].End() }
func (x *CmdExpr) End() token.Pos            { return x.Args[len(x.Args)-1].End() }
func (x *ArrayExpr) End() token.Pos          { return x.Rparen }
func (x *ArrayConstructor) End() token.Pos   { return x.Elems[len(x.Elems)-1].End() }
func (x *ConstrainedVarExpr) End() token.Pos { return x.Rbrack }
func (x *SelectExpr) End() token.Pos         { return x.Sel.End() }
func (x *MemberAccess) End() token.Pos {
	if x.Long {
		return x.Rbrace
	}
	return x.Y.End()
}
func (x *StaticMemberAccess) End() token.Pos { return x.Y.End() }

func (*Parameter) exprNode()          {}
func (*ConstrainedVarExpr) exprNode() {}
func (*BinaryExpr) exprNode()         {}
func (*CallExpr) exprNode()           {}
func (*Ident) exprNode()              {}
func (*List) exprNode()               {}
func (*Lit) exprNode()                {}
func (*BasicLit) exprNode()           {}
func (*VarExpr) exprNode()            {}
func (*BlockExpr) exprNode()          {}
func (*SubExpr) exprNode()            {}
func (*HashTable) exprNode()          {}
func (*HashEntry) exprNode()          {}
func (*IndexExpr) exprNode()          {}
func (*IndicesExpr) exprNode()        {}
func (*ExplictExpr) exprNode()        {}
func (*ModExpr) exprNode()            {}
func (*IncDecExpr) exprNode()         {}
func (*CmdExpr) exprNode()            {}
func (*ArrayExpr) exprNode()          {}
func (*ArrayConstructor) exprNode()   {}
func (*SelectExpr) exprNode()         {}
func (*MemberAccess) exprNode()       {}
func (*StaticMemberAccess) exprNode() {}

type (
	BlcokStmt struct {
		List []Stmt
	}

	ExprStmt struct {
		X Expr
	}

	AssignStmt struct {
		Lhs    Expr
		Assign token.Pos // position of '='
		Rhs    Expr
	}

	TryStmt struct {
		Try         token.Pos // position of 'try'
		Body        *BlcokStmt
		Catches     []*CatchClause
		Finally     token.Pos // position of 'finally'
		FinallyBody *BlcokStmt
	}

	CatchClause struct {
		Catch  token.Pos // position of 'catch'
		Lbrack token.Pos // position of '['
		Type   Expr
		Rbrack token.Pos // position of ']'
		Body   *BlcokStmt
	}

	TrapStmt struct {
		Trap token.Pos // position of 'trap'
		Body *BlcokStmt
	}

	FuncStmt struct {
		Func Expr
		Recv []Expr
		Body *BlcokStmt
	}

	ForeachStmt struct {
		Foreach token.Pos // position of 'foreach'
		Elm     Expr
		In      token.Pos // position of 'in'
		Values  Expr
		Body    *BlcokStmt
	}

	ForStmt struct {
		For   token.Pos // position of 'for'
		Init  Expr
		Semi  token.Pos // position of ';'
		Cond  Expr
		Semi2 token.Pos // position of ';'
		Post  Expr
		Body  *BlcokStmt
	}

	IfStmt struct {
		If     token.Pos // position of 'if'
		LParen token.Pos // position of '('
		Cond   Expr
		RParen token.Pos // position of ')'
		Body   *BlcokStmt
		Else   Stmt
	}

	ReturnStmt struct {
		X Expr
	}

	ThrowStmt struct {
		Throw token.Pos // position of 'throw'
		X     Expr
	}

	// :label
	LabelStmt struct {
		Colon token.Pos // position of ':'
		X     Expr
	}

	DoWhileStmt struct {
		Do     token.Pos // position of 'do'
		Body   *BlcokStmt
		While  token.Pos // position of 'while'
		Lparen token.Pos // position of '('
		Cond   Expr
		Rparen token.Pos // position of ')'
	}

	WhileStmt struct {
		While  token.Pos // position of 'while'
		Lparen token.Pos // position of '('
		Cond   Expr
		Rparen token.Pos // position of ')'
		Body   *BlcokStmt
	}

	ContinueStmt struct{}

	SwitchStmt struct {
		Switch        token.Pos // position of 'switch'
		Pattern       SwitchPattern
		Casesensitive bool
		Lparen        token.Pos // position of '('
		Value         Expr
		Rparen        token.Pos // position of ')'
		Lbrace        token.Pos // position of '{'
		Cases         []*CaseClause
		Default       *BlcokStmt
		Rbrace        token.Pos // position of '}'
	}

	CaseClause struct {
		Cond Expr
		Body *BlcokStmt
	}

	DataStmt struct {
		Data token.Pos // position of 'data'
		Recv []Expr
		Body *BlcokStmt
	}

	BreakStmt struct {
		Break token.Pos
		Label *Ident
	}

	DynamicparamStmt struct {
		Dynamicparam token.Pos // position of 'dynamicparam'
		Body         *BlcokStmt
	}

	ParamBlock struct {
		Param  token.Pos // position of 'param'
		Lparen token.Pos // position of '('
		Params []Expr
		Rparen token.Pos // position of ')'
	}
)

type SwitchPattern int

func (p SwitchPattern) String() string {
	return []string{
		"",
		"-regex",
		"-wildcard",
		"-exact",
	}[p]
}

const (
	SwitchPatternNone SwitchPattern = iota
	SwitchPatternRegex
	SwitchPatternWildcard
	SwitchPatternExact
)

func (s *BlcokStmt) Pos() token.Pos        { return s.List[0].Pos() }
func (x *ExprStmt) Pos() token.Pos         { return x.X.Pos() }
func (x *AssignStmt) Pos() token.Pos       { return x.Lhs.Pos() }
func (x *IfStmt) Pos() token.Pos           { return x.If }
func (x *LabelStmt) Pos() token.Pos        { return x.Colon }
func (x *ForStmt) Pos() token.Pos          { return x.For }
func (x *DoWhileStmt) Pos() token.Pos      { return x.Do }
func (x *WhileStmt) Pos() token.Pos        { return x.While }
func (x *ForeachStmt) Pos() token.Pos      { return x.Foreach }
func (x *ThrowStmt) Pos() token.Pos        { return x.Throw }
func (x *SwitchStmt) Pos() token.Pos       { return x.Switch }
func (x *CaseClause) Pos() token.Pos       { return x.Cond.Pos() }
func (x *TryStmt) Pos() token.Pos          { return x.Try }
func (x *CatchClause) Pos() token.Pos      { return x.Catch }
func (x *TrapStmt) Pos() token.Pos         { return x.Trap }
func (x *DataStmt) Pos() token.Pos         { return x.Data }
func (x *DynamicparamStmt) Pos() token.Pos { return x.Dynamicparam }

func (s *BlcokStmt) End() token.Pos  { return s.List[len(s.List)-1].End() }
func (x *ExprStmt) End() token.Pos   { return x.X.End() }
func (x *AssignStmt) End() token.Pos { return x.Rhs.End() }
func (x *IfStmt) End() token.Pos {
	if x.Else == nil {
		return x.RParen
	}
	return x.Else.End()
}
func (x *LabelStmt) End() token.Pos        { return x.X.End() }
func (x *ForStmt) End() token.Pos          { return x.Body.End() }
func (x *DoWhileStmt) End() token.Pos      { return x.Rparen }
func (x *WhileStmt) End() token.Pos        { return x.Body.End() }
func (x *ForeachStmt) End() token.Pos      { return x.Body.End() }
func (x *ThrowStmt) End() token.Pos        { return x.X.End() }
func (x *SwitchStmt) End() token.Pos       { return x.Rbrace }
func (x *CaseClause) End() token.Pos       { return x.Body.End() }
func (x *TryStmt) End() token.Pos          { return x.FinallyBody.End() }
func (x *CatchClause) End() token.Pos      { return x.Body.End() }
func (x *TrapStmt) End() token.Pos         { return x.Body.End() }
func (x *DataStmt) End() token.Pos         { return x.Body.End() }
func (x *DynamicparamStmt) End() token.Pos { return x.Body.End() }

func (*BlcokStmt) stmtNode()        {}
func (*ExprStmt) stmtNode()         {}
func (*AssignStmt) stmtNode()       {}
func (*FuncStmt) stmtNode()         {}
func (*IfStmt) stmtNode()           {}
func (*ForStmt) stmtNode()          {}
func (*ForeachStmt) stmtNode()      {}
func (*DoWhileStmt) stmtNode()      {}
func (*WhileStmt) stmtNode()        {}
func (*ReturnStmt) stmtNode()       {}
func (*ThrowStmt) stmtNode()        {}
func (*LabelStmt) stmtNode()        {}
func (*ContinueStmt) stmtNode()     {}
func (*SwitchStmt) stmtNode()       {}
func (*CaseClause) stmtNode()       {}
func (*DataStmt) stmtNode()         {}
func (*TryStmt) stmtNode()          {}
func (*CatchClause) stmtNode()      {}
func (*TrapStmt) stmtNode()         {}
func (*DynamicparamStmt) stmtNode() {}

type File struct {
	Docs  *CommentGroup
	Name  *Ident
	Stmts []Stmt
}

func (x *File) Pos() token.Pos {
	if len(x.Stmts) > 0 {
		return x.Stmts[0].Pos()
	}
	return token.NoPos
}

func (x *File) End() token.Pos {
	if len(x.Stmts) > 0 {
		return x.Stmts[len(x.Stmts)-1].End()
	}
	return token.NoPos
}

func (*File) stmtNode() {}
