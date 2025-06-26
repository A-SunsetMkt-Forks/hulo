// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package build

import (
	"fmt"
	"sync"
	"testing"

	"github.com/hulo-lang/hulo/internal/container"
	"github.com/hulo-lang/hulo/internal/util"
)

func TestNewAllocator(t *testing.T) {
	tests := []struct {
		name  string
		opts  []AllocatorOption
		check func(*testing.T, *Allocator)
	}{
		{
			name: "default allocator",
			opts: nil,
			check: func(t *testing.T, a *Allocator) {
				if a.prefix != "_v" {
					t.Errorf("expected default prefix '_v', got '%s'", a.prefix)
				}
				if a.counter != 0 {
					t.Errorf("expected default counter 0, got %d", a.counter)
				}
				if a.lock == nil {
					t.Error("expected non-nil lock")
				}
			},
		},
		{
			name: "custom prefix",
			opts: []AllocatorOption{WithPrefix("_var_")},
			check: func(t *testing.T, a *Allocator) {
				if a.prefix != "_var_" {
					t.Errorf("expected prefix '_var_', got '%s'", a.prefix)
				}
			},
		},
		{
			name: "custom counter",
			opts: []AllocatorOption{WithInitialCounter(100)},
			check: func(t *testing.T, a *Allocator) {
				if a.counter != 100 {
					t.Errorf("expected counter 100, got %d", a.counter)
				}
			},
		},
		{
			name: "custom lock",
			opts: []AllocatorOption{WithLock(&util.NoOpLocker{})},
			check: func(t *testing.T, a *Allocator) {
				if a.lock == nil {
					t.Error("expected non-nil lock")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			allocator := NewAllocator(tt.opts...)
			tt.check(t, allocator)
		})
	}
}

func TestAllocName(t *testing.T) {
	allocator := NewAllocator()

	tests := []struct {
		name     string
		original string
		want     string
	}{
		{
			name:     "first allocation",
			original: "count",
			want:     "_v1",
		},
		{
			name:     "second allocation",
			original: "total",
			want:     "_v2",
		},
		{
			name:     "duplicate allocation",
			original: "count",
			want:     "_v1", // should return the same name
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := allocator.AllocName(tt.original)
			if got != tt.want {
				t.Errorf("AllocName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetOriginalName(t *testing.T) {
	allocator := NewAllocator()
	original := "count"
	generated := allocator.AllocName(original)

	tests := []struct {
		name      string
		generated string
		want      string
	}{
		{
			name:      "existing name",
			generated: generated,
			want:      original,
		},
		{
			name:      "non-existent name",
			generated: "_v999",
			want:      "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := allocator.GetOriginalName(tt.generated)
			if got != tt.want {
				t.Errorf("GetOriginalName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetGeneratedName(t *testing.T) {
	allocator := NewAllocator()
	original := "count"
	generated := allocator.AllocName(original)

	tests := []struct {
		name     string
		original string
		want     string
	}{
		{
			name:     "existing name",
			original: original,
			want:     generated,
		},
		{
			name:     "non-existent name",
			original: "nonexistent",
			want:     "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := allocator.GetGeneratedName(tt.original)
			if got != tt.want {
				t.Errorf("GetGeneratedName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConcurrentAllocation(t *testing.T) {
	allocator := NewAllocator()
	var wg sync.WaitGroup
	mu := sync.Mutex{}
	seen := container.NewMapSet[string]()
	const numGoroutines = 10 // reduce the number of goroutines

	// Spawn multiple goroutines to allocate names concurrently
	for i := 0; i < numGoroutines; i++ {
		i := i // Create a new variable for each iteration
		wg.Add(1)
		go func() {
			defer wg.Done()
			// Use different original names for each goroutine
			originalName := fmt.Sprintf("test_%d", i)
			name := allocator.AllocName(originalName)

			mu.Lock()
			if seen.Contains(name) {
				t.Errorf("duplicate name generated: %s", name)
			}
			seen.Add(name)
			mu.Unlock()
		}()
	}

	// Wait for all goroutines to complete
	wg.Wait()

	// Verify we got exactly numGoroutines unique names
	if seen.Size() != numGoroutines {
		t.Errorf("expected %d unique names, got %d", numGoroutines, seen.Size())
	}
}

func TestAllocatorWithCustomPrefix(t *testing.T) {
	allocator := NewAllocator(WithPrefix("_var_"))
	name := allocator.AllocName("test")
	if name != "_var_1" {
		t.Errorf("expected '_var_1', got '%s'", name)
	}
}

func TestAllocatorWithCustomCounter(t *testing.T) {
	allocator := NewAllocator(WithInitialCounter(100))
	name := allocator.AllocName("test")
	if name != "_v101" {
		t.Errorf("expected '_v101', got '%s'", name)
	}
}
