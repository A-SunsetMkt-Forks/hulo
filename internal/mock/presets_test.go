// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package mock

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestNewMockPresets(t *testing.T) {
	presets := NewMockPresets()

	if presets == nil {
		t.Fatal("NewMockPresets returned nil")
	}
}

func TestMockPresets_SimpleMathProject(t *testing.T) {
	presets := NewMockPresets()
	project := presets.SimpleMathProject()

	if project == nil {
		t.Fatal("SimpleMathProject returned nil")
	}

	if project.Name != "simple-math" {
		t.Errorf("Expected name 'simple-math', got '%s'", project.Name)
	}

	// 检查模块
	if mathModule, exists := project.GetModule("math"); !exists {
		t.Error("Math module should exist")
	} else {
		if mathModule.Version != "1.0.0" {
			t.Errorf("Expected version '1.0.0', got '%s'", mathModule.Version)
		}

		// 检查文件内容
		if content, exists := mathModule.Files["index.hl"]; !exists {
			t.Error("Math module should have index.hl file")
		} else {
			expectedFunctions := []string{
				"pub fn add",
				"pub fn sub",
				"pub fn mul",
				"pub fn div",
			}
			for _, fn := range expectedFunctions {
				if !strings.Contains(content, fn) {
					t.Errorf("Math module should contain function: %s", fn)
				}
			}
		}
	}
}

func TestMockPresets_ComplexDependencyProject(t *testing.T) {
	presets := NewMockPresets()
	project := presets.ComplexDependencyProject()

	if project == nil {
		t.Fatal("ComplexDependencyProject returned nil")
	}

	if project.Name != "complex-deps" {
		t.Errorf("Expected name 'complex-deps', got '%s'", project.Name)
	}

	// 检查模块数量
	if len(project.Modules) != 2 {
		t.Errorf("Expected 2 modules, got %d", len(project.Modules))
	}

	// 检查 utils 模块
	if utilsModule, exists := project.GetModule("utils"); !exists {
		t.Error("Utils module should exist")
	} else {
		if content, exists := utilsModule.Files["string.hl"]; !exists {
			t.Error("Utils module should have string.hl file")
		} else {
			if !strings.Contains(content, "pub fn concat") {
				t.Error("Utils module should contain concat function")
			}
			if !strings.Contains(content, "pub fn upper") {
				t.Error("Utils module should contain upper function")
			}
		}

		if content, exists := utilsModule.Files["array.hl"]; !exists {
			t.Error("Utils module should have array.hl file")
		} else {
			if !strings.Contains(content, "pub fn join") {
				t.Error("Utils module should contain join function")
			}
		}
	}

	// 检查 app 模块
	if appModule, exists := project.GetModule("app"); !exists {
		t.Error("App module should exist")
	} else {
		if content, exists := appModule.Files["main.hl"]; !exists {
			t.Error("App module should have main.hl file")
		} else {
			if !strings.Contains(content, "import \"utils\"") {
				t.Error("App module should import utils")
			}
			if !strings.Contains(content, "fn main()") {
				t.Error("App module should have main function")
			}
		}

		if version, exists := appModule.Dependencies["utils"]; !exists {
			t.Error("App module should have utils dependency")
		} else if version != "1.0.0" {
			t.Errorf("Expected utils dependency version '1.0.0', got '%s'", version)
		}
	}
}

func TestMockPresets_HelloWorldProject(t *testing.T) {
	presets := NewMockPresets()
	project := presets.HelloWorldProject()

	if project == nil {
		t.Fatal("HelloWorldProject returned nil")
	}

	if project.Name != "hello-world" {
		t.Errorf("Expected name 'hello-world', got '%s'", project.Name)
	}

	// 检查模块数量
	if len(project.Modules) != 2 {
		t.Errorf("Expected 2 modules, got %d", len(project.Modules))
	}

	// 检查 hello-world 模块
	if helloModule, exists := project.GetModule("hulo-lang/hello-world"); !exists {
		t.Error("Hello world module should exist")
	} else {
		if content, exists := helloModule.Files["utils_calculate.hl"]; !exists {
			t.Error("Hello world module should have utils_calculate.hl file")
		} else {
			if !strings.Contains(content, "pub fn add") {
				t.Error("Hello world module should contain add function")
			}
			if !strings.Contains(content, "pub fn greet") {
				t.Error("Hello world module should contain greet function")
			}
		}
	}

	// 检查 math 模块
	if mathModule, exists := project.GetModule("core/math"); !exists {
		t.Error("Math module should exist")
	} else {
		if content, exists := mathModule.Files["index.hl"]; !exists {
			t.Error("Math module should have index.hl file")
		} else {
			if !strings.Contains(content, "pub fn add") {
				t.Error("Math module should contain add function")
			}
			if !strings.Contains(content, "pub fn mul") {
				t.Error("Math module should contain mul function")
			}
		}
	}

	// 检查工作目录文件
	expectedMainPath := filepath.Join(project.WorkDir, "src/main.hl")
	if !project.FileSystem.FileExists(expectedMainPath) {
		t.Error("Main file should exist in work directory")
	}

	// 检查项目配置
	expectedConfigPath := filepath.Join(project.WorkDir, "hulo.pkg.yaml")
	if !project.FileSystem.FileExists(expectedConfigPath) {
		t.Error("Project config should exist")
	}
}

func TestMockPresets_MultiVersionProject(t *testing.T) {
	presets := NewMockPresets()
	project := presets.MultiVersionProject()

	if project == nil {
		t.Fatal("MultiVersionProject returned nil")
	}

	if project.Name != "multi-version" {
		t.Errorf("Expected name 'multi-version', got '%s'", project.Name)
	}

	// 检查模块数量（应该有3个：lib@1.0.0, lib@2.0.0, app@1.0.0）
	if len(project.Modules) != 3 {
		t.Errorf("Expected 3 modules, got %d", len(project.Modules))
	}

	// 检查 lib@1.0.0 模块
	if lib1Module, exists := project.GetModule("lib"); exists {
		if lib1Module.Version != "1.0.0" {
			t.Errorf("Expected lib module version '1.0.0', got '%s'", lib1Module.Version)
		}
		if content, exists := lib1Module.Files["index.hl"]; !exists {
			t.Error("Lib module should have index.hl file")
		} else {
			if !strings.Contains(content, "version() => \"1.0.0\"") {
				t.Error("Lib module should contain version function returning 1.0.0")
			}
			if !strings.Contains(content, "old_function") {
				t.Error("Lib module should contain old_function")
			}
		}
	} else {
		t.Error("Lib module should exist")
	}

	// 检查 app 模块
	if appModule, exists := project.GetModule("app"); !exists {
		t.Error("App module should exist")
	} else {
		if version, exists := appModule.Dependencies["lib"]; !exists {
			t.Error("App module should have lib dependency")
		} else if version != "1.0.0" {
			t.Errorf("Expected lib dependency version '1.0.0', got '%s'", version)
		}
	}
}

func TestMockPresets_CircularDependencyProject(t *testing.T) {
	presets := NewMockPresets()
	project := presets.CircularDependencyProject()

	if project == nil {
		t.Fatal("CircularDependencyProject returned nil")
	}

	if project.Name != "circular-deps" {
		t.Errorf("Expected name 'circular-deps', got '%s'", project.Name)
	}

	// 检查模块数量
	if len(project.Modules) != 2 {
		t.Errorf("Expected 2 modules, got %d", len(project.Modules))
	}

	// 检查 module-a
	if moduleA, exists := project.GetModule("module-a"); !exists {
		t.Error("Module-a should exist")
	} else {
		if content, exists := moduleA.Files["index.hl"]; !exists {
			t.Error("Module-a should have index.hl file")
		} else {
			if !strings.Contains(content, "import \"module-b\"") {
				t.Error("Module-a should import module-b")
			}
			if !strings.Contains(content, "a_function") {
				t.Error("Module-a should contain a_function")
			}
		}

		if version, exists := moduleA.Dependencies["module-b"]; !exists {
			t.Error("Module-a should have module-b dependency")
		} else if version != "1.0.0" {
			t.Errorf("Expected module-b dependency version '1.0.0', got '%s'", version)
		}
	}

	// 检查 module-b
	if moduleB, exists := project.GetModule("module-b"); !exists {
		t.Error("Module-b should exist")
	} else {
		if content, exists := moduleB.Files["index.hl"]; !exists {
			t.Error("Module-b should have index.hl file")
		} else {
			if !strings.Contains(content, "import \"module-a\"") {
				t.Error("Module-b should import module-a")
			}
			if !strings.Contains(content, "b_function") {
				t.Error("Module-b should contain b_function")
			}
		}

		if version, exists := moduleB.Dependencies["module-a"]; !exists {
			t.Error("Module-b should have module-a dependency")
		} else if version != "1.0.0" {
			t.Errorf("Expected module-a dependency version '1.0.0', got '%s'", version)
		}
	}
}
