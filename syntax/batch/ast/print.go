// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package ast

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/hulo-lang/hulo/syntax/batch/token"
)

var _ Visitor = (*prettyPrinter)(nil)

type prettyPrinter struct {
	output      io.Writer
	indent      int
	indentSpace string
}

// write indents and writes a string to the output.
func (p *prettyPrinter) write(s string) {
	fmt.Fprint(p.output, s)
}

// indentWrite writes an indented string.
func (p *prettyPrinter) indentWrite(s string) {
	fmt.Fprint(p.output, strings.Repeat(p.indentSpace, p.indent))
	fmt.Fprint(p.output, s)
}

func (p *prettyPrinter) Visit(node Node) Visitor {
	switch node := node.(type) {
	case *File:
		return p.visitFile(node)
	case *CommentGroup:
		return p.visitCommentGroup(node)
	case *Comment:
		return p.visitComment(node)
	case *IfStmt:
		return p.visitIfStmt(node)
	case *BlockStmt:
		return p.visitBlockStmt(node)
	case *Word:
		return p.visitWord(node)
	case *UnaryExpr:
		return p.visitUnaryExpr(node)
	case *BinaryExpr:
		return p.visitBinaryExpr(node)
	case *CallExpr:
		return p.visitCallExpr(node)
	case *Lit:
		return p.visitLit(node)
	case *SglQuote:
		return p.visitSglQuote(node)
	case *DblQuote:
		return p.visitDblQuote(node)
	case *Command:
		return p.visitCommand(node)
	default:
		panic("unsupported node type: " + fmt.Sprintf("%T", node))
	}
}

func (p *prettyPrinter) visitCommand(node *Command) Visitor {
	Walk(p, node.Name)
	p.write(" ")
	for i, recv := range node.Recv {
		Walk(p, recv)
		if i < len(node.Recv)-1 {
			p.write(" ")
		}
	}
	return nil
}

func (p *prettyPrinter) visitWord(node *Word) Visitor {
	for i, part := range node.Parts {
		Walk(p, part)
		if i < len(node.Parts)-1 {
			p.write(" ")
		}
	}
	return nil
}

func (p *prettyPrinter) visitUnaryExpr(node *UnaryExpr) Visitor {
	p.write(node.Op.String())
	Walk(p, node.X)
	return nil
}

func (p *prettyPrinter) visitCallExpr(node *CallExpr) Visitor {
	Walk(p, node.Fun)
	p.write(" ")
	for i, recv := range node.Recv {
		Walk(p, recv)
		if i < len(node.Recv)-1 {
			p.write(" ")
		}
	}
	return nil
}

func (p *prettyPrinter) visitLit(node *Lit) Visitor {
	p.write(node.Val)
	return nil
}

func (p *prettyPrinter) visitSglQuote(node *SglQuote) Visitor {
	p.write("%%")
	Walk(p, node.Val)
	return nil
}

func (p *prettyPrinter) visitDblQuote(node *DblQuote) Visitor {
	if node.DelayedExpansion {
		p.write("!")
	} else {
		p.write("%%")
	}
	Walk(p, node.Val)
	if node.DelayedExpansion {
		p.write("!")
	} else {
		p.write("%%")
	}
	return nil
}

func (p *prettyPrinter) visitBlockStmt(node *BlockStmt) Visitor {
	for _, stmt := range node.List {
		Walk(p, stmt)
		p.write("\n")
	}
	return nil
}

func (p *prettyPrinter) visitBinaryExpr(node *BinaryExpr) Visitor {
	Walk(p, node.X)
	p.write(" ")
	p.write(node.Op.String())
	p.write(" ")
	Walk(p, node.Y)
	return nil
}

func (p *prettyPrinter) visitFile(node *File) Visitor {
	for _, doc := range node.Docs {
		Walk(p, doc)
	}
	for _, stmt := range node.Stmts {
		Walk(p, stmt)
	}
	return nil
}

func (p *prettyPrinter) visitCommentGroup(node *CommentGroup) Visitor {
	for _, comment := range node.Comments {
		Walk(p, comment)
		p.write("\n")
	}
	return nil
}

func (p *prettyPrinter) visitComment(node *Comment) Visitor {
	if node.Tok == token.DOUBLE_COLON {
		p.write("::")
	} else {
		p.write("REM ")
	}
	p.write(node.Text)
	return nil
}

func (p *prettyPrinter) visitIfStmt(node *IfStmt) Visitor {
	p.indentWrite("if ")
	Walk(p, node.Cond)
	p.write(" (\n")

	Walk(p, node.Body)

	p.indentWrite(")\n")

	for node.Else != nil {
		switch el := node.Else.(type) {
		case *IfStmt:
			p.indentWrite("else if ")
			Walk(p, el.Cond)
			p.write("(\n")
			Walk(p, el.Body)
			p.indentWrite(")\n")
			node.Else = el.Else
		case *BlockStmt:
			p.indentWrite("else (\n")
			Walk(p, el)
			p.indentWrite(")\n")
			node.Else = nil
		}
	}
	return nil
}

func Print(node Node) {
	p := &prettyPrinter{output: os.Stdout, indentSpace: "  "}
	Walk(p, node)
}

func Write(node Node, w io.Writer) {
	print(node, "", w)
}

func print(node Node, ident string, w io.Writer) {
	switch n := node.(type) {
	case *IfStmt:
		fmt.Fprintf(w, "%sif %s (\n", ident, n.Cond)
		print(n.Body, ident+"  ", w)
		fmt.Print(ident + ")")
		for n.Else != nil {
			fmt.Print(" ")
			switch el := n.Else.(type) {
			case *IfStmt:
				fmt.Fprintf(w, "%selse if %s (\n", ident, el.Cond)
				print(el.Body, ident+"  ", w)
				fmt.Fprintf(w, "%s)\n", ident)
				n.Else = el.Else
			case *BlockStmt:
				fmt.Fprintf(w, "%selse (\n", ident)
				print(el, ident+"  ", w)
				fmt.Fprintf(w, "%s)\n", ident)
				n.Else = nil
			}
		}
		fmt.Fprintln(w)

	case *ForStmt:
		fmt.Fprintf(w, "%sfor %s in %s do (\n", ident, n.X, n.List)
		print(n.Body, ident+"  ", w)
		fmt.Fprintf(w, "%s)\n", ident)
	case *ExprStmt:
		fmt.Fprintf(w, "%s%s\n", ident, n.X)
	case *AssignStmt:
		fmt.Fprintf(w, "%sset %s=%s\n", ident, n.Lhs, n.Rhs)
	case *CallStmt:
		fmt.Fprintf(w, "%scall :%s", ident, n.Name)
		if len(n.Recv) > 0 {
			fmt.Print(" ")
			for i, recv := range n.Recv {
				if i > 0 {
					fmt.Print(" ")
				}
				fmt.Print(recv)
			}
		}
		fmt.Fprintln(w)
	case *FuncDecl:
		fmt.Fprintf(w, "%s:%s\n", ident, n.Name)
		print(n.Body, ident+"  ", w)
	case *BlockStmt:
		for _, stmt := range n.List {
			print(stmt, ident, w)
		}
	case *GotoStmt:
		fmt.Fprintf(w, "%sgoto :%s\n", ident, n.Label)
	case *LabelStmt:
		fmt.Fprintf(w, "%s:%s\n", ident, n.Name)
	case *Command:
		fmt.Fprintf(w, "%s%s", ident, n.Name)
		if len(n.Recv) > 0 {
			fmt.Fprint(w, " ")
			for i, recv := range n.Recv {
				if i > 0 {
					fmt.Fprint(w, " ")
				}
				fmt.Fprint(w, recv)
			}
		}
		fmt.Fprintln(w)
	case *Comment:
		if n.Tok == token.DOUBLE_COLON {
			fmt.Fprintf(w, "%s::%s\n", ident, n.Text)
		} else {
			fmt.Fprintf(w, "%sREM %s\n", ident, n.Text)
		}
	case *CommentGroup:
		for _, comment := range n.Comments {
			print(comment, ident, w)
		}
	default:
		fmt.Fprintf(w, "%T\n", node)
	}
}
