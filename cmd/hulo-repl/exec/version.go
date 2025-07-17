// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package exec

import (
	"fmt"

	"github.com/hulo-lang/hulo/cmd/meta"
	"github.com/opencommand/tinge"
)

var _ Executor = (*VersionExecutor)(nil)

type VersionExecutor struct{}

func (e *VersionExecutor) CanHandle(cmd string) bool {
	return cmd == "version"
}

func (e *VersionExecutor) Execute(cmd string) error {
	fmt.Println(tinge.Styled().
		Newline().
		Indent(2).
		With(tinge.Bold, tinge.Italic, tinge.Green).
		Text("Hulo-REPL").
		Space().
		Green(meta.Version).
		Newline().
		String())
	// GOOS, GOARCH
	return nil
}
