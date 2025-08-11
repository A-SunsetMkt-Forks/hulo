// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package transpiler

type SymbolType int

const (
	SymbolVariable SymbolType = iota
	SymbolFunction
	SymbolClass
	SymbolConstant
	SymbolArray
	SymbolObject
)

type Scope struct {
	Type ScopeType // 作用域类型
}

type ScopeStack struct {
	scopes []*Scope
}

func (ss *ScopeStack) Push(scope *Scope) {
	ss.scopes = append(ss.scopes, scope)
}

func (ss *ScopeStack) Pop() {
	if len(ss.scopes) > 0 {
		ss.scopes = ss.scopes[:len(ss.scopes)-1]
	}
}

func (ss *ScopeStack) Current() *Scope {
	if len(ss.scopes) == 0 {
		return nil
	}
	return ss.scopes[len(ss.scopes)-1]
}

type ScopeType int

const (
	GlobalScope ScopeType = iota
	FunctionScope
	ClassScope
	BlockScope
	LoopScope
	AssignScope
)
