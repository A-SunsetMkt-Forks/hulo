// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package config

const FILE = "huloc.yaml"

type Huloc struct {
	CompilerOptions CompilerOptions `json:"compilerOptions"`
	Main            string          `json:"main"`
	Language        Language        `json:"language"`
	Include         []string        `json:"include"`
	Exclude         []string        `json:"exclude"`
}

type CompilerOptions struct {
	Bash  *BashOptions  `json:"bash"`
	Batch *BatchOptions `json:"batch"`
}

type Language string

const (
	L_BASH  Language = "bash"
	L_BATCH Language = "batch"
)
