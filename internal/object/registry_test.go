// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package object

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOverloadedFunction(t *testing.T) {
	reg := NewRegistry()

	// add(num, num) -> num
	addSignature1 := &FunctionSignature{
		positionalParams: []*Parameter{
			{name: "a", typ: numberType},
			{name: "b", typ: numberType},
		},
		returnType: numberType,
		builtin: func(args ...Value) Value {
			a := args[0].(*NumberValue)
			b := args[1].(*NumberValue)
			return &NumberValue{Value: a.Value.Add(a.Value, b.Value)}
		},
	}

	reg.RegisterFunction("add", addSignature1)

	// add(str, str) -> str
	addSignature2 := &FunctionSignature{
		positionalParams: []*Parameter{
			{name: "a", typ: stringType},
			{name: "b", typ: stringType},
		},
		returnType: stringType,
		builtin: func(args ...Value) Value {
			a := args[0].(*StringValue)
			b := args[1].(*StringValue)
			return &StringValue{Value: a.Value + b.Value}
		},
	}
	reg.RegisterFunction("add", addSignature2)

	// add(num, str) -> str
	addSignature3 := &FunctionSignature{
		positionalParams: []*Parameter{
			{name: "a", typ: numberType},
			{name: "b", typ: stringType},
		},
		returnType: stringType,
		builtin: func(args ...Value) Value {
			a := args[0].(*NumberValue)
			b := args[1].(*StringValue)
			return &StringValue{Value: a.Text() + b.Value}
		},
	}
	reg.RegisterFunction("add", addSignature3)

	addFn, exists := reg.LookupFunction("add")
	assert.True(t, exists)

	result1, err := addFn.Call(
		[]Value{NewNumberValue("10"), NewNumberValue("20")},
		nil,
	)
	assert.NoError(t, err)
	assert.Equal(t, NewNumberValue("30"), result1)

	result2, err := addFn.Call(
		[]Value{&StringValue{Value: "Hello"}, &StringValue{Value: "World"}},
		nil,
	)
	assert.NoError(t, err)
	assert.Equal(t, &StringValue{Value: "HelloWorld"}, result2)

	result3, err := addFn.Call(
		[]Value{NewNumberValue("42"), &StringValue{Value: " is the answer"}},
		nil,
	)
	assert.NoError(t, err)
	assert.Equal(t, &StringValue{Value: "42 is the answer"}, result3)
}

func TestComplexFunction(t *testing.T) {
	reg := NewRegistry()

	// fn f(s: str = "default", i: num, args: ...any, {required ok: bool, name: str = "user1"})
	sig := &FunctionSignature{
		positionalParams: []*Parameter{
			{name: "s", typ: stringType, defaultValue: &StringValue{Value: "default"}},
			{name: "i", typ: numberType},
			{name: "args", typ: anyType, variadic: true},
		},
		namedParams: []*Parameter{
			{name: "ok", typ: boolType, required: true},
			{name: "name", typ: stringType, optional: true, defaultValue: &StringValue{Value: "user1"}},
		},
		builtin: func(args ...Value) Value {
			s := args[0].(*StringValue).Value
			i := args[1].(*NumberValue).Value
			// args[2:] 是可变参数
			ok := args[len(args)-2].(*BoolValue).Value
			name := args[len(args)-1].(*StringValue).Value

			return &StringValue{Value: fmt.Sprintf("s=%s, i=%s, ok=%v, name=%s", s, i, ok, name)}
		},
	}

	reg.RegisterFunction("f", sig)

	fn, ok := reg.LookupFunction("f")
	assert.True(t, ok)

	result, err := fn.Call(
		[]Value{
			&StringValue{Value: "hello"},
			NewNumberValue("42"),
		},
		map[string]Value{
			"ok":   &BoolValue{Value: true},
			"name": &StringValue{Value: "alice"},
		},
	)

	assert.NoError(t, err)
	assert.Equal(t, "s=hello, i=42, ok=true, name=alice", result.Text())

	result2, err := fn.Call(
		[]Value{
			NewNumberValue("100"),
		},
		map[string]Value{
			"ok": &BoolValue{Value: false},
		},
	)

	assert.NoError(t, err)
	assert.Equal(t, "s=default, i=100, ok=false, name=user1", result2.Text())
}
