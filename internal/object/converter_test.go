// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package object

import (
	"testing"

	"github.com/hulo-lang/hulo/syntax/hulo/ast"
	"github.com/hulo-lang/hulo/syntax/hulo/token"
	"github.com/stretchr/testify/assert"
)

// ==================== 基础字面量转换测试 ====================

func TestConvertBasicLiterals(t *testing.T) {
	converter := &ASTConverter{}

	tests := []struct {
		name string
		expr ast.Node
		want Value
	}{
		{
			name: "string literal",
			expr: &ast.StringLiteral{Value: "Hello, World!"},
			want: &StringValue{Value: "Hello, World!"},
		},
		{
			name: "numeric literal",
			expr: &ast.NumericLiteral{Value: "123.45"},
			want: NewNumberValue("123.45"),
		},
		{
			name: "true literal",
			expr: &ast.TrueLiteral{},
			want: TRUE,
		},
		{
			name: "false literal",
			expr: &ast.FalseLiteral{},
			want: FALSE,
		},
		{
			name: "null literal",
			expr: &ast.NullLiteral{},
			want: NULL,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := converter.ConvertValue(test.expr)
			assert.NoError(t, err)
			assert.Equal(t, test.want, got)
		})
	}
}

// ==================== 复合字面量转换测试 ====================

func TestConvertArrayLiteral(t *testing.T) {
	converter := &ASTConverter{}

	// 创建数组字面量 - 使用相同类型的元素
	arrayExpr := &ast.ArrayLiteralExpr{
		Lbrack: token.Pos(1),
		Elems: []ast.Expr{
			&ast.StringLiteral{Value: "hello"},
			&ast.StringLiteral{Value: "world"},
			&ast.StringLiteral{Value: "test"},
		},
		Rbrack: token.Pos(20),
	}

	value, err := converter.ConvertValue(arrayExpr)
	assert.NoError(t, err)
	assert.IsType(t, &ArrayValue{}, value)

	arrayValue := value.(*ArrayValue)
	assert.Len(t, arrayValue.Elements, 3)
	assert.Equal(t, "hello", arrayValue.Elements[0].(*StringValue).Value)
	assert.Equal(t, "world", arrayValue.Elements[1].(*StringValue).Value)
	assert.Equal(t, "test", arrayValue.Elements[2].(*StringValue).Value)
}

func TestConvertMixedTypeArrayLiteral(t *testing.T) {
	converter := &ASTConverter{}

	// 创建混合类型的数组字面量 - 应该推断为any[]
	arrayExpr := &ast.ArrayLiteralExpr{
		Lbrack: token.Pos(1),
		Elems: []ast.Expr{
			&ast.StringLiteral{Value: "hello"},
			&ast.NumericLiteral{Value: "42"},
			&ast.TrueLiteral{},
		},
		Rbrack: token.Pos(20),
	}

	value, err := converter.ConvertValue(arrayExpr)
	assert.NoError(t, err)
	assert.IsType(t, &ArrayValue{}, value)

	arrayValue := value.(*ArrayValue)
	assert.Len(t, arrayValue.Elements, 3)
	assert.Equal(t, "any", arrayValue.Type().(*ArrayType).ElementType().Name())

	// 验证元素值
	assert.Equal(t, "hello", arrayValue.Elements[0].(*StringValue).Value)
	assert.Equal(t, "42", arrayValue.Elements[1].(*NumberValue).Value.String())
	assert.Equal(t, TRUE, arrayValue.Elements[2])
}

func TestConvertObjectLiteral(t *testing.T) {
	converter := &ASTConverter{}

	// create object literal
	objectExpr := &ast.ObjectLiteralExpr{
		Lbrace: token.Pos(1),
		Props: []ast.Expr{
			&ast.KeyValueExpr{
				Key:   &ast.Ident{Name: "name"},
				Colon: token.Pos(10),
				Value: &ast.StringLiteral{Value: "John"},
			},
			&ast.KeyValueExpr{
				Key:   &ast.Ident{Name: "age"},
				Colon: token.Pos(20),
				Value: &ast.NumericLiteral{Value: "30"},
			},
		},
		Rbrace: token.Pos(30),
	}

	value, err := converter.ConvertValue(objectExpr)
	assert.NoError(t, err)
	assert.IsType(t, &MapValue{}, value)

	mapValue := value.(*MapValue)
	assert.Len(t, mapValue.Pairs, 2)

	// check first key-value pair
	assert.Equal(t, "name", mapValue.Pairs[0].Key.(*StringValue).Value)
	assert.Equal(t, "John", mapValue.Pairs[0].Value.(*StringValue).Value)

	// check second key-value pair
	assert.Equal(t, "age", mapValue.Pairs[1].Key.(*StringValue).Value)
	assert.Equal(t, "30", mapValue.Pairs[1].Value.(*NumberValue).Value.String())
}

func TestConvertSetLiteral(t *testing.T) {
	converter := &ASTConverter{}

	// create set literal
	setExpr := &ast.SetLiteralExpr{
		Lbrace: token.Pos(1),
		Elems: []ast.Expr{
			&ast.StringLiteral{Value: "apple"},
			&ast.StringLiteral{Value: "banana"},
			&ast.StringLiteral{Value: "cherry"},
		},
		Rbrace: token.Pos(20),
	}

	value, err := converter.ConvertValue(setExpr)
	assert.NoError(t, err)
	assert.IsType(t, &SetValue{}, value)

	setValue := value.(*SetValue)
	assert.Len(t, setValue.Elements(), 3)

	// check if elements exist
	hasApple := false
	hasBanana := false
	hasCherry := false

	for _, elem := range setValue.Elements() {
		switch elem.(*StringValue).Value {
		case "apple":
			hasApple = true
		case "banana":
			hasBanana = true
		case "cherry":
			hasCherry = true
		}
	}

	assert.True(t, hasApple)
	assert.True(t, hasBanana)
	assert.True(t, hasCherry)
}

func TestConvertTupleLiteral(t *testing.T) {
	converter := &ASTConverter{}

	// 创建元组字面量
	tupleExpr := &ast.TupleLiteralExpr{
		Lbrack: token.Pos(1),
		Elems: []ast.Expr{
			&ast.StringLiteral{Value: "hello"},
			&ast.NumericLiteral{Value: "42"},
			&ast.TrueLiteral{},
		},
		Rbrack: token.Pos(20),
	}

	value, err := converter.ConvertValue(tupleExpr)
	assert.NoError(t, err)
	assert.IsType(t, &TupleValue{}, value)

	tupleValue := value.(*TupleValue)
	assert.Len(t, tupleValue.Elements(), 3)
	assert.Equal(t, "hello", tupleValue.Elements()[0].(*StringValue).Value)
	assert.Equal(t, "42", tupleValue.Elements()[1].(*NumberValue).Value.String())
	assert.Equal(t, TRUE, tupleValue.Elements()[2])
}

// ==================== 类型转换测试 ====================

func TestConvertBasicTypes(t *testing.T) {
	converter := &ASTConverter{}

	tests := []struct {
		name string
		expr ast.Node
		want Type
	}{
		{
			name: "string type",
			expr: &ast.TypeReference{Name: &ast.Ident{Name: "str"}},
			want: stringType,
		},
		{
			name: "number type",
			expr: &ast.TypeReference{Name: &ast.Ident{Name: "num"}},
			want: numberType,
		},
		{
			name: "bool type",
			expr: &ast.TypeReference{Name: &ast.Ident{Name: "bool"}},
			want: boolType,
		},
		{
			name: "any type",
			expr: &ast.TypeReference{Name: &ast.Ident{Name: "any"}},
			want: anyType,
		},
		{
			name: "any literal",
			expr: &ast.AnyLiteral{},
			want: anyType,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := converter.ConvertType(test.expr)
			assert.NoError(t, err)
			assert.Equal(t, test.want, got)
		})
	}
}

func TestConvertArrayType(t *testing.T) {
	converter := &ASTConverter{}

	// 创建数组类型
	arrayTypeExpr := &ast.TypeReference{
		Name: &ast.Ident{Name: "array"},
		TypeParams: []ast.Expr{
			&ast.TypeReference{Name: &ast.Ident{Name: "str"}},
		},
	}

	typ, err := converter.ConvertType(arrayTypeExpr)
	assert.NoError(t, err)
	assert.IsType(t, &ArrayType{}, typ)

	arrayType := typ.(*ArrayType)
	assert.Equal(t, stringType, arrayType.ElementType())
}

func TestConvertMapType(t *testing.T) {
	converter := &ASTConverter{}

	// 创建映射类型
	mapTypeExpr := &ast.TypeReference{
		Name: &ast.Ident{Name: "map"},
		TypeParams: []ast.Expr{
			&ast.TypeReference{Name: &ast.Ident{Name: "str"}},
			&ast.TypeReference{Name: &ast.Ident{Name: "num"}},
		},
	}

	typ, err := converter.ConvertType(mapTypeExpr)
	assert.NoError(t, err)
	assert.IsType(t, &MapType{}, typ)

	mapType := typ.(*MapType)
	assert.Equal(t, stringType, mapType.KeyType())
	assert.Equal(t, numberType, mapType.ValueType())
}

func TestConvertFunctionType(t *testing.T) {
	converter := &ASTConverter{}

	// 创建函数类型
	funcTypeExpr := &ast.TypeReference{
		Name: &ast.Ident{Name: "function"},
		TypeParams: []ast.Expr{
			&ast.TypeReference{Name: &ast.Ident{Name: "str"}},
			&ast.TypeReference{Name: &ast.Ident{Name: "num"}},
			&ast.TypeReference{Name: &ast.Ident{Name: "bool"}},
		},
	}

	typ, err := converter.ConvertType(funcTypeExpr)
	assert.NoError(t, err)
	assert.IsType(t, &FunctionType{}, typ)

	funcType := typ.(*FunctionType)
	assert.Len(t, funcType.signatures, 1)
	sig := funcType.signatures[0]
	assert.Len(t, sig.positionalParams, 2)
	assert.Equal(t, stringType, sig.positionalParams[0].typ)
	assert.Equal(t, numberType, sig.positionalParams[1].typ)
	assert.Equal(t, boolType, sig.returnType)
}

func TestConvertTupleType(t *testing.T) {
	converter := &ASTConverter{}

	// 创建元组类型
	tupleTypeExpr := &ast.TypeReference{
		Name: &ast.Ident{Name: "tuple"},
		TypeParams: []ast.Expr{
			&ast.TypeReference{Name: &ast.Ident{Name: "str"}},
			&ast.TypeReference{Name: &ast.Ident{Name: "num"}},
			&ast.TypeReference{Name: &ast.Ident{Name: "bool"}},
		},
	}

	typ, err := converter.ConvertType(tupleTypeExpr)
	assert.NoError(t, err)
	assert.IsType(t, &TupleType{}, typ)

	tupleType := typ.(*TupleType)
	assert.Len(t, tupleType.elementTypes, 3)
	assert.Equal(t, stringType, tupleType.elementTypes[0])
	assert.Equal(t, numberType, tupleType.elementTypes[1])
	assert.Equal(t, boolType, tupleType.elementTypes[2])
}

func TestConvertSetType(t *testing.T) {
	converter := &ASTConverter{}

	// 创建集合类型
	setTypeExpr := &ast.TypeReference{
		Name: &ast.Ident{Name: "set"},
		TypeParams: []ast.Expr{
			&ast.TypeReference{Name: &ast.Ident{Name: "str"}},
		},
	}

	typ, err := converter.ConvertType(setTypeExpr)
	assert.NoError(t, err)
	assert.IsType(t, &SetType{}, typ)

	setType := typ.(*SetType)
	assert.Equal(t, stringType, setType.elementType)
}

func TestConvertUnionType(t *testing.T) {
	converter := &ASTConverter{}

	// 创建联合类型
	unionTypeExpr := &ast.TypeReference{
		Name: &ast.Ident{Name: "union"},
		TypeParams: []ast.Expr{
			&ast.TypeReference{Name: &ast.Ident{Name: "str"}},
			&ast.TypeReference{Name: &ast.Ident{Name: "num"}},
		},
	}

	typ, err := converter.ConvertType(unionTypeExpr)
	assert.NoError(t, err)
	assert.IsType(t, &UnionType{}, typ)

	unionType := typ.(*UnionType)
	assert.Len(t, unionType.types, 2)
	assert.Equal(t, stringType, unionType.types[0])
	assert.Equal(t, numberType, unionType.types[1])
}

func TestConvertIntersectionType(t *testing.T) {
	converter := &ASTConverter{}

	// 创建交集类型
	intersectionTypeExpr := &ast.TypeReference{
		Name: &ast.Ident{Name: "intersection"},
		TypeParams: []ast.Expr{
			&ast.TypeReference{Name: &ast.Ident{Name: "str"}},
			&ast.TypeReference{Name: &ast.Ident{Name: "num"}},
		},
	}

	typ, err := converter.ConvertType(intersectionTypeExpr)
	assert.NoError(t, err)
	assert.IsType(t, &IntersectionType{}, typ)

	intersectionType := typ.(*IntersectionType)
	assert.Len(t, intersectionType.types, 2)
	assert.Equal(t, stringType, intersectionType.types[0])
	assert.Equal(t, numberType, intersectionType.types[1])
}

func TestConvertNullableType(t *testing.T) {
	converter := &ASTConverter{}

	// 创建可空类型
	nullableTypeExpr := &ast.TypeReference{
		Name: &ast.Ident{Name: "nullable"},
		TypeParams: []ast.Expr{
			&ast.TypeReference{Name: &ast.Ident{Name: "str"}},
		},
	}

	typ, err := converter.ConvertType(nullableTypeExpr)
	assert.NoError(t, err)
	assert.IsType(t, &NullableType{}, typ)

	nullableType := typ.(*NullableType)
	assert.Equal(t, stringType, nullableType.baseType)
}

// ==================== 修饰符转换测试 ====================

func TestConvertModifiers(t *testing.T) {
	converter := &ASTConverter{}

	tests := []struct {
		name      string
		modifiers []ast.Modifier
		want      []FieldModifier
	}{
		{
			name: "public modifier",
			modifiers: []ast.Modifier{
				&ast.PubModifier{Pub: token.Pos(1)},
			},
			want: []FieldModifier{FieldModifierPub},
		},
		{
			name: "const modifier",
			modifiers: []ast.Modifier{
				&ast.ConstModifier{Const: token.Pos(1)},
			},
			want: []FieldModifier{FieldModifierReadonly},
		},
		{
			name: "final modifier",
			modifiers: []ast.Modifier{
				&ast.FinalModifier{Final: token.Pos(1)},
			},
			want: []FieldModifier{FieldModifierRequired},
		},
		{
			name: "multiple modifiers",
			modifiers: []ast.Modifier{
				&ast.PubModifier{Pub: token.Pos(1)},
				&ast.ConstModifier{Const: token.Pos(2)},
			},
			want: []FieldModifier{FieldModifierPub, FieldModifierReadonly},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := converter.ConvertModifiers(test.modifiers)
			assert.Equal(t, test.want, got)
		})
	}
}

// ==================== 字段转换测试 ====================

func TestConvertField(t *testing.T) {
	converter := &ASTConverter{}

	// 创建字段
	field := &ast.Field{
		Modifiers: []ast.Modifier{
			&ast.PubModifier{Pub: token.Pos(1)},
		},
		Name:   &ast.Ident{Name: "name"},
		Colon:  token.Pos(10),
		Type:   &ast.TypeReference{Name: &ast.Ident{Name: "str"}},
		Assign: token.Pos(15),
		Value:  &ast.StringLiteral{Value: "default"},
	}

	fieldObj, err := converter.ConvertField(field)
	assert.NoError(t, err)
	assert.NotNil(t, fieldObj)

	// 检查字段类型
	assert.Equal(t, stringType, fieldObj.ft)

	// 检查修饰符
	assert.Len(t, fieldObj.mods, 1)
	assert.Equal(t, FieldModifierPub, fieldObj.mods[0])

	// 检查默认值
	assert.True(t, fieldObj.HasDefaultValue())
	defaultValue := fieldObj.GetDefaultValue()
	assert.Equal(t, "default", defaultValue.(*StringValue).Value)
}

func TestConvertFieldWithoutDefault(t *testing.T) {
	converter := &ASTConverter{}

	// 创建没有默认值的字段
	field := &ast.Field{
		Name:  &ast.Ident{Name: "age"},
		Colon: token.Pos(10),
		Type:  &ast.TypeReference{Name: &ast.Ident{Name: "num"}},
	}

	fieldObj, err := converter.ConvertField(field)
	assert.NoError(t, err)
	assert.NotNil(t, fieldObj)

	// 检查字段类型
	assert.Equal(t, numberType, fieldObj.ft)

	// 检查没有默认值
	assert.False(t, fieldObj.HasDefaultValue())
	assert.Nil(t, fieldObj.GetDefaultValue())
}

// ==================== 参数转换测试 ====================

func TestConvertParameter(t *testing.T) {
	converter := &ASTConverter{}

	// 创建参数
	param := &ast.Parameter{
		Modifier: &ast.RequiredModifier{Required: token.Pos(1)},
		Name:     &ast.Ident{Name: "x"},
		Colon:    token.Pos(10),
		Type:     &ast.TypeReference{Name: &ast.Ident{Name: "num"}},
		Assign:   token.Pos(15),
		Value:    &ast.NumericLiteral{Value: "0"},
	}

	paramObj, err := converter.ConvertParameter(param)
	assert.NoError(t, err)
	assert.NotNil(t, paramObj)

	// 检查参数属性
	assert.Equal(t, "x", paramObj.name)
	assert.Equal(t, numberType, paramObj.typ)
	assert.True(t, paramObj.required)
	assert.False(t, paramObj.optional)
	assert.False(t, paramObj.variadic)
	assert.True(t, paramObj.isNamed)

	// 检查默认值
	assert.NotNil(t, paramObj.defaultValue)
	assert.Equal(t, "0", paramObj.defaultValue.(*NumberValue).Value.String())
}

func TestConvertVariadicParameter(t *testing.T) {
	converter := &ASTConverter{}

	// create variadic parameter
	param := &ast.Parameter{
		Modifier: &ast.EllipsisModifier{Ellipsis: token.Pos(1)},
		Name:     &ast.Ident{Name: "args"},
		Colon:    token.Pos(10),
		Type:     &ast.TypeReference{Name: &ast.Ident{Name: "any"}},
	}

	paramObj, err := converter.ConvertParameter(param)
	assert.NoError(t, err)
	assert.NotNil(t, paramObj)

	// check parameter attributes
	assert.Equal(t, "args", paramObj.name)
	assert.Equal(t, anyType, paramObj.typ)
	assert.True(t, paramObj.variadic)
	assert.False(t, paramObj.required)
	assert.True(t, paramObj.optional)
}

// ==================== 错误处理测试 ====================

func TestConvertInvalidType(t *testing.T) {
	converter := &ASTConverter{}

	// 测试无效的类型
	invalidType := &ast.Ident{Name: "invalid_type"}

	_, err := converter.ConvertType(invalidType)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unknown type")
}

func TestConvertInvalidValue(t *testing.T) {
	converter := &ASTConverter{}

	// 测试无效的值（使用一个不存在的节点类型）
	invalidValue := &ast.Ident{Name: "invalid_value"}

	_, err := converter.ConvertValue(invalidValue)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported node type")
}

func TestConvertFieldWithInvalidType(t *testing.T) {
	converter := &ASTConverter{}

	// 创建有无效类型的字段
	field := &ast.Field{
		Name:  &ast.Ident{Name: "test"},
		Colon: token.Pos(10),
		Type:  &ast.Ident{Name: "invalid_type"},
	}

	_, err := converter.ConvertField(field)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to convert field type")
}

func TestConvertFieldWithInvalidDefaultValue(t *testing.T) {
	converter := &ASTConverter{}

	// 创建有无效默认值的字段
	field := &ast.Field{
		Name:   &ast.Ident{Name: "test"},
		Colon:  token.Pos(10),
		Type:   &ast.Ident{Name: "str"},
		Assign: token.Pos(15),
		Value:  &ast.Ident{Name: "invalid_value"},
	}

	_, err := converter.ConvertField(field)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to convert field default value")
}

// ==================== 边界情况测试 ====================

func TestConvertEmptyArray(t *testing.T) {
	converter := &ASTConverter{}

	// 创建空数组
	arrayExpr := &ast.ArrayLiteralExpr{
		Lbrack: token.Pos(1),
		Elems:  []ast.Expr{},
		Rbrack: token.Pos(5),
	}

	value, err := converter.ConvertValue(arrayExpr)
	assert.NoError(t, err)
	assert.IsType(t, &ArrayValue{}, value)

	arrayValue := value.(*ArrayValue)
	assert.Len(t, arrayValue.Elements, 0)
}

func TestConvertEmptyObject(t *testing.T) {
	converter := &ASTConverter{}

	// 创建空对象
	objectExpr := &ast.ObjectLiteralExpr{
		Lbrace: token.Pos(1),
		Props:  []ast.Expr{},
		Rbrace: token.Pos(5),
	}

	value, err := converter.ConvertValue(objectExpr)
	assert.NoError(t, err)
	assert.IsType(t, &MapValue{}, value)

	mapValue := value.(*MapValue)
	assert.Len(t, mapValue.Pairs, 0)
}

func TestConvertEmptySet(t *testing.T) {
	converter := &ASTConverter{}

	// 创建空集合
	setExpr := &ast.SetLiteralExpr{
		Lbrace: token.Pos(1),
		Elems:  []ast.Expr{},
		Rbrace: token.Pos(5),
	}

	value, err := converter.ConvertValue(setExpr)
	assert.NoError(t, err)
	assert.IsType(t, &SetValue{}, value)

	setValue := value.(*SetValue)
	assert.Len(t, setValue.Elements(), 0)
}

func TestConvertEmptyTuple(t *testing.T) {
	converter := &ASTConverter{}

	// 创建空元组
	tupleExpr := &ast.TupleLiteralExpr{
		Lbrack: token.Pos(1),
		Elems:  []ast.Expr{},
		Rbrack: token.Pos(5),
	}

	value, err := converter.ConvertValue(tupleExpr)
	assert.NoError(t, err)
	assert.IsType(t, &TupleValue{}, value)

	tupleValue := value.(*TupleValue)
	assert.Len(t, tupleValue.Elements(), 0)
}

// ==================== 性能测试 ====================

func BenchmarkConvertStringLiteral(b *testing.B) {
	converter := &ASTConverter{}
	expr := &ast.StringLiteral{Value: "Hello, World!"}

	for b.Loop() {
		_, err := converter.ConvertValue(expr)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkConvertArrayLiteral(b *testing.B) {
	converter := &ASTConverter{}
	expr := &ast.ArrayLiteralExpr{
		Elems: []ast.Expr{
			&ast.StringLiteral{Value: "hello"},
			&ast.NumericLiteral{Value: "42"},
			&ast.TrueLiteral{},
		},
	}

	for b.Loop() {
		_, err := converter.ConvertValue(expr)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkConvertObjectLiteral(b *testing.B) {
	converter := &ASTConverter{}
	expr := &ast.ObjectLiteralExpr{
		Props: []ast.Expr{
			&ast.KeyValueExpr{
				Key:   &ast.Ident{Name: "name"},
				Value: &ast.StringLiteral{Value: "John"},
			},
			&ast.KeyValueExpr{
				Key:   &ast.Ident{Name: "age"},
				Value: &ast.NumericLiteral{Value: "30"},
			},
		},
	}

	for b.Loop() {
		_, err := converter.ConvertValue(expr)
		if err != nil {
			b.Fatal(err)
		}
	}
}
