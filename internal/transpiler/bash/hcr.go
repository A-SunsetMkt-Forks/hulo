// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package bash

import (
	"github.com/hulo-lang/hulo/internal/transpiler"
	bast "github.com/hulo-lang/hulo/syntax/bash/ast"
	"github.com/hulo-lang/hulo/syntax/bash/astutil"
	hast "github.com/hulo-lang/hulo/syntax/hulo/ast"
)

type RuleID uint

//go:generate stringer -type=RuleID -linecomment -output=hcr_string.go
const (
	RuleIllegal     RuleID = iota // illegal
	RuleMutilString               // mutil_string
	RuleBoolFormat                // bool_format
)

type BoolNumberCodegen struct{}

func (c *BoolNumberCodegen) Strategy() string {
	return "number"
}

func (c *BoolNumberCodegen) Apply(transpiler transpiler.Transpiler[bast.Node], node hast.Node) (bast.Node, error) {
	if _, ok := node.(*hast.TrueLiteral); ok {
		return astutil.Word("1"), nil
	}
	return astutil.Word("0"), nil
}

type BoolStringCodegen struct{}

func (c *BoolStringCodegen) Strategy() string {
	return "string"
}

func (c *BoolStringCodegen) Apply(transpiler transpiler.Transpiler[bast.Node], node hast.Node) (bast.Node, error) {
	if _, ok := node.(*hast.TrueLiteral); ok {
		return astutil.Word(`"true"`), nil
	}
	return astutil.Word(`"false"`), nil
}

type BoolCommandCodegen struct{}

func (c *BoolCommandCodegen) Strategy() string {
	return "command"
}

func (c *BoolCommandCodegen) Apply(transpiler transpiler.Transpiler[bast.Node], node hast.Node) (bast.Node, error) {
	if _, ok := node.(*hast.TrueLiteral); ok {
		return astutil.Word("true"), nil
	}
	return astutil.Word("false"), nil
}
