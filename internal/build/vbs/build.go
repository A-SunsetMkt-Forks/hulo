// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package build

import (
	"fmt"

	"github.com/hulo-lang/hulo/internal/config"
	hast "github.com/hulo-lang/hulo/syntax/hulo/ast"
	htok "github.com/hulo-lang/hulo/syntax/hulo/token"
	vast "github.com/hulo-lang/hulo/syntax/vbs/ast"
	vtok "github.com/hulo-lang/hulo/syntax/vbs/token"
)

func Translate(opts *config.VBScriptOptions, node hast.Node) (vast.Node, error) {
	vnode := translate2VBScript(opts, node)
	return vnode, nil
}

func translate2VBScript(opts *config.VBScriptOptions, node hast.Node) vast.Node {
	switch node := node.(type) {
	case *hast.File:
		docs := make([]*vast.CommentGroup, len(node.Docs))
		for i, d := range node.Docs {
			docs[i] = translate2VBScript(opts, d).(*vast.CommentGroup)
		}

		stmts := []vast.Stmt{}
		for _, s := range node.Stmts {
			stmt := translate2VBScript(opts, s)
			if wrapper, ok := stmt.(*CompilerWrapper); ok {
				stmts = append(stmts, wrapper.Result...)
			} else {
				stmts = append(stmts, stmt.(vast.Stmt))
			}
		}

		return &vast.File{
			Doc:   docs,
			Stmts: stmts,
		}
	case *hast.CommentGroup:
		docs := make([]*vast.Comment, len(node.List))
		var tok vtok.Token
		if opts.CommentSyntax == "'" {
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
			X: translate2VBScript(opts, node.X).(vast.Expr),
		}
	case *hast.CallExpr:
		fun := translate2VBScript(opts, node.Fun).(*vast.Ident)

		recv := []vast.Expr{}
		for _, r := range node.Recv {
			recv = append(recv, translate2VBScript(opts, r).(vast.Expr))
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
		lhs := translate2VBScript(opts, node.Lhs).(vast.Expr)
		rhs := translate2VBScript(opts, node.Rhs).(vast.Expr)
		if node.Scope == htok.CONST {
			return &vast.ConstStmt{
				Lhs: lhs.(*vast.Ident),
				Rhs: rhs.(*vast.BasicLit),
			}
		}
		return &CompilerWrapper{
			Result: []vast.Stmt{
				&vast.DimDecl{
					List: []vast.Expr{lhs},
				},
				&vast.AssignStmt{
					Lhs: lhs,
					Rhs: rhs,
				},
			},
		}
	case *hast.IfStmt:
		cond := translate2VBScript(opts, node.Cond).(vast.Expr)
		body := translate2VBScript(opts, node.Body).(*vast.BlockStmt)
		var elseBody *vast.BlockStmt
		if node.Else != nil {
			elseBody = translate2VBScript(opts, node.Else).(*vast.BlockStmt)
		}
		return &vast.IfStmt{
			Cond: cond,
			Body: body,
			Else: elseBody,
		}
	case *hast.BlockStmt:
		stmts := []vast.Stmt{}
		for _, s := range node.List {
			stmt := translate2VBScript(opts, s)
			stmts = append(stmts, stmt.(vast.Stmt))
		}
		return &vast.BlockStmt{List: stmts}
	case *hast.BinaryExpr:
		x := translate2VBScript(opts, node.X).(vast.Expr)
		y := translate2VBScript(opts, node.Y).(vast.Expr)
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
		body := translate2VBScript(opts, node.Body).(*vast.BlockStmt)
		var (
			cond vast.Expr
			tok  vtok.Token
		)
		if node.Cond != nil {
			tok = vtok.WHILE
			cond = translate2VBScript(opts, node.Cond).(vast.Expr)
		}
		return &vast.DoLoopStmt{
			Pre:  true,
			Tok:  tok,
			Cond: cond,
			Body: body,
		}
	case *hast.DoWhileStmt:
		body := translate2VBScript(opts, node.Body).(*vast.BlockStmt)
		cond := translate2VBScript(opts, node.Cond).(vast.Expr)
		return &vast.DoLoopStmt{
			Body: body,
			Tok:  vtok.WHILE,
			Cond: cond,
		}
	case *hast.ForStmt:
		init := translate2VBScript(opts, node.Init).(vast.Expr)
		cond := translate2VBScript(opts, node.Cond).(vast.Expr)
		// post := translate2VBScript(opts, node.Post).(vast.Expr)
		body := translate2VBScript(opts, node.Body).(*vast.BlockStmt)
		return &vast.ForNextStmt{
			Start: init,
			End_:  cond,
			Body:  body,
		}
	// case *hast.RangeStmt:
	default:
		fmt.Printf("%T\n", node)
	}
	return nil
}

type CompilerWrapper struct {
	Result []vast.Stmt
}

func (*CompilerWrapper) Pos() vtok.Pos { return vtok.NoPos }
func (*CompilerWrapper) End() vtok.Pos { return vtok.NoPos }
