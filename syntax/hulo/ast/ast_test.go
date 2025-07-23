// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package ast

import (
	"testing"

	"github.com/hulo-lang/hulo/syntax/hulo/token"
	"github.com/stretchr/testify/assert"
)

func TestAssignStmt(t *testing.T) {
	Print(&AssignStmt{
		Scope: token.LET,
		Lhs:   &Ident{Name: "x"},
		Rhs:   &NumericLiteral{Value: "123"},
	})

	Print(&AssignStmt{
		Lhs: &Ident{Name: "x"},
		Tok: token.COLON_ASSIGN,
		Rhs: &NumericLiteral{Value: "123"},
	})
}

func TestIfStmt(t *testing.T) {
	Print(&IfStmt{
		Cond: &BinaryExpr{
			X:  &RefExpr{X: &Ident{Name: "x"}},
			Op: token.LT,
			Y:  &NumericLiteral{Value: "10"},
		},
		Body: &BlockStmt{
			List: []Stmt{
				&AssignStmt{
					Scope: token.LET,
					Lhs:   &Ident{Name: "x"},
					Rhs:   &NumericLiteral{Value: "10"},
				},
			},
		},
		Else: &IfStmt{
			Cond: &BinaryExpr{
				X:  &RefExpr{X: &Ident{Name: "x"}},
				Op: token.LT,
				Y:  &NumericLiteral{Value: "20"},
			},
			Body: &BlockStmt{
				List: []Stmt{
					&AssignStmt{
						Scope: token.LET,
						Lhs:   &Ident{Name: "x"},
						Rhs:   &NumericLiteral{Value: "20"},
					},
				},
			},
			Else: &BlockStmt{
				List: []Stmt{
					&AssignStmt{
						Scope: token.LET,
						Lhs:   &Ident{Name: "x"},
						Rhs:   &NumericLiteral{Value: "20"},
					},
				},
			},
		},
	})
}

func TestTypeDecl(t *testing.T) {
	Print(&File{
		Stmts: []Stmt{
			&TypeDecl{
				Name: &Ident{Name: "User"},
				Value: &TypeLiteral{
					Members: []Expr{
						&KeyValueExpr{
							Key:   &Ident{Name: "name"},
							Value: &TypeReference{Name: &Ident{Name: "str"}},
						},
						&KeyValueExpr{
							Key:   &Ident{Name: "age"},
							Value: &TypeReference{Name: &Ident{Name: "num"}},
						},
					},
				},
			},
		},
	})
}

func TestDeclareDecl(t *testing.T) {
	Print(&File{
		Stmts: []Stmt{
			&DeclareDecl{
				X: &ComptimeStmt{
					Body: &AssignStmt{
						Scope: token.CONST,
						Lhs:   &Ident{Name: "os"},
						Rhs:   &StringLiteral{Value: "windows"},
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
										&TrueLiteral{},
									},
								},
								Rhs: &NumericLiteral{Value: "10"},
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

func TestClassDecl(t *testing.T) {
	Print(&File{
		Stmts: []Stmt{
			&ClassDecl{
				Name: &Ident{Name: "Constants"},
				Ctors: []*ConstructorDecl{
					{
						Modifiers: []Modifier{&ConstModifier{Const: token.Pos(10)}},
						ClsName:   &Ident{Name: "Constants"},
						Recv: []Expr{
							&SelectExpr{
								X: &RefExpr{X: &Ident{Name: "this"}},
								Y: &Ident{Name: "PI"},
							},
							&SelectExpr{
								X: &RefExpr{X: &Ident{Name: "this"}},
								Y: &Ident{Name: "E"},
							},
						},
					},
					{
						Modifiers: []Modifier{&ConstModifier{Const: token.Pos(10)}},
						ClsName:   &Ident{Name: "Constants"},
						Recv: []Expr{
							&Parameter{Name: &Ident{Name: "pi"}, Type: &TypeReference{Name: &Ident{Name: "num"}}},
							&Parameter{Name: &Ident{Name: "e"}, Type: &TypeReference{Name: &Ident{Name: "num"}}},
						},
						InitFields: []Expr{
							&BinaryExpr{
								X:  &SelectExpr{X: &RefExpr{X: &Ident{Name: "this"}}, Y: &Ident{Name: "PI"}},
								Op: token.ASSIGN,
								Y:  &RefExpr{X: &Ident{Name: "pi"}},
							},
							&BinaryExpr{
								X:  &SelectExpr{X: &RefExpr{X: &Ident{Name: "this"}}, Y: &Ident{Name: "E"}},
								Op: token.ASSIGN,
								Y:  &RefExpr{X: &Ident{Name: "e"}},
							},
						},
					},
					{
						ClsName: &Ident{Name: "Constants"},
						Name:    &Ident{Name: "zero"},
						Recv: []Expr{
							&NamedParameters{
								Params: []Expr{
									&Parameter{Name: &Ident{Name: "pi"}, Type: &TypeReference{Name: &Ident{Name: "num"}}},
									&Parameter{Name: &Ident{Name: "e"}, Type: &TypeReference{Name: &Ident{Name: "num"}}},
								},
							},
						},
					},
				},
				Fields: &FieldList{
					List: []*Field{
						{
							Modifiers: []Modifier{&FinalModifier{Final: token.Pos(10)}},
							Name:      &Ident{Name: "PI"},
							Type:      &TypeReference{Name: &Ident{Name: "num"}},
							Value:     &NumericLiteral{Value: "3.14"},
						},
						{
							Modifiers: []Modifier{&FinalModifier{Final: token.Pos(10)}},
							Name:      &Ident{Name: "E"},
							Type:      &TypeReference{Name: &Ident{Name: "num"}},
							Value: &NumericLiteral{
								Value: "2.718281828459045",
							},
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
					Y:  &StringLiteral{Value: "Alice"},
				},
			},
			Y: &BinaryExpr{
				X:  &Ident{Name: "age"},
				Op: token.ASSIGN,
				Y:  &NumericLiteral{Value: "25"},
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

func TestExtensionDecl(t *testing.T) {
	Print(&File{
		Stmts: []Stmt{
			&ExtensionDecl{
				Name: &Ident{Name: "str"},
				Body: &ExtensionClass{
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
				Methods: []*FuncDecl{
					{
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
								Modifier: &EllipsisModifier{},
								Name:     &Ident{Name: "args"},
								Type:     &AnyLiteral{},
							},
							&NamedParameters{
								Params: []Expr{
									&Parameter{
										Modifier: &RequiredModifier{},
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

func TestEnumDecl(t *testing.T) {
	Print(&File{
		Stmts: []Stmt{
			// Basic enum
			&EnumDecl{
				Name: &Ident{Name: "Status"},
				Body: &BasicEnumBody{
					Values: []*EnumValue{
						{Name: &Ident{Name: "Pending"}},
						{Name: &Ident{Name: "Approved"}},
						{Name: &Ident{Name: "Rejected"}},
					},
				},
			},
			// Enum with modifiers
			&EnumDecl{
				Modifiers: []Modifier{
					&PubModifier{},
				},
				Name: &Ident{Name: "HttpCode"},
				Body: &BasicEnumBody{
					Values: []*EnumValue{
						{
							Name:  &Ident{Name: "OK"},
							Value: &NumericLiteral{Value: "200"},
						},
						{
							Name:  &Ident{Name: "NotFound"},
							Value: &NumericLiteral{Value: "404"},
						},
						{
							Name:  &Ident{Name: "ServerError"},
							Value: &NumericLiteral{Value: "500"},
						},
					},
				},
			},
			// Associated value enum
			&EnumDecl{
				Name: &Ident{Name: "Protocol"},
				Body: &AssociatedEnumBody{
					Fields: &FieldList{
						List: []*Field{
							{
								Name: &Ident{Name: "port"},
								Type: &Ident{Name: "num"},
							},
						},
					},
					Values: []*EnumValue{
						{
							Name: &Ident{Name: "tcp"},
							Data: []Expr{&NumericLiteral{Value: "6"}},
						},
						{
							Name: &Ident{Name: "udp"},
							Data: []Expr{&NumericLiteral{Value: "17"}},
						},
					},
				},
			},
			// ADT enum with const modifier
			&EnumDecl{
				Modifiers: []Modifier{
					&ConstModifier{Const: token.Pos(25)},
				},
				Name: &Ident{Name: "NetworkPacket"},
				Body: &ADTEnumBody{
					Variants: []*EnumVariant{
						{
							Name: &Ident{Name: "TCP"},
							Fields: &FieldList{
								List: []*Field{
									{
										Name: &Ident{Name: "src_port"},
										Type: &Ident{Name: "num"},
									},
									{
										Name: &Ident{Name: "dst_port"},
										Type: &Ident{Name: "num"},
									},
								},
							},
						},
						{
							Name: &Ident{Name: "UDP"},
							Fields: &FieldList{
								List: []*Field{
									{
										Name: &Ident{Name: "port"},
										Type: &Ident{Name: "num"},
									},
									{
										Name: &Ident{Name: "payload"},
										Type: &Ident{Name: "str"},
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

func TestModifiers(t *testing.T) {
	tests := []struct {
		name     string
		input    Modifier
		expected ModifierKind
		pos      token.Pos
	}{
		{
			name:     "FinalModifier",
			input:    &FinalModifier{Final: token.Pos(10)},
			expected: ModKindFinal,
			pos:      token.Pos(10),
		},
		{
			name:     "ConstModifier",
			input:    &ConstModifier{Const: token.Pos(20)},
			expected: ModKindConst,
			pos:      token.Pos(20),
		},
		{
			name:     "PubModifier",
			input:    &PubModifier{Pub: token.Pos(30)},
			expected: ModKindPub,
			pos:      token.Pos(30),
		},
		{
			name:     "StaticModifier",
			input:    &StaticModifier{Static: token.Pos(40)},
			expected: ModKindStatic,
			pos:      token.Pos(40),
		},
		{
			name:     "RequiredModifier",
			input:    &RequiredModifier{Required: token.Pos(50)},
			expected: ModKindRequired,
			pos:      token.Pos(50),
		},
		{
			name:     "ReadonlyModifier",
			input:    &ReadonlyModifier{Readonly: token.Pos(60)},
			expected: ModKindReadonly,
			pos:      token.Pos(60),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test modifier kind
			assert.Equal(t, tt.expected, tt.input.Kind())

			// Test position information
			assert.Equal(t, tt.pos, tt.input.Pos())

			// Test type assertions
			switch tt.expected {
			case ModKindFinal:
				assert.IsType(t, &FinalModifier{}, tt.input)
			case ModKindConst:
				assert.IsType(t, &ConstModifier{}, tt.input)
			case ModKindPub:
				assert.IsType(t, &PubModifier{}, tt.input)
			case ModKindStatic:
				assert.IsType(t, &StaticModifier{}, tt.input)
			case ModKindRequired:
				assert.IsType(t, &RequiredModifier{}, tt.input)
			case ModKindReadonly:
				assert.IsType(t, &ReadonlyModifier{}, tt.input)
			}
		})
	}

	// Test modifier array functionality
	t.Run("ModifierArray", func(t *testing.T) {
		modifiers := []Modifier{
			&FinalModifier{Final: token.Pos(10)},
			&ConstModifier{Const: token.Pos(20)},
			&PubModifier{Pub: token.Pos(30)},
		}

		assert.Equal(t, 3, len(modifiers))

		// Verify each modifier in the array
		expectedKinds := []ModifierKind{ModKindFinal, ModKindConst, ModKindPub}
		for i, mod := range modifiers {
			assert.Equal(t, expectedKinds[i], mod.Kind())
		}
	})
}
