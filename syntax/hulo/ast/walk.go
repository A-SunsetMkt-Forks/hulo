package ast

type Visitor interface {
	Visit(node Node) (w Visitor)
}

func Walk(v Visitor, node Node) {
	if node == nil {
		return
	}

	if v = v.Visit(node); v == nil {
		return
	}

	switch n := node.(type) {
	case *File:
		for _, stmt := range n.Stmts {
			Walk(v, stmt)
		}
	case *FuncDecl:
		if n.Docs != nil {
			Walk(v, n.Docs)
		}
		for _, dec := range n.Decs {
			Walk(v, dec)
		}

		if n.Name != nil {
			Walk(v, n.Name)
		}

		for _, tp := range n.TypeParams {
			Walk(v, tp)
		}

		for _, param := range n.Recv {
			Walk(v, param)
		}

		if n.Type != nil {
			Walk(v, n.Type)
		}

		if n.Body != nil {
			Walk(v, n.Body)
		}

	case *AssignStmt:
		if n.Docs != nil {
			Walk(v, n.Docs)
		}

		Walk(v, n.Lhs)

		if n.Type != nil {
			Walk(v, n.Type)
		}

		Walk(v, n.Rhs)
	}

	v.Visit(nil)
}
