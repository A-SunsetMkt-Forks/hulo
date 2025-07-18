package transpiler

import (
	hast "github.com/hulo-lang/hulo/syntax/hulo/ast"
)

type Strategy[T any] interface {
	Name() string
	Apply(root hast.Node, node hast.Node) (T, error)
}
