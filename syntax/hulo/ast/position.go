// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package ast

import (
	"github.com/hulo-lang/hulo/syntax/hulo/token"
)

// FileSet 管理多个文件的位置信息
type FileSet struct {
	files map[string]*FileInfo
}

// FileInfo 存储单个文件的位置信息
type FileInfo struct {
	Name     string
	Content  string
	Lines    []int // 每行的起始偏移量
	LineEnds []int // 每行的结束偏移量
}

// NewFileSet 创建新的文件集
func NewFileSet() *FileSet {
	return &FileSet{
		files: make(map[string]*FileInfo),
	}
}

// AddFile 添加文件到文件集
func (fs *FileSet) AddFile(name, content string) {
	lines, lineEnds := calculateLineOffsets(content)
	fs.files[name] = &FileInfo{
		Name:     name,
		Content:  content,
		Lines:    lines,
		LineEnds: lineEnds,
	}
}

// GetFile 获取文件信息
func (fs *FileSet) GetFile(name string) *FileInfo {
	return fs.files[name]
}

// calculateLineOffsets 计算每行的起始和结束偏移量
func calculateLineOffsets(content string) (lines, lineEnds []int) {
	lines = []int{0} // 第一行从偏移量0开始
	lineEnds = []int{}

	for i, char := range content {
		if char == '\n' {
			lines = append(lines, i+1) // 下一行的起始位置
			lineEnds = append(lineEnds, i) // 当前行的结束位置
		}
	}

	// 处理最后一行（如果没有换行符）
	if len(content) > 0 {
		lineEnds = append(lineEnds, len(content)-1)
	}

	return lines, lineEnds
}

// PosToLineColumn 将 token.Pos 转换为行号和列号
func (fi *FileInfo) PosToLineColumn(pos token.Pos) (line, column int) {
	if !pos.IsValid() {
		return 1, 1
	}

	posInt := int(pos)

	// 查找对应的行号
	for i, lineStart := range fi.Lines {
		if posInt < lineStart {
			line = i
			break
		}
		line = i + 1
	}

	// 如果位置超出范围，使用最后一行
	if line > len(fi.Lines) {
		line = len(fi.Lines)
	}

	// 计算列号
	if line > 0 && line <= len(fi.Lines) {
		lineStart := fi.Lines[line-1]
		column = posInt - lineStart + 1
	} else {
		column = 1
	}

	return line, column
}

// LineColumnToPos 将行号和列号转换为 token.Pos
func (fi *FileInfo) LineColumnToPos(line, column int) token.Pos {
	if line <= 0 || line > len(fi.Lines) {
		return token.NoPos
	}

	lineStart := fi.Lines[line-1]
	offset := lineStart + column - 1

	if offset >= len(fi.Content) {
		return token.NoPos
	}

	return token.Pos(offset)
}

// GetLineContent 获取指定行的内容
func (fi *FileInfo) GetLineContent(line int) string {
	if line <= 0 || line > len(fi.Lines) {
		return ""
	}

	start := fi.Lines[line-1]
	var end int

	if line < len(fi.Lines) {
		end = fi.Lines[line] - 1 // 不包括换行符
	} else {
		end = len(fi.Content)
	}

	if start >= len(fi.Content) || end > len(fi.Content) || start >= end {
		return ""
	}

	return fi.Content[start:end]
}

// GetNodeRange 获取节点的范围信息
func (fi *FileInfo) GetNodeRange(node Node) (startLine, startCol, endLine, endCol int) {
	if node == nil {
		return 1, 1, 1, 1
	}

	startLine, startCol = fi.PosToLineColumn(node.Pos())
	endLine, endCol = fi.PosToLineColumn(node.End())

	return startLine, startCol, endLine, endCol
}

// GetNodeText 获取节点对应的源代码文本
func (fi *FileInfo) GetNodeText(node Node) string {
	if node == nil {
		return ""
	}

	start := node.Pos()
	end := node.End()

	if !start.IsValid() || !end.IsValid() {
		return ""
	}

	startInt := int(start)
	endInt := int(end)

	if startInt >= len(fi.Content) || endInt > len(fi.Content) || startInt >= endInt {
		return ""
	}

	return fi.Content[startInt:endInt]
}

// GetContext 获取节点周围的上下文
func (fi *FileInfo) GetContext(node Node, contextLines int) []string {
	if node == nil {
		return nil
	}

	startLine, _, _, _ := fi.GetNodeRange(node)

	var lines []string
	start := max(1, startLine-contextLines)
	end := min(len(fi.Lines), startLine+contextLines)

	for i := start; i <= end; i++ {
		lineContent := fi.GetLineContent(i)
		if lineContent != "" {
			lines = append(lines, lineContent)
		}
	}

	return lines
}

// 辅助函数
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
