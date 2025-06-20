package vfs

import (
	"io"
	"io/fs"
)

// FileSystem represents a virtual file system interface
type FileSystem interface {
	// Basic file operations
	Open(name string) (File, error)
	ReadFile(name string) ([]byte, error)
	WriteFile(name string, data []byte) error
	ReadDir(name string) ([]DirEntry, error)

	// File existence and metadata
	Exist(name string) bool
	Stat(name string) (FileInfo, error)

	// Directory operations
	Mkdir(name string, perm fs.FileMode) error
	MkdirAll(name string, perm fs.FileMode) error
	Remove(name string) error
	RemoveAll(name string) error

	// Path operations
	Join(elem ...string) string
	Clean(path string) string
	Abs(path string) (string, error)
	Rel(basepath, targpath string) (string, error)

	// File system specific operations
	Scheme() string
	Close() error
}

// File represents a file in the virtual file system
type File interface {
	fs.File
	io.ReadWriteCloser
	io.Seeker
}

// DirEntry represents a directory entry
type DirEntry interface {
	fs.DirEntry
}

// FileInfo represents file metadata
type FileInfo interface {
	fs.FileInfo
}

// VFSConfig holds configuration for the file system
type VFSConfig struct {
	// Root directory for the file system
	Root string

	// Search paths for modules
	SearchPaths []string

	// Cache settings
	EnableCache bool
	CacheSize   int

	// Environment variables
	Env map[string]string

	// File system specific options
	Options map[string]interface{}
}

// NewVFSConfig creates a new VFS configuration with defaults
func NewVFSConfig() VFSConfig {
	return VFSConfig{
		SearchPaths: []string{},
		EnableCache: true,
		CacheSize:   1000,
		Env:         make(map[string]string),
		Options:     make(map[string]interface{}),
	}
}

// NewFileSystem creates a file system from a URI
func NewFileSystem(uri string, config VFSConfig) (FileSystem, error) {
	scheme := parseScheme(uri)

	switch scheme {
	case "file":
		return NewLocalFileSystem(uri, config)
	case "mem":
		return NewMemoryFileSystem(uri, config)
	case "http", "https":
		return NewHTTPFileSystem(uri, config)
	case "s3":
		return NewS3FileSystem(uri, config)
	case "zip":
		return NewZipFileSystem(uri, config)
	default:
		return NewLocalFileSystem(uri, config)
	}
}

// parseScheme extracts the scheme from a URI
func parseScheme(uri string) string {
	for i, char := range uri {
		if char == ':' {
			return uri[:i]
		}
		if char == '/' {
			break
		}
	}
	return "file"
}
