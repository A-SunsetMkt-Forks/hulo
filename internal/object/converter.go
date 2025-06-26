// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package object

import (
	"fmt"
	"strings"

	"github.com/hulo-lang/hulo/syntax/hulo/ast"
)

type Converter interface {
	ConvertType(node ast.Node) (Type, error)
	ConvertValue(node ast.Node) (Value, error)
	ConvertModifiers(modifiers []ast.Modifier) []FieldModifier
	ConvertField(field *ast.Field) (*Field, error)
	ConvertParameter(param *ast.Parameter) (*Parameter, error)
}

var _ Converter = (*ASTConverter)(nil)

type ASTConverter struct{}

// ConvertValue 将AST节点转换为Value
func (c *ASTConverter) ConvertValue(node ast.Node) (Value, error) {
	switch n := node.(type) {
	// 基本字面量
	case *ast.StringLiteral:
		return &StringValue{Value: n.Value}, nil
	case *ast.NumericLiteral:
		return NewNumberValue(n.Value), nil
	case *ast.TrueLiteral:
		return TRUE, nil
	case *ast.FalseLiteral:
		return FALSE, nil
	case *ast.NullLiteral:
		return NULL, nil

	// 数组字面量
	case *ast.ArrayLiteralExpr:
		return c.convertArrayLiteral(n)

	// 对象字面量
	case *ast.ObjectLiteralExpr:
		return c.convertObjectLiteral(n)

	// 命名对象字面量
	case *ast.NamedObjectLiteralExpr:
		return c.convertNamedObjectLiteral(n)

	// 集合字面量
	case *ast.SetLiteralExpr:
		return c.convertSetLiteral(n)

	// 元组字面量
	case *ast.TupleLiteralExpr:
		return c.convertTupleLiteral(n)

	// 命名元组字面量
	case *ast.NamedTupleLiteralExpr:
		return c.convertNamedTupleLiteral(n)

	// 标识符（变量引用）
	case *ast.Ident:
		return c.convertIdentifier(n)
		return nil, fmt.Errorf("unsupported value type: %T", n)

	default:
		return nil, fmt.Errorf("unsupported value type: %T", n)
	}
}

// ConvertType 将AST节点转换为Type
func (c *ASTConverter) ConvertType(node ast.Node) (Type, error) {
	switch n := node.(type) {
	// 基本类型
	case *ast.Ident:
		// return c.convertBasicType(n)
		return nil, fmt.Errorf("unsupported type: %T", n)

	// 类型引用（包括泛型）
	case *ast.TypeReference:
		return c.convertTypeReference(n)

	// 数组类型
	case *ast.ArrayType:
		return c.convertArrayType(n)

	// 函数类型
	case *ast.FunctionType:
		return c.convertFunctionType(n)

	// 元组类型
	case *ast.TupleType:
		return c.convertTupleType(n)

	// 集合类型
	case *ast.SetType:
		return c.convertSetType(n)

	// 联合类型
	case *ast.UnionType:
		return c.convertUnionType(n)

	// 交集类型
	case *ast.IntersectionType:
		return c.convertIntersectionType(n)

	// 可空类型
	case *ast.NullableType:
		return c.convertNullableType(n)

	// 条件类型
	case *ast.ConditionalType:
		return c.convertConditionalType(n)

	// 类型字面量
	case *ast.TypeLiteral:
		return c.convertTypeLiteral(n)

	// 基本字面量类型
	case *ast.AnyLiteral:
		return anyType, nil
	case *ast.NumericLiteral:
		return numberType, nil
	case *ast.StringLiteral:
		return stringType, nil
	case *ast.TrueLiteral, *ast.FalseLiteral:
		return boolType, nil

	default:
		return nil, fmt.Errorf("unsupported type: %T", n)
	}
}

// ==================== 值转换辅助方法 ====================

func (c *ASTConverter) convertArrayLiteral(expr *ast.ArrayLiteralExpr) (Value, error) {
	if len(expr.Elems) == 0 {
		// 空数组，使用any类型
		return NewArrayValue(anyType), nil
	}

	// 转换第一个元素来确定类型
	firstElem, err := c.ConvertValue(expr.Elems[0])
	if err != nil {
		return nil, fmt.Errorf("failed to convert first array element: %w", err)
	}

	elementType := firstElem.Type()
	elements := make([]Value, len(expr.Elems))
	elements[0] = firstElem

	// 转换其余元素
	for i := 1; i < len(expr.Elems); i++ {
		elem, err := c.ConvertValue(expr.Elems[i])
		if err != nil {
			return nil, fmt.Errorf("failed to convert array element %d: %w", i, err)
		}

		// 检查类型兼容性
		if !elem.Type().AssignableTo(elementType) {
			// 如果类型不兼容，将整个数组推断为any[]
			elementType = anyType
			// 重新转换所有元素为any类型
			for j := range i {
				elements[j] = c.convertToAnyType(elements[j])
			}
		}
		elements[i] = elem
	}

	return NewArrayValue(elementType, elements...), nil
}

// convertToAnyType 将值转换为any类型（如果还不是any类型）
func (c *ASTConverter) convertToAnyType(value Value) Value {
	if value.Type().Name() == "any" {
		return value
	}
	// 这里简化处理，实际应该创建一个包装器
	return value
}

// convertObjectLiteral converts an object literal to a Value
func (c *ASTConverter) convertObjectLiteral(expr *ast.ObjectLiteralExpr) (Value, error) {
	if len(expr.Props) == 0 {
		// empty object
		return NewMapValue(anyType, anyType), nil
	}

	// determine key and value types
	var keyType, valueType Type
	pairs := make([]*KeyValuePair, 0, len(expr.Props))

	for i, prop := range expr.Props {
		keyValueExpr, ok := prop.(*ast.KeyValueExpr)
		if !ok {
			return nil, fmt.Errorf("invalid property in object literal: %T", prop)
		}

		// convert key
		key, err := c.ConvertValue(keyValueExpr.Key)
		if err != nil {
			return nil, fmt.Errorf("failed to convert object key: %w", err)
		}

		// convert value
		value, err := c.ConvertValue(keyValueExpr.Value)
		if err != nil {
			return nil, fmt.Errorf("failed to convert object value: %w", err)
		}

		// determine types
		if i == 0 {
			keyType = key.Type()
			valueType = value.Type()
		} else {
			// 检查类型兼容性
			if !key.Type().AssignableTo(keyType) {
				keyType = anyType
				// 重新转换所有键为any类型
				for j := 0; j < i; j++ {
					pairs[j].Key = c.convertToAnyType(pairs[j].Key)
				}
			}
			if !value.Type().AssignableTo(valueType) {
				valueType = anyType
				// 重新转换所有值为any类型
				for j := 0; j < i; j++ {
					pairs[j].Value = c.convertToAnyType(pairs[j].Value)
				}
			}
		}

		pairs = append(pairs, &KeyValuePair{Key: key, Value: value})
	}

	mapValue := NewMapValue(keyType, valueType)
	for _, pair := range pairs {
		mapValue.Set(pair.Key, pair.Value)
	}

	return mapValue, nil
}

func (c *ASTConverter) convertNamedObjectLiteral(expr *ast.NamedObjectLiteralExpr) (Value, error) {
	// 这里需要根据类型名创建相应的类实例
	// 简化实现，创建通用对象
	return c.convertObjectLiteral(&ast.ObjectLiteralExpr{
		Props: expr.Props,
	})
}

func (c *ASTConverter) convertSetLiteral(expr *ast.SetLiteralExpr) (Value, error) {
	if len(expr.Elems) == 0 {
		// 空集合，使用any类型
		return NewSetValue(NewSetType(anyType)), nil
	}

	// 转换第一个元素来确定类型
	firstElem, err := c.ConvertValue(expr.Elems[0])
	if err != nil {
		return nil, fmt.Errorf("failed to convert first set element: %w", err)
	}

	elementType := firstElem.Type()
	elements := make([]Value, len(expr.Elems))
	elements[0] = firstElem

	// 转换其余元素
	for i := 1; i < len(expr.Elems); i++ {
		elem, err := c.ConvertValue(expr.Elems[i])
		if err != nil {
			return nil, fmt.Errorf("failed to convert set element %d: %w", i, err)
		}

		// 检查类型兼容性
		if !elem.Type().AssignableTo(elementType) {
			return nil, fmt.Errorf("incompatible types in set: %s vs %s",
				elem.Type().Name(), elementType.Name())
		}

		elements[i] = elem
	}

	return NewSetValue(NewSetType(elementType), elements...), nil
}

func (c *ASTConverter) convertTupleLiteral(expr *ast.TupleLiteralExpr) (Value, error) {
	if len(expr.Elems) == 0 {
		// 空元组
		return NewTupleValue(NewTupleType()), nil
	}

	// 转换所有元素
	elements := make([]Value, len(expr.Elems))
	elementTypes := make([]Type, len(expr.Elems))

	for i, elem := range expr.Elems {
		value, err := c.ConvertValue(elem)
		if err != nil {
			return nil, fmt.Errorf("failed to convert tuple element %d: %w", i, err)
		}
		elements[i] = value
		elementTypes[i] = value.Type()
	}

	tupleType := NewTupleType(elementTypes...)
	return NewTupleValue(tupleType, elements...), nil
}

func (c *ASTConverter) convertNamedTupleLiteral(expr *ast.NamedTupleLiteralExpr) (Value, error) {
	// 命名元组需要字段名和类型信息
	// 简化实现，转换为普通元组
	return c.convertTupleLiteral(&ast.TupleLiteralExpr{
		Elems: expr.Elems,
	})
}

func (c *ASTConverter) convertIdentifier(ident *ast.Ident) (Value, error) {
	// 对于对象字面量的键，直接转换为字符串值
	// 对于其他标识符，需要符号表支持
	return &StringValue{Value: ident.Name}, nil
}

// ==================== 类型转换辅助方法 ====================

func (c *ASTConverter) convertBasicType(ident *ast.Ident) (Type, error) {
	switch strings.ToLower(ident.Name) {
	case "str", "string":
		return stringType, nil
	case "num", "number":
		return numberType, nil
	case "bool", "boolean":
		return boolType, nil
	case "any":
		return anyType, nil
	case "void":
		return voidType, nil
	case "null":
		return nullType, nil
	case "error":
		return errorType, nil
	default:
		// 可能是用户定义的类型，需要从类型注册表中查找
		// 简化实现，返回anyType
		return anyType, nil
	}
}

func (c *ASTConverter) convertTypeReference(ref *ast.TypeReference) (Type, error) {
	// 获取类型名称
	var typeName string
	if ident, ok := ref.Name.(*ast.Ident); ok {
		typeName = strings.ToLower(ident.Name)
	} else {
		return nil, fmt.Errorf("invalid type name: %T", ref.Name)
	}

	// 处理基础类型
	switch typeName {
	case "str", "string":
		return stringType, nil
	case "num", "number":
		return numberType, nil
	case "bool", "boolean":
		return boolType, nil
	case "any":
		return anyType, nil
	case "void":
		return voidType, nil
	case "null":
		return nullType, nil
	case "error":
		return errorType, nil
	}

	// 处理特殊类型
	switch typeName {
	case "map":
		// Map 类型需要两个类型参数
		if len(ref.TypeParams) == 2 {
			keyType, err := c.ConvertType(ref.TypeParams[0])
			if err != nil {
				return nil, fmt.Errorf("failed to convert map key type: %w", err)
			}
			valueType, err := c.ConvertType(ref.TypeParams[1])
			if err != nil {
				return nil, fmt.Errorf("failed to convert map value type: %w", err)
			}
			return NewMapType(keyType, valueType), nil
		}
		return nil, fmt.Errorf("map type requires exactly 2 type parameters")

	case "array":
		// Array 类型需要一个类型参数
		if len(ref.TypeParams) == 1 {
			elementType, err := c.ConvertType(ref.TypeParams[0])
			if err != nil {
				return nil, fmt.Errorf("failed to convert array element type: %w", err)
			}
			return NewArrayType(elementType), nil
		}
		return nil, fmt.Errorf("array type requires exactly 1 type parameter")

	case "set":
		// Set 类型需要一个类型参数
		if len(ref.TypeParams) == 1 {
			elementType, err := c.ConvertType(ref.TypeParams[0])
			if err != nil {
				return nil, fmt.Errorf("failed to convert set element type: %w", err)
			}
			return NewSetType(elementType), nil
		}
		return nil, fmt.Errorf("set type requires exactly 1 type parameter")

	case "tuple":
		// Tuple 类型可以有多个类型参数
		if len(ref.TypeParams) == 0 {
			return NewTupleType(), nil
		}

		elementTypes := make([]Type, len(ref.TypeParams))
		for i, param := range ref.TypeParams {
			elementType, err := c.ConvertType(param)
			if err != nil {
				return nil, fmt.Errorf("failed to convert tuple element type %d: %w", i, err)
			}
			elementTypes[i] = elementType
		}
		return NewTupleType(elementTypes...), nil

	case "function":
		// Function 类型需要参数类型和返回类型
		if len(ref.TypeParams) < 2 {
			return nil, fmt.Errorf("function type requires at least 2 type parameters (params and return type)")
		}

		// 最后一个参数是返回类型
		returnType, err := c.ConvertType(ref.TypeParams[len(ref.TypeParams)-1])
		if err != nil {
			return nil, fmt.Errorf("failed to convert function return type: %w", err)
		}

		// 前面的参数是参数类型
		paramTypes := make([]Type, len(ref.TypeParams)-1)
		for i, param := range ref.TypeParams[:len(ref.TypeParams)-1] {
			paramType, err := c.ConvertType(param)
			if err != nil {
				return nil, fmt.Errorf("failed to convert function parameter type %d: %w", i, err)
			}
			paramTypes[i] = paramType
		}

		// 创建参数列表
		positionalParams := make([]*Parameter, len(paramTypes))
		for i, paramType := range paramTypes {
			positionalParams[i] = &Parameter{
				name:     fmt.Sprintf("param%d", i),
				typ:      paramType,
				optional: false,
				variadic: false,
				isNamed:  false,
				required: true,
			}
		}

		// 创建函数类型
		return &FunctionType{
			name: "function",
			signatures: []*FunctionSignature{
				{
					positionalParams: positionalParams,
					returnType:       returnType,
				},
			},
		}, nil

	case "union":
		// Union 类型可以有多个类型参数
		if len(ref.TypeParams) == 0 {
			return anyType, nil
		}

		types := make([]Type, len(ref.TypeParams))
		for i, param := range ref.TypeParams {
			convertedType, err := c.ConvertType(param)
			if err != nil {
				return nil, fmt.Errorf("failed to convert union type %d: %w", i, err)
			}
			types[i] = convertedType
		}
		return NewUnionType(types...), nil

	case "intersection":
		// Intersection 类型可以有多个类型参数
		if len(ref.TypeParams) == 0 {
			return anyType, nil
		}

		types := make([]Type, len(ref.TypeParams))
		for i, param := range ref.TypeParams {
			convertedType, err := c.ConvertType(param)
			if err != nil {
				return nil, fmt.Errorf("failed to convert intersection type %d: %w", i, err)
			}
			types[i] = convertedType
		}
		return NewIntersectionType(types...), nil

	case "nullable":
		// Nullable 类型需要一个类型参数
		if len(ref.TypeParams) == 1 {
			baseType, err := c.ConvertType(ref.TypeParams[0])
			if err != nil {
				return nil, fmt.Errorf("failed to convert nullable base type: %w", err)
			}
			return NewNullableType(baseType), nil
		}
		return nil, fmt.Errorf("nullable type requires exactly 1 type parameter")
	}

	// 处理用户定义的类型（可能是泛型类）
	// 这里需要从类型注册表中查找
	// 简化实现，返回anyType
	return anyType, nil
}

func (c *ASTConverter) convertArrayType(arrayType *ast.ArrayType) (Type, error) {
	elementType, err := c.ConvertType(arrayType.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to convert array element type: %w", err)
	}
	return NewArrayType(elementType), nil
}

func (c *ASTConverter) convertFunctionType(funcType *ast.FunctionType) (Type, error) {
	// 转换参数类型
	paramTypes := make([]Type, len(funcType.Recv))
	for i, param := range funcType.Recv {
		paramType, err := c.ConvertType(param)
		if err != nil {
			return nil, fmt.Errorf("failed to convert function parameter type %d: %w", i, err)
		}
		paramTypes[i] = paramType
	}

	// 转换返回类型
	returnType, err := c.ConvertType(funcType.RetVal)
	if err != nil {
		return nil, fmt.Errorf("failed to convert function return type: %w", err)
	}

	// 创建参数列表
	positionalParams := make([]*Parameter, len(paramTypes))
	for i, paramType := range paramTypes {
		positionalParams[i] = &Parameter{
			name:     fmt.Sprintf("param%d", i),
			typ:      paramType,
			optional: false,
			variadic: false,
			isNamed:  false,
			required: true,
		}
	}

	// 创建函数类型
	return &FunctionType{
		name: "function",
		signatures: []*FunctionSignature{
			{
				positionalParams: positionalParams,
				returnType:       returnType,
			},
		},
	}, nil
}

func (c *ASTConverter) convertTupleType(tupleType *ast.TupleType) (Type, error) {
	elementTypes := make([]Type, len(tupleType.Types))
	for i, elemType := range tupleType.Types {
		convertedType, err := c.ConvertType(elemType)
		if err != nil {
			return nil, fmt.Errorf("failed to convert tuple element type %d: %w", i, err)
		}
		elementTypes[i] = convertedType
	}
	return NewTupleType(elementTypes...), nil
}

func (c *ASTConverter) convertSetType(setType *ast.SetType) (Type, error) {
	if len(setType.Types) == 0 {
		return NewSetType(anyType), nil
	}

	elementType, err := c.ConvertType(setType.Types[0])
	if err != nil {
		return nil, fmt.Errorf("failed to convert set element type: %w", err)
	}
	return NewSetType(elementType), nil
}

func (c *ASTConverter) convertUnionType(unionType *ast.UnionType) (Type, error) {
	// 联合类型转换
	if len(unionType.Types) == 0 {
		return anyType, nil
	}

	// 转换所有类型
	types := make([]Type, len(unionType.Types))
	for i, t := range unionType.Types {
		convertedType, err := c.ConvertType(t)
		if err != nil {
			return nil, fmt.Errorf("failed to convert union type %d: %w", i, err)
		}
		types[i] = convertedType
	}

	return NewUnionType(types...), nil
}

func (c *ASTConverter) convertIntersectionType(intersectionType *ast.IntersectionType) (Type, error) {
	// 交集类型转换
	if len(intersectionType.Types) == 0 {
		return anyType, nil
	}

	// 转换所有类型
	types := make([]Type, len(intersectionType.Types))
	for i, t := range intersectionType.Types {
		convertedType, err := c.ConvertType(t)
		if err != nil {
			return nil, fmt.Errorf("failed to convert intersection type %d: %w", i, err)
		}
		types[i] = convertedType
	}

	return NewIntersectionType(types...), nil
}

func (c *ASTConverter) convertNullableType(nullableType *ast.NullableType) (Type, error) {
	baseType, err := c.ConvertType(nullableType.X)
	if err != nil {
		return nil, err
	}
	return NewNullableType(baseType), nil
}

func (c *ASTConverter) convertConditionalType(conditionalType *ast.ConditionalType) (Type, error) {
	// 条件类型转换
	// 简化实现，返回真值类型
	return c.ConvertType(conditionalType.TrueType)
}

func (c *ASTConverter) convertTypeLiteral(typeLiteral *ast.TypeLiteral) (Type, error) {
	obj := NewObjectType("{}", O_OBJ)
	for i, member := range typeLiteral.Members {
		field, err := c.ConvertType(member)
		if err != nil {
			return nil, fmt.Errorf("failed to convert type literal member %d: %w", i, err)
		}
		obj.AddField(field.Name(), &Field{ft: field})
	}
	return obj, nil
}

// ==================== 工具方法 ====================

// ConvertModifiers 转换修饰符
func (c *ASTConverter) ConvertModifiers(modifiers []ast.Modifier) []FieldModifier {
	var fieldModifiers []FieldModifier
	for _, mod := range modifiers {
		switch mod.Kind() {
		case ast.ModKindRequired:
			fieldModifiers = append(fieldModifiers, FieldModifierRequired)
		case ast.ModKindReadonly:
			fieldModifiers = append(fieldModifiers, FieldModifierReadonly)
		case ast.ModKindStatic:
			fieldModifiers = append(fieldModifiers, FieldModifierStatic)
		case ast.ModKindPub:
			fieldModifiers = append(fieldModifiers, FieldModifierPub)
		case ast.ModKindConst:
			fieldModifiers = append(fieldModifiers, FieldModifierReadonly)
		case ast.ModKindFinal:
			fieldModifiers = append(fieldModifiers, FieldModifierRequired)
		}
	}
	return fieldModifiers
}

// ConvertField 转换字段
func (c *ASTConverter) ConvertField(field *ast.Field) (*Field, error) {
	fieldType, err := c.ConvertType(field.Type)
	if err != nil {
		return nil, fmt.Errorf("failed to convert field type: %w", err)
	}

	modifiers := c.ConvertModifiers(field.Modifiers)

	result := &Field{
		ft:   fieldType,
		mods: modifiers,
	}

	// 转换默认值
	if field.Value != nil {
		defaultValue, err := c.ConvertValue(field.Value)
		if err != nil {
			return nil, fmt.Errorf("failed to convert field default value: %w", err)
		}
		result.SetDefaultValue(defaultValue)
	}

	return result, nil
}

// ConvertParameter 转换参数
func (c *ASTConverter) ConvertParameter(param *ast.Parameter) (*Parameter, error) {
	paramType, err := c.ConvertType(param.Type)
	if err != nil {
		return nil, fmt.Errorf("failed to convert parameter type: %w", err)
	}

	result := &Parameter{
		name:     param.Name.Name,
		typ:      paramType,
		optional: false, // 根据修饰符确定
		variadic: false, // 根据修饰符确定
		isNamed:  false, // 根据上下文确定
		required: false, // 根据修饰符确定
	}

	// 处理修饰符
	switch param.Modifier.Kind() {
	case ast.ModKindRequired:
		result.required = true
	case ast.ModKindEllipsis:
		result.variadic = true
		result.optional = true
	}

	return result, nil
}
