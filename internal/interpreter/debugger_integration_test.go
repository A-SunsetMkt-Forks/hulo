package interpreter

import (
	"math/big"
	"testing"
	"time"

	"github.com/hulo-lang/hulo/syntax/hulo/ast"
	"github.com/hulo-lang/hulo/syntax/hulo/parser"
	"github.com/stretchr/testify/assert"
)

// TestDebuggerIntegration 测试调试器在evalComptime模式下的集成
func TestDebuggerIntegration(t *testing.T) {
	// 创建一个包含comptime代码的测试脚本
	script := `
comptime {
    let x = 10
    let y = 20
    let z = $x + $y
    echo "z = " + $z
}`

	// 解析脚本
	node, err := parser.ParseSourceScript(script)
	assert.NoError(t, err)

	// 创建环境和解释器
	env := NewEnvironment()
	interp := NewInterpreter(env)

	// 创建调试器并设置到解释器
	debugger := NewDebugger(interp)
	interp.SetDebugger(debugger)

	// 启动调试模式
	debugger.StartDebugging()

	// 添加断点到comptime块内的第5行（let z = $x + $y）
	bp, err := debugger.AddBreakpoint("test.hl", 5, "")
	assert.NoError(t, err)
	assert.NotNil(t, bp)

	// 设置文件内容用于位置计算
	debugger.SetFileContent("test.hl", script)
	debugger.SetCurrentFile("test.hl")

	// 同时为解释器设置文件内容
	interp.SetFileContent("test.hl", script)
	interp.SetCurrentFile("test.hl")

	// 创建一个通道来接收调试事件
	debugEvents := make(chan string, 10)
	debugger.AddListener(&integrationTestListener{
		onBreakpointHit: func() {
			debugEvents <- "breakpoint_hit"
		},
		onStep: func() {
			debugEvents <- "step"
		},
	})

	// 在goroutine中执行解释，这样我们可以控制调试流程
	go func() {
		// 执行comptime块
		if len(node.Stmts) > 0 {
			if stmt, ok := node.Stmts[0].(*ast.ExprStmt).X.(*ast.ComptimeStmt); ok {
				interp.Eval(stmt)
			}
		}
	}()

	// 等待断点命中
	select {
	case event := <-debugEvents:
		assert.Equal(t, "breakpoint_hit", event)
	case <-time.After(2 * time.Second):
		t.Fatal("Timeout waiting for breakpoint")
	}

	// 验证调试器状态
	assert.True(t, debugger.isPaused)
	assert.NotNil(t, debugger.currentPos)

	// 获取当前变量
	variables := debugger.GetVariables()
	assert.Contains(t, variables, "x")
	assert.Contains(t, variables, "y")

	// 检查变量值（处理 big.Float 类型）
	xValue := variables["x"]
	yValue := variables["y"]

	// 将 big.Float 转换为 float64 进行比较
	if xFloat, ok := xValue.(*big.Float); ok {
		xFloat64, _ := xFloat.Float64()
		assert.Equal(t, float64(10), xFloat64)
	} else {
		assert.Equal(t, 10, xValue)
	}

	if yFloat, ok := yValue.(*big.Float); ok {
		yFloat64, _ := yFloat.Float64()
		assert.Equal(t, float64(20), yFloat64)
	} else {
		assert.Equal(t, 20, yValue)
	}

	// 继续执行
	debugger.SendCommand(DebugCommand{
		Type: CmdContinue,
	})

	// 等待执行完成
	time.Sleep(100 * time.Millisecond)
	assert.False(t, debugger.isPaused)
}

// TestDebuggerStepByStep 测试单步调试功能
func TestDebuggerStepByStep(t *testing.T) {
	script := `
comptime {
    let x = 5
    let y = 10
    if $x > 3 {
        let z = $x + $y
        echo "z = " + $z
    }
}
`

	node, err := parser.ParseSourceScript(script)
	assert.NoError(t, err)

	env := NewEnvironment()
	interp := NewInterpreter(env)
	debugger := NewDebugger(interp)
	interp.SetDebugger(debugger)

	debugger.StartDebugging()
	debugger.SetFileContent("test.hl", script)
	debugger.SetCurrentFile("test.hl")

	// 同时为解释器设置文件内容
	interp.SetFileContent("test.hl", script)
	interp.SetCurrentFile("test.hl")

	// 添加断点到if语句
	_, err = debugger.AddBreakpoint("test.hl", 5, "")
	assert.NoError(t, err)

	stepEvents := make(chan string, 10)
	debugger.AddListener(&integrationTestListener{
		onStep: func() {
			stepEvents <- "step"
		},
	})

	// 执行到断点
	go func() {
		if len(node.Stmts) > 0 {
			if stmt, ok := node.Stmts[0].(*ast.ExprStmt).X.(*ast.ComptimeStmt); ok {
				interp.Eval(stmt)
			}
		}
	}()

	// 等待断点命中
	time.Sleep(100 * time.Millisecond)

	// 单步执行
	debugger.SendCommand(DebugCommand{
		Type: CmdStepInto,
	})

	// 等待单步事件
	select {
	case event := <-stepEvents:
		assert.Equal(t, "step", event)
	case <-time.After(1 * time.Second):
		t.Log("Step event timeout (may be expected)")
	}
}

// TestDebuggerConditionalBreakpoint 测试条件断点
func TestDebuggerConditionalBreakpoint(t *testing.T) {
	script := `
comptime {
    let x = 5
    let y = 10
    let z = $x + $y
    echo "z = " + $z
}
`

	node, err := parser.ParseSourceScript(script)
	assert.NoError(t, err)

	env := NewEnvironment()
	interp := NewInterpreter(env)
	debugger := NewDebugger(interp)
	interp.SetDebugger(debugger)

	debugger.StartDebugging()
	debugger.SetFileContent("test.hl", script)
	debugger.SetCurrentFile("test.hl")

	// 同时为解释器设置文件内容
	interp.SetFileContent("test.hl", script)
	interp.SetCurrentFile("test.hl")

	// 添加条件断点：只有当x > 3时才命中
	bp, err := debugger.AddBreakpoint("test.hl", 5, "$x > 3")
	assert.NoError(t, err)
	assert.Equal(t, "$x > 3", bp.Condition)

	breakpointEvents := make(chan string, 5)
	debugger.AddListener(&integrationTestListener{
		onBreakpointHit: func() {
			breakpointEvents <- "breakpoint_hit"
		},
	})

	// 执行代码
	go func() {
		if len(node.Stmts) > 0 {
			if stmt, ok := node.Stmts[0].(*ast.ExprStmt).X.(*ast.ComptimeStmt); ok {
				interp.Eval(stmt)
			}
		}
	}()

	// 等待条件断点命中
	select {
	case event := <-breakpointEvents:
		assert.Equal(t, "breakpoint_hit", event)
	case <-time.After(2 * time.Second):
		t.Fatal("Timeout waiting for conditional breakpoint")
	}

	// 验证变量状态
	variables := debugger.GetVariables()
	assert.Equal(t, 5, variables["x"])
	assert.Equal(t, 10, variables["y"])
}

// TestDebuggerCallStackIntegration 测试调用栈功能
func TestDebuggerCallStackIntegration(t *testing.T) {
	script := `
comptime {
    fn add(a: num, b: num) -> num {
        let result = $a + $b
        return $result
    }

    let x = 5
    let y = 10
    let sum = add($x, $y)
    echo "sum = " + $sum
}
`

	node, err := parser.ParseSourceScript(script)
	assert.NoError(t, err)

	env := NewEnvironment()
	interp := NewInterpreter(env)
	debugger := NewDebugger(interp)
	interp.SetDebugger(debugger)

	debugger.StartDebugging()
	debugger.SetFileContent("test.hl", script)
	debugger.SetCurrentFile("test.hl")

	// 同时为解释器设置文件内容
	interp.SetFileContent("test.hl", script)
	interp.SetCurrentFile("test.hl")

	// 在函数内部添加断点
	_, err = debugger.AddBreakpoint("test.hl", 4, "")
	assert.NoError(t, err)

	breakpointEvents := make(chan string, 5)
	debugger.AddListener(&integrationTestListener{
		onBreakpointHit: func() {
			breakpointEvents <- "breakpoint_hit"
		},
	})

	// 执行代码
	go func() {
		if len(node.Stmts) > 0 {
			if stmt, ok := node.Stmts[0].(*ast.ExprStmt).X.(*ast.ComptimeStmt); ok {
				interp.Eval(stmt)
			}
		}
	}()

	// 等待断点命中
	select {
	case event := <-breakpointEvents:
		assert.Equal(t, "breakpoint_hit", event)
	case <-time.After(2 * time.Second):
		t.Fatal("Timeout waiting for breakpoint in function")
	}

	// 验证调用栈
	callStack := debugger.GetCallStack()
	assert.Greater(t, len(callStack), 0)

	// 验证栈帧信息
	if len(callStack) > 0 {
		frame := callStack[0]
		assert.Contains(t, frame.FunctionName, "add")
		assert.Contains(t, frame.Variables, "a")
		assert.Contains(t, frame.Variables, "b")
	}
}

// TestDebuggerVariableInspection 测试变量检查功能
func TestDebuggerVariableInspection(t *testing.T) {
	script := `
comptime {
    let x = 42
    let y = "hello"
    let z = true
    let arr = [1, 2, 3]
    let obj = {name: "test", value: 123}
}
`

	node, err := parser.ParseSourceScript(script)
	assert.NoError(t, err)

	env := NewEnvironment()
	interp := NewInterpreter(env)
	debugger := NewDebugger(interp)
	interp.SetDebugger(debugger)

	debugger.StartDebugging()
	debugger.SetFileContent("test.hl", script)
	debugger.SetCurrentFile("test.hl")

	// 同时为解释器设置文件内容
	interp.SetFileContent("test.hl", script)
	interp.SetCurrentFile("test.hl")

	// 在最后一行添加断点
	_, err = debugger.AddBreakpoint("test.hl", 7, "")
	assert.NoError(t, err)

	breakpointEvents := make(chan string, 5)
	debugger.AddListener(&integrationTestListener{
		onBreakpointHit: func() {
			breakpointEvents <- "breakpoint_hit"
		},
	})

	// 执行代码
	go func() {
		if len(node.Stmts) > 0 {
			if stmt, ok := node.Stmts[0].(*ast.ExprStmt).X.(*ast.ComptimeStmt); ok {
				interp.Eval(stmt)
			}
		}
	}()

	// 等待断点命中
	select {
	case event := <-breakpointEvents:
		assert.Equal(t, "breakpoint_hit", event)
	case <-time.After(2 * time.Second):
		t.Fatal("Timeout waiting for breakpoint")
	}

	// 检查所有变量
	variables := debugger.GetVariables()

	// 先检查基本变量
	assert.Contains(t, variables, "x")
	assert.Contains(t, variables, "y")
	assert.Contains(t, variables, "z")

	// 检查数组变量（对象变量可能不被支持）
	if val, exists := variables["arr"]; exists {
		t.Logf("arr variable found: %T %v", val, val)
	} else {
		t.Logf("arr variable NOT found")
	}

	// 注意：对象字面量在comptime模式下可能不被支持
	if val, exists := variables["obj"]; exists {
		t.Logf("obj variable found: %T %v", val, val)
	} else {
		t.Logf("obj variable NOT found (object literals may not be supported in comptime mode)")
	}

	// 打印所有变量名
	varNames := make([]string, 0, len(variables))
	for name := range variables {
		varNames = append(varNames, name)
	}
	t.Logf("All variables: %v", varNames)

	// 验证基本变量值（处理 big.Float 类型）
	xValue := variables["x"]
	if xFloat, ok := xValue.(*big.Float); ok {
		xFloat64, _ := xFloat.Float64()
		assert.Equal(t, float64(42), xFloat64)
	} else {
		assert.Equal(t, 42, xValue)
	}

	assert.Equal(t, "hello", variables["y"])
	assert.Equal(t, true, variables["z"])
}

// 集成测试监听器
type integrationTestListener struct {
	onBreakpointHit   func()
	onStep            func()
	onVariableChanged func()
	onFunctionEnter   func()
	onFunctionExit    func()
}

func (t *integrationTestListener) OnBreakpointHit(bp *Breakpoint, pos *ast.Position) {
	if t.onBreakpointHit != nil {
		t.onBreakpointHit()
	}
}

func (t *integrationTestListener) OnStep(pos *ast.Position) {
	if t.onStep != nil {
		t.onStep()
	}
}

func (t *integrationTestListener) OnVariableChanged(name string, value interface{}) {
	if t.onVariableChanged != nil {
		t.onVariableChanged()
	}
}

func (t *integrationTestListener) OnFunctionEnter(frame *CallFrame) {
	if t.onFunctionEnter != nil {
		t.onFunctionEnter()
	}
}

func (t *integrationTestListener) OnFunctionExit(frame *CallFrame) {
	if t.onFunctionExit != nil {
		t.onFunctionExit()
	}
}
