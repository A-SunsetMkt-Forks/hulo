package ast

import (
	"testing"

	"github.com/hulo-lang/hulo/internal/token"
)

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
