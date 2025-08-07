// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package mock

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestNewMockProject(t *testing.T) {
	project := NewMockProject("test-project")

	if project == nil {
		t.Fatal("NewMockProject returned nil")
	}

	if project.Name != "test-project" {
		t.Errorf("Expected name 'test-project', got '%s'", project.Name)
	}

	if project.Modules == nil {
		t.Error("Modules map is nil")
	}

	if project.FileSystem == nil {
		t.Error("FileSystem is nil")
	}
}

func TestMockProject_SetupWorkDir(t *testing.T) {
	project := NewMockProject("test-project")
	project.SetupWorkDir()

	expectedWorkDir := "./testdata"
	if project.WorkDir != expectedWorkDir {
		t.Errorf("Expected WorkDir '%s', got '%s'", expectedWorkDir, project.WorkDir)
	}

	expectedHuloModules := filepath.Join(expectedWorkDir, "hulo_modules")
	if project.HuloModules != expectedHuloModules {
		t.Errorf("Expected HuloModules '%s', got '%s'", expectedHuloModules, project.HuloModules)
	}

	expectedHuloPath := filepath.Join(expectedWorkDir, "hulo_path")
	if project.HuloPath != expectedHuloPath {
		t.Errorf("Expected HuloPath '%s', got '%s'", expectedHuloPath, project.HuloPath)
	}

	// 检查目录是否被创建
	if !project.FileSystem.DirExists(expectedWorkDir) {
		t.Error("WorkDir should be created")
	}
	if !project.FileSystem.DirExists(expectedHuloModules) {
		t.Error("HuloModules should be created")
	}
	if !project.FileSystem.DirExists(expectedHuloPath) {
		t.Error("HuloPath should be created")
	}
}

func TestMockProject_AddModule(t *testing.T) {
	project := NewMockProject("test-project")
	project.SetupWorkDir()

	// 创建模块
	module := NewMockModule("test-module", "1.0.0")
	module.AddFile("index.hl", "pub fn main() => println('hello')")
	module.AddDependency("math", "1.0.0")

	// 添加模块到项目
	project.AddModule(module)

	project.AddModuleFiles(module)

	// 检查模块是否被添加
	if retrievedModule, exists := project.GetModule("test-module"); !exists {
		t.Error("Module should exist after AddModule")
	} else if retrievedModule != module {
		t.Error("Retrieved module should be the same as added module")
	}

	// 检查文件是否被添加到文件系统
	expectedFilePath := filepath.Join(project.HuloModules, "test-module", "1.0.0", "index.hl")
	if !project.FileSystem.FileExists(expectedFilePath) {
		t.Error("Module file should be added to filesystem")
	}

	// 检查配置文件是否被创建
	expectedConfigPath := filepath.Join(project.HuloModules, "test-module", "1.0.0", "hulo.pkg.yaml")
	if !project.FileSystem.FileExists(expectedConfigPath) {
		t.Error("Module config should be created")
	}

	// 检查配置文件内容
	configFile, exists := project.FileSystem.GetFile(expectedConfigPath)
	if !exists {
		t.Fatal("Config file should exist")
	}

	if !strings.Contains(configFile.Content, "version: 1.0.0") {
		t.Error("Config should contain version")
	}
	if !strings.Contains(configFile.Content, "main: ./index.hl") {
		t.Error("Config should contain main file")
	}
}

func TestMockProject_AddWorkDirFile(t *testing.T) {
	project := NewMockProject("test-project")
	project.SetupWorkDir()

	// 添加工作目录文件
	project.AddWorkDirFile("src/main.hl", "fn main() => println('hello')")

	// 检查文件是否被创建
	expectedPath := filepath.Join(project.WorkDir, "src/main.hl")
	if !project.FileSystem.FileExists(expectedPath) {
		t.Error("WorkDir file should be created")
	}

	// 检查文件内容
	file, exists := project.FileSystem.GetFile(expectedPath)
	if !exists {
		t.Fatal("File should exist")
	}

	if file.Content != "fn main() => println('hello')" {
		t.Errorf("Expected content 'fn main() => println('hello')', got '%s'", file.Content)
	}
}

func TestMockProject_AddHuloPathFile(t *testing.T) {
	project := NewMockProject("test-project")
	project.SetupWorkDir()

	// 添加 HuloPath 文件
	project.AddHuloPathFile("core/math/index.hl", "pub fn add(a: num, b: num) => $a + $b")

	// 检查文件是否被创建
	expectedPath := filepath.Join(project.HuloPath, "core/math/index.hl")
	if !project.FileSystem.FileExists(expectedPath) {
		t.Error("HuloPath file should be created")
	}

	// 检查文件内容
	file, exists := project.FileSystem.GetFile(expectedPath)
	if !exists {
		t.Fatal("File should exist")
	}

	if file.Content != "pub fn add(a: num, b: num) => $a + $b" {
		t.Errorf("Expected content 'pub fn add(a: num, b: num) => $a + $b', got '%s'", file.Content)
	}
}

func TestMockProject_GetProjectConfig(t *testing.T) {
	project := NewMockProject("test-project")
	project.SetupWorkDir()

	// 添加带依赖的模块
	module := NewMockModule("test-module", "1.0.0")
	module.AddDependency("math", "1.0.0")
	module.AddDependency("utils", "2.0.0")
	project.AddModule(module)

	// 获取项目配置
	config := project.GetProjectConfig()

	// 检查配置内容
	if !strings.Contains(config, "main: ./src/main.hl") {
		t.Error("Config should contain main file")
	}
	if !strings.Contains(config, "math: 1.0.0") {
		t.Error("Config should contain math dependency")
	}
	if !strings.Contains(config, "utils: 2.0.0") {
		t.Error("Config should contain utils dependency")
	}
}

func TestMockProject_SetProjectConfig(t *testing.T) {
	project := NewMockProject("test-project")
	project.SetupWorkDir()

	// 设置项目配置
	configContent := `
dependencies:
  math: 1.0.0
  utils: 2.0.0
main: ./src/main.hl`

	project.SetProjectConfig(configContent)

	// 检查配置文件是否被创建
	expectedPath := filepath.Join(project.WorkDir, "hulo.pkg.yaml")
	if !project.FileSystem.FileExists(expectedPath) {
		t.Error("Project config should be created")
	}

	// 检查文件内容
	file, exists := project.FileSystem.GetFile(expectedPath)
	if !exists {
		t.Fatal("Config file should exist")
	}

	if file.Content != configContent {
		t.Errorf("Expected config content '%s', got '%s'", configContent, file.Content)
	}
}
