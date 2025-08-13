// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package container

import (
	"testing"
)

func BenchmarkNewArrayStack(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		NewArrayStack[int]()
	}
}

func BenchmarkStack_Push(b *testing.B) {
	stack := NewArrayStack[int]()
	b.ReportAllocs()


	for i := 0; b.Loop(); i++ {
		stack.Push(i)
	}
}

func BenchmarkStack_Pop(b *testing.B) {
	stack := NewArrayStack[int]()
	// Pre-fill the stack
	for i := 0; b.Loop(); i++ {
		stack.Push(i)
	}

	b.ReportAllocs()


	for b.Loop() {
		stack.Pop()
	}
}

func BenchmarkStack_Peek(b *testing.B) {
	stack := NewArrayStack[int]()
	stack.Push(1)
	stack.Push(2)
	stack.Push(3)

	b.ReportAllocs()


	for b.Loop() {
		stack.Peek()
	}
}

func BenchmarkStack_IsEmpty(b *testing.B) {
	stack := NewArrayStack[int]()
	stack.Push(1)

	b.ReportAllocs()


	for b.Loop() {
		stack.IsEmpty()
	}
}

func BenchmarkStack_Size(b *testing.B) {
	stack := NewArrayStack[int]()
	stack.Push(1)
	stack.Push(2)
	stack.Push(3)

	b.ReportAllocs()


	for b.Loop() {
		stack.Size()
	}
}

func BenchmarkStack_Clear(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		stack := NewArrayStack[int]()
		stack.Push(1)
		stack.Push(2)
		stack.Push(3)
		stack.Clear()
	}
}

func BenchmarkStack_Clone(b *testing.B) {
	stack := NewArrayStack[int]()
	stack.Push(1)
	stack.Push(2)
	stack.Push(3)

	b.ReportAllocs()


	for b.Loop() {
		stack.Clone()
	}
}

func BenchmarkStack_PushPop(b *testing.B) {
	stack := NewArrayStack[int]()
	b.ReportAllocs()


	for i := 0; b.Loop(); i++ {
		stack.Push(i)
		if i%2 == 0 {
			stack.Pop()
		}
	}
}

func BenchmarkStack_PushPeek(b *testing.B) {
	stack := NewArrayStack[int]()
	b.ReportAllocs()


	for i := 0; b.Loop(); i++ {
		stack.Push(i)
		stack.Peek()
	}
}

func BenchmarkStack_LargeStack(b *testing.B) {
	stack := NewArrayStack[int]()
	// Pre-fill with large number of items
	for i := range 10000 {
		stack.Push(i)
	}

	b.ReportAllocs()


	for i := 0; b.Loop(); i++ {
		stack.Push(i)
		stack.Pop()
	}
}

func BenchmarkStack_StringType(b *testing.B) {
	stack := NewArrayStack[string]()
	b.ReportAllocs()


	for b.Loop() {
		stack.Push("test")
	}
}

func BenchmarkStack_StructType(b *testing.B) {
	type Person struct {
		Name string
		Age  int
	}

	stack := NewArrayStack[Person]()
	person := Person{Name: "Alice", Age: 30}

	b.ReportAllocs()


	for b.Loop() {
		stack.Push(person)
	}
}

func BenchmarkStack_PointerType(b *testing.B) {
	stack := NewArrayStack[*int]()
	val := 42

	b.ReportAllocs()


	for b.Loop() {
		stack.Push(&val)
	}
}

func BenchmarkStack_InterfaceType(b *testing.B) {
	stack := NewArrayStack[any]()
	b.ReportAllocs()


	for i := 0; b.Loop(); i++ {
		stack.Push(i)
	}
}

func BenchmarkStack_MixedOperations(b *testing.B) {
	stack := NewArrayStack[int]()
	b.ReportAllocs()


	for i := 0; b.Loop(); i++ {
		stack.Push(i)
		stack.Peek()
		if i%3 == 0 {
			stack.Pop()
		}
		stack.Size()
		if i%5 == 0 {
			stack.IsEmpty()
		}
	}
}

func BenchmarkStack_GrowAndShrink(b *testing.B) {
	stack := NewArrayStack[int]()
	b.ReportAllocs()


	for b.Loop() {
		// Grow the stack
		for j := range 100 {
			stack.Push(j)
		}
		// Shrink the stack
		for range 100 {
			stack.Pop()
		}
	}
}

func BenchmarkStack_CloneLargeStack(b *testing.B) {
	stack := NewArrayStack[int]()
	// Pre-fill with large number of items
	for i := range 1000 {
		stack.Push(i)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		stack.Clone()
	}
}

func BenchmarkStack_ClearLargeStack(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		stack := NewArrayStack[int]()
		// Pre-fill with large number of items
		for j := 0; j < 1000; j++ {
			stack.Push(j)
		}
		stack.Clear()
	}
}

func BenchmarkStack_EmptyOperations(b *testing.B) {
	stack := NewArrayStack[int]()
	b.ReportAllocs()


	for b.Loop() {
		stack.IsEmpty()
		stack.Size()
		stack.Peek()
		stack.Pop()
	}
}

func BenchmarkStack_SingleItemOperations(b *testing.B) {
	stack := NewArrayStack[int]()
	stack.Push(1)

	b.ReportAllocs()


	for b.Loop() {
		stack.IsEmpty()
		stack.Size()
		stack.Peek()
	}
}

func BenchmarkStack_AlternatingPushPop(b *testing.B) {
	stack := NewArrayStack[int]()
	b.ReportAllocs()


	for i := 0; b.Loop(); i++ {
		if i%2 == 0 {
			stack.Push(i)
		} else {
			stack.Pop()
		}
	}
}

func BenchmarkArrayStack_Push(b *testing.B) {
	stack := NewArrayStack[int]()

	for i := 0; b.Loop(); i++ {
		stack.Push(i)
	}
}

func BenchmarkLinkedStack_Push(b *testing.B) {
	stack := NewLinkedStack[int]()

	for i := 0; b.Loop(); i++ {
		stack.Push(i)
	}
}

func BenchmarkArrayStack_Pop(b *testing.B) {
	stack := NewArrayStack[int]()
	for i := 0; b.Loop(); i++ {
		stack.Push(i)
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		stack.Pop()
	}
}

func BenchmarkLinkedStack_Pop(b *testing.B) {
	stack := NewLinkedStack[int]()
	for i := 0; b.Loop(); i++ {
		stack.Push(i)
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		stack.Pop()
	}
}
