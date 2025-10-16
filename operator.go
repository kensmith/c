package main

import (
	"crypto/rand"
	"fmt"
	"math"
	"math/big"
	"slices"

	"github.com/eclesh/welford"
)

/*
* Operator needs to remember:
* - the stack
* - the string which invokes it
* - the function to invoke
*   - this should return ([]float64, error)
* - help text describing it for help display
 */

type (
	Floats       []float64
	Operator     func(*Stack) (float64, error)
	OperatorMap  map[string]Operator
	OperatorFunc func(*Stack) (Floats, error)

	Op struct {
		name string
		doc  string
		f    OperatorFunc
	}

	OpMap map[string]Op
)

func wrapUnaryOp(name, doc string, f func(float64) float64) Op {
	return Op{
		name,
		doc,
		func(stack *Stack) (Floats, error) {
			top, err := stack.Pop()
			if err != nil {
				return nil, err
			}
			result := f(top)
			return Floats{result}, nil
		},
	}
}

var (
	factorialOp = Op{
		"!",
		"factorial (kind of, by way of gamma function)",
		func(stack *Stack) (Floats, error) {
			top, err := stack.Pop()
			if err != nil {
				return nil, err
			}
			return Floats{top * math.Gamma(top)}, nil
		},
	}

	mulOp = Op{
		"*",
		"multiplication",
		func(stack *Stack) (Floats, error) {
			elems, err := stack.PopR(2)
			if err != nil {
				return nil, err
			}
			return Floats{elems[0] * elems[1]}, nil
		},
	}

	plusOp = Op{
		"+",
		"addition",
		func(stack *Stack) (Floats, error) {
			elems, err := stack.PopR(2)
			if err != nil {
				return nil, err
			}
			return Floats{elems[0] + elems[1]}, nil
		},
	}

	incrOp = Op{
		"++",
		"increment",
		func(stack *Stack) (Floats, error) {
			top, err := stack.Pop()
			if err != nil {
				return nil, err
			}
			return Floats{top + 1.0}, nil
		},
	}

	minusOp = Op{
		"-",
		"subtraction",
		func(stack *Stack) (Floats, error) {
			elems, err := stack.PopR(2)
			if err != nil {
				return nil, err
			}
			return Floats{elems[0] - elems[1]}, nil
		},
	}

	decrOp = Op{
		"--",
		"decrement",
		func(stack *Stack) (Floats, error) {
			top, err := stack.Pop()
			if err != nil {
				return nil, err
			}

			return Floats{top - 1.0}, nil
		},
	}

	divOp = Op{
		"/",
		"division",
		func(stack *Stack) (Floats, error) {
			elems, err := stack.PopR(2)
			if err != nil {
				return nil, err
			}
			return Floats{elems[0] / elems[1]}, nil
		},
	}

	leftShiftOp = Op{
		"<<",
		"left shift",
		func(stack *Stack) (Floats, error) {
			elems, err := stack.PopR(2)
			if err != nil {
				return nil, err
			}
			return Floats{elems[0] * math.Pow(2, elems[1])}, nil
		},
	}

	rightShiftOp = Op{
		">>",
		"right shift",
		func(stack *Stack) (Floats, error) {
			elems, err := stack.PopR(2)
			if err != nil {
				return nil, err
			}
			return Floats{elems[0] / math.Pow(2, elems[1])}, nil
		},
	}

	pow10Op = Op{
		"pow10",
		"10^stack.Top()",
		func(stack *Stack) (Floats, error) {
			top, err := stack.Pop()
			if err != nil {
				return nil, err
			}
			result := math.Pow10(int(top))
			return Floats{result}, nil
		},
	}

	randOp = Op{
		"r",
		fmt.Sprintf("random number from 0 to %d", _defaultMaxRand),
		func(stack *Stack) (Floats, error) {
			result, err := rand.Int(rand.Reader, big.NewInt(int64(_defaultMaxRand)))
			if err != nil {
				return nil, err
			}
			return Floats{float64(result.Int64())}, nil
		},
	}

	randNOp = Op{
		"rn",
		"random number from 0 to stack.Top()",
		func(stack *Stack) (Floats, error) {
			top, err := stack.Pop()
			if err != nil {
				return nil, err
			}
			result, err := rand.Int(rand.Reader, big.NewInt(int64(top)))
			if err != nil {
				return nil, err
			}
			return Floats{float64(result.Int64())}, nil
		},
	}

	signbitOp = Op{
		"signbit",
		"the sign bit of the number, 1 means negative, 0 positive",
		func(stack *Stack) (Floats, error) {
			top := stack.Top()
			if math.Signbit(top) {
				return Floats{1.0}, nil
			}
			return Floats{0.0}, nil
		},
	}

	swapOp = Op{
		"sw",
		"swap the top two elements",
		func(stack *Stack) (Floats, error) {
			err := stack.Swap()
			if err != nil {
				return nil, err
			}
			return nil, nil
		},
	}
)

func NewOpMap() OpMap {
	ops := OpMap{
		"!":       factorialOp,
		"*":       mulOp,
		"+":       plusOp,
		"++":      incrOp,
		"-":       minusOp,
		"--":      decrOp,
		"/":       divOp,
		"<<":      leftShiftOp,
		">>":      rightShiftOp,
		"abs":     wrapUnaryOp("abs", "absolute value", math.Abs),
		"pow10":   pow10Op,
		"r":       randOp,
		"rn":      randNOp,
		"signbit": signbitOp,
		"sw":      swapOp,
		"swa":     swapOp,
		"swap":    swapOp,
	}

	return ops
}

func NewOperatorMap() OperatorMap {
	operators := OperatorMap{
		"neg": func(stack *Stack) (float64, error) {
			top, err := stack.Pop()
			if err != nil {
				return 0.0, err
			}
			return -top, nil
		},
		"ilogb": func(stack *Stack) (float64, error) {
			top, err := stack.Pop()
			if err != nil {
				return 0.0, err
			}
			return float64(math.Ilogb(top)), nil
		},
		"isinf": func(stack *Stack) (float64, error) {
			top, err := stack.Pop()
			if err != nil {
				return 0.0, err
			}
			if math.IsInf(top, 1) {
				return 1.0, nil
			}
			return 0.0, nil
		},
		"isninf": func(stack *Stack) (float64, error) {
			top, err := stack.Pop()
			if err != nil {
				return 0.0, err
			}
			if math.IsInf(top, -1) {
				return 1.0, nil
			}
			return 0.0, nil
		},
		"isnan": func(stack *Stack) (float64, error) {
			top, err := stack.Pop()
			if err != nil {
				return 0.0, err
			}
			if math.IsNaN(top) {
				return 1.0, nil
			}
			return 0.0, nil
		},
		"jn": func(stack *Stack) (float64, error) {
			elems, err := stack.PopR(2)
			if err != nil {
				return 0.0, err
			}
			return math.Jn(int(elems[0]), elems[1]), nil
		},
		"yn": func(stack *Stack) (float64, error) {
			elems, err := stack.PopR(2)
			if err != nil {
				return 0.0, err
			}
			return math.Yn(int(elems[0]), elems[1]), nil
		},
		"mil": func(stack *Stack) (float64, error) {
			elems, err := stack.PopR(2)
			if err != nil {
				return 0.0, fmt.Errorf("distance_yds speed_mph")
			}
			distance_yds := elems[0]
			speed_mph := elems[1]
			speed_yps := speed_mph * 1760 / 3600
			mrads_per_s := 1000 * math.Atan(speed_yps/distance_yds)
			return mrads_per_s, nil
		},
		"mph": func(stack *Stack) (float64, error) {
			elems, err := stack.PopR(2)
			if err != nil {
				return 0.0, fmt.Errorf("distance_yds mrads_per_s")
			}
			distance_yds := elems[0]
			mrads_per_s := elems[1]
			rads_per_s := mrads_per_s / 1000.0
			displacement_per_s := math.Tan(rads_per_s)
			speed_yps := distance_yds * displacement_per_s
			speed_mph := speed_yps * 3600.0 / 1760.0

			return speed_mph, nil
		},
		"sum": func(stack *Stack) (float64, error) {
			result := 0.0
			for !stack.Empty() {
				n, err := stack.Pop()
				if err != nil {
					break
				}
				result += n
			}
			return result, nil
		},
		"avg": func(stack *Stack) (float64, error) {
			stats := welford.New()
			arr := stack.Copy()
			for _, n := range arr {
				stats.Add(n)
			}
			return stats.Mean(), nil
		},
		"sd": func(stack *Stack) (float64, error) {
			if stack.Len() <= 1 {
				return 0.0, nil
			}
			stats := welford.New()
			arr := stack.Copy()
			for _, n := range arr {
				stats.Add(n)
			}
			return stats.Stddev(), nil
		},
		"var": func(stack *Stack) (float64, error) {
			stats := welford.New()
			arr := stack.Copy()
			for _, n := range arr {
				stats.Add(n)
			}
			return stats.Variance(), nil
		},
		"max": func(stack *Stack) (float64, error) {
			stats := welford.New()
			arr := stack.Copy()
			for _, n := range arr {
				stats.Add(n)
			}
			return stats.Max(), nil
		},
		"min": func(stack *Stack) (float64, error) {
			stats := welford.New()
			arr := stack.Copy()
			for _, n := range arr {
				stats.Add(n)
			}
			return stats.Min(), nil
		},
		"lor": func(stack *Stack) (float64, error) {
			// lorentz
			top, err := stack.Pop()
			if err != nil {
				return 0.0, err
			}
			top_sq := top * top
			c := constants["c"]
			c_sq := c * c
			return 1.0 / math.Sqrt(1-top_sq/c_sq), nil
		},
		"fc": func(stack *Stack) (float64, error) {
			top, err := stack.Pop()
			if err != nil {
				return 0.0, err
			}
			return (top - 32) * 5 / 9, nil
		},
		"cf": func(stack *Stack) (float64, error) {
			top, err := stack.Pop()
			if err != nil {
				return 0.0, err
			}
			return (top * 9 / 5) + 32, nil
		},
		"fm": func(stack *Stack) (float64, error) {
			top, err := stack.Pop()
			if err != nil {
				return 0.0, err
			}
			return top / _ftPerM, nil
		},
		"mf": func(stack *Stack) (float64, error) {
			top, err := stack.Pop()
			if err != nil {
				return 0.0, err
			}
			return top * _ftPerM, nil
		},
		"fj": func(stack *Stack) (float64, error) {
			top, err := stack.Pop()
			if err != nil {
				return 0.0, err
			}
			return top * _jPerFtLb, nil
		},
		"jf": func(stack *Stack) (float64, error) {
			top, err := stack.Pop()
			if err != nil {
				return 0.0, err
			}
			return top / _jPerFtLb, nil
		},
		"gl": func(stack *Stack) (float64, error) {
			top, err := stack.Pop()
			if err != nil {
				return 0.0, err
			}
			return top * _lPerGal, nil
		},
		"lg": func(stack *Stack) (float64, error) {
			top, err := stack.Pop()
			if err != nil {
				return 0.0, err
			}
			return top / _lPerGal, nil
		},
		"pk": func(stack *Stack) (float64, error) {
			top, err := stack.Pop()
			if err != nil {
				return 0.0, err
			}
			return top * _pPerKg, nil
		},
		"kp": func(stack *Stack) (float64, error) {
			top, err := stack.Pop()
			if err != nil {
				return 0.0, err
			}
			return top / _pPerKg, nil
		},
		"hw": func(stack *Stack) (float64, error) {
			top, err := stack.Pop()
			if err != nil {
				return 0.0, err
			}
			return top * _wPerHp, nil
		},
		"wh": func(stack *Stack) (float64, error) {
			top, err := stack.Pop()
			if err != nil {
				return 0.0, err
			}
			return top / _wPerHp, nil
		},
		"pas": func(stack *Stack) (float64, error) {
			// pasteurization time in seconds for a given temperature in fahrenheit
			// derived from a curve fit of data
			top, err := stack.Pop()
			if err != nil {
				return 0.0, err
			}
			return math.Exp(top*-0.231) * 1.23e15 * 60, nil
		},
		"pr": func(stack *Stack) (float64, error) {
			// base atmospheric pressure in inHg for a given elevation in feet
			// 29.9212524*pow(1-pow(10, -5)*2.25577*(x/3.280839895), 5.25588)
			// from https://www.engineeringtoolbox.com/air-altitude-pressure-d_462.html
			top, err := stack.Pop()
			if err != nil {
				return 0.0, err
			}
			return 29.9212524 * math.Pow(1-math.Pow(10, -5)*2.25577*(top/3.280839895), 5.25588), nil
		},
	}

	installUnaryFunctions(operators)
	installBinaryFunctions(operators)
	installTernaryFunctions(operators)
	installConstants(operators)

	return operators
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
	"mod":       math.Mod,
	"nextafter": math.Nextafter,
	"pow":       math.Pow,
	"remainder": math.Remainder,
}

var ternaryFunctions = map[string]func(float64, float64, float64) float64{
	"fma": math.FMA,
}

func wrapUnaryFunction(unary func(float64) float64) Operator {
	return func(stack *Stack) (float64, error) {
		top, err := stack.Pop()
		if err != nil {
			return 0.0, err
		}
		return unary(top), nil
	}
}

func installUnaryFunctions(operators OperatorMap) {
	for name, uFunc := range unaryFunctions {
		operators[name] = wrapUnaryFunction(uFunc)
	}
}

func installBinaryFunctions(operators OperatorMap) {
	for name, bFunc := range binaryFunctions {
		operators[name] = func(stack *Stack) (float64, error) {
			elems, err := stack.PopR(2)
			if err != nil {
				return 0.0, err
			}
			return bFunc(elems[0], elems[1]), nil
		}
	}
}

func installTernaryFunctions(operators OperatorMap) {
	for name, tFunc := range ternaryFunctions {
		operators[name] = func(stack *Stack) (float64, error) {
			elems, err := stack.PopR(3)
			if err != nil {
				return 0.0, err
			}
			return tFunc(elems[0], elems[1], elems[2]), nil
		}
	}
}

func tryOp(line string, stack *Stack, ops OpMap) error {
	rawOp, ok := ops[line]
	if ok {
		results, err := rawOp.f(stack)
		if err != nil {
			return err
		}
		for _, result := range results {
			stack.Push(result)
		}
		return nil
	}
	return fmt.Errorf("no op '%s'", line)
}

func tryOperator(line string, stack *Stack, operators OperatorMap) error {
	rawOperator, ok := operators[line]
	if ok {
		result, err := rawOperator(stack)
		if err != nil {
			fmt.Println(err)
		} else {
			stack.Push(result)
		}
		return nil
	}
	return fmt.Errorf("no operator '%s'", line)
}

func showHelp(operators OperatorMap) {
	commands := make([]string, 0, len(operators))
	for key := range operators {
		commands = append(commands, key)
	}
	slices.Sort(commands)
	for _, command := range commands {
		fmt.Println(command)
	}
}
