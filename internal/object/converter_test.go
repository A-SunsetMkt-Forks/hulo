// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package object

import (
	"math/big"
	"testing"

	"github.com/hulo-lang/hulo/syntax/hulo/ast"
	"github.com/stretchr/testify/assert"
)

func TestASTConverter(t *testing.T) {
	converter := &ASTConverter{}

	tests := []struct {
		name string
		expr ast.Node
		want Value
	}{
		{
			name: "string literal",
			expr: &ast.StringLiteral{Value: "Hello, World!"},
			want: &StringValue{Value: "Hello, World!"},
		},
		{
			name: "number literal",
			expr: &ast.NumericLiteral{Value: "123"},
			want: &NumberValue{Value: big.NewFloat(123)},
		},
		{
			name: "true literal",
			expr: &ast.TrueLiteral{},
			want: TRUE,
		},
		{
			name: "false literal",
			expr: &ast.FalseLiteral{},
			want: FALSE,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := converter.ConvertValue(test.expr)
			assert.NoError(t, err)
			assert.Equal(t, test.want, got)
		})
	}
}
