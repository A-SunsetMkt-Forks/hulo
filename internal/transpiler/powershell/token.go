// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package pwsh

import (
	htok "github.com/hulo-lang/hulo/syntax/hulo/token"
	pstok "github.com/hulo-lang/hulo/syntax/powershell/token"
)

var psTokenMap = map[htok.Token]pstok.Token{
	// 比较运算符
	htok.EQ:  pstok.EQ,
	htok.NEQ: pstok.NE,
	htok.LT:  pstok.LT,
	htok.LE:  pstok.LE,
	htok.GT:  pstok.GT,
	htok.GE:  pstok.GE,
	// 逻辑运算符
	htok.AND: pstok.AND,
	htok.OR:  pstok.OR,
	// 算术运算符
	htok.PLUS:     pstok.ADD,
	htok.MINUS:    pstok.SUB,
	htok.ASTERISK: pstok.MUL,
	htok.SLASH:    pstok.DIV,
	// 赋值
	htok.ASSIGN: pstok.ASSIGN,
	htok.DEC:    pstok.DEC,
	htok.INC:    pstok.INC,
}

func psToken(tok htok.Token) pstok.Token {
	if t, ok := psTokenMap[tok]; ok {
		return t
	}
	return 0 // 默认返回 0
}
