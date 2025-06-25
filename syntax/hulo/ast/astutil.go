// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package ast

func Identifier(name string) *Ident {
	return &Ident{Name: name}
}

func Num(val string) *NumericLiteral {
	return &NumericLiteral{Value: val}
}

func Str(val string) *StringLiteral {
	return &StringLiteral{Value: val}
}

func Bool(val bool) Expr {
	if val {
		return &TrueLiteral{}
	}
	return &FalseLiteral{}
}
