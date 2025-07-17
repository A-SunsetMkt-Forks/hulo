// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package build_test

import (
	"fmt"
	"strings"
	"testing"

	build "github.com/hulo-lang/hulo/internal/build/vbs"
	"github.com/hulo-lang/hulo/internal/config"
	"github.com/hulo-lang/hulo/internal/vfs/memvfs"
	vast "github.com/hulo-lang/hulo/syntax/vbs/ast"
	"github.com/stretchr/testify/assert"
)

func transpileToVBScript(t *testing.T, script string) {
	fs := memvfs.New()
	fs.WriteFile("test.hl", []byte(script), 0644)

	// hast.Print(node)
	bnode, err := build.Transpile(&config.Huloc{Main: "test.hl"}, fs, ".", ".")
	assert.NoError(t, err)

	for _, v := range bnode {
		fmt.Println(v)
	}
}

func TestCommandStmt(t *testing.T) {
	script := `echo "Hello, World!" 3.14 true`
	transpileToVBScript(t, script)
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
	transpileToVBScript(t, script)
}

func TestAssign(t *testing.T) {
	script := `let a = 10
	var b = 3.14
	const c = "Hello, World!"
	$d := true`
	transpileToVBScript(t, script)
}

func TestIf(t *testing.T) {
	script := `if $a > 10 {
		echo "a is greater than 10"
	} else {
		echo "a is less than or equal to 10"
	}`
	transpileToVBScript(t, script)
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
	transpileToVBScript(t, script)
}

func TestStandardWhile(t *testing.T) {
	script := `loop $a < 2 { MsgBox $a; $a++; }`
	transpileToVBScript(t, script)
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
	transpileToVBScript(t, script)
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
	transpileToVBScript(t, script)
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
	transpileToVBScript(t, script)
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
	transpileToVBScript(t, script)
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
	transpileToVBScript(t, script)
}

func TestStringInterpolation(t *testing.T) {
	script := `echo "Hello, $name!"`
	transpileToVBScript(t, script)
}

func TestStringEscape(t *testing.T) {
	script := `echo "Hello, \"World\"!"`
	transpileToVBScript(t, script)
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
	transpileToVBScript(t, script)
}

func TestFieldModifiers(t *testing.T) {
	script := `
		class TestClass {
			name: str
			pub age: num
			pub className: str = "TestClass"
		}
	`
	transpileToVBScript(t, script)
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
	transpileToVBScript(t, script)
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
	transpileToVBScript(t, script)
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
	transpileToVBScript(t, script)
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
	transpileToVBScript(t, script)
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

	scripts := map[string]string{
		"utils.hl": `
			pub fn Add(a: num, b: num) {
				return $a + $b
			}

			pub fn Multiply(a: num, b: num) {
				return $a * $b
			}

			const PI = 3.14159
		`,
		"main.hl": `
			import "utils"

			let result = Add(5, 7)
			echo $result
		`,
	}

	for path, content := range scripts {
		err := memFS.WriteFile(path, []byte(content), 0644)
		assert.NoError(t, err)
	}

	// 使用模块管理器进行转换
	results, err := build.Transpile(&config.Huloc{Main: "main.hl", CompilerOptions: config.CompilerOptions{VBScript: &config.VBScriptOptions{CommentSyntax: "'"}}}, memFS, ".", ".")
	assert.NoError(t, err)

	for _, v := range results {
		fmt.Println(v)
	}
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
		"std/unsafe/vbs/import.vbs": `
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

			const PI = 3.14159
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
			import * from "std/unsafe/vbs/index"

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

	results, err := build.Transpile(&config.Huloc{Main: "main.hl", CompilerOptions: config.CompilerOptions{VBScript: &config.VBScriptOptions{CommentSyntax: "'"}}}, memFS, ".", ".")
	assert.NoError(t, err)

	for k, v := range results {
		if strings.HasPrefix(k, "std") {
			continue
		}
		fmt.Println(k)
		fmt.Println("--------------------------------")
		fmt.Println(v)
	}
}

func TestSimpleMultiplication(t *testing.T) {
	script := `$a * $b`
	transpileToVBScript(t, script)
}
