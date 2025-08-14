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
	case *Field:
		if n.Name != nil {
			Walk(v, n.Name)
		}

	// Declarations
	case *DimDecl:
		walkExprList(v, n.List)
		if n.Set != nil {
			Walk(v, n.Set)
		}
	case *ReDimDecl:
		walkExprList(v, n.List)
	case *ClassDecl:
		if n.Name != nil {
			Walk(v, n.Name)
		}
		for _, stmt := range n.Stmts {
			Walk(v, stmt)
		}
	case *FuncDecl:
		if n.Name != nil {
			Walk(v, n.Name)
		}
		for _, param := range n.Recv {
			Walk(v, param)
		}
		if n.Body != nil {
			Walk(v, n.Body)
		}
	case *SubDecl:
		if n.Name != nil {
			Walk(v, n.Name)
		}
		for _, param := range n.Recv {
			Walk(v, param)
		}
		if n.Body != nil {
			Walk(v, n.Body)
		}
	case *PropertyGetStmt:
		if n.Name != nil {
			Walk(v, n.Name)
		}
		for _, param := range n.Recv {
			Walk(v, param)
		}
		if n.Body != nil {
			Walk(v, n.Body)
		}
	case *PropertySetStmt:
		if n.Name != nil {
			Walk(v, n.Name)
		}
		for _, param := range n.Recv {
			Walk(v, param)
		}
		if n.Body != nil {
			Walk(v, n.Body)
		}
	case *PropertyLetStmt:
		if n.Name != nil {
			Walk(v, n.Name)
		}
		for _, param := range n.Recv {
			Walk(v, param)
		}
		if n.Body != nil {
			Walk(v, n.Body)
		}

	// Statements
	case *AssignStmt:
		if n.Lhs != nil {
			Walk(v, n.Lhs)
		}
		if n.Rhs != nil {
			Walk(v, n.Rhs)
		}
	case *SetStmt:
		if n.Lhs != nil {
			Walk(v, n.Lhs)
		}
		if n.Rhs != nil {
			Walk(v, n.Rhs)
		}
	case *ConstStmt:
		if n.Lhs != nil {
			Walk(v, n.Lhs)
		}
		if n.Rhs != nil {
			Walk(v, n.Rhs)
		}
	case *PrivateStmt:
		walkExprList(v, n.List)
	case *PublicStmt:
		walkExprList(v, n.List)
	case *ExprStmt:
		if n.X != nil {
			Walk(v, n.X)
		}
	case *IfStmt:
		if n.Cond != nil {
			Walk(v, n.Cond)
		}
		if n.Body != nil {
			Walk(v, n.Body)
		}
		if n.Else != nil {
			Walk(v, n.Else)
		}
	case *BlockStmt:
		for _, stmt := range n.List {
			Walk(v, stmt)
		}
	case *ForEachStmt:
		if n.Elem != nil {
			Walk(v, n.Elem)
		}
		if n.Group != nil {
			Walk(v, n.Group)
		}
		if n.Body != nil {
			Walk(v, n.Body)
		}
		if n.Stmt != nil {
			Walk(v, n.Stmt)
		}
	case *ForNextStmt:
		if n.Start != nil {
			Walk(v, n.Start)
		}
		if n.End_ != nil {
			Walk(v, n.End_)
		}
		if n.Step != nil {
			Walk(v, n.Step)
		}
		if n.Body != nil {
			Walk(v, n.Body)
		}
	case *WhileWendStmt:
		if n.Cond != nil {
			Walk(v, n.Cond)
		}
		if n.Body != nil {
			Walk(v, n.Body)
		}
	case *DoLoopStmt:
		if n.Cond != nil {
			Walk(v, n.Cond)
		}
		if n.Body != nil {
			Walk(v, n.Body)
		}
	case *CallStmt:
		if n.Name != nil {
			Walk(v, n.Name)
		}
		walkExprList(v, n.Recv)
	case *ExitStmt:
		// nothing to do
	case *SelectStmt:
		if n.Var != nil {
			Walk(v, n.Var)
		}
		for _, case_ := range n.Cases {
			if case_.Cond != nil {
				Walk(v, case_.Cond)
			}
			if case_.Body != nil {
				Walk(v, case_.Body)
			}
		}
		if n.Else != nil {
			if n.Else.Cond != nil {
				Walk(v, n.Else.Cond)
			}
			if n.Else.Body != nil {
				Walk(v, n.Else.Body)
			}
		}
	case *WithStmt:
		if n.Cond != nil {
			Walk(v, n.Cond)
		}
		if n.Body != nil {
			Walk(v, n.Body)
		}
	case *OnErrorStmt:
		// OnErrorGoto and OnErrorResume are not nodes, just positions
	case *StopStmt:
		// nothing to do
	case *RandomizeStmt:
		// nothing to do
	case *OptionStmt:
		// nothing to do
	case *EraseStmt:
		if n.X != nil {
			Walk(v, n.X)
		}
	case *ExecuteStmt:
		// nothing to do

	// Expressions
	case *Ident:
		// nothing to do
	case *BasicLit:
		// nothing to do
	case *SelectorExpr:
		if n.X != nil {
			Walk(v, n.X)
		}
		if n.Sel != nil {
			Walk(v, n.Sel)
		}
	case *BinaryExpr:
		if n.X != nil {
			Walk(v, n.X)
		}
		if n.Y != nil {
			Walk(v, n.Y)
		}
	case *CallExpr:
		if n.Func != nil {
			Walk(v, n.Func)
		}
		walkExprList(v, n.Recv)
	case *CmdExpr:
		if n.Cmd != nil {
			Walk(v, n.Cmd)
		}
		walkExprList(v, n.Recv)
	case *IndexExpr:
		if n.X != nil {
			Walk(v, n.X)
		}
		if n.Index != nil {
			Walk(v, n.Index)
		}
	case *IndexListExpr:
		if n.X != nil {
			Walk(v, n.X)
		}
		walkExprList(v, n.Indices)
	case *NewExpr:
		if n.X != nil {
			Walk(v, n.X)
		}

	// Files
	case *File:
		for _, stmt := range n.Stmts {
			Walk(v, stmt)
		}
	}

	v.Visit(nil)
}

func walkExprList(v Visitor, list []Expr) {
	for _, x := range list {
		Walk(v, x)
	}
}

type predicate func(Node) bool

func (f predicate) Visit(node Node) Visitor {
	if f(node) {
		return f
	}
	return nil
}

// WalkIf walks the AST if the predicate returns true.
func WalkIf(node Node, f predicate) {
	Walk(predicate(f), node)
}
