package interpreter

import (
	"fmt"

	"github.com/hulo-lang/hulo/internal/object"
)

func echo(args ...object.Value) object.Value {
	for i, arg := range args {
		if i > 0 {
			fmt.Print(" ")
		}
		fmt.Print(arg.Text())
	}
	fmt.Println()
	return nil
}

var builtin = map[string]object.Value{
	"echo": object.NewFunctionValue(object.NewFunctionBuilder("echo").
		WithBuiltin(echo).
		WithVariadicParameter("args", object.GetAnyType()).
		WithReturnType(object.GetVoidType()).
		Build()),
	"print": object.NewFunctionValue(object.NewFunctionBuilder("print").
		WithBuiltin(echo).
		WithVariadicParameter("args", object.GetAnyType()).
		WithReturnType(object.GetVoidType()).
		Build()),
	"nop": object.NewFunctionValue(object.NewFunctionBuilder("nop").
		WithBuiltin(func(args ...object.Value) object.Value { return nil }).
		WithReturnType(object.GetVoidType()).
		Build()),
}
