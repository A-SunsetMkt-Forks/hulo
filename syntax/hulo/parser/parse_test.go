// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package parser

import (
	"fmt"
	"os"
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
	node, err := ParseSourceScript("Write-Host 'Hello World!'")
	assert.NoError(t, err)
	ast.Print(node)
}

func TestParseSourceLoop(t *testing.T) {
	node, err := ParseSourceScript("loop $a < 2 { MsgBox $a; $a++; }")
	assert.NoError(t, err)
	ast.Print(node)
}

func TestParseIfStmt(t *testing.T) {
	node, err := ParseSourceFile("./testdata/if.hl", OptionTracerASTTree(os.Stdout))
	assert.NoError(t, err)
	ast.Print(node)
}

func TestParseFuncDecl(t *testing.T) {
	node, err := ParseSourceFile("./testdata/func.hl", OptionTracerASTTree(os.Stdout))
	assert.NoError(t, err)
	ast.Print(node)
}

func TestParseClassDecl(t *testing.T) {
	node, err := ParseSourceFile("./testdata/class.hl", OptionTracerASTTree(os.Stdout))
	assert.NoError(t, err)
	ast.Print(node)
}

func TestParseComptime(t *testing.T) {
	node, err := ParseSourceFile("./testdata/comptime.hl", OptionTracerASTTree(os.Stdout))
	assert.NoError(t, err)
	ast.Print(node)
}

func TestParseImport(t *testing.T) {
	node, err := ParseSourceFile("./testdata/import.hl", OptionTracerASTTree(os.Stdout))
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

			// The first memberAccess is the command name
			if cmdExpr.MemberAccess() != nil {
				assert.Equal(t, tt.expectCmd, cmdExpr.MemberAccess().GetText())
			}
		})
	}
}

func TestParseDeclare(t *testing.T) {
	node, err := ParseSourceFile("./testdata/declare.hl", OptionTracerASTTree(os.Stdout))
	assert.NoError(t, err)
	ast.Print(node)
}

func TestParseTypeDecl(t *testing.T) {
	node, err := ParseSourceFile("./testdata/type.hl", OptionTracerASTTree(os.Stdout))
	assert.NoError(t, err)
	ast.Print(node)
}

func TestParseEnumDecl(t *testing.T) {
	node, err := ParseSourceFile("./testdata/enum.hl", OptionTracerASTTree(os.Stdout))
	assert.NoError(t, err)
	ast.Print(node)
}

func TestParseAccess(t *testing.T) {
	node, err := ParseSourceFile("./testdata/access.hl", OptionTracerASTTree(os.Stdout))
	assert.NoError(t, err)
	ast.Print(node)
}

func TestParseMatch(t *testing.T) {
	node, err := ParseSourceFile("./testdata/match.hl", OptionTracerASTTree(os.Stdout))
	assert.NoError(t, err)
	ast.Print(node)
}

func TestParseSelector(t *testing.T) {
	node, err := ParseSourceFile("./testdata/selector.hl", OptionTracerASTTree(os.Stdout))
	assert.NoError(t, err)
	ast.Print(node)
}

func TestParseUnsafe(t *testing.T) {
	node, err := ParseSourceFile("./testdata/unsafe.hl", OptionTracerASTTree(os.Stdout))
	assert.NoError(t, err)
	ast.Print(node)
}

func TestParseData(t *testing.T) {
	node, err := ParseSourceFile("./testdata/data.hl", OptionTracerASTTree(os.Stdout))
	assert.NoError(t, err)
	ast.Print(node)
}

func TestParseDate(t *testing.T) {
	node, err := ParseSourceFile("./testdata/date.hl", OptionTracerASTTree(os.Stdout))
	assert.NoError(t, err)
	ast.Print(node)
}

func TestParseEllipsis(t *testing.T) {
	node, err := ParseSourceScript("type Function = (...args: any[]) -> any", OptionTracerASTTree(os.Stdout))
	assert.NoError(t, err)
	ast.Print(node)
}

func TestParseFunction(t *testing.T) {
	node, err := ParseSourceScript(`// Creates and returns a reference to an Automation object.
declare fn CreateObject(
    // The name of the application providing the object.
    servername: String,
    // The type or class of the object to create.
    typename: String,
    // The name of the network server where the object is to be created. This feature is available in version 5.1 or later.
    location?: String
) -> Object`, OptionTracerASTTree(os.Stdout))
	assert.NoError(t, err)
	ast.Print(node)
}

// TODO: no support
func TestParseNumbericReturn(t *testing.T) {
	node, err := ParseSourceScript(`declare fn Sgn(x: num) -> -1 | 0 | 1`, OptionTracerASTTree(os.Stdout))
	assert.NoError(t, err)
	ast.Print(node)
}

func TestParserCallExpr(t *testing.T) {
	node, err := ParseSourceScript(`let result = utils.calculate(5, 3)`, OptionTracerASTTree(os.Stdout))
	assert.NoError(t, err)
	ast.Inspect(node, os.Stdout)
}

func TestParserChainedCallExpr(t *testing.T) {
	node, err := ParseSourceScript(`let result = utils.calculate(5, 3).to_str().len()`, OptionTracerASTTree(os.Stdout))
	assert.NoError(t, err)
	ast.Inspect(node, os.Stdout)
}

func TestParserChainedCallExprWithMod(t *testing.T) {
	node, err := ParseSourceScript(`let result = utils.a::b::c::User.create("John", 25)..name="ansurfen"..age=30;`, OptionTracerASTTree(os.Stdout), OptionDisableTracer())
	assert.NoError(t, err)
	ast.Inspect(node, os.Stdout)
}

func TestParserMethodCall(t *testing.T) {
	node, err := ParseSourceScript(`let u = User("John", 20); echo $u.to_str(); $u.greet("Jane");`, OptionTracerASTTree(os.Stdout), OptionDisableTracer())
	assert.NoError(t, err)
	ast.Inspect(node, os.Stdout)
}

func TestParseCommand(t *testing.T) {
	node, err := ParseSourceScript(`echo -e "Hello, World!"`, OptionTracerASTTree(os.Stdout))
	assert.NoError(t, err)
	ast.Inspect(node, os.Stdout)
}

func TestParseNestedCommand(t *testing.T) {
	node, err := ParseSourceScript(`my-cmd "a" sub-command "b"`, OptionTracerASTTree(os.Stdout))
	assert.NoError(t, err)
	ast.Inspect(node, os.Stdout)
}

func TestParseCommandWithOption(t *testing.T) {
	node, err := ParseSourceScript(`my-cmd -e --ignore-any-error -it`, OptionTracerASTTree(os.Stdout))
	assert.NoError(t, err)
	ast.Inspect(node, os.Stdout)
}

func TestParseCommandWithComplexOptions(t *testing.T) {
	node, err := ParseSourceScript(`my-cmd -e { name: "abc", age: 10 } -x [1, 2, 3]`, OptionTracerASTTree(os.Stdout))
	assert.NoError(t, err)
	ast.Inspect(node, os.Stdout)
}

func TestParseComplexCommand(t *testing.T) {
	node, err := ParseSourceScript(`docker run -it ubuntu bash -c "echo hello"`, OptionTracerASTTree(os.Stdout))
	assert.NoError(t, err)
	ast.Inspect(node, os.Stdout)
}

func TestParseCommandWithBuiltinArgs(t *testing.T) {
	node, err := ParseSourceScript(`container.std::my-command run -it ubuntu --- -f "abc"`, OptionTracerASTTree(os.Stdout))
	assert.NoError(t, err)
	ast.Inspect(node, os.Stdout)
}

func TestParsePipeline(t *testing.T) {
	node, err := ParseSourceScript(`container.std::my-command run -it ubuntu | grep "hello"`, OptionTracerASTTree(os.Stdout))
	assert.NoError(t, err)
	ast.Inspect(node, os.Stdout)
}

func TestParseDockerRun(t *testing.T) {
	node, err := ParseSourceScript(`container.std::docker run -it ubuntu`, OptionTracerASTTree(os.Stdout))
	assert.NoError(t, err)
	ast.Inspect(node, os.Stdout)
}

func TestParseUltraCommand(t *testing.T) {
	node, err := ParseSourceScript(`container.std::my-command run -it ubuntu --ignore-any-error --output "" -- -f ""
| grep "hello" && cd "~" && echo "world" &`, OptionTracerASTTree(os.Stdout))
	assert.NoError(t, err)
	ast.Inspect(node, os.Stdout)
}
