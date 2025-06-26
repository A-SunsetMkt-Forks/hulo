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
	afero.Fs

	ReadFile(path string) ([]byte, error)
	ReadDir(path string) ([]os.FileInfo, error)
	ReadAll(r io.Reader) ([]byte, error)

	WriteFile(path string, content []byte, perm os.FileMode) error
	Exists(path string) bool

	Walk(root string, walkFn filepath.WalkFunc) error

	Glob(pattern string) ([]string, error)
}
