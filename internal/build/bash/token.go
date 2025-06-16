// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package build

import (
	btok "github.com/hulo-lang/hulo/syntax/bash/token"
	htok "github.com/hulo-lang/hulo/syntax/hulo/token"
)

var bashTokenMap = map[htok.Token]btok.Token{
	htok.EQ: btok.Equal,
	// htok.NE: btok.NE,
	htok.LT: btok.LessEqual,
	// htok.LE: btok.LE,
	htok.GT: btok.TsGtr,
	// htok.GE: btok.GE,
	htok.AND: btok.And,
	htok.OR:  btok.Or,
}

func Token(tok htok.Token) btok.Token {
	if t, ok := bashTokenMap[tok]; ok {
		return t
	}
	return btok.Illegal
}
