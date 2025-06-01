package comptime

import (
	"github.com/hulo-lang/hulo/internal/object"
	"github.com/hulo-lang/hulo/internal/vfs"
)

type Context struct {
	global map[string]object.Value
	local  map[string]object.Value
	parent *Context

	mem vfs.FileSystem
	os  vfs.FileSystem
}

func DefaultContext() *Context {
	return &Context{}
}

func (c *Context) Get(name string) (object.Value, bool) {
	obj, ok := c.local[name]
	if !ok {
		obj, ok = c.global[name]
		if !ok && c.parent != nil {
			obj, ok = c.parent.Get(name)
		}
	}
	return obj, ok
}

func (c *Context) Fork() *Context {
	return c
}
