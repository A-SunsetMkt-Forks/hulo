// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package ast

func Literal(val string) *Lit {
	return &Lit{Val: val}
}

func Words(parts ...Expr) *Word {
	return &Word{Parts: parts}
}

func CmdExpression(name string, recv ...Expr) *CmdExpr {
	return &CmdExpr{Name: Literal(name), Recv: recv}
}

func CallStatement(name string, recv ...Expr) *CallStmt {
	return &CallStmt{Name: name, Recv: recv}
}

func AssignStatement(lhs, rhs Expr) *AssignStmt {
	return &AssignStmt{Lhs: lhs, Rhs: rhs}
}

func SingleQuote(val Expr) *SglQuote {
	return &SglQuote{Val: val}
}

func DoubleQuote(val Expr) *DblQuote {
	return &DblQuote{Val: val}
}
