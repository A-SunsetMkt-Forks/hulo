// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package object

import (
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/hulo-lang/hulo/syntax/hulo/ast"
)

type ObjKind int

func (o ObjKind) Equal(x ObjKind) bool {
	return o == x
}

const (
	O_NUM ObjKind = iota
	O_STR
	O_BOOL
	O_FUNC
	O_NULL
	O_BUILTIN
	O_LITERAL
	O_ARR
	O_MAP
	O_RET
	O_QUOTE
	O_TRAIT
	O_ENUM
	O_CLASS
	O_OBJ
)

type Type interface {
	Name() string

	Text() string

	New(values ...Value) Value

	// Implements reports whether the type implements the interface type u.
	Implements(u Type) bool

	AssignableTo(u Type) bool

	ConvertibleTo(u Type) bool

	Kind() ObjKind

	NumMethod() int

	Method(i int) Method

	MethodByName(name string) Method

	NumField() int

	Field(i int) Type

	FieldByName(name string) Type
}

type Value interface {
	Type() Type

	Text() string

	Interface() any
}

type Function interface {
	Signature() []Type
	Call(args ...Value) Value
	Match(args []Value) bool

	// NumIn returns a function type's input parameter count.
	// It panics if the type's Kind is not Func.
	NumIn() int

	// In returns the type of a function type's i'th input parameter.
	// It panics if the type's Kind is not Func.
	// It panics if i is not in the range [0, NumIn()).
	In(i int) Type
}

type Method interface {
	Type
	Function
}

type ObjectType struct {
	name    string
	kind    ObjKind
	pkgPath string
	methods map[string]Method
	fields  map[string]Type
}

func NewObjectType(name string, kind ObjKind, pkgPath string) *ObjectType {
	return &ObjectType{
		name:    name,
		kind:    kind,
		pkgPath: pkgPath,
		methods: make(map[string]Method),
		fields:  make(map[string]Type),
	}
}

func (t *ObjectType) Name() string    { return t.name }
func (t *ObjectType) Kind() ObjKind   { return t.kind }
func (t *ObjectType) PkgPath() string { return t.pkgPath }

func (t *ObjectType) NumMethod() int {
	return len(t.methods)
}

func (t *ObjectType) Method(i int) Method {
	if i < 0 || i >= len(t.methods) {
		return nil
	}
	methods := make([]Method, 0, len(t.methods))
	for _, method := range t.methods {
		methods = append(methods, method)
	}
	return methods[i]
}

func (t *ObjectType) MethodByName(name string) Method {
	return t.methods[name]
}

func (t *ObjectType) NumField() int {
	return len(t.fields)
}

func (t *ObjectType) Field(i int) Type {
	if i < 0 || i >= len(t.fields) {
		return nil
	}
	fields := make([]Type, 0, len(t.fields))
	for _, method := range t.fields {
		fields = append(fields, method)
	}
	return fields[i]
}

func (t *ObjectType) FieldByName(name string) Type {
	return t.fields[name]
}

func (t *ObjectType) Implements(u Type) bool {
	return false
}

func (t *ObjectType) AssignableTo(u Type) bool {
	return false
}

func (t *ObjectType) ConvertibleTo(u Type) bool {
	return false
}

func (t *ObjectType) New(values ...Value) Value {
	return nil
}

func (t *ObjectType) Text() string {
	return t.name
}

type NullValue struct{}

func (n *NullValue) Text() string {
	return "null"
}

func (n *NullValue) Interface() any {
	return nil
}

func (n *NullValue) Type() Type {
	return nil
}

var NULL = &NullValue{}

type NumberValue struct {
	// TODO 根据精度选择存储模型
	Value *big.Float
}

func (n *NumberValue) Text() string {
	return n.Value.String()
}

func (n *NumberValue) Interface() any {
	return n.Value
}

func (n *NumberValue) Type() Type {
	return NumberType
}

var NumberType = NewObjectType("number", O_NUM, "std")

type String struct {
	Value string
}

func (s *String) Text() string {
	return s.Value
}

func (s *String) Interface() any {
	return s.Value
}

func (s *String) Type() Type {
	return StringType
}

var StringType = NewObjectType("string", O_STR, "std")

type Boolean struct {
	Value bool
}

func (b *Boolean) Text() string {
	return strconv.FormatBool(b.Value)
}

func (b *Boolean) Interface() any {
	return b.Value
}

func (b *Boolean) Type() Type {
	return BooleanType
}

var BooleanType = NewObjectType("bool", O_BOOL, "std")

var TRUE = &Boolean{Value: true}
var FALSE = &Boolean{Value: false}

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
	return ErrorType
}

var ErrorType = NewObjectType("error", O_OBJ, "std")

type OverloadedFunction struct {
	name  string
	funcs []Function
}

func (o *OverloadedFunction) Text() string {
	return "overloaded function"
}

func (o *OverloadedFunction) Interface() any {
	return o.funcs
}

func (o *OverloadedFunction) Type() Type {
	return FunctionType
}

var FunctionType = NewObjectType("function", O_FUNC, "std")

type BuiltinFunction func(args ...Value) Value

func (b *BuiltinFunction) Text() string {
	return "builtin function"
}

func (b *BuiltinFunction) Interface() any {
	return b
}

func (b *BuiltinFunction) Type() Type {
	return FunctionType
}

type UserFunction struct {
}

type OperatorFunction struct{}

type Operator int

const (
	OpIllegal      Operator = iota
	OpAdd                   // +
	OpSub                   // -
	OpMul                   // *
	OpDiv                   // /
	OpMod                   // %
	OpConcat                // ..
	OpAnd                   // &&
	OpOr                    // ||
	OpEqual                 // ==
	OpNot                   // !
	OpLess                  // <
	OpLessEqual             // <=
	OpGreater               // >
	OpGreaterEqual          // >=
	OpAssign                // =
	OpIndex                 // []
	OpCall                  // ()

	Opa // a""
	Opb // b""
	Opc // c""
	Opd // d""
	Ope // e""
	Opf // f""
	Opg // g""
	Oph // h""
	Opi // i""
	Opj // j""
	Opk // k""
	Opl // l""
	Opm // m""
	Opn // n""
	Opo // o""
	Opp // p""
	Opq // q""
	Opr // r""
	Ops // s""
	Ott // t""
	Opu // u""
	Opv // v""
	Opw // w""
	Opz // z""

	OpA // A""
	OpB // B""
	OpC // C""
	OpD // D""
	OpE // E""
	OpF // F""
	OpG // G""
	OpH // H""
	OpI // I""
	OpJ // J""
	OpK // K""
	OpL // L""
	OpM // M""
	OpN // N""
	OpO // O""
	OpP // P""
	OpQ // Q""
	OpR // R""
	OpS // S""
	OpT // T""
	OpU // U""
	OpV // V""
	OpW // W""
	OpX // X""
	OpY // Y""
	OpZ // Z""
)

type GenericType struct {
	name        string
	typeParams  []string // e.g. T, U
	constraints map[string]Type
	template    ast.Node
}

func (gt *GenericType) Instantiate(typeArgs []Type) (*InstantiatedType, error) {
	if len(typeArgs) != len(gt.typeParams) {
		return nil, fmt.Errorf("type argument count mismatch")
	}

	typeMap := make(map[string]Type)
	for i, param := range gt.typeParams {
		typeMap[param] = typeArgs[i]
	}

	// 替换模板中的类型参数

	return &InstantiatedType{
		generic:  gt,
		typeArgs: typeArgs,
	}, nil
}

type InstantiatedType struct {
	generic  *GenericType
	typeArgs []Type
}

type GenericFunction struct {
	name       string
	typeParams []string
	signatures map[string]Function
}

type Constraint interface {
	SatisfiedBy(Type) bool
}

type InterfaceConstraint struct {
	interfaceType Type
}

func (ic *InterfaceConstraint) SatisfiedBy(t Type) bool {
	return t.Implements(ic.interfaceType)
}

type OperatorConstraint struct {
	op        Operator
	rightType Type
}

func (oc *OperatorConstraint) SatisfiedBy(t Type) bool {
	return true
}

type CompositeConstraint struct {
	constraints []Constraint
}

func (cc *CompositeConstraint) SatisfiedBy(t Type) bool {
	for _, constraint := range cc.constraints {
		if !constraint.SatisfiedBy(t) {
			return false
		}
	}
	return true
}

type ClassType struct {
	*ObjectType
	astDecl     *ast.ClassDecl
	ctors       []Constructor
	isGeneric   bool
	typeParams  []string              // e.g. T, U
	constraints map[string]Type       // e.g. T extends num | str
	instances   map[string]*ClassType // 已实例化的类型缓存
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
		ObjectType:  NewObjectType(ct.Name()+"<"+typeKey+">", O_CLASS, ct.PkgPath()),
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

type Constructor interface {
	Function
	Name() string
	IsNamed() bool // 是否为命名构造函数
}

type ClassInstance struct {
	clsType *ClassType
	fields  map[string]Value
}

func (ci *ClassInstance) Text() string {
	return fmt.Sprintf("%s{...}", ci.clsType.Name())
}

func (ci *ClassInstance) Interface() any {
	return ci
}

func (ci *ClassInstance) Type() Type {
	return ci.clsType
}

func (ci *ClassInstance) GetField(name string) Value {
	return ci.fields[name]
}

func (ci *ClassInstance) SetField(name string, value Value) {
	ci.fields[name] = value
}

type ClassConstructor struct {
	name      string
	astDecl   *ast.ConstructorDecl
	clsType   *ClassType
	signature []Type
}

func (cc *ClassConstructor) Call(args ...Value) Value {
	inst := &ClassInstance{
		clsType: cc.clsType,
		fields:  make(map[string]Value),
	}

	for _, initExpr := range cc.astDecl.InitFields {
		// 核心是executeExpr求值
		fmt.Println(initExpr)
	}

	if cc.astDecl.Body != nil {
		// executeBlockStmt(cc.astDecl.Body, inst)
	}

	return inst
}

func ConvertClassDecl(decl *ast.ClassDecl) *ClassType {
	clsType := &ClassType{
		ObjectType: NewObjectType(decl.Name.Name, O_CLASS, "std"),
		astDecl:    decl,
		ctors:      make([]Constructor, 0),
	}

	// // 1. 处理字段
	// for _, field := range classDecl.Fields.List {
	//     fieldType := convertTypeExpr(field.Type) // 需要实现类型转换
	//     classType.fields[field.Name.Name] = fieldType
	// }

	// // 2. 处理构造函数
	// for _, ctor := range classDecl.Ctors {
	//     constructor := &ClassConstructor{
	//         name:      ctor.Name.Name,
	//         astDecl:   ctor,
	//         classType: classType,
	//     }
	//     classType.constructors = append(classType.constructors, constructor)
	// }

	// // 3. 处理方法
	// for _, method := range classDecl.Methods {
	//     methodType := convertFuncDeclToMethod(method, classType) // 需要实现
	//     classType.methods[method.Name.Name] = methodType
	// }

	return clsType
}
