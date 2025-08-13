// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package container

import (
	"testing"
)

func TestNewUint64Bitmap(t *testing.T) {
	tests := []struct {
		name     string
		initial  uint64
		expected uint64
	}{
		{"zero", 0, 0},
		{"one", 1, 1},
		{"max", 0xFFFFFFFFFFFFFFFF, 0xFFFFFFFFFFFFFFFF},
		{"random", 0x123456789ABCDEF0, 0x123456789ABCDEF0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bm := NewUint64Bitmap(tt.initial)
			if got := bm.Get(); got != tt.expected {
				t.Errorf("NewUint64Bitmap() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestBitmap_Set(t *testing.T) {
	bm := NewUint64Bitmap(0)

	// Test setting individual bits
	bm.Set(1)
	if !bm.IsSet(1) {
		t.Error("Bit 1 should be set")
	}

	bm.Set(2)
	if !bm.IsSet(2) {
		t.Error("Bit 2 should be set")
	}

	// Test setting the same bit multiple times
	bm.Set(1)
	if !bm.IsSet(1) {
		t.Error("Bit 1 should still be set after setting again")
	}

	// Test setting multiple bits
	bm.Set(4)
	bm.Set(8)
	if !bm.IsSet(4) || !bm.IsSet(8) {
		t.Error("Bits 4 and 8 should be set")
	}
}

func TestBitmap_Clear(t *testing.T) {
	bm := NewUint64Bitmap(0xFFFFFFFFFFFFFFFF)

	// Test clearing individual bits
	bm.Clear(1)
	if bm.IsSet(1) {
		t.Error("Bit 1 should be cleared")
	}

	bm.Clear(2)
	if bm.IsSet(2) {
		t.Error("Bit 2 should be cleared")
	}

	// Test clearing the same bit multiple times
	bm.Clear(1)
	if bm.IsSet(1) {
		t.Error("Bit 1 should still be cleared after clearing again")
	}
}

func TestBitmap_Toggle(t *testing.T) {
	bm := NewUint64Bitmap(0)

	// Test toggling from 0 to 1
	bm.Toggle(1)
	if !bm.IsSet(1) {
		t.Error("Bit 1 should be set after toggle")
	}

	// Test toggling from 1 to 0
	bm.Toggle(1)
	if bm.IsSet(1) {
		t.Error("Bit 1 should be cleared after second toggle")
	}

	// Test toggling multiple times
	bm.Toggle(1)
	bm.Toggle(1)
	if bm.IsSet(1) {
		t.Error("Bit 1 should be cleared after two toggles")
	}
}

func TestBitmap_IsSet(t *testing.T) {
	bm := NewUint64Bitmap(0x5) // 0101 in binary

	if !bm.IsSet(1) {
		t.Error("Bit 1 should be set")
	}
	if bm.IsSet(2) {
		t.Error("Bit 2 should not be set")
	}
	if !bm.IsSet(4) {
		t.Error("Bit 4 should be set")
	}
	if bm.IsSet(8) {
		t.Error("Bit 8 should not be set")
	}
}

func TestBitmap_IsAnySet(t *testing.T) {
	bm := NewUint64Bitmap(0x5) // 0101 in binary

	if !bm.IsAnySet(0x3) { // 0011
		t.Error("Should detect any set bit in 0x3")
	}
	if !bm.IsAnySet(0x4) { // 0100
		t.Error("Should detect any set bit in 0x4")
	}
	if bm.IsAnySet(0x2) { // 0010
		t.Error("Should not detect any set bit in 0x2")
	}
}

func TestBitmap_IsAllSet(t *testing.T) {
	bm := NewUint64Bitmap(0x7) // 0111 in binary

	if !bm.IsAllSet(0x5) { // 0101
		t.Error("Should detect all bits set in 0x5")
	}
	if bm.IsAllSet(0xF) { // 1111
		t.Error("Should not detect all bits set in 0xF")
	}
}

func TestBitmap_SetAll(t *testing.T) {
	bm := NewUint64Bitmap(0)
	bm.SetAll()

	if !bm.IsFull() {
		t.Error("Bitmap should be full after SetAll")
	}
	if bm.Get() != 0xFFFFFFFFFFFFFFFF {
		t.Error("Bitmap should have all bits set")
	}
}

func TestBitmap_ClearAll(t *testing.T) {
	bm := NewUint64Bitmap(0xFFFFFFFFFFFFFFFF)
	bm.ClearAll()

	if !bm.IsEmpty() {
		t.Error("Bitmap should be empty after ClearAll")
	}
	if bm.Get() != 0 {
		t.Error("Bitmap should have no bits set")
	}
}

func TestBitmap_Count(t *testing.T) {
	tests := []struct {
		name     string
		value    uint64
		expected int
	}{
		{"zero", 0, 0},
		{"one", 1, 1},
		{"two", 3, 2},
		{"three", 7, 3},
		{"all", 0xFFFFFFFFFFFFFFFF, 64},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bm := NewUint64Bitmap(tt.value)
			if got := bm.Count(); got != tt.expected {
				t.Errorf("Count() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestBitmap_IsEmpty(t *testing.T) {
	empty := NewUint64Bitmap(0)
	if !empty.IsEmpty() {
		t.Error("Empty bitmap should be empty")
	}

	nonEmpty := NewUint64Bitmap(1)
	if nonEmpty.IsEmpty() {
		t.Error("Non-empty bitmap should not be empty")
	}
}

func TestBitmap_IsFull(t *testing.T) {
	full := NewUint64Bitmap(0xFFFFFFFFFFFFFFFF)
	if !full.IsFull() {
		t.Error("Full bitmap should be full")
	}

	nonFull := NewUint64Bitmap(1)
	if nonFull.IsFull() {
		t.Error("Non-full bitmap should not be full")
	}
}

func TestBitmap_And(t *testing.T) {
	bm1 := NewUint64Bitmap(0x5) // 0101
	bm2 := NewUint64Bitmap(0x3) // 0011

	bm1.And(bm2)
	if bm1.Get() != 0x1 { // 0001
		t.Errorf("And() = %v, want %v", bm1.Get(), uint64(0x1))
	}
}

func TestBitmap_Or(t *testing.T) {
	bm1 := NewUint64Bitmap(0x5) // 0101
	bm2 := NewUint64Bitmap(0x3) // 0011

	bm1.Or(bm2)
	if bm1.Get() != 0x7 { // 0111
		t.Errorf("Or() = %v, want %v", bm1.Get(), uint64(0x7))
	}
}

func TestBitmap_Xor(t *testing.T) {
	bm1 := NewUint64Bitmap(0x5) // 0101
	bm2 := NewUint64Bitmap(0x3) // 0011

	bm1.Xor(bm2)
	if bm1.Get() != 0x6 { // 0110
		t.Errorf("Xor() = %v, want %v", bm1.Get(), uint64(0x6))
	}
}

func TestBitmap_Not(t *testing.T) {
	bm := NewUint64Bitmap(0x5) // 0101
	bm.Not()

	expected := ^uint64(0x5)
	if bm.Get() != expected {
		t.Errorf("Not() = %v, want %v", bm.Get(), expected)
	}
}

func TestBitmap_Clone(t *testing.T) {
	original := NewUint64Bitmap(0x123456789ABCDEF0)
	clone := original.Clone()

	if !original.Equals(clone) {
		t.Error("Clone should be equal to original")
	}

	// Modify original, clone should remain unchanged
	original.Set(1)
	if original.Equals(clone) {
		t.Error("Clone should not be affected by original modification")
	}
}

func TestBitmap_Equals(t *testing.T) {
	bm1 := NewUint64Bitmap(0x123456789ABCDEF0)
	bm2 := NewUint64Bitmap(0x123456789ABCDEF0)
	bm3 := NewUint64Bitmap(0x123456789ABCDEF1)

	if !bm1.Equals(bm2) {
		t.Error("Equal bitmaps should be equal")
	}
	if bm1.Equals(bm3) {
		t.Error("Different bitmaps should not be equal")
	}
}

func TestBitmap_String(t *testing.T) {
	bm := NewUint64Bitmap(0x123456789ABCDEF0)
	str := bm.String()
	expected := "Uint64Bitmap(0x123456789ABCDEF0)"

	if str != expected {
		t.Errorf("String() = %v, want %v", str, expected)
	}
}

func TestBitmap_ToBinaryString(t *testing.T) {
	bm := NewUint64Bitmap(0x5) // 0101
	binary := bm.ToBinaryString()

	// Should be 64 characters long
	if len(binary) != 64 {
		t.Errorf("ToBinaryString() length = %v, want %v", len(binary), 64)
	}

	// Should end with "0101"
	if !endsWith(binary, "0101") {
		t.Errorf("ToBinaryString() = %v, should end with 0101", binary)
	}
}

func TestBitmap_ToHexString(t *testing.T) {
	bm := NewUint64Bitmap(0x123456789ABCDEF0)
	hex := bm.ToHexString()
	expected := "0x123456789ABCDEF0"

	if hex != expected {
		t.Errorf("ToHexString() = %v, want %v", hex, expected)
	}
}

func TestBitmap_GetSetBits(t *testing.T) {
	bm := NewUint64Bitmap(0x5) // 0101
	bits := bm.GetSetBits()

	expected := []int{0, 2}
	if !sliceEqual(bits, expected) {
		t.Errorf("GetSetBits() = %v, want %v", bits, expected)
	}
}

func TestBitmap_SetBits(t *testing.T) {
	bm := NewUint64Bitmap(0)
	bm.SetBits(1, 2, 4)

	if !bm.IsSet(1) || !bm.IsSet(2) || !bm.IsSet(4) {
		t.Error("All specified bits should be set")
	}
}

func TestBitmap_ClearBits(t *testing.T) {
	bm := NewUint64Bitmap(0x7) // 0111
	bm.ClearBits(1, 2)

	if bm.IsSet(1) || bm.IsSet(2) {
		t.Error("Specified bits should be cleared")
	}
	if !bm.IsSet(4) {
		t.Error("Unspecified bits should remain set")
	}
}

func TestBitmap_ToggleBits(t *testing.T) {
	bm := NewUint64Bitmap(0x5) // 0101
	bm.ToggleBits(1, 2)

	// Bit 1: 1 -> 0, Bit 2: 0 -> 1
	if bm.IsSet(1) || !bm.IsSet(2) {
		t.Error("Bits should be toggled correctly")
	}
}

func TestBitmap_IsSubset(t *testing.T) {
	bm1 := NewUint64Bitmap(0x5) // 0101
	bm2 := NewUint64Bitmap(0x7) // 0111

	if !bm1.IsSubset(bm2) {
		t.Error("0x5 should be a subset of 0x7")
	}
	if bm2.IsSubset(bm1) {
		t.Error("0x7 should not be a subset of 0x5")
	}
}

func TestBitmap_IsSuperset(t *testing.T) {
	bm1 := NewUint64Bitmap(0x7) // 0111
	bm2 := NewUint64Bitmap(0x5) // 0101

	if !bm1.IsSuperset(bm2) {
		t.Error("0x7 should be a superset of 0x5")
	}
	if bm2.IsSuperset(bm1) {
		t.Error("0x5 should not be a superset of 0x7")
	}
}

func TestBitmap_Intersection(t *testing.T) {
	bm1 := NewUint64Bitmap(0x5) // 0101
	bm2 := NewUint64Bitmap(0x3) // 0011

	result := bm1.Intersection(bm2)
	if result.Get() != 0x1 { // 0001
		t.Errorf("Intersection() = %v, want %v", result.Get(), uint64(0x1))
	}
}

func TestBitmap_Union(t *testing.T) {
	bm1 := NewUint64Bitmap(0x5) // 0101
	bm2 := NewUint64Bitmap(0x3) // 0011

	result := bm1.Union(bm2)
	if result.Get() != 0x7 { // 0111
		t.Errorf("Union() = %v, want %v", result.Get(), uint64(0x7))
	}
}

func TestBitmap_Difference(t *testing.T) {
	bm1 := NewUint64Bitmap(0x5) // 0101
	bm2 := NewUint64Bitmap(0x3) // 0011

	result := bm1.Difference(bm2)
	if result.Get() != 0x4 { // 0100
		t.Errorf("Difference() = %v, want %v", result.Get(), uint64(0x4))
	}
}

func TestBitmap_SymmetricDifference(t *testing.T) {
	bm1 := NewUint64Bitmap(0x5) // 0101
	bm2 := NewUint64Bitmap(0x3) // 0011

	result := bm1.SymmetricDifference(bm2)
	if result.Get() != 0x6 { // 0110
		t.Errorf("SymmetricDifference() = %v, want %v", result.Get(), uint64(0x6))
	}
}

func TestBitmap_ShiftLeft(t *testing.T) {
	bm := NewUint64Bitmap(0x5) // 0101
	bm.ShiftLeft(2)

	if bm.Get() != 0x14 { // 10100
		t.Errorf("ShiftLeft(2) = %v, want %v", bm.Get(), uint64(0x14))
	}

	// Test shift by 64 or more
	bm.ShiftLeft(64)
	if bm.Get() != 0 {
		t.Error("ShiftLeft(64) should result in 0")
	}
}

func TestBitmap_ShiftRight(t *testing.T) {
	bm := NewUint64Bitmap(0x14) // 10100
	bm.ShiftRight(2)

	if bm.Get() != 0x5 { // 0101
		t.Errorf("ShiftRight(2) = %v, want %v", bm.Get(), uint64(0x5))
	}

	// Test shift by 64 or more
	bm.ShiftRight(64)
	if bm.Get() != 0 {
		t.Error("ShiftRight(64) should result in 0")
	}
}

func TestBitmap_RotateLeft(t *testing.T) {
	bm := NewUint64Bitmap(0x5) // 0101
	bm.RotateLeft(2)

	// 0101 rotated left by 2 should be 010100...0101
	expected := uint64(0x5) << 2
	if bm.Get() != expected {
		t.Errorf("RotateLeft(2) = %v, want %v", bm.Get(), expected)
	}
}

func TestBitmap_RotateRight(t *testing.T) {
	bm := NewUint64Bitmap(0x14) // 10100
	bm.RotateRight(2)

	// 10100 rotated right by 2 should be 00101
	expected := uint64(0x5)
	if bm.Get() != expected {
		t.Errorf("RotateRight(2) = %v, want %v", bm.Get(), expected)
	}
}

func TestBitmap_GetLowestSetBit(t *testing.T) {
	bm := NewUint64Bitmap(0x14) // 10100
	lowest := bm.GetLowestSetBit()

	if lowest != 2 {
		t.Errorf("GetLowestSetBit() = %v, want %v", lowest, 2)
	}

	// Test with no bits set
	empty := NewUint64Bitmap(0)
	if empty.GetLowestSetBit() != -1 {
		t.Error("GetLowestSetBit() should return -1 for empty bitmap")
	}
}

func TestBitmap_GetHighestSetBit(t *testing.T) {
	bm := NewUint64Bitmap(0x14) // 10100
	highest := bm.GetHighestSetBit()

	if highest != 4 {
		t.Errorf("GetHighestSetBit() = %v, want %v", highest, 4)
	}

	// Test with no bits set
	empty := NewUint64Bitmap(0)
	if empty.GetHighestSetBit() != -1 {
		t.Error("GetHighestSetBit() should return -1 for empty bitmap")
	}
}

func TestBitmap_NextSetBit(t *testing.T) {
	bm := NewUint64Bitmap(0x14) // 10100

	next := bm.NextSetBit(1)
	if next != 2 {
		t.Errorf("NextSetBit(1) = %v, want %v", next, 2)
	}

	next = bm.NextSetBit(2)
	if next != 4 {
		t.Errorf("NextSetBit(2) = %v, want %v", next, 4)
	}

	// Test beyond last set bit
	next = bm.NextSetBit(4)
	if next != -1 {
		t.Errorf("NextSetBit(4) = %v, want %v", next, -1)
	}
}

func TestBitmap_PrevSetBit(t *testing.T) {
	bm := NewUint64Bitmap(0x14) // 10100

	prev := bm.PrevSetBit(5)
	if prev != 4 {
		t.Errorf("PrevSetBit(5) = %v, want %v", prev, 4)
	}

	prev = bm.PrevSetBit(3)
	if prev != 2 {
		t.Errorf("PrevSetBit(3) = %v, want %v", prev, 2)
	}

	// Test before first set bit
	prev = bm.PrevSetBit(1)
	if prev != -1 {
		t.Errorf("PrevSetBit(1) = %v, want %v", prev, -1)
	}
}

// Helper functions
func endsWith(s, suffix string) bool {
	if len(s) < len(suffix) {
		return false
	}
	return s[len(s)-len(suffix):] == suffix
}

func sliceEqual(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
