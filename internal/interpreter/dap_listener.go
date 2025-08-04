// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package interpreter

import (
	"fmt"

	"github.com/hulo-lang/hulo/syntax/hulo/ast"
)

// DAPDebugListener DAP 调试监听器
type DAPDebugListener struct {
	adapter *DebugAdapter
	conn    *DAPConnection
}

// NewDAPDebugListener 创建 DAP 调试监听器
func NewDAPDebugListener(adapter *DebugAdapter) *DAPDebugListener {
	return &DAPDebugListener{
		adapter: adapter,
	}
}

// OnBreakpointHit 断点命中事件
func (d *DAPDebugListener) OnBreakpointHit(bp *Breakpoint, pos *ast.Position) {
	// 广播断点命中事件到所有连接
	d.adapter.broadcastEvent("breakpoint", map[string]any{
		"reason":     "changed",
		"breakpoint": bp,
	})

	// 广播停止事件到所有连接
	d.adapter.broadcastEvent("stopped", map[string]any{
		"reason":            "breakpoint",
		"description":       fmt.Sprintf("Breakpoint hit at %s:%d", pos.File, pos.Line),
		"threadId":          1,
		"allThreadsStopped": true,
	})
}

// OnStep 单步执行事件
func (d *DAPDebugListener) OnStep(pos *ast.Position) {
	// 广播停止事件到所有连接
	d.adapter.broadcastEvent("stopped", map[string]any{
		"reason":            "step",
		"description":       fmt.Sprintf("Stepped to %s:%d", pos.File, pos.Line),
		"threadId":          1,
		"allThreadsStopped": true,
	})
}

// OnVariableChanged 变量变化事件
func (d *DAPDebugListener) OnVariableChanged(name string, value any) {
	// 广播输出事件到所有连接
	d.adapter.broadcastEvent("output", map[string]any{
		"category": "console",
		"output":   fmt.Sprintf("Variable %s = %v", name, value),
	})
}

// OnFunctionEnter 函数进入事件
func (d *DAPDebugListener) OnFunctionEnter(frame *CallFrame) {
	// 广播输出事件到所有连接
	d.adapter.broadcastEvent("output", map[string]any{
		"category": "console",
		"output":   fmt.Sprintf("Entering function: %s", frame.FunctionName),
	})
}

// OnFunctionExit 函数退出事件
func (d *DAPDebugListener) OnFunctionExit(frame *CallFrame) {
	// 广播输出事件到所有连接
	d.adapter.broadcastEvent("output", map[string]any{
		"category": "console",
		"output":   fmt.Sprintf("Exiting function: %s", frame.FunctionName),
	})
}
