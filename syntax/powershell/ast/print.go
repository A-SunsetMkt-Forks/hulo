// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package ast

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/hulo-lang/hulo/syntax/powershell/token"
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
	if node == nil {
		return nil
	}

	switch node := node.(type) {
	case *File:
		return p.visitFile(node)
	case *BlockStmt:
		return p.visitBlockStmt(node)
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
	case *CommaExpr:
		return p.visitCommaExpr(node)
	case *BinaryExpr:
		return p.visitBinaryExpr(node)
	case *SelectExpr:
		return p.visitSelectExpr(node)
	case *CastExpr:
		return p.visitConstrainedVarExpr(node)
	case *Lit:
		return p.visitLit(node)
	case *StringLit:
		return p.visitStringLit(node)
	case *MultiStringLit:
		return p.visitMultiStringLit(node)
	case *NumericLit:
		return p.visitNumericLit(node)
	case *Ident:
		p.write(node.Name)
		return nil
	case *RangeExpr:
		return p.visitRangeExpr(node)
	case *GroupExpr:
		return p.visitGroupExpr(node)
	case *RedirectExpr:
		return p.visitRedirectExpr(node)
	case *BlockExpr:
		return p.visitBlockExpr(node)
	case *ParamBlock:
		return p.visitParamBlock(node)
	case *BoolLit:
		return p.visitBoolLit(node)
	case *SubExpr:
		return p.visitSubExpr(node)
	case *TypeLit:
		return p.visitTypeLit(node)
	case *ContinueStmt:
		return p.visitContinueStmt(node)
	case *BreakStmt:
		return p.visitBreakStmt(node)
	case *WhileStmt:
		return p.visitWhileStmt(node)
	case *ForStmt:
		return p.visitForStmt(node)
	case *ForeachStmt:
		return p.visitForeachStmt(node)
	case *DoWhileStmt:
		return p.visitDoWhileStmt(node)
	case *DoUntilStmt:
		return p.visitDoUntilStmt(node)
	case *LabelStmt:
		return p.visitLabelStmt(node)
	case *IfStmt:
		return p.visitIfStmt(node)
	case *ThrowStmt:
		return p.visitThrowStmt(node)
	case *ReturnStmt:
		return p.visitReturnStmt(node)
	case *WorkflowDecl:
		return p.visitWorkflowDecl(node)
	case *ParallelDecl:
		return p.visitParallelDecl(node)
	case *SequenceDecl:
		return p.visitSequenceDecl(node)
	case *InlinescriptDecl:
		return p.visitInlinescriptDecl(node)
	case *ClassDecl:
		return p.visitClassDecl(node)
	case *EnumDecl:
		return p.visitEnumDecl(node)
	case *EnumKeyValue:
		return p.visitEnumKeyValue(node)
	case *PropertyDecl:
		return p.visitPropertyDecl(node)
	case *MethodDecl:
		return p.visitMethodDecl(node)
	case *ConstructorDecl:
		return p.visitConstructorDecl(node)
	default:
		panic("unsupported node type: " + fmt.Sprintf("%T", node))
	}
}

func (p *prettyPrinter) visitConstructorDecl(node *ConstructorDecl) Visitor {
	p.indentWrite("")
	Walk(p, node.Name)
	p.write("(")
	p.visitExprs(node.Params, ", ")
	p.write(") ")
	Walk(p, node.Body)
	return nil
}

func (p *prettyPrinter) visitMethodDecl(node *MethodDecl) Visitor {
	p.indentWrite("")
	if node.Static {
		p.write("static ")
	}
	Walk(p, node.Type)
	p.write(" ")
	Walk(p, node.Name)
	p.write("(")
	p.visitExprs(node.Params, ", ")
	p.write(") ")
	Walk(p, node.Body)
	p.write("\n")
	return nil
}

func (p *prettyPrinter) visitPropertyDecl(node *PropertyDecl) Visitor {
	p.indentWrite("")
	if node.Mod != ModNone {
		p.write(node.Mod.String())
		p.write(" ")
	}
	Walk(p, node.Type)

	p.write(" $")
	Walk(p, node.Name)

	if node.Value != nil {
		p.write(" = ")
		Walk(p, node.Value)
	}

	p.write("\n")

	return nil
}

func (p *prettyPrinter) visitEnumKeyValue(node *EnumKeyValue) Visitor {
	p.indentWrite("")
	Walk(p, node.Label)
	if node.Value != nil {
		p.write(" = ")
		Walk(p, node.Value)
	}
	p.write("\n")
	return nil
}

func (p *prettyPrinter) visitEnumDecl(node *EnumDecl) Visitor {
	p.indentWrite("")
	for _, attr := range node.Attrs {
		Walk(p, attr)
	}

	p.write("enum ")
	Walk(p, node.Name)
	if node.UnderlyingType != nil {
		p.write(" : ")
		Walk(p, node.UnderlyingType)
	}
	p.write(" {\n")

	p.indent++

	for _, kv := range node.List {
		Walk(p, kv)
	}

	p.indent--

	p.indentWrite("}")

	return nil
}

func (p *prettyPrinter) visitClassDecl(node *ClassDecl) Visitor {
	p.indentWrite("class ")
	Walk(p, node.Name)

	p.write(" {\n")

	p.indent++
	for _, prop := range node.Properties {
		Walk(p, prop)
	}

	if len(node.Properties) > 0 {
		p.write("\n")
	}

	for _, ctor := range node.Ctors {
		Walk(p, ctor)
	}

	if len(node.Ctors) > 0 {
		p.write("\n")
	}

	for _, method := range node.Methods {
		Walk(p, method)
	}

	p.indent--

	p.write("}\n")

	return nil
}

func (p *prettyPrinter) visitThrowStmt(node *ThrowStmt) Visitor {
	p.indentWrite("throw ")
	Walk(p, node.X)
	p.write("\n")
	return nil
}

func (p *prettyPrinter) visitReturnStmt(node *ReturnStmt) Visitor {
	p.indentWrite("return ")
	Walk(p, node.X)
	p.write("\n")
	return nil
}

func (p *prettyPrinter) visitIfStmt(node *IfStmt) Visitor {
	p.indentWrite("if (")
	Walk(p, node.Cond)
	p.write(") ")

	Walk(p, node.Body)

	for node.Else != nil {
		switch el := node.Else.(type) {
		case *IfStmt:
			p.indentWrite("else if ")
			Walk(p, el.Cond)

			Walk(p, el.Body)

			node.Else = el.Else
		case *BlockStmt:
			p.indentWrite("else")
			Walk(p, el)
			node.Else = nil
		}
	}
	return nil
}

func (p *prettyPrinter) visitParallelDecl(node *ParallelDecl) Visitor {
	p.indentWrite("parallel ")
	Walk(p, node.Body)
	return nil
}

func (p *prettyPrinter) visitSequenceDecl(node *SequenceDecl) Visitor {
	p.indentWrite("sequence ")
	Walk(p, node.Body)
	return nil
}

func (p *prettyPrinter) visitInlinescriptDecl(node *InlinescriptDecl) Visitor {
	p.indentWrite("inlinescript ")
	Walk(p, node.Body)
	return nil
}

func (p *prettyPrinter) visitWorkflowDecl(node *WorkflowDecl) Visitor {
	p.indentWrite("workflow ")
	if node.Name != nil {
		Walk(p, node.Name)
		p.write(" ")
	}
	Walk(p, node.Body)
	return nil
}

func (p *prettyPrinter) visitForeachStmt(node *ForeachStmt) Visitor {
	p.indentWrite("foreach (")
	Walk(p, node.Elm)
	p.write(" in ")
	Walk(p, node.Elms)
	p.write(") ")
	Walk(p, node.Body)
	return nil
}

func (p *prettyPrinter) visitWhileStmt(node *WhileStmt) Visitor {
	p.indentWrite("while (")
	Walk(p, node.Cond)
	p.write(") ")
	Walk(p, node.Body)
	return nil
}

func (p *prettyPrinter) visitDoUntilStmt(node *DoUntilStmt) Visitor {
	p.indentWrite("do ")
	Walk(p, node.Body)
	p.write(" until (")
	Walk(p, node.Cond)
	p.write(")\n")
	return nil
}

func (p *prettyPrinter) visitDoWhileStmt(node *DoWhileStmt) Visitor {
	p.indentWrite("do ")
	Walk(p, node.Body)
	p.write(" while (")
	Walk(p, node.Cond)
	p.write(")\n")
	return nil
}

func (p *prettyPrinter) visitForStmt(node *ForStmt) Visitor {
	p.indentWrite("for (")
	Walk(p, node.Init)
	p.write("; ")
	Walk(p, node.Cond)
	p.write("; ")
	Walk(p, node.Post)
	p.write(") ")

	Walk(p, node.Body)
	return nil
}

func (p *prettyPrinter) visitLabelStmt(node *LabelStmt) Visitor {
	p.indentWrite(":")
	Walk(p, node.X)
	p.write("\n")
	return nil
}

func (p *prettyPrinter) visitBreakStmt(*BreakStmt) Visitor {
	p.indentWrite("break\n")
	return nil
}

func (p *prettyPrinter) visitContinueStmt(node *ContinueStmt) Visitor {
	p.indentWrite("continue")
	if node.Label != nil {
		p.write(" ")
		Walk(p, node.Label)
	}
	p.write("\n")
	return nil
}

func (p *prettyPrinter) visitTypeLit(node *TypeLit) Visitor {
	p.write("[")
	Walk(p, node.Name)
	p.write("]")
	return nil
}

func (p *prettyPrinter) visitSubExpr(node *SubExpr) Visitor {
	p.write("$( ")
	Walk(p, node.X)
	p.write(" )")
	return nil
}

func (p *prettyPrinter) visitBoolLit(node *BoolLit) Visitor {
	if node.Val {
		p.write("$true")
	} else {
		p.write("$false")
	}
	return nil
}

func (p *prettyPrinter) visitBlockExpr(node *BlockExpr) Visitor {
	p.write("{")
	if len(node.List) > 0 {
		p.write(" ")
	}
	p.visitExprs(node.List)
	if len(node.List) > 0 {
		p.write(" ")
	}
	p.write("}")
	return nil
}

func (p *prettyPrinter) visitRedirectExpr(node *RedirectExpr) Visitor {
	Walk(p, node.X)
	p.write(" ")
	p.write(node.CtrOp.String())
	p.write(" ")
	Walk(p, node.Y)
	return nil
}

func (p *prettyPrinter) visitGroupExpr(node *GroupExpr) Visitor {
	if node.Sep == token.Illegal {
		node.Sep = token.COMMA
	}
	sep := node.Sep.String() + " "
	p.write("(")
	p.visitExprs(node.Elems, sep)
	p.write(")")
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

func (p *prettyPrinter) visitRangeExpr(node *RangeExpr) Visitor {
	Walk(p, node.X)
	p.write("..")
	Walk(p, node.Y)
	return nil
}

func (p *prettyPrinter) visitStringLit(node *StringLit) Visitor {
	p.write("\"")
	p.write(node.Val)
	p.write("\"")
	return nil
}

func (p *prettyPrinter) visitMultiStringLit(node *MultiStringLit) Visitor {
	p.write("@\"")
	p.write(node.Val)
	p.write("\"@")
	return nil
}

func (p *prettyPrinter) visitNumericLit(node *NumericLit) Visitor {
	p.write(node.Val)
	return nil
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
	Walk(p, node.X)
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

func (p *prettyPrinter) visitConstrainedVarExpr(node *CastExpr) Visitor {
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

func (p *prettyPrinter) visitCommaExpr(node *CommaExpr) Visitor {
	p.visitExprs(node.Elems, ",")
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
		p.visitExprs(node.Args)
	}
	return nil
}

func (p *prettyPrinter) visitCallExpr(node *CallExpr) Visitor {
	Walk(p, node.Func)
	p.write("(")
	p.visitExprs(node.Recv, ", ")
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

func (p *prettyPrinter) visitBlockStmt(node *BlockStmt) Visitor {
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
	p.visitExprs(node.Recv, ", ")
	p.write(")")
	p.write("]")
	return nil
}

func (p *prettyPrinter) visitFuncDecl(node *FuncDecl) Visitor {
	p.indentWrite("Function ")
	Walk(p, node.Name)

	Walk(p, node.Body)

	return nil
}

func (p *prettyPrinter) visitParamBlock(node *ParamBlock) Visitor {
	for _, attr := range node.Attributes {
		p.indentWrite("")
		Walk(p, attr)
		p.write("\n")
	}

	p.indentWrite("Param(\n")

	p.indent++
	for i, param := range node.Params {
		for _, attr := range param.Attrs {
			p.indentWrite("")
			Walk(p, attr)
			p.write("\n")
		}
		p.indentWrite("")
		Walk(p, param)
		if i < len(node.Params)-1 {
			p.write("\n")
		}
	}
	p.indent--

	p.write("\n")
	p.indentWrite(")\n")
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
	p.visitExprs(node.Indices, ", ")
	p.write("]")
	return nil
}

func (p *prettyPrinter) visitAssignStmt(node *AssignStmt) Visitor {
	p.indentWrite("")
	if node.Tok == token.Illegal {
		node.Tok = token.ASSIGN
	}
	Walk(p, node.Lhs)
	p.write(" ")
	p.write(node.Tok.String())
	p.write(" ")
	Walk(p, node.Rhs)
	p.write("\n")
	return nil
}

func (p *prettyPrinter) visitFile(node *File) Visitor {
	for _, doc := range node.Docs {
		Walk(p, doc)
		p.write("\n")
	}

	for _, stmt := range node.Stmts {
		Walk(p, stmt)
	}
	return nil
}

func (p *prettyPrinter) visitHashTable(node *HashTable) Visitor {
	p.write("@{")
	for i, entry := range node.Entries {
		Walk(p, entry)
		if i < len(node.Entries)-1 {
			p.write("; ")
		}
	}
	p.write("}")
	if len(node.Entries) > 0 {
		p.write("\n")
	}
	return nil
}

func (p *prettyPrinter) visitHashEntry(node *HashEntry) Visitor {
	Walk(p, node.Key)
	p.write(" = ")
	Walk(p, node.Value)
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

func Write(node Node, output io.Writer) {
	p := prettyPrinter{output: output, indentSpace: "  "}
	p.Visit(node)
}

func Print(node Node) {
	Write(node, os.Stdout)
}

func String(node Node) string {
	var buf bytes.Buffer
	Write(node, &buf)
	return buf.String()
}
