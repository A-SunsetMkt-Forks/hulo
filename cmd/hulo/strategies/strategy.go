// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package strategies

import "sort"

// Parameters is the parameters for the hulo command
type Parameters struct {
	// Langs is the language to compile.
	// if not set, all languages will be compiled
	Langs []string
	// Verbose is the Verbose mode
	Verbose bool
	// Version is the Version of the compiler
	Version bool
	// OutDir is the directory to write the output
	OutDir string

	// Args is the arguments for the strategy
	Args []string
}

// Strategy is the interface for the strategies
type Strategy interface {
	// CanHandle checks if the strategy can handle the given parameters
	CanHandle(params *Parameters) bool
	// Execute executes the strategy with the given parameters and arguments
	Execute(params *Parameters, args []string) error
	// Priority returns the priority of the strategy and the higher the number, the higher the priority
	Priority() int
}

// StrategyRegistry is the registry for the strategies
type StrategyRegistry struct {
	strategies []Strategy
}

func NewStrategyRegistry() *StrategyRegistry {
	return &StrategyRegistry{
		strategies: []Strategy{},
	}
}

func (r *StrategyRegistry) Register(s Strategy) {
	r.strategies = append(r.strategies, s)
}

func (r *StrategyRegistry) FindStrategy(params *Parameters) Strategy {
	strategies := r.strategies
	sort.Slice(strategies, func(i, j int) bool {
		return strategies[i].Priority() < strategies[j].Priority()
	})

	for _, strategy := range strategies {
		if strategy.CanHandle(params) {
			return strategy
		}
	}

	return nil
}

const (
	PriorityCompile int = iota
	PriorityVersion
)
