package build

type CompileRuleFunc func() error

type HCRDispatcher struct {
	rules map[string]CompileRuleFunc
}

func (d *HCRDispatcher) Register(name string, fn CompileRuleFunc) {
	d.rules[name] = fn
}
