// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package parser_test

import (
	"bytes"

	"os"
	"testing"

	"github.com/hulo-lang/hulo/syntax/bash/ast"
	"github.com/hulo-lang/hulo/syntax/bash/parser"
	"github.com/stretchr/testify/assert"
	"mvdan.cc/sh/v3/syntax"
)

func TestParser(t *testing.T) {
	parser := syntax.NewParser()
	var buf bytes.Buffer
	buf.WriteString(`#!/bin/bash

echo "Hello, World!"
`)
	file, err := parser.Parse(&buf, "test.sh")
	assert.NoError(t, err)
	syntax.DebugPrint(os.Stdout, file)
}

func TestConvert(t *testing.T) {
	parser := parser.NewParser()
	file, err := parser.Parse(`#!/bin/bash

echo "Hello, World!"
`)
	assert.NoError(t, err)
	ast.Print(file)
}
