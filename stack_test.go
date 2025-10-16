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
	results, err := stack.PopN(3)
	assert.Nil(t, err)
	assert.Equal(t, []float64{3, 2, 1}, results)
	stack.Push(1)
	stack.Push(2)
	stack.Push(3)
	results, err = stack.PopN(2)
	assert.Nil(t, err)
	assert.Equal(t, []float64{3, 2}, results)
}

func TestPopR(t *testing.T) {
	stack := NewStack()
	stack.Push(1)
	stack.Push(2)
	stack.Push(3)
	results, err := stack.PopR(3)
	assert.Nil(t, err)
	assert.Equal(t, []float64{1, 2, 3}, results)
	stack.Push(1)
	stack.Push(2)
	stack.Push(3)
	results, err = stack.PopR(2)
	assert.Nil(t, err)
	assert.Equal(t, []float64{2, 3}, results)
}

func TestSwap(t *testing.T) {
	stack := NewStack()
	stack.Push(1)
	stack.Push(2)
	stack.Push(3)
	_ = stack.Swap()
	results, err := stack.PopN(2)
	assert.Nil(t, err)
	assert.Equal(t, []float64{2, 3}, results)
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

func TestStringer(t *testing.T) {
	stack := NewStack()
	stack.Push(1)
	stack.Push(2)
	stack.Push(3)
	assert.Equal(t, "[ 1  2  3 ]", stack.String())
}

func TestSort(t *testing.T) {
	stack := NewStack()
	stack.Push(3)
	stack.Push(1)
	stack.Push(2)
	stack.Sort()
	assert.Equal(t, "[ 1  2  3 ]", stack.String())
}

func TestCopy(t *testing.T) {
	stack := NewStack()
	stack.Push(3)
	stack.Push(1)
	stack.Push(2)
	arr := stack.Copy()
	assert.Equal(t, []float64{3, 1, 2}, arr)
}
