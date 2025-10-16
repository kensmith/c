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
	assert.GreaterOrEqual(t, stack.Top(), 0.0)
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

func TestMax(t *testing.T) {
	stack := NewStack()
	ops := NewOps()
	stack.Push(1000)

	for range 10 {
		stack.Push(99)
		err := ops.Run("rn", stack)
		assert.Nil(t, err)
	}

	err := ops.Run("max", stack)
	assert.Nil(t, err)
	assertClose(t, 1000, stack.Top())
}

func TestMin(t *testing.T) {
	stack := NewStack()
	ops := NewOps()
	stack.Push(-1000)

	for range 10 {
		stack.Push(99)
		err := ops.Run("rn", stack)
		assert.Nil(t, err)
	}

	err := ops.Run("min", stack)
	assert.Nil(t, err)
	assertClose(t, -1000, stack.Top())
}

func TestLor(t *testing.T) {
	stack := NewStack()
	ops := NewOps()
	stack.Push(constants["c"] * 0.999)
	err := ops.Run("lor", stack)
	assert.Nil(t, err)
	assertClose(t, 22.36627204212937, stack.Top())
}

func TestMissingOp(t *testing.T) {
	stack := NewStack()
	ops := NewOps()
	stack.Push(10)
	err := ops.Run("missing op", stack)
	assert.NotNil(t, err)
	assertClose(t, 10, stack.Top())
}

func TestFc(t *testing.T) {
	stack := NewStack()
	ops := NewOps()
	stack.Push(68)
	err := ops.Run("fc", stack)
	assert.Nil(t, err)
	assertClose(t, 20, stack.Top())
}

func TestCf(t *testing.T) {
	stack := NewStack()
	ops := NewOps()
	stack.Push(20)
	err := ops.Run("cf", stack)
	assert.Nil(t, err)
	assertClose(t, 68, stack.Top())
}

func TestFm(t *testing.T) {
	stack := NewStack()
	ops := NewOps()
	stack.Push(10)
	err := ops.Run("fm", stack)
	assert.Nil(t, err)
	assertClose(t, 3.048, stack.Top())
}

func TestMf(t *testing.T) {
	stack := NewStack()
	ops := NewOps()
	stack.Push(3.048)
	err := ops.Run("mf", stack)
	assert.Nil(t, err)
	assertClose(t, 10, stack.Top())
}

func TestFj(t *testing.T) {
	stack := NewStack()
	ops := NewOps()
	stack.Push(3000)
	err := ops.Run("fj", stack)
	assert.Nil(t, err)
	assertClose(t, 4067.453844994201, stack.Top())
}

func TestJf(t *testing.T) {
	stack := NewStack()
	ops := NewOps()
	stack.Push(4067.453844994201)
	err := ops.Run("jf", stack)
	assert.Nil(t, err)
	assertClose(t, 3000, stack.Top())
}

func TestGl(t *testing.T) {
	stack := NewStack()
	ops := NewOps()
	stack.Push(1)
	err := ops.Run("gl", stack)
	assert.Nil(t, err)
	assertClose(t, 3.785411783999999890, stack.Top())
}

func TestLg(t *testing.T) {
	stack := NewStack()
	ops := NewOps()
	stack.Push(3.785411783999999890)
	err := ops.Run("lg", stack)
	assert.Nil(t, err)
	assertClose(t, 1, stack.Top())
}

func TestPk(t *testing.T) {
	stack := NewStack()
	ops := NewOps()
	stack.Push(175)
	err := ops.Run("pk", stack)
	assert.Nil(t, err)
	assertClose(t, 79.37866475, stack.Top())
}

func TestKp(t *testing.T) {
	stack := NewStack()
	ops := NewOps()
	stack.Push(79.37866475)
	err := ops.Run("kp", stack)
	assert.Nil(t, err)
	assertClose(t, 175, stack.Top())
}

func TestHw(t *testing.T) {
	stack := NewStack()
	ops := NewOps()
	stack.Push(1)
	err := ops.Run("hw", stack)
	assert.Nil(t, err)
	assertClose(t, 745.699872, stack.Top())
}
