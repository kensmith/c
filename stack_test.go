package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPush(t *testing.T) {
	stack := NewStack()
	stack.Push(10)
	stack.Push(10)
	stack.Push(10)
	assert.Equal(t, stack.Len(), 3, nil)
}

func TestPop(t *testing.T) {
	stack := NewStack()
	stack.Push(1)
	stack.Push(2)
	stack.Push(3)
	results, err := stack.Pop(3)
	assert.Nil(t, err)
	assert.Equal(t, []float64{3, 2, 1}, results, nil)
	stack.Push(1)
	stack.Push(2)
	stack.Push(3)
	results, err = stack.Pop(2)
	assert.Nil(t, err)
	assert.Equal(t, []float64{3, 2}, results, nil)
}

func TestPopR(t *testing.T) {
	stack := NewStack()
	stack.Push(1)
	stack.Push(2)
	stack.Push(3)
	results, err := stack.PopR(3)
	assert.Nil(t, err)
	assert.Equal(t, []float64{1, 2, 3}, results, nil)
	stack.Push(1)
	stack.Push(2)
	stack.Push(3)
	results, err = stack.PopR(2)
	assert.Nil(t, err)
	assert.Equal(t, []float64{2, 3}, results, nil)
}

func TestSwap(t *testing.T) {
	stack := NewStack()
	stack.Push(1)
	stack.Push(2)
	stack.Push(3)
	stack.Swap()
	results, err := stack.Pop(2)
	assert.Nil(t, err)
	assert.Equal(t, []float64{2, 3}, results, nil)
}

func TestClear(t *testing.T) {
	stack := NewStack()
	stack.Push(1)
	stack.Push(2)
	stack.Push(3)
	assert.Equal(t, stack.Len(), 3, nil)
	stack.Clear()
	assert.Equal(t, stack.Len(), 0, nil)
}
