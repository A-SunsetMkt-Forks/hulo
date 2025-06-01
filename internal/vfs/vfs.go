package vfs

import (
	"github.com/hulo-lang/hulo/internal/object"
)

type FileSystem interface {
	ImportAndGet(pkg string) (object.Value, bool)
	Get(name string) (object.Value, bool)
	Put()
	Import(filename string) (any, error)

	Exist(filename string) bool
}
