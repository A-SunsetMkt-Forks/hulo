// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package container

import (
	"fmt"
)

// Bitmap defines the interface for bitmap operations
type Bitmap interface {
	Set(bit uint64)
	Clear(bit uint64)
	Toggle(bit uint64)
	IsSet(bit uint64) bool
	IsAnySet(bits uint64) bool
	IsAllSet(bits uint64) bool
	Get() uint64
	SetAll()
	ClearAll()
	Count() int
	IsEmpty() bool
	IsFull() bool
	And(other Bitmap)
	Or(other Bitmap)
	Xor(other Bitmap)
	Not()
	Clone() Bitmap
	Equals(other Bitmap) bool
	String() string
	ToBinaryString() string
	ToHexString() string
	GetSetBits() []int
	SetBits(bits ...uint64)
	ClearBits(bits ...uint64)
	ToggleBits(bits ...uint64)
	IsSubset(other Bitmap) bool
	IsSuperset(other Bitmap) bool
	Intersection(other Bitmap) Bitmap
	Union(other Bitmap) Bitmap
	Difference(other Bitmap) Bitmap
	SymmetricDifference(other Bitmap) Bitmap
	ShiftLeft(n int)
	ShiftRight(n int)
	RotateLeft(n int)
	RotateRight(n int)
	GetLowestSetBit() int
	GetHighestSetBit() int
	NextSetBit(from int) int
	PrevSetBit(from int) int
}

// Uint64Bitmap represents a bitmap implementation using uint64
type Uint64Bitmap struct {
	bits uint64
}

// Ensure Uint64Bitmap implements Bitmap interface
var _ Bitmap = (*Uint64Bitmap)(nil)

// NewUint64Bitmap creates a new bitmap with initial value
func NewUint64Bitmap(initial uint64) Bitmap {
	return &Uint64Bitmap{bits: initial}
}

// Set enables a specific bit
func (b *Uint64Bitmap) Set(bit uint64) {
	b.bits |= bit
}

// Clear disables a specific bit
func (b *Uint64Bitmap) Clear(bit uint64) {
	b.bits &^= bit
}

// Toggle toggles a specific bit
func (b *Uint64Bitmap) Toggle(bit uint64) {
	b.bits ^= bit
}

// IsSet checks if a specific bit is enabled
func (b *Uint64Bitmap) IsSet(bit uint64) bool {
	return (b.bits & bit) != 0
}

// IsAnySet checks if any of the given bits are enabled
func (b *Uint64Bitmap) IsAnySet(bits uint64) bool {
	return (b.bits & bits) != 0
}

// IsAllSet checks if all of the given bits are enabled
func (b *Uint64Bitmap) IsAllSet(bits uint64) bool {
	return (b.bits & bits) == bits
}

// Get returns the current bitmap value
func (b *Uint64Bitmap) Get() uint64 {
	return b.bits
}

// SetAll enables all bits
func (b *Uint64Bitmap) SetAll() {
	b.bits = 0xFFFFFFFFFFFFFFFF
}

// ClearAll disables all bits
func (b *Uint64Bitmap) ClearAll() {
	b.bits = 0
}

// Count returns the number of set bits
func (b *Uint64Bitmap) Count() int {
	count := 0
	bits := b.bits
	for bits != 0 {
		count += int(bits & 1)
		bits >>= 1
	}
	return count
}

// IsEmpty returns true if no bits are set
func (b *Uint64Bitmap) IsEmpty() bool {
	return b.bits == 0
}

// IsFull returns true if all bits are set
func (b *Uint64Bitmap) IsFull() bool {
	return b.bits == 0xFFFFFFFFFFFFFFFF
}

// And performs bitwise AND with another bitmap
func (b *Uint64Bitmap) And(other Bitmap) {
	b.bits &= other.Get()
}

// Or performs bitwise OR with another bitmap
func (b *Uint64Bitmap) Or(other Bitmap) {
	b.bits |= other.Get()
}

// Xor performs bitwise XOR with another bitmap
func (b *Uint64Bitmap) Xor(other Bitmap) {
	b.bits ^= other.Get()
}

// Not performs bitwise NOT
func (b *Uint64Bitmap) Not() {
	b.bits = ^b.bits
}

// Clone creates a copy of the bitmap
func (b *Uint64Bitmap) Clone() Bitmap {
	return &Uint64Bitmap{bits: b.bits}
}

// Equals checks if two bitmaps are equal
func (b *Uint64Bitmap) Equals(other Bitmap) bool {
	return b.bits == other.Get()
}

// String returns a string representation of the bitmap
func (b *Uint64Bitmap) String() string {
	return fmt.Sprintf("Uint64Bitmap(0x%016X)", b.bits)
}

// ToBinaryString returns binary representation
func (b *Uint64Bitmap) ToBinaryString() string {
	return fmt.Sprintf("%064b", b.bits)
}

// ToHexString returns hexadecimal representation
func (b *Uint64Bitmap) ToHexString() string {
	return fmt.Sprintf("0x%016X", b.bits)
}

// GetSetBits returns a slice of bit positions that are set
func (b *Uint64Bitmap) GetSetBits() []int {
	var positions []int
	bits := b.bits
	pos := 0

	for bits != 0 {
		if bits&1 != 0 {
			positions = append(positions, pos)
		}
		bits >>= 1
		pos++
	}

	return positions
}

// SetBits sets multiple bits at once
func (b *Uint64Bitmap) SetBits(bits ...uint64) {
	for _, bit := range bits {
		b.Set(bit)
	}
}

// ClearBits clears multiple bits at once
func (b *Uint64Bitmap) ClearBits(bits ...uint64) {
	for _, bit := range bits {
		b.Clear(bit)
	}
}

// ToggleBits toggles multiple bits at once
func (b *Uint64Bitmap) ToggleBits(bits ...uint64) {
	for _, bit := range bits {
		b.Toggle(bit)
	}
}

// IsSubset checks if this bitmap is a subset of another
func (b *Uint64Bitmap) IsSubset(other Bitmap) bool {
	return (b.bits & other.Get()) == b.bits
}

// IsSuperset checks if this bitmap is a superset of another
func (b *Uint64Bitmap) IsSuperset(other Bitmap) bool {
	return (b.bits & other.Get()) == other.Get()
}

// Intersection returns the intersection with another bitmap
func (b *Uint64Bitmap) Intersection(other Bitmap) Bitmap {
	return &Uint64Bitmap{bits: b.bits & other.Get()}
}

// Union returns the union with another bitmap
func (b *Uint64Bitmap) Union(other Bitmap) Bitmap {
	return &Uint64Bitmap{bits: b.bits | other.Get()}
}

// Difference returns the difference with another bitmap
func (b *Uint64Bitmap) Difference(other Bitmap) Bitmap {
	return &Uint64Bitmap{bits: b.bits &^ other.Get()}
}

// SymmetricDifference returns the symmetric difference with another bitmap
func (b *Uint64Bitmap) SymmetricDifference(other Bitmap) Bitmap {
	return &Uint64Bitmap{bits: b.bits ^ other.Get()}
}

// ShiftLeft shifts bits to the left by n positions
func (b *Uint64Bitmap) ShiftLeft(n int) {
	if n >= 64 {
		b.bits = 0
	} else if n > 0 {
		b.bits <<= uint(n)
	}
}

// ShiftRight shifts bits to the right by n positions
func (b *Uint64Bitmap) ShiftRight(n int) {
	if n >= 64 {
		b.bits = 0
	} else if n > 0 {
		b.bits >>= uint(n)
	}
}

// RotateLeft rotates bits to the left by n positions
func (b *Uint64Bitmap) RotateLeft(n int) {
	n = n % 64
	if n < 0 {
		n += 64
	}
	if n == 0 {
		return
	}

	b.bits = (b.bits << uint(n)) | (b.bits >> uint(64-n))
}

// RotateRight rotates bits to the right by n positions
func (b *Uint64Bitmap) RotateRight(n int) {
	b.RotateLeft(-n)
}

// GetLowestSetBit returns the position of the lowest set bit
func (b *Uint64Bitmap) GetLowestSetBit() int {
	if b.bits == 0 {
		return -1
	}
	return int(b.bits & -b.bits)
}

// GetHighestSetBit returns the position of the highest set bit
func (b *Uint64Bitmap) GetHighestSetBit() int {
	if b.bits == 0 {
		return -1
	}

	pos := 0
	bits := b.bits
	for bits > 1 {
		bits >>= 1
		pos++
	}
	return pos
}

// NextSetBit returns the next set bit after the given position
func (b *Uint64Bitmap) NextSetBit(from int) int {
	if from >= 63 {
		return -1
	}

	mask := uint64(1) << uint(from+1)
	bits := b.bits &^ (mask - 1)

	if bits == 0 {
		return -1
	}

	return from + 1 + int(bits&-bits)
}

// PrevSetBit returns the previous set bit before the given position
func (b *Uint64Bitmap) PrevSetBit(from int) int {
	if from <= 0 {
		return -1
	}

	mask := uint64(1) << uint(from)
	bits := b.bits & (mask - 1)

	if bits == 0 {
		return -1
	}

	return int(bits&-bits) - 1
}
