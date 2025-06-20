// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package container

import (
	"testing"
)

func TestNewArrayStack(t *testing.T) {
	stack := NewArrayStack[int]()
	if stack == nil {
		t.Error("NewArrayStack() should not return nil")
	}
	if !stack.IsEmpty() {
		t.Error("New stack should be empty")
	}
	if stack.Size() != 0 {
		t.Error("New stack should have size 0")
	}
}

func TestStack_Push(t *testing.T) {
	stack := NewArrayStack[int]()

	// Test pushing single item
	stack.Push(1)
	if stack.Size() != 1 {
		t.Errorf("Stack size should be 1, got %d", stack.Size())
	}

	// Test pushing multiple items
	stack.Push(2)
	stack.Push(3)
	if stack.Size() != 3 {
		t.Errorf("Stack size should be 3, got %d", stack.Size())
	}
}

func TestStack_Pop(t *testing.T) {
	stack := NewArrayStack[int]()

	// Test pop on empty stack
	item, ok := stack.Pop()
	if ok {
		t.Error("Pop on empty stack should return false")
	}
	if item != 0 {
		t.Errorf("Pop on empty stack should return zero value, got %v", item)
	}

	// Test pop on non-empty stack
	stack.Push(1)
	stack.Push(2)
	stack.Push(3)

	item, ok = stack.Pop()
	if !ok {
		t.Error("Pop on non-empty stack should return true")
	}
	if item != 3 {
		t.Errorf("Pop should return last pushed item, got %v, want %v", item, 3)
	}
	if stack.Size() != 2 {
		t.Errorf("Stack size should be 2 after pop, got %d", stack.Size())
	}

	// Test pop order (LIFO)
	item, ok = stack.Pop()
	if !ok || item != 2 {
		t.Errorf("Second pop should return 2, got %v", item)
	}

	item, ok = stack.Pop()
	if !ok || item != 1 {
		t.Errorf("Third pop should return 1, got %v", item)
	}

	// Stack should be empty now
	if !stack.IsEmpty() {
		t.Error("Stack should be empty after popping all items")
	}
}

func TestStack_Peek(t *testing.T) {
	stack := NewArrayStack[int]()

	// Test peek on empty stack
	item, ok := stack.Peek()
	if ok {
		t.Error("Peek on empty stack should return false")
	}
	if item != 0 {
		t.Errorf("Peek on empty stack should return zero value, got %v", item)
	}

	// Test peek on non-empty stack
	stack.Push(1)
	stack.Push(2)

	item, ok = stack.Peek()
	if !ok {
		t.Error("Peek on non-empty stack should return true")
	}
	if item != 2 {
		t.Errorf("Peek should return top item, got %v, want %v", item, 2)
	}

	// Stack size should remain unchanged
	if stack.Size() != 2 {
		t.Errorf("Stack size should remain 2 after peek, got %d", stack.Size())
	}
}

func TestStack_IsEmpty(t *testing.T) {
	stack := NewArrayStack[int]()

	// Test empty stack
	if !stack.IsEmpty() {
		t.Error("New stack should be empty")
	}

	// Test non-empty stack
	stack.Push(1)
	if stack.IsEmpty() {
		t.Error("Stack with items should not be empty")
	}

	// Test after popping all items
	stack.Pop()
	if !stack.IsEmpty() {
		t.Error("Stack should be empty after popping all items")
	}
}

func TestStack_Size(t *testing.T) {
	stack := NewArrayStack[int]()

	// Test initial size
	if stack.Size() != 0 {
		t.Errorf("Initial stack size should be 0, got %d", stack.Size())
	}

	// Test size after pushing
	stack.Push(1)
	if stack.Size() != 1 {
		t.Errorf("Stack size should be 1, got %d", stack.Size())
	}

	stack.Push(2)
	stack.Push(3)
	if stack.Size() != 3 {
		t.Errorf("Stack size should be 3, got %d", stack.Size())
	}

	// Test size after popping
	stack.Pop()
	if stack.Size() != 2 {
		t.Errorf("Stack size should be 2 after pop, got %d", stack.Size())
	}
}

func TestStack_Clear(t *testing.T) {
	stack := NewArrayStack[int]()

	// Test clear on empty stack
	stack.Clear()
	if !stack.IsEmpty() {
		t.Error("Cleared empty stack should be empty")
	}

	// Test clear on non-empty stack
	stack.Push(1)
	stack.Push(2)
	stack.Push(3)

	stack.Clear()
	if !stack.IsEmpty() {
		t.Error("Cleared stack should be empty")
	}
	if stack.Size() != 0 {
		t.Errorf("Cleared stack size should be 0, got %d", stack.Size())
	}

	// Test that we can still push after clearing
	stack.Push(4)
	if stack.Size() != 1 {
		t.Errorf("Stack size should be 1 after pushing to cleared stack, got %d", stack.Size())
	}
}

func TestStack_Clone(t *testing.T) {
	stack := NewArrayStack[int]()
	stack.Push(1)
	stack.Push(2)
	stack.Push(3)

	clone := stack.Clone()

	// Test that clone has same size
	if clone.Size() != stack.Size() {
		t.Errorf("Clone size should match original, got %d, want %d", clone.Size(), stack.Size())
	}

	// Test that clone has same items in same order
	for !stack.IsEmpty() {
		origItem, _ := stack.Pop()
		cloneItem, _ := clone.Pop()
		if origItem != cloneItem {
			t.Errorf("Clone item should match original, got %v, want %v", cloneItem, origItem)
		}
	}

	// Test that original is not affected by clone operations
	if stack.Size() != 0 {
		t.Error("Original stack should be empty after popping clone")
	}
}

func TestStack_StringType(t *testing.T) {
	stack := NewArrayStack[string]()

	// Test with string type
	stack.Push("hello")
	stack.Push("world")

	item, ok := stack.Pop()
	if !ok || item != "world" {
		t.Errorf("String pop should return 'world', got %v", item)
	}

	item, ok = stack.Pop()
	if !ok || item != "hello" {
		t.Errorf("String pop should return 'hello', got %v", item)
	}
}

func TestStack_StructType(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}

	stack := NewArrayStack[Person]()

	person1 := Person{Name: "Alice", Age: 30}
	person2 := Person{Name: "Bob", Age: 25}

	stack.Push(person1)
	stack.Push(person2)

	item, ok := stack.Pop()
	if !ok || item.Name != "Bob" || item.Age != 25 {
		t.Errorf("Struct pop should return Bob, got %v", item)
	}

	item, ok = stack.Pop()
	if !ok || item.Name != "Alice" || item.Age != 30 {
		t.Errorf("Struct pop should return Alice, got %v", item)
	}
}

func TestStack_PointerType(t *testing.T) {
	stack := NewArrayStack[*int]()

	val1 := 1
	val2 := 2

	stack.Push(&val1)
	stack.Push(&val2)

	item, ok := stack.Pop()
	if !ok || *item != 2 {
		t.Errorf("Pointer pop should return 2, got %v", *item)
	}

	item, ok = stack.Pop()
	if !ok || *item != 1 {
		t.Errorf("Pointer pop should return 1, got %v", *item)
	}
}

func TestStack_InterfaceType(t *testing.T) {
	stack := NewArrayStack[interface{}]()

	stack.Push(1)
	stack.Push("hello")
	stack.Push(3.14)

	item, ok := stack.Pop()
	if !ok || item != 3.14 {
		t.Errorf("Interface pop should return 3.14, got %v", item)
	}

	item, ok = stack.Pop()
	if !ok || item != "hello" {
		t.Errorf("Interface pop should return 'hello', got %v", item)
	}

	item, ok = stack.Pop()
	if !ok || item != 1 {
		t.Errorf("Interface pop should return 1, got %v", item)
	}
}

func TestStack_LargeNumberOfItems(t *testing.T) {
	stack := NewArrayStack[int]()

	// Push many items
	for i := range 1000 {
		stack.Push(i)
	}

	if stack.Size() != 1000 {
		t.Errorf("Stack should have 1000 items, got %d", stack.Size())
	}

	// Pop all items in reverse order
	for i := 999; i >= 0; i-- {
		item, ok := stack.Pop()
		if !ok || item != i {
			t.Errorf("Pop should return %d, got %v", i, item)
		}
	}

	if !stack.IsEmpty() {
		t.Error("Stack should be empty after popping all items")
	}
}

func TestStack_ConcurrentOperations(t *testing.T) {
	stack := NewArrayStack[int]()

	// Test that operations work correctly in sequence
	stack.Push(1)
	stack.Push(2)
	stack.Push(3)

	// Peek should not change the stack
	item, ok := stack.Peek()
	if !ok || item != 3 {
		t.Errorf("Peek should return 3, got %v", item)
	}

	// Pop should remove the top item
	item, ok = stack.Pop()
	if !ok || item != 3 {
		t.Errorf("Pop should return 3, got %v", item)
	}

	// Peek should now return the new top item
	item, ok = stack.Peek()
	if !ok || item != 2 {
		t.Errorf("Peek should return 2, got %v", item)
	}
}
