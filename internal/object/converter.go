// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package object

import (
	"fmt"

	"github.com/hulo-lang/hulo/syntax/hulo/ast"
)

type Converter interface {
	ConvertType(node ast.Node) (Type, error)
	ConvertValue(node ast.Node) (Value, error)
}

type ASTConverter struct {
}

func (c *ASTConverter) ConvertValue(node ast.Node) (Value, error) {
	switch n := node.(type) {
	case *ast.StringLiteral:
		return &StringValue{Value: n.Value}, nil
	case *ast.NumericLiteral:
		return NewNumberValue(n.Value), nil
	case *ast.TrueLiteral:
		return TRUE, nil
	case *ast.FalseLiteral:
		return FALSE, nil
	case *ast.NullLiteral:
		return NULL, nil
	case *ast.ArrayType:
	case *ast.ObjectLiteralExpr:
	case *ast.ArrayLiteralExpr:
	case *ast.SetLiteralExpr:
	default:
		return nil, fmt.Errorf("unsupported type: %T", n)
	}
	return nil, nil
}
