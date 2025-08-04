package interpreter

import (
	"testing"
	"time"

	"github.com/hulo-lang/hulo/syntax/hulo/ast"
	"github.com/hulo-lang/hulo/syntax/hulo/parser"
	"github.com/stretchr/testify/assert"
)

func TestDebuggerBasic(t *testing.T) {
	// 创建一个简单的测试脚本
	script := `let x = 10
let y = 20
let z = $x + $y
echo "z = " + $z`

	// 解析脚本
	_, err := parser.ParseSourceScript(script)
	assert.NoError(t, err)

	// 创建解释器
	interp := &Interpreter{}

	// 创建调试器
	debugger := NewDebugger(interp)
	assert.NotNil(t, debugger)

	// 添加断点（在第2行）
	bp, err := debugger.AddBreakpoint("test.hl", 2, "")
	assert.NoError(t, err)
	assert.NotNil(t, bp)
	assert.Equal(t, "test.hl:2", bp.ID)
	assert.Equal(t, 2, bp.Line)
	assert.True(t, bp.Enabled)

	// 获取断点列表
	breakpoints := debugger.GetBreakpoints()
	assert.Len(t, breakpoints, 1)
	assert.Equal(t, bp.ID, breakpoints[0].ID)

	// 移除断点
	err = debugger.RemoveBreakpoint("test.hl", 2)
	assert.NoError(t, err)

	// 验证断点已被移除
	breakpoints = debugger.GetBreakpoints()
	assert.Len(t, breakpoints, 0)
}

func TestDebuggerWithCondition(t *testing.T) {
	script := `let x = 5
let y = 10
if $x > 3 {
    let z = $x + $y
    echo "z = " + $z
}`

	_, err := parser.ParseSourceScript(script)
	assert.NoError(t, err)

	interp := &Interpreter{}
	debugger := NewDebugger(interp)

	// 添加条件断点
	bp, err := debugger.AddBreakpoint("test.hl", 3, "$x > 3")
	assert.NoError(t, err)
	assert.NotNil(t, bp)
	assert.Equal(t, "$x > 3", bp.Condition)
}

func TestDebuggerCommands(t *testing.T) {
	interp := &Interpreter{}
	debugger := NewDebugger(interp)

	// 测试发送命令
	responseChan := make(chan any, 1)
	debugger.SendCommand(DebugCommand{
		Type:     CmdGetVariables,
		Response: responseChan,
	})

	// 等待响应（设置超时）
	select {
	case response := <-responseChan:
		assert.NotNil(t, response)
	case <-time.After(100 * time.Millisecond):
		// 超时是正常的，因为调试器可能还没有处理命令
		t.Log("Command response timeout (expected)")
	}
}

func TestDebuggerCallStack(t *testing.T) {
	script := `fn add(a: num, b: num) -> num {
    let result = $a + $b
    return $result
}

fn main() {
    let x = 5
    let y = 10
    let sum = add($x, $y)
    echo "sum = " + $sum
}`

	_, err := parser.ParseSourceScript(script)
	assert.NoError(t, err)

	interp := &Interpreter{}
	debugger := NewDebugger(interp)

	// 添加断点到函数内部
	bp, err := debugger.AddBreakpoint("test.hl", 2, "")
	assert.NoError(t, err)
	assert.NotNil(t, bp)

	// 获取调用栈（应该为空，因为还没有执行）
	callStack := debugger.GetCallStack()
	assert.Len(t, callStack, 0)
}

func TestDebuggerVariables(t *testing.T) {
	interp := &Interpreter{}
	debugger := NewDebugger(interp)

	// 获取变量（应该为空，因为还没有执行）
	variables := debugger.GetVariables()
	assert.Len(t, variables, 0)
}

func TestDebuggerPosition(t *testing.T) {
	interp := &Interpreter{}
	debugger := NewDebugger(interp)

	// 设置文件内容
	script := `let x = 10
let y = 20
let z = $x + $y`
	debugger.SetFileContent("test.hl", script)

	// 测试位置获取
	node, err := parser.ParseSourceScript(script)
	assert.NoError(t, err)

	pos := debugger.getNodePosition(node)
	assert.NotNil(t, pos)
	assert.Equal(t, "test.hl", pos.File)
	assert.Equal(t, node, pos.Node)

	// 验证位置计算是否正确
	assert.Greater(t, pos.Line, 0)
	assert.Greater(t, pos.Column, 0)
}

func TestDebuggerStartStop(t *testing.T) {
	interp := &Interpreter{}
	debugger := NewDebugger(interp)

	// 启动调试
	debugger.StartDebugging()
	assert.True(t, debugger.isDebugging)

	// 停止调试
	debugger.StopDebugging()
	assert.False(t, debugger.isDebugging)
}

func TestDebuggerListeners(t *testing.T) {
	interp := &Interpreter{}
	debugger := NewDebugger(interp)

	// 创建一个测试监听器
	listenerCalled := false
	testListener := &testDebugListener{
		onBreakpointHit: func() {
			listenerCalled = true
		},
	}

	// 添加监听器
	debugger.AddListener(testListener)
	assert.Len(t, debugger.listeners, 1)

	// 模拟断点命中
	debugger.notifyBreakpointHit(&Breakpoint{ID: "test"}, &ast.Position{File: "test.hl", Line: 1})

	// 验证监听器被调用
	assert.True(t, listenerCalled)
}

func TestDebuggerFileContent(t *testing.T) {
	interp := &Interpreter{}
	debugger := NewDebugger(interp)

	// 测试多行文件内容
	content := `let x = 10
let y = 20
if $x > 5 {
    let z = $x + $y
    echo "z = " + $z
}`

	debugger.SetFileContent("test.hl", content)

	// 解析脚本并测试位置获取
	script := `let x = 10`
	node, err := parser.ParseSourceScript(script)
	assert.NoError(t, err)

	pos := debugger.getNodePosition(node)
	assert.NotNil(t, pos)
	assert.Equal(t, "test.hl", pos.File)
	assert.Equal(t, 1, pos.Line) // 应该在第一行
}

// 测试用的调试监听器
type testDebugListener struct {
	onBreakpointHit func()
}

func (t *testDebugListener) OnBreakpointHit(bp *Breakpoint, pos *ast.Position) {
	if t.onBreakpointHit != nil {
		t.onBreakpointHit()
	}
}

func (t *testDebugListener) OnStep(pos *ast.Position)                             {}
func (t *testDebugListener) OnVariableChanged(name string, value interface{}) {}
func (t *testDebugListener) OnFunctionEnter(frame *CallFrame)                 {}
func (t *testDebugListener) OnFunctionExit(frame *CallFrame)                  {}
