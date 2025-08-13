// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package linker

import (
	"path/filepath"
	"strconv"
	"strings"

	"maps"

	"github.com/hulo-lang/hulo/internal/vfs"
)

type UnkownSymbol interface {
	Name() string
	Source() string
	Link(symbol *LinkableSymbol) error
}

type BeginEnd struct {
	Begin, End string
}

type Linker struct {
	fs        vfs.VFS
	linkables map[string]*LinkableFile
	listeners map[string][]BeginEnd
}

func NewLinker(fs vfs.VFS) *Linker {
	return &Linker{
		fs:        fs,
		linkables: make(map[string]*LinkableFile),
		listeners: make(map[string][]BeginEnd),
	}
}

func (l *Linker) Read(file string) error {
	content, err := l.fs.ReadFile(file)
	if err != nil {
		return err
	}
	symbols := make(map[string]*LinkableSymbol)
	ext := filepath.Ext(file)

	// 处理文件内容，提取符号
	if listeners, exists := l.listeners[ext]; exists {
		for _, beginEnd := range listeners {
			extractedSymbols := l.extractSymbols(string(content), beginEnd)
			maps.Copy(symbols, extractedSymbols)
		}
	}

	l.linkables[file] = &LinkableFile{
		content: string(content),
		symbols: symbols,
	}
	return nil
}

func (l *Linker) Listen(ext string, be ...BeginEnd) {
	l.listeners[ext] = append(l.listeners[ext], be...)
}

func (l *Linker) Load(file string) *LinkableFile {
	if linkable, exists := l.linkables[file]; exists {
		return linkable
	}
	return nil
}

// extractSymbols 从文件内容中提取符号
func (l *Linker) extractSymbols(content string, beginEnd BeginEnd) map[string]*LinkableSymbol {
	symbols := make(map[string]*LinkableSymbol)
	lines := strings.Split(content, "\n")

	var currentSymbol *LinkableSymbol
	var symbolName string
	var symbolContent []string
	var inSymbol bool

	for i, line := range lines {
		trimmedLine := strings.TrimSpace(line)

		// 检查是否开始符号定义
		if strings.HasPrefix(trimmedLine, beginEnd.Begin) {
			// 提取符号名称和参数
			rest := strings.TrimSpace(trimmedLine[len(beginEnd.Begin):])
			parsed := l.parseSymbolDefinition(rest)

			if parsed.name != "" {
				symbolName = parsed.name
				currentSymbol = &LinkableSymbol{
					tags: parsed.tags,
				}
				symbolContent = []string{}
				inSymbol = true
			}
			continue
		}

		// 检查是否结束符号定义
		if inSymbol && strings.HasPrefix(trimmedLine, beginEnd.End) {
			if currentSymbol != nil && symbolName != "" {
				currentSymbol.text = strings.Join(symbolContent, "\n")
				currentSymbol.tags["start_line"] = strconv.Itoa(i - len(symbolContent))
				currentSymbol.tags["end_line"] = strconv.Itoa(i)
				symbols[symbolName] = currentSymbol
			}
			currentSymbol = nil
			symbolName = ""
			symbolContent = []string{}
			inSymbol = false
			continue
		}

		// 如果在符号内部，收集内容
		if inSymbol && currentSymbol != nil {
			symbolContent = append(symbolContent, line)
		}
	}

	return symbols
}

// parseSymbolDefinition 解析符号定义行，提取名称和标签
func (l *Linker) parseSymbolDefinition(line string) struct {
	name string
	tags map[string]string
} {
	result := struct {
		name string
		tags map[string]string
	}{
		tags: make(map[string]string),
	}

	// 分割字段
	fields := strings.Fields(line)
	if len(fields) == 0 {
		return result
	}

	// 第一个字段是符号名称
	result.name = fields[0]

	// 解析剩余的参数
	for i := 1; i < len(fields); i++ {
		field := fields[i]

		// 处理带引号的参数，如 file="a.vbs"
		if strings.Contains(field, "=") {
			parts := strings.SplitN(field, "=", 2)
			if len(parts) == 2 {
				key := parts[0]
				value := parts[1]

				// 去除引号
				if (strings.HasPrefix(value, `"`) && strings.HasSuffix(value, `"`)) ||
					(strings.HasPrefix(value, `'`) && strings.HasSuffix(value, `'`)) {
					value = value[1 : len(value)-1]
				}

				result.tags[key] = value
			}
		} else {
			// 处理布尔标志，如 inline
			result.tags[field] = "true"
		}
	}

	return result
}

type LinkableFile struct {
	content string
	symbols map[string]*LinkableSymbol
}

func (lf *LinkableFile) Lookup(name string) *LinkableSymbol {
	if symbol, exists := lf.symbols[name]; exists {
		return symbol
	}
	return nil
}

type LinkableSymbol struct {
	text string
	tags map[string]string
}

func (ls *LinkableSymbol) Text() string {
	return ls.text
}

func (ls *LinkableSymbol) GetTag(name string) string {
	return ls.tags[name]
}

func (ls *LinkableSymbol) HasTag(name string) bool {
	_, exists := ls.tags[name]
	return exists
}
