// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package interpreter

import (
	"fmt"
	"math/big"
	"strconv"

	"github.com/hulo-lang/hulo/internal/core"
	"github.com/hulo-lang/hulo/internal/object"
	"github.com/hulo-lang/hulo/syntax/hulo/ast"
	"github.com/hulo-lang/hulo/syntax/hulo/token"
)

// FunctionEvaluatorImpl 函数求值器实现
type FunctionEvaluatorImpl struct {
	interpreter *Interpreter
}

// 确保 FunctionEvaluatorImpl 实现 core.FunctionEvaluator 接口
var _ core.FunctionEvaluator = (*FunctionEvaluatorImpl)(nil)

// NewFunctionEvaluator 创建函数求值器
func NewFunctionEvaluator(interpreter *Interpreter) *FunctionEvaluatorImpl {
	return &FunctionEvaluatorImpl{
		interpreter: interpreter,
	}
}

// EvaluateFunction 求值函数
func (fe *FunctionEvaluatorImpl) EvaluateFunction(funcDecl *ast.FuncDecl, args []interface{}) (interface{}, error) {
	// 1. 创建函数执行环境
	funcEnv := fe.interpreter.env.Fork()

	// 2. 绑定参数到环境
	if err := fe.bindParameters(funcDecl, args, funcEnv); err != nil {
		return nil, fmt.Errorf("failed to bind parameters: %w", err)
	}

	// 3. 临时切换环境
	originalEnv := fe.interpreter.env
	fe.interpreter.env = funcEnv
	defer func() {
		fe.interpreter.env = originalEnv
	}()

	// 4. 执行函数体
	if funcDecl.Body != nil {
		result := fe.interpreter.Eval(funcDecl.Body)
		if result != nil {
			// 将结果转换为 interface{}
			return fe.convertResultToInterface(result)
		}
	}

	// 5. 如果没有函数体或没有返回值，返回 nil
	return nil, nil
}

// bindParameters 绑定参数到环境
func (fe *FunctionEvaluatorImpl) bindParameters(funcDecl *ast.FuncDecl, args []interface{}, env *Environment) error {
	// 1. 处理位置参数
	for i, param := range funcDecl.Recv {
		if i >= len(args) {
			return fmt.Errorf("not enough arguments for function %s", funcDecl.Name.Name)
		}

		// 获取参数名
		var paramName string
		if ident, ok := param.(*ast.Ident); ok {
			paramName = ident.Name
		} else {
			paramName = fmt.Sprintf("param%d", i)
		}

		// 将 interface{} 转换为 object.Value
		value := fe.convertInterfaceToValue(args[i])

		// 绑定参数到环境
		env.Set(paramName, value)
	}

	return nil
}

// convertResultToInterface 将求值结果转换为 interface{}
func (fe *FunctionEvaluatorImpl) convertResultToInterface(result ast.Node) (interface{}, error) {
	switch result := result.(type) {
	case *ast.BasicLit:
		return fe.convertBasicLitToInterface(result)
	case *ast.Ident:
		// 从环境中获取变量值
		if value, found := fe.interpreter.env.GetValue(result.Name); found {
			return value.Interface(), nil
		}
		return nil, nil
	default:
		// 对于复杂的表达式，可能需要递归求值
		return nil, nil
	}
}

// convertBasicLitToInterface 转换基本字面量为 interface{}
func (fe *FunctionEvaluatorImpl) convertBasicLitToInterface(lit *ast.BasicLit) (interface{}, error) {
	switch lit.Kind {
	case token.NUM:
		// 尝试解析为数字
		if f, err := strconv.ParseFloat(lit.Value, 64); err == nil {
			return f, nil
		}
		return lit.Value, nil
	case token.STR:
		return lit.Value, nil
	case token.TRUE:
		return true, nil
	case token.FALSE:
		return false, nil
	case token.NULL:
		return nil, nil
	default:
		return nil, fmt.Errorf("unsupported literal type: %v", lit.Kind)
	}
}

// convertInterfaceToValue 将 interface{} 转换为 object.Value
func (fe *FunctionEvaluatorImpl) convertInterfaceToValue(v interface{}) object.Value {
	switch val := v.(type) {
	case string:
		return &object.StringValue{Value: val}
	case int:
		num := &object.NumberValue{Value: big.NewFloat(float64(val))}
		return num
	case float64:
		num := &object.NumberValue{Value: big.NewFloat(val)}
		return num
	case bool:
		if val {
			return object.TRUE
		}
		return object.FALSE
	case nil:
		return object.NULL
	default:
		// 对于复杂类型，返回 null
		return object.NULL
	}
}
