// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package container

import (
	"cmp"
)

// HeapType defines the type of heap
type HeapType int

const (
	MinHeap HeapType = iota
	MaxHeap
)

// Heap defines the interface for heap operations
type Heap[T any] interface {
	Push(item T)
	Pop() (T, bool)
	Peek() (T, bool)
	IsEmpty() bool
	Size() int
	Clear()
	Clone() Heap[T]
	Type() HeapType
}

// SliceHeap represents a heap implementation using a slice
type SliceHeap[T any] struct {
	items []T
	less  func(a, b T) bool
	typ   HeapType
}

// Ensure SliceHeap implements Heap interface
var _ Heap[any] = (*SliceHeap[any])(nil)

// NewMinHeap creates a new min heap
func NewMinHeap[T cmp.Ordered]() Heap[T] {
	return &SliceHeap[T]{
		items: make([]T, 0),
		less:  cmp.Less[T],
		typ:   MinHeap,
	}
}

// NewMaxHeap creates a new max heap
func NewMaxHeap[T cmp.Ordered]() Heap[T] {
	return &SliceHeap[T]{
		items: make([]T, 0),
		less:  func(a, b T) bool { return cmp.Less(b, a) },
		typ:   MaxHeap,
	}
}

// NewHeap creates a new heap with custom comparison function
func NewHeap[T any](less func(a, b T) bool, typ HeapType) Heap[T] {
	return &SliceHeap[T]{
		items: make([]T, 0),
		less:  less,
		typ:   typ,
	}
}

// Push adds an item to the heap
func (h *SliceHeap[T]) Push(item T) {
	h.items = append(h.items, item)
	h.up(len(h.items) - 1)
}

// Pop removes and returns the top item from the heap
// Returns zero value and false if heap is empty
func (h *SliceHeap[T]) Pop() (T, bool) {
	var zero T
	if len(h.items) == 0 {
		return zero, false
	}

	item := h.items[0]
	lastIdx := len(h.items) - 1
	h.items[0] = h.items[lastIdx]
	h.items = h.items[:lastIdx]

	if len(h.items) > 0 {
		h.down(0)
	}

	return item, true
}

// Peek returns the top item without removing it
// Returns zero value and false if heap is empty
func (h *SliceHeap[T]) Peek() (T, bool) {
	var zero T
	if len(h.items) == 0 {
		return zero, false
	}
	return h.items[0], true
}

// IsEmpty returns true if the heap has no items
func (h *SliceHeap[T]) IsEmpty() bool {
	return len(h.items) == 0
}

// Size returns the number of items in the heap
func (h *SliceHeap[T]) Size() int {
	return len(h.items)
}

// Clear removes all items from the heap
func (h *SliceHeap[T]) Clear() {
	h.items = h.items[:0]
}

// Clone creates a copy of the heap
func (h *SliceHeap[T]) Clone() Heap[T] {
	clone := &SliceHeap[T]{
		items: make([]T, len(h.items)),
		less:  h.less,
		typ:   h.typ,
	}
	copy(clone.items, h.items)
	return clone
}

// Type returns the type of the heap (MinHeap or MaxHeap)
func (h *SliceHeap[T]) Type() HeapType {
	return h.typ
}

// up moves an item up in the heap to maintain heap property
func (h *SliceHeap[T]) up(index int) {
	for {
		parent := (index - 1) / 2
		if index == 0 || !h.less(h.items[index], h.items[parent]) {
			break
		}
		h.items[index], h.items[parent] = h.items[parent], h.items[index]
		index = parent
	}
}

// down moves an item down in the heap to maintain heap property
func (h *SliceHeap[T]) down(index int) {
	for {
		left := 2*index + 1
		right := 2*index + 2
		smallest := index

		if left < len(h.items) && h.less(h.items[left], h.items[smallest]) {
			smallest = left
		}
		if right < len(h.items) && h.less(h.items[right], h.items[smallest]) {
			smallest = right
		}

		if smallest == index {
			break
		}

		h.items[index], h.items[smallest] = h.items[smallest], h.items[index]
		index = smallest
	}
}

// heapify converts a slice into a heap in O(n) time
func (h *SliceHeap[T]) heapify() {
	for i := len(h.items)/2 - 1; i >= 0; i-- {
		h.down(i)
	}
}

// Heapify creates a heap from a slice of items
func Heapify[T cmp.Ordered](items []T, typ HeapType) Heap[T] {
	var less func(a, b T) bool
	if typ == MinHeap {
		less = cmp.Less[T]
	} else {
		less = func(a, b T) bool { return cmp.Less(b, a) }
	}

	heap := &SliceHeap[T]{
		items: make([]T, len(items)),
		less:  less,
		typ:   typ,
	}
	copy(heap.items, items)
	heap.heapify()
	return heap
}

// HeapifyCustom creates a heap from a slice of items with custom comparison
func HeapifyCustom[T any](items []T, less func(a, b T) bool, typ HeapType) Heap[T] {
	heap := &SliceHeap[T]{
		items: make([]T, len(items)),
		less:  less,
		typ:   typ,
	}
	copy(heap.items, items)
	heap.heapify()
	return heap
}
