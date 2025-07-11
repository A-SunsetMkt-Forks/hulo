// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package build_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	build "github.com/hulo-lang/hulo/internal/build/vbs"
	"github.com/hulo-lang/hulo/internal/config"
	"github.com/hulo-lang/hulo/internal/vfs/memvfs"
	"github.com/hulo-lang/hulo/syntax/hulo/parser"
	vast "github.com/hulo-lang/hulo/syntax/vbs/ast"
	"github.com/stretchr/testify/assert"
)

func TestCommandStmt(t *testing.T) {
	script := `echo "Hello, World!" 3.14 true`
	node, err := parser.ParseSourceScript(script)
	assert.NoError(t, err)
	// hast.Print(node)
	bnode, err := build.TranspileToVBScript(&config.VBScriptOptions{}, node)
	assert.NoError(t, err)
	vast.Print(bnode)
}

func TestComment(t *testing.T) {
	script := `// this is a comment

	/*
	 * this
	 *    is a multi
	 * comment
	 */

	/**
	 * this
	 *    is a multi
	 * comment
	 */`
	node, err := parser.ParseSourceScript(script)
	assert.NoError(t, err)
	// hast.Print(node)
	bnode, err := build.TranspileToVBScript(&config.VBScriptOptions{
		CommentSyntax: "'",
	}, node)
	assert.NoError(t, err)
	vast.Print(bnode)
}

func TestAssign(t *testing.T) {
	script := `let a = 10
	var b = 3.14
	const c = "Hello, World!"
	$d := true`
	node, err := parser.ParseSourceScript(script)
	assert.NoError(t, err)
	// hast.Print(node)
	bnode, err := build.TranspileToVBScript(&config.VBScriptOptions{}, node)
	assert.NoError(t, err)
	vast.Print(bnode)
}

func TestIf(t *testing.T) {
	script := `if $a > 10 {
		echo "a is greater than 10"
	} else {
		echo "a is less than or equal to 10"
	}`
	node, err := parser.ParseSourceScript(script)
	assert.NoError(t, err)
	// hast.Print(node)
	bnode, err := build.TranspileToVBScript(&config.VBScriptOptions{}, node)
	assert.NoError(t, err)
	vast.Print(bnode)
}

func TestNestedIf(t *testing.T) {
	script := `if $a > 10 {
		echo "a is greater than 10"
		if $b > 20 {
			echo "b is greater than 20"
		} else {
			echo "b is less than or equal to 20"
		}
	} else if $a < 0 {
		echo "a is less than 0"
	} else {
		echo "a is 0"
	}`
	node, err := parser.ParseSourceScript(script)
	assert.NoError(t, err)
	// hast.Print(node)
	bnode, err := build.TranspileToVBScript(&config.VBScriptOptions{}, node)
	assert.NoError(t, err)
	vast.Print(bnode)
}

func TestStandardWhile(t *testing.T) {
	script := `loop $a < 2 { MsgBox $a; $a++; }`
	node, err := parser.ParseSourceScript(script)
	assert.NoError(t, err)
	// hast.Print(node)
	bnode, err := build.TranspileToVBScript(&config.VBScriptOptions{}, node)
	assert.NoError(t, err)
	vast.Print(bnode)
}

func TestWhile(t *testing.T) {
	script := `
	loop {
		echo "Hello, World!"
	}

	loop $a == true {
		echo "Hello, World!"
	}

	do {
		echo "Hello, World!"
	} loop ($a > 10)`
	node, err := parser.ParseSourceScript(script, parser.OptionTracerASTTree(os.Stdout))
	assert.NoError(t, err)
	// hast.Print(node)
	bnode, err := build.TranspileToVBScript(&config.VBScriptOptions{}, node)
	assert.NoError(t, err)
	vast.Print(bnode)
}

// TODO j-- error
func TestFor(t *testing.T) {
	script := `
	declare fn MsgBox(message: str);
	loop $i := 0; $i < 10; $i++ {
		MsgBox "Hello, World!"
		loop $j := 10; $j > 0; $j-- {
			MsgBox "Hello, World!"
		}
	}`
	node, err := parser.ParseSourceScript(script, parser.OptionTracerASTTree(os.Stdout))
	assert.NoError(t, err)
	// hast.Print(node)
	bnode, err := build.TranspileToVBScript(&config.VBScriptOptions{}, node)
	assert.NoError(t, err)
	vast.Print(bnode)
}

func TestForIn(t *testing.T) {
	script := `
	declare fn MsgBox(message: str);
	let arr: list<num> = [1, 3.14, 5.0, 0.7]
	loop $item in $arr {
		MsgBox $item
	}

	loop $i in [0, 1, 2] {
		MsgBox $i
	}`
	node, err := parser.ParseSourceScript(script, parser.OptionTracerASTTree(os.Stdout))
	assert.NoError(t, err)
	// hast.Print(node)
	bnode, err := build.TranspileToVBScript(&config.VBScriptOptions{}, node)
	assert.NoError(t, err)
	vast.Print(bnode)
}

func TestForOf(t *testing.T) {
	script := `
declare fn MsgBox(message: str);
let config: map<str, str> = {"host": "localhost", "port": "8080"}
loop ($key, $value) of $config {
    MsgBox "$key = $value"
}
loop ($key, _) of $config {
    MsgBox $key
}

loop (_, $value) of $config {
    MsgBox $value
}

loop $key of $config {
    MsgBox $key
}
`
	node, err := parser.ParseSourceScript(script, parser.OptionTracerASTTree(os.Stdout))
	assert.NoError(t, err)
	// hast.Print(node)
	bnode, err := build.TranspileToVBScript(&config.VBScriptOptions{}, node)
	assert.NoError(t, err)
	vast.Print(bnode)
}

func TestFunc(t *testing.T) {
	script := `
	fn helloWorld() {
		echo "Hello, World!"
	}

	fn sayHello(name: str) {
		echo "Hello, $name!"
	}

	fn add(a: num, b: num) {
		return $a + $b
	}

	fn multiply(a: num, b: num) -> num {
		return $a * $b
	}
	`
	node, err := parser.ParseSourceScript(script, parser.OptionTracerASTTree(os.Stdout))
	assert.NoError(t, err)
	// hast.Print(node)
	bnode, err := build.TranspileToVBScript(&config.VBScriptOptions{}, node)
	assert.NoError(t, err)
	vast.Print(bnode)
}

func TestStringInterpolation(t *testing.T) {
	script := `echo "Hello, $name!"`
	node, err := parser.ParseSourceScript(script)
	assert.NoError(t, err)

	bnode, err := build.TranspileToVBScript(&config.VBScriptOptions{}, node)
	assert.NoError(t, err)
	vast.Print(bnode)
}

func TestStringEscape(t *testing.T) {
	script := `echo "Hello, \"World\"!"`
	node, err := parser.ParseSourceScript(script)
	assert.NoError(t, err)

	bnode, err := build.TranspileToVBScript(&config.VBScriptOptions{}, node)
	assert.NoError(t, err)
	vast.Print(bnode)
}

func TestClassDecl(t *testing.T) {
	script := `
		declare fn MsgBox(message: str)

		pub class Person {
			name: str
			age: num

			pub fn get_name() => $name
			pub fn get_age() => $age

			pub fn to_str() -> str {
				return "Person(name: $name, age: $age)"
			}

			pub fn greet() {
				MsgBox "Hello, my name is $name and I am $age years old."
			}

			pub fn greet_with_age(age: num) {
				MsgBox "Hello, my name is $name and I am $age years old."
			}
		}

		let p = Person()
		$p.name = "Tom"
		$p.age = 30
		$p.greet()

		let p2 = Person("Jerry", 20)
		$p2.greet()
		$p2.greet_with_age(10)
	`
	node, err := parser.ParseSourceScript(script, parser.OptionTracerASTTree(os.Stdout))
	assert.NoError(t, err)
	// hast.Print(node)
	bnode, err := build.TranspileToVBScript(&config.VBScriptOptions{}, node)
	assert.NoError(t, err)
	vast.Print(bnode)
}

func TestFieldModifiers(t *testing.T) {
	script := `
		class TestClass {
			name: str
			pub age: num
			pub className: str = "TestClass"
		}
	`
	node, err := parser.ParseSourceScript(script)
	assert.NoError(t, err)

	bnode, err := build.TranspileToVBScript(&config.VBScriptOptions{}, node)
	assert.NoError(t, err)

	fmt.Println("=== VBScript AST ===")
	vast.Print(bnode)
}

func TestMatch(t *testing.T) {
	script := `
let n = 10
match $n {
    10 => println("ten"),
    20 => println("twenty"),
    _ => println("unknown")
}

let status = "success"
match $status {
    "success" => println("Operation completed"),
    "error" => println("Operation failed"),
    "pending" => println("Operation in progress"),
    _ => println("Unknown status")
}`
	node, err := parser.ParseSourceScript(script, parser.OptionTracerASTTree(os.Stdout))
	assert.NoError(t, err)
	// hast.Print(node)
	bnode, err := build.TranspileToVBScript(&config.VBScriptOptions{}, node)
	assert.NoError(t, err)
	vast.Print(bnode)
}

func TestDeclare(t *testing.T) {
	script := `
		declare {
			const vbCr = "\r"
		}

		declare fn InputBox(prompt: str, title?: str, default?: str, xpos?: num, ypos?: num, helpfile?: str, context?: num) -> any;

		InputBox("Hello, World!", "InputBox", "default", 100, 100, "helpfile", 100)
		echo "Hello, World!"
		`
	node, err := parser.ParseSourceScript(script, parser.OptionTracerASTTree(os.Stdout))
	assert.NoError(t, err)
	// hast.Print(node)
	bnode, err := build.TranspileToVBScript(&config.VBScriptOptions{}, node)
	assert.NoError(t, err)
	vast.Print(bnode)
}

func TestImport(t *testing.T) {
	script := `
		import "unsafe/vbs"
		import "utils"

		MsgBox Add(5, 7)
	`
	node, err := parser.ParseSourceScript(script, parser.OptionTracerASTTree(os.Stdout))
	assert.NoError(t, err)
	// hast.Print(node)
	bnode, err := build.TranspileToVBScript(&config.VBScriptOptions{}, node)
	assert.NoError(t, err)
	vast.Print(bnode)
}

func TestEnumDecl(t *testing.T) {
	script := `
enum Status {
	Pending,
	Approved,
	Rejected
}

enum HttpCode {
    OK = 200,
    NotFound = 404,
    ServerError = 500,
    // 自动赋值为 501
    GatewayTimeout
}

enum Direction {
    North = "N",
    South = "S",
    East = "E",
    West = "W"
}

enum Config {
    RetryCount = 3,
    Timeout = "30s",
    EnableLogging = true
}

declare fn MsgBox(message: str);

MsgBox Status::Pending;
MsgBox HttpCode::OK;
MsgBox Direction::North;
MsgBox Config::RetryCount;
MsgBox Direction::North;`
	node, err := parser.ParseSourceScript(script, parser.OptionTracerASTTree(os.Stdout))
	assert.NoError(t, err)
	// hast.Print(node)
	bnode, err := build.TranspileToVBScript(&config.VBScriptOptions{}, node)
	assert.NoError(t, err)
	vast.Print(bnode)
}

func TestUnsafe(t *testing.T) {
	script := `
unsafe {
MsgBox "Hello, World!"
Dim count
}
		extern count: num
		echo $count
	`
	node, err := parser.ParseSourceScript(script, parser.OptionTracerASTTree(os.Stdout))
	assert.NoError(t, err)
	// hast.Print(node)
	bnode, err := build.TranspileToVBScript(&config.VBScriptOptions{}, node)
	assert.NoError(t, err)
	vast.Print(bnode)
}

func TestParseStringInterpolation(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []build.StringPart
	}{
		{
			name:  "simple variable",
			input: "$name",
			expected: []build.StringPart{
				{
					Text:       "",
					IsVariable: true,
					Expr:       &vast.Ident{Name: "name"},
				},
			},
		},
		{
			name:  "braced variable",
			input: "${name}",
			expected: []build.StringPart{
				{
					Text:       "",
					IsVariable: true,
					Expr:       &vast.Ident{Name: "name"},
				},
			},
		},
		{
			name:  "text only",
			input: "Hello, World!",
			expected: []build.StringPart{
				{
					Text:       "Hello, World!",
					IsVariable: false,
				},
			},
		},
		{
			name:  "text with simple variable",
			input: "Hello, $name!",
			expected: []build.StringPart{
				{
					Text:       "Hello, ",
					IsVariable: false,
				},
				{
					Text:       "",
					IsVariable: true,
					Expr:       &vast.Ident{Name: "name"},
				},
				{
					Text:       "!",
					IsVariable: false,
				},
			},
		},
		{
			name:  "text with braced variable",
			input: "Hello, ${name}!",
			expected: []build.StringPart{
				{
					Text:       "Hello, ",
					IsVariable: false,
				},
				{
					Text:       "",
					IsVariable: true,
					Expr:       &vast.Ident{Name: "name"},
				},
				{
					Text:       "!",
					IsVariable: false,
				},
			},
		},
		{
			name:  "multiple simple variables",
			input: "$firstName $lastName",
			expected: []build.StringPart{
				{
					Text:       "",
					IsVariable: true,
					Expr:       &vast.Ident{Name: "firstName"},
				},
				{
					Text:       " ",
					IsVariable: false,
				},
				{
					Text:       "",
					IsVariable: true,
					Expr:       &vast.Ident{Name: "lastName"},
				},
			},
		},
		{
			name:  "multiple braced variables",
			input: "${firstName} ${lastName}",
			expected: []build.StringPart{
				{
					Text:       "",
					IsVariable: true,
					Expr:       &vast.Ident{Name: "firstName"},
				},
				{
					Text:       " ",
					IsVariable: false,
				},
				{
					Text:       "",
					IsVariable: true,
					Expr:       &vast.Ident{Name: "lastName"},
				},
			},
		},
		{
			name:  "mixed simple and braced variables",
			input: "$firstName ${lastName}",
			expected: []build.StringPart{
				{
					Text:       "",
					IsVariable: true,
					Expr:       &vast.Ident{Name: "firstName"},
				},
				{
					Text:       " ",
					IsVariable: false,
				},
				{
					Text:       "",
					IsVariable: true,
					Expr:       &vast.Ident{Name: "lastName"},
				},
			},
		},
		{
			name:  "variable with underscore",
			input: "$user_name",
			expected: []build.StringPart{
				{
					Text:       "",
					IsVariable: true,
					Expr:       &vast.Ident{Name: "user_name"},
				},
			},
		},
		{
			name:  "variable with numbers",
			input: "$user123",
			expected: []build.StringPart{
				{
					Text:       "",
					IsVariable: true,
					Expr:       &vast.Ident{Name: "user123"},
				},
			},
		},
		{
			name:  "variable with mixed case",
			input: "$UserName",
			expected: []build.StringPart{
				{
					Text:       "",
					IsVariable: true,
					Expr:       &vast.Ident{Name: "UserName"},
				},
			},
		},
		{
			name:  "braced variable with underscore",
			input: "${user_name}",
			expected: []build.StringPart{
				{
					Text:       "",
					IsVariable: true,
					Expr:       &vast.Ident{Name: "user_name"},
				},
			},
		},
		{
			name:  "braced variable with numbers",
			input: "${user123}",
			expected: []build.StringPart{
				{
					Text:       "",
					IsVariable: true,
					Expr:       &vast.Ident{Name: "user123"},
				},
			},
		},
		{
			name:  "braced variable with mixed case",
			input: "${UserName}",
			expected: []build.StringPart{
				{
					Text:       "",
					IsVariable: true,
					Expr:       &vast.Ident{Name: "UserName"},
				},
			},
		},
		{
			name:  "complex interpolation",
			input: "Welcome, $firstName ${lastName}! Your ID is $user123.",
			expected: []build.StringPart{
				{
					Text:       "Welcome, ",
					IsVariable: false,
				},
				{
					Text:       "",
					IsVariable: true,
					Expr:       &vast.Ident{Name: "firstName"},
				},
				{
					Text:       " ",
					IsVariable: false,
				},
				{
					Text:       "",
					IsVariable: true,
					Expr:       &vast.Ident{Name: "lastName"},
				},
				{
					Text:       "! Your ID is ",
					IsVariable: false,
				},
				{
					Text:       "",
					IsVariable: true,
					Expr:       &vast.Ident{Name: "user123"},
				},
				{
					Text:       ".",
					IsVariable: false,
				},
			},
		},
		{
			name:  "dollar sign in text",
			input: "Price: $10.99",
			expected: []build.StringPart{
				{
					Text:       "Price: ",
					IsVariable: false,
				},
				{
					Text:       "",
					IsVariable: true,
					Expr:       &vast.Ident{Name: "10"},
				},
				{
					Text:       ".99",
					IsVariable: false,
				},
			},
		},
		{
			name:  "dollar sign followed by non-identifier",
			input: "Price: $!invalid",
			expected: []build.StringPart{
				{
					Text:       "Price: ",
					IsVariable: false,
				},
				{
					Text:       "",
					IsVariable: true,
					Expr:       &vast.Ident{Name: ""},
				},
				{
					Text:       "!invalid",
					IsVariable: false,
				},
			},
		},
		{
			name:  "incomplete braced variable",
			input: "Hello ${name",
			expected: []build.StringPart{
				{
					Text:       "Hello ",
					IsVariable: false,
				},
				{
					Text:       "",
					IsVariable: true,
					Expr:       &vast.Ident{Name: "name"},
				},
			},
		},
		{
			name:  "empty braced variable",
			input: "Hello ${}",
			expected: []build.StringPart{
				{
					Text:       "Hello ",
					IsVariable: false,
				},
				{
					Text:       "",
					IsVariable: true,
					Expr:       &vast.Ident{Name: ""},
				},
			},
		},
		{
			name:  "variable at start",
			input: "$name is here",
			expected: []build.StringPart{
				{
					Text:       "",
					IsVariable: true,
					Expr:       &vast.Ident{Name: "name"},
				},
				{
					Text:       " is here",
					IsVariable: false,
				},
			},
		},
		{
			name:  "variable at end",
			input: "Hello $name",
			expected: []build.StringPart{
				{
					Text:       "Hello ",
					IsVariable: false,
				},
				{
					Text:       "",
					IsVariable: true,
					Expr:       &vast.Ident{Name: "name"},
				},
			},
		},
		{
			name:  "braced variable at start",
			input: "${name} is here",
			expected: []build.StringPart{
				{
					Text:       "",
					IsVariable: true,
					Expr:       &vast.Ident{Name: "name"},
				},
				{
					Text:       " is here",
					IsVariable: false,
				},
			},
		},
		{
			name:  "braced variable at end",
			input: "Hello ${name}",
			expected: []build.StringPart{
				{
					Text:       "Hello ",
					IsVariable: false,
				},
				{
					Text:       "",
					IsVariable: true,
					Expr:       &vast.Ident{Name: "name"},
				},
			},
		},
		{
			name:  "consecutive variables",
			input: "$a$b$c",
			expected: []build.StringPart{
				{
					Text:       "",
					IsVariable: true,
					Expr:       &vast.Ident{Name: "a"},
				},
				{
					Text:       "",
					IsVariable: true,
					Expr:       &vast.Ident{Name: "b"},
				},
				{
					Text:       "",
					IsVariable: true,
					Expr:       &vast.Ident{Name: "c"},
				},
			},
		},
		{
			name:  "consecutive braced variables",
			input: "${a}${b}${c}",
			expected: []build.StringPart{
				{
					Text:       "",
					IsVariable: true,
					Expr:       &vast.Ident{Name: "a"},
				},
				{
					Text:       "",
					IsVariable: true,
					Expr:       &vast.Ident{Name: "b"},
				},
				{
					Text:       "",
					IsVariable: true,
					Expr:       &vast.Ident{Name: "c"},
				},
			},
		},
		{
			name:  "braced variable with special characters",
			input: "${user-name}",
			expected: []build.StringPart{
				{
					Text:       "",
					IsVariable: true,
					Expr:       &vast.Ident{Name: "user-name"},
				},
			},
		},
		{
			name:  "braced variable with spaces",
			input: "${user name}",
			expected: []build.StringPart{
				{
					Text:       "",
					IsVariable: true,
					Expr:       &vast.Ident{Name: "user name"},
				},
			},
		},
		{
			name:  "braced variable with dots",
			input: "${user.name}",
			expected: []build.StringPart{
				{
					Text:       "",
					IsVariable: true,
					Expr:       &vast.Ident{Name: "user.name"},
				},
			},
		},
		{
			name:     "empty string",
			input:    "",
			expected: []build.StringPart{},
		},
		{
			name:  "single dollar sign",
			input: "$",
			expected: []build.StringPart{
				{
					Text:       "",
					IsVariable: true,
					Expr:       &vast.Ident{Name: ""},
				},
			},
		},
		{
			name:  "dollar sign at end",
			input: "Hello $",
			expected: []build.StringPart{
				{
					Text:       "Hello ",
					IsVariable: false,
				},
				{
					Text:       "",
					IsVariable: true,
					Expr:       &vast.Ident{Name: ""},
				},
			},
		},
		{
			name:  "dollar sign followed by space",
			input: "Hello $ world",
			expected: []build.StringPart{
				{
					Text:       "Hello ",
					IsVariable: false,
				},
				{
					Text:       "",
					IsVariable: true,
					Expr:       &vast.Ident{Name: ""},
				},
				{
					Text:       " world",
					IsVariable: false,
				},
			},
		},
		{
			name:  "dollar sign followed by punctuation",
			input: "Hello $, world",
			expected: []build.StringPart{
				{
					Text:       "Hello ",
					IsVariable: false,
				},
				{
					Text:       "",
					IsVariable: true,
					Expr:       &vast.Ident{Name: ""},
				},
				{
					Text:       ", world",
					IsVariable: false,
				},
			},
		},
		{
			name:  "dollar sign followed by valid identifier",
			input: "Hello $a, world",
			expected: []build.StringPart{
				{
					Text:       "Hello ",
					IsVariable: false,
				},
				{
					Text:       "",
					IsVariable: true,
					Expr:       &vast.Ident{Name: "a"},
				},
				{
					Text:       ", world",
					IsVariable: false,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := build.ParseStringInterpolation(tt.input)

			// Check length
			assert.Equal(t, len(tt.expected), len(result), "Expected %d parts, got %d", len(tt.expected), len(result))

			// Check each part
			for i, expectedPart := range tt.expected {
				if i >= len(result) {
					t.Errorf("Missing part %d", i)
					continue
				}

				actualPart := result[i]

				// Check text
				assert.Equal(t, expectedPart.Text, actualPart.Text, "Part %d text mismatch", i)

				// Check IsVariable flag
				assert.Equal(t, expectedPart.IsVariable, actualPart.IsVariable, "Part %d IsVariable flag mismatch", i)

				// Check expression if it's a variable
				if expectedPart.IsVariable {
					assert.NotNil(t, actualPart.Expr, "Part %d should have expression", i)

					expectedIdent, ok := expectedPart.Expr.(*vast.Ident)
					assert.True(t, ok, "Part %d expected expression to be *vast.Ident", i)

					actualIdent, ok := actualPart.Expr.(*vast.Ident)
					assert.True(t, ok, "Part %d actual expression should be *vast.Ident", i)

					assert.Equal(t, expectedIdent.Name, actualIdent.Name, "Part %d variable name mismatch", i)
				} else {
					assert.Nil(t, actualPart.Expr, "Part %d should not have expression", i)
				}
			}
		})
	}
}

func TestModuleImport(t *testing.T) {
	// 创建内存文件系统
	memFS := memvfs.New()

	// 创建测试模块文件
	moduleContent := `
		pub fn Add(a: num, b: num) {
			return $a + $b
		}

		pub fn Multiply(a: num, b: num) {
			return $a * $b
		}

		const PI = 3.14159
	`

	// 写入模块文件
	err := memFS.WriteFile("utils.hl", []byte(moduleContent), 0644)
	assert.NoError(t, err)

	// 创建主文件
	mainContent := `
		import "utils"

		let result = Add(5, 7)
		echo $result
	`

	// 解析主文件
	node, err := parser.ParseSourceScript(mainContent)
	assert.NoError(t, err)

	// 使用模块管理器进行转换
	bnode, err := build.TranspileToVBScriptWithModules(&config.VBScriptOptions{}, node, memFS, ".")
	assert.NoError(t, err)

	// 打印结果
	vast.Print(bnode)
}

func TestTranspileWithModules(t *testing.T) {
	memFS := memvfs.New()

	testFiles := map[string]string{
		"std/unsafe/vbs/fmt.hl": `
			declare fn Print(message: str)
		`,
		"std/unsafe/vbs/type.hl": `
			declare {
				type String = str
			}

			import * from "fmt"
			declare fn CvtString(s: str) -> String
		`,
		"std/unsafe/vbs/io.hl": `
			import * from "type"
			// MsgBox function
			declare fn MsgBox(message: String)
		`,
		"std/unsafe/vbs/index.hl": `
			import { Import } from "import.vbs"
			import * from "io"

			declare fn Import(path: str)
		`,
		"import.vbs": `
Function Import(modulePath)
	' Runtime import function
	Set fso = CreateObject("Scripting.FileSystemObject")
	If fso.FileExists(modulePath) Then
		ExecuteGlobal fso.OpenTextFile(modulePath, 1).ReadAll()
	End If
End Function
`,
		"math.hl": `
			pub fn Add(a: num, b: num) -> num {
				return $a + $b
			}

			pub fn Multiply(a: num, b: num) -> num {
				return $a * $b
			}

			pub const PI = 3.14159
		`,
		"utils.hl": `
			import "math"

			pub fn Calculate(a: num, b: num) -> num {
				let sum = math.Add($a, $b)
				let product = math.Multiply($a, $b)
				return $sum + $product
			}

			// pub fn PI() -> num {
			// 	return math.$PI
			// }
		`,
		"main.hl": `
			import "utils"
			import * from "utils"

			let result = utils.Calculate(5, 3)
			MsgBox $result;
			MsgBox Calculate(5, 3);
			CvtString("Hello, World!")
			Print("Hello, World!")
		`,
	}

	for path, content := range testFiles {
		err := memFS.WriteFile(path, []byte(content), 0644)
		assert.NoError(t, err)
	}

	results, err := build.Transpile(&config.VBScriptOptions{}, "main.hl", memFS, ".", ".")
	assert.NoError(t, err)

	assert.NotNil(t, results)
	assert.Contains(t, results, "main.hl")
	assert.Contains(t, results, "utils.hl")
	assert.Contains(t, results, "math.hl")

	for file, code := range results {
		file = strings.Replace(file, ".hl", ".vbs", 1)
		fmt.Printf("=== %s ===\n", file)
		fmt.Println(code)
		fmt.Println()
	}
}

func TestSimpleMultiplication(t *testing.T) {
	script := `$a * $b`
	node, err := parser.ParseSourceScript(script)
	assert.NoError(t, err)

	bnode, err := build.TranspileToVBScript(&config.VBScriptOptions{}, node)
	assert.NoError(t, err)
	vast.Print(bnode)
}
