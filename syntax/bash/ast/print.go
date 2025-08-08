// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package ast

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/hulo-lang/hulo/syntax/bash/token"
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
	case *FuncDecl:
		p.write("*ast.FuncDecl {\n")
		p.indent++
		if n.Name != nil {
			p.indentWrite("Name: ")
			Walk(p, n.Name)
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
		if n.X != nil {
			p.indentWrite("X: ")
			Walk(p, n.X)
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
	case *WhileStmt:
		p.write("*ast.WhileStmt {\n")
		p.indent++
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
		p.indent--
		p.indentWrite("}")
	case *UntilStmt:
		p.write("*ast.UntilStmt {\n")
		p.indent++
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
		p.indent--
		p.indentWrite("}")
	case *ForStmt:
		p.write("*ast.ForStmt {\n")
		p.indent++
		if n.Init != nil {
			p.indentWrite("Init: ")
			Walk(p, n.Init)
			p.write("\n")
		}
		if n.Cond != nil {
			p.indentWrite("Cond: ")
			Walk(p, n.Cond)
			p.write("\n")
		}
		if n.Post != nil {
			p.indentWrite("Post: ")
			Walk(p, n.Post)
			p.write("\n")
		}
		if n.Body != nil {
			p.indentWrite("Body: ")
			Walk(p, n.Body)
			p.write("\n")
		}
		p.indent--
		p.indentWrite("}")
	case *ForInStmt:
		p.write("*ast.ForInStmt {\n")
		p.indent++
		if n.Var != nil {
			p.indentWrite("Var: ")
			Walk(p, n.Var)
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
	case *SelectStmt:
		p.write("*ast.SelectStmt {\n")
		p.indent++
		if n.Var != nil {
			p.indentWrite("Var: ")
			Walk(p, n.Var)
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
	case *IfStmt:
		p.write("*ast.IfStmt {\n")
		p.indent++
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
	case *CaseStmt:
		p.write("*ast.CaseStmt {\n")
		p.indent++
		if n.X != nil {
			p.indentWrite("X: ")
			Walk(p, n.X)
			p.write("\n")
		}
		if len(n.Patterns) > 0 {
			p.indentWrite("Patterns: [\n")
			p.indent++
			for i, pattern := range n.Patterns {
				p.indentWrite(fmt.Sprintf("%d: *ast.CasePattern {\n", i))
				p.indent++
				if len(pattern.Conds) > 0 {
					p.indentWrite("Conds: [\n")
					p.indent++
					for j, cond := range pattern.Conds {
						p.indentWrite(fmt.Sprintf("%d: ", j))
						Walk(p, cond)
						p.write("\n")
					}
					p.indent--
					p.indentWrite("]\n")
				}
				if pattern.Body != nil {
					p.indentWrite("Body: ")
					Walk(p, pattern.Body)
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
			Walk(p, n.Else)
			p.write("\n")
		}
		p.indent--
		p.indentWrite("}")
	case *AssignStmt:
		p.write("*ast.AssignStmt {\n")
		p.indent++
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
		p.indentWrite(fmt.Sprintf("Local: %t\n", n.Local.IsValid()))
		p.indent--
		p.indentWrite("}")
	case *PipelineExpr:
		p.write("*ast.PipelineExpr {\n")
		p.indent++
		p.indentWrite(fmt.Sprintf("CtrOp: %s\n", n.CtrOp))
		if len(n.Cmds) > 0 {
			p.indentWrite("Cmds: [\n")
			p.indent++
			for i, cmd := range n.Cmds {
				p.indentWrite(fmt.Sprintf("%d: ", i))
				Walk(p, cmd)
				p.write("\n")
			}
			p.indent--
			p.indentWrite("]\n")
		}
		p.indent--
		p.indentWrite("}")
	case *Redirect:
		p.write("*ast.Redirect {\n")
		p.indent++
		if n.N != nil {
			p.indentWrite("N: ")
			Walk(p, n.N)
			p.write("\n")
		}
		p.indentWrite(fmt.Sprintf("CtrOp: %s\n", n.CtrOp))
		if n.Word != nil {
			p.indentWrite("Word: ")
			Walk(p, n.Word)
			p.write("\n")
		}
		p.indent--
		p.indentWrite("}")
	case *Word:
		p.write(fmt.Sprintf("*ast.Word (Val: %q)", n.Val))
	case *Ident:
		p.write(fmt.Sprintf("*ast.Ident (Name: %q)", n.Name))
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
	case *CmdExpr:
		p.write("*ast.CmdExpr {\n")
		p.indent++
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
	case *TestExpr:
		p.write("*ast.TestExpr {\n")
		p.indent++
		if n.X != nil {
			p.indentWrite("X: ")
			Walk(p, n.X)
			p.write("\n")
		}
		p.indent--
		p.indentWrite("}")
	case *ExtendedTestExpr:
		p.write("*ast.ExtendedTestExpr {\n")
		p.indent++
		if n.X != nil {
			p.indentWrite("X: ")
			Walk(p, n.X)
			p.write("\n")
		}
		p.indent--
		p.indentWrite("}")
	case *ArithEvalExpr:
		p.write("*ast.ArithEvalExpr {\n")
		p.indent++
		if n.X != nil {
			p.indentWrite("X: ")
			Walk(p, n.X)
			p.write("\n")
		}
		p.indent--
		p.indentWrite("}")
	case *CmdSubst:
		p.write("*ast.CmdSubst {\n")
		p.indent++
		p.indentWrite(fmt.Sprintf("Tok: %s\n", n.Tok))
		if n.X != nil {
			p.indentWrite("X: ")
			Walk(p, n.X)
			p.write("\n")
		}
		p.indent--
		p.indentWrite("}")
	case *ProcSubst:
		p.write("*ast.ProcSubst {\n")
		p.indent++
		p.indentWrite(fmt.Sprintf("CtrOp: %s\n", n.CtrOp))
		if n.X != nil {
			p.indentWrite("X: ")
			Walk(p, n.X)
			p.write("\n")
		}
		p.indent--
		p.indentWrite("}")
	case *ArithExpr:
		p.write("*ast.ArithExpr {\n")
		p.indent++
		if n.X != nil {
			p.indentWrite("X: ")
			Walk(p, n.X)
			p.write("\n")
		}
		p.indent--
		p.indentWrite("}")
	case *VarExpExpr:
		p.write("*ast.VarExpExpr {\n")
		p.indent++
		if n.X != nil {
			p.indentWrite("X: ")
			Walk(p, n.X)
			p.write("\n")
		}
		p.indent--
		p.indentWrite("}")
	case *ParamExpExpr:
		p.write("*ast.ParamExpExpr {\n")
		p.indent++
		if n.Var != nil {
			p.indentWrite("Var: ")
			Walk(p, n.Var)
			p.write("\n")
		}
		if n.ParamExp != nil {
			p.indentWrite("ParamExp: ")
			p.write(fmt.Sprintf("%T", n.ParamExp))
			p.write("\n")
		}
		p.indent--
		p.indentWrite("}")
	case *IndexExpr:
		p.write("*ast.IndexExpr {\n")
		p.indent++
		if n.X != nil {
			p.indentWrite("X: ")
			Walk(p, n.X)
			p.write("\n")
		}
		if n.Y != nil {
			p.indentWrite("Y: ")
			Walk(p, n.Y)
			p.write("\n")
		}
		p.indent--
		p.indentWrite("}")
	case *ArrExpr:
		p.write("*ast.ArrExpr {\n")
		p.indent++
		if len(n.Vars) > 0 {
			p.indentWrite("Vars: [\n")
			p.indent++
			for i, v := range n.Vars {
				p.indentWrite(fmt.Sprintf("%d: ", i))
				Walk(p, v)
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
	case *CmdListExpr:
		p.write("*ast.CmdListExpr {\n")
		p.indent++
		p.indentWrite(fmt.Sprintf("CtrOp: %s\n", n.CtrOp))
		if len(n.Cmds) > 0 {
			p.indentWrite("Cmds: [\n")
			p.indent++
			for i, cmd := range n.Cmds {
				p.indentWrite(fmt.Sprintf("%d: ", i))
				Walk(p, cmd)
				p.write("\n")
			}
			p.indent--
			p.indentWrite("]\n")
		}
		p.indent--
		p.indentWrite("}")
	case *CmdGroup:
		p.write("*ast.CmdGroup {\n")
		p.indent++
		p.indentWrite(fmt.Sprintf("Op: %s\n", n.Op))
		if len(n.List) > 0 {
			p.indentWrite("List: [\n")
			p.indent++
			for i, item := range n.List {
				p.indentWrite(fmt.Sprintf("%d: ", i))
				Walk(p, item)
				p.write("\n")
			}
			p.indent--
			p.indentWrite("]\n")
		}
		p.indent--
		p.indentWrite("}")
	case *Comment:
		p.write(fmt.Sprintf("*ast.Comment (Text: %q)", n.Text))
	case *CommentGroup:
		p.write("*ast.CommentGroup {\n")
		p.indent++
		if len(n.List) > 0 {
			p.indentWrite("List: [\n")
			p.indent++
			for i, comment := range n.List {
				p.indentWrite(fmt.Sprintf("%d: ", i))
				Walk(p, comment)
				p.write("\n")
			}
			p.indent--
			p.indentWrite("]\n")
		}
		p.indent--
		p.indentWrite("}")
	case *ReturnStmt:
		p.write("*ast.ReturnStmt {\n")
		p.indent++
		if n.X != nil {
			p.indentWrite("X: ")
			Walk(p, n.X)
			p.write("\n")
		}
		p.indent--
		p.indentWrite("}")
	default:
		p.write(fmt.Sprintf("unknown node type: %T", node))
	}
	return nil
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
	switch node := node.(type) {
	case *File:
		return p.visitFile(node)
	case *FuncDecl:
		return p.visitFuncDecl(node)
	case *ExprStmt:
		return p.visitExprStmt(node)
	case *BlockStmt:
		return p.visitBlockStmt(node)
	case *WhileStmt:
		return p.visitWhileStmt(node)
	case *UntilStmt:
		return p.visitUntilStmt(node)
	case *ForStmt:
		return p.visitForStmt(node)
	case *ForInStmt:
		return p.visitForInStmt(node)
	case *SelectStmt:
		return p.visitSelectStmt(node)
	case *IfStmt:
		return p.visitIfStmt(node)
	case *CaseStmt:
		return p.visitCaseStmt(node)
	case *AssignStmt:
		return p.visitAssignStmt(node)
	case *PipelineExpr:
		return p.visitPipelineExpr(node)
	case *Redirect:
		return p.visitRedirect(node)
	/// Expressions
	case *Word:
		p.write(node.Val)
		return nil
	case *Ident:
		p.write(node.Name)
		return nil
	case *BinaryExpr:
		return p.visitBinaryExpr(node)
	case *CmdExpr:
		return p.visitCallExpr(node)
	case *TestExpr:
		return p.visitTestExpr(node)
	case *ExtendedTestExpr:
		return p.visitExtendedTestExpr(node)
	case *ArithEvalExpr:
		return p.visitArithEvalExpr(node)
	case *CmdSubst:
		return p.visitCmdSubst(node)
	case *ProcSubst:
		return p.visitProcSubst(node)
	case *ArithExpr:
		return p.visitArithExpr(node)
	case *VarExpExpr:
		p.write("$")
		Walk(p, node.X)
		return nil
	case *ParamExpExpr:
		return p.visitParamExpExpr(node)
	case *IndexExpr:
		Walk(p, node.X)
		p.write("[")
		Walk(p, node.Y)
		p.write("]")
		return nil
	case *ArrExpr:
		p.write("(")
		p.visitExprs(node.Vars)
		p.write(")")
		return nil
	case *UnaryExpr:
		p.write(node.Op.String())
		Walk(p, node.X)
		return nil
	case *CmdListExpr:
		return p.visitCmdListExpr(node)
	case *CmdGroup:
		return p.visitCmdGroup(node)
	case *Comment:
		return p.visitComment(node)
	case *CommentGroup:
		return p.visitCommentGroup(node)
	case *ReturnStmt:
		return p.visitReturnStmt(node)
	default:
		panic("unsupported node type: " + fmt.Sprintf("%T", node))
	}
}

func (p *prettyPrinter) visitReturnStmt(node *ReturnStmt) Visitor {
	p.indentWrite("return ")
	Walk(p, node.X)
	p.write("\n")
	return nil
}

func (p *prettyPrinter) visitComment(node *Comment) Visitor {
	p.indentWrite("#")
	p.write(node.Text)
	p.write("\n")
	return nil
}

func (p *prettyPrinter) visitCommentGroup(node *CommentGroup) Visitor {
	for _, c := range node.List {
		Walk(p, c)
	}
	return nil
}

func (p *prettyPrinter) visitRedirect(node *Redirect) Visitor {
	if node.N != nil {
		Walk(p, node.N)
	}
	p.write(node.CtrOp.String())
	Walk(p, node.Word)
	return nil
}

func (p *prettyPrinter) visitCmdGroup(node *CmdGroup) Visitor {
	p.indentWrite("")
	if node.Op == token.LeftBrace {
		p.write("{")
	} else {
		p.write("(")
	}
	p.visitExprs(node.List)
	if node.Op == token.LeftBrace {
		p.write("}")
	} else {
		p.write(")")
	}
	return nil
}

func (p *prettyPrinter) visitCmdListExpr(node *CmdListExpr) Visitor {
	p.indentWrite("")
	for i, cmd := range node.Cmds {
		Walk(p, cmd)
		if i < len(node.Cmds)-1 {
			p.write(fmt.Sprintf(" %s ", node.CtrOp))
		}
	}
	return nil
}

func (p *prettyPrinter) visitFile(node *File) Visitor {
	for _, s := range node.Stmts {
		Walk(p, s)
	}
	return nil
}

func (p *prettyPrinter) visitFuncDecl(node *FuncDecl) Visitor {
	p.indentWrite("function ")
	Walk(p, node.Name)
	p.write("() {\n")
	Walk(p, node.Body)
	p.indentWrite("}\n")
	return nil
}

func (p *prettyPrinter) visitExprStmt(node *ExprStmt) Visitor {
	p.indentWrite("")
	Walk(p, node.X)
	p.write("\n")
	return nil
}

func (p *prettyPrinter) visitBlockStmt(node *BlockStmt) Visitor {
	p.indent++
	for _, s := range node.List {
		Walk(p, s)
	}
	p.indent--
	return nil
}

func (p *prettyPrinter) visitWhileStmt(node *WhileStmt) Visitor {
	p.indentWrite("while ")
	Walk(p, node.Cond)
	p.write("; do\n")
	Walk(p, node.Body)
	p.indentWrite("done\n")
	return nil
}

func (p *prettyPrinter) visitUntilStmt(node *UntilStmt) Visitor {
	p.indentWrite("until ")
	Walk(p, node.Cond)
	p.write("; do\n")
	Walk(p, node.Body)
	p.indentWrite("done\n")
	return nil
}

func (p *prettyPrinter) visitForStmt(node *ForStmt) Visitor {
	p.indentWrite("for (( ")
	if node.Init != nil {
		if init, ok := node.Init.(*AssignStmt); ok {
			Walk(p, init.Lhs)
			p.write("=")
			Walk(p, init.Rhs)
		} else {
			Walk(p, node.Init)
		}
	}
	p.write("; ")
	if node.Cond != nil {
		Walk(p, node.Cond)
	}
	p.write("; ")
	if node.Post != nil {
		if post, ok := node.Post.(*AssignStmt); ok {
			Walk(p, post.Lhs)
			p.write("=")
			Walk(p, post.Rhs)
		} else {
			Walk(p, node.Post)
		}
	}
	p.write(" )); do\n")

	Walk(p, node.Body)

	p.indentWrite("done\n")
	return nil
}

func (p *prettyPrinter) visitForInStmt(node *ForInStmt) Visitor {
	p.indentWrite("for ")
	Walk(p, node.Var)
	p.write(" in ")
	Walk(p, node.List)
	p.write("; do\n")

	Walk(p, node.Body)
	p.indentWrite("done\n")
	return nil
}

func (p *prettyPrinter) visitSelectStmt(node *SelectStmt) Visitor {
	p.indentWrite("select ")
	Walk(p, node.Var)
	p.write(" in ")
	Walk(p, node.List)
	p.write("; do\n")

	Walk(p, node.Body)
	p.indentWrite("done\n")
	return nil
}

func (p *prettyPrinter) visitCaseStmt(node *CaseStmt) Visitor {
	p.indentWrite("case ")
	Walk(p, node.X)
	p.write(" in\n")

	for _, pattern := range node.Patterns {
		p.indentWrite("  ")
		p.visitExprs(pattern.Conds, "|")
		p.write(")\n")
		p.indent++
		Walk(p, pattern.Body)
		p.indent--
		p.indentWrite("    ;;\n")
	}

	if node.Else != nil {
		p.indentWrite("  *)\n")
		p.indent++
		Walk(p, node.Else)
		p.indent--
		p.indentWrite("    ;;\n")
	}

	p.indentWrite("esac\n")
	return nil
}

func (p *prettyPrinter) visitAssignStmt(node *AssignStmt) Visitor {
	p.indentWrite("")
	if node.Local.IsValid() {
		p.write("local ")
	}
	Walk(p, node.Lhs)
	p.write("=")
	Walk(p, node.Rhs)
	p.write("\n")
	return nil
}

func (p *prettyPrinter) visitArithExpr(node *ArithExpr) Visitor {
	p.write("$(( ")
	Walk(p, node.X)
	p.write(" ))")
	return nil
}

func (p *prettyPrinter) visitProcSubst(node *ProcSubst) Visitor {
	p.write(node.CtrOp.String())
	Walk(p, node.X)
	p.write(")")
	return nil
}

func (p *prettyPrinter) visitCmdSubst(node *CmdSubst) Visitor {
	p.write(node.Tok.String())
	Walk(p, node.X)
	switch node.Tok {
	case token.DollParen:
		p.write(")")
	case token.BckQuote:
		p.write("`")
	case token.DollBrace:
		p.write("}")
	}
	return nil
}

func (p *prettyPrinter) visitArithEvalExpr(node *ArithEvalExpr) Visitor {
	p.write("(( ")
	Walk(p, node.X)
	p.write(" ))")
	return nil
}

func (p *prettyPrinter) visitExtendedTestExpr(node *ExtendedTestExpr) Visitor {
	p.write("[[ ")
	Walk(p, node.X)
	p.write(" ]]")
	return nil
}

func (p *prettyPrinter) visitBinaryExpr(node *BinaryExpr) Visitor {
	Walk(p, node.X)

	switch node.Op {
	case token.TsReMatch, token.TsNempStr, token.TsEmpStr,
		token.TsNewer, token.TsOlder, token.TsDevIno,
		token.TsEql, token.TsNeq, token.TsLeq,
		token.TsGeq, token.TsLss, token.TsGtr:
		p.write(fmt.Sprintf(" %s ", node.Op))
	default:
		p.write(node.Op.String())
	}
	Walk(p, node.Y)
	return nil
}

func (p *prettyPrinter) visitTestExpr(node *TestExpr) Visitor {
	p.write("[ ")
	Walk(p, node.X)
	p.write(" ]")
	return nil
}

func (p *prettyPrinter) visitCallExpr(node *CmdExpr) Visitor {
	Walk(p, node.Name)

	if len(node.Recv) > 0 {
		p.write(" ")
		p.visitExprs(node.Recv)
	}
	return nil
}

func (p *prettyPrinter) visitIfStmt(node *IfStmt) Visitor {
	p.indentWrite("if ")
	Walk(p, node.Cond)
	p.write("; then\n")

	Walk(p, node.Body)

	for node.Else != nil {
		switch el := node.Else.(type) {
		case *IfStmt:
			p.indentWrite("elif ")
			Walk(p, el.Cond)
			p.write("; then\n")

			Walk(p, el.Body)

			node.Else = el.Else
		case *BlockStmt:
			p.indentWrite("else\n")
			Walk(p, el)
			node.Else = nil
		}
	}

	p.indentWrite("fi\n")
	return nil
}

func (p *prettyPrinter) visitParamExpExpr(node *ParamExpExpr) Visitor {
	switch exp := node.ParamExp.(type) {
	case *DefaultValExp:
		p.write("${")
		Walk(p, node.Var)
		p.write(":-")
		Walk(p, exp.Val)
		p.write("}")
	case *DefaultValAssignExp:
		p.write("${")
		Walk(p, node.Var)
		p.write(":=")
		Walk(p, exp.Val)
		p.write("}")
	case *NonNullCheckExp:
		p.write("${")
		Walk(p, node.Var)
		p.write(":?")
		Walk(p, exp.Val)
		p.write("}")
	case *NonNullExp:
		p.write("${")
		Walk(p, node.Var)
		p.write(":+")
		Walk(p, exp.Val)
		p.write("}")
	case *PrefixExp:
		p.write("${!")
		Walk(p, node.Var)
		p.write("*}")
	case *PrefixArrayExp:
		p.write("${!")
		Walk(p, node.Var)
		p.write("@}")
	case *ArrayIndexExp:
		p.write("${!")
		Walk(p, node.Var)
		if exp.Tok == token.Star {
			p.write("[*]}")
		} else {
			p.write("[@]}")
		}
	case *LengthExp:
		p.write("${#")
		Walk(p, node.Var)
		p.write("}")
	case *DelPrefix:
		p.write("${")
		Walk(p, node.Var)
		if exp.Longest {
			p.write("##")
		} else {
			p.write("#")
		}
		Walk(p, exp.Val)
		p.write("}")
	case *DelSuffix:
		p.write("${")
		Walk(p, node.Var)
		if exp.Longest {
			p.write("%%%%")
		} else {
			p.write("%%")
		}
		Walk(p, exp.Val)
		p.write("}")
	case *SubstringExp:
		p.write("${")
		Walk(p, node.Var)
		p.write(":")
		if exp.Offset != exp.Length {
			p.write(fmt.Sprintf("%d:%d", exp.Offset, exp.Length))
		} else {
			p.write(fmt.Sprintf("%d", exp.Offset))
		}
		p.write("}")
	case *ReplaceExp:
		p.write("${")
		Walk(p, node.Var)
		p.write("/")
		p.write(exp.Old)
		p.write("/")
		p.write(exp.New)
		p.write("}")
	case *ReplacePrefixExp:
		p.write("${")
		Walk(p, node.Var)
		p.write("/#")
		p.write(exp.Old)
		p.write("/")
		p.write(exp.New)
		p.write("}")
	case *ReplaceSuffixExp:
		p.write("${")
		Walk(p, node.Var)
		p.write("/%%")
		p.write(exp.Old)
		p.write("/")
		p.write(exp.New)
		p.write("}")
	case *CaseConversionExp:
		p.write("${")
		Walk(p, node.Var)
		if exp.FirstChar && exp.ToUpper {
			p.write("^")
		} else if !exp.FirstChar && exp.ToUpper {
			p.write("^^")
		} else if exp.FirstChar && !exp.ToUpper {
			p.write(",")
		} else {
			p.write(",,")
		}
		p.write("}")
	case *OperatorExp:
		p.write("${")
		Walk(p, node.Var)
		p.write("@")
		p.write(string(exp.Op))
		p.write("}")
	default:
		p.write("${")
		Walk(p, node.Var)
		p.write("}")
	}
	return nil
}

func (p *prettyPrinter) visitExprs(exprs []Expr, sep ...string) {
	sepStr := " "
	if len(sep) > 0 {
		sepStr = sep[0]
	}
	for i, e := range exprs {
		Walk(p, e)
		if i < len(exprs)-1 {
			p.write(sepStr)
		}
	}
}

func (p *prettyPrinter) visitPipelineExpr(node *PipelineExpr) Visitor {
	for i, cmd := range node.Cmds {
		Walk(p, cmd)
		if i < len(node.Cmds)-1 {
			p.write(fmt.Sprintf(" %s ", node.CtrOp))
		}
	}
	return nil
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
