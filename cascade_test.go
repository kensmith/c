package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func assertClose(t *testing.T, expected, result float64) {
	const epsilon = 1e-9
	assert.InDelta(t, expected, result, epsilon)
}

func TestPlus(t *testing.T) {
	stack := NewStack()
	operators := NewOperatorMap()
	stack.Push(1234)
	stack.Push(2345)
	err := tryCascade("+", stack, operators)
	assert.Nil(t, err)
	result, err := stack.Pop()
	assertClose(t, 3579.0, result)
	assert.Nil(t, err)
}

func TestMinus(t *testing.T) {
	stack := NewStack()
	operators := NewOperatorMap()
	stack.Push(1234)
	stack.Push(2345)
	err := tryCascade("-", stack, operators)
	assert.Nil(t, err)
	result, err := stack.Pop()
	assertClose(t, -1111.0, result)
	assert.Nil(t, err)
}
