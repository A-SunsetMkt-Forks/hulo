package transpiler_test

import (
	"fmt"
	"testing"

	"github.com/hulo-lang/hulo/internal/config"
	build "github.com/hulo-lang/hulo/internal/transpiler/powershell"
	"github.com/hulo-lang/hulo/internal/vfs/memvfs"
	"github.com/stretchr/testify/assert"
)

func TestBuildModule(t *testing.T) {
	fs := memvfs.New()
	testFiles := map[string]string{
		"core/unsafe/bash/index.hl": `
declare fn echo(message: str);`,
		"my_script.hl": `
			pub fn greet() {
				echo "Hello, World!"
			}
		`,
		"main.hl": `
			import "my_script"

			my_script.greet()
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

func TestConvertSelectExpr(t *testing.T) {
	fs := memvfs.New()
	testFiles := map[string]string{
		"math.hl": `
// Math module for Bash

const PI=3.14159

pub fn sqrt(x: num) -> num {
	return $x * $x
}`,
		"main.hl": `
			import "math"

			echo($PI)
			echo(math.sqrt(16))

			$p := "hello"
			echo($p)
			echo($p.length())
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

func TestTanspileLoop(t *testing.T) {
	fs := memvfs.New()
	testFiles := map[string]string{
		"main.hl": `
			loop $i := 0; $i < 10; $i++ {
				echo $i
			}
			let a = 0
			loop $a < 2 {
				echo $a;
				$a++;
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

func TestTanspileScope(t *testing.T) {
	fs := memvfs.New()
	testFiles := map[string]string{
		"main.hl": `
			let a = 10
			var b = 3.14
			const c = "Hello World!"
			$d := true
			// {
			// 	let a = 3
			// }
			fn test() {
				let a = 5
				echo $a
			}
			echo $a
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

func TestTanspileIf(t *testing.T) {
	fs := memvfs.New()
	testFiles := map[string]string{
		"main.hl": `
			let a = read("Input a:")
			let b = read("Input b:")

			if $a > 10 {
				echo "a is greater than 10"
				if $b > 20 {
					echo "b is greater than 20"
				} else {
					echo "b is less than or equal to 20"
				}
			} else if $a < 0 {
				echo "a is less than 0"
			} else {
				echo "a is between 0 and 10"
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

			echo(math.div(10, 2))
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

func TestTranspileComment(t *testing.T) {
	fs := memvfs.New()
	testFiles := map[string]string{
		"main.hl": `
// this is a comment

/*
 * this
 *    is a multi
 * comment
 */

/**
 * this
 *    is a multi
 * comment
 */`,
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

func TestTranspileFunc(t *testing.T) {
	fs := memvfs.New()
	testFiles := map[string]string{
		"main.hl": `
fn helloWorld() {
	echo "Hello, World!"
}

fn sayHello(name: str) -> void {
	echo "Hello, $name!"
}

fn add(a: num, b: num) => $a + $b

fn multiply(a: num, b: num) -> num {
	return $a * $b
}

helloWorld;
sayHello "Hulo";
echo add(1, 2);
echo multiply(1, 2);
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

func TestTranspileForIn(t *testing.T) {
	fs := memvfs.New()
	testFiles := map[string]string{
		"main.hl": `
let arr: list<num> = [1, 2, 3, 4, 5]

loop $item in $arr {
    echo $item
}

loop $i in [0, 1, 2] {
	echo $i
}`,
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

func TestTranspileForOf(t *testing.T) {
	fs := memvfs.New()
	testFiles := map[string]string{
		"main.hl": `
let config: map<str, str> = {"host": "localhost", "port": "8080"}
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
}`,
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

func TestTranspileDoWhile(t *testing.T) {
	fs := memvfs.New()
	testFiles := map[string]string{
		"main.hl": `
			let a = 15
			do {
				echo "Hello, World!";
				$a--;
			} loop ($a > 0);
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

func TestTranspileClass(t *testing.T) {
	fs := memvfs.New()
	testFiles := map[string]string{
		"main.hl": `
			class User {
				pub name: str
				pub age: num

				pub fn to_str() -> str {
					return "User(name: $name, age: $age)"
				}

				pub fn greet(other: str) {
					echo "Hello, $other! I'm $name."
				}
			}

			let u = User("John", 20)
			echo $u.to_str();
			$u.greet("Jane");
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
