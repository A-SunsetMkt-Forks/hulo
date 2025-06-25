// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package object

import (
	"fmt"
	"strings"

	"github.com/hulo-lang/hulo/syntax/hulo/ast"
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

type StringType struct {
	*ObjectType
}

func NewStringType() *StringType {
	strType := &StringType{
		ObjectType: NewObjectType("string", O_STR),
	}

	strType.addStringMethods()

	return strType
}

func (s *StringType) addStringMethods() {
	s.AddMethod("length", nil)
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
		ObjectType: NewObjectType("number", O_NUM),
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
	astDecl     *ast.ClassDecl
	ctors       []Constructor
	isGeneric   bool
	typeParams  []string              // e.g. T, U
	constraints map[string]Type       // e.g. T extends num | str
	instances   map[string]*ClassType // 已实例化的类型缓存
	registry    *Registry
}

func NewGenericClass(name string, typeParams []string, constraints map[string]Type) *ClassType {
	return &ClassType{
		ObjectType:  NewObjectType(name, O_CLASS),
		isGeneric:   true,
		typeParams:  typeParams,
		constraints: constraints,
		instances:   make(map[string]*ClassType),
	}
}

func NewClassType(name string) *ClassType {
	return &ClassType{
		ObjectType:  NewObjectType(name, O_CLASS),
		isGeneric:   false,
		typeParams:  nil,
		constraints: nil,
		instances:   nil,
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

func (ct *ClassType) Instantiate(typeArgs []Type) (*ClassType, error) {
	if !ct.isGeneric {
		return ct, nil
	}

	if len(typeArgs) != len(ct.typeParams) {
		return nil, fmt.Errorf("type argument count mismatch: expected %d, got %d",
			len(ct.typeParams), len(typeArgs))
	}

	// 检查类型约束
	for i, param := range ct.typeParams {
		if constraint, exists := ct.constraints[param]; exists {
			if !typeArgs[i].Implements(constraint) {
				return nil, fmt.Errorf("type %s does not satisfy constraint for %s",
					typeArgs[i].Name(), param)
			}
		}
	}

	// 生成实例化类型的唯一标识
	typeKey := ct.generateTypeKey(typeArgs)

	// 检查缓存
	if instance, exists := ct.instances[typeKey]; exists {
		return instance, nil
	}

	// 创建新的实例化类型
	instance := &ClassType{
		ObjectType:  NewObjectType(ct.Name()+"<"+typeKey+">", O_CLASS),
		astDecl:     ct.astDecl,
		isGeneric:   false, // 实例化后不再是泛型
		typeParams:  nil,
		constraints: nil,
		instances:   nil,
	}

	// 替换AST中的类型参数
	instance.ctors = ct.instantiateConstructors(typeArgs)

	// 缓存实例
	ct.instances[typeKey] = instance

	return instance, nil
}

func (ct *ClassType) generateTypeKey(typeArgs []Type) string {
	keys := make([]string, len(typeArgs))
	for i, t := range typeArgs {
		keys[i] = t.Name()
	}
	return strings.Join(keys, ",")
}

func (ct *ClassType) instantiateConstructors(typeArgs []Type) []Constructor {
	typeMap := make(map[string]Type)
	for i, param := range ct.typeParams {
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
	// 要比较函数签名
	return nil
}

func (ct *ClassType) getDefaultConstructor() Constructor {
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

var _ Type = &FunctionType{}

type FunctionType struct {
	name        string
	signatures  []*FunctionSignature
	isGeneric   bool
	typeParams  []string
	constraints map[string]Type
	visible     bool
}

func NewGenericFunction(name string, typeParams []string, constraints map[string]Type) *FunctionType {
	return &FunctionType{
		name:        name,
		signatures:  make([]*FunctionSignature, 0),
		isGeneric:   true,
		typeParams:  typeParams,
		constraints: constraints,
		visible:     true,
	}
}

func (ft *FunctionType) Name() string {
	return ft.name
}

func (ft *FunctionType) Kind() ObjKind {
	return O_FUNC
}

func (ft *FunctionType) Text() string {
	return "fn"
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

func (ft *FunctionType) Call(args []Value, namedArgs map[string]Value) (Value, error) {
	// 1. 匹配最佳签名
	signature, err := ft.MatchSignature(args, namedArgs)
	if err != nil {
		return nil, fmt.Errorf("function call failed: %v", err)
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
		// eval 求值
		// return ft.executeFunction(signature.body, allArgs)
	}

	return nil, fmt.Errorf("function has no implementation")
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

type OperatorOverload struct {
	operator   Operator
	leftType   Type
	rightType  Type
	returnType Type
	function   *FunctionType
}
