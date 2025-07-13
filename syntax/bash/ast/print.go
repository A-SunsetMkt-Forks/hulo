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

func (p *prettyPrinter) print(a ...any) {
	fmt.Fprint(p.output, a...)
}

func (p *prettyPrinter) printf(format string, a ...any) {
	fmt.Fprintf(p.output, strings.Repeat("  ", p.indent)+format, a...)
}

func (p *prettyPrinter) println(a ...any) {
	fmt.Fprint(p.output, strings.Repeat("  ", p.indent))
	fmt.Fprintln(p.output, a...)
}

func (p *prettyPrinter) Visit(node Node) Visitor {
	switch node := node.(type) {
	case *File:
		for _, s := range node.Stmts {
			Walk(p, s)
		}
	case *FuncDecl:
		p.write("function ")
		Walk(p, node.Name)
		p.write("() {\n")

		Walk(p, node.Body)

		p.write("}\n")

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
		p.printf("while %s; do\n", node.Cond)
		p.indent++
		Walk(p, node.Body)
		p.indent--
		p.println("done")

	case *UntilStmt:
		p.printf("until %s; do\n", node.Cond)
		p.indent++
		Walk(p, node.Body)
		p.indent--
		p.println("done")

	case *ForStmt:
		p.print("for (( ")
		if node.Init != nil {
			if init, ok := node.Init.(*AssignStmt); ok {
				p.printf("%s=%s", init.Lhs, init.Rhs)
			} else {
				p.printf("%s", node.Init)
			}
		}
		p.printf("; ")
		if node.Cond != nil {
			p.printf("%s", node.Cond)
		}
		p.printf("; ")
		if node.Post != nil {
			if post, ok := node.Post.(*AssignStmt); ok {
				p.printf("%s=%s", post.Lhs, post.Rhs)
			} else {
				p.printf("%s", node.Post)
			}
		}
		p.println(")); do")
		p.indent++
		Walk(p, node.Body)
		p.indent--
		p.println("done")

	case *ForInStmt:
		p.printf("for %s in %s; do\n", node.Var, node.List)
		p.indent++
		Walk(p, node.Body)
		p.indent--
		p.println("done")

	case *SelectStmt:
		p.printf("select %s in %s; do\n", node.Var, node.List)
		p.indent++
		Walk(p, node.Body)
		p.indent--
		p.println("done")

	case *IfStmt:
		p.printf("if %s; then\n", node.Cond)
		p.indent++
		Walk(p, node.Body)
		p.indent--
		for node.Else != nil {
			switch el := node.Else.(type) {
			case *IfStmt:
				p.printf("elif %s; then\n", el.Cond)
				p.indent++
				Walk(p, el.Body)
				p.indent--
				node.Else = el.Else
			case *BlockStmt:
				p.println("else")
				Walk(p, el)
				node.Else = nil
			}
		}

		p.println("fi")

	case *CaseStmt:
		p.printf("case %s in\n", node.X)

		for _, pattern := range node.Patterns {
			p.printf("  %s)\n", pattern.Conds)
			p.indent += 2
			Walk(p, pattern.Body)
			p.indent -= 2
			p.println(";;")
		}

		if node.Else != nil {
			p.println("  *)")
			p.indent += 2
			Walk(p, node.Else)
			p.indent -= 2
		}

		p.println("esac")

	case *AssignStmt:
		if node.Local.IsValid() {
			p.printf("local %s=%s\n", node.Lhs, node.Rhs)
		} else {
			p.printf("%s=%s\n", node.Lhs, node.Rhs)
		}
	case *BreakStmt:
		p.println("break")
	case *ContinueStmt:
		p.println("continue")

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

func (p *prettyPrinter) visitExprs(exprs []Expr) {
	for i, e := range exprs {
		Walk(p, e)
		if i < len(exprs)-1 {
			p.write(" ")
		}
	}
}

func Print(node Node) {
	Walk(&prettyPrinter{indent: 0, output: os.Stdout, indentSpace: "  "}, node)
}

func String(node Node) string {
	buf := &strings.Builder{}
	Walk(&prettyPrinter{indent: 0, output: buf, indentSpace: "  "}, node)
	return buf.String()
}
