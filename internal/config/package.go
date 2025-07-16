// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package config

const HuloPkgFileName = "hulo.pkg.yaml"

type HuloPkg struct {
	Name         string            `yaml:"name"`
	Version      string            `yaml:"version"`
	Description  string            `yaml:"description"`
	Author       string            `yaml:"author"`
	License      string            `yaml:"license"`
	Repository   string            `yaml:"repository"`
	Homepage     string            `yaml:"homepage"`
	Keywords     []string          `yaml:"keywords"`
	Scripts      map[string]string `yaml:"scripts"`
	Dependencies map[string]Dependency      `yaml:"dependencies"`
}

func NewHuloPkg() *HuloPkg {
	return &HuloPkg{
		Name:         "untitled",
		Version:      "0.1.0",
		Description:  "A Hulo package",
		Author:       "",
		License:      "MIT",
		Repository:   "",
		Homepage:     "",
		Keywords:     []string{},
		Scripts:      make(map[string]string),
		Dependencies: make(map[string]Dependency),
	}
}

const HuloPkgLockFileName = "hulo.pkg.lock.yaml"

type HuloPkgLock struct {
	Dependencies []Dependency `yaml:"dependencies"`
}
