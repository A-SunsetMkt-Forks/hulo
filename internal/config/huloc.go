// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package config

import "github.com/go-playground/validator/v10"

// NAME is the default configuration file name.
const NAME = "huloc.yaml"

// Huloc is the configuration for the Hulo compiler.
type Huloc struct {
	CompilerOptions CompilerOptions `json:"compilerOptions"`
	Main            string          `json:"main" validate:"required,endswith=.hl"`
	Include         []string        `json:"include"`
	Exclude         []string        `json:"exclude"`
	HuloPath        string          `json:"hulopath"`
	// EnableMangle is the option to enable variable name mangling.
	// If true, the variable name will be mangled to a random string, like `_scope_0_1`.
	// If false, the variable name will be used as is.
	EnableMangle bool          `json:"enableMangle"`
	Parser       ParserOptions `json:"parserOptions"`
	// OutDir is the output directory.
	OutDir string `json:"outDir"`
	// Targets is the targets to compile.
	Targets []string `json:"targets"`
}

// Validate validates the Huloc configuration.
func (c *Huloc) Validate() error {
	validate := validator.New()
	return validate.Struct(c)
}

// CompilerOptions is the configuration for the compiler.
type CompilerOptions struct {
	Bash     *BashOptions     `json:"bash"`
	Batch    *BatchOptions    `json:"batch"`
	VBScript *VBScriptOptions `json:"vbs"`
}

const (
	L_BASH     = "bash"
	L_BATCH    = "batch"
	L_VBSCRIPT = "vbs"
)

// ParserOptions is the options for the parser.
type ParserOptions struct {
	// ShowASTTree is the option to show the AST tree.
	// If the value is "stdout", the AST tree will be shown to the standard output.
	// If the value is "file", the AST tree will be shown to the file.
	// If the value is "none", the AST tree will not be shown.
	ShowASTTree string `json:"showASTTree"`
	// EnableTracer is the option to enable the tracer.
	EnableTracer bool `json:"enableTracer"`
	// DisableTiming is the option to disable the timing.
	DisableTiming bool `json:"disableTiming"`
	// WatchNode is the option to watch the node.
	WatchNode []string `json:"watchNode"`
}
