// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package interpreter

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"sync"

	"github.com/caarlos0/log"
)

// DebugAdapter DAP 适配器
type DebugAdapter struct {
	debugger        *Debugger
	listener        net.Listener
	connections     map[string]*DAPConnection
	connMutex       sync.RWMutex
	connectionCount int

	// DAP 消息处理
	requestHandlers map[string]RequestHandler
	eventHandlers   map[string]EventHandler

	// 调试会话
	sessionID string
	isRunning bool

	// 上下文
	ctx    context.Context
	cancel context.CancelFunc
}

// DAPConnection DAP 连接
type DAPConnection struct {
	id        string
	conn      net.Conn
	encoder   *json.Encoder
	decoder   *json.Decoder
	sendChan  chan DAPMessage
	closeChan chan struct{}
	adapter   *DebugAdapter
}

// DAPMessage DAP 消息
type DAPMessage struct {
	Type    string          `json:"type"`
	Seq     int             `json:"seq"`
	Command string          `json:"command,omitempty"`
	Event   string          `json:"event,omitempty"`
	Body    json.RawMessage `json:"body,omitempty"`
	Request json.RawMessage `json:"request,omitempty"`
	Success bool            `json:"success,omitempty"`
	Message string          `json:"message,omitempty"`
}

// RequestHandler 请求处理器
type RequestHandler func(*DAPConnection, json.RawMessage) (interface{}, error)

// EventHandler 事件处理器
type EventHandler func(*DAPConnection, json.RawMessage) error

// NewDebugAdapter 创建 DAP 适配器
func NewDebugAdapter(debugger *Debugger) *DebugAdapter {
	ctx, cancel := context.WithCancel(context.Background())

	adapter := &DebugAdapter{
		debugger:        debugger,
		connections:     make(map[string]*DAPConnection),
		requestHandlers: make(map[string]RequestHandler),
		eventHandlers:   make(map[string]EventHandler),
		ctx:             ctx,
		cancel:          cancel,
	}

	// 注册请求处理器
	adapter.registerRequestHandlers()

	// 注册事件处理器
	adapter.registerEventHandlers()

	// 添加 DAP 监听器到调试器
	debugger.AddListener(NewDAPDebugListener(adapter))

	return adapter
}

// Start 启动 DAP 服务器
func (da *DebugAdapter) Start(address string) error {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return fmt.Errorf("failed to start DAP server: %w", err)
	}

	da.listener = listener
	da.isRunning = true

	log.Infof("DAP server started on %s", address)

	go da.acceptConnections()

	return nil
}

// Stop 停止 DAP 服务器
func (da *DebugAdapter) Stop() error {
	da.isRunning = false
	da.cancel()

	if da.listener != nil {
		return da.listener.Close()
	}

	return nil
}

// acceptConnections 接受连接
func (da *DebugAdapter) acceptConnections() {
	for da.isRunning {
		conn, err := da.listener.Accept()
		if err != nil {
			if da.isRunning {
				log.WithError(err).Error("Accept connection error")
			}
			continue
		}

		go da.handleConnection(conn)
	}
}

// handleConnection 处理连接
func (da *DebugAdapter) handleConnection(conn net.Conn) {
	da.connMutex.Lock()
	da.connectionCount++
	connID := fmt.Sprintf("conn_%d", da.connectionCount)
	da.connMutex.Unlock()

	dapConn := &DAPConnection{
		id:        connID,
		conn:      conn,
		encoder:   json.NewEncoder(conn),
		decoder:   json.NewDecoder(conn),
		sendChan:  make(chan DAPMessage, 100),
		closeChan: make(chan struct{}),
		adapter:   da,
	}

	// 添加到连接列表
	da.connMutex.Lock()
	da.connections[connID] = dapConn
	da.connMutex.Unlock()

	defer func() {
		da.connMutex.Lock()
		delete(da.connections, connID)
		da.connMutex.Unlock()
		conn.Close()
		log.Infof("Connection %s closed", connID)
	}()

	log.Infof("New DAP connection: %s from %s", connID, conn.RemoteAddr())

	// 启动消息发送协程
	go da.sendMessages(dapConn)

	// 处理接收的消息
	for {
		var msg DAPMessage
		if err := dapConn.decoder.Decode(&msg); err != nil {
			log.WithError(err).Errorf("Failed to decode message from %s", connID)
			break
		}

		go da.handleMessage(dapConn, msg)
	}
}

// sendMessages 发送消息
func (da *DebugAdapter) sendMessages(conn *DAPConnection) {
	for {
		select {
		case msg := <-conn.sendChan:
			if err := conn.encoder.Encode(msg); err != nil {
				log.WithError(err).Errorf("Failed to send message to %s", conn.id)
				return
			}
		case <-conn.closeChan:
			return
		case <-da.ctx.Done():
			return
		}
	}
}

// handleMessage 处理消息
func (da *DebugAdapter) handleMessage(conn *DAPConnection, msg DAPMessage) {
	log.Debugf("Received message from %s: %s", conn.id, msg.Type)

	switch msg.Type {
	case "request":
		da.handleRequest(conn, msg)
	case "event":
		da.handleEvent(conn, msg)
	default:
		log.Warnf("Unknown message type: %s", msg.Type)
	}
}

// handleRequest 处理请求
func (da *DebugAdapter) handleRequest(conn *DAPConnection, msg DAPMessage) {
	handler, exists := da.requestHandlers[msg.Command]
	if !exists {
		da.sendErrorResponse(conn, msg.Seq, "Unknown command: "+msg.Command)
		return
	}

	result, err := handler(conn, msg.Body)
	if err != nil {
		da.sendErrorResponse(conn, msg.Seq, err.Error())
		return
	}

	da.sendSuccessResponse(conn, msg.Seq, result)
}

// handleEvent 处理事件
func (da *DebugAdapter) handleEvent(conn *DAPConnection, msg DAPMessage) {
	handler, exists := da.eventHandlers[msg.Event]
	if !exists {
		log.Warnf("Unknown event: %s", msg.Event)
		return
	}

	if err := handler(conn, msg.Body); err != nil {
		log.WithError(err).Errorf("Failed to handle event: %s", msg.Event)
	}
}

// sendSuccessResponse 发送成功响应
func (da *DebugAdapter) sendSuccessResponse(conn *DAPConnection, seq int, body interface{}) {
	response := DAPMessage{
		Type:    "response",
		Seq:     seq,
		Success: true,
	}

	if body != nil {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			log.WithError(err).Error("Failed to marshal response body")
			return
		}
		response.Body = bodyBytes
	}

	conn.sendChan <- response
}

// sendErrorResponse 发送错误响应
func (da *DebugAdapter) sendErrorResponse(conn *DAPConnection, seq int, message string) {
	response := DAPMessage{
		Type:    "response",
		Seq:     seq,
		Success: false,
		Message: message,
	}

	conn.sendChan <- response
}

// sendEvent 发送事件
func (da *DebugAdapter) sendEvent(conn *DAPConnection, event string, body interface{}) error {
	msg := DAPMessage{
		Type:  "event",
		Event: event,
	}

	if body != nil {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return err
		}
		msg.Body = bodyBytes
	}

	conn.sendChan <- msg
	return nil
}

// broadcastEvent 广播事件到所有连接
func (da *DebugAdapter) broadcastEvent(event string, body interface{}) {
	da.connMutex.RLock()
	defer da.connMutex.RUnlock()

	for _, conn := range da.connections {
		da.sendEvent(conn, event, body)
	}
}

// registerRequestHandlers 注册请求处理器
func (da *DebugAdapter) registerRequestHandlers() {
	// 初始化请求
	da.requestHandlers["initialize"] = da.handleInitialize
	da.requestHandlers["launch"] = da.handleLaunch
	da.requestHandlers["attach"] = da.handleAttach
	da.requestHandlers["disconnect"] = da.handleDisconnect
	da.requestHandlers["terminate"] = da.handleTerminate

	// 断点管理
	da.requestHandlers["setBreakpoints"] = da.handleSetBreakpoints
	da.requestHandlers["setFunctionBreakpoints"] = da.handleSetFunctionBreakpoints
	da.requestHandlers["setExceptionBreakpoints"] = da.handleSetExceptionBreakpoints

	// 调试控制
	da.requestHandlers["continue"] = da.handleContinue
	da.requestHandlers["next"] = da.handleNext
	da.requestHandlers["stepIn"] = da.handleStepIn
	da.requestHandlers["stepOut"] = da.handleStepOut
	da.requestHandlers["pause"] = da.handlePause

	// 变量和栈
	da.requestHandlers["variables"] = da.handleVariables
	da.requestHandlers["stackTrace"] = da.handleStackTrace
	da.requestHandlers["scopes"] = da.handleScopes
	da.requestHandlers["evaluate"] = da.handleEvaluate
	da.requestHandlers["setVariable"] = da.handleSetVariable

	// 线程
	da.requestHandlers["threads"] = da.handleThreads
}

// registerEventHandlers 注册事件处理器
func (da *DebugAdapter) registerEventHandlers() {
	// 暂时为空，可以根据需要添加事件处理器
}
