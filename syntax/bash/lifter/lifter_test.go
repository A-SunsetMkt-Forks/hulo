// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package lifter

import (
	"testing"

	"github.com/hulo-lang/hulo/syntax/hulo/ast"
	"github.com/stretchr/testify/assert"
)

func TestLifter(t *testing.T) {
	lifter := NewLifter()
	script := `echo "Hello, World!"`
	node, err := lifter.Lift(script)
	assert.NoError(t, err)
	ast.Print(node)
}
