// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
//
// Package ast declares the types used to represent syntax trees for bash scripts.
// Bash reference: https://www.gnu.org/software/bash/manual/bash.html
package ast

import (
	"github.com/hulo-lang/hulo/syntax/bash/token"
)

// ----------------------------------------------------------------------------
// Node interfaces

// All node types implement the Node interface.
// A Node represents any node in the bash syntax tree.
type Node interface {
	Pos() token.Pos // position of first character belonging to the node
	End() token.Pos // position of first character immediately after the node
}

// All statement nodes implement the Stmt interface.
// A Stmt represents any statement in the bash syntax tree.
type Stmt interface {
	Node
	stmtNode()
}

// All expression nodes implement the Expr interface.
// An Expr represents any expression in the bash syntax tree.
type Expr interface {
	Node
	exprNode()
}

// ----------------------------------------------------------------------------
// Comments

// A Comment node represents a single #-style comment.
// The Text field contains the comment text without the leading # character.
type Comment struct {
	Hash token.Pos // position of "#" character
	Text string    // comment text (excluding '#' and trailing newline)
}

func (c *Comment) Pos() token.Pos { return c.Hash }
func (c *Comment) End() token.Pos { return token.Pos(int(c.Hash) + len(c.Text)) }
func (c *Comment) stmtNode()      {}

// A CommentGroup represents a sequence of comments
// with no other tokens and no empty lines between.
type CommentGroup struct {
	List []*Comment // len(List) > 0
}

func (g *CommentGroup) Pos() token.Pos { return g.List[0].Pos() }
func (g *CommentGroup) End() token.Pos { return g.List[len(g.List)-1].End() }
func (c *CommentGroup) stmtNode()      {}

// ----------------------------------------------------------------------------
// Function declarations

// A FuncDecl node represents a function declaration.
// e.g., function_name() { ... } or function function_name { ... }
type FuncDecl struct {
	Function token.Pos  // position of "function" keyword, if any
	Name     *Ident     // function name
	Lparen   token.Pos  // position of "("
	Rparen   token.Pos  // position of ")"
	Body     *BlockStmt // function body
}

func (d *FuncDecl) Pos() token.Pos { return d.Function }
func (d *FuncDecl) End() token.Pos { return d.Body.Closing }
func (*FuncDecl) stmtNode()        {}

// ----------------------------------------------------------------------------
// Statements

type (
	// An AssignStmt node represents a variable assignment.
	// e.g., local var=value or var=value
	AssignStmt struct {
		Local  token.Pos // position of "local" keyword, if any
		Lhs    Expr      // left-hand side (variable name)
		Assign token.Pos // position of "="
		Rhs    Expr      // right-hand side (value)
	}

	// A BlockStmt node represents a braced statement list.
	// e.g., { command1; command2; }
	BlockStmt struct {
		Tok     token.Token // Token.NONE | Token.LBRACE
		Opening token.Pos   // position of "{"
		List    []Stmt      // list of statements
		Closing token.Pos   // position of "}"
	}

	// An ExprStmt node represents a standalone expression.
	// e.g., echo "hello" or ls -la
	ExprStmt struct {
		X Expr // expression
	}

	// A ReturnStmt node represents a return statement.
	// e.g., return 0 or return $?
	ReturnStmt struct {
		Return token.Pos // position of "return"
		X      Expr      // return value expression
	}

	// A WhileStmt node represents a while loop.
	// e.g., while condition; do commands; done
	WhileStmt struct {
		While token.Pos  // position of "while"
		Cond  Expr       // condition expression
		Semi  token.Pos  // position of ";"
		Do    token.Pos  // position of "do"
		Body  *BlockStmt // loop body
		Done  token.Pos  // position of "done"
	}

	// An UntilStmt node represents an until loop.
	// e.g., until condition; do commands; done
	UntilStmt struct {
		Until token.Pos  // position of "until"
		Cond  Expr       // condition expression
		Semi  token.Pos  // position of ";"
		Do    token.Pos  // position of "do"
		Body  *BlockStmt // loop body
		Done  token.Pos  // position of "done"
	}

	// A ForInStmt node represents a for-in loop.
	// e.g., for var in list; do commands; done
	ForInStmt struct {
		For  token.Pos  // position of "for"
		Var  Expr       // variable name
		In   token.Pos  // position of "in"
		List Expr       // list expression
		Semi token.Pos  // position of ";"
		Do   token.Pos  // position of "do"
		Body *BlockStmt // loop body
		Done token.Pos  // position of "done"
	}

	// A ForStmt node represents a C-style for loop.
	// e.g., for ((i=0; i<10; i++)); do commands; done
	ForStmt struct {
		For    token.Pos  // position of "for"
		Lparen token.Pos  // position of "(("
		Init   Node       // initialization expression
		Semi1  token.Pos  // position of ";"
		Cond   Expr       // condition expression
		Semi2  token.Pos  // position of ";"
		Post   Node       // post-iteration expression
		Rparen token.Pos  // position of "))"
		Do     token.Pos  // position of "do"
		Body   *BlockStmt // loop body
		Done   token.Pos  // position of "done"
	}

	// An IfStmt node represents an if statement.
	// e.g., if condition; then commands; else commands; fi
	IfStmt struct {
		If   token.Pos  // position of "if"
		Cond Expr       // condition expression
		Semi token.Pos  // position of ";"
		Then token.Pos  // position of "then"
		Body *BlockStmt // if body
		Else Stmt       // else clause (can be IfStmt for elif)
		Fi   token.Pos  // position of "fi"
	}

	// A CaseStmt node represents a case statement.
	// e.g., case $var in pattern1) commands;; pattern2) commands;; esac
	CaseStmt struct {
		Case     token.Pos     // position of "case"
		X        Expr          // expression to match
		In       token.Pos     // position of "in"
		Patterns []*CaseClause // list of case clauses
		Else     *BlockStmt    // else clause (optional)
		Esac     token.Pos     // position of "esac"
	}

	// A CaseClause represents a single case pattern and its commands.
	CaseClause struct {
		Conds  []Expr     // pattern expressions
		Rparen token.Pos  // position of ")"
		Body   *BlockStmt // commands to execute
		Semi   token.Pos  // position of ";;"
	}

	// A SelectStmt node represents a select statement.
	// e.g., select var in list; do commands; done
	SelectStmt struct {
		Select token.Pos  // position of "select"
		Var    Expr       // variable name
		In     token.Pos  // position of "in"
		List   Expr       // list expression
		Semi   token.Pos  // position of ";"
		Do     token.Pos  // position of "do"
		Body   *BlockStmt // loop body
		Done   token.Pos  // position of "done"
	}
)

// Position methods for statements
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
func (s *ExprStmt) Pos() token.Pos   { return s.X.Pos() }
func (s *ReturnStmt) Pos() token.Pos { return s.Return }
func (s *WhileStmt) Pos() token.Pos  { return s.While }
func (s *UntilStmt) Pos() token.Pos  { return s.Until }
func (s *ForInStmt) Pos() token.Pos  { return s.For }
func (s *ForStmt) Pos() token.Pos    { return s.For }
func (s *IfStmt) Pos() token.Pos     { return s.If }
func (s *CaseStmt) Pos() token.Pos   { return s.Case }
func (s *SelectStmt) Pos() token.Pos { return s.Select }

// End position methods for statements
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
func (s *ExprStmt) End() token.Pos   { return s.X.End() }
func (s *ReturnStmt) End() token.Pos { return s.X.End() }
func (s *WhileStmt) End() token.Pos  { return s.Done }
func (s *UntilStmt) End() token.Pos  { return s.Done }
func (s *ForInStmt) End() token.Pos  { return s.Done }
func (s *ForStmt) End() token.Pos    { return s.Done }
func (s *IfStmt) End() token.Pos     { return s.Fi }
func (s *CaseStmt) End() token.Pos   { return s.Esac }
func (s *SelectStmt) End() token.Pos { return s.Done }

// Statement node type assertions
func (*AssignStmt) stmtNode() {}
func (*BlockStmt) stmtNode()  {}
func (*ExprStmt) stmtNode()   {}
func (*ReturnStmt) stmtNode() {}
func (*WhileStmt) stmtNode()  {}
func (*UntilStmt) stmtNode()  {}
func (*ForInStmt) stmtNode()  {}
func (*ForStmt) stmtNode()    {}
func (*IfStmt) stmtNode()     {}
func (*CaseStmt) stmtNode()   {}
func (*SelectStmt) stmtNode() {}

// ----------------------------------------------------------------------------
// Expressions

type (
	// A BinaryExpr node represents a binary expression.
	// e.g., $a -eq $b or "string" = "other"
	BinaryExpr struct {
		Compress bool        // true if spaces are compressed around operator
		X        Expr        // left operand
		OpPos    token.Pos   // position of Op
		Op       token.Token // operator (EQ, NE, LT, GT, etc.)
		Y        Expr        // right operand
	}

	// A CmdExpr node represents a command expression.
	// e.g., echo "hello" or ls -la
	CmdExpr struct {
		Name *Ident // command name
		Recv []Expr // command arguments
	}

	// A CmdListExpr node represents a command list with control operators.
	// e.g., cmd1 && cmd2 or cmd1 || cmd2
	CmdListExpr struct {
		CtrOp token.Token // Token.OrOr or Token.AndAnd
		Cmds  []Expr      // list of commands
	}

	// An Ident node represents an identifier.
	// e.g., variable names, function names
	Ident struct {
		NamePos token.Pos // identifier position
		Name    string    // identifier name
	}

	// A Word node represents a word (string literal or variable).
	// e.g., "hello world" or $variable
	Word struct {
		Val    string    // word value
		ValPos token.Pos // word position
	}

	// A TestExpr node represents a test command expression.
	// e.g., [ -f file ] or [ $a -eq $b ]
	TestExpr struct {
		Lbrack token.Pos // position of "["
		X      Expr      // test expression
		Rbrack token.Pos // position of "]"
	}

	// An ExtendedTestExpr node represents an extended test command expression.
	// e.g., [[ -f file ]] or [[ $a == $b ]]
	ExtendedTestExpr struct {
		Lbrack token.Pos // position of "[["
		X      Expr      // test expression
		Rbrack token.Pos // position of "]]"
	}

	// An ArithEvalExpr node represents an arithmetic evaluation expression.
	// e.g., ((a + b)) or ((i++))
	ArithEvalExpr struct {
		Lparen token.Pos // position of "(("
		X      Expr      // arithmetic expression
		Rparen token.Pos // position of "))"
	}

	// A CmdGroup node represents a command group expression.
	// e.g., { cmd1; cmd2; } or (cmd1; cmd2)
	CmdGroup struct {
		Op      token.Token // Token.LBRACE or Token.LPAREN
		Opening token.Pos   // position of opening brace/paren
		List    []Expr      // list of commands
		Closing token.Pos   // position of closing brace/paren
	}

	// A CmdSubst node represents a command substitution expression.
	// e.g., $(command) or `command` or ${command}
	CmdSubst struct {
		Dollar  token.Pos   // position of "$"
		Tok     token.Token // Token.DollParen or Token.BckQuote or Token.DollBrace
		Opening token.Pos   // position of opening delimiter
		X       Expr        // command expression
		Closing token.Pos   // position of closing delimiter
	}

	// A ProcSubst node represents a process substitution expression.
	// e.g., <(command) or >(command)
	ProcSubst struct {
		CtrOp   token.Token // Token.CmdIn or Token.CmdOut
		Opening token.Pos   // position of opening delimiter
		X       Expr        // command expression
		Closing token.Pos   // position of closing delimiter
	}

	// An ArithExpr node represents an arithmetic expansion expression.
	// e.g., $((a + b))
	ArithExpr struct {
		Dollar token.Pos // position of "$"
		Lparen token.Pos // position of "(("
		X      Expr      // arithmetic expression
		Rparen token.Pos // position of "))"
	}

	// A VarExpExpr node represents a variable expansion expression.
	// e.g., $variable
	VarExpExpr struct {
		Dollar token.Pos // position of "$"
		X      Expr      // variable name
	}

	// An IndexExpr node represents an array index expression.
	// e.g., array[index]
	IndexExpr struct {
		X      Expr      // array expression
		Lbrack token.Pos // position of "["
		Y      Expr      // index expression
		Rbrack token.Pos // position of "]"
	}

	// A ParamExpExpr node represents a parameter expansion expression.
	// e.g., ${var:-default} or ${var#prefix}
	ParamExpExpr struct {
		Dollar   token.Pos // position of "$"
		Lbrace   token.Pos // position of "{"
		Var      Expr      // variable name
		ParamExp ParamExp  // parameter expansion operation
		Rbrace   token.Pos // position of "}"
	}

	// ParamExp represents a parameter expansion operation.
	ParamExp interface {
		paramExpNode()
	}

	// DefaultValExp represents a default value expansion.
	// e.g., ${var:-default}
	DefaultValExp struct {
		Colon token.Pos // position of ":"
		Sub   token.Pos // position of "-"
		Val   Expr      // default value
	}

	// DefaultValAssignExp represents a default value assignment expansion.
	// e.g., ${var:=default}
	DefaultValAssignExp struct {
		Colon  token.Pos // position of ":"
		Assign token.Pos // position of "="
		Val    Expr      // default value
	}

	// NonNullCheckExp represents a non-null check expansion.
	// e.g., ${var:?error_message}
	NonNullCheckExp struct {
		Colon token.Pos // position of ":"
		Quest token.Pos // position of "?"
		Val   Expr      // error message
	}

	// NonNullExp represents a non-null expansion.
	// e.g., ${var:+value}
	NonNullExp struct {
		Colon token.Pos // position of ":"
		Add   token.Pos // position of "+"
		Val   Expr      // value to substitute
	}

	// PrefixExp represents a prefix expansion.
	// e.g., ${!var*}
	PrefixExp struct {
		Bitnot token.Pos // position of '!'
		Mul    token.Pos // position of '*'
	}

	// PrefixArrayExp represents a prefix array expansion.
	// e.g., ${!var@}
	PrefixArrayExp struct {
		Bitnot token.Pos // position of '!'
		At     token.Pos // position of '@'
	}

	// ArrayIndexExp represents an array index expansion.
	// e.g., ${!var[*]} or ${!var[@]}
	ArrayIndexExp struct {
		Bitnot   token.Pos   // position of '!'
		Lbracket token.Pos   // position of '['
		Tok      token.Token // Token.AT or Token.MUL
		TokPos   token.Pos   // position of token
		Rbracket token.Pos   // position of ']'
	}

	// DelPrefix represents a prefix deletion expansion.
	// e.g., ${var#prefix} or ${var##prefix}
	DelPrefix struct {
		Longest bool      // true for ## (longest match)
		Hash    token.Pos // position of "#" or "##"
		Val     Expr      // prefix to delete
	}

	// DelSuffix represents a suffix deletion expansion.
	// e.g., ${var%suffix} or ${var%%suffix}
	DelSuffix struct {
		Longest bool      // true for %% (longest match)
		Mod     token.Pos // position of "%" or "%%"
		Val     Expr      // suffix to delete
	}

	// SubstringExp represents a substring expansion.
	// e.g., ${var:offset} or ${var:offset:length}
	SubstringExp struct {
		Colon1 token.Pos // position of ":"
		Offset int       // offset value
		Colon2 token.Pos // position of ":" (for length)
		Length int       // length value
	}

	// ReplaceExp represents a pattern replacement expansion.
	// e.g., ${var/old/new} or ${var//old/new}
	ReplaceExp struct {
		All  bool      // true for // (global replacement)
		Div1 token.Pos // position of "/" or "//"
		Old  string    // pattern to replace
		Div2 token.Pos // position of "/"
		New  string    // replacement string
	}

	// ReplacePrefixExp represents a prefix replacement expansion.
	// e.g., ${var/#old/new}
	ReplacePrefixExp struct {
		Div1 token.Pos // position of "/"
		Hash token.Pos // position of "#"
		Old  string    // prefix pattern
		Div2 token.Pos // position of "/"
		New  string    // replacement string
	}

	// ReplaceSuffixExp represents a suffix replacement expansion.
	// e.g., ${var/%old/new}
	ReplaceSuffixExp struct {
		Div1 token.Pos // position of "/"
		Mod  token.Pos // position of "%"
		Old  string    // suffix pattern
		Div2 token.Pos // position of "/"
		New  string    // replacement string
	}

	// LengthExp represents a length expansion.
	// e.g., ${#var}
	LengthExp struct {
		Hash token.Pos // position of "#"
	}

	// CaseConversionExp represents a case conversion expansion.
	// e.g., ${var,} or ${var^} or ${var,,} or ${var^^}
	CaseConversionExp struct {
		FirstChar bool // true for single character conversion
		ToUpper   bool // true for uppercase conversion
	}

	// OperatorExp represents an operator expansion.
	// e.g., ${var@U} or ${var@L}
	OperatorExp struct {
		At token.Pos   // position of "@"
		Op ExpOperator // expansion operator
	}

	// An ArrExpr node represents an array expression.
	// e.g., (item1 item2 item3)
	ArrExpr struct {
		Lparen token.Pos // position of "("
		Vars   []Expr    // array elements
		Rparen token.Pos // position of ")"
	}

	// A UnaryExpr node represents a unary expression.
	// e.g., !condition or -n string
	UnaryExpr struct {
		OpPos token.Pos   // position of Op
		Op    token.Token // operator
		X     Expr        // operand
	}

	// A PipelineExpr node represents a pipeline expression.
	// e.g., cmd1 | cmd2 or cmd1 |& cmd2
	PipelineExpr struct {
		CtrOp token.Token // Token.Or or Token.OrAnd
		Cmds  []Expr      // list of commands
	}

	// A Redirect represents an I/O redirection.
	// e.g., >file, <file, 2>&1
	Redirect struct {
		N     Expr        // file descriptor number (optional)
		CtrOp token.Token // redirection operator
		OpPos token.Pos   // position of operator
		Word  Expr        // target file or descriptor
	}
)

// Position methods for expressions
func (x *BinaryExpr) Pos() token.Pos  { return x.X.Pos() }
func (x *CmdExpr) Pos() token.Pos     { return x.Name.NamePos }
func (x *CmdListExpr) Pos() token.Pos { return x.Cmds[0].Pos() }
func (x *CmdGroup) Pos() token.Pos    { return x.Opening }
func (x *Redirect) Pos() token.Pos {
	if x.N != nil {
		return x.N.Pos()
	}
	return x.OpPos
}
func (x *Ident) Pos() token.Pos            { return x.NamePos }
func (x *Word) Pos() token.Pos             { return x.ValPos }
func (x *TestExpr) Pos() token.Pos         { return x.Lbrack }
func (x *ExtendedTestExpr) Pos() token.Pos { return x.Lbrack }
func (x *ArithEvalExpr) Pos() token.Pos    { return x.Lparen }
func (x *CmdSubst) Pos() token.Pos         { return x.Dollar }
func (x *ProcSubst) Pos() token.Pos        { return x.Opening }
func (x *ArithExpr) Pos() token.Pos        { return x.Dollar }
func (x *VarExpExpr) Pos() token.Pos       { return x.Dollar }
func (x *ParamExpExpr) Pos() token.Pos     { return x.Dollar }
func (x *IndexExpr) Pos() token.Pos        { return x.X.Pos() }
func (x *ArrExpr) Pos() token.Pos          { return x.Lparen }
func (x *UnaryExpr) Pos() token.Pos        { return x.OpPos }
func (x *PipelineExpr) Pos() token.Pos     { return x.Cmds[0].Pos() }

// End position methods for expressions
func (x *BinaryExpr) End() token.Pos { return x.Y.End() }
func (x *CmdExpr) End() token.Pos {
	if len(x.Recv) > 0 {
		return x.Recv[len(x.Recv)-1].End()
	}
	return token.NoPos
}
func (x *CmdListExpr) End() token.Pos {
	if len(x.Cmds) > 0 {
		return x.Cmds[len(x.Cmds)-1].End()
	}
	return token.NoPos
}
func (x *CmdGroup) End() token.Pos         { return x.Closing }
func (x *Redirect) End() token.Pos         { return x.Word.End() }
func (x *Ident) End() token.Pos            { return token.Pos(int(x.NamePos) + len(x.Name)) }
func (x *Word) End() token.Pos             { return token.Pos(int(x.ValPos) + len(x.Val)) }
func (x *TestExpr) End() token.Pos         { return x.Rbrack }
func (x *ExtendedTestExpr) End() token.Pos { return x.Rbrack }
func (x *ArithEvalExpr) End() token.Pos    { return x.Rparen }
func (x *CmdSubst) End() token.Pos         { return x.Closing }
func (x *ProcSubst) End() token.Pos        { return x.Closing }
func (x *ArithExpr) End() token.Pos        { return x.Rparen }
func (x *VarExpExpr) End() token.Pos       { return x.X.End() }
func (x *ParamExpExpr) End() token.Pos     { return x.Rbrace }
func (x *IndexExpr) End() token.Pos        { return x.Rbrack }
func (x *ArrExpr) End() token.Pos          { return x.Rparen }
func (x *UnaryExpr) End() token.Pos        { return x.X.End() }
func (x *PipelineExpr) End() token.Pos     { return x.Cmds[len(x.Cmds)-1].End() }

// Expression node type assertions
func (*BinaryExpr) exprNode()       {}
func (*CmdExpr) exprNode()          {}
func (*CmdListExpr) exprNode()      {}
func (*CmdGroup) exprNode()         {}
func (*Redirect) exprNode()         {}
func (*Ident) exprNode()            {}
func (*Word) exprNode()             {}
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
func (*UnaryExpr) exprNode()        {}
func (*PipelineExpr) exprNode()     {}

// Parameter expansion node type assertions
func (*DefaultValExp) paramExpNode()       {}
func (*DefaultValAssignExp) paramExpNode() {}
func (*NonNullCheckExp) paramExpNode()     {}
func (*NonNullExp) paramExpNode()          {}
func (*PrefixExp) paramExpNode()           {}
func (*PrefixArrayExp) paramExpNode()      {}
func (*ArrayIndexExp) paramExpNode()       {}
func (*DelPrefix) paramExpNode()           {}
func (*DelSuffix) paramExpNode()           {}
func (*SubstringExp) paramExpNode()        {}
func (*ReplaceExp) paramExpNode()          {}
func (*ReplacePrefixExp) paramExpNode()    {}
func (*ReplaceSuffixExp) paramExpNode()    {}
func (*LengthExp) paramExpNode()           {}
func (*CaseConversionExp) paramExpNode()   {}
func (*OperatorExp) paramExpNode()         {}

// ----------------------------------------------------------------------------
// Expansion operators

// ExpOperator represents a parameter expansion operator.
type ExpOperator string

const (
	// ExpOperatorU converts lowercase alphabetic characters to uppercase.
	ExpOperatorU = "U"

	// ExpOperatoru converts the first character to uppercase if alphabetic.
	ExpOperatoru = "u"

	// ExpOperatorL converts uppercase alphabetic characters to lowercase.
	ExpOperatorL = "L"

	// ExpOperatorQ quotes the value in a format that can be reused as input.
	ExpOperatorQ = "Q"

	// ExpOperatorE expands backslash escape sequences.
	ExpOperatorE = "E"

	// ExpOperatorP expands the value as a prompt string.
	ExpOperatorP = "P"

	// ExpOperatorA produces an assignment statement that recreates the parameter.
	ExpOperatorA = "A"

	// ExpOperatorK produces a quoted version of the value.
	ExpOperatorK = "K"

	// ExpOperatora produces flag values representing parameter attributes.
	ExpOperatora = "a"

	// ExpOperatork expands keys and values to separate words.
	ExpOperatork = "k"
)

// ----------------------------------------------------------------------------
// File structure

// A File node represents a bash source file.
type File struct {
	Doc *CommentGroup // associated documentation

	Stmts []Stmt // list of statements
}

func (*File) Pos() token.Pos { return token.NoPos }
func (*File) End() token.Pos { return token.NoPos }
