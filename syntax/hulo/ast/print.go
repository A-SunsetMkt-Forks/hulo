// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package ast

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/hulo-lang/hulo/syntax/hulo/token"
)

// TODO: 只打印类型节点
type printer struct{}

// prettyPrinter holds the state for printing, primarily the output.
type prettyPrinter struct {
	output      io.Writer
	indent      int
	indentSpace string
}

var _ Visitor = (*prettyPrinter)(nil)

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
	indentStr := strings.Repeat(p.indentSpace, p.indent)

	switch n := node.(type) {
	case *File:
		return p.visitFile(n)
	case *FuncDecl:
		return p.visitFuncDecl(n)
	case *ConstructorDecl:
		return p.visitConstructorDecl(n)
	case *BlockStmt:
		return p.visitBlockStmt(n)
	case *AssignStmt:
		return p.visitAssignStmt(n)
	case *Ident:
		p.write(n.Name)
		return nil
	case *NumericLiteral:
		p.write(n.Value)
		return nil
	case *StringLiteral:
		p.write(fmt.Sprintf(`"%s"`, n.Value))
		return nil
	case *CmdStmt:
		return p.visitCmdStmt(n)
	case *EnumDecl:
		return p.visitEnumDecl(n)
	case *ClassDecl:
		return p.visitClassDecl(n)
	case *TraitDecl:
		return p.visitTraitDecl(n)
	case *TypeDecl:
		return p.visitTypeDecl(n)
	case *NullLiteral:
		p.write("null")
		return nil
	case *AnyLiteral:
		p.write("any")
		return nil
	case *TrueLiteral:
		p.write("true")
		return nil
	case *FalseLiteral:
		p.write("false")
		return nil
	case *UnionType:
		return p.visitUnionType(n)
	case *IntersectionType:
		return p.visitIntersectionType(n)
	case *TypeParameter:
		return p.visitTypeParameter(n)
	case *TypeReference:
		return p.visitTypeReference(n)
	case *Parameter:
		return p.visitParameter(n)
	case *NamedParameters:
		return p.visitNamedParameters(n)
	case *NullableType:
		Walk(p, n.X)
		p.write("?")
		return nil
	case *ExtensionDecl:
		return p.visitExtensionDecl(n)
	case *OperatorDecl:
		return p.visitOperatorDecl(n)
	case *CallExpr:
		return p.visitCallExpr(n)
	case *RefExpr:
		return p.visitRefExpr(n)
	case *ReturnStmt:
		return p.visitReturnStmt(n)
	case *ThrowStmt:
		return p.visitThrowStmt(n)
	case *TryStmt:
		return p.visitTryStmt(n)
	case *CatchClause:
		return p.visitCatchClause(n)
	case *FinallyStmt:
		return p.visitFinallyStmt(n)
	case *UnaryExpr:
		return p.visitUnaryExpr(n)
	case *SelectExpr:
		return p.visitSelectExpr(n)
	case *ModAccessExpr:
		return p.visitModAccessExpr(n)
	case *NewDelExpr:
		return p.visitNewDelExpr(n)
	case *IfStmt:
		return p.visitIfStmt(n)
	case *DoWhileStmt:
		return p.visitDoWhileStmt(n)
	case *WhileStmt:
		return p.visitWhileStmt(n)
	case *ForeachStmt:
		return p.visitForeachStmt(n)
	case *ForInStmt:
		return p.visitForInStmt(n)
	case *ForStmt:
		return p.visitForStmt(n)
	case *Decorator:
		return p.visitDecorator(n)
	case *BinaryExpr:
		return p.visitBinaryExpr(n)
	case *ModDecl:
		return p.visitModDecl(n)
	case *UseDecl:
		return p.visitUseDecl(n)
	case *ArrayLiteralExpr:
		return p.visitArrayLiteralExpr(n)
	case *ExprStmt:
		return p.visitExprStmt(n)
	case *IndexExpr:
		return p.visitIndexExpr(n)
	case *SliceExpr:
		return p.visitSliceExpr(n)
	case *IncDecExpr:
		return p.visitIncDecExpr(n)
	case *CascadeExpr:
		return p.visitCascadeExpr(n)
	case *DeclareDecl:
		return p.visitDeclareDecl(n)
	case *ComptimeStmt:
		return p.visitComptimeStmt(n)
	case *ObjectLiteralExpr:
		return p.visitObjectLiteralExpr(n)
	case *KeyValueExpr:
		return p.visitKeyValueExpr(n)
	case *NamedObjectLiteralExpr:
		return p.visitNamedObjectLiteralExpr(n)
	case *TypeLiteral:
		return p.visitTypeLiteral(n)
	case *Import:
		return p.visitImport(n)
	default:
		fmt.Fprintf(p.output, "%s%T\n", indentStr, n)
		panic("unsupport")
	}
	fmt.Println()
	return p
}

func (p *prettyPrinter) visitNamedObjectLiteralExpr(n *NamedObjectLiteralExpr) Visitor {
	Walk(p, n.Name)
	if len(n.Props) == 0 {
		p.write("{}")
	} else {
		p.write("{ ")
		p.visitExprList(n.Props)
		p.write(" }")
	}
	return nil
}

func (p *prettyPrinter) visitKeyValueExpr(n *KeyValueExpr) Visitor {
	Walk(p, n.Key)
	p.write(": ")
	Walk(p, n.Value)
	return nil
}

func (p *prettyPrinter) visitObjectLiteralExpr(n *ObjectLiteralExpr) Visitor {
	if len(n.Props) == 0 {
		p.write("{}")
	} else {
		p.write("{ ")
		p.visitExprList(n.Props)
		p.write(" }")
	}
	return nil
}

func (p *prettyPrinter) visitComptimeStmt(n *ComptimeStmt) Visitor {
	Walk(p, n.X)
	p.write("\n")
	return nil
}

func (p *prettyPrinter) visitDeclareDecl(n *DeclareDecl) Visitor {
	p.indentWrite("declare ")
	Walk(p, n.X)
	return nil
}

func (p *prettyPrinter) visitCascadeExpr(n *CascadeExpr) Visitor {
	Walk(p, n.X)
	p.write("..")
	Walk(p, n.Y)
	return nil
}

func (p *prettyPrinter) visitIncDecExpr(n *IncDecExpr) Visitor {
	if n.Pre {
		p.write(n.Tok.String())
	}
	Walk(p, n.X)
	if !n.Pre {
		p.write(n.Tok.String())
	}
	return nil
}

func (p *prettyPrinter) visitSliceExpr(n *SliceExpr) Visitor {
	Walk(p, n.X)
	p.write("[")
	switch {
	case n.Low != nil && n.High != nil && n.Max != nil:
		Walk(p, n.Low)
		p.write("..")
		Walk(p, n.High)
		p.write("..")
		Walk(p, n.Max)
	case n.Low != nil && n.High != nil:
		Walk(p, n.Low)
		p.write("..")
		Walk(p, n.High)
	case n.Low != nil:
		Walk(p, n.Low)
		p.write("..")
	case n.High != nil:
		p.write("..")
		Walk(p, n.High)
	}
	p.write("]")
	return nil
}

func (p *prettyPrinter) visitIndexExpr(n *IndexExpr) Visitor {
	Walk(p, n.X)
	p.write("[")
	Walk(p, n.Index)
	p.write("]")
	return nil
}

func (p *prettyPrinter) visitExprStmt(n *ExprStmt) Visitor {
	p.indentWrite("")
	Walk(p, n.X)
	p.write("\n")
	return nil
}

func (p *prettyPrinter) visitArrayLiteralExpr(n *ArrayLiteralExpr) Visitor {
	p.write("[")
	p.visitExprList(n.Elems)
	p.write("]")
	return nil
}

func (p *prettyPrinter) visitUseDecl(n *UseDecl) Visitor {
	p.indentWrite("use ")
	Walk(p, n.Lhs)

	if n.Rhs != nil {
		p.write(" = ")
		Walk(p, n.Rhs)
	}
	p.write("\n")
	return nil
}

func (p *prettyPrinter) visitModDecl(n *ModDecl) Visitor {
	p.indentWrite("mod ")
	Walk(p, n.Name)
	Walk(p, n.Body)
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

func (p *prettyPrinter) visitDecorator(n *Decorator) Visitor {
	p.indentWrite("@")
	Walk(p, n.Name)
	if len(n.Recv) > 0 {
		p.write("(")
		p.visitExprList(n.Recv)
		p.write(")")
	}
	p.write("\n")
	return nil
}

func (p *prettyPrinter) visitCmdStmt(n *CmdStmt) Visitor {
	p.indentWrite("")
	Walk(p, n.Name)
	if n.Recv != nil {
		p.write(" ")
		p.visitExprList(n.Recv)
	}
	p.write("\n")
	return nil
}

func (p *prettyPrinter) visitIfStmt(n *IfStmt) Visitor {
	p.indentWrite("if ")
	Walk(p, n.Cond)

	Walk(p, n.Body)

	for n.Else != nil {
		switch el := n.Else.(type) {
		case *IfStmt:
			p.indentWrite("else if ")
			Walk(p, el.Cond)

			Walk(p, el.Body)

			n.Else = el.Else
		case *BlockStmt:
			p.indentWrite("else")
			Walk(p, el)
			n.Else = nil
		}
	}

	return nil
}

func (p *prettyPrinter) visitNewDelExpr(n *NewDelExpr) Visitor {
	p.write(n.Tok.String())
	p.write(" ")
	Walk(p, n.X)
	return nil
}

func (p *prettyPrinter) visitModAccessExpr(n *ModAccessExpr) Visitor {
	Walk(p, n.X)
	p.write("::")
	Walk(p, n.Y)
	return nil
}

func (p *prettyPrinter) visitSelectExpr(n *SelectExpr) Visitor {
	Walk(p, n.X)
	p.write(".")
	Walk(p, n.Y)
	return nil
}

func (p *prettyPrinter) visitUnaryExpr(n *UnaryExpr) Visitor {
	p.write(n.Op.String())
	Walk(p, n.X)
	return nil
}

func (p *prettyPrinter) visitReturnStmt(n *ReturnStmt) Visitor {
	p.indentWrite("return")
	if n.X != nil {
		p.write(" ")
		Walk(p, n.X)
	}
	p.write("\n")
	return nil
}

func (p *prettyPrinter) visitRefExpr(n *RefExpr) Visitor {
	p.write("$")
	Walk(p, n.X)
	return nil
}

func (p *prettyPrinter) visitCallExpr(n *CallExpr) Visitor {
	Walk(p, n.Fun)
	p.write("(")
	p.visitExprList(n.Recv)
	p.write(")")
	return nil
}

func (p *prettyPrinter) visitDecorators(n []*Decorator) {
	for _, dec := range n {
		Walk(p, dec)
	}
}

func (p *prettyPrinter) visitOperatorDecl(n *OperatorDecl) Visitor {
	if n.Decs != nil {
		p.visitDecorators(n.Decs)
	}

	p.indentWrite("operator ")
	p.write(n.Name.String())

	if n.Params != nil {
		p.write("(")
		p.visitExprList(n.Params)
		p.write(")")
	}

	if n.Type != nil {
		p.write(" -> ")
		Walk(p, n.Type)
	}

	if n.Body != nil {
		Walk(p, n.Body)
	}

	return nil
}

func (p *prettyPrinter) visitBlockStmt(n *BlockStmt) Visitor {
	p.write(" {\n")
	p.indent++
	for _, stmt := range n.List {
		Walk(p, stmt)
	}
	p.indent--
	p.indentWrite("}\n")
	return nil
}

func (p *prettyPrinter) visitExtensionDecl(n *ExtensionDecl) Visitor {
	p.indentWrite("extension ")

	switch body := n.Body.(type) {
	case *ExtensionClass:
		p.write("class ")
		Walk(p, n.Name)
		if body.Body != nil {
			Walk(p, body.Body)
		}
	case *ExtensionEnum:
		p.write("enum ")
		Walk(p, n.Name)
		if body.Body != nil {
			Walk(p, body.Body)
		}
	case *ExtensionTrait:
		p.write("trait ")
		Walk(p, n.Name)
		if body.Body != nil {
			Walk(p, body.Body)
		}
	case *ExtensionType:
		p.write("type ")
		Walk(p, n.Name)
		if body.Body != nil {
			Walk(p, body.Body)
		}
	case *ExtensionMod:
		p.write("mod ")
		Walk(p, n.Name)
		if body.Body != nil {
			Walk(p, body.Body)
		}
	}
	return nil
}

func (p *prettyPrinter) visitNamedParameters(n *NamedParameters) Visitor {
	p.write("{")
	p.visitExprList(n.Params)
	p.write("}")
	return nil
}

func (p *prettyPrinter) visitIntersectionType(n *IntersectionType) Visitor {
	p.visitExprList(n.Types, " & ")
	return nil
}

func (p *prettyPrinter) visitUnionType(n *UnionType) Visitor {
	p.visitExprList(n.Types, " | ")
	return nil
}

func (p *prettyPrinter) visitParameter(n *Parameter) Visitor {
	if n.Modifier != nil {
		p.visitModifier(n.Modifier)
	}
	Walk(p, n.Name)

	if n.Type != nil {
		p.write(": ")
		Walk(p, n.Type)
	}

	if n.Value != nil {
		p.write(" = ")
		Walk(p, n.Value)
	}

	return nil
}

func (p *prettyPrinter) visitConstructorDecl(n *ConstructorDecl) Visitor {
	if n.Docs != nil {
		Walk(p, n.Docs)
	}

	for _, dec := range n.Decs {
		Walk(p, dec)
	}

	p.indentWrite("")

	if n.Modifiers != nil {
		p.visitModifiers(n.Modifiers)
	}

	Walk(p, n.ClsName)

	if n.Name != nil {
		p.write(".")
		Walk(p, n.Name)
	}

	p.write("(")
	p.visitExprList(n.Recv)
	p.write(")")

	if n.InitFields != nil {
		p.write(" ")
		p.visitExprList(n.InitFields)
	}

	if n.Body != nil {
		Walk(p, n.Body)
	} else {
		p.write(" {}\n")
	}

	return nil
}

func (p *prettyPrinter) visitFuncDecl(n *FuncDecl) Visitor {
	if n.Docs != nil {
		Walk(p, n.Docs)
	}

	for _, dec := range n.Decs {
		Walk(p, dec)
	}

	p.indentWrite("")

	if n.Modifiers != nil {
		p.visitModifiers(n.Modifiers)
	}

	p.write("fn ")
	if n.Name != nil {
		Walk(p, n.Name)
	}

	if n.TypeParams != nil {
		p.visitTypeParams(n.TypeParams)
	}

	p.write("(")
	for i, param := range n.Recv {
		if i > 0 {
			p.write(", ")
		}
		Walk(p, param)
	}
	p.write(")")

	if n.Throw {
		p.write(" throws")
	}

	if n.Type != nil {
		p.write(" -> ")
		Walk(p, n.Type)
	}

	if n.Body != nil {
		Walk(p, n.Body)
	}

	return nil
}

func (p *prettyPrinter) visitTypeReference(n *TypeReference) Visitor {
	Walk(p, n.Name)
	if n.TypeParams != nil {
		p.visitTypeParams(n.TypeParams)
	}
	return nil
}

func (p *prettyPrinter) visitTypeParameter(n *TypeParameter) Visitor {
	Walk(p, n.Name)
	if n.Constraints != nil {
		p.write(" extends ")
		p.visitExprList(n.Constraints, " + ")
	}

	return nil
}

func (p *prettyPrinter) visitAssignStmt(n *AssignStmt) Visitor {
	if n.Tok == token.COLON_ASSIGN {
		p.indentWrite("let")
	} else {
		p.indentWrite(n.Scope.String())
	}
	p.write(" ")
	Walk(p, n.Lhs)

	if n.Type != nil {
		p.write(": ")
		Walk(p, n.Type)
	}

	p.write(" = ")
	Walk(p, n.Rhs)
	p.write("\n")

	return nil
}

func (p *prettyPrinter) visitModifiers(ms []Modifier) {
	for _, m := range ms {
		p.visitModifier(m)
	}
}

func (p *prettyPrinter) visitModifier(m Modifier) {
	switch m.(type) {
	case *FinalModifier:
		p.write("final ")
	case *ConstModifier:
		p.write("const ")
	case *PubModifier:
		p.write("pub ")
	case *StaticModifier:
		p.write("static ")
	case *RequiredModifier:
		p.write("required ")
	case *ReadonlyModifier:
		p.write("readonly ")
	case *EllipsisModifier:
		p.write("...")
	}
}

func (p *prettyPrinter) visitTypeParams(typeParams []Expr) {
	p.write("<")
	p.visitExprList(typeParams)
	p.write(">")
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

func (p *prettyPrinter) visitTypeDecl(n *TypeDecl) Visitor {
	p.indentWrite("type ")
	Walk(p, n.Name)
	p.write(" = ")
	Walk(p, n.Value)
	return nil
}

func (p *prettyPrinter) visitTraitDecl(n *TraitDecl) Visitor {
	p.indentWrite("trait ")
	return nil
}

func (p *prettyPrinter) visitField(n *Field) Visitor {
	p.indentWrite("")
	if n.Modifiers != nil {
		p.visitModifiers(n.Modifiers)
	}
	Walk(p, n.Name)
	if n.Type != nil {
		p.write(": ")
		Walk(p, n.Type)
	}
	if n.Value != nil {
		p.write(" = ")
		Walk(p, n.Value)
	}
	return nil
}

func (p *prettyPrinter) visitFieldList(fields []*Field) {
	for _, field := range fields {
		p.visitField(field)
		p.write("\n")
	}
}

func (p *prettyPrinter) visitClassDecl(n *ClassDecl) Visitor {
	indentStr := strings.Repeat("  ", p.indent)

	p.write(indentStr)

	// 打印修饰符
	if n.Modifiers != nil {
		p.visitModifiers(n.Modifiers)
	}

	p.write("class ")
	Walk(p, n.Name)
	if n.TypeParams != nil {
		p.visitTypeParams(n.TypeParams)
	}
	p.write(" {\n")

	p.indent++

	if n.Fields != nil {
		p.visitFieldList(n.Fields.List)
	}

	for _, ctor := range n.Ctors {
		Walk(p, ctor)
	}

	for _, method := range n.Methods {
		Walk(p, method)
	}
	p.indent--

	p.write(indentStr)
	p.write("}\n")
	return nil
}

func (p *prettyPrinter) visitEnumDecl(n *EnumDecl) Visitor {
	indentStr := strings.Repeat("  ", p.indent)

	p.write(indentStr)

	if n.Modifiers != nil {
		p.visitModifiers(n.Modifiers)
	}
	fmt.Printf("enum %s", n.Name)

	// Print type parameters
	if n.TypeParams != nil {
		p.visitTypeParams(n.TypeParams)
	}

	// Print enum body
	switch body := n.Body.(type) {
	case *BasicEnumBody:
		p.write(" {\n")
		for i, value := range body.Values {
			fmt.Printf("%s  %s", indentStr, value.Name)
			if value.Value != nil {
				fmt.Printf(" = %s", value.Value)
			}
			if i < len(body.Values)-1 {
				fmt.Print(",")
			}
			fmt.Println()
		}
		fmt.Printf("%s}\n", indentStr)
	case *AssociatedEnumBody:
		p.write(" {\n")
		if body.Fields != nil {
			for _, field := range body.Fields.List {
				fmt.Printf("%s  %s: %s\n", indentStr, field.Name, field.Type)
			}
		}
		for i, value := range body.Values {
			fmt.Printf("%s  %s", indentStr, value.Name)
			if value.Data != nil {
				fmt.Print("(")
				p.visitExprList(value.Data)
				p.write(")")
			}
			if i < len(body.Values)-1 {
				p.write(",")
			}
			fmt.Println()
		}
		fmt.Printf("%s}\n", indentStr)
	case *ADTEnumBody:
		p.write(" {\n")
		for i, variant := range body.Variants {
			fmt.Printf("%s  %s", indentStr, variant.Name)
			if variant.Fields != nil {
				p.write(" {")
				for j, field := range variant.Fields.List {
					if j > 0 {
						p.write(", ")
					}
					fmt.Printf("%s: %s", field.Name, field.Type)
				}
				p.write("}")
			}
			if i < len(body.Variants)-1 {
				p.write(",")
			}
			fmt.Println()
		}
		if body.Methods != nil {
			for _, method := range body.Methods {
				print(method, indentStr+"  ")
			}
		}
		fmt.Printf("%s}\n", indentStr)
	}
	return nil
}

func (p *prettyPrinter) visitFile(n *File) Visitor {
	for _, stmt := range n.Stmts {
		Walk(p, stmt)
	}
	return nil
}

func (p *prettyPrinter) visitTryStmt(n *TryStmt) Visitor {
	p.indentWrite("try ")
	Walk(p, n.Body)

	for _, catch := range n.Catches {
		Walk(p, catch)
	}

	if n.Finally != nil {
		Walk(p, n.Finally)
	}
	return nil
}

func (p *prettyPrinter) visitCatchClause(n *CatchClause) Visitor {
	p.indentWrite("catch ")
	if n.Cond != nil {
		p.write("(")
		Walk(p, n.Cond)
		p.write(")")
	}
	Walk(p, n.Body)
	return nil
}

func (p *prettyPrinter) visitFinallyStmt(n *FinallyStmt) Visitor {
	p.indentWrite("finally ")
	Walk(p, n.Body)
	return nil
}

func (p *prettyPrinter) visitThrowStmt(n *ThrowStmt) Visitor {
	p.indentWrite("throw")
	p.write(" ")
	Walk(p, n.X)
	return nil
}

func (p *prettyPrinter) visitWhileStmt(n *WhileStmt) Visitor {
	p.indentWrite("loop ")
	Walk(p, n.Cond)

	Walk(p, n.Body)
	return nil
}

func (p *prettyPrinter) visitDoWhileStmt(n *DoWhileStmt) Visitor {
	p.indentWrite("do ")
	Walk(p, n.Body)
	p.write(" loop ")
	p.write("(")
	Walk(p, n.Cond)
	p.write(")")
	return nil
}

func (p *prettyPrinter) visitForStmt(n *ForStmt) Visitor {
	p.indentWrite("loop ")
	if n.Init != nil {
		if as, ok := n.Init.(*AssignStmt); ok {
			Walk(p, as.Lhs)
			p.write(" := ")
			Walk(p, as.Rhs)
		} else {
			Walk(p, n.Init)
		}
	}
	p.write("; ")
	Walk(p, n.Cond)
	p.write("; ")
	Walk(p, n.Post)

	Walk(p, n.Body)
	return nil
}

func (p *prettyPrinter) visitForInStmt(n *ForInStmt) Visitor {
	p.indentWrite("loop ")
	Walk(p, n.Index)
	p.write(" in ")
	if n.RangeExpr.Start != nil {
		Walk(p, n.RangeExpr.Start)
	}

	if n.RangeExpr.End_ != nil {
		p.write("..")
		Walk(p, n.RangeExpr.End_)
	}
	if n.RangeExpr.Step != nil {
		p.write("..")
		Walk(p, n.RangeExpr.Step)
	}
	Walk(p, n.Body)
	return nil
}

func (p *prettyPrinter) visitForeachStmt(n *ForeachStmt) Visitor {
	p.indentWrite("loop (")
	if n.Index != nil {
		Walk(p, n.Index)
	}
	if n.Value != nil {
		p.write(", ")
		Walk(p, n.Value)
	}
	p.write(")")
	if n.Var != nil {
		p.write(" ")
		p.write(n.Tok.String())
		p.write(" ")
		Walk(p, n.Var)
	}
	Walk(p, n.Body)
	return nil
}

func (p *prettyPrinter) visitTypeLiteral(n *TypeLiteral) Visitor {
	p.write("{ ")
	p.visitExprList(n.Members)
	p.write(" }")
	return nil
}

func (p *prettyPrinter) visitImport(n *Import) Visitor {
	p.indentWrite("import ")

	// Handle different types of imports
	if n.ImportSingle != nil {
		p.write(fmt.Sprintf(`"%s"`, n.ImportSingle.Path))
		if n.ImportSingle.Alias != "" {
			p.write(" as ")
			p.write(n.ImportSingle.Alias)
		}
	} else if n.ImportAll != nil {
		p.write("*")
		if n.ImportAll.Alias != "" {
			p.write(" as ")
			p.write(n.ImportAll.Alias)
		}
		p.write(" from ")
		p.write(fmt.Sprintf(`"%s"`, n.ImportAll.Path))
	} else if n.ImportMulti != nil {
		p.write("{ ")
		for i, field := range n.ImportMulti.List {
			if i > 0 {
				p.write(", ")
			}
			p.write(field.Field)
			if field.Alias != "" {
				p.write(" as ")
				p.write(field.Alias)
			}
		}
		p.write(" } from ")
		p.write(fmt.Sprintf(`"%s"`, n.ImportMulti.Path))
	}

	p.write("\n")
	return nil
}

type PrettyPrinterOption func(*prettyPrinter) error

func WithIndentSpace(space string) PrettyPrinterOption {
	return func(p *prettyPrinter) error {
		p.indentSpace = space
		return nil
	}
}

func Print(node Node, options ...PrettyPrinterOption) error {
	printer := prettyPrinter{output: os.Stdout, indentSpace: "  "}

	for _, opt := range options {
		err := opt(&printer)
		if err != nil {
			return err
		}
	}
	printer.Visit(node)

	return nil
}
