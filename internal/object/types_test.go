// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package object

import (
	"fmt"
	"strings"
	"testing"

	"math/big"

	"github.com/stretchr/testify/assert"
)

func TestChainCall(t *testing.T) {
	// "abcdef".length()
	abcdef := &StringValue{Value: "abcdef"}

	if obj, ok := abcdef.Type().(Object); ok {
		result, err := obj.MethodByName("length").Call([]Value{abcdef}, nil)
		assert.NoError(t, err)
		assert.Equal(t, result.Type().Kind(), O_NUM)
	} else {
		t.Fatal("StringType should implement Object interface")
	}

	// str.length("abcdef")
	var strObj Object = stringType
	strObj.MethodByName("length").Call([]Value{abcdef}, nil)
	// $a := "abcdef"
	// $a.length()
}

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
	/*
		class List<T> {
			items: T[]

			fn add<T>(item: T) -> void
		}
	*/
	listType := NewGenericClass("List", []string{"T"}, nil)

	listType.AddField("items", &Field{
		ft: &ArrayType{
			ObjectType:  NewObjectType("T[]", O_ARR),
			elementType: NewTypeParameter("T", nil),
		},
		mods: []FieldModifier{FieldModifierRequired},
	})

	addMethod := NewGenericFunction("add", []string{"T"}, nil)
	addMethod.signatures = []*FunctionSignature{
		{
			positionalParams: []*Parameter{
				{name: "item", typ: NewTypeParameter("T", nil)},
			},
			returnType: voidType,
		},
	}
	listType.AddMethod("add", addMethod)

	// 实例化 List<num>
	numberListType, err := listType.Instantiate([]Type{numberType})
	assert.NoError(t, err)

	// 实例化 List<str>
	stringListType, err := listType.Instantiate([]Type{stringType})
	assert.NoError(t, err)

	// 验证类型不同
	assert.NotEqual(t, numberListType.Name(), stringListType.Name())

	numberField := numberListType.fields["items"]
	stringField := stringListType.fields["items"]

	// 检查字段类型是否正确实例化
	numberArrayType, ok := numberField.ft.(*ArrayType)
	assert.True(t, ok)
	assert.Equal(t, "num", numberArrayType.elementType.Name())

	stringArrayType, ok := stringField.ft.(*ArrayType)
	assert.True(t, ok)
	assert.Equal(t, "str", stringArrayType.elementType.Name())
}

func TestArrayType(t *testing.T) {
	// 测试创建数组类型
	intArrayType := NewArrayType(numberType)
	assert.Equal(t, "num[]", intArrayType.Name())
	assert.Equal(t, O_ARR, intArrayType.Kind())
	assert.Equal(t, numberType, intArrayType.ElementType())

	// 测试数组值
	intArray := NewArrayValue(numberType,
		&NumberValue{Value: big.NewFloat(1)},
		&NumberValue{Value: big.NewFloat(2)},
		&NumberValue{Value: big.NewFloat(3)},
	)
	assert.Equal(t, "[1, 2, 3]", intArray.Text())
	assert.Equal(t, 3, intArray.Length())

	// 测试数组方法
	lengthMethod := intArrayType.MethodByName("length")
	assert.NotNil(t, lengthMethod)

	// 测试length方法
	result, err := lengthMethod.Call([]Value{intArray}, nil)
	assert.NoError(t, err)
	assert.Equal(t, &NumberValue{Value: big.NewFloat(3)}, result)

	// 测试push方法
	pushMethod := intArrayType.MethodByName("push")
	assert.NotNil(t, pushMethod)

	result, err = pushMethod.Call([]Value{intArray, &NumberValue{Value: big.NewFloat(4)}}, nil)
	assert.NoError(t, err)
	assert.Equal(t, 4, intArray.Length())
	assert.Equal(t, "[1, 2, 3, 4]", intArray.Text())

	// 测试pop方法
	popMethod := intArrayType.MethodByName("pop")
	assert.NotNil(t, popMethod)

	result, err = popMethod.Call([]Value{intArray}, nil)
	assert.NoError(t, err)
	assert.Equal(t, &NumberValue{Value: big.NewFloat(4)}, result)
	assert.Equal(t, 3, intArray.Length())
}

func TestMapType(t *testing.T) {
	// 测试创建map类型
	stringIntMapType := NewMapType(stringType, numberType)
	assert.Equal(t, "map<str, num>", stringIntMapType.Name())
	assert.Equal(t, O_MAP, stringIntMapType.Kind())
	assert.Equal(t, stringType, stringIntMapType.KeyType())
	assert.Equal(t, numberType, stringIntMapType.ValueType())

	// 测试map值
	stringIntMap := NewMapValue(stringType, numberType)
	stringIntMap.Set(&StringValue{Value: "a"}, &NumberValue{Value: big.NewFloat(1)})
	stringIntMap.Set(&StringValue{Value: "b"}, &NumberValue{Value: big.NewFloat(2)})

	assert.Equal(t, "{a: 1, b: 2}", stringIntMap.Text())
	assert.Equal(t, 2, stringIntMap.Length())

	// 测试map方法
	lengthMethod := stringIntMapType.MethodByName("length")
	assert.NotNil(t, lengthMethod)

	// 测试length方法
	result, err := lengthMethod.Call([]Value{stringIntMap}, nil)
	assert.NoError(t, err)
	assert.Equal(t, &NumberValue{Value: big.NewFloat(2)}, result)

	// 测试has方法
	hasMethod := stringIntMapType.MethodByName("has")
	assert.NotNil(t, hasMethod)

	result, err = hasMethod.Call([]Value{stringIntMap, &StringValue{Value: "a"}}, nil)
	assert.NoError(t, err)
	assert.Equal(t, &BoolValue{Value: true}, result)

	result, err = hasMethod.Call([]Value{stringIntMap, &StringValue{Value: "c"}}, nil)
	assert.NoError(t, err)
	assert.Equal(t, &BoolValue{Value: false}, result)

	// 测试get方法
	getMethod := stringIntMapType.MethodByName("get")
	assert.NotNil(t, getMethod)

	result, err = getMethod.Call([]Value{stringIntMap, &StringValue{Value: "a"}}, nil)
	assert.NoError(t, err)
	assert.Equal(t, &NumberValue{Value: big.NewFloat(1)}, result)

	result, err = getMethod.Call([]Value{stringIntMap, &StringValue{Value: "c"}}, nil)
	assert.NoError(t, err)
	assert.Equal(t, NULL, result)

	// 测试insert方法
	insertMethod := stringIntMapType.MethodByName("insert")
	assert.NotNil(t, insertMethod)

	result, err = insertMethod.Call([]Value{stringIntMap, &StringValue{Value: "c"}, &NumberValue{Value: big.NewFloat(3)}}, nil)
	assert.NoError(t, err)
	assert.Equal(t, NULL, result) // 新插入的key返回null
	assert.Equal(t, 3, stringIntMap.Length())

	// 测试更新已存在的key
	result, err = insertMethod.Call([]Value{stringIntMap, &StringValue{Value: "a"}, &NumberValue{Value: big.NewFloat(10)}}, nil)
	assert.NoError(t, err)
	assert.Equal(t, &NumberValue{Value: big.NewFloat(1)}, result) // 返回旧值
	assert.Equal(t, "{a: 10, b: 2, c: 3}", stringIntMap.Text())
}

func TestArraySlice(t *testing.T) {
	intArray := NewArrayValue(numberType)
	for i := 1; i <= 5; i++ {
		intArray.Append(&NumberValue{Value: big.NewFloat(float64(i))})
	}

	intArrayType := NewArrayType(numberType)
	sliceMethod := intArrayType.MethodByName("slice")
	assert.NotNil(t, sliceMethod)

	// 测试基本切片
	result, err := sliceMethod.Call([]Value{intArray, &NumberValue{Value: big.NewFloat(1)}, &NumberValue{Value: big.NewFloat(3)}}, nil)
	assert.NoError(t, err)
	slicedArray := result.(*ArrayValue)
	assert.Equal(t, "[2, 3]", slicedArray.Text())

	// 测试负索引
	result, err = sliceMethod.Call([]Value{intArray, &NumberValue{Value: big.NewFloat(-3)}, &NumberValue{Value: big.NewFloat(-1)}}, nil)
	assert.NoError(t, err)
	slicedArray = result.(*ArrayValue)
	assert.Equal(t, "[3, 4]", slicedArray.Text())
}

func TestArrayJoin(t *testing.T) {
	strArray := NewArrayValue(stringType,
		&StringValue{Value: "hello"},
		&StringValue{Value: "world"},
		&StringValue{Value: "test"},
	)

	strArrayType := NewArrayType(stringType)
	joinMethod := strArrayType.MethodByName("join")
	assert.NotNil(t, joinMethod)

	// 测试默认分隔符
	result, err := joinMethod.Call([]Value{strArray}, nil)
	assert.NoError(t, err)
	assert.Equal(t, &StringValue{Value: "hello,world,test"}, result)

	// 测试自定义分隔符
	result, err = joinMethod.Call([]Value{strArray, &StringValue{Value: " - "}}, nil)
	assert.NoError(t, err)
	assert.Equal(t, &StringValue{Value: "hello - world - test"}, result)
}

func TestTraitSystem(t *testing.T) {
	// 创建Service trait
	serviceTrait := NewTraitType("Service")

	// 添加字段
	serviceTrait.AddField("name", &Field{
		ft:   stringType,
		mods: []FieldModifier{FieldModifierReadonly},
	})
	serviceTrait.AddField("port", &Field{
		ft:   numberType,
		mods: []FieldModifier{FieldModifierReadonly},
	})

	// 添加方法
	serviceTrait.AddMethod("run", &FunctionType{
		name: "run",
		signatures: []*FunctionSignature{
			{
				positionalParams: []*Parameter{
					{name: "self", typ: serviceTrait},
				},
				returnType: voidType,
			},
		},
	})

	// 注册trait
	GlobalTraitRegistry.RegisterTrait(serviceTrait)

	// 创建GameService类
	gameServiceType := NewClassType("GameService")

	// 为GameService实现Service trait
	impl := NewTraitImpl(serviceTrait, gameServiceType)

	// 实现字段
	impl.AddField("name", &Field{
		ft:   stringType,
		mods: []FieldModifier{FieldModifierReadonly},
	})
	impl.AddField("port", &Field{
		ft:   numberType,
		mods: []FieldModifier{FieldModifierReadonly},
	})

	// 实现方法
	impl.AddMethod("run", &FunctionType{
		name: "run",
		signatures: []*FunctionSignature{
			{
				positionalParams: []*Parameter{
					{name: "self", typ: gameServiceType},
				},
				returnType: voidType,
				builtin: func(args ...Value) Value {
					return NULL
				},
			},
		},
	})

	// 注册实现
	GlobalTraitRegistry.RegisterImpl(impl)

	// 测试trait实现检查
	assert.True(t, ImplementsTrait(gameServiceType, "Service"))
	assert.False(t, ImplementsTrait(stringType, "Service"))

	// 测试获取trait
	retrievedTrait := GlobalTraitRegistry.GetTrait("Service")
	assert.NotNil(t, retrievedTrait)
	assert.Equal(t, "Service", retrievedTrait.Name())

	// 测试获取实现
	impls := GlobalTraitRegistry.GetImpl("Service", "GameService")
	assert.Len(t, impls, 1)
	assert.Equal(t, gameServiceType, impls[0].implType)
}

func TestTupleType(t *testing.T) {
	// 测试匿名元组
	personType := NewTupleTypeWithMethods(stringType, numberType)
	assert.Equal(t, "[str, num]", personType.Name())
	assert.Equal(t, 2, personType.Length())

	person := NewTupleValue(personType,
		&StringValue{Value: "Alice"},
		&NumberValue{Value: big.NewFloat(30)})

	assert.Equal(t, "[Alice, 30]", person.Text())
	assert.Equal(t, "Alice", person.Get(0).Text())
	assert.Equal(t, "30", person.Get(1).Text())

	// 测试命名元组
	pointType := NewNamedTupleTypeWithMethods(
		[]string{"x", "y"},
		[]Type{numberType, numberType})
	assert.Equal(t, "(x: num, y: num)", pointType.Name())

	point := NewTupleValue(pointType,
		&NumberValue{Value: big.NewFloat(10)},
		&NumberValue{Value: big.NewFloat(20)})

	assert.Equal(t, "(x: 10, y: 20)", point.Text())
	assert.Equal(t, "10", point.GetByName("x").Text())
	assert.Equal(t, "20", point.GetByName("y").Text())
	assert.Equal(t, -1, pointType.GetFieldIndex("z")) // 不存在的字段

	// 测试更新元素
	newPerson, err := person.With(0, &StringValue{Value: "Bob"})
	assert.NoError(t, err)
	assert.Equal(t, "[Bob, 30]", newPerson.Text())
	assert.Equal(t, "[Alice, 30]", person.Text()) // 原元组不变

	// 测试类型错误
	_, err = person.With(0, &NumberValue{Value: big.NewFloat(100)})
	assert.Error(t, err) // 类型不匹配

	// 测试索引越界
	_, err = person.With(10, &StringValue{Value: "test"})
	assert.Error(t, err)

	// 测试解构
	elements := person.Destructure()
	assert.Len(t, elements, 2)
	assert.Equal(t, "Alice", elements[0].Text())
	assert.Equal(t, "30", elements[1].Text())

	// 测试元组方法
	length := personType.MethodByName("length")
	assert.NotNil(t, length)

	get := personType.MethodByName("get")
	assert.NotNil(t, get)

	with := personType.MethodByName("with")
	assert.NotNil(t, with)

	destructure := personType.MethodByName("destructure")
	assert.NotNil(t, destructure)
}

func TestTupleTypeCompatibility(t *testing.T) {
	// 测试类型兼容性
	tuple1 := NewTupleType(stringType, numberType)
	tuple2 := NewTupleType(stringType, numberType)
	tuple3 := NewTupleType(stringType, stringType)

	assert.True(t, tuple1.AssignableTo(tuple2))
	assert.False(t, tuple1.AssignableTo(tuple3))

	// 测试长度不同的元组
	tuple4 := NewTupleType(stringType, numberType, boolType)
	assert.False(t, tuple1.AssignableTo(tuple4))
}

// ==================== 集合类型系统测试 ====================

func TestSetType(t *testing.T) {
	// 测试集合类型创建
	numberSetType := NewSetType(numberType)
	assert.Equal(t, "set<num>", numberSetType.Name())
	assert.Equal(t, O_SET, numberSetType.Kind())
	assert.Equal(t, numberType, numberSetType.ElementType())

	stringSetType := NewSetType(stringType)
	assert.Equal(t, "set<str>", stringSetType.Name())
	assert.Equal(t, O_SET, stringSetType.Kind())
	assert.Equal(t, stringType, stringSetType.ElementType())
}

func TestSetValue(t *testing.T) {
	// 测试集合值创建
	numberSetType := NewSetType(numberType)
	ids := NewSetValue(numberSetType,
		&NumberValue{Value: big.NewFloat(1)},
		&NumberValue{Value: big.NewFloat(2)},
		&NumberValue{Value: big.NewFloat(3)})

	assert.Equal(t, numberSetType, ids.Type())
	assert.Equal(t, 3, ids.Length())
	assert.Equal(t, "{1, 2, 3}", ids.Text())
	assert.False(t, ids.IsEmpty())

	// 测试空集合
	emptySet := NewSetValue(numberSetType)
	assert.Equal(t, 0, emptySet.Length())
	assert.Equal(t, "{}", emptySet.Text())
	assert.True(t, emptySet.IsEmpty())
}

func TestSetSetOperations(t *testing.T) {
	numberSetType := NewSetType(numberType)

	setA := NewSetValue(numberSetType,
		&NumberValue{Value: big.NewFloat(1)},
		&NumberValue{Value: big.NewFloat(2)},
		&NumberValue{Value: big.NewFloat(3)})

	setB := NewSetValue(numberSetType,
		&NumberValue{Value: big.NewFloat(3)},
		&NumberValue{Value: big.NewFloat(4)},
		&NumberValue{Value: big.NewFloat(5)})

	// 测试并集
	union, err := setA.Union(setB)
	assert.NoError(t, err)
	assert.Equal(t, "{1, 2, 3, 4, 5}", union.Text())

	// 测试交集
	intersect, err := setA.Intersect(setB)
	assert.NoError(t, err)
	assert.Equal(t, "{3}", intersect.Text())

	// 测试差集
	difference, err := setA.Difference(setB)
	assert.NoError(t, err)
	assert.Equal(t, "{1, 2}", difference.Text())
}

func TestSetContains(t *testing.T) {
	numberSetType := NewSetType(numberType)
	ids := NewSetValue(numberSetType,
		&NumberValue{Value: big.NewFloat(1)},
		&NumberValue{Value: big.NewFloat(2)},
		&NumberValue{Value: big.NewFloat(3)})

	// 测试存在检查
	has2, err := ids.Contains(&NumberValue{Value: big.NewFloat(2)})
	assert.NoError(t, err)
	assert.True(t, has2)

	has5, err := ids.Contains(&NumberValue{Value: big.NewFloat(5)})
	assert.NoError(t, err)
	assert.False(t, has5)

	// 测试类型错误
	stringSetType := NewSetType(stringType)
	flags := NewSetValue(stringSetType,
		&StringValue{Value: "read"},
		&StringValue{Value: "write"})

	// 尝试在字符串集合中检查数字
	_, err = flags.Contains(&NumberValue{Value: big.NewFloat(1)})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot check num in set<str>")
}

func TestSetTypeCompatibility(t *testing.T) {
	numberSetType := NewSetType(numberType)
	stringSetType := NewSetType(stringType)

	// 测试类型兼容性
	assert.True(t, numberSetType.AssignableTo(numberSetType))
	assert.False(t, numberSetType.AssignableTo(stringSetType))
	assert.False(t, numberSetType.ConvertibleTo(stringSetType))
	assert.False(t, numberSetType.Implements(stringSetType))
}

func TestSetMethods(t *testing.T) {
	numberSetType := NewSetType(numberType)
	ids := NewSetValue(numberSetType,
		&NumberValue{Value: big.NewFloat(1)},
		&NumberValue{Value: big.NewFloat(2)},
		&NumberValue{Value: big.NewFloat(3)})

	// 测试 length 方法
	lengthMethod := numberSetType.MethodByName("length")
	assert.NotNil(t, lengthMethod)

	result, err := lengthMethod.Call([]Value{ids}, nil)
	assert.NoError(t, err)
	assert.Equal(t, "3", result.Text())

	// 测试 add 方法
	addMethod := numberSetType.MethodByName("add")
	assert.NotNil(t, addMethod)

	result, err = addMethod.Call([]Value{ids, &NumberValue{Value: big.NewFloat(4)}}, nil)
	assert.NoError(t, err)
	assert.Equal(t, "{1, 2, 3, 4}", result.Text())

	// 测试 has 方法
	hasMethod := numberSetType.MethodByName("has")
	assert.NotNil(t, hasMethod)

	result, err = hasMethod.Call([]Value{ids, &NumberValue{Value: big.NewFloat(2)}}, nil)
	assert.NoError(t, err)
	assert.Equal(t, "true", result.Text())

	result, err = hasMethod.Call([]Value{ids, &NumberValue{Value: big.NewFloat(5)}}, nil)
	assert.NoError(t, err)
	assert.Equal(t, "false", result.Text())
}

func TestSetToArray(t *testing.T) {
	numberSetType := NewSetType(numberType)
	ids := NewSetValue(numberSetType,
		&NumberValue{Value: big.NewFloat(1)},
		&NumberValue{Value: big.NewFloat(2)},
		&NumberValue{Value: big.NewFloat(3)})

	// 测试转换为数组
	array := ids.ToArray()
	assert.Equal(t, 3, len(array.Elements))
	assert.Equal(t, numberType, array.arrayType.ElementType())

	// 验证数组元素
	elementTexts := make([]string, len(array.Elements))
	for i, elem := range array.Elements {
		elementTexts[i] = elem.Text()
	}
	assert.Contains(t, elementTexts, "1")
	assert.Contains(t, elementTexts, "2")
	assert.Contains(t, elementTexts, "3")
}

func TestSetClear(t *testing.T) {
	numberSetType := NewSetType(numberType)
	ids := NewSetValue(numberSetType,
		&NumberValue{Value: big.NewFloat(1)},
		&NumberValue{Value: big.NewFloat(2)},
		&NumberValue{Value: big.NewFloat(3)})

	assert.Equal(t, 3, ids.Length())

	// 测试清空
	ids.Clear()
	assert.Equal(t, 0, ids.Length())
	assert.Equal(t, "{}", ids.Text())
	assert.True(t, ids.IsEmpty())
}

func TestSetTypeSafety(t *testing.T) {
	numberSetType := NewSetType(numberType)
	ids := NewSetValue(numberSetType,
		&NumberValue{Value: big.NewFloat(1)},
		&NumberValue{Value: big.NewFloat(2)},
		&NumberValue{Value: big.NewFloat(3)})

	// 测试类型安全 - 尝试添加错误类型的元素
	err := ids.Add(&StringValue{Value: "hello"})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot add str to set<num>")

	// 测试类型安全 - 尝试删除错误类型的元素
	err = ids.Remove(&StringValue{Value: "hello"})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "element hello not found in set")
}

func TestSetOperationsWithDifferentTypes(t *testing.T) {
	numberSetType := NewSetType(numberType)
	stringSetType := NewSetType(stringType)

	setA := NewSetValue(numberSetType,
		&NumberValue{Value: big.NewFloat(1)},
		&NumberValue{Value: big.NewFloat(2)})

	setB := NewSetValue(stringSetType,
		&StringValue{Value: "hello"},
		&StringValue{Value: "world"})

	// 测试不同类型集合的运算
	_, err := setA.Union(setB)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot union sets with incompatible element types")

	_, err = setA.Intersect(setB)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot intersect sets with incompatible element types")

	_, err = setA.Difference(setB)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot compute difference of sets with incompatible element types")
}

// ==================== 枚举类型系统测试 ====================

func TestNumericEnum(t *testing.T) {
	// 测试标准数字枚举
	statusEnum := NewNumericEnum("Status", []string{"Pending", "Approved", "Rejected"})

	assert.Equal(t, "Status", statusEnum.Name())
	assert.Equal(t, O_ENUM, statusEnum.Kind())
	assert.Equal(t, 3, len(statusEnum.GetValues()))

	// 测试枚举值
	pending := statusEnum.GetValue("Pending")
	assert.NotNil(t, pending)
	assert.Equal(t, "Pending", pending.Name())
	assert.Equal(t, "0", pending.value.Text())

	approved := statusEnum.GetValue("Approved")
	assert.NotNil(t, approved)
	assert.Equal(t, "Approved", approved.Name())
	assert.Equal(t, "1", approved.value.Text())

	rejected := statusEnum.GetValue("Rejected")
	assert.NotNil(t, rejected)
	assert.Equal(t, "Rejected", rejected.Name())
	assert.Equal(t, "2", rejected.value.Text())
}

func TestExplicitNumericEnum(t *testing.T) {
	// 测试显式数字枚举
	httpCodeEnum := NewExplicitNumericEnum("HttpCode", map[string]int{
		"OK":             200,
		"NotFound":       404,
		"ServerError":    500,
		"GatewayTimeout": 501,
	})

	assert.Equal(t, "HttpCode", httpCodeEnum.Name())
	assert.Equal(t, 4, len(httpCodeEnum.GetValues()))

	// 测试枚举值
	ok := httpCodeEnum.GetValue("OK")
	assert.NotNil(t, ok)
	assert.Equal(t, "200", ok.value.Text())

	notFound := httpCodeEnum.GetValue("NotFound")
	assert.NotNil(t, notFound)
	assert.Equal(t, "404", notFound.value.Text())

	serverError := httpCodeEnum.GetValue("ServerError")
	assert.NotNil(t, serverError)
	assert.Equal(t, "500", serverError.value.Text())

	gatewayTimeout := httpCodeEnum.GetValue("GatewayTimeout")
	assert.NotNil(t, gatewayTimeout)
	assert.Equal(t, "501", gatewayTimeout.value.Text())
}

func TestStringEnum(t *testing.T) {
	// 测试字符串枚举
	directionEnum := NewStringEnum("Direction", map[string]string{
		"North": "N",
		"South": "S",
		"East":  "E",
		"West":  "W",
	})

	assert.Equal(t, "Direction", directionEnum.Name())
	assert.Equal(t, 4, len(directionEnum.GetValues()))

	// 测试枚举值
	north := directionEnum.GetValue("North")
	assert.NotNil(t, north)
	assert.Equal(t, "N", north.value.Text())

	south := directionEnum.GetValue("South")
	assert.NotNil(t, south)
	assert.Equal(t, "S", south.value.Text())

	east := directionEnum.GetValue("East")
	assert.NotNil(t, east)
	assert.Equal(t, "E", east.value.Text())

	west := directionEnum.GetValue("West")
	assert.NotNil(t, west)
	assert.Equal(t, "W", west.value.Text())
}

func TestMixedEnum(t *testing.T) {
	// 测试混合类型枚举
	configEnum := NewMixedEnum("Config", map[string]Value{
		"RetryCount":    &NumberValue{Value: big.NewFloat(3)},
		"Timeout":       &StringValue{Value: "30s"},
		"EnableLogging": &BoolValue{Value: true},
	})

	assert.Equal(t, "Config", configEnum.Name())
	assert.Equal(t, 3, len(configEnum.GetValues()))

	// 测试枚举值
	retryCount := configEnum.GetValue("RetryCount")
	assert.NotNil(t, retryCount)
	assert.Equal(t, "3", retryCount.value.Text())

	timeout := configEnum.GetValue("Timeout")
	assert.NotNil(t, timeout)
	assert.Equal(t, "30s", timeout.value.Text())

	enableLogging := configEnum.GetValue("EnableLogging")
	assert.NotNil(t, enableLogging)
	assert.Equal(t, "true", enableLogging.value.Text())
}

func TestAssociatedEnum(t *testing.T) {
	// 测试关联值枚举
	protocolEnum := NewAssociatedEnum("Protocol", map[string]*Field{
		"port": {ft: numberType, mods: []FieldModifier{FieldModifierRequired}},
	}, map[string]Value{
		"tcp": &NumberValue{Value: big.NewFloat(6)},
		"udp": &NumberValue{Value: big.NewFloat(17)},
	})

	assert.Equal(t, "Protocol", protocolEnum.Name())
	assert.Equal(t, 2, len(protocolEnum.GetValues()))

	// 测试枚举值
	tcp := protocolEnum.GetValue("tcp")
	assert.NotNil(t, tcp)
	assert.Equal(t, "6", tcp.value.Text())

	// 测试关联字段
	tcp.SetField("port", &NumberValue{Value: big.NewFloat(6)})
	assert.Equal(t, "6", tcp.GetField("port").Text())

	udp := protocolEnum.GetValue("udp")
	assert.NotNil(t, udp)
	assert.Equal(t, "17", udp.value.Text())

	udp.SetField("port", &NumberValue{Value: big.NewFloat(17)})
	assert.Equal(t, "17", udp.GetField("port").Text())
}

func TestADTEnum(t *testing.T) {
	// 测试代数数据类型枚举
	networkPacketEnum := NewADTEnum("NetworkPacket", map[string]map[string]Type{
		"TCP": {
			"src_port": numberType,
			"dst_port": numberType,
		},
		"UDP": {
			"port":    numberType,
			"payload": stringType,
		},
	})

	assert.Equal(t, "NetworkPacket", networkPacketEnum.Name())
	assert.Equal(t, 2, len(networkPacketEnum.GetValues()))
	assert.True(t, networkPacketEnum.IsADT())

	// 测试枚举值
	tcpPacket := networkPacketEnum.GetValue("TCP")
	assert.NotNil(t, tcpPacket)
	assert.Equal(t, "TCP", tcpPacket.value.Text())

	// 测试ADT字段
	tcpPacket.SetADTField("src_port", &NumberValue{Value: big.NewFloat(8080)})
	tcpPacket.SetADTField("dst_port", &NumberValue{Value: big.NewFloat(80)})

	assert.Equal(t, "8080", tcpPacket.GetADTField("src_port").Text())
	assert.Equal(t, "80", tcpPacket.GetADTField("dst_port").Text())
	assert.Equal(t, "TCP{src_port: 8080, dst_port: 80}", tcpPacket.Text())

	udpPacket := networkPacketEnum.GetValue("UDP")
	assert.NotNil(t, udpPacket)
	assert.Equal(t, "UDP", udpPacket.value.Text())

	udpPacket.SetADTField("port", &NumberValue{Value: big.NewFloat(53)})
	udpPacket.SetADTField("payload", &StringValue{Value: "DNS query"})

	assert.Equal(t, "53", udpPacket.GetADTField("port").Text())
	assert.Equal(t, "DNS query", udpPacket.GetADTField("payload").Text())
	assert.Equal(t, "UDP{port: 53, payload: DNS query}", udpPacket.Text())
}

func TestEnumMethods(t *testing.T) {
	// 测试枚举方法
	statusEnum := NewNumericEnum("Status", []string{"Pending", "Approved", "Rejected"})
	pending := statusEnum.GetValue("Pending")

	// 测试 values 方法
	valuesMethod := statusEnum.MethodByName("values")
	assert.NotNil(t, valuesMethod)

	result, err := valuesMethod.Call([]Value{pending}, nil)
	assert.NoError(t, err)
	assert.NotNil(t, result)

	// 测试 from_value 方法
	fromValueMethod := statusEnum.MethodByName("from_value")
	assert.NotNil(t, fromValueMethod)

	result, err = fromValueMethod.Call([]Value{pending, &NumberValue{Value: big.NewFloat(1)}}, nil)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Approved", result.(*EnumVariant).Name())

	// 测试 is_valid 方法
	isValidMethod := statusEnum.MethodByName("is_valid")
	assert.NotNil(t, isValidMethod)

	result, err = isValidMethod.Call([]Value{pending, &NumberValue{Value: big.NewFloat(1)}}, nil)
	assert.NoError(t, err)
	assert.Equal(t, "true", result.Text())

	result, err = isValidMethod.Call([]Value{pending, &NumberValue{Value: big.NewFloat(99)}}, nil)
	assert.NoError(t, err)
	assert.Equal(t, "false", result.Text())
}

func TestEnumTypeCompatibility(t *testing.T) {
	// 测试枚举类型兼容性
	statusEnum1 := NewNumericEnum("Status", []string{"Pending", "Approved"})
	statusEnum2 := NewNumericEnum("Status", []string{"Pending", "Approved"})
	directionEnum := NewStringEnum("Direction", map[string]string{"North": "N"})

	// 相同名称的枚举应该兼容
	assert.True(t, statusEnum1.AssignableTo(statusEnum2))

	// 不同名称的枚举不兼容
	assert.False(t, statusEnum1.AssignableTo(directionEnum))
	assert.False(t, statusEnum1.ConvertibleTo(directionEnum))
	assert.False(t, statusEnum1.Implements(directionEnum))
}

func TestEnumValueOperations(t *testing.T) {
	// 测试枚举值操作
	statusEnum := NewNumericEnum("Status", []string{"Pending", "Approved", "Rejected"})
	pending := statusEnum.GetValue("Pending")

	// 测试基本属性
	assert.Equal(t, "Pending", pending.Name())
	assert.Equal(t, "0", pending.Value().Text())
	assert.Equal(t, statusEnum, pending.Type())
	assert.Equal(t, "Pending", pending.Text())

	// 测试接口方法
	interfaceValue := pending.Interface()
	assert.NotNil(t, interfaceValue)

	// 测试字段操作
	pending.SetField("custom_field", &StringValue{Value: "custom_value"})
	assert.Equal(t, "custom_value", pending.GetField("custom_field").Text())

	// 测试不存在的字段
	assert.Equal(t, NULL, pending.GetField("non_existent"))
}

func TestEnumValueADTOperations(t *testing.T) {
	// 测试ADT枚举值操作
	networkPacketEnum := NewADTEnum("NetworkPacket", map[string]map[string]Type{
		"TCP": {"src_port": numberType, "dst_port": numberType},
	})

	tcpPacket := networkPacketEnum.GetValue("TCP")

	// 测试ADT字段操作
	tcpPacket.SetADTField("src_port", &NumberValue{Value: big.NewFloat(8080)})
	tcpPacket.SetADTField("dst_port", &NumberValue{Value: big.NewFloat(80)})

	assert.Equal(t, "8080", tcpPacket.GetADTField("src_port").Text())
	assert.Equal(t, "80", tcpPacket.GetADTField("dst_port").Text())

	// 测试不存在的ADT字段
	assert.Equal(t, NULL, tcpPacket.GetADTField("non_existent"))

	// 测试ADT字符串表示
	assert.Equal(t, "TCP{src_port: 8080, dst_port: 80}", tcpPacket.Text())

	// 测试空ADT字段的字符串表示
	// 创建一个新的枚举实例来确保变体是空的
	emptyNetworkEnum := NewADTEnum("EmptyNetwork", map[string]map[string]Type{
		"TCP": {"src_port": numberType, "dst_port": numberType},
	})
	emptyTCPPacket := emptyNetworkEnum.GetValue("TCP")
	assert.Equal(t, "TCP", emptyTCPPacket.Text())
}

// 测试枚举边界情况和错误处理
func TestEnumEdgeCases(t *testing.T) {
	// 测试空枚举
	emptyEnum := NewEnumType("Empty")
	assert.Equal(t, 0, len(emptyEnum.GetVariants()))
	assert.Nil(t, emptyEnum.GetVariant("NonExistent"))

	// 测试重复变体名
	duplicateEnum := NewEnumType("Duplicate")
	variant1 := NewEnumVariant("Same", &NumberValue{Value: big.NewFloat(1)})
	variant2 := NewEnumVariant("Same", &NumberValue{Value: big.NewFloat(2)})

	duplicateEnum.AddVariant(variant1)
	duplicateEnum.AddVariant(variant2) // 应该覆盖前一个

	result := duplicateEnum.GetVariant("Same")
	assert.NotNil(t, result)
	assert.Equal(t, "2", result.Value().Text()) // 应该是后添加的值
}

// 测试枚举类型安全性
func TestEnumTypeSafety(t *testing.T) {
	// 测试类型不匹配的字段设置
	statusEnum := NewNumericEnum("Status", []string{"Pending", "Approved"})
	pending := statusEnum.GetValue("Pending")

	// 设置错误类型的字段应该被允许（运行时检查）
	pending.SetField("wrong_type", &StringValue{Value: "should_be_number"})
	assert.Equal(t, "should_be_number", pending.GetField("wrong_type").Text())

	// 测试ADT字段类型检查
	networkEnum := NewADTEnum("Network", map[string]map[string]Type{
		"TCP": {"port": numberType},
	})
	tcpVariant := networkEnum.GetVariant("TCP")

	// 设置正确类型的ADT字段
	tcpVariant.SetADTField("port", &NumberValue{Value: big.NewFloat(80)})
	assert.Equal(t, "80", tcpVariant.GetADTField("port").Text())

	// 设置错误类型的ADT字段（应该被允许，但类型不匹配）
	tcpVariant.SetADTField("port", &StringValue{Value: "should_be_number"})
	assert.Equal(t, "should_be_number", tcpVariant.GetADTField("port").Text())
}

// 测试枚举序列化和反序列化
func TestEnumSerialization(t *testing.T) {
	// 测试枚举值的文本表示
	statusEnum := NewNumericEnum("Status", []string{"Pending", "Approved", "Rejected"})

	pending := statusEnum.GetValue("Pending")
	assert.Equal(t, "Pending", pending.Text())

	approved := statusEnum.GetValue("Approved")
	assert.Equal(t, "Approved", approved.Text())

	// 测试关联值枚举的文本表示
	protocolEnum := NewAssociatedEnum("Protocol", map[string]*Field{
		"port": {ft: numberType, mods: []FieldModifier{FieldModifierRequired}},
	}, map[string]Value{
		"tcp": &NumberValue{Value: big.NewFloat(6)},
	})

	tcp := protocolEnum.GetValue("tcp")
	tcp.SetField("port", &NumberValue{Value: big.NewFloat(8080)})
	assert.Equal(t, "tcp(port: 8080)", tcp.Text())

	// 测试ADT枚举的文本表示
	networkEnum := NewADTEnum("Network", map[string]map[string]Type{
		"TCP": {"src_port": numberType, "dst_port": numberType},
	})

	tcpPacket := networkEnum.GetVariant("TCP")
	tcpPacket.SetADTField("src_port", &NumberValue{Value: big.NewFloat(8080)})
	tcpPacket.SetADTField("dst_port", &NumberValue{Value: big.NewFloat(80)})
	assert.Equal(t, "TCP{src_port: 8080, dst_port: 80}", tcpPacket.Text())
}

// 测试枚举工厂函数
func TestEnumFactoryFunctions(t *testing.T) {
	// 测试 NewNumericEnum
	numericEnum := NewNumericEnum("Test", []string{"A", "B", "C"})
	assert.Equal(t, 3, len(numericEnum.GetVariants()))
	assert.Equal(t, "0", numericEnum.GetValue("A").Value().Text())
	assert.Equal(t, "1", numericEnum.GetValue("B").Value().Text())
	assert.Equal(t, "2", numericEnum.GetValue("C").Value().Text())

	// 测试 NewExplicitNumericEnum
	explicitEnum := NewExplicitNumericEnum("Test", map[string]int{
		"X": 10,
		"Y": 20,
		"Z": 30,
	})
	assert.Equal(t, 3, len(explicitEnum.GetVariants()))
	assert.Equal(t, "10", explicitEnum.GetValue("X").Value().Text())
	assert.Equal(t, "20", explicitEnum.GetValue("Y").Value().Text())
	assert.Equal(t, "30", explicitEnum.GetValue("Z").Value().Text())

	// 测试 NewStringEnum
	stringEnum := NewStringEnum("Test", map[string]string{
		"Alpha": "A",
		"Beta":  "B",
	})
	assert.Equal(t, 2, len(stringEnum.GetVariants()))
	assert.Equal(t, "A", stringEnum.GetValue("Alpha").Value().Text())
	assert.Equal(t, "B", stringEnum.GetValue("Beta").Value().Text())

	// 测试 NewMixedEnum
	mixedEnum := NewMixedEnum("Test", map[string]Value{
		"Number": &NumberValue{Value: big.NewFloat(42)},
		"String": &StringValue{Value: "hello"},
		"Bool":   &BoolValue{Value: true},
	})
	assert.Equal(t, 3, len(mixedEnum.GetVariants()))
	assert.Equal(t, "42", mixedEnum.GetValue("Number").Value().Text())
	assert.Equal(t, "hello", mixedEnum.GetValue("String").Value().Text())
	assert.Equal(t, "true", mixedEnum.GetValue("Bool").Value().Text())
}

// 测试枚举变体的完整功能
func TestEnumVariantFullFunctionality(t *testing.T) {
	// 创建枚举变体
	variant := NewEnumVariant("TestVariant", &NumberValue{Value: big.NewFloat(123)})

	// 测试基本属性
	assert.Equal(t, "TestVariant", variant.Name())
	assert.Equal(t, "123", variant.Value().Text())
	assert.Equal(t, "TestVariant", variant.Text())

	// 测试接口方法
	interfaceValue := variant.Interface()
	assert.NotNil(t, interfaceValue)
	assert.Equal(t, variant, interfaceValue)

	// 测试字段操作
	variant.SetField("custom_field", &StringValue{Value: "custom_value"})
	assert.Equal(t, "custom_value", variant.GetField("custom_field").Text())

	// 测试不存在的字段
	assert.Equal(t, NULL, variant.GetField("non_existent"))

	// 测试ADT字段操作
	variant.SetADTField("adt_field", &NumberValue{Value: big.NewFloat(456)})
	assert.Equal(t, "456", variant.GetADTField("adt_field").Text())

	// 测试不存在的ADT字段
	assert.Equal(t, NULL, variant.GetADTField("non_existent_adt"))

	// 测试ADT字段类型
	variant.AddADTFieldType("typed_field", stringType)
	// 这里只是添加类型定义，不涉及值操作
}

// 测试枚举类型的方法调用
func TestEnumMethodCalls(t *testing.T) {
	statusEnum := NewNumericEnum("Status", []string{"Pending", "Approved", "Rejected"})
	pending := statusEnum.GetValue("Pending")

	// 测试 values 方法
	valuesMethod := statusEnum.MethodByName("values")
	assert.NotNil(t, valuesMethod)

	result, err := valuesMethod.Call([]Value{pending}, nil)
	assert.NoError(t, err)
	assert.NotNil(t, result)

	// 验证返回的是数组
	arrayResult, ok := result.(*ArrayValue)
	assert.True(t, ok)
	assert.Equal(t, 3, len(arrayResult.Elements))

	// 测试 from_value 方法
	fromValueMethod := statusEnum.MethodByName("from_value")
	assert.NotNil(t, fromValueMethod)

	// 测试有效值
	result, err = fromValueMethod.Call([]Value{pending, &NumberValue{Value: big.NewFloat(1)}}, nil)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Approved", result.(*EnumVariant).Name())

	// 测试无效值
	result, err = fromValueMethod.Call([]Value{pending, &NumberValue{Value: big.NewFloat(99)}}, nil)
	assert.NoError(t, err)
	assert.Equal(t, NULL, result)

	// 测试 is_valid 方法
	isValidMethod := statusEnum.MethodByName("is_valid")
	assert.NotNil(t, isValidMethod)

	// 测试有效值
	result, err = isValidMethod.Call([]Value{pending, &NumberValue{Value: big.NewFloat(1)}}, nil)
	assert.NoError(t, err)
	assert.Equal(t, "true", result.Text())

	// 测试无效值
	result, err = isValidMethod.Call([]Value{pending, &NumberValue{Value: big.NewFloat(99)}}, nil)
	assert.NoError(t, err)
	assert.Equal(t, "false", result.Text())
}

// 测试枚举类型兼容性和转换
func TestEnumCompatibilityAndConversion(t *testing.T) {
	// 创建两个相同名称的枚举
	enum1 := NewNumericEnum("Status", []string{"Pending", "Approved"})
	enum2 := NewNumericEnum("Status", []string{"Pending", "Approved"})

	// 创建不同名称的枚举
	differentEnum := NewNumericEnum("Different", []string{"Pending", "Approved"})

	// 测试类型兼容性
	assert.True(t, enum1.AssignableTo(enum2))
	assert.False(t, enum1.AssignableTo(differentEnum))

	// 测试类型转换
	assert.False(t, enum1.ConvertibleTo(differentEnum))
	assert.False(t, enum1.ConvertibleTo(stringType))

	// 测试接口实现
	assert.False(t, enum1.Implements(differentEnum))
	assert.False(t, enum1.Implements(stringType))
}

// 测试枚举变体的字符串表示
func TestEnumVariantStringRepresentation(t *testing.T) {
	// 测试普通枚举变体
	simpleVariant := NewEnumVariant("Simple", &NumberValue{Value: big.NewFloat(42)})
	assert.Equal(t, "Simple", simpleVariant.Text())

	// 测试有关联字段的枚举变体（Dart风格）
	associatedVariant := NewEnumVariant("Associated", &NumberValue{Value: big.NewFloat(6)})
	associatedVariant.SetField("port", &NumberValue{Value: big.NewFloat(8080)})
	associatedVariant.SetField("host", &StringValue{Value: "localhost"})
	assert.Equal(t, "Associated(port: 8080, host: localhost)", associatedVariant.Text())

	// 测试ADT枚举变体（Rust风格）
	// 创建一个ADT枚举类型
	adtEnum := NewADTEnum("TestADT", map[string]map[string]Type{
		"ADT": {
			"field1": numberType,
			"field2": stringType,
		},
	})

	// 获取ADT变体并设置字段
	adtVariant := adtEnum.GetVariant("ADT")
	adtVariant.SetADTField("field1", &NumberValue{Value: big.NewFloat(123)})
	adtVariant.SetADTField("field2", &StringValue{Value: "test"})
	assert.Equal(t, "ADT{field1: 123, field2: test}", adtVariant.Text())

	// 测试空字段的枚举变体
	emptyVariant := NewEnumVariant("Empty", &StringValue{Value: "Empty"})
	assert.Equal(t, "Empty", emptyVariant.Text())
}

// 测试anyType的函数匹配
func TestAnyTypeFunctionMatching(t *testing.T) {
	// 创建一个接受anyType参数的函数
	anyTypeFunction := &FunctionType{
		name: "test_any",
		signatures: []*FunctionSignature{
			{
				positionalParams: []*Parameter{
					{name: "param", typ: anyType},
				},
				returnType: stringType,
				builtin: func(args ...Value) Value {
					return &StringValue{Value: "any_type_result"}
				},
			},
		},
	}

	// 创建一个接受具体类型的函数
	specificFunction := &FunctionType{
		name: "test_specific",
		signatures: []*FunctionSignature{
			{
				positionalParams: []*Parameter{
					{name: "param", typ: stringType},
				},
				returnType: stringType,
				builtin: func(args ...Value) Value {
					return &StringValue{Value: "specific_result"}
				},
			},
		},
	}

	// 测试字符串参数应该优先匹配具体类型函数
	stringArg := &StringValue{Value: "test"}

	// 测试anyType函数匹配
	result, err := anyTypeFunction.Call([]Value{stringArg}, nil)
	assert.NoError(t, err)
	assert.Equal(t, "any_type_result", result.Text())

	// 测试具体类型函数匹配
	result, err = specificFunction.Call([]Value{stringArg}, nil)
	assert.NoError(t, err)
	assert.Equal(t, "specific_result", result.Text())

	// 测试数字参数只能匹配anyType函数
	numberArg := &NumberValue{Value: big.NewFloat(42)}

	result, err = anyTypeFunction.Call([]Value{numberArg}, nil)
	assert.NoError(t, err)
	assert.Equal(t, "any_type_result", result.Text())

	// 数字参数不应该匹配字符串类型函数
	_, err = specificFunction.Call([]Value{numberArg}, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no matching signature found")
}

// 测试函数匹配的优先级
func TestFunctionMatchingPriority(t *testing.T) {
	// 创建一个有多个签名的函数
	multiSignatureFunction := &FunctionType{
		name: "test_priority",
		signatures: []*FunctionSignature{
			// 精确匹配：+100分
			{
				positionalParams: []*Parameter{
					{name: "param", typ: stringType},
				},
				returnType: stringType,
				builtin: func(args ...Value) Value {
					return &StringValue{Value: "exact_match"}
				},
			},
			// anyType匹配：+10分
			{
				positionalParams: []*Parameter{
					{name: "param", typ: anyType},
				},
				returnType: stringType,
				builtin: func(args ...Value) Value {
					return &StringValue{Value: "any_match"}
				},
			},
		},
	}

	// 字符串参数应该匹配精确签名
	stringArg := &StringValue{Value: "test"}
	result, err := multiSignatureFunction.Call([]Value{stringArg}, nil)
	assert.NoError(t, err)
	assert.Equal(t, "exact_match", result.Text())

	// 数字参数应该匹配anyType签名
	numberArg := &NumberValue{Value: big.NewFloat(42)}
	result, err = multiSignatureFunction.Call([]Value{numberArg}, nil)
	assert.NoError(t, err)
	assert.Equal(t, "any_match", result.Text())
}

// 测试静态方法和继承
func TestStaticMethodsAndInheritance(t *testing.T) {
	// 1. 测试静态方法
	counterClass := NewClassType("Counter")

	// 添加静态字段
	counterClass.AddStaticField("count", &Field{
		ft:   numberType,
		mods: []FieldModifier{FieldModifierStatic},
	})

	// 设置初始值
	counterClass.SetStaticValue("count", &NumberValue{Value: big.NewFloat(0)})

	// 添加静态方法
	counterClass.AddStaticMethod("increment", &FunctionType{
		name: "increment",
		signatures: []*FunctionSignature{
			{
				positionalParams: []*Parameter{},
				returnType:       numberType,
				builtin: func(args ...Value) Value {
					currentCount := counterClass.GetStaticValue("count").(*NumberValue)
					newCount := &NumberValue{Value: big.NewFloat(0)}
					newCount.Value.Add(currentCount.Value, big.NewFloat(1))
					counterClass.SetStaticValue("count", newCount)
					return newCount
				},
			},
		},
	})

	// 测试静态方法调用
	incrementMethod := counterClass.GetStaticMethod("increment")
	assert.NotNil(t, incrementMethod)

	result1, err := incrementMethod.Call([]Value{}, nil)
	assert.NoError(t, err)
	assert.Equal(t, "1", result1.Text())

	result2, err := incrementMethod.Call([]Value{}, nil)
	assert.NoError(t, err)
	assert.Equal(t, "2", result2.Text())

	// 测试静态字段访问
	staticCount := counterClass.GetStaticValue("count")
	assert.Equal(t, "2", staticCount.Text())

	// 2. 测试继承
	pointClass := NewClassType("Point")
	pointClass.AddField("x", &Field{
		ft:   numberType,
		mods: []FieldModifier{FieldModifierRequired},
	})
	pointClass.AddField("y", &Field{
		ft:   numberType,
		mods: []FieldModifier{FieldModifierRequired},
	})

	// 添加方法
	pointClass.AddMethod("to_str", &FunctionType{
		name: "to_str",
		signatures: []*FunctionSignature{
			{
				positionalParams: []*Parameter{
					{name: "self", typ: pointClass},
				},
				returnType: stringType,
				builtin: func(args ...Value) Value {
					point := args[0].(*ClassInstance)
					x := point.GetField("x").(*NumberValue)
					y := point.GetField("y").(*NumberValue)
					return &StringValue{Value: fmt.Sprintf("Point(%s, %s)", x.Text(), y.Text())}
				},
			},
		},
	})

	// 创建子类
	colorPointClass := NewClassType("ColorPoint")
	colorPointClass.SetParent(pointClass)
	colorPointClass.AddField("color", &Field{
		ft:   stringType,
		mods: []FieldModifier{FieldModifierRequired},
	})

	// 重写方法
	colorPointClass.OverrideMethod("to_str", &FunctionType{
		name: "to_str",
		signatures: []*FunctionSignature{
			{
				positionalParams: []*Parameter{
					{name: "self", typ: colorPointClass},
				},
				returnType: stringType,
				builtin: func(args ...Value) Value {
					colorPoint := args[0].(*ClassInstance)
					x := colorPoint.GetField("x").(*NumberValue)
					y := colorPoint.GetField("y").(*NumberValue)
					color := colorPoint.GetField("color").(*StringValue)
					return &StringValue{Value: fmt.Sprintf("ColorPoint(%s, %s, %s)", x.Text(), y.Text(), color.Value)}
				},
			},
		},
	})

	// 测试继承关系
	assert.True(t, colorPointClass.IsSubclassOf(pointClass))
	assert.Nil(t, pointClass.GetParent())
	assert.Equal(t, pointClass, colorPointClass.GetParent())

	// 测试方法重写
	point := &ClassInstance{
		clsType: pointClass,
		fields: map[string]Value{
			"x": &NumberValue{Value: big.NewFloat(10)},
			"y": &NumberValue{Value: big.NewFloat(20)},
		},
	}

	colorPoint := &ClassInstance{
		clsType: colorPointClass,
		fields: map[string]Value{
			"x":     &NumberValue{Value: big.NewFloat(30)},
			"y":     &NumberValue{Value: big.NewFloat(40)},
			"color": &StringValue{Value: "red"},
		},
	}

	// 测试父类方法
	pointToStr := pointClass.MethodByName("to_str")
	pointResult, err := pointToStr.Call([]Value{point}, nil)
	assert.NoError(t, err)
	assert.Equal(t, "Point(10, 20)", pointResult.Text())

	// 测试重写的方法
	colorPointToStr := colorPointClass.MethodByName("to_str")
	colorPointResult, err := colorPointToStr.Call([]Value{colorPoint}, nil)
	assert.NoError(t, err)
	assert.Equal(t, "ColorPoint(30, 40, red)", colorPointResult.Text())

	// 测试字段访问
	assert.Equal(t, "10", point.GetField("x").Text())
	assert.Equal(t, "30", colorPoint.GetField("x").Text())
	assert.Equal(t, "red", colorPoint.GetField("color").Text())
}

// 测试静态字段访问
func TestStaticFieldAccess(t *testing.T) {
	class := NewClassType("TestClass")

	// 添加静态字段
	class.AddStaticField("staticField", &Field{
		ft:   stringType,
		mods: []FieldModifier{FieldModifierStatic},
	})

	// 设置静态字段值
	class.SetStaticValue("staticField", &StringValue{Value: "static value"})

	// 测试静态字段访问
	staticValue := class.GetStaticValue("staticField")
	assert.Equal(t, "static value", staticValue.Text())

	// 测试不存在的静态字段
	nonExistent := class.GetStaticValue("nonExistent")
	assert.Equal(t, NULL, nonExistent)
}

// 测试方法重写
func TestMethodOverride(t *testing.T) {
	parentClass := NewClassType("Parent")
	childClass := NewClassType("Child")
	childClass.SetParent(parentClass)

	// 父类方法
	parentClass.AddMethod("method", &FunctionType{
		name: "method",
		signatures: []*FunctionSignature{
			{
				positionalParams: []*Parameter{
					{name: "self", typ: parentClass},
				},
				returnType: stringType,
				builtin: func(args ...Value) Value {
					return &StringValue{Value: "parent method"}
				},
			},
		},
	})

	// 子类重写方法
	childClass.OverrideMethod("method", &FunctionType{
		name: "method",
		signatures: []*FunctionSignature{
			{
				positionalParams: []*Parameter{
					{name: "self", typ: childClass},
				},
				returnType: stringType,
				builtin: func(args ...Value) Value {
					return &StringValue{Value: "child method"}
				},
			},
		},
	})

	// 测试方法调用
	parentInstance := &ClassInstance{clsType: parentClass, fields: make(map[string]Value)}
	childInstance := &ClassInstance{clsType: childClass, fields: make(map[string]Value)}

	parentMethod := parentClass.MethodByName("method")
	childMethod := childClass.MethodByName("method")

	parentResult, _ := parentMethod.Call([]Value{parentInstance}, nil)
	childResult, _ := childMethod.Call([]Value{childInstance}, nil)

	assert.Equal(t, "parent method", parentResult.Text())
	assert.Equal(t, "child method", childResult.Text())
}

// 测试trait、泛型和继承联动
func TestTraitGenericInheritanceIntegration(t *testing.T) {
	// 清理全局注册表，确保测试独立性
	GlobalTraitRegistry = NewTraitRegistry()

	// 1. 定义基础trait
	displayTrait := NewTraitType("Display")
	displayTrait.AddMethod("to_str", &FunctionType{
		name: "to_str",
		signatures: []*FunctionSignature{
			{
				positionalParams: []*Parameter{
					{name: "self", typ: displayTrait},
				},
				returnType: stringType,
			},
		},
	})

	cloneTrait := NewTraitType("Clone")
	cloneTrait.AddMethod("clone", &FunctionType{
		name: "clone",
		signatures: []*FunctionSignature{
			{
				positionalParams: []*Parameter{
					{name: "self", typ: cloneTrait},
				},
				returnType: anyType,
			},
		},
	})

	// 注册trait
	GlobalTraitRegistry.RegisterTrait(displayTrait)
	GlobalTraitRegistry.RegisterTrait(cloneTrait)

	// 2. 创建基础类
	pointClass := NewClassType("Point")
	pointClass.AddField("x", &Field{
		ft:   numberType,
		mods: []FieldModifier{FieldModifierRequired},
	})
	pointClass.AddField("y", &Field{
		ft:   numberType,
		mods: []FieldModifier{FieldModifierRequired},
	})

	// 为Point实现Display trait
	pointDisplayImpl := NewTraitImpl(displayTrait, pointClass)
	pointDisplayImpl.AddMethod("to_str", &FunctionType{
		name: "to_str",
		signatures: []*FunctionSignature{
			{
				positionalParams: []*Parameter{
					{name: "self", typ: pointClass},
				},
				returnType: stringType,
				builtin: func(args ...Value) Value {
					point := args[0].(*ClassInstance)
					x := point.GetField("x").(*NumberValue)
					y := point.GetField("y").(*NumberValue)
					return &StringValue{Value: fmt.Sprintf("Point(%s, %s)", x.Text(), y.Text())}
				},
			},
		},
	})
	GlobalTraitRegistry.RegisterImpl(pointDisplayImpl)

	// 为Point实现Clone trait
	pointCloneImpl := NewTraitImpl(cloneTrait, pointClass)
	pointCloneImpl.AddMethod("clone", &FunctionType{
		name: "clone",
		signatures: []*FunctionSignature{
			{
				positionalParams: []*Parameter{
					{name: "self", typ: pointClass},
				},
				returnType: pointClass,
				builtin: func(args ...Value) Value {
					point := args[0].(*ClassInstance)
					return &ClassInstance{
						clsType: pointClass,
						fields: map[string]Value{
							"x": point.GetField("x"),
							"y": point.GetField("y"),
						},
					}
				},
			},
		},
	})
	GlobalTraitRegistry.RegisterImpl(pointCloneImpl)

	// 3. 创建泛型容器类，带有trait约束
	containerClass := NewGenericClass("Container", []string{"T"}, nil)

	// 添加trait约束：T必须实现Display和Clone
	containerClass.AddTraitConstraint("T", displayTrait)
	containerClass.AddTraitConstraint("T", cloneTrait)

	// 添加字段
	containerClass.AddField("item", &Field{
		ft:   anyType,
		mods: []FieldModifier{FieldModifierRequired},
	})

	// 添加方法
	containerClass.AddMethod("display_item", &FunctionType{
		name: "display_item",
		signatures: []*FunctionSignature{
			{
				positionalParams: []*Parameter{
					{name: "self", typ: containerClass},
				},
				returnType: stringType,
				builtin: func(args ...Value) Value {
					container := args[0].(*ClassInstance)
					item := container.GetField("item")

					if classInstance, ok := item.(*ClassInstance); ok {
						// 通过trait registry查找Display trait的实现
						impls := GlobalTraitRegistry.GetImpl("Display", classInstance.clsType.Name())
						if len(impls) > 0 {
							displayMethod := impls[0].GetMethod("to_str")
							if displayMethod != nil {
								result, _ := displayMethod.Call([]Value{classInstance}, nil)
								return result
							}
						}
					}
					return &StringValue{Value: "unknown"}
				},
			},
		},
	})

	// 4. 创建继承的类
	colorPointClass := NewClassType("ColorPoint")
	colorPointClass.SetParent(pointClass)
	colorPointClass.AddField("color", &Field{
		ft:   stringType,
		mods: []FieldModifier{FieldModifierRequired},
	})

	// 为ColorPoint实现Display trait（重写父类实现）
	colorPointDisplayImpl := NewTraitImpl(displayTrait, colorPointClass)
	colorPointDisplayImpl.AddMethod("to_str", &FunctionType{
		name: "to_str",
		signatures: []*FunctionSignature{
			{
				positionalParams: []*Parameter{
					{name: "self", typ: colorPointClass},
				},
				returnType: stringType,
				builtin: func(args ...Value) Value {
					colorPoint := args[0].(*ClassInstance)
					x := colorPoint.GetField("x").(*NumberValue)
					y := colorPoint.GetField("y").(*NumberValue)
					color := colorPoint.GetField("color").(*StringValue)
					return &StringValue{Value: fmt.Sprintf("ColorPoint(%s, %s, %s)", x.Text(), y.Text(), color.Value)}
				},
			},
		},
	})
	GlobalTraitRegistry.RegisterImpl(colorPointDisplayImpl)

	// 为ColorPoint实现Clone trait
	colorPointCloneImpl := NewTraitImpl(cloneTrait, colorPointClass)
	colorPointCloneImpl.AddMethod("clone", &FunctionType{
		name: "clone",
		signatures: []*FunctionSignature{
			{
				positionalParams: []*Parameter{
					{name: "self", typ: colorPointClass},
				},
				returnType: colorPointClass,
				builtin: func(args ...Value) Value {
					colorPoint := args[0].(*ClassInstance)
					return &ClassInstance{
						clsType: colorPointClass,
						fields: map[string]Value{
							"x":     colorPoint.GetField("x"),
							"y":     colorPoint.GetField("y"),
							"color": colorPoint.GetField("color"),
						},
					}
				},
			},
		},
	})
	GlobalTraitRegistry.RegisterImpl(colorPointCloneImpl)

	// 5. 测试trait约束检查
	t.Run("Trait约束检查", func(t *testing.T) {
		// 检查Point是否满足trait约束
		err := containerClass.CheckTraitConstraints([]Type{pointClass})
		assert.NoError(t, err)

		// 检查ColorPoint是否满足trait约束
		err = containerClass.CheckTraitConstraints([]Type{colorPointClass})
		assert.NoError(t, err)

		// 检查不满足约束的类型
		simpleClass := NewClassType("Simple")
		err = containerClass.CheckTraitConstraints([]Type{simpleClass})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "does not satisfy trait constraint")
	})

	// 6. 测试泛型实例化
	t.Run("泛型实例化", func(t *testing.T) {
		// 实例化Container<Point>
		pointContainerInstance, err := containerClass.InstantiateWithTraitConstraints([]Type{pointClass})
		assert.NoError(t, err)
		assert.NotNil(t, pointContainerInstance)

		// 实例化Container<ColorPoint>
		colorPointContainerInstance, err := containerClass.InstantiateWithTraitConstraints([]Type{colorPointClass})
		assert.NoError(t, err)
		assert.NotNil(t, colorPointContainerInstance)
	})

	// 7. 测试继承和trait实现
	t.Run("继承和trait实现", func(t *testing.T) {
		// 检查继承关系
		assert.True(t, colorPointClass.IsSubclassOf(pointClass))

		// 检查trait实现
		assert.True(t, pointClass.ImplementsTrait(displayTrait))
		assert.True(t, pointClass.ImplementsTrait(cloneTrait))
		assert.True(t, colorPointClass.ImplementsTrait(displayTrait))
		assert.True(t, colorPointClass.ImplementsTrait(cloneTrait))
	})

	// 8. 测试实际使用
	t.Run("实际使用", func(t *testing.T) {
		// 创建实例
		point := &ClassInstance{
			clsType: pointClass,
			fields: map[string]Value{
				"x": &NumberValue{Value: big.NewFloat(10)},
				"y": &NumberValue{Value: big.NewFloat(20)},
			},
		}

		colorPoint := &ClassInstance{
			clsType: colorPointClass,
			fields: map[string]Value{
				"x":     &NumberValue{Value: big.NewFloat(30)},
				"y":     &NumberValue{Value: big.NewFloat(40)},
				"color": &StringValue{Value: "red"},
			},
		}

		// 测试trait方法调用 - 通过trait实现调用
		pointDisplayImpl := GlobalTraitRegistry.GetImpl("Display", "Point")[0]
		pointDisplayMethod := pointDisplayImpl.GetMethod("to_str")
		assert.NotNil(t, pointDisplayMethod)
		pointDisplayResult, _ := pointDisplayMethod.Call([]Value{point}, nil)
		assert.Equal(t, "Point(10, 20)", pointDisplayResult.Text())

		colorPointDisplayImpl := GlobalTraitRegistry.GetImpl("Display", "ColorPoint")[0]
		colorPointDisplayMethod := colorPointDisplayImpl.GetMethod("to_str")
		assert.NotNil(t, colorPointDisplayMethod)
		colorPointDisplayResult, _ := colorPointDisplayMethod.Call([]Value{colorPoint}, nil)
		assert.Equal(t, "ColorPoint(30, 40, red)", colorPointDisplayResult.Text())

		// 测试容器方法
		containerInstance, _ := containerClass.InstantiateWithTraitConstraints([]Type{pointClass})
		pointContainer := &ClassInstance{
			clsType: containerInstance,
			fields: map[string]Value{
				"item": point,
			},
		}

		displayMethod := containerInstance.MethodByName("display_item")
		assert.NotNil(t, displayMethod)
		displayResult, _ := displayMethod.Call([]Value{pointContainer}, nil)
		assert.Equal(t, "Point(10, 20)", displayResult.Text())
	})
}

// ==================== 泛型系统单元测试 ====================

func TestGenericBase(t *testing.T) {
	// 测试GenericBase基础功能
	gb := NewGenericBase([]string{"T", "U"}, map[string]Constraint{
		"T": NewInterfaceConstraint(stringType),
		"U": NewInterfaceConstraint(numberType),
	})

	if !gb.IsGeneric() {
		t.Error("GenericBase should be generic")
	}

	if len(gb.TypeParameters()) != 2 {
		t.Error("GenericBase should have 2 type parameters")
	}

	if len(gb.Constraints()) != 2 {
		t.Error("GenericBase should have 2 constraints")
	}

	// 测试约束检查 - 正确的类型应该通过
	err := gb.CheckConstraints([]Type{stringType, numberType})
	if err != nil {
		t.Errorf("Constraint check should pass: %v", err)
	}

	// 测试约束检查 - 错误的类型应该失败
	err = gb.CheckConstraints([]Type{numberType, stringType})
	if err == nil {
		t.Error("Constraint check should fail")
	}
}

func TestClassTypeGeneric(t *testing.T) {
	// 测试泛型类
	containerClass := NewGenericClass("Container", []string{"T"}, map[string]Constraint{
		"T": NewInterfaceConstraint(stringType),
	})

	if !containerClass.IsGeneric() {
		t.Error("Container class should be generic")
	}

	// 测试实例化
	instance, err := containerClass.Instantiate([]Type{stringType})
	if err != nil {
		t.Errorf("Instantiation should succeed: %v", err)
	}

	if instance == nil {
		t.Error("Instance should not be nil")
	}

	// 测试约束违反
	_, err = containerClass.Instantiate([]Type{numberType})
	if err == nil {
		t.Error("Instantiation should fail with wrong type")
	}
}

func TestFunctionTypeGeneric(t *testing.T) {
	// 测试泛型函数
	identityFunc := NewGenericFunction("identity", []string{"T"}, nil)
	identityFunc.signatures = []*FunctionSignature{
		{
			positionalParams: []*Parameter{
				{name: "x", typ: NewTypeParameter("T", nil)},
			},
			returnType: NewTypeParameter("T", nil),
		},
	}

	if !identityFunc.IsGeneric() {
		t.Error("Identity function should be generic")
	}

	if len(identityFunc.TypeParameters()) != 1 {
		t.Error("Identity function should have 1 type parameter")
	}
}

func TestTypeParameter(t *testing.T) {
	// 测试类型参数
	tp := NewTypeParameter("T", map[string]Constraint{
		"T": NewInterfaceConstraint(stringType),
	})

	if tp.Name() != "T" {
		t.Errorf("Type parameter name should be 'T', got '%s'", tp.Name())
	}

	if tp.Kind() != O_TYPE_PARAM {
		t.Errorf("Type parameter kind should be O_TYPE_PARAM, got %v", tp.Kind())
	}

	if tp.Text() != "T" {
		t.Errorf("Type parameter text should be 'T', got '%s'", tp.Text())
	}

	// 类型参数不能直接赋值
	if tp.AssignableTo(stringType) {
		t.Error("Type parameter should not be assignable to concrete type")
	}
}

func TestConstraintSystem(t *testing.T) {
	// 测试约束系统
	// 对应的 Hulo 源代码：
	// trait Display {
	//     fn to_str(self) -> str;
	// }
	//
	// trait Clone {
	//     fn clone(self) -> Self;
	// }
	//
	// trait Printable {
	//     fn print(self);
	// }
	//
	// // 为 String 实现 Display
	// impl Display for String {
	//     fn to_str(self) -> str { self }
	// }
	//
	// // 为 String 实现 Clone
	// impl Clone for String {
	//     fn clone(self) -> String { self }
	// }
	//
	// // 为 String 实现 Printable
	// impl Printable for String {
	//     fn print(self) { println(self) }
	// }
	//
	// // 泛型函数，要求 T 实现 Display
	// fn print_item<T: Display>(item: T) {
	//     println(item.to_str())
	// }
	//
	// // 泛型函数，要求 T 同时实现 Display 和 Clone
	// fn process_item<T: Display + Clone>(item: T) {
	//     let cloned = item.clone();
	//     println(item.to_str());
	//     println(cloned.to_str());
	// }
	//
	// // 泛型函数，要求 T 实现 Display 或 Clone
	// fn handle_item<T: Display | Clone>(item: T) {
	//     // 根据实际类型调用相应方法
	// }

	displayTrait := NewTraitType("Display")
	cloneTrait := NewTraitType("Clone")
	printableTrait := NewTraitType("Printable")

	// 注册trait
	GlobalTraitRegistry.RegisterTrait(displayTrait)
	GlobalTraitRegistry.RegisterTrait(cloneTrait)
	GlobalTraitRegistry.RegisterTrait(printableTrait)

	// 创建 String 类
	stringClass := NewClassType("String")

	// 为 String 实现 Display trait
	stringDisplayImpl := NewTraitImpl(displayTrait, stringClass)
	GlobalTraitRegistry.RegisterImpl(stringDisplayImpl)

	// 为 String 实现 Clone trait
	stringCloneImpl := NewTraitImpl(cloneTrait, stringClass)
	GlobalTraitRegistry.RegisterImpl(stringCloneImpl)

	// 为 String 实现 Printable trait
	stringPrintableImpl := NewTraitImpl(printableTrait, stringClass)
	GlobalTraitRegistry.RegisterImpl(stringPrintableImpl)

	// 测试1: String 满足 Display 约束
	// 对应: fn print_item<T: Display>(item: T)
	displayConstraint := NewInterfaceConstraint(displayTrait)
	if !displayConstraint.SatisfiedBy(stringClass) {
		t.Error("String should satisfy Display constraint")
	}

	// 测试2: String 满足复合约束 Display + Clone
	// 对应: fn process_item<T: Display + Clone>(item: T)
	compositeConstraint := NewCompositeConstraint(
		NewInterfaceConstraint(displayTrait),
		NewInterfaceConstraint(cloneTrait),
	)
	if !compositeConstraint.SatisfiedBy(stringClass) {
		t.Error("String should satisfy Display + Clone constraint")
	}

	// 测试3: String 满足联合约束 Display | Clone
	// 对应: fn handle_item<T: Display | Clone>(item: T)
	unionConstraint := NewUnionConstraint(
		NewInterfaceConstraint(displayTrait),
		NewInterfaceConstraint(cloneTrait),
	)
	if !unionConstraint.SatisfiedBy(stringClass) {
		t.Error("String should satisfy Display | Clone constraint")
	}

	// 测试4: 创建 Number 类，只实现 Display
	numberClass := NewClassType("Number")
	numberDisplayImpl := NewTraitImpl(displayTrait, numberClass)
	GlobalTraitRegistry.RegisterImpl(numberDisplayImpl)

	// Number 满足 Display 约束
	if !displayConstraint.SatisfiedBy(numberClass) {
		t.Error("Number should satisfy Display constraint")
	}

	// Number 不满足 Display + Clone 约束
	if compositeConstraint.SatisfiedBy(numberClass) {
		t.Error("Number should NOT satisfy Display + Clone constraint")
	}

	// Number 满足 Display | Clone 约束（因为满足 Display）
	if !unionConstraint.SatisfiedBy(numberClass) {
		t.Error("Number should satisfy Display | Clone constraint")
	}

	// 测试5: 运算符约束
	// 对应: fn add<T: Add>(a: T, b: T) -> T
	operatorConstraint := NewOperatorConstraint(OpAdd, numberType, numberType)
	if operatorConstraint.String() != "op" {
		t.Errorf("Operator constraint string should be 'op', got '%s'", operatorConstraint.String())
	}

	// 测试6: 约束字符串表示
	expectedCompositeStr := "implements Display AND implements Clone"
	if compositeConstraint.String() != expectedCompositeStr {
		t.Errorf("Composite constraint string should be '%s', got '%s'", expectedCompositeStr, compositeConstraint.String())
	}

	expectedUnionStr := "implements Display OR implements Clone"
	if unionConstraint.String() != expectedUnionStr {
		t.Errorf("Union constraint string should be '%s', got '%s'", expectedUnionStr, unionConstraint.String())
	}
}

func TestGenericInstantiation(t *testing.T) {
	// 测试泛型实例化
	containerClass := NewGenericClass("Container", []string{"T"}, nil)

	// 添加字段
	containerClass.AddField("item", &Field{
		ft:   NewTypeParameter("T", nil),
		mods: []FieldModifier{FieldModifierRequired},
	})

	// 实例化
	instance, err := containerClass.Instantiate([]Type{stringType})
	if err != nil {
		t.Errorf("Instantiation should succeed: %v", err)
	}

	// 检查实例化后的字段类型
	field := instance.FieldByName("item")
	if field == nil {
		t.Error("Instance should have 'item' field")
	}

	if field.Name() != "str" {
		t.Errorf("Field type should be 'str', got '%s'", field.Name())
	}
}

func TestConstraintViolation(t *testing.T) {
	// 测试约束违反
	containerClass := NewGenericClass("Container", []string{"T"}, map[string]Constraint{
		"T": NewInterfaceConstraint(stringType),
	})

	// 尝试用不满足约束的类型实例化
	_, err := containerClass.Instantiate([]Type{numberType})
	if err == nil {
		t.Error("Should fail with constraint violation")
	}

	expectedError := "type num does not satisfy constraint"
	if !strings.Contains(err.Error(), expectedError) {
		t.Errorf("Error should contain '%s', got '%s'", expectedError, err.Error())
	}
}

func TestGenericMethod(t *testing.T) {
	// 测试泛型方法
	containerClass := NewGenericClass("Container", []string{"T"}, nil)

	// 添加泛型方法
	containerClass.AddMethod("add", &FunctionType{
		name: "add",
		signatures: []*FunctionSignature{
			{
				positionalParams: []*Parameter{
					{name: "self", typ: containerClass},
					{name: "item", typ: NewTypeParameter("T", nil)},
				},
				returnType: voidType,
			},
		},
	})

	// 实例化
	instance, err := containerClass.Instantiate([]Type{stringType})
	if err != nil {
		t.Errorf("Instantiation should succeed: %v", err)
	}

	// 检查方法
	method := instance.MethodByName("add")
	if method == nil {
		t.Error("Instance should have 'add' method")
	}
}

func TestGenericCaching(t *testing.T) {
	// 测试泛型缓存
	containerClass := NewGenericClass("Container", []string{"T"}, nil)

	// 第一次实例化
	instance1, err := containerClass.Instantiate([]Type{stringType})
	if err != nil {
		t.Errorf("First instantiation should succeed: %v", err)
	}

	// 第二次实例化（应该使用缓存）
	instance2, err := containerClass.Instantiate([]Type{stringType})
	if err != nil {
		t.Errorf("Second instantiation should succeed: %v", err)
	}

	// 应该是同一个实例
	if instance1 != instance2 {
		t.Error("Cached instances should be the same")
	}
}

func TestComplexGenericConstraints(t *testing.T) {
	// 测试复杂约束
	displayTrait := NewTraitType("Display")
	cloneTrait := NewTraitType("Clone")

	// 创建复杂约束：T: Display + Clone
	complexConstraint := NewCompositeConstraint(
		NewInterfaceConstraint(displayTrait),
		NewInterfaceConstraint(cloneTrait),
	)

	containerClass := NewGenericClass("Container", []string{"T"}, map[string]Constraint{
		"T": complexConstraint,
	})

	// 测试约束字符串表示
	constraintStr := complexConstraint.String()
	expectedStr := "implements Display AND implements Clone"
	if constraintStr != expectedStr {
		t.Errorf("Constraint string should be '%s', got '%s'", expectedStr, constraintStr)
	}

	// 测试约束检查
	err := containerClass.CheckConstraints([]Type{displayTrait})
	if err == nil {
		t.Error("Should fail with complex constraint")
	}
}

func TestGenericInheritance(t *testing.T) {
	// 测试泛型继承
	baseClass := NewGenericClass("Base", []string{"T"}, nil)
	derivedClass := NewGenericClass("Derived", []string{"T"}, nil)
	derivedClass.SetParent(baseClass)

	if !derivedClass.IsSubclassOf(baseClass) {
		t.Error("Derived should be subclass of Base")
	}

	// 实例化派生类
	instance, err := derivedClass.Instantiate([]Type{stringType})
	if err != nil {
		t.Errorf("Derived class instantiation should succeed: %v", err)
	}

	if instance == nil {
		t.Error("Instance should not be nil")
	}
}
