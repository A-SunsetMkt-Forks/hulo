// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package suggest

import "github.com/c-bata/go-prompt"

var KeyWords = []prompt.Suggest{
	// module
	{Text: "mod", Description: "Module declaration"},
	{Text: "use", Description: "Use declaration"},
	{Text: "import", Description: "Import declaration"},
	{Text: "from", Description: "Import from specific module"},

	// type system
	{Text: "type", Description: "Type declaration"},
	{Text: "typeof", Description: "Type of expression"},
	{Text: "as", Description: "Type casting"},

	// control flow
	{Text: "if", Description: "If statement"},
	{Text: "else", Description: "Else statement"},
	{Text: "match", Description: "Match statement"},
	{Text: "do", Description: "Do block"},
	{Text: "loop", Description: "Loop statement"},
	{Text: "in", Description: "In operator"},
	{Text: "range", Description: "Range operator"},
	{Text: "continue", Description: "Continue statement"},
	{Text: "break", Description: "Break statement"},

	// modifiers
	{Text: "let", Description: "Let declaration (immutable)"},
	{Text: "var", Description: "Variable declaration (mutable)"},
	{Text: "const", Description: "Constant declaration"},
	{Text: "static", Description: "Static declaration"},
	{Text: "final", Description: "Final declaration"},
	{Text: "pub", Description: "Public visibility"},
	{Text: "required", Description: "Required field"},

	// exception handling
	{Text: "try", Description: "Try block"},
	{Text: "catch", Description: "Catch block"},
	{Text: "finally", Description: "Finally block"},
	{Text: "throw", Description: "Throw statement"},
	{Text: "throws", Description: "Throws declaration"},

	// function
	{Text: "fn", Description: "Function declaration"},
	{Text: "operator", Description: "Operator overloading"},
	{Text: "return", Description: "Return statement"},

	// class and object
	{Text: "enum", Description: "Enum declaration"},
	{Text: "class", Description: "Class declaration"},
	{Text: "trait", Description: "Trait declaration"},
	{Text: "impl", Description: "Implementation block"},
	{Text: "for", Description: "For loop"},
	{Text: "this", Description: "This reference"},
	{Text: "super", Description: "Super reference"},
	{Text: "extend", Description: "Extend declaration"},

	// special
	{Text: "declare", Description: "Declare statement"},
	{Text: "defer", Description: "Defer statement"},
	{Text: "comptime", Description: "Compile-time execution"},
	{Text: "when", Description: "When condition"},
	{Text: "unsafe", Description: "Unsafe block"},

	// type
	{Text: "num", Description: "Number type"},
	{Text: "str", Description: "String type"},
	{Text: "bool", Description: "Boolean type"},
	{Text: "any", Description: "Any type"},

	// literals
	{Text: "true", Description: "Boolean true literal"},
	{Text: "false", Description: "Boolean false literal"},
	{Text: "null", Description: "Null literal"},

	// operators
	{Text: "new", Description: "New operator"},
	{Text: "delete", Description: "Delete operator"},
}

var Commands = []prompt.Suggest{
	{Text: "help", Description: "Show help information"},
	{Text: "exit", Description: "Exit the REPL"},
	{Text: "quit", Description: "Exit the REPL"},
	{Text: "clear", Description: "Clear the screen"},
	{Text: "config", Description: "Show current configuration"},
	{Text: "version", Description: "Show version information"},
}
