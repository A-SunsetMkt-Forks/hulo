// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package parser

import (
	"fmt"
	"os"
	"strings"
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

func TestParseFunctionType(t *testing.T) {
	node, err := ParseSourceScript(`fn add(a: num, b: num) -> num { return a + b }`, OptionTracerASTTree(os.Stdout))
	assert.NoError(t, err)
	ast.Inspect(node, os.Stdout)
}

func TestParseChainedMethodCall(t *testing.T) {
	node, err := ParseSourceScript(`$p.to_str(10, "Hello World").to_num(true)`, OptionTracerASTTree(os.Stdout))
	assert.NoError(t, err)
	ast.Inspect(node, os.Stdout)
}

func TestDeferCallExpr(t *testing.T) {
	node, err := ParseSourceScript(`defer echo(1)`, OptionTracerASTTree(os.Stdout))
	assert.NoError(t, err)
	ast.Inspect(node, os.Stdout)
}

func TestDeferBlockStmt(t *testing.T) {
	node, err := ParseSourceScript(`defer { echo(1) }`, OptionTracerASTTree(os.Stdout))
	assert.NoError(t, err)
	ast.Inspect(node, os.Stdout)
}

func TestGenerics(t *testing.T) {
	node, err := ParseSourceScript(`fn add<T, U>(a: T, b: U) -> T { return $a + $b }`, OptionTracerASTTree(os.Stdout))
	assert.NoError(t, err)
	ast.Inspect(node, os.Stdout)
}

func TestGenericConstraints(t *testing.T) {
	tests := []struct {
		name string
		code string
	}{
		{
			name: "Single constraint",
			code: `fn process<T extends Addable>(item: T) -> T { return $item + $item }`,
		},
		{
			name: "Multiple constraints (intersection)",
			code: `fn process<T extends Addable & Comparable>(item: T) -> T { if $item > 0 { return $item + $item } return $item }`,
		},
		{
			name: "Union constraints",
			code: `fn process<T extends str | num>(item: T) -> T { return $item }`,
		},
		{
			name: "Conditional type",
			code: `fn process<T extends (T extends str ? HasLength : HasValue)>(item: T) -> T { return $item }`,
		},
		{
			name: "Mapped type",
			code: `fn process<T extends { [K in keyof T]: T[K] extends str ? T[K] : never }>(item: T) -> T { return $item }`,
		},
		{
			name: "Keyof type",
			code: `fn process<T extends keyof User>(item: T) -> T { return $item }`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node, err := ParseSourceScript(tt.code)
			assert.NoError(t, err)
			ast.Inspect(node, os.Stdout)
		})
	}
}

func TestInferType(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "infer type in conditional",
			input:    "T extends U ? infer P : never",
			expected: "*ast.ConditionalType {\n  CheckType: *ast.TypeReference {\n    Name: *ast.Ident (Name: \"T\")\n  }\n  ExtendsType: *ast.TypeReference {\n    Name: *ast.Ident (Name: \"U\")\n  }\n  TrueType: *ast.InferType {\n    Type: *ast.TypeReference {\n      Name: *ast.Ident (Name: \"P\")\n    }\n  }\n  FalseType: *ast.TypeReference {\n    Name: *ast.Ident (Name: \"never\")\n  }\n}",
		},
		{
			name:     "infer type with never",
			input:    "T extends U ? infer never : never",
			expected: "*ast.ConditionalType {\n  CheckType: *ast.TypeReference {\n    Name: *ast.Ident (Name: \"T\")\n  }\n  ExtendsType: *ast.TypeReference {\n    Name: *ast.Ident (Name: \"U\")\n  }\n  TrueType: *ast.InferType {\n    Type: *ast.TypeReference {\n      Name: *ast.Ident (Name: \"never\")\n    }\n  }\n  FalseType: *ast.TypeReference {\n    Name: *ast.Ident (Name: \"never\")\n  }\n}",
		},
		{
			name:     "infer type with null",
			input:    "T extends U ? infer null : never",
			expected: "*ast.ConditionalType {\n  CheckType: *ast.TypeReference {\n    Name: *ast.Ident (Name: \"T\")\n  }\n  ExtendsType: *ast.TypeReference {\n    Name: *ast.Ident (Name: \"U\")\n  }\n  TrueType: *ast.InferType {\n    Type: *ast.TypeReference {\n      Name: *ast.Ident (Name: \"null\")\n    }\n  }\n  FalseType: *ast.TypeReference {\n    Name: *ast.Ident (Name: \"never\")\n  }\n}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a type declaration to test the conditional type
			input := "type Test = " + tt.input

			analyzer, err := NewAnalyzer(antlr.NewInputStream(input))
			if err != nil {
				t.Fatalf("Failed to create analyzer: %v", err)
			}

			file := analyzer.file
			if file == nil {
				t.Fatal("File is nil")
			}

			// Get the first statement
			if len(file.AllStatement()) == 0 {
				t.Fatal("No statements found")
			}

			stmt := file.Statement(0)
			if stmt == nil {
				t.Fatal("First statement is nil")
			}

			// Parse the statement
			result := analyzer.VisitStatement(stmt.(*generated.StatementContext))
			if result == nil {
				t.Fatal("VisitStatement returned nil")
			}

			// Get the type value from the type declaration
			typeDecl, ok := result.(*ast.TypeDecl)
			if !ok {
				t.Fatalf("Expected TypeDecl, got %T", result)
			}

			// Convert to string for comparison
			var buf strings.Builder
			ast.Inspect(typeDecl.Value, &buf)

			actual := strings.TrimSpace(buf.String())
			expected := strings.TrimSpace(tt.expected)

			if actual != expected {
				t.Errorf("Expected:\n%s\n\nGot:\n%s", expected, actual)
			}
		})
	}
}

func TestUtilityTypes(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "If conditional type",
			input:    "C extends true ? T : F",
			expected: "*ast.ConditionalType {\n  CheckType: *ast.TypeReference {\n    Name: *ast.Ident (Name: \"C\")\n  }\n  ExtendsType: *ast.TrueLiteral\n  TrueType: *ast.TypeReference {\n    Name: *ast.Ident (Name: \"T\")\n  }\n  FalseType: *ast.TypeReference {\n    Name: *ast.Ident (Name: \"F\")\n  }\n}",
		},
		{
			name:     "NonNullable type",
			input:    "T extends null ? never : T",
			expected: "*ast.ConditionalType {\n  CheckType: *ast.TypeReference {\n    Name: *ast.Ident (Name: \"T\")\n  }\n  ExtendsType: *ast.TypeReference {\n    Name: *ast.Ident (Name: \"null\")\n  }\n  TrueType: *ast.TypeReference {\n    Name: *ast.Ident (Name: \"never\")\n  }\n  FalseType: *ast.TypeReference {\n    Name: *ast.Ident (Name: \"T\")\n  }\n}",
		},
		{
			name:     "Parameters type",
			input:    "T extends (...args: any[]) -> any ? infer P : never",
			expected: "*ast.ConditionalType {\n  CheckType: *ast.TypeReference {\n    Name: *ast.Ident (Name: \"T\")\n  }\n  ExtendsType: *ast.FunctionType {\n    Recv: [\n      0: *ast.Parameter {\n        Name: *ast.Ident (Name: \"args\")\n        Type: *ast.ArrayType {\n          Name: *ast.TypeReference {\n            Name: *ast.Ident (Name: \"any\")\n          }\n        }\n      }\n    ]\n    RetVal: *ast.TypeReference {\n      Name: *ast.Ident (Name: \"any\")\n    }\n  }\n  TrueType: *ast.InferType {\n    Type: *ast.TypeReference {\n      Name: *ast.Ident (Name: \"P\")\n    }\n  }\n  FalseType: *ast.TypeReference {\n    Name: *ast.Ident (Name: \"never\")\n  }\n}",
		},
		{
			name:     "ReturnType type",
			input:    "T extends (...args: any[]) -> infer R ? R : never",
			expected: "*ast.ConditionalType {\n  CheckType: *ast.TypeReference {\n    Name: *ast.Ident (Name: \"T\")\n  }\n  ExtendsType: *ast.FunctionType {\n    Recv: [\n      0: *ast.Parameter {\n        Name: *ast.Ident (Name: \"args\")\n        Type: *ast.ArrayType {\n          Name: *ast.TypeReference {\n            Name: *ast.Ident (Name: \"any\")\n          }\n        }\n      }\n    ]\n    RetVal: *ast.InferType {\n      Type: *ast.TypeReference {\n        Name: *ast.Ident (Name: \"R\")\n      }\n    }\n  }\n  TrueType: *ast.TypeReference {\n    Name: *ast.Ident (Name: \"R\")\n  }\n  FalseType: *ast.TypeReference {\n    Name: *ast.Ident (Name: \"never\")\n  }\n}",
		},
		{
			name:     "Optional type",
			input:    "T | null",
			expected: "*ast.UnionType {\n  Types: [\n    0: *ast.TypeReference {\n      Name: *ast.Ident (Name: \"T\")\n    }\n    1: *ast.TypeReference {\n      Name: *ast.Ident (Name: \"null\")\n    }\n  ]\n}",
		},
		{
			name:     "Readonly type",
			input:    "{ readonly [P in keyof T]: T[P] }",
			expected: "*ast.ObjectType {\n  Properties: []*ast.Property {\n    *ast.Property {\n      Key: *ast.MappedType {\n        KeyType: *ast.TypeReference {\n          Name: *ast.Ident (Name: \"P\")\n        }\n        Constraint: *ast.KeyofType {\n          Type: *ast.TypeReference {\n            Name: *ast.Ident (Name: \"T\")\n          }\n        }\n        ValueType: *ast.IndexAccess {\n          Object: *ast.TypeReference {\n            Name: *ast.Ident (Name: \"T\")\n          }\n          Index: *ast.TypeReference {\n            Name: *ast.Ident (Name: \"P\")\n          }\n        }\n        Readonly: true\n      }\n    }\n  }\n}",
		},
		{
			name:     "Partial type",
			input:    "{ [P in keyof T]?: T[P] }",
			expected: "*ast.ObjectType {\n  Properties: []*ast.Property {\n    *ast.Property {\n      Key: *ast.MappedType {\n        KeyType: *ast.TypeReference {\n          Name: *ast.Ident (Name: \"P\")\n        }\n        Constraint: *ast.KeyofType {\n          Type: *ast.TypeReference {\n            Name: *ast.Ident (Name: \"T\")\n          }\n        }\n        ValueType: *ast.IndexAccess {\n          Object: *ast.TypeReference {\n            Name: *ast.Ident (Name: \"T\")\n          }\n          Index: *ast.TypeReference {\n            Name: *ast.Ident (Name: \"P\")\n          }\n        }\n        Optional: true\n      }\n    }\n  }\n}",
		},
		{
			name:     "Required type",
			input:    "{ [P in keyof T]-?: T[P] }",
			expected: "*ast.ObjectType {\n  Properties: []*ast.Property {\n    *ast.Property {\n      Key: *ast.MappedType {\n        KeyType: *ast.TypeReference {\n          Name: *ast.Ident (Name: \"P\")\n        }\n        Constraint: *ast.KeyofType {\n          Type: *ast.TypeReference {\n            Name: *ast.Ident (Name: \"T\")\n          }\n        }\n        ValueType: *ast.IndexAccess {\n          Object: *ast.TypeReference {\n            Name: *ast.Ident (Name: \"T\")\n          }\n          Index: *ast.TypeReference {\n            Name: *ast.Ident (Name: \"P\")\n          }\n        }\n        Required: true\n      }\n    }\n  }\n}",
		},
		{
			name:     "Pick type",
			input:    "{ [P in K]: T[P] }",
			expected: "*ast.ObjectType {\n  Properties: []*ast.Property {\n    *ast.Property {\n      Key: *ast.MappedType {\n        KeyType: *ast.TypeReference {\n          Name: *ast.Ident (Name: \"P\")\n        }\n        Constraint: *ast.TypeReference {\n          Name: *ast.Ident (Name: \"K\")\n        }\n        ValueType: *ast.IndexAccess {\n          Object: *ast.TypeReference {\n            Name: *ast.Ident (Name: \"T\")\n          }\n          Index: *ast.TypeReference {\n            Name: *ast.Ident (Name: \"P\")\n          }\n        }\n      }\n    }\n  }\n}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a type declaration to test the utility type
			input := "type Test = " + tt.input

			analyzer, err := NewAnalyzer(antlr.NewInputStream(input))
			if err != nil {
				t.Fatalf("Failed to create analyzer: %v", err)
			}

			file := analyzer.file
			if file == nil {
				t.Fatal("File is nil")
			}

			// Get the first statement
			if len(file.AllStatement()) == 0 {
				t.Fatal("No statements found")
			}

			stmt := file.Statement(0)
			if stmt == nil {
				t.Fatal("First statement is nil")
			}

			// Parse the statement
			result := analyzer.VisitStatement(stmt.(*generated.StatementContext))
			if result == nil {
				t.Fatal("VisitStatement returned nil")
			}

			// Get the type value from the type declaration
			typeDecl, ok := result.(*ast.TypeDecl)
			if !ok {
				t.Fatalf("Expected TypeDecl, got %T", result)
			}

			// Convert to string for comparison
			var buf strings.Builder
			ast.Inspect(typeDecl.Value, &buf)

			actual := strings.TrimSpace(buf.String())
			expected := strings.TrimSpace(tt.expected)

			if actual != expected {
				t.Errorf("Expected:\n%s\n\nGot:\n%s", expected, actual)
			}
		})
	}
}

func TestDecorator(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "class decorator",
			input: `@Component
class User {
  name: str
}`,
			expected: "*ast.File {\n  Stmts: [\n    0: *ast.ClassDecl {\n      Decs: [\n        0: *ast.Decorator {\n          Name: *ast.Ident (Name: \"Component\")\n        }\n      ]\n      Name: *ast.Ident (Name: \"User\")\n      Fields: [\n        0: *ast.Field {\n          Name: *ast.Ident (Name: \"name\")\n          Type: *ast.TypeReference {\n            Name: *ast.Ident (Name: \"str\")\n          }\n        }\n      ]\n    }\n  ]\n}",
		},
		{
			name: "class decorator with arguments",
			input: `@Component("user", version: "1.0")
class User {
  name: str
}`,
			expected: "*ast.File {\n  Stmts: [\n    0: *ast.ClassDecl {\n      Decs: [\n        0: *ast.Decorator {\n          Name: *ast.Ident (Name: \"Component\")\n          Recv: [\n            0: *ast.StringLiteral (Value: \"user\")\n            1: *ast.KeyValueExpr {\n              Key: *ast.Ident (Name: \"version\")\n              Value: *ast.StringLiteral (Value: \"1.0\")\n            }\n          ]\n        }\n      ]\n      Name: *ast.Ident (Name: \"User\")\n      Fields: [\n        0: *ast.Field {\n          Name: *ast.Ident (Name: \"name\")\n          Type: *ast.TypeReference {\n            Name: *ast.Ident (Name: \"str\")\n          }\n        }\n      ]\n    }\n  ]\n}",
		},
		{
			name: "function decorator",
			input: `@Log
fn greet(name: str) -> str {
  return "Hello " + name
}`,
			expected: "*ast.File {\n  Stmts: [\n    0: *ast.FuncDecl {\n      Decs: [\n        0: *ast.Decorator {\n          Name: *ast.Ident (Name: \"Log\")\n        }\n      ]\n      Name: *ast.Ident (Name: \"greet\")\n      Recv: [\n        0: *ast.Parameter {\n          Name: *ast.Ident (Name: \"name\")\n          Type: *ast.TypeReference {\n            Name: *ast.Ident (Name: \"str\")\n          }\n        }\n      ]\n      Type: *ast.TypeReference {\n        Name: *ast.Ident (Name: \"str\")\n      }\n    }\n  ]\n}",
		},
		{
			name: "function decorator with arguments",
			input: `@Log(level: "info")
fn greet(name: str) -> str {
  return "Hello " + name
}`,
			expected: "*ast.File {\n  Stmts: [\n    0: *ast.FuncDecl {\n      Decs: [\n        0: *ast.Decorator {\n          Name: *ast.Ident (Name: \"Log\")\n          Recv: [\n            0: *ast.KeyValueExpr {\n              Key: *ast.Ident (Name: \"level\")\n              Value: *ast.StringLiteral (Value: \"info\")\n            }\n          ]\n        }\n      ]\n      Name: *ast.Ident (Name: \"greet\")\n      Recv: [\n        0: *ast.Parameter {\n          Name: *ast.Ident (Name: \"name\")\n          Type: *ast.TypeReference {\n            Name: *ast.Ident (Name: \"str\")\n          }\n        }\n      ]\n      Type: *ast.TypeReference {\n        Name: *ast.Ident (Name: \"str\")\n      }\n    }\n  ]\n}",
		},
		{
			name: "method decorator",
			input: `class User {
  @Validate
  fn setAge(age: num) {
    this.age = age
  }
}`,
			expected: "*ast.File {\n  Stmts: [\n    0: *ast.ClassDecl {\n      Name: *ast.Ident (Name: \"User\")\n      Methods: [\n        0: *ast.FuncDecl {\n          Decs: [\n            0: *ast.Decorator {\n              Name: *ast.Ident (Name: \"Validate\")\n            }\n          ]\n          Name: *ast.Ident (Name: \"setAge\")\n          Recv: [\n            0: *ast.Parameter {\n              Name: *ast.Ident (Name: \"age\")\n              Type: *ast.TypeReference {\n                Name: *ast.Ident (Name: \"num\")\n              }\n            }\n          ]\n        }\n      ]\n    }\n  ]\n}",
		},
		{
			name: "multiple method decorators",
			input: `class User {
  @Component()
  @ABC()
  fn to_str() {
  }
}`,
			expected: "*ast.File {\n  Stmts: [\n    0: *ast.ClassDecl {\n      Name: *ast.Ident (Name: \"User\")\n      Methods: [\n        0: *ast.FuncDecl {\n          Decs: [\n            0: *ast.Decorator {\n              Name: *ast.Ident (Name: \"Component\")\n            }\n            1: *ast.Decorator {\n              Name: *ast.Ident (Name: \"ABC\")\n            }\n          ]\n          Name: *ast.Ident (Name: \"to_str\")\n        }\n      ]\n    }\n  ]\n}",
		},
		{
			name: "field decorator",
			input: `class User {
  @Readonly
  data: str
}`,
			expected: "*ast.File {\n  Stmts: [\n    0: *ast.ClassDecl {\n      Name: *ast.Ident (Name: \"User\")\n      Fields: [\n        0: *ast.Field {\n          Decs: [\n            0: *ast.Decorator {\n              Name: *ast.Ident (Name: \"Readonly\")\n            }\n          ]\n          Name: *ast.Ident (Name: \"data\")\n          Type: *ast.TypeReference {\n            Name: *ast.Ident (Name: \"str\")\n          }\n        }\n      ]\n    }\n  ]\n}",
		},
		{
			name: "multiple class decorators",
			input: `@Component
@Injectable
class Service {
  data: str
}`,
			expected: "*ast.File {\n  Stmts: [\n    0: *ast.ClassDecl {\n      Decs: [\n        0: *ast.Decorator {\n          Name: *ast.Ident (Name: \"Component\")\n        }\n        1: *ast.Decorator {\n          Name: *ast.Ident (Name: \"Injectable\")\n        }\n      ]\n      Name: *ast.Ident (Name: \"Service\")\n      Fields: [\n        0: *ast.Field {\n          Name: *ast.Ident (Name: \"data\")\n          Type: *ast.TypeReference {\n            Name: *ast.Ident (Name: \"str\")\n          }\n        }\n      ]\n    }\n  ]\n}",
		},
		{
			name: "enum decorator",
			input: `@Serializable
enum Status {
  Pending,
  Approved,
  Rejected
}`,
			expected: "*ast.File {\n  Stmts: [\n    0: *ast.EnumDecl {\n      Decs: [\n        0: *ast.Decorator {\n          Name: *ast.Ident (Name: \"Serializable\")\n        }\n      ]\n      Name: *ast.Ident (Name: \"Status\")\n      Body: [\n        0: *ast.EnumValue {\n          Name: *ast.Ident (Name: \"Pending\")\n        }\n        1: *ast.EnumValue {\n          Name: *ast.Ident (Name: \"Approved\")\n        }\n        2: *ast.EnumValue {\n          Name: *ast.Ident (Name: \"Rejected\")\n        }\n      ]\n    }\n  ]\n}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			analyzer, err := NewAnalyzer(antlr.NewInputStream(tt.input))
			if err != nil {
				t.Fatalf("Failed to create analyzer: %v", err)
			}

			file := analyzer.file
			if file == nil {
				t.Fatal("File is nil")
			}

			// Parse the file
			result := analyzer.VisitFile(file.(*generated.FileContext))
			if result == nil {
				t.Fatal("VisitFile returned nil")
			}

			// Convert to string for comparison
			var buf strings.Builder
			ast.Inspect(result.(ast.Node), &buf)

			actual := strings.TrimSpace(buf.String())
			expected := strings.TrimSpace(tt.expected)

			if actual != expected {
				t.Errorf("Expected:\n%s\n\nGot:\n%s", expected, actual)
			}
		})
	}
}
