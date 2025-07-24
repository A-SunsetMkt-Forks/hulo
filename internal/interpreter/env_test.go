package interpreter

import (
	"testing"

	"github.com/hulo-lang/hulo/internal/object"
	"github.com/hulo-lang/hulo/syntax/hulo/token"
)

func TestEnvironment_DeclareAndGet(t *testing.T) {
	env := NewEnvironment()

	// 测试 const 变量
	constValue := &object.NumberValue{}
	constValue.Value.SetString("42")
	err := env.Declare("PI", constValue, token.CONST)
	if err != nil {
		t.Errorf("Failed to declare const variable: %v", err)
	}

	// 测试 let 变量
	letValue := &object.StringValue{Value: "hello"}
	err = env.Declare("message", letValue, token.LET)
	if err != nil {
		t.Errorf("Failed to declare let variable: %v", err)
	}

	// 测试 var 变量
	varValue := &object.BoolValue{Value: true}
	err = env.Declare("flag", varValue, token.VAR)
	if err != nil {
		t.Errorf("Failed to declare var variable: %v", err)
	}

	// 测试获取变量
	if value, ok := env.GetValue("PI"); !ok || value != constValue {
		t.Errorf("Failed to get const variable PI")
	}

	if value, ok := env.GetValue("message"); !ok || value != letValue {
		t.Errorf("Failed to get let variable message")
	}

	if value, ok := env.GetValue("flag"); !ok || value != varValue {
		t.Errorf("Failed to get var variable flag")
	}

	// 测试作用域
	if scope, ok := env.GetScope("PI"); !ok || scope != token.CONST {
		t.Errorf("Expected PI to have CONST scope, got %v", scope)
	}

	if scope, ok := env.GetScope("message"); !ok || scope != token.LET {
		t.Errorf("Expected message to have LET scope, got %v", scope)
	}

	if scope, ok := env.GetScope("flag"); !ok || scope != token.VAR {
		t.Errorf("Expected flag to have VAR scope, got %v", scope)
	}

	// 测试只读属性
	if !env.IsReadOnly("PI") {
		t.Errorf("Expected PI to be read-only")
	}

	if env.IsReadOnly("message") {
		t.Errorf("Expected message to be writable")
	}

	if env.IsReadOnly("flag") {
		t.Errorf("Expected flag to be writable")
	}
}

func TestEnvironment_Reassign(t *testing.T) {
	env := NewEnvironment()

	// 声明变量
	initialValue := &object.NumberValue{}
	initialValue.Value.SetString("10")
	err := env.Declare("x", initialValue, token.VAR)
	if err != nil {
		t.Errorf("Failed to declare variable: %v", err)
	}

	// 重新赋值
	newValue := &object.NumberValue{}
	newValue.Value.SetString("20")
	err = env.Assign("x", newValue)
	if err != nil {
		t.Errorf("Failed to reassign variable: %v", err)
	}

	// 验证新值
	if value, ok := env.GetValue("x"); !ok || value != newValue {
		t.Errorf("Failed to get reassigned value")
	}
}

func TestEnvironment_ConstReassign(t *testing.T) {
	env := NewEnvironment()

	// 声明 const 变量
	constValue := &object.NumberValue{}
	constValue.Value.SetString("42")
	err := env.Declare("PI", constValue, token.CONST)
	if err != nil {
		t.Errorf("Failed to declare const variable: %v", err)
	}

	// 尝试重新赋值（应该失败）
	newValue := &object.NumberValue{}
	newValue.Value.SetString("3.14")
	err = env.Assign("PI", newValue)
	if err == nil {
		t.Errorf("Expected error when reassigning const variable")
	}

	// 验证原值未改变
	if value, ok := env.GetValue("PI"); !ok || value != constValue {
		t.Errorf("Const variable value should not change")
	}
}

func TestEnvironment_DuplicateDeclaration(t *testing.T) {
	env := NewEnvironment()

	// 第一次声明
	value1 := &object.NumberValue{}
	value1.Value.SetString("1")
	err := env.Declare("x", value1, token.VAR)
	if err != nil {
		t.Errorf("Failed to declare variable: %v", err)
	}

	// 第二次声明（应该失败）
	value2 := &object.NumberValue{}
	value2.Value.SetString("2")
	err = env.Declare("x", value2, token.LET)
	if err == nil {
		t.Errorf("Expected error when declaring duplicate variable")
	}

	// 验证原值未改变
	if value, ok := env.GetValue("x"); !ok || value != value1 {
		t.Errorf("Original variable value should not change")
	}
}

func TestEnvironment_OuterScope(t *testing.T) {
	outer := NewEnvironment()
	inner := outer.Fork()

	// 在外层声明变量
	outerValue := &object.StringValue{Value: "outer"}
	err := outer.Declare("x", outerValue, token.VAR)
	if err != nil {
		t.Errorf("Failed to declare variable in outer scope: %v", err)
	}

	// 在内层获取变量
	if value, ok := inner.GetValue("x"); !ok || value != outerValue {
		t.Errorf("Failed to get variable from outer scope")
	}

	// 在内层重新赋值
	innerValue := &object.StringValue{Value: "inner"}
	err = inner.Assign("x", innerValue)
	if err != nil {
		t.Errorf("Failed to reassign variable in inner scope: %v", err)
	}

	// 验证内层的新值
	if value, ok := inner.GetValue("x"); !ok || value != innerValue {
		t.Errorf("Failed to get reassigned value in inner scope")
	}

	// 验证外层值也改变了（因为内层修改了外层的变量）
	if value, ok := outer.GetValue("x"); !ok || value != innerValue {
		t.Errorf("Outer scope value should be updated when inner scope modifies it")
	}
}
