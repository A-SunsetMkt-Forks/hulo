package ast

import (
	"fmt"

	"github.com/hulo-lang/hulo/syntax/hulo/token"
)

func Print(node Node) {
	print(node, "")
}

func print(node Node, ident string) {
	// fmt.Printf("%T\n", stmt)
	switch node := node.(type) {
	case *File:
		for _, decl := range node.Decls {
			print(decl, ident+"  ")
		}
		for _, ss := range node.Stmts {
			print(ss, ident+"  ")
		}
	case *Annotation:
		fmt.Printf("%s%s%s\n", ident, token.AT, node.Name)
	case *ExprStmt:
		fmt.Printf("%s%s\n", ident, node.X)
	case *TryStmt:
		fmt.Println(ident + "try {")
		for _, ss := range node.Body.List {
			print(ss, ident+"  ")
		}
		fmt.Print(ident + "}")
		for _, catch := range node.Catches {
			print(catch, ident)
		}
		if node.Finally != nil {
			print(node.Finally, ident)
		}
	case *CatchStmt:
		fmt.Print("catch ")
		if node.Cond != nil {
			fmt.Printf("(%s)", node.Cond)
		}
		fmt.Println("{")
		for _, ss := range node.Body.List {
			print(ss, ident+"  ")
		}
		fmt.Println(ident + "}")
	case *FinallyStmt:
		fmt.Print("finally {")
		for _, ss := range node.Body.List {
			print(ss, ident+"  ")
		}
		fmt.Println(ident + "}")
	case *IfStmt:
		fmt.Printf("%sif %s {\n", ident, node.Cond)
		print(node.Body, ident+"  ")
		fmt.Println(ident + "}")

		for node.Else != nil {
			switch el := node.Else.(type) {
			case *IfStmt:
				fmt.Printf("%selse if %s {\n", ident, el.Cond)
				print(el.Body, ident+"  ")
				fmt.Println(ident + "}")
				node.Else = el.Else
			case *BlockStmt:
				fmt.Printf("%selse {\n", ident)
				print(el, ident+"  ")
				fmt.Println(ident + "}")
				node.Else = nil
			}
		}

	case *ThrowStmt:
		fmt.Printf("%s%s %s\n", ident, token.THROW, node.X)
	case *ReturnStmt:
		fmt.Printf("%s%s %s\n", ident, token.RETURN, node.X)
	case *AssignStmt:
		if node.Tok == token.COLON_ASSIGN {
			fmt.Printf("%s$%s := %s\n", ident, node.Lhs, node.Rhs)
		} else {
			fmt.Printf("%s%s %s = %s\n", ident, node.Scope, node.Lhs, node.Rhs)
		}
	case *BlockStmt:
		for _, ss := range node.List {
			print(ss, ident)
		}
	case *WhileStmt:
		fmt.Printf("%sloop {\n", ident)
		print(node.Body, ident+"  ")
		fmt.Println(ident + "}")
	case *DoWhileStmt:
		fmt.Printf("%sdo {\n", ident)
		print(node.Body, ident+"  ")
		fmt.Printf("%s} loop (%s)\n", ident, node.Cond)
	case *ForStmt:
		fmt.Printf("%sloop ", ident)
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
		print(node.Body, ident+"  ")
		fmt.Println(ident + "}")
	case *RangeStmt:
		fmt.Printf("%sloop %s in range(", ident, node.Index)
		if node.RangeClauseExpr.Start != nil {
			fmt.Print(node.RangeClauseExpr.Start)
		}
		fmt.Print(", ")
		if node.RangeClauseExpr.End != nil {
			fmt.Print(node.RangeClauseExpr.End)
		}
		if node.RangeClauseExpr.Step != nil {
			fmt.Printf(", %s", node.RangeClauseExpr.Step)
		}
		fmt.Println(") {")
		print(node.Body, ident+"  ")
		fmt.Println(ident + "}")
	case *ForeachStmt:
		fmt.Printf("%sloop (", ident)
		if node.Index != nil {
			fmt.Print(node.Index)
		}
		if node.Value != nil {
			fmt.Printf(", %s", node.Value)
		}
		fmt.Printf(") in %s {\n", node.Var)
		print(node.Body, ident+"  ")
		fmt.Println(ident + "}")
	case *FuncDecl:
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

		// Print parameters
		if node.Recv != nil && len(node.Recv.List) > 0 {
			fmt.Print("(")
			for i, field := range node.Recv.List {
				if i > 0 {
					fmt.Print(", ")
				}
				if field.Name != nil {
					fmt.Print(field.Name.Name)
				}
				if field.Type != nil {
					fmt.Printf(": %s", field.Type)
				}
				if field.Value != nil {
					fmt.Printf(" = %s", field.Value)
				}
			}
			fmt.Print(")")
		} else {
			fmt.Print("()")
		}

		// Print function body
		if node.Body != nil {
			fmt.Println(" {")
			for _, stmt := range node.Body.List {
				print(stmt, ident+"  ")
			}
			fmt.Print(ident + "}")
		} else {
			fmt.Print(" {}")
		}
		fmt.Println()
	case *ComptimeStmt:
		fmt.Println(ident + "comptime {")
		switch x := node.X.(type) {
		case *BlockStmt:
			print(x, ident+"  ")
		}
		fmt.Println(ident + "}")
	default:
		fmt.Printf("%T\n", node)
	}
}
