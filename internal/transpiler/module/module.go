// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package transpiler

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/caarlos0/log"
	"github.com/hulo-lang/hulo/internal/config"
	"github.com/hulo-lang/hulo/internal/container"
	"github.com/hulo-lang/hulo/internal/vfs"
	bast "github.com/hulo-lang/hulo/syntax/batch/ast"
	hast "github.com/hulo-lang/hulo/syntax/hulo/ast"
	"github.com/hulo-lang/hulo/syntax/hulo/parser"
	"gopkg.in/yaml.v3"
)

type Module struct {
	Name         string // 模块名，要处理别名的情况 import "time" as t
	Path         string // 模块绝对路径
	AST          *hast.File
	Exports      map[string]*ExportInfo
	Imports      []*ImportInfo
	Dependencies []string
	Transpiled   *bast.File
	Symbols      *SymbolTable
	State        ModuleState
	Pkg          *config.HuloPkg
}

// 在 main 模块下 查询他的依赖 utils_calculate, add, m_add
func (m *Module) ResolveImportedSymbols(key string) error {
	// 传递 key 进来必须能够映射出来

	// 先查当前模板的符号表有没有

	// 查当前模块的依赖
	return nil
}

type ModuleState int

const (
	ModuleStateUnresolved ModuleState = iota
	ModuleStateResolved
	ModuleStateMangled
)

func (m *Module) IsMangled() bool {
	return m.State == ModuleStateMangled
}

type ImportInfo struct {
	ModulePath string
	SymbolName []string
	Alias      string
	Kind       ImportKind
}

type ImportKind int

const (
	ImportSingle ImportKind = iota
	ImportMulti
	ImportAll
)

type SymbolTable struct{}

type ExportInfo struct{}

// 从 main 入口开始解析
func ResolveAllDependencies(resolver *DependecyResolver, mainFile string) error {

	// 递归解析所有依赖，会收集依赖建立 Imports, 默认情况下 mainFile 是 ./main.hl
	if err := resolver.resolveRecursive(nil, ".", mainFile); err != nil {
		return err
	}

	// 混淆叶节点
	for _, module := range resolver.modules {
		if len(module.Dependencies) == 0 {
			if err := MangleModule(module); err != nil {
				return fmt.Errorf("mangle module: %w", err)
			}
		}
	}

	// 按拓扑顺序混淆其他模块
	for _, module := range resolver.modules {
		if len(module.Dependencies) > 0 {
			// 确保所有依赖都已混淆
			for _, dep := range module.Dependencies {
				if !resolver.modules[dep].IsMangled() {
					return fmt.Errorf("dependency %s not mangled", dep)
				}
			}

			// 混淆当前模块
			if err := MangleModule(module); err != nil {
				return err
			}
		}
	}

	// 混淆完可以建立符号表了

	// math.add() -> _math_add()
	// m.add() -> _math_add()
	// add() -> _math_add()
	// a() -> _math_add()
	// 所有导出 -> 直接映射

	// 对于不同的输出方式，需要建立引用

	return nil
}

func MangleModule(module *Module) error {
	return nil
}

type DependecyResolver struct {
	modules         map[string]*Module
	visited         container.Set[string]
	stack           container.Set[string]
	order           []string
	current         *Module
	pkgs            map[string]*config.HuloPkg
	options         *config.Huloc
	fs              vfs.VFS
	huloModulesPath string
	huloPath        string
}

func NewDependecyResolver() *DependecyResolver {
	return &DependecyResolver{}
}

// 返回绝对路径
func (r *DependecyResolver) resolvePath(parent string, path string) (string, error) {
	log.WithField("parent", parent).Infof("resolvePath: %s", path)
	if r.isRelativePath(path) {
		ret, err := filepath.Abs(filepath.Join(filepath.Dir(parent), path))
		if err != nil {
			return "", fmt.Errorf("failed to resolve path: %w", err)
		}
		if r.fs.Exists(ret) {
			return ret, nil
		}
		if r.fs.Exists(ret + ".hl") {
			return ret + ".hl", nil
		}
		if r.fs.Exists(filepath.Join(ret, "index.hl")) {
			return filepath.Join(ret, "index.hl"), nil
		}
		return "", fmt.Errorf("no hulo.pkg.yaml found for %s", ret)
	}
	if r.isCoreLibPath(path) {
		return filepath.Join(r.huloPath, path, "index.hl"), nil
	}

	// hulo-lang/hello-world hulo-lang/hello-world/abc
	// 要区分是 abc.hl abc/index.hl 以及 hello-world 的入口点在什么地方
	// 根据版本信息，以及suffix 以及 main 表示文件夹入口 确认到底哪个文件
	if ownerRepoVer, suffix, ok := r.isModulePath(path); ok {
		if pkg, ok := r.pkgs[ownerRepoVer]; ok {
			return filepath.Join(ownerRepoVer, pkg.Main, suffix), nil
		}
		pkgPath := filepath.Join(r.huloModulesPath, ownerRepoVer, config.HuloPkgFileName)
		if !r.fs.Exists(pkgPath) {
			return "", fmt.Errorf("no hulo.pkg.yaml found for %s", pkgPath)
		}
		content, err := r.fs.ReadFile(pkgPath)
		if err != nil {
			return "", fmt.Errorf("failed to read pkg file: %w", err)
		}
		pkg := &config.HuloPkg{}
		if err := yaml.Unmarshal(content, pkg); err != nil {
			return "", fmt.Errorf("failed to unmarshal pkg file: %w", err)
		}
		r.pkgs[ownerRepoVer] = pkg
		pend := filepath.Join(r.huloModulesPath, ownerRepoVer, pkg.Main, suffix)
		if r.fs.Exists(filepath.Join(pend, "index.hl")) {
			return filepath.Join(pend, "index.hl"), nil
		}
		return "", fmt.Errorf("no index.hl found for %s", pend)
	}

	return "", fmt.Errorf("invalid path: %s", path)
}

func (r *DependecyResolver) isCoreLibPath(path string) bool {
	return r.fs.Exists(filepath.Join(r.huloPath, path, "index.hl"))
}

func (r *DependecyResolver) isRelativePath(path string) bool {
	return strings.HasPrefix(path, ".")
}

func (r *DependecyResolver) isModulePath(path string) (string, string, bool) {
	if r.current.Pkg == nil {
		return "", "", false
	}

	for name, version := range r.current.Pkg.Dependencies {
		if strings.HasPrefix(path, name) {
			symbolName := strings.TrimPrefix(path, name)
			return filepath.Join(name, version), symbolName, true
		}
	}

	return "", "", false
}

func (r *DependecyResolver) resolveRecursive(parentPkg *Module, parent, filepath string) error {
	absPath, err := r.resolvePath(parent, filepath)
	if err != nil {
		return err
	}
	log.IncreasePadding()
	defer log.DecreasePadding()

	if r.visited.Contains(absPath) {
		return nil
	}

	if r.stack.Contains(absPath) {
		return fmt.Errorf("circular dependency detected: %s", absPath)
	}

	r.visited.Add(absPath)
	r.stack.Add(absPath)
	defer func() { r.stack.Remove(absPath) }()

	log.Infof("loadModule: %s", absPath)
	module, err := r.loadModule(parentPkg, absPath, filepath)
	if err != nil {
		return err
	}
	r.current = module

	// 递归解析依赖
	for _, dep := range module.Dependencies {
		if err := r.resolveRecursive(module, absPath, dep); err != nil {
			return err
		}
	}

	r.order = append(r.order, absPath)

	return nil
}

func (r *DependecyResolver) findNearestPkgYaml(filePath string) (string, error) {
	dir := filepath.Dir(filePath)
	for {
		pkgPath := filepath.Join(dir, config.HuloPkgFileName)
		if r.fs.Exists(pkgPath) {
			return pkgPath, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			// 已经到达根目录
			break
		}
		dir = parent
	}
	return "", fmt.Errorf("no hulo.pkg.yaml found for %s", filePath)
}

// 收集 importInfo 以及 载入配置文件
func (r *DependecyResolver) loadModule(parentPkg *Module, absPath string, symbolName string) (*Module, error) {
	if module, ok := r.modules[absPath]; ok {
		return module, nil
	}

	content, err := r.fs.ReadFile(absPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", absPath, err)
	}

	opts := []parser.ParserOptions{}
	if len(r.options.Parser.ShowASTTree) > 0 {
		switch r.options.Parser.ShowASTTree {
		case "stdout":
			opts = append(opts, parser.OptionTracerASTTree(os.Stdout))
		case "stderr":
			opts = append(opts, parser.OptionTracerASTTree(os.Stderr))
		case "file":
			file, err := r.fs.Open(r.options.Parser.ShowASTTree)
			if err != nil {
				return nil, fmt.Errorf("failed to open file %s: %w", r.options.Parser.ShowASTTree, err)
			}
			defer file.Close()
			opts = append(opts, parser.OptionTracerASTTree(file))
		}
	}
	if !r.options.Parser.EnableTracer {
		opts = append(opts, parser.OptionDisableTracer())
	}
	if r.options.Parser.DisableTiming {
		opts = append(opts, parser.OptionTracerDisableTiming())
	}

	log.WithField("symbolName", symbolName).Infof("parse file: %s", absPath)
	ast, err := parser.ParseSourceScript(string(content), opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to parse file %s: %w", absPath, err)
	}

	imports, dependencies := r.extractImports(ast, absPath)
	log.WithField("dependencies", dependencies).Infof("extractImports: %s", absPath)

	module := &Module{
		Path:         absPath,
		Imports:      imports,
		Dependencies: dependencies,
	}

	if r.isRelativePath(symbolName) {
		pkgPath, err := r.findNearestPkgYaml(absPath)
		if err != nil {
			return nil, fmt.Errorf("failed to find nearest pkg file: %w", err)
		}
		content, err := r.fs.ReadFile(pkgPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read pkg file: %w", err)
		}
		module.Pkg = &config.HuloPkg{}
		if err := yaml.Unmarshal(content, module.Pkg); err != nil {
			return nil, fmt.Errorf("failed to unmarshal pkg file: %w", err)
		}

		r.pkgs[pkgPath] = module.Pkg
	} else if _, _, ok := r.isModulePath(symbolName); ok {
		// 按照 owner/repo/version 的路径去拿包下面的 hulo.pkg.yaml 文件
		pkgPath, err := r.findNearestPkgYaml(absPath)
		if err != nil {
			return nil, fmt.Errorf("failed to find nearest pkg file: %w", err)
		}
		// pkgPath 出来后要读取这个文件
		content, err := r.fs.ReadFile(pkgPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read pkg file: %w", err)
		}

		module.Pkg = &config.HuloPkg{}
		if err := yaml.Unmarshal(content, module.Pkg); err != nil {
			return nil, fmt.Errorf("failed to unmarshal pkg file: %w", err)
		}

		r.pkgs[pkgPath] = module.Pkg
	}

	return module, nil
}

func (r *DependecyResolver) extractImports(ast *hast.File, currentModulePath string) ([]*ImportInfo, []string) {
	imports := []*ImportInfo{}
	dependencies := []string{}
	dependencySet := container.NewMapSet[string]()

	for _, stmt := range ast.Stmts {
		if importStmt, ok := stmt.(*hast.Import); ok {
			var importInfo *ImportInfo

			switch {
			case importStmt.ImportSingle != nil:
				importInfo = &ImportInfo{
					ModulePath: importStmt.ImportSingle.Path,
					Alias:      importStmt.ImportSingle.Alias,
					Kind:       ImportSingle,
				}
			case importStmt.ImportMulti != nil:
				var symbols []string
				for _, field := range importStmt.ImportMulti.List {
					symbols = append(symbols, field.Field)
				}
				importInfo = &ImportInfo{
					ModulePath: importStmt.ImportMulti.Path,
					SymbolName: symbols,
					Kind:       ImportMulti,
				}
			case importStmt.ImportAll != nil:
				importInfo = &ImportInfo{
					ModulePath: importStmt.ImportAll.Path,
					Alias:      importStmt.ImportAll.Alias,
					Kind:       ImportAll,
				}
			}

			if importInfo != nil {
				imports = append(imports, importInfo)
				dependencySet.Add(importInfo.ModulePath)
			}
		}
	}

	for _, dep := range dependencySet.Items() {
		dependencies = append(dependencies, dep)
	}

	return imports, dependencies
}

func (r *DependecyResolver) parseSymbolName(symbolName string) (string, string) {
	return "", ""
}

// 检查循环依赖
func (r *DependecyResolver) detectCycles() (*Cycles, error) {
	return nil, nil
}

func (r *DependecyResolver) topologicalSort() error {
	return nil
}

type Cycles struct{}
