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

// Visit visits the AST node.
func (p *prettyPrinter) Visit(node Node) Visitor {
	switch node := node.(type) {
	case *File:
		return p.visitFile(node)
	case *BlcokStmt:
		return p.visitBlcokStmt(node)
	case *FuncDecl:
		return p.visitFuncDecl(node)
	case *ProcessDecl:
		return p.visitProcessDecl(node)
	case *Attribute:
		return p.visitAttribute(node)
	case *ExprStmt:
		return p.visitExprStmt(node)
	case *CommentGroup:
		return p.visitCommentGroup(node)
	case *SingleLineComment:
		return p.visitSingleLineComment(node)
	case *DelimitedComment:
		return p.visitDelimitedComment(node)
	case *AssignStmt:
		return p.visitAssignStmt(node)
	case *SwitchStmt:
		return p.visitSwitchStmt(node)
	case *CaseClause:
		return p.visitCaseClause(node)
	case *TryStmt:
		return p.visitTryStmt(node)
	case *CatchClause:
		return p.visitCatchClause(node)
	case *TrapStmt:
		return p.visitTrapStmt(node)
	case *DataStmt:
		return p.visitDataStmt(node)
	case *DynamicparamStmt:
		return p.visitDynamicparamStmt(node)
	case *HashTable:
		return p.visitHashTable(node)
	case *HashEntry:
		return p.visitHashEntry(node)
	case *IncDecExpr:
		return p.visitIncDecExpr(node)
	case *IndexExpr:
		return p.visitIndexExpr(node)
	case *IndicesExpr:
		return p.visitIndicesExpr(node)
	case *VarExpr:
		return p.visitVarExpr(node)
	case *Parameter:
		return p.visitParameter(node)
	case *CallExpr:
		return p.visitCallExpr(node)
	case *CmdExpr:
		return p.visitCmdExpr(node)
	case *MemberAccess:
		return p.visitMemberAccess(node)
	case *StaticMemberAccess:
		return p.visitStaticMemberAccess(node)
	case *ArrayExpr:
		return p.visitArrayExpr(node)
	case *ArrayConstructor:
		return p.visitArrayConstructor(node)
	case *SelectExpr:
		return p.visitSelectExpr(node)
	case *ConstrainedVarExpr:
		return p.visitConstrainedVarExpr(node)
	case *Lit:
		return p.visitLit(node)
	case *Ident:
		p.write(node.Name)
		return nil
	default:
		panic("unsupported node type: " + fmt.Sprintf("%T", node))
	}
}

func (p *prettyPrinter) visitMemberAccess(node *MemberAccess) Visitor {
	p.write("$")
	if node.Long {
		p.write("{")
	}
	Walk(p, node.X)
	p.write(":")
	Walk(p, node.Y)
	if node.Long {
		p.write("}")
	}
	return nil
}

func (p *prettyPrinter) visitStaticMemberAccess(node *StaticMemberAccess) Visitor {
	p.write("[")
	Walk(p, node.X)
	p.write("]")
	p.write("::")
	Walk(p, node.Y)
	return nil
}

func (p *prettyPrinter) visitSelectExpr(node *SelectExpr) Visitor {
	Walk(p, node.X)
	p.write(".")
	Walk(p, node.Sel)
	return nil
}

func (p *prettyPrinter) visitConstrainedVarExpr(node *ConstrainedVarExpr) Visitor {
	p.write("[")
	Walk(p, node.Type)
	p.write("]")
	Walk(p, node.X)
	return nil
}

func (p *prettyPrinter) visitArrayExpr(node *ArrayExpr) Visitor {
	p.write("@(")
	for i, elem := range node.Elems {
		Walk(p, elem)
		if i < len(node.Elems)-1 {
			p.write(", ")
		}
	}
	p.write(")")
	return nil
}

func (p *prettyPrinter) visitArrayConstructor(node *ArrayConstructor) Visitor {
	for i, elem := range node.Elems {
		Walk(p, elem)
		if i < len(node.Elems)-1 {
			p.write(", ")
		}
	}
	return nil
}

func (p *prettyPrinter) visitDynamicparamStmt(node *DynamicparamStmt) Visitor {
	p.indentWrite("dynamicparam ")
	Walk(p, node.Body)
	return nil
}

func (p *prettyPrinter) visitDataStmt(node *DataStmt) Visitor {
	p.indentWrite("data ")
	for _, recv := range node.Recv {
		Walk(p, recv)
		p.write(" ")
	}
	Walk(p, node.Body)
	return nil
}

func (p *prettyPrinter) visitTrapStmt(node *TrapStmt) Visitor {
	p.indentWrite("trap ")
	Walk(p, node.Body)
	return nil
}

func (p *prettyPrinter) visitTryStmt(node *TryStmt) Visitor {
	p.indentWrite("try ")
	Walk(p, node.Body)

	for _, catch := range node.Catches {
		Walk(p, catch)
	}

	if node.FinallyBody != nil {
		p.indentWrite("finally ")
		Walk(p, node.FinallyBody)
	}
	return nil
}

func (p *prettyPrinter) visitCatchClause(node *CatchClause) Visitor {
	p.indentWrite("catch ")

	if node.Type != nil {
		p.write("[")
		Walk(p, node.Type)
		p.write("] ")
	}

	Walk(p, node.Body)

	return nil
}

func (p *prettyPrinter) visitSwitchStmt(node *SwitchStmt) Visitor {
	p.indentWrite("Switch ")
	if node.Pattern != SwitchPatternNone {
		p.write(node.Pattern.String())
		p.write(" ")
	}
	p.write("(")
	Walk(p, node.Value)
	p.write(") {\n")

	p.indent++
	for _, caseClause := range node.Cases {
		Walk(p, caseClause)
	}
	if node.Default != nil {
		p.indentWrite("default ")
		Walk(p, node.Default)
	}
	p.indent--
	p.indentWrite("}\n")
	return nil
}

func (p *prettyPrinter) visitCaseClause(node *CaseClause) Visitor {
	p.indentWrite("")
	Walk(p, node.Cond)
	p.write(" ")
	Walk(p, node.Body)
	return nil
}

func (p *prettyPrinter) visitLit(node *Lit) Visitor {
	p.write(node.Val)
	return nil
}

func (p *prettyPrinter) visitCmdExpr(node *CmdExpr) Visitor {
	Walk(p, node.Cmd)
	if len(node.Args) > 0 {
		p.write(" ")
		for i, arg := range node.Args {
			Walk(p, arg)
			if i < len(node.Args)-1 {
				p.write(" ")
			}
		}
	}
	return nil
}

func (p *prettyPrinter) visitCallExpr(node *CallExpr) Visitor {
	Walk(p, node.Func)
	p.write("(")
	for i, arg := range node.Recv {
		Walk(p, arg)
		if i < len(node.Recv)-1 {
			p.write(", ")
		}
	}
	p.write(")")
	return nil
}

func (p *prettyPrinter) visitExprStmt(node *ExprStmt) Visitor {
	p.indentWrite("")
	Walk(p, node.X)
	p.write("\n")
	return nil
}

func (p *prettyPrinter) visitProcessDecl(node *ProcessDecl) Visitor {
	p.indentWrite("Process ")
	Walk(p, node.Body)
	return nil
}

func (p *prettyPrinter) visitBlcokStmt(node *BlcokStmt) Visitor {
	p.write("{\n")
	p.indent++
	for _, stmt := range node.List {
		Walk(p, stmt)
	}
	p.indent--
	p.indentWrite("}\n")
	return nil
}

func (p *prettyPrinter) visitParameter(node *Parameter) Visitor {
	Walk(p, node.X)
	return nil
}

func (p *prettyPrinter) visitAttribute(node *Attribute) Visitor {
	p.write("[")
	Walk(p, node.Name)
	p.write("(")
	for i, recv := range node.Recv {
		Walk(p, recv)
		if i < len(node.Recv)-1 {
			p.write(", ")
		}
	}
	p.write(")")
	p.write("]")
	return nil
}

func (p *prettyPrinter) visitFuncDecl(node *FuncDecl) Visitor {
	p.indentWrite("Function ")
	Walk(p, node.Name)

	p.write(" {\n")

	for _, attribute := range node.Attributes {
		p.indentWrite(" ")
		Walk(p, attribute)
		p.write("\n")
	}

	p.indentWrite("Param(\n")
	for i, param := range node.Params {
		Walk(p, param)
		if i < len(node.Params)-1 {
			p.write(", ")
		}
	}
	p.write("\n")
	p.indentWrite(")\n")

	Walk(p, node.Body)
	p.write("}\n")
	return nil
}

func (p *prettyPrinter) visitCommentGroup(node *CommentGroup) Visitor {
	for _, comment := range node.List {
		Walk(p, comment)
		p.write("\n")
	}
	return nil
}

func (p *prettyPrinter) visitSingleLineComment(node *SingleLineComment) Visitor {
	p.write("#")
	p.write(node.Text)
	return nil
}

func (p *prettyPrinter) visitDelimitedComment(node *DelimitedComment) Visitor {
	p.write("<#")
	p.write(node.Text)
	p.write("#>")
	return nil
}

func (p *prettyPrinter) visitVarExpr(node *VarExpr) Visitor {
	p.write("$")
	Walk(p, node.X)
	return nil
}

func (p *prettyPrinter) visitIncDecExpr(node *IncDecExpr) Visitor {
	if node.Pre {
		p.write(node.Tok.String())
	}
	Walk(p, node.X)
	if !node.Pre {
		p.write(node.Tok.String())
	}
	return nil
}

func (p *prettyPrinter) visitIndexExpr(node *IndexExpr) Visitor {
	Walk(p, node.X)
	p.write("[")
	Walk(p, node.Index)
	p.write("]")
	return nil
}

func (p *prettyPrinter) visitIndicesExpr(node *IndicesExpr) Visitor {
	Walk(p, node.X)
	p.write("[")
	for i, index := range node.Indices {
		Walk(p, index)
		if i < len(node.Indices)-1 {
			p.write(", ")
		}
	}
	p.write("]")
	return nil
}

func (p *prettyPrinter) visitAssignStmt(node *AssignStmt) Visitor {
	Walk(p, node.Lhs)
	p.write(" = ")
	Walk(p, node.Rhs)
	return nil
}

func (p *prettyPrinter) visitFile(node *File) Visitor {
	for _, stmt := range node.Stmts {
		Walk(p, stmt)
	}
	return nil
}

func (p *prettyPrinter) visitHashTable(node *HashTable) Visitor {
	p.write("@{\n")
	for i, entry := range node.Entries {
		Walk(p, entry)
		if i < len(node.Entries)-1 {
			p.write("; ")
		}
	}
	p.write("}\n")
	return nil
}

func (p *prettyPrinter) visitHashEntry(node *HashEntry) Visitor {
	Walk(p, node.Key)
	p.write(" = ")
	Walk(p, node.Value)
	return nil
}

func Print(node Node) {
	p := prettyPrinter{output: os.Stdout, indentSpace: "  "}
	p.Visit(node)
}
