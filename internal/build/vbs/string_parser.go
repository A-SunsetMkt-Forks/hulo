// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package build

import (
	"strings"

	vast "github.com/hulo-lang/hulo/syntax/vbs/ast"
)

type StringPart struct {
	Text       string
	IsVariable bool
	Expr       vast.Expr
}

func ParseStringInterpolation(s string) []StringPart {
	var parts []StringPart
	var current strings.Builder
	var i int

	for i < len(s) {
		if s[i] == '$' {
			// 保存当前文本部分
			if current.Len() > 0 {
				parts = append(parts, StringPart{
					Text:       current.String(),
					IsVariable: false,
				})
				current.Reset()
			}

			i++ // 跳过 $

			if i < len(s) && s[i] == '{' {
				// ${name} 形式
				i++ // 跳过 {
				var varName strings.Builder
				for i < len(s) && s[i] != '}' {
					varName.WriteByte(s[i])
					i++
				}
				if i < len(s) {
					i++ // 跳过 }
				}

				parts = append(parts, StringPart{
					Text:       "",
					IsVariable: true,
					Expr: &vast.Ident{
						Name: varName.String(),
					},
				})
			} else {
				// $name 形式
				var varName strings.Builder
				for i < len(s) && (s[i] == '_' || (s[i] >= 'a' && s[i] <= 'z') || (s[i] >= 'A' && s[i] <= 'Z') || (s[i] >= '0' && s[i] <= '9')) {
					varName.WriteByte(s[i])
					i++
				}

				parts = append(parts, StringPart{
					Text:       "",
					IsVariable: true,
					Expr: &vast.Ident{
						Name: varName.String(),
					},
				})
			}
		} else {
			current.WriteByte(s[i])
			i++
		}
	}

	// 添加最后的文本部分
	if current.Len() > 0 {
		parts = append(parts, StringPart{
			Text:       current.String(),
			IsVariable: false,
		})
	}

	return parts
}
