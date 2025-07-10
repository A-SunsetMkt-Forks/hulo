// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package vfs

import (
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/afero"
)

// VFS is a virtual file system interface.
// It extends the afero.Fs interface to provide additional file system operations.
type VFS interface {
	// afero.Fs is the base interface that provides basic file system operations.
	afero.Fs

	// ReadFile reads the file at the given path and returns its contents.
	ReadFile(path string) ([]byte, error)

	// ReadDir reads the directory at the given path and returns its contents.
	ReadDir(path string) ([]os.FileInfo, error)

	// ReadAll reads all the data from the given reader and returns its contents.
	ReadAll(r io.Reader) ([]byte, error)

	// WriteFile writes the given content to the file at the given path.
	WriteFile(path string, content []byte, perm os.FileMode) error

	// Exists checks if the file or directory exists at the given path.
	Exists(path string) bool

	// Walk traverses the file system starting at the given root directory,
	// calling walkFn for each file or directory in the tree, including root.
	Walk(root string, walkFn filepath.WalkFunc) error

	// Glob returns all files matching the given pattern.
	Glob(pattern string) ([]string, error)
}
