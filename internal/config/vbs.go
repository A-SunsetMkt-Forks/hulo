// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package config

// VBScriptOptions is the configuration for the VBScript compiler.
type VBScriptOptions struct {
	CommentSyntax string `json:"commentSyntax" validate:"oneof=rem single_quote"`
}
