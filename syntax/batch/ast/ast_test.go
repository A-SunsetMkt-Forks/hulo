// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package ast

import (
	"testing"

	"github.com/hulo-lang/hulo/syntax/batch/token"
)

func TestAssign(t *testing.T) {
	Print(&File{
		Stmts: []Stmt{
			&ExprStmt{
				X: Words(Literal("@echo"), Literal("off")),
			},
			&ExprStmt{
				X: CmdExpression("SETLOCAL"),
			},
			&CallStmt{
				Name: "Display",
				Recv: []Expr{
					Words(Literal("5"), Literal(","), Literal("10")),
				},
			},
			&ExprStmt{
				X: &CallExpr{
					Fun: Literal("EXIT"),
					Recv: []Expr{
						&UnaryExpr{
							Op: token.DIV,
							X: &Lit{
								Val: "B",
							},
						},
						&DblQuote{
							Val: &Lit{
								Val: "ERRORLEVEL",
							},
						},
					},
				},
			},
			&FuncDecl{
				Name: "Display",
				Body: &BlockStmt{
					List: []Stmt{
						&ExprStmt{
							X: &CmdExpr{
								Name: &Lit{
									Val: "echo",
								},
								Recv: []Expr{
									&Lit{Val: "The value of parameter 1 is"},
									&SglQuote{
										Val: &Lit{
											Val: "~1",
										},
									},
								},
							},
						},
						&ExprStmt{
							X: &CmdExpr{
								Name: &Lit{
									Val: "echo",
								},
								Recv: []Expr{
									&Lit{Val: "The value of parameter 1 is"},
									&SglQuote{
										Val: &Lit{
											Val: "~2",
										},
									},
								},
							},
						},
						&ExprStmt{
							X: &CmdExpr{
								Name: Literal("EXIT"),
								Recv: []Expr{
									&UnaryExpr{
										Op: token.DIV,
										X: &Lit{
											Val: "B",
										},
									},
									Literal("0"),
								},
							},
						},
					},
				},
			},
		},
	})
}

func TestBatchFeatures(t *testing.T) {
	tests := []struct {
		name string
		file *File
	}{
		{
			name: "Basic IF conditions",
			file: &File{
				Stmts: []Stmt{
					&IfStmt{
						Cond: &BinaryExpr{
							X:  &DblQuote{Val: &Lit{Val: "var"}},
							Op: token.DOUBLE_ASSIGN,
							Y:  &Lit{Val: "value"},
						},
						Body: &BlockStmt{
							List: []Stmt{
								&ExprStmt{
									X: &CmdExpr{
										Name: &Lit{Val: "ECHO"},
										Recv: []Expr{&Lit{Val: "Match"}},
									},
								},
							},
						},
						Else: &BlockStmt{
							List: []Stmt{
								&ExprStmt{
									X: &CmdExpr{
										Name: &Lit{Val: "ECHO"},
										Recv: []Expr{&Lit{Val: "No Match"}},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "FOR loop with file iteration",
			file: &File{
				Stmts: []Stmt{
					&ForStmt{
						X:    &SglQuote{Val: &SglQuote{Val: &Lit{Val: "f"}}},
						List: &Lit{Val: "(*.txt)"},
						Body: &BlockStmt{
							List: []Stmt{
								&ExprStmt{
									X: &CmdExpr{
										Name: &Lit{Val: "echo"},
										Recv: []Expr{&SglQuote{Val: &SglQuote{Val: &Lit{Val: "f"}}}},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "GOTO and labels",
			file: &File{
				Stmts: []Stmt{
					&IfStmt{
						Cond: &BinaryExpr{
							X:  &DblQuote{Val: &Lit{Val: "1"}},
							Op: token.DOUBLE_ASSIGN,
							Y:  &Lit{Val: "start"},
						},
						Body: &BlockStmt{
							List: []Stmt{
								&GotoStmt{Label: "start"},
							},
						},
					},
					&LabelStmt{Name: "start"},
					&ExprStmt{
						X: &CmdExpr{
							Name: &Lit{Val: "ECHO"},
							Recv: []Expr{&Lit{Val: "Started"}},
						},
					},
					&GotoStmt{Label: "EOF"},
				},
			},
		},
		{
			name: "Function call and definition",
			file: &File{
				Stmts: []Stmt{
					&CallStmt{
						Name: "myFunction",
						Recv: []Expr{
							&Lit{Val: "arg1"},
							&Lit{Val: "arg2"},
						},
					},
					&GotoStmt{Label: "EOF"},
					&FuncDecl{
						Name: "myFunction",
						Body: &BlockStmt{
							List: []Stmt{
								&ExprStmt{
									X: &CmdExpr{
										Name: &Lit{Val: "ECHO"},
										Recv: []Expr{&Lit{Val: "Argument1: "}, &SglQuote{Val: &Lit{Val: "1"}}},
									},
								},
								&GotoStmt{Label: "EOF"},
							},
						},
					},
				},
			},
		},
		{
			name: "EXIT commands",
			file: &File{
				Stmts: []Stmt{
					&ExprStmt{
						X: &CmdExpr{
							Name: &Lit{Val: "EXIT"},
							Recv: []Expr{&Lit{Val: "1"}},
						},
					},
					&ExprStmt{
						X: &CmdExpr{
							Name: &Lit{Val: "EXIT"},
							Recv: []Expr{
								&UnaryExpr{
									Op: token.DIV,
									X:  &Lit{Val: "B"},
								},
								&Lit{Val: "0"},
							},
						},
					},
				},
			},
		},
		{
			name: "Nested IF and FOR",
			file: &File{
				Stmts: []Stmt{
					&IfStmt{
						Cond: &BinaryExpr{
							X:  &DblQuote{Val: &Lit{Val: "a"}},
							Op: token.EQU,
							Y:  &Lit{Val: "1"},
						},
						Body: &BlockStmt{
							List: []Stmt{
								&ForStmt{
									X:    &SglQuote{Val: &SglQuote{Val: &Lit{Val: "i"}}},
									List: &Lit{Val: "(1 2 3)"},
									Body: &BlockStmt{
										List: []Stmt{
											&ExprStmt{
												X: &CmdExpr{
													Name: &Lit{Val: "ECHO"},
													Recv: []Expr{
														&Lit{Val: "Inside Loop "},
														&SglQuote{Val: &Lit{Val: "i"}},
													},
												},
											},
										},
									},
								},
							},
						},
						Else: &BlockStmt{
							List: []Stmt{
								&ExprStmt{
									X: &CmdExpr{
										Name: &Lit{Val: "ECHO"},
										Recv: []Expr{&Lit{Val: "Not Matched"}},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Delayed expansion",
			file: &File{
				Stmts: []Stmt{
					&ExprStmt{
						X: &CmdExpr{
							Name: &Lit{Val: "setlocal"},
							Recv: []Expr{&Lit{Val: "enabledelayedexpansion"}},
						},
					},
					&AssignStmt{
						Lhs: &Lit{Val: "VAR"},
						Rhs: &Lit{
							Val: "10",
						},
					},
					&ForStmt{
						X:    &SglQuote{Val: &SglQuote{Val: &Lit{Val: "i"}}},
						List: &Lit{Val: "(1 2 3)"},
						Body: &BlockStmt{
							List: []Stmt{
								&AssignStmt{
									Lhs: &Lit{Val: "VAR"},
									Rhs: &SglQuote{Val: &SglQuote{Val: &Lit{Val: "i"}}},
								},
								&ExprStmt{
									X: &CmdExpr{
										Name: &Lit{Val: "echo"},
										Recv: []Expr{&DblQuote{
											DelayedExpansion: true,
											Val:              &Lit{Val: "VAR"},
										}},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Print(tt.file)
		})
	}
}

func TestIfStmt(t *testing.T) {
	Print(&IfStmt{
		Cond: &BinaryExpr{
			X:  &DblQuote{Val: &Lit{Val: "1"}},
			Op: token.DOUBLE_ASSIGN,
			Y:  &Lit{Val: "start"},
		},
		Body: &IfStmt{
			Body: &BlockStmt{
				List: []Stmt{
					&ExprStmt{
						X: &CmdExpr{
							Name: &Lit{Val: "echo"},
						},
					},
				},
			},
			Else: &BlockStmt{
				List: []Stmt{
					&ExprStmt{
						X: &CmdExpr{
							Name: &Lit{Val: "echo"},
						},
					},
				},
			},
		},
	})
}
