// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package linker

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hulo-lang/hulo/internal/vfs/memvfs"
	vast "github.com/hulo-lang/hulo/syntax/vbs/ast"
	"github.com/stretchr/testify/assert"
)

const runtimeScript = `
begin print1
	Function Println1(a)
		println(a)
	End Function
ending

' HULO_LINK_BEGIN print2
	Function Println2(a)
		println(a)
	End Function
' HULO_LINK_END

REM HULO_LINK_BEGIN print3
	Function Println3(a)
		println(a)
	End Function
REM HULO_LINK_END

begin print4
	Function Println4(a)
		println(a)
	End Function
ending
`

func TestLinkerListen(t *testing.T) {
	fs := memvfs.New()
	fs.WriteFile("runtime.vbs", []byte(runtimeScript), 0644)
	linker := NewLinker(fs)
	linker.Listen(".vbs",
		BeginEnd{"' HULO_LINK_BEGIN", "' HULO_LINK_END"},     // 监听 ' HULO_LINK_BEGIN
		BeginEnd{"REM HULO_LINK_BEGIN", "REM HULO_LINK_END"}, // 监听 REM HULO_LINK_BEGIN
		BeginEnd{"begin", "ending"},                          // 监听自定义分隔符 BEGIN
	)
	err := linker.Read("runtime.vbs")
	assert.NoError(t, err)
	linkable := linker.Load("runtime.vbs")
	assert.Equal(t, linkable.content, runtimeScript)
}

func TestLink(t *testing.T) {
	linker := NewLinker(nil)
	// 按照文件名链接
	// @hulo-link import.vbs
	// 按照符号名称链接
	// @hulo-link import.vbs Import
	err := linker.Read("import.vbs")
	assert.NoError(t, err)
	err = linker.Read("math.vbs")
	assert.NoError(t, err)
	importLinkable := linker.Load("import.vbs")
	symbol := importLinkable.Lookup("Import")
	fmt.Println(symbol.tags["inline"])
}

var vbsFile = &vast.File{
	Stmts: []vast.Stmt{
		&vast.ExprStmt{
			X: &UnknownSymbol{
				Ident: &vast.Ident{
					Name: "' @hulo-link import.vbs Import",
				},
				once: true,
			},
		},
	},
}

type UnknownSymbol struct {
	*vast.Ident
	once bool
}

func TestLinkAST(t *testing.T) {
	_ = NewLinker(nil)

	for _, doc := range vbsFile.Doc {
		for _, comment := range doc.List {
			if strings.HasPrefix(comment.Text, "' @hulo-link") {
				// 去除这条注释, 然后插入代码块

			}
		}
	}
}

func TestLinkWithRelativePath(t *testing.T) {
	_ = NewLinker(nil)
}

func TestParseSymbolDefinition(t *testing.T) {
	linker := NewLinker(nil)

	tests := []struct {
		name     string
		input    string
		expected struct {
			name string
			tags map[string]string
		}
	}{
		{
			name:  "simple symbol",
			input: "Import",
			expected: struct {
				name string
				tags map[string]string
			}{
				name: "Import",
				tags: map[string]string{},
			},
		},
		{
			name:  "symbol with inline flag",
			input: "Import inline",
			expected: struct {
				name string
				tags map[string]string
			}{
				name: "Import",
				tags: map[string]string{
					"inline": "true",
				},
			},
		},
		{
			name:  "symbol with file parameter",
			input: `Import file="a.vbs"`,
			expected: struct {
				name string
				tags map[string]string
			}{
				name: "Import",
				tags: map[string]string{
					"file": "a.vbs",
				},
			},
		},
		{
			name:  "symbol with multiple parameters",
			input: `Import inline file="a.vbs" once`,
			expected: struct {
				name string
				tags map[string]string
			}{
				name: "Import",
				tags: map[string]string{
					"inline": "true",
					"file":   "a.vbs",
					"once":   "true",
				},
			},
		},
		{
			name:  "symbol with single quotes",
			input: `Import file='b.vbs'`,
			expected: struct {
				name string
				tags map[string]string
			}{
				name: "Import",
				tags: map[string]string{
					"file": "b.vbs",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := linker.parseSymbolDefinition(tt.input)
			assert.Equal(t, tt.expected.name, result.name)
			assert.Equal(t, tt.expected.tags, result.tags)
		})
	}
}

func TestExtractSymbols(t *testing.T) {
	linker := NewLinker(nil)

	testContent := `' HULO_LINK_BEGIN Import inline file="a.vbs"
Function Import()
    println("Import function")
End Function
' HULO_LINK_END

' HULO_LINK_BEGIN Math
Function Add(a, b)
    return a + b
End Function
' HULO_LINK_END

' HULO_LINK_BEGIN Utils file="utils.vbs" once
Function Format()
    println("Format function")
End Function
' HULO_LINK_END`

	beginEnd := BeginEnd{
		Begin: "' HULO_LINK_BEGIN",
		End:   "' HULO_LINK_END",
	}

	symbols := linker.extractSymbols(testContent, beginEnd)

	// 验证提取了3个符号
	assert.Equal(t, 3, len(symbols))

	// 验证Import符号
	importSymbol := symbols["Import"]
	assert.NotNil(t, importSymbol)
	assert.Equal(t, "true", importSymbol.tags["inline"])
	assert.Equal(t, "a.vbs", importSymbol.tags["file"])
	assert.Contains(t, importSymbol.tags["content"], "Function Import()")

	// 验证Math符号
	mathSymbol := symbols["Math"]
	assert.NotNil(t, mathSymbol)
	assert.Contains(t, mathSymbol.tags["content"], "Function Add(a, b)")

	// 验证Utils符号
	utilsSymbol := symbols["Utils"]
	assert.NotNil(t, utilsSymbol)
	assert.Equal(t, "utils.vbs", utilsSymbol.tags["file"])
	assert.Equal(t, "true", utilsSymbol.tags["once"])
	assert.Contains(t, utilsSymbol.tags["content"], "Function Format()")
}

func TestExtractSymbolsWithDifferentFormats(t *testing.T) {
	linker := NewLinker(nil)

	testContent := `begin Import inline file="a.vbs"
Function Import()
    println("Import function")
End Function
ending

REM HULO_LINK_BEGIN Math file="math.vbs"
Function Add(a, b)
    return a + b
End Function
REM HULO_LINK_END`

	// 测试自定义分隔符
	beginEnd1 := BeginEnd{
		Begin: "begin",
		End:   "ending",
	}

	symbols1 := linker.extractSymbols(testContent, beginEnd1)
	assert.Equal(t, 1, len(symbols1))

	importSymbol := symbols1["Import"]
	assert.NotNil(t, importSymbol)
	assert.Equal(t, "true", importSymbol.tags["inline"])
	assert.Equal(t, "a.vbs", importSymbol.tags["file"])

	// 测试REM格式
	beginEnd2 := BeginEnd{
		Begin: "REM HULO_LINK_BEGIN",
		End:   "REM HULO_LINK_END",
	}

	symbols2 := linker.extractSymbols(testContent, beginEnd2)
	assert.Equal(t, 1, len(symbols2))

	mathSymbol := symbols2["Math"]
	assert.NotNil(t, mathSymbol)
	assert.Equal(t, "math.vbs", mathSymbol.tags["file"])
}

func TestLinkerReadAndLoad(t *testing.T) {
	fs := memvfs.New()

	testContent := `' HULO_LINK_BEGIN Import inline file="a.vbs"
Function Import()
    println("Import function")
End Function
' HULO_LINK_END`

	fs.WriteFile("test.vbs", []byte(testContent), 0644)

	linker := NewLinker(fs)
	linker.Listen(".vbs", BeginEnd{
		Begin: "' HULO_LINK_BEGIN",
		End:   "' HULO_LINK_END",
	})

	err := linker.Read("test.vbs")
	assert.NoError(t, err)

	linkable := linker.Load("test.vbs")
	assert.NotNil(t, linkable)
	assert.Equal(t, testContent, linkable.content)

	// 验证符号查找
	symbol := linkable.Lookup("Import")
	assert.NotNil(t, symbol)
	assert.True(t, symbol.HasTag("inline"))
	assert.Equal(t, "a.vbs", symbol.GetTag("file"))
	assert.Contains(t, symbol.GetTag("content"), "Function Import()")
}

// @hulo-link import.vbs Import once 链接一次 如果多次链接 会报错

// 内联的情况
// ' HULO_LINK_BEGIN Import inline

// ' HULO_LINK_END

// 单独创建文件的情况 如果两个符号输出一样的话 会合并成一个文件里面的
// ' HULO_LINK_BEGIN abs file=fs.vbs

// ' HULO_LINK_BEGIN rel file=fs.vbs
