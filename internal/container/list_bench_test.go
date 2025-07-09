// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package container

import (
	"fmt"
	"testing"
)

func BenchmarkNewSinglyLinkedList(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		NewSinglyLinkedList[int]()
	}
}

func BenchmarkNewDoublyLinkedList(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		NewDoublyLinkedList[int]()
	}
}

func BenchmarkSinglyLinkedList_InsertFront(b *testing.B) {
	list := NewSinglyLinkedList[int]()
	b.ReportAllocs()

	for i := 0; b.Loop(); i++ {
		list.InsertFront(i)
	}
}

func BenchmarkDoublyLinkedList_InsertFront(b *testing.B) {
	list := NewDoublyLinkedList[int]()
	b.ReportAllocs()

	for i := 0; b.Loop(); i++ {
		list.InsertFront(i)
	}
}

func BenchmarkSinglyLinkedList_InsertBack(b *testing.B) {
	list := NewSinglyLinkedList[int]()
	b.ReportAllocs()

	for i := 0; b.Loop(); i++ {
		list.InsertBack(i)
	}
}

func BenchmarkDoublyLinkedList_InsertBack(b *testing.B) {
	list := NewDoublyLinkedList[int]()
	b.ReportAllocs()

	for i := 0; b.Loop(); i++ {
		list.InsertBack(i)
	}
}

func BenchmarkSinglyLinkedList_InsertAt(b *testing.B) {
	list := NewSinglyLinkedList[int]()
	b.ReportAllocs()

	for i := 0; b.Loop(); i++ {
		list.InsertAt(i%10, i)
	}
}

func BenchmarkDoublyLinkedList_InsertAt(b *testing.B) {
	list := NewDoublyLinkedList[int]()
	b.ReportAllocs()

	for i := 0; b.Loop(); i++ {
		list.InsertAt(i%10, i)
	}
}

func BenchmarkSinglyLinkedList_RemoveFront(b *testing.B) {
	list := NewSinglyLinkedList[int]()
	// Pre-fill the list
	for i := 0; i < b.N; i++ {
		list.InsertBack(i)
	}

	b.ReportAllocs()

	for b.Loop() {
		list.RemoveFront()
	}
}

func BenchmarkDoublyLinkedList_RemoveFront(b *testing.B) {
	list := NewDoublyLinkedList[int]()
	// Pre-fill the list
	for i := 0; i < b.N; i++ {
		list.InsertBack(i)
	}

	b.ReportAllocs()

	for b.Loop() {
		list.RemoveFront()
	}
}

func BenchmarkSinglyLinkedList_RemoveBack(b *testing.B) {
	list := NewSinglyLinkedList[int]()
	// Pre-fill the list
	for i := 0; i < b.N; i++ {
		list.InsertBack(i)
	}

	b.ReportAllocs()

	for b.Loop() {
		list.RemoveBack()
	}
}

func BenchmarkDoublyLinkedList_RemoveBack(b *testing.B) {
	list := NewDoublyLinkedList[int]()
	// Pre-fill the list
	for i := 0; i < b.N; i++ {
		list.InsertBack(i)
	}

	b.ReportAllocs()

	for b.Loop() {
		list.RemoveBack()
	}
}

func BenchmarkSinglyLinkedList_RemoveAt(b *testing.B) {
	list := NewSinglyLinkedList[int]()
	// Pre-fill the list
	for i := 0; i < b.N; i++ {
		list.InsertBack(i)
	}

	b.ReportAllocs()

	for i := 0; b.Loop(); i++ {
		list.RemoveAt(i % list.Size())
	}
}

func BenchmarkDoublyLinkedList_RemoveAt(b *testing.B) {
	list := NewDoublyLinkedList[int]()
	// Pre-fill the list
	for i := 0; i < b.N; i++ {
		list.InsertBack(i)
	}

	b.ReportAllocs()

	for i := 0; b.Loop(); i++ {
		list.RemoveAt(i % list.Size())
	}
}

func BenchmarkSinglyLinkedList_GetAt(b *testing.B) {
	list := NewSinglyLinkedList[int]()
	// Pre-fill the list
	for i := 0; i < 1000; i++ {
		list.InsertBack(i)
	}

	b.ReportAllocs()

	for i := 0; b.Loop(); i++ {
		list.GetAt(i % list.Size())
	}
}

func BenchmarkDoublyLinkedList_GetAt(b *testing.B) {
	list := NewDoublyLinkedList[int]()
	// Pre-fill the list
	for i := 0; i < 1000; i++ {
		list.InsertBack(i)
	}

	b.ReportAllocs()

	for i := 0; b.Loop(); i++ {
		list.GetAt(i % list.Size())
	}
}

func BenchmarkSinglyLinkedList_SetAt(b *testing.B) {
	list := NewSinglyLinkedList[int]()
	// Pre-fill the list
	for i := 0; i < 1000; i++ {
		list.InsertBack(i)
	}

	b.ReportAllocs()

	for i := 0; b.Loop(); i++ {
		list.SetAt(i%list.Size(), i)
	}
}

func BenchmarkDoublyLinkedList_SetAt(b *testing.B) {
	list := NewDoublyLinkedList[int]()
	// Pre-fill the list
	for i := 0; i < 1000; i++ {
		list.InsertBack(i)
	}

	b.ReportAllocs()

	for i := 0; b.Loop(); i++ {
		list.SetAt(i%list.Size(), i)
	}
}

func BenchmarkSinglyLinkedList_Front(b *testing.B) {
	list := NewSinglyLinkedList[int]()
	list.InsertBack(1)
	list.InsertBack(2)
	list.InsertBack(3)

	b.ReportAllocs()

	for b.Loop() {
		list.Front()
	}
}

func BenchmarkDoublyLinkedList_Front(b *testing.B) {
	list := NewDoublyLinkedList[int]()
	list.InsertBack(1)
	list.InsertBack(2)
	list.InsertBack(3)

	b.ReportAllocs()

	for b.Loop() {
		list.Front()
	}
}

func BenchmarkSinglyLinkedList_Back(b *testing.B) {
	list := NewSinglyLinkedList[int]()
	list.InsertBack(1)
	list.InsertBack(2)
	list.InsertBack(3)

	b.ReportAllocs()

	for b.Loop() {
		list.Back()
	}
}

func BenchmarkDoublyLinkedList_Back(b *testing.B) {
	list := NewDoublyLinkedList[int]()
	list.InsertBack(1)
	list.InsertBack(2)
	list.InsertBack(3)

	b.ReportAllocs()

	for b.Loop() {
		list.Back()
	}
}

func BenchmarkSinglyLinkedList_IsEmpty(b *testing.B) {
	list := NewSinglyLinkedList[int]()
	list.InsertBack(1)

	b.ReportAllocs()

	for b.Loop() {
		list.IsEmpty()
	}
}

func BenchmarkDoublyLinkedList_IsEmpty(b *testing.B) {
	list := NewDoublyLinkedList[int]()
	list.InsertBack(1)

	b.ReportAllocs()

	for b.Loop() {
		list.IsEmpty()
	}
}

func BenchmarkSinglyLinkedList_Size(b *testing.B) {
	list := NewSinglyLinkedList[int]()
	list.InsertBack(1)
	list.InsertBack(2)
	list.InsertBack(3)

	b.ReportAllocs()

	for b.Loop() {
		list.Size()
	}
}

func BenchmarkDoublyLinkedList_Size(b *testing.B) {
	list := NewDoublyLinkedList[int]()
	list.InsertBack(1)
	list.InsertBack(2)
	list.InsertBack(3)

	b.ReportAllocs()

	for b.Loop() {
		list.Size()
	}
}

func BenchmarkSinglyLinkedList_Clear(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		list := NewSinglyLinkedList[int]()
		list.InsertBack(1)
		list.InsertBack(2)
		list.InsertBack(3)
		list.Clear()
	}
}

func BenchmarkDoublyLinkedList_Clear(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		list := NewDoublyLinkedList[int]()
		list.InsertBack(1)
		list.InsertBack(2)
		list.InsertBack(3)
		list.Clear()
	}
}

func BenchmarkSinglyLinkedList_Clone(b *testing.B) {
	list := NewSinglyLinkedList[int]()
	list.InsertBack(1)
	list.InsertBack(2)
	list.InsertBack(3)

	b.ReportAllocs()

	for b.Loop() {
		list.Clone()
	}
}

func BenchmarkDoublyLinkedList_Clone(b *testing.B) {
	list := NewDoublyLinkedList[int]()
	list.InsertBack(1)
	list.InsertBack(2)
	list.InsertBack(3)

	b.ReportAllocs()

	for b.Loop() {
		list.Clone()
	}
}

func BenchmarkSinglyLinkedList_ForEach(b *testing.B) {
	list := NewSinglyLinkedList[int]()
	for i := 0; i < 100; i++ {
		list.InsertBack(i)
	}

	b.ReportAllocs()

	for b.Loop() {
		list.ForEach(func(item int) bool {
			return true
		})
	}
}

func BenchmarkDoublyLinkedList_ForEach(b *testing.B) {
	list := NewDoublyLinkedList[int]()
	for i := 0; i < 100; i++ {
		list.InsertBack(i)
	}

	b.ReportAllocs()

	for b.Loop() {
		list.ForEach(func(item int) bool {
			return true
		})
	}
}

func BenchmarkSinglyLinkedList_InsertFrontRemoveFront(b *testing.B) {
	list := NewSinglyLinkedList[int]()
	b.ReportAllocs()

	for i := 0; b.Loop(); i++ {
		list.InsertFront(i)
		if i%2 == 0 {
			list.RemoveFront()
		}
	}
}

func BenchmarkDoublyLinkedList_InsertFrontRemoveFront(b *testing.B) {
	list := NewDoublyLinkedList[int]()
	b.ReportAllocs()

	for i := 0; b.Loop(); i++ {
		list.InsertFront(i)
		if i%2 == 0 {
			list.RemoveFront()
		}
	}
}

func BenchmarkSinglyLinkedList_InsertBackRemoveBack(b *testing.B) {
	list := NewSinglyLinkedList[int]()
	b.ReportAllocs()

	for i := 0; b.Loop(); i++ {
		list.InsertBack(i)
		if i%2 == 0 {
			list.RemoveBack()
		}
	}
}

func BenchmarkDoublyLinkedList_InsertBackRemoveBack(b *testing.B) {
	list := NewDoublyLinkedList[int]()
	b.ReportAllocs()

	for i := 0; b.Loop(); i++ {
		list.InsertBack(i)
		if i%2 == 0 {
			list.RemoveBack()
		}
	}
}

func BenchmarkSinglyLinkedList_LargeList(b *testing.B) {
	list := NewSinglyLinkedList[int]()
	// Pre-fill with large number of items
	for i := range 10000 {
		list.InsertBack(i)
	}

	b.ReportAllocs()

	for i := 0; b.Loop(); i++ {
		list.InsertAt(i%100, i)
		list.RemoveAt(i % list.Size())
	}
}

func BenchmarkDoublyLinkedList_LargeList(b *testing.B) {
	list := NewDoublyLinkedList[int]()
	// Pre-fill with large number of items
	for i := range 10000 {
		list.InsertBack(i)
	}

	b.ReportAllocs()

	for i := 0; b.Loop(); i++ {
		list.InsertAt(i%100, i)
		list.RemoveAt(i % list.Size())
	}
}

func BenchmarkSinglyLinkedList_StringType(b *testing.B) {
	list := NewSinglyLinkedList[string]()
	b.ReportAllocs()

	for b.Loop() {
		list.InsertBack("test")
	}
}

func BenchmarkDoublyLinkedList_StringType(b *testing.B) {
	list := NewDoublyLinkedList[string]()
	b.ReportAllocs()

	for b.Loop() {
		list.InsertBack("test")
	}
}

func BenchmarkSinglyLinkedList_StructType(b *testing.B) {
	type Person struct {
		Name string
		Age  int
	}

	list := NewSinglyLinkedList[Person]()
	person := Person{Name: "Alice", Age: 30}

	b.ReportAllocs()

	for b.Loop() {
		list.InsertBack(person)
	}
}

func BenchmarkDoublyLinkedList_StructType(b *testing.B) {
	type Person struct {
		Name string
		Age  int
	}

	list := NewDoublyLinkedList[Person]()
	person := Person{Name: "Alice", Age: 30}

	b.ReportAllocs()

	for b.Loop() {
		list.InsertBack(person)
	}
}

func BenchmarkSinglyLinkedList_MixedOperations(b *testing.B) {
	list := NewSinglyLinkedList[int]()
	b.ReportAllocs()

	for i := 0; b.Loop(); i++ {
		list.InsertFront(i)
		list.InsertBack(i)
		list.Front()
		list.Back()
		if i%3 == 0 {
			list.RemoveFront()
		}
		if i%5 == 0 {
			list.RemoveBack()
		}
		list.Size()
		if i%7 == 0 {
			list.IsEmpty()
		}
	}
}

func BenchmarkDoublyLinkedList_MixedOperations(b *testing.B) {
	list := NewDoublyLinkedList[int]()
	b.ReportAllocs()

	for i := 0; b.Loop(); i++ {
		list.InsertFront(i)
		list.InsertBack(i)
		list.Front()
		list.Back()
		if i%3 == 0 {
			list.RemoveFront()
		}
		if i%5 == 0 {
			list.RemoveBack()
		}
		list.Size()
		if i%7 == 0 {
			list.IsEmpty()
		}
	}
}

func BenchmarkSinglyLinkedList_GrowAndShrink(b *testing.B) {
	list := NewSinglyLinkedList[int]()
	b.ReportAllocs()

	for b.Loop() {
		// Grow the list
		for j := range 100 {
			list.InsertBack(j)
		}
		// Shrink the list
		for range 100 {
			list.RemoveBack()
		}
	}
}

func BenchmarkDoublyLinkedList_GrowAndShrink(b *testing.B) {
	list := NewDoublyLinkedList[int]()
	b.ReportAllocs()

	for b.Loop() {
		// Grow the list
		for j := range 100 {
			list.InsertBack(j)
		}
		// Shrink the list
		for range 100 {
			list.RemoveBack()
		}
	}
}

func BenchmarkSinglyLinkedList_CloneLargeList(b *testing.B) {
	list := NewSinglyLinkedList[int]()
	// Pre-fill with large number of items
	for i := range 1000 {
		list.InsertBack(i)
	}

	b.ReportAllocs()

	for b.Loop() {
		list.Clone()
	}
}

func BenchmarkDoublyLinkedList_CloneLargeList(b *testing.B) {
	list := NewDoublyLinkedList[int]()
	// Pre-fill with large number of items
	for i := range 1000 {
		list.InsertBack(i)
	}

	b.ReportAllocs()

	for b.Loop() {
		list.Clone()
	}
}

func BenchmarkSinglyLinkedList_ClearLargeList(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		list := NewSinglyLinkedList[int]()
		// Pre-fill with large number of items
		for j := range 1000 {
			list.InsertBack(j)
		}
		list.Clear()
	}
}

func BenchmarkDoublyLinkedList_ClearLargeList(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		list := NewDoublyLinkedList[int]()
		// Pre-fill with large number of items
		for j := range 1000 {
			list.InsertBack(j)
		}
		list.Clear()
	}
}

func BenchmarkSinglyLinkedList_EmptyOperations(b *testing.B) {
	list := NewSinglyLinkedList[int]()
	b.ReportAllocs()

	for b.Loop() {
		list.IsEmpty()
		list.Size()
		list.Front()
		list.Back()
		list.RemoveFront()
		list.RemoveBack()
	}
}

func BenchmarkDoublyLinkedList_EmptyOperations(b *testing.B) {
	list := NewDoublyLinkedList[int]()
	b.ReportAllocs()

	for b.Loop() {
		list.IsEmpty()
		list.Size()
		list.Front()
		list.Back()
		list.RemoveFront()
		list.RemoveBack()
	}
}

func BenchmarkSinglyLinkedList_SingleItemOperations(b *testing.B) {
	list := NewSinglyLinkedList[int]()
	list.InsertBack(1)

	b.ReportAllocs()

	for b.Loop() {
		list.IsEmpty()
		list.Size()
		list.Front()
		list.Back()
	}
}

func BenchmarkDoublyLinkedList_SingleItemOperations(b *testing.B) {
	list := NewDoublyLinkedList[int]()
	list.InsertBack(1)

	b.ReportAllocs()

	for b.Loop() {
		list.IsEmpty()
		list.Size()
		list.Front()
		list.Back()
	}
}

func BenchmarkSinglyLinkedList_AlternatingInsertRemove(b *testing.B) {
	list := NewSinglyLinkedList[int]()
	b.ReportAllocs()

	for i := 0; b.Loop(); i++ {
		if i%2 == 0 {
			list.InsertFront(i)
		} else {
			list.RemoveFront()
		}
	}
}

func BenchmarkDoublyLinkedList_AlternatingInsertRemove(b *testing.B) {
	list := NewDoublyLinkedList[int]()
	b.ReportAllocs()

	for i := 0; b.Loop(); i++ {
		if i%2 == 0 {
			list.InsertFront(i)
		} else {
			list.RemoveFront()
		}
	}
}

func BenchmarkSinglyLinkedList_ReverseOrderInsert(b *testing.B) {
	list := NewSinglyLinkedList[int]()
	b.ReportAllocs()

	for b.Loop() {
		list.Clear()
		// Insert in reverse order
		for i := 999; i >= 0; i-- {
			list.InsertFront(i)
		}
	}
}

func BenchmarkDoublyLinkedList_ReverseOrderInsert(b *testing.B) {
	list := NewDoublyLinkedList[int]()
	b.ReportAllocs()

	for b.Loop() {
		list.Clear()
		// Insert in reverse order
		for i := 999; i >= 0; i-- {
			list.InsertFront(i)
		}
	}
}

func BenchmarkSinglyLinkedList_SortedOrderInsert(b *testing.B) {
	list := NewSinglyLinkedList[int]()
	b.ReportAllocs()

	for b.Loop() {
		list.Clear()
		// Insert in sorted order
		for i := range 1000 {
			list.InsertBack(i)
		}
	}
}

func BenchmarkDoublyLinkedList_SortedOrderInsert(b *testing.B) {
	list := NewDoublyLinkedList[int]()
	b.ReportAllocs()

	for b.Loop() {
		list.Clear()
		// Insert in sorted order
		for i := range 1000 {
			list.InsertBack(i)
		}
	}
}

func BenchmarkSinglyLinkedList_RandomOrderInsert(b *testing.B) {
	list := NewSinglyLinkedList[int]()
	b.ReportAllocs()

	for b.Loop() {
		list.Clear()
		// Insert in pseudo-random order
		for i := range 1000 {
			size := list.Size()
			if size == 0 {
				list.InsertAt(0, i)
			} else {
				list.InsertAt((i*7)%size, i)
			}
		}
	}
}

func BenchmarkDoublyLinkedList_RandomOrderInsert(b *testing.B) {
	list := NewDoublyLinkedList[int]()
	b.ReportAllocs()

	for b.Loop() {
		list.Clear()
		// Insert in pseudo-random order
		for i := range 1000 {
			size := list.Size()
			if size == 0 {
				list.InsertAt(0, i)
			} else {
				list.InsertAt((i*7)%size, i)
			}
		}
	}
}

func BenchmarkSinglyLinkedList_ExtractAll(b *testing.B) {
	list := NewSinglyLinkedList[int]()
	// Pre-fill list
	for i := range 1000 {
		list.InsertBack(i)
	}

	b.ReportAllocs()

	for b.Loop() {
		// Extract all items
		for !list.IsEmpty() {
			list.RemoveFront()
		}
		// Re-fill for next iteration
		for i := range 1000 {
			list.InsertBack(i)
		}
	}
}

func BenchmarkDoublyLinkedList_ExtractAll(b *testing.B) {
	list := NewDoublyLinkedList[int]()
	// Pre-fill list
	for i := range 1000 {
		list.InsertBack(i)
	}

	b.ReportAllocs()

	for b.Loop() {
		// Extract all items
		for !list.IsEmpty() {
			list.RemoveFront()
		}
		// Re-fill for next iteration
		for i := range 1000 {
			list.InsertBack(i)
		}
	}
}

func BenchmarkSinglyLinkedList_Append(b *testing.B) {
	list := NewSinglyLinkedList[string]()
	b.ReportAllocs()

	for b.Loop() {
		list.Append("a", "b", "c", "d", "e")
	}
}

func BenchmarkDoublyLinkedList_Append(b *testing.B) {
	list := NewDoublyLinkedList[string]()
	b.ReportAllocs()

	for b.Loop() {
		list.Append("a", "b", "c", "d", "e")
	}
}

func BenchmarkSinglyLinkedList_ToSlice(b *testing.B) {
	list := NewSinglyLinkedList[string]()
	// Pre-fill the list
	for i := 0; i < 1000; i++ {
		list.InsertBack(fmt.Sprintf("item_%d", i))
	}

	b.ReportAllocs()

	for b.Loop() {
		list.ToSlice()
	}
}

func BenchmarkDoublyLinkedList_ToSlice(b *testing.B) {
	list := NewDoublyLinkedList[string]()
	// Pre-fill the list
	for i := 0; i < 1000; i++ {
		list.InsertBack(fmt.Sprintf("item_%d", i))
	}

	b.ReportAllocs()

	for b.Loop() {
		list.ToSlice()
	}
}

func BenchmarkSinglyLinkedList_AppendList(b *testing.B) {
	list1 := NewSinglyLinkedList[string]()
	list2 := NewSinglyLinkedList[string]()

	// Pre-fill list2
	for i := 0; i < 100; i++ {
		list2.InsertBack(fmt.Sprintf("item_%d", i))
	}

	b.ReportAllocs()

	for b.Loop() {
		list1.Clear()
		list1.AppendList(list2)
	}
}

func BenchmarkDoublyLinkedList_AppendList(b *testing.B) {
	list1 := NewDoublyLinkedList[string]()
	list2 := NewDoublyLinkedList[string]()

	// Pre-fill list2
	for i := 0; i < 100; i++ {
		list2.InsertBack(fmt.Sprintf("item_%d", i))
	}

	b.ReportAllocs()

	for b.Loop() {
		list1.Clear()
		list1.AppendList(list2)
	}
}

func BenchmarkSinglyLinkedList_BufferWorkflow(b *testing.B) {
	list := NewSinglyLinkedList[string]()
	b.ReportAllocs()

	for b.Loop() {
		list.Clear()
		list.Append("stmt1", "stmt2", "stmt3")
		list.InsertAt(1, "inserted")
		result := list.ToSlice()
		_ = result
	}
}

func BenchmarkDoublyLinkedList_BufferWorkflow(b *testing.B) {
	list := NewDoublyLinkedList[string]()
	b.ReportAllocs()

	for b.Loop() {
		list.Clear()
		list.Append("stmt1", "stmt2", "stmt3")
		list.InsertAt(1, "inserted")
		result := list.ToSlice()
		_ = result
	}
}
