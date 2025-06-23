package ast

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/hulo-lang/hulo/syntax/hulo/token"
)

// TODO: 只打印类型节点
type printer struct {}

type prettyPrinter struct {
	output io.Writer
	indent int
}

func (p *prettyPrinter) Visit(node Node) Visitor {
	if node == nil {
		p.indent--
		return nil
	}

	indentStr := strings.Repeat("  ", p.indent)

	switch n := node.(type) {
	case *File:
		for _, stmt := range n.Stmts {
			Walk(p, stmt)
		}
	case *FuncDecl:
		if n.Docs != nil {
			Walk(p, n.Docs)
		}

		for _, dec := range n.Decs {
			Walk(p, dec)
		}

		if n.Name != nil {
			Walk(p, n.Name)
		}

		for _, tp := range n.TypeParams {
			Walk(p, tp)
		}

		for _, param := range n.Recv {
			Walk(p, param)
		}

		for n.Type != nil {
			Walk(p, n.Type)
		}

		if n.Body != nil {
			Walk(p, n.Body)
		}

		return nil
	case *BlockStmt:
		return p
	default:
		fmt.Fprintf(p.output, "%s%T", indentStr, n)
	}
	fmt.Println()

	p.indent++
	return p
}

func Print(node Node) {
	// print(node, "")
	printer := prettyPrinter{output: os.Stdout}
	printer.Visit(node)
}

func String(node Node) string {
	buf := &strings.Builder{}
	printer := prettyPrinter{output: buf}
	printer.Visit(node)
	return buf.String()
}

func Write(reader io.Writer, node Node) {
	printer := prettyPrinter{output: reader}
	printer.Visit(node)
}

func print(node Node, indent string) {
	// fmt.Printf("%T\n", stmt)
	switch node := node.(type) {
	case *File:
		for _, decl := range node.Decls {
			print(decl, indent+"  ")
		}
		for _, ss := range node.Stmts {
			print(ss, indent+"  ")
		}
	case *Decorator:
		fmt.Printf("%s%s%s\n", indent, token.AT, node.Name)
	case *ExprStmt:
		fmt.Printf("%s%s\n", indent, node.X)
	case *TryStmt:
		fmt.Println(indent + "try {")
		for _, ss := range node.Body.List {
			print(ss, indent+"  ")
		}
		fmt.Print(indent + "}")
		for _, catch := range node.Catches {
			print(catch, indent)
		}
		if node.Finally != nil {
			print(node.Finally, indent)
		}
	case *CatchClause:
		fmt.Print("catch ")
		if node.Cond != nil {
			fmt.Printf("(%s)", node.Cond)
		}
		fmt.Println("{")
		for _, ss := range node.Body.List {
			print(ss, indent+"  ")
		}
		fmt.Println(indent + "}")
	case *FinallyStmt:
		fmt.Print("finally {")
		for _, ss := range node.Body.List {
			print(ss, indent+"  ")
		}
		fmt.Println(indent + "}")
	case *IfStmt:
		fmt.Printf("%sif %s {\n", indent, node.Cond)
		print(node.Body, indent+"  ")
		fmt.Println(indent + "}")

		for node.Else != nil {
			switch el := node.Else.(type) {
			case *IfStmt:
				fmt.Printf("%selse if %s {\n", indent, el.Cond)
				print(el.Body, indent+"  ")
				fmt.Println(indent + "}")
				node.Else = el.Else
			case *BlockStmt:
				fmt.Printf("%selse {\n", indent)
				print(el, indent+"  ")
				fmt.Println(indent + "}")
				node.Else = nil
			}
		}

	case *ThrowStmt:
		fmt.Printf("%s%s %s\n", indent, token.THROW, node.X)
	case *ReturnStmt:
		fmt.Printf("%s%s %s\n", indent, token.RETURN, node.X)
	case *AssignStmt:
		if node.Tok == token.COLON_ASSIGN {
			fmt.Printf("%s$%s := %s\n", indent, node.Lhs, node.Rhs)
		} else {
			if node.Type != nil {
				fmt.Printf("%s%s %s: %s = %s\n", indent, node.Scope, node.Lhs, node.Type, node.Rhs)
			} else {
				fmt.Printf("%s%s %s = %s\n", indent, node.Scope, node.Lhs, node.Rhs)
			}
		}
	case *BlockStmt:
		for _, ss := range node.List {
			print(ss, indent)
		}
	case *WhileStmt:
		fmt.Printf("%sloop {\n", indent)
		print(node.Body, indent+"  ")
		fmt.Println(indent + "}")
	case *DoWhileStmt:
		fmt.Printf("%sdo {\n", indent)
		print(node.Body, indent+"  ")
		fmt.Printf("%s} loop (%s)\n", indent, node.Cond)
	case *ForStmt:
		fmt.Printf("%sloop ", indent)
		if node.Init != nil {
			if as, ok := node.Init.(*AssignStmt); ok {
				fmt.Printf("%s := %s", as.Lhs, as.Rhs)
			} else {
				print(node.Init, "")
			}
		}
		fmt.Print("; ")
		if node.Cond != nil {
			fmt.Print(node.Cond)
		}
		fmt.Print("; ")
		if node.Post != nil {
			fmt.Print(node.Post)
		}
		fmt.Println(" {")
		print(node.Body, indent+"  ")
		fmt.Println(indent + "}")
	case *RangeStmt:
		fmt.Printf("%sloop %s in range(", indent, node.Index)
		if node.RangeExpr.Start != nil {
			fmt.Print(node.RangeExpr.Start)
		}
		fmt.Print(", ")
		if node.RangeExpr.End_ != nil {
			fmt.Print(node.RangeExpr.End_)
		}
		if node.RangeExpr.Step != nil {
			fmt.Printf(", %s", node.RangeExpr.Step)
		}
		fmt.Println(") {")
		print(node.Body, indent+"  ")
		fmt.Println(indent + "}")
	case *ForeachStmt:
		fmt.Printf("%sloop (", indent)
		if node.Index != nil {
			fmt.Print(node.Index)
		}
		if node.Value != nil {
			fmt.Printf(", %s", node.Value)
		}
		fmt.Printf(") in %s {\n", node.Var)
		print(node.Body, indent+"  ")
		fmt.Println(indent + "}")
	case *FuncDecl:
		fmt.Print(indent)
		// Print modifiers
		if node.Mod != nil {
			if node.Mod.Pub.IsValid() {
				fmt.Print("pub ")
			}
			if node.Mod.Const.IsValid() {
				fmt.Print("const ")
			}
			if node.Mod.Static.IsValid() {
				fmt.Print("static ")
			}
		}

		// Print function name
		fmt.Printf("fn %s", node.Name)

		if node.TypeParams != nil {
			expr := []string{}
			for _, e := range node.TypeParams {
				expr = append(expr, e.String())
			}
			fmt.Printf("<%s>", strings.Join(expr, ", "))
		}

		// Print parameters
		if node.Recv != nil && len(node.Recv) > 0 {
			fmt.Print("(")
			expr := []string{}
			for _, field := range node.Recv {
				expr = append(expr, field.String())
			}
			fmt.Print(strings.Join(expr, ", "))
			fmt.Print(")")
		} else {
			fmt.Print("()")
		}

		if node.Throw {
			fmt.Print(" throws")
		}

		if node.Type != nil {
			fmt.Printf(" -> %s", node.Type)
		} else {
			fmt.Print(" -> void")
		}

		// Print function body
		if node.Body != nil {
			fmt.Println(" {")
			for _, stmt := range node.Body.List {
				print(stmt, indent+"  ")
			}
			fmt.Print(indent + "}")
		} else {
			fmt.Print(" {}")
		}
		fmt.Println()
	case *ComptimeStmt:
		fmt.Println(indent + "comptime {")
		switch x := node.X.(type) {
		case *BlockStmt:
			print(x, indent+"  ")
		}
		fmt.Println(indent + "}")
	case *ModDecl:
		fmt.Printf("%smod %s {\n", indent, node.Name)
		for _, decl := range node.Body.List {
			print(decl, indent+"  ")
		}
		fmt.Println(indent + "}")
	case *UseDecl:
		if node.Rhs != nil {
			fmt.Printf("%suse %s = %s\n", indent, node.Lhs, node.Rhs)
		} else {
			fmt.Printf("%suse %s\n", indent, node.Lhs)
		}
	case *ExtensionDecl:
		switch {
		case node.ExtensionClass != nil:
			fmt.Printf("%sextension class %s {\n", indent, node.Name)
			if node.ExtensionClass.Body != nil {
				for _, decl := range node.ExtensionClass.Body.List {
					print(decl, indent+"  ")
				}
			}
			fmt.Printf("%s}\n", indent)
		}
	case *OperatorDecl:
		for _, dec := range node.Decs {
			fields := []string{}
			if dec.Recv != nil {
				for _, field := range dec.Recv.List {
					fields = append(fields, field.Name.String())
				}
			}
			fmt.Printf("%s@%s(%s)\n", indent, dec.Name, strings.Join(fields, ", "))
		}
		params := []string{}
		for _, p := range node.Params {
			params = append(params, p.String())
		}
		fmt.Printf("%soperator %s(%s)", indent, node.Name, strings.Join(params, ", "))
		if node.Type != nil {
			fmt.Printf(" -> %s", node.Type)
		}
		fmt.Printf(" {\n")
		if node.Body != nil {
			for _, decl := range node.Body.List {
				print(decl, indent+"  ")
			}
		}
		fmt.Printf("%s}\n", indent)
	case *CmdStmt:
		expr := []string{}
		for _, e := range node.Recv {
			expr = append(expr, e.String())
		}
		fmt.Printf("%s%s %s\n", indent, node.Name, strings.Join(expr, " "))
	case *ClassDecl:
		fmt.Printf("%sclass %s", indent, node.Name)
		if node.TypeParams != nil {
			expr := []string{}
			for _, e := range node.TypeParams {
				expr = append(expr, e.String())
			}
			fmt.Printf("<%s>", strings.Join(expr, ", "))
		}
		fmt.Printf(" {\n")

		for _, m := range node.Methods {
			print(m, indent+"  ")
		}
		fmt.Printf("%s}\n", indent)
	default:
		fmt.Printf("%T\n", node)
	}
}
