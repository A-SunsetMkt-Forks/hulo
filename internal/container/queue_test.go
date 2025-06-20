// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package container_test

import (
	"testing"

	"github.com/hulo-lang/hulo/internal/container"
	"github.com/stretchr/testify/assert"
)

func TestSliceQueue_New(t *testing.T) {
	q := container.NewSliceQueue[int]()
	assert.NotNil(t, q)
	assert.True(t, q.IsEmpty())
	assert.Equal(t, 0, q.Size())
}

func TestSliceQueue_Enqueue(t *testing.T) {
	q := container.NewSliceQueue[int]()
	q.Enqueue(10)
	assert.False(t, q.IsEmpty())
	assert.Equal(t, 1, q.Size())

	q.Enqueue(20)
	assert.Equal(t, 2, q.Size())
}

func TestSliceQueue_Dequeue(t *testing.T) {
	q := container.NewSliceQueue[string]()
	q.Enqueue("hello")
	q.Enqueue("world")

	item, ok := q.Dequeue()
	assert.True(t, ok)
	assert.Equal(t, "hello", item)
	assert.Equal(t, 1, q.Size())

	item, ok = q.Dequeue()
	assert.True(t, ok)
	assert.Equal(t, "world", item)
	assert.Equal(t, 0, q.Size())
	assert.True(t, q.IsEmpty())
}

func TestSliceQueue_Dequeue_Empty(t *testing.T) {
	q := container.NewSliceQueue[int]()
	item, ok := q.Dequeue()
	assert.False(t, ok)
	assert.Equal(t, 0, item) // zero value for int
	assert.True(t, q.IsEmpty())
}

func TestSliceQueue_Peek(t *testing.T) {
	q := container.NewSliceQueue[int]()
	q.Enqueue(1)
	q.Enqueue(2)

	item, ok := q.Peek()
	assert.True(t, ok)
	assert.Equal(t, 1, item)
	assert.Equal(t, 2, q.Size()) // Peek should not change size
}

func TestSliceQueue_Peek_Empty(t *testing.T) {
	q := container.NewSliceQueue[int]()
	item, ok := q.Peek()
	assert.False(t, ok)
	assert.Equal(t, 0, item) // zero value for int
}

func TestSliceQueue_Clear(t *testing.T) {
	q := container.NewSliceQueue[int]()
	q.Enqueue(100)
	q.Enqueue(200)
	assert.False(t, q.IsEmpty())

	q.Clear()
	assert.True(t, q.IsEmpty())
	assert.Equal(t, 0, q.Size())
}

func TestSliceQueue_Clone(t *testing.T) {
	q1 := container.NewSliceQueue[int]()
	q1.Enqueue(1)
	q1.Enqueue(2)

	q2 := q1.Clone()
	assert.Equal(t, q1.Size(), q2.Size())

	// Dequeue from original, should not affect clone
	item1, ok1 := q1.Dequeue()
	assert.True(t, ok1)
	assert.Equal(t, 1, item1)
	assert.Equal(t, 1, q1.Size())
	assert.Equal(t, 2, q2.Size())

	// Dequeue from clone, should not affect original
	item2, ok2 := q2.Dequeue()
	assert.True(t, ok2)
	assert.Equal(t, 1, item2)
	assert.Equal(t, 1, q1.Size())
	assert.Equal(t, 1, q2.Size())
}
