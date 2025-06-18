package ast

import (
	"fmt"
	"io"
	"os"

	"github.com/hulo-lang/hulo/syntax/batch/token"
)

func Print(node Node) {
	print(node, "", os.Stdout)
}

func Write(node Node, w io.Writer) {
	print(node, "", w)
}

func print(node Node, ident string, w io.Writer) {
	switch n := node.(type) {
	case *File:
		for _, doc := range n.Docs {
			print(doc, ident, w)
		}

		for _, stmt := range n.Stmts {
			print(stmt, ident, w)
		}
	case *IfStmt:
		fmt.Fprintf(w, "%sif %s (\n", ident, n.Cond)
		print(n.Body, ident+"  ", w)
		fmt.Print(ident + ")")
		for n.Else != nil {
			fmt.Print(" ")
			switch el := n.Else.(type) {
			case *IfStmt:
				fmt.Fprintf(w, "%selse if %s (\n", ident, el.Cond)
				print(el.Body, ident+"  ", w)
				fmt.Fprintf(w, "%s)\n", ident)
				n.Else = el.Else
			case *BlockStmt:
				fmt.Fprintf(w, "%selse (\n", ident)
				print(el, ident+"  ", w)
				fmt.Fprintf(w, "%s)\n", ident)
				n.Else = nil
			}
		}
		fmt.Fprintln(w)

	case *ForStmt:
		fmt.Fprintf(w, "%sfor %s in %s do (\n", ident, n.X, n.List)
		print(n.Body, ident+"  ", w)
		fmt.Fprintf(w, "%s)\n", ident)
	case *ExprStmt:
		fmt.Fprintf(w, "%s%s\n", ident, n.X)
	case *AssignStmt:
		fmt.Fprintf(w, "%sset %s=%s\n", ident, n.Lhs, n.Rhs)
	case *CallStmt:
		fmt.Fprintf(w, "%scall :%s", ident, n.Name)
		if len(n.Recv) > 0 {
			fmt.Print(" ")
			for i, recv := range n.Recv {
				if i > 0 {
					fmt.Print(" ")
				}
				fmt.Print(recv)
			}
		}
		fmt.Fprintln(w)
	case *FuncDecl:
		fmt.Fprintf(w, "%s:%s\n", ident, n.Name)
		print(n.Body, ident+"  ", w)
	case *BlockStmt:
		for _, stmt := range n.List {
			print(stmt, ident, w)
		}
	case *GotoStmt:
		fmt.Fprintf(w, "%sgoto :%s\n", ident, n.Label)
	case *LabelStmt:
		fmt.Fprintf(w, "%s:%s\n", ident, n.Name)
	case *Command:
		fmt.Fprintf(w, "%s%s", ident, n.Name)
		if len(n.Recv) > 0 {
			fmt.Fprint(w, " ")
			for i, recv := range n.Recv {
				if i > 0 {
					fmt.Fprint(w, " ")
				}
				fmt.Fprint(w, recv)
			}
		}
		fmt.Fprintln(w)
	case *Comment:
		if n.Tok == token.DOUBLE_COLON {
			fmt.Fprintf(w, "%s::%s\n", ident, n.Text)
		} else {
			fmt.Fprintf(w, "%sREM %s\n", ident, n.Text)
		}
	case *CommentGroup:
		for _, comment := range n.Comments {
			print(comment, ident, w)
		}
	default:
		fmt.Fprintf(w, "%T\n", node)
	}
}
