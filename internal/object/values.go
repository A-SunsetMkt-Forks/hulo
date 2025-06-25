// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package object

import (
	"math/big"
	"strconv"
)

var (
	NULL  = &NullValue{}
	TRUE  = &BoolValue{Value: true}
	FALSE = &BoolValue{Value: false}
)

// NullValue represents the null value.
type NullValue struct{}

func (n *NullValue) Text() string {
	return "null"
}

func (n *NullValue) Interface() any {
	return nil
}

func (n *NullValue) Type() Type {
	// TODO: return voidType
	return nil
}

// NumberValue represents the number value.
type NumberValue struct {
	// TODO 根据精度选择存储模型
	Value *big.Float
}

// NewNumberValue creates a new number value from the given string.
func NewNumberValue(value string) *NumberValue {
	num, _ := strconv.ParseFloat(value, 64)
	return &NumberValue{Value: big.NewFloat(num)}
}

func (n *NumberValue) Text() string {
	return n.Value.String()
}

func (n *NumberValue) Interface() any {
	return n.Value
}

func (n *NumberValue) Type() Type {
	return numberType
}

// BoolValue represents the boolean value.
type BoolValue struct {
	Value bool
}

func (b *BoolValue) Text() string {
	return strconv.FormatBool(b.Value)
}

func (b *BoolValue) Interface() any {
	return b.Value
}

func (b *BoolValue) Type() Type {
	return boolType
}

// ErrorValue represents the error value.
type ErrorValue struct {
	Value string
}

func (e *ErrorValue) Text() string {
	return e.Value
}

func (e *ErrorValue) Interface() any {
	return e.Value
}

func (e *ErrorValue) Type() Type {
	return errorType
}

// StringValue represents the string value.
type StringValue struct {
	Value string
}

func (s *StringValue) Text() string {
	return s.Value
}

func (s *StringValue) Interface() any {
	return s.Value
}

func (s *StringValue) Type() Type {
	return stringType
}
