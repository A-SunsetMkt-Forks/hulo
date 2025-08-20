// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package object

import (
	"testing"

	"github.com/hulo-lang/hulo/syntax/hulo/ast"
	"github.com/stretchr/testify/assert"
)

func TestBasicType(t *testing.T) {
	// assume there is a script as follows:
	// let a: str | num = 10
	// let b: str | num = "hello"
	// $a = $b
	converter := &ASTConverter{}
	a, err := converter.ConvertValue(&ast.NumericLiteral{Value: "10"})
	assert.NoError(t, err)
	b, err := converter.ConvertValue(&ast.StringLiteral{Value: "hello"})
	assert.NoError(t, err)
	assert.False(t, a.Type().AssignableTo(b.Type()))
}

func TestClassType(t *testing.T) {

}

func TestEnumType(t *testing.T) {

}

func TestGenericType(t *testing.T) {

}

func TestFunctionType(t *testing.T) {
}

func TestUtilityTypes(t *testing.T) {

}
