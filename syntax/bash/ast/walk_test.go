// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package ast_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hulo-lang/hulo/syntax/bash/ast"
	"github.com/hulo-lang/hulo/syntax/bash/token"
	"github.com/stretchr/testify/assert"
)

// testVisitor 是一个测试用的访问者
type testVisitor struct {
	visitFunc func(ast.Node) ast.Visitor
}

func (v *testVisitor) Visit(node ast.Node) ast.Visitor {
	return v.visitFunc(node)
}

func TestWalkBasic(t *testing.T) {
	// 创建一个简单的AST
	file := &ast.File{
		Stmts: []ast.Stmt{
			&ast.ExprStmt{
				X: &ast.CmdExpr{
					Name: &ast.Ident{Name: "echo"},
					Recv: []ast.Expr{&ast.Word{Val: "Hello"}},
				},
			},
			&ast.AssignStmt{
				Lhs: &ast.Ident{Name: "var"},
				Rhs: &ast.Word{Val: "value"},
			},
		},
	}

	// 测试基础遍历
	var visited []string
	var visitor ast.Visitor
	visitor = &testVisitor{
		visitFunc: func(node ast.Node) ast.Visitor {
			switch node.(type) {
			case *ast.File:
				visited = append(visited, "File")
			case *ast.ExprStmt:
				visited = append(visited, "ExprStmt")
			case *ast.CmdExpr:
				visited = append(visited, "CmdExpr")
			case *ast.Ident:
				visited = append(visited, "Ident")
			case *ast.Word:
				visited = append(visited, "Word")
			case *ast.AssignStmt:
				visited = append(visited, "AssignStmt")
			}
			return visitor
		},
	}

	ast.Walk(visitor, file)

	expected := []string{
		"File", "ExprStmt", "CmdExpr", "Ident", "Word", "AssignStmt", "Ident", "Word",
	}
	assert.Equal(t, expected, visited)
}

func TestWalkNil(t *testing.T) {
	// 测试遍历nil节点
	var visited []string
	var visitor ast.Visitor
	visitor = &testVisitor{
		visitFunc: func(node ast.Node) ast.Visitor {
			if node == nil {
				visited = append(visited, "nil")
			} else {
				visited = append(visited, "not-nil")
			}
			return visitor
		},
	}

	ast.Walk(visitor, nil)
	assert.Equal(t, []string{"nil"}, visited)
}

func TestWalkIf(t *testing.T) {
	// 创建AST
	file := &ast.File{
		Stmts: []ast.Stmt{
			&ast.ExprStmt{
				X: &ast.CmdExpr{
					Name: &ast.Ident{Name: "echo"},
					Recv: []ast.Expr{&ast.Word{Val: "Hello"}},
				},
			},
			&ast.AssignStmt{
				Lhs: &ast.Ident{Name: "var"},
				Rhs: &ast.Word{Val: "value"},
			},
		},
	}

	// 只遍历包含CmdExpr的路径
	var visited []string
	ast.WalkIf(file, func(node ast.Node) bool {
		switch node.(type) {
		case *ast.File, *ast.ExprStmt, *ast.CmdExpr, *ast.Ident, *ast.Word:
			// 允许遍历这些节点类型
			visited = append(visited, fmt.Sprintf("%T", node))
			return true
		default:
			// 不遍历其他类型的节点
			return false
		}
	})

	// 应该访问File, ExprStmt, CmdExpr, Ident, Word
	assert.Contains(t, visited, "*ast.File")
	assert.Contains(t, visited, "*ast.ExprStmt")
	assert.Contains(t, visited, "*ast.CmdExpr")
	assert.Contains(t, visited, "*ast.Ident")
	assert.Contains(t, visited, "*ast.Word")
	assert.NotContains(t, visited, "*ast.AssignStmt")
}

func TestWalkIfWithMultipleConditions(t *testing.T) {
	// 创建AST
	file := &ast.File{
		Stmts: []ast.Stmt{
			&ast.ExprStmt{
				X: &ast.CmdExpr{
					Name: &ast.Ident{Name: "echo"},
					Recv: []ast.Expr{&ast.Word{Val: "Hello"}},
				},
			},
			&ast.AssignStmt{
				Lhs: &ast.Ident{Name: "var"},
				Rhs: &ast.Word{Val: "value"},
			},
			&ast.ExprStmt{
				X: &ast.CmdExpr{
					Name: &ast.Ident{Name: "ls"},
					Recv: []ast.Expr{&ast.Word{Val: "-la"}},
				},
			},
		},
	}

	// 只遍历包含CmdExpr的路径
	var visited []string
	ast.WalkIf(file, func(node ast.Node) bool {
		switch node.(type) {
		case *ast.File, *ast.ExprStmt, *ast.CmdExpr, *ast.Ident, *ast.Word:
			// 允许遍历这些节点类型
			visited = append(visited, fmt.Sprintf("%T", node))
			return true
		default:
			// 不遍历其他类型的节点
			return false
		}
	})

	// 应该访问File, ExprStmt, CmdExpr, Ident, Word
	assert.Contains(t, visited, "*ast.File")
	assert.Contains(t, visited, "*ast.ExprStmt")
	assert.Contains(t, visited, "*ast.CmdExpr")
	assert.Contains(t, visited, "*ast.Ident")
	assert.Contains(t, visited, "*ast.Word")
	assert.NotContains(t, visited, "*ast.AssignStmt")

	// 应该有2个CmdExpr
	cmdExprCount := 0
	for _, v := range visited {
		if v == "*ast.CmdExpr" {
			cmdExprCount++
		}
	}
	assert.Equal(t, 2, cmdExprCount)
}

func TestWalkComplexStructure(t *testing.T) {
	// 创建一个复杂的AST结构
	file := &ast.File{
		Stmts: []ast.Stmt{
			&ast.FuncDecl{
				Name: &ast.Ident{Name: "main"},
				Body: &ast.BlockStmt{
					List: []ast.Stmt{
						&ast.IfStmt{
							Cond: &ast.TestExpr{
								X: &ast.BinaryExpr{
									X:  &ast.VarExpExpr{X: &ast.Ident{Name: "var"}},
									Op: token.Equal,
									Y:  &ast.Word{Val: "test"},
								},
							},
							Body: &ast.BlockStmt{
								List: []ast.Stmt{
									&ast.ExprStmt{
										X: &ast.CmdExpr{
											Name: &ast.Ident{Name: "echo"},
											Recv: []ast.Expr{&ast.Word{Val: "Match"}},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	var nodeTypes []string
	var visitor ast.Visitor
	visitor = &testVisitor{
		visitFunc: func(node ast.Node) ast.Visitor {
			switch node.(type) {
			case *ast.File:
				nodeTypes = append(nodeTypes, "File")
			case *ast.FuncDecl:
				nodeTypes = append(nodeTypes, "FuncDecl")
			case *ast.BlockStmt:
				nodeTypes = append(nodeTypes, "BlockStmt")
			case *ast.IfStmt:
				nodeTypes = append(nodeTypes, "IfStmt")
			case *ast.TestExpr:
				nodeTypes = append(nodeTypes, "TestExpr")
			case *ast.BinaryExpr:
				nodeTypes = append(nodeTypes, "BinaryExpr")
			case *ast.VarExpExpr:
				nodeTypes = append(nodeTypes, "VarExpExpr")
			case *ast.Ident:
				nodeTypes = append(nodeTypes, "Ident")
			case *ast.Word:
				nodeTypes = append(nodeTypes, "Word")
			case *ast.ExprStmt:
				nodeTypes = append(nodeTypes, "ExprStmt")
			case *ast.CmdExpr:
				nodeTypes = append(nodeTypes, "CmdExpr")
			}
			return visitor
		},
	}

	ast.Walk(visitor, file)

	// 验证所有节点类型都被访问
	expectedTypes := []string{
		"File", "FuncDecl", "Ident", "BlockStmt", "IfStmt", "TestExpr", "BinaryExpr",
		"VarExpExpr", "Ident", "Word", "BlockStmt", "ExprStmt", "CmdExpr", "Ident", "Word",
	}
	assert.Equal(t, expectedTypes, nodeTypes)
}

func TestWalkEarlyTermination(t *testing.T) {
	// 创建AST
	file := &ast.File{
		Stmts: []ast.Stmt{
			&ast.ExprStmt{
				X: &ast.CmdExpr{
					Name: &ast.Ident{Name: "echo"},
					Recv: []ast.Expr{&ast.Word{Val: "First"}},
				},
			},
			&ast.ExprStmt{
				X: &ast.CmdExpr{
					Name: &ast.Ident{Name: "echo"},
					Recv: []ast.Expr{&ast.Word{Val: "Second"}},
				},
			},
			&ast.ExprStmt{
				X: &ast.CmdExpr{
					Name: &ast.Ident{Name: "echo"},
					Recv: []ast.Expr{&ast.Word{Val: "Third"}},
				},
			},
		},
	}

	var visited []string
	var shouldStop bool
	var visitor ast.Visitor
	visitor = &testVisitor{
		visitFunc: func(node ast.Node) ast.Visitor {
			if shouldStop {
				return nil
			}

			if cmdExpr, ok := node.(*ast.CmdExpr); ok {
				visited = append(visited, cmdExpr.Name.Name)
				// 在第二个echo后停止遍历
				if len(visited) == 2 {
					shouldStop = true
					return nil
				}
			}
			return visitor
		},
	}

	ast.Walk(visitor, file)
	assert.Equal(t, []string{"echo", "echo"}, visited)
	assert.Len(t, visited, 2)
}

func TestWalkEmptyFile(t *testing.T) {
	file := &ast.File{
		Stmts: []ast.Stmt{},
	}

	var visited []string
	var visitor ast.Visitor
	visitor = &testVisitor{
		visitFunc: func(node ast.Node) ast.Visitor {
			switch node.(type) {
			case *ast.File:
				visited = append(visited, "File")
			}
			return visitor
		},
	}

	ast.Walk(visitor, file)
	assert.Equal(t, []string{"File"}, visited)
}

func TestWalkSingleNode(t *testing.T) {
	node := &ast.Ident{Name: "test"}

	var visited []string
	var visitor ast.Visitor
	visitor = &testVisitor{
		visitFunc: func(node ast.Node) ast.Visitor {
			switch node.(type) {
			case *ast.Ident:
				visited = append(visited, "Ident")
			}
			return visitor
		},
	}

	ast.Walk(visitor, node)
	assert.Equal(t, []string{"Ident"}, visited)
}

func TestWalkWithInspect(t *testing.T) {
	// 测试Walk与Inspect的集成
	file := &ast.File{
		Stmts: []ast.Stmt{
			&ast.ExprStmt{
				X: &ast.CmdExpr{
					Name: &ast.Ident{Name: "echo"},
					Recv: []ast.Expr{&ast.Word{Val: "Hello"}},
				},
			},
		},
	}

	var buf strings.Builder
	ast.Inspect(file, &buf)

	output := buf.String()
	assert.Contains(t, output, "*ast.File")
	assert.Contains(t, output, "*ast.ExprStmt")
	assert.Contains(t, output, "*ast.CmdExpr")
	assert.Contains(t, output, "*ast.Ident")
	assert.Contains(t, output, "*ast.Word")
}

func TestWalkWithPrint(t *testing.T) {
	// 测试Walk与Print的集成
	file := &ast.File{
		Stmts: []ast.Stmt{
			&ast.ExprStmt{
				X: &ast.CmdExpr{
					Name: &ast.Ident{Name: "echo"},
					Recv: []ast.Expr{&ast.Word{Val: "Hello"}},
				},
			},
		},
	}

	var buf strings.Builder
	ast.Write(file, &buf)

	output := buf.String()
	assert.Contains(t, output, "echo")
	assert.Contains(t, output, "Hello")
}

func TestWalkWithInvalidNode(t *testing.T) {
	// 测试处理无效节点
	var visited []string
	var visitor ast.Visitor
	visitor = &testVisitor{
		visitFunc: func(node ast.Node) ast.Visitor {
			if node == nil {
				visited = append(visited, "nil")
			} else {
				visited = append(visited, "valid")
			}
			return visitor
		},
	}

	// 传递nil节点
	ast.Walk(visitor, nil)
	assert.Equal(t, []string{"nil"}, visited)
}
