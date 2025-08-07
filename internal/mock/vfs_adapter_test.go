// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package mock

import (
	"fmt"
	"strings"
	"testing"
)

func TestVFSAdapter_ToVFS(t *testing.T) {
	// 创建 Mock 项目
	project := NewMockBuilder().
		WithName("test-project").
		WithModule("math", "1.0.0").
		WithFile("index.hl", `pub fn add(a: num, b: num) => $a + $b`).
		WithMain("./index.hl").
		End().
		Build()

	// 添加工作目录文件
	project.AddWorkDirFile("src/main.hl", `fn main() => println("hello")`)

	// 调试：打印所有文件路径
	fmt.Println("=== Debug: All files in MockFileSystem ===")
	for _, file := range project.FileSystem.GetAllFiles() {
		fmt.Printf("File: %s\n", file)
	}

	// 调试：打印模块信息
	fmt.Println("=== Debug: Module info ===")
	for name, module := range project.Modules {
		fmt.Printf("Module: %s, Version: %s\n", name, module.Version)
		fmt.Printf("Files: %v\n", module.Files)
	}

	// 转换为 VFS
	vfs := project.ToVFS()

	// 检查文件是否存在
	if !vfs.Exists("./testdata/src/main.hl") {
		t.Error("Main file should exist in VFS")
	}

	if !vfs.Exists("./testdata/hulo_modules/math/1.0.0/index.hl") {
		t.Error("Module file should exist in VFS")
	}

	// 检查文件内容
	content, err := vfs.ReadFile("./testdata/src/main.hl")
	if err != nil {
		t.Fatalf("Failed to read main file: %v", err)
	}

	if !strings.Contains(string(content), "fn main()") {
		t.Error("Main file should contain main function")
	}
}

func TestCreateTestVFS(t *testing.T) {
	// 测试快速创建 VFS 的辅助函数
	vfs := CreateSimpleMathVFS()

	if vfs == nil {
		t.Fatal("CreateSimpleMathVFS should return VFS")
	}

	// 检查文件是否存在
	modulePath := "./testdata/hulo_modules/math/1.0.0/index.hl"
	if !vfs.Exists(modulePath) {
		t.Error("Math module should exist")
	}
}
