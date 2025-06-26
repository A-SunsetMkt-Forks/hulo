// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package module

import (
	"fmt"
	"sync"
)

// DependencyResolver 依赖解析器
type DependencyResolver struct {
	modules map[string]*ModuleImpl
	visited map[string]bool
	stack   map[string]bool
	order   []string
	mutex   sync.RWMutex
}

// NewDependencyResolver 创建依赖解析器
func NewDependencyResolver() *DependencyResolver {
	return &DependencyResolver{
		modules: make(map[string]*ModuleImpl),
		visited: make(map[string]bool),
		stack:   make(map[string]bool),
		order:   make([]string, 0),
	}
}

// AddModule 添加模块
func (dr *DependencyResolver) AddModule(module *ModuleImpl) {
	dr.mutex.Lock()
	defer dr.mutex.Unlock()
	dr.modules[module.Name()] = module
}

// DetectCycles 检测循环依赖
func (dr *DependencyResolver) DetectCycles() ([]string, error) {
	dr.mutex.Lock()
	defer dr.mutex.Unlock()

	// 重置状态
	dr.visited = make(map[string]bool)
	dr.stack = make(map[string]bool)
	dr.order = make([]string, 0)

	// 检测循环依赖
	for name := range dr.modules {
		if !dr.visited[name] {
			if dr.hasCycle(name) {
				return dr.getCycle(), nil
			}
		}
	}

	return nil, nil
}

// hasCycle 检测是否有循环依赖
func (dr *DependencyResolver) hasCycle(name string) bool {
	if dr.stack[name] {
		return true // 发现循环
	}

	if dr.visited[name] {
		return false // 已访问过
	}

	dr.visited[name] = true
	dr.stack[name] = true

	module := dr.modules[name]
	if module != nil {
		for depName := range module.GetDependencies() {
			if dr.hasCycle(depName) {
				return true
			}
		}
	}

	dr.stack[name] = false
	dr.order = append(dr.order, name)
	return false
}

// getCycle 获取循环依赖路径
func (dr *DependencyResolver) getCycle() []string {
	// 简化实现，返回所有在栈中的模块
	cycle := make([]string, 0)
	for name, inStack := range dr.stack {
		if inStack {
			cycle = append(cycle, name)
		}
	}
	return cycle
}

// TopologicalSort 拓扑排序
func (dr *DependencyResolver) TopologicalSort() ([]string, error) {
	dr.mutex.Lock()
	defer dr.mutex.Unlock()

	// 计算入度
	inDegree := make(map[string]int)
	for name := range dr.modules {
		inDegree[name] = 0
	}

	// 计算每个模块的入度
	for _, module := range dr.modules {
		for depName := range module.GetDependencies() {
			inDegree[depName]++
		}
	}

	// 找到入度为0的节点
	queue := make([]string, 0)
	for name, degree := range inDegree {
		if degree == 0 {
			queue = append(queue, name)
		}
	}

	// 拓扑排序
	result := make([]string, 0)
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		result = append(result, current)

		// 更新依赖当前模块的模块的入度
		currentModule := dr.modules[current]
		if currentModule != nil {
			for depName := range currentModule.dependents {
				inDegree[depName]--
				if inDegree[depName] == 0 {
					queue = append(queue, depName)
				}
			}
		}
	}

	// 检查是否有循环依赖
	if len(result) != len(dr.modules) {
		return nil, fmt.Errorf("circular dependency detected")
	}

	return result, nil
}

// ResolveDependencies 解析所有依赖
func (dr *DependencyResolver) ResolveDependencies() error {
	// 检测循环依赖
	cycles, err := dr.DetectCycles()
	if err != nil {
		return err
	}

	if len(cycles) > 0 {
		return fmt.Errorf("circular dependency detected: %v", cycles)
	}

	// 拓扑排序
	order, err := dr.TopologicalSort()
	if err != nil {
		return err
	}

	// 按顺序初始化模块
	for _, moduleName := range order {
		module := dr.modules[moduleName]
		if module != nil {
			module.SetState(ModuleInitialized)
		}
	}

	return nil
}

// BreakCycle 打破循环依赖
func (dr *DependencyResolver) BreakCycle(cycle []string) error {
	if len(cycle) < 2 {
		return fmt.Errorf("invalid cycle: %v", cycle)
	}

	// 简单的循环打破策略：移除最后一个依赖
	lastModule := dr.modules[cycle[len(cycle)-1]]
	firstModule := dr.modules[cycle[0]]

	if lastModule != nil && firstModule != nil {
		// 移除依赖
		delete(lastModule.deps, firstModule.name)
		delete(firstModule.dependents, lastModule.name)
	}

	return nil
}
