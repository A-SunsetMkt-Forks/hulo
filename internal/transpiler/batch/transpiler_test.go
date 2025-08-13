// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package batch_test

import (
	"fmt"
	"testing"

	"github.com/hulo-lang/hulo/internal/config"
	build "github.com/hulo-lang/hulo/internal/transpiler/batch"
	"github.com/hulo-lang/hulo/internal/vfs/memvfs"
	"github.com/stretchr/testify/assert"
)

func TestCommandStmt(t *testing.T) {
	script := `echo "Hello, World!" 3.14 true`
	fs := memvfs.New()
	err := fs.WriteFile("main.hl", []byte(script), 0644)
	assert.NoError(t, err)

	// hast.Print(node)
	results, err := build.Transpile(&config.Huloc{Main: "./main.hl", CompilerOptions: config.DefaultCompilerOptions()}, fs, "")
	assert.NoError(t, err)

	assert.Equal(t, 1, len(results))
	assert.NotNil(t, results["main.bat"])
	assert.Contains(t, results["main.bat"], `@echo off`)
	assert.Contains(t, results["main.bat"], `echo Hello, World! 3.14 1`)
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
	fs := memvfs.New()
	err := fs.WriteFile("main.hl", []byte(script), 0644)
	assert.NoError(t, err)

	results, err := build.Transpile(&config.Huloc{Main: "./main.hl", CompilerOptions: config.DefaultCompilerOptions()}, fs, "")
	assert.NoError(t, err)

	for file, code := range results {
		fmt.Printf("=== %s ===\n", file)
		fmt.Println(code)
		fmt.Println()
	}
}

func TestLinker(t *testing.T) {
	fs := memvfs.New()
	testFiles := map[string]string{
		"math.bat": `REM HULO_LINK_BEGIN abs
:abs
setlocal enabledelayedexpansion

set "input=%~1"
set /a "absValue=%input%"

if !input! lss 0 (
    set /a "absValue=-!input!"
)

(endlocal & set %2=%absValue%)
exit /b
REM HULO_LINK_END`,
		"main.hl": `// @hulo:link math.bat abs`,
	}
	for path, content := range testFiles {
		err := fs.WriteFile(path, []byte(content), 0644)
		assert.NoError(t, err)
	}

	// hast.Print(node)
	results, err := build.Transpile(&config.Huloc{Main: "./main.hl", CompilerOptions: config.DefaultCompilerOptions()}, fs, "")
	assert.NoError(t, err)

	for file, code := range results {
		fmt.Printf("=== %s ===\n", file)
		fmt.Println(code)
		fmt.Println()
	}
}

func TestAssign(t *testing.T) {
	script := `let a = 10
	var b = 3.14
	const c = "Hello, World!"
	$d := true`
	fs := memvfs.New()
	err := fs.WriteFile("main.hl", []byte(script), 0644)
	assert.NoError(t, err)

	results, err := build.Transpile(&config.Huloc{Main: "./main.hl"}, fs, "")
	assert.NoError(t, err)

	for file, code := range results {
		fmt.Printf("=== %s ===\n", file)
		fmt.Println(code)
		fmt.Println()
	}
}

func TestIf(t *testing.T) {
	script := `
	$a := 20
	if $a > 10 {
		echo "a is greater than 10"
	} else {
		echo "a is less than or equal to 10"
	}`
	fs := memvfs.New()
	err := fs.WriteFile("main.hl", []byte(script), 0644)
	assert.NoError(t, err)

	results, err := build.Transpile(&config.Huloc{Main: "./main.hl", CompilerOptions: config.DefaultCompilerOptions()}, fs, "")
	assert.NoError(t, err)

	for file, code := range results {
		fmt.Printf("=== %s ===\n", file)
		fmt.Println(code)
		fmt.Println()
	}
}

func TestLoop(t *testing.T) {
	script := `loop $a < 10 {
		echo "Hello, World!"
	}

	do {
		if $a == 5 {
			continue
		}
		echo "Hello, World!"
	} loop ($a > 10)

	loop $i := 0; $i < 10; $i++ {
		echo "Hello, World!"
		if $i > 5 {
			break
		}
	}`
	fs := memvfs.New()
	err := fs.WriteFile("main.hl", []byte(script), 0644)
	assert.NoError(t, err)

	results, err := build.Transpile(&config.Huloc{Main: "./main.hl", CompilerOptions: config.DefaultCompilerOptions()}, fs, "")
	assert.NoError(t, err)

	for file, code := range results {
		fmt.Printf("=== %s ===\n", file)
		fmt.Println(code)
		fmt.Println()
	}
}

func TestTranspileForIn(t *testing.T) {
	script := `let arr: list<num> = [1, 2, 3, 4, 5]

loop $item in $arr {
    echo $item
}

loop $i in [0, 1, 2] {
	echo $i
}`
	fs := memvfs.New()
	err := fs.WriteFile("main.hl", []byte(script), 0644)
	assert.NoError(t, err)

	results, err := build.Transpile(&config.Huloc{Main: "./main.hl", CompilerOptions: config.DefaultCompilerOptions()}, fs, "")
	assert.NoError(t, err)

	for file, code := range results {
		fmt.Printf("=== %s ===\n", file)
		fmt.Println(code)
		fmt.Println()
	}
}

func TestTranspileUnsafe(t *testing.T) {
	script := `unsafe """echo Hello, World!"""`
	fs := memvfs.New()
	err := fs.WriteFile("main.hl", []byte(script), 0644)
	assert.NoError(t, err)

	results, err := build.Transpile(&config.Huloc{Main: "./main.hl", CompilerOptions: config.DefaultCompilerOptions()}, fs, "")
	assert.NoError(t, err)

	for file, code := range results {
		fmt.Printf("=== %s ===\n", file)
		fmt.Println(code)
		fmt.Println()
	}
}

func TestTranspileUnsafeWithVariable(t *testing.T) {
	script := `comptime let a = 10
	unsafe """{% loop i in [1, 2, 3] %} echo $i {{ a }} {% endloop %}"""
	`
	fs := memvfs.New()
	err := fs.WriteFile("main.hl", []byte(script), 0644)
	assert.NoError(t, err)

	results, err := build.Transpile(&config.Huloc{Main: "./main.hl", CompilerOptions: config.DefaultCompilerOptions()}, fs, "")
	assert.NoError(t, err)

	for file, code := range results {
		fmt.Printf("=== %s ===\n", file)
		fmt.Println(code)
		fmt.Println()
	}
}

func TestTranspileForOf(t *testing.T) {
	script := `let config: map<str, str> = {"host": "localhost", "port": "8080"}
loop ($key, $value) of $config {
    echo "$key = $value"
}
loop ($key, _) of $config {
    echo $key
}

loop (_, $value) of $config {
    echo $value
}

loop $key of $config {
    echo $key
}`
	fs := memvfs.New()
	err := fs.WriteFile("main.hl", []byte(script), 0644)
	assert.NoError(t, err)

	results, err := build.Transpile(&config.Huloc{Main: "main.hl", HuloPath: "."}, fs, ".")
	assert.NoError(t, err)

	for file, code := range results {
		fmt.Printf("=== %s ===\n", file)
		fmt.Println(code)
		fmt.Println()
	}
}

func TestTranspileImport(t *testing.T) {
	fs := memvfs.New()
	testFiles := map[string]string{
		"math.hl": `
			pub fn add(a: num, b: num) -> num {
				return $a + $b
			}

			pub fn sub_(a: num, b: num) -> num {
				return $a - $b
			}

			pub fn mul(a: num, b: num) -> num {
				return $a * $b
			}

			pub fn div(a: num, b: num) -> num {
				return $a / $b
			}
		`,
		"main.hl": `
			import "math"
			import * from "math"

			echo(math.add(1, 2))
			echo(div(10, 2))
		`,
	}

	for path, content := range testFiles {
		fs.WriteFile(path, []byte(content), 0644)
	}

	results, err := build.Transpile(&config.Huloc{Main: "main.hl", HuloPath: "."}, fs, ".")
	assert.NoError(t, err)

	for file, code := range results {
		fmt.Printf("=== %s ===\n", file)
		fmt.Println(code)
		fmt.Println()
	}
}

func TestTranspileMatch(t *testing.T) {
	fs := memvfs.New()
	testFiles := map[string]string{
		"main.hl": `
			let n = read("Input a number:");

			match $n {
				0 => echo "The number is 0.",
				1 => echo "The number is 1.",
				_ => echo "The number is not 0 or 1."
			}
		`,
	}

	for path, content := range testFiles {
		fs.WriteFile(path, []byte(content), 0644)
	}

	results, err := build.Transpile(&config.Huloc{Main: "main.hl", HuloPath: "."}, fs, ".")
	assert.NoError(t, err)

	for file, code := range results {
		fmt.Printf("=== %s ===\n", file)
		fmt.Println(code)
		fmt.Println()
	}
}
