// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package ast

import "github.com/hulo-lang/hulo/syntax/bash/token"

func Literal(v string) *Lit {
	return &Lit{Val: v}
}

func Identifier(name string) *Ident {
	return &Ident{Name: name}
}

func BinaryExpression(x Expr, op token.Token, y Expr) *BinaryExpr {
	return &BinaryExpr{
		X:  x,
		Op: op,
		Y:  y,
	}
}

func CallExpression(name string, args ...Expr) *CallExpr {
	return &CallExpr{
		Func: Identifier(name),
		Recv: args,
	}
}
