// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package object

import (
	"fmt"
	"sync"

	"github.com/hulo-lang/hulo/internal/util"
)

type Registry struct {
	functions map[string]*FunctionType
	types     map[string]Type
	overloads map[Operator][]*OperatorOverload
	mutex     sync.Locker
}

func NewRegistry() *Registry {
	return &Registry{
		functions: make(map[string]*FunctionType),
		types:     make(map[string]Type),
		mutex:     &util.NoOpLocker{},
	}
}

func (r *Registry) RegisterFunction(name string, signature *FunctionSignature) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if r.functions[name] == nil {
		r.functions[name] = &FunctionType{
			name:       name,
			signatures: make([]*FunctionSignature, 0),
		}
	}

	// 检查签名冲突
	for _, existingSig := range r.functions[name].signatures {
		if r.signaturesConflict(existingSig, signature) {
			return fmt.Errorf("function signature conflict for %s", name)
		}
	}

	r.functions[name].signatures = append(r.functions[name].signatures, signature)
	return nil
}

func (r *Registry) RegisterOperator(op Operator, left, right, returnType Type, fn *FunctionType) {
	r.overloads[op] = append(r.overloads[op], &OperatorOverload{
		operator:   op,
		leftType:   left,
		rightType:  right,
		returnType: returnType,
		function:   fn,
	})
}

// 检查签名冲突
func (r *Registry) signaturesConflict(sig1, sig2 *FunctionSignature) bool {
	if len(sig1.positionalParams) != len(sig2.positionalParams) {
		return false
	}

	// 检查位置参数类型是否完全相同
	for i, param1 := range sig1.positionalParams {
		param2 := sig2.positionalParams[i]
		if param1.typ.Name() != param2.typ.Name() {
			return false
		}
	}

	// 检查命名参数类型是否完全相同
	for i, param1 := range sig1.namedParams {
		param2 := sig2.namedParams[i]
		if param1.typ.Name() != param2.typ.Name() {
			return false
		}
	}

	return true
}

func (r *Registry) LookupFunction(name string) (*FunctionType, bool) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	fn, exists := r.functions[name]
	return fn, exists
}

func (r *Registry) LookupOperator(op Operator, left, right Value) (*OperatorOverload, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	for _, overload := range r.overloads[op] {
		if left.Type().AssignableTo(overload.leftType) &&
			right.Type().AssignableTo(overload.rightType) {
			return overload, nil
		}
	}
	return nil, fmt.Errorf("no operator overlaod found for %v", op)
}
