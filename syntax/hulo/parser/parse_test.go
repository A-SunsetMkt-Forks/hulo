package parser

import (
	"fmt"
	"testing"

	"github.com/antlr4-go/antlr/v4"
	"github.com/hulo-lang/hulo/syntax/hulo/parser/generated"
	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	stream, err := antlr.NewFileStream("./testdata/bool.hl")
	assert.NoError(t, err)

	lexer := generated.NewhuloLexer(stream)
	tokens := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	parser := generated.NewhuloParser(tokens)
	fmt.Println(parser.File().ToStringTree(nil, parser))
}
