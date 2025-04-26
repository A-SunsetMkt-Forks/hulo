package object

import "github.com/hulo-lang/hulo/internal/ast"

type ObjKind int

func (o ObjKind) Equal(x ObjKind) bool {
	return o == x
}

const (
	O_NUM ObjKind = iota
	O_STR
	O_BOOL
	O_FUNC
	O_NULL
	O_BUILTIN
	O_LITERAL
	O_ARR
	O_MAP
	O_RET
	O_QUOTE
	O_TRAIT
	O_ENUM
	O_CLASS
	O_OBJ
)

type Type interface {
	Decl() ast.Decl

	Modifier() ast.Modifier

	PkgPath() string

	Name() string

	Text() string

	New(values ...Value) Value

	// Implements reports whether the type implements the interface type u.
	Implements(u Type) bool

	Kind() ObjKind

	NumMethod() int

	Method(i int) Method

	MethodByName(name string) Method

	NumField() int

	Field(i int) Type

	FieldByName(name string) Type
}

type Value interface {
	// Name returns identifier
	Name() string

	Type() Type

	Text() string
}

type Method interface {
	Type

	Call(values ...Value) []Value

	// NumIn returns a function type's input parameter count.
	// It panics if the type's Kind is not Func.
	NumIn() int

	// NumOut returns a function type's output parameter count.
	// It panics if the type's Kind is not Func.
	NumOut() int

	// In returns the type of a function type's i'th input parameter.
	// It panics if the type's Kind is not Func.
	// It panics if i is not in the range [0, NumIn()).
	In(i int) Type

	// Out returns the type of a function type's i'th output parameter.
	// It panics if the type's Kind is not Func.
	// It panics if i is not in the range [0, NumOut()).
	Out(i int) Type

	IsCallable() bool
}
