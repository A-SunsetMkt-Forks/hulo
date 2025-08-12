// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package config

// BatchOptions is the configuration for the Batch compiler.
type BatchOptions struct {
	CommentSyntax string `yaml:"comment_syntax" validate:"oneof=rem double_colon"`
	BoolFormat    string `yaml:"bool_format" validate:"oneof=number string cmd"`
}

func DefaultBatchOptions() *BatchOptions {
	return &BatchOptions{
		CommentSyntax: "rem",
		BoolFormat:    "number",
	}
}
