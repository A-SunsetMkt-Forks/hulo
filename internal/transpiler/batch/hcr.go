package transpiler

import (
	"fmt"

	bast "github.com/hulo-lang/hulo/syntax/batch/ast"
	hast "github.com/hulo-lang/hulo/syntax/hulo/ast"
)

type HuloComptimeRule interface {
	Convert(node hast.Node) bast.Node
	Name() string
}

type HCRManager struct {
	strategies map[string][]HuloComptimeRule
}

func (m *HCRManager) Register(name string, rules ...HuloComptimeRule) {
	m.strategies[name] = append(m.strategies[name], rules...)
}

func (m *HCRManager) Invoke(name string, node hast.Node) (any, error) {
	rules := m.strategies[name]
	for _, rule := range rules {
		if rule.Name() != name {
			continue
		}
		result := rule.Convert(node)
		if result != nil {
			return result, nil
		}
	}
	return nil, fmt.Errorf("no rule found for %s", name)
}

func (m *HCRManager) Lookup(name string) []HuloComptimeRule {
	return m.strategies[name]
}

type BoolAsNumberConvertor struct{}

func (c *BoolAsNumberConvertor) Strategy() string {
	return "bool-as-number"
}

func (c *BoolAsNumberConvertor) Convert(node hast.Node) bast.Node {
	return nil
}

type BoolAsStringConvertor struct{}

func (c *BoolAsStringConvertor) Strategy() string {
	return "bool-as-string"
}

func (c *BoolAsStringConvertor) Convert(node hast.Node) bast.Node {
	return nil
}

type BoolAsCmdConvertor struct{}

func (c *BoolAsCmdConvertor) Strategy() string {
	return "bool-as-cmd"
}

func (c *BoolAsCmdConvertor) Convert(node hast.Node) bast.Node {
	return nil
}

type MultiStringConvertor struct{}
