// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package interpreter

import (
	"testing"

	"github.com/hulo-lang/hulo/internal/object"
	"github.com/hulo-lang/hulo/internal/vfs"
	"github.com/hulo-lang/hulo/internal/vfs/memvfs"
	"github.com/hulo-lang/hulo/syntax/hulo/ast"
	hast "github.com/hulo-lang/hulo/syntax/hulo/ast"
	"github.com/hulo-lang/hulo/syntax/hulo/parser"
	"github.com/hulo-lang/hulo/syntax/hulo/token"
	"github.com/stretchr/testify/assert"
)

func TestEval(t *testing.T) {
	script := `loop {
			echo "Hello, World!"
		}

		do {
			echo "Hello, World!"
		} loop ($a > 10)

		loop $i := 0; $i < 10; $i++ {
			echo "Hello, World!"
		}

		comptime {
			echo "Hello, World!"
		}`
	var node hast.Node
	var err error
	node, err = parser.ParseSourceScript(script)
	assert.NoError(t, err)
	hast.Print(node)
	interp := &Interpreter{}
	node = interp.Eval(node)
	hast.Print(node)
}

func TestEvalWithVFS(t *testing.T) {
	var memFS vfs.VFS = memvfs.New()
	files := map[string]string{
		"main.hl": `import { PI } from "math"

fn main() {
    echo "PI = " + $PI.to_str()
}`,
		"math.hl": `pub const PI = 3.14159

pub fn add(a: num, b: num) -> num {
    return $a + $b
}

pub fn multiply(a: num, b: num) -> num {
    return $a * $b
}`,
		"utils.hl": `pub fn format_string(str: str) -> str {
    return $str.to_upper()
}

pub fn reverse_string(str: str) -> str {
    // 简单的字符串反转
    $result = ""
    loop $i in $str.length() - 1 .. -1 .. -1 {
        $result += ${str[i]}
    }
    return $result
}`,
	}
	for filename, content := range files {
		err := memFS.WriteFile(filename, []byte(content), 0644)
		assert.NoError(t, err)
	}
	// interp := &Interpreter{
	// 	fs:  memFS,
	// 	cvt: &object.ASTConverter{},
	// }
	// err := interp.ExecuteMain("main.hl")
	// assert.NoError(t, err)
}

func TestEvalComptimeWhen(t *testing.T) {
	// log.SetLevel(log.DebugLevel)
	script := `comptime when $TARGET == "ps" {
		Write-Host "Hello, PowerShell"
	} else when $TARGET == "bat" {
		echo "Hello, Batch"
	} else when $TARGET == "bash" {
		echo "Hello, Bash"
	} else when $TARGET == "vbs" {
		MsgBox "Hello, VBScript"
	}`
	var node hast.Node
	var err error
	node, err = parser.ParseSourceScript(script)
	assert.NoError(t, err)
	// hast.Print(node)
	env := NewEnvironment()
	env.SetWithScope("TARGET", &object.StringValue{Value: "ps"}, token.CONST, true)
	interp := &Interpreter{
		env: env,
		cvt: &object.ASTConverter{},
	}
	node = interp.Eval(node)
	hast.Print(node)
}

func TestEvalComptimeAssign(t *testing.T) {
	// log.SetLevel(log.DebugLevel)
	script := `let x = comptime {
		let a = 1
		let b = $a + 2
		echo "b is" $b
		if $b == 3 {
			echo "b is 3"
			if $a == 1 {
				echo "a is 1";
				nop
			}
		} else {
			echo "b is not 3";
			nop
		}
		return $b
	}`
	var node hast.Node
	var err error
	node, err = parser.ParseSourceScript(script)
	assert.NoError(t, err)
	// hast.Print(node)
	interp := &Interpreter{
		env: NewEnvironment(),
		cvt: &object.ASTConverter{},
	}
	node = interp.Eval(node)
	hast.Print(node)
}

func TestEvalComptimeFor(t *testing.T) {
	// log.SetLevel(log.DebugLevel)
	script := `let x = comptime {
			let a = 0
			loop $i := 0; $i < 10; $i++ {
				echo "i is" $i;
				$a += $i;
			}
			return $a
		}`
	var node hast.Node
	var err error
	node, err = parser.ParseSourceScript(script)
	assert.NoError(t, err)
	// hast.Print(node)
	interp := &Interpreter{
		env: NewEnvironment(),
		cvt: &object.ASTConverter{},
	}
	node = interp.Eval(node)
	hast.Print(node)
}

func TestEvalComptimeWhile(t *testing.T) {
	// log.SetLevel(log.DebugLevel)
	script := `let x = comptime {
			let a = 0
			loop $a < 10 {
				echo "a is" $a;
				$a++;
			}
			return $a
		}`
	var node hast.Node
	var err error
	node, err = parser.ParseSourceScript(script)
	assert.NoError(t, err)
	// hast.Print(node)
	interp := &Interpreter{
		env: NewEnvironment(),
		cvt: &object.ASTConverter{},
	}
	node = interp.Eval(node)
	hast.Print(node)
}

func TestEvalComptimeDoWhile(t *testing.T) {
	// log.SetLevel(log.DebugLevel)
	script := `let x = comptime {
			let a = 15
			do {
				MsgBox "Hello, World!";
				$a--;
			} loop ($a > 10)
			return $a
		}`
	var node hast.Node
	var err error
	node, err = parser.ParseSourceScript(script)
	assert.NoError(t, err)
	// hast.Print(node)
	interp := &Interpreter{
		env: NewEnvironment(),
		cvt: &object.ASTConverter{},
	}
	node = interp.Eval(node)
	hast.Print(node)
}

func TestComptimeMatch(t *testing.T) {
	// log.SetLevel(log.DebugLevel)
	script := `let x = comptime {
		let a = 1
		match $a {
			1 => { echo "a is 1"; nop },
			2 => { echo "a is 2"; nop },
			_ => { echo "a is not 1 or 2"; nop },
		}
		return $a
	}`
	var node hast.Node
	var err error
	node, err = parser.ParseSourceScript(script)
	assert.NoError(t, err)
	// hast.Print(node)
	interp := &Interpreter{
		env: NewEnvironment(),
		cvt: &object.ASTConverter{},
	}
	node = interp.Eval(node)
	hast.Print(node)
}

func TestComptimeFunc(t *testing.T) {
	script := `let x = comptime {
		fn add(a: num, b: num) -> num {
			return $a + $b
		}

		fn add(a: str, b: str) -> str {
			return $a + $b
		}

		echo add(1, 2);
		echo add("Hello, ", "World!");
		return add("Hello, ", "World!")
	}`
	var node hast.Node
	var err error
	node, err = parser.ParseSourceScript(script)
	assert.NoError(t, err)
	// hast.Print(node)
	interp := &Interpreter{
		env: NewEnvironment(),
		cvt: &object.ASTConverter{},
	}
	node = interp.Eval(node)
	hast.Print(node)
}

func TestComptimeComment(t *testing.T) {
	script := `let x = comptime {
		// this is a comment
		1 + 2
	}`
	var node hast.Node
	var err error
	node, err = parser.ParseSourceScript(script)
	assert.NoError(t, err)
	// hast.Print(node)
	interp := &Interpreter{
		env: NewEnvironment(),
		cvt: &object.ASTConverter{},
	}
	node = interp.Eval(node)
	hast.Print(node)
}

func TestClassDecl(t *testing.T) {
	script := `
comptime {
	pub class Person {
		name: str
		age: num = 0
		pub className: str = "Person"

		pub fn get_name() => $name
		pub fn get_age() => $age

		pub fn to_str() -> str {
			return "Person(name: $name, age: $age)"
		}

		pub fn greet(name: str) {
			echo "Hello $name, my name is ${this.name}"
		}
	}

	let p = Person()
	echo $p.to_str();
	$p.greet("John")
}`
	var node hast.Node
	var err error
	node, err = parser.ParseSourceScript(script)
	assert.NoError(t, err)
	// hast.Print(node)
	interp := &Interpreter{
		env:     NewEnvironment(),
		cvt:     &object.ASTConverter{},
		fileSet: ast.NewFileSet(),
	}
	node = interp.Eval(node)
	hast.Print(node)
}
