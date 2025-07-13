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
		for _, s := range node.Stmts {
			Walk(p, s)
		}
	case *FuncDecl:
		p.indentWrite("function ")
		Walk(p, node.Name)
		p.write("() {\n")
		Walk(p, node.Body)
		p.indentWrite("}\n")

	case *ExprStmt:
		p.indentWrite("")
		Walk(p, node.X)
		p.write("\n")
		return nil

	case *BlockStmt:
		p.indent++
		for _, s := range node.List {
			Walk(p, s)
		}
		p.indent--

	case *WhileStmt:
		p.indentWrite("while ")
		Walk(p, node.Cond)
		p.write("; do\n")
		Walk(p, node.Body)
		p.indentWrite("done\n")

	case *UntilStmt:
		p.indentWrite("until ")
		Walk(p, node.Cond)
		p.write("; do\n")
		Walk(p, node.Body)
		p.indentWrite("done\n")

	case *ForStmt:
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
		p.write(")); do\n")

		Walk(p, node.Body)

		p.indentWrite("done\n")

	case *ForInStmt:
		p.indentWrite("for ")
		Walk(p, node.Var)
		p.write(" in ")
		Walk(p, node.List)
		p.write("; do\n")

		Walk(p, node.Body)
		p.indentWrite("done\n")

	case *SelectStmt:
		p.indentWrite("select ")
		Walk(p, node.Var)
		p.write(" in ")
		Walk(p, node.List)
		p.write("; do\n")

		Walk(p, node.Body)
		p.indentWrite("done\n")

	case *IfStmt:
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
				p.write("else\n")
				Walk(p, el)
				node.Else = nil
			}
		}

		p.indentWrite("fi\n")

	case *CaseStmt:
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
			p.indentWrite("  ;;\n")
		}

		if node.Else != nil {
			p.indentWrite("  *)\n")
			p.indent++
			Walk(p, node.Else)
			p.indent--
			p.indentWrite("  ;;\n")
		}

		p.indentWrite("esac\n")

	case *AssignStmt:
		p.indentWrite("")
		if node.Local.IsValid() {
			p.write("local ")
		}
		Walk(p, node.Lhs)
		p.write("=")
		Walk(p, node.Rhs)
		p.write("\n")
		return nil
	case *BreakStmt:
		p.indentWrite("break\n")
		return nil
	case *ContinueStmt:
		p.indentWrite("continue\n")
		return nil
	/// Expressions
	case *Lit:
		p.write(node.Val)
		return nil
	case *Ident:
		p.write(node.Name)
		return nil
	case *BinaryExpr:
		Walk(p, node.X)

		switch node.Op {
		case token.TsLss, token.TsGtr:
			p.write(fmt.Sprintf(" %s ", node.Op))
		default:
			p.write(node.Op.String())
		}
		Walk(p, node.Y)
		return nil
	case *CallExpr:
		Walk(p, node.Func)

		if len(node.Recv) > 0 {
			p.write(" ")
			p.visitExprs(node.Recv)
		}
		return nil
	case *TestExpr:
		p.write("[ ")
		Walk(p, node.X)
		p.write(" ]")
		return nil
	case *ExtendedTestExpr:
		p.write("[[ ")
		Walk(p, node.X)
		p.write(" ]]")
		return nil
	case *ArithEvalExpr:
		p.write("(( ")
		Walk(p, node.X)
		p.write(" ))")
		return nil
	case *CmdSubst:
		if node.Tok == token.LeftParen {
			p.write("$( ")
		} else {
			p.write("` ")
		}
		Walk(p, node.X)
		if node.Tok == token.LeftParen {
			p.write(" )")
		} else {
			p.write(" `")
		}
		return nil
	case *ProcSubst:
		if node.Tok == token.RdrIn {
			p.write("<( ")
		} else {
			p.write(">( ")
		}
		Walk(p, node.X)
		p.write(" )")
		return nil
	case *ArithExpr:
		p.write("$(( ")
		Walk(p, node.X)
		p.write(" ))")
		return nil
	case *VarExpExpr:
		p.write("$")
		Walk(p, node.X)
		return nil
	case *ParamExpExpr:
		switch {
		case node.DefaultValExp != nil:
			p.write("${")
			Walk(p, node.Var)
			p.write(":-")
			Walk(p, node.DefaultValExp.Val)
			p.write("}")
			return nil
		case node.DefaultValAssignExp != nil:
			p.write("${")
			Walk(p, node.Var)
			p.write(":=")
			Walk(p, node.DefaultValAssignExp.Val)
			p.write("}")
			return nil
		case node.NonNullCheckExp != nil:
			p.write("${")
			Walk(p, node.Var)
			p.write(":?")
			Walk(p, node.NonNullCheckExp.Val)
			p.write("}")
			return nil
		case node.NonNullExp != nil:
			p.write("${")
			Walk(p, node.Var)
			p.write(":+")
			Walk(p, node.NonNullExp.Val)
			p.write("}")
			return nil
		case node.PrefixExp != nil:
			p.write("${!")
			Walk(p, node.Var)
			p.write("*}")
			return nil
		case node.PrefixArrayExp != nil:
			p.write("${!")
			Walk(p, node.Var)
			p.write("@}")
			return nil
		case node.ArrayIndexExp != nil:
			p.write("${!")
			Walk(p, node.Var)
			if node.Tok == token.Star {
				p.write("[*]}")
			} else {
				p.write("[@]}")
			}
			return nil
		case node.LengthExp != nil:
			p.write("${#")
			Walk(p, node.Var)
			p.write("}")
			return nil
		case node.DelPrefix != nil:
			p.write("${")
			Walk(p, node.Var)
			if node.DelPrefix.Longest {
				p.write("##")
			} else {
				p.write("#")
			}
			Walk(p, node.DelPrefix.Val)
			p.write("}")
			return nil
		case node.DelSuffix != nil:
			p.write("${")
			Walk(p, node.Var)
			if node.DelSuffix.Longest {
				p.write("%%%%")
			} else {
				p.write("%%")
			}
			Walk(p, node.DelSuffix.Val)
			p.write("}")
			return nil
		case node.SubstringExp != nil:
			p.write("${")
			Walk(p, node.Var)
			p.write(":")
			if node.SubstringExp.Offset != node.SubstringExp.Length {
				p.write(fmt.Sprintf("%d:%d", node.SubstringExp.Offset, node.SubstringExp.Length))
			} else {
				p.write(fmt.Sprintf("%d", node.SubstringExp.Offset))
			}
			p.write("}")
			return nil
		case node.ReplaceExp != nil:
			p.write("${")
			Walk(p, node.Var)
			p.write("/")
			p.write(node.ReplaceExp.Old)
			p.write("/")
			p.write(node.ReplaceExp.New)
			p.write("}")
			return nil
		case node.ReplacePrefixExp != nil:
			p.write("${")
			Walk(p, node.Var)
			p.write("/#")
			p.write(node.ReplacePrefixExp.Old)
			p.write("/")
			p.write(node.ReplacePrefixExp.New)
			p.write("}")
			return nil
		case node.ReplaceSuffixExp != nil:
			p.write("${")
			Walk(p, node.Var)
			p.write("/%%")
			p.write(node.ReplaceSuffixExp.Old)
			p.write("/")
			p.write(node.ReplaceSuffixExp.New)
			p.write("}")
			return nil
		case node.CaseConversionExp != nil:
			p.write("${")
			Walk(p, node.Var)
			if node.CaseConversionExp.FirstChar && node.CaseConversionExp.ToUpper {
				p.write("^")
			} else if !node.CaseConversionExp.FirstChar && node.CaseConversionExp.ToUpper {
				p.write("^^")
			} else if node.CaseConversionExp.FirstChar && !node.CaseConversionExp.ToUpper {
				p.write(",")
			} else {
				p.write(",,")
			}
			p.write("}")
			return nil
		case node.OperatorExp != nil:
			p.write("${")
			Walk(p, node.Var)
			p.write("@")
			p.write(string(node.OperatorExp.Op))
			p.write("}")
			return nil
		default:
			p.write("${")
			Walk(p, node.Var)
			p.write("}")
			return nil
		}
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
	default:
		panic("unsupported node type: " + fmt.Sprintf("%T", node))
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

// Print prints the AST to the standard output.
func Print(node Node) {
	Walk(&prettyPrinter{indent: 0, output: os.Stdout, indentSpace: "  "}, node)
}

// String returns the AST as a string.
func String(node Node) string {
	buf := &strings.Builder{}
	Walk(&prettyPrinter{indent: 0, output: buf, indentSpace: "  "}, node)
	return buf.String()
}
