// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package parser

import (
	"github.com/hulo-lang/hulo/syntax/hulo/ast"
)

type SymbolInfo struct {
	Name string
	Node ast.Node
}

type SymbolTable struct{}

// TODO 变量提升？作用域？exports表和 bindings表
