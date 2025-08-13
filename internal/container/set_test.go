// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package container

import (
	"testing"
)

func TestNewMapSet(t *testing.T) {
	set := NewMapSet[int]()
	if set == nil {
		t.Error("NewMapSet() should not return nil")
	}
	if !set.IsEmpty() {
		t.Error("New set should be empty")
	}
	if set.Size() != 0 {
		t.Error("New set should have size 0")
	}
}

func TestSet_Add(t *testing.T) {
	set := NewMapSet[int]()

	// Test adding single item
	set.Add(1)
	if set.Size() != 1 {
		t.Errorf("Set size should be 1, got %d", set.Size())
	}
	if !set.Contains(1) {
		t.Error("Set should contain added item")
	}

	// Test adding multiple items
	set.Add(2)
	set.Add(3)
	if set.Size() != 3 {
		t.Errorf("Set size should be 3, got %d", set.Size())
	}

	// Test adding duplicate item
	set.Add(1)
	if set.Size() != 3 {
		t.Errorf("Set size should remain 3 after adding duplicate, got %d", set.Size())
	}
}

func TestSet_Remove(t *testing.T) {
	set := NewMapSet[int]()

	// Test removing from empty set
	set.Remove(1)
	if set.Size() != 0 {
		t.Error("Removing from empty set should not change size")
	}

	// Test removing existing item
	set.Add(1)
	set.Add(2)
	set.Remove(1)
	if set.Size() != 1 {
		t.Errorf("Set size should be 1 after removal, got %d", set.Size())
	}
	if set.Contains(1) {
		t.Error("Set should not contain removed item")
	}
	if !set.Contains(2) {
		t.Error("Set should still contain non-removed item")
	}

	// Test removing non-existent item
	set.Remove(3)
	if set.Size() != 1 {
		t.Error("Removing non-existent item should not change size")
	}
}

func TestSet_Contains(t *testing.T) {
	set := NewMapSet[int]()

	// Test empty set
	if set.Contains(1) {
		t.Error("Empty set should not contain any items")
	}

	// Test with added items
	set.Add(1)
	set.Add(2)

	if !set.Contains(1) {
		t.Error("Set should contain added item 1")
	}
	if !set.Contains(2) {
		t.Error("Set should contain added item 2")
	}
	if set.Contains(3) {
		t.Error("Set should not contain non-added item 3")
	}
}

func TestSet_IsEmpty(t *testing.T) {
	set := NewMapSet[int]()

	// Test empty set
	if !set.IsEmpty() {
		t.Error("New set should be empty")
	}

	// Test non-empty set
	set.Add(1)
	if set.IsEmpty() {
		t.Error("Set with items should not be empty")
	}

	// Test after removing all items
	set.Remove(1)
	if !set.IsEmpty() {
		t.Error("Set should be empty after removing all items")
	}
}

func TestSet_Size(t *testing.T) {
	set := NewMapSet[int]()

	// Test initial size
	if set.Size() != 0 {
		t.Errorf("Initial set size should be 0, got %d", set.Size())
	}

	// Test size after adding
	set.Add(1)
	if set.Size() != 1 {
		t.Errorf("Set size should be 1, got %d", set.Size())
	}

	set.Add(2)
	set.Add(3)
	if set.Size() != 3 {
		t.Errorf("Set size should be 3, got %d", set.Size())
	}

	// Test size after removing
	set.Remove(2)
	if set.Size() != 2 {
		t.Errorf("Set size should be 2 after removal, got %d", set.Size())
	}

	// Test size after adding duplicate
	set.Add(1)
	if set.Size() != 2 {
		t.Errorf("Set size should remain 2 after adding duplicate, got %d", set.Size())
	}
}

func TestSet_Clear(t *testing.T) {
	set := NewMapSet[int]()

	// Test clear on empty set
	set.Clear()
	if !set.IsEmpty() {
		t.Error("Cleared empty set should be empty")
	}

	// Test clear on non-empty set
	set.Add(1)
	set.Add(2)
	set.Add(3)

	set.Clear()
	if !set.IsEmpty() {
		t.Error("Cleared set should be empty")
	}
	if set.Size() != 0 {
		t.Errorf("Cleared set size should be 0, got %d", set.Size())
	}

	// Test that we can still add after clearing
	set.Add(4)
	if set.Size() != 1 {
		t.Errorf("Set size should be 1 after adding to cleared set, got %d", set.Size())
	}
}

func TestSet_Items(t *testing.T) {
	set := NewMapSet[int]()

	// Test empty set
	items := set.Items()
	if len(items) != 0 {
		t.Errorf("Empty set should return empty slice, got %v", items)
	}

	// Test with items
	set.Add(1)
	set.Add(2)
	set.Add(3)

	items = set.Items()
	if len(items) != 3 {
		t.Errorf("Set should return 3 items, got %d", len(items))
	}

	// Check that all items are present (order doesn't matter)
	expected := map[int]bool{1: true, 2: true, 3: true}
	for _, item := range items {
		if !expected[item] {
			t.Errorf("Unexpected item in set: %v", item)
		}
	}
}

func TestSet_Union(t *testing.T) {
	set1 := NewMapSet[int]()
	set2 := NewMapSet[int]()

	set1.Add(1)
	set1.Add(2)
	set2.Add(2)
	set2.Add(3)

	union := set1.Union(set2)

	if union.Size() != 3 {
		t.Errorf("Union should have 3 items, got %d", union.Size())
	}

	// Check that all items are present
	if !union.Contains(1) || !union.Contains(2) || !union.Contains(3) {
		t.Error("Union should contain all items from both sets")
	}
}

func TestSet_Intersection(t *testing.T) {
	set1 := NewMapSet[int]()
	set2 := NewMapSet[int]()

	set1.Add(1)
	set1.Add(2)
	set1.Add(3)
	set2.Add(2)
	set2.Add(3)
	set2.Add(4)

	intersection := set1.Intersection(set2)

	if intersection.Size() != 2 {
		t.Errorf("Intersection should have 2 items, got %d", intersection.Size())
	}

	// Check that only common items are present
	if !intersection.Contains(2) || !intersection.Contains(3) {
		t.Error("Intersection should contain only common items")
	}
	if intersection.Contains(1) || intersection.Contains(4) {
		t.Error("Intersection should not contain items from only one set")
	}
}

func TestSet_Difference(t *testing.T) {
	set1 := NewMapSet[int]()
	set2 := NewMapSet[int]()

	set1.Add(1)
	set1.Add(2)
	set1.Add(3)
	set2.Add(2)
	set2.Add(3)
	set2.Add(4)

	difference := set1.Difference(set2)

	if difference.Size() != 1 {
		t.Errorf("Difference should have 1 item, got %d", difference.Size())
	}

	// Check that only items in set1 but not in set2 are present
	if !difference.Contains(1) {
		t.Error("Difference should contain items only in first set")
	}
	if difference.Contains(2) || difference.Contains(3) || difference.Contains(4) {
		t.Error("Difference should not contain items in second set")
	}
}

func TestSet_Clone(t *testing.T) {
	set := NewMapSet[int]()
	set.Add(1)
	set.Add(2)
	set.Add(3)

	clone := set.Clone()

	// Test that clone has same size
	if clone.Size() != set.Size() {
		t.Errorf("Clone size should match original, got %d, want %d", clone.Size(), set.Size())
	}

	// Test that clone contains all items
	for _, item := range set.Items() {
		if !clone.Contains(item) {
			t.Errorf("Clone should contain item %v", item)
		}
	}

	// Test that original is not affected by clone modifications
	clone.Add(4)
	if set.Contains(4) {
		t.Error("Original set should not be affected by clone modifications")
	}
}

func TestSet_StringType(t *testing.T) {
	set := NewMapSet[string]()

	set.Add("hello")
	set.Add("world")

	if set.Size() != 2 {
		t.Errorf("String set should have 2 items, got %d", set.Size())
	}

	if !set.Contains("hello") || !set.Contains("world") {
		t.Error("String set should contain added items")
	}

	set.Remove("hello")
	if set.Contains("hello") {
		t.Error("String set should not contain removed item")
	}
}

func TestSet_StructType(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}

	set := NewMapSet[Person]()

	person1 := Person{Name: "Alice", Age: 30}
	person2 := Person{Name: "Bob", Age: 25}

	set.Add(person1)
	set.Add(person2)

	if set.Size() != 2 {
		t.Errorf("Struct set should have 2 items, got %d", set.Size())
	}

	if !set.Contains(person1) || !set.Contains(person2) {
		t.Error("Struct set should contain added items")
	}
}

func TestSet_PointerType(t *testing.T) {
	set := NewMapSet[*int]()

	val1 := 1
	val2 := 2

	set.Add(&val1)
	set.Add(&val2)

	if set.Size() != 2 {
		t.Errorf("Pointer set should have 2 items, got %d", set.Size())
	}

	if !set.Contains(&val1) || !set.Contains(&val2) {
		t.Error("Pointer set should contain added items")
	}
}

func TestSet_LargeNumberOfItems(t *testing.T) {
	set := NewMapSet[int]()

	// Add many items
	for i := 0; i < 1000; i++ {
		set.Add(i)
	}

	if set.Size() != 1000 {
		t.Errorf("Set should have 1000 items, got %d", set.Size())
	}

	// Check that all items are present
	for i := 0; i < 1000; i++ {
		if !set.Contains(i) {
			t.Errorf("Set should contain item %d", i)
		}
	}

	// Remove all items
	for i := 0; i < 1000; i++ {
		set.Remove(i)
	}

	if !set.IsEmpty() {
		t.Error("Set should be empty after removing all items")
	}
}

func TestSet_EmptySetOperations(t *testing.T) {
	set1 := NewMapSet[int]()
	set2 := NewMapSet[int]()

	// Test union with empty sets
	union := set1.Union(set2)
	if !union.IsEmpty() {
		t.Error("Union of empty sets should be empty")
	}

	// Test intersection with empty sets
	intersection := set1.Intersection(set2)
	if !intersection.IsEmpty() {
		t.Error("Intersection of empty sets should be empty")
	}

	// Test difference with empty sets
	difference := set1.Difference(set2)
	if !difference.IsEmpty() {
		t.Error("Difference of empty sets should be empty")
	}
}

func TestSet_DisjointSets(t *testing.T) {
	set1 := NewMapSet[int]()
	set2 := NewMapSet[int]()

	set1.Add(1)
	set1.Add(2)
	set2.Add(3)
	set2.Add(4)

	// Test intersection of disjoint sets
	intersection := set1.Intersection(set2)
	if !intersection.IsEmpty() {
		t.Error("Intersection of disjoint sets should be empty")
	}

	// Test union of disjoint sets
	union := set1.Union(set2)
	if union.Size() != 4 {
		t.Errorf("Union of disjoint sets should have 4 items, got %d", union.Size())
	}
}

func TestSet_SubsetOperations(t *testing.T) {
	set1 := NewMapSet[int]()
	set2 := NewMapSet[int]()

	set1.Add(1)
	set1.Add(2)
	set2.Add(1)
	set2.Add(2)
	set2.Add(3)

	// Test difference when set1 is subset of set2
	difference := set1.Difference(set2)
	if !difference.IsEmpty() {
		t.Error("Difference when first set is subset should be empty")
	}

	// Test intersection when set1 is subset of set2
	intersection := set1.Intersection(set2)
	if intersection.Size() != 2 {
		t.Errorf("Intersection when first set is subset should have 2 items, got %d", intersection.Size())
	}
}
