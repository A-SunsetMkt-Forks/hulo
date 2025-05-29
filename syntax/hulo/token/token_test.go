package token_test

import (
	"testing"

	"github.com/hulo-lang/hulo/syntax/hulo/token"
	"github.com/stretchr/testify/assert"
)

func TestTokenString(t *testing.T) {
	tests := []struct {
		tok      token.Token
		expected string
	}{
		{token.ILLEGAL, "ILLEGAL"},
		{token.EOF, "EOF"},
		{token.IDENT, "IDENT"},
		{token.PUB, "pub"},
		{token.STATIC, "static"},
	}

	for _, test := range tests {
		assert.Equal(t, test.expected, test.tok.String())
	}
}
