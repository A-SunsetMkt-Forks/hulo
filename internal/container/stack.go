// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package container

// Stack defines the interface for stack operations
type Stack[T any] interface {
	Push(item T)
	Pop() (T, bool)
	Peek() (T, bool)
	IsEmpty() bool
	Size() int
	Clear()
	Clone() Stack[T]
}

// ArrayStack represents a stack implementation using a slice
type ArrayStack[T any] struct {
	items []T
}

// Ensure ArrayStack implements Stack interface
var _ Stack[any] = (*ArrayStack[any])(nil)

// NewArrayStack creates a new empty stack
func NewArrayStack[T any]() Stack[T] {
	return &ArrayStack[T]{
		items: make([]T, 0),
	}
}

// Push adds an item to the top of the stack
func (s *ArrayStack[T]) Push(item T) {
	s.items = append(s.items, item)
}

// Pop removes and returns the top item from the stack
// Returns zero value and false if stack is empty
func (s *ArrayStack[T]) Pop() (T, bool) {
	var zero T
	if len(s.items) == 0 {
		return zero, false
	}

	lastIdx := len(s.items) - 1
	item := s.items[lastIdx]
	s.items = s.items[:lastIdx]
	return item, true
}

// Peek returns the top item without removing it
// Returns zero value and false if stack is empty
func (s *ArrayStack[T]) Peek() (T, bool) {
	var zero T
	if len(s.items) == 0 {
		return zero, false
	}
	return s.items[len(s.items)-1], true
}

// IsEmpty returns true if the stack has no items
func (s *ArrayStack[T]) IsEmpty() bool {
	return len(s.items) == 0
}

// Size returns the number of items in the stack
func (s *ArrayStack[T]) Size() int {
	return len(s.items)
}

// Clear removes all items from the stack
func (s *ArrayStack[T]) Clear() {
	s.items = s.items[:0]
}

// Clone creates a copy of the stack
func (s *ArrayStack[T]) Clone() Stack[T] {
	clone := &ArrayStack[T]{
		items: make([]T, len(s.items)),
	}
	copy(clone.items, s.items)
	return clone
}
