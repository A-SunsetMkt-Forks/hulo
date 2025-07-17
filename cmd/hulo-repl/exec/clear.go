// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package exec

import "fmt"

var _ Executor = (*ClearExecutor)(nil)

type ClearExecutor struct{}

func (e *ClearExecutor) CanHandle(cmd string) bool {
	return cmd == "clear"
}

func (e *ClearExecutor) Execute(cmd string) error {
	fmt.Print("\033[H\033[2J")
	return nil
}
