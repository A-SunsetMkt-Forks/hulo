// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package mock

type MockBuilder struct {
	project *MockProject
}

func NewMockBuilder() *MockBuilder {
	return &MockBuilder{
		project: NewMockProject("test-project"),
	}
}

// 链式调用构建
func (mb *MockBuilder) WithModule(name, version string) *MockModuleBuilder {
	module := NewMockModule(name, version)
	mb.project.AddModule(module)
	return &MockModuleBuilder{
		module:  module,
		builder: mb,
	}
}

func (mb *MockBuilder) WithWorkDir(path string) *MockBuilder {
	mb.project.WorkDir = path
	return mb
}

func (mb *MockBuilder) WithName(name string) *MockBuilder {
	mb.project.Name = name
	return mb
}

func (mb *MockBuilder) WithHuloModules(path string) *MockBuilder {
	mb.project.HuloModules = path
	return mb
}

func (mb *MockBuilder) WithHuloPath(path string) *MockBuilder {
	mb.project.HuloPath = path
	return mb
}

func (mb *MockBuilder) Build() *MockProject {
	// 确保工作目录已设置
	if mb.project.WorkDir == "" {
		mb.project.SetupWorkDir()
	}

	// 添加所有模块文件到文件系统（因为现在路径已经设置好了）
	for _, module := range mb.project.Modules {
		mb.project.AddModuleFiles(module)
	}

	return mb.project
}

// 模块构建器
type MockModuleBuilder struct {
	module  *MockModule
	builder *MockBuilder
}

func (mmb *MockModuleBuilder) WithFile(path, content string) *MockModuleBuilder {
	mmb.module.AddFile(path, content)
	return mmb
}

func (mmb *MockModuleBuilder) WithDependency(name, version string) *MockModuleBuilder {
	mmb.module.AddDependency(name, version)
	return mmb
}

func (mmb *MockModuleBuilder) WithMain(main string) *MockModuleBuilder {
	mmb.module.Config.Main = main
	return mmb
}

func (mmb *MockModuleBuilder) End() *MockBuilder {
	return mmb.builder
}
