package interpreter

import (
	"fmt"

	"github.com/hulo-lang/hulo/internal/object"
	"github.com/hulo-lang/hulo/syntax/hulo/ast"
	"github.com/hulo-lang/hulo/syntax/hulo/token"
)

type Variable struct {
	Name     string
	DeclType object.Type
	Value    object.Value
	Scope    token.Token
	Node     ast.Node
}

func (v *Variable) Assign(newValue object.Value) error {
	if !v.CanAssign(newValue) {
		return fmt.Errorf("cannot assign %s to variable %s of type %s",
			newValue.Type().Name(), v.Name, v.DeclType.Name())
	}

	if v.Scope == token.CONST && v.Value != nil {
		return fmt.Errorf("cannot reassign const variable %s", v.Name)
	}

	v.Value = newValue
	return nil
}

func (v *Variable) CanAssign(newValue object.Value) bool {
	return newValue.Type().AssignableTo(v.DeclType)
}
