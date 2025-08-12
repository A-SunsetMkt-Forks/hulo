// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package optimizer

import (
	"strconv"

	"github.com/hulo-lang/hulo/syntax/hulo/ast"
	"github.com/hulo-lang/hulo/syntax/hulo/token"
)

// TODO import 的时候未使用的变量优化

type Optimizer struct {
	scissors []Scissor // 各种剪枝工具
}

// NewOptimizer 创建新的优化器实例
func NewOptimizer() *Optimizer {
	return &Optimizer{
		scissors: []Scissor{
			&ConstantFoldingScissor{},
			&DeadCodeScissor{},
			&UnusedVariableScissor{},
		},
	}
}

type Scissor interface {
	Name() string
	ShouldCut(node ast.Node) bool
	Cut(node ast.Node) ast.Node
}

// Optimize 优化AST
func (o *Optimizer) Optimize(node ast.Node) ast.Node {
	if node == nil {
		return nil
	}

	// 先对子节点进行优化
	o.optimizeChildren(node)

	// 然后检查当前节点是否需要剪枝
	for _, scissor := range o.scissors {
		if scissor.ShouldCut(node) {
			return scissor.Cut(node)
		}
	}

	return node
}

// optimizeChildren 递归优化所有子节点
func (o *Optimizer) optimizeChildren(node ast.Node) {
	switch n := node.(type) {
	case *ast.BlockStmt:
		o.optimizeBlockStmt(n)
	case *ast.IfStmt:
		o.optimizeIfStmt(n)
	case *ast.WhileStmt:
		o.optimizeWhileStmt(n)
	case *ast.ForStmt:
		o.optimizeForStmt(n)
	case *ast.ForInStmt:
		o.optimizeForInStmt(n)
	case *ast.ForeachStmt:
		o.optimizeForeachStmt(n)
	case *ast.DoWhileStmt:
		o.optimizeDoWhileStmt(n)
	case *ast.MatchStmt:
		o.optimizeMatchStmt(n)
	case *ast.TryStmt:
		o.optimizeTryStmt(n)
	case *ast.FuncDecl:
		o.optimizeFuncDecl(n)
	case *ast.ClassDecl:
		o.optimizeClassDecl(n)
	case *ast.EnumDecl:
		o.optimizeEnumDecl(n)
	case *ast.BinaryExpr:
		o.optimizeBinaryExpr(n)
	case *ast.CallExpr:
		o.optimizeCallExpr(n)
	case *ast.AssignStmt:
		o.optimizeAssignStmt(n)
	}
}

// 各种语句的优化方法
func (o *Optimizer) optimizeBlockStmt(stmt *ast.BlockStmt) {
	var optimizeds []ast.Stmt
	for _, s := range stmt.List {
		if optimized := o.Optimize(s); optimized != nil {
			optimizeds = append(optimizeds, optimized.(ast.Stmt))
		}
	}
	stmt.List = optimizeds
}

func (o *Optimizer) optimizeIfStmt(stmt *ast.IfStmt) {
	stmt.Cond = o.Optimize(stmt.Cond).(ast.Expr)
	stmt.Body = o.Optimize(stmt.Body).(*ast.BlockStmt)
	if stmt.Else != nil {
		stmt.Else = o.Optimize(stmt.Else).(ast.Stmt)
	}
}

func (o *Optimizer) optimizeWhileStmt(stmt *ast.WhileStmt) {
	stmt.Cond = o.Optimize(stmt.Cond).(ast.Expr)
	stmt.Body = o.Optimize(stmt.Body).(*ast.BlockStmt)
}

func (o *Optimizer) optimizeForStmt(stmt *ast.ForStmt) {
	if stmt.Init != nil {
		stmt.Init = o.Optimize(stmt.Init).(ast.Stmt)
	}
	if stmt.Cond != nil {
		stmt.Cond = o.Optimize(stmt.Cond).(ast.Expr)
	}
	if stmt.Post != nil {
		stmt.Post = o.Optimize(stmt.Post).(ast.Expr)
	}
	stmt.Body = o.Optimize(stmt.Body).(*ast.BlockStmt)
}

func (o *Optimizer) optimizeForInStmt(stmt *ast.ForInStmt) {
	stmt.RangeExpr.Start = o.Optimize(stmt.RangeExpr.Start).(ast.Expr)
	stmt.RangeExpr.End_ = o.Optimize(stmt.RangeExpr.End_).(ast.Expr)
	if stmt.RangeExpr.Step != nil {
		stmt.RangeExpr.Step = o.Optimize(stmt.RangeExpr.Step).(ast.Expr)
	}
	stmt.Body = o.Optimize(stmt.Body).(*ast.BlockStmt)
}

func (o *Optimizer) optimizeForeachStmt(stmt *ast.ForeachStmt) {
	stmt.Index = o.Optimize(stmt.Index).(ast.Expr)
	stmt.Value = o.Optimize(stmt.Value).(ast.Expr)
	stmt.Var = o.Optimize(stmt.Var).(ast.Expr)
	stmt.Body = o.Optimize(stmt.Body).(*ast.BlockStmt)
}

func (o *Optimizer) optimizeDoWhileStmt(stmt *ast.DoWhileStmt) {
	stmt.Body = o.Optimize(stmt.Body).(*ast.BlockStmt)
	stmt.Cond = o.Optimize(stmt.Cond).(ast.Expr)
}

func (o *Optimizer) optimizeMatchStmt(stmt *ast.MatchStmt) {
	stmt.Expr = o.Optimize(stmt.Expr).(ast.Expr)
	for _, c := range stmt.Cases {
		c.Cond = o.Optimize(c.Cond).(ast.Expr)
		c.Body = o.Optimize(c.Body).(*ast.BlockStmt)
	}
	if stmt.Default != nil {
		stmt.Default.Body = o.Optimize(stmt.Default.Body).(*ast.BlockStmt)
	}
}

func (o *Optimizer) optimizeTryStmt(stmt *ast.TryStmt) {
	stmt.Body = o.Optimize(stmt.Body).(*ast.BlockStmt)
	for _, c := range stmt.Catches {
		c.Cond = o.Optimize(c.Cond).(ast.Expr)
		c.Body = o.Optimize(c.Body).(*ast.BlockStmt)
	}
	if stmt.Finally != nil {
		stmt.Finally.Body = o.Optimize(stmt.Finally.Body).(*ast.BlockStmt)
	}
}

func (o *Optimizer) optimizeFuncDecl(decl *ast.FuncDecl) {
	for i, param := range decl.Recv {
		decl.Recv[i] = o.Optimize(param).(ast.Expr)
	}
	if decl.Type != nil {
		decl.Type = o.Optimize(decl.Type).(ast.Expr)
	}
	if decl.Body != nil {
		decl.Body = o.Optimize(decl.Body).(*ast.BlockStmt)
	}
}

func (o *Optimizer) optimizeClassDecl(decl *ast.ClassDecl) {
	if decl.Fields != nil {
		for _, field := range decl.Fields.List {
			if field.Type != nil {
				field.Type = o.Optimize(field.Type).(ast.Expr)
			}
			if field.Value != nil {
				field.Value = o.Optimize(field.Value).(ast.Expr)
			}
		}
	}
	for _, ctor := range decl.Ctors {
		for i, param := range ctor.Recv {
			ctor.Recv[i] = o.Optimize(param).(ast.Expr)
		}
		if ctor.Body != nil {
			ctor.Body = o.Optimize(ctor.Body).(*ast.BlockStmt)
		}
	}
	for _, method := range decl.Methods {
		o.optimizeFuncDecl(method)
	}
}

func (o *Optimizer) optimizeEnumDecl(decl *ast.EnumDecl) {
	// 枚举声明的优化逻辑
}

func (o *Optimizer) optimizeBinaryExpr(expr *ast.BinaryExpr) {
	expr.X = o.Optimize(expr.X).(ast.Expr)
	expr.Y = o.Optimize(expr.Y).(ast.Expr)
}

func (o *Optimizer) optimizeCallExpr(expr *ast.CallExpr) {
	expr.Fun = o.Optimize(expr.Fun).(ast.Expr)
	for i, arg := range expr.Recv {
		expr.Recv[i] = o.Optimize(arg).(ast.Expr)
	}
}

func (o *Optimizer) optimizeAssignStmt(stmt *ast.AssignStmt) {
	stmt.Lhs = o.Optimize(stmt.Lhs).(ast.Expr)
	if stmt.Type != nil {
		stmt.Type = o.Optimize(stmt.Type).(ast.Expr)
	}
	stmt.Rhs = o.Optimize(stmt.Rhs).(ast.Expr)
}

// 常量折叠剪刀
type ConstantFoldingScissor struct{}

func (s *ConstantFoldingScissor) Name() string {
	return "ConstantFolding"
}

func (s *ConstantFoldingScissor) ShouldCut(node ast.Node) bool {
	return isConstantExpression(node)
}

func (s *ConstantFoldingScissor) Cut(node ast.Node) ast.Node {
	return calculateConstant(node)
}

// 死代码消除剪刀
type DeadCodeScissor struct{}

func (s *DeadCodeScissor) Name() string {
	return "DeadCode"
}

func (s *DeadCodeScissor) ShouldCut(node ast.Node) bool {
	return isDeadCode(node)
}

func (s *DeadCodeScissor) Cut(node ast.Node) ast.Node {
	return nil // 直接删除死代码
}

// 未使用变量消除剪刀
type UnusedVariableScissor struct{}

func (s *UnusedVariableScissor) Name() string {
	return "UnusedVariable"
}

func (s *UnusedVariableScissor) ShouldCut(node ast.Node) bool {
	return isUnusedVariable(node)
}

func (s *UnusedVariableScissor) Cut(node ast.Node) ast.Node {
	return nil // 直接删除未使用的变量
}

// 判断是否为常量表达式
func isConstantExpression(node ast.Node) bool {
	switch n := node.(type) {
	case *ast.BasicLit:
		return true
	case *ast.NumericLiteral, *ast.StringLiteral, *ast.TrueLiteral, *ast.FalseLiteral:
		return true
	case *ast.BinaryExpr:
		return isConstantExpression(n.X) && isConstantExpression(n.Y)
	case *ast.UnaryExpr:
		return isConstantExpression(n.X)
	}
	return false
}

// 计算常量表达式的值
func calculateConstant(node ast.Node) ast.Node {
	switch n := node.(type) {
	case *ast.BasicLit:
		return n
	case *ast.NumericLiteral:
		return n
	case *ast.StringLiteral:
		return n
	case *ast.TrueLiteral:
		return n
	case *ast.FalseLiteral:
		return n
	case *ast.BinaryExpr:
		return calculateBinaryConstant(n)
	case *ast.UnaryExpr:
		return calculateUnaryConstant(n)
	}
	return node
}

// 计算二元运算的常量值
func calculateBinaryConstant(expr *ast.BinaryExpr) ast.Node {
	left := calculateConstant(expr.X)
	right := calculateConstant(expr.Y)

	// 获取数值
	leftVal, leftOk := getNumericValue(left)
	rightVal, rightOk := getNumericValue(right)

	if !leftOk || !rightOk {
		return expr // 无法计算，返回原表达式
	}

	var result float64
	switch expr.Op {
	case token.PLUS:
		result = leftVal + rightVal
	case token.MINUS:
		result = leftVal - rightVal
	case token.ASTERISK:
		result = leftVal * rightVal
	case token.SLASH:
		if rightVal == 0 {
			return expr // 除零错误，返回原表达式
		}
		result = leftVal / rightVal
	case token.MOD:
		if rightVal == 0 {
			return expr
		}
		result = float64(int(leftVal) % int(rightVal))
	default:
		return expr
	}

	// 创建新的数值字面量
	return &ast.NumericLiteral{
		ValuePos: expr.Pos(),
		Value:    strconv.FormatFloat(result, 'g', -1, 64),
	}
}

// 计算一元运算的常量值
func calculateUnaryConstant(expr *ast.UnaryExpr) ast.Node {
	operand := calculateConstant(expr.X)

	val, ok := getNumericValue(operand)
	if !ok {
		return expr
	}

	var result float64
	switch expr.Op {
	case token.PLUS:
		result = val
	case token.MINUS:
		result = -val
	default:
		return expr
	}

	return &ast.NumericLiteral{
		ValuePos: expr.Pos(),
		Value:    strconv.FormatFloat(result, 'g', -1, 64),
	}
}

// 从节点中提取数值
func getNumericValue(node ast.Node) (float64, bool) {
	switch n := node.(type) {
	case *ast.BasicLit:
		if n.Kind == token.NUM {
			val, err := strconv.ParseFloat(n.Value, 64)
			return val, err == nil
		}
	case *ast.NumericLiteral:
		val, err := strconv.ParseFloat(n.Value, 64)
		return val, err == nil
	}
	return 0, false
}

// 判断是否为死代码
func isDeadCode(node ast.Node) bool {
	switch n := node.(type) {
	case *ast.IfStmt:
		// 检查条件是否为常量false
		if isConstantFalse(n.Cond) {
			return true
		}
	case *ast.WhileStmt:
		// 检查条件是否为常量false
		if isConstantFalse(n.Cond) {
			return true
		}
	case *ast.ForStmt:
		// 检查条件是否为常量false
		if n.Cond != nil && isConstantFalse(n.Cond) {
			return true
		}
	case *ast.ExprStmt:
		// 检查是否为无副作用的表达式
		if hasNoSideEffects(n.X) {
			return true
		}
	}
	return false
}

// 判断表达式是否为常量false
func isConstantFalse(expr ast.Expr) bool {
	switch e := expr.(type) {
	case *ast.FalseLiteral:
		return true
	case *ast.BasicLit:
		return e.Kind == token.FALSE
	}
	return false
}

// 判断表达式是否有副作用
func hasNoSideEffects(expr ast.Expr) bool {
	switch e := expr.(type) {
	case *ast.Ident, *ast.BasicLit, *ast.NumericLiteral, *ast.StringLiteral:
		return true
	case *ast.BinaryExpr:
		return hasNoSideEffects(e.X) && hasNoSideEffects(e.Y)
	case *ast.UnaryExpr:
		return hasNoSideEffects(e.X)
	case *ast.CallExpr:
		// 函数调用通常有副作用，除非是纯函数
		return false
	}
	return false
}

// 判断是否为未使用的变量
func isUnusedVariable(node ast.Node) bool {
	switch n := node.(type) {
	case *ast.AssignStmt:
		// 检查变量是否被使用
		if ident, ok := n.Lhs.(*ast.Ident); ok {
			return !isVariableUsed(ident.Name, n)
		}
	}
	return false
}

// 检查变量是否被使用（简化版本）
func isVariableUsed(varName string, node ast.Node) bool {
	// 这里需要实现一个完整的变量使用分析
	// 为了简化，我们假设所有变量都被使用
	return true
}
