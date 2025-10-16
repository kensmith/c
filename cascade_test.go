package main

import (
	//"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func assertClose(t *testing.T, expected, result float64) {
	const epsilon = 1e-9
	assert.InDelta(t, expected, result, epsilon)
}

func TestPlus(t *testing.T) {
	stack := NewStack()
	ops := NewOpMap()
	stack.Push(1234)
	stack.Push(2345)
	err := tryOpCascade("+", stack, ops)
	assert.Nil(t, err)
	assertClose(t, 3579, stack.Top())
}

func TestSwapOp(t *testing.T) {
	stack := NewStack()
	ops := NewOpMap()
	stack.Push(1234)
	stack.Push(2345)
	err := tryOpCascade("sw", stack, ops)
	assert.Nil(t, err)
	assertClose(t, 1234, stack.Top())
}

func TestMinus(t *testing.T) {
	stack := NewStack()
	ops := NewOpMap()
	stack.Push(1234)
	stack.Push(2345)
	err := tryOpCascade("-", stack, ops)
	assert.Nil(t, err)
	assertClose(t, -1111, stack.Top())
}

func TestMul(t *testing.T) {
	stack := NewStack()
	ops := NewOpMap()
	stack.Push(1234)
	stack.Push(2345)
	err := tryOpCascade("*", stack, ops)
	assert.Nil(t, err)
	assertClose(t, 2893730, stack.Top())
}

func TestDiv(t *testing.T) {
	stack := NewStack()
	ops := NewOpMap()
	stack.Push(1234)
	stack.Push(2345)
	err := tryOpCascade("/", stack, ops)
	assert.Nil(t, err)
	assertClose(t, 0.5262260127931769, stack.Top())
}

func TestAbs(t *testing.T) {
	stack := NewStack()
	ops := NewOpMap()
	stack.Push(-1234)
	err := tryOpCascade("abs", stack, ops)
	assert.Nil(t, err)
	assertClose(t, 1234, stack.Top())
}

func TestLeftShift(t *testing.T) {
	stack := NewStack()
	ops := NewOpMap()
	stack.Push(1234)
	stack.Push(2)
	err := tryOpCascade("<<", stack, ops)
	assert.Nil(t, err)
	assertClose(t, 4936, stack.Top())
}

/*
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

func TestMul(t *testing.T) {
	stack := NewStack()
	operators := NewOperatorMap()
	stack.Push(1234)
	stack.Push(2345)
	err := tryCascade("*", stack, operators)
	assert.Nil(t, err)
	result, err := stack.Pop()
	assertClose(t, 2893730.0, result)
	assert.Nil(t, err)
}

func TestDiv(t *testing.T) {
	stack := NewStack()
	operators := NewOperatorMap()
	stack.Push(1234)
	stack.Push(2345)
	err := tryCascade("/", stack, operators)
	assert.Nil(t, err)
	result, err := stack.Pop()
	assertClose(t, 0.526226012793, result)
	assert.Nil(t, err)
}

func TestLshft(t *testing.T) {
	stack := NewStack()
	operators := NewOperatorMap()
	stack.Push(1234)
	stack.Push(1)
	err := tryCascade("<<", stack, operators)
	assert.Nil(t, err)
	result, err := stack.Pop()
	assertClose(t, 2468, result)
	assert.Nil(t, err)
}

func TestRshft(t *testing.T) {
	stack := NewStack()
	operators := NewOperatorMap()
	stack.Push(1234)
	stack.Push(1)
	err := tryCascade(">>", stack, operators)
	assert.Nil(t, err)
	result, err := stack.Pop()
	assertClose(t, 617, result)
	assert.Nil(t, err)
}

func TestFact(t *testing.T) {
	stack := NewStack()
	operators := NewOperatorMap()
	stack.Push(10)
	err := tryCascade("!", stack, operators)
	assert.Nil(t, err)
	result, err := stack.Pop()
	assertClose(t, 3628800.0, result)
	assert.Nil(t, err)
}

func TestIncr(t *testing.T) {
	stack := NewStack()
	operators := NewOperatorMap()
	stack.Push(10)
	err := tryCascade("++", stack, operators)
	assert.Nil(t, err)
	result, err := stack.Pop()
	assertClose(t, 11.0, result)
	assert.Nil(t, err)
}

func TestDecr(t *testing.T) {
	stack := NewStack()
	operators := NewOperatorMap()
	stack.Push(10)
	err := tryCascade("--", stack, operators)
	assert.Nil(t, err)
	result, err := stack.Pop()
	assertClose(t, 9.0, result)
	assert.Nil(t, err)
}

func TestRand(t *testing.T) {
	stack := NewStack()
	operators := NewOperatorMap()
	err := tryCascade("r", stack, operators)
	assert.Nil(t, err)
	lhs, err := stack.Pop()
	assert.Nil(t, err)
	err = tryCascade("r", stack, operators)
	assert.Nil(t, err)
	rhs, err := stack.Pop()
	assert.NotEqual(t, lhs, rhs)
	assert.Nil(t, err)
}

func TestRandN(t *testing.T) {
	stack := NewStack()
	operators := NewOperatorMap()
	stack.Push(math.MaxInt32)
	err := tryCascade("rn", stack, operators)
	assert.Nil(t, err)
	lhs, err := stack.Pop()
	assert.Nil(t, err)
	stack.Push(math.MaxInt32)
	err = tryCascade("rn", stack, operators)
	assert.Nil(t, err)
	rhs, err := stack.Pop()
	assert.NotEqual(t, lhs, rhs)
	assert.Nil(t, err)
}

func TestPow10(t *testing.T) {
	stack := NewStack()
	operators := NewOperatorMap()
	stack.Push(11)
	err := tryCascade("pow10", stack, operators)
	assert.Nil(t, err)
	result, err := stack.Pop()
	assertClose(t, 1e11, result)
	assert.Nil(t, err)
}

func TestSignbit(t *testing.T) {
	stack := NewStack()
	operators := NewOperatorMap()
	stack.Push(11)
	err := tryCascade("signbit", stack, operators)
	assert.Nil(t, err)
	result, err := stack.Pop()
	assertClose(t, 0, result)
	assert.Nil(t, err)
	stack.Push(-11)
	err = tryCascade("signbit", stack, operators)
	assert.Nil(t, err)
	result, err = stack.Pop()
	assertClose(t, 1, result)
	assert.Nil(t, err)
}

func TestNeg(t *testing.T) {
	stack := NewStack()
	operators := NewOperatorMap()
	stack.Push(11)
	err := tryCascade("neg", stack, operators)
	assert.Nil(t, err)
	result, err := stack.Pop()
	assertClose(t, -11, result)
	assert.Nil(t, err)
}

func TestIlogb(t *testing.T) {
	stack := NewStack()
	operators := NewOperatorMap()
	stack.Push(11)
	err := tryCascade("ilogb", stack, operators)
	assert.Nil(t, err)
	result, err := stack.Pop()
	assertClose(t, 3, result)
	assert.Nil(t, err)
}

func TestIsInf(t *testing.T) {
	stack := NewStack()
	operators := NewOperatorMap()
	err := tryCascade("inf", stack, operators)
	assert.Nil(t, err)
	err = tryCascade("isinf", stack, operators)
	assert.Nil(t, err)
	result, err := stack.Pop()
	assertClose(t, 1, result)
	assert.Nil(t, err)
}

func TestIsNinf(t *testing.T) {
	stack := NewStack()
	operators := NewOperatorMap()
	err := tryCascade("ninf", stack, operators)
	assert.Nil(t, err)
	err = tryCascade("isninf", stack, operators)
	assert.Nil(t, err)
	result, err := stack.Pop()
	assertClose(t, 1, result)
	assert.Nil(t, err)
}

func TestIsNan(t *testing.T) {
	stack := NewStack()
	operators := NewOperatorMap()
	err := tryCascade("nan", stack, operators)
	assert.Nil(t, err)
	err = tryCascade("isnan", stack, operators)
	assert.Nil(t, err)
	result, err := stack.Pop()
	assertClose(t, 1, result)
	assert.Nil(t, err)
}

func TestMil(t *testing.T) {
	stack := NewStack()
	operators := NewOperatorMap()
	stack.Push(440)
	stack.Push(2)
	err := tryCascade("mil", stack, operators)
	assert.Nil(t, err)
	result, err := stack.Pop()
	assertClose(t, 2.22221856425409, result)
	assert.Nil(t, err)
}

func TestMph(t *testing.T) {
	stack := NewStack()
	operators := NewOperatorMap()
	stack.Push(440)
	stack.Push(2.22221856425409)
	err := tryCascade("mph", stack, operators)
	assert.Nil(t, err)
	result, err := stack.Pop()
	assertClose(t, 2.0, result)
	assert.Nil(t, err)
}

func TestSum(t *testing.T) {
	stack := NewStack()
	operators := NewOperatorMap()
	for i := range 5 {
		stack.Push(float64(i))
	}
	err := tryCascade("sum", stack, operators)
	assert.Nil(t, err)
	result, err := stack.Pop()
	assertClose(t, 10.0, result)
	assert.Nil(t, err)
}

func TestAvg(t *testing.T) {
	stack := NewStack()
	operators := NewOperatorMap()
	for i := range 6 {
		stack.Push(float64(i))
	}
	assert.Equal(t, "[ 0  1  2  3  4  5 ]", stack.String())
	err := tryCascade("avg", stack, operators)
	assert.Nil(t, err)
	result, err := stack.Pop()
	assertClose(t, 2.5, result)
	assert.Nil(t, err)
}

func TestSd(t *testing.T) {
	stack := NewStack()
	operators := NewOperatorMap()
	for i := range 6 {
		stack.Push(float64(i))
	}
	assert.Equal(t, "[ 0  1  2  3  4  5 ]", stack.String())
	err := tryCascade("sd", stack, operators)
	assert.Nil(t, err)
	result, err := stack.Pop()
	assertClose(t, 1.8708286933869707, result)
	assert.Nil(t, err)
}

func TestVar(t *testing.T) {
	stack := NewStack()
	operators := NewOperatorMap()
	for i := range 6 {
		stack.Push(float64(i))
	}
	assert.Equal(t, "[ 0  1  2  3  4  5 ]", stack.String())
	err := tryCascade("var", stack, operators)
	assert.Nil(t, err)
	result, err := stack.Pop()
	assertClose(t, 3.5, result)
	assert.Nil(t, err)
}

func TestMax(t *testing.T) {
	stack := NewStack()
	operators := NewOperatorMap()
	for i := range 6 {
		stack.Push(float64(i))
	}
	assert.Equal(t, "[ 0  1  2  3  4  5 ]", stack.String())
	err := tryCascade("max", stack, operators)
	assert.Nil(t, err)
	result, err := stack.Pop()
	assertClose(t, 5, result)
	assert.Nil(t, err)
}

func TestMin(t *testing.T) {
	stack := NewStack()
	operators := NewOperatorMap()
	for i := range 6 {
		stack.Push(float64(i))
	}
	assert.Equal(t, "[ 0  1  2  3  4  5 ]", stack.String())
	err := tryCascade("min", stack, operators)
	assert.Nil(t, err)
	result, err := stack.Pop()
	assertClose(t, 0, result)
	assert.Nil(t, err)
}

func TestLor(t *testing.T) {
	stack := NewStack()
	operators := NewOperatorMap()
	stack.Push(0.9999)
	err := tryCascade("c", stack, operators)
	assert.Nil(t, err)
	err = tryCascade("*", stack, operators)
	assert.Nil(t, err)
	err = tryCascade("lor", stack, operators)
	assert.Nil(t, err)
	result, err := stack.Pop()
	assertClose(t, 70.71244595187527, result)
	assert.Nil(t, err)
}

func TestCf(t *testing.T) {
	stack := NewStack()
	operators := NewOperatorMap()
	stack.Push(20)
	tryCascade("cf", stack, operators)
	result, _ := stack.Pop()
	assertClose(t, 68, result)
}

func TestFc(t *testing.T) {
	stack := NewStack()
	operators := NewOperatorMap()
	stack.Push(68)
	tryCascade("fc", stack, operators)
	result, _ := stack.Pop()
	assertClose(t, 20, result)
}

func TestFm(t *testing.T) {
	stack := NewStack()
	operators := NewOperatorMap()
	stack.Push(3)
	tryCascade("fm", stack, operators)
	result, _ := stack.Pop()
	assertClose(t, 0.9144000000036575, result)
}

func TestMf(t *testing.T) {
	stack := NewStack()
	operators := NewOperatorMap()
	stack.Push(0.9144000000036575)
	tryCascade("mf", stack, operators)
	result, _ := stack.Pop()
	assertClose(t, 3, result)
}

func TestFj(t *testing.T) {
	stack := NewStack()
	operators := NewOperatorMap()
	stack.Push(3000)
	tryCascade("fj", stack, operators)
	result, _ := stack.Pop()
	assertClose(t, 4067.453844994201, result)
}

func TestJf(t *testing.T) {
	stack := NewStack()
	operators := NewOperatorMap()
	stack.Push(4067.453844994201)
	tryCascade("jf", stack, operators)
	result, _ := stack.Pop()
	assertClose(t, 3000, result)
}

func TestGl(t *testing.T) {
	stack := NewStack()
	operators := NewOperatorMap()
	stack.Push(1)
	tryCascade("gl", stack, operators)
	result, _ := stack.Pop()
	assertClose(t, 3.785411784, result)
}

func TestLg(t *testing.T) {
	stack := NewStack()
	operators := NewOperatorMap()
	stack.Push(3.785411784)
	tryCascade("lg", stack, operators)
	result, _ := stack.Pop()
	assertClose(t, 1, result)
}

func TestPk(t *testing.T) {
	stack := NewStack()
	operators := NewOperatorMap()
	stack.Push(175)
	tryCascade("pk", stack, operators)
	result, _ := stack.Pop()
	assertClose(t, 79.37866475, result)
}

func TestKp(t *testing.T) {
	stack := NewStack()
	operators := NewOperatorMap()
	stack.Push(79.37866475)
	tryCascade("kp", stack, operators)
	result, _ := stack.Pop()
	assertClose(t, 175, result)
}

func TestHw(t *testing.T) {
	stack := NewStack()
	operators := NewOperatorMap()
	stack.Push(385)
	tryCascade("hw", stack, operators)
	result, _ := stack.Pop()
	assertClose(t, 287094.45072, result)
}

func TestWh(t *testing.T) {
	stack := NewStack()
	operators := NewOperatorMap()
	stack.Push(287094.45072)
	tryCascade("wh", stack, operators)
	result, _ := stack.Pop()
	assertClose(t, 385, result)
}

func TestPas(t *testing.T) {
	stack := NewStack()
	operators := NewOperatorMap()

	stack.Push(126)
	tryCascade("pas", stack, operators)
	result, _ := stack.Pop()
	assertClose(t, 16884.226578295245, result)

	stack.Push(132)
	tryCascade("pas", stack, operators)
	result, _ = stack.Pop()
	assertClose(t, 4222.299342426824, result)

	stack.Push(165)
	tryCascade("pas", stack, operators)
	result, _ = stack.Pop()
	assertClose(t, 2.065010118739787, result)
}

func TestPr(t *testing.T) {
	stack := NewStack()
	operators := NewOperatorMap()

	stack.Push(0)
	tryCascade("pr", stack, operators)
	result, _ := stack.Pop()
	assertClose(t, 29.9212524, result)

	stack.Push(1000)
	tryCascade("pr", stack, operators)
	result, _ = stack.Pop()
	assertClose(t, 28.85568264604063, result)

	stack.Push(10000)
	tryCascade("pr", stack, operators)
	result, _ = stack.Pop()
	assertClose(t, 20.576973132332288, result)

	stack.Push(35000)
	tryCascade("pr", stack, operators)
	result, _ = stack.Pop()
	assertClose(t, 7.040615836647221, result)
}
*/
