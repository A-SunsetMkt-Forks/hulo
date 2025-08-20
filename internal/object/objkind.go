// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package object

//go:generate stringer -type=ObjKind -linecomment
type ObjKind int

// Equal reports whether o and x represent the same kind.
func (o ObjKind) Equal(x ObjKind) bool {
	return o == x
}

const (
	// Basic types

	O_NUM    ObjKind = iota // num
	O_STR                   // str
	O_BOOL                  // bool
	O_SYMBOL                // symbol

	// Composite types

	O_ARR   // arr
	O_TUPLE // tuple
	O_SET   // set
	O_OBJ   // object
	O_MAP   // map
	O_FUNC  // func

	// Special types

	O_NULL  // null
	O_ANY   // any
	O_VOID  // void
	O_NEVER // never
	O_ERROR // error

	// Traits and classes

	O_TRAIT // trait
	O_CLASS // class
	O_ENUM  // enum

	O_BUILTIN
	O_LITERAL
	O_RET
	O_QUOTE

	// A | B
	O_UNION // union
	// A & B
	O_INTERSECTION // intersection
	// T?
	O_NULLABLE // nullable

	// Decorators
	O_DECORATOR // decorator

	// Generic types
	O_TYPE_PARAM // type_param
)
