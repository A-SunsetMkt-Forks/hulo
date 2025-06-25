// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package object

import (
	"fmt"
	"slices"
)

// TypeUtils 类型工具接口
type TypeUtils interface {
	// 基础类型工具
	Partial(t Type) (Type, error)
	Required(t Type) (Type, error)
	Optional(t Type) Type
	Readonly(t Type) (Type, error)
	Mutable(t Type) (Type, error)

	// 对象类型工具
	Pick(t Type, keys ...string) (Type, error)
	Omit(t Type, keys ...string) (Type, error)
	Record(keyType, valueType Type) Type

	// 联合类型工具
	Exclude(unionType, excludeType Type) (Type, error)
	Extract(unionType, extractType Type) (Type, error)
	NonNullable(t Type) (Type, error)

	// 条件类型工具
	If(condition bool, trueType, falseType Type) Type
	Infer(template Type, constraints map[string]Type) (Type, error)

	// 映射类型工具
	Map(t Type, mapper func(Type) Type) (Type, error)
	Filter(t Type, predicate func(Type) bool) (Type, error)

	// 高级类型工具
	DeepPartial(t Type) (Type, error)
	DeepRequired(t Type) (Type, error)
	DeepReadonly(t Type) (Type, error)
	DeepMutable(t Type) (Type, error)

	// 类型组合工具
	Merge(t1, t2 Type) (Type, error)
	Intersect(t1, t2 Type) (Type, error)
	Union(types ...Type) Type

	// 类型查询工具
	Keys(t Type) ([]string, error)
	Values(t Type) ([]Type, error)
	Entries(t Type) (map[string]Type, error)
}

// DefaultTypeUtils 默认类型工具实现
type DefaultTypeUtils struct{}

// NewTypeUtils 创建类型工具实例
func NewTypeUtils() TypeUtils {
	return &DefaultTypeUtils{}
}

// ==================== 基础类型工具 ====================

// Partial<T> - 使所有字段变为可选
func (tu *DefaultTypeUtils) Partial(t Type) (Type, error) {
	switch objType := t.(type) {
	case *ObjectType:
		newFields := make(map[string]*Field)
		for name, field := range objType.fields {
			newField := &Field{
				ft:   field.ft,
				mods: slices.Clone(field.mods),
			}
			// 移除 Required 修饰符，添加 Optional 修饰符
			newField.mods = removeModifier(newField.mods, FieldModifierRequired)
			newFields[name] = newField
		}
		return &ObjectType{
			name:    objType.name,
			kind:    objType.kind,
			methods: objType.methods,
			fields:  newFields,
		}, nil
	default:
		return nil, fmt.Errorf("Partial only works on object types, got %s", t.Name())
	}
}

// Required<T> - 使所有字段变为必需
func (tu *DefaultTypeUtils) Required(t Type) (Type, error) {
	switch objType := t.(type) {
	case *ObjectType:
		newFields := make(map[string]*Field)
		for name, field := range objType.fields {
			newField := &Field{
				ft:   field.ft,
				mods: slices.Clone(field.mods),
			}
			// 添加 Required 修饰符
			if !hasModifier(newField.mods, FieldModifierRequired) {
				newField.mods = append(newField.mods, FieldModifierRequired)
			}
			newFields[name] = newField
		}
		return &ObjectType{
			name:    objType.name,
			kind:    objType.kind,
			methods: objType.methods,
			fields:  newFields,
		}, nil
	default:
		return nil, fmt.Errorf("Required only works on object types, got %s", t.Name())
	}
}

// Optional<T> - 使类型变为可选（T | undefined）
func (tu *DefaultTypeUtils) Optional(t Type) Type {
	return NewUnionType(t, nullType)
}

// Readonly<T> - 使所有字段变为只读
func (tu *DefaultTypeUtils) Readonly(t Type) (Type, error) {
	switch objType := t.(type) {
	case *ObjectType:
		newFields := make(map[string]*Field)
		for name, field := range objType.fields {
			newField := &Field{
				ft:   field.ft,
				mods: slices.Clone(field.mods),
			}
			// 添加 Readonly 修饰符
			if !hasModifier(newField.mods, FieldModifierReadonly) {
				newField.mods = append(newField.mods, FieldModifierReadonly)
			}
			newFields[name] = newField
		}
		return &ObjectType{
			name:    objType.name,
			kind:    objType.kind,
			methods: objType.methods,
			fields:  newFields,
		}, nil
	default:
		return nil, fmt.Errorf("Readonly only works on object types, got %s", t.Name())
	}
}

// Mutable<T> - 使所有字段变为可变
func (tu *DefaultTypeUtils) Mutable(t Type) (Type, error) {
	switch objType := t.(type) {
	case *ObjectType:
		newFields := make(map[string]*Field)
		for name, field := range objType.fields {
			newField := &Field{
				ft:   field.ft,
				mods: slices.Clone(field.mods),
			}
			// 移除 Readonly 修饰符
			newField.mods = removeModifier(newField.mods, FieldModifierReadonly)
			newFields[name] = newField
		}
		return &ObjectType{
			name:    objType.name,
			kind:    objType.kind,
			methods: objType.methods,
			fields:  newFields,
		}, nil
	default:
		return nil, fmt.Errorf("Mutable only works on object types, got %s", t.Name())
	}
}

// ==================== 对象类型工具 ====================

// Pick<T, K> - 从对象类型中选择指定字段
func (tu *DefaultTypeUtils) Pick(t Type, keys ...string) (Type, error) {
	switch objType := t.(type) {
	case *ObjectType:
		newFields := make(map[string]*Field)
		for _, key := range keys {
			if field, exists := objType.fields[key]; exists {
				newFields[key] = field
			}
		}
		return &ObjectType{
			name:    objType.name,
			kind:    objType.kind,
			methods: make(map[string]Method),
			fields:  newFields,
		}, nil
	default:
		return nil, fmt.Errorf("Pick only works on object types, got %s", t.Name())
	}
}

// Omit<T, K> - 从对象类型中排除指定字段
func (tu *DefaultTypeUtils) Omit(t Type, keys ...string) (Type, error) {
	switch objType := t.(type) {
	case *ObjectType:
		newFields := make(map[string]*Field)
		keySet := make(map[string]bool)
		for _, key := range keys {
			keySet[key] = true
		}

		for name, field := range objType.fields {
			if !keySet[name] {
				newFields[name] = field
			}
		}
		return &ObjectType{
			name:    objType.name,
			kind:    objType.kind,
			methods: objType.methods,
			fields:  newFields,
		}, nil
	default:
		return nil, fmt.Errorf("Omit only works on object types, got %s", t.Name())
	}
}

// Record<K, V> - 创建键值对类型
func (tu *DefaultTypeUtils) Record(keyType, valueType Type) Type {
	return &ObjectType{
		name:    "Record",
		kind:    O_OBJ,
		methods: make(map[string]Method),
		fields: map[string]*Field{
			"[key: string]": {
				ft:   valueType,
				mods: []FieldModifier{FieldModifierRequired},
			},
		},
	}
}

// ==================== 联合类型工具 ====================

// Exclude<T, U> - 从联合类型中排除指定类型
func (tu *DefaultTypeUtils) Exclude(unionType, excludeType Type) (Type, error) {
	switch ut := unionType.(type) {
	case *UnionType:
		newTypes := make([]Type, 0)
		for _, t := range ut.types {
			if !t.AssignableTo(excludeType) {
				newTypes = append(newTypes, t)
			}
		}
		if len(newTypes) == 0 {
			return neverType, nil
		}
		if len(newTypes) == 1 {
			return newTypes[0], nil
		}
		return NewUnionType(newTypes...), nil
	default:
		return nil, fmt.Errorf("Exclude only works on union types, got %s", unionType.Name())
	}
}

// Extract<T, U> - 从联合类型中提取指定类型
func (tu *DefaultTypeUtils) Extract(unionType, extractType Type) (Type, error) {
	switch ut := unionType.(type) {
	case *UnionType:
		newTypes := make([]Type, 0)
		for _, t := range ut.types {
			if t.AssignableTo(extractType) {
				newTypes = append(newTypes, t)
			}
		}
		if len(newTypes) == 0 {
			return neverType, nil
		}
		if len(newTypes) == 1 {
			return newTypes[0], nil
		}
		return NewUnionType(newTypes...), nil
	default:
		return nil, fmt.Errorf("Extract only works on union types, got %s", unionType.Name())
	}
}

// NonNullable<T> - 排除 null 和 undefined
func (tu *DefaultTypeUtils) NonNullable(t Type) (Type, error) {
	return tu.Exclude(t, nullType)
}

// ==================== 条件类型工具 ====================

// If<C, T, F> - 条件类型
func (tu *DefaultTypeUtils) If(condition bool, trueType, falseType Type) Type {
	if condition {
		return trueType
	}
	return falseType
}

// Infer<T, C> - 类型推断
func (tu *DefaultTypeUtils) Infer(template Type, constraints map[string]Type) (Type, error) {
	// 这里需要实现类型参数替换逻辑
	// 简化实现：直接返回模板类型
	return template, nil
}

// ==================== 映射类型工具 ====================

// Map<T, F> - 映射类型
func (tu *DefaultTypeUtils) Map(t Type, mapper func(Type) Type) (Type, error) {
	switch objType := t.(type) {
	case *ObjectType:
		newFields := make(map[string]*Field)
		for name, field := range objType.fields {
			newType := mapper(field.ft)
			newFields[name] = &Field{
				ft:   newType,
				mods: slices.Clone(field.mods),
			}
		}
		return &ObjectType{
			name:    objType.name,
			kind:    objType.kind,
			methods: objType.methods,
			fields:  newFields,
		}, nil
	default:
		return nil, fmt.Errorf("Map only works on object types, got %s", t.Name())
	}
}

// Filter<T, P> - 过滤类型
func (tu *DefaultTypeUtils) Filter(t Type, predicate func(Type) bool) (Type, error) {
	switch objType := t.(type) {
	case *ObjectType:
		newFields := make(map[string]*Field)
		for name, field := range objType.fields {
			if predicate(field.ft) {
				newFields[name] = field
			}
		}
		return &ObjectType{
			name:    objType.name,
			kind:    objType.kind,
			methods: objType.methods,
			fields:  newFields,
		}, nil
	default:
		return nil, fmt.Errorf("Filter only works on object types, got %s", t.Name())
	}
}

// ==================== 高级类型工具 ====================

// DeepPartial<T> - 深度可选
func (tu *DefaultTypeUtils) DeepPartial(t Type) (Type, error) {
	switch objType := t.(type) {
	case *ObjectType:
		newFields := make(map[string]*Field)
		for name, field := range objType.fields {
			deepType, err := tu.DeepPartial(field.ft)
			if err != nil {
				deepType = field.ft
			}
			newField := &Field{
				ft:   deepType,
				mods: append([]FieldModifier{}, field.mods...),
			}
			// 移除 Required 修饰符
			newField.mods = removeModifier(newField.mods, FieldModifierRequired)
			newFields[name] = newField
		}
		return &ObjectType{
			name:    objType.name,
			kind:    objType.kind,
			methods: objType.methods,
			fields:  newFields,
		}, nil
	default:
		return t, nil
	}
}

// DeepRequired<T> - 深度必需
func (tu *DefaultTypeUtils) DeepRequired(t Type) (Type, error) {
	switch objType := t.(type) {
	case *ObjectType:
		newFields := make(map[string]*Field)
		for name, field := range objType.fields {
			deepType, err := tu.DeepRequired(field.ft)
			if err != nil {
				deepType = field.ft
			}
			newField := &Field{
				ft:   deepType,
				mods: append([]FieldModifier{}, field.mods...),
			}
			// 添加 Required 修饰符
			if !hasModifier(newField.mods, FieldModifierRequired) {
				newField.mods = append(newField.mods, FieldModifierRequired)
			}
			newFields[name] = newField
		}
		return &ObjectType{
			name:    objType.name,
			kind:    objType.kind,
			methods: objType.methods,
			fields:  newFields,
		}, nil
	default:
		return t, nil
	}
}

// DeepReadonly<T> - 深度只读
func (tu *DefaultTypeUtils) DeepReadonly(t Type) (Type, error) {
	switch objType := t.(type) {
	case *ObjectType:
		newFields := make(map[string]*Field)
		for name, field := range objType.fields {
			deepType, err := tu.DeepReadonly(field.ft)
			if err != nil {
				deepType = field.ft
			}
			newField := &Field{
				ft:   deepType,
				mods: append([]FieldModifier{}, field.mods...),
			}
			// 添加 Readonly 修饰符
			if !hasModifier(newField.mods, FieldModifierReadonly) {
				newField.mods = append(newField.mods, FieldModifierReadonly)
			}
			newFields[name] = newField
		}
		return &ObjectType{
			name:    objType.name,
			kind:    objType.kind,
			methods: objType.methods,
			fields:  newFields,
		}, nil
	default:
		return t, nil
	}
}

// DeepMutable<T> - 深度可变
func (tu *DefaultTypeUtils) DeepMutable(t Type) (Type, error) {
	switch objType := t.(type) {
	case *ObjectType:
		newFields := make(map[string]*Field)
		for name, field := range objType.fields {
			deepType, err := tu.DeepMutable(field.ft)
			if err != nil {
				deepType = field.ft
			}
			newField := &Field{
				ft:   deepType,
				mods: append([]FieldModifier{}, field.mods...),
			}
			// 移除 Readonly 修饰符
			newField.mods = removeModifier(newField.mods, FieldModifierReadonly)
			newFields[name] = newField
		}
		return &ObjectType{
			name:    objType.name,
			kind:    objType.kind,
			methods: objType.methods,
			fields:  newFields,
		}, nil
	default:
		return t, nil
	}
}

// ==================== 类型组合工具 ====================

// Merge<T1, T2> - 合并两个对象类型
func (tu *DefaultTypeUtils) Merge(t1, t2 Type) (Type, error) {
	obj1, ok1 := t1.(*ObjectType)
	obj2, ok2 := t2.(*ObjectType)

	if !ok1 || !ok2 {
		return nil, fmt.Errorf("Merge only works on object types")
	}

	newFields := make(map[string]*Field)

	// 复制第一个对象的字段
	for name, field := range obj1.fields {
		newFields[name] = field
	}

	// 合并第二个对象的字段（覆盖同名字段）
	for name, field := range obj2.fields {
		newFields[name] = field
	}

	return &ObjectType{
		name:    "Merged",
		kind:    O_OBJ,
		methods: obj1.methods, // 简化：只保留第一个对象的方法
		fields:  newFields,
	}, nil
}

// Intersect<T1, T2> - 类型交集
func (tu *DefaultTypeUtils) Intersect(t1, t2 Type) (Type, error) {
	// 简化实现：返回第一个类型
	// 实际实现需要复杂的类型交集计算
	return t1, nil
}

// Union - 创建联合类型
func (tu *DefaultTypeUtils) Union(types ...Type) Type {
	return NewUnionType(types...)
}

// ==================== 类型查询工具 ====================

// Keys<T> - 获取对象类型的键
func (tu *DefaultTypeUtils) Keys(t Type) ([]string, error) {
	switch objType := t.(type) {
	case *ObjectType:
		keys := make([]string, 0, len(objType.fields))
		for key := range objType.fields {
			keys = append(keys, key)
		}
		return keys, nil
	default:
		return nil, fmt.Errorf("Keys only works on object types, got %s", t.Name())
	}
}

// Values<T> - 获取对象类型的值
func (tu *DefaultTypeUtils) Values(t Type) ([]Type, error) {
	switch objType := t.(type) {
	case *ObjectType:
		values := make([]Type, 0, len(objType.fields))
		for _, field := range objType.fields {
			values = append(values, field.ft)
		}
		return values, nil
	default:
		return nil, fmt.Errorf("Values only works on object types, got %s", t.Name())
	}
}

// Entries<T> - 获取对象类型的键值对
func (tu *DefaultTypeUtils) Entries(t Type) (map[string]Type, error) {
	switch objType := t.(type) {
	case *ObjectType:
		entries := make(map[string]Type)
		for name, field := range objType.fields {
			entries[name] = field.ft
		}
		return entries, nil
	default:
		return nil, fmt.Errorf("Entries only works on object types, got %s", t.Name())
	}
}

// ==================== 辅助函数 ====================

// 辅助函数
func hasModifier(mods []FieldModifier, mod FieldModifier) bool {
	for _, m := range mods {
		if m == mod {
			return true
		}
	}
	return false
}

func removeModifier(mods []FieldModifier, mod FieldModifier) []FieldModifier {
	result := make([]FieldModifier, 0, len(mods))
	for _, m := range mods {
		if m != mod {
			result = append(result, m)
		}
	}
	return result
}

// 全局类型工具实例
var GlobalTypeUtils TypeUtils = NewTypeUtils()

// 便捷函数
func Partial(t Type) (Type, error)              { return GlobalTypeUtils.Partial(t) }
func Required(t Type) (Type, error)             { return GlobalTypeUtils.Required(t) }
func Optional(t Type) Type                      { return GlobalTypeUtils.Optional(t) }
func Readonly(t Type) (Type, error)             { return GlobalTypeUtils.Readonly(t) }
func Mutable(t Type) (Type, error)              { return GlobalTypeUtils.Mutable(t) }
func Pick(t Type, keys ...string) (Type, error) { return GlobalTypeUtils.Pick(t, keys...) }
func Omit(t Type, keys ...string) (Type, error) { return GlobalTypeUtils.Omit(t, keys...) }
func Record(keyType, valueType Type) Type       { return GlobalTypeUtils.Record(keyType, valueType) }
func Exclude(unionType, excludeType Type) (Type, error) {
	return GlobalTypeUtils.Exclude(unionType, excludeType)
}
func Extract(unionType, extractType Type) (Type, error) {
	return GlobalTypeUtils.Extract(unionType, extractType)
}
func NonNullable(t Type) (Type, error) { return GlobalTypeUtils.NonNullable(t) }
func If(condition bool, trueType, falseType Type) Type {
	return GlobalTypeUtils.If(condition, trueType, falseType)
}
func Infer(template Type, constraints map[string]Type) (Type, error) {
	return GlobalTypeUtils.Infer(template, constraints)
}
func Map(t Type, mapper func(Type) Type) (Type, error) { return GlobalTypeUtils.Map(t, mapper) }
func Filter(t Type, predicate func(Type) bool) (Type, error) {
	return GlobalTypeUtils.Filter(t, predicate)
}
func DeepPartial(t Type) (Type, error)        { return GlobalTypeUtils.DeepPartial(t) }
func DeepRequired(t Type) (Type, error)       { return GlobalTypeUtils.DeepRequired(t) }
func DeepReadonly(t Type) (Type, error)       { return GlobalTypeUtils.DeepReadonly(t) }
func DeepMutable(t Type) (Type, error)        { return GlobalTypeUtils.DeepMutable(t) }
func Merge(t1, t2 Type) (Type, error)         { return GlobalTypeUtils.Merge(t1, t2) }
func Intersect(t1, t2 Type) (Type, error)     { return GlobalTypeUtils.Intersect(t1, t2) }
func Union(types ...Type) Type                { return GlobalTypeUtils.Union(types...) }
func Keys(t Type) ([]string, error)           { return GlobalTypeUtils.Keys(t) }
func Values(t Type) ([]Type, error)           { return GlobalTypeUtils.Values(t) }
func Entries(t Type) (map[string]Type, error) { return GlobalTypeUtils.Entries(t) }
