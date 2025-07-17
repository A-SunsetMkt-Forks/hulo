// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package exec

import (
	"fmt"

	"github.com/opencommand/tinge"
)

var _ Executor = (*HelpExecutor)(nil)

type HelpExecutor struct{}

func (e *HelpExecutor) CanHandle(cmd string) bool {
	return cmd == "help"
}

func (e *HelpExecutor) Execute(cmd string) error {
	printHelp()
	return nil
}

func printHelp() {
	fmt.Println(tinge.Styled().
		Newline().
		With(tinge.Bold, tinge.Blue).
		Text("🚀 Hulo REPL Help").
		Newline().
		String())

	fmt.Println(tinge.Styled().
		With(tinge.Bold).
		Text("Commands:").
		Newline().
		String())

	fmt.Println(tinge.Styled().
		Space(2).
		With(tinge.Green).
		Text("help").
		Space(3).
		Text("Show this help message").
		Newline().
		String())

	fmt.Println(tinge.Styled().
		Space(2).
		With(tinge.Green).
		Text("exit").
		Space(3).
		Text("Exit the REPL").
		Newline().
		String())

	fmt.Println(tinge.Styled().
		Space(2).
		With(tinge.Green).
		Text("quit").
		Space(3).
		Text("Exit the REPL").
		Newline().
		String())

	fmt.Println(tinge.Styled().
		Space(2).
		With(tinge.Green).
		Text("clear").
		Space(3).
		Text("Clear the screen").
		Newline().
		String())

	fmt.Println(tinge.Styled().
		Space(2).
		With(tinge.Green).
		Text("config").
		Space(3).
		Text("Show current configuration").
		Newline().
		String())

	fmt.Println(tinge.Styled().
		Space(2).
		With(tinge.Green).
		Text("version").
		Space(3).
		Text("Show version information").
		Newline().
		String())

	fmt.Println()

	fmt.Println(tinge.Styled().
		With(tinge.Bold).
		Text("Examples:").
		Newline().
		String())

	fmt.Println(tinge.Styled().
		Space(2).
		With(tinge.Yellow).
		Text("let x: num = 42").
		Newline().
		String())

	fmt.Println(tinge.Styled().
		Space(2).
		With(tinge.Yellow).
		Text("echo \"Hello, World!\"").
		Newline().
		String())

	fmt.Println(tinge.Styled().
		Space(2).
		With(tinge.Yellow).
		Text("fn add(a: num, b: num) -> num { return a + b }").
		Newline().
		String())

	fmt.Println(tinge.Styled().
		Space(2).
		With(tinge.Yellow).
		Text("if x > 10 { echo \"x is large\" }").
		Newline().
		String())

	fmt.Println(tinge.Styled().
		Space(2).
		With(tinge.Yellow).
		Text("import * from \"std\"").
		Newline().
		String())

	fmt.Println()

	fmt.Println(tinge.Styled().
		With(tinge.Italic).
		Text("💡 Use Tab for autocompletion and Ctrl+C to exit.").
		Newline().
		String())
}
