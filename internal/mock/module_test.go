// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package mock

import (
	"testing"
)

func TestNewMockModule(t *testing.T) {
	module := NewMockModule("test-module", "1.0.0")

	if module == nil {
		t.Fatal("NewMockModule returned nil")
	}

	if module.Name != "test-module" {
		t.Errorf("Expected name 'test-module', got '%s'", module.Name)
	}

	if module.Version != "1.0.0" {
		t.Errorf("Expected version '1.0.0', got '%s'", module.Version)
	}

	if module.Dependencies == nil {
		t.Error("Dependencies map is nil")
	}

	if module.Files == nil {
		t.Error("Files map is nil")
	}

	if module.Config == nil {
		t.Error("Config is nil")
	}
}

func TestMockModule_AddFile(t *testing.T) {
	module := NewMockModule("test-module", "1.0.0")

	// 添加文件
	module.AddFile("index.hl", "pub fn main() => println('hello')")

	// 检查文件是否存在
	if content, exists := module.Files["index.hl"]; !exists {
		t.Error("File should exist after AddFile")
	} else if content != "pub fn main() => println('hello')" {
		t.Errorf("Expected content 'pub fn main() => println('hello')', got '%s'", content)
	}

	// 添加多个文件
	module.AddFile("utils.hl", "pub fn helper() => 'helper'")

	if len(module.Files) != 2 {
		t.Errorf("Expected 2 files, got %d", len(module.Files))
	}
}

func TestMockModule_AddDependency(t *testing.T) {
	module := NewMockModule("test-module", "1.0.0")

	// 添加依赖
	module.AddDependency("math", "1.0.0")
	module.AddDependency("utils", "2.0.0")

	// 检查依赖
	if version, exists := module.Dependencies["math"]; !exists {
		t.Error("Dependency 'math' should exist")
	} else if version != "1.0.0" {
		t.Errorf("Expected version '1.0.0' for math, got '%s'", version)
	}

	if version, exists := module.Dependencies["utils"]; !exists {
		t.Error("Dependency 'utils' should exist")
	} else if version != "2.0.0" {
		t.Errorf("Expected version '2.0.0' for utils, got '%s'", version)
	}

	if len(module.Dependencies) != 2 {
		t.Errorf("Expected 2 dependencies, got %d", len(module.Dependencies))
	}
}

func TestMockModule_Config(t *testing.T) {
	module := NewMockModule("test-module", "1.0.0")

	// 设置配置
	module.Config.Main = "./main.hl"
	module.Config.Dependencies = map[string]string{
		"math": "1.0.0",
	}

	// 检查配置
	if module.Config.Main != "./main.hl" {
		t.Errorf("Expected main './main.hl', got '%s'", module.Config.Main)
	}

	if module.Config.Version != "1.0.0" {
		t.Errorf("Expected version '1.0.0', got '%s'", module.Config.Version)
	}

	if len(module.Config.Dependencies) != 1 {
		t.Errorf("Expected 1 dependency in config, got %d", len(module.Config.Dependencies))
	}
}
