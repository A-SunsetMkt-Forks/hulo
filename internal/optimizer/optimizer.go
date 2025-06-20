package optimizer

import "github.com/hulo-lang/hulo/syntax/hulo/ast"

type Optimizer struct {
	scissors []Scissor // 各种剪枝工具
}

type Scissor interface {
	Name() string
	ShouldCut(node ast.Node) bool
	Cut(node ast.Node) ast.Node
}

// 常量折叠剪刀
type ConstantFoldingScissor struct{}

func (s *ConstantFoldingScissor) ShouldCut(node ast.Node) bool {
	// 判断是否是可以折叠的常量表达式
	return isConstantExpression(node)
}

func (s *ConstantFoldingScissor) Cut(node ast.Node) ast.Node {
	// 把 1+2*3 剪成 7
	return calculateConstant(node)
}

// 死代码消除剪刀
type DeadCodeScissor struct{}

func (s *DeadCodeScissor) ShouldCut(node ast.Node) bool {
	// 判断是否是死代码
	return isDeadCode(node)
}

func (s *DeadCodeScissor) Cut(node ast.Node) ast.Node {
	// 直接删除这个节点
	return nil
}

func isConstantExpression(node ast.Node) bool {
	return false
}

func calculateConstant(node ast.Node) ast.Node {
	return nil
}

func isDeadCode(node ast.Node) bool {
	return false
}
