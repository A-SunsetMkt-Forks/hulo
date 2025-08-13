// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package container provides generic data structures and algorithms for the Hulo language.
//
// This package includes various container implementations such as:
//
//   - Stack: A last-in-first-out (LIFO) data structure with two implementations:
//   - ArrayStack: Stack implementation using a slice for better performance with small to medium datasets
//   - LinkedStack: Stack implementation using a linked list for memory efficiency and dynamic growth
//
// Each container type is generic and can work with any data type, providing type safety
// at compile time while maintaining runtime efficiency.
//
// Example usage:
//
//	// Create a new array-based stack
//	stack := NewArrayStack[int]()
//	stack.Push(1)
//	stack.Push(2)
//	stack.Push(3)
//
//	// Pop items (LIFO order)
//	item, ok := stack.Pop() // Returns 3, true
//	item, ok = stack.Pop()  // Returns 2, true
//	item, ok = stack.Pop()  // Returns 1, true
//	item, ok = stack.Pop()  // Returns 0, false (empty stack)
//
//	// Create a linked stack for memory efficiency
//	linkedStack := NewLinkedStack[string]()
//	linkedStack.Push("hello")
//	linkedStack.Push("world")
//
//	// Peek without removing
//	top, ok := linkedStack.Peek() // Returns "world", true
//
//	// Check stack properties
//	isEmpty := linkedStack.IsEmpty() // false
//	size := linkedStack.Size()       // 2
//
//	// Clear the stack
//	linkedStack.Clear()
//
//	// Clone a stack
//	original := NewArrayStack[int]()
//	original.Push(1)
//	original.Push(2)
//	clone := original.Clone()
//
// The package is designed to be thread-safe for individual operations but not for
// concurrent access to the same container instance. For concurrent access, external
// synchronization should be used.
package container
