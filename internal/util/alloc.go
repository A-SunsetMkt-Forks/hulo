// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package util

import (
	"fmt"
	"sync"
)

// AllocatorOption defines a function that configures an Allocator
type AllocatorOption func(*Allocator)

// WithLock sets the lock implementation for the Allocator
//
// Example:
//
//	allocator := NewAllocator(WithLock(&util.NoOpLocker{}))
func WithLock(lock sync.Locker) AllocatorOption {
	return func(a *Allocator) {
		a.lock = lock
	}
}

// WithPrefix sets the prefix for generated variable names
//
// Example:
//
//	allocator := NewAllocator(WithPrefix("_var_"))
func WithPrefix(prefix string) AllocatorOption {
	return func(a *Allocator) {
		a.prefix = prefix
	}
}

// WithInitialCounter sets the initial counter value
//
// Example:
//
//	allocator := NewAllocator(WithInitialCounter(100))
func WithInitialCounter(counter uint64) AllocatorOption {
	return func(a *Allocator) {
		a.counter = counter
	}
}

// Allocator manages variable name allocation and mapping
//
// It provides a way to generate unique variable names and maintain
// a mapping between original names and generated names. This is useful
// for code generation where you need to ensure variable names don't conflict
// and potentially obfuscate the original names.
//
// Example:
//
//	allocator := NewAllocator(
//		WithLock(&util.NoOpLocker{}),
//		WithPrefix("_var_"),
//		WithInitialCounter(100),
//	)
//	name := allocator.AllocName("userCount")  // might return "_var_101"
type Allocator struct {
	// counter for generating unique identifiers
	counter uint64
	// lock for thread safety
	lock sync.Locker
	// mapping from original names to generated names
	nameMap map[string]string
	// reverse mapping from generated names to original names
	reverseMap map[string]string
	// prefix for generated variable names
	prefix string
}

// NewAllocator creates a new Allocator with the given options
//
// Example:
//
//	allocator := NewAllocator(
//		WithLock(&util.NoOpLocker{}),
//		WithPrefix("_var_"),
//		WithInitialCounter(100),
//	)
func NewAllocator(opts ...AllocatorOption) *Allocator {
	// Create allocator with default values
	allocator := &Allocator{
		lock:       &NoOpLocker{}, // Default lock implementation
		nameMap:    make(map[string]string),
		reverseMap: make(map[string]string),
		prefix:     "_v", // Default prefix
	}

	// Apply all options
	for _, opt := range opts {
		opt(allocator)
	}

	return allocator
}

// AllocName allocates a new unique name for the given original name
//
// If the original name has already been allocated, it returns the previously
// generated name. Otherwise, it generates a new unique name.
//
// Example:
//
//	name := allocator.AllocName("count")  // might return "_v1"
//	name2 := allocator.AllocName("count") // returns the same "_v1"
func (a *Allocator) AllocName(originalName string) string {
	a.lock.Lock()
	defer a.lock.Unlock()

	// Check if we already have a mapping for this name
	if generated, exists := a.nameMap[originalName]; exists {
		return generated
	}

	// Generate new name
	a.counter++
	generatedName := fmt.Sprintf("%s%d", a.prefix, a.counter)

	// Store the mapping
	a.nameMap[originalName] = generatedName
	a.reverseMap[generatedName] = originalName

	return generatedName
}

// GetOriginalName returns the original name for a generated name
//
// Returns empty string if the generated name is not found.
//
// Example:
//
//	origName := allocator.GetOriginalName("_v1")  // returns "count"
func (a *Allocator) GetOriginalName(generatedName string) string {
	a.lock.Lock()
	defer a.lock.Unlock()
	return a.reverseMap[generatedName]
}

// GetGeneratedName returns the generated name for an original name
//
// Returns empty string if the original name is not found.
//
// Example:
//
//	genName := allocator.GetGeneratedName("count")  // returns "_v1"
func (a *Allocator) GetGeneratedName(originalName string) string {
	a.lock.Lock()
	defer a.lock.Unlock()
	return a.nameMap[originalName]
}
