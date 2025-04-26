package comptime

type Context struct{}

func DefaultContext() *Context {
	return &Context{}
}
