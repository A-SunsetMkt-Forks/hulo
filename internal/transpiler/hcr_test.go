package transpiler

import (
	"testing"

	"github.com/hulo-lang/hulo/internal/linker"
	"github.com/hulo-lang/hulo/syntax/hulo/ast"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockRule 是一个用于测试的模拟规则
type MockRule struct {
	strategy string
	applyFn  func(Transpiler[string], ast.Node) (string, error)
}

func (m *MockRule) Strategy() string {
	return m.strategy
}

func (m *MockRule) Apply(transpiler Transpiler[string], node ast.Node) (string, error) {
	if m.applyFn != nil {
		return m.applyFn(transpiler, node)
	}
	return "mock_result", nil
}

// MockTranspiler 是一个用于测试的模拟转换器
type MockTranspiler struct{}

func (m *MockTranspiler) Convert(node ast.Node) string {
	return "converted"
}

func (m *MockTranspiler) GetTargetExt() string {
	return ".mock"
}

func (m *MockTranspiler) GetTargetName() string {
	return "mock"
}

func (m *MockTranspiler) UnresolvedSymbols() map[string][]linker.UnkownSymbol {
	return make(map[string][]linker.UnkownSymbol)
}

// MockRuleID 是一个用于测试的规则ID
type MockRuleID string

func (m MockRuleID) String() string {
	return string(m)
}

func TestNewHCRDispatcher(t *testing.T) {
	dispatcher := NewHCRDispatcher[string]()

	assert.NotNil(t, dispatcher)
	assert.NotNil(t, dispatcher.rules)
	assert.NotNil(t, dispatcher.cached)
	assert.Empty(t, dispatcher.rules)
	assert.Empty(t, dispatcher.cached)
}

func TestHCRDispatcher_Register(t *testing.T) {
	dispatcher := NewHCRDispatcher[string]()

	rule1 := &MockRule{strategy: "strategy1"}
	rule2 := &MockRule{strategy: "strategy2"}

	ruleID := MockRuleID("test_rule")

	// 测试注册单个规则
	dispatcher.Register(ruleID, rule1)
	assert.Len(t, dispatcher.rules[ruleID], 1)
	assert.Equal(t, rule1, dispatcher.rules[ruleID][0])

	// 测试注册多个规则
	dispatcher.Register(ruleID, rule2)
	assert.Len(t, dispatcher.rules[ruleID], 2)
	assert.Equal(t, rule1, dispatcher.rules[ruleID][0])
	assert.Equal(t, rule2, dispatcher.rules[ruleID][1])

	// 测试注册到不同的规则ID
	ruleID2 := MockRuleID("test_rule2")
	dispatcher.Register(ruleID2, rule1)
	assert.Len(t, dispatcher.rules[ruleID2], 1)
	assert.Len(t, dispatcher.rules[ruleID], 2) // 原有规则不受影响
}

func TestHCRDispatcher_Bind(t *testing.T) {
	dispatcher := NewHCRDispatcher[string]()

	rule1 := &MockRule{strategy: "strategy1"}
	rule2 := &MockRule{strategy: "strategy2"}
	ruleID := MockRuleID("test_rule")

	// 注册规则
	dispatcher.Register(ruleID, rule1, rule2)

	// 测试绑定存在的策略
	err := dispatcher.Bind(ruleID, "strategy1")
	assert.NoError(t, err)
	assert.Equal(t, rule1, dispatcher.cached[ruleID])

	// 测试重复绑定应该失败
	err = dispatcher.Bind(ruleID, "strategy2")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already bound")
	assert.Equal(t, rule1, dispatcher.cached[ruleID]) // 应该保持原来的规则

	// 测试绑定不存在的策略
	ruleID2 := MockRuleID("test_rule2")
	err = dispatcher.Bind(ruleID2, "nonexistent")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestHCRDispatcher_Get(t *testing.T) {
	dispatcher := NewHCRDispatcher[string]()

	rule1 := &MockRule{strategy: "strategy1"}
	ruleID := MockRuleID("test_rule")

	// 测试获取未绑定的规则
	rule, err := dispatcher.Get(ruleID)
	assert.Error(t, err)
	assert.Nil(t, rule)
	assert.Contains(t, err.Error(), "not found or not bound")

	// 注册并绑定规则
	dispatcher.Register(ruleID, rule1)
	err = dispatcher.Bind(ruleID, "strategy1")
	require.NoError(t, err)

	// 测试获取已绑定的规则
	rule, err = dispatcher.Get(ruleID)
	assert.NoError(t, err)
	assert.Equal(t, rule1, rule)
}

func TestHCRDispatcher_Integration(t *testing.T) {
	dispatcher := NewHCRDispatcher[string]()

	// 创建多个规则
	rule1 := &MockRule{strategy: "number"}
	rule2 := &MockRule{strategy: "string"}
	rule3 := &MockRule{strategy: "cmd"}

	ruleID := MockRuleID("bool_format")

	// 注册所有规则
	dispatcher.Register(ruleID, rule1, rule2, rule3)

	// 验证规则已注册
	assert.Len(t, dispatcher.rules[ruleID], 3)

	// 绑定到 number 策略
	err := dispatcher.Bind(ruleID, "number")
	assert.NoError(t, err)

	// 获取并验证规则
	rule, err := dispatcher.Get(ruleID)
	assert.NoError(t, err)
	assert.Equal(t, rule1, rule)

	// 测试规则的应用
	mockTranspiler := &MockTranspiler{}
	result, err := rule.Apply(mockTranspiler, nil)
	assert.NoError(t, err)
	assert.Equal(t, "mock_result", result)
}

func TestHCRDispatcher_WithCustomApply(t *testing.T) {
	dispatcher := NewHCRDispatcher[string]()

	// 创建带有自定义 Apply 函数的规则
	customRule := &MockRule{
		strategy: "custom",
		applyFn: func(transpiler Transpiler[string], node ast.Node) (string, error) {
			return "custom_result", nil
		},
	}

	ruleID := MockRuleID("custom_rule")
	dispatcher.Register(ruleID, customRule)

	err := dispatcher.Bind(ruleID, "custom")
	assert.NoError(t, err)

	rule, err := dispatcher.Get(ruleID)
	assert.NoError(t, err)

	mockTranspiler := &MockTranspiler{}
	result, err := rule.Apply(mockTranspiler, nil)
	assert.NoError(t, err)
	assert.Equal(t, "custom_result", result)
}

func TestHCRDispatcher_EmptyRules(t *testing.T) {
	dispatcher := NewHCRDispatcher[string]()
	ruleID := MockRuleID("empty_rule")

	// 测试绑定空规则列表
	err := dispatcher.Bind(ruleID, "any_strategy")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestHCRDispatcher_MultipleRuleIDs(t *testing.T) {
	dispatcher := NewHCRDispatcher[string]()

	rule1 := &MockRule{strategy: "strategy1"}
	rule2 := &MockRule{strategy: "strategy2"}

	ruleID1 := MockRuleID("rule1")
	ruleID2 := MockRuleID("rule2")

	// 为不同的规则ID注册不同的规则
	dispatcher.Register(ruleID1, rule1)
	dispatcher.Register(ruleID2, rule2)

	// 绑定规则
	err1 := dispatcher.Bind(ruleID1, "strategy1")
	err2 := dispatcher.Bind(ruleID2, "strategy2")

	assert.NoError(t, err1)
	assert.NoError(t, err2)

	// 验证规则独立
	rule1Result, err := dispatcher.Get(ruleID1)
	assert.NoError(t, err)
	assert.Equal(t, rule1, rule1Result)

	rule2Result, err := dispatcher.Get(ruleID2)
	assert.NoError(t, err)
	assert.Equal(t, rule2, rule2Result)
}
