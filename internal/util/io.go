// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package util

import (
	"os"

	"sigs.k8s.io/yaml"
)

func LoadConfigure[T any](path string) (T, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return *new(T), err
	}
	var config T
	err = yaml.Unmarshal(content, &config)
	if err != nil {
		return *new(T), err
	}
	return config, nil
}
