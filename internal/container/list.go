// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package container

// List defines the interface for linked list operations.
type List[T any] interface {
	// Basic operations
	InsertFront(item T)
	InsertBack(item T)
	InsertAt(index int, item T) bool
	RemoveFront() (T, bool)
	RemoveBack() (T, bool)
	RemoveAt(index int) (T, bool)
	GetAt(index int) (T, bool)
	SetAt(index int, item T) bool

	// Utility operations
	IsEmpty() bool
	Size() int
	Clear()
	Clone() List[T]

	// Iterator operations
	Front() (T, bool)
	Back() (T, bool)
	ForEach(fn func(item T) bool)

	// Slice operations
	Append(items ...T)
	ToSlice() []T
	AppendList(other List[T])
}

// Node represents a node in a singly linked list.
type Node[T any] struct {
	Value T
	Next  *Node[T]
}

// SinglyLinkedList represents a singly linked list implementation.
type SinglyLinkedList[T any] struct {
	head *Node[T]
	size int
}

// Ensure SinglyLinkedList implements List interface.
var _ List[any] = (*SinglyLinkedList[any])(nil)

// NewSinglyLinkedList creates a new empty singly linked list.
func NewSinglyLinkedList[T any]() List[T] {
	return &SinglyLinkedList[T]{
		head: nil,
		size: 0,
	}
}

// InsertFront adds an item to the front of the list.
func (l *SinglyLinkedList[T]) InsertFront(item T) {
	newNode := &Node[T]{
		Value: item,
		Next:  l.head,
	}
	l.head = newNode
	l.size++
}

// InsertBack adds an item to the back of the list.
func (l *SinglyLinkedList[T]) InsertBack(item T) {
	newNode := &Node[T]{
		Value: item,
		Next:  nil,
	}

	if l.head == nil {
		l.head = newNode
	} else {
		current := l.head
		for current.Next != nil {
			current = current.Next
		}
		current.Next = newNode
	}
	l.size++
}

// InsertAt inserts an item at the specified index.
func (l *SinglyLinkedList[T]) InsertAt(index int, item T) bool {
	if index < 0 || index > l.size {
		return false
	}

	if index == 0 {
		l.InsertFront(item)
		return true
	}

	if index == l.size {
		l.InsertBack(item)
		return true
	}

	newNode := &Node[T]{
		Value: item,
		Next:  nil,
	}

	current := l.head
	for i := 0; i < index-1; i++ {
		current = current.Next
	}

	newNode.Next = current.Next
	current.Next = newNode
	l.size++
	return true
}

// RemoveFront removes and returns the front item.
func (l *SinglyLinkedList[T]) RemoveFront() (T, bool) {
	var zero T
	if l.head == nil {
		return zero, false
	}

	item := l.head.Value
	l.head = l.head.Next
	l.size--
	return item, true
}

// RemoveBack removes and returns the back item.
func (l *SinglyLinkedList[T]) RemoveBack() (T, bool) {
	var zero T
	if l.head == nil {
		return zero, false
	}

	if l.head.Next == nil {
		item := l.head.Value
		l.head = nil
		l.size = 0
		return item, true
	}

	current := l.head
	for current.Next.Next != nil {
		current = current.Next
	}

	item := current.Next.Value
	current.Next = nil
	l.size--
	return item, true
}

// RemoveAt removes and returns the item at the specified index.
func (l *SinglyLinkedList[T]) RemoveAt(index int) (T, bool) {
	var zero T
	if index < 0 || index >= l.size {
		return zero, false
	}

	if index == 0 {
		return l.RemoveFront()
	}

	current := l.head
	for range index - 1 {
		current = current.Next
	}

	item := current.Next.Value
	current.Next = current.Next.Next
	l.size--
	return item, true
}

// GetAt returns the item at the specified index.
func (l *SinglyLinkedList[T]) GetAt(index int) (T, bool) {
	var zero T
	if index < 0 || index >= l.size {
		return zero, false
	}

	current := l.head
	for i := 0; i < index; i++ {
		current = current.Next
	}

	return current.Value, true
}

// SetAt sets the item at the specified index.
func (l *SinglyLinkedList[T]) SetAt(index int, item T) bool {
	if index < 0 || index >= l.size {
		return false
	}

	current := l.head
	for range index {
		current = current.Next
	}

	current.Value = item
	return true
}

// IsEmpty returns true if the list has no items.
func (l *SinglyLinkedList[T]) IsEmpty() bool {
	return l.size == 0
}

// Size returns the number of items in the list.
func (l *SinglyLinkedList[T]) Size() int {
	return l.size
}

// Clear removes all items from the list.
func (l *SinglyLinkedList[T]) Clear() {
	l.head = nil
	l.size = 0
}

// Clone creates a shallow copy of the list.
func (l *SinglyLinkedList[T]) Clone() List[T] {
	clone := NewSinglyLinkedList[T]()
	current := l.head
	for current != nil {
		clone.InsertBack(current.Value)
		current = current.Next
	}
	return clone
}

// Front returns the front item without removing it.
func (l *SinglyLinkedList[T]) Front() (T, bool) {
	var zero T
	if l.head == nil {
		return zero, false
	}
	return l.head.Value, true
}

// Back returns the back item without removing it.
func (l *SinglyLinkedList[T]) Back() (T, bool) {
	var zero T
	if l.head == nil {
		return zero, false
	}

	current := l.head
	for current.Next != nil {
		current = current.Next
	}
	return current.Value, true
}

// ForEach iterates over all items in the list.
func (l *SinglyLinkedList[T]) ForEach(fn func(item T) bool) {
	current := l.head
	for current != nil {
		if !fn(current.Value) {
			break
		}
		current = current.Next
	}
}

// Append adds multiple items to the back of the list.
func (l *SinglyLinkedList[T]) Append(items ...T) {
	for _, item := range items {
		l.InsertBack(item)
	}
}

// ToSlice converts the list to a slice.
func (l *SinglyLinkedList[T]) ToSlice() []T {
	result := make([]T, l.size)
	current := l.head
	for i := 0; i < l.size; i++ {
		result[i] = current.Value
		current = current.Next
	}
	return result
}

// AppendList appends all items from another list to this list.
func (l *SinglyLinkedList[T]) AppendList(other List[T]) {
	other.ForEach(func(item T) bool {
		l.InsertBack(item)
		return true
	})
}

// DoublyNode represents a node in a doubly linked list.
type DoublyNode[T any] struct {
	Value T
	Prev  *DoublyNode[T]
	Next  *DoublyNode[T]
}

// DoublyLinkedList represents a doubly linked list implementation.
type DoublyLinkedList[T any] struct {
	head *DoublyNode[T]
	tail *DoublyNode[T]
	size int
}

// Ensure DoublyLinkedList implements List interface.
var _ List[any] = (*DoublyLinkedList[any])(nil)

// NewDoublyLinkedList creates a new empty doubly linked list.
func NewDoublyLinkedList[T any]() List[T] {
	return &DoublyLinkedList[T]{
		head: nil,
		tail: nil,
		size: 0,
	}
}

// InsertFront adds an item to the front of the list.
func (l *DoublyLinkedList[T]) InsertFront(item T) {
	newNode := &DoublyNode[T]{
		Value: item,
		Prev:  nil,
		Next:  l.head,
	}

	if l.head != nil {
		l.head.Prev = newNode
	} else {
		l.tail = newNode
	}

	l.head = newNode
	l.size++
}

// InsertBack adds an item to the back of the list.
func (l *DoublyLinkedList[T]) InsertBack(item T) {
	newNode := &DoublyNode[T]{
		Value: item,
		Prev:  l.tail,
		Next:  nil,
	}

	if l.tail != nil {
		l.tail.Next = newNode
	} else {
		l.head = newNode
	}

	l.tail = newNode
	l.size++
}

// InsertAt inserts an item at the specified index.
func (l *DoublyLinkedList[T]) InsertAt(index int, item T) bool {
	if index < 0 || index > l.size {
		return false
	}

	if index == 0 {
		l.InsertFront(item)
		return true
	}

	if index == l.size {
		l.InsertBack(item)
		return true
	}

	// Find the node at the target position
	var target *DoublyNode[T]
	if index <= l.size/2 {
		// Search from head
		target = l.head
		for range index {
			target = target.Next
		}
	} else {
		// Search from tail
		target = l.tail
		for i := l.size - 1; i > index; i-- {
			target = target.Prev
		}
	}

	newNode := &DoublyNode[T]{
		Value: item,
		Prev:  target.Prev,
		Next:  target,
	}

	target.Prev.Next = newNode
	target.Prev = newNode
	l.size++
	return true
}

// RemoveFront removes and returns the front item.
func (l *DoublyLinkedList[T]) RemoveFront() (T, bool) {
	var zero T
	if l.head == nil {
		return zero, false
	}

	item := l.head.Value
	l.head = l.head.Next

	if l.head != nil {
		l.head.Prev = nil
	} else {
		l.tail = nil
	}

	l.size--
	return item, true
}

// RemoveBack removes and returns the back item.
func (l *DoublyLinkedList[T]) RemoveBack() (T, bool) {
	var zero T
	if l.tail == nil {
		return zero, false
	}

	item := l.tail.Value
	l.tail = l.tail.Prev

	if l.tail != nil {
		l.tail.Next = nil
	} else {
		l.head = nil
	}

	l.size--
	return item, true
}

// RemoveAt removes and returns the item at the specified index.
func (l *DoublyLinkedList[T]) RemoveAt(index int) (T, bool) {
	var zero T
	if index < 0 || index >= l.size {
		return zero, false
	}

	if index == 0 {
		return l.RemoveFront()
	}

	if index == l.size-1 {
		return l.RemoveBack()
	}

	// Find the node to remove
	var target *DoublyNode[T]
	if index <= l.size/2 {
		// Search from head
		target = l.head
		for range index {
			target = target.Next
		}
	} else {
		// Search from tail
		target = l.tail
		for i := l.size - 1; i > index; i-- {
			target = target.Prev
		}
	}

	item := target.Value
	target.Prev.Next = target.Next
	target.Next.Prev = target.Prev
	l.size--
	return item, true
}

// GetAt returns the item at the specified index.
func (l *DoublyLinkedList[T]) GetAt(index int) (T, bool) {
	var zero T
	if index < 0 || index >= l.size {
		return zero, false
	}

	var target *DoublyNode[T]
	if index <= l.size/2 {
		// Search from head
		target = l.head
		for range index {
			target = target.Next
		}
	} else {
		// Search from tail
		target = l.tail
		for i := l.size - 1; i > index; i-- {
			target = target.Prev
		}
	}

	return target.Value, true
}

// SetAt sets the item at the specified index.
func (l *DoublyLinkedList[T]) SetAt(index int, item T) bool {
	if index < 0 || index >= l.size {
		return false
	}

	var target *DoublyNode[T]
	if index <= l.size/2 {
		// Search from head
		target = l.head
		for range index {
			target = target.Next
		}
	} else {
		// Search from tail
		target = l.tail
		for i := l.size - 1; i > index; i-- {
			target = target.Prev
		}
	}

	target.Value = item
	return true
}

// IsEmpty returns true if the list has no items.
func (l *DoublyLinkedList[T]) IsEmpty() bool {
	return l.size == 0
}

// Size returns the number of items in the list.
func (l *DoublyLinkedList[T]) Size() int {
	return l.size
}

// Clear removes all items from the list.
func (l *DoublyLinkedList[T]) Clear() {
	l.head = nil
	l.tail = nil
	l.size = 0
}

// Clone creates a shallow copy of the list.
func (l *DoublyLinkedList[T]) Clone() List[T] {
	clone := NewDoublyLinkedList[T]()
	current := l.head
	for current != nil {
		clone.InsertBack(current.Value)
		current = current.Next
	}
	return clone
}

// Front returns the front item without removing it.
func (l *DoublyLinkedList[T]) Front() (T, bool) {
	var zero T
	if l.head == nil {
		return zero, false
	}
	return l.head.Value, true
}

// Back returns the back item without removing it.
func (l *DoublyLinkedList[T]) Back() (T, bool) {
	var zero T
	if l.tail == nil {
		return zero, false
	}
	return l.tail.Value, true
}

// ForEach iterates over all items in the list.
func (l *DoublyLinkedList[T]) ForEach(fn func(item T) bool) {
	current := l.head
	for current != nil {
		if !fn(current.Value) {
			break
		}
		current = current.Next
	}
}

// Append adds multiple items to the back of the list.
func (l *DoublyLinkedList[T]) Append(items ...T) {
	for _, item := range items {
		l.InsertBack(item)
	}
}

// ToSlice converts the list to a slice.
func (l *DoublyLinkedList[T]) ToSlice() []T {
	result := make([]T, l.size)
	current := l.head
	for i := 0; i < l.size; i++ {
		result[i] = current.Value
		current = current.Next
	}
	return result
}

// AppendList appends all items from another list to this list.
func (l *DoublyLinkedList[T]) AppendList(other List[T]) {
	other.ForEach(func(item T) bool {
		l.InsertBack(item)
		return true
	})
}
