package ast

import (
	"fmt"
	"strings"

	"github.com/hulo-lang/hulo/internal/token"
)

func Print(stmt Stmt) {
	print(stmt, "")
}

func print(stmt Stmt, ident string) {
	switch s := stmt.(type) {
	case *File:
		for _, ss := range s.Stmts {
			print(ss, ident+"  ")
		}
	case *Macro:
		fmt.Println(ident + "@" + identString(s.Name))
	case *ExprStmt:
		fmt.Println(ident + Expr2String(s.X))
	case *TryStmt:
		fmt.Println(ident + "try {")
		for _, ss := range s.Body.List {
			print(ss, ident+"  ")
		}
		fmt.Print(ident + "}")
		for _, catch := range s.Catches {
			print(catch, ident)
		}
		if s.Finally != nil {
			print(s.Finally, ident)
		}
	case *CatchStmt:
		fmt.Print("catch ")
		if s.Cond != nil {
			fmt.Printf("(%s)", Expr2String(s.Cond))
		}
		fmt.Println("{")
		for _, ss := range s.Body.List {
			print(ss, ident+"  ")
		}
		fmt.Println(ident + "}")
	case *FinallyStmt:
		fmt.Print("finally {")
		for _, ss := range s.Body.List {
			print(ss, ident+"  ")
		}
		fmt.Println(ident + "}")
	case *ThrowStmt:
		fmt.Println(ident+"throw", Expr2String(s.X))
	default:
		fmt.Printf("%T\n", s)
	}
}

func Expr2String(expr Expr) string {
	switch e := expr.(type) {
	case *Ident:
		return e.Name
	case *BasicLit:
		return e.Value
	case *BinaryExpr:
		switch e.Op {
		case token.COLON:
			return fmt.Sprintf("%s: %s", Expr2String(e.X), Expr2String(e.Y))
		default:
			return fmt.Sprintf("%s%s%s", Expr2String(e.X), e.Op, Expr2String(e.Y))
		}
	case *CallExpr:
		args := strings.Join(getExprs(e.Recv), ", ")
		return fmt.Sprintf("%s(%s)", Expr2String(e.Fun), args)
	case *NewDelExpr:
		return fmt.Sprintf("new %s", Expr2String(e.X))
	case *TypeAnnotation:
		res := []string{}
		for _, t := range e.List {
			res = append(res, Expr2String(t))
		}
		return fmt.Sprint(strings.Join(res, " | "))
	default:
		fmt.Printf("%T\n", e)
	}
	return ""
}

func getExprs(es []Expr) []string {
	ret := []string{}
	for _, e := range es {
		ret = append(ret, Expr2String(e))
	}
	return ret
}

func identString(e Expr) string {
	return e.(*Ident).Name
}
