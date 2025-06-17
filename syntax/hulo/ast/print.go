package ast

import (
	"fmt"

	"github.com/hulo-lang/hulo/syntax/hulo/token"
)

func Print(stmt Stmt) {
	print(stmt, "")
}

func print(stmt Stmt, ident string) {
	// fmt.Printf("%T\n", stmt)
	switch node := stmt.(type) {
	case *File:
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
			print(node.Post, "")
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
	default:
		fmt.Printf("%T\n", node)
	}
}
