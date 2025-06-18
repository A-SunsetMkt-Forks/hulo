package ast

import (
	"fmt"
	"testing"

	"github.com/hulo-lang/hulo/syntax/batch/token"
)

func TestAssign(t *testing.T) {
	fmt.Println(IfStmt{
		Cond: &Word{
			Parts: []Expr{
				&UnaryExpr{
					Op: token.DIV,
					X: &Lit{
						Val: "f",
					},
				},
			},
		},
	})
}
