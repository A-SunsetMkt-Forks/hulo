// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package container_test

import (
	"testing"

	"github.com/hulo-lang/hulo/internal/container"
)

func BenchmarkNewSliceQueue(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		container.NewSliceQueue[int]()
	}
}

func BenchmarkQueue_Enqueue(b *testing.B) {
	q := container.NewSliceQueue[int]()
	b.ReportAllocs()

	for i := 0; b.Loop(); i++ {
		q.Enqueue(i)
	}
}

func BenchmarkQueue_Dequeue(b *testing.B) {
	q := container.NewSliceQueue[int]()
	// Pre-fill the queue
	for i := 0; b.Loop(); i++ {
		q.Enqueue(i)
	}
	b.ReportAllocs()

	for b.Loop() {
		q.Dequeue()
	}
}

func BenchmarkQueue_Peek(b *testing.B) {
	q := container.NewSliceQueue[int]()
	q.Enqueue(1)
	q.Enqueue(2)
	q.Enqueue(3)
	b.ReportAllocs()

	for b.Loop() {
		q.Peek()
	}
}

func BenchmarkQueue_IsEmpty(b *testing.B) {
	q := container.NewSliceQueue[int]()
	q.Enqueue(1)
	b.ReportAllocs()

	for b.Loop() {
		q.IsEmpty()
	}
}

func BenchmarkQueue_Size(b *testing.B) {
	q := container.NewSliceQueue[int]()
	q.Enqueue(1)
	q.Enqueue(2)
	q.Enqueue(3)
	b.ReportAllocs()

	for b.Loop() {
		q.Size()
	}
}

func BenchmarkQueue_Clear(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		q := container.NewSliceQueue[int]()
		q.Enqueue(1)
		q.Enqueue(2)
		q.Enqueue(3)
		q.Clear()
	}
}

func BenchmarkQueue_Clone(b *testing.B) {
	q := container.NewSliceQueue[int]()
	q.Enqueue(1)
	q.Enqueue(2)
	q.Enqueue(3)
	b.ReportAllocs()

	for b.Loop() {
		q.Clone()
	}
}

func BenchmarkQueue_EnqueueDequeue(b *testing.B) {
	q := container.NewSliceQueue[int]()
	b.ReportAllocs()

	for i := 0; b.Loop(); i++ {
		q.Enqueue(i)
		q.Dequeue()
	}
}

func BenchmarkQueue_LargeQueue(b *testing.B) {
	q := container.NewSliceQueue[int]()
	// Pre-fill with large number of items
	for i := 0; i < 10000; i++ {
		q.Enqueue(i)
	}
	b.ReportAllocs()

	for i := 0; b.Loop(); i++ {
		q.Enqueue(i)
		q.Dequeue()
	}
}

func BenchmarkQueue_StringType(b *testing.B) {
	q := container.NewSliceQueue[string]()
	b.ReportAllocs()

	for i := 0; b.Loop(); i++ {
		q.Enqueue("test")
	}
}

func BenchmarkQueue_StructType(b *testing.B) {
	type Item struct {
		ID   int
		Data string
	}
	q := container.NewSliceQueue[Item]()
	item := Item{ID: 1, Data: "data"}
	b.ReportAllocs()

	for b.Loop() {
		q.Enqueue(item)
	}
}

func BenchmarkQueue_MixedOperations(b *testing.B) {
	q := container.NewSliceQueue[int]()
	b.ReportAllocs()

	for i := 0; b.Loop(); i++ {
		q.Enqueue(i)
		q.Peek()
		if i%3 == 0 {
			q.Dequeue()
		}
		q.Size()
		if i%5 == 0 {
			q.IsEmpty()
		}
	}
}

func BenchmarkQueue_GrowAndShrink(b *testing.B) {
	q := container.NewSliceQueue[int]()
	b.ReportAllocs()

	for b.Loop() {
		// Grow the queue
		for j := range 100 {
			q.Enqueue(j)
		}
		// Shrink the queue
		for range 100 {
			q.Dequeue()
		}
	}
}
