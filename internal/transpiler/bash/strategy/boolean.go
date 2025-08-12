package strategy

import (
	bast "github.com/hulo-lang/hulo/syntax/bash/ast"
	hast "github.com/hulo-lang/hulo/syntax/hulo/ast"
)

// var (
// 	_ transpiler.Strategy[bast.Node] = (*BooleanAsNumberStrategy)(nil)
// 	_ transpiler.Strategy[bast.Node] = (*BooleanAsStringStrategy)(nil)
// 	_ transpiler.Strategy[bast.Node] = (*BooleanAsCommandStrategy)(nil)
// )

type BooleanAsNumberStrategy struct{}

/**
 * Transform boolean to number
 *
 * If the boolean is true, it will be transformed to number 1, otherwise it will be transformed to number 0
 *
 * For example:
 *
 *
 * if [ "$a" = "b" ]; then
 *   echo "a is b"
 * else
 *   echo "a is not b"
 * fi
 *
 *
 */
func (s *BooleanAsNumberStrategy) Name() string {
	return "number"
}

func (s *BooleanAsNumberStrategy) Apply(root hast.Node, node hast.Node) (bast.Node, error) {
	return nil, nil
}

// BooleanAsStringStrategy converts boolean values to strings
//
// If the boolean value is true, it will be converted to the string "true",
// otherwise it will be converted to the string "false"
//
// Example:
//
//	if [ "$a" = "true" ]; then
//	  echo "a is true"
//	else
//	  echo "a is false"
//	fi
type BooleanAsStringStrategy struct{}

func (s *BooleanAsStringStrategy) Name() string {
	return "string"
}

func (s *BooleanAsStringStrategy) Apply(root hast.Node, node hast.Node) (bast.Node, error) {
	return nil, nil
}

type BooleanAsCommandStrategy struct{}

func (s *BooleanAsCommandStrategy) Name() string {
	return "command"
}

func (s *BooleanAsCommandStrategy) Apply(root hast.Node, node hast.Node) (bast.Node, error) {
	return nil, nil
}

type NumberBoolEncoder struct{}

// BoolCodegen is the interface for the boolean codegen
type BoolCodegen interface {
	// RenderVal renders the value of the boolean
	//
	// Hulo:
	//  let a = true
	//  let b = false
	//
	// Target code:
	//  a=1
	//  b=0
	RenderVal(node hast.Expr) (bast.Node, error)
	// RenderExpr renders the expression of the boolean
	//
	// Hulo:
	//  if true {
	//		echo "this is true"
	//  } else {
	//		echo "this is false"
	//  }
	//
	// Target code:
	//  if [ "1" = "1" ]; then
	//    echo "this is true"
	//  else
	//    echo "this is false"
	//  fi
	RenderExpr(node hast.Expr) (bast.Node, error)
}
