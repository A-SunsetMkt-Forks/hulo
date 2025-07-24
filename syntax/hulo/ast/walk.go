// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
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

	case *ClassDecl:
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
		if n.Parent != nil {
			Walk(v, n.Parent)
		}
		if n.Fields != nil {
			Walk(v, n.Fields)
		}
		for _, method := range n.Methods {
			Walk(v, method)
		}

	case *TypeParameter:
		if n.Name != nil {
			Walk(v, n.Name)
		}
		for _, constraint := range n.Constraints {
			Walk(v, constraint)
		}

	case *TypeReference:
		if n.Name != nil {
			Walk(v, n.Name)
		}
		for _, tp := range n.TypeParams {
			Walk(v, tp)
		}

	case *UnionType:
		for _, t := range n.Types {
			Walk(v, t)
		}

	case *IntersectionType:
		for _, t := range n.Types {
			Walk(v, t)
		}

	case *Parameter:
		if n.Name != nil {
			Walk(v, n.Name)
		}
		if n.Type != nil {
			Walk(v, n.Type)
		}
		if n.Value != nil {
			Walk(v, n.Value)
		}

	case *NamedParameters:
		for _, param := range n.Params {
			Walk(v, param)
		}

	case *NullableType:
		if n.X != nil {
			Walk(v, n.X)
		}

	case *AnyLiteral, *StringLiteral, *NumericLiteral,
		*TrueLiteral, *FalseLiteral, *Ident, *Comment:
		// 叶子节点，无需遍历子节点

	case *CommentGroup:
		for _, cmt := range n.List {
			Walk(v, cmt)
		}

	case *BlockStmt:
		for _, stmt := range n.List {
			Walk(v, stmt)
		}
	}

	v.Visit(nil)
}
