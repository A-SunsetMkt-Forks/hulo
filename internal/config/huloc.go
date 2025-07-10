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
	Language        Language        `json:"language"`
	Include         []string        `json:"include"`
	Exclude         []string        `json:"exclude"`
	HULOPATH        string          `json:"hulopath"`
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

type Language string

const (
	L_BASH     Language = "bash"
	L_BATCH    Language = "batch"
	L_VBSCRIPT Language = "vbs"
)
