package parser

import (
	"errors"

	"github.com/antlr4-go/antlr/v4"
	"github.com/caarlos0/log"
	"github.com/hulo-lang/hulo/syntax/hulo/ast"
)

type ParseOptions struct {
	Channel int
}

func ParseSourceFile(filename string, opt ParseOptions) (*ast.File, error) {
	stream, err := antlr.NewFileStream(filename)
	if err != nil {
		return nil, err
	}

	lex := NewhuloLexer(stream)
	tokens := antlr.NewCommonTokenStream(lex, opt.Channel)
	parser := NewhuloParser(tokens)
	tree := parser.File()
	visitor := &Visitor{}

	log.Infof("visiting %s", filename)

	if file, ok := accept[*ast.File](tree, visitor); ok {
		file.Name = &ast.Ident{Name: filename}
		return file, nil
	}

	return nil, errors.New("fail to parse file")
}
