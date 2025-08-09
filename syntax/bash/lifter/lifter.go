// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package lifter

import (
	"fmt"

	bast "github.com/hulo-lang/hulo/syntax/bash/ast"
	"github.com/hulo-lang/hulo/syntax/bash/parser"
	hast "github.com/hulo-lang/hulo/syntax/hulo/ast"
)

type Lifter struct {
	parser *parser.Parser
}

func NewLifter() *Lifter {
	return &Lifter{
		parser: parser.NewParser(),
	}
}

func (l *Lifter) Lift(script string) (hast.Node, error) {
	file, err := l.parser.Parse(script)
	if err != nil {
		return nil, err
	}
	return l.convert(file), nil
}

func (l *Lifter) convert(node bast.Node) hast.Node {
	switch node := node.(type) {
	case *bast.File:
		return l.convertFile(node)
	case *bast.ExprStmt:
		return l.convertExprStmt(node)
	case *bast.CmdExpr:
		return l.convertCmdExpr(node)
	case *bast.Ident:
		return &hast.Ident{Name: node.Name}
	case *bast.Word:
		return &hast.Ident{Name: node.Val}
	case *bast.FuncDecl:
		return l.convertFuncDecl(node)
	case *bast.AssignStmt:
		return l.convertAssignStmt(node)
	default:
		panic(fmt.Sprintf("unknown node type: %T", node))
	}
}

func (l *Lifter) convertFile(file *bast.File) hast.Node {
	ret := &hast.File{}
	for _, stmt := range file.Stmts {
		ret.Stmts = append(ret.Stmts, l.convert(stmt).(hast.Stmt))
	}
	return ret
}

func (l *Lifter) convertExprStmt(stmt *bast.ExprStmt) hast.Node {
	return &hast.ExprStmt{
		X: l.convert(stmt.X).(hast.Expr),
	}
}

func (l *Lifter) convertCmdExpr(cmd *bast.CmdExpr) hast.Node {
	fun := l.convert(cmd.Name).(hast.Expr)
	args := make([]hast.Expr, 0, len(cmd.Recv))
	for _, arg := range cmd.Recv {
		args = append(args, l.convert(arg).(hast.Expr))
	}
	return &hast.CallExpr{
		Fun:  fun,
		Recv: args,
	}
}

func (l *Lifter) convertFuncDecl(decl *bast.FuncDecl) hast.Node {
	body := &hast.BlockStmt{}
	for _, stmt := range decl.Body.List {
		body.List = append(body.List, l.convert(stmt).(hast.Stmt))
	}
	return &hast.FuncDecl{
		Name: &hast.Ident{Name: decl.Name.Name},
		Body: body,
	}
}

func (l *Lifter) convertAssignStmt(stmt *bast.AssignStmt) hast.Node {
	return &hast.AssignStmt{
		Lhs: l.convert(stmt.Lhs).(hast.Expr),
		Rhs: l.convert(stmt.Rhs).(hast.Expr),
	}
}
