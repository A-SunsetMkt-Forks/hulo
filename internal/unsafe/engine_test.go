// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package unsafe_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hulo-lang/hulo/internal/unsafe"
	"github.com/stretchr/testify/assert"
)

func TestParseLoopStmt(t *testing.T) {
	unsafe.Parse(`{% loop item in items %}
	{# this is a comment #}
	echo {{item}}
{% endloop %}

echo {{ sum 1 2 | to_str }}`)
}

func TestParseIfStmt(t *testing.T) {
	unsafe.Parse(`{% if true %}
	echo hello
{% endif %}`)
}

func TestHelloWorld(t *testing.T) {
	engine := unsafe.NewTemplateEngine()
	engine.SetVariable("name", "World")
	result, err := engine.Execute(`Hello {{ name }}`)
	assert.NoError(t, err)
	fmt.Println(result)
}

func TestTemplateEngine(t *testing.T) {
	engine := unsafe.NewTemplateEngine()
	engine.SetVariable("items", []string{"hello", "world"})
	result, err := engine.Execute(`{% loop item in items %}
	{# this is a comment #}
	echo {{item}}
{% endloop %}

echo {{ sum 1 2 | to_str }}`)
	assert.NoError(t, err)
	fmt.Println(result)
}

func TestPipelines(t *testing.T) {
	engine := unsafe.NewTemplateEngine()
	engine.SetFunction("add", func(args ...any) any {
		// Smart type conversion
		toFloat := func(v any) float64 {
			switch val := v.(type) {
			case int:
				return float64(val)
			case float64:
				return val
			case string:
				if f, err := strconv.ParseFloat(val, 64); err == nil {
					return f
				}
			}
			return 0
		}

		if len(args) < 2 {
			return 0
		}
		return toFloat(args[0]) + toFloat(args[1])
	})
	engine.SetFunction("to_str", func(args ...any) string {
		return fmt.Sprintf("%v", args[0])
	})
	result, err := engine.Execute(`{{ sum 1 2 | add 3 | add -3 | to_str }}`)
	assert.NoError(t, err)
	fmt.Println(result)
}

func TestMap(t *testing.T) {
	engine := unsafe.NewTemplateEngine()
	engine.SetVariable("items", map[any]any{"6": true, true: "hello"})
	result, err := engine.Execute(`{% loop key, value of items %}
	{# this is a comment #}
	echo {{key}} {{value}}
{% endloop %}

echo {{ items | to_str }}`)
	assert.NoError(t, err)
	fmt.Println(result)
}

func TestMacroWithoutArgs(t *testing.T) {
	engine := unsafe.NewTemplateEngine()
	result, err := engine.Execute(`{% macro greet %}Hello, World!{% endmacro %}

echo {{ greet() }}`)
	assert.NoError(t, err)
	fmt.Println(result)
}

func TestMacroWithArgs(t *testing.T) {
	engine := unsafe.NewTemplateEngine()
	result, err := engine.Execute(`{% macro greet(name) %}Hello, {{name}}!{% endmacro %}

echo {{ greet "World" }}`)
	assert.NoError(t, err)
	fmt.Println(result)
}

func TestParseIfElseStmt(t *testing.T) {
	engine := unsafe.NewTemplateEngine()
	result, err := engine.Execute(`{% if true %}
	echo hello
{% else %}
	echo world
{% endif %}`)
	assert.NoError(t, err)
	fmt.Println(result)
}

func TestParseNestIfStmt(t *testing.T) {
	engine := unsafe.NewTemplateEngine()
	result, err := engine.Execute(`{% if true %}
	echo hello1
	{% if true %}
		echo hello2
	{% else %}
		echo world2
	{% endif %}
{% else %}
	echo world1
{% endif %}`)
	assert.NoError(t, err)
	fmt.Println(result)
}

func TestOperators(t *testing.T) {
	engine := unsafe.NewTemplateEngine()

	// Test comparison operators
	result, err := engine.Execute(`{% if 5 > 3 %}
	echo "5 is greater than 3"
{% endif %}`)
	assert.NoError(t, err)
	fmt.Println("Test 5 > 3:", result)

	result, err = engine.Execute(`{% if 3 < 5 %}
	echo "3 is less than 5"
{% endif %}`)
	assert.NoError(t, err)
	fmt.Println("Test 3 < 5:", result)

	result, err = engine.Execute(`{% if 5 == 5 %}
	echo "5 equals 5"
{% endif %}`)
	assert.NoError(t, err)
	fmt.Println("Test 5 == 5:", result)

	result, err = engine.Execute(`{% if 5 != 3 %}
	echo "5 is not equal to 3"
{% endif %}`)
	assert.NoError(t, err)
	fmt.Println("Test 5 != 3:", result)

	// Test logical operators
	result, err = engine.Execute(`{% if true && true %}
	echo "Both conditions are true"
{% endif %}`)
	assert.NoError(t, err)
	fmt.Println("Test true && true:", result)

	result, err = engine.Execute(`{% if true || false %}
	echo "At least one condition is true"
{% endif %}`)
	assert.NoError(t, err)
	fmt.Println("Test true || false:", result)

	result, err = engine.Execute(`{% if !false %}
	echo "Not false is true"
{% endif %}`)
	assert.NoError(t, err)
	fmt.Println("Test !false:", result)

	// Test complex expressions
	result, err = engine.Execute(`{% if 5 > 3 && 10 < 20 %}
	echo "Complex condition is true"
{% endif %}`)
	assert.NoError(t, err)
	fmt.Println("Test complex condition:", result)
}
