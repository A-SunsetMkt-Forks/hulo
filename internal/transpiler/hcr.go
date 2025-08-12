// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package transpiler

import (
	"fmt"

	"github.com/hulo-lang/hulo/syntax/hulo/ast"
)

// HuloCompileRule is the interface for the compile rule.
type HuloCompileRule[T any] interface {
	// Strategy returns the strategy of the rule.
	Strategy() string
	// Apply applies the rule to the node.
	Apply(transpiler Transpiler[T], node ast.Node) (T, error)
}

// HCRDispatcher is the dispatcher for the compile rule.
// It is used to dispatch the compile rule to the correct rule.
type HCRDispatcher[T any] struct {
	rules  map[RuleID][]HuloCompileRule[T]
	cached map[RuleID]HuloCompileRule[T]
}

func NewHCRDispatcher[T any]() *HCRDispatcher[T] {
	return &HCRDispatcher[T]{
		rules:  make(map[RuleID][]HuloCompileRule[T]),
		cached: make(map[RuleID]HuloCompileRule[T]),
	}
}

func (hcrd *HCRDispatcher[T]) Register(name RuleID, rules ...HuloCompileRule[T]) {
	hcrd.rules[name] = append(hcrd.rules[name], rules...)
}

func (hcrd *HCRDispatcher[T]) Get(key RuleID) (HuloCompileRule[T], error) {
	rule, ok := hcrd.cached[key]
	if ok {
		return rule, nil
	}
	return nil, fmt.Errorf("rule %s not found or not bound", key)
}

func (hcrd *HCRDispatcher[T]) Bind(key RuleID, value string) error {
	_, ok := hcrd.cached[key]
	if ok {
		return fmt.Errorf("rule %s already bound", key)
	}
	for _, rule := range hcrd.rules[key] {
		if rule.Strategy() == value {
			hcrd.cached[key] = rule
			return nil
		}
	}
	return fmt.Errorf("rule %s not found", key)
}

type RuleID interface {
	String() string
}
