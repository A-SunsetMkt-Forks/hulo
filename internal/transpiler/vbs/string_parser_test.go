// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package transpiler_test

import (
	"testing"

	build "github.com/hulo-lang/hulo/internal/transpiler/vbs"
	vast "github.com/hulo-lang/hulo/syntax/vbs/ast"
	"github.com/stretchr/testify/assert"
)

func TestParseStringInterpolation(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []build.StringPart
	}{
		{
			name:  "simple variable",
			input: "$name",
			expected: []build.StringPart{
				{
					Text:       "",
					IsVariable: true,
					Expr:       &vast.Ident{Name: "name"},
				},
			},
		},
		{
			name:  "braced variable",
			input: "${name}",
			expected: []build.StringPart{
				{
					Text:       "",
					IsVariable: true,
					Expr:       &vast.Ident{Name: "name"},
				},
			},
		},
		{
			name:  "text only",
			input: "Hello, World!",
			expected: []build.StringPart{
				{
					Text:       "Hello, World!",
					IsVariable: false,
				},
			},
		},
		{
			name:  "text with simple variable",
			input: "Hello, $name!",
			expected: []build.StringPart{
				{
					Text:       "Hello, ",
					IsVariable: false,
				},
				{
					Text:       "",
					IsVariable: true,
					Expr:       &vast.Ident{Name: "name"},
				},
				{
					Text:       "!",
					IsVariable: false,
				},
			},
		},
		{
			name:  "text with braced variable",
			input: "Hello, ${name}!",
			expected: []build.StringPart{
				{
					Text:       "Hello, ",
					IsVariable: false,
				},
				{
					Text:       "",
					IsVariable: true,
					Expr:       &vast.Ident{Name: "name"},
				},
				{
					Text:       "!",
					IsVariable: false,
				},
			},
		},
		{
			name:  "multiple simple variables",
			input: "$firstName $lastName",
			expected: []build.StringPart{
				{
					Text:       "",
					IsVariable: true,
					Expr:       &vast.Ident{Name: "firstName"},
				},
				{
					Text:       " ",
					IsVariable: false,
				},
				{
					Text:       "",
					IsVariable: true,
					Expr:       &vast.Ident{Name: "lastName"},
				},
			},
		},
		{
			name:  "multiple braced variables",
			input: "${firstName} ${lastName}",
			expected: []build.StringPart{
				{
					Text:       "",
					IsVariable: true,
					Expr:       &vast.Ident{Name: "firstName"},
				},
				{
					Text:       " ",
					IsVariable: false,
				},
				{
					Text:       "",
					IsVariable: true,
					Expr:       &vast.Ident{Name: "lastName"},
				},
			},
		},
		{
			name:  "mixed simple and braced variables",
			input: "$firstName ${lastName}",
			expected: []build.StringPart{
				{
					Text:       "",
					IsVariable: true,
					Expr:       &vast.Ident{Name: "firstName"},
				},
				{
					Text:       " ",
					IsVariable: false,
				},
				{
					Text:       "",
					IsVariable: true,
					Expr:       &vast.Ident{Name: "lastName"},
				},
			},
		},
		{
			name:  "variable with underscore",
			input: "$user_name",
			expected: []build.StringPart{
				{
					Text:       "",
					IsVariable: true,
					Expr:       &vast.Ident{Name: "user_name"},
				},
			},
		},
		{
			name:  "variable with numbers",
			input: "$user123",
			expected: []build.StringPart{
				{
					Text:       "",
					IsVariable: true,
					Expr:       &vast.Ident{Name: "user123"},
				},
			},
		},
		{
			name:  "variable with mixed case",
			input: "$UserName",
			expected: []build.StringPart{
				{
					Text:       "",
					IsVariable: true,
					Expr:       &vast.Ident{Name: "UserName"},
				},
			},
		},
		{
			name:  "braced variable with underscore",
			input: "${user_name}",
			expected: []build.StringPart{
				{
					Text:       "",
					IsVariable: true,
					Expr:       &vast.Ident{Name: "user_name"},
				},
			},
		},
		{
			name:  "braced variable with numbers",
			input: "${user123}",
			expected: []build.StringPart{
				{
					Text:       "",
					IsVariable: true,
					Expr:       &vast.Ident{Name: "user123"},
				},
			},
		},
		{
			name:  "braced variable with mixed case",
			input: "${UserName}",
			expected: []build.StringPart{
				{
					Text:       "",
					IsVariable: true,
					Expr:       &vast.Ident{Name: "UserName"},
				},
			},
		},
		{
			name:  "complex interpolation",
			input: "Welcome, $firstName ${lastName}! Your ID is $user123.",
			expected: []build.StringPart{
				{
					Text:       "Welcome, ",
					IsVariable: false,
				},
				{
					Text:       "",
					IsVariable: true,
					Expr:       &vast.Ident{Name: "firstName"},
				},
				{
					Text:       " ",
					IsVariable: false,
				},
				{
					Text:       "",
					IsVariable: true,
					Expr:       &vast.Ident{Name: "lastName"},
				},
				{
					Text:       "! Your ID is ",
					IsVariable: false,
				},
				{
					Text:       "",
					IsVariable: true,
					Expr:       &vast.Ident{Name: "user123"},
				},
				{
					Text:       ".",
					IsVariable: false,
				},
			},
		},
		{
			name:  "dollar sign in text",
			input: "Price: $10.99",
			expected: []build.StringPart{
				{
					Text:       "Price: ",
					IsVariable: false,
				},
				{
					Text:       "",
					IsVariable: true,
					Expr:       &vast.Ident{Name: "10"},
				},
				{
					Text:       ".99",
					IsVariable: false,
				},
			},
		},
		{
			name:  "dollar sign followed by non-identifier",
			input: "Price: $!invalid",
			expected: []build.StringPart{
				{
					Text:       "Price: ",
					IsVariable: false,
				},
				{
					Text:       "",
					IsVariable: true,
					Expr:       &vast.Ident{Name: ""},
				},
				{
					Text:       "!invalid",
					IsVariable: false,
				},
			},
		},
		{
			name:  "incomplete braced variable",
			input: "Hello ${name",
			expected: []build.StringPart{
				{
					Text:       "Hello ",
					IsVariable: false,
				},
				{
					Text:       "",
					IsVariable: true,
					Expr:       &vast.Ident{Name: "name"},
				},
			},
		},
		{
			name:  "empty braced variable",
			input: "Hello ${}",
			expected: []build.StringPart{
				{
					Text:       "Hello ",
					IsVariable: false,
				},
				{
					Text:       "",
					IsVariable: true,
					Expr:       &vast.Ident{Name: ""},
				},
			},
		},
		{
			name:  "variable at start",
			input: "$name is here",
			expected: []build.StringPart{
				{
					Text:       "",
					IsVariable: true,
					Expr:       &vast.Ident{Name: "name"},
				},
				{
					Text:       " is here",
					IsVariable: false,
				},
			},
		},
		{
			name:  "variable at end",
			input: "Hello $name",
			expected: []build.StringPart{
				{
					Text:       "Hello ",
					IsVariable: false,
				},
				{
					Text:       "",
					IsVariable: true,
					Expr:       &vast.Ident{Name: "name"},
				},
			},
		},
		{
			name:  "braced variable at start",
			input: "${name} is here",
			expected: []build.StringPart{
				{
					Text:       "",
					IsVariable: true,
					Expr:       &vast.Ident{Name: "name"},
				},
				{
					Text:       " is here",
					IsVariable: false,
				},
			},
		},
		{
			name:  "braced variable at end",
			input: "Hello ${name}",
			expected: []build.StringPart{
				{
					Text:       "Hello ",
					IsVariable: false,
				},
				{
					Text:       "",
					IsVariable: true,
					Expr:       &vast.Ident{Name: "name"},
				},
			},
		},
		{
			name:  "consecutive variables",
			input: "$a$b$c",
			expected: []build.StringPart{
				{
					Text:       "",
					IsVariable: true,
					Expr:       &vast.Ident{Name: "a"},
				},
				{
					Text:       "",
					IsVariable: true,
					Expr:       &vast.Ident{Name: "b"},
				},
				{
					Text:       "",
					IsVariable: true,
					Expr:       &vast.Ident{Name: "c"},
				},
			},
		},
		{
			name:  "consecutive braced variables",
			input: "${a}${b}${c}",
			expected: []build.StringPart{
				{
					Text:       "",
					IsVariable: true,
					Expr:       &vast.Ident{Name: "a"},
				},
				{
					Text:       "",
					IsVariable: true,
					Expr:       &vast.Ident{Name: "b"},
				},
				{
					Text:       "",
					IsVariable: true,
					Expr:       &vast.Ident{Name: "c"},
				},
			},
		},
		{
			name:  "braced variable with special characters",
			input: "${user-name}",
			expected: []build.StringPart{
				{
					Text:       "",
					IsVariable: true,
					Expr:       &vast.Ident{Name: "user-name"},
				},
			},
		},
		{
			name:  "braced variable with spaces",
			input: "${user name}",
			expected: []build.StringPart{
				{
					Text:       "",
					IsVariable: true,
					Expr:       &vast.Ident{Name: "user name"},
				},
			},
		},
		{
			name:  "braced variable with dots",
			input: "${user.name}",
			expected: []build.StringPart{
				{
					Text:       "",
					IsVariable: true,
					Expr:       &vast.Ident{Name: "user.name"},
				},
			},
		},
		{
			name:     "empty string",
			input:    "",
			expected: []build.StringPart{},
		},
		{
			name:  "single dollar sign",
			input: "$",
			expected: []build.StringPart{
				{
					Text:       "",
					IsVariable: true,
					Expr:       &vast.Ident{Name: ""},
				},
			},
		},
		{
			name:  "dollar sign at end",
			input: "Hello $",
			expected: []build.StringPart{
				{
					Text:       "Hello ",
					IsVariable: false,
				},
				{
					Text:       "",
					IsVariable: true,
					Expr:       &vast.Ident{Name: ""},
				},
			},
		},
		{
			name:  "dollar sign followed by space",
			input: "Hello $ world",
			expected: []build.StringPart{
				{
					Text:       "Hello ",
					IsVariable: false,
				},
				{
					Text:       "",
					IsVariable: true,
					Expr:       &vast.Ident{Name: ""},
				},
				{
					Text:       " world",
					IsVariable: false,
				},
			},
		},
		{
			name:  "dollar sign followed by punctuation",
			input: "Hello $, world",
			expected: []build.StringPart{
				{
					Text:       "Hello ",
					IsVariable: false,
				},
				{
					Text:       "",
					IsVariable: true,
					Expr:       &vast.Ident{Name: ""},
				},
				{
					Text:       ", world",
					IsVariable: false,
				},
			},
		},
		{
			name:  "dollar sign followed by valid identifier",
			input: "Hello $a, world",
			expected: []build.StringPart{
				{
					Text:       "Hello ",
					IsVariable: false,
				},
				{
					Text:       "",
					IsVariable: true,
					Expr:       &vast.Ident{Name: "a"},
				},
				{
					Text:       ", world",
					IsVariable: false,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := build.ParseStringInterpolation(tt.input)

			// Check length
			assert.Equal(t, len(tt.expected), len(result), "Expected %d parts, got %d", len(tt.expected), len(result))

			// Check each part
			for i, expectedPart := range tt.expected {
				if i >= len(result) {
					t.Errorf("Missing part %d", i)
					continue
				}

				actualPart := result[i]

				// Check text
				assert.Equal(t, expectedPart.Text, actualPart.Text, "Part %d text mismatch", i)

				// Check IsVariable flag
				assert.Equal(t, expectedPart.IsVariable, actualPart.IsVariable, "Part %d IsVariable flag mismatch", i)

				// Check expression if it's a variable
				if expectedPart.IsVariable {
					assert.NotNil(t, actualPart.Expr, "Part %d should have expression", i)

					expectedIdent, ok := expectedPart.Expr.(*vast.Ident)
					assert.True(t, ok, "Part %d expected expression to be *vast.Ident", i)

					actualIdent, ok := actualPart.Expr.(*vast.Ident)
					assert.True(t, ok, "Part %d actual expression should be *vast.Ident", i)

					assert.Equal(t, expectedIdent.Name, actualIdent.Name, "Part %d variable name mismatch", i)
				} else {
					assert.Nil(t, actualPart.Expr, "Part %d should not have expression", i)
				}
			}
		})
	}
}
