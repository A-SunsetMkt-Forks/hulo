// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package util

import "sync"

var _ sync.Locker = (*NoOpLocker)(nil)

// NoOpLocker is a virtual implementation of Locker
//
// This is a no-op implementation that does nothing when Lock() or Unlock() is called.
// It's useful for testing or when you want to maintain the locking interface
// without actual synchronization.
type NoOpLocker struct{}

// Lock is a no-op implementation
func (n *NoOpLocker) Lock()   {}

// Unlock is a no-op implementation
func (n *NoOpLocker) Unlock() {}
