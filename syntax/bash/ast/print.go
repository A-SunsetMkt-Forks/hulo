// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package ast

import (
	"fmt"
	"io"
	"os"
	"strings"
)

var _ Visitor = (*printer)(nil)

type printer struct {
	output io.Writer
	ident  int
}

func (p *printer) print(a ...any) {
	fmt.Fprint(p.output, a...)
}

func (p *printer) printf(format string, a ...any) {
	fmt.Fprintf(p.output, strings.Repeat("  ", p.ident)+format, a...)
}

func (p *printer) println(a ...any) {
	fmt.Fprint(p.output, strings.Repeat("  ", p.ident))
	fmt.Fprintln(p.output, a...)
}

func (p *printer) Visit(node Node) Visitor {
	switch node := node.(type) {
	case *File:
		for _, d := range node.Decls {
			Walk(p, d)
		}

		for _, s := range node.Stmts {
			Walk(p, s)
		}
	case *FuncDecl:
		p.printf("function %s {\n", node.Name)
		p.ident++
		Walk(p, node.Body)
		p.ident--
		p.println("}")

	case *BlockStmt:
		for _, s := range node.List {
			Walk(p, s)
		}

	case *WhileStmt:
		p.printf("while %s; do\n", node.Cond)
		p.ident++
		Walk(p, node.Body)
		p.ident--
		p.println("done")

	case *UntilStmt:
		p.printf("until %s; do\n", node.Cond)
		p.ident++
		Walk(p, node.Body)
		p.ident--
		p.println("done")

	case *ForStmt:
		p.printf("for ((%s; %s; %s)); do\n", node.Init, node.Cond, node.Post)
		p.ident++
		Walk(p, node.Body)
		p.ident--
		p.println("done")

	case *ForInStmt:
		p.printf("for %s in %s; do\n", node.Var, node.List)
		p.ident++
		Walk(p, node.Body)
		p.ident--
		p.println("done")

	case *SelectStmt:
		p.printf("select %s in %s; do\n", node.Var, node.List)
		p.ident++
		Walk(p, node.Body)
		p.ident--
		p.println("done")

	case *IfStmt:
		p.printf("if %s; then\n", node.Cond)
		p.ident++
		Walk(p, node.Body)
		p.ident--
		for node.Else != nil {
			switch el := node.Else.(type) {
			case *IfStmt:
				p.printf("elif %s; then\n", el.Cond)
				p.ident++
				Walk(p, el.Body)
				p.ident--
				node.Else = el.Else
			case *BlockStmt:
				p.println("else")
				Walk(p, el)
				node.Else = nil
			}
		}

		p.println("fi")

	case *CaseStmt:
		p.printf("case %s in\n", node.X)

		for _, pattern := range node.Patterns {
			p.printf("  %s)\n", pattern.Conds)
			p.ident += 2
			Walk(p, pattern.Body)
			p.ident -= 2
			p.println(";;")
		}

		if node.Else != nil {
			p.println("  *)")
			p.ident += 2
			Walk(p, node.Else)
			p.ident -= 2
		}

		p.println("esac")

	case *AssignStmt:
		if node.Local.IsValid() {
			p.printf("local %s=%s\n", node.Lhs, node.Rhs)
		} else {
			p.printf("%s=%s\n", node.Lhs, node.Rhs)
		}
	case *BreakStmt:
		p.println("break")
	case *ContinueStmt:
		p.println("continue")
	}
	return nil
}

func Print(node Node) {
	Walk(&printer{ident: 0, output: os.Stdout}, node)
}

func String(node Node) string {
	buf := &strings.Builder{}
	Walk(&printer{ident: 0, output: buf}, node)
	return buf.String()
}
