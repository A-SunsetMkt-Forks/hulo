package cachedvfs

import (
	"github.com/hulo-lang/hulo/syntax/hulo/ast"
	"github.com/hulo-lang/hulo/syntax/hulo/parser"
	"github.com/hulo-lang/hulo/internal/vfs"
)

func FS() vfs.FileSystem {
	return nil
}

type CachedVFS struct {
	dirs map[string]*ast.File
	opt  *parser.ParseOptions
}

func (fs *CachedVFS) Import(filename string) error {
	_, ok := fs.dirs[filename]
	if ok {
		return nil
	}
	file, err := parser.ParseSourceFile(filename, *fs.opt)
	if err != nil {
		return err
	}
	fs.dirs[filename] = file
	return nil
}
