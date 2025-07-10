// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package parser

import (
	"fmt"
	"io"

	"github.com/caarlos0/log"
	"github.com/hulo-lang/hulo/syntax/hulo/ast"
)

type ParserOptions func(*Analyzer) error

func OptionDisableTracer() ParserOptions {
	return func(a *Analyzer) error {
		a.Tracer.writer = io.Discard
		return nil
	}
}

func OptionTracerLevel(level log.Level) ParserOptions {
	return func(a *Analyzer) error {
		log.SetLevel(level)
		return nil
	}
}

func OptionTracerLevelOf(s string) ParserOptions {
	return func(a *Analyzer) error {
		l, err := log.ParseLevel(s)
		if err != nil {
			return err
		}
		log.SetLevel(l)
		return nil
	}
}

func OptionTracerIgnore(nodes ...string) ParserOptions {
	return func(a *Analyzer) error {
		if a.Tracer != nil {
			a.Tracer.IgnoreNodes(nodes...)
		}
		return nil
	}
}

func OptionTracerDisableTiming() ParserOptions {
	return func(a *Analyzer) error {
		if a.Tracer != nil {
			a.Tracer.DisableTiming()
		}
		return nil
	}
}

func OptionTracerWatchNode(nodeType string, onEnter func(ast.Node, Position), onExit func(ast.Node, any, error)) ParserOptions {
	return func(a *Analyzer) error {
		if a.Tracer != nil {
			watcher := &defaultNodeWatcher{
				onEnter: onEnter,
				onExit:  onExit,
			}
			a.Tracer.WatchNode(nodeType, watcher)
		}
		return nil
	}
}

func OptionTracerWriter(w io.Writer) ParserOptions {
	return func(a *Analyzer) error {
		a.writer = w
		return nil
	}
}

func OptionTracerASTTree(w io.Writer) ParserOptions {
	return func(a *Analyzer) error {
		fmt.Fprintln(w, a.file.ToStringTree(nil, a.parser))
		return nil
	}
}
