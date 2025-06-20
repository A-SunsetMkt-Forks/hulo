package lib

import (
	"github.com/hulo-lang/hulo/internal/object"
)

var Math = map[string]object.Value{
	// 存储 builtin 实现的
	// "Random"
}

func init() {
	// 拉取 HULOPATH 路径中的实现 存储到 Math 里面

	// 引入后就有 ast 了
	// 然后 eval 一下 ast 加入到 Math 中

}
