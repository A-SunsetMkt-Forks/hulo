// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package container

// Stack defines the interface for stack operations
type Stack[T any] interface {
	// Push adds an item to the top of the stack
	Push(item T)
	// Pop removes and returns the top item from the stack
	Pop() (T, bool)
	// Peek returns the top item without removing it
	Peek() (T, bool)
	// IsEmpty returns true if the stack has no items
	IsEmpty() bool
	// Size returns the number of items in the stack
	Size() int
	// Clear removes all items from the stack
	Clear()
	// Clone creates a copy of the stack
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

// node represents a node in the linked stack
type node[T any] struct {
	value T
	next  *node[T]
}

// LinkedStack represents a stack implementation using a linked list
type LinkedStack[T any] struct {
	top  *node[T]
	size int
}

// Ensure LinkedStack implements Stack interface
var _ Stack[any] = (*LinkedStack[any])(nil)

// NewLinkedStack creates a new empty linked stack
func NewLinkedStack[T any]() Stack[T] {
	return &LinkedStack[T]{
		top:  nil,
		size: 0,
	}
}

// Push adds an item to the top of the stack
func (s *LinkedStack[T]) Push(item T) {
	newNode := &node[T]{
		value: item,
		next:  s.top,
	}
	s.top = newNode
	s.size++
}

// Pop removes and returns the top item from the stack
// Returns zero value and false if stack is empty
func (s *LinkedStack[T]) Pop() (T, bool) {
	var zero T
	if s.top == nil {
		return zero, false
	}

	item := s.top.value
	s.top = s.top.next
	s.size--
	return item, true
}

// Peek returns the top item without removing it
// Returns zero value and false if stack is empty
func (s *LinkedStack[T]) Peek() (T, bool) {
	var zero T
	if s.top == nil {
		return zero, false
	}
	return s.top.value, true
}

// IsEmpty returns true if the stack has no items
func (s *LinkedStack[T]) IsEmpty() bool {
	return s.top == nil
}

// Size returns the number of items in the stack
func (s *LinkedStack[T]) Size() int {
	return s.size
}

// Clear removes all items from the stack
func (s *LinkedStack[T]) Clear() {
	s.top = nil
	s.size = 0
}

// Clone creates a copy of the stack
func (s *LinkedStack[T]) Clone() Stack[T] {
	clone := &LinkedStack[T]{
		top:  nil,
		size: s.size,
	}

	if s.top == nil {
		return clone
	}

	// Create a map to store nodes by their position
	nodes := make([]*node[T], s.size)
	current := s.top
	for i := range s.size {
		nodes[i] = &node[T]{
			value: current.value,
			next:  nil,
		}
		current = current.next
	}

	// Link the nodes in reverse order (to maintain stack order)
	for i := len(nodes) - 1; i > 0; i-- {
		nodes[i].next = nodes[i-1]
	}
	clone.top = nodes[len(nodes)-1]

	return clone
}
