package main

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func assertClose(t *testing.T, expected, result float64) {
	const epsilon = 1e-9
	assert.InDelta(t, expected, result, epsilon)
}

func TestPlus(t *testing.T) {
	stack := NewStack()
	ops := NewOps()
	stack.Push(1234)
	stack.Push(2345)
	err := ops.Run("+", stack)
	assert.Nil(t, err)
	assertClose(t, 3579, stack.Top())
}

func TestSwapOp(t *testing.T) {
	stack := NewStack()
	ops := NewOps()
	stack.Push(1234)
	stack.Push(2345)
	err := ops.Run("sw", stack)
	assert.Nil(t, err)
	assertClose(t, 1234, stack.Top())
}

func TestMinus(t *testing.T) {
	stack := NewStack()
	ops := NewOps()
	stack.Push(1234)
	stack.Push(2345)
	err := ops.Run("-", stack)
	assert.Nil(t, err)
	assertClose(t, -1111, stack.Top())
}

func TestMul(t *testing.T) {
	stack := NewStack()
	ops := NewOps()
	stack.Push(1234)
	stack.Push(2345)
	err := ops.Run("*", stack)
	assert.Nil(t, err)
	assertClose(t, 2893730, stack.Top())
}

func TestDiv(t *testing.T) {
	stack := NewStack()
	ops := NewOps()
	stack.Push(1234)
	stack.Push(2345)
	err := ops.Run("/", stack)
	assert.Nil(t, err)
	assertClose(t, 0.5262260127931769, stack.Top())
}

func TestAbs(t *testing.T) {
	stack := NewStack()
	ops := NewOps()
	stack.Push(-1234)
	err := ops.Run("abs", stack)
	assert.Nil(t, err)
	assertClose(t, 1234, stack.Top())
}

func TestLeftShift(t *testing.T) {
	stack := NewStack()
	ops := NewOps()
	stack.Push(1234)
	stack.Push(2)
	err := ops.Run("<<", stack)
	assert.Nil(t, err)
	assertClose(t, 4936, stack.Top())
}

func TestRightShift(t *testing.T) {
	stack := NewStack()
	ops := NewOps()
	stack.Push(1234)
	stack.Push(2)
	err := ops.Run(">>", stack)
	assert.Nil(t, err)
	assertClose(t, 308.5, stack.Top())
}

func TestFactorial(t *testing.T) {
	stack := NewStack()
	ops := NewOps()
	stack.Push(10)
	err := ops.Run("!", stack)
	assert.Nil(t, err)
	assertClose(t, 3628800, stack.Top())
}

func TestIncr(t *testing.T) {
	stack := NewStack()
	ops := NewOps()
	stack.Push(10)
	err := ops.Run("++", stack)
	assert.Nil(t, err)
	assertClose(t, 11, stack.Top())
}

func TestDecr(t *testing.T) {
	stack := NewStack()
	ops := NewOps()
	stack.Push(10)
	err := ops.Run("--", stack)
	assert.Nil(t, err)
	assertClose(t, 9, stack.Top())
}

func TestRandN(t *testing.T) {
	stack := NewStack()
	ops := NewOps()
	stack.Push(100)
	err := ops.Run("rn", stack)
	assert.Nil(t, err)
	assert.Greater(t, stack.Top(), 0.0)
	assert.Less(t, stack.Top(), 100.0)
}

func TestRand(t *testing.T) {
	stack := NewStack()
	ops := NewOps()
	err := ops.Run("r", stack)
	assert.Nil(t, err)
	assert.Greater(t, stack.Top(), 0.0)
	assert.Less(t, stack.Top(), float64(_defaultMaxRand))
}

func TestPow10(t *testing.T) {
	stack := NewStack()
	ops := NewOps()
	stack.Push(10)
	err := ops.Run("pow10", stack)
	assert.Nil(t, err)
	assertClose(t, 1e10, stack.Top())
}

func TestSignbit(t *testing.T) {
	stack := NewStack()
	ops := NewOps()
	stack.Push(-10)
	err := ops.Run("signbit", stack)
	assert.Nil(t, err)
	assertClose(t, 1, stack.Top())
}

func TestNeg(t *testing.T) {
	stack := NewStack()
	ops := NewOps()
	stack.Push(-10)
	err := ops.Run("neg", stack)
	assert.Nil(t, err)
	assertClose(t, 10, stack.Top())
}

func TestIlogb(t *testing.T) {
	stack := NewStack()
	ops := NewOps()
	stack.Push(-10)
	err := ops.Run("ilogb", stack)
	assert.Nil(t, err)
	assertClose(t, 3, stack.Top())
}

func TestIsInf(t *testing.T) {
	stack := NewStack()
	ops := NewOps()

	stack.Push(math.Inf(1))
	err := ops.Run("isinf", stack)
	assert.Nil(t, err)
	assertClose(t, 1, stack.Top())

	stack.Push(math.Inf(-1))
	err = ops.Run("isinf", stack)
	assert.Nil(t, err)
	assertClose(t, 0, stack.Top())
}

func TestIsNInf(t *testing.T) {
	stack := NewStack()
	ops := NewOps()

	stack.Push(math.Inf(1))
	err := ops.Run("isninf", stack)
	assert.Nil(t, err)
	assertClose(t, 0, stack.Top())

	stack.Push(math.Inf(-1))
	err = ops.Run("isninf", stack)
	assert.Nil(t, err)
	assertClose(t, 1, stack.Top())
}

func TestIsNan(t *testing.T) {
	stack := NewStack()
	ops := NewOps()

	stack.Push(math.NaN())
	err := ops.Run("isnan", stack)
	assert.Nil(t, err)
	assertClose(t, 1, stack.Top())

	stack.Push(0)
	err = ops.Run("isnan", stack)
	assert.Nil(t, err)
	assertClose(t, 0, stack.Top())
}

func TestJn(t *testing.T) {
	stack := NewStack()
	ops := NewOps()

	stack.Push(1)
	stack.Push(2)
	err := ops.Run("jn", stack)
	assert.Nil(t, err)
	assertClose(t, 0.11490348493190049, stack.Top())
}

func TestMil(t *testing.T) {
	stack := NewStack()
	ops := NewOps()
	stack.Push(440)
	stack.Push(2)
	err := ops.Run("mil", stack)
	assert.Nil(t, err)
	assertClose(t, 2.22221856425409, stack.Top())
}

func TestMPH(t *testing.T) {
	stack := NewStack()
	ops := NewOps()
	stack.Push(440)
	stack.Push(2.22221856425409)
	err := ops.Run("mph", stack)
	assert.Nil(t, err)
	assertClose(t, 2, stack.Top())
}

func TestSum(t *testing.T) {
	stack := NewStack()
	ops := NewOps()
	for i := range 5 {
		stack.Push(float64(i))
	}
	err := ops.Run("sum", stack)
	assert.Nil(t, err)
	assertClose(t, 10, stack.Top())
	assert.Equal(t, stack.Len(), 6)
}

func TestAvg(t *testing.T) {
	stack := NewStack()
	ops := NewOps()
	for i := range 6 {
		stack.Push(float64(i))
	}
	err := ops.Run("avg", stack)
	assert.Nil(t, err)
	assertClose(t, 2.5, stack.Top())
	assert.Equal(t, stack.Len(), 7)
}

func TestSd(t *testing.T) {
	stack := NewStack()
	ops := NewOps()
	for i := range 6 {
		stack.Push(float64(i))
	}
	err := ops.Run("sd", stack)
	assert.Nil(t, err)
	assertClose(t, 1.8708286933869707, stack.Top())
}

func TestVar(t *testing.T) {
	stack := NewStack()
	ops := NewOps()
	for i := range 6 {
		stack.Push(float64(i))
	}
	err := ops.Run("var", stack)
	assert.Nil(t, err)
	assertClose(t, 3.5, stack.Top())
}
