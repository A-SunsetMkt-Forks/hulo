// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package vbs

import (
	htok "github.com/hulo-lang/hulo/syntax/hulo/token"
	vtok "github.com/hulo-lang/hulo/syntax/vbs/token"
)

var vbsTokenMap = map[htok.Token]vtok.Token{
	htok.STR:      vtok.STRING,
	htok.NUM:      vtok.DOUBLE,
	htok.TRUE:     vtok.TRUE,
	htok.FALSE:    vtok.FALSE,
	htok.EQ:       vtok.EQ,
	htok.NEQ:      vtok.IEQ,
	htok.LT:       vtok.LT,
	htok.GT:       vtok.GT,
	htok.LE:       vtok.LTQ,
	htok.GE:       vtok.GTQ,
	htok.AND:      vtok.AND,
	htok.OR:       vtok.OR,
	htok.NOT:      vtok.NOT,
	htok.PLUS:     vtok.ADD,
	htok.MINUS:    vtok.SUB,
	htok.ASTERISK: vtok.MUL,
	htok.SLASH:    vtok.DIV,
}

func Token(tok htok.Token) vtok.Token {
	if t, ok := vbsTokenMap[tok]; ok {
		return t
	}
	return vtok.ILLEGAL
}
