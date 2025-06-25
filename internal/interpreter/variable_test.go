package interpreter

import (
	"testing"

	"github.com/hulo-lang/hulo/internal/object"
	"github.com/hulo-lang/hulo/syntax/hulo/ast"
	"github.com/hulo-lang/hulo/syntax/hulo/token"
	"github.com/stretchr/testify/assert"
)

func TestVariableAssignment(t *testing.T) {
	// let a: str | num = 10; let b: str | num = "hello"; $a = $b
	unionType := object.NewUnionType(object.NewNumberType(), object.NewStringType())

	a := &Variable{
		Name:     "a",
		DeclType: unionType,
		Value:    object.NewNumberValue("10"),
		Scope:    token.LET,
	}

	b := &Variable{
		Name:     "b",
		DeclType: unionType,
		Scope:    token.LET,
	}

	converter := &object.ASTConverter{}
	aValue, err := converter.ConvertValue(&ast.NumericLiteral{Value: "10"})
	assert.NoError(t, err)
	bValue, err := converter.ConvertValue(&ast.StringLiteral{Value: "hello"})
	assert.NoError(t, err)

	err = a.Assign(aValue)
    assert.NoError(t, err)

	err = b.Assign(bValue)
    assert.NoError(t, err)

	err = a.Assign(b.Value)
	assert.NoError(t, err)
}
