// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package interpreter

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/caarlos0/log"
)

// StartDAPServer 启动 DAP 服务器
func StartDAPServer(interpreter *Interpreter, address string) error {
	// 创建调试器
	debugger := NewDebugger(interpreter)
	interpreter.debugger = debugger

	// 创建 DAP 适配器
	adapter := NewDebugAdapter(debugger)

	// 启动 DAP 服务器
	if err := adapter.Start(address); err != nil {
		return fmt.Errorf("failed to start DAP server: %w", err)
	}

	log.Infof("DAP server started on %s", address)
	log.Info("Press Ctrl+C to stop the server")

	// 等待中断信号
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan

	// 停止服务器
	log.Info("Stopping DAP server...")
	if err := adapter.Stop(); err != nil {
		log.WithError(err).Error("Failed to stop DAP server")
	}

	return nil
}

// StartDAPServerWithFile 启动 DAP 服务器并加载文件
func StartDAPServerWithFile(filePath string, address string) error {
	// 创建环境
	env := NewEnvironment()

	// 创建解释器
	interpreter := NewInterpreter(env)

	// 启动 DAP 服务器
	return StartDAPServer(interpreter, address)
}
