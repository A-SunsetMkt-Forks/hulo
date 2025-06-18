// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package parser

import "github.com/hulo-lang/hulo/syntax/hulo/ast"

type SymbolInfo struct {
	Name     string
	Declared bool
	Node     ast.Node
}
