// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package ast

// A Visitor's Visit method is invoked for each node encountered by Walk.
// If the result visitor w is not nil, Walk visits each of the children
// of node with the visitor w, followed by a call of w.Visit(nil).
type Visitor interface {
	Visit(node Node) (w Visitor)
}

func Walk(v Visitor, node Node) {
	if v = v.Visit(node); v == nil {
		return
	}

	// walk children
	// (the order of the cases matches the order
	// of the corresponding node types in ast.go)
	switch n := node.(type) {
	// Comments and fields
	case *Comment:
		// nothing to do
	case *CommentGroup:
		for _, c := range n.List {
			Walk(v, c)
		}
	case *Ident:
		// nothing to do
	case *DimDecl:
		Walk(v, n)

	case *ReDimDecl:
		Walk(v, n)

	case *ClassDecl:
		Walk(v, n)

	case *FuncDecl:
		Walk(v, n)

	case *PropertyDecl:
		Walk(v, n)

	case *SubDecl:
		Walk(v, n)

	case *AssignStmt:
		Walk(v, n)

	case *ExprStmt:
		Walk(v, n)

	case *MemberStmt:
		Walk(v, n)

	case *IfStmt:
		Walk(v, n)

	case *ForEachStmt:
		Walk(v, n)

	case *ForNextStmt:
		Walk(v, n)

	case *WhileWendStmt:
		Walk(v, n)

	case *DoLoopStmt:
		Walk(v, n)

	case *CallStmt:
		Walk(v, n)

	case *ExitStmt:
		Walk(v, n)

	case *SelectStmt:
		Walk(v, n)

	case *WithStmt:
		Walk(v, n)

	case *OnErrorStmt:
		Walk(v, n)

	case *StopStmt:
		Walk(v, n)

	case *RandomizeStmt:
		Walk(v, n)

	case *OptionStmt:
		Walk(v, n)
	}
}

func walkExprList(v Visitor, list []Expr) {
	for _, x := range list {
		Walk(v, x)
	}
}
