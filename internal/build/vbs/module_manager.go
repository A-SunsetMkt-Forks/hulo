// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package build

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/hulo-lang/hulo/internal/vfs"
	hast "github.com/hulo-lang/hulo/syntax/hulo/ast"
	"github.com/hulo-lang/hulo/syntax/hulo/parser"
	htok "github.com/hulo-lang/hulo/syntax/hulo/token"
)

// ModuleManager 管理模块的导入和依赖
type ModuleManager struct {
	vfs           vfs.VFS
	basePath      string
	importedFiles map[string]*ImportedModule
	moduleCache   map[string]*ModuleInfo
}

// ImportedModule 表示已导入的模块
type ImportedModule struct {
	Path    string                 // 模块路径
	AST     interface{}            // 解析后的AST
	Exports map[string]interface{} // 导出的符号
	Imports []string               // 该模块的导入列表
}

// ModuleInfo 表示模块信息
type ModuleInfo struct {
	Path         string
	Content      []byte
	AST          interface{}
	Transpiled   interface{} // 转换后的代码
	Dependencies []string
}

// NewModuleManager 创建新的模块管理器
func NewModuleManager(vfs vfs.VFS, basePath string) *ModuleManager {
	return &ModuleManager{
		vfs:           vfs,
		basePath:      basePath,
		importedFiles: make(map[string]*ImportedModule),
		moduleCache:   make(map[string]*ModuleInfo),
	}
}

// ImportModule 导入指定路径的模块
func (mm *ModuleManager) ImportModule(importPath string) (*ImportedModule, error) {
	// 检查是否已经导入过
	if imported, exists := mm.importedFiles[importPath]; exists {
		return imported, nil
	}

	// 解析模块路径
	resolvedPath, err := mm.resolveModulePath(importPath)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve module path %s: %w", importPath, err)
	}

	// 检查缓存
	if cached, exists := mm.moduleCache[resolvedPath]; exists {
		imported := &ImportedModule{
			Path:    resolvedPath,
			AST:     cached.AST,
			Exports: mm.extractExports(cached.AST),
			Imports: cached.Dependencies,
		}
		mm.importedFiles[importPath] = imported
		return imported, nil
	}

	// 读取文件内容
	content, err := mm.vfs.ReadFile(resolvedPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read module file %s: %w", resolvedPath, err)
	}

	// 解析AST
	ast, err := parser.ParseSourceScript(string(content))
	if err != nil {
		return nil, fmt.Errorf("failed to parse module %s: %w", resolvedPath, err)
	}

	// 提取导入依赖
	dependencies := mm.extractImports(ast)

	// 缓存模块信息
	moduleInfo := &ModuleInfo{
		Path:         resolvedPath,
		Content:      content,
		AST:          ast,
		Dependencies: dependencies,
	}
	mm.moduleCache[resolvedPath] = moduleInfo

	// 创建导入模块对象
	imported := &ImportedModule{
		Path:    resolvedPath,
		AST:     ast,
		Exports: mm.extractExports(ast),
		Imports: dependencies,
	}

	// 递归处理依赖
	for _, dep := range dependencies {
		if _, err := mm.ImportModule(dep); err != nil {
			return nil, fmt.Errorf("failed to import dependency %s: %w", dep, err)
		}
	}

	mm.importedFiles[importPath] = imported
	return imported, nil
}

// resolveModulePath 解析模块路径
func (mm *ModuleManager) resolveModulePath(importPath string) (string, error) {
	// 移除可能的文件扩展名
	importPath = strings.TrimSuffix(importPath, ".hl")

	// 尝试不同的文件扩展名
	extensions := []string{".hl", ".vbs", ""}

	for _, ext := range extensions {
		// 相对路径
		relativePath := importPath + ext
		if mm.vfs.Exists(relativePath) {
			return relativePath, nil
		}

		// 相对于basePath的路径
		basePath := filepath.Join(mm.basePath, importPath+ext)
		if mm.vfs.Exists(basePath) {
			return basePath, nil
		}

		// 标准库路径
		stdPath := filepath.Join("internal/std", importPath+ext)
		if mm.vfs.Exists(stdPath) {
			return stdPath, nil
		}
	}

	return "", fmt.Errorf("module not found: %s", importPath)
}

// extractImports 从AST中提取导入语句
func (mm *ModuleManager) extractImports(ast interface{}) []string {
	var imports []string

	// 类型断言为File节点
	if file, ok := ast.(*hast.File); ok {
		for _, stmt := range file.Stmts {
			if importStmt, ok := stmt.(*hast.Import); ok {
				if importStmt.ImportSingle != nil {
					imports = append(imports, importStmt.ImportSingle.Path)
				}
			}
		}
	}

	return imports
}

// extractExports 从AST中提取导出的符号
func (mm *ModuleManager) extractExports(ast interface{}) map[string]interface{} {
	exports := make(map[string]interface{})

	// 类型断言为File节点
	if file, ok := ast.(*hast.File); ok {
		for _, stmt := range file.Stmts {
			switch s := stmt.(type) {
			case *hast.FuncDecl:
				// 检查是否有pub修饰符
				for _, modifier := range s.Modifiers {
					if _, ok := modifier.(*hast.PubModifier); ok {
						exports[s.Name.Name] = s
						break
					}
				}
				// 如果没有修饰符，也认为是导出的（默认导出）
				if len(s.Modifiers) == 0 {
					exports[s.Name.Name] = s
				}

			case *hast.ClassDecl:
				// 检查是否有pub修饰符
				if s.Pub.IsValid() {
					exports[s.Name.Name] = s
				}

			case *hast.AssignStmt:
				// 检查是否是常量声明
				if s.Scope == htok.CONST {
					if ident, ok := s.Lhs.(*hast.Ident); ok {
						exports[ident.Name] = s
					}
				}

			case *hast.EnumDecl:
				// 枚举默认导出
				exports[s.Name.Name] = s
			}
		}
	}

	return exports
}

// GetAllImportedModules 获取所有已导入的模块
func (mm *ModuleManager) GetAllImportedModules() map[string]*ImportedModule {
	return mm.importedFiles
}

// GetModuleAST 获取指定模块的AST
func (mm *ModuleManager) GetModuleAST(importPath string) (interface{}, bool) {
	if imported, exists := mm.importedFiles[importPath]; exists {
		return imported.AST, true
	}
	return nil, false
}

// GetModuleExports 获取指定模块的导出
func (mm *ModuleManager) GetModuleExports(importPath string) (map[string]interface{}, bool) {
	if imported, exists := mm.importedFiles[importPath]; exists {
		return imported.Exports, true
	}
	return nil, false
}
