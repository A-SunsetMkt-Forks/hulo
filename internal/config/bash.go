// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package config

// BashOptions is the configuration for the Bash compiler.
type BashOptions struct {
	MultiString string `yaml:"multi_string"`
	BoolFormat  string `yaml:"bool_format" validate:"oneof=number string command"`
}

func DefaultBashOptions() *BashOptions {
	return &BashOptions{
		MultiString: "off",
		BoolFormat:  "number",
	}
}
