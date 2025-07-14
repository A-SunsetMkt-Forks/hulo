// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package ast

import (
	"strconv"

	"github.com/hulo-lang/hulo/syntax/bash/token"
)

func Literal(v string) *Word {
	return &Word{Val: v}
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

func CmdExpression(name string, args ...Expr) *CmdExpr {
	return &CmdExpr{
		Name: Identifier(name),
		Recv: args,
	}
}

func Option(name string) *Ident {
	return &Ident{Name: name}
}

/// Bourne Shell Builtins

func NoOpCmd(args ...Expr) *CmdExpr {
	return CmdExpression(":", args...)
}

func SourceCmd(file string) *CmdExpr {
	return CmdExpression("source", Literal(file))
}

func Break(n ...int) *CmdExpr {
	if len(n) == 0 {
		return CmdExpression("break")
	}
	return CmdExpression("break", Literal(strconv.Itoa(n[0])))
}

func Continue(n ...int) *CmdExpr {
	if len(n) == 0 {
		return CmdExpression("continue")
	}
	return CmdExpression("continue", Literal(strconv.Itoa(n[0])))
}

func Exit(n ...int) *CmdExpr {
	if len(n) == 0 {
		return CmdExpression("exit")
	}
	return CmdExpression("exit", Literal(strconv.Itoa(n[0])))
}

func ExecCmd(args ...Expr) *CmdExpr {
	return CmdExpression("exec", args...)
}

func True() *CmdExpr {
	return CmdExpression("true")
}

func False() *CmdExpr {
	return CmdExpression("false")
}
