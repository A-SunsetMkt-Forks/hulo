// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package object

import (
	"fmt"
	"strings"

	"slices"

	"github.com/hulo-lang/hulo/syntax/hulo/ast"
)

type Type interface {
	Name() string

	Text() string

	Kind() ObjKind

	// New(values ...Value) Value

	// Implements reports whether the type implements the interface type u.
	Implements(u Type) bool

	AssignableTo(u Type) bool

	ConvertibleTo(u Type) bool

	// NumMethod() int

	// Method(i int) Method

	// MethodByName(name string) Method

	// NumField() int

	// Field(i int) Type

	// FieldByName(name string) Type
}

// TypeRelation represents the relation between two types.
type TypeRelation interface {
	// IsAssignableTo reports whether a type can be assigned to the type u.
	IsAssignableTo(u Type) bool
	// IsSubtypeOf reports whether a type is a subtype of the type u.
	IsSubtypeOf(u Type) bool
	// IsSupertypeOf reports whether a type is a supertype of the type u.
	IsSupertypeOf(u Type) bool
	// IsCompatibleWith reports whether a type is compatible with the type u.
	IsCompatibleWith(u Type) bool
}

// TypeOperator represents the operator between two types.
type TypeOperator interface {
	// Unify returns the most specific type that is a supertype of both types.
	Unify(other Type) (Type, error)
	// Intersect returns the most specific type that is a subtype of both types.
	Intersect(other Type) (Type, error)
	// Union returns the most specific type that is a supertype of both types.
	Union(other Type) Type
}

type Value interface {
	Type() Type // 获取值的类型

	Text() string // 值的文本表示

	Interface() any // 获取底层 Go 值
}

type SourceLocation interface {
	File() string
	Line() int
	Col() int
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
}

// Type variance
type Variance int

const (
	Invariant Variance = iota
	Covariant
	Contravariant
	Bivariant
)

//go:generate stringer -type=FieldModifier -linecomment
type FieldModifier int

const (
	FieldModifierNone     FieldModifier = iota // none
	FieldModifierRequired                      // required
	FieldModifierOptional                      // ?
	FieldModifierReadonly                      // readonly
	FieldModifierStatic                        // static
	FieldModifierPub                           // pub
)

type Field struct {
	ft   Type
	mods []FieldModifier

	getter Function
	setter Function

	decorators []*ast.Decorator
}

func (f *Field) AddModifier(mod FieldModifier) *Field {
	f.mods = append(f.mods, mod)
	return f
}

func (f *Field) HasModifier(mod FieldModifier) bool {
	return slices.Contains(f.mods, mod)
}

func (f *Field) IsRequired() bool {
	return f.HasModifier(FieldModifierRequired)
}

type ObjectType struct {
	name      string
	kind      ObjKind
	methods   map[string]Method
	fields    map[string]*Field
	operators map[Operator]*FunctionType
	// e.g. str("Hello World")
	callSignature *Function
	// e.g. new str("Hello World")
	constructSignature *Function
}

func NewObjectType(name string, kind ObjKind) *ObjectType {
	return &ObjectType{
		name:      name,
		kind:      kind,
		methods:   make(map[string]Method),
		fields:    make(map[string]*Field),
		operators: make(map[Operator]*FunctionType),
	}
}

func (o *ObjectType) AddMethod(name string, signature Method) {
	o.methods[name] = signature
}

func (t *ObjectType) Name() string  { return t.name }
func (t *ObjectType) Kind() ObjKind { return t.kind }

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
	fields := make([]*Field, 0, len(t.fields))
	for _, field := range t.fields {
		fields = append(fields, field)
	}
	return fields[i].ft
}

func (t *ObjectType) FieldByName(name string) Type {
	return t.fields[name].ft
}

func (t *ObjectType) Implements(u Type) bool {
	return false
}

func (t *ObjectType) AssignableTo(u Type) bool {
	return t.Kind() == u.Kind()
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

type BuiltinFunction func(args ...Value) Value

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
	template    Type            // 模板类型
	instances   map[string]Type // 已实例化的类型缓存
}

// AssignableTo implements Type.
func (gt *GenericType) AssignableTo(u Type) bool {
	panic("unimplemented")
}

// ConvertibleTo implements Type.
func (gt *GenericType) ConvertibleTo(u Type) bool {
	panic("unimplemented")
}

// Implements implements Type.
func (gt *GenericType) Implements(u Type) bool {
	panic("unimplemented")
}

// Kind implements Type.
func (gt *GenericType) Kind() ObjKind {
	panic("unimplemented")
}

// Name implements Type.
func (gt *GenericType) Name() string {
	panic("unimplemented")
}

// Text implements Type.
func (gt *GenericType) Text() string {
	panic("unimplemented")
}

type TypeParameter struct {
	name       string
	constraint Type
	variance   Variance
}

func (gt *GenericType) Instantiate(typeArgs []Type) (Type, error) {
	// 1. 检查类型参数数量
	if len(typeArgs) != len(gt.typeParams) {
		return nil, fmt.Errorf("type argument count mismatch: expected %d, got %d",
			len(gt.typeParams), len(typeArgs))
	}

	// 2. 检查类型约束
	for i, param := range gt.typeParams {
		if constraint, exists := gt.constraints[param]; exists {
			if !typeArgs[i].AssignableTo(constraint) {
				return nil, fmt.Errorf("type %s does not satisfy constraint for %s",
					typeArgs[i].Name(), param)
			}
		}
	}

	// 3. 生成类型键
	typeKey := gt.generateTypeKey(typeArgs)

	// 4. 检查缓存
	if instance, exists := gt.instances[typeKey]; exists {
		return instance, nil
	}

	// 5. 创建实例化类型
	instance := gt.substituteType(gt.template, gt.createTypeMap(typeArgs))

	// 6. 缓存实例
	gt.instances[typeKey] = instance

	return instance, nil
}

func (gt *GenericType) createTypeMap(typeArgs []Type) map[string]Type {
	typeMap := make(map[string]Type)
	for i, param := range gt.typeParams {
		typeMap[param] = typeArgs[i]
	}
	return typeMap
}

func (gt *GenericType) generateTypeKey(typeArgs []Type) string {
	keys := make([]string, len(typeArgs))
	for i, t := range typeArgs {
		keys[i] = t.Name()
	}
	return strings.Join(keys, ",")
}

func (gt *GenericType) substituteType(t Type, typeMap map[string]Type) Type {
	switch v := t.(type) {
	case *GenericType:
		// 递归替换泛型类型
		return gt.substituteGenericType(v, typeMap)

	case *ClassType:
		// 替换类类型中的泛型参数
		return gt.substituteClassType(v, typeMap)

	case *FunctionType:
		// 替换函数类型中的泛型参数
		return gt.substituteFunctionType(v, typeMap)

	case *UnionType:
		// 替换联合类型中的泛型参数
		return gt.substituteUnionType(v, typeMap)

	default:
		return t
	}
}

func (gt *GenericType) substituteGenericType(genType *GenericType, typeMap map[string]Type) Type {
	// 检查是否是类型参数
	if replacement, exists := typeMap[genType.name]; exists {
		return replacement
	}

	// 递归替换内部类型
	newTemplate := gt.substituteType(genType.template, typeMap)
	return &GenericType{
		name:       genType.name,
		typeParams: genType.typeParams,
		template:   newTemplate,
		instances:  make(map[string]Type),
	}
}

func (gt *GenericType) substituteClassType(clsType *ClassType, typeMap map[string]Type) Type {
	newClass := &ClassType{
		ObjectType: clsType.ObjectType,
		// name:       clsType.name,
		// fields:     make(map[string]*Field),
		// methods:    make(map[string]*FunctionType),
		// operators:  make(map[Operator]*FunctionType),
		isGeneric:  false, // 实例化后不再是泛型
	}

	// 替换字段类型
	for name, field := range clsType.fields {
		newClass.fields[name] = &Field{
			ft:   gt.substituteType(field.ft, typeMap),
			mods: field.mods,
		}
	}

	// 替换方法类型
	for name, method := range clsType.methods {
		newClass.methods[name] = gt.substituteFunctionType(method.(*FunctionType), typeMap)
	}

	// 替换运算符类型
	for op, fn := range clsType.operators {
		newClass.operators[op] = gt.substituteFunctionType(fn, typeMap)
	}

	return newClass
}

func (gt *GenericType) substituteFunctionType(fnType *FunctionType, typeMap map[string]Type) *FunctionType {
	newFn := &FunctionType{
		name:       fnType.name,
		signatures: make([]*FunctionSignature, 0),
		isGeneric:  false,
	}

	// 替换函数签名
	for _, sig := range fnType.signatures {
		newSig := &FunctionSignature{
			positionalParams: make([]*Parameter, 0),
			namedParams:      make([]*Parameter, 0),
			returnType:       gt.substituteType(sig.returnType, typeMap),
			isOperator:       sig.isOperator,
			operator:         sig.operator,
			body:             sig.body,
			builtin:          sig.builtin,
		}

		// 替换位置参数
		for _, param := range sig.positionalParams {
			newSig.positionalParams = append(newSig.positionalParams, &Parameter{
				name:         param.name,
				typ:          gt.substituteType(param.typ, typeMap),
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
				typ:          gt.substituteType(param.typ, typeMap),
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
				typ:          gt.substituteType(sig.variadicParam.typ, typeMap),
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

func (gt *GenericType) substituteUnionType(unionType *UnionType, typeMap map[string]Type) Type {
	newTypes := make([]Type, len(unionType.types))
	for i, t := range unionType.types {
		newTypes[i] = gt.substituteType(t, typeMap)
	}
	return NewUnionType(newTypes...)
}

type InstantiatedType struct {
	generic  *GenericType
	typeArgs []Type
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
		ObjectType: NewObjectType(decl.Name.Name, O_CLASS),
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

type PrimitiveType struct {
	name string
	kind ObjKind
}

var _ Type = (*PrimitiveType)(nil)

func NewPrimitiveType(name string, kind ObjKind) *PrimitiveType {
	return &PrimitiveType{
		name: name,
		kind: kind,
	}
}

func (pt *PrimitiveType) Name() string {
	return pt.name
}

func (pt *PrimitiveType) Kind() ObjKind {
	return pt.kind
}

func (pt *PrimitiveType) Text() string {
	return pt.name
}

func (pt *PrimitiveType) Implements(u Type) bool { return false }

func (pt *PrimitiveType) AssignableTo(u Type) bool { return pt.Kind() == u.Kind() }

func (pt *PrimitiveType) ConvertibleTo(u Type) bool { return pt.kind.Equal(u.Kind()) }
