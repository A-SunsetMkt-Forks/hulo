// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package osfs

import (
	"io"
	"os"
	"path/filepath"

	"github.com/hulo-lang/hulo/internal/vfs"
	"github.com/spf13/afero"
)

var _ vfs.VFS = (*OSVFS)(nil)

type OSVFS struct {
	afero.Fs
}

func New() *OSVFS {
	return &OSVFS{
		Fs: afero.NewOsFs(),
	}
}

func (o *OSVFS) ReadFile(path string) ([]byte, error) {
	return afero.ReadFile(o.Fs, path)
}

func (o *OSVFS) ReadDir(path string) ([]os.FileInfo, error) {
	return afero.ReadDir(o.Fs, path)
}

func (o *OSVFS) ReadAll(r io.Reader) ([]byte, error) {
	return afero.ReadAll(r)
}

func (o *OSVFS) WriteFile(path string, content []byte, perm os.FileMode) error {
	return afero.WriteFile(o.Fs, path, content, perm)
}

func (o *OSVFS) Exists(path string) bool {
	_, err := o.Fs.Stat(path)
	return err == nil
}

func (o *OSVFS) Walk(root string, walkFn filepath.WalkFunc) error {
	return afero.Walk(o.Fs, root, walkFn)
}

func (o *OSVFS) Glob(pattern string) ([]string, error) {
	return afero.Glob(o.Fs, pattern)
}
