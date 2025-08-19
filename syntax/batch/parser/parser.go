// syntax/batch/parser/parser.go
// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package parser

import (
	"fmt"
	"github.com/hulo-lang/hulo/syntax/batch/ast"
	"github.com/hulo-lang/hulo/syntax/batch/lexer"
	"github.com/hulo-lang/hulo/syntax/batch/token"
)

type Parser struct {
	l         *lexer.Lexer
	curToken  token.TokenInfo
	peekToken token.TokenInfo
	errors    []string
}

func NewParser(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *ast.File {
	program := &ast.File{}

	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Stmts = append(program.Stmts, stmt)
		}
		p.nextToken()
	}

	return program
}

func (p *Parser) parseStatement() ast.Stmt {
	switch p.curToken.Type {
	case token.IF:
		return p.parseIfStatement()
	case token.FOR:
		return p.parseForStatement()
	case token.CALL:
		return p.parseCallStatement()
	case token.GOTO:
		return p.parseGotoStatement()
	case token.COLON:
		return p.parseLabelStatement()
	case token.REM, token.DOUBLE_COLON:
		return p.parseComment()
	case token.AT:
		return p.parseAtCommand()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseIfStatement() *ast.IfStmt {
	stmt := &ast.IfStmt{If: p.curToken.Pos}

	p.nextToken() // 跳过if

	// 解析条件
	stmt.Cond = p.parseExpression()

	// 检查括号
	if p.curToken.Type == token.LPAREN {
		stmt.Lparen = p.curToken.Pos
		p.nextToken() // 跳过左括号

		// 解析if体
		stmt.Body = p.parseBlockStatement()

		if p.curToken.Type == token.RPAREN {
			stmt.Rparen = p.curToken.Pos
			p.nextToken() // 跳过右括号
		}
	} else {
		// 没有括号，解析单个语句
		stmt.Body = p.parseStatement()
	}

	// 检查else
	if p.curToken.Type == token.ELSE {
		p.nextToken() // 跳过else
		stmt.Else = p.parseStatement()
	}

	return stmt
}

func (p *Parser) parseForStatement() *ast.ForStmt {
	stmt := &ast.ForStmt{For: p.curToken.Pos}

	p.nextToken() // 跳过for

	if p.curToken.Type == token.IDENT {
		stmt.X = p.parseExpression()
	}

	if p.curToken.Type == token.IN {
		stmt.In = p.curToken.Pos
		p.nextToken() // 跳过in
		stmt.List = p.parseExpression()
	}

	if p.curToken.Type == token.DO {
		stmt.Do = p.curToken.Pos
		p.nextToken() // 跳过do
	}

	if p.curToken.Type == token.LPAREN {
		p.nextToken() // 跳过左括号
		stmt.Body = p.parseBlockStatement()

		if p.curToken.Type == token.RPAREN {
			stmt.Rparen = p.curToken.Pos
			p.nextToken() // 跳过右括号
		}
	}

	return stmt
}

func (p *Parser) parseCallStatement() *ast.CallStmt {
	stmt := &ast.CallStmt{Call: p.curToken.Pos}

	p.nextToken() // 跳过call

	if p.curToken.Type == token.COLON {
		stmt.Colon = p.curToken.Pos
		stmt.IsFile = false
		p.nextToken() // 跳过冒号
	} else {
		stmt.IsFile = true
	}

	if p.curToken.Type == token.IDENT {
		stmt.Name = p.curToken.Literal
		p.nextToken()
	}

	// 解析参数
	for p.curToken.Type != token.SEMI && p.curToken.Type != token.EOF {
		stmt.Recv = append(stmt.Recv, p.parseExpression())
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseGotoStatement() *ast.GotoStmt {
	stmt := &ast.GotoStmt{GoTo: p.curToken.Pos}

	p.nextToken() // 跳过goto

	if p.curToken.Type == token.COLON {
		stmt.Colon = p.curToken.Pos
		p.nextToken() // 跳过冒号
	}

	if p.curToken.Type == token.IDENT {
		stmt.Label = p.curToken.Literal
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseLabelStatement() *ast.LabelStmt {
	stmt := &ast.LabelStmt{Colon: p.curToken.Pos}

	p.nextToken() // 跳过冒号

	if p.curToken.Type == token.IDENT {
		stmt.Name = p.curToken.Literal
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseComment() *ast.Comment {
	comment := &ast.Comment{
		TokPos: p.curToken.Pos,
		Tok:    p.curToken.Type,
		Text:   p.curToken.Literal,
	}
	return comment
}

func (p *Parser) parseAtCommand() *ast.ExprStmt {
	p.nextToken() // 跳过@
	expr := p.parseExpression()
	return &ast.ExprStmt{X: expr}
}

func (p *Parser) parseExpressionStatement() *ast.ExprStmt {
	stmt := &ast.ExprStmt{X: p.parseExpression()}
	return stmt
}

func (p *Parser) parseExpression() ast.Expr {
	switch p.curToken.Type {
	case token.IDENT, token.LITERAL:
		lit := &ast.Lit{
			ValPos: p.curToken.Pos,
			Val:    p.curToken.Literal,
		}
		return lit
	case token.LPAREN:
		p.nextToken() // 跳过左括号
		expr := p.parseExpression()
		if p.curToken.Type == token.RPAREN {
			p.nextToken() // 跳过右括号
		}
		return expr
	default:
		// 对于其他情况，返回字面量
		lit := &ast.Lit{
			ValPos: p.curToken.Pos,
			Val:    p.curToken.Literal,
		}
		return lit
	}
}

func (p *Parser) parseBlockStatement() *ast.BlockStmt {
	block := &ast.BlockStmt{}

	// 如果没有左大括号，就解析单个语句
	if p.curToken.Type != token.LBRACE {
		stmt := p.parseStatement()
		if stmt != nil {
			block.List = append(block.List, stmt)
		}
		return block
	}

	p.nextToken() // 跳过左大括号

	for p.curToken.Type != token.RBRACE && p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			block.List = append(block.List, stmt)
		}
		p.nextToken()
	}

	return block
}

func (p *Parser) expectPeek(t token.Token) bool {
	if p.peekToken.Type == t {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

func (p *Parser) peekError(t token.Token) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) addError(msg string) {
	p.errors = append(p.errors, msg)
}

func (p *Parser) Errors() []string {
	return p.errors
}
