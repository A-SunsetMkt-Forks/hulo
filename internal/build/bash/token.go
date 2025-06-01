// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package build

import (
	btok "github.com/hulo-lang/hulo/syntax/bash/token"
	htok "github.com/hulo-lang/hulo/syntax/hulo/token"
)

var bashTokenMap = map[htok.Token]btok.Token{
	htok.EQ: btok.EQ,
	// htok.NE: btok.NE,
	htok.LT: btok.LT,
	// htok.LE: btok.LE,
	htok.GT: btok.GT,
	// htok.GE: btok.GE,
	htok.AND: btok.AND,
	htok.OR: btok.OR,
}

func Token(tok htok.Token) btok.Token {
	if t, ok := bashTokenMap[tok]; ok {
		return t
	}
	return btok.NONE
}
