// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package mock

type MockModule struct {
	Name         string
	Version      string
	Dependencies map[string]string
	Files        map[string]string
	Config       *MockModuleConfig
}

type MockModuleConfig struct {
	Main         string            `yaml:"main"`
	Dependencies map[string]string `yaml:"dependencies"`
	Version      string            `yaml:"version"`
}

// 创建 Mock 模块
func NewMockModule(name, version string) *MockModule {
	return &MockModule{
		Name:         name,
		Version:      version,
		Dependencies: make(map[string]string),
		Files:        make(map[string]string),
		Config: &MockModuleConfig{
			Version: version,
		},
	}
}

// 添加文件到模块
func (mm *MockModule) AddFile(relativePath, content string) {
	mm.Files[relativePath] = content
}

// 设置依赖
func (mm *MockModule) AddDependency(name, version string) {
	mm.Dependencies[name] = version
}
