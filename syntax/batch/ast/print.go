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
	case *ExprStmt:
		return p.visitExprStmt(node)
	case *IfStmt:
		return p.visitIfStmt(node)
	case *ForStmt:
		return p.visitForStmt(node)
	case *AssignStmt:
		return p.visitAssignStmt(node)
	case *CallStmt:
		return p.visitCallStmt(node)
	case *FuncDecl:
		return p.visitFuncDecl(node)
	case *BlockStmt:
		return p.visitBlockStmt(node)
	case *GotoStmt:
		return p.visitGotoStmt(node)
	case *LabelStmt:
		return p.visitLabelStmt(node)
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
	case *CmdExpr:
		return p.visitCommand(node)
	default:
		panic("unsupported node type: " + fmt.Sprintf("%T", node))
	}
}

func (p *prettyPrinter) visitGotoStmt(node *GotoStmt) Visitor {
	p.indentWrite("goto :")
	p.write(node.Label)
	p.write("\n")
	return nil
}

func (p *prettyPrinter) visitLabelStmt(node *LabelStmt) Visitor {
	p.indentWrite(":")
	p.write(node.Name)
	p.write("\n")
	return nil
}

func (p *prettyPrinter) visitFuncDecl(node *FuncDecl) Visitor {
	p.indentWrite(":")
	p.write(node.Name)
	p.write("\n")
	Walk(p, node.Body)
	return nil
}

func (p *prettyPrinter) visitCallStmt(node *CallStmt) Visitor {
	p.indentWrite("call ")
	if !node.IsFile {
		p.write(":")
	}
	p.write(node.Name)
	if len(node.Recv) > 0 {
		p.write(" ")
	}
	p.visitExprs(node.Recv)
	p.write("\n")
	return nil
}

func (p *prettyPrinter) visitForStmt(node *ForStmt) Visitor {
	p.indentWrite("for ")
	Walk(p, node.X)
	p.write(" in ")
	Walk(p, node.List)
	p.write(" do (\n")
	Walk(p, node.Body)
	p.indentWrite(")\n")
	return nil
}

func (p *prettyPrinter) visitAssignStmt(node *AssignStmt) Visitor {
	p.indentWrite("set ")
	if node.Opt != nil {
		Walk(p, node.Opt)
		p.write(" ")
	}
	Walk(p, node.Lhs)
	p.write("=")
	Walk(p, node.Rhs)
	p.write("\n")
	return nil
}

func (p *prettyPrinter) visitExprStmt(node *ExprStmt) Visitor {
	p.indentWrite("")
	Walk(p, node.X)
	p.write("\n")
	return nil
}

func (p *prettyPrinter) visitCommand(node *CmdExpr) Visitor {
	Walk(p, node.Name)
	if len(node.Recv) > 0 {
		p.write(" ")
	}
	p.visitExprs(node.Recv)
	return nil
}

func (p *prettyPrinter) visitWord(node *Word) Visitor {
	p.visitExprs(node.Parts)
	return nil
}

func (p *prettyPrinter) visitUnaryExpr(node *UnaryExpr) Visitor {
	p.write(node.Op.String())
	Walk(p, node.X)
	return nil
}

func (p *prettyPrinter) visitCallExpr(node *CallExpr) Visitor {
	Walk(p, node.Fun)
	if len(node.Recv) > 0 {
		p.write(" ")
	}
	p.visitExprs(node.Recv)
	return nil
}

func (p *prettyPrinter) visitLit(node *Lit) Visitor {
	p.write(node.Val)
	return nil
}

func (p *prettyPrinter) visitSglQuote(node *SglQuote) Visitor {
	p.write("%")
	Walk(p, node.Val)
	return nil
}

func (p *prettyPrinter) visitDblQuote(node *DblQuote) Visitor {
	if node.DelayedExpansion {
		p.write("!")
	} else {
		p.write("%")
	}
	Walk(p, node.Val)
	if node.DelayedExpansion {
		p.write("!")
	} else {
		p.write("%")
	}
	return nil
}

func (p *prettyPrinter) visitBlockStmt(node *BlockStmt) Visitor {
	p.indent++
	for _, stmt := range node.List {
		Walk(p, stmt)
	}
	p.indent--
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
		p.indentWrite("::")
	} else {
		p.indentWrite("REM ")
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
			p.write(" (\n")
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

func (p *prettyPrinter) visitExprs(exprs []Expr, sep ...string) {
	sepStr := " "
	if len(sep) > 0 {
		sepStr = sep[0]
	}
	for i, expr := range exprs {
		Walk(p, expr)
		if i < len(exprs)-1 {
			p.write(sepStr)
		}
	}
}

func Print(node Node) {
	p := &prettyPrinter{output: os.Stdout, indentSpace: "  "}
	Walk(p, node)
}

func Write(node Node, w io.Writer) {
	p := &prettyPrinter{output: w, indentSpace: "  "}
	Walk(p, node)
}

func String(node Node) string {
	var buf strings.Builder
	Write(node, &buf)
	return buf.String()
}
