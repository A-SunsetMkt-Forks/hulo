// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package interpreter

import (
	"encoding/json"
	"fmt"
)

// DAP 请求和响应结构体

// InitializeRequest 初始化请求
type InitializeRequest struct {
	ClientID                     string `json:"clientID,omitempty"`
	ClientName                   string `json:"clientName,omitempty"`
	AdapterID                    string `json:"adapterID,omitempty"`
	LinesStartAt1                bool   `json:"linesStartAt1,omitempty"`
	ColumnsStartAt1              bool   `json:"columnsStartAt1,omitempty"`
	PathFormat                   string `json:"pathFormat,omitempty"`
	SupportsVariableType         bool   `json:"supportsVariableType,omitempty"`
	SupportsVariablePaging       bool   `json:"supportsVariablePaging,omitempty"`
	SupportsRunInTerminalRequest bool   `json:"supportsRunInTerminalRequest,omitempty"`
	Locale                       string `json:"locale,omitempty"`
}

// InitializeResponse 初始化响应
type InitializeResponse struct {
	SupportsConfigurationDoneRequest      bool `json:"supportsConfigurationDoneRequest"`
	SupportsFunctionBreakpoints           bool `json:"supportsFunctionBreakpoints"`
	SupportsConditionalBreakpoints        bool `json:"supportsConditionalBreakpoints"`
	SupportsHitConditionalBreakpoints     bool `json:"supportsHitConditionalBreakpoints"`
	SupportsEvaluateForHovers             bool `json:"supportsEvaluateForHovers"`
	SupportsStepBack                      bool `json:"supportsStepBack"`
	SupportsSetVariable                   bool `json:"supportsSetVariable"`
	SupportsRestartFrame                  bool `json:"supportsRestartFrame"`
	SupportsGotoTargetsRequest            bool `json:"supportsGotoTargetsRequest"`
	SupportsStepInTargetsRequest          bool `json:"supportsStepInTargetsRequest"`
	SupportsCompletionsRequest            bool `json:"supportsCompletionsRequest"`
	SupportsModulesRequest                bool `json:"supportsModulesRequest"`
	SupportsRestartRequest                bool `json:"supportsRestartRequest"`
	SupportsExceptionOptions              bool `json:"supportsExceptionOptions"`
	SupportsValueFormattingOptions        bool `json:"supportsValueFormattingOptions"`
	SupportsExceptionInfoRequest          bool `json:"supportsExceptionInfoRequest"`
	SupportsTerminateDebuggee             bool `json:"supportsTerminateDebuggee"`
	SupportsDelayedStackTraceLoading      bool `json:"supportsDelayedStackTraceLoading"`
	SupportsLoadedSourcesRequest          bool `json:"supportsLoadedSourcesRequest"`
	SupportsLogPoints                     bool `json:"supportsLogPoints"`
	SupportsTerminateThreadsRequest       bool `json:"supportsTerminateThreadsRequest"`
	SupportsSetExpression                 bool `json:"supportsSetExpression"`
	SupportsTerminateRequest              bool `json:"supportsTerminateRequest"`
	SupportsDataBreakpoints               bool `json:"supportsDataBreakpoints"`
	SupportsReadMemoryRequest             bool `json:"supportsReadMemoryRequest"`
	SupportsWriteMemoryRequest            bool `json:"supportsWriteMemoryRequest"`
	SupportsDisassembleRequest            bool `json:"supportsDisassembleRequest"`
	SupportsCancelRequest                 bool `json:"supportsCancelRequest"`
	SupportsBreakpointLocationsRequest    bool `json:"supportsBreakpointLocationsRequest"`
	SupportsClipboardContext              bool `json:"supportsClipboardContext"`
	SupportsSteppingGranularity           bool `json:"supportsSteppingGranularity"`
	SupportsInstructionBreakpoints        bool `json:"supportsInstructionBreakpoints"`
	SupportsExceptionFilterOptions        bool `json:"supportsExceptionFilterOptions"`
	SupportsSingleThreadExecutionRequests bool `json:"supportsSingleThreadExecutionRequests"`
}

// LaunchRequest 启动请求
type LaunchRequest struct {
	Program        string            `json:"program"`
	Args           []string          `json:"args,omitempty"`
	Cwd            string            `json:"cwd,omitempty"`
	Env            map[string]string `json:"env,omitempty"`
	Console        string            `json:"console,omitempty"`
	StopOnEntry    bool              `json:"stopOnEntry,omitempty"`
	DebugServer    int               `json:"debugServer,omitempty"`
	AdditionalData json.RawMessage   `json:"__additionalData,omitempty"`
}

// SetBreakpointsRequest 设置断点请求
type SetBreakpointsRequest struct {
	Source         Source             `json:"source"`
	Breakpoints    []SourceBreakpoint `json:"breakpoints,omitempty"`
	Lines          []int              `json:"lines,omitempty"`
	SourceModified bool               `json:"sourceModified,omitempty"`
}

// Source 源文件信息
type Source struct {
	Name             string          `json:"name,omitempty"`
	Path             string          `json:"path,omitempty"`
	SourceReference  int             `json:"sourceReference,omitempty"`
	PresentationHint string          `json:"presentationHint,omitempty"`
	Origin           string          `json:"origin,omitempty"`
	AdapterData      json.RawMessage `json:"adapterData,omitempty"`
	Checksums        []Checksum      `json:"checksums,omitempty"`
}

// SourceBreakpoint 源断点
type SourceBreakpoint struct {
	Line         int    `json:"line"`
	Column       int    `json:"column,omitempty"`
	Condition    string `json:"condition,omitempty"`
	HitCondition string `json:"hitCondition,omitempty"`
	LogMessage   string `json:"logMessage,omitempty"`
}

// Checksum 校验和
type Checksum struct {
	Algorithm string `json:"algorithm"`
	Checksum  string `json:"checksum"`
}

// SetBreakpointsResponse 设置断点响应
type SetBreakpointsResponse struct {
	Breakpoints []Breakpoint `json:"breakpoints"`
}

// DAPBreakpoint DAP 断点
type DAPBreakpoint struct {
	ID                   int    `json:"id"`
	Verified             bool   `json:"verified"`
	Message              string `json:"message,omitempty"`
	Line                 int    `json:"line,omitempty"`
	Column               int    `json:"column,omitempty"`
	EndLine              int    `json:"endLine,omitempty"`
	EndColumn            int    `json:"endColumn,omitempty"`
	InstructionReference string `json:"instructionReference,omitempty"`
	Offset               int    `json:"offset,omitempty"`
}

// ContinueRequest 继续请求
type ContinueRequest struct {
	ThreadId int `json:"threadId"`
}

// ContinueResponse 继续响应
type ContinueResponse struct {
	AllThreadsContinued bool `json:"allThreadsContinued"`
}

// NextRequest 下一步请求
type NextRequest struct {
	ThreadId     int    `json:"threadId"`
	SingleThread bool   `json:"singleThread,omitempty"`
	Granularity  string `json:"granularity,omitempty"`
}

// StepInRequest 步入请求
type StepInRequest struct {
	ThreadId     int  `json:"threadId"`
	TargetId     int  `json:"targetId,omitempty"`
	SingleThread bool `json:"singleThread,omitempty"`
}

// StepOutRequest 步出请求
type StepOutRequest struct {
	ThreadId     int  `json:"threadId"`
	SingleThread bool `json:"singleThread,omitempty"`
}

// PauseRequest 暂停请求
type PauseRequest struct {
	ThreadId int `json:"threadId"`
}

// StackTraceRequest 栈跟踪请求
type StackTraceRequest struct {
	ThreadId   int          `json:"threadId"`
	StartFrame int          `json:"startFrame,omitempty"`
	Levels     int          `json:"levels,omitempty"`
	Format     *ValueFormat `json:"format,omitempty"`
}

// StackTraceResponse 栈跟踪响应
type StackTraceResponse struct {
	StackFrames []StackFrame `json:"stackFrames"`
	TotalFrames int          `json:"totalFrames,omitempty"`
}

// StackFrame 栈帧
type StackFrame struct {
	ID                          int         `json:"id"`
	Name                        string      `json:"name"`
	Source                      *Source     `json:"source,omitempty"`
	Line                        int         `json:"line"`
	Column                      int         `json:"column"`
	EndLine                     int         `json:"endLine,omitempty"`
	EndColumn                   int         `json:"endColumn,omitempty"`
	CanRestart                  bool        `json:"canRestart,omitempty"`
	InstructionPointerReference string      `json:"instructionPointerReference,omitempty"`
	ModuleId                    interface{} `json:"moduleId,omitempty"`
	PresentationHint            string      `json:"presentationHint,omitempty"`
}

// VariablesRequest 变量请求
type VariablesRequest struct {
	VariablesReference int          `json:"variablesReference"`
	Filter             string       `json:"filter,omitempty"`
	Start              int          `json:"start,omitempty"`
	Count              int          `json:"count,omitempty"`
	Format             *ValueFormat `json:"format,omitempty"`
}

// VariablesResponse 变量响应
type VariablesResponse struct {
	Variables []DAPVariable `json:"variables"`
}

// DAPVariable 变量
type DAPVariable struct {
	Name               string                    `json:"name"`
	Value              string                    `json:"value"`
	Type               string                    `json:"type,omitempty"`
	PresentationHint   *VariablePresentationHint `json:"presentationHint,omitempty"`
	EvaluateName       string                    `json:"evaluateName,omitempty"`
	VariablesReference int                       `json:"variablesReference"`
	NamedVariables     int                       `json:"namedVariables,omitempty"`
	IndexedVariables   int                       `json:"indexedVariables,omitempty"`
	MemoryReference    string                    `json:"memoryReference,omitempty"`
}

// VariablePresentationHint 变量表示提示
type VariablePresentationHint struct {
	Kind       string   `json:"kind,omitempty"`
	Attributes []string `json:"attributes,omitempty"`
	Visibility string   `json:"visibility,omitempty"`
	Lazy       bool     `json:"lazy,omitempty"`
}

// ValueFormat 值格式
type ValueFormat struct {
	Hex bool `json:"hex,omitempty"`
}

// EvaluateRequest 求值请求
type EvaluateRequest struct {
	Expression string       `json:"expression"`
	FrameId    int          `json:"frameId,omitempty"`
	Context    string       `json:"context,omitempty"`
	Format     *ValueFormat `json:"format,omitempty"`
}

// EvaluateResponse 求值响应
type EvaluateResponse struct {
	Result             string                    `json:"result"`
	Type               string                    `json:"type,omitempty"`
	PresentationHint   *VariablePresentationHint `json:"presentationHint,omitempty"`
	VariablesReference int                       `json:"variablesReference"`
	NamedVariables     int                       `json:"namedVariables,omitempty"`
	IndexedVariables   int                       `json:"indexedVariables,omitempty"`
	MemoryReference    string                    `json:"memoryReference,omitempty"`
}

// ThreadsResponse 线程响应
type ThreadsResponse struct {
	Threads []Thread `json:"threads"`
}

// Thread 线程
type Thread struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// DAP 请求处理器实现

// handleInitialize 处理初始化请求
func (da *DebugAdapter) handleInitialize(conn *DAPConnection, body json.RawMessage) (interface{}, error) {
	var req InitializeRequest
	if err := json.Unmarshal(body, &req); err != nil {
		return nil, err
	}

	response := &InitializeResponse{
		SupportsConfigurationDoneRequest:      true,
		SupportsFunctionBreakpoints:           true,
		SupportsConditionalBreakpoints:        true,
		SupportsHitConditionalBreakpoints:     true,
		SupportsEvaluateForHovers:             true,
		SupportsSetVariable:                   true,
		SupportsRestartFrame:                  true,
		SupportsCompletionsRequest:            true,
		SupportsExceptionOptions:              true,
		SupportsValueFormattingOptions:        true,
		SupportsExceptionInfoRequest:          true,
		SupportsTerminateDebuggee:             true,
		SupportsLoadedSourcesRequest:          true,
		SupportsLogPoints:                     true,
		SupportsTerminateThreadsRequest:       true,
		SupportsSetExpression:                 true,
		SupportsTerminateRequest:              true,
		SupportsDataBreakpoints:               true,
		SupportsCancelRequest:                 true,
		SupportsBreakpointLocationsRequest:    true,
		SupportsClipboardContext:              true,
		SupportsSteppingGranularity:           true,
		SupportsInstructionBreakpoints:        true,
		SupportsExceptionFilterOptions:        true,
		SupportsSingleThreadExecutionRequests: true,
	}

	return response, nil
}

// handleLaunch 处理启动请求
func (da *DebugAdapter) handleLaunch(conn *DAPConnection, body json.RawMessage) (interface{}, error) {
	var req LaunchRequest
	if err := json.Unmarshal(body, &req); err != nil {
		return nil, err
	}

	// 启动调试会话
	da.debugger.StartDebugging()

	// 如果设置了 stopOnEntry，发送停止事件
	if req.StopOnEntry {
		da.sendStoppedEvent(conn, "entry", "Stopped on entry")
	}

	return nil, nil
}

// handleAttach 处理附加请求
func (da *DebugAdapter) handleAttach(conn *DAPConnection, body json.RawMessage) (interface{}, error) {
	// 启动调试会话
	da.debugger.StartDebugging()
	return nil, nil
}

// handleDisconnect 处理断开连接请求
func (da *DebugAdapter) handleDisconnect(conn *DAPConnection, body json.RawMessage) (interface{}, error) {
	da.debugger.StopDebugging()
	return nil, nil
}

// handleTerminate 处理终止请求
func (da *DebugAdapter) handleTerminate(conn *DAPConnection, body json.RawMessage) (interface{}, error) {
	da.debugger.StopDebugging()
	return nil, nil
}

// handleSetBreakpoints 处理设置断点请求
func (da *DebugAdapter) handleSetBreakpoints(conn *DAPConnection, body json.RawMessage) (interface{}, error) {
	var req SetBreakpointsRequest
	if err := json.Unmarshal(body, &req); err != nil {
		return nil, err
	}

	var breakpoints []DAPBreakpoint

	for i, bp := range req.Breakpoints {
		// 添加断点到调试器
		debugBp, err := da.debugger.AddBreakpoint(req.Source.Path, bp.Line, bp.Condition)
		if err != nil {
			// 创建失败的断点
			breakpoints = append(breakpoints, DAPBreakpoint{
				ID:       i + 1,
				Line:     bp.Line,
				Column:   bp.Column,
				Verified: false,
				Message:  err.Error(),
			})
			continue
		}

		// 创建成功的断点
		breakpoints = append(breakpoints, DAPBreakpoint{
			ID:       i + 1,
			Line:     debugBp.Line,
			Column:   debugBp.Column,
			Verified: true,
		})
	}

	// Convert DAPBreakpoint to Breakpoint
	var convertedBreakpoints []Breakpoint
	for _, bp := range breakpoints {
		convertedBreakpoints = append(convertedBreakpoints, Breakpoint{
			ID:      fmt.Sprintf("%d", bp.ID),
			Line:    bp.Line,
			Column:  bp.Column,
			Enabled: bp.Verified,
		})
	}

	return &SetBreakpointsResponse{Breakpoints: convertedBreakpoints}, nil
}

// handleSetFunctionBreakpoints 处理函数断点请求
func (da *DebugAdapter) handleSetFunctionBreakpoints(conn *DAPConnection, body json.RawMessage) (interface{}, error) {
	// 暂时返回空响应
	return map[string]interface{}{"breakpoints": []interface{}{}}, nil
}

// handleSetExceptionBreakpoints 处理异常断点请求
func (da *DebugAdapter) handleSetExceptionBreakpoints(conn *DAPConnection, body json.RawMessage) (interface{}, error) {
	// 暂时返回空响应
	return map[string]interface{}{}, nil
}

// handleContinue 处理继续请求
func (da *DebugAdapter) handleContinue(conn *DAPConnection, body json.RawMessage) (interface{}, error) {
	var req ContinueRequest
	if err := json.Unmarshal(body, &req); err != nil {
		return nil, err
	}

	// 发送继续命令到调试器
	da.debugger.SendCommand(DebugCommand{
		Type: CmdContinue,
	})

	return &ContinueResponse{AllThreadsContinued: true}, nil
}

// handleNext 处理下一步请求
func (da *DebugAdapter) handleNext(conn *DAPConnection, body json.RawMessage) (interface{}, error) {
	var req NextRequest
	if err := json.Unmarshal(body, &req); err != nil {
		return nil, err
	}

	// 发送下一步命令到调试器
	da.debugger.SendCommand(DebugCommand{
		Type: CmdStepOver,
	})

	return nil, nil
}

// handleStepIn 处理步入请求
func (da *DebugAdapter) handleStepIn(conn *DAPConnection, body json.RawMessage) (interface{}, error) {
	var req StepInRequest
	if err := json.Unmarshal(body, &req); err != nil {
		return nil, err
	}

	// 发送步入命令到调试器
	da.debugger.SendCommand(DebugCommand{
		Type: CmdStepInto,
	})

	return nil, nil
}

// handleStepOut 处理步出请求
func (da *DebugAdapter) handleStepOut(conn *DAPConnection, body json.RawMessage) (interface{}, error) {
	var req StepOutRequest
	if err := json.Unmarshal(body, &req); err != nil {
		return nil, err
	}

	// 发送步出命令到调试器
	da.debugger.SendCommand(DebugCommand{
		Type: CmdStepOut,
	})

	return nil, nil
}

// handlePause 处理暂停请求
func (da *DebugAdapter) handlePause(conn *DAPConnection, body json.RawMessage) (interface{}, error) {
	var req PauseRequest
	if err := json.Unmarshal(body, &req); err != nil {
		return nil, err
	}

	// 发送暂停命令到调试器
	da.debugger.SendCommand(DebugCommand{
		Type: CmdPause,
	})

	return nil, nil
}

// handleStackTrace 处理栈跟踪请求
func (da *DebugAdapter) handleStackTrace(conn *DAPConnection, body json.RawMessage) (interface{}, error) {
	var req StackTraceRequest
	if err := json.Unmarshal(body, &req); err != nil {
		return nil, err
	}

	// 获取调用栈
	callStack := da.debugger.GetCallStack()

	var stackFrames []StackFrame
	for i, frame := range callStack {
		stackFrames = append(stackFrames, StackFrame{
			ID:     i,
			Name:   frame.FunctionName,
			Source: &Source{Path: frame.File},
			Line:   frame.Line,
			Column: 1,
		})
	}

	return &StackTraceResponse{
		StackFrames: stackFrames,
		TotalFrames: len(stackFrames),
	}, nil
}

// handleVariables 处理变量请求
func (da *DebugAdapter) handleVariables(conn *DAPConnection, body json.RawMessage) (interface{}, error) {
	var req VariablesRequest
	if err := json.Unmarshal(body, &req); err != nil {
		return nil, err
	}

	// 获取变量
	variables := da.debugger.GetVariables()

	var vars []DAPVariable
	for name, value := range variables {
		vars = append(vars, DAPVariable{
			Name:               name,
			Value:              fmt.Sprintf("%v", value),
			VariablesReference: 0,
		})
	}

	return &VariablesResponse{Variables: vars}, nil
}

// handleScopes 处理作用域请求
func (da *DebugAdapter) handleScopes(conn *DAPConnection, body json.RawMessage) (interface{}, error) {
	// 暂时返回空响应
	return map[string]interface{}{"scopes": []interface{}{}}, nil
}

// handleEvaluate 处理求值请求
func (da *DebugAdapter) handleEvaluate(conn *DAPConnection, body json.RawMessage) (interface{}, error) {
	var req EvaluateRequest
	if err := json.Unmarshal(body, &req); err != nil {
		return nil, err
	}

	// 发送求值命令到调试器
	responseChan := make(chan interface{}, 1)
	da.debugger.SendCommand(DebugCommand{
		Type:     CmdEvaluate,
		Data:     req.Expression,
		Response: responseChan,
	})

	// 等待响应
	result := <-responseChan

	return &EvaluateResponse{
		Result:             fmt.Sprintf("%v", result),
		VariablesReference: 0,
	}, nil
}

// handleSetVariable 处理设置变量请求
func (da *DebugAdapter) handleSetVariable(conn *DAPConnection, body json.RawMessage) (interface{}, error) {
	// 暂时返回空响应
	return map[string]interface{}{}, nil
}

// handleThreads 处理线程请求
func (da *DebugAdapter) handleThreads(conn *DAPConnection, body json.RawMessage) (interface{}, error) {
	// 返回主线程
	threads := []Thread{
		{
			ID:   1,
			Name: "Main Thread",
		},
	}

	return &ThreadsResponse{Threads: threads}, nil
}

// sendStoppedEvent 发送停止事件
func (da *DebugAdapter) sendStoppedEvent(conn *DAPConnection, reason, description string) error {
	body := map[string]interface{}{
		"reason":            reason,
		"description":       description,
		"threadId":          1,
		"allThreadsStopped": true,
	}

	return da.sendEvent(conn, "stopped", body)
}
