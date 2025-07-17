// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package config

const HuloRCFileName = "hulorc.yaml"

type HuloRC struct {
	Registry    Registry `yaml:"registry"`
	Network     Network  `yaml:"network"`
	HuloModules string   `yaml:"hulo_modules"`
}

type Registry struct {
	Default string            `yaml:"default"`
	Mirrors map[string]string `yaml:"mirrors"`
}

type Network struct {
	Timeout     int `yaml:"timeout"`
	Retries     int `yaml:"retries"`
	Concurrency int `yaml:"concurrency"`
}

// TODO: is required?
// cache:
//   dir: ~/.local/.hulo_modules
//   max_size: 1GB
//   ttl: 30d
