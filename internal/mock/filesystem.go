// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package mock

import (
	"path/filepath"
	"time"
)

type MockFileSystem struct {
	files map[string]*MockFile
	dirs  map[string]*MockDir
}

type MockFile struct {
	Path    string
	Content string
	ModTime time.Time
	Size    int64
}

type MockDir struct {
	Path     string
	Children []string // 子文件/目录路径
}

// 创建 Mock 文件系统
func NewMockFileSystem() *MockFileSystem {
	return &MockFileSystem{
		files: make(map[string]*MockFile),
		dirs:  make(map[string]*MockDir),
	}
}

// 添加文件
func (mfs *MockFileSystem) AddFile(path, content string) {
	// 标准化路径分隔符
	normalizedPath := filepath.ToSlash(path)
	mfs.files[normalizedPath] = &MockFile{
		Path:    normalizedPath,
		Content: content,
		ModTime: time.Now(),
		Size:    int64(len(content)),
	}

	// 创建目录结构
	mfs.ensureDir(normalizedPath)
}

// 添加目录
func (mfs *MockFileSystem) AddDir(path string) {
	// 标准化路径分隔符
	normalizedPath := filepath.ToSlash(path)
	mfs.dirs[normalizedPath] = &MockDir{
		Path:     normalizedPath,
		Children: make([]string, 0),
	}
}

// 确保目录存在
func (mfs *MockFileSystem) ensureDir(filePath string) {
	dir := filepath.Dir(filePath)
	if dir == "." || dir == "/" || dir == filePath || dir == "" {
		return
	}

	// 递归创建父目录
	mfs.ensureDir(dir)

	// 添加当前目录
	if _, exists := mfs.dirs[dir]; !exists {
		mfs.AddDir(dir)
	}
}

// 获取文件内容
func (mfs *MockFileSystem) GetFile(path string) (*MockFile, bool) {
	file, exists := mfs.files[path]
	return file, exists
}

// 获取目录
func (mfs *MockFileSystem) GetDir(path string) (*MockDir, bool) {
	dir, exists := mfs.dirs[path]
	return dir, exists
}

// 检查文件是否存在
func (mfs *MockFileSystem) FileExists(path string) bool {
	_, exists := mfs.files[path]
	return exists
}

// 检查目录是否存在
func (mfs *MockFileSystem) DirExists(path string) bool {
	_, exists := mfs.dirs[path]
	return exists
}

// 列出目录内容
func (mfs *MockFileSystem) ListDir(path string) []string {
	if dir, exists := mfs.dirs[path]; exists {
		return dir.Children
	}
	return []string{}
}

// 获取所有文件路径
func (mfs *MockFileSystem) GetAllFiles() []string {
	files := make([]string, 0, len(mfs.files))
	for path := range mfs.files {
		files = append(files, path)
	}
	return files
}

// 获取所有目录路径
func (mfs *MockFileSystem) GetAllDirs() []string {
	dirs := make([]string, 0, len(mfs.dirs))
	for path := range mfs.dirs {
		dirs = append(dirs, path)
	}
	return dirs
}
