// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package object

import (
	"fmt"

	"slices"

	"github.com/hulo-lang/hulo/internal/core"
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

type Object interface {
	Type
	NumMethod() int
	Method(i int) Method
	MethodByName(name string) Method
	NumField() int
	Field(i int) Type
	FieldByName(name string) Type
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
	Call(args []Value, namedArgs map[string]Value, evaluator ...core.FunctionEvaluator) (Value, error)
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

	// 字段的默认值
	defaultValue Value
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

// SetDefaultValue 设置字段的默认值
func (f *Field) SetDefaultValue(value Value) {
	f.defaultValue = value
}

// GetDefaultValue 获取字段的默认值
func (f *Field) GetDefaultValue() Value {
	return f.defaultValue
}

// HasDefaultValue 检查字段是否有默认值
func (f *Field) HasDefaultValue() bool {
	return f.defaultValue != nil
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

func (o *ObjectType) AddField(name string, field *Field) {
	o.fields[name] = field
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

//go:generate stringer -type=Operator -linecomment
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

type TypeParameter struct {
	name        string
	constraints map[string]Constraint // 类型参数约束
}

func NewTypeParameter(name string, constraints map[string]Constraint) *TypeParameter {
	return &TypeParameter{
		name:        name,
		constraints: constraints,
	}
}

func (tp *TypeParameter) Name() string {
	return tp.name
}

func (tp *TypeParameter) Kind() ObjKind {
	return O_TYPE_PARAM
}

func (tp *TypeParameter) Text() string {
	return tp.name
}

func (tp *TypeParameter) AssignableTo(u Type) bool {
	return false // 类型参数不能直接赋值
}

func (tp *TypeParameter) ConvertibleTo(u Type) bool {
	return false
}

func (tp *TypeParameter) Implements(u Type) bool {
	return false
}

// SatisfiesConstraint 检查类型参数是否满足约束
func (tp *TypeParameter) SatisfiesConstraint(constraint Constraint) bool {
	// 类型参数本身不检查约束，在实例化时检查
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
	// 首先检查实例字段
	if value, exists := ci.fields[name]; exists {
		return value
	}

	// 检查静态字段
	if staticValue := ci.clsType.GetStaticValue(name); staticValue != NULL {
		return staticValue
	}

	// 检查父类的字段（递归）
	if ci.clsType.GetParent() != nil {
		// 这里需要创建一个父类的实例来访问字段
		// 简化处理，直接返回NULL
		return NULL
	}

	return NULL
}

func (ci *ClassInstance) SetField(name string, value Value) {
	// 检查是否为静态字段
	if ci.clsType.GetStaticField(name) != nil {
		ci.clsType.SetStaticValue(name, value)
		return
	}

	// 设置实例字段
	ci.fields[name] = value
}

// GetMethod 获取方法，支持继承和方法重写
func (ci *ClassInstance) GetMethod(name string) Method {
	return ci.clsType.MethodByName(name)
}

// CallMethod 调用方法
func (ci *ClassInstance) CallMethod(name string, args []Value, namedArgs map[string]Value) (Value, error) {
	method := ci.GetMethod(name)
	if method == nil {
		return NULL, fmt.Errorf("method %s not found", name)
	}

	// 将实例作为第一个参数传递
	allArgs := append([]Value{ci}, args...)
	return method.Call(allArgs, namedArgs)
}

type ClassConstructor struct {
	name      string
	astDecl   *ast.ConstructorDecl
	clsType   *ClassType
	signature []Type
	isNamed   bool
}

func (cc *ClassConstructor) Name() string {
	return cc.name
}

func (cc *ClassConstructor) IsNamed() bool {
	return cc.isNamed
}

func (cc *ClassConstructor) Signature() []Type {
	return cc.signature
}

func (cc *ClassConstructor) Call(args ...Value) Value {
	inst := &ClassInstance{
		clsType: cc.clsType,
		fields:  make(map[string]Value),
	}

	// 执行字段初始化
	if cc.astDecl != nil {
		for _, initExpr := range cc.astDecl.InitFields {
			// TODO: 执行字段初始化表达式
			// 这里需要实现表达式求值
			_ = initExpr
		}

		if cc.astDecl.Body != nil {
			// TODO: 执行构造函数体
			// 这里需要实现语句块执行
		}
	}

	return inst
}

func (cc *ClassConstructor) Match(args []Value) bool {
	if len(args) != len(cc.signature) {
		return false
	}

	for i, arg := range args {
		if !arg.Type().AssignableTo(cc.signature[i]) {
			return false
		}
	}

	return true
}

func (cc *ClassConstructor) NumIn() int {
	return len(cc.signature)
}

func (cc *ClassConstructor) In(i int) Type {
	if i < 0 || i >= len(cc.signature) {
		return nil
	}
	return cc.signature[i]
}

// NewClassConstructor 创建新的构造函数
func NewClassConstructor(name string, clsType *ClassType, signature []Type, isNamed bool) *ClassConstructor {
	return &ClassConstructor{
		name:      name,
		clsType:   clsType,
		signature: signature,
		isNamed:   isNamed,
	}
}

// NewDefaultConstructor 创建默认构造函数
func NewDefaultConstructor(clsType *ClassType) *ClassConstructor {
	return &ClassConstructor{
		name:      "default",
		clsType:   clsType,
		signature: []Type{},
		isNamed:   false,
	}
}

func ConvertClassDecl(decl *ast.ClassDecl) *ClassType {
	// 检查是否为泛型类
	var isGeneric bool
	var typeParams []string
	var constraints map[string]Constraint

	if decl.TypeParams != nil {
		isGeneric = true
		typeParams = make([]string, len(decl.TypeParams))
		constraints = make(map[string]Constraint)

		for i, param := range decl.TypeParams {
			if typeParam, ok := param.(*ast.TypeParameter); ok {
				if ident, ok := typeParam.Name.(*ast.Ident); ok {
					typeParams[i] = ident.Name
					if typeParam.Constraints != nil {
						// 转换约束类型为Constraint接口
						var constraintList []Constraint
						for _, constraintExpr := range typeParam.Constraints {
							constraintType := convertTypeExpr(constraintExpr)
							// 将Type转换为Constraint
							constraint := NewInterfaceConstraint(constraintType)
							constraintList = append(constraintList, constraint)
						}

						// 如果有多个约束，创建复合约束
						if len(constraintList) == 1 {
							constraints[ident.Name] = constraintList[0]
						} else if len(constraintList) > 1 {
							constraints[ident.Name] = NewCompositeConstraint(constraintList...)
						}
					}
				}
			}
		}
	}

	var clsType *ClassType
	if isGeneric {
		clsType = NewGenericClass(decl.Name.Name, typeParams, constraints)
	} else {
		clsType = NewClassType(decl.Name.Name)
	}

	clsType.astDecl = decl

	// 处理继承
	if decl.Parent != nil {
		// TODO: 解析父类并设置继承关系
		// parentType := resolveType(decl.Parent.Name)
		// clsType.SetParent(parentType)
	}

	// 处理字段
	if decl.Fields != nil {
		for _, field := range decl.Fields.List {
			fieldType := convertTypeExpr(field.Type)
			fieldMods := convertFieldModifiers(field.Modifiers)

			// 检查是否为静态字段
			isStatic := false
			for _, mod := range fieldMods {
				if mod == FieldModifierStatic {
					isStatic = true
					break
				}
			}

			if isStatic {
				clsType.AddStaticField(field.Name.Name, &Field{
					ft:   fieldType,
					mods: fieldMods,
				})
			} else {
				clsType.AddField(field.Name.Name, &Field{
					ft:   fieldType,
					mods: fieldMods,
				})
			}
		}
	}

	// 处理构造函数
	if decl.Ctors != nil {
		for _, ctor := range decl.Ctors {
			signature := convertConstructorSignature(ctor)
			isNamed := ctor.Name.Name != decl.Name.Name

			constructor := &ClassConstructor{
				name:      ctor.Name.Name,
				astDecl:   ctor,
				clsType:   clsType,
				signature: signature,
				isNamed:   isNamed,
			}

			clsType.AddConstructor(constructor)
		}
	}

	// 如果没有构造函数，添加默认构造函数
	if len(clsType.ctors) == 0 {
		defaultCtor := NewDefaultConstructor(clsType)
		clsType.AddConstructor(defaultCtor)
	}

	// 处理方法
	if decl.Methods != nil {
		for _, method := range decl.Methods {
			methodType := convertMethodDecl(method, clsType)

			// 检查是否为静态方法
			isStatic := false
			for _, mod := range method.Modifiers {
				if mod.Kind() == ast.ModKindStatic {
					isStatic = true
					break
				}
			}

			if isStatic {
				clsType.AddStaticMethod(method.Name.Name, methodType.(*FunctionType))
			} else {
				// 检查是否为重写方法
				isOverride := false
				for _, dec := range method.Decs {
					if dec.Name.Name == "override" {
						isOverride = true
						break
					}
				}

				if isOverride {
					clsType.OverrideMethod(method.Name.Name, methodType.(*FunctionType))
				} else {
					clsType.AddMethod(method.Name.Name, methodType)
				}
			}
		}
	}

	return clsType
}

// convertTypeExpr 转换类型表达式
func convertTypeExpr(expr ast.Expr) Type {
	// TODO: 实现类型表达式转换
	// 这里需要根据AST节点类型进行转换
	return anyType
}

// convertFieldModifiers 转换字段修饰符
func convertFieldModifiers(modifiers []ast.Modifier) []FieldModifier {
	var mods []FieldModifier
	for _, mod := range modifiers {
		switch mod.Kind() {
		case ast.ModKindRequired:
			mods = append(mods, FieldModifierRequired)
		case ast.ModKindReadonly:
			mods = append(mods, FieldModifierReadonly)
		case ast.ModKindStatic:
			mods = append(mods, FieldModifierStatic)
		case ast.ModKindPub:
			mods = append(mods, FieldModifierPub)
		}
	}
	return mods
}

// convertConstructorSignature 转换构造函数签名
func convertConstructorSignature(ctor *ast.ConstructorDecl) []Type {
	var signature []Type
	if ctor.Recv != nil {
		for _, param := range ctor.Recv {
			if param, ok := param.(*ast.Parameter); ok {
				paramType := convertTypeExpr(param.Type)
				signature = append(signature, paramType)
			}
		}
	}
	return signature
}

// convertMethodDecl 转换方法声明
func convertMethodDecl(method *ast.FuncDecl, clsType *ClassType) Method {
	// TODO: 实现方法声明转换
	// 这里需要创建FunctionType
	return &FunctionType{
		name:       method.Name.Name,
		signatures: []*FunctionSignature{},
	}
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
