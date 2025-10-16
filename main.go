package main

import (
	"crypto/rand"
	"fmt"
	"math"
	"math/big"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/adrg/xdg"
	"github.com/chzyer/readline"
	"github.com/expr-lang/expr"
)

const (
	_ftPerM         = 3.280839895
	_jPerFtLb       = 1.3558179483314004
	_lPerGal        = 3.785411784
	_pPerKg         = 0.45359237
	_wPerHp         = 745.699872
	_defaultMaxRand = math.MaxInt16
)

var (
	_histDirname  = filepath.Join(xdg.StateHome, "github.com", "kensmith", "c")
	_histFilename = filepath.Join(_histDirname, "history")
)

type Operator func(*Stack) (float64, error)

var operators = map[string]Operator{
	"+": func(stack *Stack) (float64, error) {
		elems, err := stack.PopR(2)
		if err != nil {
			return 0.0, err
		}
		return elems[0] + elems[1], nil
	},
	"-": func(stack *Stack) (float64, error) {
		elems, err := stack.PopR(2)
		if err != nil {
			return 0.0, err
		}
		return elems[0] - elems[1], nil
	},
	"*": func(stack *Stack) (float64, error) {
		elems, err := stack.PopR(2)
		if err != nil {
			return 0.0, err
		}
		return elems[0] * elems[1], nil
	},
	"/": func(stack *Stack) (float64, error) {
		elems, err := stack.PopR(2)
		if err != nil {
			return 0.0, err
		}
		return elems[0] / elems[1], nil
	},
	"<<": func(stack *Stack) (float64, error) {
		elems, err := stack.PopR(2)
		if err != nil {
			return 0.0, err
		}
		return elems[0] * math.Pow(2, elems[1]), nil
	},
	">>": func(stack *Stack) (float64, error) {
		elems, err := stack.PopR(2)
		if err != nil {
			return 0.0, err
		}
		return elems[0] / math.Pow(2, elems[1]), nil
	},
	"!": func(stack *Stack) (float64, error) {
		top, err := stack.Pop()
		if err != nil {
			return 0.0, err
		}
		return top * math.Gamma(top), nil
	},
	"++": func(stack *Stack) (float64, error) {
		top, err := stack.Pop()
		if err != nil {
			return 0.0, err
		}
		return top + 1, nil
	},
	"--": func(stack *Stack) (float64, error) {
		top, err := stack.Pop()
		if err != nil {
			return 0.0, err
		}
		return top - 1, nil
	},
	"rn": func(stack *Stack) (float64, error) {
		top, err := stack.Pop()
		if err != nil {
			return 0.0, err
		}
		result, err := rand.Int(rand.Reader, big.NewInt(int64(top)))
		if err != nil {
			return 0.0, err
		}
		return float64(result.Int64()), nil
	},
	"r": func(stack *Stack) (float64, error) {
		result, err := rand.Int(rand.Reader, big.NewInt(int64(_defaultMaxRand)))
		if err != nil {
			return 0.0, err
		}
		return float64(result.Int64()), nil
	},
	"pow10": func(stack *Stack) (float64, error) {
		top, err := stack.Pop()
		if err != nil {
			return 0.0, err
		}
		return math.Pow10(int(top)), nil
	},
	"signbit": func(stack *Stack) (float64, error) {
		top, err := stack.Pop()
		if err != nil {
			return 0.0, err
		}
		if math.Signbit(top) {
			return 1, nil
		}

		return 0, nil
	},
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
		result := 0.0
		size := stack.Len()
		for !stack.Empty() {
			n, err := stack.Pop()
			if err != nil {
				break
			}
			result += n
		}
		return result / float64(size), nil
	},
	"max": func(stack *Stack) (float64, error) {
		return stack.Max(), nil
	},
	"min": func(stack *Stack) (float64, error) {
		return stack.Min(), nil
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

var constants = map[string]float64{
	"c":       299792458,
	"e":       math.E,
	"inf":     math.Inf(1),
	"ln10":    math.Ln10,
	"ln2":     math.Ln2,
	"log20e":  math.Log10E,
	"log2e":   math.Log2E,
	"nan":     math.NaN(),
	"ninf":    math.Inf(-1),
	"phi":     math.Phi,
	"pi":      math.Pi,
	"sqrt2":   math.Sqrt2,
	"sqrte":   math.SqrtE,
	"sqrtphi": math.SqrtPhi,
	"sqrtpi":  math.SqrtPi,
}

func main() {
	for name, c := range constants {
		operators[name] = func(stack *Stack) (float64, error) {
			return c, nil
		}
	}

	for name, uFunc := range unaryFunctions {
		operators[name] = func(stack *Stack) (float64, error) {
			top, err := stack.Pop()
			if err != nil {
				return 0.0, err
			}
			return uFunc(top), nil
		}
	}

	for name, bFunc := range binaryFunctions {
		operators[name] = func(stack *Stack) (float64, error) {
			elems, err := stack.PopR(2)
			if err != nil {
				return 0.0, err
			}
			return bFunc(elems[0], elems[1]), nil
		}
	}

	for name, tFunc := range ternaryFunctions {
		operators[name] = func(stack *Stack) (float64, error) {
			elems, err := stack.PopR(3)
			if err != nil {
				return 0.0, err
			}
			return tFunc(elems[0], elems[1], elems[2]), nil
		}
	}

	var shell *readline.Instance
	err := os.MkdirAll(_histDirname, 0o750)
	if err != nil {
		fmt.Printf("history disabled due to inability to create directory: %s", _histDirname)
		shell, err = readline.New("[  ]> ")
		if err != nil {
			panic(err)
		}
	} else {
		shell, err = readline.NewEx(&readline.Config{
			Prompt:      "[  ]> ",
			HistoryFile: _histFilename,
		})
	}
	if err != nil {
		panic(err)
	}
	defer func() {
		err := shell.Close()
		if err != nil {
			fmt.Println(err)
		}
	}()

	stack := Stack{}
	env := map[string]any{
		"stack": &stack,
		"s":     &stack,
	}

	for {
		line, err := shell.Readline()
		if err != nil {
			// normal exit due to ctrl-c, ctrl-d
			return
		}
		commasRemoved := strings.ReplaceAll(line, ",", "")
		lineTrimmed := strings.TrimSpace(commasRemoved)
		rawOperator, ok := operators[lineTrimmed]
		if ok {
			result, err := rawOperator(&stack)
			if err != nil {
				fmt.Println(err)
			} else {
				stack.Push(result)
			}
		} else {
			switch lineTrimmed {
			case "sw":
				fallthrough
			case "swap":
				err := stack.Swap()
				if err != nil {
					fmt.Println(err)
					continue
				}
			case "p":
				fallthrough
			case "pop":
				_, err := stack.Pop()
				if err != nil {
					fmt.Println(err)
					continue
				}
			case "sort":
				stack.Sort()
			case "cl":
				fallthrough
			case "clr":
				fallthrough
			case "clear":
				stack.Clear()
			case "q":
				return
			default:
				output, err := expr.Eval(lineTrimmed, env)
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
				stack.Push(value)
			}
		}
		shell.SetPrompt(stack.String() + "> ")
	}
}
