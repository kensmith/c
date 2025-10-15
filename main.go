package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/chzyer/readline"
	"github.com/expr-lang/expr"
)

func pop1(stack *[]float64) (float64, error) {
	size := len(*stack)
	if size < 1 {
		return 0.0, fmt.Errorf("insufficient stack")
	}
	top := (*stack)[size-1]
	*stack = (*stack)[:size-1]
	return top, nil
}

func pop2(stack *[]float64) (float64, float64, error) {
	size := len(*stack)
	if size < 2 {
		return 0.0, 0.0, fmt.Errorf("insufficient stack")
	}
	rhs := (*stack)[size-1]
	lhs := (*stack)[size-2]
	*stack = (*stack)[:size-2]
	return lhs, rhs, nil
}

type Operator func(*[]float64) (float64, error)

var operators = map[string]Operator{
	"+": func(stack *[]float64) (float64, error) {
		lhs, rhs, err := pop2(stack)
		if err != nil {
			return 0.0, err
		}
		return lhs + rhs, nil
	},
	"-": func(stack *[]float64) (float64, error) {
		lhs, rhs, err := pop2(stack)
		if err != nil {
			return 0.0, err
		}
		return lhs - rhs, nil
	},
	"*": func(stack *[]float64) (float64, error) {
		lhs, rhs, err := pop2(stack)
		if err != nil {
			return 0.0, err
		}
		return lhs * rhs, nil
	},
	"/": func(stack *[]float64) (float64, error) {
		lhs, rhs, err := pop2(stack)
		if err != nil {
			return 0.0, err
		}
		return lhs / rhs, nil
	},
	"<<": func(stack *[]float64) (float64, error) {
		lhs, rhs, err := pop2(stack)
		if err != nil {
			return 0.0, err
		}
		return lhs * math.Pow(2, rhs), nil
	},
	">>": func(stack *[]float64) (float64, error) {
		lhs, rhs, err := pop2(stack)
		if err != nil {
			return 0.0, err
		}
		return lhs / math.Pow(2, rhs), nil
	},
	"!": func(stack *[]float64) (float64, error) {
		top, err := pop1(stack)
		if err != nil {
			return 0.0, err
		}
		return top * math.Gamma(top), nil
	},
	"++": func(stack *[]float64) (float64, error) {
		top, err := pop1(stack)
		if err != nil {
			return 0.0, err
		}
		return top + 1, nil
	},
	"--": func(stack *[]float64) (float64, error) {
		top, err := pop1(stack)
		if err != nil {
			return 0.0, err
		}
		return top - 1, nil
	},
	"pow10": func(stack *[]float64) (float64, error) {
		top, err := pop1(stack)
		if err != nil {
			return 0.0, err
		}
		return math.Pow10(int(top)), nil
	},
	"signbit": func(stack *[]float64) (float64, error) {
		top, err := pop1(stack)
		if err != nil {
			return 0.0, err
		}
		if math.Signbit(top) {
			return 1, nil
		}

		return 0, nil
	},
	"ilogb": func(stack *[]float64) (float64, error) {
		top, err := pop1(stack)
		if err != nil {
			return 0.0, err
		}
		return float64(math.Ilogb(top)), nil
	},
	/*
	   func Sincos(x float64) (sin, cos float64)

	   func Copysign(f, sign float64) float64
	   func FMA(x, y, z float64) float64
	   func Float64bits(f float64) uint64
	   func Float64frombits(b uint64) float64
	   func Frexp(f float64) (frac float64, exp int)
	   func Inf(sign int) float64
	   func IsInf(f float64, sign int) bool
	   func IsNaN(f float64) (is bool)
	   func Jn(n int, x float64) float64
	   func Ldexp(frac float64, exp int) float64
	   func Modf(f float64) (int float64, frac float64)
	   func NaN() float64
	   func Pow10(n int) float64
	   func Yn(n int, x float64) float64
	*/

}

var unaryFunctions = map[string]func(float64) float64{
	"abs":         math.Abs,
	"acos":        math.Acos,
	"acosh":       math.Acosh,
	"asin":        math.Asin,
	"asinh":       math.Asinh,
	"atan":        math.Atan,
	"cbrt":        math.Cbrt,
	"ceil":        math.Ceil,
	"cos":         math.Cos,
	"cosh":        math.Cosh,
	"erf":         math.Erf,
	"erfc":        math.Erfc,
	"erfcinv":     math.Erfcinv,
	"erfinv":      math.Erfinv,
	"exp":         math.Exp,
	"exp2":        math.Exp2,
	"expm1":       math.Expm1,
	"floor":       math.Floor,
	"gamma":       math.Gamma,
	"j0":          math.J0,
	"j1":          math.J1,
	"log":         math.Log,
	"log10":       math.Log10,
	"log1p":       math.Log1p,
	"log2":        math.Log2,
	"logb":        math.Logb,
	"round":       math.Round,
	"roundtoeven": math.RoundToEven,
	"sin":         math.Sin,
	"sinh":        math.Sinh,
	"sqrt":        math.Sqrt,
	"tan":         math.Tan,
	"tanh":        math.Tanh,
	"trunc":       math.Trunc,
	"y0":          math.Y0,
	"y1":          math.Y1,
}

var binaryFunctions = map[string]func(float64, float64) float64{
	"%":         math.Mod,
	"**":        math.Pow,
	"^":         math.Pow,
	"atan2":     math.Atan2,
	"dim":       math.Dim,
	"hypot":     math.Hypot,
	"max":       math.Max,
	"min":       math.Min,
	"mod":       math.Mod,
	"nextafter": math.Nextafter,
	"pow":       math.Pow,
	"remainder": math.Remainder,
}

var constants = map[string]float64{
	"c": 299792458,
}

func main() {
	for name, uFunc := range unaryFunctions {
		operators[name] = func(stack *[]float64) (float64, error) {
			top, err := pop1(stack)
			if err != nil {
				return 0.0, err
			}
			return uFunc(top), nil
		}
	}

	for name, bFunc := range binaryFunctions {
		operators[name] = func(stack *[]float64) (float64, error) {
			lhs, rhs, err := pop2(stack)
			if err != nil {
				return 0.0, err
			}
			return bFunc(lhs, rhs), nil
		}
	}

	for name, c := range constants {
		operators[name] = func(stack *[]float64) (float64, error) {
			return c, nil
		}
	}

	shell, err := readline.NewEx(&readline.Config{
		Prompt:      "[  ]> ",
		HistoryFile: ".history",
	})
	if err != nil {
		panic(err)
	}
	defer func() {
		err := shell.Close()
		if err != nil {
			fmt.Println(err)
		}
	}()

	stack := []float64{}

	for {
		line, err := shell.Readline()
		if err != nil {
			// normal exit due to ctrl-c, ctrl-d
			return
		}
		lineTrimmed := strings.TrimSpace(line)
		rawOperator, ok := operators[lineTrimmed]
		if ok {
			result, err := rawOperator(&stack)
			if err != nil {
				fmt.Println(err)
			} else {
				stack = append(stack, result)
			}
		} else {
			output, err := expr.Eval(line, nil)
			if err != nil {
				fmt.Println(err)
				continue
			}

			outputStr := fmt.Sprintf("%v", output)
			value, err := strconv.ParseFloat(outputStr, 64)
			if err != nil {
				// not pushing to the stack
				continue
			}
			stack = append(stack, value)
		}
		stackSize := len(stack)
		var b strings.Builder
		fmt.Fprintf(&b, "[ ")
		for i, n := range stack {
			fmt.Fprintf(&b, "%g", n)
			if i < stackSize-1 {
				fmt.Fprintf(&b, "  ")
			}
		}
		fmt.Fprintf(&b, " ]> ")
		shell.SetPrompt(b.String())
	}
}
