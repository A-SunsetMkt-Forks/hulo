// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package transpiler

import (
	"github.com/hulo-lang/hulo/internal/transpiler"
	"github.com/hulo-lang/hulo/internal/transpiler/bash/strategy"
	bast "github.com/hulo-lang/hulo/syntax/bash/ast"
	hast "github.com/hulo-lang/hulo/syntax/hulo/ast"
)

type CompileRuleFunc func(node hast.Node) (bast.Node, error)

type HCRDispatcher struct {
	rules map[string]transpiler.Strategy[bast.Node]
}

func (d *HCRDispatcher) Put(fullRule string, cb transpiler.Strategy[bast.Node]) error {
	hcr, err := transpiler.ParseRule(fullRule)
	if err != nil {
		return err
	}

	d.rules[hcr.Name()] = cb

	return nil
}

func (d *HCRDispatcher) Get(ruleName string) (transpiler.Strategy[bast.Node], error) {
	hcr, err := transpiler.ParseRule(ruleName)
	if err != nil {
		return nil, err
	}

	return d.rules[hcr.Name()], nil
}

var hcrDispatcher = &HCRDispatcher{rules: make(map[string]transpiler.Strategy[bast.Node])}

func init() {
	hcrDispatcher.Put("number", &strategy.BooleanAsNumberStrategy{})
	hcrDispatcher.Put("string", &strategy.BooleanAsStringStrategy{})
	hcrDispatcher.Put("command", &strategy.BooleanAsCommandStrategy{})
}
