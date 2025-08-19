// syntax/batch/lexer/lexer.go
// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package lexer

import (
	"github.com/hulo-lang/hulo/syntax/batch/token"
	"unicode"
)

type Lexer struct {
	input   string
	pos     int
	readPos int
	ch      rune
	line    int
	col     int
}

func NewLexer(input string) *Lexer {
	l := &Lexer{input: input, line: 1, col: 0}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPos >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = rune(l.input[l.readPos])
	}
	l.pos = l.readPos
	l.readPos++
	l.col++
}

func (l *Lexer) peekChar() rune {
	if l.readPos >= len(l.input) {
		return 0
	}
	return rune(l.input[l.readPos])
}

func (l *Lexer) NextToken() token.TokenInfo {
	l.skipWhitespace()

	pos := token.Pos(l.pos)

	switch l.ch {
	case '@':
		l.readChar()
		return token.TokenInfo{Type: token.AT, Pos: pos, Literal: "@"}
	case '(':
		l.readChar()
		return token.TokenInfo{Type: token.LPAREN, Pos: pos, Literal: "("}
	case ')':
		l.readChar()
		return token.TokenInfo{Type: token.RPAREN, Pos: pos, Literal: ")"}
	case '{':
		l.readChar()
		return token.TokenInfo{Type: token.LBRACE, Pos: pos, Literal: "{"}
	case '}':
		l.readChar()
		return token.TokenInfo{Type: token.RBRACE, Pos: pos, Literal: "}"}
	case ',':
		l.readChar()
		return token.TokenInfo{Type: token.COMMA, Pos: pos, Literal: ","}
	case ';':
		l.readChar()
		return token.TokenInfo{Type: token.SEMI, Pos: pos, Literal: ";"}
	case '=':
		if l.peekChar() == '=' {
			l.readChar()
			l.readChar()
			return token.TokenInfo{Type: token.EQU, Pos: pos, Literal: "=="}
		}
		l.readChar()
		return token.TokenInfo{Type: token.ASSIGN, Pos: pos, Literal: "="}
	case '+':
		if l.peekChar() == '=' {
			l.readChar()
			l.readChar()
			return token.TokenInfo{Type: token.ADD_ASSIGN, Pos: pos, Literal: "+="}
		}
		l.readChar()
		return token.TokenInfo{Type: token.ADD, Pos: pos, Literal: "+"}
	case '-':
		if l.peekChar() == '=' {
			l.readChar()
			l.readChar()
			return token.TokenInfo{Type: token.SUB_ASSIGN, Pos: pos, Literal: "-="}
		}
		l.readChar()
		return token.TokenInfo{Type: token.SUB, Pos: pos, Literal: "-"}
	case '*':
		if l.peekChar() == '=' {
			l.readChar()
			l.readChar()
			return token.TokenInfo{Type: token.MUL_ASSIGN, Pos: pos, Literal: "*="}
		}
		l.readChar()
		return token.TokenInfo{Type: token.MUL, Pos: pos, Literal: "*"}
	case '/':
		if l.peekChar() == '=' {
			l.readChar()
			l.readChar()
			return token.TokenInfo{Type: token.DIV_ASSIGN, Pos: pos, Literal: "/="}
		}
		l.readChar()
		return token.TokenInfo{Type: token.DIV, Pos: pos, Literal: "/"}
	case '%':
		if l.peekChar() == '=' {
			l.readChar()
			l.readChar()
			return token.TokenInfo{Type: token.MOD_ASSIGN, Pos: pos, Literal: "%="}
		}
		l.readChar()
		return token.TokenInfo{Type: token.MOD, Pos: pos, Literal: "%"}
	case '&':
		if l.peekChar() == '&' {
			l.readChar()
			l.readChar()
			return token.TokenInfo{Type: token.AND, Pos: pos, Literal: "&&"}
		}
		l.readChar()
		return token.TokenInfo{Type: token.BITAND, Pos: pos, Literal: "&"}
	case '|':
		if l.peekChar() == '|' {
			l.readChar()
			l.readChar()
			return token.TokenInfo{Type: token.OR, Pos: pos, Literal: "||"}
		}
		l.readChar()
		return token.TokenInfo{Type: token.BITOR, Pos: pos, Literal: "|"}
	case '^':
		l.readChar()
		return token.TokenInfo{Type: token.BITXOR, Pos: pos, Literal: "^"}
	case '!':
		l.readChar()
		return token.TokenInfo{Type: token.NOT, Pos: pos, Literal: "!"}
	case '~':
		l.readChar()
		return token.TokenInfo{Type: token.TILDE, Pos: pos, Literal: "~"}
	case '>':
		if l.peekChar() == '>' {
			l.readChar()
			l.readChar()
			return token.TokenInfo{Type: token.AppOut, Pos: pos, Literal: ">>"}
		}
		l.readChar()
		return token.TokenInfo{Type: token.RdrOut, Pos: pos, Literal: ">"}
	case '<':
		l.readChar()
		return token.TokenInfo{Type: token.RdrIn, Pos: pos, Literal: "<"}
	case ':':
		if l.peekChar() == ':' {
			l.readChar()
			l.readChar()
			return token.TokenInfo{Type: token.DOUBLE_COLON, Pos: pos, Literal: "::"}
		}
		l.readChar()
		return token.TokenInfo{Type: token.COLON, Pos: pos, Literal: ":"}
	case '"':
		return l.readString()
	case 0:
		return token.TokenInfo{Type: token.EOF, Pos: pos, Literal: ""}
	default:
		if isLetter(l.ch) {
			ident := l.readIdentifier()
			tokType := l.lookupIdent(ident)
			return token.TokenInfo{Type: tokType, Pos: pos, Literal: ident}
		} else if isDigit(l.ch) {
			number := l.readNumber()
			return token.TokenInfo{Type: token.LITERAL, Pos: pos, Literal: number}
		} else {
			illegal := string(l.ch)
			l.readChar()
			return token.TokenInfo{Type: token.ILLEGAL, Pos: pos, Literal: illegal}
		}
	}
}

func (l *Lexer) readString() token.TokenInfo {
	pos := token.Pos(l.pos)
	l.readChar() // 跳过开始的引号

	var str string
	for l.ch != '"' && l.ch != 0 {
		str += string(l.ch)
		l.readChar()
	}

	if l.ch == '"' {
		l.readChar() // 跳过结束的引号
	}

	return token.TokenInfo{Type: token.LITERAL, Pos: pos, Literal: "\"" + str + "\""}
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		if l.ch == '\n' {
			l.line++
			l.col = 0
		}
		l.readChar()
	}
}

func (l *Lexer) readIdentifier() string {
	start := l.pos
	for isLetter(l.ch) || isDigit(l.ch) || l.ch == '_' {
		l.readChar()
	}
	return l.input[start:l.pos]
}

func (l *Lexer) readNumber() string {
	start := l.pos
	for isDigit(l.ch) || l.ch == '.' {
		l.readChar()
	}
	return l.input[start:l.pos]
}

func (l *Lexer) lookupIdent(ident string) token.Token {
	switch ident {
	case "if":
		return token.IF
	case "else":
		return token.ELSE
	case "for":
		return token.FOR
	case "do":
		return token.DO
	case "in":
		return token.IN
	case "call":
		return token.CALL
	case "goto":
		return token.GOTO
	case "rem":
		return token.REM
	default:
		return token.IDENT
	}
}

func isLetter(ch rune) bool {
	return unicode.IsLetter(ch) || ch == '_'
}

func isDigit(ch rune) bool {
	return unicode.IsDigit(ch)
}
