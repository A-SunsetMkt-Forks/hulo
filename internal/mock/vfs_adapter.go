// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package mock

import (
	"os"
	"path/filepath"

	"github.com/hulo-lang/hulo/internal/vfs"
	"github.com/hulo-lang/hulo/internal/vfs/memvfs"
)

// VFSAdapter 将 Mock 项目转换为 VFS 实例
type VFSAdapter struct {
	project *MockProject
	vfs     vfs.VFS
}

// NewVFSAdapter 创建 VFS 适配器
func NewVFSAdapter(project *MockProject) *VFSAdapter {
	return &VFSAdapter{
		project: project,
		vfs:     memvfs.New(),
	}
}

// ToVFS 将 Mock 项目转换为 VFS 实例
func (va *VFSAdapter) ToVFS() vfs.VFS {
	// 将 Mock 文件系统中的所有文件写入 VFS
	for _, file := range va.project.FileSystem.GetAllFiles() {
		mockFile, exists := va.project.FileSystem.GetFile(file)
		if !exists {
			continue
		}

		// 标准化路径分隔符
		normalizedPath := filepath.ToSlash(file)

		// 写入文件到 VFS
		err := va.vfs.WriteFile(normalizedPath, []byte(mockFile.Content), os.FileMode(0644))
		if err != nil {
			// 如果写入失败，尝试创建目录
			dir := filepath.Dir(normalizedPath)
			va.ensureDir(dir)
			va.vfs.WriteFile(normalizedPath, []byte(mockFile.Content), os.FileMode(0644))
		}
	}

	return va.vfs
}

// ToRealFS 将 Mock 项目写入真实文件系统
func (va *VFSAdapter) ToRealFS(basePath string) error {
	// 将 Mock 文件系统中的所有文件写入真实文件系统
	for _, file := range va.project.FileSystem.GetAllFiles() {
		mockFile, exists := va.project.FileSystem.GetFile(file)
		if !exists {
			continue
		}

		// 标准化路径分隔符
		normalizedPath := filepath.ToSlash(file)

		// 构建真实路径
		realPath := filepath.Join(basePath, normalizedPath)

		// 确保目录存在
		dir := filepath.Dir(realPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}

		// 写入文件
		if err := os.WriteFile(realPath, []byte(mockFile.Content), 0644); err != nil {
			return err
		}
	}

	return nil
}

// ensureDir 确保目录存在
func (va *VFSAdapter) ensureDir(dir string) {
	// 递归创建目录
	parent := filepath.Dir(dir)
	if parent != "." && parent != "/" {
		va.ensureDir(parent)
	}

	// 创建当前目录（如果不存在）
	va.vfs.MkdirAll(dir, 0755)
}

// MockProject 的扩展方法
func (mp *MockProject) ToVFS() vfs.VFS {
	adapter := NewVFSAdapter(mp)
	return adapter.ToVFS()
}

func (mp *MockProject) ToRealFS(basePath string) error {
	adapter := NewVFSAdapter(mp)
	return adapter.ToRealFS(basePath)
}

// 创建测试用的 VFS 实例
func CreateTestVFS(project *MockProject) vfs.VFS {
	return project.ToVFS()
}

// 从预设创建测试 VFS
func CreateSimpleMathVFS() vfs.VFS {
	presets := NewMockPresets()
	project := presets.SimpleMathProject()
	return project.ToVFS()
}

func CreateHelloWorldVFS() vfs.VFS {
	presets := NewMockPresets()
	project := presets.HelloWorldProject()
	return project.ToVFS()
}

func CreateComplexDepsVFS() vfs.VFS {
	presets := NewMockPresets()
	project := presets.ComplexDependencyProject()
	return project.ToVFS()
}
