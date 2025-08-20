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

func Inspect(node Node, output io.Writer) {
	printer := &printer{output: output, indent: 0}
	printer.visit(node)
}

type printer struct {
	output io.Writer
	indent int
}

func (p *printer) write(s string) {
	fmt.Fprint(p.output, s)
}

func (p *printer) indentWrite(s string) {
	fmt.Fprint(p.output, strings.Repeat("  ", p.indent))
	fmt.Fprint(p.output, s)
}

func (p *printer) visit(node Node) {
	if node == nil {
		p.write("nil")
		return
	}

	switch n := node.(type) {
	case *File:
		p.write("*ast.File {\n")
		p.indent++
		if len(n.Docs) > 0 {
			p.indentWrite("Docs: ")
			p.visit(n.Docs[0])
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
	case *Ident:
		p.write(fmt.Sprintf("*ast.Ident (Name: %q)", n.Name))
	case *NumericLiteral:
		p.write(fmt.Sprintf("*ast.NumericLiteral (Value: %q)", n.Value))
	case *StringLiteral:
		p.write(fmt.Sprintf("*ast.StringLiteral (Value: %q)", n.Value))
	case *TrueLiteral:
		p.write("*ast.TrueLiteral")
	case *FalseLiteral:
		p.write("*ast.FalseLiteral")
	case *NullLiteral:
		p.write("*ast.NullLiteral")
	case *AnyLiteral:
		p.write("*ast.AnyLiteral")
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
	case *FuncDecl:
		p.write("*ast.FuncDecl {\n")
		p.indent++
		if len(n.Decs) > 0 {
			p.indentWrite("Decs: [\n")
			p.indent++
			for i, dec := range n.Decs {
				p.indentWrite(fmt.Sprintf("%d: ", i))
				p.visit(dec)
				p.write("\n")
			}
			p.indent--
			p.indentWrite("]\n")
		}
		if n.Name != nil {
			p.indentWrite("Name: ")
			p.visit(n.Name)
			p.write("\n")
		}
		if len(n.TypeParams) > 0 {
			p.indentWrite("TypeParams: [\n")
			p.indent++
			for i, param := range n.TypeParams {
				p.indentWrite(fmt.Sprintf("%d: ", i))
				p.visit(param)
				p.write("\n")
			}
			p.indent--
			p.indentWrite("]\n")
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
		if n.Type != nil {
			p.indentWrite("Type: ")
			p.visit(n.Type)
			p.write("\n")
		}
		if n.Body != nil {
			p.indentWrite("Body: ")
			p.visit(n.Body)
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
	case *AssignStmt:
		p.write("*ast.AssignStmt {\n")
		p.indent++
		if n.Lhs != nil {
			p.indentWrite("Lhs: ")
			p.visit(n.Lhs)
			p.write("\n")
		}
		if n.Tok != 0 {
			p.indentWrite(fmt.Sprintf("Tok: %s\n", n.Tok))
		}
		if n.Rhs != nil {
			p.indentWrite("Rhs: ")
			p.visit(n.Rhs)
			p.write("\n")
		}
		p.indent--
		p.indentWrite("}")
	case *CallExpr:
		p.write("*ast.CallExpr {\n")
		p.indent++
		if n.Fun != nil {
			p.indentWrite("Fun: ")
			p.visit(n.Fun)
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
	case *RefExpr:
		p.write("*ast.RefExpr {\n")
		p.indent++
		if n.X != nil {
			p.indentWrite("X: ")
			p.visit(n.X)
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
	case *UnaryExpr:
		p.write("*ast.UnaryExpr {\n")
		p.indent++
		p.indentWrite(fmt.Sprintf("Op: %s\n", n.Op))
		if n.X != nil {
			p.indentWrite("X: ")
			p.visit(n.X)
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
	case *WhileStmt:
		p.write("*ast.WhileStmt {\n")
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
	case *ForStmt:
		p.write("*ast.ForStmt {\n")
		p.indent++
		if n.Init != nil {
			p.indentWrite("Init: ")
			p.visit(n.Init)
			p.write("\n")
		}
		if n.Cond != nil {
			p.indentWrite("Cond: ")
			p.visit(n.Cond)
			p.write("\n")
		}
		if n.Post != nil {
			p.indentWrite("Post: ")
			p.visit(n.Post)
			p.write("\n")
		}
		if n.Body != nil {
			p.indentWrite("Body: ")
			p.visit(n.Body)
			p.write("\n")
		}
		p.indent--
		p.indentWrite("}")
	case *ForeachStmt:
		p.write("*ast.ForeachStmt {\n")
		p.indent++
		if n.Index != nil {
			p.indentWrite("Index: ")
			p.visit(n.Index)
			p.write("\n")
		}
		if n.Value != nil {
			p.indentWrite("Value: ")
			p.visit(n.Value)
			p.write("\n")
		}
		if n.Var != nil {
			p.indentWrite("Var: ")
			p.visit(n.Var)
			p.write("\n")
		}
		if n.Body != nil {
			p.indentWrite("Body: ")
			p.visit(n.Body)
			p.write("\n")
		}
		p.indent--
		p.indentWrite("}")
	case *ReturnStmt:
		p.write("*ast.ReturnStmt {\n")
		p.indent++
		if n.X != nil {
			p.indentWrite("X: ")
			p.visit(n.X)
			p.write("\n")
		}
		p.indent--
		p.indentWrite("}")
	case *DeferStmt:
		p.write("*ast.DeferStmt {\n")
		p.indent++
		if n.X != nil {
			p.indentWrite("X: ")
			p.visit(n.X)
			p.write("\n")
		}
		p.indent--
		p.indentWrite("}")
	case *Decorator:
		p.write("*ast.Decorator {\n")
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
	case *ArrayLiteralExpr:
		p.write("*ast.ArrayLiteralExpr {\n")
		p.indent++
		if len(n.Elems) > 0 {
			p.indentWrite("Elems: [\n")
			p.indent++
			for i, elem := range n.Elems {
				p.indentWrite(fmt.Sprintf("%d: ", i))
				p.visit(elem)
				p.write("\n")
			}
			p.indent--
			p.indentWrite("]\n")
		}
		p.indent--
		p.indentWrite("}")
	case *ObjectLiteralExpr:
		p.write("*ast.ObjectLiteralExpr {\n")
		p.indent++
		if len(n.Props) > 0 {
			p.indentWrite("Props: [\n")
			p.indent++
			for i, prop := range n.Props {
				p.indentWrite(fmt.Sprintf("%d: ", i))
				p.visit(prop)
				p.write("\n")
			}
			p.indent--
			p.indentWrite("]\n")
		}
		p.indent--
		p.indentWrite("}")
	case *KeyValueExpr:
		p.write("*ast.KeyValueExpr {\n")
		p.indent++
		if n.Key != nil {
			p.indentWrite("Key: ")
			p.visit(n.Key)
			p.write("\n")
		}
		if n.Value != nil {
			p.indentWrite("Value: ")
			p.visit(n.Value)
			p.write("\n")
		}
		p.indent--
		p.indentWrite("}")
	case *SelectExpr:
		p.write("*ast.SelectExpr {\n")
		p.indent++
		if n.X != nil {
			p.indentWrite("X: ")
			p.visit(n.X)
			p.write("\n")
		}
		if n.Y != nil {
			p.indentWrite("Y: ")
			p.visit(n.Y)
			p.write("\n")
		}
		p.indent--
		p.indentWrite("}")
	case *ModAccessExpr:
		p.write("*ast.ModAccessExpr {\n")
		p.indent++
		if n.X != nil {
			p.indentWrite("X: ")
			p.visit(n.X)
			p.write("\n")
		}
		if n.Y != nil {
			p.indentWrite("Y: ")
			p.visit(n.Y)
			p.write("\n")
		}
		p.indent--
		p.indentWrite("}")
	case *IncDecExpr:
		p.write("*ast.IncDecExpr {\n")
		p.indent++
		if n.X != nil {
			p.indentWrite("X: ")
			p.visit(n.X)
			p.write("\n")
		}
		p.indentWrite(fmt.Sprintf("Tok: %s\n", n.Tok))
		p.indentWrite(fmt.Sprintf("Pre: %t\n", n.Pre))
		p.indent--
		p.indentWrite("}")
	case *ClassDecl:
		p.write("*ast.ClassDecl {\n")
		p.indent++
		if len(n.Decs) > 0 {
			p.indentWrite("Decs: [\n")
			p.indent++
			for i, dec := range n.Decs {
				p.indentWrite(fmt.Sprintf("%d: ", i))
				p.visit(dec)
				p.write("\n")
			}
			p.indent--
			p.indentWrite("]\n")
		}
		if n.Name != nil {
			p.indentWrite("Name: ")
			p.visit(n.Name)
			p.write("\n")
		}
		if n.Fields != nil && len(n.Fields.List) > 0 {
			p.indentWrite("Fields: [\n")
			p.indent++
			for i, field := range n.Fields.List {
				p.indentWrite(fmt.Sprintf("%d: *ast.Field {\n", i))
				p.indent++
				if len(field.Decs) > 0 {
					p.indentWrite("Decs: [\n")
					p.indent++
					for j, dec := range field.Decs {
						p.indentWrite(fmt.Sprintf("%d: ", j))
						p.visit(dec)
						p.write("\n")
					}
					p.indent--
					p.indentWrite("]\n")
				}
				if field.Name != nil {
					p.indentWrite("Name: ")
					p.visit(field.Name)
					p.write("\n")
				}
				if field.Type != nil {
					p.indentWrite("Type: ")
					p.visit(field.Type)
					p.write("\n")
				}
				if field.Value != nil {
					p.indentWrite("Value: ")
					p.visit(field.Value)
					p.write("\n")
				}
				p.indent--
				p.indentWrite("}\n")
			}
			p.indent--
			p.indentWrite("]\n")
		}
		if len(n.Methods) > 0 {
			p.indentWrite("Methods: [\n")
			p.indent++
			for i, method := range n.Methods {
				p.indentWrite(fmt.Sprintf("%d: ", i))
				p.visit(method)
				p.write("\n")
			}
			p.indent--
			p.indentWrite("]\n")
		}
		p.indent--
		p.indentWrite("}")
	case *EnumDecl:
		p.write("*ast.EnumDecl {\n")
		p.indent++
		if len(n.Decs) > 0 {
			p.indentWrite("Decs: [\n")
			p.indent++
			for i, dec := range n.Decs {
				p.indentWrite(fmt.Sprintf("%d: ", i))
				p.visit(dec)
				p.write("\n")
			}
			p.indent--
			p.indentWrite("]\n")
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
	case *MatchStmt:
		p.write("*ast.MatchStmt {\n")
		p.indent++
		if n.Expr != nil {
			p.indentWrite("Expr: ")
			p.visit(n.Expr)
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
		if n.Default != nil {
			p.indentWrite("Default: *ast.CaseClause {\n")
			p.indent++
			if n.Default.Cond != nil {
				p.indentWrite("Cond: ")
				p.visit(n.Default.Cond)
				p.write("\n")
			}
			if n.Default.Body != nil {
				p.indentWrite("Body: ")
				p.visit(n.Default.Body)
				p.write("\n")
			}
			p.indent--
			p.indentWrite("}\n")
		}
		p.indent--
		p.indentWrite("}")
	case *DeclareDecl:
		p.write("*ast.DeclareDecl {\n")
		p.indent++
		if n.X != nil {
			p.indentWrite("X: ")
			p.visit(n.X)
			p.write("\n")
		}
		p.indent--
		p.indentWrite("}")
	case *Import:
		p.write("*ast.Import {\n")
		p.indent++
		if n.ImportSingle != nil {
			p.indentWrite(fmt.Sprintf("ImportSingle: *ast.ImportSingle (Path: %q, Alias: %q)\n", n.ImportSingle.Path, n.ImportSingle.Alias))
		}
		if n.ImportAll != nil {
			p.indentWrite(fmt.Sprintf("ImportAll: *ast.ImportAll (Path: %q, Alias: %q)\n", n.ImportAll.Path, n.ImportAll.Alias))
		}
		if n.ImportMulti != nil {
			p.indentWrite("*ast.ImportMulti {\n")
			p.indent++
			if len(n.ImportMulti.List) > 0 {
				p.indentWrite("List: [\n")
				p.indent++
				for i, item := range n.ImportMulti.List {
					p.indentWrite(fmt.Sprintf("%d: *ast.ImportField (Field: %q, Alias: %q)\n", i, item.Field, item.Alias))
				}
				p.indent--
				p.indentWrite("]\n")
			}
			p.indentWrite(fmt.Sprintf("Path: %q\n", n.ImportMulti.Path))
			p.indent--
			p.indentWrite("}\n")
		}
		p.indent--
		p.indentWrite("}")
	case *UnsafeExpr:
		p.write(fmt.Sprintf("*ast.UnsafeStmt (Text: %q)", n.Text))
	case *ExternDecl:
		p.write("*ast.ExternDecl {\n")
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
	case *Parameter:
		p.write("*ast.Parameter {\n")
		p.indent++
		if n.Name != nil {
			p.indentWrite("Name: ")
			p.visit(n.Name)
			p.write("\n")
		}
		if n.Type != nil {
			p.indentWrite("Type: ")
			p.visit(n.Type)
			p.write("\n")
		}
		if n.Value != nil {
			p.indentWrite("Value: ")
			p.visit(n.Value)
			p.write("\n")
		}
		p.indent--
		p.indentWrite("}")
	case *DoWhileStmt:
		p.write("*ast.DoWhileStmt {\n")
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
	case *ForInStmt:
		p.write("*ast.ForInStmt {\n")
		p.indent++
		if n.Index != nil {
			p.indentWrite("Index: ")
			p.visit(n.Index)
			p.write("\n")
		}
		if n.RangeExpr.Start != nil {
			p.indentWrite("RangeExpr.Start: ")
			p.visit(n.RangeExpr.Start)
			p.write("\n")
		}
		if n.RangeExpr.End_ != nil {
			p.indentWrite("RangeExpr.End: ")
			p.visit(n.RangeExpr.End_)
			p.write("\n")
		}
		if n.RangeExpr.Step != nil {
			p.indentWrite("RangeExpr.Step: ")
			p.visit(n.RangeExpr.Step)
			p.write("\n")
		}
		if n.Body != nil {
			p.indentWrite("Body: ")
			p.visit(n.Body)
			p.write("\n")
		}
		p.indent--
		p.indentWrite("}")
	case *CmdStmt:
		p.write("*ast.CmdStmt {\n")
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
	case *CmdExpr:
		p.write("*ast.CmdExpr {\n")
		p.indent++
		if n.Cmd != nil {
			p.indentWrite("Cmd: ")
			p.visit(n.Cmd)
			p.write("\n")
		}
		if len(n.Args) > 0 {
			p.indentWrite("Args: [\n")
			p.indent++
			for i, arg := range n.Args {
				p.indentWrite(fmt.Sprintf("%d: ", i))
				p.visit(arg)
				p.write("\n")
			}
			p.indent--
			p.indentWrite("]\n")
		}
		if n.IsAsync {
			p.indentWrite("IsAsync: true\n")
		}
		if len(n.BuiltinArgs) > 0 {
			p.indentWrite("BuiltinArgs: [\n")
			p.indent++
			for i, arg := range n.BuiltinArgs {
				p.indentWrite(fmt.Sprintf("%d: ", i))
				p.visit(arg)
				p.write("\n")
			}
			p.indent--
			p.indentWrite("]\n")
		}
		p.indent--
		p.indentWrite("}")
	case *TypeReference:
		p.write("*ast.TypeReference {\n")
		p.indent++
		if n.Name != nil {
			p.indentWrite("Name: ")
			p.visit(n.Name)
			p.write("\n")
		}
		if len(n.TypeParams) > 0 {
			p.indentWrite("TypeParams: [\n")
			p.indent++
			for i, param := range n.TypeParams {
				p.indentWrite(fmt.Sprintf("%d: ", i))
				p.visit(param)
				p.write("\n")
			}
			p.indent--
			p.indentWrite("]\n")
		}
		p.indent--
		p.indentWrite("}")
	case *ConditionalType:
		p.write("*ast.ConditionalType {\n")
		p.indent++
		if n.CheckType != nil {
			p.indentWrite("CheckType: ")
			p.visit(n.CheckType)
			p.write("\n")
		}
		if n.ExtendsType != nil {
			p.indentWrite("ExtendsType: ")
			p.visit(n.ExtendsType)
			p.write("\n")
		}
		if n.TrueType != nil {
			p.indentWrite("TrueType: ")
			p.visit(n.TrueType)
			p.write("\n")
		}
		if n.FalseType != nil {
			p.indentWrite("FalseType: ")
			p.visit(n.FalseType)
			p.write("\n")
		}
		p.indent--
		p.indentWrite("}")
	case *MappedType:
		p.write("*ast.MappedType {\n")
		p.indent++
		if n.KeyName != nil {
			p.indentWrite("KeyName: ")
			p.visit(n.KeyName)
			p.write("\n")
		}
		if n.KeyofType != nil {
			p.indentWrite("KeyofType: ")
			p.visit(n.KeyofType)
			p.write("\n")
		}
		if n.ValueType != nil {
			p.indentWrite("ValueType: ")
			p.visit(n.ValueType)
			p.write("\n")
		}
		if n.Readonly {
			p.indentWrite("Readonly: true\n")
		}
		if n.Optional {
			p.indentWrite("Optional: true\n")
		}
		if n.Required {
			p.indentWrite("Required: true\n")
		}
		p.indent--
		p.indentWrite("}")
	case *KeyofType:
		p.write("*ast.KeyofType {\n")
		p.indent++
		if n.Type != nil {
			p.indentWrite("Type: ")
			p.visit(n.Type)
			p.write("\n")
		}
		p.indent--
		p.indentWrite("}")
	case *InferType:
		p.write("*ast.InferType {\n")
		p.indent++
		if n.X != nil {
			p.indentWrite("Type: ")
			p.visit(n.X)
			p.write("\n")
		}
		p.indent--
		p.indentWrite("}")
	case *UnionType:
		p.write("*ast.UnionType {\n")
		p.indent++
		if len(n.Types) > 0 {
			p.indentWrite("Types: [\n")
			p.indent++
			for i, typ := range n.Types {
				p.indentWrite(fmt.Sprintf("%d: ", i))
				p.visit(typ)
				p.write("\n")
			}
			p.indent--
			p.indentWrite("]\n")
		}
		p.indent--
		p.indentWrite("}")
	case *FunctionType:
		p.write("*ast.FunctionType {\n")
		p.indent++
		if len(n.Recv) > 0 {
			p.indentWrite("Recv: [\n")
			p.indent++
			for i, recv := range n.Recv {
				p.indentWrite(fmt.Sprintf("%d: ", i))
				p.visit(recv)
				p.write("\n")
			}
			p.indent--
			p.indentWrite("]\n")
		}
		if n.RetVal != nil {
			p.indentWrite("RetVal: ")
			p.visit(n.RetVal)
			p.write("\n")
		}
		p.indent--
		p.indentWrite("}")
	case *ArrayType:
		p.write("*ast.ArrayType {\n")
		p.indent++
		if n.Name != nil {
			p.indentWrite("Name: ")
			p.visit(n.Name)
			p.write("\n")
		}
		p.indent--
		p.indentWrite("}")
	case *TypeParameter:
		p.write("*ast.TypeParameter {\n")
		p.indent++
		if n.Name != nil {
			p.indentWrite("Name: ")
			p.visit(n.Name)
			p.write("\n")
		}
		if len(n.Constraints) > 0 {
			p.indentWrite("Constraints: [\n")
			p.indent++
			for i, constraint := range n.Constraints {
				p.indentWrite(fmt.Sprintf("%d: ", i))
				p.visit(constraint)
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
}

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
	case *CmdExpr:
		return p.visitCmdExpr(n)
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
	case *ConditionalType:
		return p.visitConditionalType(n)
	case *MappedType:
		return p.visitMappedType(n)
	case *KeyofType:
		return p.visitKeyofType(n)
	case *InferType:
		return p.visitInferType(n)
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
	case *DeferStmt:
		return p.visitDeferStmt(n)
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
	case *FunctionType:
		return p.visitFunctionType(n)
	case *TupleType:
		return p.visitTupleType(n)
	case *ArrayType:
		return p.visitArrayType(n)
	case *Import:
		return p.visitImport(n)
	case *MatchStmt:
		return p.visitMatchStmt(n)
	case *ConditionalExpr:
		return p.visitConditionalExpr(n)
	case *UnsafeExpr:
		return p.visitUnsafeStmt(n)
	case *ExternDecl:
		return p.visitExternDecl(n)
	case *CommentGroup:
		return p.visitCommentGroup(n)
	case *Comment:
		return p.visitComment(n)
	default:
		panic("unsupported node type: " + fmt.Sprintf("%T", n))
	}
}

func (p *prettyPrinter) visitCommentGroup(n *CommentGroup) Visitor {
	for _, cmt := range n.List {
		Walk(p, cmt)
	}
	return nil
}

func (p *prettyPrinter) visitComment(n *Comment) Visitor {
	p.write("//")
	p.write(n.Text)
	p.write("\n")
	return nil
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
	p.indentWrite("comptime")
	if n.Cond != nil {
		p.write(" when ")
		Walk(p, n.Cond)
	}
	Walk(p, n.Body)
	for n.Else != nil {
		switch el := n.Else.(type) {
		case *ComptimeStmt:
			p.write("else")
			if el.Cond != nil {
				p.write(" when ")
				Walk(p, el.Cond)
			}
			Walk(p, el.Body)
			n.Else = el.Else
		}
	}
	p.write("\n")
	return nil
}

func (p *prettyPrinter) visitDeclareDecl(n *DeclareDecl) Visitor {
	Walk(p, n.Docs)
	p.indentWrite("declare ")
	if fx, ok := n.X.(*FuncDecl); ok {
		fx.Body = &BlockStmt{}
	}
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

func (p *prettyPrinter) visitCmdExpr(n *CmdExpr) Visitor {
	Walk(p, n.Cmd)
	if len(n.Args) > 0 {
		p.write(" ")
		p.visitExprList(n.Args)
	}
	if n.IsAsync {
		p.write(" &")
	}
	if len(n.BuiltinArgs) > 0 {
		p.write(" --- ")
		p.visitExprList(n.BuiltinArgs)
	}
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

func (p *prettyPrinter) visitDeferStmt(n *DeferStmt) Visitor {
	p.indentWrite("defer ")
	if n.X != nil {
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
	if len(n.List) == 0 {
		p.write(" {}\n")
		return nil
	}
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
	p.indentWrite("")
	if n.Docs != nil {
		Walk(p, n.Docs)
	}
	if n.Scope != 0 {
		p.write(n.Scope.String())
		p.write(" ")
	}
	if n.Lhs != nil {
		Walk(p, n.Lhs)
	}
	if n.Type != nil {
		p.write(": ")
		Walk(p, n.Type)
	}
	if n.Tok != 0 {
		p.write(" ")
		p.write(n.Tok.String())
		p.write(" ")
	}
	if n.Rhs != nil {
		Walk(p, n.Rhs)
	}
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
	p.write("\n")
	return nil
}

func (p *prettyPrinter) visitTraitDecl(_ *TraitDecl) Visitor {
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
	fmt.Printf("enum %s", n.Name.Name)

	// Print type parameters
	if n.TypeParams != nil {
		p.visitTypeParams(n.TypeParams)
	}

	p.indent++
	// Print enum body
	switch body := n.Body.(type) {
	case *BasicEnumBody:
		p.write(" {\n")
		for i, value := range body.Values {
			p.indentWrite("")
			Walk(p, value.Name)
			if value.Value != nil {
				p.write(" = ")
				Walk(p, value.Value)
			}
			if i < len(body.Values)-1 {
				p.write(",")
			}
			p.write("\n")
		}
		p.write(indentStr)
		p.write("}\n")
	case *AssociatedEnumBody:
		p.write(" {\n")
		if body.Fields != nil {
			for _, field := range body.Fields.List {
				p.indentWrite("")
				p.visitModifiers(field.Modifiers)
				if field.Name != nil {
					Walk(p, field.Name)
				}
				if field.Type != nil {
					p.write(": ")
					Walk(p, field.Type)
				}
			}
			p.write(";\n")
		}
		for i, value := range body.Values {
			p.indentWrite("")
			Walk(p, value.Name)
			if value.Data != nil {
				p.write("(")
				p.visitExprList(value.Data)
				p.write(")")
			}
			if i < len(body.Values)-1 {
				p.write(",")
			} else {
				p.write(";")
			}
			p.write("\n")
		}
		if len(body.Methods) > 0 {
			p.write(";\n")
			for _, method := range body.Methods {
				Walk(p, method)
			}
		}
		p.write(indentStr)
		p.write("}\n")
	case *ADTEnumBody:
		p.write(" {\n")
		for i, variant := range body.Variants {
			fmt.Printf("%s  %s", indentStr, variant.Name.Name)
			if variant.Fields != nil {
				p.write(" {")
				for j, field := range variant.Fields.List {
					if j > 0 {
						p.write(", ")
					}
					fmt.Printf("%s: %s", field.Name.Name, field.Type)
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
		p.write(indentStr)
		p.write("}\n")
	}
	p.indent--

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

func (p *prettyPrinter) visitFunctionType(n *FunctionType) Visitor {
	p.write("(")
	p.visitExprList(n.Recv)
	p.write(") -> ")
	if n.RetVal != nil {
		Walk(p, n.RetVal)
	}
	return nil
}

func (p *prettyPrinter) visitTupleType(n *TupleType) Visitor {
	p.write("[")
	p.visitExprList(n.Types)
	p.write("]")
	return nil
}

func (p *prettyPrinter) visitArrayType(n *ArrayType) Visitor {
	Walk(p, n.Name)
	p.write("[]")
	return nil
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

// MatchStmt pretty printer
func (p *prettyPrinter) visitMatchStmt(n *MatchStmt) Visitor {
	p.indentWrite("match ")
	if n.Expr != nil {
		Walk(p, n.Expr)
	}
	p.write(" {\n")
	p.indent++
	for _, c := range n.Cases {
		p.visitCaseClause(c)
	}
	if n.Default != nil {
		p.visitCaseClause(n.Default)
	}
	p.indent--
	p.indentWrite("}\n")
	return nil
}

// CaseClause pretty printer
func (p *prettyPrinter) visitCaseClause(n *CaseClause) Visitor {
	p.indentWrite("")
	if n.Cond != nil {
		Walk(p, n.Cond)
	}
	p.write(" =>")
	if n.Body != nil {
		Walk(p, n.Body)
	}
	return nil
}

func (p *prettyPrinter) visitConditionalExpr(n *ConditionalExpr) Visitor {
	Walk(p, n.Cond)
	p.write(" ? ")
	Walk(p, n.WhenTrue)
	p.write(" : ")
	Walk(p, n.WhneFalse)
	return nil
}

// visitUnsafeStmt pretty printer for unsafe statements
func (p *prettyPrinter) visitUnsafeStmt(n *UnsafeExpr) Visitor {
	p.indentWrite("unsafe {\n")
	p.indent++

	// Print the unsafe code as a comment or raw text
	if n.Text != "" {
		// Split the text into lines and print each line
		lines := strings.Split(n.Text, "\n")
		for _, line := range lines {
			if strings.TrimSpace(line) != "" {
				p.indentWrite(line + "\n")
			}
		}
	}

	p.indent--
	p.indentWrite("}\n")
	return nil
}

// visitExternDecl pretty printer for extern declarations
func (p *prettyPrinter) visitExternDecl(n *ExternDecl) Visitor {
	p.indentWrite("extern ")

	// Print all extern items
	for i, item := range n.List {
		if i > 0 {
			p.write(", ")
		}

		// Handle KeyValueExpr which represents extern items
		if kv, ok := item.(*KeyValueExpr); ok {
			Walk(p, kv.Key) // Print the identifier
			p.write(": ")
			Walk(p, kv.Value) // Print the type
		} else {
			Walk(p, item)
		}
	}

	p.write("\n")
	return nil
}

func (p *prettyPrinter) visitConditionalType(n *ConditionalType) Visitor {
	p.indentWrite("ConditionalType {")
	if n.CheckType != nil {
		p.write(" CheckType: ")
		Walk(p, n.CheckType)
	}
	if n.ExtendsType != nil {
		p.write(" ExtendsType: ")
		Walk(p, n.ExtendsType)
	}
	if n.TrueType != nil {
		p.write(" TrueType: ")
		Walk(p, n.TrueType)
	}
	if n.FalseType != nil {
		p.write(" FalseType: ")
		Walk(p, n.FalseType)
	}
	p.write("}")
	return nil
}

func (p *prettyPrinter) visitMappedType(n *MappedType) Visitor {
	p.indentWrite("MappedType {")
	if n.KeyName != nil {
		p.write(" KeyName: ")
		Walk(p, n.KeyName)
	}
	if n.KeyofType != nil {
		p.write(" KeyofType: ")
		Walk(p, n.KeyofType)
	}
	if n.ValueType != nil {
		p.write(" ValueType: ")
		Walk(p, n.ValueType)
	}
	if n.Readonly {
		p.write(" Readonly: true")
	}
	if n.Optional {
		p.write(" Optional: true")
	}
	if n.Required {
		p.write(" Required: true")
	}
	p.write("}")
	return nil
}

func (p *prettyPrinter) visitKeyofType(n *KeyofType) Visitor {
	p.indentWrite("KeyofType {")
	if n.Type != nil {
		p.write(" Type: ")
		Walk(p, n.Type)
	}
	p.write("}")
	return nil
}

func (p *prettyPrinter) visitInferType(n *InferType) Visitor {
	p.indentWrite("InferType {")
	if n.X != nil {
		p.write(" Type: ")
		Walk(p, n.X)
	}
	p.write("}")
	return nil
}
