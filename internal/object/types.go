// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package object

import (
	"fmt"
	"math/big"
	"strings"
	"sync"

	"slices"

	"github.com/hulo-lang/hulo/internal/core"
	"github.com/hulo-lang/hulo/syntax/hulo/ast"
	"github.com/hulo-lang/hulo/syntax/hulo/token"
)

var (
	anyType    = NewPrimitiveType("any", O_ANY)
	voidType   = NewPrimitiveType("void", O_VOID)
	neverType  = NewPrimitiveType("never", O_NEVER)
	nullType   = NewPrimitiveType("null", O_NULL)
	errorType  = NewErrorType()
	stringType = NewStringType()
	numberType = NewNumberType()
	boolType   = NewBoolType()
)

func GetAnyType() Type    { return anyType }
func GetVoidType() Type   { return voidType }
func GetNeverType() Type  { return neverType }
func GetNullType() Type   { return nullType }
func GetErrorType() Type  { return errorType }
func GetStringType() Type { return stringType }
func GetNumberType() Type { return numberType }
func GetBoolType() Type   { return boolType }

type StringType struct {
	*ObjectType
}

func NewStringType() *StringType {
	strType := &StringType{
		ObjectType: NewObjectType("str", O_STR),
	}

	strType.addStringMethods()

	return strType
}

func (s *StringType) addStringMethods() {
	s.AddMethod("length", &FunctionType{
		name: "length",
		signatures: []*FunctionSignature{
			{
				positionalParams: []*Parameter{
					{
						name: "s",
						typ:  s,
					},
				},
				returnType: numberType,
				builtin: func(args ...Value) Value {
					str := args[0]
					return &NumberValue{Value: big.NewFloat(float64(len(str.Text())))}
				},
			},
		},
	})
	s.AddMethod("substring", nil)
	s.AddMethod("split", nil)
	s.AddMethod("replace", nil)
	s.AddMethod("to_lower", nil)
	s.AddMethod("to_upper", nil)
	s.AddMethod("trim", nil)
	s.AddMethod("trim_left", nil)
	s.AddMethod("trim_right", nil)
}

type NumberType struct {
	*ObjectType
}

func NewNumberType() *NumberType {
	numType := &NumberType{
		ObjectType: NewObjectType("num", O_NUM),
	}

	numType.addNumberMethods()

	numType.addNumberFields()

	return numType
}

func (n *NumberType) addNumberMethods() {
	n.AddMethod("abs", nil)
	n.AddMethod("ceil", nil)
	n.AddMethod("floor", nil)
	n.AddMethod("round", nil)
	n.AddMethod("trunc", nil)
	n.AddMethod("pow", nil)
	n.AddMethod("sqrt", nil)
	n.AddMethod("cbrt", nil)
	n.AddMethod("exp", nil)
	n.AddMethod("log", nil)
	n.AddMethod("log10", nil)
	n.AddMethod("log2", nil)
	n.AddMethod("sin", nil)
	n.AddMethod("cos", nil)
	n.AddMethod("tan", nil)
	n.AddMethod("asin", nil)
	n.AddMethod("acos", nil)
	n.AddMethod("atan", nil)
	n.AddMethod("sinh", nil)
	n.AddMethod("cosh", nil)
	n.AddMethod("tanh", nil)
	n.AddMethod("asinh", nil)
	n.AddMethod("acosh", nil)
	n.AddMethod("atanh", nil)
	n.AddMethod("to_str", nil)
}

func (n *NumberType) addNumberFields() {

}

type BoolType struct {
	*ObjectType
}

func NewBoolType() *BoolType {
	boolType := &BoolType{
		ObjectType: NewObjectType("bool", O_BOOL),
	}

	boolType.addBoolMethods()

	return boolType
}

func (b *BoolType) addBoolMethods() {
	b.AddMethod("to_str", nil)
	b.AddMethod("to_num", nil)
}

type ErrorType struct {
	*ObjectType
}

func NewErrorType() *ErrorType {
	errorType := &ErrorType{
		ObjectType: NewObjectType("error", O_ERROR),
	}

	errorType.addErrorMethods()

	errorType.addErrorFields()

	return errorType
}

func (e *ErrorType) addErrorMethods() {
	e.AddMethod("to_str", nil)
}

func (e *ErrorType) addErrorFields() {

}

type ClassType struct {
	*ObjectType
	*GenericBase // 嵌入泛型基础结构
	astDecl      *ast.ClassDecl
	ctors        []Constructor
	registry     *Registry

	// 静态成员支持
	staticFields  map[string]*Field
	staticMethods map[string]*FunctionType
	staticValues  map[string]Value // 静态字段的值

	// 继承支持
	parent     *ClassType   // 父类
	parentName string       // 父类名称
	interfaces []*TraitType // 实现的接口

	// 方法重写支持
	overriddenMethods map[string]*FunctionType // 重写的方法

	// 泛型trait约束支持
	traitConstraints map[string][]*TraitType // 类型参数的trait约束，如 T: Display + Clone
}

func NewGenericClass(name string, typeParams []string, constraints map[string]Constraint) *ClassType {
	return &ClassType{
		ObjectType:        NewObjectType(name, O_CLASS),
		GenericBase:       NewGenericBase(typeParams, constraints),
		astDecl:           nil,
		ctors:             make([]Constructor, 0),
		registry:          nil,
		staticFields:      make(map[string]*Field),
		staticMethods:     make(map[string]*FunctionType),
		staticValues:      make(map[string]Value),
		interfaces:        make([]*TraitType, 0),
		overriddenMethods: make(map[string]*FunctionType),
		traitConstraints:  make(map[string][]*TraitType),
	}
}

func NewClassType(name string) *ClassType {
	return &ClassType{
		ObjectType:        NewObjectType(name, O_CLASS),
		GenericBase:       NewGenericBase(nil, nil),
		astDecl:           nil,
		ctors:             make([]Constructor, 0),
		registry:          nil,
		staticFields:      make(map[string]*Field),
		staticMethods:     make(map[string]*FunctionType),
		staticValues:      make(map[string]Value),
		interfaces:        make([]*TraitType, 0),
		overriddenMethods: make(map[string]*FunctionType),
		traitConstraints:  make(map[string][]*TraitType),
	}
}

func (ct *ClassType) AddOperator(op Operator, fn *FunctionType) {
	ct.operators[op] = fn

	// 同时注册到全局运算符注册表
	if ct.registry != nil {
		// 为所有可能的右操作数类型注册重载
		// 这里简化处理，实际需要更复杂的类型推断
		ct.registry.RegisterOperator(op, ct, anyType, fn.signatures[0].returnType, fn)
	}
}

func (ct *ClassType) AddField(name string, filed *Field) {
	ct.fields[name] = filed
}

// AddStaticField 添加静态字段
func (ct *ClassType) AddStaticField(name string, field *Field) {
	ct.staticFields[name] = field
}

// AddStaticMethod 添加静态方法
func (ct *ClassType) AddStaticMethod(name string, method *FunctionType) {
	ct.staticMethods[name] = method
}

// GetStaticField 获取静态字段
func (ct *ClassType) GetStaticField(name string) *Field {
	return ct.staticFields[name]
}

// GetStaticMethod 获取静态方法
func (ct *ClassType) GetStaticMethod(name string) *FunctionType {
	return ct.staticMethods[name]
}

// SetStaticValue 设置静态字段的值
func (ct *ClassType) SetStaticValue(name string, value Value) {
	ct.staticValues[name] = value
}

// GetStaticValue 获取静态字段的值
func (ct *ClassType) GetStaticValue(name string) Value {
	if value, exists := ct.staticValues[name]; exists {
		return value
	}
	return NULL
}

// SetParent 设置父类
func (ct *ClassType) SetParent(parent *ClassType) {
	ct.parent = parent
	ct.parentName = parent.Name()
}

// GetParent 获取父类
func (ct *ClassType) GetParent() *ClassType {
	return ct.parent
}

// AddInterface 添加实现的接口
func (ct *ClassType) AddInterface(iface *TraitType) {
	ct.interfaces = append(ct.interfaces, iface)
}

// GetInterfaces 获取所有实现的接口
func (ct *ClassType) GetInterfaces() []*TraitType {
	return ct.interfaces
}

// OverrideMethod 重写方法
func (ct *ClassType) OverrideMethod(name string, method *FunctionType) {
	ct.overriddenMethods[name] = method
}

// GetOverriddenMethod 获取重写的方法
func (ct *ClassType) GetOverriddenMethod(name string) *FunctionType {
	return ct.overriddenMethods[name]
}

// IsSubclassOf 检查是否为指定类的子类
func (ct *ClassType) IsSubclassOf(parentType *ClassType) bool {
	if ct.parent == nil {
		return false
	}
	if ct.parent == parentType {
		return true
	}
	return ct.parent.IsSubclassOf(parentType)
}

// ImplementsTrait 检查是否实现了指定trait
func (ct *ClassType) ImplementsTrait(iface *TraitType) bool {
	// 首先检查是否有直接的trait实现
	impls := GlobalTraitRegistry.GetImpl(iface.Name(), ct.Name())
	if len(impls) > 0 {
		return true
	}

	// 检查当前类的接口列表
	if slices.Contains(ct.interfaces, iface) {
		return true
	}

	// 检查父类是否实现了该接口
	if ct.parent != nil {
		return ct.parent.ImplementsTrait(iface)
	}

	return false
}

// Implements 实现Type接口的Implements方法
func (ct *ClassType) Implements(t Type) bool {
	if trait, ok := t.(*TraitType); ok {
		return ct.ImplementsTrait(trait)
	}
	return false
}

// MethodByName 重写方法查找，支持继承和方法重写
func (ct *ClassType) MethodByName(name string) Method {
	// 首先检查重写的方法
	if overridden, exists := ct.overriddenMethods[name]; exists {
		return overridden
	}

	// 检查当前类的方法
	if method, exists := ct.methods[name]; exists {
		return method
	}

	// 检查静态方法
	if staticMethod, exists := ct.staticMethods[name]; exists {
		return staticMethod
	}

	// 检查父类的方法
	if ct.parent != nil {
		return ct.parent.MethodByName(name)
	}

	return nil
}

// FieldByName 重写字段查找，支持继承
func (ct *ClassType) FieldByName(name string) Type {
	// 检查当前类的字段
	if field, exists := ct.fields[name]; exists {
		return field.ft
	}

	// 检查静态字段
	if staticField, exists := ct.staticFields[name]; exists {
		return staticField.ft
	}

	// 检查父类的字段
	if ct.parent != nil {
		return ct.parent.FieldByName(name)
	}

	return nil
}

func (ct *ClassType) Instantiate(typeArgs []Type) (*ClassType, error) {
	if !ct.IsGeneric() {
		return ct, nil
	}

	if len(typeArgs) != len(ct.TypeParameters()) {
		return nil, fmt.Errorf("type argument count mismatch: expected %d, got %d",
			len(ct.TypeParameters()), len(typeArgs))
	}

	// 检查类型约束
	if err := ct.CheckConstraints(typeArgs); err != nil {
		return nil, err
	}

	// 生成实例化类型的唯一标识
	typeKey := ct.generateTypeKey(typeArgs)

	// 检查缓存
	if instance := ct.getCachedInstance(typeKey); instance != nil {
		if classInstance, ok := instance.(*ClassType); ok {
			return classInstance, nil
		}
	}

	// 创建类型映射
	typeMap := make(map[string]Type)
	for i, param := range ct.TypeParameters() {
		typeMap[param] = typeArgs[i]
	}

	// 创建新的实例化类型
	instance := &ClassType{
		ObjectType:        NewObjectType(ct.Name()+"<"+typeKey+">", O_CLASS),
		GenericBase:       NewGenericBase(nil, nil), // 实例化后不再是泛型
		astDecl:           ct.astDecl,
		ctors:             make([]Constructor, 0),
		registry:          nil,
		staticFields:      make(map[string]*Field),
		staticMethods:     make(map[string]*FunctionType),
		staticValues:      make(map[string]Value),
		interfaces:        make([]*TraitType, 0),
		overriddenMethods: make(map[string]*FunctionType),
		traitConstraints:  make(map[string][]*TraitType),
	}

	// 实例化字段
	for name, field := range ct.fields {
		instance.fields[name] = &Field{
			ft:   ct.substituteType(field.ft, typeMap),
			mods: field.mods,
		}
	}

	// 实例化方法
	for name, method := range ct.methods {
		if fnType, ok := method.(*FunctionType); ok {
			instance.methods[name] = ct.substituteFunctionType(fnType, typeMap)
		} else {
			instance.methods[name] = method
		}
	}

	// 实例化运算符
	for op, fn := range ct.operators {
		instance.operators[op] = ct.substituteFunctionType(fn, typeMap)
	}

	// 实例化构造函数
	instance.ctors = ct.instantiateConstructors(typeArgs)

	// 缓存实例
	ct.cacheInstance(typeKey, instance)

	return instance, nil
}

// TODO
type TypeSubstituter interface {
	SubstituteType(t Type, typeMap map[string]Type) Type
}

// substituteType 替换类型中的泛型参数
func (ct *ClassType) substituteType(t Type, typeMap map[string]Type) Type {
	switch v := t.(type) {
	case *TypeParameter:
		// 检查是否是类型参数
		if replacement, exists := typeMap[v.Name()]; exists {
			return replacement
		}
		return t
	case *ArrayType:
		return &ArrayType{
			ObjectType:  NewObjectType(ct.substituteType(v.elementType, typeMap).Name()+"[]", O_ARR),
			elementType: ct.substituteType(v.elementType, typeMap),
		}
	case *MapType:
		return &MapType{
			ObjectType: NewObjectType(fmt.Sprintf("map<%s, %s>",
				ct.substituteType(v.keyType, typeMap).Name(),
				ct.substituteType(v.valueType, typeMap).Name()), O_MAP),
			keyType:   ct.substituteType(v.keyType, typeMap),
			valueType: ct.substituteType(v.valueType, typeMap),
		}
	default:
		return t
	}
}

// substituteFunctionType 替换函数类型中的泛型参数
func (ct *ClassType) substituteFunctionType(fnType *FunctionType, typeMap map[string]Type) *FunctionType {
	newFn := &FunctionType{
		name:        fnType.name,
		signatures:  make([]*FunctionSignature, 0),
		GenericBase: NewGenericBase(nil, nil), // 实例化后不再是泛型
		visible:     fnType.visible,
	}

	// 替换函数签名
	for _, sig := range fnType.signatures {
		newSig := &FunctionSignature{
			positionalParams: make([]*Parameter, 0),
			namedParams:      make([]*Parameter, 0),
			returnType:       ct.substituteType(sig.returnType, typeMap),
			isOperator:       sig.isOperator,
			operator:         sig.operator,
			body:             sig.body,
			builtin:          sig.builtin,
		}

		// 替换位置参数
		for _, param := range sig.positionalParams {
			newSig.positionalParams = append(newSig.positionalParams, &Parameter{
				name:         param.name,
				typ:          ct.substituteType(param.typ, typeMap),
				optional:     param.optional,
				defaultValue: param.defaultValue,
				variadic:     param.variadic,
				isNamed:      param.isNamed,
				required:     param.required,
			})
		}

		// 替换命名参数
		for _, param := range sig.namedParams {
			newSig.namedParams = append(newSig.namedParams, &Parameter{
				name:         param.name,
				typ:          ct.substituteType(param.typ, typeMap),
				optional:     param.optional,
				defaultValue: param.defaultValue,
				variadic:     param.variadic,
				isNamed:      param.isNamed,
				required:     param.required,
			})
		}

		// 替换可变参数
		if sig.variadicParam != nil {
			newSig.variadicParam = &Parameter{
				name:         sig.variadicParam.name,
				typ:          ct.substituteType(sig.variadicParam.typ, typeMap),
				optional:     sig.variadicParam.optional,
				defaultValue: sig.variadicParam.defaultValue,
				variadic:     sig.variadicParam.variadic,
				isNamed:      sig.variadicParam.isNamed,
				required:     sig.variadicParam.required,
			}
		}

		newFn.signatures = append(newFn.signatures, newSig)
	}

	return newFn
}

func (ct *ClassType) instantiateConstructors(typeArgs []Type) []Constructor {
	typeMap := make(map[string]Type)
	for i, param := range ct.TypeParameters() {
		typeMap[param] = typeArgs[i]
	}

	ctors := make([]Constructor, len(ct.ctors))
	for i, ctor := range ct.ctors {
		ctors[i] = ct.instantiateConstructor(ctor, typeMap)
	}
	return ctors
}

func (ct *ClassType) instantiateConstructor(ctor Constructor, typeMap map[string]Type) Constructor {
	// 这里需要根据具体的构造函数实现来替换类型参数
	// 例如，替换参数类型、返回类型等
	// return &ClassConstructor{
	// 	name:      ctor.Name(),
	// 	astDecl:   ctor.(*ClassConstructor).astDecl, // 需要类型断言
	// 	clsType:   ct,
	// 	signature: ct.instantiateSignature(ctor.Signature(), typeMap),
	// }
	return nil
}

func (ct *ClassType) instantiateSignature(signature []Type, typeMap map[string]Type) []Type {
	return nil
}

func (ct *ClassType) New(values ...Value) Value {
	inst := &ClassInstance{
		clsType: ct,
		fields:  make(map[string]Value),
	}

	ctor := ct.selectConstructor(values...)
	if ctor == nil {
		ctor = ct.getDefaultConstructor()
	}

	if ctor != nil {
		result := ctor.Call(values...)
		if initData, ok := result.(*ClassInstance); ok {
			inst.fields = initData.fields
		}
	}

	return inst
}

func (ct *ClassType) selectConstructor(values ...Value) Constructor {
	var bestMatch Constructor
	var bestScore int = -1

	for _, ctor := range ct.ctors {
		score := ct.calculateConstructorMatchScore(ctor, values)
		if score > bestScore {
			bestScore = score
			bestMatch = ctor
		}
	}

	return bestMatch
}

func (ct *ClassType) calculateConstructorMatchScore(ctor Constructor, values []Value) int {
	// 获取构造函数的签名
	signature := ctor.Signature()

	// 检查参数数量
	if len(values) != len(signature) {
		return -1
	}

	score := 0
	for i, value := range values {
		paramType := signature[i]
		valueType := value.Type()

		// 计算类型匹配分数
		if valueType.Name() == paramType.Name() {
			score += 100 // 精确匹配
		} else if valueType.AssignableTo(paramType) {
			score += 50 // 类型兼容
		} else if valueType.ConvertibleTo(paramType) {
			score += 25 // 可转换
		} else {
			return -1 // 不匹配
		}
	}

	return score
}

func (ct *ClassType) getDefaultConstructor() Constructor {
	// 查找无参数的默认构造函数
	for _, ctor := range ct.ctors {
		if len(ctor.Signature()) == 0 {
			return ctor
		}
	}
	return nil
}

// AddConstructor 添加构造函数
func (ct *ClassType) AddConstructor(ctor Constructor) {
	ct.ctors = append(ct.ctors, ctor)
}

// GetConstructors 获取所有构造函数
func (ct *ClassType) GetConstructors() []Constructor {
	return ct.ctors
}

// GetConstructorByName 根据名称获取命名构造函数
func (ct *ClassType) GetConstructorByName(name string) Constructor {
	for _, ctor := range ct.ctors {
		if ctor.Name() == name {
			return ctor
		}
	}
	return nil
}

var _ Type = &UnionType{}

type UnionType struct {
	types []Type
}

func NewUnionType(types ...Type) *UnionType {
	return &UnionType{
		types: types,
	}
}

func (u *UnionType) Name() string {
	return "union"
}

func (u *UnionType) Text() string {
	names := make([]string, len(u.types))
	for i, t := range u.types {
		names[i] = t.Name()
	}
	return strings.Join(names, " | ")
}

func (u *UnionType) New(values ...Value) Value {
	return nil
}

func (u *UnionType) Kind() ObjKind {
	return O_UNION
}

func (u *UnionType) AssignableTo(t Type) bool {
	for _, t := range u.types {
		if t.AssignableTo(t) {
			return true
		}
	}
	return false
}

func (u *UnionType) ConvertibleTo(t Type) bool {
	return false
}

func (u *UnionType) Implements(t Type) bool {
	return false
}

type IntersectionType struct {
	types []Type
}

func NewIntersectionType(types ...Type) *IntersectionType {
	return &IntersectionType{
		types: types,
	}
}

func (i *IntersectionType) Name() string {
	return "intersection"
}

func (i *IntersectionType) Text() string {
	names := make([]string, len(i.types))
	for i, t := range i.types {
		names[i] = t.Name()
	}
	return strings.Join(names, " | ")
}

func (i *IntersectionType) New(values ...Value) Value {
	return nil
}

func (i *IntersectionType) Kind() ObjKind {
	return O_INTERSECTION
}

func (i *IntersectionType) AssignableTo(t Type) bool {
	ret := true
	for _, t := range i.types {
		if !t.AssignableTo(t) {
			ret = false
			break
		}
	}
	return ret
}

func (i *IntersectionType) ConvertibleTo(t Type) bool {
	return false
}

func (i *IntersectionType) Implements(t Type) bool {
	return false
}

// NullableType represents a nullable type (e.g., T?)
type NullableType struct {
	baseType Type
}

func NewNullableType(baseType Type) *NullableType {
	return &NullableType{
		baseType: baseType,
	}
}

func (n *NullableType) Name() string {
	return n.baseType.Name() + "?"
}

func (n *NullableType) Text() string {
	return n.baseType.Text() + "?"
}

func (n *NullableType) New(values ...Value) Value {
	return nil
}

func (n *NullableType) Kind() ObjKind {
	return O_NULLABLE
}

func (n *NullableType) AssignableTo(t Type) bool {
	// 可空类型可以赋值给任何类型（因为null可以赋值给任何类型）
	// 或者赋值给相同的可空类型
	if nullable, ok := t.(*NullableType); ok {
		return n.baseType.AssignableTo(nullable.baseType)
	}
	return n.baseType.AssignableTo(t)
}

func (n *NullableType) ConvertibleTo(t Type) bool {
	return n.baseType.ConvertibleTo(t)
}

func (n *NullableType) Implements(t Type) bool {
	return n.baseType.Implements(t)
}

func (n *NullableType) BaseType() Type {
	return n.baseType
}

var _ Type = &FunctionType{}
var _ Method = &FunctionType{}

type FunctionType struct {
	*GenericBase // 嵌入泛型基础结构
	name         string
	signatures   []*FunctionSignature
	visible      bool
}

func NewGenericFunction(name string, typeParams []string, constraints map[string]Constraint) *FunctionType {
	return &FunctionType{
		name:        name,
		signatures:  make([]*FunctionSignature, 0),
		GenericBase: NewGenericBase(typeParams, constraints),
		visible:     true,
	}
}

// FunctionBuilder 函数构建器
type FunctionBuilder struct {
	name          string
	builtin       BuiltinFunction
	params        []*Parameter
	variadicParam *Parameter
	returnType    Type
	body          *ast.FuncDecl
}

// NewFunctionBuilder 创建一个新的函数构建器
func NewFunctionBuilder(name string) *FunctionBuilder {
	return &FunctionBuilder{
		name:       name,
		params:     make([]*Parameter, 0),
		returnType: GetVoidType(), // 默认返回类型
	}
}

// WithBuiltin 设置内置函数实现
func (fb *FunctionBuilder) WithBuiltin(builtin BuiltinFunction) *FunctionBuilder {
	fb.builtin = builtin
	return fb
}

// WithParameter 添加位置参数
func (fb *FunctionBuilder) WithParameter(name string, typ Type) *FunctionBuilder {
	fb.params = append(fb.params, NewParameter(name, typ))
	return fb
}

// WithOptionalParameter 添加可选参数
func (fb *FunctionBuilder) WithOptionalParameter(name string, typ Type, defaultValue Value) *FunctionBuilder {
	fb.params = append(fb.params, NewOptionalParameter(name, typ, defaultValue))
	return fb
}

// WithVariadicParameter 设置可变参数
func (fb *FunctionBuilder) WithVariadicParameter(name string, typ Type) *FunctionBuilder {
	fb.variadicParam = NewVariadicParameter(name, typ)
	return fb
}

// WithReturnType 设置返回类型
func (fb *FunctionBuilder) WithReturnType(returnType Type) *FunctionBuilder {
	fb.returnType = returnType
	return fb
}

// WithBody 设置函数体（用于非内置函数）
func (fb *FunctionBuilder) WithBody(body *ast.FuncDecl) *FunctionBuilder {
	fb.body = body
	return fb
}

// Build 构建FunctionType
func (fb *FunctionBuilder) Build() *FunctionType {
	signature := &FunctionSignature{
		positionalParams: fb.params,
		variadicParam:    fb.variadicParam,
		returnType:       fb.returnType,
		builtin:          fb.builtin,
		body:             fb.body,
	}

	return &FunctionType{
		name:        fb.name,
		signatures:  []*FunctionSignature{signature},
		visible:     true,
		GenericBase: nil,
	}
}

func (ft *FunctionType) Name() string {
	return ft.name
}

func (ft *FunctionType) Kind() ObjKind {
	return O_FUNC
}

func (ft *FunctionType) Text() string {
	typeParams := ""
	if ft.GenericBase != nil {
		if len(ft.GenericBase.typeParams) > 0 {
			typeParams = "<" + strings.Join(ft.GenericBase.typeParams, ", ") + ">"
		}
	}
	params := ""
	if len(ft.signatures) > 0 {
		for i, param := range ft.signatures[0].positionalParams {
			params += param.name + ":" + param.typ.Text()
			if i < len(ft.signatures[0].positionalParams)-1 {
				params += ", "
			}
		}
		if ft.signatures[0].variadicParam != nil {
			params += ft.signatures[0].variadicParam.name + ":..." + ft.signatures[0].variadicParam.typ.Text() + ", "
		}
	}
	return fmt.Sprintf("fn %s%s(%s) -> %s", ft.name, typeParams, params, ft.signatures[0].returnType.Text())
}

func (ft *FunctionType) AssignableTo(t Type) bool {
	return false
}

func (ft *FunctionType) ConvertibleTo(t Type) bool {
	return false
}

func (ft *FunctionType) Implements(t Type) bool {
	return false
}

func (ft *FunctionType) MatchSignature(args []Value, namedArgs map[string]Value) (*FunctionSignature, error) {
	var bestMatch *FunctionSignature
	var bestScore int = -1

	for _, sig := range ft.signatures {
		score, err := ft.calculateMatchScore(args, namedArgs, sig)
		if err != nil {
			continue
		}

		if score > bestScore {
			bestScore = score
			bestMatch = sig
		}
	}

	if bestMatch == nil {
		return nil, fmt.Errorf("no matching signature found for function %s", ft.name)
	}

	return bestMatch, nil
}

func (ft *FunctionType) calculateMatchScore(args []Value, namedArgs map[string]Value, sig *FunctionSignature) (int, error) {
	score := 0

	// 1. 检查位置参数
	posParamCount := len(sig.positionalParams)
	if len(args) < posParamCount {
		return -1, fmt.Errorf("insufficient positional arguments")
	}

	for i, param := range sig.positionalParams {
		if i >= len(args) {
			if !param.optional {
				return -1, fmt.Errorf("missing required parameter %s", param.name)
			}
			continue
		}

		argType := args[i].Type()
		paramScore := ft.calculateTypeMatchScore(argType, param.typ)
		if paramScore < 0 {
			return -1, fmt.Errorf("argument %d cannot be assigned to parameter %s", i, param.name)
		}
		score += paramScore
	}

	// 2. 检查可变参数
	if sig.variadicParam != nil {
		for i := posParamCount; i < len(args); i++ {
			argType := args[i].Type()
			if !argType.AssignableTo(sig.variadicParam.typ) {
				return -1, fmt.Errorf("variadic argument %d cannot be assigned to parameter %s", i, sig.variadicParam.name)
			}
			score += 50 // 可变参数匹配
		}
	} else if len(args) > posParamCount {
		return -1, fmt.Errorf("too many arguments")
	}

	// 3. 检查命名参数
	for _, param := range sig.namedParams {
		if value, exists := namedArgs[param.name]; exists {
			argType := value.Type()
			paramScore := ft.calculateTypeMatchScore(argType, param.typ)
			if paramScore < 0 {
				return -1, fmt.Errorf("named argument %s cannot be assigned to parameter %s", param.name, param.name)
			}
			score += paramScore
		} else if param.required {
			return -1, fmt.Errorf("missing required named parameter %s", param.name)
		}
		// 可选参数有默认值，不扣分
	}

	return score, nil
}

func (ft *FunctionType) calculateTypeMatchScore(from, to Type) int {
	// anyType 匹配：+10分（最低优先级，但允许匹配）
	if to.Name() == "any" {
		return 10
	}

	// 精确匹配：+100分
	if from.Name() == to.Name() {
		return 100
	}

	// 类型兼容：+50分
	if from.AssignableTo(to) {
		return 50
	}

	// 可转换：+25分
	if from.ConvertibleTo(to) {
		return 25
	}

	return -1
}

func (ft *FunctionType) Call(args []Value, namedArgs map[string]Value, evaluator ...core.FunctionEvaluator) (Value, error) {
	// 1. 匹配函数签名
	signature, err := ft.MatchSignature(args, namedArgs)
	if err != nil {
		return nil, fmt.Errorf("function signature mismatch: %v", err)
	}

	// 2. 构建完整参数列表
	allArgs, err := ft.buildCompleteArgs(args, namedArgs, signature)
	if err != nil {
		return nil, fmt.Errorf("argument preparation failed: %v", err)
	}

	// 3. 执行函数
	if signature.builtin != nil {
		return signature.builtin(allArgs...), nil
	} else if signature.body != nil {
		// 检查是否有求值器
		if len(evaluator) > 0 {
			if fnEvaluator, ok := evaluator[0].(core.FunctionEvaluator); ok {
				// 转换参数为 interface{}
				interfaceArgs := make([]any, len(allArgs))
				for i, arg := range allArgs {
					interfaceArgs[i] = arg
				}

				// 使用求值器执行函数
				result, err := fnEvaluator.EvaluateFunction(signature.body, interfaceArgs)
				if err != nil {
					return nil, fmt.Errorf("function evaluation failed: %w", err)
				}

				// 转换结果为 Value
				return ft.convertResultToValue(result)
			}
		}
		return nil, fmt.Errorf("function execution requires evaluator")
	}

	return nil, fmt.Errorf("function has no implementation")
}

// convertResultToValue 将 interface{} 转换为 Value
func (ft *FunctionType) convertResultToValue(result interface{}) (Value, error) {
	// 简单的类型转换实现
	switch v := result.(type) {
	case string:
		return &StringValue{Value: v}, nil
	case int:
		num := &NumberValue{Value: big.NewFloat(float64(v))}
		return num, nil
	case float64:
		num := &NumberValue{Value: big.NewFloat(v)}
		return num, nil
	case bool:
		if v {
			return TRUE, nil
		}
		return FALSE, nil
	case nil:
		return NULL, nil
	default:
		// 对于复杂类型，可能需要更复杂的转换逻辑
		return NULL, fmt.Errorf("unsupported result type: %T", result)
	}
}

func (ft *FunctionType) buildCompleteArgs(args []Value, namedArgs map[string]Value, sig *FunctionSignature) ([]Value, error) {
	allArgs := make([]Value, 0)

	// 1. 处理位置参数
	for i, param := range sig.positionalParams {
		if i < len(args) {
			// 有传入的参数
			converted, err := ft.convertValue(args[i], param.typ)
			if err != nil {
				return nil, err
			}
			allArgs = append(allArgs, converted)
		} else if param.optional {
			// 使用默认值
			allArgs = append(allArgs, param.defaultValue)
		} else {
			return nil, fmt.Errorf("missing required parameter %s", param.name)
		}
	}

	// 2. 处理可变参数
	if sig.variadicParam != nil {
		for i := len(sig.positionalParams); i < len(args); i++ {
			converted, err := ft.convertValue(args[i], sig.variadicParam.typ)
			if err != nil {
				return nil, err
			}
			allArgs = append(allArgs, converted)
		}
	}

	// 3. 处理命名参数
	for _, param := range sig.namedParams {
		if value, exists := namedArgs[param.name]; exists {
			converted, err := ft.convertValue(value, param.typ)
			if err != nil {
				return nil, err
			}
			allArgs = append(allArgs, converted)
		} else if param.optional {
			allArgs = append(allArgs, param.defaultValue)
		} else {
			return nil, fmt.Errorf("missing required named parameter %s", param.name)
		}
	}

	return allArgs, nil
}

func (ft *FunctionType) convertValue(v Value, targetType Type) (Value, error) {
	// 实现类型转换逻辑
	switch targetType.Name() {
	case "string":
		return &StringValue{Value: v.Text()}, nil
	case "number":
		// 实现数字转换
		return v, nil
	default:
		return v, nil
	}
}

type FunctionSignature struct {
	// 位置参数
	positionalParams []*Parameter
	// 命名参数（解构参数）
	namedParams []*Parameter
	// 可变参数
	variadicParam *Parameter
	returnType    Type
	// 是否为运算符重载
	isOperator bool
	operator   Operator
	// 函数实体
	body *ast.FuncDecl
	// 内置函数
	builtin BuiltinFunction
}

type Parameter struct {
	name         string
	typ          Type
	optional     bool
	defaultValue Value
	variadic     bool
	// 命名参数相关
	isNamed  bool
	required bool
}

// NewParameter 创建一个新的参数
func NewParameter(name string, typ Type) *Parameter {
	return &Parameter{
		name:     name,
		typ:      typ,
		required: true,
	}
}

// NewOptionalParameter 创建一个可选参数
func NewOptionalParameter(name string, typ Type, defaultValue Value) *Parameter {
	return &Parameter{
		name:         name,
		typ:          typ,
		optional:     true,
		defaultValue: defaultValue,
		required:     false,
	}
}

// NewVariadicParameter 创建一个可变参数
func NewVariadicParameter(name string, typ Type) *Parameter {
	return &Parameter{
		name:     name,
		typ:      typ,
		variadic: true,
		required: true,
	}
}

type OperatorOverload struct {
	operator   Operator
	leftType   Type
	rightType  Type
	returnType Type
	function   *FunctionType
}

// ==================== 数组类型 ====================

type ArrayType struct {
	*ObjectType
	elementType Type
}

func NewArrayType(elementType Type) *ArrayType {
	arrayType := &ArrayType{
		ObjectType:  NewObjectType(elementType.Name()+"[]", O_ARR),
		elementType: elementType,
	}

	arrayType.addArrayMethods()

	return arrayType
}

func (a *ArrayType) ElementType() Type {
	return a.elementType
}

func (a *ArrayType) addArrayMethods() {
	// 添加数组的常用方法
	a.AddMethod("length", &FunctionType{
		name: "length",
		signatures: []*FunctionSignature{
			{
				positionalParams: []*Parameter{
					{
						name: "arr",
						typ:  a,
					},
				},
				returnType: numberType,
				builtin: func(args ...Value) Value {
					arr := args[0].(*ArrayValue)
					return &NumberValue{Value: big.NewFloat(float64(len(arr.Elements)))}
				},
			},
		},
	})

	a.AddMethod("push", &FunctionType{
		name: "push",
		signatures: []*FunctionSignature{
			{
				positionalParams: []*Parameter{
					{
						name: "arr",
						typ:  a,
					},
					{
						name: "element",
						typ:  a.elementType,
					},
				},
				returnType: a,
				builtin: func(args ...Value) Value {
					arr := args[0].(*ArrayValue)
					element := args[1]
					arr.Elements = append(arr.Elements, element)
					return arr
				},
			},
		},
	})

	a.AddMethod("pop", &FunctionType{
		name: "pop",
		signatures: []*FunctionSignature{
			{
				positionalParams: []*Parameter{
					{
						name: "arr",
						typ:  a,
					},
				},
				returnType: a.elementType,
				builtin: func(args ...Value) Value {
					arr := args[0].(*ArrayValue)
					if len(arr.Elements) == 0 {
						return NULL
					}
					last := arr.Elements[len(arr.Elements)-1]
					arr.Elements = arr.Elements[:len(arr.Elements)-1]
					return last
				},
			},
		},
	})

	a.AddMethod("slice", &FunctionType{
		name: "slice",
		signatures: []*FunctionSignature{
			{
				positionalParams: []*Parameter{
					{
						name: "arr",
						typ:  a,
					},
					{
						name:         "start",
						typ:          numberType,
						optional:     true,
						defaultValue: &NumberValue{Value: big.NewFloat(0)},
					},
					{
						name:         "end",
						typ:          numberType,
						optional:     true,
						defaultValue: &NumberValue{Value: big.NewFloat(-1)},
					},
				},
				returnType: a,
				builtin: func(args ...Value) Value {
					arr := args[0].(*ArrayValue)
					start, _ := args[1].(*NumberValue).Value.Int64()
					end, _ := args[2].(*NumberValue).Value.Int64()
					startInt := int(start)
					endInt := int(end)

					if endInt == -1 {
						endInt = len(arr.Elements)
					}

					if startInt < 0 {
						startInt = len(arr.Elements) + startInt
					}
					if endInt < 0 {
						endInt = len(arr.Elements) + endInt
					}

					if startInt >= len(arr.Elements) || endInt <= startInt {
						return &ArrayValue{Elements: []Value{}, arrayType: a}
					}

					if endInt > len(arr.Elements) {
						endInt = len(arr.Elements)
					}

					return &ArrayValue{
						Elements:  arr.Elements[startInt:endInt],
						arrayType: a,
					}
				},
			},
		},
	})

	a.AddMethod("join", &FunctionType{
		name: "join",
		signatures: []*FunctionSignature{
			{
				positionalParams: []*Parameter{
					{
						name: "arr",
						typ:  a,
					},
					{
						name:         "separator",
						typ:          stringType,
						optional:     true,
						defaultValue: &StringValue{Value: ","},
					},
				},
				returnType: stringType,
				builtin: func(args ...Value) Value {
					arr := args[0].(*ArrayValue)
					separator := args[1].(*StringValue).Value

					parts := make([]string, len(arr.Elements))
					for i, elem := range arr.Elements {
						parts[i] = elem.Text()
					}

					return &StringValue{Value: strings.Join(parts, separator)}
				},
			},
		},
	})

	a.AddMethod("map", &FunctionType{
		name: "map",
		signatures: []*FunctionSignature{
			{
				positionalParams: []*Parameter{
					{
						name: "arr",
						typ:  a,
					},
					{
						name: "fn",
						typ: &FunctionType{
							name: "mapper",
							signatures: []*FunctionSignature{
								{
									positionalParams: []*Parameter{
										{name: "item", typ: a.elementType},
										{name: "index", typ: numberType},
									},
									returnType: anyType,
								},
							},
						},
					},
				},
				returnType: &ArrayType{
					ObjectType:  NewObjectType("any[]", O_ARR),
					elementType: anyType,
				},
				builtin: func(args ...Value) Value {
					arr := args[0].(*ArrayValue)
					mapper := args[1].(*FunctionValue)

					result := make([]Value, len(arr.Elements))
					for i, elem := range arr.Elements {
						mapped, _ := mapper.Call([]Value{elem, &NumberValue{Value: big.NewFloat(float64(i))}}, nil)
						result[i] = mapped
					}

					return &ArrayValue{
						Elements: result,
						arrayType: &ArrayType{
							ObjectType:  NewObjectType("any[]", O_ARR),
							elementType: anyType,
						},
					}
				},
			},
		},
	})

	a.AddMethod("filter", &FunctionType{
		name: "filter",
		signatures: []*FunctionSignature{
			{
				positionalParams: []*Parameter{
					{
						name: "arr",
						typ:  a,
					},
					{
						name: "fn",
						typ: &FunctionType{
							name: "predicate",
							signatures: []*FunctionSignature{
								{
									positionalParams: []*Parameter{
										{name: "item", typ: a.elementType},
										{name: "index", typ: numberType},
									},
									returnType: boolType,
								},
							},
						},
					},
				},
				returnType: a,
				builtin: func(args ...Value) Value {
					arr := args[0].(*ArrayValue)
					predicate := args[1].(*FunctionValue)

					result := make([]Value, 0)
					for i, elem := range arr.Elements {
						keep, _ := predicate.Call([]Value{elem, &NumberValue{Value: big.NewFloat(float64(i))}}, nil)
						if keep.(*BoolValue).Value {
							result = append(result, elem)
						}
					}

					return &ArrayValue{
						Elements:  result,
						arrayType: a,
					}
				},
			},
		},
	})

	a.AddMethod("reduce", &FunctionType{
		name: "reduce",
		signatures: []*FunctionSignature{
			{
				positionalParams: []*Parameter{
					{
						name: "arr",
						typ:  a,
					},
					{
						name: "fn",
						typ: &FunctionType{
							name: "reducer",
							signatures: []*FunctionSignature{
								{
									positionalParams: []*Parameter{
										{name: "accumulator", typ: anyType},
										{name: "item", typ: a.elementType},
										{name: "index", typ: numberType},
									},
									returnType: anyType,
								},
							},
						},
					},
					{
						name:         "initial",
						typ:          anyType,
						optional:     true,
						defaultValue: NULL,
					},
				},
				returnType: anyType,
				builtin: func(args ...Value) Value {
					arr := args[0].(*ArrayValue)
					reducer := args[1].(*FunctionValue)
					initial := args[2]

					if len(arr.Elements) == 0 {
						return initial
					}

					accumulator := initial
					if initial == NULL && len(arr.Elements) > 0 {
						accumulator = arr.Elements[0]
						arr.Elements = arr.Elements[1:]
					}

					for i, elem := range arr.Elements {
						accumulator, _ = reducer.Call([]Value{accumulator, elem, &NumberValue{Value: big.NewFloat(float64(i))}}, nil)
					}

					return accumulator
				},
			},
		},
	})
}

func (a *ArrayType) AssignableTo(t Type) bool {
	if other, ok := t.(*ArrayType); ok {
		return a.elementType.AssignableTo(other.elementType)
	}
	return false
}

func (a *ArrayType) ConvertibleTo(t Type) bool {
	return false
}

func (a *ArrayType) Implements(t Type) bool {
	return false
}

// ==================== Map类型 ====================

type MapType struct {
	*ObjectType
	keyType   Type
	valueType Type
}

func NewMapType(keyType, valueType Type) *MapType {
	mapType := &MapType{
		ObjectType: NewObjectType(fmt.Sprintf("map<%s, %s>", keyType.Name(), valueType.Name()), O_MAP),
		keyType:    keyType,
		valueType:  valueType,
	}

	mapType.addMapMethods()

	return mapType
}

func (m *MapType) KeyType() Type {
	return m.keyType
}

func (m *MapType) ValueType() Type {
	return m.valueType
}

func (m *MapType) addMapMethods() {
	// 添加map的常用方法
	m.AddMethod("length", &FunctionType{
		name: "length",
		signatures: []*FunctionSignature{
			{
				positionalParams: []*Parameter{
					{
						name: "map",
						typ:  m,
					},
				},
				returnType: numberType,
				builtin: func(args ...Value) Value {
					mapVal := args[0].(*MapValue)
					return &NumberValue{Value: big.NewFloat(float64(len(mapVal.Pairs)))}
				},
			},
		},
	})

	m.AddMethod("is_empty", &FunctionType{
		name: "is_empty",
		signatures: []*FunctionSignature{
			{
				positionalParams: []*Parameter{
					{
						name: "map",
						typ:  m,
					},
				},
				returnType: boolType,
				builtin: func(args ...Value) Value {
					mapVal := args[0].(*MapValue)
					return &BoolValue{Value: len(mapVal.Pairs) == 0}
				},
			},
		},
	})

	m.AddMethod("has", &FunctionType{
		name: "has",
		signatures: []*FunctionSignature{
			{
				positionalParams: []*Parameter{
					{
						name: "map",
						typ:  m,
					},
					{
						name: "key",
						typ:  m.keyType,
					},
				},
				returnType: boolType,
				builtin: func(args ...Value) Value {
					mapVal := args[0].(*MapValue)
					key := args[1]

					for _, pair := range mapVal.Pairs {
						if pair.Key.Text() == key.Text() {
							return &BoolValue{Value: true}
						}
					}
					return &BoolValue{Value: false}
				},
			},
		},
	})

	m.AddMethod("get", &FunctionType{
		name: "get",
		signatures: []*FunctionSignature{
			{
				positionalParams: []*Parameter{
					{
						name: "map",
						typ:  m,
					},
					{
						name: "key",
						typ:  m.keyType,
					},
				},
				returnType: NewUnionType(m.valueType, nullType),
				builtin: func(args ...Value) Value {
					mapVal := args[0].(*MapValue)
					key := args[1]

					for _, pair := range mapVal.Pairs {
						if pair.Key.Text() == key.Text() {
							return pair.Value
						}
					}
					return NULL
				},
			},
		},
	})

	m.AddMethod("insert", &FunctionType{
		name: "insert",
		signatures: []*FunctionSignature{
			{
				positionalParams: []*Parameter{
					{
						name: "map",
						typ:  m,
					},
					{
						name: "key",
						typ:  m.keyType,
					},
					{
						name: "value",
						typ:  m.valueType,
					},
				},
				returnType: NewUnionType(m.valueType, nullType),
				builtin: func(args ...Value) Value {
					mapVal := args[0].(*MapValue)
					key := args[1]
					value := args[2]

					// 查找是否已存在该key
					for _, pair := range mapVal.Pairs {
						if pair.Key.Text() == key.Text() {
							oldValue := pair.Value
							pair.Value = value
							return oldValue
						}
					}

					// 添加新的键值对
					mapVal.Pairs = append(mapVal.Pairs, &KeyValuePair{
						Key:   key,
						Value: value,
					})

					return NULL
				},
			},
		},
	})

	m.AddMethod("remove", &FunctionType{
		name: "remove",
		signatures: []*FunctionSignature{
			{
				positionalParams: []*Parameter{
					{
						name: "map",
						typ:  m,
					},
					{
						name: "key",
						typ:  m.keyType,
					},
				},
				returnType: NewUnionType(m.valueType, nullType),
				builtin: func(args ...Value) Value {
					mapVal := args[0].(*MapValue)
					key := args[1]

					for i, pair := range mapVal.Pairs {
						if pair.Key.Text() == key.Text() {
							removedValue := pair.Value
							mapVal.Pairs = append(mapVal.Pairs[:i], mapVal.Pairs[i+1:]...)
							return removedValue
						}
					}

					return NULL
				},
			},
		},
	})

	m.AddMethod("clear", &FunctionType{
		name: "clear",
		signatures: []*FunctionSignature{
			{
				positionalParams: []*Parameter{
					{
						name: "map",
						typ:  m,
					},
				},
				returnType: voidType,
				builtin: func(args ...Value) Value {
					mapVal := args[0].(*MapValue)
					mapVal.Pairs = []*KeyValuePair{}
					return NULL
				},
			},
		},
	})

	m.AddMethod("keys", &FunctionType{
		name: "keys",
		signatures: []*FunctionSignature{
			{
				positionalParams: []*Parameter{
					{
						name: "map",
						typ:  m,
					},
				},
				returnType: &ArrayType{
					ObjectType:  NewObjectType(m.keyType.Name()+"[]", O_ARR),
					elementType: m.keyType,
				},
				builtin: func(args ...Value) Value {
					mapVal := args[0].(*MapValue)
					keys := make([]Value, len(mapVal.Pairs))
					for i, pair := range mapVal.Pairs {
						keys[i] = pair.Key
					}
					return &ArrayValue{
						Elements: keys,
						arrayType: &ArrayType{
							ObjectType:  NewObjectType(m.keyType.Name()+"[]", O_ARR),
							elementType: m.keyType,
						},
					}
				},
			},
		},
	})

	m.AddMethod("values", &FunctionType{
		name: "values",
		signatures: []*FunctionSignature{
			{
				positionalParams: []*Parameter{
					{
						name: "map",
						typ:  m,
					},
				},
				returnType: &ArrayType{
					ObjectType:  NewObjectType(m.valueType.Name()+"[]", O_ARR),
					elementType: m.valueType,
				},
				builtin: func(args ...Value) Value {
					mapVal := args[0].(*MapValue)
					values := make([]Value, len(mapVal.Pairs))
					for i, pair := range mapVal.Pairs {
						values[i] = pair.Value
					}
					return &ArrayValue{
						Elements: values,
						arrayType: &ArrayType{
							ObjectType:  NewObjectType(m.valueType.Name()+"[]", O_ARR),
							elementType: m.valueType,
						},
					}
				},
			},
		},
	})

	m.AddMethod("to_str", &FunctionType{
		name: "to_str",
		signatures: []*FunctionSignature{
			{
				positionalParams: []*Parameter{
					{
						name: "map",
						typ:  m,
					},
				},
				returnType: stringType,
				builtin: func(args ...Value) Value {
					mapVal := args[0].(*MapValue)

					pairs := make([]string, len(mapVal.Pairs))
					for i, pair := range mapVal.Pairs {
						pairs[i] = fmt.Sprintf("%s: %s", pair.Key.Text(), pair.Value.Text())
					}

					return &StringValue{Value: fmt.Sprintf("{%s}", strings.Join(pairs, ", "))}
				},
			},
		},
	})
}

func (m *MapType) AssignableTo(t Type) bool {
	if other, ok := t.(*MapType); ok {
		return m.keyType.AssignableTo(other.keyType) && m.valueType.AssignableTo(other.valueType)
	}
	return false
}

func (m *MapType) ConvertibleTo(t Type) bool {
	return false
}

func (m *MapType) Implements(t Type) bool {
	return false
}

// ==================== Trait系统 ====================

// TraitType 表示trait类型
type TraitType struct {
	*ObjectType
	name        string
	methods     map[string]*FunctionType
	fields      map[string]*Field
	operators   map[Operator]*FunctionType
	defaultImpl map[string]*FunctionType // 默认实现
}

func NewTraitType(name string) *TraitType {
	return &TraitType{
		ObjectType:  NewObjectType(name, O_TRAIT),
		name:        name,
		methods:     make(map[string]*FunctionType),
		fields:      make(map[string]*Field),
		operators:   make(map[Operator]*FunctionType),
		defaultImpl: make(map[string]*FunctionType),
	}
}

func (t *TraitType) Name() string {
	return t.name
}

func (t *TraitType) Kind() ObjKind {
	return O_TRAIT
}

func (t *TraitType) Text() string {
	return fmt.Sprintf("trait %s", t.name)
}

func (t *TraitType) AssignableTo(u Type) bool {
	return false // trait不能直接赋值
}

func (t *TraitType) ConvertibleTo(u Type) bool {
	return false
}

func (t *TraitType) Implements(u Type) bool {
	// 如果u是另一个trait，检查是否有实现关系
	if otherTrait, ok := u.(*TraitType); ok {
		// 检查是否有trait实现关系
		impls := GlobalTraitRegistry.GetImpl(otherTrait.Name(), t.Name())
		return len(impls) > 0
	}

	// trait不能实现具体类型
	return false
}

// AddMethod 添加trait方法
func (t *TraitType) AddMethod(name string, method *FunctionType) {
	t.methods[name] = method
}

// AddField 添加trait字段
func (t *TraitType) AddField(name string, field *Field) {
	t.fields[name] = field
}

// AddOperator 添加trait运算符
func (t *TraitType) AddOperator(op Operator, fn *FunctionType) {
	t.operators[op] = fn
}

// AddDefaultImpl 添加默认实现
func (t *TraitType) AddDefaultImpl(name string, impl *FunctionType) {
	t.defaultImpl[name] = impl
}

// GetDefaultImpl 获取默认实现
func (t *TraitType) GetDefaultImpl(name string) *FunctionType {
	return t.defaultImpl[name]
}

// TraitImpl 表示trait的实现
type TraitImpl struct {
	trait     *TraitType
	implType  Type
	methods   map[string]*FunctionType
	fields    map[string]*Field
	operators map[Operator]*FunctionType
}

func NewTraitImpl(trait *TraitType, implType Type) *TraitImpl {
	return &TraitImpl{
		trait:     trait,
		implType:  implType,
		methods:   make(map[string]*FunctionType),
		fields:    make(map[string]*Field),
		operators: make(map[Operator]*FunctionType),
	}
}

func (ti *TraitImpl) AddMethod(name string, method *FunctionType) {
	ti.methods[name] = method
}

func (ti *TraitImpl) AddField(name string, field *Field) {
	ti.fields[name] = field
}

func (ti *TraitImpl) AddOperator(op Operator, fn *FunctionType) {
	ti.operators[op] = fn
}

// GetMethod 获取方法实现，优先返回自定义实现，否则返回默认实现
func (ti *TraitImpl) GetMethod(name string) *FunctionType {
	if method, exists := ti.methods[name]; exists {
		return method
	}
	return ti.trait.GetDefaultImpl(name)
}

// TraitRegistry trait注册表
type TraitRegistry struct {
	traits map[string]*TraitType
	impls  map[string][]*TraitImpl // key: "TraitName:TypeName"
	mutex  sync.RWMutex
}

func NewTraitRegistry() *TraitRegistry {
	return &TraitRegistry{
		traits: make(map[string]*TraitType),
		impls:  make(map[string][]*TraitImpl),
	}
}

func (tr *TraitRegistry) RegisterTrait(trait *TraitType) {
	tr.mutex.Lock()
	defer tr.mutex.Unlock()
	tr.traits[trait.Name()] = trait
}

func (tr *TraitRegistry) RegisterImpl(impl *TraitImpl) {
	tr.mutex.Lock()
	defer tr.mutex.Unlock()
	key := fmt.Sprintf("%s:%s", impl.trait.Name(), impl.implType.Name())
	tr.impls[key] = append(tr.impls[key], impl)
}

func (tr *TraitRegistry) GetTrait(name string) *TraitType {
	tr.mutex.RLock()
	defer tr.mutex.RUnlock()
	return tr.traits[name]
}

func (tr *TraitRegistry) GetImpl(traitName, typeName string) []*TraitImpl {
	tr.mutex.RLock()
	defer tr.mutex.RUnlock()
	key := fmt.Sprintf("%s:%s", traitName, typeName)
	return tr.impls[key]
}

// 全局trait注册表
var GlobalTraitRegistry = NewTraitRegistry()

// ConvertTraitDecl 转换trait声明
func ConvertTraitDecl(decl *ast.TraitDecl) *TraitType {
	trait := NewTraitType(decl.Name.Name)

	// 处理字段
	if decl.Fields != nil {
		for _, field := range decl.Fields.List {
			fieldType := convertTypeExpr(field.Type)
			fieldMods := convertFieldModifiers(field.Modifiers)

			trait.AddField(field.Name.Name, &Field{
				ft:   fieldType,
				mods: fieldMods,
			})
		}
	}

	// TODO: 处理方法（需要从FieldList中解析方法）
	// trait的方法通常也是作为字段定义的，但类型是函数类型

	return trait
}

// ConvertImplDecl 转换impl声明
func ConvertImplDecl(decl *ast.ImplDecl) *TraitImpl {
	// 获取trait
	traitName := decl.Trait.Name
	trait := GlobalTraitRegistry.GetTrait(traitName)
	if trait == nil {
		panic(fmt.Sprintf("trait %s not found", traitName))
	}

	// 获取实现类型
	var implType Type
	if decl.ImplDeclBody != nil {
		implType = convertTypeExpr(&ast.TypeReference{Name: decl.ImplDeclBody.Class})
	} else {
		// 处理多类型实现
		implType = anyType // 简化处理
	}

	impl := NewTraitImpl(trait, implType)

	// TODO: 从BlockStmt中解析字段、方法和运算符实现
	// 这里需要分析实现体中的语句

	return impl
}

// convertOperatorDecl 转换运算符声明
func convertOperatorDecl(op *ast.OperatorDecl) *FunctionType {
	// 将token.Token转换为Operator
	operator := tokenToOperator(op.Name)

	return &FunctionType{
		name: fmt.Sprintf("operator_%s", op.Name),
		signatures: []*FunctionSignature{
			{
				positionalParams: []*Parameter{
					{name: "self", typ: anyType},
					{name: "other", typ: anyType},
				},
				returnType: anyType,
				isOperator: true,
				operator:   operator,
				body:       nil, // 暂时设为nil，后续需要处理
			},
		},
	}
}

// tokenToOperator 将token.Token转换为Operator
func tokenToOperator(tok token.Token) Operator {
	switch tok {
	case token.PLUS:
		return OpAdd
	case token.MINUS:
		return OpSub
	case token.ASTERISK:
		return OpMul
	case token.SLASH:
		return OpDiv
	case token.MOD:
		return OpMod
	case token.CONCAT:
		return OpConcat
	case token.AND:
		return OpAnd
	case token.OR:
		return OpOr
	case token.EQ:
		return OpEqual
	case token.NOT:
		return OpNot
	case token.LT:
		return OpLess
	case token.LE:
		return OpLessEqual
	case token.GT:
		return OpGreater
	case token.GE:
		return OpGreaterEqual
	case token.ASSIGN:
		return OpAssign
	case token.LBRACKET:
		return OpIndex
	case token.LPAREN:
		return OpCall
	default:
		return OpIllegal
	}
}

// 示例：Service trait
func ExampleServiceTrait() {
	// 定义Service trait
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

	// 添加运算符
	serviceTrait.AddOperator(OpGreater, &FunctionType{
		name: "operator_>",
		signatures: []*FunctionSignature{
			{
				positionalParams: []*Parameter{
					{name: "self", typ: serviceTrait},
					{name: "other", typ: serviceTrait},
				},
				returnType: boolType,
				isOperator: true,
				operator:   OpGreater,
			},
		},
	})

	// 注册trait
	GlobalTraitRegistry.RegisterTrait(serviceTrait)

	// 为GameService实现Service trait
	gameServiceType := NewClassType("GameService")
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

	// 实现运算符
	impl.AddOperator(OpGreater, &FunctionType{
		name: "operator_>",
		signatures: []*FunctionSignature{
			{
				positionalParams: []*Parameter{
					{name: "self", typ: gameServiceType},
					{name: "other", typ: serviceTrait},
				},
				returnType: boolType,
				isOperator: true,
				operator:   OpGreater,
				builtin: func(args ...Value) Value {
					return &BoolValue{Value: true}
				},
			},
		},
	})

	// 注册实现
	GlobalTraitRegistry.RegisterImpl(impl)
}

// 检查类型是否实现了trait
func ImplementsTrait(t Type, traitName string) bool {
	trait := GlobalTraitRegistry.GetTrait(traitName)
	if trait == nil {
		return false
	}

	impls := GlobalTraitRegistry.GetImpl(traitName, t.Name())
	return len(impls) > 0
}

// ==================== 元组类型系统 ====================

// TupleType 表示元组类型
type TupleType struct {
	*ObjectType
	elementTypes []Type
	fieldNames   []string // 命名元组的字段名，空字符串表示匿名字段
}

func NewTupleType(elementTypes ...Type) *TupleType {
	return &TupleType{
		ObjectType:   NewObjectType("tuple", O_TUPLE),
		elementTypes: elementTypes,
		fieldNames:   make([]string, len(elementTypes)),
	}
}

func NewNamedTupleType(fieldNames []string, elementTypes []Type) *TupleType {
	if len(fieldNames) != len(elementTypes) {
		panic("field names and element types must have the same length")
	}

	return &TupleType{
		ObjectType:   NewObjectType("named_tuple", O_TUPLE),
		elementTypes: elementTypes,
		fieldNames:   fieldNames,
	}
}

func (t *TupleType) Name() string {
	if t.isNamed() {
		return t.namedTupleString()
	}
	return t.anonymousTupleString()
}

func (t *TupleType) Kind() ObjKind {
	return O_TUPLE
}

func (t *TupleType) Text() string {
	return t.Name()
}

func (t *TupleType) AssignableTo(u Type) bool {
	if other, ok := u.(*TupleType); ok {
		if len(t.elementTypes) != len(other.elementTypes) {
			return false
		}

		for i, elemType := range t.elementTypes {
			if !elemType.AssignableTo(other.elementTypes[i]) {
				return false
			}
		}
		return true
	}
	return false
}

func (t *TupleType) ConvertibleTo(u Type) bool {
	return false
}

func (t *TupleType) Implements(u Type) bool {
	return false
}

func (t *TupleType) ElementTypes() []Type {
	return t.elementTypes
}

func (t *TupleType) FieldNames() []string {
	return t.fieldNames
}

func (t *TupleType) Length() int {
	return len(t.elementTypes)
}

func (t *TupleType) isNamed() bool {
	for _, name := range t.fieldNames {
		if name != "" {
			return true
		}
	}
	return false
}

func (t *TupleType) anonymousTupleString() string {
	parts := make([]string, len(t.elementTypes))
	for i, elemType := range t.elementTypes {
		parts[i] = elemType.Name()
	}
	return fmt.Sprintf("[%s]", strings.Join(parts, ", "))
}

func (t *TupleType) namedTupleString() string {
	parts := make([]string, len(t.elementTypes))
	for i, elemType := range t.elementTypes {
		if t.fieldNames[i] != "" {
			parts[i] = fmt.Sprintf("%s: %s", t.fieldNames[i], elemType.Name())
		} else {
			parts[i] = elemType.Name()
		}
	}
	return fmt.Sprintf("(%s)", strings.Join(parts, ", "))
}

// GetElementType 获取指定位置的元素类型
func (t *TupleType) GetElementType(index int) Type {
	if index < 0 || index >= len(t.elementTypes) {
		return nil
	}
	return t.elementTypes[index]
}

// GetFieldIndex 根据字段名获取索引
func (t *TupleType) GetFieldIndex(fieldName string) int {
	for i, name := range t.fieldNames {
		if name == fieldName {
			return i
		}
	}
	return -1
}

// 添加元组类型的方法
func (t *TupleType) addTupleMethods() {
	// length 方法
	t.AddMethod("length", &FunctionType{
		name: "length",
		signatures: []*FunctionSignature{
			{
				positionalParams: []*Parameter{
					{name: "tuple", typ: t},
				},
				returnType: numberType,
				builtin: func(args ...Value) Value {
					tuple := args[0].(*TupleValue)
					return &NumberValue{Value: big.NewFloat(float64(tuple.Length()))}
				},
			},
		},
	})

	// get 方法 - 通过索引获取元素
	t.AddMethod("get", &FunctionType{
		name: "get",
		signatures: []*FunctionSignature{
			{
				positionalParams: []*Parameter{
					{name: "tuple", typ: t},
					{name: "index", typ: numberType},
				},
				returnType: anyType,
				builtin: func(args ...Value) Value {
					tuple := args[0].(*TupleValue)
					index, _ := args[1].(*NumberValue).Value.Int64()
					return tuple.Get(int(index))
				},
			},
		},
	})

	// with 方法 - 创建新元组，更新指定位置
	t.AddMethod("with", &FunctionType{
		name: "with",
		signatures: []*FunctionSignature{
			{
				positionalParams: []*Parameter{
					{name: "tuple", typ: t},
					{name: "index", typ: numberType},
					{name: "value", typ: anyType},
				},
				returnType: t,
				builtin: func(args ...Value) Value {
					tuple := args[0].(*TupleValue)
					index, _ := args[1].(*NumberValue).Value.Int64()
					value := args[2]

					newTuple, err := tuple.With(int(index), value)
					if err != nil {
						return NULL
					}
					return newTuple
				},
			},
		},
	})

	// destructure 方法 - 解构元组
	t.AddMethod("destructure", &FunctionType{
		name: "destructure",
		signatures: []*FunctionSignature{
			{
				positionalParams: []*Parameter{
					{name: "tuple", typ: t},
				},
				returnType: &ArrayType{
					ObjectType:  NewObjectType("any[]", O_ARR),
					elementType: anyType,
				},
				builtin: func(args ...Value) Value {
					tuple := args[0].(*TupleValue)
					elements := tuple.Destructure()
					return &ArrayValue{
						Elements: elements,
						arrayType: &ArrayType{
							ObjectType:  NewObjectType("any[]", O_ARR),
							elementType: anyType,
						},
					}
				},
			},
		},
	})
}

// 创建元组类型的工厂函数
func NewTupleTypeWithMethods(elementTypes ...Type) *TupleType {
	tupleType := NewTupleType(elementTypes...)
	tupleType.addTupleMethods()
	return tupleType
}

func NewNamedTupleTypeWithMethods(fieldNames []string, elementTypes []Type) *TupleType {
	tupleType := NewNamedTupleType(fieldNames, elementTypes)
	tupleType.addTupleMethods()
	return tupleType
}

// ==================== 集合类型系统 ====================

// SetType 表示集合类型
type SetType struct {
	*ObjectType
	elementType Type
}

func NewSetType(elementType Type) *SetType {
	setType := &SetType{
		ObjectType:  NewObjectType(fmt.Sprintf("set<%s>", elementType.Name()), O_SET),
		elementType: elementType,
	}

	setType.addSetMethods()

	return setType
}

func (s *SetType) Name() string {
	return fmt.Sprintf("set<%s>", s.elementType.Name())
}

func (s *SetType) Kind() ObjKind {
	return O_SET
}

func (s *SetType) Text() string {
	return s.Name()
}

func (s *SetType) AssignableTo(u Type) bool {
	if other, ok := u.(*SetType); ok {
		return s.elementType.AssignableTo(other.elementType)
	}
	return false
}

func (s *SetType) ConvertibleTo(u Type) bool {
	return false
}

func (s *SetType) Implements(u Type) bool {
	return false
}

func (s *SetType) ElementType() Type {
	return s.elementType
}

// 添加集合类型的方法
func (s *SetType) addSetMethods() {
	// length 方法
	s.AddMethod("length", &FunctionType{
		name: "length",
		signatures: []*FunctionSignature{
			{
				positionalParams: []*Parameter{
					{name: "set", typ: s},
				},
				returnType: numberType,
				builtin: func(args ...Value) Value {
					set := args[0].(*SetValue)
					return &NumberValue{Value: big.NewFloat(float64(set.Length()))}
				},
			},
		},
	})

	// add 方法
	s.AddMethod("add", &FunctionType{
		name: "add",
		signatures: []*FunctionSignature{
			{
				positionalParams: []*Parameter{
					{name: "set", typ: s},
					{name: "element", typ: s.elementType},
				},
				returnType: s,
				builtin: func(args ...Value) Value {
					set := args[0].(*SetValue)
					element := args[1]

					err := set.Add(element)
					if err != nil {
						return NULL
					}
					return set
				},
			},
		},
	})

	// remove 方法
	s.AddMethod("remove", &FunctionType{
		name: "remove",
		signatures: []*FunctionSignature{
			{
				positionalParams: []*Parameter{
					{name: "set", typ: s},
					{name: "element", typ: s.elementType},
				},
				returnType: s,
				builtin: func(args ...Value) Value {
					set := args[0].(*SetValue)
					element := args[1]

					err := set.Remove(element)
					if err != nil {
						return NULL
					}
					return set
				},
			},
		},
	})

	// has 方法
	s.AddMethod("has", &FunctionType{
		name: "has",
		signatures: []*FunctionSignature{
			{
				positionalParams: []*Parameter{
					{name: "set", typ: s},
					{name: "element", typ: s.elementType},
				},
				returnType: boolType,
				builtin: func(args ...Value) Value {
					set := args[0].(*SetValue)
					element := args[1]

					has, err := set.Contains(element)
					if err != nil {
						return &BoolValue{Value: false}
					}
					return &BoolValue{Value: has}
				},
			},
		},
	})

	// union 方法
	s.AddMethod("union", &FunctionType{
		name: "union",
		signatures: []*FunctionSignature{
			{
				positionalParams: []*Parameter{
					{name: "set", typ: s},
					{name: "other", typ: s},
				},
				returnType: s,
				builtin: func(args ...Value) Value {
					set := args[0].(*SetValue)
					other := args[1].(*SetValue)

					result, err := set.Union(other)
					if err != nil {
						return NULL
					}
					return result
				},
			},
		},
	})

	// intersect 方法
	s.AddMethod("intersect", &FunctionType{
		name: "intersect",
		signatures: []*FunctionSignature{
			{
				positionalParams: []*Parameter{
					{name: "set", typ: s},
					{name: "other", typ: s},
				},
				returnType: s,
				builtin: func(args ...Value) Value {
					set := args[0].(*SetValue)
					other := args[1].(*SetValue)

					result, err := set.Intersect(other)
					if err != nil {
						return NULL
					}
					return result
				},
			},
		},
	})

	// difference 方法
	s.AddMethod("difference", &FunctionType{
		name: "difference",
		signatures: []*FunctionSignature{
			{
				positionalParams: []*Parameter{
					{name: "set", typ: s},
					{name: "other", typ: s},
				},
				returnType: s,
				builtin: func(args ...Value) Value {
					set := args[0].(*SetValue)
					other := args[1].(*SetValue)

					result, err := set.Difference(other)
					if err != nil {
						return NULL
					}
					return result
				},
			},
		},
	})

	// clear 方法
	s.AddMethod("clear", &FunctionType{
		name: "clear",
		signatures: []*FunctionSignature{
			{
				positionalParams: []*Parameter{
					{name: "set", typ: s},
				},
				returnType: voidType,
				builtin: func(args ...Value) Value {
					set := args[0].(*SetValue)
					set.Clear()
					return NULL
				},
			},
		},
	})

	// to_array 方法
	s.AddMethod("to_array", &FunctionType{
		name: "to_array",
		signatures: []*FunctionSignature{
			{
				positionalParams: []*Parameter{
					{name: "set", typ: s},
				},
				returnType: NewArrayType(s.elementType),
				builtin: func(args ...Value) Value {
					set := args[0].(*SetValue)
					return set.ToArray()
				},
			},
		},
	})
}

// ==================== 枚举类型系统 ====================

// EnumType 表示枚举类型
type EnumType struct {
	*ObjectType
	name         string
	variants     map[string]*EnumVariant // 枚举变体
	variantList  []*EnumVariant
	fields       map[string]*Field // 关联字段（Dart风格）
	isADT        bool              // 是否为代数数据类型（Rust风格）
	discriminant string            // 判别字段名
}

func NewEnumType(name string) *EnumType {
	enumType := &EnumType{
		ObjectType:  NewObjectType(name, O_ENUM),
		name:        name,
		variants:    make(map[string]*EnumVariant),
		variantList: make([]*EnumVariant, 0),
		fields:      make(map[string]*Field),
		isADT:       false,
	}

	enumType.addEnumMethods()

	return enumType
}

func (e *EnumType) Name() string {
	return e.name
}

func (e *EnumType) Kind() ObjKind {
	return O_ENUM
}

func (e *EnumType) Text() string {
	return e.name
}

func (e *EnumType) AssignableTo(u Type) bool {
	if other, ok := u.(*EnumType); ok {
		return e.name == other.name
	}
	return false
}

func (e *EnumType) ConvertibleTo(u Type) bool {
	return false
}

func (e *EnumType) Implements(u Type) bool {
	return false
}

// AddVariant 添加枚举变体
func (e *EnumType) AddVariant(variant *EnumVariant) {
	e.variants[variant.name] = variant
	e.variantList = append(e.variantList, variant)
	variant.enumType = e
}

// GetVariant 获取枚举变体
func (e *EnumType) GetVariant(name string) *EnumVariant {
	return e.variants[name]
}

// GetVariants 获取所有枚举变体
func (e *EnumType) GetVariants() []*EnumVariant {
	return e.variantList
}

// AddField 添加关联字段（Dart风格）
func (e *EnumType) AddField(name string, field *Field) {
	e.fields[name] = field
}

// IsADT 是否为代数数据类型
func (e *EnumType) IsADT() bool {
	return e.isADT
}

// SetADT 设置为代数数据类型
func (e *EnumType) SetADT(discriminant string) {
	e.isADT = true
	e.discriminant = discriminant
}

// 添加枚举类型的方法
func (e *EnumType) addEnumMethods() {
	// values 方法 - 获取所有枚举变体
	e.AddMethod("values", &FunctionType{
		name: "values",
		signatures: []*FunctionSignature{
			{
				positionalParams: []*Parameter{
					{name: "enum", typ: anyType}, // 改为anyType以接受枚举变体
				},
				returnType: &ArrayType{
					ObjectType:  NewObjectType("enum_variant[]", O_ARR),
					elementType: e,
				},
				builtin: func(args ...Value) Value {
					enumValue := args[0].(*EnumVariant)
					enum := enumValue.enumType
					variants := enum.GetVariants()

					result := make([]Value, len(variants))
					for i, variant := range variants {
						result[i] = variant
					}

					return &ArrayValue{
						Elements: result,
						arrayType: &ArrayType{
							ObjectType:  NewObjectType("enum_variant[]", O_ARR),
							elementType: e,
						},
					}
				},
			},
		},
	})

	// from_value 方法 - 根据值获取枚举
	e.AddMethod("from_value", &FunctionType{
		name: "from_value",
		signatures: []*FunctionSignature{
			{
				positionalParams: []*Parameter{
					{name: "enum", typ: anyType}, // 改为anyType以接受枚举变体
					{name: "value", typ: anyType},
				},
				returnType: NewUnionType(e, nullType),
				builtin: func(args ...Value) Value {
					enumValue := args[0].(*EnumVariant)
					enum := enumValue.enumType
					value := args[1]

					for _, variant := range enum.GetVariants() {
						if variant.value.Text() == value.Text() {
							return variant
						}
					}
					return NULL
				},
			},
		},
	})

	// is_valid 方法 - 检查值是否为有效枚举
	e.AddMethod("is_valid", &FunctionType{
		name: "is_valid",
		signatures: []*FunctionSignature{
			{
				positionalParams: []*Parameter{
					{name: "enum", typ: anyType}, // 改为anyType以接受枚举变体
					{name: "value", typ: anyType},
				},
				returnType: boolType,
				builtin: func(args ...Value) Value {
					enumValue := args[0].(*EnumVariant)
					enum := enumValue.enumType
					value := args[1]

					for _, variant := range enum.GetVariants() {
						if variant.value.Text() == value.Text() {
							return &BoolValue{Value: true}
						}
					}
					return &BoolValue{Value: false}
				},
			},
		},
	})
}

// 创建标准数字枚举
func NewNumericEnum(name string, valueNames []string) *EnumType {
	enumType := NewEnumType(name)

	for i, valueName := range valueNames {
		variant := NewEnumVariant(valueName, &NumberValue{Value: big.NewFloat(float64(i))})
		enumType.AddVariant(variant)
	}

	return enumType
}

// 创建显式数字枚举
func NewExplicitNumericEnum(name string, valueMap map[string]int) *EnumType {
	enumType := NewEnumType(name)

	for valueName, numValue := range valueMap {
		variant := NewEnumVariant(valueName, &NumberValue{Value: big.NewFloat(float64(numValue))})
		enumType.AddVariant(variant)
	}

	return enumType
}

// 创建字符串枚举
func NewStringEnum(name string, valueMap map[string]string) *EnumType {
	enumType := NewEnumType(name)

	for valueName, strValue := range valueMap {
		variant := NewEnumVariant(valueName, &StringValue{Value: strValue})
		enumType.AddVariant(variant)
	}

	return enumType
}

// 创建混合类型枚举
func NewMixedEnum(name string, valueMap map[string]Value) *EnumType {
	enumType := NewEnumType(name)

	for valueName, value := range valueMap {
		variant := NewEnumVariant(valueName, value)
		enumType.AddVariant(variant)
	}

	return enumType
}

// 创建Dart风格的关联值枚举
func NewAssociatedEnum(name string, fields map[string]*Field, valueMap map[string]Value) *EnumType {
	enumType := NewEnumType(name)

	// 添加关联字段
	for fieldName, field := range fields {
		enumType.AddField(fieldName, field)
	}

	// 添加枚举变体
	for valueName, value := range valueMap {
		variant := NewEnumVariant(valueName, value)
		enumType.AddVariant(variant)
	}

	return enumType
}

// 创建Rust风格的代数数据类型枚举
func NewADTEnum(name string, variants map[string]map[string]Type) *EnumType {
	enumType := NewEnumType(name)
	enumType.SetADT("variant")

	for variantName, fields := range variants {
		// 创建变体
		variant := NewEnumVariant(variantName, &StringValue{Value: variantName})

		// 添加ADT字段类型
		for fieldName, fieldType := range fields {
			variant.AddADTFieldType(fieldName, fieldType)
		}

		enumType.AddVariant(variant)
	}

	return enumType
}

// 示例：枚举类型使用
func ExampleEnumUsage() {
	// 1. 标准数字枚举
	statusEnum := NewNumericEnum("Status", []string{"Pending", "Approved", "Rejected"})
	fmt.Printf("Status enum: %s\n", statusEnum.Name())

	pending := statusEnum.GetVariant("Pending")
	fmt.Printf("Pending value: %s\n", pending.value.Text()) // 0

	// 2. 显式数字枚举
	httpCodeEnum := NewExplicitNumericEnum("HttpCode", map[string]int{
		"OK":             200,
		"NotFound":       404,
		"ServerError":    500,
		"GatewayTimeout": 501,
	})

	ok := httpCodeEnum.GetVariant("OK")
	fmt.Printf("OK value: %s\n", ok.value.Text()) // 200

	// 3. 字符串枚举
	directionEnum := NewStringEnum("Direction", map[string]string{
		"North": "N",
		"South": "S",
		"East":  "E",
		"West":  "W",
	})

	north := directionEnum.GetVariant("North")
	fmt.Printf("North value: %s\n", north.value.Text()) // N

	// 4. 混合类型枚举
	configEnum := NewMixedEnum("Config", map[string]Value{
		"RetryCount":    &NumberValue{Value: big.NewFloat(3)},
		"Timeout":       &StringValue{Value: "30s"},
		"EnableLogging": &BoolValue{Value: true},
	})

	retryCount := configEnum.GetVariant("RetryCount")
	fmt.Printf("RetryCount value: %s\n", retryCount.value.Text()) // 3

	// 5. Dart风格的关联值枚举
	protocolEnum := NewAssociatedEnum("Protocol", map[string]*Field{
		"port": {ft: numberType, mods: []FieldModifier{FieldModifierRequired}},
	}, map[string]Value{
		"tcp": &NumberValue{Value: big.NewFloat(6)},
		"udp": &NumberValue{Value: big.NewFloat(17)},
	})

	tcp := protocolEnum.GetVariant("tcp")
	tcp.SetField("port", &NumberValue{Value: big.NewFloat(6)})
	fmt.Printf("TCP: %s\n", tcp.Text())                       // tcp(port: 6)
	fmt.Printf("TCP port: %s\n", tcp.GetField("port").Text()) // 6

	// 6. Rust风格的代数数据类型枚举
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

	tcpPacket := networkPacketEnum.GetVariant("TCP")
	tcpPacket.SetADTField("src_port", &NumberValue{Value: big.NewFloat(8080)})
	tcpPacket.SetADTField("dst_port", &NumberValue{Value: big.NewFloat(80)})
	fmt.Printf("TCP packet: %s\n", tcpPacket.Text()) // TCP{src_port: 8080, dst_port: 80}

	udpPacket := networkPacketEnum.GetVariant("UDP")
	udpPacket.SetADTField("port", &NumberValue{Value: big.NewFloat(53)})
	udpPacket.SetADTField("payload", &StringValue{Value: "DNS query"})
	fmt.Printf("UDP packet: %s\n", udpPacket.Text()) // UDP{port: 53, payload: DNS query}
}

// GetValue 获取枚举变体（别名，保持向后兼容性）
func (e *EnumType) GetValue(name string) *EnumVariant {
	return e.GetVariant(name)
}

// GetValues 获取所有枚举变体（别名，保持向后兼容性）
func (e *EnumType) GetValues() []*EnumVariant {
	return e.GetVariants()
}

// 示例：静态方法和继承
func ExampleStaticMethodsAndInheritance() {
	// 1. 创建Counter类（静态方法示例）
	counterClass := NewClassType("Counter")

	// 添加静态字段
	counterClass.AddStaticField("count", &Field{
		ft:   numberType,
		mods: []FieldModifier{FieldModifierStatic},
	})

	// 设置静态字段的初始值
	counterClass.SetStaticValue("count", &NumberValue{Value: big.NewFloat(0)})

	// 添加静态方法
	counterClass.AddStaticMethod("increment", &FunctionType{
		name: "increment",
		signatures: []*FunctionSignature{
			{
				positionalParams: []*Parameter{},
				returnType:       numberType,
				builtin: func(args ...Value) Value {
					// 获取当前静态值
					currentCount := counterClass.GetStaticValue("count").(*NumberValue)
					newCount := &NumberValue{Value: big.NewFloat(0)}
					newCount.Value.Add(currentCount.Value, big.NewFloat(1))

					// 更新静态值
					counterClass.SetStaticValue("count", newCount)
					return newCount
				},
			},
		},
	})

	// 2. 创建Point类（继承示例）
	pointClass := NewClassType("Point")

	// 添加字段
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

	// 3. 创建ColorPoint类（继承Point）
	colorPointClass := NewClassType("ColorPoint")
	colorPointClass.SetParent(pointClass)

	// 添加新字段
	colorPointClass.AddField("color", &Field{
		ft:   stringType,
		mods: []FieldModifier{FieldModifierRequired},
	})

	// 重写to_str方法
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

	// 4. 使用示例
	fmt.Println("=== 静态方法示例 ===")

	// 调用静态方法
	incrementMethod := counterClass.GetStaticMethod("increment")
	result1, _ := incrementMethod.Call([]Value{}, nil)
	fmt.Printf("Counter.increment(): %s\n", result1.Text()) // 1

	result2, _ := incrementMethod.Call([]Value{}, nil)
	fmt.Printf("Counter.increment(): %s\n", result2.Text()) // 2

	// 访问静态字段
	staticCount := counterClass.GetStaticValue("count")
	fmt.Printf("Counter.count: %s\n", staticCount.Text()) // 2

	fmt.Println("\n=== 继承示例 ===")

	// 创建Point实例
	point := &ClassInstance{
		clsType: pointClass,
		fields: map[string]Value{
			"x": &NumberValue{Value: big.NewFloat(10)},
			"y": &NumberValue{Value: big.NewFloat(20)},
		},
	}

	// 调用Point的方法
	pointToStr := pointClass.MethodByName("to_str")
	pointResult, _ := pointToStr.Call([]Value{point}, nil)
	fmt.Printf("Point.to_str(): %s\n", pointResult.Text()) // Point(10, 20)

	// 创建ColorPoint实例
	colorPoint := &ClassInstance{
		clsType: colorPointClass,
		fields: map[string]Value{
			"x":     &NumberValue{Value: big.NewFloat(30)},
			"y":     &NumberValue{Value: big.NewFloat(40)},
			"color": &StringValue{Value: "red"},
		},
	}

	// 调用重写的方法
	colorPointToStr := colorPointClass.MethodByName("to_str")
	colorPointResult, _ := colorPointToStr.Call([]Value{colorPoint}, nil)
	fmt.Printf("ColorPoint.to_str(): %s\n", colorPointResult.Text()) // ColorPoint(30, 40, red)

	// 检查继承关系
	fmt.Printf("ColorPoint is subclass of Point: %t\n", colorPointClass.IsSubclassOf(pointClass)) // true
}

// AddTraitConstraint 为类型参数添加trait约束
func (ct *ClassType) AddTraitConstraint(typeParam string, trait *TraitType) {
	if ct.traitConstraints[typeParam] == nil {
		ct.traitConstraints[typeParam] = make([]*TraitType, 0)
	}
	ct.traitConstraints[typeParam] = append(ct.traitConstraints[typeParam], trait)
}

// GetTraitConstraints 获取类型参数的trait约束
func (ct *ClassType) GetTraitConstraints(typeParam string) []*TraitType {
	return ct.traitConstraints[typeParam]
}

// CheckTraitConstraints 检查类型参数是否满足trait约束
func (ct *ClassType) CheckTraitConstraints(typeArgs []Type) error {
	if !ct.IsGeneric() {
		return nil
	}

	for i, typeArg := range typeArgs {
		if i >= len(ct.TypeParameters()) {
			break
		}
		typeParam := ct.TypeParameters()[i]
		constraints := ct.GetTraitConstraints(typeParam)

		for _, constraint := range constraints {
			if !ct.satisfiesTraitConstraint(typeArg, constraint) {
				return fmt.Errorf("type %s does not satisfy trait constraint %s",
					typeArg.Name(), constraint.Name())
			}
		}
	}
	return nil
}

// satisfiesTraitConstraint 检查类型是否满足trait约束
func (ct *ClassType) satisfiesTraitConstraint(t Type, trait *TraitType) bool {
	// 检查类型是否实现了该trait
	if classType, ok := t.(*ClassType); ok {
		return classType.ImplementsTrait(trait)
	}

	// 检查是否有trait实现
	impls := GlobalTraitRegistry.GetImpl(trait.Name(), t.Name())
	return len(impls) > 0
}

// 重写Instantiate方法以支持trait约束检查
func (ct *ClassType) InstantiateWithTraitConstraints(typeArgs []Type) (*ClassType, error) {
	// 检查trait约束
	if err := ct.CheckTraitConstraints(typeArgs); err != nil {
		return nil, err
	}

	// 调用原有的实例化逻辑
	return ct.Instantiate(typeArgs)
}

// 示例：trait、泛型和继承联动
func ExampleTraitGenericInheritanceIntegration() {
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
		ft:   anyType, // 这里应该是泛型类型T，简化处理
		mods: []FieldModifier{FieldModifierRequired},
	})

	// 添加方法
	containerClass.AddMethod("get_item", &FunctionType{
		name: "get_item",
		signatures: []*FunctionSignature{
			{
				positionalParams: []*Parameter{
					{name: "self", typ: containerClass},
				},
				returnType: anyType,
				builtin: func(args ...Value) Value {
					container := args[0].(*ClassInstance)
					return container.GetField("item")
				},
			},
		},
	})

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

					// 调用item的to_str方法（通过trait约束保证存在）
					if classInstance, ok := item.(*ClassInstance); ok {
						if method := classInstance.clsType.MethodByName("to_str"); method != nil {
							result, _ := method.Call([]Value{classInstance}, nil)
							return result
						}
					}
					return &StringValue{Value: "unknown"}
				},
			},
		},
	})

	// 4. 创建继承的泛型类
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

	// 5. 使用示例
	fmt.Println("=== Trait、泛型、继承联动示例 ===")

	// 创建Point实例
	point := &ClassInstance{
		clsType: pointClass,
		fields: map[string]Value{
			"x": &NumberValue{Value: big.NewFloat(10)},
			"y": &NumberValue{Value: big.NewFloat(20)},
		},
	}

	// 创建ColorPoint实例
	colorPoint := &ClassInstance{
		clsType: colorPointClass,
		fields: map[string]Value{
			"x":     &NumberValue{Value: big.NewFloat(30)},
			"y":     &NumberValue{Value: big.NewFloat(40)},
			"color": &StringValue{Value: "red"},
		},
	}

	// 测试trait实现
	fmt.Println("Point Display trait:")
	pointDisplayMethod := pointClass.MethodByName("to_str")
	pointDisplayResult, _ := pointDisplayMethod.Call([]Value{point}, nil)
	fmt.Printf("  %s\n", pointDisplayResult.Text())

	fmt.Println("ColorPoint Display trait (inherited):")
	colorPointDisplayMethod := colorPointClass.MethodByName("to_str")
	colorPointDisplayResult, _ := colorPointDisplayMethod.Call([]Value{colorPoint}, nil)
	fmt.Printf("  %s\n", colorPointDisplayResult.Text())

	// 测试泛型容器实例化
	fmt.Println("\n泛型容器实例化:")

	// 实例化Container<Point>
	pointContainerInstance, err := containerClass.InstantiateWithTraitConstraints([]Type{pointClass})
	if err != nil {
		fmt.Printf("Container<Point> 实例化失败: %v\n", err)
	} else {
		fmt.Printf("Container<Point> 实例化成功\n")

		// 创建容器实例
		pointContainer := &ClassInstance{
			clsType: pointContainerInstance,
			fields: map[string]Value{
				"item": point,
			},
		}

		// 测试容器方法
		displayMethod := pointContainerInstance.MethodByName("display_item")
		displayResult, _ := displayMethod.Call([]Value{pointContainer}, nil)
		fmt.Printf("  Container<Point>.display_item(): %s\n", displayResult.Text())
	}

	// 实例化Container<ColorPoint>
	colorPointContainerInstance, err := containerClass.InstantiateWithTraitConstraints([]Type{colorPointClass})
	if err != nil {
		fmt.Printf("Container<ColorPoint> 实例化失败: %v\n", err)
	} else {
		fmt.Printf("Container<ColorPoint> 实例化成功\n")

		// 创建容器实例
		colorPointContainer := &ClassInstance{
			clsType: colorPointContainerInstance,
			fields: map[string]Value{
				"item": colorPoint,
			},
		}

		// 测试容器方法
		displayMethod := colorPointContainerInstance.MethodByName("display_item")
		displayResult, _ := displayMethod.Call([]Value{colorPointContainer}, nil)
		fmt.Printf("  Container<ColorPoint>.display_item(): %s\n", displayResult.Text())
	}

	// 测试继承关系
	fmt.Printf("\n继承关系检查:\n")
	fmt.Printf("  ColorPoint is subclass of Point: %t\n", colorPointClass.IsSubclassOf(pointClass))
	fmt.Printf("  Point implements Display: %t\n", pointClass.ImplementsTrait(displayTrait))
	fmt.Printf("  ColorPoint implements Display: %t\n", colorPointClass.ImplementsTrait(displayTrait))
	fmt.Printf("  Point implements Clone: %t\n", pointClass.ImplementsTrait(cloneTrait))
	fmt.Printf("  ColorPoint implements Clone: %t\n", colorPointClass.ImplementsTrait(cloneTrait))
}

// 示例：泛型trait约束继承
func ExampleGenericTraitConstraintInheritance() {
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

	// 2. 定义泛型trait
	comparableTrait := NewTraitType("Comparable")
	comparableTrait.AddMethod("compare", &FunctionType{
		name: "compare",
		signatures: []*FunctionSignature{
			{
				positionalParams: []*Parameter{
					{name: "self", typ: comparableTrait},
					{name: "other", typ: comparableTrait},
				},
				returnType: numberType,
			},
		},
	})

	// 3. 定义继承trait约束的泛型类
	// class SortedContainer<T: Display + Comparable> extends Container<T>
	containerClass := NewGenericClass("Container", []string{"T"}, nil)
	containerClass.AddTraitConstraint("T", displayTrait)

	sortedContainerClass := NewGenericClass("SortedContainer", []string{"T"}, nil)
	sortedContainerClass.SetParent(containerClass)
	sortedContainerClass.AddTraitConstraint("T", displayTrait)
	sortedContainerClass.AddTraitConstraint("T", comparableTrait)

	// 4. 创建满足约束的类
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

	// 为Point实现Comparable trait
	pointComparableImpl := NewTraitImpl(comparableTrait, pointClass)
	pointComparableImpl.AddMethod("compare", &FunctionType{
		name: "compare",
		signatures: []*FunctionSignature{
			{
				positionalParams: []*Parameter{
					{name: "self", typ: pointClass},
					{name: "other", typ: pointClass},
				},
				returnType: numberType,
				builtin: func(args ...Value) Value {
					self := args[0].(*ClassInstance)
					other := args[1].(*ClassInstance)

					selfX := self.GetField("x").(*NumberValue)
					otherX := other.GetField("x").(*NumberValue)

					if selfX.Value.Cmp(otherX.Value) < 0 {
						return &NumberValue{Value: big.NewFloat(-1)}
					} else if selfX.Value.Cmp(otherX.Value) > 0 {
						return &NumberValue{Value: big.NewFloat(1)}
					}
					return &NumberValue{Value: big.NewFloat(0)}
				},
			},
		},
	})
	GlobalTraitRegistry.RegisterImpl(pointComparableImpl)

	// 5. 测试泛型约束继承
	fmt.Println("=== 泛型trait约束继承示例 ===")

	// 测试Container<Point>实例化
	_, err := containerClass.InstantiateWithTraitConstraints([]Type{pointClass})
	if err != nil {
		fmt.Printf("Container<Point> 实例化失败: %v\n", err)
	} else {
		fmt.Printf("Container<Point> 实例化成功\n")
	}

	// 测试SortedContainer<Point>实例化
	_, err = sortedContainerClass.InstantiateWithTraitConstraints([]Type{pointClass})
	if err != nil {
		fmt.Printf("SortedContainer<Point> 实例化失败: %v\n", err)
	} else {
		fmt.Printf("SortedContainer<Point> 实例化成功\n")
	}

	// 测试不满足约束的类型
	simpleClass := NewClassType("Simple")
	_, err = sortedContainerClass.InstantiateWithTraitConstraints([]Type{simpleClass})
	if err != nil {
		fmt.Printf("SortedContainer<Simple> 实例化失败（预期）: %v\n", err)
	} else {
		fmt.Printf("SortedContainer<Simple> 实例化成功（意外）\n")
	}

	// 6. 测试继承关系
	fmt.Printf("\n继承关系检查:\n")
	fmt.Printf("  SortedContainer<T> is subclass of Container<T>: %t\n",
		sortedContainerClass.IsSubclassOf(containerClass))

	// 7. 测试trait约束
	fmt.Printf("\nTrait约束检查:\n")
	fmt.Printf("  Point implements Display: %t\n", pointClass.ImplementsTrait(displayTrait))
	fmt.Printf("  Point implements Comparable: %t\n", pointClass.ImplementsTrait(comparableTrait))
	fmt.Printf("  Simple implements Display: %t\n", simpleClass.ImplementsTrait(displayTrait))
	fmt.Printf("  Simple implements Comparable: %t\n", simpleClass.ImplementsTrait(comparableTrait))
}

// 装饰器相关类型定义
type DecoratorTarget int

const (
	DecoratorTargetClass DecoratorTarget = iota
	DecoratorTargetMethod
	DecoratorTargetField
	DecoratorTargetFunction
	DecoratorTargetParameter
	DecoratorTargetEnum
	DecoratorTargetTrait
)

// DecoratorProcessor 装饰器处理器接口
type DecoratorProcessor interface {
	Process(target any, args []Value) (any, error)
}

// DecoratorType 装饰器类型
type DecoratorType struct {
	*ObjectType
	name      string
	target    DecoratorTarget
	processor DecoratorProcessor
}

// NewDecoratorType 创建新的装饰器类型
func NewDecoratorType(name string, target DecoratorTarget, processor DecoratorProcessor) *DecoratorType {
	return &DecoratorType{
		ObjectType: NewObjectType(name, O_DECORATOR),
		name:       name,
		target:     target,
		processor:  processor,
	}
}

// Process 处理装饰器
func (d *DecoratorType) Process(target any, args []Value) (any, error) {
	if d.processor != nil {
		return d.processor.Process(target, args)
	}
	return target, nil
}

// Target 获取装饰器目标类型
func (d *DecoratorType) Target() DecoratorTarget {
	return d.target
}

// 装饰器注册表
type DecoratorRegistry struct {
	decorators map[string]*DecoratorType
	mutex      sync.RWMutex
}

// NewDecoratorRegistry 创建装饰器注册表
func NewDecoratorRegistry() *DecoratorRegistry {
	return &DecoratorRegistry{
		decorators: make(map[string]*DecoratorType),
	}
}

// RegisterDecorator 注册装饰器
func (dr *DecoratorRegistry) RegisterDecorator(decorator *DecoratorType) {
	dr.mutex.Lock()
	defer dr.mutex.Unlock()
	dr.decorators[decorator.name] = decorator
}

// GetDecorator 获取装饰器
func (dr *DecoratorRegistry) GetDecorator(name string) *DecoratorType {
	dr.mutex.RLock()
	defer dr.mutex.RUnlock()
	return dr.decorators[name]
}

// GlobalDecoratorRegistry 全局装饰器注册表
var GlobalDecoratorRegistry = NewDecoratorRegistry()

// AddConstraint 添加类型参数约束
func (ct *ClassType) AddConstraint(typeParam string, constraint Constraint) {
	if ct.GenericBase.constraints == nil {
		ct.GenericBase.constraints = make(map[string]Constraint)
	}
	ct.GenericBase.constraints[typeParam] = constraint
}

// AddTypeConstraint 添加类型约束（便捷方法）
func (ct *ClassType) AddTypeConstraint(typeParam string, targetType Type) {
	ct.AddConstraint(typeParam, NewInterfaceConstraint(targetType))
}

// AddOperatorConstraint 添加运算符约束（便捷方法）
func (ct *ClassType) AddOperatorConstraint(typeParam string, op Operator, rightType, returnType Type) {
	ct.AddConstraint(typeParam, NewOperatorConstraint(op, rightType, returnType))
}

// AddCompositeConstraint 添加复合约束（便捷方法）
func (ct *ClassType) AddCompositeConstraint(typeParam string, constraints ...Constraint) {
	ct.AddConstraint(typeParam, NewCompositeConstraint(constraints...))
}

// AddUnionConstraint 添加联合约束（便捷方法）
func (ct *ClassType) AddUnionConstraint(typeParam string, constraints ...Constraint) {
	ct.AddConstraint(typeParam, NewUnionConstraint(constraints...))
}

// GetConstraint 获取类型参数的约束
func (ct *ClassType) GetConstraint(typeParam string) Constraint {
	return ct.Constraints()[typeParam]
}

// HasConstraint 检查类型参数是否有约束
func (ct *ClassType) HasConstraint(typeParam string) bool {
	_, exists := ct.Constraints()[typeParam]
	return exists
}

// CheckConstraints 检查类型参数是否满足所有约束
func (ct *ClassType) CheckConstraints(typeArgs []Type) error {
	if !ct.isGeneric || ct.constraints == nil {
		return nil
	}

	for i, param := range ct.typeParams {
		if i >= len(typeArgs) {
			break
		}
		if constraint, exists := ct.constraints[param]; exists {
			if !constraint.SatisfiedBy(typeArgs[i]) {
				return fmt.Errorf("type %s does not satisfy constraint %s for parameter %s",
					typeArgs[i].Name(), constraint.String(), param)
			}
		}
	}
	return nil
}

// 示例：将Constraint系统集成到泛型系统中
func ExampleConstraintIntegration() {
	fmt.Println("=== Constraint系统集成示例 ===")

	// 1. 创建基础trait
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

	// 2. 创建约束
	displayConstraint := NewInterfaceConstraint(displayTrait)
	cloneConstraint := NewInterfaceConstraint(cloneTrait)

	// 复合约束：T 必须同时实现 Display 和 Clone
	displayAndCloneConstraint := NewCompositeConstraint(displayConstraint, cloneConstraint)

	// 联合约束：T 必须实现 Display 或 Clone
	_ = NewUnionConstraint(displayConstraint, cloneConstraint)

	// 3. 使用不同的方式创建泛型类

	// 方式1：使用NewGenericClass + AddConstraint
	fmt.Println("\n1. 使用 NewGenericClass + AddConstraint:")
	container1 := NewGenericClass("Container", []string{"T"}, nil)
	container1.AddConstraint("T", displayAndCloneConstraint)

	// 方式2：使用便捷方法
	fmt.Println("\n2. 使用便捷方法:")
	container2 := NewGenericClass("Container", []string{"T"}, nil)
	container2.AddCompositeConstraint("T", displayConstraint, cloneConstraint)

	// 方式3：使用NewConstrainedGenericClass
	fmt.Println("\n3. 使用 NewConstrainedGenericClass:")
	_ = NewConstrainedGenericClass("Container", []string{"T"}, map[string]Constraint{
		"T": displayAndCloneConstraint,
	})

	// 4. 创建满足约束的类
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

	// 5. 测试约束检查
	fmt.Println("\n4. 测试约束检查:")

	// 测试满足约束的类型
	_, err := container1.Instantiate([]Type{pointClass})
	if err != nil {
		fmt.Printf("  Point 实例化失败: %v\n", err)
	} else {
		fmt.Printf("  Point 实例化成功\n")
	}

	// 测试不满足约束的类型
	simpleClass := NewClassType("Simple")
	_, err = container1.Instantiate([]Type{simpleClass})
	if err != nil {
		fmt.Printf("  Simple 实例化失败（预期）: %v\n", err)
	} else {
		fmt.Printf("  Simple 实例化成功（意外）\n")
	}

	// 6. 测试不同类型的约束
	fmt.Println("\n5. 测试不同类型的约束:")

	// 测试数值约束
	numberContainer := NewGenericClass("NumberContainer", []string{"T"}, nil)
	numberContainer.AddConstraint("T", NumericConstraint)

	_, err = numberContainer.Instantiate([]Type{numberType})
	if err != nil {
		fmt.Printf("  NumberContainer<number> 实例化失败: %v\n", err)
	} else {
		fmt.Printf("  NumberContainer<number> 实例化成功\n")
	}

	_, err = numberContainer.Instantiate([]Type{stringType})
	if err != nil {
		fmt.Printf("  NumberContainer<string> 实例化失败（预期）: %v\n", err)
	} else {
		fmt.Printf("  NumberContainer<string> 实例化成功（意外）\n")
	}

	// 测试运算符约束
	comparableContainer := NewGenericClass("ComparableContainer", []string{"T"}, nil)
	comparableContainer.AddConstraint("T", ComparableConstraint)

	_, err = comparableContainer.Instantiate([]Type{numberType})
	if err != nil {
		fmt.Printf("  ComparableContainer<number> 实例化失败: %v\n", err)
	} else {
		fmt.Printf("  ComparableContainer<number> 实例化成功\n")
	}

	// 7. 测试约束查询
	fmt.Println("\n6. 测试约束查询:")

	fmt.Printf("  container1 有约束 T: %t\n", container1.HasConstraint("T"))
	fmt.Printf("  container1 的约束 T: %s\n", container1.GetConstraint("T").String())

	// 8. 测试约束检查方法
	fmt.Println("\n7. 测试约束检查方法:")

	err = container1.CheckConstraints([]Type{pointClass})
	if err != nil {
		fmt.Printf("  约束检查失败: %v\n", err)
	} else {
		fmt.Printf("  约束检查通过\n")
	}

	err = container1.CheckConstraints([]Type{simpleClass})
	if err != nil {
		fmt.Printf("  约束检查失败（预期）: %v\n", err)
	} else {
		fmt.Printf("  约束检查通过（意外）\n")
	}
}

// 示例：高级约束组合
func ExampleAdvancedConstraintCombination() {
	fmt.Println("=== 高级约束组合示例 ===")

	// 1. 创建复杂的约束组合
	// T: (Display + Clone) | (Numeric + Comparable)
	displayTrait := NewTraitType("Display")
	cloneTrait := NewTraitType("Clone")

	displayConstraint := NewInterfaceConstraint(displayTrait)
	cloneConstraint := NewInterfaceConstraint(cloneTrait)
	numericConstraint := NumericConstraint
	comparableConstraint := ComparableConstraint

	// 组合1：Display + Clone
	displayAndClone := NewCompositeConstraint(displayConstraint, cloneConstraint)

	// 组合2：Numeric + Comparable
	numericAndComparable := NewCompositeConstraint(numericConstraint, comparableConstraint)

	// 最终约束：组合1 OR 组合2
	finalConstraint := NewUnionConstraint(displayAndClone, numericAndComparable)

	// 2. 创建使用复杂约束的泛型类
	advancedContainer := NewGenericClass("AdvancedContainer", []string{"T"}, nil)
	advancedContainer.AddConstraint("T", finalConstraint)

	// 3. 测试不同类型的实例化
	fmt.Printf("复杂约束: %s\n", finalConstraint.String())

	// 测试满足组合1的类型
	displayCloneClass := NewClassType("DisplayClone")
	// ... 实现Display和Clone trait

	// 测试满足组合2的类型
	_ = numberType // 已经满足Numeric和Comparable

	// 测试不满足任何组合的类型
	simpleClass := NewClassType("Simple")

	// 这里可以添加实际的测试逻辑
	_ = advancedContainer
	_ = displayCloneClass
	_ = simpleClass
}

// ==================== 统一泛型接口 ====================

// Generic 统一的泛型接口
type Generic interface {
	Type
	// 泛型相关方法
	IsGeneric() bool
	TypeParameters() []string
	Constraints() map[string]Constraint
	Instantiate(typeArgs []Type) (Type, error)
	CheckConstraints(typeArgs []Type) error
}

// GenericBase 泛型基础结构，可以被其他类型嵌入
type GenericBase struct {
	isGeneric   bool
	typeParams  []string
	constraints map[string]Constraint
	instances   map[string]Type
}

func NewGenericBase(typeParams []string, constraints map[string]Constraint) *GenericBase {
	return &GenericBase{
		isGeneric:   len(typeParams) > 0,
		typeParams:  typeParams,
		constraints: constraints,
		instances:   make(map[string]Type),
	}
}

func (gb *GenericBase) IsGeneric() bool {
	return gb.isGeneric
}

func (gb *GenericBase) TypeParameters() []string {
	return gb.typeParams
}

func (gb *GenericBase) Constraints() map[string]Constraint {
	return gb.constraints
}

func (gb *GenericBase) CheckConstraints(typeArgs []Type) error {
	if !gb.isGeneric || gb.constraints == nil {
		return nil
	}

	for i, param := range gb.typeParams {
		if i >= len(typeArgs) {
			break
		}
		if constraint, exists := gb.constraints[param]; exists {
			if !constraint.SatisfiedBy(typeArgs[i]) {
				return fmt.Errorf("type %s does not satisfy constraint %s for parameter %s",
					typeArgs[i].Name(), constraint.String(), param)
			}
		}
	}
	return nil
}

func (gb *GenericBase) generateTypeKey(typeArgs []Type) string {
	keys := make([]string, len(typeArgs))
	for i, t := range typeArgs {
		keys[i] = t.Name()
	}
	return strings.Join(keys, ",")
}

func (gb *GenericBase) getCachedInstance(typeKey string) Type {
	if instance, exists := gb.instances[typeKey]; exists {
		return instance
	}
	return nil
}

func (gb *GenericBase) cacheInstance(typeKey string, instance Type) {
	gb.instances[typeKey] = instance
}

func (ft *FunctionType) AppendSignatures(sigs []*FunctionSignature) {
	ft.signatures = append(ft.signatures, sigs...)
}

func (ft *FunctionType) Signatures() []*FunctionSignature {
	return ft.signatures
}
