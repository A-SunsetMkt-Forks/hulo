// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package parser

import (
	"errors"
	"io"
	"time"

	"github.com/antlr4-go/antlr/v4"
	"github.com/caarlos0/log"
	"github.com/hulo-lang/hulo/syntax/hulo/ast"
)

var (
	ErrFailToParseFile   = errors.New("fail to parse file")
	ErrFailToParseInput  = errors.New("fail to parse input")
	ErrFailToParseReader = errors.New("fail to parse reader")
)

func ParseSourceFile(filename string, opts ...ParserOptions) (*ast.File, error) {
	stream, err := antlr.NewFileStream(filename)
	if err != nil {
		return nil, err
	}

	analyzer, err := NewAnalyzer(stream, opts...)
	if err != nil {
		return nil, err
	}

	log.Debugf("visiting %s", filename)
	start := time.Now()

	if file, ok := accept[*ast.File](analyzer.file, analyzer); ok {
		file.Name = &ast.Ident{Name: filename}
		return file, nil
	}

	elapsed := time.Since(start)
	log.Debugf("visited %s in %.2fs", filename, elapsed.Seconds())

	return nil, ErrFailToParseFile
}

func ParseSourceScript(input string, opts ...ParserOptions) (*ast.File, error) {
	stream := antlr.NewInputStream(input)

	log.Debug("parsing input")
	start := time.Now()

	analyzer, err := NewAnalyzer(stream, opts...)
	if err != nil {
		return nil, err
	}

	if file, ok := accept[*ast.File](analyzer.file, analyzer); ok {
		file.Name = &ast.Ident{Name: "<source>"}
		return file, nil
	}

	elapsed := time.Since(start)
	log.Debugf("parsed input in %.2fs", elapsed.Seconds())

	return nil, ErrFailToParseInput
}

func ParseReader(reader io.Reader, opts ...ParserOptions) (*ast.File, error) {
	stream := antlr.NewIoStream(reader)

	log.Debug("parsing reader")
	start := time.Now()

	analyzer, err := NewAnalyzer(stream, opts...)
	if err != nil {
		return nil, err
	}

	if file, ok := accept[*ast.File](analyzer.file, analyzer); ok {
		file.Name = &ast.Ident{Name: "<reader>"}
		return file, nil
	}

	elapsed := time.Since(start)
	log.Debugf("parsed reader in %.2fs", elapsed.Seconds())

	return nil, ErrFailToParseReader
}
