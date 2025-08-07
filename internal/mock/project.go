// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package mock

import (
	"fmt"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type MockProject struct {
	Name        string
	Modules     map[string]*MockModule
	FileSystem  *MockFileSystem
	WorkDir     string
	HuloModules string
	HuloPath    string
}

// 创建 Mock 项目
func NewMockProject(name string) *MockProject {
	return &MockProject{
		Name:       name,
		Modules:    make(map[string]*MockModule),
		FileSystem: NewMockFileSystem(),
	}
}

// 设置工作目录结构
func (mp *MockProject) SetupWorkDir() {
	mp.WorkDir = "./testdata"
	mp.HuloModules = filepath.Join(mp.WorkDir, "hulo_modules")
	mp.HuloPath = filepath.Join(mp.WorkDir, "hulo_path")

	// 创建目录结构
	mp.FileSystem.AddDir(mp.WorkDir)
	mp.FileSystem.AddDir(mp.HuloModules)
	mp.FileSystem.AddDir(mp.HuloPath)
}

// 添加模块到项目
func (mp *MockProject) AddModule(module *MockModule) {
	mp.Modules[module.Name] = module
}

// 添加模块文件到文件系统
func (mp *MockProject) AddModuleFiles(module *MockModule) {
	// 将模块文件添加到文件系统
	modulePath := filepath.Join(mp.HuloModules, module.Name, module.Version)
	for relativePath, content := range module.Files {
		fullPath := filepath.Join(modulePath, relativePath)
		mp.FileSystem.AddFile(fullPath, content)
	}

	// 添加配置文件
	configPath := filepath.Join(modulePath, "hulo.pkg.yaml")
	configContent := mp.generateModuleConfig(module)
	mp.FileSystem.AddFile(configPath, configContent)
}

// 生成模块配置文件内容
func (mp *MockProject) generateModuleConfig(module *MockModule) string {
	config := &MockModuleConfig{
		Main:         module.Config.Main,
		Dependencies: module.Dependencies,
		Version:      module.Version,
	}

	// 如果没有设置 Main，使用默认值
	if config.Main == "" {
		config.Main = "./index.hl"
	}

	yamlData, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Sprintf("version: %s\nmain: %s", module.Version, config.Main)
	}

	return string(yamlData)
}

// 获取模块
func (mp *MockProject) GetModule(name string) (*MockModule, bool) {
	module, exists := mp.Modules[name]
	return module, exists
}

// 添加工作目录文件
func (mp *MockProject) AddWorkDirFile(relativePath, content string) {
	fullPath := filepath.Join(mp.WorkDir, relativePath)
	mp.FileSystem.AddFile(fullPath, content)
}

// 添加 HuloPath 文件
func (mp *MockProject) AddHuloPathFile(relativePath, content string) {
	fullPath := filepath.Join(mp.HuloPath, relativePath)
	mp.FileSystem.AddFile(fullPath, content)
}

// 获取项目配置
func (mp *MockProject) GetProjectConfig() string {
	config := map[string]interface{}{
		"dependencies": make(map[string]string),
		"main":         "./src/main.hl",
	}

	// 收集所有模块的依赖
	for _, module := range mp.Modules {
		for depName, depVersion := range module.Dependencies {
			deps := config["dependencies"].(map[string]string)
			deps[depName] = depVersion
		}
	}

	yamlData, err := yaml.Marshal(config)
	if err != nil {
		return "main: ./src/main.hl"
	}

	return string(yamlData)
}

// 设置项目配置
func (mp *MockProject) SetProjectConfig(configContent string) {
	configPath := filepath.Join(mp.WorkDir, "hulo.pkg.yaml")
	mp.FileSystem.AddFile(configPath, configContent)
}
