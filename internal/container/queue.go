// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package container

// Queue defines the interface for queue operations (FIFO).
type Queue[T any] interface {
	Enqueue(item T)
	Dequeue() (T, bool)
	Peek() (T, bool)
	IsEmpty() bool
	Size() int
	Clear()
	Clone() Queue[T]
}

// SliceQueue represents a queue implementation using a slice.
type SliceQueue[T any] struct {
	items []T
}

// Ensure SliceQueue implements Queue interface.
var _ Queue[any] = (*SliceQueue[any])(nil)

// NewSliceQueue creates a new empty queue.
func NewSliceQueue[T any]() Queue[T] {
	return &SliceQueue[T]{
		items: make([]T, 0),
	}
}

// Enqueue adds an item to the end of the queue.
func (q *SliceQueue[T]) Enqueue(item T) {
	q.items = append(q.items, item)
}

// Dequeue removes and returns the item from the front of the queue.
// Returns zero value and false if the queue is empty.
func (q *SliceQueue[T]) Dequeue() (T, bool) {
	var zero T
	if len(q.items) == 0 {
		return zero, false
	}

	item := q.items[0]
	q.items = q.items[1:]
	return item, true
}

// Peek returns the front item without removing it.
// Returns zero value and false if the queue is empty.
func (q *SliceQueue[T]) Peek() (T, bool) {
	var zero T
	if len(q.items) == 0 {
		return zero, false
	}
	return q.items[0], true
}

// IsEmpty returns true if the queue has no items.
func (q *SliceQueue[T]) IsEmpty() bool {
	return len(q.items) == 0
}

// Size returns the number of items in the queue.
func (q *SliceQueue[T]) Size() int {
	return len(q.items)
}

// Clear removes all items from the queue.
func (q *SliceQueue[T]) Clear() {
	q.items = make([]T, 0)
}

// Clone creates a shallow copy of the queue.
func (q *SliceQueue[T]) Clone() Queue[T] {
	clone := &SliceQueue[T]{
		items: make([]T, len(q.items)),
	}
	copy(clone.items, q.items)
	return clone
}
