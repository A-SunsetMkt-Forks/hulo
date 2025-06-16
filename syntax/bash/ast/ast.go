// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package ast

import (
	"fmt"
	"strings"

	"github.com/hulo-lang/hulo/syntax/bash/token"
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

// All declaration nodes implement the Decl interface.
type Decl interface {
	Node
	declNode()
}

// All expression nodes implement the Expr interface.
type Expr interface {
	Node
	exprNode()
	String() string
}

// ----------------------------------------------------------------------------
// Comments

// A Comment node represents a single //-style or /*-style comment.
//

type Comment struct {
	Hash token.Pos
	Text string
}

func (c *Comment) Pos() token.Pos { return c.Hash }
func (c *Comment) End() token.Pos { return token.Pos(int(c.Hash) + len(c.Text)) }

type CommentGroup struct {
	List []*Comment
}

func (g *CommentGroup) Pos() token.Pos { return g.List[0].Pos() }
func (g *CommentGroup) End() token.Pos { return g.List[len(g.List)-1].End() }

type FuncDecl struct {
	Function token.Pos
	Name     *Ident
	Lparen   token.Pos
	Rparen   token.Pos
	Body     *BlockStmt
}

func (d *FuncDecl) Pos() token.Pos { return d.Function }

func (d *FuncDecl) End() token.Pos { return d.Body.Closing }

func (*FuncDecl) declNode() {}

type (
	AssignStmt struct {
		Local  token.Pos // position of "local"
		Lhs    Expr
		Assign token.Pos // position of "="
		Rhs    Expr
	}

	BlockStmt struct {
		Tok     token.Token // Token.NONE | Token.LBRACE
		Opening token.Pos
		List    []Stmt
		Closing token.Pos
	}

	ExprStmt struct {
		X Expr
	}

	ReturnStmt struct {
		Return token.Pos // position of "return"
		X      Expr
	}

	BreakStmt struct {
		Break token.Pos // position of "break"
	}

	ContinueStmt struct {
		Continue token.Pos // position of "continue"
	}

	WhileStmt struct {
		While token.Pos // position of "while"
		Cond  Expr
		Semi  token.Pos // position of ";"
		Do    token.Pos // position of "do"
		Body  *BlockStmt
		Done  token.Pos // position of "done"
	}

	UntilStmt struct {
		Until token.Pos // position of "until"
		Cond  Expr
		Semi  token.Pos // position of ";"
		Do    token.Pos // position of "do"
		Body  *BlockStmt
		Done  token.Pos // position of "done"
	}

	// for name [ [in [words …] ] ; ] do commands; done
	ForInStmt struct {
		For  token.Pos // position of "for"
		Var  Expr
		In   token.Pos // position of "in"
		List Expr
		Semi token.Pos // position of ";"
		Do   token.Pos // position of "do"
		Body *BlockStmt
		Done token.Pos // position of "done"
	}

	// for (( expr1 ; expr2 ; expr3 )) ; do commands ; done
	ForStmt struct {
		For    token.Pos // position of "for"
		Lparen token.Pos // position of "(("
		Init   Node
		Semi1  token.Pos // position of ";"
		Cond   Expr
		Semi2  token.Pos // position of ";"
		Post   Node
		Rparen token.Pos // position of "))"
		Do     token.Pos // position of "do"
		Body   *BlockStmt
		Done   token.Pos // position of "done"
	}

	IfStmt struct {
		If   token.Pos // position of "if"
		Cond Expr
		Semi token.Pos // position of ";"
		Then token.Pos // position of "then"
		Body *BlockStmt
		Else Stmt
		Fi   token.Pos // position of "fi"
	}

	CaseStmt struct {
		Case     token.Pos // position of "case"
		X        Expr
		In       token.Pos // position of "in"
		Patterns []*CaseClause
		Else     *BlockStmt
		Esac     token.Pos // position of "esac"
	}

	CaseClause struct {
		Conds  []Expr
		Rparen token.Pos // position of ")"
		Body   *BlockStmt
		Semi   token.Pos // position of ";;"
	}

	SelectStmt struct {
		Select token.Pos // position of "select"
		Var    Expr
		In     token.Pos // position of "in"
		List   Expr
		Semi   token.Pos // position of ";"
		Do     token.Pos // position of "do"
		Body   *BlockStmt
		Done   token.Pos // position of "done"
	}
)

func (s *AssignStmt) Pos() token.Pos { return s.Lhs.Pos() }
func (s *BlockStmt) Pos() token.Pos {
	if s.Opening.IsValid() {
		return s.Opening
	}
	if len(s.List) > 0 {
		return s.List[0].Pos()
	}
	return token.NoPos
}
func (s *ExprStmt) Pos() token.Pos     { return s.X.Pos() }
func (s *ReturnStmt) Pos() token.Pos   { return s.Return }
func (s *BreakStmt) Pos() token.Pos    { return s.Break }
func (s *ContinueStmt) Pos() token.Pos { return s.Continue }
func (s *WhileStmt) Pos() token.Pos    { return s.While }
func (s *UntilStmt) Pos() token.Pos    { return s.Until }
func (s *ForInStmt) Pos() token.Pos    { return s.For }
func (s *ForStmt) Pos() token.Pos      { return s.For }
func (s *IfStmt) Pos() token.Pos       { return s.If }
func (s *CaseStmt) Pos() token.Pos     { return s.Case }
func (s *SelectStmt) Pos() token.Pos   { return s.Select }

func (s *AssignStmt) End() token.Pos { return s.Rhs.End() }
func (s *BlockStmt) End() token.Pos {
	if s.Closing.IsValid() {
		return s.Closing
	}
	if len(s.List) > 0 {
		return s.List[len(s.List)-1].Pos()
	}
	return token.NoPos
}
func (s *ExprStmt) End() token.Pos     { return s.X.End() }
func (s *ReturnStmt) End() token.Pos   { return s.X.End() }
func (s *BreakStmt) End() token.Pos    { return s.Break }
func (s *ContinueStmt) End() token.Pos { return s.Continue }
func (s *WhileStmt) End() token.Pos    { return s.Done }
func (s *UntilStmt) End() token.Pos    { return s.Done }
func (s *ForInStmt) End() token.Pos    { return s.Done }
func (s *ForStmt) End() token.Pos      { return s.Done }
func (s *IfStmt) End() token.Pos       { return s.Fi }
func (s *CaseStmt) End() token.Pos     { return s.Esac }
func (s *SelectStmt) End() token.Pos   { return s.Done }

func (*AssignStmt) stmtNode()   {}
func (*BlockStmt) stmtNode()    {}
func (*ExprStmt) stmtNode()     {}
func (*ReturnStmt) stmtNode()   {}
func (*BreakStmt) stmtNode()    {}
func (*ContinueStmt) stmtNode() {}
func (*WhileStmt) stmtNode()    {}
func (*UntilStmt) stmtNode()    {}
func (*ForInStmt) stmtNode()    {}
func (*ForStmt) stmtNode()      {}
func (*IfStmt) stmtNode()       {}
func (*CaseStmt) stmtNode()     {}
func (*SelectStmt) stmtNode()   {}

type (
	// Word string

	// A BinaryExpr node represents a binary expression.
	BinaryExpr struct {
		Compress bool
		X        Expr        // left operand
		OpPos    token.Pos   // position of Op
		Op       token.Token // operator
		Y        Expr        // right operand
	}

	// A CallExpr node represents a call expression.
	CallExpr struct {
		Func *Ident
		Recv []Expr
	}

	// A Ident node represents an identifier expression.
	Ident struct {
		NamePos token.Pos
		Name    string
	}

	// A BasicLit node represents a literal of basic type.
	BasicLit struct {
		Kind     token.Token // Token.Empty | Token.Null | Token.Boolean | Token.Byte | Token.Integer | Token.Currency | Token.Long | Token.Single | Token.Double | Token.Date | Token.String | Token.Object | Token.Error
		Value    string
		ValuePos token.Pos // literal position
	}

	// [ ]
	//
	// A TestExpr node represents a basic test command expression.
	TestExpr struct {
		Lbrack token.Pos // position of "["
		X      Expr
		Rbrack token.Pos // position of "]"
	}

	// [[ ]]
	//
	// An ExtendedTestExpr node represents an extended test command expression.
	ExtendedTestExpr struct {
		Lbrack token.Pos // position of "[["
		X      Expr
		Rbrack token.Pos // position of "]]"
	}

	// (( ))
	//
	// An ArithEvalExpr node represents an arithmetic evaluation expression.
	ArithEvalExpr struct {
		Lparen token.Pos // position of "(("
		X      Expr
		Rparen token.Pos // position of "))"
	}

	// Command Grouping: { }
	//
	// A CmdGroup node represents a command group expression.
	CmdGroup struct {
		Lbrace token.Pos // position of "{"
		List   []Stmt
		Rbrace token.Pos // position of "}"
	}

	// Command Substitution: $( ) or ` `
	//
	// A CmdSubst node represents a command substitution expression.
	CmdSubst struct {
		Dollar  token.Pos
		Tok     token.Token // Token.LPAREN or Token.BACK_QUOTE
		Opening token.Pos
		X       Expr
		Closing token.Pos
	}

	// Process Substitution: <( ) or >( )
	//
	// A ProcSubst node represents a process substitution expression.
	ProcSubst struct {
		Tok    token.Token // Token.GT or Token.LT
		TokPos token.Pos
		Lparen token.Pos // position of "("
		X      Expr
		Rparen token.Pos // position of ")"
	}

	// Arithmetic Expansion: $(())
	//
	// An ArithExpr node represents an arithmetic expansion expression.
	ArithExpr struct {
		Dollar token.Pos // position of "$"
		Lparen token.Pos // position of "(("
		X      Expr
		Rparen token.Pos // position of "))"
	}

	// Variable Expansion: $a
	//
	// A VarExpExpr node represents an variable expansion expression.
	VarExpExpr struct {
		Dollar token.Pos // position of "$"
		X      Expr
	}

	// X[Y]
	//
	// A IndexExpr node represents index expression.
	IndexExpr struct {
		X      Expr
		Lbrack token.Pos // position of "["
		Y      Expr
		Rbrack token.Pos // position of "]"
	}

	// Parameter Expandsion: ${}
	//
	// A ParamExpExpr node represents an parameter expansion expression.
	ParamExpExpr struct {
		Dollar               token.Pos // position of "$"
		Lbrace               token.Pos // position of "{"
		Var                  Expr
		*DefaultValExp                 // ${var:-val}
		*DefaultValAssignExp           // ${var:=val}
		*NonNullCheckExp               // ${var:?val}
		*NonNullExp                    // ${var:+val}
		*PrefixExp                     // ${!var*}
		*PrefixArrayExp                // ${!var@}
		*ArrayIndexExp                 // ${!var[*]} or ${!var[@]}
		*LengthExp                     // ${#var}
		*DelPrefix                     // ${var#val} or ${var##val}
		*DelSuffix                     // ${var%val} or ${var%%val}
		*SubstringExp                  // ${var:offset} or ${var:offset:length}
		*ReplaceExp                    // ${var/old/new}
		*ReplacePrefixExp              // ${var/#old/new}
		*ReplaceSuffixExp              // ${var/%old/new}
		*CaseConversionExp             // ${var,} or ${var^} or ${var,,} or ${var^^}
		*OperatorExp                   // ${var@op}
		Rbrace               token.Pos // position of "}"
	}

	// If parameter is unset or null, the expansion of word is substituted.
	// Otherwise, the value of parameter is substituted.
	DefaultValExp struct {
		Colon token.Pos
		Sub   token.Pos
		Val   Expr
	}

	// If parameter is unset or null, the expansion of word is assigned to parameter.
	// The value of parameter is then substituted.
	// Positional parameters and special parameters may not be assigned to in this way.
	DefaultValAssignExp struct {
		Colon  token.Pos
		Assign token.Pos
		Val    Expr
	}

	// If parameter is null or unset, the expansion of word (or a message to that effect if word is not present)
	// is written to the standard error and the shell, if it is not interactive, exits.
	// Otherwise, the value of parameter is substituted.
	NonNullCheckExp struct {
		Colon token.Pos
		Quest token.Pos
		Val   Expr
	}

	// If parameter is null or unset, nothing is substituted,
	// otherwise the expansion of word is substituted.
	NonNullExp struct {
		Colon token.Pos
		Add   token.Pos
		Val   Expr
	}

	PrefixExp struct {
		Bitnot token.Pos // position of '!'
		Mul    token.Pos // position of '*'
	}

	PrefixArrayExp struct {
		Bitnot token.Pos // position of '!'
		At     token.Pos // position of '@'
	}

	ArrayIndexExp struct {
		Bitnot   token.Pos   // position of '!'
		Lbracket token.Pos   // position of '['
		Tok      token.Token // Token.AT or Token.MUL
		TokPos   token.Pos
		Rbracket token.Pos // position of ']'
	}

	DelPrefix struct {
		Longest bool
		Hash    token.Pos // position of "#" or "##"
		Val     Expr
	}

	DelSuffix struct {
		Longest bool
		Mod     token.Pos // position of "%" or "%%"
		Val     Expr
	}

	SubstringExp struct {
		Colon1 token.Pos // position of ":"
		Offset int
		Colon2 token.Pos // position of ":"
		Length int
	}

	ReplaceExp struct {
		All  bool
		Div1 token.Pos // position of "/" or "//"
		Old  string
		Div2 token.Pos // position of "/"
		New  string
	}

	ReplacePrefixExp struct {
		Div1 token.Pos // position of "/"
		Hash token.Pos // position of "#"
		Old  string
		Div2 token.Pos // position of "/"
		New  string
	}

	ReplaceSuffixExp struct {
		Div1 token.Pos // position of "/"
		Mod  token.Pos // position of "%"
		Old  string
		Div2 token.Pos // position of "/"
		New  string
	}

	LengthExp struct {
		Hash token.Pos // position of "#"
	}

	CaseConversionExp struct {
		FirstChar bool
		ToUpper   bool
	}

	OperatorExp struct {
		At token.Pos // position of "@"
		Op ExpOperator
	}

	// (VAR1 VAR2 ... VARN)
	ArrExpr struct {
		Lparen token.Pos // position of "("
		Vars   []Expr
		Rparen token.Pos // position of ")"
	}
)

// func (x Word) Pos() token.Pos              { return token.NoPos }
func (x *BinaryExpr) Pos() token.Pos       { return x.X.Pos() }
func (x *CallExpr) Pos() token.Pos         { return x.Func.NamePos }
func (x *Ident) Pos() token.Pos            { return x.NamePos }
func (x *BasicLit) Pos() token.Pos         { return x.ValuePos }
func (x *TestExpr) Pos() token.Pos         { return x.Lbrack }
func (x *ExtendedTestExpr) Pos() token.Pos { return x.Lbrack }
func (x *ArithEvalExpr) Pos() token.Pos    { return x.Lparen }
func (x *CmdSubst) Pos() token.Pos         { return x.Dollar }
func (x *ProcSubst) Pos() token.Pos        { return x.TokPos }
func (x *ArithExpr) Pos() token.Pos        { return x.Dollar }
func (x *VarExpExpr) Pos() token.Pos       { return x.Dollar }
func (x *ParamExpExpr) Pos() token.Pos     { return x.Dollar }
func (x *IndexExpr) Pos() token.Pos        { return x.X.Pos() }
func (x *ArrExpr) Pos() token.Pos          { return x.Lparen }

// func (x Word) End() token.Pos        { return token.NoPos }
func (x *BinaryExpr) End() token.Pos { return x.Y.End() }
func (x *CallExpr) End() token.Pos {
	if len(x.Recv) > 0 {
		return x.Recv[len(x.Recv)-1].End()
	}
	return token.NoPos
}
func (x *Ident) End() token.Pos            { return token.Pos(int(x.NamePos) + len(x.Name)) }
func (x *BasicLit) End() token.Pos         { return token.Pos(int(x.ValuePos) + len(x.Value)) }
func (x *TestExpr) End() token.Pos         { return x.Rbrack }
func (x *ExtendedTestExpr) End() token.Pos { return x.Rbrack }
func (x *ArithEvalExpr) End() token.Pos    { return x.Rparen }
func (x *CmdSubst) End() token.Pos         { return x.Closing }
func (x *ProcSubst) End() token.Pos        { return x.Rparen }
func (x *ArithExpr) End() token.Pos        { return x.Rparen }
func (x *VarExpExpr) End() token.Pos       { return x.X.End() }
func (x *ParamExpExpr) End() token.Pos     { return x.Rbrace }
func (x *IndexExpr) End() token.Pos        { return x.Rbrack }
func (x *ArrExpr) End() token.Pos          { return x.Rparen }

// func (Word) exprNode()              {}
func (*BinaryExpr) exprNode()       {}
func (*CallExpr) exprNode()         {}
func (*Ident) exprNode()            {}
func (*BasicLit) exprNode()         {}
func (*TestExpr) exprNode()         {}
func (*ExtendedTestExpr) exprNode() {}
func (*ArithEvalExpr) exprNode()    {}
func (*CmdSubst) exprNode()         {}
func (*ProcSubst) exprNode()        {}
func (*ArithExpr) exprNode()        {}
func (*VarExpExpr) exprNode()       {}
func (*ParamExpExpr) exprNode()     {}
func (*IndexExpr) exprNode()        {}
func (*ArrExpr) exprNode()          {}

func (e *BinaryExpr) String() string {
	switch e.Op {
	case token.TsLss, token.TsGtr:
		return fmt.Sprintf("%s %s %s", e.X, e.Op, e.Y)
	default:
		return fmt.Sprintf("%s%s%s", e.X, e.Op, e.Y)
	}
}

func (e *Ident) String() string {
	return e.Name
}

func (e *BasicLit) String() string {
	if e.Kind == token.STRING {
		return fmt.Sprintf(`"%s"`, e.Value)
	}
	return e.Value
}

func (e *CallExpr) String() string {
	if len(e.Recv) == 0 {
		return e.Func.String()
	}
	recv := []string{}
	for _, e := range e.Recv {
		recv = append(recv, e.String())
	}
	return fmt.Sprintf("%s %s", e.Func, strings.Join(recv, " "))
}

func (e *TestExpr) String() string {
	return fmt.Sprintf("[ %s ]", e.X)
}

func (e *ExtendedTestExpr) String() string {
	return fmt.Sprintf("[[ %s ]]", e.X)
}

func (e *ArithEvalExpr) String() string {
	return fmt.Sprintf("(( %s ))", e.X)
}

func (e *CmdSubst) String() string {
	if e.Tok == token.LeftParen {
		return fmt.Sprintf("$( %s )", e.X)
	}
	return fmt.Sprintf("` %s `", e.X)
}

func (e *ProcSubst) String() string {
	if e.Tok == token.RdrIn {
		return fmt.Sprintf("<( %s )", e.X)
	}
	return fmt.Sprintf(">( %s )", e.X)
}

func (e *ArithExpr) String() string {
	return fmt.Sprintf("$(( %s ))", e.X)
}

func (e *VarExpExpr) String() string {
	return fmt.Sprintf("$%s", e.X)
}

func (e *ParamExpExpr) String() string {
	switch {
	case e.DefaultValExp != nil:
		return fmt.Sprintf("${%s:-%s}", e.Var, e.DefaultValExp.Val)
	case e.DefaultValAssignExp != nil:
		return fmt.Sprintf("${%s:=%s}", e.Var, e.DefaultValAssignExp.Val)
	case e.NonNullCheckExp != nil:
		return fmt.Sprintf("${%s:?%s}", e.Var, e.NonNullCheckExp.Val)
	case e.NonNullExp != nil:
		return fmt.Sprintf("${%s:+%s}", e.Var, e.NonNullExp.Val)
	case e.PrefixExp != nil:
		return fmt.Sprintf("${!%s*}", e.Var)
	case e.PrefixArrayExp != nil:
		return fmt.Sprintf("${!%s@}", e.Var)
	case e.ArrayIndexExp != nil:
		if e.Tok == token.Star {
			return fmt.Sprintf("${!%s[*]}", e.Var)
		}
		return fmt.Sprintf("${!%s[@]}", e.Var)
	case e.LengthExp != nil:
		return fmt.Sprintf("${#%s}", e.Var)
	case e.DelPrefix != nil:
		if e.DelPrefix.Longest {
			return fmt.Sprintf("${%s##%s}", e.Var, e.DelPrefix.Val)
		}
		return fmt.Sprintf("${%s#%s}", e.Var, e.DelPrefix.Val)
	case e.DelSuffix != nil:
		if e.DelSuffix.Longest {
			return fmt.Sprintf("${%s%%%%%s}", e.Var, e.DelPrefix.Val)
		}
		return fmt.Sprintf("${%s%%%s}", e.Var, e.DelPrefix.Val)
	case e.SubstringExp != nil:
		if e.SubstringExp.Offset != e.SubstringExp.Length {
			return fmt.Sprintf("${%s:%d:%d}", e.Var, e.SubstringExp.Offset, e.SubstringExp.Length)
		}
		return fmt.Sprintf("${%s:%d}", e.Var, e.SubstringExp.Offset)
	case e.ReplaceExp != nil:
		return fmt.Sprintf("${%s/%s/%s}", e.Var, e.ReplaceExp.Old, e.ReplaceExp.New)
	case e.ReplacePrefixExp != nil:
		return fmt.Sprintf("${%s/#%s/%s}", e.Var, e.ReplacePrefixExp.Old, e.ReplacePrefixExp.New)
	case e.ReplaceSuffixExp != nil:
		return fmt.Sprintf("${%s/%%%s/%s}", e.Var, e.ReplaceSuffixExp.Old, e.ReplaceSuffixExp.New)
	case e.CaseConversionExp != nil:
		if e.CaseConversionExp.FirstChar && e.CaseConversionExp.ToUpper {
			return fmt.Sprintf("${%s^}", e.Var)
		} else if !e.CaseConversionExp.FirstChar && e.CaseConversionExp.ToUpper {
			return fmt.Sprintf("${%s^^}", e.Var)
		} else if e.CaseConversionExp.FirstChar && !e.CaseConversionExp.ToUpper {
			return fmt.Sprintf("${%s,}", e.Var)
		} else {
			return fmt.Sprintf("${%s,,}", e.Var)
		}
	case e.OperatorExp != nil:
		return fmt.Sprintf("${%s@%s}", e.Var, e.OperatorExp.Op)
	}
	return fmt.Sprintf("${%s}", e.Var)
}

func (e *IndexExpr) String() string {
	return fmt.Sprintf("%s[%s]", e.X, e.Y)
}

func (e *ArrExpr) String() string {
	vars := []string{}
	for _, v := range e.Vars {
		vars = append(vars, v.String())
	}
	return fmt.Sprintf("(%s)", strings.Join(vars, " "))
}

type ExpOperator string

const (
	// The expansion is a string that is the value of parameter
	// with lowercase alphabetic characters converted to uppercase.
	ExpOperatorU = "U"

	// The expansion is a string that is the value of parameter
	// with the first character converted to uppercase,
	// if it is alphabetic.
	ExpOperatoru = "u"

	// The expansion is a string that is the value of parameter
	// with uppercase alphabetic characters converted to lowercase.
	ExpOperatorL = "L"

	// The expansion is a string that is the value of parameter quoted
	// in a format that can be reused as input.
	ExpOperatorQ = "Q"

	// The expansion is a string that is the value of parameter
	// with backslash escape sequences expanded as with the $'…' quoting mechanism.
	ExpOperatorE = "E"

	// The expansion is a string that is the result of expanding the value of parameter
	// as if it were a prompt string (see Controlling the Prompt).
	ExpOperatorP = "P"

	// The expansion is a string in the form of an assignment statement or declare command that,
	// if evaluated, will recreate parameter with its attributes and value.
	ExpOperatorA = "A"

	// Produces a possibly-quoted version of the value of parameter,
	// except that it prints the values of indexed and
	// associative arrays as a sequence of quoted key-value pairs (see Arrays).
	ExpOperatorK = "K"

	// The expansion is a string consisting of flag values representing parameter’s attributes.
	ExpOperatora = "a"

	// Like the ‘K’ transformation, but expands the keys and values of indexed and
	// associative arrays to separate words after word splitting.
	ExpOperatork = "k"
)

type File struct {
	Doc *CommentGroup

	Stmts []Stmt
	Decls []Decl
}

func (*File) Pos() token.Pos { return token.NoPos }
func (*File) End() token.Pos { return token.NoPos }
