// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package mock

import (
	"testing"
)

func TestNewMockFileSystem(t *testing.T) {
	fs := NewMockFileSystem()
	if fs == nil {
		t.Fatal("NewMockFileSystem returned nil")
	}
	if fs.files == nil {
		t.Error("files map is nil")
	}
	if fs.dirs == nil {
		t.Error("dirs map is nil")
	}
}

func TestMockFileSystem_AddFile(t *testing.T) {
	fs := NewMockFileSystem()

	// 添加文件
	fs.AddFile("/test/file.txt", "hello world")

	// 检查文件是否存在
	if !fs.FileExists("/test/file.txt") {
		t.Error("File should exist after AddFile")
	}

	// 检查文件内容
	file, exists := fs.GetFile("/test/file.txt")
	if !exists {
		t.Fatal("GetFile should return file")
	}
	if file.Content != "hello world" {
		t.Errorf("Expected content 'hello world', got '%s'", file.Content)
	}
	if file.Path != "/test/file.txt" {
		t.Errorf("Expected path '/test/file.txt', got '%s'", file.Path)
	}
}

func TestMockFileSystem_AddDir(t *testing.T) {
	fs := NewMockFileSystem()

	// 添加目录
	fs.AddDir("/test/dir")

	// 检查目录是否存在
	if !fs.DirExists("/test/dir") {
		t.Error("Directory should exist after AddDir")
	}

	// 检查目录
	dir, exists := fs.GetDir("/test/dir")
	if !exists {
		t.Fatal("GetDir should return directory")
	}
	if dir.Path != "/test/dir" {
		t.Errorf("Expected path '/test/dir', got '%s'", dir.Path)
	}
}

func TestMockFileSystem_EnsureDir(t *testing.T) {
	fs := NewMockFileSystem()

	// 添加文件，应该自动创建父目录
	fs.AddFile("/parent/child/file.txt", "content")

	// 检查父目录是否被创建
	if !fs.DirExists("/parent") {
		t.Error("Parent directory should be created automatically")
	}
	if !fs.DirExists("/parent/child") {
		t.Error("Child directory should be created automatically")
	}

	// 检查文件是否存在
	if !fs.FileExists("/parent/child/file.txt") {
		t.Error("File should exist")
	}
}

func TestMockFileSystem_FileExists(t *testing.T) {
	fs := NewMockFileSystem()

	// 添加文件
	fs.AddFile("/test.txt", "content")

	// 测试存在
	if !fs.FileExists("/test.txt") {
		t.Error("File should exist")
	}

	// 测试不存在
	if fs.FileExists("/nonexistent.txt") {
		t.Error("File should not exist")
	}
}

func TestMockFileSystem_DirExists(t *testing.T) {
	fs := NewMockFileSystem()

	// 添加目录
	fs.AddDir("/test")

	// 测试存在
	if !fs.DirExists("/test") {
		t.Error("Directory should exist")
	}

	// 测试不存在
	if fs.DirExists("/nonexistent") {
		t.Error("Directory should not exist")
	}
}

func TestMockFileSystem_GetAllFiles(t *testing.T) {
	fs := NewMockFileSystem()

	// 添加多个文件
	fs.AddFile("/file1.txt", "content1")
	fs.AddFile("/file2.txt", "content2")
	fs.AddFile("/dir/file3.txt", "content3")

	// 获取所有文件
	files := fs.GetAllFiles()

	// 检查文件数量
	expectedCount := 3
	if len(files) != expectedCount {
		t.Errorf("Expected %d files, got %d", expectedCount, len(files))
	}

	// 检查文件路径
	expectedFiles := map[string]bool{
		"/file1.txt":     true,
		"/file2.txt":     true,
		"/dir/file3.txt": true,
	}

	for _, file := range files {
		if !expectedFiles[file] {
			t.Errorf("Unexpected file: %s", file)
		}
	}
}

func TestMockFileSystem_GetAllDirs(t *testing.T) {
	fs := NewMockFileSystem()

	// 添加多个目录
	fs.AddDir("/dir1")
	fs.AddDir("/dir2")
	fs.AddDir("/dir1/subdir")

	// 添加文件（会创建父目录）
	fs.AddFile("/dir3/file.txt", "content")

	// 获取所有目录
	dirs := fs.GetAllDirs()

	// 检查目录数量（至少应该有4个：dir1, dir2, dir1/subdir, dir3）
	if len(dirs) < 4 {
		t.Errorf("Expected at least 4 directories, got %d", len(dirs))
	}

	// 检查关键目录是否存在
	expectedDirs := map[string]bool{
		"/dir1":        true,
		"/dir2":        true,
		"/dir1/subdir": true,
		"/dir3":        true,
	}

	for _, dir := range dirs {
		if expectedDirs[dir] {
			delete(expectedDirs, dir)
		}
	}

	// 检查是否所有期望的目录都找到了
	for dir := range expectedDirs {
		t.Errorf("Expected directory not found: %s", dir)
	}
}
