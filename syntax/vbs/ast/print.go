// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package ast

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/hulo-lang/hulo/syntax/vbs/token"
)

var _ Visitor = (*printer)(nil)

type printer struct {
	output io.Writer
	ident  string
}

func (p *printer) print(a ...any) (n int, err error) {
	return fmt.Fprint(p.output, a...)
}

func (p *printer) printf(format string, a ...any) (n int, err error) {
	return fmt.Fprintf(p.output, format, a...)
}

func (p *printer) println(a ...any) (n int, err error) {
	return fmt.Fprintln(p.output, a...)
}

func (p *printer) writeString(a any) *printer {
	fmt.Fprint(p.output, a)
	return p
}
func (p *printer) spacedString(a ...any) *printer {
	for i, v := range a {
		if i > 0 {
			fmt.Fprint(p.output, " ")
		}
		fmt.Fprint(p.output, v)
	}
	return p
}
func (p *printer) indent() *printer {
	fmt.Fprint(p.output, p.ident)
	return p
}
func (p *printer) space() *printer {
	fmt.Fprint(p.output, " ")
	return p
}

func (p *printer) Visit(node Node) Visitor {
	switch n := node.(type) {
	case *File:
		for _, d := range n.Doc {
			Walk(p, d)
		}

		for _, s := range n.Stmts {
			Walk(p, s)
		}

	case *CommentGroup:
		for _, c := range n.List {
			if !c.Tok.IsValid() {
				c.Tok = token.REM
			}
			p.printf("%s%s %s\n", p.ident, c.Tok, c.Text)
		}

	case *DimDecl:
		list := []string{}
		for _, l := range n.List {
			list = append(list, l.String())
		}
		p.printf("%s %s", p.ident+"Dim", strings.Join(list, ", "))
		if n.Colon.IsValid() {
			p.printf(": %s %s = %s", token.SET, n.Set.Lhs, n.Set.Rhs)
		}
		p.println()

	case *ReDimDecl:
		p.indent().writeString(token.REDIM).space()
		if n.Preserve.IsValid() {
			p.writeString(token.PRESERVE).space()
		}
		list := []string{}
		for _, l := range n.List {
			list = append(list, l.String())
		}
		p.println(strings.Join(list, ", "))

	case *ClassDecl:
		if n.Mod.IsValid() {
			p.printf("%s%s ", p.ident, n.Mod)
		}
		p.printf(p.ident+"Class %s\n", n.Name)

		temp := p.ident
		p.ident += "  "

		for _, s := range n.Stmts {
			Walk(p, s)
		}
		p.ident = temp

		p.println(p.ident + "End Class")

	case *SubDecl:
		ident := p.ident
		if n.Mod.IsValid() {
			p.printf("%s%s ", ident, n.Mod)
			ident = ""
		}
		p.printf(ident+"Sub %s\n", n.Name)

		for _, s := range n.Body.List {
			temp := p.ident
			p.ident += "  "
			Walk(p, s)
			p.ident = temp
		}
		p.indent().spacedString(token.END, token.SUB)
		p.println()

	case *FuncDecl:
		ident := p.ident
		if n.Mod.IsValid() {
			p.print(ident + "Public ")
			ident = ""
		}
		if n.Default.IsValid() {
			if n.Mod == token.PRIVATE {
				panic("Used only with the Public keyword in a Class block to indicate that the Function procedure is the default method for the class.")
			}
			p.print("Default ")
		}
		list := []string{}
		for _, r := range n.Recv {
			if r.TokPos.IsValid() {
				list = append(list, fmt.Sprintf("%s %s", r.Tok, r.Name.Name))
			} else {
				list = append(list, r.Name.Name)
			}
		}
		p.printf(ident+"Function %s(%s)\n", n.Name, strings.Join(list, ", "))

		for _, s := range n.Body.List {
			temp := p.ident
			p.ident += "  "
			Walk(p, s)
			p.ident = temp
		}

		p.println(p.ident + "End Function")

	case *PropertyGetStmt:
		ident := p.ident
		if n.Mod.IsValid() {
			p.printf("%s%s ", ident, n.Mod)
			ident = ""
		}
		list := []string{}
		for _, r := range n.Recv {
			if r.TokPos.IsValid() {
				list = append(list, fmt.Sprintf("%s %s", r.Tok, r.Name.Name))
			} else {
				list = append(list, r.Name.Name)
			}
		}
		p.printf(ident+"Property Get %s(%s)\n", n.Name, strings.Join(list, ", "))

		for _, s := range n.Body.List {
			temp := p.ident
			p.ident += "  "
			Walk(p, s)
			p.ident = temp
		}

		p.println(p.ident + "End Property")

	case *IfStmt:
		p.printf(p.ident+"If %s Then\n", n.Cond)

		for _, s := range n.Body.List {
			temp := p.ident
			p.ident += "  "
			Walk(p, s)
			p.ident = temp
		}

		for n.Else != nil {
			switch el := n.Else.(type) {
			case *IfStmt:
				p.println(p.ident+"ElseIf", el.Cond, "Then")
				Walk(p, el.Body)
				n.Else = el.Else
			case *BlockStmt:
				p.println(p.ident + "Else")
				Walk(p, el)
				n.Else = nil
			}
		}

		p.println(p.ident + "End If")
	case *BlockStmt:
		for _, s := range n.List {
			temp := p.ident
			p.ident += "  "
			Walk(p, s)
			p.ident = temp
		}

	case *ExprStmt:
		p.printf(p.ident+"%s\n", n.X)

	case *PrivateStmt:
		list := []string{}
		for _, l := range n.List {
			list = append(list, l.String())
		}
		p.printf(p.ident+"Private %s\n", strings.Join(list, ", "))

	case *PublicStmt:
		list := []string{}
		for _, l := range n.List {
			list = append(list, l.String())
		}
		p.printf(p.ident+"Public %s\n", strings.Join(list, ", "))

	case *ConstStmt:
		p.printf(p.ident+"%s %s = %s\n", token.CONST, n.Lhs, n.Rhs)

	case *SetStmt:
		p.printf(p.ident+"%s %s = %s\n", token.SET, n.Lhs, n.Rhs)

	case *AssignStmt:
		p.printf(p.ident+"%s = %s\n", n.Lhs, n.Rhs)

	case *ForEachStmt:
		p.printf(p.ident+"For Each %s In %s\n", n.Elem, n.Group)
		for _, s := range n.Body.List {
			temp := p.ident
			p.ident += "  "
			Walk(p, s)
			p.ident = temp
		}
		p.print(p.ident + "Next")
		if n.Stmt != nil {
			p.print(" ")
			Walk(p, n)
		} else {
			p.println()
		}

	case *ForNextStmt:
		p.indent().spacedString(token.FOR, n.Start, token.TO, n.End_)
		if n.Step != nil {
			p.space().spacedString(token.STEP, n.Step).println()
		}
		for _, s := range n.Body.List {
			temp := p.ident
			p.ident += "  "
			Walk(p, s)
			p.ident = temp
		}
		p.println(p.ident + "Next")

	case *WhileWendStmt:
		p.printf(p.ident+"While %s\n", n.Cond)
		for _, s := range n.Body.List {
			temp := p.ident
			p.ident += "  "
			Walk(p, s)
			p.ident = temp
		}
		p.println(p.ident + "Wend")

	case *DoLoopStmt:
		if n.Pre {
			p.print(p.ident + "Do ")
			if n.Cond != nil {
				p.printf("%s %s", n.Tok, n.Cond)
			}
			p.println()
			for _, s := range n.Body.List {
				temp := p.ident
				p.ident += "  "
				Walk(p, s)
				p.ident = temp
			}
			p.println(p.ident + "Loop")
		} else {
			p.printf("%sDo\n", p.ident)
			for _, s := range n.Body.List {
				temp := p.ident
				p.ident += "  "
				Walk(p, s)
				p.ident = temp
			}
			p.print(p.ident + "Loop ")
			if n.Cond != nil {
				p.printf("%s %s", n.Tok, n.Cond)
			}
			p.println()
		}

	case *CallStmt:
		recv := []string{}
		for _, r := range n.Recv {
			recv = append(recv, r.String())
		}
		p.printf(p.ident+"%s %s %s\n", token.CALL, n.Name, strings.Join(recv, ", "))

	case *ExitStmt:
		if n.Tok == token.ILLEGAL {
			p.println(p.ident + "Exit")
		} else {
			p.printf(p.ident+"Exit %s\n", n.Tok)
		}

	case *SelectStmt:
		p.printf(p.ident+"Select Case %s\n", n.Var)
		for _, c := range n.Cases {
			p.printf(p.ident+"  Case %s\n", c.Cond)
			for _, s := range c.Body.List {
				temp := p.ident
				p.ident += "    "
				Walk(p, s)
				p.ident = temp
			}
		}
		if n.Else != nil {
			p.println(p.ident + "  Case Else")
			for _, s := range n.Else.Body.List {
				temp := p.ident
				p.ident += "    "
				Walk(p, s)
				p.ident = temp
			}
		}
		p.println(p.ident + "End Select")

	case *WithStmt:
		p.indent().spacedString(token.WITH, n.Cond).println()
		for _, s := range n.Body.List {
			temp := p.ident
			p.ident += "  "
			Walk(p, s)
			p.ident = temp
		}
		p.indent().spacedString(token.END, token.WITH)
		p.println()

	case *OnErrorStmt:
		if n.OnErrorGoto != nil {
			p.indent().spacedString(token.ON, token.ERROR, token.GOTO, "0")
		}
		if n.OnErrorResume != nil {
			p.indent().spacedString(token.ON, token.ERROR, token.RESUME, token.NEXT)
		}
		p.println()
	case *StopStmt:
		p.println(p.ident + "Stop")
	case *RandomizeStmt:
		p.println(p.ident + "Randomize")
	case *OptionStmt:
		p.println(p.ident + "Option Explicit")
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
