package parser

import (
	"github.com/antlr4-go/antlr/v4"
	"github.com/hulo-lang/hulo/syntax/hulo/parser/generated"
)

type Visitor struct {
	*generated.BasehuloParserVisitor
}

func accept[T any](tree antlr.ParseTree, visitor antlr.ParseTreeVisitor) (T, bool) {
	if tree == nil {
		return *new(T), false
	}
	t, ok := tree.Accept(visitor).(T)
	return t, ok
}
