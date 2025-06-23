package ast

import (
	"testing"

	"github.com/hulo-lang/hulo/syntax/hulo/token"
)

func TestAssignStmt(t *testing.T) {
	Print(&AssignStmt{
		Scope: token.LET,
		Lhs:   &Ident{Name: "x"},
		Rhs: &BasicLit{
			Kind:  token.NUM,
			Value: "123",
		},
	})

	Print(&AssignStmt{
		Lhs: &Ident{Name: "x"},
		Tok: token.COLON_ASSIGN,
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
					Scope: token.LET,
					Lhs:   &Ident{Name: "x"},
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
						Scope: token.LET,
						Lhs:   &Ident{Name: "x"},
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
						Scope: token.LET,
						Lhs:   &Ident{Name: "x"},
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
		Decls: []Stmt{
			&DeclareDecl{
				X: &ComptimeStmt{
					X: &AssignStmt{
						Scope: token.CONST,
						Lhs:   &Ident{Name: "os"},
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

func TestModDecl(t *testing.T) {
	Print(&ModDecl{
		Name: &Ident{Name: "fruit"},
		Body: &BlockStmt{
			List: []Stmt{
				&ModDecl{
					Name: &Ident{Name: "apple"},
					Body: &BlockStmt{
						List: []Stmt{
							&AssignStmt{
								Scope: token.VAR,
								Lhs:   &Ident{Name: "count"},
								Type: &UnionType{
									Types: []Expr{
										&TypeReference{
											Name: &Ident{Name: "str"},
										},
										&TypeReference{
											Name: &Ident{Name: "num"},
										},
										&BasicLit{
											Kind:  token.TRUE,
											Value: "true",
										},
									},
								},
								Rhs: &BasicLit{
									Kind:  token.NUM,
									Value: "10",
								},
							},
						},
					},
				},
				&ModDecl{
					Name: &Ident{Name: "pine"},
					Body: &BlockStmt{
						List: []Stmt{
							&UseDecl{Lhs: &Ident{Name: "PI"}},
						},
					},
				},
			},
		},
	})
}

func TestCascadeExpr(t *testing.T) {
	Print(&AssignStmt{
		Scope: token.LET,
		Lhs:   &Ident{Name: "user"},
		Rhs: &CascadeExpr{
			X: &CascadeExpr{
				X: &CallExpr{
					Fun: &Ident{Name: "User"},
				},
				Y: &BinaryExpr{
					X:  &Ident{Name: "name"},
					Op: token.ASSIGN,
					Y:  &BasicLit{Kind: token.STR, Value: "Alice"},
				},
			},
			Y: &BinaryExpr{
				X:  &Ident{Name: "age"},
				Op: token.ASSIGN,
				Y:  &BasicLit{Kind: token.NUM, Value: "25"},
			},
		},
	})
}

func TestArrayExpr(t *testing.T) {
	Print(&File{
		Stmts: []Stmt{
			&AssignStmt{
				Scope: token.LET,
				Lhs:   &Ident{Name: "arr"},
				Type: &TypeReference{
					Name: &Ident{Name: "list"},
					TypeParams: []Expr{
						&TypeReference{
							Name: &Ident{Name: "num"},
						},
					},
				},
				Rhs: &ArrayLiteralExpr{
					Elems: []Expr{
						&NumericLiteral{
							Value: "10",
						},
						&NumericLiteral{
							Value: "3.14",
						},
						&NumericLiteral{
							Value: "1e3",
						},
						&NumericLiteral{
							Value: "-2.5e-4",
						},
						&NumericLiteral{
							Value: "-2.5e-4",
						},
						&NumericLiteral{
							Value: "0o77",
						},
					},
				},
			},
			&ExprStmt{
				X: &IndexExpr{
					X:     &StringLiteral{Value: "abcde"},
					Index: &NumericLiteral{Value: "3"},
				},
			},
			&ExprStmt{
				X: &SliceExpr{
					X:    &StringLiteral{Value: "abcde"},
					Low:  &NumericLiteral{Value: "4"},
					High: &NumericLiteral{Value: "0"},
				},
			},
			&ExprStmt{
				X: &UnaryExpr{
					Op: token.HASH,
					X:  &StringLiteral{Value: "abc"},
				},
			},
			&ExprStmt{
				X: &BinaryExpr{
					X:  &StringLiteral{Value: "a"},
					Op: token.IN,
					Y:  &StringLiteral{Value: "abc"},
				},
			},
			&ExprStmt{
				X: &UnaryExpr{
					Op: token.Opr,
					X: &StringLiteral{
						Value: "C:\\Users\\name",
					},
				},
			},
		},
	})
}

func TestExtends(t *testing.T) {
	Print(&File{
		Stmts: []Stmt{
			&ExtensionDecl{
				Name: &Ident{Name: "str"},
				ExtensionClass: &ExtensionClass{
					Body: &BlockStmt{
						List: []Stmt{
							&OperatorDecl{
								Name: token.Opf,
								Decs: []*Decorator{
									{
										Name: &Ident{Name: "override"},
									},
								},
								Params: []Expr{
									&Parameter{Name: &Ident{Name: "s"}, Type: &TypeReference{Name: &Ident{Name: "str"}}},
								},
								Type: &TypeReference{
									Name: &Ident{Name: "File"},
								},
								Body: &BlockStmt{
									List: []Stmt{
										&AssignStmt{
											Scope: token.LET,
											Lhs:   &Ident{Name: "f"},
											Rhs: &CallExpr{
												Fun:  &Ident{Name: "File"},
												Recv: []Expr{&RefExpr{X: &Ident{Name: "s"}}},
											},
										},
										&CmdStmt{
											Name: &Ident{Name: "echo"},
											Recv: []Expr{&StringLiteral{Value: "creating ${f.name} object"}},
										},
										&ReturnStmt{
											X: &RefExpr{X: &Ident{Name: "f"}},
										},
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

func TestMap(t *testing.T) {
	Print(&File{
		Stmts: []Stmt{
			&AssignStmt{
				Scope: token.LET,
				Lhs:   &Ident{Name: "m"},
				Type: &TypeReference{
					Name: &Ident{Name: "map"},
					TypeParams: []Expr{
						&TypeReference{Name: &Ident{Name: "str"}},
						&TypeReference{Name: &Ident{Name: "any"}},
					},
				},
				Rhs: &ObjectLiteralExpr{
					Props: []Expr{
						&KeyValueExpr{Key: &StringLiteral{Value: "a"}, Value: &NumericLiteral{Value: "10"}},
						&KeyValueExpr{Key: &StringLiteral{Value: "b"}, Value: &ObjectLiteralExpr{}},
						&KeyValueExpr{Key: &StringLiteral{Value: "c"}, Value: &ArrayLiteralExpr{}},
					},
				},
			},
			&AssignStmt{
				Scope: token.LET,
				Lhs:   &Ident{Name: "m"},
				Type:  &AnyLiteral{},
				Rhs: &NamedObjectLiteralExpr{
					Name: &Ident{Name: "User"},
					Props: []Expr{
						&KeyValueExpr{Key: &Ident{Name: "name"}, Value: &StringLiteral{Value: "ansurfen"}},
					},
				},
			},
		},
	})
}

func TestGenericDecl(t *testing.T) {
	Print(&File{
		Stmts: []Stmt{
			&ClassDecl{
				Name: &Ident{Name: "FileSystem"},
				TypeParams: []Expr{
					&TypeParameter{
						Name: &Ident{Name: "T"},
						Constraints: []Expr{
							&TypeReference{Name: &Ident{Name: "Readable"}},
							&TypeReference{Name: &Ident{Name: "Writable"}},
						},
					},
					&TypeParameter{Name: &Ident{Name: "U"}},
				},
				Methods: []Stmt{
					&FuncDecl{
						Name: &Ident{Name: "to_str"},
						TypeParams: []Expr{
							&TypeParameter{Name: &Ident{Name: "T"}},
							&TypeParameter{Name: &Ident{Name: "U"}},
						},
						Recv: []Expr{
							&Parameter{
								Name: &Ident{Name: "a"},
								Type: &TypeReference{Name: &Ident{Name: "T"}},
							},
							&Parameter{
								Name:  &Ident{Name: "b"},
								Type:  &UnionType{Types: []Expr{&TypeReference{Name: &Ident{Name: "str"}}, &TypeReference{Name: &Ident{Name: "num"}}}},
								Value: &StringLiteral{Value: "default"},
							},
							&Parameter{
								Variadic: true,
								Name:     &Ident{Name: "args"},
								Type:     &AnyLiteral{},
							},
							&NamedParameters{
								Params: []Expr{
									&Parameter{
										Required: true,
										Name:     &Ident{Name: "ok"},
										Type:     &TypeReference{Name: &Ident{Name: "bool"}},
									},
									&Parameter{
										Name:  &Ident{Name: "name"},
										Type:  &NullableType{X: &TypeReference{Name: &Ident{Name: "str"}}},
										Value: &StringLiteral{Value: "user1"},
									},
								},
							},
						},
						Type:  &TypeReference{Name: &Ident{Name: "U"}},
						Throw: true,
					},
				},
			},
		},
	})
}

func TestUseDecl(t *testing.T) {
	Print(&File{
		Stmts: []Stmt{
			&UseDecl{
				Lhs: &Ident{Name: "grep"},
				Rhs: &TypeReference{
					Name: &Ident{Name: "If"},
					TypeParams: []Expr{
						&BinaryExpr{
							X:  &RefExpr{X: &Ident{Name: "os"}},
							Op: token.EQ,
							Y:  &StringLiteral{Value: "windows"},
						},
						&IntersectionType{
							Types: []Expr{&TypeReference{Name: &Ident{Name: "find"}}, &TypeReference{Name: &Ident{Name: "findstr"}}},
						},
						&TypeReference{
							Name: &Ident{Name: "If"},
							TypeParams: []Expr{
								&BinaryExpr{
									X:  &RefExpr{X: &Ident{Name: "os"}},
									Op: token.EQ,
									Y:  &StringLiteral{Value: "posix"},
								},
								&TypeReference{Name: &Ident{Name: "grep"}},
								&NullLiteral{},
							},
						},
					},
				},
			},
		},
	})
}
