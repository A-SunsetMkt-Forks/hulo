// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package container

import (
	"testing"
)

func BenchmarkNewMapSet(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		NewMapSet[int]()
	}
}

func BenchmarkSet_Add(b *testing.B) {
	set := NewMapSet[int]()
	b.ReportAllocs()


	for i := 0; b.Loop(); i++ {
		set.Add(i)
	}
}

func BenchmarkSet_AddDuplicate(b *testing.B) {
	set := NewMapSet[int]()
	set.Add(1)
	b.ReportAllocs()


	for b.Loop() {
		set.Add(1)
	}
}

func BenchmarkSet_Remove(b *testing.B) {
	set := NewMapSet[int]()
	// Pre-fill the set
	for i := 0; b.Loop(); i++ {
		set.Add(i)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		set.Remove(i)
	}
}

func BenchmarkSet_RemoveNonExistent(b *testing.B) {
	set := NewMapSet[int]()
	set.Add(1)
	b.ReportAllocs()


	for i := 0; b.Loop(); i++ {
		set.Remove(i + 1000)
	}
}

func BenchmarkSet_Contains(b *testing.B) {
	set := NewMapSet[int]()
	set.Add(1)
	set.Add(2)
	set.Add(3)
	b.ReportAllocs()


	for i := 0; b.Loop(); i++ {
		set.Contains(i % 4)
	}
}

func BenchmarkSet_IsEmpty(b *testing.B) {
	set := NewMapSet[int]()
	set.Add(1)
	b.ReportAllocs()


	for b.Loop() {
		set.IsEmpty()
	}
}

func BenchmarkSet_Size(b *testing.B) {
	set := NewMapSet[int]()
	set.Add(1)
	set.Add(2)
	set.Add(3)
	b.ReportAllocs()


	for b.Loop() {
		set.Size()
	}
}

func BenchmarkSet_Clear(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		set := NewMapSet[int]()
		set.Add(1)
		set.Add(2)
		set.Add(3)
		set.Clear()
	}
}

func BenchmarkSet_Items(b *testing.B) {
	set := NewMapSet[int]()
	set.Add(1)
	set.Add(2)
	set.Add(3)
	set.Add(4)
	set.Add(5)
	b.ReportAllocs()


	for b.Loop() {
		set.Items()
	}
}

func BenchmarkSet_Union(b *testing.B) {
	set1 := NewMapSet[int]()
	set2 := NewMapSet[int]()

	set1.Add(1)
	set1.Add(2)
	set1.Add(3)
	set2.Add(2)
	set2.Add(3)
	set2.Add(4)

	b.ReportAllocs()


	for b.Loop() {
		set1.Union(set2)
	}
}

func BenchmarkSet_Intersection(b *testing.B) {
	set1 := NewMapSet[int]()
	set2 := NewMapSet[int]()

	set1.Add(1)
	set1.Add(2)
	set1.Add(3)
	set2.Add(2)
	set2.Add(3)
	set2.Add(4)

	b.ReportAllocs()


	for b.Loop() {
		set1.Intersection(set2)
	}
}

func BenchmarkSet_Difference(b *testing.B) {
	set1 := NewMapSet[int]()
	set2 := NewMapSet[int]()

	set1.Add(1)
	set1.Add(2)
	set1.Add(3)
	set2.Add(2)
	set2.Add(3)
	set2.Add(4)

	b.ReportAllocs()


	for b.Loop() {
		set1.Difference(set2)
	}
}

func BenchmarkSet_Clone(b *testing.B) {
	set := NewMapSet[int]()
	set.Add(1)
	set.Add(2)
	set.Add(3)
	b.ReportAllocs()


	for b.Loop() {
		set.Clone()
	}
}

func BenchmarkSet_AddRemove(b *testing.B) {
	set := NewMapSet[int]()
	b.ReportAllocs()


	for i := 0; b.Loop(); i++ {
		set.Add(i)
		if i%2 == 0 {
			set.Remove(i)
		}
	}
}

func BenchmarkSet_LargeSet(b *testing.B) {
	set := NewMapSet[int]()
	// Pre-fill with large number of items
	for i := range 10000 {
		set.Add(i)
	}

	b.ReportAllocs()


	for i := 0; b.Loop(); i++ {
		set.Add(i)
		set.Remove(i)
	}
}

func BenchmarkSet_StringType(b *testing.B) {
	set := NewMapSet[string]()
	b.ReportAllocs()


	for b.Loop() {
		set.Add("test")
	}
}

func BenchmarkSet_StructType(b *testing.B) {
	type Person struct {
		Name string
		Age  int
	}

	set := NewMapSet[Person]()
	person := Person{Name: "Alice", Age: 30}

	b.ReportAllocs()


	for b.Loop() {
		set.Add(person)
	}
}

func BenchmarkSet_PointerType(b *testing.B) {
	set := NewMapSet[*int]()
	val := 42

	b.ReportAllocs()


	for b.Loop() {
		set.Add(&val)
	}
}

func BenchmarkSet_MixedOperations(b *testing.B) {
	set := NewMapSet[int]()
	b.ReportAllocs()


	for i := 0; b.Loop(); i++ {
		set.Add(i)
		set.Contains(i)
		if i%3 == 0 {
			set.Remove(i)
		}
		set.Size()
		if i%5 == 0 {
			set.IsEmpty()
		}
	}
}

func BenchmarkSet_GrowAndShrink(b *testing.B) {
	set := NewMapSet[int]()
	b.ReportAllocs()


	for b.Loop() {
		// Grow the set
		for j := range 100 {
			set.Add(j)
		}
		// Shrink the set
		for j := range 100 {
			set.Remove(j)
		}
	}
}

func BenchmarkSet_CloneLargeSet(b *testing.B) {
	set := NewMapSet[int]()
	// Pre-fill with large number of items
	for i := range 1000 {
		set.Add(i)
	}

	b.ReportAllocs()


	for b.Loop() {
		set.Clone()
	}
}

func BenchmarkSet_ClearLargeSet(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		set := NewMapSet[int]()
		// Pre-fill with large number of items
		for j := range 1000 {
			set.Add(j)
		}
		set.Clear()
	}
}

func BenchmarkSet_EmptySetOperations(b *testing.B) {
	set := NewMapSet[int]()
	b.ReportAllocs()


	for i := 0; b.Loop(); i++ {
		set.IsEmpty()
		set.Size()
		set.Contains(i)
		set.Remove(i)
	}
}

func BenchmarkSet_SingleItemOperations(b *testing.B) {
	set := NewMapSet[int]()
	set.Add(1)

	b.ReportAllocs()


	for b.Loop() {
		set.IsEmpty()
		set.Size()
		set.Contains(1)
	}
}

func BenchmarkSet_DisjointSetOperations(b *testing.B) {
	set1 := NewMapSet[int]()
	set2 := NewMapSet[int]()

	// Create disjoint sets
	for i := range 100 {
		set1.Add(i)
		set2.Add(i + 100)
	}

	b.ReportAllocs()


	for b.Loop() {
		set1.Union(set2)
		set1.Intersection(set2)
		set1.Difference(set2)
	}
}

func BenchmarkSet_OverlappingSetOperations(b *testing.B) {
	set1 := NewMapSet[int]()
	set2 := NewMapSet[int]()

	// Create overlapping sets
	for i := range 100 {
		set1.Add(i)
		if i%2 == 0 {
			set2.Add(i)
		}
	}

	b.ReportAllocs()


	for b.Loop() {
		set1.Union(set2)
		set1.Intersection(set2)
		set1.Difference(set2)
	}
}

func BenchmarkSet_ItemsLargeSet(b *testing.B) {
	set := NewMapSet[int]()
	// Pre-fill with large number of items
	for i := range 1000 {
		set.Add(i)
	}

	b.ReportAllocs()


	for b.Loop() {
		set.Items()
	}
}

func BenchmarkSet_AlternatingAddRemove(b *testing.B) {
	set := NewMapSet[int]()
	b.ReportAllocs()


	for i := 0; b.Loop(); i++ {
		if i%2 == 0 {
			set.Add(i)
		} else {
			set.Remove(i - 1)
		}
	}
}
