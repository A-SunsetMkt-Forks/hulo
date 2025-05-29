package build

import (
	"errors"

	"github.com/Masterminds/semver/v3"
	"github.com/hulo-lang/hulo/internal/build"
	bast "github.com/hulo-lang/hulo/syntax/bash/ast"
	hast "github.com/hulo-lang/hulo/syntax/hulo/ast"
)

type CompileRuleFunc func(node hast.Node) (bast.Node, error)

type HCRDispatcher struct {
	rules map[string]map[*semver.Constraints]CompileRuleFunc
}

func (d *HCRDispatcher) Put(fullRule string, cb CompileRuleFunc) error {
	hcr, err := build.ParseRule(fullRule)
	if err != nil {
		return err
	}

	c, err := semver.NewConstraint(hcr.Version().Original())
	if err != nil {
		return err
	}

	d.rules[hcr.Name()][c] = cb

	return nil
}

func (d *HCRDispatcher) Get(ruleName string) (CompileRuleFunc, error) {
	hcr, err := build.ParseRule(ruleName)
	if err != nil {
		return nil, err
	}

	for rule, fn := range d.rules[hcr.Name()] {
		if ok, _ := rule.Validate(hcr.Version()); ok {
			return fn, nil
		}
	}
	return nil, errors.New("rule not found")
}

var hcrDispatcher = &HCRDispatcher{rules: make(map[string]map[*semver.Constraints]CompileRuleFunc)}
