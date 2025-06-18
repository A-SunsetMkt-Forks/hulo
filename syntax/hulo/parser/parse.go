// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package parser

import (
	"errors"
	"time"

	"github.com/antlr4-go/antlr/v4"
	"github.com/caarlos0/log"
	"github.com/hulo-lang/hulo/syntax/hulo/ast"
	"github.com/hulo-lang/hulo/syntax/hulo/parser/generated"
)

type ParseOptions struct {
	Channel int
}

func ParseSourceFile(filename string, opt ParseOptions) (*ast.File, error) {
	stream, err := antlr.NewFileStream(filename)
	if err != nil {
		return nil, err
	}

	lex := generated.NewhuloLexer(stream)
	tokens := antlr.NewCommonTokenStream(lex, opt.Channel)
	parser := generated.NewhuloParser(tokens)
	tree := parser.File()
	visitor := &Visitor{}

	log.Infof("visiting %s", filename)
	start := time.Now()

	if file, ok := accept[*ast.File](tree, visitor); ok {
		file.Name = &ast.Ident{Name: filename}
		return file, nil
	}

	elapsed := time.Since(start)
	log.Infof("visited %s in %.2fs", filename, elapsed.Seconds())

	return nil, errors.New("fail to parse file")
}

func ParseSourceScript(input string, opt ParseOptions) (*ast.File, error) {
	stream := antlr.NewInputStream(input)

	log.Infof("parsing input")
	start := time.Now()

	lex := generated.NewhuloLexer(stream)
	tokens := antlr.NewCommonTokenStream(lex, opt.Channel)
	parser := generated.NewhuloParser(tokens)
	tree := parser.File()
	visitor := &Visitor{}
	// fmt.Println(tree.ToStringTree(nil, parser))

	if file, ok := accept[*ast.File](tree, visitor); ok {
		file.Name = &ast.Ident{Name: "string"}
		return file, nil
	}

	elapsed := time.Since(start)
	log.Infof("parsed input in %.2fs", elapsed.Seconds())

	return nil, errors.New("fail to parse input")
}
