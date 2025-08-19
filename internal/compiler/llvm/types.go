// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package llvm

import "github.com/llir/llvm/ir/types"

var typeMapping = map[string]types.Type{
	"num":  types.I32,
	"str":  types.NewPointer(types.I8),
	"bool": types.I1,
	"void": types.Void,
}
