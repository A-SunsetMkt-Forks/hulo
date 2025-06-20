// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package parser

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/hulo-lang/hulo/internal/container"
	"github.com/hulo-lang/hulo/syntax/hulo/ast"
)

type Position struct {
	Line   int // line number, starting at 1
	Column int // column number, starting at 1 (byte count)
}

type TraceLevel int

const (
	DebugNone TraceLevel = iota
	DebugParseTree
	DebugAST
	DebugSymbols
	DebugAll
)

type Tracer struct {
	level         TraceLevel
	writer        io.Writer
	depth         int
	stack         container.Stack[nodeInfo]
	timings       map[string]time.Duration
	startTime     time.Time
	ignoredNodes  container.Set[string]  // nodes to ignore
	disableTiming bool                   // whether to disable timing
	nodeWatchers  map[string]NodeWatcher // node watchers
}

type nodeInfo struct {
	name     string
	nodeType string
	pos      Position
	start    time.Time
}

// NodeWatcher defines the node watcher interface
type NodeWatcher interface {
	OnEnter(node ast.Node, pos Position)
	OnExit(node ast.Node, result any, err error)
}

// defaultNodeWatcher is the default implementation of NodeWatcher
type defaultNodeWatcher struct {
	onEnter func(node ast.Node, pos Position)
	onExit  func(node ast.Node, result any, err error)
}

func (w *defaultNodeWatcher) OnEnter(node ast.Node, pos Position) {
	if w.onEnter != nil {
		w.onEnter(node, pos)
	}
}

func (w *defaultNodeWatcher) OnExit(node ast.Node, result any, err error) {
	if w.onExit != nil {
		w.onExit(node, result, err)
	}
}

func NewTracer() *Tracer {
	return &Tracer{
		writer:       os.Stdout,
		timings:      make(map[string]time.Duration),
		level:        DebugAll,
		stack:        container.NewArrayStack[nodeInfo](),
		ignoredNodes: container.NewMapSet[string](),
		nodeWatchers: make(map[string]NodeWatcher),
	}
}

func (d *Tracer) SetOutput(w io.Writer) {
	d.writer = w
}

func (d *Tracer) SetLevel(level TraceLevel) {
	d.level = level
}

func (d *Tracer) indent() string {
	if d.stack.IsEmpty() {
		return ""
	}
	return strings.Repeat("│   ", d.stack.Size()-1) + "├── "
}

// TODO handles error at nodeType's type
func (d *Tracer) Enter(name string, nodeType string, pos Position) {
	if d == nil || d.writer == nil {
		return
	}

	if d.level == DebugNone {
		return
	}

	// Check if this node should be ignored
	if d.ignoredNodes.Contains(name) {
		return
	}

	info := nodeInfo{
		name:     name,
		nodeType: nodeType,
		pos:      pos,
		start:    time.Now(),
	}
	d.stack.Push(info)

	// Print enter information
	fmt.Fprintf(d.writer, "%senter: %s (%s at %d:%d)\n",
		d.indent(),
		name,
		nodeType,
		pos.Line,
		pos.Column,
	)

	// Call node watcher
	if watcher, ok := d.nodeWatchers[name]; ok {
		watcher.OnEnter(nil, pos) // TODO: pass actual ast.Node
	}
}

func (d *Tracer) Exit(result any, err error) {
	if d == nil || d.writer == nil || d.stack.IsEmpty() {
		return
	}

	if d.level == DebugNone {
		return
	}

	info, _ := d.stack.Peek()

	// Check if this node should be ignored
	if d.ignoredNodes.Contains(info.name) || d.ignoredNodes.Contains(info.nodeType) {
		d.stack.Pop()
		return
	}

	indent := strings.Repeat("│   ", d.stack.Size()-1) + "└── "

	// Calculate duration
	var duration time.Duration
	if !d.disableTiming {
		duration = time.Since(info.start)
		d.timings[info.name] = duration
	}

	// Print exit information
	if err != nil {
		fmt.Fprintf(d.writer, "%sexit ❌ %s: %v\n", indent, info.name, err)
		d.printErrorContext()
	} else {
		resultStr := d.formatResult(result)
		durationStr := ""
		if !d.disableTiming {
			durationStr = fmt.Sprintf(" (%s)", duration)
		}
		fmt.Fprintf(d.writer, "%sexit ✓ %s → %s%s\n",
			indent,
			info.name,
			resultStr,
			durationStr,
		)
	}

	// Call node watcher
	if watcher, ok := d.nodeWatchers[info.name]; ok {
		watcher.OnExit(nil, result, err) // TODO: pass actual ast.Node
	}

	d.stack.Pop()

	// If this is the outermost node and timing is enabled, print summary
	if d.stack.IsEmpty() && !d.disableTiming {
		d.printSummary()
	}
}

func (d *Tracer) EmitError(err error) {
	if d == nil || d.writer == nil {
		return
	}
	info, _ := d.stack.Peek()
	indent := strings.Repeat("│   ", d.stack.Size()-1) + "└── "
	fmt.Fprintf(d.writer, "%sexit ❌ %s: %v\n", indent, info.name, err)
	d.printErrorContext()
}

func (d *Tracer) formatResult(result any) string {
	switch v := result.(type) {
	case *ast.BasicLit:
		return fmt.Sprintf("%s(%s)", v.Kind, v.Value)
	case *ast.Ident:
		return fmt.Sprintf("Identifier(%s)", v.Name)
	case *ast.BinaryExpr:
		return fmt.Sprintf("BinaryExpr(%s)", v.Op)
	case nil:
		return "nil"
	default:
		return fmt.Sprintf("%T", v)
	}
}

func (d *Tracer) printErrorContext() {
	fmt.Fprintf(d.writer, "\n[ERROR STACK]\n")

	// Clone the stack for traversal
	stackClone := d.stack.Clone()
	var stackItems []nodeInfo

	// Extract all items from the cloned stack
	for !stackClone.IsEmpty() {
		if item, ok := stackClone.Pop(); ok {
			stackItems = append(stackItems, item) // Append to get reverse order
		}
	}

	// Print error stack in reverse call order (most recent call first)
	for i, info := range stackItems {
		prefix := "├── "
		if i == len(stackItems)-1 {
			prefix = "└── "
		}
		fmt.Fprintf(d.writer, "%s%s (%s at line %d, col %d)\n",
			prefix,
			info.name,
			info.nodeType,
			info.pos.Line,
			info.pos.Column,
		)
	}
	fmt.Fprintln(d.writer)
}

// timingEntry represents a timing entry for the heap
type timingEntry struct {
	name     string
	duration time.Duration
}

// timingEntryLess compares timing entries for a MIN heap (shorter duration is "less")
func timingEntryLess(a, b timingEntry) bool {
	return a.duration < b.duration
}

func (d *Tracer) printSummary() {
	fmt.Fprintf(d.writer, "\n[SUMMARY]\n")

	// Print total time
	var total time.Duration
	for _, t := range d.timings {
		total += t
	}
	fmt.Fprintf(d.writer, "├── Total time: %s\n", total)

	// --- OPTIMIZED TOP-K ALGORITHM (O(N log K)) ---
	maxCount := min(3, len(d.timings))
	if maxCount == 0 {
		// Also handles the case where d.timings is empty
		fmt.Fprintln(d.writer)
		return
	}

	// Use a min-heap to find the top K slowest nodes efficiently.
	// The heap will store the K slowest nodes found so far.
	// The top of the min-heap is the *smallest* of the K slowest nodes.
	heap := container.NewHeap(timingEntryLess, container.MinHeap)

	for name, duration := range d.timings {
		entry := timingEntry{name: name, duration: duration}
		if heap.Size() < maxCount {
			heap.Push(entry)
		} else {
			// If the current entry is slower than the "smallest of the slowest" in our heap,
			// replace that smallest one.
			if peek, ok := heap.Peek(); ok && entry.duration > peek.duration {
				heap.Pop()       // Remove the smallest of the top K
				heap.Push(entry) // Add the new, slower one
			}
		}
	}

	// The heap now contains the top K slowest nodes, but in min-heap order (smallest at the top).
	// To print from slowest to fastest, we pop them all into a slice and then iterate that slice backwards.
	slowestNodes := make([]timingEntry, heap.Size())
	for i := heap.Size() - 1; i >= 0; i-- {
		entry, _ := heap.Pop()
		slowestNodes[i] = entry
	}

	fmt.Fprintf(d.writer, "└── Slowest nodes:\n")
	for i, entry := range slowestNodes {
		prefix := "    ├──"
		if i == len(slowestNodes)-1 {
			prefix = "    └──"
		}
		fmt.Fprintf(d.writer, "%s %s: %s\n", prefix, entry.name, entry.duration)
	}
	fmt.Fprintln(d.writer)
}

func (d *Tracer) IgnoreNodes(nodes ...string) {
	for _, node := range nodes {
		d.ignoredNodes.Add(node)
	}
}

func (d *Tracer) DisableTiming() {
	d.disableTiming = true
}

func (d *Tracer) WatchNode(nodeType string, watcher NodeWatcher) {
	d.nodeWatchers[nodeType] = watcher
}
