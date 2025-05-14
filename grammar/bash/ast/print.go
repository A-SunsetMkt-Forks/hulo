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
		p.println(token.FUNCTION, n.Name.Text(), token.LPAREN+token.RPAREN, token.LBRACE)
		temp := p.ident
		p.ident += "  "
		p.block(n.Body.List)
		p.ident = temp
		println(token.RBRACE)

	case *AssignStmt:
		if n.Local.IsValid() {
			p.println_c(token.LOCAL, token.SPACE, n.Lhs.Text(), token.ASSIGN, n.Rhs.Text())
		} else {
			p.println_c(n.Lhs.Text(), token.ASSIGN, n.Rhs.Text())
		}

	case *ReturnStmt:
		p.println(token.RETURN, n.X.Text())

	case *BreakStmt:
		p.println(token.BREAK)

	case *ContinueStmt:
		p.println(token.CONTINUE)

	case *WhileStmt:
		p.println(token.WHILE, n.Cond.Text()).
			println(token.DO)

		temp := p.ident
		p.ident += "  "
		p.block(n.Body.List)
		p.ident = temp

		p.println(token.DONE)

	case *UntilStmt:
		p.println(token.UNTIL, n.Cond.Text()).
			println(token.DO)

		temp := p.ident
		p.ident += "  "
		p.block(n.Body.List)
		p.ident = temp

		p.println(token.DONE)

	case *ForeachStmt:
		group := []string{}
		for _, g := range n.Group {
			group = append(group, g.Text())
		}
		p.println(token.FOR, n.Elem.Text(), token.IN, strings.Join(group, " "), token.SEMI).
			println(token.DO)

		temp := p.ident
		p.ident += "  "
		p.block(n.Body.List)
		p.ident = temp

		p.println(token.DONE)

	case *ForStmt:
		p.println_c(token.FOR, token.SPACE, token.DOUBLE_LPAREN,
			n.Init.Text(), token.SEMI, token.SPACE,
			n.Cond.Text(), token.SEMI, token.SPACE,
			n.Post.Text(), token.SPACE, token.DOUBLE_RPAREN,
			token.SEMI).
			println(token.DO)

		temp := p.ident
		p.ident += "  "
		p.block(n.Body.List)
		p.ident = temp

		p.println(token.DONE)

	case *IfStmt:
		p.println_c(token.IF, token.SPACE, n.Cond.Text(), token.SEMI, token.SPACE, token.THEN)

		temp := p.ident
		p.ident += "  "
		p.block(n.Body.List)
		p.ident = temp

		for n.Else != nil {
			switch el := n.Else.(type) {
			case *IfStmt:

			case *BlockStmt:
				fmt.Fprint(p.output, p.ident+"else")
				p.print(el)
				n.Else = nil
			}
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
		p.println(token.CASE, n.Var.Text(), token.IN)

		for _, c := range n.Cases {
			cs := []string{}
			for _, cond := range c.Conds {
				cs = append(cs, cond.Text())
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
		p.println(token.SELECT, n.Elem.Text(), token.IN, n.Group.Text(), token.SEMI).
			println(token.DO)

		temp := p.ident
		p.ident += "  "
		p.block(n.Body.List)
		p.ident = temp

		p.println(token.DONE)

	case *ExprStmt:
		p.println(n.X.Text())
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
