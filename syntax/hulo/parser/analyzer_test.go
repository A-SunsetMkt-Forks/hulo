package parser

import (
	"testing"

	"github.com/antlr4-go/antlr/v4"
	"github.com/hulo-lang/hulo/syntax/hulo/ast"
	"github.com/hulo-lang/hulo/syntax/hulo/token"
)

func TestDeclareStatement(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected ast.Node
	}{
		{
			name:  "declare block",
			input: "declare { let x = 10 }",
			expected: &ast.DeclareDecl{
				X: &ast.BlockStmt{
					List: []ast.Stmt{
						&ast.AssignStmt{
							Scope:    token.LET,
							ScopePos: token.Pos(9), // position of "let"
							Lhs: &ast.Ident{
								NamePos: token.Pos(13), // position of "x"
								Name:    "x",
							},
							Tok: token.ASSIGN,
							Rhs: &ast.NumericLiteral{
								Value:    "10",
								ValuePos: token.Pos(17), // position of "10"
							},
						},
					},
				},
			},
		},
		{
			name:  "declare class",
			input: "declare class User { name: str }",
			expected: &ast.DeclareDecl{
				X: &ast.ClassDecl{
					Class: token.Pos(9), // position of "class"
					Name: &ast.Ident{
						NamePos: token.Pos(15), // position of "User"
						Name:    "User",
					},
					Fields: &ast.FieldList{
						List: []*ast.Field{
							{
								Name: &ast.Ident{
									NamePos: token.Pos(22), // position of "name"
									Name:    "name",
								},
								Colon: token.Pos(26), // position of ":"
								Type: &ast.TypeReference{
									Name: &ast.Ident{
										NamePos: token.Pos(28), // position of "str"
										Name:    "str",
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name:  "declare function",
			input: "declare fn hello() {}",
			expected: &ast.DeclareDecl{
				X: &ast.FuncDecl{
					Fn: token.Pos(9), // position of "fn"
					Name: &ast.Ident{
						NamePos: token.Pos(12), // position of "hello"
						Name:    "hello",
					},
					Body: &ast.BlockStmt{
						Lbrace: token.Pos(18), // position of "{"
						Rbrace: token.Pos(19), // position of "}"
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create input stream
			input := antlr.NewInputStream(tt.input)

			// Create analyzer
			analyzer, err := NewAnalyzer(input)
			if err != nil {
				t.Fatalf("Failed to create analyzer: %v", err)
			}

			// Parse the input
			result := analyzer.file.Accept(analyzer)
			if result == nil {
				t.Fatalf("Failed to parse input")
			}

			// Get the first statement (should be our declare statement)
			file, ok := result.(*ast.File)
			if !ok {
				t.Fatalf("Expected *ast.File, got %T", result)
			}

			if len(file.Stmts) == 0 {
				t.Fatalf("No statements found in file")
			}

			declareStmt, ok := file.Stmts[0].(*ast.DeclareDecl)
			if !ok {
				t.Fatalf("Expected *ast.DeclareDecl, got %T", file.Stmts[0])
			}

			// Basic validation
			if declareStmt.Declare == token.NoPos {
				t.Error("Declare position should not be NoPos")
			}

			if declareStmt.X == nil {
				t.Error("Declare content should not be nil")
			}

			// Check the type of the declared content
			switch tt.expected.(*ast.DeclareDecl).X.(type) {
			case *ast.BlockStmt:
				if _, ok := declareStmt.X.(*ast.BlockStmt); !ok {
					t.Errorf("Expected *ast.BlockStmt, got %T", declareStmt.X)
				}
			case *ast.ClassDecl:
				if _, ok := declareStmt.X.(*ast.ClassDecl); !ok {
					t.Errorf("Expected *ast.ClassDecl, got %T", declareStmt.X)
				}
			case *ast.FuncDecl:
				if _, ok := declareStmt.X.(*ast.FuncDecl); !ok {
					t.Errorf("Expected *ast.FuncDecl, got %T", declareStmt.X)
				}
			}
		})
	}
}
