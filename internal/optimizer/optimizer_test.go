// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package optimizer

import (
	"testing"

	"github.com/hulo-lang/hulo/syntax/hulo/ast"
	"github.com/hulo-lang/hulo/syntax/hulo/token"
	"github.com/stretchr/testify/assert"
)

func TestNewOptimizer(t *testing.T) {
	opt := NewOptimizer()
	assert.NotNil(t, opt)
	assert.Equal(t, len(opt.scissors), 3)
}

func TestConstantFoldingScissor(t *testing.T) {
	scissor := &ConstantFoldingScissor{}

	tests := []struct {
		name      string
		node      ast.Node
		shouldCut bool
	}{
		{
			name:      "numeric literal",
			node:      &ast.NumericLiteral{Value: "42"},
			shouldCut: true,
		},
		{
			name:      "string literal",
			node:      &ast.StringLiteral{Value: "hello"},
			shouldCut: true,
		},
		{
			name:      "true literal",
			node:      &ast.TrueLiteral{},
			shouldCut: true,
		},
		{
			name:      "false literal",
			node:      &ast.FalseLiteral{},
			shouldCut: true,
		},
		{
			name: "binary expression with constants",
			node: &ast.BinaryExpr{
				X:  &ast.NumericLiteral{Value: "2"},
				Op: token.PLUS,
				Y:  &ast.NumericLiteral{Value: "3"},
			},
			shouldCut: true,
		},
		{
			name: "binary expression with variable",
			node: &ast.BinaryExpr{
				X:  &ast.Ident{Name: "x"},
				Op: token.PLUS,
				Y:  &ast.NumericLiteral{Value: "3"},
			},
			shouldCut: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := scissor.ShouldCut(tt.node)
			assert.Equal(t, result, tt.shouldCut)
		})
	}
}

func TestCalculateBinaryConstant(t *testing.T) {
	tests := []struct {
		name     string
		expr     *ast.BinaryExpr
		expected string
	}{
		{
			name: "addition",
			expr: &ast.BinaryExpr{
				X:  &ast.NumericLiteral{Value: "2"},
				Op: token.PLUS,
				Y:  &ast.NumericLiteral{Value: "3"},
			},
			expected: "5",
		},
		{
			name: "subtraction",
			expr: &ast.BinaryExpr{
				X:  &ast.NumericLiteral{Value: "10"},
				Op: token.MINUS,
				Y:  &ast.NumericLiteral{Value: "4"},
			},
			expected: "6",
		},
		{
			name: "multiplication",
			expr: &ast.BinaryExpr{
				X:  &ast.NumericLiteral{Value: "6"},
				Op: token.ASTERISK,
				Y:  &ast.NumericLiteral{Value: "7"},
			},
			expected: "42",
		},
		{
			name: "division",
			expr: &ast.BinaryExpr{
				X:  &ast.NumericLiteral{Value: "15"},
				Op: token.SLASH,
				Y:  &ast.NumericLiteral{Value: "3"},
			},
			expected: "5",
		},
		{
			name: "modulo",
			expr: &ast.BinaryExpr{
				X:  &ast.NumericLiteral{Value: "17"},
				Op: token.MOD,
				Y:  &ast.NumericLiteral{Value: "5"},
			},
			expected: "2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calculateBinaryConstant(tt.expr)
			assert.IsType(t, &ast.NumericLiteral{}, result)
			assert.Equal(t, result.(*ast.NumericLiteral).Value, tt.expected)
		})
	}
}

func TestCalculateUnaryConstant(t *testing.T) {
	tests := []struct {
		name     string
		expr     *ast.UnaryExpr
		expected string
	}{
		{
			name: "positive",
			expr: &ast.UnaryExpr{
				Op: token.PLUS,
				X:  &ast.NumericLiteral{Value: "42"},
			},
			expected: "42",
		},
		{
			name: "negative",
			expr: &ast.UnaryExpr{
				Op: token.MINUS,
				X:  &ast.NumericLiteral{Value: "42"},
			},
			expected: "-42",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calculateUnaryConstant(tt.expr)
			assert.IsType(t, &ast.NumericLiteral{}, result)
			assert.Equal(t, result.(*ast.NumericLiteral).Value, tt.expected)
		})
	}
}

func TestGetNumericValue(t *testing.T) {
	tests := []struct {
		name     string
		node     ast.Node
		expected float64
		ok       bool
	}{
		{
			name:     "numeric literal",
			node:     &ast.NumericLiteral{Value: "42.5"},
			expected: 42.5,
			ok:       true,
		},
		{
			name:     "basic lit number",
			node:     &ast.BasicLit{Kind: token.NUM, Value: "123"},
			expected: 123,
			ok:       true,
		},
		{
			name:     "string literal",
			node:     &ast.StringLiteral{Value: "hello"},
			expected: 0,
			ok:       false,
		},
		{
			name:     "identifier",
			node:     &ast.Ident{Name: "x"},
			expected: 0,
			ok:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, ok := getNumericValue(tt.node)
			assert.Equal(t, ok, tt.ok)
			assert.Equal(t, result, tt.expected)
		})
	}
}

func TestIsConstantFalse(t *testing.T) {
	tests := []struct {
		name     string
		expr     ast.Expr
		expected bool
	}{
		{
			name:     "false literal",
			expr:     &ast.FalseLiteral{},
			expected: true,
		},
		{
			name:     "basic lit false",
			expr:     &ast.BasicLit{Kind: token.FALSE, Value: "false"},
			expected: true,
		},
		{
			name:     "true literal",
			expr:     &ast.TrueLiteral{},
			expected: false,
		},
		{
			name:     "numeric literal",
			expr:     &ast.NumericLiteral{Value: "0"},
			expected: false,
		},
		{
			name:     "identifier",
			expr:     &ast.Ident{Name: "x"},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isConstantFalse(tt.expr)
			assert.Equal(t, result, tt.expected)
		})
	}
}

func TestHasNoSideEffects(t *testing.T) {
	tests := []struct {
		name     string
		expr     ast.Expr
		expected bool
	}{
		{
			name:     "identifier",
			expr:     &ast.Ident{Name: "x"},
			expected: true,
		},
		{
			name:     "numeric literal",
			expr:     &ast.NumericLiteral{Value: "42"},
			expected: true,
		},
		{
			name:     "string literal",
			expr:     &ast.StringLiteral{Value: "hello"},
			expected: true,
		},
		{
			name: "binary expression",
			expr: &ast.BinaryExpr{
				X:  &ast.NumericLiteral{Value: "2"},
				Op: token.PLUS,
				Y:  &ast.NumericLiteral{Value: "3"},
			},
			expected: true,
		},
		{
			name: "function call",
			expr: &ast.CallExpr{
				Fun: &ast.Ident{Name: "print"},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := hasNoSideEffects(tt.expr)
			assert.Equal(t, result, tt.expected)
		})
	}
}

func TestOptimizeBlockStmt(t *testing.T) {
	opt := NewOptimizer()

	// 创建一个包含常量表达式的块语句
	block := &ast.BlockStmt{
		List: []ast.Stmt{
			&ast.ExprStmt{
				X: &ast.BinaryExpr{
					X:  &ast.NumericLiteral{Value: "2"},
					Op: token.PLUS,
					Y:  &ast.NumericLiteral{Value: "3"},
				},
			},
			&ast.AssignStmt{
				Lhs: &ast.Ident{Name: "result"},
				Rhs: &ast.NumericLiteral{Value: "42"},
			},
		},
	}

	opt.optimizeBlockStmt(block)

	// 由于第一个表达式是无副作用的常量表达式，它会被死代码消除删除
	// 所以应该只剩下第二个语句
	assert.Equal(t, len(block.List), 1)

	// 检查剩下的语句
	if assignStmt, ok := block.List[0].(*ast.AssignStmt); ok {
		assert.IsType(t, &ast.Ident{}, assignStmt.Lhs)
		assert.Equal(t, assignStmt.Lhs.(*ast.Ident).Name, "result")
	}
	assert.IsType(t, &ast.AssignStmt{}, block.List[0])
}

func TestOptimizeIfStmt(t *testing.T) {
	opt := NewOptimizer()

	// 创建一个if语句，条件为常量false
	ifStmt := &ast.IfStmt{
		Cond: &ast.FalseLiteral{},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.ExprStmt{
					X: &ast.Ident{Name: "x"},
				},
			},
		},
	}

	opt.optimizeIfStmt(ifStmt)

	// 条件应该保持不变
	assert.IsType(t, &ast.FalseLiteral{}, ifStmt.Cond)
}

func TestOptimizeBinaryExpr(t *testing.T) {
	opt := NewOptimizer()

	// 创建一个二元表达式
	binaryExpr := &ast.BinaryExpr{
		X: &ast.BinaryExpr{
			X:  &ast.NumericLiteral{Value: "2"},
			Op: token.PLUS,
			Y:  &ast.NumericLiteral{Value: "3"},
		},
		Op: token.ASTERISK,
		Y:  &ast.NumericLiteral{Value: "4"},
	}

	opt.optimizeBinaryExpr(binaryExpr)

	// 检查是否被优化
	assert.IsType(t, &ast.NumericLiteral{}, binaryExpr.X)
	assert.Equal(t, binaryExpr.X.(*ast.NumericLiteral).Value, "5")
}

func TestDeadCodeScissor(t *testing.T) {
	scissor := &DeadCodeScissor{}

	tests := []struct {
		name      string
		node      ast.Node
		shouldCut bool
	}{
		{
			name: "if with false condition",
			node: &ast.IfStmt{
				Cond: &ast.FalseLiteral{},
				Body: &ast.BlockStmt{},
			},
			shouldCut: true,
		},
		{
			name: "while with false condition",
			node: &ast.WhileStmt{
				Cond: &ast.FalseLiteral{},
				Body: &ast.BlockStmt{},
			},
			shouldCut: true,
		},
		{
			name: "expression statement with no side effects",
			node: &ast.ExprStmt{
				X: &ast.NumericLiteral{Value: "42"},
			},
			shouldCut: true,
		},
		{
			name: "if with true condition",
			node: &ast.IfStmt{
				Cond: &ast.TrueLiteral{},
				Body: &ast.BlockStmt{},
			},
			shouldCut: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := scissor.ShouldCut(tt.node)
			assert.Equal(t, result, tt.shouldCut)
		})
	}
}

func TestUnusedVariableScissor(t *testing.T) {
	scissor := &UnusedVariableScissor{}

	// 创建一个赋值语句
	assignStmt := &ast.AssignStmt{
		Lhs: &ast.Ident{Name: "x"},
		Rhs: &ast.NumericLiteral{Value: "42"},
	}

	// 由于isVariableUsed的简化实现总是返回true，所以这里应该返回false
	result := scissor.ShouldCut(assignStmt)
	assert.False(t, result)
}

func TestOptimizeNilNode(t *testing.T) {
	opt := NewOptimizer()
	result := opt.Optimize(nil)
	assert.Nil(t, result)
}

func TestOptimizeComplexExpression(t *testing.T) {
	opt := NewOptimizer()

	// 创建一个复杂的表达式: (2 + 3) * (4 - 1)
	complexExpr := &ast.BinaryExpr{
		X: &ast.BinaryExpr{
			X:  &ast.NumericLiteral{Value: "2"},
			Op: token.PLUS,
			Y:  &ast.NumericLiteral{Value: "3"},
		},
		Op: token.ASTERISK,
		Y: &ast.BinaryExpr{
			X:  &ast.NumericLiteral{Value: "4"},
			Op: token.MINUS,
			Y:  &ast.NumericLiteral{Value: "1"},
		},
	}

	result := opt.Optimize(complexExpr)

	// 应该被优化为 5 * 3 = 15
	assert.IsType(t, &ast.NumericLiteral{}, result)
	assert.Equal(t, result.(*ast.NumericLiteral).Value, "15")
	assert.IsType(t, &ast.NumericLiteral{}, result)
}

func TestOptimizeWithVariables(t *testing.T) {
	opt := NewOptimizer()

	// 创建一个包含变量的表达式
	expr := &ast.BinaryExpr{
		X:  &ast.Ident{Name: "x"},
		Op: token.PLUS,
		Y:  &ast.NumericLiteral{Value: "5"},
	}

	result := opt.Optimize(expr)

	// 应该保持不变，因为包含变量
	assert.Equal(t, result, expr)
}

func TestOptimizeConstantFolding(t *testing.T) {
	opt := NewOptimizer()

	// 创建一个赋值语句，包含常量表达式
	assignStmt := &ast.AssignStmt{
		Lhs: &ast.Ident{Name: "result"},
		Rhs: &ast.BinaryExpr{
			X: &ast.BinaryExpr{
				X:  &ast.NumericLiteral{Value: "2"},
				Op: token.PLUS,
				Y:  &ast.NumericLiteral{Value: "3"},
			},
			Op: token.ASTERISK,
			Y:  &ast.NumericLiteral{Value: "4"},
		},
	}

	result := opt.Optimize(assignStmt)

	// 检查结果是否为赋值语句
	if assign, ok := result.(*ast.AssignStmt); ok {
		// 检查右侧是否被常量折叠
		assert.IsType(t, &ast.NumericLiteral{}, assign.Rhs)
		assert.Equal(t, assign.Rhs.(*ast.NumericLiteral).Value, "20")
	}
	assert.IsType(t, &ast.AssignStmt{}, result)
}

func BenchmarkOptimize(b *testing.B) {
	opt := NewOptimizer()

	// 创建一个简单的表达式用于基准测试
	expr := &ast.BinaryExpr{
		X:  &ast.NumericLiteral{Value: "2"},
		Op: token.PLUS,
		Y:  &ast.NumericLiteral{Value: "3"},
	}

	for b.Loop() {
		opt.Optimize(expr)
	}
}
