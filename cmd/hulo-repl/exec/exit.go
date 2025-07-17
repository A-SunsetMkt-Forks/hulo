// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package exec

import (
	"fmt"
	"os"
)

var _ Executor = (*ExitExecutor)(nil)

type ExitExecutor struct{}

func (e *ExitExecutor) CanHandle(cmd string) bool {
	return cmd == "exit" || cmd == "quit"
}

func (e *ExitExecutor) Execute(cmd string) error {
	fmt.Print("\r\033[K")
	os.Exit(0)
	return nil
}
