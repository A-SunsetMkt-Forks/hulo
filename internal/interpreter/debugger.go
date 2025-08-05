// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package interpreter

import (
	"context"
	"fmt"
	"sync"

	"github.com/caarlos0/log"
	"github.com/hulo-lang/hulo/syntax/hulo/ast"
	"github.com/hulo-lang/hulo/syntax/hulo/token"
)

type Debugger struct {
	interpreter *Interpreter

	// 断点管理
	breakpoints     map[string]map[int]*Breakpoint // file -> line -> breakpoint
	breakpointMutex sync.RWMutex

	// 调试状态
	isDebugging bool
	isPaused    bool
	currentPos  *ast.Position

	// 通信通道
	pauseChan   chan struct{}
	resumeChan  chan DebugCommand
	commandChan chan DebugCommand

	// 调试信息
	callStack []*CallFrame
	variables map[string]interface{}

	// 监听器
	listeners []DebugListener

	// 文件位置信息
	fileSet *ast.FileSet

	// 上下文
	ctx    context.Context
	cancel context.CancelFunc
}

// CallFrame 调用栈帧
type CallFrame struct {
	FunctionName string
	File         string
	Line         int
	Variables    map[string]interface{}
}

// Breakpoint 断点
type Breakpoint struct {
	ID          string
	File        string
	Line        int
	Column      int
	Condition   string // 条件断点表达式
	HitCount    int    // 命中次数
	MaxHitCount int    // 最大命中次数
	Enabled     bool
	Temporary   bool // 临时断点
}

// DebugCommand 调试命令
type DebugCommand struct {
	Type     DebugCommandType
	Data     interface{}
	Response chan interface{}
}

type DebugCommandType int

const (
	CmdContinue DebugCommandType = iota
	CmdStepInto
	CmdStepOver
	CmdStepOut
	CmdPause
	CmdResume
	CmdAddBreakpoint
	CmdRemoveBreakpoint
	CmdGetVariables
	CmdGetCallStack
	CmdEvaluate
)

type DebugListener interface {
	OnBreakpointHit(bp *Breakpoint, pos *ast.Position)
	OnStep(pos *ast.Position)
	OnVariableChanged(name string, value interface{})
	OnFunctionEnter(frame *CallFrame)
	OnFunctionExit(frame *CallFrame)
}

func NewDebugger(interpreter *Interpreter) *Debugger {
	ctx, cancel := context.WithCancel(context.Background())

	return &Debugger{
		interpreter: interpreter,
		breakpoints: make(map[string]map[int]*Breakpoint),
		callStack:   make([]*CallFrame, 0),
		variables:   make(map[string]any),
		listeners:   make([]DebugListener, 0),
		fileSet:     ast.NewFileSet(),
		pauseChan:   make(chan struct{}, 1),
		resumeChan:  make(chan DebugCommand, 1),
		commandChan: make(chan DebugCommand, 10),
		ctx:         ctx,
		cancel:      cancel,
	}
}

func (d *Debugger) AddBreakpoint(file string, line int, condition string) (*Breakpoint, error) {
	d.breakpointMutex.Lock()
	defer d.breakpointMutex.Unlock()

	if d.breakpoints[file] == nil {
		d.breakpoints[file] = make(map[int]*Breakpoint)
	}

	bp := &Breakpoint{
		ID:        fmt.Sprintf("%s:%d", file, line),
		File:      file,
		Line:      line,
		Condition: condition,
		Enabled:   true,
	}

	d.breakpoints[file][line] = bp
	return bp, nil
}

func (d *Debugger) RemoveBreakpoint(file string, line int) error {
	d.breakpointMutex.Lock()
	defer d.breakpointMutex.Unlock()

	if fileBreakpoints, exists := d.breakpoints[file]; exists {
		delete(fileBreakpoints, line)
	}

	return nil
}

func (d *Debugger) GetBreakpoints() []*Breakpoint {
	d.breakpointMutex.RLock()
	defer d.breakpointMutex.RUnlock()

	var result []*Breakpoint
	for _, fileBreakpoints := range d.breakpoints {
		for _, bp := range fileBreakpoints {
			result = append(result, bp)
		}
	}
	return result
}

type DebugState struct{}

func (d *Debugger) shouldBreak(node ast.Node) bool {
	if !d.isDebugging {
		return false
	}

	pos := d.getNodePosition(node)
	if pos == nil {
		return false
	}

	d.breakpointMutex.RLock()
	fileBreakpoints, exists := d.breakpoints[pos.File]
	d.breakpointMutex.RUnlock()

	if !exists {
		return false
	}

	bp, exists := fileBreakpoints[pos.Line]
	if !exists || !bp.Enabled {
		return false
	}

	// 检查命中次数
	if bp.MaxHitCount > 0 && bp.HitCount >= bp.MaxHitCount {
		return false
	}

	// 检查条件断点
	if bp.Condition != "" {
		if !d.evaluateCondition(bp.Condition, node) {
			return false
		}
	}

	// 更新命中次数
	bp.HitCount++

	// 捕获当前解释器状态
	d.currentPos = pos
	d.variables = d.interpreter.GetCurrentEnvironment().GetAll()
	d.callStack = d.interpreter.GetInterpreterCallStack()

	// 通知监听器
	d.notifyBreakpointHit(bp, pos)

	// 暂停执行
	d.pause()

	return true
}

func (d *Debugger) getNodePosition(node ast.Node) *ast.Position {
	if node == nil {
		return nil
	}

	return &ast.Position{
		File:   d.getCurrentFile(),
		Line:   d.getLineFromPos(node.Pos()),
		Column: d.getColumnFromPos(node.Pos()),
		Node:   node,
	}
}

func (d *Debugger) evaluateCondition(condition string, node ast.Node) bool {
	// 这里需要实现条件表达式的求值
	// 可以使用你的解释器来求值条件表达式
	// 暂时返回 true
	return true
}

func (d *Debugger) StartDebugging() {
	d.isDebugging = true
	d.isPaused = false

	// 启动调试循环
	go d.debugLoop()
}

func (d *Debugger) StopDebugging() {
	d.isDebugging = false
	d.isPaused = false
	d.cancel()
}

func (d *Debugger) debugLoop() {
	for {
		select {
		case <-d.ctx.Done():
			return
		case cmd := <-d.commandChan:
			d.handleCommand(cmd)
		case <-d.pauseChan:
			d.waitForResume()
		}
	}
}

func (d *Debugger) waitForResume() {
	d.isPaused = true

	// 通知所有监听器
	for _, listener := range d.listeners {
		listener.OnBreakpointHit(nil, d.currentPos)
	}

	// 等待恢复命令
	select {
	case cmd := <-d.resumeChan:
		d.isPaused = false
		d.handleCommand(cmd)
	case <-d.ctx.Done():
		return
	}
}

func (d *Debugger) handleCommand(cmd DebugCommand) {
	switch cmd.Type {
	case CmdContinue:
		d.continue_()
	case CmdStepInto:
		d.stepInto()
	case CmdStepOver:
		d.stepOver()
	case CmdStepOut:
		d.stepOut()
	case CmdPause:
		d.pause()
	case CmdResume:
		d.resume()
	case CmdAddBreakpoint:
		d.handleAddBreakpoint(cmd)
	case CmdRemoveBreakpoint:
		d.handleRemoveBreakpoint(cmd)
	case CmdGetVariables:
		d.handleGetVariables(cmd)
	case CmdGetCallStack:
		d.handleGetCallStack(cmd)
	case CmdEvaluate:
		d.handleEvaluate(cmd)
	}

	// 发送响应
	if cmd.Response != nil {
		cmd.Response <- nil
	}
}

func (d *Debugger) pause() {
	d.isPaused = true
	d.notifyClient(&DebugState{})

	// 等待恢复命令
	cmd := <-d.resumeChan
	// 处理恢复命令
	d.handleCommand(cmd)
}

func (d *Debugger) notifyClient(*DebugState) {}

func (d *Debugger) stepInto() {}
func (d *Debugger) stepOver() {}
func (d *Debugger) stepOut()  {}
func (d *Debugger) continue_() {
	// 发送恢复命令到 resumeChan
	select {
	case d.resumeChan <- DebugCommand{Type: CmdContinue}:
	default:
	}
}
func (d *Debugger) inspectVariable() {}

func (d *Debugger) AddListener(listener DebugListener) {
	d.listeners = append(d.listeners, listener)
}

func (d *Debugger) SendCommand(cmd DebugCommand) {
	d.commandChan <- cmd
}

func (d *Debugger) GetCallStack() []*CallFrame {
	return d.callStack
}

func (d *Debugger) GetVariables() map[string]interface{} {
	return d.variables
}

func (d *Debugger) getCurrentFile() string {
	if len(d.callStack) > 0 {
		return d.callStack[len(d.callStack)-1].File
	}
	// 如果没有调用栈，返回第一个断点的文件名
	d.breakpointMutex.RLock()
	defer d.breakpointMutex.RUnlock()

	for filename := range d.breakpoints {
		return filename
	}
	return "unknown"
}

// SetFileContent 设置文件内容，用于位置计算
func (d *Debugger) SetFileContent(filename, content string) {
	d.fileSet.AddFile(filename, content)
}

// SetCurrentFile 设置当前文件名
func (d *Debugger) SetCurrentFile(filename string) {
	// 创建一个默认的调用栈帧来设置当前文件
	if len(d.callStack) == 0 {
		frame := &CallFrame{
			FunctionName: "main",
			File:         filename,
			Line:         1,
			Variables:    make(map[string]interface{}),
		}
		d.callStack = append(d.callStack, frame)
	} else {
		d.callStack[0].File = filename
	}
}

func (d *Debugger) getLineFromPos(pos token.Pos) int {
	if !pos.IsValid() {
		return 1
	}

	file := d.getCurrentFile()
	fi := d.fileSet.GetFile(file)
	if fi == nil {
		return 1 // 如果文件不存在，返回默认值
	}

	line, _ := fi.PosToLineColumn(pos)
	return line
}

func (d *Debugger) getColumnFromPos(pos token.Pos) int {
	if !pos.IsValid() {
		return 1
	}

	file := d.getCurrentFile()
	fi := d.fileSet.GetFile(file)
	if fi == nil {
		return 1 // 如果文件不存在，返回默认值
	}

	_, column := fi.PosToLineColumn(pos)
	return column
}

func (d *Debugger) notifyBreakpointHit(bp *Breakpoint, pos *ast.Position) {
	for _, listener := range d.listeners {
		listener.OnBreakpointHit(bp, pos)
	}
}

func (d *Debugger) handleAddBreakpoint(cmd DebugCommand) {
	if data, ok := cmd.Data.(map[string]interface{}); ok {
		file := data["file"].(string)
		line := int(data["line"].(float64))
		condition := ""
		if cond, exists := data["condition"]; exists {
			condition = cond.(string)
		}

		_, err := d.AddBreakpoint(file, line, condition)
		if err != nil {
			log.WithError(err).Error("Failed to add breakpoint")
		} else {
			log.Infof("Added breakpoint at %s:%d", file, line)
		}
	}
}

func (d *Debugger) handleRemoveBreakpoint(cmd DebugCommand) {
	if data, ok := cmd.Data.(map[string]interface{}); ok {
		file := data["file"].(string)
		line := int(data["line"].(float64))

		err := d.RemoveBreakpoint(file, line)
		if err != nil {
			log.WithError(err).Error("Failed to remove breakpoint")
		} else {
			log.Infof("Removed breakpoint at %s:%d", file, line)
		}
	}
}

func (d *Debugger) handleGetVariables(cmd DebugCommand) {
	variables := d.GetVariables()
	if cmd.Response != nil {
		cmd.Response <- variables
	}
}

func (d *Debugger) handleGetCallStack(cmd DebugCommand) {
	callStack := d.GetCallStack()
	if cmd.Response != nil {
		cmd.Response <- callStack
	}
}

func (d *Debugger) handleEvaluate(cmd DebugCommand) {
	if expression, ok := cmd.Data.(string); ok {
		// 这里需要实现表达式求值
		// 暂时返回表达式本身
		if cmd.Response != nil {
			cmd.Response <- expression
		}
	}
}

func (d *Debugger) resume() {
	d.isPaused = false
}
