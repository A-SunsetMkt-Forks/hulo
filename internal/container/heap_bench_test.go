// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package container

import (
	"testing"
)

func BenchmarkNewMinHeap(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		NewMinHeap[int]()
	}
}

func BenchmarkNewMaxHeap(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		NewMaxHeap[int]()
	}
}

func BenchmarkMinHeap_Push(b *testing.B) {
	heap := NewMinHeap[int]()
	b.ReportAllocs()

	for i := 0; b.Loop(); i++ {
		heap.Push(i)
	}
}

func BenchmarkMaxHeap_Push(b *testing.B) {
	heap := NewMaxHeap[int]()
	b.ReportAllocs()

	for i := 0; b.Loop(); i++ {
		heap.Push(i)
	}
}

func BenchmarkMinHeap_Pop(b *testing.B) {
	heap := NewMinHeap[int]()
	// Pre-fill the heap
	for i := 0; i < b.N; i++ {
		heap.Push(i)
	}

	b.ReportAllocs()

	for b.Loop() {
		heap.Pop()
	}
}

func BenchmarkMaxHeap_Pop(b *testing.B) {
	heap := NewMaxHeap[int]()
	// Pre-fill the heap
	for i := 0; i < b.N; i++ {
		heap.Push(i)
	}

	b.ReportAllocs()

	for b.Loop() {
		heap.Pop()
	}
}

func BenchmarkMinHeap_Peek(b *testing.B) {
	heap := NewMinHeap[int]()
	heap.Push(1)
	heap.Push(2)
	heap.Push(3)

	b.ReportAllocs()

	for b.Loop() {
		heap.Peek()
	}
}

func BenchmarkMaxHeap_Peek(b *testing.B) {
	heap := NewMaxHeap[int]()
	heap.Push(1)
	heap.Push(2)
	heap.Push(3)

	b.ReportAllocs()

	for b.Loop() {
		heap.Peek()
	}
}

func BenchmarkHeap_IsEmpty(b *testing.B) {
	heap := NewMinHeap[int]()
	heap.Push(1)

	b.ReportAllocs()

	for b.Loop() {
		heap.IsEmpty()
	}
}

func BenchmarkHeap_Size(b *testing.B) {
	heap := NewMinHeap[int]()
	heap.Push(1)
	heap.Push(2)
	heap.Push(3)

	b.ReportAllocs()

	for b.Loop() {
		heap.Size()
	}
}

func BenchmarkHeap_Clear(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		heap := NewMinHeap[int]()
		heap.Push(1)
		heap.Push(2)
		heap.Push(3)
		heap.Clear()
	}
}

func BenchmarkHeap_Clone(b *testing.B) {
	heap := NewMinHeap[int]()
	heap.Push(1)
	heap.Push(2)
	heap.Push(3)

	b.ReportAllocs()

	for b.Loop() {
		heap.Clone()
	}
}

func BenchmarkMinHeap_PushPop(b *testing.B) {
	heap := NewMinHeap[int]()
	b.ReportAllocs()

	for i := 0; b.Loop(); i++ {
		heap.Push(i)
		if i%2 == 0 {
			heap.Pop()
		}
	}
}

func BenchmarkMaxHeap_PushPop(b *testing.B) {
	heap := NewMaxHeap[int]()
	b.ReportAllocs()

	for i := 0; b.Loop(); i++ {
		heap.Push(i)
		if i%2 == 0 {
			heap.Pop()
		}
	}
}

func BenchmarkMinHeap_PushPeek(b *testing.B) {
	heap := NewMinHeap[int]()
	b.ReportAllocs()

	for i := 0; b.Loop(); i++ {
		heap.Push(i)
		heap.Peek()
	}
}

func BenchmarkMaxHeap_PushPeek(b *testing.B) {
	heap := NewMaxHeap[int]()
	b.ReportAllocs()

	for i := 0; b.Loop(); i++ {
		heap.Push(i)
		heap.Peek()
	}
}

func BenchmarkHeap_LargeHeap(b *testing.B) {
	heap := NewMinHeap[int]()
	// Pre-fill with large number of items
	for i := range 10000 {
		heap.Push(i)
	}

	b.ReportAllocs()

	for i := 0; b.Loop(); i++ {
		heap.Push(i)
		heap.Pop()
	}
}

func BenchmarkHeap_StringType(b *testing.B) {
	heap := NewMinHeap[string]()
	b.ReportAllocs()

	for b.Loop() {
		heap.Push("test")
	}
}

func BenchmarkHeap_StructType(b *testing.B) {
	type Person struct {
		Name string
		Age  int
	}

	heap := NewHeap(func(a, b Person) bool {
		return a.Age < b.Age
	}, MinHeap)
	person := Person{Name: "Alice", Age: 30}

	b.ReportAllocs()

	for b.Loop() {
		heap.Push(person)
	}
}

func BenchmarkHeap_MixedOperations(b *testing.B) {
	heap := NewMinHeap[int]()
	b.ReportAllocs()

	for i := 0; b.Loop(); i++ {
		heap.Push(i)
		heap.Peek()
		if i%3 == 0 {
			heap.Pop()
		}
		heap.Size()
		if i%5 == 0 {
			heap.IsEmpty()
		}
	}
}

func BenchmarkHeap_GrowAndShrink(b *testing.B) {
	heap := NewMinHeap[int]()
	b.ReportAllocs()

	for b.Loop() {
		// Grow the heap
		for j := range 100 {
			heap.Push(j)
		}
		// Shrink the heap
		for range 100 {
			heap.Pop()
		}
	}
}

func BenchmarkHeap_CloneLargeHeap(b *testing.B) {
	heap := NewMinHeap[int]()
	// Pre-fill with large number of items
	for i := range 1000 {
		heap.Push(i)
	}

	b.ReportAllocs()

	for b.Loop() {
		heap.Clone()
	}
}

func BenchmarkHeap_ClearLargeHeap(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		heap := NewMinHeap[int]()
		// Pre-fill with large number of items
		for j := range 1000 {
			heap.Push(j)
		}
		heap.Clear()
	}
}

func BenchmarkHeap_EmptyOperations(b *testing.B) {
	heap := NewMinHeap[int]()
	b.ReportAllocs()

	for b.Loop() {
		heap.IsEmpty()
		heap.Size()
		heap.Peek()
		heap.Pop()
	}
}

func BenchmarkHeap_SingleItemOperations(b *testing.B) {
	heap := NewMinHeap[int]()
	heap.Push(1)

	b.ReportAllocs()

	for b.Loop() {
		heap.IsEmpty()
		heap.Size()
		heap.Peek()
	}
}

func BenchmarkHeap_AlternatingPushPop(b *testing.B) {
	heap := NewMinHeap[int]()
	b.ReportAllocs()

	for i := 0; b.Loop(); i++ {
		if i%2 == 0 {
			heap.Push(i)
		} else {
			heap.Pop()
		}
	}
}

func BenchmarkHeapify(b *testing.B) {
	items := make([]int, 1000)
	for i := range 1000 {
		items[i] = i
	}

	b.ReportAllocs()

	for b.Loop() {
		Heapify(items, MinHeap)
	}
}

func BenchmarkHeapifyCustom(b *testing.B) {
	type Task struct {
		Name     string
		Priority int
	}

	tasks := make([]Task, 1000)
	for i := range 1000 {
		tasks[i] = Task{Name: "Task", Priority: i}
	}

	b.ReportAllocs()

	for b.Loop() {
		HeapifyCustom(tasks, func(a, b Task) bool {
			return a.Priority > b.Priority
		}, MaxHeap)
	}
}

func BenchmarkHeap_ReverseOrderPush(b *testing.B) {
	heap := NewMinHeap[int]()
	b.ReportAllocs()

	for b.Loop() {
		heap.Clear()
		// Push in reverse order to test heap property maintenance
		for i := 999; i >= 0; i-- {
			heap.Push(i)
		}
	}
}

func BenchmarkHeap_SortedOrderPush(b *testing.B) {
	heap := NewMinHeap[int]()
	b.ReportAllocs()

	for b.Loop() {
		heap.Clear()
		// Push in sorted order
		for i := range 1000 {
			heap.Push(i)
		}
	}
}

func BenchmarkHeap_RandomOrderPush(b *testing.B) {
	heap := NewMinHeap[int]()
	b.ReportAllocs()

	for b.Loop() {
		heap.Clear()
		// Push in pseudo-random order
		for i := range 1000 {
			heap.Push((i * 7) % 1000)
		}
	}
}

func BenchmarkHeap_ExtractAll(b *testing.B) {
	heap := NewMinHeap[int]()
	// Pre-fill heap
	for i := range 1000 {
		heap.Push(i)
	}

	b.ReportAllocs()

	for b.Loop() {
		// Extract all items
		for !heap.IsEmpty() {
			heap.Pop()
		}
		// Re-fill for next iteration
		for i := range 1000 {
			heap.Push(i)
		}
	}
}
