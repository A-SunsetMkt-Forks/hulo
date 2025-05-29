package build

import (
	"fmt"

	"github.com/hulo-lang/hulo/internal/config"
	bast "github.com/hulo-lang/hulo/syntax/bash/ast"
	hast "github.com/hulo-lang/hulo/syntax/hulo/ast"
)

func Translate(opts *config.BashOptions, node hast.Node) (bast.Node, error) {
	bnode := translate2Bash(opts, node)

	fmt.Println("write to file system", bnode)
	return bnode, nil
}

func translate2Bash(opts *config.BashOptions, node hast.Node) bast.Node {
	switch node := node.(type) {
	case *hast.File:
		decls := make([]bast.Decl, len(node.Decls))
		for i, d := range node.Decls {
			decls[i] = translate2Bash(opts, d).(bast.Decl)
		}

		stmts := make([]bast.Stmt, len(node.Stmts))
		for i, s := range node.Stmts {
			stmts[i] = translate2Bash(opts, s).(bast.Stmt)
		}

		return &bast.File{
			Decls: decls,
			Stmts: stmts,
		}
	case *hast.IfStmt:
		fmt.Println(opts.Boolean)
		return &bast.IfStmt{
			// TODO HCR
			Cond: &bast.TestExpr{
				X: translate2Bash(opts, node.Cond).(bast.Expr),
			},
			Body: translate2Bash(opts, node.Body).(*bast.BlockStmt),
		}
	case *hast.BlockStmt:
		return &bast.BlockStmt{}
	case *hast.BinaryExpr:
		return &bast.BinaryExpr{
			X:  translate2Bash(opts, node.X).(bast.Expr),
			Op: Token(node.Op),
			Y:  translate2Bash(opts, node.Y).(bast.Expr),
		}
	case *hast.RefExpr:
		return &bast.VarExpExpr{
			X: translate2Bash(opts, node.X).(*bast.Ident),
		}
	case *hast.Ident:
		return &bast.Ident{
			Name: node.Name,
		}
	case *hast.BasicLit:
		return &bast.BasicLit{
			Kind:  Token(node.Kind),
			Value: node.Value,
		}
	}

	return nil
}
