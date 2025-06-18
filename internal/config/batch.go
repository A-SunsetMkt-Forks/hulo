// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package config

// BatchOptions is the configuration for the Batch compiler.
type BatchOptions struct {
	CommentSyntax string `json:"commentSyntax" validate:"oneof=rem double_colon"`
}
