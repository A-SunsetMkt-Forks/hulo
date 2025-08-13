// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package container

// Set defines the interface for set operations
type Set[T comparable] interface {
	Add(item T)
	Remove(item T)
	Contains(item T) bool
	IsEmpty() bool
	Size() int
	Clear()
	Items() []T
	Union(other Set[T]) Set[T]
	Intersection(other Set[T]) Set[T]
	Difference(other Set[T]) Set[T]
	Clone() Set[T]
}

// MapSet represents a set implementation using a map
type MapSet[T comparable] struct {
	items map[T]struct{}
}

// Ensure MapSet implements Setter interface
var _ Set[any] = (*MapSet[any])(nil)

// NewMapSet creates a new empty set
func NewMapSet[T comparable]() Set[T] {
	return &MapSet[T]{
		items: make(map[T]struct{}),
	}
}

// Add adds an item to the set
func (s *MapSet[T]) Add(item T) {
	s.items[item] = struct{}{}
}

// Remove removes an item from the set
func (s *MapSet[T]) Remove(item T) {
	delete(s.items, item)
}

// Contains returns true if the item is in the set
func (s *MapSet[T]) Contains(item T) bool {
	_, exists := s.items[item]
	return exists
}

// IsEmpty returns true if the set has no items
func (s *MapSet[T]) IsEmpty() bool {
	return len(s.items) == 0
}

// Size returns the number of items in the set
func (s *MapSet[T]) Size() int {
	return len(s.items)
}

// Clear removes all items from the set
func (s *MapSet[T]) Clear() {
	s.items = make(map[T]struct{})
}

// Items returns all items in the set as a slice
func (s *MapSet[T]) Items() []T {
	items := make([]T, 0, len(s.items))
	for item := range s.items {
		items = append(items, item)
	}
	return items
}

// Union returns a new set containing all items from both sets
func (s *MapSet[T]) Union(other Set[T]) Set[T] {
	result := NewMapSet[T]()
	for _, item := range s.Items() {
		result.Add(item)
	}
	for _, item := range other.Items() {
		result.Add(item)
	}
	return result
}

// Intersection returns a new set containing items that exist in both sets
func (s *MapSet[T]) Intersection(other Set[T]) Set[T] {
	result := NewMapSet[T]()
	for _, item := range s.Items() {
		if other.Contains(item) {
			result.Add(item)
		}
	}
	return result
}

// Difference returns a new set containing items that exist in this set but not in the other
func (s *MapSet[T]) Difference(other Set[T]) Set[T] {
	result := NewMapSet[T]()
	for _, item := range s.Items() {
		if !other.Contains(item) {
			result.Add(item)
		}
	}
	return result
}

// Clone returns a new set with the same items as this set
func (s *MapSet[T]) Clone() Set[T] {
	result := NewMapSet[T]()
	for item := range s.items {
		result.Add(item)
	}
	return result
}
