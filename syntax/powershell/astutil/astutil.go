// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package astutil

import (
	"github.com/hulo-lang/hulo/syntax/powershell/ast"
	"github.com/hulo-lang/hulo/syntax/powershell/token"
)

func BinaryExpr(lhs ast.Expr, op token.Token, rhs ast.Expr) *ast.BinaryExpr {
	return &ast.BinaryExpr{X: lhs, Op: op, Y: rhs}
}

func Ident(name string) *ast.Ident {
	return &ast.Ident{Name: name}
}

func Lit(s string) *ast.Lit {
	return &ast.Lit{Val: s}
}

func StringLit(s string) *ast.StringLit {
	return &ast.StringLit{Val: s}
}

func NumericLit(s string) *ast.NumericLit {
	return &ast.NumericLit{Val: s}
}

func BoolLit(b bool) *ast.BoolLit {
	return &ast.BoolLit{Val: b}
}

func VarExpr(x ast.Expr) *ast.VarExpr {
	return &ast.VarExpr{X: x}
}

func BlockExpr(list ...ast.Expr) *ast.BlockExpr {
	return &ast.BlockExpr{List: list}
}
