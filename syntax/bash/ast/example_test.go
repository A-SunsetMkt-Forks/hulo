// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package ast

import (
	"fmt"
	"os"
	"strings"

	"github.com/hulo-lang/hulo/syntax/bash/token"
)

// ExampleInspect_basic demonstrates basic usage of the Inspect function
// with a simple AST node.
func ExampleInspect_basic() {
	// Create a simple word node
	node := &Word{Val: "hello world"}

	// Inspect the node and print to stdout
	Inspect(node, os.Stdout)

	// Output:
	// *ast.Word (Val: "hello world")
}

// ExampleInspect_file demonstrates inspecting a complete file AST
// with multiple statements.
func ExampleInspect_file() {
	// Create a file with multiple statements
	file := &File{
		Stmts: []Stmt{
			&AssignStmt{
				Lhs: &Ident{Name: "greeting"},
				Rhs: &Word{Val: "Hello, World!"},
			},
			&ExprStmt{
				X: &CmdExpr{
					Name: &Ident{Name: "echo"},
					Recv: []Expr{
						&Word{Val: "Hello"},
						&Word{Val: "World"},
					},
				},
			},
		},
	}

	// Inspect the file AST
	Inspect(file, os.Stdout)

	// Output:
	// *ast.File {
	//   Stmts: [
	//     0: *ast.AssignStmt {
	//       Lhs: *ast.Ident (Name: "greeting")
	//       Rhs: *ast.Word (Val: "Hello, World!")
	//       Local: false
	//     }
	//     1: *ast.ExprStmt {
	//       X: *ast.CmdExpr {
	//         Name: *ast.Ident (Name: "echo")
	//         Recv: [
	//           0: *ast.Word (Val: "Hello")
	//           1: *ast.Word (Val: "World")
	//         ]
	//       }
	//     }
	//   ]
	// }
}

// ExampleInspect_function demonstrates inspecting a function declaration
// with a complex body.
func ExampleInspect_function() {
	// Create a function declaration
	funcDecl := &FuncDecl{
		Name: &Ident{Name: "greet"},
		Body: &BlockStmt{
			List: []Stmt{
				&IfStmt{
					Cond: &BinaryExpr{
						X:  &Ident{Name: "name"},
						Op: token.Equal,
						Y:  &Word{Val: ""},
					},
					Body: &BlockStmt{
						List: []Stmt{
							&ExprStmt{
								X: &CmdExpr{
									Name: &Ident{Name: "echo"},
									Recv: []Expr{
										&Word{Val: "Hello, Anonymous!"},
									},
								},
							},
						},
					},
				},
				&ExprStmt{
					X: &CmdExpr{
						Name: &Ident{Name: "echo"},
						Recv: []Expr{
							&Word{Val: "Hello,"},
							&Ident{Name: "name"},
						},
					},
				},
			},
		},
	}

	// Inspect the function declaration
	Inspect(funcDecl, os.Stdout)

	// Output:
	// *ast.FuncDecl {
	//   Name: *ast.Ident (Name: "greet")
	//   Body: *ast.BlockStmt {
	//     List: [
	//       0: *ast.IfStmt {
	//         Cond: *ast.BinaryExpr {
	//           X: *ast.Ident (Name: "name")
	//           Op: ==
	//           Y: *ast.Word (Val: "")
	//         }
	//         Body: *ast.BlockStmt {
	//           List: [
	//             0: *ast.ExprStmt {
	//               X: *ast.CmdExpr {
	//                 Name: *ast.Ident (Name: "echo")
	//                 Recv: [
	//                   0: *ast.Word (Val: "Hello, Anonymous!")
	//                 ]
	//               }
	//             }
	//           ]
	//         }
	//       }
	//       1: *ast.ExprStmt {
	//         X: *ast.CmdExpr {
	//           Name: *ast.Ident (Name: "echo")
	//           Recv: [
	//             0: *ast.Word (Val: "Hello,")
	//             1: *ast.Ident (Name: "name")
	//           ]
	//         }
	//       }
	//     ]
	//   }
	// }
}

// ExampleInspect_loop demonstrates inspecting loop statements
// like while, for, and for-in loops.
func ExampleInspect_loop() {
	// Create a while loop
	whileLoop := &WhileStmt{
		Cond: &BinaryExpr{
			X:  &Ident{Name: "counter"},
			Op: token.TsLss,
			Y:  &Word{Val: "10"},
		},
		Body: &BlockStmt{
			List: []Stmt{
				&ExprStmt{
					X: &CmdExpr{
						Name: &Ident{Name: "echo"},
						Recv: []Expr{
							&Ident{Name: "counter"},
						},
					},
				},
				&AssignStmt{
					Lhs: &Ident{Name: "counter"},
					Rhs: &BinaryExpr{
						X:  &Ident{Name: "counter"},
						Op: token.Plus,
						Y:  &Word{Val: "1"},
					},
				},
			},
		},
	}

	// Inspect the while loop
	Inspect(whileLoop, os.Stdout)

	// Output:
	// *ast.WhileStmt {
	//   Cond: *ast.BinaryExpr {
	//     X: *ast.Ident (Name: "counter")
	//     Op: -lt
	//     Y: *ast.Word (Val: "10")
	//   }
	//   Body: *ast.BlockStmt {
	//     List: [
	//       0: *ast.ExprStmt {
	//         X: *ast.CmdExpr {
	//           Name: *ast.Ident (Name: "echo")
	//           Recv: [
	//             0: *ast.Ident (Name: "counter")
	//           ]
	//         }
	//       }
	//       1: *ast.AssignStmt {
	//         Lhs: *ast.Ident (Name: "counter")
	//         Rhs: *ast.BinaryExpr {
	//           X: *ast.Ident (Name: "counter")
	//           Op: +
	//           Y: *ast.Word (Val: "1")
	//         }
	//         Local: false
	//       }
	//     ]
	//   }
	// }
}

// ExampleInspect_expression demonstrates inspecting various expressions
// including binary expressions, command expressions, and variable expansions.
func ExampleInspect_expression() {
	// Create a complex expression
	expr := &BinaryExpr{
		X: &CmdExpr{
			Name: &Ident{Name: "echo"},
			Recv: []Expr{
				&VarExpExpr{
					X: &Ident{Name: "USER"},
				},
			},
		},
		Op: token.AndAnd,
		Y: &CmdExpr{
			Name: &Ident{Name: "echo"},
			Recv: []Expr{
				&Word{Val: "Command succeeded"},
			},
		},
	}

	// Inspect the expression
	Inspect(expr, os.Stdout)

	// Output:
	// *ast.BinaryExpr {
	//   X: *ast.CmdExpr {
	//     Name: *ast.Ident (Name: "echo")
	//     Recv: [
	//       0: *ast.VarExpExpr {
	//         X: *ast.Ident (Name: "USER")
	//       }
	//     ]
	//   }
	//   Op: &&
	//   Y: *ast.CmdExpr {
	//     Name: *ast.Ident (Name: "echo")
	//     Recv: [
	//       0: *ast.Word (Val: "Command succeeded")
	//     ]
	//   }
	// }
}

// ExampleInspect_stringBuilder demonstrates how to capture the output
// of Inspect in a string for further processing.
func ExampleInspect_stringBuilder() {
	// Create a simple AST node
	node := &AssignStmt{
		Lhs: &Ident{Name: "variable"},
		Rhs: &Word{Val: "value"},
	}

	// Use strings.Builder to capture the output
	var buf strings.Builder
	Inspect(node, &buf)

	// Get the string representation
	result := buf.String()
	fmt.Println("Inspect output:")
	fmt.Println(result)

	// Output:
	// Inspect output:
	// *ast.AssignStmt {
	//   Lhs: *ast.Ident (Name: "variable")
	//   Rhs: *ast.Word (Val: "value")
	//   Local: false
	// }
}

// ExampleInspect_nil demonstrates how Inspect handles nil nodes.
func ExampleInspect_nil() {
	// Inspect a nil node
	Inspect(nil, os.Stdout)

	// Output:
	// nil
}

// ExampleInspect_complex demonstrates a more complex AST structure
// with nested statements and expressions.
func ExampleInspect_complex() {
	// Create a complex AST representing a bash script
	script := &File{
		Stmts: []Stmt{
			&CommentGroup{
				List: []*Comment{
					{Text: " This is a sample bash script"},
				},
			},
			&FuncDecl{
				Name: &Ident{Name: "main"},
				Body: &BlockStmt{
					List: []Stmt{
						&AssignStmt{
							Lhs: &Ident{Name: "count"},
							Rhs: &Word{Val: "0"},
						},
						&ForInStmt{
							Var:  &Ident{Name: "item"},
							List: &Word{Val: "apple banana cherry"},
							Body: &BlockStmt{
								List: []Stmt{
									&ExprStmt{
										X: &CmdExpr{
											Name: &Ident{Name: "echo"},
											Recv: []Expr{
												&Word{Val: "Processing:"},
												&Ident{Name: "item"},
											},
										},
									},
									&AssignStmt{
										Lhs: &Ident{Name: "count"},
										Rhs: &BinaryExpr{
											X:  &Ident{Name: "count"},
											Op: token.Plus,
											Y:  &Word{Val: "1"},
										},
									},
								},
							},
						},
						&ExprStmt{
							X: &CmdExpr{
								Name: &Ident{Name: "echo"},
								Recv: []Expr{
									&Word{Val: "Total items:"},
									&Ident{Name: "count"},
								},
							},
						},
					},
				},
			},
		},
	}

	// Inspect the complex script
	Inspect(script, os.Stdout)

	// Output:
	// *ast.File {
	//   Stmts: [
	//     0: *ast.CommentGroup {
	//       List: [
	//         0: *ast.Comment (Text: " This is a sample bash script")
	//       ]
	//     }
	//     1: *ast.FuncDecl {
	//       Name: *ast.Ident (Name: "main")
	//       Body: *ast.BlockStmt {
	//         List: [
	//           0: *ast.AssignStmt {
	//             Lhs: *ast.Ident (Name: "count")
	//             Rhs: *ast.Word (Val: "0")
	//             Local: false
	//           }
	//           1: *ast.ForInStmt {
	//             Var: *ast.Ident (Name: "item")
	//             List: *ast.Word (Val: "apple banana cherry")
	//             Body: *ast.BlockStmt {
	//               List: [
	//                 0: *ast.ExprStmt {
	//                   X: *ast.CmdExpr {
	//                     Name: *ast.Ident (Name: "echo")
	//                     Recv: [
	//                       0: *ast.Word (Val: "Processing:")
	//                       1: *ast.Ident (Name: "item")
	//                     ]
	//                   }
	//                 }
	//                 1: *ast.AssignStmt {
	//                   Lhs: *ast.Ident (Name: "count")
	//                   Rhs: *ast.BinaryExpr {
	//                     X: *ast.Ident (Name: "count")
	//                     Op: +
	//                     Y: *ast.Word (Val: "1")
	//                   }
	//                   Local: false
	//                 }
	//               ]
	//             }
	//           }
	//           2: *ast.ExprStmt {
	//             X: *ast.CmdExpr {
	//               Name: *ast.Ident (Name: "echo")
	//               Recv: [
	//                 0: *ast.Word (Val: "Total items:")
	//                 1: *ast.Ident (Name: "count")
	//               ]
	//             }
	//           }
	//         ]
	//       }
	//     }
	//   ]
	// }
}
