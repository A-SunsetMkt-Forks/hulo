package comptime

import (
	"github.com/hulo-lang/hulo/internal/ast"
	"github.com/hulo-lang/hulo/internal/object"
)

func Eval(ctx *Context, node ast.Node) object.Value {
	switch node := node.(type) {
	case *ast.File:
	case *ast.FuncDecl:
		return evalFuncDecl(ctx, node)
	}
	return nil
}

func evalFuncDecl(ctx *Context, node *ast.FuncDecl) object.Value {
	return nil
}
