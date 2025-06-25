// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package object

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVectorClass(t *testing.T) {
	vectorType := NewClassType("Vector")

	vectorType.AddField("x", &Field{ft: numberType})
	vectorType.AddField("y", &Field{ft: numberType})

	addOp := &FunctionType{
		name: "add",
		signatures: []*FunctionSignature{
			{
				positionalParams: []*Parameter{
					{name: "self", typ: vectorType},
					{name: "other", typ: vectorType},
				},
				returnType: vectorType,
				isOperator: true,
				operator:   OpAdd,
				builtin: func(args ...Value) Value {
					v1 := args[0].(*ClassInstance)
					v2 := args[1].(*ClassInstance)

					x1 := v1.GetField("x").(*NumberValue)
					y1 := v1.GetField("y").(*NumberValue)
					x2 := v2.GetField("x").(*NumberValue)
					y2 := v2.GetField("y").(*NumberValue)

					result := &ClassInstance{
						clsType: vectorType,
						fields: map[string]Value{
							"x": &NumberValue{Value: x1.Value.Add(x1.Value, x2.Value)},
							"y": &NumberValue{Value: y1.Value.Add(y1.Value, y2.Value)},
						},
					}

					return result
				},
			},
		},
	}

	vectorType.AddOperator(OpAdd, addOp)

	mulOp := &FunctionType{
		name: "mul",
		signatures: []*FunctionSignature{
			{
				positionalParams: []*Parameter{
					{name: "self", typ: vectorType},
					{name: "scalar", typ: numberType},
				},
				returnType: vectorType,
				isOperator: true,
				operator:   OpMul,
				builtin: func(args ...Value) Value {
					v := args[0].(*ClassInstance)
					scalar := args[1].(*NumberValue)

					x := v.GetField("x").(*NumberValue)
					y := v.GetField("y").(*NumberValue)

					result := &ClassInstance{
						clsType: vectorType,
						fields: map[string]Value{
							"x": &NumberValue{Value: x.Value.Mul(x.Value, scalar.Value)},
							"y": &NumberValue{Value: y.Value.Mul(y.Value, scalar.Value)},
						},
					}

					return result
				},
			},
		},
	}

	vectorType.AddOperator(OpMul, mulOp)

	v1 := &ClassInstance{
		clsType: vectorType,
		fields: map[string]Value{
			"x": NewNumberValue("1"),
			"y": NewNumberValue("2"),
		},
	}

	v2 := &ClassInstance{
		clsType: vectorType,
		fields: map[string]Value{
			"x": NewNumberValue("3"),
			"y": NewNumberValue("4"),
		},
	}
	assert.NotEqual(t, v1, v2)
	// 测试向量加法
	// result, err := ExecuteOperator(OpAdd, v1, v2)
	var result any

	resultVec := result.(*ClassInstance)
	x := resultVec.GetField("x").(*NumberValue)
	y := resultVec.GetField("y").(*NumberValue)

	assert.Equal(t, "4", x.Text())
	assert.Equal(t, "6", y.Text())

	// 测试向量 * 标量
	// scalar := NewNumberValue("2")
	//   result2, err := registry.ExecuteOperator(OpMul, v1, scalar)
	//   assert.NoError(t, err)

	var result2 any
	resultVec2 := result2.(*ClassInstance)
	x2 := resultVec2.GetField("x").(*NumberValue)
	y2 := resultVec2.GetField("y").(*NumberValue)

	assert.Equal(t, "2", x2.Text())
	assert.Equal(t, "4", y2.Text())
}

func TestGenericList(t *testing.T) {
	// List<T>
	listType := NewGenericClass("List", []string{"T"}, nil)

	listType.AddField("items", &Field{
		ft:   &GenericType{name: "T"},
		mods: []FieldModifier{FieldModifierRequired},
	})

	addMethod := NewGenericFunction("add", []string{"T"}, nil)
	addMethod.signatures = []*FunctionSignature{
		{
			positionalParams: []*Parameter{
				{name: "item", typ: &GenericType{name: "T"}},
			},
			returnType: voidType,
		},
	}
	listType.AddMethod("add", addMethod)

	// 实例化 List<number>
	numberListType, err := listType.Instantiate([]Type{numberType})
	assert.NoError(t, err)

	// 实例化 List<string>
	stringListType, err := listType.Instantiate([]Type{stringType})
	assert.NoError(t, err)

	// 验证类型不同
	assert.NotEqual(t, numberListType.Name(), stringListType.Name())

	numberField := numberListType.fields["items"]
	stringField := stringListType.fields["items"]

	assert.Equal(t, "number", numberField.ft.Name())
	assert.Equal(t, "string", stringField.ft.Name())
}
