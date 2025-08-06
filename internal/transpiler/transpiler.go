package transpiler

import "github.com/hulo-lang/hulo/syntax/hulo/ast"

type Transpiler[T any] interface {
	// 转译主文件，会返回所有关联的文件，然后输出未解析的符号信息让linker去解析
	Transpile(mainFile string) (files map[string]T, unresolvedSymbols []T, err error)

	Convert(node ast.Node) (T, error)

	// 获取目标语言类型
	GetTargetType() T

	// 获取目标语言扩展文件名
	GetTargetExt() string

	// 获取目标语言名称
	GetTargetName() string
}
