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

// Inspect prints a detailed representation of the AST node to the output writer.
func Inspect(node Node, output io.Writer) {
	printer := &inspectPrinter{output: output, indent: 0}
	Walk(printer, node)
}

var _ Visitor = (*inspectPrinter)(nil)

// inspectPrinter is a visitor that prints the AST node in a detailed format.
type inspectPrinter struct {
	output io.Writer
	indent int
}

func (p *inspectPrinter) write(s string) {
	fmt.Fprint(p.output, s)
}

func (p *inspectPrinter) indentWrite(s string) {
	fmt.Fprint(p.output, strings.Repeat("  ", p.indent))
	fmt.Fprint(p.output, s)
}

// Visit implements the Visitor interface for inspectPrinter
func (p *inspectPrinter) Visit(node Node) Visitor {
	if node == nil {
		p.write("nil")
		return nil
	}

	switch n := node.(type) {
	case *File:
		p.write("*ast.File {\n")
		p.indent++
		if len(n.Docs) > 0 {
			p.indentWrite("Docs: [\n")
			p.indent++
			for i, doc := range n.Docs {
				p.indentWrite(fmt.Sprintf("%d: ", i))
				Walk(p, doc)
				p.write("\n")
			}
			p.indent--
			p.indentWrite("]\n")
		}
		if len(n.Stmts) > 0 {
			p.indentWrite("Stmts: [\n")
			p.indent++
			for i, stmt := range n.Stmts {
				p.indentWrite(fmt.Sprintf("%d: ", i))
				Walk(p, stmt)
				p.write("\n")
			}
			p.indent--
			p.indentWrite("]\n")
		}
		p.indent--
		p.indentWrite("}")
	case *CommentGroup:
		p.write("*ast.CommentGroup {\n")
		p.indent++
		p.indentWrite("Comments: [\n")
		p.indent++
		for i, comment := range n.Comments {
			p.indentWrite(fmt.Sprintf("%d: ", i))
			Walk(p, comment)
			p.write("\n")
		}
		p.indent--
		p.indentWrite("]\n")
		p.indent--
		p.indentWrite("}")
	case *Comment:
		p.write(fmt.Sprintf("*ast.Comment (Tok: %s, Text: %q)", n.Tok, n.Text))
	case *IfStmt:
		p.write("*ast.IfStmt {\n")
		p.indent++
		if n.Docs != nil {
			p.indentWrite("Docs: ")
			Walk(p, n.Docs)
			p.write("\n")
		}
		if n.Cond != nil {
			p.indentWrite("Cond: ")
			Walk(p, n.Cond)
			p.write("\n")
		}
		if n.Body != nil {
			p.indentWrite("Body: ")
			Walk(p, n.Body)
			p.write("\n")
		}
		if n.Else != nil {
			p.indentWrite("Else: ")
			Walk(p, n.Else)
			p.write("\n")
		}
		p.indent--
		p.indentWrite("}")
	case *ForStmt:
		p.write("*ast.ForStmt {\n")
		p.indent++
		if n.Docs != nil {
			p.indentWrite("Docs: ")
			Walk(p, n.Docs)
			p.write("\n")
		}
		if n.X != nil {
			p.indentWrite("X: ")
			Walk(p, n.X)
			p.write("\n")
		}
		if n.List != nil {
			p.indentWrite("List: ")
			Walk(p, n.List)
			p.write("\n")
		}
		if n.Body != nil {
			p.indentWrite("Body: ")
			Walk(p, n.Body)
			p.write("\n")
		}
		p.indent--
		p.indentWrite("}")
	case *ExprStmt:
		p.write("*ast.ExprStmt {\n")
		p.indent++
		if n.Docs != nil {
			p.indentWrite("Docs: ")
			Walk(p, n.Docs)
			p.write("\n")
		}
		if n.X != nil {
			p.indentWrite("X: ")
			Walk(p, n.X)
			p.write("\n")
		}
		p.indent--
		p.indentWrite("}")
	case *AssignStmt:
		p.write("*ast.AssignStmt {\n")
		p.indent++
		if n.Docs != nil {
			p.indentWrite("Docs: ")
			Walk(p, n.Docs)
			p.write("\n")
		}
		if n.Opt != nil {
			p.indentWrite("Opt: ")
			Walk(p, n.Opt)
			p.write("\n")
		}
		if n.Lhs != nil {
			p.indentWrite("Lhs: ")
			Walk(p, n.Lhs)
			p.write("\n")
		}
		if n.Rhs != nil {
			p.indentWrite("Rhs: ")
			Walk(p, n.Rhs)
			p.write("\n")
		}
		p.indent--
		p.indentWrite("}")
	case *CallStmt:
		p.write("*ast.CallStmt {\n")
		p.indent++
		if n.Docs != nil {
			p.indentWrite("Docs: ")
			Walk(p, n.Docs)
			p.write("\n")
		}
		p.indentWrite(fmt.Sprintf("Name: %q\n", n.Name))
		p.indentWrite(fmt.Sprintf("IsFile: %t\n", n.IsFile))
		if len(n.Recv) > 0 {
			p.indentWrite("Recv: [\n")
			p.indent++
			for i, arg := range n.Recv {
				p.indentWrite(fmt.Sprintf("%d: ", i))
				Walk(p, arg)
				p.write("\n")
			}
			p.indent--
			p.indentWrite("]\n")
		}
		p.indent--
		p.indentWrite("}")
	case *FuncDecl:
		p.write("*ast.FuncDecl {\n")
		p.indent++
		if n.Docs != nil {
			p.indentWrite("Docs: ")
			Walk(p, n.Docs)
			p.write("\n")
		}
		p.indentWrite(fmt.Sprintf("Name: %q\n", n.Name))
		if n.Body != nil {
			p.indentWrite("Body: ")
			Walk(p, n.Body)
			p.write("\n")
		}
		p.indent--
		p.indentWrite("}")
	case *BlockStmt:
		p.write("*ast.BlockStmt {\n")
		p.indent++
		if len(n.List) > 0 {
			p.indentWrite("List: [\n")
			p.indent++
			for i, stmt := range n.List {
				p.indentWrite(fmt.Sprintf("%d: ", i))
				Walk(p, stmt)
				p.write("\n")
			}
			p.indent--
			p.indentWrite("]\n")
		}
		p.indent--
		p.indentWrite("}")
	case *GotoStmt:
		p.write("*ast.GotoStmt {\n")
		p.indent++
		if n.Docs != nil {
			p.indentWrite("Docs: ")
			Walk(p, n.Docs)
			p.write("\n")
		}
		p.indentWrite(fmt.Sprintf("Label: %q\n", n.Label))
		p.indent--
		p.indentWrite("}")
	case *LabelStmt:
		p.write("*ast.LabelStmt {\n")
		p.indent++
		if n.Docs != nil {
			p.indentWrite("Docs: ")
			Walk(p, n.Docs)
			p.write("\n")
		}
		p.indentWrite(fmt.Sprintf("Name: %q\n", n.Name))
		p.indent--
		p.indentWrite("}")
	case *Word:
		p.write("*ast.Word {\n")
		p.indent++
		if len(n.Parts) > 0 {
			p.indentWrite("Parts: [\n")
			p.indent++
			for i, part := range n.Parts {
				p.indentWrite(fmt.Sprintf("%d: ", i))
				Walk(p, part)
				p.write("\n")
			}
			p.indent--
			p.indentWrite("]\n")
		}
		p.indent--
		p.indentWrite("}")
	case *UnaryExpr:
		p.write("*ast.UnaryExpr {\n")
		p.indent++
		p.indentWrite(fmt.Sprintf("Op: %s\n", n.Op))
		if n.X != nil {
			p.indentWrite("X: ")
			Walk(p, n.X)
			p.write("\n")
		}
		p.indent--
		p.indentWrite("}")
	case *BinaryExpr:
		p.write("*ast.BinaryExpr {\n")
		p.indent++
		if n.X != nil {
			p.indentWrite("X: ")
			Walk(p, n.X)
			p.write("\n")
		}
		p.indentWrite(fmt.Sprintf("Op: %s\n", n.Op))
		if n.Y != nil {
			p.indentWrite("Y: ")
			Walk(p, n.Y)
			p.write("\n")
		}
		p.indent--
		p.indentWrite("}")
	case *CallExpr:
		p.write("*ast.CallExpr {\n")
		p.indent++
		if n.Fun != nil {
			p.indentWrite("Fun: ")
			Walk(p, n.Fun)
			p.write("\n")
		}
		if len(n.Recv) > 0 {
			p.indentWrite("Recv: [\n")
			p.indent++
			for i, arg := range n.Recv {
				p.indentWrite(fmt.Sprintf("%d: ", i))
				Walk(p, arg)
				p.write("\n")
			}
			p.indent--
			p.indentWrite("]\n")
		}
		p.indent--
		p.indentWrite("}")
	case *Lit:
		p.write(fmt.Sprintf("*ast.Lit (Value: %q)", n.Val))
	case *SglQuote:
		p.write("*ast.SglQuote {\n")
		p.indent++
		if n.Val != nil {
			p.indentWrite("Val: ")
			Walk(p, n.Val)
			p.write("\n")
		}
		p.indent--
		p.indentWrite("}")
	case *DblQuote:
		p.write("*ast.DblQuote {\n")
		p.indent++
		p.indentWrite(fmt.Sprintf("DelayedExpansion: %t\n", n.DelayedExpansion))
		if n.Val != nil {
			p.indentWrite("Val: ")
			Walk(p, n.Val)
			p.write("\n")
		}
		p.indent--
		p.indentWrite("}")
	case *CmdExpr:
		p.write("*ast.CmdExpr {\n")
		p.indent++
		if n.Docs != nil {
			p.indentWrite("Docs: ")
			Walk(p, n.Docs)
			p.write("\n")
		}
		if n.Name != nil {
			p.indentWrite("Name: ")
			Walk(p, n.Name)
			p.write("\n")
		}
		if len(n.Recv) > 0 {
			p.indentWrite("Recv: [\n")
			p.indent++
			for i, arg := range n.Recv {
				p.indentWrite(fmt.Sprintf("%d: ", i))
				Walk(p, arg)
				p.write("\n")
			}
			p.indent--
			p.indentWrite("]\n")
		}
		p.indent--
		p.indentWrite("}")
	default:
		p.write(fmt.Sprintf("unknown node type: %T", node))
	}
	return nil
}

var _ Visitor = (*prettyPrinter)(nil)

// prettyPrinter is a visitor that prints the AST node in a pretty format.
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
	if node.Op == token.NOT {
		p.write(" ")
	}
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

// Print prints the string representation of the AST node to the standard output.
func Print(node Node) {
	p := &prettyPrinter{output: os.Stdout, indentSpace: "  "}
	Walk(p, node)
}

// Write writes the string representation of the AST node to the writer.
func Write(node Node, w io.Writer) {
	p := &prettyPrinter{output: w, indentSpace: "  "}
	Walk(p, node)
}

// String returns the string representation of the AST node.
func String(node Node) string {
	var buf strings.Builder
	Write(node, &buf)
	return buf.String()
}
