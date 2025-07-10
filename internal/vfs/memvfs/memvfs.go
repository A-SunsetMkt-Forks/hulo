// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package memvfs

import (
	"io"
	"os"
	"path/filepath"

	"github.com/hulo-lang/hulo/internal/vfs"
	"github.com/spf13/afero"
)

var _ vfs.VFS = (*MemVFS)(nil)

type MemVFS struct {
	afero.Fs
}

func New() *MemVFS {
	return &MemVFS{
		Fs: afero.NewMemMapFs(),
	}
}

func (m *MemVFS) ReadFile(path string) ([]byte, error) {
	return afero.ReadFile(m.Fs, path)
}

func (m *MemVFS) ReadDir(path string) ([]os.FileInfo, error) {
	return afero.ReadDir(m.Fs, path)
}

func (m *MemVFS) ReadAll(r io.Reader) ([]byte, error) {
	return afero.ReadAll(r)
}

func (m *MemVFS) WriteFile(path string, content []byte, perm os.FileMode) error {
	return afero.WriteFile(m.Fs, path, content, perm)
}

func (m *MemVFS) Exists(path string) bool {
	_, err := m.Fs.Stat(path)
	return err == nil
}

func (m *MemVFS) Walk(root string, walkFn filepath.WalkFunc) error {
	return afero.Walk(m.Fs, root, walkFn)
}

func (m *MemVFS) Glob(pattern string) ([]string, error) {
	return afero.Glob(m.Fs, pattern)
}
