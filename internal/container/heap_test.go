// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package container

import (
	"testing"
)

func TestNewMinHeap(t *testing.T) {
	heap := NewMinHeap[int]()
	if heap == nil {
		t.Error("NewMinHeap() should not return nil")
	}
	if !heap.IsEmpty() {
		t.Error("New min heap should be empty")
	}
	if heap.Size() != 0 {
		t.Error("New min heap should have size 0")
	}
	if heap.Type() != MinHeap {
		t.Error("New min heap should have MinHeap type")
	}
}

func TestNewMaxHeap(t *testing.T) {
	heap := NewMaxHeap[int]()
	if heap == nil {
		t.Error("NewMaxHeap() should not return nil")
	}
	if !heap.IsEmpty() {
		t.Error("New max heap should be empty")
	}
	if heap.Size() != 0 {
		t.Error("New max heap should have size 0")
	}
	if heap.Type() != MaxHeap {
		t.Error("New max heap should have MaxHeap type")
	}
}

func TestMinHeap_Push(t *testing.T) {
	heap := NewMinHeap[int]()

	// Test pushing single item
	heap.Push(3)
	if heap.Size() != 1 {
		t.Errorf("Min heap size should be 1, got %d", heap.Size())
	}

	// Test pushing multiple items
	heap.Push(1)
	heap.Push(2)
	if heap.Size() != 3 {
		t.Errorf("Min heap size should be 3, got %d", heap.Size())
	}
}

func TestMaxHeap_Push(t *testing.T) {
	heap := NewMaxHeap[int]()

	// Test pushing single item
	heap.Push(1)
	if heap.Size() != 1 {
		t.Errorf("Max heap size should be 1, got %d", heap.Size())
	}

	// Test pushing multiple items
	heap.Push(3)
	heap.Push(2)
	if heap.Size() != 3 {
		t.Errorf("Max heap size should be 3, got %d", heap.Size())
	}
}

func TestMinHeap_Pop(t *testing.T) {
	heap := NewMinHeap[int]()

	// Test pop on empty heap
	item, ok := heap.Pop()
	if ok {
		t.Error("Pop on empty min heap should return false")
	}
	if item != 0 {
		t.Errorf("Pop on empty min heap should return zero value, got %v", item)
	}

	// Test pop on non-empty heap
	heap.Push(3)
	heap.Push(1)
	heap.Push(2)

	item, ok = heap.Pop()
	if !ok {
		t.Error("Pop on non-empty min heap should return true")
	}
	if item != 1 {
		t.Errorf("Min heap pop should return smallest item, got %v, want %v", item, 1)
	}
	if heap.Size() != 2 {
		t.Errorf("Min heap size should be 2 after pop, got %d", heap.Size())
	}

	// Test pop order (ascending)
	item, ok = heap.Pop()
	if !ok || item != 2 {
		t.Errorf("Second min heap pop should return 2, got %v", item)
	}

	item, ok = heap.Pop()
	if !ok || item != 3 {
		t.Errorf("Third min heap pop should return 3, got %v", item)
	}

	// Heap should be empty now
	if !heap.IsEmpty() {
		t.Error("Min heap should be empty after popping all items")
	}
}

func TestMaxHeap_Pop(t *testing.T) {
	heap := NewMaxHeap[int]()

	// Test pop on empty heap
	item, ok := heap.Pop()
	if ok {
		t.Error("Pop on empty max heap should return false")
	}
	if item != 0 {
		t.Errorf("Pop on empty max heap should return zero value, got %v", item)
	}

	// Test pop on non-empty heap
	heap.Push(1)
	heap.Push(3)
	heap.Push(2)

	item, ok = heap.Pop()
	if !ok {
		t.Error("Pop on non-empty max heap should return true")
	}
	if item != 3 {
		t.Errorf("Max heap pop should return largest item, got %v, want %v", item, 3)
	}
	if heap.Size() != 2 {
		t.Errorf("Max heap size should be 2 after pop, got %d", heap.Size())
	}

	// Test pop order (descending)
	item, ok = heap.Pop()
	if !ok || item != 2 {
		t.Errorf("Second max heap pop should return 2, got %v", item)
	}

	item, ok = heap.Pop()
	if !ok || item != 1 {
		t.Errorf("Third max heap pop should return 1, got %v", item)
	}

	// Heap should be empty now
	if !heap.IsEmpty() {
		t.Error("Max heap should be empty after popping all items")
	}
}

func TestMinHeap_Peek(t *testing.T) {
	heap := NewMinHeap[int]()

	// Test peek on empty heap
	item, ok := heap.Peek()
	if ok {
		t.Error("Peek on empty min heap should return false")
	}
	if item != 0 {
		t.Errorf("Peek on empty min heap should return zero value, got %v", item)
	}

	// Test peek on non-empty heap
	heap.Push(3)
	heap.Push(1)
	heap.Push(2)

	item, ok = heap.Peek()
	if !ok {
		t.Error("Peek on non-empty min heap should return true")
	}
	if item != 1 {
		t.Errorf("Min heap peek should return smallest item, got %v, want %v", item, 1)
	}

	// Heap size should remain unchanged
	if heap.Size() != 3 {
		t.Errorf("Min heap size should remain 3 after peek, got %d", heap.Size())
	}
}

func TestMaxHeap_Peek(t *testing.T) {
	heap := NewMaxHeap[int]()

	// Test peek on empty heap
	item, ok := heap.Peek()
	if ok {
		t.Error("Peek on empty max heap should return false")
	}
	if item != 0 {
		t.Errorf("Peek on empty max heap should return zero value, got %v", item)
	}

	// Test peek on non-empty heap
	heap.Push(1)
	heap.Push(3)
	heap.Push(2)

	item, ok = heap.Peek()
	if !ok {
		t.Error("Peek on non-empty max heap should return true")
	}
	if item != 3 {
		t.Errorf("Max heap peek should return largest item, got %v, want %v", item, 3)
	}

	// Heap size should remain unchanged
	if heap.Size() != 3 {
		t.Errorf("Max heap size should remain 3 after peek, got %d", heap.Size())
	}
}

func TestHeap_IsEmpty(t *testing.T) {
	heap := NewMinHeap[int]()

	// Test empty heap
	if !heap.IsEmpty() {
		t.Error("New heap should be empty")
	}

	// Test non-empty heap
	heap.Push(1)
	if heap.IsEmpty() {
		t.Error("Heap with items should not be empty")
	}

	// Test after popping all items
	heap.Pop()
	if !heap.IsEmpty() {
		t.Error("Heap should be empty after popping all items")
	}
}

func TestHeap_Size(t *testing.T) {
	heap := NewMinHeap[int]()

	// Test initial size
	if heap.Size() != 0 {
		t.Errorf("Initial heap size should be 0, got %d", heap.Size())
	}

	// Test size after pushing
	heap.Push(1)
	if heap.Size() != 1 {
		t.Errorf("Heap size should be 1, got %d", heap.Size())
	}

	heap.Push(2)
	heap.Push(3)
	if heap.Size() != 3 {
		t.Errorf("Heap size should be 3, got %d", heap.Size())
	}

	// Test size after popping
	heap.Pop()
	if heap.Size() != 2 {
		t.Errorf("Heap size should be 2 after pop, got %d", heap.Size())
	}
}

func TestHeap_Clear(t *testing.T) {
	heap := NewMinHeap[int]()

	// Test clear on empty heap
	heap.Clear()
	if !heap.IsEmpty() {
		t.Error("Cleared empty heap should be empty")
	}

	// Test clear on non-empty heap
	heap.Push(1)
	heap.Push(2)
	heap.Push(3)

	heap.Clear()
	if !heap.IsEmpty() {
		t.Error("Cleared heap should be empty")
	}
	if heap.Size() != 0 {
		t.Errorf("Cleared heap size should be 0, got %d", heap.Size())
	}

	// Test that we can still push after clearing
	heap.Push(4)
	if heap.Size() != 1 {
		t.Errorf("Heap size should be 1 after pushing to cleared heap, got %d", heap.Size())
	}
}

func TestHeap_Clone(t *testing.T) {
	heap := NewMinHeap[int]()
	heap.Push(3)
	heap.Push(1)
	heap.Push(2)

	clone := heap.Clone()

	// Test that clone has same size
	if clone.Size() != heap.Size() {
		t.Errorf("Clone size should match original, got %d, want %d", clone.Size(), heap.Size())
	}

	// Test that clone has same type
	if clone.Type() != heap.Type() {
		t.Errorf("Clone type should match original, got %v, want %v", clone.Type(), heap.Type())
	}

	// Test that clone produces same order
	for !heap.IsEmpty() {
		origItem, _ := heap.Pop()
		cloneItem, _ := clone.Pop()
		if origItem != cloneItem {
			t.Errorf("Clone item should match original, got %v, want %v", cloneItem, origItem)
		}
	}

	// Test that original is not affected by clone operations
	if heap.Size() != 0 {
		t.Error("Original heap should be empty after popping clone")
	}
}

func TestHeap_StringType(t *testing.T) {
	heap := NewMinHeap[string]()

	// Test with string type
	heap.Push("c")
	heap.Push("a")
	heap.Push("b")

	item, ok := heap.Pop()
	if !ok || item != "a" {
		t.Errorf("String min heap pop should return 'a', got %v", item)
	}

	item, ok = heap.Pop()
	if !ok || item != "b" {
		t.Errorf("String min heap pop should return 'b', got %v", item)
	}

	item, ok = heap.Pop()
	if !ok || item != "c" {
		t.Errorf("String min heap pop should return 'c', got %v", item)
	}
}

func TestHeap_StructType(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}

	// Create min heap based on age
	heap := NewHeap(func(a, b Person) bool {
		return a.Age < b.Age
	}, MinHeap)

	person1 := Person{Name: "Alice", Age: 30}
	person2 := Person{Name: "Bob", Age: 25}
	person3 := Person{Name: "Charlie", Age: 35}

	heap.Push(person1)
	heap.Push(person2)
	heap.Push(person3)

	item, ok := heap.Pop()
	if !ok || item.Name != "Bob" || item.Age != 25 {
		t.Errorf("Struct heap pop should return Bob, got %v", item)
	}

	item, ok = heap.Pop()
	if !ok || item.Name != "Alice" || item.Age != 30 {
		t.Errorf("Struct heap pop should return Alice, got %v", item)
	}

	item, ok = heap.Pop()
	if !ok || item.Name != "Charlie" || item.Age != 35 {
		t.Errorf("Struct heap pop should return Charlie, got %v", item)
	}
}

func TestHeap_LargeNumberOfItems(t *testing.T) {
	heap := NewMinHeap[int]()

	// Push many items in reverse order
	for i := 999; i >= 0; i-- {
		heap.Push(i)
	}

	if heap.Size() != 1000 {
		t.Errorf("Heap should have 1000 items, got %d", heap.Size())
	}

	// Pop all items in ascending order
	for i := 0; i < 1000; i++ {
		item, ok := heap.Pop()
		if !ok || item != i {
			t.Errorf("Pop should return %d, got %v", i, item)
		}
	}

	if !heap.IsEmpty() {
		t.Error("Heap should be empty after popping all items")
	}
}

func TestHeapify(t *testing.T) {
	items := []int{3, 1, 4, 1, 5, 9, 2, 6}

	// Test min heapify
	minHeap := Heapify(items, MinHeap)
	if minHeap.Size() != 8 {
		t.Errorf("Min heapify should have 8 items, got %d", minHeap.Size())
	}

	// Check that items come out in ascending order
	expected := []int{1, 1, 2, 3, 4, 5, 6, 9}
	for i, exp := range expected {
		item, ok := minHeap.Pop()
		if !ok || item != exp {
			t.Errorf("Min heapify pop %d should return %d, got %v", i, exp, item)
		}
	}

	// Test max heapify
	maxHeap := Heapify(items, MaxHeap)
	if maxHeap.Size() != 8 {
		t.Errorf("Max heapify should have 8 items, got %d", maxHeap.Size())
	}

	// Check that items come out in descending order
	expected = []int{9, 6, 5, 4, 3, 2, 1, 1}
	for i, exp := range expected {
		item, ok := maxHeap.Pop()
		if !ok || item != exp {
			t.Errorf("Max heapify pop %d should return %d, got %v", i, exp, item)
		}
	}
}

func TestHeapifyCustom(t *testing.T) {
	type Task struct {
		Name     string
		Priority int
	}

	tasks := []Task{
		{"Low", 1},
		{"High", 3},
		{"Medium", 2},
		{"Critical", 4},
	}

	// Create heap based on priority (higher priority first)
	heap := HeapifyCustom(tasks, func(a, b Task) bool {
		return a.Priority > b.Priority
	}, MaxHeap)

	if heap.Size() != 4 {
		t.Errorf("Custom heapify should have 4 items, got %d", heap.Size())
	}

	// Check that items come out in priority order
	expected := []string{"Critical", "High", "Medium", "Low"}
	for i, exp := range expected {
		item, ok := heap.Pop()
		if !ok || item.Name != exp {
			t.Errorf("Custom heapify pop %d should return %s, got %v", i, exp, item.Name)
		}
	}
}

func TestHeap_ConcurrentOperations(t *testing.T) {
	heap := NewMinHeap[int]()

	// Test that operations work correctly in sequence
	heap.Push(3)
	heap.Push(1)
	heap.Push(2)

	// Peek should not change the heap
	item, ok := heap.Peek()
	if !ok || item != 1 {
		t.Errorf("Peek should return 1, got %v", item)
	}

	// Pop should remove the top item
	item, ok = heap.Pop()
	if !ok || item != 1 {
		t.Errorf("Pop should return 1, got %v", item)
	}

	// Peek should now return the new top item
	item, ok = heap.Peek()
	if !ok || item != 2 {
		t.Errorf("Peek should return 2, got %v", item)
	}
}

func TestHeap_Type(t *testing.T) {
	minHeap := NewMinHeap[int]()
	if minHeap.Type() != MinHeap {
		t.Errorf("Min heap type should be MinHeap, got %v", minHeap.Type())
	}

	maxHeap := NewMaxHeap[int]()
	if maxHeap.Type() != MaxHeap {
		t.Errorf("Max heap type should be MaxHeap, got %v", maxHeap.Type())
	}
}
