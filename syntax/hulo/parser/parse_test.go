package parser

import (
	"fmt"
	"testing"

	"github.com/antlr4-go/antlr/v4"
	"github.com/hulo-lang/hulo/syntax/hulo/ast"
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

func TestParseSourceFile(t *testing.T) {
	node, err := ParseSourceScript("Write-Host 'abc'", ParseOptions{})
	assert.NoError(t, err)
	ast.Print(node)
}

func TestIdentifier(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		expectFail bool
		expectCmd  string // optional: expected command name
	}{
		{
			name:      "ValidCommand_WriteHost",
			input:     "Write-Host",
			expectCmd: "Write-Host",
		},
		{
			name:      "ValidCommand_GetItem",
			input:     "Get-Item",
			expectCmd: "Get-Item",
		},
		{
			name:       "InvalidToken",
			input:      "123-Invalid",
			expectFail: true,
		},
		{
			name:       "EmptyInput",
			input:      "",
			expectFail: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stream := antlr.NewInputStream(tt.input)
			lexer := generated.NewhuloLexer(stream)
			tokens := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
			parser := generated.NewhuloParser(tokens)

			// optional: enable error listener capture
			parser.RemoveErrorListeners()
			errListener := &antlr.DiagnosticErrorListener{}
			parser.AddErrorListener(errListener)

			file := parser.File()

			if tt.expectFail {
				// 基本判断：没有 statement，表示解析失败
				// if len(file.AllStatement()) > 1 {
				// 	t.Errorf("expected failure, but got statements: %v", file.AllStatement())
				// }
				return
			}

			stmts := file.AllStatement()
			assert.Len(t, stmts, 1)

			stmt := stmts[0]
			exprStmt := stmt.ExpressionStatement()
			assert.NotNil(t, exprStmt)

			cmdExpr := exprStmt.Expression().CommandExpression()
			assert.NotNil(t, cmdExpr)

			member := cmdExpr.MemberAccess()
			assert.Equal(t, tt.expectCmd, member.GetText())
		})
	}
}
