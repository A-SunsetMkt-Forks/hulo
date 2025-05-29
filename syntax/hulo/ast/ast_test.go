package ast

import (
	"testing"

	"github.com/hulo-lang/hulo/syntax/hulo/token"
)

func TestAssignStmt(t *testing.T) {
	Print(&AssignStmt{
		Tok: token.LET,
		Lhs: &Ident{Name: "x"},
		Rhs: &BasicLit{
			Kind:  token.NUM,
			Value: "123",
		},
	})
}

func TestIfStmt(t *testing.T) {
	Print(&IfStmt{
		Cond: &BinaryExpr{
			X:  &RefExpr{X: &Ident{Name: "x"}},
			Op: token.LT,
			Y: &BasicLit{
				Kind:  token.NUM,
				Value: "10",
			},
		},
		Body: &BlockStmt{
			List: []Stmt{
				&AssignStmt{
					Tok: token.LET,
					Lhs: &Ident{Name: "x"},
					Rhs: &BasicLit{
						Kind:  token.NUM,
						Value: "10",
					},
				},
			},
		},
		Else: &IfStmt{
			Cond: &BinaryExpr{
				X:  &RefExpr{X: &Ident{Name: "x"}},
				Op: token.LT,
				Y: &BasicLit{
					Kind:  token.NUM,
					Value: "20",
				},
			},
			Body: &BlockStmt{
				List: []Stmt{
					&AssignStmt{
						Tok: token.LET,
						Lhs: &Ident{Name: "x"},
						Rhs: &BasicLit{
							Kind:  token.NUM,
							Value: "20",
						},
					},
				},
			},
			Else: &BlockStmt{
				List: []Stmt{
					&AssignStmt{
						Tok: token.LET,
						Lhs: &Ident{Name: "x"},
						Rhs: &BasicLit{
							Kind:  token.NUM,
							Value: "20",
						},
					},
				},
			},
		},
	})
}
func TestDeclareDecl(t *testing.T) {
	Print(&File{
		Decls: []Decl{
			&DeclareDecl{
				X: &ComptimeStmt{
					X: &AssignStmt{
						Tok: token.CONST,
						Lhs: &Ident{Name: "os"},
					},
				},
			},
			&DeclareDecl{
				X: &ClassDecl{
					Name: &Ident{Name: "User"},
				},
			},
		},
	})
}
