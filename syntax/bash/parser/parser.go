// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package parser

import (
	"bytes"

	"github.com/hulo-lang/hulo/syntax/bash/ast"
	"mvdan.cc/sh/v3/syntax"
)

type Parser struct {
	parser *syntax.Parser
}

func NewParser() *Parser {
	return &Parser{
		parser: syntax.NewParser(),
	}
}

func (p *Parser) Parse(script string) (*ast.File, error) {
	file, err := p.parser.Parse(bytes.NewReader([]byte(script)), "script.sh")
	if err != nil {
		return nil, err
	}
	return p.convert(file), nil
}

func (p *Parser) convert(file *syntax.File) *ast.File {
	ret := &ast.File{}

	syntax.Walk(file, func(node syntax.Node) bool {
		switch node := node.(type) {
		case *syntax.Stmt:
			return true
		case *syntax.CallExpr:
			cmd := &ast.CmdExpr{}
			for i, args := range node.Args {
				var ident string
				for _, part := range args.Parts {
					switch part := part.(type) {
					case *syntax.Lit:
						ident = part.Value
					case *syntax.SglQuoted:
						ident = part.Value
					case *syntax.DblQuoted:
						for _, part := range part.Parts {
							switch part := part.(type) {
							case *syntax.Lit:
								ident = part.Value
							}
						}
					}
				}
				if i == 0 {
					cmd.Name = &ast.Ident{Name: ident}
				} else {
					cmd.Recv = append(cmd.Recv, &ast.Word{Val: ident})
				}
			}
			ret.Stmts = append(ret.Stmts, &ast.ExprStmt{X: cmd})
			return true
		}
		return true
	})

	return ret
}
