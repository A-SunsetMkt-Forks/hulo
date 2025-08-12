// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package batch

import (
	btok "github.com/hulo-lang/hulo/syntax/batch/token"
	htok "github.com/hulo-lang/hulo/syntax/hulo/token"
)

var batchTokenMap = map[htok.Token]btok.Token{
	htok.LT:       btok.LSS,
	htok.GT:       btok.GTR,
	htok.LE:       btok.LEQ,
	htok.GE:       btok.GEQ,
	htok.EQ:       btok.EQU,
	htok.NEQ:      btok.NEQ,
	htok.AND:      btok.AND,
	htok.OR:       btok.OR,
	htok.NOT:      btok.NOT,
	htok.PLUS:     btok.ADD,
	htok.MINUS:    btok.SUB,
	htok.ASTERISK: btok.MUL,
	htok.SLASH:    btok.DIV,
	htok.MOD:      btok.MOD,
}

func Token(tok htok.Token) btok.Token {
	if t, ok := batchTokenMap[tok]; ok {
		return t
	}
	return btok.ILLEGAL
}
