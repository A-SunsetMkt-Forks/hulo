// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
//
// PowerShell reference: https://learn.microsoft.com/en-us/powershell/scripting/lang-spec/chapter-01?view=powershell-7.5
// And https://learn.microsoft.com/en-us/powershell/module/microsoft.powershell.core/about/about?view=powershell-7.5
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
//
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

// TODO: add help directive
// https://learn.microsoft.com/en-us/powershell/scripting/lang-spec/chapter-14?view=powershell-7.5
type HelpDirective interface {
	helpDirectiveNode()
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

	ClassDecl struct {
		Class token.Pos // position of class keyword
		Name  *Ident
		Body  *BlockStmt
	}

	EnumDecl struct {
		Enum token.Pos // position of enum keyword
		Name *Ident
		Body *BlockStmt
	}

	FuncDecl struct {
		Function token.Pos // position of function keyword
		Name     *Ident
		Body     *BlockStmt
	}

	WorkflowDecl struct {
		Workflow token.Pos // position of workflow keyword
		Name     *Ident
		Body     *BlockStmt
	}

	ParallelDecl struct {
		Parallel token.Pos // position of parallel keyword
		Body     *BlockStmt
	}

	SequenceDecl struct {
		Sequence token.Pos // position of sequence keyword
		Body     *BlockStmt
	}

	InlinescriptDecl struct {
		Inlinescript token.Pos // position of inlinescript keyword
		Body         *BlockStmt
	}

	ProcessDecl struct {
		Process token.Pos // position of process keyword
		Body    *BlockStmt
	}

	Attribute struct {
		Lparen token.Pos // position of '['
		Name   *Ident
		Recv   []Expr
		Rparen token.Pos // position of ']'
	}
)

func (d *FuncDecl) Pos() token.Pos         { return d.Function }
func (d *WorkflowDecl) Pos() token.Pos     { return d.Workflow }
func (d *ParallelDecl) Pos() token.Pos     { return d.Parallel }
func (d *SequenceDecl) Pos() token.Pos     { return d.Sequence }
func (d *InlinescriptDecl) Pos() token.Pos { return d.Inlinescript }
func (d *ProcessDecl) Pos() token.Pos      { return d.Process }
func (d *Attribute) Pos() token.Pos        { return d.Lparen }
func (d *ClassDecl) Pos() token.Pos        { return d.Class }
func (d *EnumDecl) Pos() token.Pos         { return d.Enum }

func (d *FuncDecl) End() token.Pos         { return d.Body.End() }
func (d *WorkflowDecl) End() token.Pos     { return d.Body.End() }
func (d *ParallelDecl) End() token.Pos     { return d.Body.End() }
func (d *SequenceDecl) End() token.Pos     { return d.Body.End() }
func (d *InlinescriptDecl) End() token.Pos { return d.Body.End() }
func (d *ProcessDecl) End() token.Pos      { return d.Body.End() }
func (d *Attribute) End() token.Pos        { return d.Rparen }
func (d *ClassDecl) End() token.Pos        { return d.Body.End() }
func (d *EnumDecl) End() token.Pos         { return d.Body.End() }

func (*FuncDecl) stmtNode()     {}
func (*WorkflowDecl) stmtNode() {}
func (*ProcessDecl) stmtNode()  {}
func (*Attribute) stmtNode()    {}
func (*ClassDecl) stmtNode()    {}
func (*EnumDecl) stmtNode()     {}

type (
	TypeLit struct {
		Lbrack token.Pos // position of '['
		Name   Expr
		Rbrack token.Pos // position of ']'
	}

	CastExpr struct {
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

	// [attribute-1] [attribute-2] ...
	// [type]$name
	Parameter struct {
		Attributes []*Attribute
		X          Expr
		Assign     token.Pos // position of '='
		Value      Expr
	}

	VarExpr struct {
		Dollar token.Pos // position of '$'
		X      Expr
	}

	// A BlockExpr node represents a script block expression.
	// e.g., { "Hello there" } or { param($x) $x * 2 }
	BlockExpr struct {
		Lbrace token.Pos // position of '{'
		List   []Expr    // list of expressions in the block
		Rbrace token.Pos // position of '}'
	}

	// $( ... )
	SubExpr struct {
		Dollar token.Pos // position of '$'
		Lparen token.Pos // position of '('
		X      Expr
		Rparen token.Pos // position of ')'
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

	// X::Y
	StaticMemberAccess struct {
		X        Expr
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
	CommaExpr struct {
		Elems []Expr
	}

	// 1>&2 or 1>>2
	RedirectExpr struct {
		X     Expr
		CtrOp token.Token
		OpPos token.Pos
		Y     Expr
	}

	StringLit struct {
		ValPos token.Pos
		Val    string
	}

	// @" ... "@
	MultiStringLit struct {
		LAt    token.Pos // position of '@'
		LQuote token.Pos // position of '"'
		Val    string
		Rquote token.Pos // position of '"'
		RAt    token.Pos // position of '@'
	}

	NumericLit struct {
		ValPos token.Pos
		Val    string
	}

	BoolLit struct {
		ValPos token.Pos
		Val    bool
	}

	SelectExpr struct {
		X   Expr
		Dot token.Pos // position of '.'
		Sel Expr
	}

	GroupExpr struct {
		Sep    token.Token
		Lparen token.Pos // position of '('
		Elems  []Expr
		Rparen token.Pos // position of ')'
	}

	// 1..10
	RangeExpr struct {
		X      Expr
		DblDot token.Pos // position of '..'
		Y      Expr
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
func (x *StringLit) Pos() token.Pos          { return x.ValPos }
func (x *MultiStringLit) Pos() token.Pos     { return x.LAt }
func (x *NumericLit) Pos() token.Pos         { return x.ValPos }
func (x *IncDecExpr) Pos() token.Pos         { return x.X.Pos() }
func (x *HashTable) Pos() token.Pos          { return x.At }
func (x *HashEntry) Pos() token.Pos          { return x.Key.Pos() }
func (x *IndexExpr) Pos() token.Pos          { return x.X.Pos() }
func (x *IndicesExpr) Pos() token.Pos        { return x.X.Pos() }
func (x *VarExpr) Pos() token.Pos            { return x.Dollar }
func (x *BlockExpr) Pos() token.Pos          { return x.Lbrace }
func (x *CallExpr) Pos() token.Pos           { return x.Func.Pos() }
func (x *CmdExpr) Pos() token.Pos            { return x.Cmd.Pos() }
func (x *ArrayExpr) Pos() token.Pos          { return x.At }
func (x *CommaExpr) Pos() token.Pos          { return x.Elems[0].Pos() }
func (x *CastExpr) Pos() token.Pos           { return x.Lbrack }
func (x *TypeLit) Pos() token.Pos            { return x.Lbrack }
func (x *SelectExpr) Pos() token.Pos         { return x.X.Pos() }
func (x *MemberAccess) Pos() token.Pos       { return x.Dollar }
func (x *StaticMemberAccess) Pos() token.Pos { return x.X.Pos() }
func (x *RangeExpr) Pos() token.Pos          { return x.X.Pos() }
func (x *GroupExpr) Pos() token.Pos          { return x.Lparen }
func (x *RedirectExpr) Pos() token.Pos       { return x.X.Pos() }
func (x *BoolLit) Pos() token.Pos            { return x.ValPos }
func (x *SubExpr) Pos() token.Pos            { return x.Dollar }

func (x *Parameter) End() token.Pos      { return x.X.End() }
func (x *BinaryExpr) End() token.Pos     { return x.Y.End() }
func (x *Ident) End() token.Pos          { return token.Pos(int(x.NamePos) + len(x.Name)) }
func (x *Lit) End() token.Pos            { return token.Pos(int(x.ValPos) + len(x.Val)) }
func (x *StringLit) End() token.Pos      { return token.Pos(int(x.ValPos) + len(x.Val)) }
func (x *MultiStringLit) End() token.Pos { return x.RAt }
func (x *NumericLit) End() token.Pos     { return token.Pos(int(x.ValPos) + len(x.Val)) }
func (x *IncDecExpr) End() token.Pos {
	return x.TokPos + 2
}
func (x *HashTable) End() token.Pos   { return x.Rbrace }
func (x *HashEntry) End() token.Pos   { return x.Value.End() }
func (x *IndexExpr) End() token.Pos   { return x.Rbrack }
func (x *IndicesExpr) End() token.Pos { return x.Rbrack }
func (x *VarExpr) End() token.Pos     { return x.X.End() }
func (x *BlockExpr) End() token.Pos   { return x.Rbrace }
func (x *CallExpr) End() token.Pos    { return x.Recv[len(x.Recv)-1].End() }
func (x *CmdExpr) End() token.Pos     { return x.Args[len(x.Args)-1].End() }
func (x *ArrayExpr) End() token.Pos   { return x.Rparen }
func (x *CommaExpr) End() token.Pos   { return x.Elems[len(x.Elems)-1].End() }
func (x *CastExpr) End() token.Pos    { return x.Rbrack }
func (x *TypeLit) End() token.Pos     { return x.Rbrack }
func (x *SelectExpr) End() token.Pos  { return x.Sel.End() }
func (x *MemberAccess) End() token.Pos {
	if x.Long {
		return x.Rbrace
	}
	return x.Y.End()
}
func (x *StaticMemberAccess) End() token.Pos { return x.Y.End() }
func (x *RangeExpr) End() token.Pos          { return x.Y.End() }
func (x *GroupExpr) End() token.Pos          { return x.Rparen }
func (x *RedirectExpr) End() token.Pos       { return x.Y.End() }
func (x *BoolLit) End() token.Pos {
	if x.Val {
		return token.Pos(int(x.ValPos) + 4)
	}
	return token.Pos(int(x.ValPos) + 5)
}
func (x *SubExpr) End() token.Pos { return x.Rparen }

func (*Parameter) exprNode()          {}
func (*CastExpr) exprNode()           {}
func (*BinaryExpr) exprNode()         {}
func (*CallExpr) exprNode()           {}
func (*Ident) exprNode()              {}
func (*Lit) exprNode()                {}
func (*StringLit) exprNode()          {}
func (*MultiStringLit) exprNode()     {}
func (*NumericLit) exprNode()         {}
func (*BoolLit) exprNode()            {}
func (*TypeLit) exprNode()            {}
func (*VarExpr) exprNode()            {}
func (*BlockExpr) exprNode()          {}
func (*SubExpr) exprNode()            {}
func (*HashTable) exprNode()          {}
func (*HashEntry) exprNode()          {}
func (*IndexExpr) exprNode()          {}
func (*IndicesExpr) exprNode()        {}
func (*ExplictExpr) exprNode()        {}
func (*IncDecExpr) exprNode()         {}
func (*CmdExpr) exprNode()            {}
func (*ArrayExpr) exprNode()          {}
func (*CommaExpr) exprNode()          {}
func (*SelectExpr) exprNode()         {}
func (*MemberAccess) exprNode()       {}
func (*StaticMemberAccess) exprNode() {}
func (*RangeExpr) exprNode()          {}
func (*GroupExpr) exprNode()          {}
func (*RedirectExpr) exprNode()       {}

type (
	BlockStmt struct {
		List []Stmt
	}

	ExprStmt struct {
		X Expr
	}

	AssignStmt struct {
		Lhs    Expr
		TokPos token.Pos // position of '=' or '+=', '-=', '*=', etc.
		Tok    token.Token
		Rhs    Expr
	}

	TryStmt struct {
		Try         token.Pos // position of 'try'
		Body        *BlockStmt
		Catches     []*CatchClause
		Finally     token.Pos // position of 'finally'
		FinallyBody *BlockStmt
	}

	CatchClause struct {
		Catch  token.Pos // position of 'catch'
		Lbrack token.Pos // position of '['
		Type   Expr
		Rbrack token.Pos // position of ']'
		Body   *BlockStmt
	}

	TrapStmt struct {
		Trap token.Pos // position of 'trap'
		Body *BlockStmt
	}

	FuncStmt struct {
		Func Expr
		Recv []Expr
		Body *BlockStmt
	}

	ForeachStmt struct {
		Foreach token.Pos   // position of 'foreach'
		Opt     token.Token // TODO: -Parallel
		Lparen  token.Pos   // position of '('
		Elm     Expr
		In      token.Pos // position of 'in'
		Elms    Expr
		Rparen  token.Pos // position of ')'
		Body    *BlockStmt
	}

	ForStmt struct {
		For   token.Pos // position of 'for'
		Init  Expr
		Semi  token.Pos // position of ';'
		Cond  Expr
		Semi2 token.Pos // position of ';'
		Post  Expr
		Body  *BlockStmt
	}

	IfStmt struct {
		If     token.Pos // position of 'if'
		LParen token.Pos // position of '('
		Cond   Expr
		RParen token.Pos // position of ')'
		Body   *BlockStmt
		Else   Stmt
	}

	ReturnStmt struct {
		Return token.Pos // position of 'return'
		X      Expr
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
		Body   *BlockStmt
		While  token.Pos // position of 'while'
		Lparen token.Pos // position of '('
		Cond   Expr
		Rparen token.Pos // position of ')'
	}

	DoUntilStmt struct {
		Do     token.Pos // position of 'do'
		Body   *BlockStmt
		Until  token.Pos // position of 'until'
		Lparen token.Pos // position of '('
		Cond   Expr
		Rparen token.Pos // position of ')'
	}

	WhileStmt struct {
		While  token.Pos // position of 'while'
		Lparen token.Pos // position of '('
		Cond   Expr
		Rparen token.Pos // position of ')'
		Body   *BlockStmt
	}

	// continue [label]
	ContinueStmt struct {
		Continue token.Pos // position of 'continue'
		Label    *Ident
	}

	SwitchStmt struct {
		Switch        token.Pos // position of 'switch'
		Pattern       SwitchPattern
		Casesensitive bool
		Lparen        token.Pos // position of '('
		Value         Expr
		Rparen        token.Pos // position of ')'
		Lbrace        token.Pos // position of '{'
		Cases         []*CaseClause
		Default       *BlockStmt
		Rbrace        token.Pos // position of '}'
	}

	CaseClause struct {
		Cond Expr
		Body *BlockStmt
	}

	DataStmt struct {
		Data token.Pos // position of 'data'
		Recv []Expr
		Body *BlockStmt
	}

	BreakStmt struct {
		Break token.Pos
	}

	DynamicparamStmt struct {
		Dynamicparam token.Pos // position of 'dynamicparam'
		Body         *BlockStmt
	}

	ParamBlock struct {
		Attributes []*Attribute
		Param      token.Pos // position of 'param'
		Lparen     token.Pos // position of '('
		Params     []*Parameter
		Rparen     token.Pos // position of ')'
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

func (s *BlockStmt) Pos() token.Pos        { return s.List[0].Pos() }
func (x *ExprStmt) Pos() token.Pos         { return x.X.Pos() }
func (x *AssignStmt) Pos() token.Pos       { return x.Lhs.Pos() }
func (x *IfStmt) Pos() token.Pos           { return x.If }
func (x *LabelStmt) Pos() token.Pos        { return x.Colon }
func (x *ForStmt) Pos() token.Pos          { return x.For }
func (x *DoWhileStmt) Pos() token.Pos      { return x.Do }
func (x *DoUntilStmt) Pos() token.Pos      { return x.Do }
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
func (x *ParamBlock) Pos() token.Pos {
	if len(x.Attributes) > 0 {
		return x.Attributes[0].Pos()
	}
	return x.Param
}
func (x *ContinueStmt) Pos() token.Pos { return x.Continue }
func (x *BreakStmt) Pos() token.Pos    { return x.Break }
func (x *ReturnStmt) Pos() token.Pos   { return x.Return }

func (s *BlockStmt) End() token.Pos  { return s.List[len(s.List)-1].End() }
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
func (x *DoUntilStmt) End() token.Pos      { return x.Rparen }
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
func (x *ParamBlock) End() token.Pos       { return x.Rparen }
func (x *ContinueStmt) End() token.Pos     { return x.Label.End() }
func (x *BreakStmt) End() token.Pos        { return x.Break + 5 }
func (x *ReturnStmt) End() token.Pos       { return x.Return }

func (*BlockStmt) stmtNode()        {}
func (*ExprStmt) stmtNode()         {}
func (*AssignStmt) stmtNode()       {}
func (*FuncStmt) stmtNode()         {}
func (*IfStmt) stmtNode()           {}
func (*ForStmt) stmtNode()          {}
func (*ForeachStmt) stmtNode()      {}
func (*DoWhileStmt) stmtNode()      {}
func (*DoUntilStmt) stmtNode()      {}
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
func (*ParamBlock) stmtNode()       {}
func (*BreakStmt) stmtNode()        {}
func (*ParallelDecl) stmtNode()     {}
func (*SequenceDecl) stmtNode()     {}
func (*InlinescriptDecl) stmtNode() {}

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
