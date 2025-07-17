// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package config_test

import (
	"os"
	"testing"

	"github.com/hulo-lang/hulo/internal/config"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestValidate(t *testing.T) {
	src, err := os.ReadFile("huloc.yaml")
	assert.NoError(t, err)

	var huloc config.Huloc
	err = yaml.Unmarshal(src, &huloc)
	assert.NoError(t, err)

	assert.NoError(t, huloc.Validate())
}
