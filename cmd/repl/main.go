// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package main

import (
	"fmt"

	"github.com/c-bata/go-prompt"
	"github.com/opencommand/tinge"
)

func executor(in string) {
	fmt.Println(in)
}

func completer(d prompt.Document) []prompt.Suggest {
	return prompt.FilterHasPrefix(keyWords, d.GetWordBeforeCursor(), true)
}

func main() {
	fmt.Println(tinge.Styled().
		Newline().
		Space(2).
		With(tinge.Bold, tinge.Italic, tinge.Green).
		Text("Hulo-REPL").
		Space().
		Green("v1.0.0").
		Newline().
		String())
	prompt.New(
		executor,
		completer,
		prompt.OptionPrefix(">>> "),
	).Run()
}
