// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package exec

type Executor interface {
	CanHandle(cmd string) bool
	Execute(cmd string) error
}

