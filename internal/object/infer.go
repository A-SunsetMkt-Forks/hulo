package object

import (
	"fmt"
	"maps"
	"strings"
	"sync"
)

// TypeInference 类型推断接口
type TypeInference interface {
	// InferTypeArgs 从具体类型推断泛型类型参数
	InferTypeArgs(concreteType Type, genericType Type) ([]Type, error)

	// InferFromUsage 从使用场景推断类型
	InferFromUsage(usage *TypeUsage) (map[string]Type, error)

	// InferFromConstraints 从约束推断类型
	InferFromConstraints(constraints map[string]Constraint) (map[string]Type, error)

	// InferFromAssignment 从赋值推断类型
	InferFromAssignment(targetType Type, sourceType Type) (map[string]Type, error)

	// InferFromContext 从上下文推断类型
	InferFromContext(context *TypeContext) (map[string]Type, error)

	// InferBidirectional 双向类型推断
	InferBidirectional(target Type, source Type) (map[string]Type, error)

	// InferFromPattern 从模式匹配推断类型
	InferFromPattern(pattern PatternMatch) (map[string]Type, error)

	// SolveConstraints 求解约束系统
	SolveConstraints(constraints map[string]Constraint) (map[string]Type, error)

	// InferFromOverload 从重载函数推断类型
	InferFromOverload(overloads []*FunctionType, args []Value) (map[string]Type, error)

	// UnifyTypes 类型统一算法
	UnifyTypes(t1, t2 Type) (map[string]Type, error)

	// PropagateConstraints 约束传播
	PropagateConstraints(constraints map[string]Constraint, unifier map[string]Type) (map[string]Constraint, error)

	// InferRecursive 递归类型推断
	InferRecursive(types []Type) (map[string]Type, error)

	// InferWithContext 带上下文的类型推断
	InferWithContext(target Type, context *TypeContext) (map[string]Type, error)

	// SubstituteTypes 类型替换
	SubstituteTypes(t Type, substitutions map[string]Type) Type

	// CheckConstraintSatisfaction 检查约束满足性
	CheckConstraintSatisfaction(t Type, constraint Constraint) bool

	// InferFromLiteral 从字面量推断类型
	InferFromLiteral(literal Value) Type

	// InferFromExpression 从表达式推断类型（增强版）
	InferFromExpression(expr *ExpressionContext) map[string]Type

	// InferFromBinaryOp 从二元操作推断类型
	InferFromBinaryOp(op Operator, left, right Type) (map[string]Type, error)

	// InferFromUnaryOp 从一元操作推断类型
	InferFromUnaryOp(op Operator, operand Type) (map[string]Type, error)

	// InferFromIndexOp 从索引操作推断类型
	InferFromIndexOp(container, index Type) (map[string]Type, error)

	// InferFromCall 从函数调用推断类型
	InferFromCall(funcType Type, args []Type) (map[string]Type, error)
}

// TypeUsage 表示类型使用场景
type TypeUsage struct {
	Context     string                // 使用上下文
	Arguments   []Type                // 参数类型
	ReturnType  Type                  // 返回类型
	Constraints map[string]Constraint // 约束信息
	Hints       map[string]Type       // 类型提示
}

// DefaultTypeInference 默认类型推断实现
type DefaultTypeInference struct {
	cache map[string][]Type // 推断结果缓存
	mutex sync.RWMutex
}

// NewTypeInference 创建类型推断实例
func NewTypeInference() TypeInference {
	return &DefaultTypeInference{
		cache: make(map[string][]Type),
	}
}

// InferTypeArgs 从具体类型推断泛型类型参数
func (ti *DefaultTypeInference) InferTypeArgs(concreteType Type, genericType Type) ([]Type, error) {
	// 生成缓存键
	cacheKey := fmt.Sprintf("%s->%s", concreteType.Name(), genericType.Name())

	// 检查缓存
	ti.mutex.RLock()
	if cached, exists := ti.cache[cacheKey]; exists {
		ti.mutex.RUnlock()
		return cached, nil
	}
	ti.mutex.RUnlock()

	// 执行推断
	typeArgs, err := ti.performInference(concreteType, genericType)
	if err != nil {
		return nil, err
	}

	// 缓存结果
	ti.mutex.Lock()
	ti.cache[cacheKey] = typeArgs
	ti.mutex.Unlock()

	return typeArgs, nil
}

// performInference 执行实际的类型推断
func (ti *DefaultTypeInference) performInference(concreteType Type, genericType Type) ([]Type, error) {
	switch gt := genericType.(type) {
	case *ClassType:
		return ti.inferClassTypeArgs(concreteType, gt)
	case *FunctionType:
		return ti.inferFunctionTypeArgs(concreteType, gt)
	case *ArrayType:
		return ti.inferArrayTypeArgs(concreteType, gt)
	case *MapType:
		return ti.inferMapTypeArgs(concreteType, gt)
	default:
		return nil, fmt.Errorf("unsupported generic type for inference: %T", genericType)
	}
}

// inferClassTypeArgs 推断类类型的类型参数
func (ti *DefaultTypeInference) inferClassTypeArgs(concreteType Type, genericClass *ClassType) ([]Type, error) {
	if !genericClass.IsGeneric() {
		return nil, fmt.Errorf("class %s is not generic", genericClass.Name())
	}

	// 检查是否为泛型类的实例
	if classInstance, ok := concreteType.(*ClassType); ok {
		// 尝试从类名推断
		if inferred := ti.inferFromClassName(classInstance.Name(), genericClass); len(inferred) > 0 {
			return inferred, nil
		}
	}

	// 从字段类型推断
	if inferred := ti.inferFromFields(concreteType, genericClass); len(inferred) > 0 {
		return inferred, nil
	}

	// 从方法签名推断
	if inferred := ti.inferFromMethods(concreteType, genericClass); len(inferred) > 0 {
		return inferred, nil
	}

	return nil, fmt.Errorf("cannot infer type arguments for %s from %s", genericClass.Name(), concreteType.Name())
}

// inferFunctionTypeArgs 推断函数类型的类型参数
func (ti *DefaultTypeInference) inferFunctionTypeArgs(concreteType Type, genericFunc *FunctionType) ([]Type, error) {
	if !genericFunc.IsGeneric() {
		return nil, fmt.Errorf("function %s is not generic", genericFunc.Name())
	}

	// 从函数签名推断
	if funcType, ok := concreteType.(*FunctionType); ok {
		return ti.inferFromFunctionSignature(funcType, genericFunc)
	}

	return nil, fmt.Errorf("cannot infer function type arguments from %s", concreteType.Name())
}

// inferArrayTypeArgs 推断数组类型的类型参数
func (ti *DefaultTypeInference) inferArrayTypeArgs(concreteType Type, genericArray *ArrayType) ([]Type, error) {
	if arrayType, ok := concreteType.(*ArrayType); ok {
		// 推断元素类型
		elementType := arrayType.ElementType()
		return []Type{elementType}, nil
	}
	return nil, fmt.Errorf("cannot infer array type arguments from %s", concreteType.Name())
}

// inferMapTypeArgs 推断映射类型的类型参数
func (ti *DefaultTypeInference) inferMapTypeArgs(concreteType Type, genericMap *MapType) ([]Type, error) {
	if mapType, ok := concreteType.(*MapType); ok {
		// 推断键和值类型
		keyType := mapType.KeyType()
		valueType := mapType.ValueType()
		return []Type{keyType, valueType}, nil
	}
	return nil, fmt.Errorf("cannot infer map type arguments from %s", concreteType.Name())
}

// inferFromClassName 从类名推断类型参数
func (ti *DefaultTypeInference) inferFromClassName(className string, genericClass *ClassType) []Type {
	// 解析类名中的类型参数，如 "List<String>" -> ["String"]
	// 这里简化实现，实际需要更复杂的解析
	if strings.Contains(className, "<") && strings.Contains(className, ">") {
		start := strings.Index(className, "<")
		end := strings.LastIndex(className, ">")
		if start < end {
			argsStr := className[start+1 : end]
			args := strings.Split(argsStr, ",")
			typeArgs := make([]Type, 0, len(args))

			for _, arg := range args {
				arg = strings.TrimSpace(arg)
				// 这里需要从类型注册表中查找类型
				// 简化实现，返回 anyType
				typeArgs = append(typeArgs, anyType)
			}

			return typeArgs
		}
	}
	return nil
}

// inferFromFields 从字段类型推断类型参数
func (ti *DefaultTypeInference) inferFromFields(concreteType Type, genericClass *ClassType) []Type {
	// 获取具体类型的字段
	var concreteFields map[string]*Field
	switch ct := concreteType.(type) {
	case *ClassType:
		concreteFields = ct.fields
	case *ObjectType:
		concreteFields = ct.fields
	default:
		return nil
	}

	// 获取泛型类的字段模板
	genericFields := genericClass.fields

	// 尝试匹配字段类型
	typeArgs := make([]Type, len(genericClass.TypeParameters()))
	typeMap := make(map[string]Type)

	for fieldName, genericField := range genericFields {
		if concreteField, exists := concreteFields[fieldName]; exists {
			// 推断字段类型中的类型参数
			if inferred := ti.inferTypeFromField(genericField.ft, concreteField.ft); len(inferred) > 0 {
				for param, typ := range inferred {
					typeMap[param] = typ
				}
			}
		}
	}

	// 将推断结果转换为类型参数列表
	for i, param := range genericClass.TypeParameters() {
		if typ, exists := typeMap[param]; exists {
			typeArgs[i] = typ
		} else {
			typeArgs[i] = anyType // 无法推断时使用 anyType
		}
	}

	return typeArgs
}

// inferFromMethods 从方法签名推断类型参数
func (ti *DefaultTypeInference) inferFromMethods(concreteType Type, genericClass *ClassType) []Type {
	// 获取具体类型的方法
	var concreteMethods map[string]Method
	switch ct := concreteType.(type) {
	case *ClassType:
		concreteMethods = ct.methods
	case *ObjectType:
		concreteMethods = ct.methods
	default:
		return nil
	}

	// 获取泛型类的方法模板
	genericMethods := genericClass.methods

	// 尝试匹配方法签名
	typeArgs := make([]Type, len(genericClass.TypeParameters()))
	typeMap := make(map[string]Type)

	for methodName, genericMethod := range genericMethods {
		if concreteMethod, exists := concreteMethods[methodName]; exists {
			// 推断方法签名中的类型参数
			if inferred := ti.inferTypeFromMethod(genericMethod, concreteMethod); len(inferred) > 0 {
				for param, typ := range inferred {
					typeMap[param] = typ
				}
			}
		}
	}

	// 将推断结果转换为类型参数列表
	for i, param := range genericClass.TypeParameters() {
		if typ, exists := typeMap[param]; exists {
			typeArgs[i] = typ
		} else {
			typeArgs[i] = anyType
		}
	}

	return typeArgs
}

// inferFromFunctionSignature 从函数签名推断类型参数
func (ti *DefaultTypeInference) inferFromFunctionSignature(concreteFunc *FunctionType, genericFunc *FunctionType) ([]Type, error) {
	if len(concreteFunc.signatures) == 0 || len(genericFunc.signatures) == 0 {
		return nil, fmt.Errorf("function has no signatures")
	}

	// 使用第一个签名进行推断
	concreteSig := concreteFunc.signatures[0]
	genericSig := genericFunc.signatures[0]

	typeArgs := make([]Type, len(genericFunc.TypeParameters()))
	typeMap := make(map[string]Type)

	// 从参数类型推断
	for i, genericParam := range genericSig.positionalParams {
		if i < len(concreteSig.positionalParams) {
			concreteParam := concreteSig.positionalParams[i]
			if inferred := ti.inferTypeFromField(genericParam.typ, concreteParam.typ); len(inferred) > 0 {
				for param, typ := range inferred {
					typeMap[param] = typ
				}
			}
		}
	}

	// 从返回类型推断
	if inferred := ti.inferTypeFromField(genericSig.returnType, concreteSig.returnType); len(inferred) > 0 {
		for param, typ := range inferred {
			typeMap[param] = typ
		}
	}

	// 将推断结果转换为类型参数列表
	for i, param := range genericFunc.TypeParameters() {
		if typ, exists := typeMap[param]; exists {
			typeArgs[i] = typ
		} else {
			typeArgs[i] = anyType
		}
	}

	return typeArgs, nil
}

// inferTypeFromField 从字段类型推断类型参数
func (ti *DefaultTypeInference) inferTypeFromField(genericType Type, concreteType Type) map[string]Type {
	result := make(map[string]Type)

	switch gt := genericType.(type) {
	case *TypeParameter:
		// 直接类型参数
		result[gt.Name()] = concreteType

	case *ArrayType:
		// 数组类型
		if concreteArray, ok := concreteType.(*ArrayType); ok {
			if inferred := ti.inferTypeFromField(gt.elementType, concreteArray.elementType); len(inferred) > 0 {
				for param, typ := range inferred {
					result[param] = typ
				}
			}
		}

	case *MapType:
		// 映射类型
		if concreteMap, ok := concreteType.(*MapType); ok {
			// 推断键类型
			if inferred := ti.inferTypeFromField(gt.keyType, concreteMap.keyType); len(inferred) > 0 {
				for param, typ := range inferred {
					result[param] = typ
				}
			}
			// 推断值类型
			if inferred := ti.inferTypeFromField(gt.valueType, concreteMap.valueType); len(inferred) > 0 {
				for param, typ := range inferred {
					result[param] = typ
				}
			}
		}

	case *ClassType:
		// 泛型类类型
		if concreteClass, ok := concreteType.(*ClassType); ok {
			if gt.IsGeneric() {
				// 递归推断泛型类的类型参数
				if typeArgs, err := ti.InferTypeArgs(concreteClass, gt); err == nil {
					for i, param := range gt.TypeParameters() {
						if i < len(typeArgs) {
							result[param] = typeArgs[i]
						}
					}
				}
			}
		}
	}

	return result
}

// inferTypeFromMethod 从方法签名推断类型参数
func (ti *DefaultTypeInference) inferTypeFromMethod(genericMethod Method, concreteMethod Method) map[string]Type {
	result := make(map[string]Type)

	// 这里需要更复杂的实现来比较方法签名
	// 简化实现：返回空结果
	return result
}

// InferFromUsage 从使用场景推断类型
func (ti *DefaultTypeInference) InferFromUsage(usage *TypeUsage) (map[string]Type, error) {
	result := make(map[string]Type)

	// 从参数推断
	for _, argType := range usage.Arguments {
		if inferred := ti.inferFromType(argType); len(inferred) > 0 {
			for param, typ := range inferred {
				result[param] = typ
			}
		}
	}

	// 从返回类型推断
	if usage.ReturnType != nil {
		if inferred := ti.inferFromType(usage.ReturnType); len(inferred) > 0 {
			for param, typ := range inferred {
				result[param] = typ
			}
		}
	}

	// 从约束推断
	if len(usage.Constraints) > 0 {
		if inferred, err := ti.InferFromConstraints(usage.Constraints); err == nil {
			for param, typ := range inferred {
				result[param] = typ
			}
		}
	}

	// 从提示推断
	for param, hint := range usage.Hints {
		result[param] = hint
	}

	return result, nil
}

// InferFromConstraints 从约束推断类型
func (ti *DefaultTypeInference) InferFromConstraints(constraints map[string]Constraint) (map[string]Type, error) {
	result := make(map[string]Type)

	for param, constraint := range constraints {
		if typ := ti.inferFromConstraint(constraint); typ != nil {
			result[param] = typ
		}
	}

	return result, nil
}

// InferFromAssignment 从赋值推断类型
func (ti *DefaultTypeInference) InferFromAssignment(targetType Type, sourceType Type) (map[string]Type, error) {
	result := make(map[string]Type)

	// 从目标类型推断
	if inferred := ti.inferFromType(targetType); len(inferred) > 0 {
		for param, typ := range inferred {
			result[param] = typ
		}
	}

	// 从源类型推断
	if inferred := ti.inferFromType(sourceType); len(inferred) > 0 {
		for param, typ := range inferred {
			result[param] = typ
		}
	}

	return result, nil
}

// inferFromType 从类型推断类型参数
func (ti *DefaultTypeInference) inferFromType(t Type) map[string]Type {
	result := make(map[string]Type)

	switch typ := t.(type) {
	case *TypeParameter:
		// 类型参数本身
		result[typ.Name()] = anyType

	case *ArrayType:
		// 数组元素类型
		if inferred := ti.inferFromType(typ.elementType); len(inferred) > 0 {
			for param, elementType := range inferred {
				result[param] = NewArrayType(elementType)
			}
		}

	case *MapType:
		// 映射的键值类型
		if keyInferred := ti.inferFromType(typ.keyType); len(keyInferred) > 0 {
			for param, keyType := range keyInferred {
				result[param] = keyType
			}
		}
		if valueInferred := ti.inferFromType(typ.valueType); len(valueInferred) > 0 {
			for param, valueType := range valueInferred {
				result[param] = valueType
			}
		}

	case *ClassType:
		// 泛型类
		if typ.IsGeneric() {
			for _, param := range typ.TypeParameters() {
				result[param] = anyType
			}
		}
	}

	return result
}

// inferFromConstraint 从约束推断类型
func (ti *DefaultTypeInference) inferFromConstraint(constraint Constraint) Type {
	switch c := constraint.(type) {
	case *InterfaceConstraint:
		// 从接口约束推断，返回接口类型本身
		return c.interfaceType

	case *OperatorConstraint:
		// 从运算符约束推断，返回支持该运算符的类型
		// 这里简化处理，实际需要更复杂的类型推断
		return anyType

	case *CompositeConstraint:
		// 从复合约束推断，返回满足所有约束的类型
		// 这里简化处理，返回第一个约束的类型
		if len(c.constraints) > 0 {
			return ti.inferFromConstraint(c.constraints[0])
		}
		return anyType

	case *UnionConstraint:
		// 从联合约束推断，返回满足任一约束的类型
		// 这里简化处理，返回第一个约束的类型
		if len(c.constraints) > 0 {
			return ti.inferFromConstraint(c.constraints[0])
		}
		return anyType

	default:
		return anyType
	}
}

// InferFromContext 从上下文推断类型
func (ti *DefaultTypeInference) InferFromContext(context *TypeContext) (map[string]Type, error) {
	result := make(map[string]Type)

	// 从变量声明推断
	for _, varType := range context.Variables {
		if inferred := ti.inferFromType(varType); len(inferred) > 0 {
			for param, typ := range inferred {
				result[param] = typ
			}
		}
	}

	// 从函数调用推断
	for _, call := range context.FunctionCalls {
		if inferred := ti.inferFromFunctionCall(call); len(inferred) > 0 {
			for param, typ := range inferred {
				result[param] = typ
			}
		}
	}

	// 从表达式推断
	for _, expr := range context.Expressions {
		if inferred := ti.inferFromExpression(expr); len(inferred) > 0 {
			for param, typ := range inferred {
				result[param] = typ
			}
		}
	}

	return result, nil
}

// TypeContext 表示类型推断上下文
type TypeContext struct {
	Variables     map[string]Type        // 变量类型
	FunctionCalls []*FunctionCallContext // 函数调用
	Expressions   []*ExpressionContext   // 表达式
	Constraints   map[string]Constraint  // 约束
	Hints         map[string]Type        // 类型提示
}

// FunctionCallContext 函数调用上下文
type FunctionCallContext struct {
	FunctionName string
	Arguments    []Type
	ReturnType   Type
}

// ExpressionContext 表达式上下文
type ExpressionContext struct {
	Expression string
	Type       Type
}

// inferFromFunctionCall 从函数调用推断类型
func (ti *DefaultTypeInference) inferFromFunctionCall(call *FunctionCallContext) map[string]Type {
	result := make(map[string]Type)

	// 这里需要根据函数名和参数类型进行推断
	// 简化实现，实际需要更复杂的函数重载解析
	if call.FunctionName == "identity" && len(call.Arguments) == 1 {
		// 对于 identity 函数，返回类型应该与参数类型相同
		result["T"] = call.Arguments[0]
	}

	return result
}

// inferFromExpression 从表达式推断类型
func (ti *DefaultTypeInference) inferFromExpression(expr *ExpressionContext) map[string]Type {
	result := make(map[string]Type)

	// 从表达式类型推断
	if inferred := ti.inferFromType(expr.Type); len(inferred) > 0 {
		for param, typ := range inferred {
			result[param] = typ
		}
	}

	return result
}

// InferBidirectional 双向类型推断
func (ti *DefaultTypeInference) InferBidirectional(target Type, source Type) (map[string]Type, error) {
	result := make(map[string]Type)

	// 从目标类型推断
	if targetInferred := ti.inferFromType(target); len(targetInferred) > 0 {
		for param, typ := range targetInferred {
			result[param] = typ
		}
	}

	// 从源类型推断
	if sourceInferred := ti.inferFromType(source); len(sourceInferred) > 0 {
		for param, typ := range sourceInferred {
			// 如果目标类型已经推断出该参数，检查兼容性
			if existing, exists := result[param]; exists {
				if !typ.AssignableTo(existing) && !existing.AssignableTo(typ) {
					return nil, fmt.Errorf("type inference conflict for parameter %s: %s vs %s",
						param, existing.Name(), typ.Name())
				}
				// 选择更具体的类型
				if ti.isMoreSpecific(typ, existing) {
					result[param] = typ
				}
			} else {
				result[param] = typ
			}
		}
	}

	return result, nil
}

// isMoreSpecific 检查类型是否更具体
func (ti *DefaultTypeInference) isMoreSpecific(t1, t2 Type) bool {
	// 具体类型比 anyType 更具体
	if t2.Name() == "any" && t1.Name() != "any" {
		return true
	}

	// 非泛型类型比泛型类型更具体
	if !ti.isGenericType(t1) && ti.isGenericType(t2) {
		return true
	}

	// 这里可以添加更多具体的比较规则
	return false
}

// isGenericType 检查是否为泛型类型
func (ti *DefaultTypeInference) isGenericType(t Type) bool {
	switch v := t.(type) {
	case *ClassType:
		return v.IsGeneric()
	case *FunctionType:
		return v.IsGeneric()
	case *TypeParameter:
		return true
	default:
		return false
	}
}

// InferFromPattern 从模式匹配推断类型
func (ti *DefaultTypeInference) InferFromPattern(pattern PatternMatch) (map[string]Type, error) {
	result := make(map[string]Type)

	switch p := pattern.(type) {
	case *VariablePattern:
		// 变量模式：绑定变量到类型
		result[p.VariableName] = p.Type

	case *ConstructorPattern:
		// 构造函数模式：推断构造函数的类型参数
		if inferred := ti.inferFromType(p.ConstructorType); len(inferred) > 0 {
			for param, typ := range inferred {
				result[param] = typ
			}
		}

	case *TuplePattern:
		// 元组模式：推断元组元素的类型
		for i := range p.Elements {
			if inferred := ti.inferFromType(p.TupleType.GetElementType(i)); len(inferred) > 0 {
				for param, typ := range inferred {
					result[param] = typ
				}
			}
		}

	case *ArrayPattern:
		// 数组模式：推断数组元素类型
		if arrayType, ok := p.ArrayType.(*ArrayType); ok {
			if inferred := ti.inferFromType(arrayType.ElementType()); len(inferred) > 0 {
				for param, typ := range inferred {
					result[param] = typ
				}
			}
		}
	}

	return result, nil
}

// PatternMatch 模式匹配接口
type PatternMatch interface {
	Match(value Value) bool
}

// VariablePattern 变量模式
type VariablePattern struct {
	VariableName string
	Type         Type
}

func (vp *VariablePattern) Match(value Value) bool {
	return value.Type().AssignableTo(vp.Type)
}

// ConstructorPattern 构造函数模式
type ConstructorPattern struct {
	ConstructorName string
	ConstructorType Type
	Arguments       []PatternMatch
}

func (cp *ConstructorPattern) Match(value Value) bool {
	// 简化实现，实际需要更复杂的匹配逻辑
	return value.Type().Name() == cp.ConstructorName
}

// TuplePattern 元组模式
type TuplePattern struct {
	TupleType *TupleType
	Elements  []PatternMatch
}

func (tp *TuplePattern) Match(value Value) bool {
	if tupleValue, ok := value.(*TupleValue); ok {
		elements := tupleValue.Elements()
		if len(elements) != len(tp.Elements) {
			return false
		}
		for i, elemPattern := range tp.Elements {
			if !elemPattern.Match(elements[i]) {
				return false
			}
		}
		return true
	}
	return false
}

// ArrayPattern 数组模式
type ArrayPattern struct {
	ArrayType Type
	Elements  []PatternMatch
}

func (ap *ArrayPattern) Match(value Value) bool {
	if arrayValue, ok := value.(*ArrayValue); ok {
		if len(arrayValue.Elements) != len(ap.Elements) {
			return false
		}
		for i, elemPattern := range ap.Elements {
			if !elemPattern.Match(arrayValue.Elements[i]) {
				return false
			}
		}
		return true
	}
	return false
}

// SolveConstraints 求解约束系统
func (ti *DefaultTypeInference) SolveConstraints(constraints map[string]Constraint) (map[string]Type, error) {
	result := make(map[string]Type)

	// 第一遍：处理简单约束
	for param, constraint := range constraints {
		if typ := ti.inferFromConstraint(constraint); typ != nil {
			result[param] = typ
		}
	}

	// 第二遍：处理复合约束
	for param, constraint := range constraints {
		if composite, ok := constraint.(*CompositeConstraint); ok {
			// 检查是否所有约束都能满足
			satisfied := true
			for _, subConstraint := range composite.constraints {
				if typ := ti.inferFromConstraint(subConstraint); typ == nil {
					satisfied = false
					break
				}
			}
			if satisfied {
				// 选择最具体的类型
				var mostSpecific Type
				for _, subConstraint := range composite.constraints {
					if typ := ti.inferFromConstraint(subConstraint); typ != nil {
						if mostSpecific == nil || ti.isMoreSpecific(typ, mostSpecific) {
							mostSpecific = typ
						}
					}
				}
				if mostSpecific != nil {
					result[param] = mostSpecific
				}
			}
		}
	}

	// 第三遍：处理联合约束
	for param, constraint := range constraints {
		if union, ok := constraint.(*UnionConstraint); ok {
			// 选择第一个满足的约束
			for _, subConstraint := range union.constraints {
				if typ := ti.inferFromConstraint(subConstraint); typ != nil {
					result[param] = typ
					break
				}
			}
		}
	}

	return result, nil
}

// InferFromOverload 从重载函数推断类型
func (ti *DefaultTypeInference) InferFromOverload(overloads []*FunctionType, args []Value) (map[string]Type, error) {
	result := make(map[string]Type)

	// 找到最佳匹配的重载
	var bestMatch *FunctionType
	var bestScore int = -1

	for _, overload := range overloads {
		score := ti.calculateOverloadScore(overload, args)
		if score > bestScore {
			bestScore = score
			bestMatch = overload
		}
	}

	if bestMatch != nil {
		// 从最佳匹配推断类型参数
		if inferred, err := ti.inferFromFunctionSignature(bestMatch, bestMatch); err == nil && len(inferred) > 0 {
			for i, typ := range inferred {
				if i < len(bestMatch.TypeParameters()) {
					result[bestMatch.TypeParameters()[i]] = typ
				}
			}
		}
	}

	return result, nil
}

// calculateOverloadScore 计算重载匹配分数
func (ti *DefaultTypeInference) calculateOverloadScore(overload *FunctionType, args []Value) int {
	score := 0

	if len(overload.signatures) == 0 {
		return -1
	}

	sig := overload.signatures[0]

	// 检查参数数量匹配
	if len(sig.positionalParams) != len(args) {
		return -1
	}

	// 计算类型匹配分数
	for i, param := range sig.positionalParams {
		if i < len(args) {
			argType := args[i].Type()

			// 精确匹配：+100分
			if argType.Name() == param.typ.Name() {
				score += 100
			} else if argType.AssignableTo(param.typ) {
				// 兼容匹配：+50分
				score += 50
			} else if param.typ.Name() == "any" {
				// any类型匹配：+10分
				score += 10
			} else {
				// 不匹配：-1分
				return -1
			}
		}
	}

	return score
}

// UnifyTypes 类型统一算法
func (ti *DefaultTypeInference) UnifyTypes(t1, t2 Type) (map[string]Type, error) {
	unifier := make(map[string]Type)

	if err := ti.unify(t1, t2, unifier); err != nil {
		return nil, err
	}

	return unifier, nil
}

// unify 递归统一类型
func (ti *DefaultTypeInference) unify(t1, t2 Type, unifier map[string]Type) error {
	// 如果两个类型相同，直接返回
	if t1.Name() == t2.Name() {
		return nil
	}

	// 处理类型参数
	if param1, ok := t1.(*TypeParameter); ok {
		return ti.unifyTypeParameter(param1, t2, unifier)
	}

	if param2, ok := t2.(*TypeParameter); ok {
		return ti.unifyTypeParameter(param2, t1, unifier)
	}

	// 处理数组类型
	if array1, ok1 := t1.(*ArrayType); ok1 {
		if array2, ok2 := t2.(*ArrayType); ok2 {
			return ti.unify(array1.ElementType(), array2.ElementType(), unifier)
		}
	}

	// 处理映射类型
	if map1, ok1 := t1.(*MapType); ok1 {
		if map2, ok2 := t2.(*MapType); ok2 {
			if err := ti.unify(map1.KeyType(), map2.KeyType(), unifier); err != nil {
				return err
			}
			return ti.unify(map1.ValueType(), map2.ValueType(), unifier)
		}
	}

	// 处理函数类型
	if func1, ok1 := t1.(*FunctionType); ok1 {
		if func2, ok2 := t2.(*FunctionType); ok2 {
			return ti.unifyFunctionTypes(func1, func2, unifier)
		}
	}

	// 处理类类型
	if class1, ok1 := t1.(*ClassType); ok1 {
		if class2, ok2 := t2.(*ClassType); ok2 {
			return ti.unifyClassTypes(class1, class2, unifier)
		}
	}

	// 如果类型不兼容，返回错误
	return fmt.Errorf("cannot unify types %s and %s", t1.Name(), t2.Name())
}

// unifyTypeParameter 统一类型参数
func (ti *DefaultTypeInference) unifyTypeParameter(param *TypeParameter, t Type, unifier map[string]Type) error {
	// 检查是否已经绑定
	if existing, exists := unifier[param.Name()]; exists {
		return ti.unify(existing, t, unifier)
	}

	// 检查是否出现循环引用
	if ti.occursIn(param.Name(), t, unifier) {
		return fmt.Errorf("circular type reference detected for parameter %s", param.Name())
	}

	// 绑定类型参数
	unifier[param.Name()] = t
	return nil
}

// occursIn 检查类型参数是否出现在类型中
func (ti *DefaultTypeInference) occursIn(paramName string, t Type, unifier map[string]Type) bool {
	switch v := t.(type) {
	case *TypeParameter:
		if v.Name() == paramName {
			return true
		}
		if replacement, exists := unifier[v.Name()]; exists {
			return ti.occursIn(paramName, replacement, unifier)
		}
		return false

	case *ArrayType:
		return ti.occursIn(paramName, v.ElementType(), unifier)

	case *MapType:
		return ti.occursIn(paramName, v.KeyType(), unifier) ||
			ti.occursIn(paramName, v.ValueType(), unifier)

	case *FunctionType:
		return ti.occursInFunctionType(paramName, v, unifier)

	case *ClassType:
		return ti.occursInClassType(paramName, v, unifier)

	default:
		return false
	}
}

// occursInFunctionType 检查类型参数是否出现在函数类型中
func (ti *DefaultTypeInference) occursInFunctionType(paramName string, funcType *FunctionType, unifier map[string]Type) bool {
	for _, sig := range funcType.signatures {
		// 检查参数类型
		for _, param := range sig.positionalParams {
			if ti.occursIn(paramName, param.typ, unifier) {
				return true
			}
		}

		// 检查命名参数类型
		for _, param := range sig.namedParams {
			if ti.occursIn(paramName, param.typ, unifier) {
				return true
			}
		}

		// 检查返回类型
		if ti.occursIn(paramName, sig.returnType, unifier) {
			return true
		}
	}
	return false
}

// occursInClassType 检查类型参数是否出现在类类型中
func (ti *DefaultTypeInference) occursInClassType(paramName string, classType *ClassType, unifier map[string]Type) bool {
	// 检查字段类型
	for _, field := range classType.fields {
		if ti.occursIn(paramName, field.ft, unifier) {
			return true
		}
	}

	// 检查方法类型
	for _, method := range classType.methods {
		if funcType, ok := method.(*FunctionType); ok {
			if ti.occursInFunctionType(paramName, funcType, unifier) {
				return true
			}
		}
	}

	return false
}

// unifyFunctionTypes 统一函数类型
func (ti *DefaultTypeInference) unifyFunctionTypes(func1, func2 *FunctionType, unifier map[string]Type) error {
	if len(func1.signatures) == 0 || len(func2.signatures) == 0 {
		return fmt.Errorf("function type has no signatures")
	}

	sig1 := func1.signatures[0]
	sig2 := func2.signatures[0]

	// 统一参数类型
	if len(sig1.positionalParams) != len(sig2.positionalParams) {
		return fmt.Errorf("function parameter count mismatch")
	}

	for i := range sig1.positionalParams {
		if err := ti.unify(sig1.positionalParams[i].typ, sig2.positionalParams[i].typ, unifier); err != nil {
			return err
		}
	}

	// 统一返回类型
	return ti.unify(sig1.returnType, sig2.returnType, unifier)
}

// unifyClassTypes 统一类类型
func (ti *DefaultTypeInference) unifyClassTypes(class1, class2 *ClassType, unifier map[string]Type) error {
	// 检查类名是否相同
	if class1.Name() != class2.Name() {
		return fmt.Errorf("cannot unify different class types: %s and %s", class1.Name(), class2.Name())
	}

	// 如果都是泛型类，统一类型参数
	if class1.IsGeneric() && class2.IsGeneric() {
		params1 := class1.TypeParameters()
		params2 := class2.TypeParameters()

		if len(params1) != len(params2) {
			return fmt.Errorf("generic class type parameter count mismatch")
		}

		// 这里简化处理，实际需要更复杂的类型参数统一
		for i := range params1 {
			if params1[i] != params2[i] {
				return fmt.Errorf("generic class type parameter mismatch")
			}
		}
	}

	return nil
}

// PropagateConstraints 约束传播
func (ti *DefaultTypeInference) PropagateConstraints(constraints map[string]Constraint, unifier map[string]Type) (map[string]Constraint, error) {
	propagated := make(map[string]Constraint)

	// 应用统一器到约束
	for param, constraint := range constraints {
		// 如果参数已经被统一，传播约束
		if unifiedType, exists := unifier[param]; exists {
			// 创建新的约束，要求统一后的类型满足原约束
			propagatedConstraint := ti.createPropagatedConstraint(constraint, unifiedType)
			if propagatedConstraint != nil {
				propagated[param] = propagatedConstraint
			}
		} else {
			// 保持原约束
			propagated[param] = constraint
		}
	}

	return propagated, nil
}

// createPropagatedConstraint 创建传播的约束
func (ti *DefaultTypeInference) createPropagatedConstraint(original Constraint, unifiedType Type) Constraint {
	// 这里简化实现，实际需要更复杂的约束传播逻辑
	switch c := original.(type) {
	case *InterfaceConstraint:
		// 如果统一后的类型已经实现了接口，返回空约束
		if unifiedType.Implements(c.interfaceType) {
			return nil
		}
		// 否则保持原约束
		return c

	case *OperatorConstraint:
		// 检查统一后的类型是否支持运算符
		if objType, ok := unifiedType.(*ObjectType); ok {
			if _, exists := objType.operators[c.operator]; exists {
				return nil
			}
		}
		return c

	default:
		return original
	}
}

// InferRecursive 递归类型推断
func (ti *DefaultTypeInference) InferRecursive(types []Type) (map[string]Type, error) {
	result := make(map[string]Type)

	// 构建约束系统
	constraints := make(map[string]Constraint)

	// 从所有类型中收集约束
	for _, t := range types {
		if inferred := ti.inferFromType(t); len(inferred) > 0 {
			for param, typ := range inferred {
				// 为每个类型参数创建约束
				constraints[param] = NewInterfaceConstraint(typ)
			}
		}
	}

	// 求解约束系统
	solved, err := ti.SolveConstraints(constraints)
	if err != nil {
		return nil, err
	}

	// 合并结果
	maps.Copy(result, solved)

	return result, nil
}

// InferWithContext 带上下文的类型推断
func (ti *DefaultTypeInference) InferWithContext(target Type, context *TypeContext) (map[string]Type, error) {
	result := make(map[string]Type)

	// 从目标类型推断
	if targetInferred := ti.inferFromType(target); len(targetInferred) > 0 {
		for param, typ := range targetInferred {
			result[param] = typ
		}
	}

	// 从上下文推断
	if contextInferred, err := ti.InferFromContext(context); err == nil {
		for param, typ := range contextInferred {
			// 检查冲突
			if existing, exists := result[param]; exists {
				if !typ.AssignableTo(existing) && !existing.AssignableTo(typ) {
					return nil, fmt.Errorf("type inference conflict for parameter %s: %s vs %s",
						param, existing.Name(), typ.Name())
				}
				// 选择更具体的类型
				if ti.isMoreSpecific(typ, existing) {
					result[param] = typ
				}
			} else {
				result[param] = typ
			}
		}
	}

	return result, nil
}

// SubstituteTypes 类型替换
func (ti *DefaultTypeInference) SubstituteTypes(t Type, substitutions map[string]Type) Type {
	switch v := t.(type) {
	case *TypeParameter:
		if replacement, exists := substitutions[v.Name()]; exists {
			return replacement
		}
		return t

	case *ArrayType:
		return &ArrayType{
			ObjectType:  NewObjectType(ti.SubstituteTypes(v.ElementType(), substitutions).Name()+"[]", O_ARR),
			elementType: ti.SubstituteTypes(v.ElementType(), substitutions),
		}

	case *MapType:
		return &MapType{
			ObjectType: NewObjectType(fmt.Sprintf("map<%s, %s>",
				ti.SubstituteTypes(v.KeyType(), substitutions).Name(),
				ti.SubstituteTypes(v.ValueType(), substitutions).Name()), O_MAP),
			keyType:   ti.SubstituteTypes(v.KeyType(), substitutions),
			valueType: ti.SubstituteTypes(v.ValueType(), substitutions),
		}

	case *FunctionType:
		return ti.substituteFunctionType(v, substitutions)

	case *ClassType:
		return ti.substituteClassType(v, substitutions)

	default:
		return t
	}
}

// substituteFunctionType 替换函数类型中的类型参数
func (ti *DefaultTypeInference) substituteFunctionType(funcType *FunctionType, substitutions map[string]Type) *FunctionType {
	newFunc := &FunctionType{
		name:        funcType.name,
		signatures:  make([]*FunctionSignature, len(funcType.signatures)),
		GenericBase: funcType.GenericBase,
		visible:     funcType.visible,
	}

	for i, sig := range funcType.signatures {
		newSig := &FunctionSignature{
			positionalParams: make([]*Parameter, len(sig.positionalParams)),
			namedParams:      make([]*Parameter, len(sig.namedParams)),
			returnType:       ti.SubstituteTypes(sig.returnType, substitutions),
			isOperator:       sig.isOperator,
			operator:         sig.operator,
			builtin:          sig.builtin,
		}

		// 替换位置参数
		for j, param := range sig.positionalParams {
			newSig.positionalParams[j] = &Parameter{
				name:     param.name,
				typ:      ti.SubstituteTypes(param.typ, substitutions),
				optional: param.optional,
				variadic: param.variadic,
				isNamed:  param.isNamed,
				required: param.required,
			}
		}

		// 替换命名参数
		for j, param := range sig.namedParams {
			newSig.namedParams[j] = &Parameter{
				name:     param.name,
				typ:      ti.SubstituteTypes(param.typ, substitutions),
				optional: param.optional,
				variadic: param.variadic,
				isNamed:  param.isNamed,
				required: param.required,
			}
		}

		newFunc.signatures[i] = newSig
	}

	return newFunc
}

// substituteClassType 替换类类型中的类型参数
func (ti *DefaultTypeInference) substituteClassType(classType *ClassType, substitutions map[string]Type) *ClassType {
	newClass := &ClassType{
		ObjectType:        NewObjectType(classType.Name(), O_CLASS),
		GenericBase:       classType.GenericBase,
		astDecl:           classType.astDecl,
		ctors:             classType.ctors,
		registry:          classType.registry,
		staticFields:      classType.staticFields,
		staticMethods:     classType.staticMethods,
		staticValues:      classType.staticValues,
		parent:            classType.parent,
		parentName:        classType.parentName,
		interfaces:        classType.interfaces,
		overriddenMethods: classType.overriddenMethods,
		traitConstraints:  classType.traitConstraints,
	}

	// 替换字段类型
	for name, field := range classType.fields {
		newClass.fields[name] = &Field{
			ft:         ti.SubstituteTypes(field.ft, substitutions),
			mods:       field.mods,
			getter:     field.getter,
			setter:     field.setter,
			decorators: field.decorators,
		}
	}

	// 替换方法类型
	for name, method := range classType.methods {
		if funcType, ok := method.(*FunctionType); ok {
			newClass.methods[name] = ti.substituteFunctionType(funcType, substitutions)
		} else {
			newClass.methods[name] = method
		}
	}

	// 替换运算符类型
	for op, funcType := range classType.operators {
		newClass.operators[op] = ti.substituteFunctionType(funcType, substitutions)
	}

	return newClass
}

// CheckConstraintSatisfaction 检查约束满足性
func (ti *DefaultTypeInference) CheckConstraintSatisfaction(t Type, constraint Constraint) bool {
	return constraint.SatisfiedBy(t)
}

// InferFromLiteral 从字面量推断类型
func (ti *DefaultTypeInference) InferFromLiteral(literal Value) Type {
	switch literal.(type) {
	case *NumberValue:
		return numberType
	case *StringValue:
		return stringType
	case *BoolValue:
		return boolType
	case *ArrayValue:
		if arrayValue, ok := literal.(*ArrayValue); ok {
			if len(arrayValue.Elements) > 0 {
				elementType := ti.InferFromLiteral(arrayValue.Elements[0])
				return NewArrayType(elementType)
			}
			return NewArrayType(anyType)
		}
	case *MapValue:
		if mapValue, ok := literal.(*MapValue); ok {
			if len(mapValue.Pairs) > 0 {
				keyType := ti.InferFromLiteral(mapValue.Pairs[0].Key)
				valueType := ti.InferFromLiteral(mapValue.Pairs[0].Value)
				return NewMapType(keyType, valueType)
			}
			return NewMapType(anyType, anyType)
		}
	case *TupleValue:
		if tupleValue, ok := literal.(*TupleValue); ok {
			elementTypes := make([]Type, len(tupleValue.Elements()))
			for i, elem := range tupleValue.Elements() {
				elementTypes[i] = ti.InferFromLiteral(elem)
			}
			return NewTupleType(elementTypes...)
		}
	case *SetValue:
		if setValue, ok := literal.(*SetValue); ok {
			if len(setValue.Elements()) > 0 {
				// 获取第一个元素的类型
				for _, elem := range setValue.Elements() {
					elementType := ti.InferFromLiteral(elem)
					return NewSetType(elementType)
				}
			}
			return NewSetType(anyType)
		}
	}

	return anyType
}

// InferFromExpression 从表达式推断类型（增强版）
func (ti *DefaultTypeInference) InferFromExpression(expr *ExpressionContext) map[string]Type {
	result := make(map[string]Type)

	// 从表达式类型推断
	if inferred := ti.inferFromType(expr.Type); len(inferred) > 0 {
		for param, typ := range inferred {
			result[param] = typ
		}
	}

	// 根据表达式内容进行特殊推断
	switch expr.Expression {
	case "[]":
		// 空数组字面量
		result["T"] = anyType
	case "{}":
		// 空映射字面量
		result["K"] = anyType
		result["V"] = anyType
	case "()":
		// 空元组字面量
		// 不添加类型参数
	}

	return result
}

// InferFromBinaryOp 从二元操作推断类型
func (ti *DefaultTypeInference) InferFromBinaryOp(op Operator, left, right Type) (map[string]Type, error) {
	result := make(map[string]Type)

	switch op {
	case OpAdd, OpSub, OpMul, OpDiv, OpMod:
		// 算术运算：结果类型通常是操作数的统一类型
		if unifier, err := ti.UnifyTypes(left, right); err == nil {
			for param, typ := range unifier {
				result[param] = typ
			}
		}

	case OpEqual, OpLess, OpLessEqual, OpGreater, OpGreaterEqual:
		// 比较运算：结果类型是 bool
		// 操作数类型需要兼容
		if unifier, err := ti.UnifyTypes(left, right); err == nil {
			for param, typ := range unifier {
				result[param] = typ
			}
		}

	case OpConcat:
		// 字符串连接：操作数都应该是字符串类型
		if left.Name() == "str" && right.Name() == "str" {
			// 都是字符串，不需要推断
		} else {
			// 尝试统一为字符串类型
			result["T"] = stringType
		}
	}

	return result, nil
}

// InferFromUnaryOp 从一元操作推断类型
func (ti *DefaultTypeInference) InferFromUnaryOp(op Operator, operand Type) (map[string]Type, error) {
	result := make(map[string]Type)

	switch op {
	case OpNot:
		// 逻辑非：操作数应该是 bool 类型
		if operand.Name() != "bool" {
			result["T"] = boolType
		}

	case OpSub:
		// 负号：操作数应该是数值类型
		if operand.Name() != "num" {
			result["T"] = numberType
		}
	}

	return result, nil
}

// InferFromIndexOp 从索引操作推断类型
func (ti *DefaultTypeInference) InferFromIndexOp(container, index Type) (map[string]Type, error) {
	result := make(map[string]Type)

	// 数组索引
	if arrayType, ok := container.(*ArrayType); ok {
		result["T"] = arrayType.ElementType()
		if index.Name() != "num" {
			result["I"] = numberType
		}
	}

	// 映射索引
	if mapType, ok := container.(*MapType); ok {
		result["V"] = mapType.ValueType()
		if !index.AssignableTo(mapType.KeyType()) {
			result["K"] = mapType.KeyType()
		}
	}

	return result, nil
}

// InferFromCall 从函数调用推断类型
func (ti *DefaultTypeInference) InferFromCall(funcType Type, args []Type) (map[string]Type, error) {
	result := make(map[string]Type)

	if functionType, ok := funcType.(*FunctionType); ok {
		if len(functionType.signatures) > 0 {
			sig := functionType.signatures[0]

			// 统一参数类型
			for i, arg := range args {
				if i < len(sig.positionalParams) {
					if unifier, err := ti.UnifyTypes(arg, sig.positionalParams[i].typ); err == nil {
						maps.Copy(result, unifier)
					}
				}
			}

			// 从返回类型推断
			if inferred := ti.inferFromType(sig.returnType); len(inferred) > 0 {
				maps.Copy(result, inferred)
			}
		}
	}

	return result, nil
}
