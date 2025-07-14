// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package ast_test

import (
	"testing"

	"github.com/hulo-lang/hulo/syntax/powershell/ast"
	"github.com/hulo-lang/hulo/syntax/powershell/token"
)

/*
	Function MyFunction {
	    [CmdletBinding()]
	    Param (
	        [Parameter(ValueFromPipeline=$true)]
	        [string]$Name
	    )

	    Process {
	        Write-Verbose "Processing $Name"
	        Write-Output "Hello, $Name!"
	    }
	}
*/
func TestFuncDecl(t *testing.T) {
	ast.Print(&ast.FuncDecl{
		Name: &ast.Ident{Name: "MyFunction"},
		Attributes: []*ast.Attribute{
			{Name: &ast.Ident{Name: "CmdletBinding"}},
		},
		Params: []ast.Expr{
			&ast.Parameter{
				X: &ast.ConstrainedVarExpr{
					Type: &ast.Ident{Name: "string"},
					X:    &ast.Ident{Name: "Name"},
				},
			},
		},
		Body: &ast.BlcokStmt{
			List: []ast.Stmt{
				&ast.ProcessDecl{
					Body: &ast.BlcokStmt{
						List: []ast.Stmt{
							&ast.ExprStmt{
								X: &ast.CmdExpr{
									Cmd: &ast.Ident{Name: "Write-Verbose"},
									Args: []ast.Expr{
										&ast.Lit{Val: `"Processing $Name"`},
									},
								},
							},
							&ast.ExprStmt{
								X: &ast.CmdExpr{
									Cmd: &ast.Ident{Name: "Write-Output"},
									Args: []ast.Expr{
										&ast.Lit{Val: `"Hello, $Name!"`},
									},
								},
							},
						},
					},
				},
			},
		},
	})
}

func TestHashTable(t *testing.T) {

}

/*
	try {
	    $a[$i] = 10
	    "Assignment completed without error"
	    break
	}

	catch [IndexOutOfRangeException] {
	    "Handling out-of-bounds index, >$_<`n"
	    $i = 5
	}

	catch {
	    "Caught unexpected exception"
	}

	finally {
	    # ...
	}
*/
func TestTryStmt(t *testing.T) {
	ast.Print(&ast.TryStmt{
		Body: &ast.BlcokStmt{
			List: []ast.Stmt{&ast.ExprStmt{X: &ast.Lit{Val: `"try"`}}},
		},
		Catches: []*ast.CatchClause{
			{
				Type: &ast.Ident{Name: "IndexOutOfRangeException"},
				Body: &ast.BlcokStmt{
					List: []ast.Stmt{&ast.ExprStmt{X: &ast.Lit{Val: "Handling out-of-bounds index, >$_<`n"}}},
				},
			},
			{
				Body: &ast.BlcokStmt{
					List: []ast.Stmt{&ast.ExprStmt{X: &ast.Lit{Val: "Caught unexpected exception"}}},
				},
			},
		},
		FinallyBody: &ast.BlcokStmt{
			List: []ast.Stmt{&ast.ExprStmt{X: &ast.Lit{Val: "# ..."}}},
		},
	})
}

/*
	trap {
	    "Caught unexpected exception"
	}
*/
func TestTrapStmt(t *testing.T) {
	ast.Print(&ast.TrapStmt{
		Body: &ast.BlcokStmt{List: []ast.Stmt{&ast.ExprStmt{X: &ast.Lit{Val: "# ..."}}}},
	})
}

/*
	data -SupportedCommand ConvertTo-Xml {
	    Format-Xml -Strings string1, string2, string3
	}
*/
func TestDataStmt(t *testing.T) {
	ast.Print(&ast.DataStmt{
		Recv: []ast.Expr{&ast.Ident{Name: "-SupportedCommand"}, &ast.Ident{Name: "ConvertTo-Xml"}},
		Body: &ast.BlcokStmt{List: []ast.Stmt{&ast.ExprStmt{X: &ast.CmdExpr{Cmd: &ast.Ident{Name: "Format-Xml"}, Args: []ast.Expr{&ast.Lit{Val: "string1"}, &ast.Lit{Val: "string2"}, &ast.Lit{Val: "string3"}}}}}},
	})
}

/*
	switch -Wildcard ("abc") {
	    a* { "a*, $_" }
	    ?B? { "?B? , $_" }
	    default { "default, $_" }
	}

	switch -Regex -CaseSensitive ("abc") {
	    ^a* { "a*" }
	    ^A* { "A*" }
	}

	switch (0, 1, 19, 20, 21) {
	    { $_ -lt 20 } { "-lt 20" }
	    { $_ -band 1 } { "Odd" }
	    { $_ -eq 19 } { "-eq 19" }
	    default { "default" }
	}
*/
func TestSwitchStmt(t *testing.T) {
	ast.Print(&ast.SwitchStmt{
		Pattern: ast.SwitchPatternWildcard,
		Value:   &ast.Lit{Val: `"abc"`},
		Cases: []*ast.CaseClause{
			{
				Cond: &ast.Ident{Name: "a*"},
				Body: &ast.BlcokStmt{
					List: []ast.Stmt{
						&ast.ExprStmt{
							X: &ast.Lit{Val: `"a*, $_"`},
						},
					},
				},
			},
			{
				Cond: &ast.Ident{Name: "?B?"},
				Body: &ast.BlcokStmt{
					List: []ast.Stmt{
						&ast.ExprStmt{
							X: &ast.Lit{Val: `"?B? , $_"`},
						},
					},
				},
			},
		},
		Default: &ast.BlcokStmt{
			List: []ast.Stmt{
				&ast.ExprStmt{X: &ast.Lit{Val: `"default, $_"`}},
			},
		},
	})
}

func TestAssigStmt(t *testing.T) {
	ast.Print(&ast.AssignStmt{
		Lhs: &ast.IndexExpr{
			X:     &ast.VarExpr{X: &ast.Ident{Name: "x"}},
			Index: &ast.IncDecExpr{X: &ast.VarExpr{X: &ast.Ident{Name: "i"}}, Tok: token.INC},
		},
		Rhs: &ast.IndexExpr{
			X:     &ast.VarExpr{X: &ast.Ident{Name: "x"}},
			Index: &ast.IncDecExpr{Pre: true, X: &ast.VarExpr{X: &ast.Ident{Name: "i"}}, Tok: token.DEC},
		},
	})
}

// @(10, "blue", 12.54e3, 16.30D)
func TestArrayExpr(t *testing.T) {
	ast.Print(&ast.ArrayExpr{
		Elems: []ast.Expr{&ast.Lit{Val: "10"}, &ast.Lit{Val: `"blue"`}, &ast.Lit{Val: "12.54e3"}, &ast.Lit{Val: "16.30D"}},
	})
}

func TestComments(t *testing.T) {
	ast.Print(&ast.CommentGroup{
		List: []ast.Comment{
			&ast.SingleLineComment{
				Text: "This is a single line comment",
			},
			&ast.DelimitedComment{
				Text: `This is a
	delimited comment`,
			},
		},
	})
}
