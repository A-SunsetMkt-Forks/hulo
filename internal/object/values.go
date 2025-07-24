// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package object

import (
	"fmt"
	"math/big"
	"slices"
	"strconv"
	"strings"
)

var (
	NULL  = &NullValue{}
	TRUE  = &BoolValue{Value: true}
	FALSE = &BoolValue{Value: false}
)

// NullValue represents the null value.
type NullValue struct{}

func (n *NullValue) Text() string {
	return "null"
}

func (n *NullValue) Interface() any {
	return nil
}

func (n *NullValue) Type() Type {
	// TODO: return voidType
	return nil
}

// NumberValue represents the number value.
type NumberValue struct {
	// TODO 根据精度选择存储模型
	Value *big.Float
}

// NewNumberValue creates a new number value from the given string.
func NewNumberValue(value string) *NumberValue {
	num, _ := strconv.ParseFloat(value, 64)
	return &NumberValue{Value: big.NewFloat(num)}
}

func (n *NumberValue) Text() string {
	return n.Value.String()
}

func (n *NumberValue) Interface() any {
	return n.Value
}

func (n *NumberValue) Type() Type {
	return numberType
}

func (n *NumberValue) Clone() *NumberValue {
	return &NumberValue{Value: new(big.Float).Set(n.Value)}
}

// BoolValue represents the boolean value.
type BoolValue struct {
	Value bool
}

func (b *BoolValue) Text() string {
	return strconv.FormatBool(b.Value)
}

func (b *BoolValue) Interface() any {
	return b.Value
}

func (b *BoolValue) Type() Type {
	return boolType
}

func (b *BoolValue) Clone() *BoolValue {
	return &BoolValue{Value: b.Value}
}

// ErrorValue represents the error value.
type ErrorValue struct {
	Value string
}

func (e *ErrorValue) Text() string {
	return e.Value
}

func (e *ErrorValue) Interface() any {
	return e.Value
}

func (e *ErrorValue) Type() Type {
	return errorType
}

// StringValue represents the string value.
type StringValue struct {
	Value string
}

func (s *StringValue) Text() string {
	return s.Value
}

func (s *StringValue) Interface() any {
	return s.Value
}

func (s *StringValue) Type() Type {
	return stringType
}

func (s *StringValue) Clone() *StringValue {
	return &StringValue{Value: s.Value}
}

// ==================== 数组值类型 ====================

// ArrayValue represents an array value.
type ArrayValue struct {
	Elements  []Value
	arrayType *ArrayType
}

// NewArrayValue creates a new array value with the given elements.
func NewArrayValue(elementType Type, elements ...Value) *ArrayValue {
	return &ArrayValue{
		Elements:  elements,
		arrayType: NewArrayType(elementType),
	}
}

func (a *ArrayValue) Text() string {
	if len(a.Elements) == 0 {
		return "[]"
	}

	parts := make([]string, len(a.Elements))
	for i, elem := range a.Elements {
		parts[i] = elem.Text()
	}

	return fmt.Sprintf("[%s]", strings.Join(parts, ", "))
}

func (a *ArrayValue) Interface() any {
	return a.Elements
}

func (a *ArrayValue) Type() Type {
	return a.arrayType
}

// Get returns the element at the given index.
func (a *ArrayValue) Get(index int) Value {
	if index < 0 || index >= len(a.Elements) {
		return NULL
	}
	return a.Elements[index]
}

// Set sets the element at the given index.
func (a *ArrayValue) Set(index int, value Value) {
	if index < 0 {
		return
	}

	// 扩展数组如果需要
	for len(a.Elements) <= index {
		a.Elements = append(a.Elements, NULL)
	}

	a.Elements[index] = value
}

// Append adds an element to the end of the array.
func (a *ArrayValue) Append(value Value) {
	a.Elements = append(a.Elements, value)
}

// Length returns the length of the array.
func (a *ArrayValue) Length() int {
	return len(a.Elements)
}

// ==================== Map值类型 ====================

// KeyValuePair represents a key-value pair in a map.
type KeyValuePair struct {
	Key   Value
	Value Value
}

// MapValue represents a map value.
type MapValue struct {
	Pairs   []*KeyValuePair
	mapType *MapType
}

// NewMapValue creates a new map value with the given key and value types.
func NewMapValue(keyType, valueType Type) *MapValue {
	return &MapValue{
		Pairs:   []*KeyValuePair{},
		mapType: NewMapType(keyType, valueType),
	}
}

func (m *MapValue) Text() string {
	if len(m.Pairs) == 0 {
		return "{}"
	}

	pairs := make([]string, len(m.Pairs))
	for i, pair := range m.Pairs {
		pairs[i] = fmt.Sprintf("%s: %s", pair.Key.Text(), pair.Value.Text())
	}

	return fmt.Sprintf("{%s}", strings.Join(pairs, ", "))
}

func (m *MapValue) Interface() any {
	result := make(map[string]any)
	for _, pair := range m.Pairs {
		result[pair.Key.Text()] = pair.Value.Interface()
	}
	return result
}

func (m *MapValue) Type() Type {
	return m.mapType
}

// Get returns the value associated with the given key.
func (m *MapValue) Get(key Value) Value {
	for _, pair := range m.Pairs {
		if pair.Key.Text() == key.Text() {
			return pair.Value
		}
	}
	return NULL
}

// Set sets the value associated with the given key.
func (m *MapValue) Set(key, value Value) {
	// 查找是否已存在该key
	for _, pair := range m.Pairs {
		if pair.Key.Text() == key.Text() {
			pair.Value = value
			return
		}
	}

	// 添加新的键值对
	m.Pairs = append(m.Pairs, &KeyValuePair{
		Key:   key,
		Value: value,
	})
}

// Delete removes the key-value pair with the given key.
func (m *MapValue) Delete(key Value) {
	for i, pair := range m.Pairs {
		if pair.Key.Text() == key.Text() {
			m.Pairs = slices.Delete(m.Pairs, i, i+1)
			return
		}
	}
}

// Has returns true if the map contains the given key.
func (m *MapValue) Has(key Value) bool {
	for _, pair := range m.Pairs {
		if pair.Key.Text() == key.Text() {
			return true
		}
	}
	return false
}

// Keys returns all keys in the map.
func (m *MapValue) Keys() []Value {
	keys := make([]Value, len(m.Pairs))
	for i, pair := range m.Pairs {
		keys[i] = pair.Key
	}
	return keys
}

// Values returns all values in the map.
func (m *MapValue) Values() []Value {
	values := make([]Value, len(m.Pairs))
	for i, pair := range m.Pairs {
		values[i] = pair.Value
	}
	return values
}

// Length returns the number of key-value pairs in the map.
func (m *MapValue) Length() int {
	return len(m.Pairs)
}

// Clear removes all key-value pairs from the map.
func (m *MapValue) Clear() {
	m.Pairs = []*KeyValuePair{}
}

// ==================== 函数值类型 ====================

// FunctionValue represents a function value.
type FunctionValue struct {
	functionType *FunctionType
}

// NewFunctionValue creates a new function value.
func NewFunctionValue(functionType *FunctionType) *FunctionValue {
	return &FunctionValue{
		functionType: functionType,
	}
}

func (f *FunctionValue) Text() string {
	return f.functionType.Text()
}

func (f *FunctionValue) Interface() any {
	return f.functionType
}

func (f *FunctionValue) Type() Type {
	return f.functionType
}

// Call calls the function with the given arguments.
func (f *FunctionValue) Call(args []Value, namedArgs map[string]Value) (Value, error) {
	return f.functionType.Call(args, namedArgs)
}

// TupleValue 表示元组值
type TupleValue struct {
	tupleType *TupleType
	elements  []Value
}

func NewTupleValue(tupleType *TupleType, elements ...Value) *TupleValue {
	if len(elements) != len(tupleType.elementTypes) {
		panic(fmt.Sprintf("tuple element count mismatch: expected %d, got %d",
			len(tupleType.elementTypes), len(elements)))
	}

	return &TupleValue{
		tupleType: tupleType,
		elements:  elements,
	}
}

func (tv *TupleValue) Type() Type {
	return tv.tupleType
}

func (tv *TupleValue) Text() string {
	if tv.tupleType.isNamed() {
		return tv.namedTupleString()
	}
	return tv.anonymousTupleString()
}

func (tv *TupleValue) Interface() any {
	return tv.elements
}

func (tv *TupleValue) Elements() []Value {
	return tv.elements
}

func (tv *TupleValue) Length() int {
	return len(tv.elements)
}

// Get 通过索引获取元素
func (tv *TupleValue) Get(index int) Value {
	if index < 0 || index >= len(tv.elements) {
		return NULL
	}
	return tv.elements[index]
}

// GetByName 通过字段名获取元素
func (tv *TupleValue) GetByName(fieldName string) Value {
	index := tv.tupleType.GetFieldIndex(fieldName)
	if index == -1 {
		return NULL
	}
	return tv.elements[index]
}

// Set 设置指定位置的元素
func (tv *TupleValue) Set(index int, value Value) error {
	if index < 0 || index >= len(tv.elements) {
		return fmt.Errorf("tuple index out of bounds: %d", index)
	}

	expectedType := tv.tupleType.GetElementType(index)
	if !value.Type().AssignableTo(expectedType) {
		return fmt.Errorf("cannot assign %s to %s", value.Type().Name(), expectedType.Name())
	}

	tv.elements[index] = value
	return nil
}

// With 创建新的元组，更新指定位置的元素
func (tv *TupleValue) With(index int, value Value) (*TupleValue, error) {
	if index < 0 || index >= len(tv.elements) {
		return nil, fmt.Errorf("tuple index out of bounds: %d", index)
	}

	expectedType := tv.tupleType.GetElementType(index)
	if !value.Type().AssignableTo(expectedType) {
		return nil, fmt.Errorf("cannot assign %s to %s", value.Type().Name(), expectedType.Name())
	}

	newElements := make([]Value, len(tv.elements))
	copy(newElements, tv.elements)
	newElements[index] = value

	return NewTupleValue(tv.tupleType, newElements...), nil
}

// Destructure 解构元组
func (tv *TupleValue) Destructure() []Value {
	result := make([]Value, len(tv.elements))
	copy(result, tv.elements)
	return result
}

func (tv *TupleValue) anonymousTupleString() string {
	parts := make([]string, len(tv.elements))
	for i, elem := range tv.elements {
		parts[i] = elem.Text()
	}
	return fmt.Sprintf("[%s]", strings.Join(parts, ", "))
}

func (tv *TupleValue) namedTupleString() string {
	parts := make([]string, len(tv.elements))
	for i, elem := range tv.elements {
		if tv.tupleType.fieldNames[i] != "" {
			parts[i] = fmt.Sprintf("%s: %s", tv.tupleType.fieldNames[i], elem.Text())
		} else {
			parts[i] = elem.Text()
		}
	}
	return fmt.Sprintf("(%s)", strings.Join(parts, ", "))
}

// SetValue 表示集合值
type SetValue struct {
	setType  *SetType
	elements map[string]Value // 使用字符串键来保证唯一性
}

func NewSetValue(setType *SetType, elements ...Value) *SetValue {
	set := &SetValue{
		setType:  setType,
		elements: make(map[string]Value),
	}

	// 添加初始元素
	for _, elem := range elements {
		set.Add(elem)
	}

	return set
}

func (sv *SetValue) Type() Type {
	return sv.setType
}

func (sv *SetValue) Text() string {
	if len(sv.elements) == 0 {
		return "{}"
	}

	parts := make([]string, 0, len(sv.elements))
	for _, elem := range sv.elements {
		parts = append(parts, elem.Text())
	}
	return fmt.Sprintf("{%s}", strings.Join(parts, ", "))
}

func (sv *SetValue) Interface() any {
	return sv.elements
}

func (sv *SetValue) Elements() map[string]Value {
	return sv.elements
}

func (sv *SetValue) Length() int {
	return len(sv.elements)
}

func (sv *SetValue) IsEmpty() bool {
	return len(sv.elements) == 0
}

// Add 添加元素
func (sv *SetValue) Add(element Value) error {
	// 类型检查
	if !element.Type().AssignableTo(sv.setType.elementType) {
		return fmt.Errorf("cannot add %s to set<%s>",
			element.Type().Name(), sv.setType.elementType.Name())
	}

	key := element.Text()
	sv.elements[key] = element
	return nil
}

// Remove 删除元素
func (sv *SetValue) Remove(element Value) error {
	key := element.Text()
	if _, exists := sv.elements[key]; !exists {
		return fmt.Errorf("element %s not found in set", element.Text())
	}

	delete(sv.elements, key)
	return nil
}

// Has 检查元素是否存在
func (sv *SetValue) Has(element Value) bool {
	key := element.Text()
	_, exists := sv.elements[key]
	return exists
}

// Contains 类型安全的包含检查
func (sv *SetValue) Contains(element Value) (bool, error) {
	// 类型检查
	if !element.Type().AssignableTo(sv.setType.elementType) {
		return false, fmt.Errorf("cannot check %s in set<%s>",
			element.Type().Name(), sv.setType.elementType.Name())
	}

	return sv.Has(element), nil
}

// Union 并集运算
func (sv *SetValue) Union(other *SetValue) (*SetValue, error) {
	// 类型兼容性检查
	if !sv.setType.elementType.AssignableTo(other.setType.elementType) &&
		!other.setType.elementType.AssignableTo(sv.setType.elementType) {
		return nil, fmt.Errorf("cannot union sets with incompatible element types: %s and %s",
			sv.setType.elementType.Name(), other.setType.elementType.Name())
	}

	// 使用更通用的元素类型
	var resultElementType Type
	if sv.setType.elementType.AssignableTo(other.setType.elementType) {
		resultElementType = other.setType.elementType
	} else {
		resultElementType = sv.setType.elementType
	}

	resultSetType := NewSetType(resultElementType)
	result := NewSetValue(resultSetType)

	// 添加当前集合的所有元素
	for _, elem := range sv.elements {
		result.Add(elem)
	}

	// 添加另一个集合的所有元素
	for _, elem := range other.elements {
		result.Add(elem)
	}

	return result, nil
}

// Intersect 交集运算
func (sv *SetValue) Intersect(other *SetValue) (*SetValue, error) {
	// 类型兼容性检查
	if !sv.setType.elementType.AssignableTo(other.setType.elementType) &&
		!other.setType.elementType.AssignableTo(sv.setType.elementType) {
		return nil, fmt.Errorf("cannot intersect sets with incompatible element types: %s and %s",
			sv.setType.elementType.Name(), other.setType.elementType.Name())
	}

	// 使用更通用的元素类型
	var resultElementType Type
	if sv.setType.elementType.AssignableTo(other.setType.elementType) {
		resultElementType = other.setType.elementType
	} else {
		resultElementType = sv.setType.elementType
	}

	resultSetType := NewSetType(resultElementType)
	result := NewSetValue(resultSetType)

	// 只添加两个集合都有的元素
	for _, elem := range sv.elements {
		if other.Has(elem) {
			result.Add(elem)
		}
	}

	return result, nil
}

// Difference 差集运算
func (sv *SetValue) Difference(other *SetValue) (*SetValue, error) {
	// 类型兼容性检查
	if !sv.setType.elementType.AssignableTo(other.setType.elementType) &&
		!other.setType.elementType.AssignableTo(sv.setType.elementType) {
		return nil, fmt.Errorf("cannot compute difference of sets with incompatible element types: %s and %s",
			sv.setType.elementType.Name(), other.setType.elementType.Name())
	}

	result := NewSetValue(sv.setType)

	// 添加当前集合中不在另一个集合中的元素
	for _, elem := range sv.elements {
		if !other.Has(elem) {
			result.Add(elem)
		}
	}

	return result, nil
}

// Clear 清空集合
func (sv *SetValue) Clear() {
	sv.elements = make(map[string]Value)
}

// ToArray 转换为数组
func (sv *SetValue) ToArray() *ArrayValue {
	elements := make([]Value, 0, len(sv.elements))
	for _, elem := range sv.elements {
		elements = append(elements, elem)
	}

	return &ArrayValue{
		Elements:  elements,
		arrayType: NewArrayType(sv.setType.elementType),
	}
}

// EnumVariant 表示枚举变体
type EnumVariant struct {
	name      string
	value     Value // 变体的值
	enumType  *EnumType
	fields    map[string]Value // 关联字段值（Dart风格）
	adtFields map[string]Value // ADT字段值（Rust风格）
	adtTypes  map[string]Type  // ADT字段类型（Rust风格）
}

func NewEnumVariant(name string, value Value) *EnumVariant {
	return &EnumVariant{
		name:      name,
		value:     value,
		fields:    make(map[string]Value),
		adtFields: make(map[string]Value),
		adtTypes:  make(map[string]Type),
	}
}

func (ev *EnumVariant) Name() string {
	return ev.name
}

func (ev *EnumVariant) Value() Value {
	return ev.value
}

func (ev *EnumVariant) Type() Type {
	return ev.enumType
}

func (ev *EnumVariant) Text() string {
	if ev.enumType != nil && ev.enumType.isADT {
		return ev.adtString()
	}
	// Dart风格的关联值枚举
	if len(ev.fields) > 0 {
		return ev.associatedString()
	}
	return ev.name
}

func (ev *EnumVariant) Interface() any {
	return ev
}

// SetField 设置关联字段值（Dart风格）
func (ev *EnumVariant) SetField(name string, value Value) {
	ev.fields[name] = value
}

// GetField 获取关联字段值（Dart风格）
func (ev *EnumVariant) GetField(name string) Value {
	if value, exists := ev.fields[name]; exists {
		return value
	}
	return NULL
}

// SetADTField 设置ADT字段值（Rust风格）
func (ev *EnumVariant) SetADTField(name string, value Value) {
	ev.adtFields[name] = value
}

// GetADTField 获取ADT字段值（Rust风格）
func (ev *EnumVariant) GetADTField(name string) Value {
	if value, exists := ev.adtFields[name]; exists {
		return value
	}
	return NULL
}

// AddADTFieldType 添加ADT字段类型（Rust风格）
func (ev *EnumVariant) AddADTFieldType(name string, fieldType Type) {
	ev.adtTypes[name] = fieldType
}

// Dart风格的关联值字符串表示
func (ev *EnumVariant) associatedString() string {
	if len(ev.fields) == 0 {
		return ev.name
	}

	parts := make([]string, 0, len(ev.fields))
	for name, value := range ev.fields {
		parts = append(parts, fmt.Sprintf("%s: %s", name, value.Text()))
	}
	return fmt.Sprintf("%s(%s)", ev.name, strings.Join(parts, ", "))
}

// Rust风格的ADT字符串表示
func (ev *EnumVariant) adtString() string {
	if len(ev.adtFields) == 0 {
		return ev.name
	}

	parts := make([]string, 0, len(ev.adtFields))
	for name, value := range ev.adtFields {
		parts = append(parts, fmt.Sprintf("%s: %s", name, value.Text()))
	}
	return fmt.Sprintf("%s{%s}", ev.name, strings.Join(parts, ", "))
}
