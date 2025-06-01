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

func (p *printer) Visit(node Node) Visitor {
	switch n := node.(type) {
	case *File:
		for _, d := range n.Decls {
			Walk(p, d)
		}

		for _, s := range n.Stmts {
			Walk(p, s)
		}

	case *DimDecl:
		p.printf("%s %s", p.ident+"Dim", ExprListStr(n.List))
		if n.Colon.IsValid() {
			p.printf(": Set %s = %s", ExprStr(n.Set.Lhs), ExprStr(n.Set.Rhs))
		}
		p.println()

	case *ReDimDecl:
		p.print("ReDim ")
		if n.Preserve.IsValid() {
			p.print("Preserve ")
		}
		p.println(ExprListStr(n.List))

	case *ClassDecl:
		if n.Mod.HasPublic() {
			p.print(p.ident + "Public ")
		}
		p.printf(p.ident+"Class %s\n", ExprStr(n.Name))

		temp := p.ident
		p.ident += "  "
		for _, d := range n.Decls {
			Walk(p, d)
		}

		for _, s := range n.Stmts {
			Walk(p, s)
		}
		p.ident = temp

		p.println(p.ident + "End Class")

	case *SubDecl:
		ident := p.ident
		if n.Mod.HasPublic() {
			p.print(p.ident + "Public ")
			ident = ""
		}
		p.printf(ident+"Sub %s\n", ExprStr(n.Name))

		for _, s := range n.Body.List {
			temp := p.ident
			p.ident += "  "
			Walk(p, s)
			p.ident = temp
		}

		p.println(p.ident + "End Sub")

	case *FuncDecl:
		ident := p.ident
		if n.Mod.HasPublic() {
			p.print(ident + "Public ")
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
		p.printf(ident+"Function %s(%s)\n", ExprStr(n.Name), strings.Join(list, ", "))

		for _, s := range n.Body.List {
			temp := p.ident
			p.ident += "  "
			Walk(p, s)
			p.ident = temp
		}

		p.println(p.ident + "End Function")

	case *PropertyDecl:
		ident := p.ident
		if n.Mod.HasPublic() {
			p.print(ident + "Public ")
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
		p.printf(ident+"Property %s %s(%s)\n", n.Tok, ExprStr(n.Name), strings.Join(list, ", "))

		for _, s := range n.Body.List {
			temp := p.ident
			p.ident += "  "
			Walk(p, s)
			p.ident = temp
		}

		p.println(p.ident + "End Property")

	case *IfStmt:
		p.printf(p.ident+"If %s Then\n", ExprStr(n.Cond))

		for _, s := range n.Body.List {
			temp := p.ident
			p.ident += "  "
			Walk(p, s)
			p.ident = temp
		}

		if n.ElseIf != nil {
			for _, elif := range n.ElseIf {
				p.println(p.ident+"ElseIf", ExprStr(elif.Cond), "Then")
				for _, s := range elif.Body.List {
					temp := p.ident
					p.ident += "  "
					Walk(p, s)
					p.ident = temp
				}
			}
		}

		if n.Else != nil {
			p.println(p.ident + "Else")
			for _, s := range n.Else.List {
				temp := p.ident
				p.ident += "  "
				Walk(p, s)
				p.ident = temp
			}
		}

		p.println(p.ident + "End If")

	case *ExprStmt:
		p.printf(p.ident+"%s\n", ExprStr(n.X))

	case *MemberStmt:
		if n.Mod.HasPublic() {
			p.print(p.ident + "Public ")
		}
		if n.Mod.HasPrivate() {
			p.print(p.ident + "Private ")
		}
		p.println(ExprStr(n.Name))

	case *AssignStmt:
		modifier := ""
		switch n.Tok {
		case token.SET:
			modifier = "Set "
		case token.CONST:
			modifier = "Const "
		}
		p.printf(p.ident+"%s%s = %s\n", modifier, ExprStr(n.Lhs), ExprStr(n.Rhs))

	case *ForEachStmt:
		p.printf(p.ident+"For Each %s In %s\n", ExprStr(n.Elem), ExprStr(n.Group))
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
		p.printf(p.ident+"For %s To %s", ExprStr(n.Start), ExprStr(n.End_))
		if n.Step != nil {
			p.printf(" Step %s\n", ExprStr(n.Step))
		}
		for _, s := range n.Body.List {
			temp := p.ident
			p.ident += "  "
			Walk(p, s)
			p.ident = temp
		}
		p.println(p.ident + "Next")

	case *WhileWendStmt:
		p.printf(p.ident+"While %s\n", ExprStr(n.Cond))
		for _, s := range n.Body.List {
			temp := p.ident
			p.ident += "  "
			Walk(p, s)
			p.ident = temp
		}
		p.println(p.ident + "Wend")

	case *DoLoopStmt:
		if n.Pre {
			p.printf(p.ident+"Do %s %s\n", n.Tok, ExprStr(n.Cond))
			for _, s := range n.Body.List {
				temp := p.ident
				p.ident += "  "
				Walk(p, s)
				p.ident = temp
			}
			p.println(p.ident + "Loop")
		} else {
			p.printf(p.ident + "Do\n")
			for _, s := range n.Body.List {
				temp := p.ident
				p.ident += "  "
				Walk(p, s)
				p.ident = temp
			}
			p.printf(p.ident+"Loop %s %s\n", n.Tok, ExprStr(n.Cond))
		}

	case *CallStmt:
		p.printf(p.ident+"Call %s %s\n", ExprStr(n.Name), ExprListStr(n.Recv))

	case *ExitStmt:
		if len(n.X) == 0 {
			p.println(p.ident + "Exit")
		} else {
			p.printf(p.ident+"Exit %s\n", n.X)
		}

	case *SelectStmt:
		p.printf(p.ident+"Select Case %s\n", ExprStr(n.Var))
		for _, c := range n.Cases {
			p.printf(p.ident+"  Case %s\n", ExprStr(c.Cond))
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
		p.printf(p.ident+"With %s\n", ExprStr(n.Cond))
		for _, s := range n.Body.List {
			temp := p.ident
			p.ident += "  "
			Walk(p, s)
			p.ident = temp
		}
		p.println(p.ident + "End With")

	case *OnErrorStmt:
		if n.OnErrorGoto != nil {
			p.println(p.ident + "On Error GoTo 0")
		}
		if n.OnErrorResume != nil {
			p.println(p.ident + "On Error Resume Next")
		}
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

func ExprStr(e Expr) string {
	switch e := e.(type) {
	case *Ident:
		return e.Name
	case *BasicLit:
		if e.Kind == token.STRING {
			return fmt.Sprintf(`"%s"`, e.Value)
		}
		return e.Value
	case *SelectorExpr:
		return fmt.Sprintf("%s.%s", ExprStr(e.X), ExprStr(e.Sel))
	case *BinaryExpr:
		return fmt.Sprintf("%s %s %s", ExprStr(e.X), e.Op, ExprStr(e.Y))
	case *CallExpr:
		return fmt.Sprintf("%s(%s)", ExprStr(e.Func), ExprListStr(e.Recv))
	case *IndexExpr:
		return fmt.Sprintf("%s(%s)", ExprStr(e.X), ExprStr(e.Index))
	case *IndexListExpr:
		return fmt.Sprintf("%s(%s)", ExprStr(e.X), ExprListStr(e.Indices))
	}
	return ""
}

func ExprListStr(list []Expr) string {
	res := []string{}
	for _, e := range list {
		res = append(res, ExprStr(e))
	}
	return strings.Join(res, ", ")
}
