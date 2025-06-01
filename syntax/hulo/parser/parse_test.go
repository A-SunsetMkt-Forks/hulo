package parser

import (
	"fmt"
	"testing"

	"github.com/antlr4-go/antlr/v4"
	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	stream, err := antlr.NewFileStream("./testdata/bool.hl")
	assert.NoError(t, err)

	lexer := NewhuloLexer(stream)
	tokens := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	parser := NewhuloParser(tokens)
	fmt.Println(parser.File().ToStringTree(nil, parser))
}
