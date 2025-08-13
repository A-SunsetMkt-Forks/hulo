// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package batch

type CallFrame interface {
	Kind() Caller
}

type Caller uint

const (
	CallerUnknown Caller = iota
	CallerFunction
	CallerClass
	CallerBlock
	CallerLoop
	CallerAssign
)

type LoopFrame struct {
	StartLabel string
	EndLabel   string
}

func (f *LoopFrame) Kind() Caller {
	return CallerLoop
}

func (f *LoopFrame) IsAnonymous() bool {
	return f.StartLabel == ""
}

type AssignFrame struct{}

func (*AssignFrame) Kind() Caller { return CallerAssign }

type FunctionFrame struct{}

func (*FunctionFrame) Kind() Caller { return CallerFunction }

type ClassFrame struct{}

func (*ClassFrame) Kind() Caller { return CallerClass }

type BlockFrame struct{}
