// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package module

import (
	"sync"

	"github.com/hulo-lang/hulo/internal/core"
	"maps"
)

type ModuleImpl struct {
	name       string
	path       string
	types      map[string]core.Type
	values     map[string]core.Value
	exports    map[string]*ExportInfo
	imports    map[string]*ImportInfo
	deps       map[string]*ModuleImpl
	dependents map[string]*ModuleImpl
	state      ModuleState
	initOrder  int
	mutex      sync.RWMutex
}

// ExportInfo 导出信息
type ExportInfo struct {
	name     string
	alias    string
	value    any
	kind     ExportKind
	public   bool
	internal bool
}

// ImportInfo 导入信息
type ImportInfo struct {
	name     string
	alias    string
	from     string
	kind     ImportKind
	resolved bool
}

// ExportKind 导出类型
type ExportKind int

const (
	ExportType ExportKind = iota
	ExportFunction
	ExportValue
	ExportOperator
)

// ImportKind 导入类型
type ImportKind int

const (
	ImportDefault ImportKind = iota
	ImportNamed
	ImportNamespace
	ImportAll
)

// ModuleState 模块状态
type ModuleState int

const (
	ModuleUnloaded ModuleState = iota
	ModuleLoading
	ModuleLoaded
	ModuleInitializing
	ModuleInitialized
	ModuleError
)

// NewModule 创建新模块
func NewModule(name, path string) *ModuleImpl {
	return &ModuleImpl{
		name:       name,
		path:       path,
		types:      make(map[string]core.Type),
		values:     make(map[string]core.Value),
		exports:    make(map[string]*ExportInfo),
		imports:    make(map[string]*ImportInfo),
		deps:       make(map[string]*ModuleImpl),
		dependents: make(map[string]*ModuleImpl),
		state:      ModuleUnloaded,
	}
}

// Name 获取模块名
func (m *ModuleImpl) Name() string {
	return m.name
}

// Path 获取模块路径
func (m *ModuleImpl) Path() string {
	return m.path
}

// GetType 获取类型
func (m *ModuleImpl) GetType(name string) (core.Type, bool) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	typ, ok := m.types[name]
	return typ, ok
}

// GetValue 获取值
func (m *ModuleImpl) GetValue(name string) (core.Value, bool) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	val, ok := m.values[name]
	return val, ok
}

// Export 导出
func (m *ModuleImpl) Export(name string, value interface{}) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	exportInfo := &ExportInfo{
		name:   name,
		value:  value,
		public: true,
	}

	m.exports[name] = exportInfo
	return nil
}

// Import 导入
func (m *ModuleImpl) Import(name string, alias string, from string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	importInfo := &ImportInfo{
		name:  name,
		alias: alias,
		from:  from,
		kind:  ImportNamed,
	}

	if alias != "" {
		m.imports[alias] = importInfo
	} else {
		m.imports[name] = importInfo
	}

	return nil
}

// AddDependency 添加依赖
func (m *ModuleImpl) AddDependency(dep *ModuleImpl) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.deps[dep.name] = dep
	dep.dependents[m.name] = m
}

// GetDependencies 获取依赖
func (m *ModuleImpl) GetDependencies() map[string]*ModuleImpl {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	result := make(map[string]*ModuleImpl)
	maps.Copy(result, m.deps)
	return result
}

// SetState 设置状态
func (m *ModuleImpl) SetState(state ModuleState) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.state = state
}

// GetState 获取状态
func (m *ModuleImpl) GetState() ModuleState {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.state
}
