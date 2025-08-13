// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package config

type Dependency struct {
	Version      string       `yaml:"version"`
	Resolved     string       `yaml:"resolved"`
	Commit       string       `yaml:"commit"`
	Dependencies []Dependency `yaml:"dependencies"`
	Dev          bool         `yaml:"dev"`
}
