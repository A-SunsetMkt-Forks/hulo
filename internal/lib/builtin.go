package lib

import (
	"fmt"

	"github.com/hulo-lang/hulo/internal/object"
)

var Buitins = map[string]object.Value{
	"echo": object.BuiltinFunction(func(args ...object.Value) object.Value {
		for i := range args {
			args[i].Text()
		}
		fmt.Println()
		return nil
	}),
	"to_str": object.BuiltinFunction(func(args ...object.Value) object.Value {
		return nil
	}),
}
