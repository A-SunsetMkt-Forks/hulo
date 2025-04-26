// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package ast

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/hulo-lang/hulo/grammar/bash/token"
)

var _ Visitor = (*printer)(nil)

type printer struct {
	output io.Writer
	ident  string
	temp   string
}

func (p *printer) print(toks ...any) *printer {
	if len(p.ident) != 0 {
		toks = append([]any{p.ident}, toks...)
	}
	fmt.Fprint(p.output, toks...)
	return p
}

// println prints tokens with new line
func (p *printer) println(toks ...any) *printer {
	if len(p.ident) != 0 {
		toks = append([]any{p.ident}, toks...)
	}
	fmt.Fprintln(p.output, toks...)
	return p
}

// println_c prints tokens with new line and compressing
func (p *printer) println_c(toks ...string) *printer {
	if len(p.ident) != 0 {
		toks = append([]string{p.ident}, toks...)
	}
	result := ""
	for _, tok := range toks {
		result += tok
	}
	fmt.Fprintln(p.output, result)
	return p
}

func (p *printer) block(stmts []Stmt) *printer {
	for _, s := range stmts {
		Walk(p, s)
	}
	return p
}

func (p *printer) Visit(node Node) Visitor {
	switch n := node.(type) {
	case *File:
		for _, d := range n.Decls {
			Walk(p, d)
		}

		for _, s := range n.Stmts {
			Walk(p, s)
		}

	case *FuncDecl:
		p.println(token.FUNCTION, ExprStr(n.Name), token.LPAREN+token.RPAREN, token.LBRACE)
		temp := p.ident
		p.ident += "  "
		p.block(n.Body.List)
		p.ident = temp
		println(token.RBRACE)

	case *AssignStmt:
		if n.Local {
			p.println_c(token.LOCAL, token.SPACE, ExprStr(n.Lhs), token.ASSIGN, ExprStr(n.Rhs))
		} else {
			p.println_c(ExprStr(n.Lhs), token.ASSIGN, ExprStr(n.Rhs))
		}

	case *ReturnStmt:
		p.println(token.RETURN, ExprStr(n.X))

	case *BreakStmt:
		p.println(token.BREAK)

	case *ContinueStmt:
		p.println(token.CONTINUE)

	case *WhileStmt:
		p.println(token.WHILE, ExprStr(n.Cond)).
			println(token.DO)

		temp := p.ident
		p.ident += "  "
		p.block(n.Body.List)
		p.ident = temp

		p.println(token.DONE)

	case *UntilStmt:
		p.println(token.UNTIL, ExprStr(n.Cond)).
			println(token.DO)

		temp := p.ident
		p.ident += "  "
		p.block(n.Body.List)
		p.ident = temp

		p.println(token.DONE)

	case *ForeachStmt:
		p.println(token.FOR, ExprStr(n.Elem), token.IN, ExprListStr(n.Group), token.SEMI).
			println(token.DO)

		temp := p.ident
		p.ident += "  "
		p.block(n.Body.List)
		p.ident = temp

		p.println(token.DONE)

	case *ForStmt:
		p.println_c(token.FOR, token.SPACE, token.DOUBLE_LPAREN,
			ExprStr(n.Init), token.SEMI, token.SPACE,
			ExprStr(n.Cond), token.SEMI, token.SPACE,
			ExprStr(n.Post), token.SPACE, token.DOUBLE_RPAREN,
			token.SEMI).
			println(token.DO)

		temp := p.ident
		p.ident += "  "
		p.block(n.Body.List)
		p.ident = temp

		p.println(token.DONE)

	case *IfStmt:
		p.println_c(token.IF, token.SPACE, ExprStr(n.Cond), token.SEMI, token.SPACE, token.THEN)

		temp := p.ident
		p.ident += "  "
		p.block(n.Body.List)
		p.ident = temp

		for _, elif := range n.Elif {
			p.println_c(token.ELIF, token.SPACE, ExprStr(elif.Cond), token.SEMI, token.SPACE, token.THEN)

			temp := p.ident
			p.ident += "  "
			p.block(n.Body.List)
			p.ident = temp
		}

		if n.Else != nil {
			p.println(token.ELSE)
			temp := p.ident
			p.ident += "  "
			p.block(n.Body.List)
			p.ident = temp
		}

		p.println(token.FI)

	case *CaseStmt:
		p.println(token.CASE, ExprStr(n.Var), token.IN)

		for _, c := range n.Cases {
			cs := []string{}
			for _, cond := range c.Conds {
				cs = append(cs, ExprStr(cond))
			}

			p.println_c(strings.Join(cs, " | "), token.RPAREN)

			temp := p.ident
			p.ident += "  "
			p.block(c.Body.List)
			p.ident = temp

			temp = p.ident
			p.ident += "  "
			p.println(token.DOUBLE_SEMI)
			p.ident = temp
		}

		if n.Else != nil {
			p.println_c(token.MUL, token.RPAREN)

			temp := p.ident
			p.ident += "  "
			p.block(n.Else.List)
			p.ident = temp

			temp = p.ident
			p.ident += "  "
			p.println(token.DOUBLE_SEMI)
			p.ident = temp
		}

		p.println(token.ESAC)

	case *SelectStmt:
		p.println(token.SELECT, ExprStr(n.Elem), token.IN, ExprStr(n.Group), token.SEMI).
			println(token.DO)

		temp := p.ident
		p.ident += "  "
		p.block(n.Body.List)
		p.ident = temp

		p.println(token.DONE)

	case *ExprStmt:
		p.println(ExprStr(n.X))
	}
	return nil
}

func Print(node Node) {
	Walk(&printer{ident: "", output: os.Stdout}, node)
}

func String(node Node) string {
	buf := &strings.Builder{}
	Walk(&printer{ident: "", output: buf}, node)
	return buf.String()
}

func ExprStr(e Expr) string {
	switch e := e.(type) {
	case *Ident:
		return e.Name

	case *BasicLit:
		if e.Kind == token.STRING {
			return fmt.Sprintf(`"%s"`, e.Value)
		}
		return e.Value

	case *CallExpr:
		return fmt.Sprintf("%s %s", ExprStr(e.Func), ExprListStr(e.Recv))

	case *BasicTestExpr:
		return fmt.Sprintf("[ %s ]", ExprStr(e.X))

	case *ExtendedTestExpr:
		return fmt.Sprintf("[[ %s ]]", ExprStr(e.X))

	case *ArithEvalExpr:
		return fmt.Sprintf("(( %s ))", ExprStr(e.X))

	case *CmdSubst:
		if e.Tok == token.LPAREN {
			return fmt.Sprintf("$( %s )", ExprStr(e.X))
		}
		return fmt.Sprintf("` %s `", ExprStr(e.X))

	case *ProcSubst:
		if e.Tok == token.LT {
			return fmt.Sprintf("<( %s )", ExprStr(e.X))
		}
		return fmt.Sprintf(">( %s )", ExprStr(e.X))

	case *ArithExp:
		return fmt.Sprintf("$(( %s ))", ExprStr(e.X))

	case *BinaryExpr:
		if e.Op == token.NONE {
			if e.Compress {
				return fmt.Sprintf("%s%s", ExprStr(e.X), ExprStr(e.Y))
			}
			return fmt.Sprintf("%s %s", ExprStr(e.X), ExprStr(e.Y))
		}
		if e.Compress {
			return fmt.Sprintf("%s%s%s", ExprStr(e.X), e.Op, ExprStr(e.Y))
		}
		return fmt.Sprintf("%s %s %s", ExprStr(e.X), e.Op, ExprStr(e.Y))

	case *ParamExp:
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
			if e.Tok == token.MUL {
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
	}
	return ""
}

func ExprListStr(list []Expr) string {
	res := []string{}
	for _, e := range list {
		res = append(res, ExprStr(e))
	}
	return strings.Join(res, " ")
}
