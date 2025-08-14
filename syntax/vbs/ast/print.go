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

// Inspect prints the AST structure for debugging purposes
func Inspect(node Node, output io.Writer) {
	printer := &debugPrinter{output: output, indent: 0}
	printer.visit(node)
}

type debugPrinter struct {
	output io.Writer
	indent int
}

func (p *debugPrinter) write(s string) {
	fmt.Fprint(p.output, s)
}

func (p *debugPrinter) indentWrite(s string) {
	fmt.Fprint(p.output, strings.Repeat("  ", p.indent))
	fmt.Fprint(p.output, s)
}

func (p *debugPrinter) visit(node Node) {
	if node == nil {
		p.write("nil")
		return
	}

	switch n := node.(type) {
	case *File:
		p.write("*ast.File {\n")
		p.indent++
		if len(n.Stmts) > 0 {
			p.indentWrite("Stmts: [\n")
			p.indent++
			for i, stmt := range n.Stmts {
				p.indentWrite(fmt.Sprintf("%d: ", i))
				p.visit(stmt)
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
		p.indentWrite("List: [\n")
		p.indent++
		for i, comment := range n.List {
			p.indentWrite(fmt.Sprintf("%d: ", i))
			p.visit(comment)
			p.write("\n")
		}
		p.indent--
		p.indentWrite("]\n")
		p.indent--
		p.indentWrite("}")
	case *Comment:
		p.write(fmt.Sprintf("*ast.Comment (Text: %q)", n.Text))
	case *DimDecl:
		p.write("*ast.DimDecl {\n")
		p.indent++
		if len(n.List) > 0 {
			p.indentWrite("List: [\n")
			p.indent++
			for i, item := range n.List {
				p.indentWrite(fmt.Sprintf("%d: ", i))
				p.visit(item)
				p.write("\n")
			}
			p.indent--
			p.indentWrite("]\n")
		}
		if n.Colon.IsValid() {
			p.indentWrite("Colon: true\n")
		}
		if n.Set != nil {
			p.indentWrite("Set: ")
			p.visit(n.Set)
			p.write("\n")
		}
		p.indent--
		p.indentWrite("}")
	case *ReDimDecl:
		p.write("*ast.ReDimDecl {\n")
		p.indent++
		if n.Preserve.IsValid() {
			p.indentWrite("Preserve: true\n")
		}
		if len(n.List) > 0 {
			p.indentWrite("List: [\n")
			p.indent++
			for i, item := range n.List {
				p.indentWrite(fmt.Sprintf("%d: ", i))
				p.visit(item)
				p.write("\n")
			}
			p.indent--
			p.indentWrite("]\n")
		}
		p.indent--
		p.indentWrite("}")
	case *ClassDecl:
		p.write("*ast.ClassDecl {\n")
		p.indent++
		if n.Mod.IsValid() {
			p.indentWrite(fmt.Sprintf("Mod: %s\n", n.Mod))
		}
		if n.Name != nil {
			p.indentWrite("Name: ")
			p.visit(n.Name)
			p.write("\n")
		}
		if len(n.Stmts) > 0 {
			p.indentWrite("Stmts: [\n")
			p.indent++
			for i, stmt := range n.Stmts {
				p.indentWrite(fmt.Sprintf("%d: ", i))
				p.visit(stmt)
				p.write("\n")
			}
			p.indent--
			p.indentWrite("]\n")
		}
		p.indent--
		p.indentWrite("}")
	case *SubDecl:
		p.write("*ast.SubDecl {\n")
		p.indent++
		if n.Mod.IsValid() {
			p.indentWrite(fmt.Sprintf("Mod: %s\n", n.Mod))
		}
		if n.Name != nil {
			p.indentWrite("Name: ")
			p.visit(n.Name)
			p.write("\n")
		}
		if n.Body != nil {
			p.indentWrite("Body: ")
			p.visit(n.Body)
			p.write("\n")
		}
		p.indent--
		p.indentWrite("}")
	case *FuncDecl:
		p.write("*ast.FuncDecl {\n")
		p.indent++
		if n.Mod.IsValid() {
			p.indentWrite(fmt.Sprintf("Mod: %s\n", n.Mod))
		}
		if n.Default.IsValid() {
			p.indentWrite("Default: true\n")
		}
		if n.Name != nil {
			p.indentWrite("Name: ")
			p.visit(n.Name)
			p.write("\n")
		}
		if len(n.Recv) > 0 {
			p.indentWrite("Recv: [\n")
			p.indent++
			for i, param := range n.Recv {
				p.indentWrite(fmt.Sprintf("%d: ", i))
				p.visit(param)
				p.write("\n")
			}
			p.indent--
			p.indentWrite("]\n")
		}
		if n.Body != nil {
			p.indentWrite("Body: ")
			p.visit(n.Body)
			p.write("\n")
		}
		p.indent--
		p.indentWrite("}")
	case *PropertySetStmt:
		p.write("*ast.PropertySetStmt {\n")
		p.indent++
		if n.Mod.IsValid() {
			p.indentWrite(fmt.Sprintf("Mod: %s\n", n.Mod))
		}
		if n.Name != nil {
			p.indentWrite("Name: ")
			p.visit(n.Name)
			p.write("\n")
		}
		if len(n.Recv) > 0 {
			p.indentWrite("Recv: [\n")
			p.indent++
			for i, param := range n.Recv {
				p.indentWrite(fmt.Sprintf("%d: ", i))
				p.visit(param)
				p.write("\n")
			}
			p.indent--
			p.indentWrite("]\n")
		}
		if n.Body != nil {
			p.indentWrite("Body: ")
			p.visit(n.Body)
			p.write("\n")
		}
		p.indent--
		p.indentWrite("}")
	case *PropertyLetStmt:
		p.write("*ast.PropertyLetStmt {\n")
		p.indent++
		if n.Mod.IsValid() {
			p.indentWrite(fmt.Sprintf("Mod: %s\n", n.Mod))
		}
		if n.Name != nil {
			p.indentWrite("Name: ")
			p.visit(n.Name)
			p.write("\n")
		}
		if len(n.Recv) > 0 {
			p.indentWrite("Recv: [\n")
			p.indent++
			for i, param := range n.Recv {
				p.indentWrite(fmt.Sprintf("%d: ", i))
				p.visit(param)
				p.write("\n")
			}
			p.indent--
			p.indentWrite("]\n")
		}
		if n.Body != nil {
			p.indentWrite("Body: ")
			p.visit(n.Body)
			p.write("\n")
		}
		p.indent--
		p.indentWrite("}")
	case *PropertyGetStmt:
		p.write("*ast.PropertyGetStmt {\n")
		p.indent++
		if n.Mod.IsValid() {
			p.indentWrite(fmt.Sprintf("Mod: %s\n", n.Mod))
		}
		if n.Name != nil {
			p.indentWrite("Name: ")
			p.visit(n.Name)
			p.write("\n")
		}
		if len(n.Recv) > 0 {
			p.indentWrite("Recv: [\n")
			p.indent++
			for i, param := range n.Recv {
				p.indentWrite(fmt.Sprintf("%d: ", i))
				p.visit(param)
				p.write("\n")
			}
			p.indent--
			p.indentWrite("]\n")
		}
		if n.Body != nil {
			p.indentWrite("Body: ")
			p.visit(n.Body)
			p.write("\n")
		}
		p.indent--
		p.indentWrite("}")
	case *IfStmt:
		p.write("*ast.IfStmt {\n")
		p.indent++
		if n.Cond != nil {
			p.indentWrite("Cond: ")
			p.visit(n.Cond)
			p.write("\n")
		}
		if n.Body != nil {
			p.indentWrite("Body: ")
			p.visit(n.Body)
			p.write("\n")
		}
		if n.Else != nil {
			p.indentWrite("Else: ")
			p.visit(n.Else)
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
				p.visit(stmt)
				p.write("\n")
			}
			p.indent--
			p.indentWrite("]\n")
		}
		p.indent--
		p.indentWrite("}")
	case *ExprStmt:
		p.write("*ast.ExprStmt {\n")
		p.indent++
		if n.X != nil {
			p.indentWrite("X: ")
			p.visit(n.X)
			p.write("\n")
		}
		p.indent--
		p.indentWrite("}")
	case *PrivateStmt:
		p.write("*ast.PrivateStmt {\n")
		p.indent++
		if len(n.List) > 0 {
			p.indentWrite("List: [\n")
			p.indent++
			for i, item := range n.List {
				p.indentWrite(fmt.Sprintf("%d: ", i))
				p.visit(item)
				p.write("\n")
			}
			p.indent--
			p.indentWrite("]\n")
		}
		p.indent--
		p.indentWrite("}")
	case *PublicStmt:
		p.write("*ast.PublicStmt {\n")
		p.indent++
		if len(n.List) > 0 {
			p.indentWrite("List: [\n")
			p.indent++
			for i, item := range n.List {
				p.indentWrite(fmt.Sprintf("%d: ", i))
				p.visit(item)
				p.write("\n")
			}
			p.indent--
			p.indentWrite("]\n")
		}
		p.indent--
		p.indentWrite("}")
	case *ConstStmt:
		p.write("*ast.ConstStmt {\n")
		p.indent++
		if n.Lhs != nil {
			p.indentWrite("Lhs: ")
			p.visit(n.Lhs)
			p.write("\n")
		}
		if n.Rhs != nil {
			p.indentWrite("Rhs: ")
			p.visit(n.Rhs)
			p.write("\n")
		}
		p.indent--
		p.indentWrite("}")
	case *SetStmt:
		p.write("*ast.SetStmt {\n")
		p.indent++
		if n.Lhs != nil {
			p.indentWrite("Lhs: ")
			p.visit(n.Lhs)
			p.write("\n")
		}
		if n.Rhs != nil {
			p.indentWrite("Rhs: ")
			p.visit(n.Rhs)
			p.write("\n")
		}
		p.indent--
		p.indentWrite("}")
	case *AssignStmt:
		p.write("*ast.AssignStmt {\n")
		p.indent++
		if n.Lhs != nil {
			p.indentWrite("Lhs: ")
			p.visit(n.Lhs)
			p.write("\n")
		}
		if n.Rhs != nil {
			p.indentWrite("Rhs: ")
			p.visit(n.Rhs)
			p.write("\n")
		}
		p.indent--
		p.indentWrite("}")
	case *ForEachStmt:
		p.write("*ast.ForEachStmt {\n")
		p.indent++
		if n.Elem != nil {
			p.indentWrite("Elem: ")
			p.visit(n.Elem)
			p.write("\n")
		}
		if n.Group != nil {
			p.indentWrite("Group: ")
			p.visit(n.Group)
			p.write("\n")
		}
		if n.Body != nil {
			p.indentWrite("Body: ")
			p.visit(n.Body)
			p.write("\n")
		}
		if n.Stmt != nil {
			p.indentWrite("Stmt: ")
			p.visit(n.Stmt)
			p.write("\n")
		}
		p.indent--
		p.indentWrite("}")
	case *ForNextStmt:
		p.write("*ast.ForNextStmt {\n")
		p.indent++
		if n.Start != nil {
			p.indentWrite("Start: ")
			p.visit(n.Start)
			p.write("\n")
		}
		if n.End_ != nil {
			p.indentWrite("End: ")
			p.visit(n.End_)
			p.write("\n")
		}
		if n.Step != nil {
			p.indentWrite("Step: ")
			p.visit(n.Step)
			p.write("\n")
		}
		if n.Body != nil {
			p.indentWrite("Body: ")
			p.visit(n.Body)
			p.write("\n")
		}
		p.indent--
		p.indentWrite("}")
	case *WhileWendStmt:
		p.write("*ast.WhileWendStmt {\n")
		p.indent++
		if n.Cond != nil {
			p.indentWrite("Cond: ")
			p.visit(n.Cond)
			p.write("\n")
		}
		if n.Body != nil {
			p.indentWrite("Body: ")
			p.visit(n.Body)
			p.write("\n")
		}
		p.indent--
		p.indentWrite("}")
	case *DoLoopStmt:
		p.write("*ast.DoLoopStmt {\n")
		p.indent++
		p.indentWrite(fmt.Sprintf("Pre: %t\n", n.Pre))
		if n.Cond != nil {
			p.indentWrite("Cond: ")
			p.visit(n.Cond)
			p.write("\n")
		}
		if n.Tok != token.ILLEGAL {
			p.indentWrite(fmt.Sprintf("Tok: %s\n", n.Tok))
		}
		if n.Body != nil {
			p.indentWrite("Body: ")
			p.visit(n.Body)
			p.write("\n")
		}
		p.indent--
		p.indentWrite("}")
	case *CallStmt:
		p.write("*ast.CallStmt {\n")
		p.indent++
		if n.Name != nil {
			p.indentWrite("Name: ")
			p.visit(n.Name)
			p.write("\n")
		}
		if len(n.Recv) > 0 {
			p.indentWrite("Recv: [\n")
			p.indent++
			for i, arg := range n.Recv {
				p.indentWrite(fmt.Sprintf("%d: ", i))
				p.visit(arg)
				p.write("\n")
			}
			p.indent--
			p.indentWrite("]\n")
		}
		p.indent--
		p.indentWrite("}")
	case *ExitStmt:
		p.write("*ast.ExitStmt {\n")
		p.indent++
		if n.Tok != token.ILLEGAL {
			p.indentWrite(fmt.Sprintf("Tok: %s\n", n.Tok))
		}
		p.indent--
		p.indentWrite("}")
	case *SelectStmt:
		p.write("*ast.SelectStmt {\n")
		p.indent++
		if n.Var != nil {
			p.indentWrite("Var: ")
			p.visit(n.Var)
			p.write("\n")
		}
		if len(n.Cases) > 0 {
			p.indentWrite("Cases: [\n")
			p.indent++
			for i, case_ := range n.Cases {
				p.indentWrite(fmt.Sprintf("%d: *ast.CaseClause {\n", i))
				p.indent++
				if case_.Cond != nil {
					p.indentWrite("Cond: ")
					p.visit(case_.Cond)
					p.write("\n")
				}
				if case_.Body != nil {
					p.indentWrite("Body: ")
					p.visit(case_.Body)
					p.write("\n")
				}
				p.indent--
				p.indentWrite("}\n")
			}
			p.indent--
			p.indentWrite("]\n")
		}
		if n.Else != nil {
			p.indentWrite("Else: ")
			p.visit(n.Else.Body)
			p.write("\n")
		}
		p.indent--
		p.indentWrite("}")
	case *WithStmt:
		p.write("*ast.WithStmt {\n")
		p.indent++
		if n.Cond != nil {
			p.indentWrite("Cond: ")
			p.visit(n.Cond)
			p.write("\n")
		}
		if n.Body != nil {
			p.indentWrite("Body: ")
			p.visit(n.Body)
			p.write("\n")
		}
		p.indent--
		p.indentWrite("}")
	case *OnErrorStmt:
		p.write("*ast.OnErrorStmt {\n")
		p.indent++
		if n.OnErrorGoto != nil {
			p.indentWrite("OnErrorGoto: 0")
			p.write("\n")
		}
		if n.OnErrorResume != nil {
			p.indentWrite("OnErrorResume: Next")
			p.write("\n")
		}
		p.indent--
		p.indentWrite("}")
	case *StopStmt:
		p.write("*ast.StopStmt")
	case *RandomizeStmt:
		p.write("*ast.RandomizeStmt")
	case *OptionStmt:
		p.write("*ast.OptionStmt")
	case *EraseStmt:
		p.write("*ast.EraseStmt {\n")
		p.indent++
		if n.X != nil {
			p.indentWrite("X: ")
			p.visit(n.X)
			p.write("\n")
		}
		p.indent--
		p.indentWrite("}")
	case *ExecuteStmt:
		p.write(fmt.Sprintf("*ast.ExecuteStmt (Statement: %q)", n.Statement))
	case *Ident:
		p.write(fmt.Sprintf("*ast.Ident (Name: %q)", n.Name))
	case *BasicLit:
		p.write(fmt.Sprintf("*ast.BasicLit (Kind: %s, Value: %q)", n.Kind, n.Value))
	case *SelectorExpr:
		p.write("*ast.SelectorExpr {\n")
		p.indent++
		if n.X != nil {
			p.indentWrite("X: ")
			p.visit(n.X)
			p.write("\n")
		}
		if n.Sel != nil {
			p.indentWrite("Sel: ")
			p.visit(n.Sel)
			p.write("\n")
		}
		p.indent--
		p.indentWrite("}")
	case *BinaryExpr:
		p.write("*ast.BinaryExpr {\n")
		p.indent++
		if n.X != nil {
			p.indentWrite("X: ")
			p.visit(n.X)
			p.write("\n")
		}
		p.indentWrite(fmt.Sprintf("Op: %s\n", n.Op))
		if n.Y != nil {
			p.indentWrite("Y: ")
			p.visit(n.Y)
			p.write("\n")
		}
		p.indent--
		p.indentWrite("}")
	case *CallExpr:
		p.write("*ast.CallExpr {\n")
		p.indent++
		if n.Func != nil {
			p.indentWrite("Func: ")
			p.visit(n.Func)
			p.write("\n")
		}
		if len(n.Recv) > 0 {
			p.indentWrite("Recv: [\n")
			p.indent++
			for i, arg := range n.Recv {
				p.indentWrite(fmt.Sprintf("%d: ", i))
				p.visit(arg)
				p.write("\n")
			}
			p.indent--
			p.indentWrite("]\n")
		}
		p.indent--
		p.indentWrite("}")
	case *CmdExpr:
		p.write("*ast.CmdExpr {\n")
		p.indent++
		if n.Cmd != nil {
			p.indentWrite("Cmd: ")
			p.visit(n.Cmd)
			p.write("\n")
		}
		if len(n.Recv) > 0 {
			p.indentWrite("Recv: [\n")
			p.indent++
			for i, arg := range n.Recv {
				p.indentWrite(fmt.Sprintf("%d: ", i))
				p.visit(arg)
				p.write("\n")
			}
			p.indent--
			p.indentWrite("]\n")
		}
		p.indent--
		p.indentWrite("}")
	case *IndexExpr:
		p.write("*ast.IndexExpr {\n")
		p.indent++
		if n.X != nil {
			p.indentWrite("X: ")
			p.visit(n.X)
			p.write("\n")
		}
		if n.Index != nil {
			p.indentWrite("Index: ")
			p.visit(n.Index)
			p.write("\n")
		}
		p.indent--
		p.indentWrite("}")
	case *IndexListExpr:
		p.write("*ast.IndexListExpr {\n")
		p.indent++
		if n.X != nil {
			p.indentWrite("X: ")
			p.visit(n.X)
			p.write("\n")
		}
		if len(n.Indices) > 0 {
			p.indentWrite("Indices: [\n")
			p.indent++
			for i, index := range n.Indices {
				p.indentWrite(fmt.Sprintf("%d: ", i))
				p.visit(index)
				p.write("\n")
			}
			p.indent--
			p.indentWrite("]\n")
		}
		p.indent--
		p.indentWrite("}")
	case *NewExpr:
		p.write("*ast.NewExpr {\n")
		p.indent++
		if n.X != nil {
			p.indentWrite("X: ")
			p.visit(n.X)
			p.write("\n")
		}
		p.indent--
		p.indentWrite("}")
	case *Field:
		p.write("*ast.Field {\n")
		p.indent++
		if n.Tok.IsValid() {
			p.indentWrite(fmt.Sprintf("Tok: %s\n", n.Tok))
		}
		if n.Name != nil {
			p.indentWrite("Name: ")
			p.visit(n.Name)
			p.write("\n")
		}
		p.indent--
		p.indentWrite("}")
	default:
		p.write(fmt.Sprintf("unknown node type: %T", node))
	}
}

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
	switch n := node.(type) {
	case *File:
		return p.visitFile(n)
	case *CommentGroup:
		return p.visitCommentGroup(n)
	case *Comment:
		return p.visitComment(n)
	case *DimDecl:
		return p.visitDimDecl(n)
	case *ReDimDecl:
		return p.visitReDimDecl(n)
	case *ClassDecl:
		return p.visitClassDecl(n)
	case *SubDecl:
		return p.visitSubDecl(n)
	case *FuncDecl:
		return p.visitFuncDecl(n)
	case *PropertyGetStmt:
		return p.visitPropertyGetStmt(n)
	case *PropertySetStmt:
		return p.visitPropertySetStmt(n)
	case *PropertyLetStmt:
		return p.visitPropertyLetStmt(n)
	case *IfStmt:
		return p.visitIfStmt(n)
	case *BlockStmt:
		return p.visitBlockStmt(n)
	case *ExprStmt:
		return p.visitExprStmt(n)
	case *PrivateStmt:
		return p.visitPrivateStmt(n)
	case *PublicStmt:
		return p.visitPublicStmt(n)
	case *ConstStmt:
		return p.visitConstStmt(n)
	case *SetStmt:
		return p.visitSetStmt(n)
	case *AssignStmt:
		return p.visitAssignStmt(n)
	case *ForEachStmt:
		return p.visitForEachStmt(n)
	case *ForNextStmt:
		return p.visitForNextStmt(n)
	case *WhileWendStmt:
		return p.visitWhileWendStmt(n)
	case *DoLoopStmt:
		return p.visitDoLoopStmt(n)
	case *CallStmt:
		return p.visitCallStmt(n)
	case *ExitStmt:
		return p.visitExitStmt(n)
	case *SelectStmt:
		return p.visitSelectStmt(n)
	case *WithStmt:
		return p.visitWithStmt(n)
	case *OnErrorStmt:
		return p.visitOnErrorStmt(n)
	case *StopStmt:
		return p.visitStopStmt(n)
	case *RandomizeStmt:
		return p.visitRandomizeStmt(n)
	case *OptionStmt:
		return p.visitOptionStmt(n)
	case *EraseStmt:
		return p.visitEraseStmt(n)
	case *ExecuteStmt:
		return p.visitExecuteStmt(n)
	case *Ident:
		return p.visitIdent(n)
	case *BasicLit:
		return p.visitBasicLit(n)
	case *SelectorExpr:
		return p.visitSelectorExpr(n)
	case *BinaryExpr:
		return p.visitBinaryExpr(n)
	case *CallExpr:
		return p.visitCallExpr(n)
	case *CmdExpr:
		return p.visitCmdExpr(n)
	case *IndexExpr:
		return p.visitIndexExpr(n)
	case *IndexListExpr:
		return p.visitIndexListExpr(n)
	case *NewExpr:
		return p.visitNewExpr(n)
	case *Field:
		return p.visitField(n)
	}
	return nil
}

func (p *prettyPrinter) visitNewExpr(n *NewExpr) Visitor {
	p.write("New ")
	Walk(p, n.X)
	return nil
}

func (p *prettyPrinter) visitField(n *Field) Visitor {
	if n.Tok.IsValid() {
		p.write(n.Tok.String())
		p.write(" ")
	}
	Walk(p, n.Name)
	return nil
}

func (p *prettyPrinter) visitIndexListExpr(n *IndexListExpr) Visitor {
	Walk(p, n.X)
	p.write("(")
	for i, r := range n.Indices {
		if i > 0 {
			p.write(", ")
		}
		Walk(p, r)
	}
	p.write(")")
	return nil
}

func (p *prettyPrinter) visitIndexExpr(n *IndexExpr) Visitor {
	Walk(p, n.X)
	p.write("(")
	Walk(p, n.Index)
	p.write(")")
	return nil
}

func (p *prettyPrinter) visitCmdExpr(n *CmdExpr) Visitor {
	Walk(p, n.Cmd)
	p.write(" ")
	for i, r := range n.Recv {
		if i > 0 {
			p.write(", ")
		}
		Walk(p, r)
	}
	return nil
}

func (p *prettyPrinter) visitCallExpr(n *CallExpr) Visitor {
	Walk(p, n.Func)
	p.write("(")
	for i, r := range n.Recv {
		if i > 0 {
			p.write(", ")
		}
		Walk(p, r)
	}
	p.write(")")
	return nil
}

func (p *prettyPrinter) visitIdent(n *Ident) Visitor {
	p.write(n.Name)
	return nil
}

func (p *prettyPrinter) visitBasicLit(n *BasicLit) Visitor {
	switch n.Kind {
	case token.TRUE:
		p.write("True")
	case token.FALSE:
		p.write("False")
	case token.STRING:
		p.write(fmt.Sprintf(`"%s"`, n.Value))
	default:
		p.write(n.Value)
	}
	return nil
}

func (p *prettyPrinter) visitSelectorExpr(n *SelectorExpr) Visitor {
	Walk(p, n.X)
	p.write(".")
	Walk(p, n.Sel)
	return nil
}

func (p *prettyPrinter) visitBinaryExpr(n *BinaryExpr) Visitor {
	Walk(p, n.X)
	p.write(" ")
	p.write(n.Op.String())
	p.write(" ")
	Walk(p, n.Y)
	return nil
}

func (p *prettyPrinter) visitForNextStmt(n *ForNextStmt) Visitor {
	p.indentWrite("For ")
	Walk(p, n.Start)
	p.write(" To ")
	Walk(p, n.End_)
	if n.Step != nil {
		p.write(" Step ")
		Walk(p, n.Step)
		p.write("\n")
	}
	Walk(p, n.Body)
	p.indentWrite("Next")
	p.write("\n")
	return nil
}

func (p *prettyPrinter) visitForEachStmt(n *ForEachStmt) Visitor {
	p.indentWrite("For Each ")
	Walk(p, n.Elem)
	p.write(" In ")
	Walk(p, n.Group)
	p.write("\n")
	Walk(p, n.Body)
	p.indentWrite("Next")
	if n.Stmt != nil {
		p.write(" ")
		Walk(p, n)
	} else {
		p.write("\n")
	}
	return nil
}

func (p *prettyPrinter) visitAssignStmt(n *AssignStmt) Visitor {
	p.indentWrite("")
	Walk(p, n.Lhs)
	p.write(" = ")
	Walk(p, n.Rhs)
	p.write("\n")
	return nil
}

func (p *prettyPrinter) visitSetStmt(n *SetStmt) Visitor {
	p.indentWrite("Set ")
	Walk(p, n.Lhs)
	p.write(" = ")
	Walk(p, n.Rhs)
	p.write("\n")
	return nil
}

func (p *prettyPrinter) visitPrivateStmt(n *PrivateStmt) Visitor {
	p.indentWrite("Private ")
	p.visitExprList(n.List)
	p.write("\n")
	return nil
}

func (p *prettyPrinter) visitPublicStmt(n *PublicStmt) Visitor {
	p.indentWrite("Public ")
	p.visitExprList(n.List)
	p.write("\n")
	return nil
}

func (p *prettyPrinter) visitConstStmt(n *ConstStmt) Visitor {
	p.indentWrite("Const ")
	Walk(p, n.Lhs)
	p.write(" = ")
	Walk(p, n.Rhs)
	p.write("\n")
	return nil
}

func (p *prettyPrinter) visitWhileWendStmt(n *WhileWendStmt) Visitor {
	p.indentWrite("While ")
	Walk(p, n.Cond)
	p.write("\n")
	Walk(p, n.Body)
	p.indentWrite("Wend")
	p.write("\n")
	return nil
}

func (p *prettyPrinter) visitDoLoopStmt(n *DoLoopStmt) Visitor {
	if n.Pre {
		p.indentWrite("Do ")
		if n.Cond != nil {
			p.write(n.Tok.String())
			p.write(" ")
			Walk(p, n.Cond)
		}
		p.write("\n")
		Walk(p, n.Body)
		p.indentWrite("Loop")
		p.write("\n")
	} else {
		p.indentWrite("Do\n")
		for _, s := range n.Body.List {
			temp := p.indentSpace
			p.indentSpace += "  "
			Walk(p, s)
			p.indentSpace = temp
		}
		p.indentWrite("Loop ")
		if n.Cond != nil {
			p.write(n.Tok.String())
			p.write(" ")
			Walk(p, n.Cond)
		}
		p.write("\n")
	}
	return nil
}

func (p *prettyPrinter) visitCallStmt(n *CallStmt) Visitor {
	p.indentWrite("Call ")
	Walk(p, n.Name)
	p.write(" ")
	p.visitExprList(n.Recv)
	p.write("\n")
	return nil
}

func (p *prettyPrinter) visitExitStmt(n *ExitStmt) Visitor {
	p.indentWrite("Exit")
	if n.Tok != token.ILLEGAL {
		p.write(" ")
		p.write(n.Tok.String())
	}
	p.write("\n")
	return nil
}

func (p *prettyPrinter) visitSelectStmt(n *SelectStmt) Visitor {
	p.indentWrite("Select Case ")
	Walk(p, n.Var)
	p.write("\n")
	for _, c := range n.Cases {
		p.indentWrite("  Case ")
		Walk(p, c.Cond)
		p.write("\n")
		Walk(p, c.Body)
	}
	if n.Else != nil {
		p.indentWrite("  Case Else")
		Walk(p, n.Else.Body)
	}
	p.indentWrite("End Select")
	p.write("\n")
	return nil
}

func (p *prettyPrinter) visitWithStmt(n *WithStmt) Visitor {
	p.indentWrite("With ")
	Walk(p, n.Cond)
	p.write("\n")
	Walk(p, n.Body)
	p.indentWrite("End With")
	p.write("\n")
	return nil
}

func (p *prettyPrinter) visitOnErrorStmt(n *OnErrorStmt) Visitor {
	if n.OnErrorGoto != nil {
		p.indentWrite("On Error Goto 0")
	}
	if n.OnErrorResume != nil {
		p.indentWrite("On Error Resume Next")
	}
	p.write("\n")
	return nil
}

func (p *prettyPrinter) visitStopStmt(_ *StopStmt) Visitor {
	p.indentWrite("Stop")
	p.write("\n")
	return nil
}

func (p *prettyPrinter) visitRandomizeStmt(_ *RandomizeStmt) Visitor {
	p.indentWrite("Randomize")
	p.write("\n")
	return nil
}

func (p *prettyPrinter) visitOptionStmt(_ *OptionStmt) Visitor {
	p.indentWrite("Option Explicit")
	p.write("\n")
	return nil
}

func (p *prettyPrinter) visitEraseStmt(n *EraseStmt) Visitor {
	p.indentWrite("Erase ")
	Walk(p, n.X)
	p.write("\n")
	return nil
}

func (p *prettyPrinter) visitExecuteStmt(n *ExecuteStmt) Visitor {
	p.indentWrite("Execute ")
	p.write(fmt.Sprintf(`"%s"`, n.Statement))
	p.write("\n")
	return nil
}

func (p *prettyPrinter) visitExprStmt(n *ExprStmt) Visitor {
	p.indentWrite("")
	Walk(p, n.X)
	p.write("\n")
	return nil
}

func (p *prettyPrinter) visitBlockStmt(n *BlockStmt) Visitor {
	p.indent++
	for _, s := range n.List {
		Walk(p, s)
	}
	p.indent--
	return nil
}

func (p *prettyPrinter) visitIfStmt(n *IfStmt) Visitor {
	p.indentWrite("If ")
	Walk(p, n.Cond)
	p.write(" Then\n")

	for _, s := range n.Body.List {
		temp := p.indentSpace
		p.indentSpace += "  "
		Walk(p, s)
		p.indentSpace = temp
	}

	for n.Else != nil {
		switch el := n.Else.(type) {
		case *IfStmt:
			p.indentWrite("ElseIf ")
			Walk(p, el.Cond)
			p.write(" Then\n")
			Walk(p, el.Body)
			n.Else = el.Else
		case *BlockStmt:
			p.indentWrite("Else")
			p.write("\n")
			Walk(p, el)
			n.Else = nil
		}
	}

	p.indentWrite("End If")
	p.write("\n")
	return nil
}

func (p *prettyPrinter) visitPropertyGetStmt(n *PropertyGetStmt) Visitor {
	p.indentWrite("")
	if n.Mod.IsValid() {
		p.write(n.Mod.String())
		p.write(" ")
	}
	list := []string{}
	for _, r := range n.Recv {
		if r.TokPos.IsValid() {
			list = append(list, fmt.Sprintf("%s %s", r.Tok, r.Name.Name))
		} else {
			list = append(list, r.Name.Name)
		}
	}
	p.indentWrite("Property Get ")
	Walk(p, n.Name)
	p.write("(")
	p.write(strings.Join(list, ", "))
	p.write(")\n")
	Walk(p, n.Body)

	p.indentWrite("End Property")
	p.write("\n")
	return nil
}

func (p *prettyPrinter) visitPropertySetStmt(n *PropertySetStmt) Visitor {
	p.indentWrite("")
	if n.Mod.IsValid() {
		p.write(n.Mod.String())
		p.write(" ")
	}
	list := []string{}
	for _, r := range n.Recv {
		if r.TokPos.IsValid() {
			list = append(list, fmt.Sprintf("%s %s", r.Tok, r.Name.Name))
		} else {
			list = append(list, r.Name.Name)
		}
	}
	p.indentWrite("Property Set ")
	Walk(p, n.Name)
	p.write("(")
	p.write(strings.Join(list, ", "))
	p.write(")\n")
	Walk(p, n.Body)

	p.indentWrite("End Property")
	p.write("\n")
	return nil
}

func (p *prettyPrinter) visitPropertyLetStmt(n *PropertyLetStmt) Visitor {
	p.indentWrite("")
	if n.Mod.IsValid() {
		p.write(n.Mod.String())
		p.write(" ")
	}
	list := []string{}
	for _, r := range n.Recv {
		if r.TokPos.IsValid() {
			list = append(list, fmt.Sprintf("%s %s", r.Tok, r.Name.Name))
		} else {
			list = append(list, r.Name.Name)
		}
	}
	p.indentWrite("Property Let ")
	Walk(p, n.Name)
	p.write("(")
	p.write(strings.Join(list, ", "))
	p.write(")\n")
	Walk(p, n.Body)

	p.indentWrite("End Property")
	p.write("\n")
	return nil
}

func (p *prettyPrinter) visitFuncDecl(n *FuncDecl) Visitor {
	p.indentWrite("")
	if n.Mod.IsValid() {
		p.write(n.Mod.String())
		p.write(" ")
	}
	if n.Default.IsValid() {
		if n.Mod == token.PRIVATE {
			panic("Used only with the Public keyword in a Class block to indicate that the Function procedure is the default method for the class.")
		}
		p.write("Default ")
	}
	list := []string{}
	for _, r := range n.Recv {
		if r.TokPos.IsValid() {
			list = append(list, fmt.Sprintf("%s %s", r.Tok, r.Name.Name))
		} else {
			list = append(list, r.Name.Name)
		}
	}
	p.indentWrite("Function ")
	Walk(p, n.Name)
	p.write("(")
	p.write(strings.Join(list, ", "))
	p.write(")\n")

	for _, s := range n.Body.List {
		temp := p.indentSpace
		p.indentSpace += "  "
		Walk(p, s)
		p.indentSpace = temp
	}

	p.indentWrite("End Function")
	p.write("\n")
	return nil
}

func (p *prettyPrinter) visitSubDecl(n *SubDecl) Visitor {
	p.indentWrite("")
	if n.Mod.IsValid() {
		p.write(n.Mod.String())
		p.write(" ")
	}

	p.write("Sub ")
	Walk(p, n.Name)
	p.write("\n")

	for _, s := range n.Body.List {
		temp := p.indentSpace
		p.indentSpace += "  "
		Walk(p, s)
		p.indentSpace = temp
	}
	p.indentWrite("End Sub")
	p.write("\n")
	return nil
}

func (p *prettyPrinter) visitClassDecl(n *ClassDecl) Visitor {
	p.indentWrite("")
	if n.Mod.IsValid() {
		p.write(n.Mod.String())
		p.write(" ")
	}
	p.write("Class ")
	Walk(p, n.Name)

	temp := p.indentSpace
	p.indentSpace += "  "

	for _, s := range n.Stmts {
		Walk(p, s)
	}
	p.indentSpace = temp

	p.indentWrite("End Class")
	p.write("\n")
	return nil
}

func (p *prettyPrinter) visitReDimDecl(n *ReDimDecl) Visitor {
	p.indentWrite("ReDim ")
	if n.Preserve.IsValid() {
		p.write("Preserve ")
	}
	p.visitExprList(n.List)
	p.write("\n")
	return nil
}

func (p *prettyPrinter) visitDimDecl(n *DimDecl) Visitor {
	p.indentWrite("Dim ")
	p.visitExprList(n.List)
	if n.Colon.IsValid() {
		p.write(": Set ")
		Walk(p, n.Set.Lhs)
		p.write(" = ")
		Walk(p, n.Set.Rhs)
	}
	p.write("\n")
	return nil
}

func (p *prettyPrinter) visitComment(n *Comment) Visitor {
	p.indentWrite(n.Tok.String())
	p.write(n.Text)
	p.write("\n")
	return nil
}

func (p *prettyPrinter) visitCommentGroup(n *CommentGroup) Visitor {
	for _, c := range n.List {
		Walk(p, c)
	}
	return nil
}

func (p *prettyPrinter) visitFile(n *File) Visitor {
	for _, s := range n.Stmts {
		Walk(p, s)
	}
	return nil
}

func (p *prettyPrinter) visitExprList(exprs []Expr, sep ...string) {
	defaultSep := ", "
	if len(sep) > 0 {
		defaultSep = sep[0]
	}
	for i, e := range exprs {
		if i > 0 {
			p.write(defaultSep)
		}
		Walk(p, e)
	}
}

// Write writes the AST to the output.
func Write(node Node, output io.Writer) {
	Walk(&prettyPrinter{indent: 0, output: output, indentSpace: "  "}, node)
}

// Print prints the AST to the standard output.
func Print(node Node) {
	Write(node, os.Stdout)
}

// String returns the AST as a string.
func String(node Node) string {
	buf := &strings.Builder{}
	Write(node, buf)
	return buf.String()
}
