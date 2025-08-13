package transpiler

import (
	"github.com/hulo-lang/hulo/internal/linker"
	"github.com/hulo-lang/hulo/syntax/hulo/ast"
)

type Transpiler[T any] interface {
	Convert(node ast.Node) T

	// TODO 实现各种Convert方法

	// 获取目标语言扩展文件名
	GetTargetExt() string

	// 获取目标语言名称
	GetTargetName() string

	UnresolvedSymbols() map[string][]linker.UnkownSymbol
}
