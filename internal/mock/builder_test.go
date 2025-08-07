// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package mock

import (
	"strings"
	"testing"
)

func TestNewMockBuilder(t *testing.T) {
	builder := NewMockBuilder()

	if builder == nil {
		t.Fatal("NewMockBuilder returned nil")
	}

	if builder.project == nil {
		t.Error("project is nil")
	}

	if builder.project.Name != "test-project" {
		t.Errorf("Expected project name 'test-project', got '%s'", builder.project.Name)
	}
}

func TestMockBuilder_WithName(t *testing.T) {
	builder := NewMockBuilder()

	builder.WithName("custom-project")

	if builder.project.Name != "custom-project" {
		t.Errorf("Expected name 'custom-project', got '%s'", builder.project.Name)
	}
}

func TestMockBuilder_WithWorkDir(t *testing.T) {
	builder := NewMockBuilder()

	builder.WithWorkDir("/custom/workdir")

	if builder.project.WorkDir != "/custom/workdir" {
		t.Errorf("Expected workdir '/custom/workdir', got '%s'", builder.project.WorkDir)
	}
}

func TestMockBuilder_WithHuloModules(t *testing.T) {
	builder := NewMockBuilder()

	builder.WithHuloModules("/custom/modules")

	if builder.project.HuloModules != "/custom/modules" {
		t.Errorf("Expected hulo modules '/custom/modules', got '%s'", builder.project.HuloModules)
	}
}

func TestMockBuilder_WithHuloPath(t *testing.T) {
	builder := NewMockBuilder()

	builder.WithHuloPath("/custom/path")

	if builder.project.HuloPath != "/custom/path" {
		t.Errorf("Expected hulo path '/custom/path', got '%s'", builder.project.HuloPath)
	}
}

func TestMockBuilder_WithModule(t *testing.T) {
	builder := NewMockBuilder()

	moduleBuilder := builder.WithModule("test-module", "1.0.0")

	if moduleBuilder == nil {
		t.Fatal("WithModule should return module builder")
	}

	if moduleBuilder.module == nil {
		t.Error("module is nil")
	}

	if moduleBuilder.module.Name != "test-module" {
		t.Errorf("Expected module name 'test-module', got '%s'", moduleBuilder.module.Name)
	}

	if moduleBuilder.module.Version != "1.0.0" {
		t.Errorf("Expected module version '1.0.0', got '%s'", moduleBuilder.module.Version)
	}

	// 检查模块是否被添加到项目
	if _, exists := builder.project.GetModule("test-module"); !exists {
		t.Error("Module should be added to project")
	}
}

func TestMockBuilder_Build(t *testing.T) {
	builder := NewMockBuilder()

	project := builder.Build()

	if project == nil {
		t.Fatal("Build should return project")
	}

	// 检查默认工作目录是否被设置
	if project.WorkDir == "" {
		t.Error("WorkDir should be set by default")
	}

	if project.HuloModules == "" {
		t.Error("HuloModules should be set by default")
	}

	if project.HuloPath == "" {
		t.Error("HuloPath should be set by default")
	}
}

func TestMockModuleBuilder_WithFile(t *testing.T) {
	builder := NewMockBuilder()
	moduleBuilder := builder.WithModule("test-module", "1.0.0")

	moduleBuilder.WithFile("index.hl", "pub fn main() => println('hello')")

	if content, exists := moduleBuilder.module.Files["index.hl"]; !exists {
		t.Error("File should be added to module")
	} else if content != "pub fn main() => println('hello')" {
		t.Errorf("Expected content 'pub fn main() => println('hello')', got '%s'", content)
	}

	// 测试链式调用
	result := moduleBuilder.WithFile("utils.hl", "pub fn helper() => 'helper'")
	if result != moduleBuilder {
		t.Error("WithFile should return self for chaining")
	}
}

func TestMockModuleBuilder_WithDependency(t *testing.T) {
	builder := NewMockBuilder()
	moduleBuilder := builder.WithModule("test-module", "1.0.0")

	moduleBuilder.WithDependency("math", "1.0.0")

	if version, exists := moduleBuilder.module.Dependencies["math"]; !exists {
		t.Error("Dependency should be added to module")
	} else if version != "1.0.0" {
		t.Errorf("Expected version '1.0.0', got '%s'", version)
	}

	// 测试链式调用
	result := moduleBuilder.WithDependency("utils", "2.0.0")
	if result != moduleBuilder {
		t.Error("WithDependency should return self for chaining")
	}
}

func TestMockModuleBuilder_WithMain(t *testing.T) {
	builder := NewMockBuilder()
	moduleBuilder := builder.WithModule("test-module", "1.0.0")

	moduleBuilder.WithMain("./main.hl")

	if moduleBuilder.module.Config.Main != "./main.hl" {
		t.Errorf("Expected main './main.hl', got '%s'", moduleBuilder.module.Config.Main)
	}

	// 测试链式调用
	result := moduleBuilder.WithMain("./index.hl")
	if result != moduleBuilder {
		t.Error("WithMain should return self for chaining")
	}
}

func TestMockModuleBuilder_End(t *testing.T) {
	builder := NewMockBuilder()
	moduleBuilder := builder.WithModule("test-module", "1.0.0")

	result := moduleBuilder.End()

	if result != builder {
		t.Error("End should return the parent builder")
	}
}

func TestMockBuilder_ComplexBuild(t *testing.T) {
	project := NewMockBuilder().
		WithName("complex-project").
		WithWorkDir("/custom/workdir").
		WithModule("math", "1.0.0").
		WithFile("index.hl", "pub fn add(a: num, b: num) => $a + $b").
		WithDependency("utils", "1.0.0").
		WithMain("./index.hl").
		End().
		WithModule("utils", "1.0.0").
		WithFile("string.hl", "pub fn concat(a: str, b: str) => $a + $b").
		WithMain("./string.hl").
		End().
		Build()

	// 检查项目名称
	if project.Name != "complex-project" {
		t.Errorf("Expected name 'complex-project', got '%s'", project.Name)
	}

	// 检查工作目录
	if project.WorkDir != "/custom/workdir" {
		t.Errorf("Expected workdir '/custom/workdir', got '%s'", project.WorkDir)
	}

	// 检查模块数量
	if len(project.Modules) != 2 {
		t.Errorf("Expected 2 modules, got %d", len(project.Modules))
	}

	// 检查 math 模块
	if mathModule, exists := project.GetModule("math"); !exists {
		t.Error("Math module should exist")
	} else {
		if mathModule.Version != "1.0.0" {
			t.Errorf("Expected math module version '1.0.0', got '%s'", mathModule.Version)
		}
		if content, exists := mathModule.Files["index.hl"]; !exists {
			t.Error("Math module should have index.hl file")
		} else if !strings.Contains(content, "pub fn add") {
			t.Error("Math module should contain add function")
		}
		if version, exists := mathModule.Dependencies["utils"]; !exists {
			t.Error("Math module should have utils dependency")
		} else if version != "1.0.0" {
			t.Errorf("Expected utils dependency version '1.0.0', got '%s'", version)
		}
	}

	// 检查 utils 模块
	if utilsModule, exists := project.GetModule("utils"); !exists {
		t.Error("Utils module should exist")
	} else {
		if utilsModule.Version != "1.0.0" {
			t.Errorf("Expected utils module version '1.0.0', got '%s'", utilsModule.Version)
		}
		if content, exists := utilsModule.Files["string.hl"]; !exists {
			t.Error("Utils module should have string.hl file")
		} else if !strings.Contains(content, "pub fn concat") {
			t.Error("Utils module should contain concat function")
		}
	}
}
