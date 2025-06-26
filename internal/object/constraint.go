// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package object

import (
	"fmt"
	"math/big"
	"strings"
)

type Constraint interface {
	SatisfiedBy(Type) bool
	String() string
}

type InterfaceConstraint struct {
	interfaceType Type
}

func NewInterfaceConstraint(interfaceType Type) *InterfaceConstraint {
	return &InterfaceConstraint{
		interfaceType: interfaceType,
	}
}

func (ic *InterfaceConstraint) SatisfiedBy(t Type) bool {
	// 如果约束类型是trait，使用Implements检查
	if trait, ok := ic.interfaceType.(*TraitType); ok {
		return t.Implements(trait)
	}

	// 如果约束类型是具体类型，检查类型兼容性
	return t.AssignableTo(ic.interfaceType) || t.Name() == ic.interfaceType.Name()
}

func (ic *InterfaceConstraint) String() string {
	return fmt.Sprintf("implements %s", ic.interfaceType.Name())
}

type OperatorConstraint struct {
	operator   Operator
	rightType  Type
	returnType Type
}

func NewOperatorConstraint(op Operator, rightType, returnType Type) *OperatorConstraint {
	return &OperatorConstraint{
		operator:   op,
		rightType:  rightType,
		returnType: returnType,
	}
}

func (oc *OperatorConstraint) SatisfiedBy(t Type) bool {
	if objType, ok := t.(*ObjectType); ok {
		if _, exists := objType.operators[oc.operator]; exists {
			return true
		}
	}
	return false
}

func (oc *OperatorConstraint) String() string {
	return fmt.Sprintf("supports %s", oc.operator)
}

type CompositeConstraint struct {
	constraints []Constraint
}

func NewCompositeConstraint(constraints ...Constraint) *CompositeConstraint {
	return &CompositeConstraint{
		constraints: constraints,
	}
}

func (cc *CompositeConstraint) SatisfiedBy(t Type) bool {
	for _, constraint := range cc.constraints {
		if !constraint.SatisfiedBy(t) {
			return false
		}
	}
	return true
}

func (cc *CompositeConstraint) String() string {
	parts := make([]string, len(cc.constraints))
	for i, c := range cc.constraints {
		parts[i] = c.String()
	}
	return strings.Join(parts, " AND ")
}

type UnionConstraint struct {
	constraints []Constraint
}

func NewUnionConstraint(constraints ...Constraint) *UnionConstraint {
	return &UnionConstraint{
		constraints: constraints,
	}
}

func (uc *UnionConstraint) SatisfiedBy(t Type) bool {
	for _, constraint := range uc.constraints {
		if constraint.SatisfiedBy(t) {
			return true
		}
	}
	return false
}

func (uc *UnionConstraint) String() string {
	parts := make([]string, len(uc.constraints))
	for i, c := range uc.constraints {
		parts[i] = c.String()
	}
	return strings.Join(parts, " OR ")
}

// 预定义的约束
var (
	// 数值类型约束
	NumericConstraint = NewInterfaceConstraint(numberType)

	// 可比较类型约束
	ComparableConstraint = NewCompositeConstraint(
		NewOperatorConstraint(OpEqual, anyType, boolType),
		NewOperatorConstraint(OpLess, anyType, boolType),
		NewOperatorConstraint(OpGreater, anyType, boolType),
	)

	// 可序列化类型约束
	SerializableConstraint = NewInterfaceConstraint(anyType) // 简化处理

	// 字符串类型约束
	StringConstraint = NewInterfaceConstraint(stringType)

	// 布尔类型约束
	BoolConstraint = NewInterfaceConstraint(boolType)
)

// 便捷的约束创建函数
func NewTypeConstraint(targetType Type) *InterfaceConstraint {
	return NewInterfaceConstraint(targetType)
}

func NewOperatorConstraintSimple(op Operator) *OperatorConstraint {
	return NewOperatorConstraint(op, anyType, anyType)
}

// 从Type创建Constraint的辅助函数
func TypeToConstraint(t Type) Constraint {
	return NewInterfaceConstraint(t)
}

// 从多个Type创建CompositeConstraint的辅助函数
func TypesToCompositeConstraint(types ...Type) *CompositeConstraint {
	constraints := make([]Constraint, len(types))
	for i, t := range types {
		constraints[i] = TypeToConstraint(t)
	}
	return NewCompositeConstraint(constraints...)
}

// 从多个Type创建UnionConstraint的辅助函数
func TypesToUnionConstraint(types ...Type) *UnionConstraint {
	constraints := make([]Constraint, len(types))
	for i, t := range types {
		constraints[i] = TypeToConstraint(t)
	}
	return NewUnionConstraint(constraints...)
}

// 增强的泛型类，支持Constraint
type ConstrainedGenericClass struct {
	*ClassType
	constraintMap map[string]Constraint // 类型参数约束
}

func NewConstrainedGenericClass(name string, typeParams []string, constraints map[string]Constraint) *ConstrainedGenericClass {
	return &ConstrainedGenericClass{
		ClassType:     NewGenericClass(name, typeParams, constraints),
		constraintMap: constraints,
	}
}

func (cgc *ConstrainedGenericClass) Instantiate(typeArgs []Type) (*ClassType, error) {
	// 检查约束
	for i, param := range cgc.typeParams {
		if constraint, exists := cgc.constraintMap[param]; exists {
			if !constraint.SatisfiedBy(typeArgs[i]) {
				return nil, fmt.Errorf("type %s does not satisfy constraint %s for parameter %s",
					typeArgs[i].Name(), constraint.String(), param)
			}
		}
	}

	// 调用父类的实例化方法
	return cgc.ClassType.Instantiate(typeArgs)
}

// 示例：使用Constraint的泛型类
func ExampleConstrainedGeneric() {
	// 定义一个要求数值类型的泛型类
	numberListType := NewConstrainedGenericClass("NumberList", []string{"T"}, map[string]Constraint{
		"T": NumericConstraint,
	})

	// 添加方法
	numberListType.AddMethod("sum", &FunctionType{
		name: "sum",
		signatures: []*FunctionSignature{
			{
				positionalParams: []*Parameter{
					{name: "list", typ: numberListType},
				},
				returnType: numberType,
				builtin: func(args ...Value) Value {
					// 实现求和逻辑
					return &NumberValue{Value: big.NewFloat(0)}
				},
			},
		},
	})

	// 实例化 - 成功
	_, err := numberListType.Instantiate([]Type{numberType})
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	// 实例化 - 失败（string不是数值类型）
	_, err = numberListType.Instantiate([]Type{stringType})
	if err != nil {
		fmt.Printf("Expected error: %v\n", err)
	}
}

// 示例：复合约束
func ExampleCompositeConstraint() {
	// 定义一个要求可比较且可序列化的泛型类
	comparableSerializableType := NewConstrainedGenericClass("ComparableSerializable", []string{"T"}, map[string]Constraint{
		"T": NewCompositeConstraint(
			ComparableConstraint,
			SerializableConstraint,
		),
	})

	// 这个类型要求T既支持比较操作，又支持序列化
	_ = comparableSerializableType
}

// 示例：联合约束
func ExampleUnionConstraint() {
	// 定义一个接受数值或字符串的泛型类
	flexibleType := NewConstrainedGenericClass("Flexible", []string{"T"}, map[string]Constraint{
		"T": NewUnionConstraint(
			NumericConstraint,
			StringConstraint,
		),
	})

	// 这个类型接受数值类型或字符串类型
	_ = flexibleType
}
