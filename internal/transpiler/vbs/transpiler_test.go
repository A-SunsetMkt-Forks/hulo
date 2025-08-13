// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package vbs_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hulo-lang/hulo/internal/config"
	build "github.com/hulo-lang/hulo/internal/transpiler/vbs"
	"github.com/hulo-lang/hulo/internal/vfs/memvfs"
	"github.com/stretchr/testify/assert"
)

func transpileToVBScript(t *testing.T, script string) {
	fs := memvfs.New()
	fs.WriteFile("test.hl", []byte(script), 0644)

	// hast.Print(node)
	bnode, err := build.Transpile(&config.Huloc{Main: "test.hl", HuloPath: ".", CompilerOptions: config.CompilerOptions{VBScript: &config.VBScriptOptions{CommentSyntax: "'"}}}, fs, ".")
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
	results, err := build.Transpile(&config.Huloc{Main: "main.hl", HuloPath: ".", CompilerOptions: config.CompilerOptions{VBScript: &config.VBScriptOptions{CommentSyntax: "'"}}}, memFS, ".")
	assert.NoError(t, err)

	for _, v := range results {
		fmt.Println(v)
	}
}

func TestTranspileWithModules(t *testing.T) {
	memFS := memvfs.New()

	testFiles := map[string]string{
		"core/unsafe/vbs/fmt.hl": `
			declare fn Print(message: str)
		`,
		"core/unsafe/vbs/type.hl": `
			declare {
				type String = str
			}

			import * from "fmt"
			declare fn CvtString(s: str) -> String
		`,
		"core/unsafe/vbs/io.hl": `
			import * from "type"
			// MsgBox function
			declare fn MsgBox(message: String)
		`,
		"core/unsafe/vbs/index.hl": `
			import { Import } from "import.vbs"
			import * from "io"

			declare fn Import(path: str)
		`,
		"core/unsafe/vbs/import.vbs": `
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
			import * from "core/unsafe/vbs/index"

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

	results, err := build.Transpile(&config.Huloc{Main: "main.hl", HuloPath: ".", CompilerOptions: config.CompilerOptions{VBScript: &config.VBScriptOptions{CommentSyntax: "'"}}}, memFS, ".")
	assert.NoError(t, err)

	for k, v := range results {
		if strings.HasPrefix(k, "core") {
			continue
		}
		fmt.Println(k)
		fmt.Println("--------------------------------")
		fmt.Println(v)
	}
}

// TODO: 实现的思路
// 需要对每个函数做混淆，基于 module_function_name 做混淆
// 解决函数的时候 findFunction 根据 module, function_name 查询出实际调用的名字
func TestTranspileModuleAlias(t *testing.T) {
	fs := memvfs.New()
	scritps := map[string]string{
		"utils.hl": `
			pub fn calculate(a: num, b: num) -> num {
				return $a + $b;
			}
		`,
		"main.hl": `
			import "utils" as u
			// import { calculate } from "utils"
			import { calculate as c } from "utils"

			declare fn MsgBox(message: str);

			fn calculate(a: num, b: num) -> num {
				return $a + $b;
			}

			let result = u.calculate(5, 3);
			MsgBox $result calculate(5, 3);
			MsgBox c(5, 3);
		`,
	}

	for path, content := range scritps {
		err := fs.WriteFile(path, []byte(content), 0644)
		assert.NoError(t, err)
	}

	results, err := build.Transpile(&config.Huloc{Main: "main.hl", HuloPath: ".", CompilerOptions: config.CompilerOptions{VBScript: &config.VBScriptOptions{CommentSyntax: "'"}}}, fs, ".")
	assert.NoError(t, err)

	for k, v := range results {
		fmt.Println(k)
		fmt.Println("--------------------------------")
		fmt.Println(v)
	}
}

func TestSimpleMultiplication(t *testing.T) {
	script := `$a * $b`
	transpileToVBScript(t, script)
}
