// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package build

import (
	"fmt"

	"github.com/hulo-lang/hulo/internal/build"
	"github.com/hulo-lang/hulo/internal/config"
	"github.com/hulo-lang/hulo/internal/container"
	hast "github.com/hulo-lang/hulo/syntax/hulo/ast"
	htok "github.com/hulo-lang/hulo/syntax/hulo/token"
	vast "github.com/hulo-lang/hulo/syntax/vbs/ast"
	vtok "github.com/hulo-lang/hulo/syntax/vbs/token"
)

func TranspileToVBScript(opts *config.VBScriptOptions, node hast.Node) (vast.Node, error) {
	tr := &VBScriptTranspiler{
		opts:         opts,
		scopeStack:   container.NewArrayStack[ScopeType](),
		declaredVars: container.NewMapSet[string](),
	}
	return tr.Convert(node), nil
}

type VBScriptTranspiler struct {
	opts *config.VBScriptOptions

	buffer []vast.Stmt

	scopeStack   container.Stack[ScopeType]
	declaredVars container.Set[string]

	alloc *build.Allocator
}

type ScopeType int

const (
	GlobalScope ScopeType = iota
	ForScope
	MatchScope
	IfScope
)

func (v *VBScriptTranspiler) pushScope(scope ScopeType) {
	v.scopeStack.Push(scope)
}

func (v *VBScriptTranspiler) popScope() {
	v.scopeStack.Pop()
}

func (v *VBScriptTranspiler) currentScope() (ScopeType, bool) {
	return v.scopeStack.Peek()
}

func (v *VBScriptTranspiler) Convert(node hast.Node) vast.Node {
	switch node := node.(type) {
	case *hast.File:
		return v.ConvertFile(node)
	case *hast.CommentGroup:
		docs := make([]*vast.Comment, len(node.List))
		var tok vtok.Token
		if v.opts.CommentSyntax == "'" {
			tok = vtok.SGL_QUOTE
		} else {
			tok = vtok.REM
		}

		for i, d := range node.List {
			docs[i] = &vast.Comment{
				Tok:  tok,
				Text: d.Text,
			}
		}
		return &vast.CommentGroup{List: docs}
	case *hast.ExprStmt:
		return &vast.ExprStmt{
			X: v.Convert(node.X).(vast.Expr),
		}
	case *hast.CallExpr:
		fun := v.Convert(node.Fun).(*vast.Ident)

		recv := []vast.Expr{}
		for _, r := range node.Recv {
			recv = append(recv, v.Convert(r).(vast.Expr))
		}

		return &vast.CallExpr{
			Func: fun,
			Recv: recv,
		}
	case *hast.Ident:
		return &vast.Ident{Name: node.Name}
	case *hast.BasicLit:
		return &vast.BasicLit{
			Kind:  Token(node.Kind),
			Value: node.Value,
		}
	case *hast.AssignStmt:
		return v.ConvertAssignStmt(node)
	case *hast.IfStmt:
		cond := v.Convert(node.Cond).(vast.Expr)
		body := v.Convert(node.Body).(*vast.BlockStmt)
		var elseBody *vast.BlockStmt
		if node.Else != nil {
			elseBody = v.Convert(node.Else).(*vast.BlockStmt)
		}
		return &vast.IfStmt{
			Cond: cond,
			Body: body,
			Else: elseBody,
		}
	case *hast.BlockStmt:
		return v.ConvertBlock(node)
	case *hast.BinaryExpr:
		x := v.Convert(node.X).(vast.Expr)
		y := v.Convert(node.Y).(vast.Expr)
		return &vast.BinaryExpr{
			X:  x,
			Op: Token(node.Op),
			Y:  y,
		}
	case *hast.RefExpr:
		return &vast.Ident{
			Name: node.X.String(),
		}
	case *hast.WhileStmt:
		body := v.Convert(node.Body).(*vast.BlockStmt)
		var (
			cond vast.Expr
			tok  vtok.Token
		)
		if node.Cond != nil {
			tok = vtok.WHILE
			cond = v.Convert(node.Cond).(vast.Expr)
		}
		return &vast.DoLoopStmt{
			Pre:  true,
			Tok:  tok,
			Cond: cond,
			Body: body,
		}
	case *hast.DoWhileStmt:
		body := v.Convert(node.Body).(*vast.BlockStmt)
		cond := v.Convert(node.Cond).(vast.Expr)
		return &vast.DoLoopStmt{
			Body: body,
			Tok:  vtok.WHILE,
			Cond: cond,
		}
	case *hast.ForStmt:
		return v.ConvertForStmt(node)
	// case *hast.RangeStmt:
	case *hast.IncDecExpr:
		return v.ConvertIncDecExpr(node)
	default:
		fmt.Printf("%T\n", node)
	}
	return nil
}

func (v *VBScriptTranspiler) ConvertBlock(node *hast.BlockStmt) vast.Node {
	stmts := []vast.Stmt{}
	for _, s := range node.List {
		stmt := v.Convert(s)
		stmts = append(stmts, v.Flush()...)
		stmts = append(stmts, stmt.(vast.Stmt))
	}
	return &vast.BlockStmt{List: stmts}
}

func (v *VBScriptTranspiler) ConvertIncDecExpr(node *hast.IncDecExpr) vast.Node {
	counter := v.Convert(node.X).(*vast.Ident)
	var op vtok.Token
	if node.Tok == htok.INC {
		op = vtok.ADD
	} else {
		op = vtok.SUB
	}

	return &vast.AssignStmt{
		Lhs: counter,
		Rhs: &vast.BinaryExpr{
			X:  counter,
			Op: op,
			Y: &vast.BasicLit{
				Kind:  vtok.INTEGER,
				Value: "1",
			},
		},
	}
}

func (v *VBScriptTranspiler) ConvertFile(node *hast.File) vast.Node {
	docs := make([]*vast.CommentGroup, len(node.Docs))
	for i, d := range node.Docs {
		docs[i] = v.Convert(d).(*vast.CommentGroup)
	}

	stmts := []vast.Stmt{}
	for _, s := range node.Stmts {
		stmt := v.Convert(s)

		stmts = append(stmts, v.Flush()...)

		stmts = append(stmts, stmt.(vast.Stmt))
	}

	return &vast.File{
		Doc:   docs,
		Stmts: stmts,
	}
}

func (v *VBScriptTranspiler) ConvertForStmt(node *hast.ForStmt) vast.Node {
	v.pushScope(ForScope)
	defer v.popScope()

	init := v.Convert(node.Init).(vast.Stmt)
	cond := v.Convert(node.Cond).(vast.Expr)
	post := v.Convert(node.Post).(vast.Stmt)

	stmts := []vast.Stmt{}
	body := v.Convert(node.Body).(*vast.BlockStmt)

	stmts = append(stmts, v.Flush()...)
	stmts = append(stmts, body.List...)
	stmts = append(stmts, post)

	v.Emit(init)

	return &vast.DoLoopStmt{
		Pre:  true,
		Tok:  vtok.WHILE,
		Cond: cond,
		Body: &vast.BlockStmt{List: stmts},
	}
}

func (v *VBScriptTranspiler) ConvertAssignStmt(node *hast.AssignStmt) vast.Node {
	lhs := v.Convert(node.Lhs).(vast.Expr)
	rhs := v.Convert(node.Rhs).(vast.Expr)
	if node.Scope == htok.CONST {
		return &vast.ConstStmt{
			Lhs: lhs.(*vast.Ident),
			Rhs: rhs.(*vast.BasicLit),
		}
	}

	if v.needsDimDeclaration(lhs) {
		v.Emit(&vast.DimDecl{List: []vast.Expr{lhs}})
	}
	return &vast.AssignStmt{
		Lhs: lhs,
		Rhs: rhs,
	}
}

func (v *VBScriptTranspiler) Emit(n ...vast.Stmt) {
	v.buffer = append(v.buffer, n...)
}

func (v *VBScriptTranspiler) Flush() []vast.Stmt {
	stmts := v.buffer
	v.buffer = nil
	return stmts
}

func (v *VBScriptTranspiler) needsDimDeclaration(expr vast.Expr) bool {
	// 检查当前作用域
	if scope, ok := v.currentScope(); ok {
		// 在循环作用域中不需要Dim声明
		if scope == ForScope {
			return false
		}
	}

	// 检查变量是否已经声明过
	if ident, ok := expr.(*vast.Ident); ok {
		if v.declaredVars.Contains(ident.Name) {
			return false
		}
		// 标记为已声明
		v.declaredVars.Add(ident.Name)
		return true
	}

	return false
}
