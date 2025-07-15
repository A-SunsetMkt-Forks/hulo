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
	if node == nil {
		return
	}

	if v = v.Visit(node); v == nil {
		return
	}

	// walk children
	// (the order of the cases matches the order
	// of the corresponding node types in ast.go)
	switch n := node.(type) {
	// File and comments
	case *File:
		for _, doc := range n.Docs {
			Walk(v, doc)
		}
		for _, stmt := range n.Stmts {
			Walk(v, stmt)
		}
	case *CommentGroup:
		for _, comment := range n.Comments {
			Walk(v, comment)
		}
	case *Comment:
		// nothing to do - leaf node

	// Statements
	case *ExprStmt:
		if n.Docs != nil {
			Walk(v, n.Docs)
		}
		Walk(v, n.X)
	case *IfStmt:
		if n.Docs != nil {
			Walk(v, n.Docs)
		}
		Walk(v, n.Cond)
		Walk(v, n.Body)
		if n.Else != nil {
			Walk(v, n.Else)
		}
	case *ForStmt:
		if n.Docs != nil {
			Walk(v, n.Docs)
		}
		Walk(v, n.X)
		Walk(v, n.List)
		Walk(v, n.Body)
	case *AssignStmt:
		if n.Docs != nil {
			Walk(v, n.Docs)
		}
		Walk(v, n.Lhs)
		Walk(v, n.Rhs)
	case *CallStmt:
		if n.Docs != nil {
			Walk(v, n.Docs)
		}
		walkExprList(v, n.Recv)
	case *FuncDecl:
		if n.Docs != nil {
			Walk(v, n.Docs)
		}
		Walk(v, n.Body)
	case *BlockStmt:
		for _, stmt := range n.List {
			Walk(v, stmt)
		}
	case *GotoStmt:
		if n.Docs != nil {
			Walk(v, n.Docs)
		}
		// nothing else to walk - label is a string
	case *LabelStmt:
		if n.Docs != nil {
			Walk(v, n.Docs)
		}
		// nothing else to walk - name is a string

	// Expressions
	case *Word:
		for _, part := range n.Parts {
			Walk(v, part)
		}
	case *UnaryExpr:
		Walk(v, n.X)
	case *BinaryExpr:
		Walk(v, n.X)
		Walk(v, n.Y)
	case *CallExpr:
		Walk(v, n.Fun)
		walkExprList(v, n.Recv)
	case *Lit:
		// nothing to do - leaf node
	case *SglQuote:
		Walk(v, n.Val)
	case *DblQuote:
		Walk(v, n.Val)
	case *CmdExpr:
		if n.Docs != nil {
			Walk(v, n.Docs)
		}
		Walk(v, n.Name)
		walkExprList(v, n.Recv)
	}

	v.Visit(nil)
}

// walkExprList walks a list of expressions.
func walkExprList(v Visitor, list []Expr) {
	for _, x := range list {
		Walk(v, x)
	}
}
