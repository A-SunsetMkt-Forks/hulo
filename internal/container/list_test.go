// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package container_test

import (
	"testing"

	"github.com/hulo-lang/hulo/internal/container"
	"github.com/stretchr/testify/assert"
)

// Test helper function to test both singly and doubly linked lists
func testListImplementation(t *testing.T, name string, newList func() container.List[int]) {
	t.Run(name, func(t *testing.T) {
		t.Run("New", func(t *testing.T) {
			list := newList()
			assert.NotNil(t, list)
			assert.True(t, list.IsEmpty())
			assert.Equal(t, 0, list.Size())
		})

		t.Run("InsertFront", func(t *testing.T) {
			list := newList()
			list.InsertFront(10)
			assert.False(t, list.IsEmpty())
			assert.Equal(t, 1, list.Size())

			list.InsertFront(20)
			assert.Equal(t, 2, list.Size())

			// Check order (last inserted should be at front)
			front, ok := list.Front()
			assert.True(t, ok)
			assert.Equal(t, 20, front)
		})

		t.Run("InsertBack", func(t *testing.T) {
			list := newList()
			list.InsertBack(10)
			assert.False(t, list.IsEmpty())
			assert.Equal(t, 1, list.Size())

			list.InsertBack(20)
			assert.Equal(t, 2, list.Size())

			// Check order (last inserted should be at back)
			back, ok := list.Back()
			assert.True(t, ok)
			assert.Equal(t, 20, back)
		})

		t.Run("InsertAt", func(t *testing.T) {
			list := newList()

			// Insert at beginning
			assert.True(t, list.InsertAt(0, 10))
			assert.Equal(t, 1, list.Size())

			// Insert at end
			assert.True(t, list.InsertAt(1, 30))
			assert.Equal(t, 2, list.Size())

			// Insert in middle
			assert.True(t, list.InsertAt(1, 20))
			assert.Equal(t, 3, list.Size())

			// Verify order: [10, 20, 30]
			val, ok := list.GetAt(0)
			assert.True(t, ok)
			assert.Equal(t, 10, val)

			val, ok = list.GetAt(1)
			assert.True(t, ok)
			assert.Equal(t, 20, val)

			val, ok = list.GetAt(2)
			assert.True(t, ok)
			assert.Equal(t, 30, val)

			// Test invalid index
			assert.False(t, list.InsertAt(-1, 40))
			assert.False(t, list.InsertAt(5, 40))
		})

		t.Run("RemoveFront", func(t *testing.T) {
			list := newList()
			list.InsertBack(10)
			list.InsertBack(20)

			item, ok := list.RemoveFront()
			assert.True(t, ok)
			assert.Equal(t, 10, item)
			assert.Equal(t, 1, list.Size())

			item, ok = list.RemoveFront()
			assert.True(t, ok)
			assert.Equal(t, 20, item)
			assert.Equal(t, 0, list.Size())
			assert.True(t, list.IsEmpty())
		})

		t.Run("RemoveFront_Empty", func(t *testing.T) {
			list := newList()
			item, ok := list.RemoveFront()
			assert.False(t, ok)
			assert.Equal(t, 0, item)
		})

		t.Run("RemoveBack", func(t *testing.T) {
			list := newList()
			list.InsertBack(10)
			list.InsertBack(20)

			item, ok := list.RemoveBack()
			assert.True(t, ok)
			assert.Equal(t, 20, item)
			assert.Equal(t, 1, list.Size())

			item, ok = list.RemoveBack()
			assert.True(t, ok)
			assert.Equal(t, 10, item)
			assert.Equal(t, 0, list.Size())
			assert.True(t, list.IsEmpty())
		})

		t.Run("RemoveBack_Empty", func(t *testing.T) {
			list := newList()
			item, ok := list.RemoveBack()
			assert.False(t, ok)
			assert.Equal(t, 0, item)
		})

		t.Run("RemoveAt", func(t *testing.T) {
			list := newList()
			list.InsertBack(10)
			list.InsertBack(20)
			list.InsertBack(30)

			// Remove from middle
			item, ok := list.RemoveAt(1)
			assert.True(t, ok)
			assert.Equal(t, 20, item)
			assert.Equal(t, 2, list.Size())

			// Remove from beginning
			item, ok = list.RemoveAt(0)
			assert.True(t, ok)
			assert.Equal(t, 10, item)
			assert.Equal(t, 1, list.Size())

			// Remove from end
			item, ok = list.RemoveAt(0)
			assert.True(t, ok)
			assert.Equal(t, 30, item)
			assert.Equal(t, 0, list.Size())
		})

		t.Run("RemoveAt_Invalid", func(t *testing.T) {
			list := newList()
			list.InsertBack(10)

			item, ok := list.RemoveAt(-1)
			assert.False(t, ok)
			assert.Equal(t, 0, item)

			item, ok = list.RemoveAt(1)
			assert.False(t, ok)
			assert.Equal(t, 0, item)
		})

		t.Run("GetAt", func(t *testing.T) {
			list := newList()
			list.InsertBack(10)
			list.InsertBack(20)
			list.InsertBack(30)

			val, ok := list.GetAt(0)
			assert.True(t, ok)
			assert.Equal(t, 10, val)

			val, ok = list.GetAt(1)
			assert.True(t, ok)
			assert.Equal(t, 20, val)

			val, ok = list.GetAt(2)
			assert.True(t, ok)
			assert.Equal(t, 30, val)
		})

		t.Run("GetAt_Invalid", func(t *testing.T) {
			list := newList()
			list.InsertBack(10)

			val, ok := list.GetAt(-1)
			assert.False(t, ok)
			assert.Equal(t, 0, val)

			val, ok = list.GetAt(1)
			assert.False(t, ok)
			assert.Equal(t, 0, val)
		})

		t.Run("SetAt", func(t *testing.T) {
			list := newList()
			list.InsertBack(10)
			list.InsertBack(20)

			assert.True(t, list.SetAt(0, 100))
			assert.True(t, list.SetAt(1, 200))

			val, ok := list.GetAt(0)
			assert.True(t, ok)
			assert.Equal(t, 100, val)

			val, ok = list.GetAt(1)
			assert.True(t, ok)
			assert.Equal(t, 200, val)
		})

		t.Run("SetAt_Invalid", func(t *testing.T) {
			list := newList()
			list.InsertBack(10)

			assert.False(t, list.SetAt(-1, 100))
			assert.False(t, list.SetAt(1, 100))
		})

		t.Run("Front", func(t *testing.T) {
			list := newList()
			list.InsertBack(10)
			list.InsertBack(20)

			val, ok := list.Front()
			assert.True(t, ok)
			assert.Equal(t, 10, val)
			assert.Equal(t, 2, list.Size()) // Should not change size
		})

		t.Run("Front_Empty", func(t *testing.T) {
			list := newList()
			val, ok := list.Front()
			assert.False(t, ok)
			assert.Equal(t, 0, val)
		})

		t.Run("Back", func(t *testing.T) {
			list := newList()
			list.InsertBack(10)
			list.InsertBack(20)

			val, ok := list.Back()
			assert.True(t, ok)
			assert.Equal(t, 20, val)
			assert.Equal(t, 2, list.Size()) // Should not change size
		})

		t.Run("Back_Empty", func(t *testing.T) {
			list := newList()
			val, ok := list.Back()
			assert.False(t, ok)
			assert.Equal(t, 0, val)
		})

		t.Run("Clear", func(t *testing.T) {
			list := newList()
			list.InsertBack(100)
			list.InsertBack(200)
			assert.False(t, list.IsEmpty())

			list.Clear()
			assert.True(t, list.IsEmpty())
			assert.Equal(t, 0, list.Size())
		})

		t.Run("Clone", func(t *testing.T) {
			list1 := newList()
			list1.InsertBack(1)
			list1.InsertBack(2)

			list2 := list1.Clone()
			assert.Equal(t, list1.Size(), list2.Size())

			// Remove from original, should not affect clone
			item1, ok1 := list1.RemoveFront()
			assert.True(t, ok1)
			assert.Equal(t, 1, item1)
			assert.Equal(t, 1, list1.Size())
			assert.Equal(t, 2, list2.Size())

			// Remove from clone, should not affect original
			item2, ok2 := list2.RemoveFront()
			assert.True(t, ok2)
			assert.Equal(t, 1, item2)
			assert.Equal(t, 1, list1.Size())
			assert.Equal(t, 1, list2.Size())
		})

		t.Run("ForEach", func(t *testing.T) {
			list := newList()
			list.InsertBack(1)
			list.InsertBack(2)
			list.InsertBack(3)

			var items []int
			list.ForEach(func(item int) bool {
				items = append(items, item)
				return true // continue iteration
			})

			assert.Equal(t, []int{1, 2, 3}, items)
		})

		t.Run("ForEach_Break", func(t *testing.T) {
			list := newList()
			list.InsertBack(1)
			list.InsertBack(2)
			list.InsertBack(3)

			var items []int
			list.ForEach(func(item int) bool {
				items = append(items, item)
				return item < 2 // break at 2
			})

			assert.Equal(t, []int{1, 2}, items)
		})

		t.Run("Complex_Operations", func(t *testing.T) {
			list := newList()

			// Build list: [5, 10, 15, 20]
			list.InsertBack(10)
			list.InsertBack(15)
			list.InsertFront(5)
			list.InsertAt(3, 20)

			assert.Equal(t, 4, list.Size())

			// Verify order
			val, ok := list.GetAt(0)
			assert.True(t, ok)
			assert.Equal(t, 5, val)

			val, ok = list.GetAt(1)
			assert.True(t, ok)
			assert.Equal(t, 10, val)

			val, ok = list.GetAt(2)
			assert.True(t, ok)
			assert.Equal(t, 15, val)

			val, ok = list.GetAt(3)
			assert.True(t, ok)
			assert.Equal(t, 20, val)

			// Remove and verify
			item, ok := list.RemoveAt(1)
			assert.True(t, ok)
			assert.Equal(t, 10, item)
			assert.Equal(t, 3, list.Size())

			// Verify new order: [5, 15, 20]
			val, ok = list.GetAt(0)
			assert.True(t, ok)
			assert.Equal(t, 5, val)

			val, ok = list.GetAt(1)
			assert.True(t, ok)
			assert.Equal(t, 15, val)

			val, ok = list.GetAt(2)
			assert.True(t, ok)
			assert.Equal(t, 20, val)
		})
	})
}

func TestSinglyLinkedList(t *testing.T) {
	testListImplementation(t, "SinglyLinkedList", container.NewSinglyLinkedList[int])
}

func TestDoublyLinkedList(t *testing.T) {
	testListImplementation(t, "DoublyLinkedList", container.NewDoublyLinkedList[int])
}

// Additional tests specific to doubly linked list advantages
func TestDoublyLinkedList_Advantages(t *testing.T) {
	t.Run("Efficient_Back_Operations", func(t *testing.T) {
		list := container.NewDoublyLinkedList[int]()

		// Test efficient back insertion
		for i := 0; i < 1000; i++ {
			list.InsertBack(i)
		}

		assert.Equal(t, 1000, list.Size())

		// Test efficient back removal
		for i := 999; i >= 0; i-- {
			item, ok := list.RemoveBack()
			assert.True(t, ok)
			assert.Equal(t, i, item)
		}

		assert.True(t, list.IsEmpty())
	})

	t.Run("Efficient_Index_Access", func(t *testing.T) {
		list := container.NewDoublyLinkedList[int]()

		// Fill list
		for i := 0; i < 100; i++ {
			list.InsertBack(i)
		}

		// Test access from both ends
		val, ok := list.GetAt(0)
		assert.True(t, ok)
		assert.Equal(t, 0, val)

		val, ok = list.GetAt(99)
		assert.True(t, ok)
		assert.Equal(t, 99, val)

		val, ok = list.GetAt(50)
		assert.True(t, ok)
		assert.Equal(t, 50, val)
	})
}

// String type tests
func TestList_StringType(t *testing.T) {
	t.Run("SinglyLinkedList_String", func(t *testing.T) {
		list := container.NewSinglyLinkedList[string]()
		list.InsertBack("hello")
		list.InsertBack("world")

		val, ok := list.GetAt(0)
		assert.True(t, ok)
		assert.Equal(t, "hello", val)

		val, ok = list.GetAt(1)
		assert.True(t, ok)
		assert.Equal(t, "world", val)
	})

	t.Run("DoublyLinkedList_String", func(t *testing.T) {
		list := container.NewDoublyLinkedList[string]()
		list.InsertBack("hello")
		list.InsertBack("world")

		val, ok := list.GetAt(0)
		assert.True(t, ok)
		assert.Equal(t, "hello", val)

		val, ok = list.GetAt(1)
		assert.True(t, ok)
		assert.Equal(t, "world", val)
	})
}

// Struct type tests
func TestList_StructType(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}

	t.Run("SinglyLinkedList_Struct", func(t *testing.T) {
		list := container.NewSinglyLinkedList[Person]()
		list.InsertBack(Person{Name: "Alice", Age: 30})
		list.InsertBack(Person{Name: "Bob", Age: 25})

		val, ok := list.GetAt(0)
		assert.True(t, ok)
		assert.Equal(t, "Alice", val.Name)
		assert.Equal(t, 30, val.Age)

		val, ok = list.GetAt(1)
		assert.True(t, ok)
		assert.Equal(t, "Bob", val.Name)
		assert.Equal(t, 25, val.Age)
	})

	t.Run("DoublyLinkedList_Struct", func(t *testing.T) {
		list := container.NewDoublyLinkedList[Person]()
		list.InsertBack(Person{Name: "Alice", Age: 30})
		list.InsertBack(Person{Name: "Bob", Age: 25})

		val, ok := list.GetAt(0)
		assert.True(t, ok)
		assert.Equal(t, "Alice", val.Name)
		assert.Equal(t, 30, val.Age)

		val, ok = list.GetAt(1)
		assert.True(t, ok)
		assert.Equal(t, "Bob", val.Name)
		assert.Equal(t, 25, val.Age)
	})
}

// Buffer operations tests
func TestList_BufferOperations(t *testing.T) {
	t.Run("SinglyLinkedList_Buffer", func(t *testing.T) {
		testBufferOperations(t, container.NewSinglyLinkedList[string])
	})

	t.Run("DoublyLinkedList_Buffer", func(t *testing.T) {
		testBufferOperations(t, container.NewDoublyLinkedList[string])
	})
}

func testBufferOperations(t *testing.T, newList func() container.List[string]) {
	t.Run("Append", func(t *testing.T) {
		list := newList()
		list.Append("a", "b", "c")

		assert.Equal(t, 3, list.Size())

		val, ok := list.GetAt(0)
		assert.True(t, ok)
		assert.Equal(t, "a", val)

		val, ok = list.GetAt(1)
		assert.True(t, ok)
		assert.Equal(t, "b", val)

		val, ok = list.GetAt(2)
		assert.True(t, ok)
		assert.Equal(t, "c", val)
	})

	t.Run("ToSlice", func(t *testing.T) {
		list := newList()
		list.Append("x", "y", "z")

		slice := list.ToSlice()
		assert.Equal(t, []string{"x", "y", "z"}, slice)
		assert.Equal(t, 3, len(slice))
	})

	t.Run("AppendList", func(t *testing.T) {
		list1 := newList()
		list1.Append("a", "b")

		list2 := newList()
		list2.Append("c", "d")

		list1.AppendList(list2)

		assert.Equal(t, 4, list1.Size())
		slice := list1.ToSlice()
		assert.Equal(t, []string{"a", "b", "c", "d"}, slice)
	})

	t.Run("AppendList_Empty", func(t *testing.T) {
		list1 := newList()
		list1.Append("a", "b")

		list2 := newList() // 空列表

		list1.AppendList(list2)

		assert.Equal(t, 2, list1.Size())
		slice := list1.ToSlice()
		assert.Equal(t, []string{"a", "b"}, slice)
	})

	t.Run("AppendList_ToEmpty", func(t *testing.T) {
		list1 := newList() // 空列表

		list2 := newList()
		list2.Append("c", "d")

		list1.AppendList(list2)

		assert.Equal(t, 2, list1.Size())
		slice := list1.ToSlice()
		assert.Equal(t, []string{"c", "d"}, slice)
	})

	t.Run("BufferWorkflow", func(t *testing.T) {
		list := newList()

		// 模拟buffer工作流程
		list.Append("stmt1", "stmt2")
		assert.Equal(t, 2, list.Size())

		// 插入语句
		list.InsertAt(1, "inserted")
		assert.Equal(t, 3, list.Size())

		// 转换为slice
		result := list.ToSlice()
		assert.Equal(t, []string{"stmt1", "inserted", "stmt2"}, result)

		// 清空
		list.Clear()
		assert.True(t, list.IsEmpty())
		assert.Equal(t, 0, list.Size())
	})
}
