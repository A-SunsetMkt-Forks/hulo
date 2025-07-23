package interpreter

import (
	"testing"

	"github.com/hulo-lang/hulo/internal/object"
	"github.com/hulo-lang/hulo/internal/vfs"
	"github.com/hulo-lang/hulo/internal/vfs/memvfs"
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

func TestEvalComptime(t *testing.T) {
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
