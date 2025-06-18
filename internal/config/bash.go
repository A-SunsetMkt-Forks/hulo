// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package config

// BashOptions is the configuration for the Bash compiler.
type BashOptions struct {
	MultiString   string `json:"multiString"`
	BooleanFormat string `json:"booleanFormat" validate:"oneof=number string command"`
}
