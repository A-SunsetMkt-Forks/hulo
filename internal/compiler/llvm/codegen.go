// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package llvm

import (
	"fmt"

	"github.com/hulo-lang/hulo/internal/config"
	"github.com/hulo-lang/hulo/syntax/hulo/ast"
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

type LLVMCodegen struct {
	opts        *config.LLVMOpts
	module      *ir.Module
	currentFunc *ir.Func
	locals      map[string]*ir.InstAlloca
	blockStack  []*ir.Block
}

func NewLLVMCodegen(opts *config.LLVMOpts) *LLVMCodegen {
	return &LLVMCodegen{
		opts:       opts,
		module:     ir.NewModule(),
		locals:     make(map[string]*ir.InstAlloca),
		blockStack: make([]*ir.Block, 0),
	}
}

func (c *LLVMCodegen) CodegenFile(node *ast.File) (*ir.Module, error) {
	// 添加标准库函数声明
	c.addStandardLibrary()

	// 处理文件中的所有语句
	for _, stmt := range node.Stmts {
		if err := c.codegenStmt(stmt); err != nil {
			return nil, err
		}
	}

	return c.module, nil
}

func (c *LLVMCodegen) CodegenBlockStmt(node *ast.BlockStmt) (*ir.Block, error) {
	if c.currentFunc == nil {
		return nil, fmt.Errorf("no current function")
	}

	block := c.currentFunc.NewBlock("")
	c.pushBlock(block)
	defer c.popBlock()

	// 处理块中的每个语句
	for _, stmt := range node.List {
		if err := c.codegenStmt(stmt); err != nil {
			return nil, err
		}
	}

	return block, nil
}

func (c *LLVMCodegen) CodegenIfStmt(node *ast.IfStmt) (*ir.Block, error) {
	if c.currentFunc == nil {
		return nil, fmt.Errorf("no current function")
	}

	// 生成条件表达式
	cond, err := c.codegenExpr(node.Cond)
	if err != nil {
		return nil, err
	}

	// 创建基本块
	thenBlock := c.currentFunc.NewBlock("then")
	elseBlock := c.currentFunc.NewBlock("else")
	mergeBlock := c.currentFunc.NewBlock("merge")

	// 条件跳转
	currentBlock := c.getCurrentBlock()
	currentBlock.NewCondBr(cond, thenBlock, elseBlock)

	// 生成then块
	c.pushBlock(thenBlock)
	if err := c.codegenStmt(node.Body); err != nil {
		return nil, err
	}
	thenBlock.NewBr(mergeBlock)
	c.popBlock()

	// 生成else块
	if node.Else != nil {
		c.pushBlock(elseBlock)
		if err := c.codegenStmt(node.Else); err != nil {
			return nil, err
		}
		elseBlock.NewBr(mergeBlock)
		c.popBlock()
	} else {
		elseBlock.NewBr(mergeBlock)
	}

	c.pushBlock(mergeBlock)
	return mergeBlock, nil
}

func (c *LLVMCodegen) CodegenForStmt(node *ast.ForStmt) ([]*ir.Block, error) {
	if c.currentFunc == nil {
		return nil, fmt.Errorf("no current function")
	}

	// 创建循环的基本块
	initBlock := c.currentFunc.NewBlock("for.init")
	condBlock := c.currentFunc.NewBlock("for.cond")
	bodyBlock := c.currentFunc.NewBlock("for.body")
	postBlock := c.currentFunc.NewBlock("for.post")
	exitBlock := c.currentFunc.NewBlock("for.exit")

	// 生成初始化语句
	if node.Init != nil {
		c.pushBlock(initBlock)
		if err := c.codegenStmt(node.Init); err != nil {
			return nil, err
		}
		initBlock.NewBr(condBlock)
		c.popBlock()
	} else {
		initBlock.NewBr(condBlock)
	}

	// 生成条件表达式
	c.pushBlock(condBlock)
	cond, err := c.codegenExpr(node.Cond)
	if err != nil {
		return nil, err
	}
	condBlock.NewCondBr(cond, bodyBlock, exitBlock)
	c.popBlock()

	// 生成循环体
	c.pushBlock(bodyBlock)
	if err := c.codegenStmt(node.Body); err != nil {
		return nil, err
	}
	bodyBlock.NewBr(postBlock)
	c.popBlock()

	// 生成后置表达式
	if node.Post != nil {
		c.pushBlock(postBlock)
		if _, err := c.codegenExpr(node.Post); err != nil {
			return nil, err
		}
		postBlock.NewBr(condBlock)
		c.popBlock()
	} else {
		postBlock.NewBr(condBlock)
	}

	c.pushBlock(exitBlock)
	return []*ir.Block{initBlock, condBlock, bodyBlock, postBlock, exitBlock}, nil
}

func (c *LLVMCodegen) CodegenForInStmt(node *ast.ForInStmt) ([]*ir.Block, error) {
	if c.currentFunc == nil {
		return nil, fmt.Errorf("no current function")
	}

	// 创建循环的基本块
	initBlock := c.currentFunc.NewBlock("forin.init")
	condBlock := c.currentFunc.NewBlock("forin.cond")
	bodyBlock := c.currentFunc.NewBlock("forin.body")
	exitBlock := c.currentFunc.NewBlock("forin.exit")

	// 分配循环变量
	varName := node.Index.Name
	varAlloca := initBlock.NewAlloca(types.I32)
	c.locals[varName] = varAlloca

	// rge := node.RangeExpr

	// 初始化循环变量
	initBlock.NewStore(constant.NewInt(types.I32, 0), varAlloca)
	initBlock.NewBr(condBlock)

	// 生成条件检查 - 简化版本
	c.pushBlock(condBlock)
	index := condBlock.NewLoad(types.I32, varAlloca)

	// 假设循环10次
	maxIterations := constant.NewInt(types.I32, 10)
	condition := condBlock.NewICmp(enum.IPredSLT, index, maxIterations)
	condBlock.NewCondBr(condition, bodyBlock, exitBlock)
	c.popBlock()

	// 生成循环体
	c.pushBlock(bodyBlock)
	if err := c.codegenStmt(node.Body); err != nil {
		return nil, err
	}
	// 递增索引
	newIndex := bodyBlock.NewAdd(index, constant.NewInt(types.I32, 1))
	bodyBlock.NewStore(newIndex, varAlloca)
	bodyBlock.NewBr(condBlock)
	c.popBlock()

	c.pushBlock(exitBlock)
	return []*ir.Block{initBlock, condBlock, bodyBlock, exitBlock}, nil
}

func (c *LLVMCodegen) CodegenFuncDecl(node *ast.FuncDecl) (*ir.Func, error) {
	// 确定函数返回类型
	var returnType types.Type = types.Void

	if node.Type != nil {
		if t, err := c.convertType(node.Type); err == nil {
			returnType = t
		}
	}

	// 创建函数参数
	var params []*ir.Param
	for _, param := range node.Recv {
		if ident, ok := param.(*ast.Ident); ok {
			paramType := types.I32 // 默认类型，实际应该从AST中获取
			params = append(params, ir.NewParam(ident.Name, paramType))
		}
	}

	// 创建函数
	fn := c.module.NewFunc(node.Name.Name, returnType, params...)

	// 保存当前函数
	oldFunc := c.currentFunc
	c.currentFunc = fn
	defer func() { c.currentFunc = oldFunc }()

	// 创建函数体
	entryBlock := fn.NewBlock("")
	c.pushBlock(entryBlock)

	// 处理函数体
	if node.Body != nil {
		if err := c.codegenStmt(node.Body); err != nil {
			return nil, err
		}
	}

	// 如果没有显式的return语句，添加默认return
	if returnType != types.Void {
		entryBlock.NewRet(constant.NewInt(types.I32, 0))
	} else {
		entryBlock.NewRet(nil)
	}

	c.popBlock()
	return fn, nil
}

func (c *LLVMCodegen) CodegenExprStmt(node *ast.ExprStmt) (*ir.Module, error) {
	if c.currentFunc == nil {
		return nil, fmt.Errorf("no current function")
	}

	_, err := c.codegenExpr(node.X)
	return c.module, err
}

func (c *LLVMCodegen) CodegenAssignStmt(node *ast.AssignStmt) (*ir.Module, error) {
	if c.currentFunc == nil {
		return nil, fmt.Errorf("no current function")
	}

	// 生成右值表达式
	value, err := c.codegenExpr(node.Rhs)
	if err != nil {
		return nil, err
	}

	// 处理左值
	if ident, ok := node.Lhs.(*ast.Ident); ok {
		currentBlock := c.getCurrentBlock()
		if currentBlock == nil {
			return nil, fmt.Errorf("no current block")
		}

		// 检查是否已存在局部变量
		if alloca, exists := c.locals[ident.Name]; exists {
			currentBlock.NewStore(value, alloca)
		} else {
			// 创建新的局部变量 - 修正：在基本块中分配
			alloca := currentBlock.NewAlloca(value.Type())
			c.locals[ident.Name] = alloca
			currentBlock.NewStore(value, alloca)
		}
	}

	return c.module, nil
}

func (c *LLVMCodegen) CodegenReturnStmt(node *ast.ReturnStmt) (*ir.Module, error) {
	if c.currentFunc == nil {
		return nil, fmt.Errorf("no current function")
	}

	if node.X != nil {
		value, err := c.codegenExpr(node.X)
		if err != nil {
			return nil, err
		}
		c.getCurrentBlock().NewRet(value)
	} else {
		c.getCurrentBlock().NewRet(nil)
	}

	return c.module, nil
}

func (c *LLVMCodegen) Codegen(node ast.Node) (any, error) {
	switch node := node.(type) {
	case *ast.File:
		return c.CodegenFile(node)
	case *ast.BlockStmt:
		return c.CodegenBlockStmt(node)
	case *ast.IfStmt:
		return c.CodegenIfStmt(node)
	case *ast.ForStmt:
		return c.CodegenForStmt(node)
	case *ast.ForInStmt:
		return c.CodegenForInStmt(node)
	case *ast.FuncDecl:
		return c.CodegenFuncDecl(node)
	case *ast.ExprStmt:
		return c.CodegenExprStmt(node)
	case *ast.AssignStmt:
		return c.CodegenAssignStmt(node)
	case *ast.ReturnStmt:
		return c.CodegenReturnStmt(node)
	}
	return nil, fmt.Errorf("unsupported node type: %T", node)
}

// 辅助方法
func (c *LLVMCodegen) codegenStmt(stmt ast.Stmt) error {
	_, err := c.Codegen(stmt)
	return err
}

func (c *LLVMCodegen) codegenExpr(_ ast.Expr) (value.Value, error) {
	// 这里需要实现表达式代码生成
	// 简化版本，返回默认值
	return constant.NewInt(types.I32, 0), nil
}

func (c *LLVMCodegen) convertType(_ ast.Expr) (types.Type, error) {
	// 这里需要实现类型转换
	// 简化版本，返回默认类型
	return types.I32, nil
}

func (c *LLVMCodegen) addStandardLibrary() {
	// 添加printf函数声明
	printf := c.module.NewFunc("printf", types.I32, ir.NewParam("", types.NewPointer(types.I8)))
	printf.Sig.Variadic = true

	// 添加其他标准库函数
	c.module.NewFunc("malloc", types.NewPointer(types.I8), ir.NewParam("", types.I64))
	c.module.NewFunc("free", types.Void, ir.NewParam("", types.NewPointer(types.I8)))
}

func (c *LLVMCodegen) pushBlock(block *ir.Block) {
	c.blockStack = append(c.blockStack, block)
}

func (c *LLVMCodegen) popBlock() {
	if len(c.blockStack) > 0 {
		c.blockStack = c.blockStack[:len(c.blockStack)-1]
	}
}

func (c *LLVMCodegen) getCurrentBlock() *ir.Block {
	if len(c.blockStack) > 0 {
		return c.blockStack[len(c.blockStack)-1]
	}
	return nil
}
