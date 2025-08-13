// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package container

import (
	"testing"
)

func BenchmarkNewUint64Bitmap(b *testing.B) {
	b.ReportAllocs()
	for i := 0; b.Loop(); i++ {
		NewUint64Bitmap(uint64(i))
	}
}

func BenchmarkBitmap_Set(b *testing.B) {
	bm := NewUint64Bitmap(0)
	b.ReportAllocs()


	for i := 0; b.Loop(); i++ {
		bm.Set(uint64(i % 64))
	}
}

func BenchmarkBitmap_Clear(b *testing.B) {
	bm := NewUint64Bitmap(0xFFFFFFFFFFFFFFFF)
	b.ReportAllocs()


	for i := 0; b.Loop(); i++ {
		bm.Clear(uint64(i % 64))
	}
}

func BenchmarkBitmap_Toggle(b *testing.B) {
	bm := NewUint64Bitmap(0)
	b.ReportAllocs()


	for i := 0; b.Loop(); i++ {
		bm.Toggle(uint64(i % 64))
	}
}

func BenchmarkBitmap_IsSet(b *testing.B) {
	bm := NewUint64Bitmap(0x123456789ABCDEF0)
	b.ReportAllocs()


	for i := 0; b.Loop(); i++ {
		bm.IsSet(uint64(i % 64))
	}
}

func BenchmarkBitmap_IsAnySet(b *testing.B) {
	bm := NewUint64Bitmap(0x123456789ABCDEF0)
	b.ReportAllocs()


	for i := 0; b.Loop(); i++ {
		bm.IsAnySet(uint64(i % 64))
	}
}

func BenchmarkBitmap_IsAllSet(b *testing.B) {
	bm := NewUint64Bitmap(0x123456789ABCDEF0)
	b.ReportAllocs()


	for i := 0; b.Loop(); i++ {
		bm.IsAllSet(uint64(i % 64))
	}
}

func BenchmarkBitmap_SetAll(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		bm := NewUint64Bitmap(0)
		bm.SetAll()
	}
}

func BenchmarkBitmap_ClearAll(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		bm := NewUint64Bitmap(0xFFFFFFFFFFFFFFFF)
		bm.ClearAll()
	}
}

func BenchmarkBitmap_Count(b *testing.B) {
	bm := NewUint64Bitmap(0x123456789ABCDEF0)
	b.ReportAllocs()


	for b.Loop() {
		bm.Count()
	}
}

func BenchmarkBitmap_IsEmpty(b *testing.B) {
	bm := NewUint64Bitmap(0x123456789ABCDEF0)
	b.ReportAllocs()


	for b.Loop() {
		bm.IsEmpty()
	}
}

func BenchmarkBitmap_IsFull(b *testing.B) {
	bm := NewUint64Bitmap(0x123456789ABCDEF0)
	b.ReportAllocs()


	for b.Loop() {
		bm.IsFull()
	}
}

func BenchmarkBitmap_And(b *testing.B) {
	bm1 := NewUint64Bitmap(0x123456789ABCDEF0)
	bm2 := NewUint64Bitmap(0xFEDCBA9876543210)
	b.ReportAllocs()


	for b.Loop() {
		bm1.And(bm2)
	}
}

func BenchmarkBitmap_Or(b *testing.B) {
	bm1 := NewUint64Bitmap(0x123456789ABCDEF0)
	bm2 := NewUint64Bitmap(0xFEDCBA9876543210)
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		bm1.Or(bm2)
	}
}

func BenchmarkBitmap_Xor(b *testing.B) {
	bm1 := NewUint64Bitmap(0x123456789ABCDEF0)
	bm2 := NewUint64Bitmap(0xFEDCBA9876543210)
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		bm1.Xor(bm2)
	}
}

func BenchmarkBitmap_Not(b *testing.B) {
	bm := NewUint64Bitmap(0x123456789ABCDEF0)
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		bm.Not()
	}
}

func BenchmarkBitmap_Clone(b *testing.B) {
	bm := NewUint64Bitmap(0x123456789ABCDEF0)
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		bm.Clone()
	}
}

func BenchmarkBitmap_Equals(b *testing.B) {
	bm1 := NewUint64Bitmap(0x123456789ABCDEF0)
	bm2 := NewUint64Bitmap(0x123456789ABCDEF0)
	bm3 := NewUint64Bitmap(0xFEDCBA9876543210)
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		bm1.Equals(bm2)
		bm1.Equals(bm3)
	}
}

func BenchmarkBitmap_String(b *testing.B) {
	bm := NewUint64Bitmap(0x123456789ABCDEF0)
	b.ReportAllocs()


	for b.Loop() {
		_ = bm.String()
	}
}

func BenchmarkBitmap_ToBinaryString(b *testing.B) {
	bm := NewUint64Bitmap(0x123456789ABCDEF0)
	b.ReportAllocs()


	for b.Loop() {
		bm.ToBinaryString()
	}
}

func BenchmarkBitmap_ToHexString(b *testing.B) {
	bm := NewUint64Bitmap(0x123456789ABCDEF0)
	b.ReportAllocs()


	for b.Loop() {
		bm.ToHexString()
	}
}

func BenchmarkBitmap_GetSetBits(b *testing.B) {
	bm := NewUint64Bitmap(0x123456789ABCDEF0)
	b.ReportAllocs()


	for b.Loop() {
		bm.GetSetBits()
	}
}

func BenchmarkBitmap_SetBits(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		bm := NewUint64Bitmap(0)
		bm.SetBits(1, 2, 4, 8, 16, 32)
	}
}

func BenchmarkBitmap_ClearBits(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		bm := NewUint64Bitmap(0xFFFFFFFFFFFFFFFF)
		bm.ClearBits(1, 2, 4, 8, 16, 32)
	}
}

func BenchmarkBitmap_ToggleBits(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		bm := NewUint64Bitmap(0x123456789ABCDEF0)
		bm.ToggleBits(1, 2, 4, 8, 16, 32)
	}
}

func BenchmarkBitmap_IsSubset(b *testing.B) {
	bm1 := NewUint64Bitmap(0x5)
	bm2 := NewUint64Bitmap(0x7)
	bm3 := NewUint64Bitmap(0x3)
	b.ReportAllocs()


	for b.Loop() {
		bm1.IsSubset(bm2)
		bm1.IsSubset(bm3)
	}
}

func BenchmarkBitmap_IsSuperset(b *testing.B) {
	bm1 := NewUint64Bitmap(0x7)
	bm2 := NewUint64Bitmap(0x5)
	bm3 := NewUint64Bitmap(0xF)
	b.ReportAllocs()


	for b.Loop() {
		bm1.IsSuperset(bm2)
		bm1.IsSuperset(bm3)
	}
}

func BenchmarkBitmap_Intersection(b *testing.B) {
	bm1 := NewUint64Bitmap(0x123456789ABCDEF0)
	bm2 := NewUint64Bitmap(0xFEDCBA9876543210)
	b.ReportAllocs()


	for b.Loop() {
		bm1.Intersection(bm2)
	}
}

func BenchmarkBitmap_Union(b *testing.B) {
	bm1 := NewUint64Bitmap(0x123456789ABCDEF0)
	bm2 := NewUint64Bitmap(0xFEDCBA9876543210)
	b.ReportAllocs()


	for b.Loop() {
		bm1.Union(bm2)
	}
}

func BenchmarkBitmap_Difference(b *testing.B) {
	bm1 := NewUint64Bitmap(0x123456789ABCDEF0)
	bm2 := NewUint64Bitmap(0xFEDCBA9876543210)
	b.ReportAllocs()


	for b.Loop() {
		bm1.Difference(bm2)
	}
}

func BenchmarkBitmap_SymmetricDifference(b *testing.B) {
	bm1 := NewUint64Bitmap(0x123456789ABCDEF0)
	bm2 := NewUint64Bitmap(0xFEDCBA9876543210)
	b.ReportAllocs()


	for b.Loop() {
		bm1.SymmetricDifference(bm2)
	}
}

func BenchmarkBitmap_ShiftLeft(b *testing.B) {
	bm := NewUint64Bitmap(0x123456789ABCDEF0)
	b.ReportAllocs()


	for i := 0; b.Loop(); i++ {
		bm.ShiftLeft(i % 64)
	}
}

func BenchmarkBitmap_ShiftRight(b *testing.B) {
	bm := NewUint64Bitmap(0x123456789ABCDEF0)
	b.ReportAllocs()


	for i := 0; b.Loop(); i++ {
		bm.ShiftRight(i % 64)
	}
}

func BenchmarkBitmap_RotateLeft(b *testing.B) {
	bm := NewUint64Bitmap(0x123456789ABCDEF0)
	b.ReportAllocs()


	for i := 0; b.Loop(); i++ {
		bm.RotateLeft(i % 64)
	}
}

func BenchmarkBitmap_RotateRight(b *testing.B) {
	bm := NewUint64Bitmap(0x123456789ABCDEF0)
	b.ReportAllocs()


	for i := 0; b.Loop(); i++ {
		bm.RotateRight(i % 64)
	}
}

func BenchmarkBitmap_GetLowestSetBit(b *testing.B) {
	bm := NewUint64Bitmap(0x123456789ABCDEF0)
	b.ReportAllocs()


	for b.Loop() {
		bm.GetLowestSetBit()
	}
}

func BenchmarkBitmap_GetHighestSetBit(b *testing.B) {
	bm := NewUint64Bitmap(0x123456789ABCDEF0)
	b.ReportAllocs()


	for b.Loop() {
		bm.GetHighestSetBit()
	}
}

func BenchmarkBitmap_NextSetBit(b *testing.B) {
	bm := NewUint64Bitmap(0x123456789ABCDEF0)
	b.ReportAllocs()


	for i := 0; b.Loop(); i++ {
		bm.NextSetBit(i % 63)
	}
}

func BenchmarkBitmap_PrevSetBit(b *testing.B) {
	bm := NewUint64Bitmap(0x123456789ABCDEF0)
	b.ReportAllocs()


	for i := 0; b.Loop(); i++ {
		bm.PrevSetBit((i % 63) + 1)
	}
}

// Benchmark for multiple operations in sequence
func BenchmarkBitmap_MixedOperations(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		bm := NewUint64Bitmap(0)
		bm.SetBits(1, 2, 4, 8, 16, 32)
		bm.Toggle(64)
		bm.Clear(2)
		bm.And(NewUint64Bitmap(0xFFFFFFFF))
		bm.Or(NewUint64Bitmap(0x100000000))
		bm.Count()
		bm.GetSetBits()
	}
}
