// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package exec

import (
	"fmt"

	"github.com/opencommand/tinge"
)

var _ Executor = (*ConfigExecutor)(nil)

type ConfigExecutor struct{}

func (e *ConfigExecutor) CanHandle(cmd string) bool {
	return cmd == "config"
}

func (e *ConfigExecutor) Execute(cmd string) error {
	printConfig()
	return nil
}

func printConfig() {
	fmt.Println(tinge.Styled().
		With(tinge.Bold, tinge.Blue).
		Text("⚙️  Current Configuration").
		Newline().
		String())
}
