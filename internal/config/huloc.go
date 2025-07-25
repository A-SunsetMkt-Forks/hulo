// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package config

import "github.com/go-playground/validator/v10"

// HulocFileName is the default configuration file name.
const HulocFileName = "huloc.yaml"

// Huloc is the configuration for the Hulo compiler.
type Huloc struct {
	CompilerOptions CompilerOptions `yaml:"compiler_options"`
	Main            string          `yaml:"main" validate:"required,endswith=.hl"`
	Include         []string        `yaml:"include"`
	Exclude         []string        `yaml:"exclude"`
	HuloPath        string          `yaml:"hulopath"`
	// EnableMangle is the option to enable variable name mangling.
	// If true, the variable name will be mangled to a random string, like `_scope_0_1`.
	// If false, the variable name will be used as is.
	EnableMangle bool          `yaml:"enable_mangle"`
	Parser       ParserOptions `yaml:"parser_options"`
	// OutDir is the output directory.
	OutDir string `yaml:"out_dir"`
	// Targets is the targets to compile.
	Targets []string `yaml:"targets"`
}

// Validate validates the Huloc configuration.
func (c *Huloc) Validate() error {
	validate := validator.New()
	return validate.Struct(c)
}

// CompilerOptions is the configuration for the compiler.
type CompilerOptions struct {
	Bash     *BashOptions     `yaml:"bash"`
	Batch    *BatchOptions    `yaml:"batch"`
	VBScript *VBScriptOptions `yaml:"vbs"`
}

// ParserOptions is the options for the parser.
type ParserOptions struct {
	// ShowASTTree is the option to show the AST tree.
	// If the value is "stdout", the AST tree will be shown to the standard output.
	// If the value is "file", the AST tree will be shown to the file.
	// If the value is "none", the AST tree will not be shown.
	ShowASTTree string `yaml:"show_ast_tree"`
	// EnableTracer is the option to enable the tracer.
	EnableTracer bool `yaml:"enable_tracer"`
	// DisableTiming is the option to disable the timing.
	DisableTiming bool `yaml:"disable_timing"`
	// WatchNode is the option to watch the node.
	WatchNode []string `yaml:"watch_node"`
}
