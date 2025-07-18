// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package exec

import (
	"fmt"

	"github.com/hulo-lang/hulo/internal/config"
	"github.com/hulo-lang/hulo/internal/transpiler/bash"
	"github.com/hulo-lang/hulo/internal/vfs"
	"github.com/hulo-lang/hulo/internal/vfs/memvfs"
)

var _ Executor = (*DebugExecutor)(nil)

type DebugExecutor struct {
	Transpiler config.DebugSettingsTranspiler
	Parser     config.DebugSettingsParser
	Analyzer   config.DebugSettingsAnalyzer
	Target     string
	fs         vfs.VFS
}

func NewDebugExecutor(fs vfs.VFS) *DebugExecutor {
	return &DebugExecutor{
		fs: memvfs.New(),
	}
}

func (e *DebugExecutor) CanHandle(cmd string) bool {
	return true
}

func (e *DebugExecutor) Execute(cmd string) error {
	e.fs.WriteFile("main.hl", []byte(cmd), 0644)
	res, err := transpiler.Transpile(&config.Huloc{Main: "main.hl", HuloPath: "."}, e.fs, ".")
	if err != nil {
		return err
	}
	for _, code := range res {
		fmt.Println(code)
	}
	return nil
}
