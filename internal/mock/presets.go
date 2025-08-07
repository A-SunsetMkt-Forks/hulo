// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package mock

// MockPresets 提供常用的测试场景
type MockPresets struct{}

// NewMockPresets 创建 Mock 预设实例
func NewMockPresets() *MockPresets {
	return &MockPresets{}
}

// SimpleMathProject 创建简单的数学模块项目
func (mp *MockPresets) SimpleMathProject() *MockProject {
	return NewMockBuilder().
		WithName("simple-math").
		WithModule("math", "1.0.0").
		WithFile("index.hl", `
pub fn add(a: num, b: num) => $a + $b
pub fn sub(a: num, b: num) => $a - $b
pub fn mul(a: num, b: num) => $a * $b
pub fn div(a: num, b: num) => $a / $b`).
		WithMain("./index.hl").
		End().
		Build()
}

// ComplexDependencyProject 创建复杂依赖项目
func (mp *MockPresets) ComplexDependencyProject() *MockProject {
	return NewMockBuilder().
		WithName("complex-deps").
		WithModule("utils", "1.0.0").
		WithFile("string.hl", `
pub fn concat(a: str, b: str) => $a + $b
pub fn upper(s: str) => $s.to_upper()`).
		WithFile("array.hl", `
pub fn join(arr: array<str>, sep: str) => $arr.join($sep)`).
		WithMain("./string.hl").
		End().
		WithModule("app", "1.0.0").
		WithFile("main.hl", `
import "utils"
import "utils/array"

fn main() {
	let result = utils.concat("Hello", "World")
	println(result)
}`).
		WithDependency("utils", "1.0.0").
		WithMain("./main.hl").
		End().
		Build()
}

// HelloWorldProject 创建 Hello World 项目
func (mp *MockPresets) HelloWorldProject() *MockProject {
	project := NewMockBuilder().
		WithName("hello-world").
		WithModule("hulo-lang/hello-world", "1.0.0").
		WithFile("index.hl", `
import "hulo-lang/hello-world"`).
		WithFile("utils_calculate.hl", `
pub fn add(a: num, b: num) -> num {
	return $a + $b
}

pub fn greet(name: str) -> str {
	return "Hello, " + $name + "!"
}`).
		WithMain("./index.hl").
		End().
		WithModule("core/math", "1.0.0").
		WithFile("index.hl", `
pub fn add(a: num, b: num) => $a + $b
pub fn mul(a: num, b: num) => $a * $b`).
		WithMain("./index.hl").
		End().
		Build()

	// 添加工作目录文件
	project.AddWorkDirFile("src/main.hl", `
import "hulo-lang/hello-world" as h
import "hulo-lang/hello-world/utils_calculate" as uc
import * from "math"
import "math"

fn main() {
	let result = uc.add(1, 2)
	let greeting = h.greet("World")
	println(result)
	println(greeting)
}`)

	// 设置项目配置
	project.SetProjectConfig(`
dependencies:
	hulo-lang/hello-world: 1.0.0
main: ./src/main.hl`)

	return project
}

// MultiVersionProject 创建多版本项目
func (mp *MockPresets) MultiVersionProject() *MockProject {
	return NewMockBuilder().
		WithName("multi-version").
		WithModule("lib", "1.0.0").
		WithFile("index.hl", `
pub fn version() => "1.0.0"
pub fn old_function() => "deprecated"`).
		WithMain("./index.hl").
		End().
		WithModule("lib", "2.0.0").
		WithFile("index.hl", `
pub fn version() => "2.0.0"
pub fn new_function() => "new feature"`).
		WithMain("./index.hl").
		End().
		WithModule("app", "1.0.0").
		WithFile("main.hl", `
import "lib"

fn main() {
	println(lib.version())
}`).
		WithDependency("lib", "1.0.0").
		WithMain("./main.hl").
		End().
		Build()
}

// CircularDependencyProject 创建循环依赖项目（用于测试错误处理）
func (mp *MockPresets) CircularDependencyProject() *MockProject {
	return NewMockBuilder().
		WithName("circular-deps").
		WithModule("module-a", "1.0.0").
		WithFile("index.hl", `
import "module-b"
pub fn a_function() => "from module a"`).
		WithDependency("module-b", "1.0.0").
		WithMain("./index.hl").
		End().
		WithModule("module-b", "1.0.0").
		WithFile("index.hl", `
import "module-a"
pub fn b_function() => "from module b"`).
		WithDependency("module-a", "1.0.0").
		WithMain("./index.hl").
		End().
		Build()
}
