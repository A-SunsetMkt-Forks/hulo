// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package astutil

import "github.com/hulo-lang/hulo/syntax/hulo/ast"

func Ident(name string) *ast.Ident {
	return &ast.Ident{Name: name}
}

func Num(val string) *ast.NumericLiteral {
	return &ast.NumericLiteral{Value: val}
}

func Str(val string) *ast.StringLiteral {
	return &ast.StringLiteral{Value: val}
}

func Bool(val bool) ast.Expr {
	if val {
		return &ast.TrueLiteral{}
	}
	return &ast.FalseLiteral{}
}

func Set(elems ...ast.Expr) *ast.SetLiteralExpr {
	return &ast.SetLiteralExpr{Elems: elems}
}

func Tuple(elems ...ast.Expr) *ast.TupleLiteralExpr {
	return &ast.TupleLiteralExpr{Elems: elems}
}

func NamedTuple(elems ...ast.Expr) *ast.NamedTupleLiteralExpr {
	return &ast.NamedTupleLiteralExpr{Elems: elems}
}

func Array(elems ...ast.Expr) *ast.ArrayLiteralExpr {
	return &ast.ArrayLiteralExpr{Elems: elems}
}

func Object(elems ...ast.Expr) *ast.ObjectLiteralExpr {
	return &ast.ObjectLiteralExpr{Props: elems}
}

func NamedObject(name string, elems ...ast.Expr) *ast.NamedObjectLiteralExpr {
	return &ast.NamedObjectLiteralExpr{Name: &ast.Ident{Name: name}, Props: elems}
}

var NULL = &ast.NullLiteral{}

func CallExpr(fun string, recv ...ast.Expr) *ast.CallExpr {
	return &ast.CallExpr{Fun: &ast.Ident{Name: fun}, Recv: recv}
}

func CmdExpr(cmd string, args ...ast.Expr) *ast.CmdExpr {
	return &ast.CmdExpr{Cmd: &ast.Ident{Name: cmd}, Args: args}
}
