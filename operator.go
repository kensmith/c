package main

import (
	"crypto/rand"
	"fmt"
	"math"
	"math/big"
	"slices"
	"strings"

	"github.com/eclesh/welford"
)

// TODO add remaining math pkg funcs

/*
* Operator needs to remember:
* - the stack
* - the string which invokes it
* - the function to invoke
*   - this should return ([]float64, error)
* - help text describing it for help display
 */

type (
	Operator    func(*Stack) (float64, error)
	OperatorMap map[string]Operator

	Floats []float64

	OperatorFunc func(*Stack) (Floats, error)

	Op struct {
		name string
		doc  string
		f    OperatorFunc
	}

	OpMap map[string]Op

	Ops struct {
		opmap OpMap
	}
)

func (o *Ops) Run(line string, stack *Stack) error {
	op, ok := o.opmap[line]
	if ok {
		results, err := op.f(stack)
		if err != nil {
			return err
		}
		for _, result := range results {
			stack.Push(result)
		}
		return nil
	}
	return fmt.Errorf("no operator '%s'", line)
}

func NewOps() *Ops {
	ops := Ops{
		opmap: OpMap{
			"!":           factorialOp,
			"%":           wrapBinaryOp("%", "floating-point remainder of x/y", math.Mod),
			"*":           mulOp,
			"**":          wrapBinaryOp("**", "x^y, the base-x exponential of y", math.Pow),
			"+":           plusOp,
			"++":          incrOp,
			"-":           minusOp,
			"--":          decrOp,
			"/":           divOp,
			"<<":          leftShiftOp,
			">>":          rightShiftOp,
			"^":           wrapBinaryOp("^", "x^y, the base-x exponential of y", math.Pow),
			"abs":         wrapUnaryOp("abs", "absolute value", math.Abs),
			"acos":        wrapUnaryOp("acos", "arccosine, in radians", math.Acos),
			"acosh":       wrapUnaryOp("acosh", "inverse hyperbolic cosine", math.Acosh),
			"asin":        wrapUnaryOp("asin", "arcsine", math.Asin),
			"asinh":       wrapUnaryOp("asinh", "inverse hyperbolic sine ", math.Asinh),
			"atan":        wrapUnaryOp("atan", "arctangent", math.Atan),
			"atan2":       wrapBinaryOp("atan2", "tangent of y/x", math.Atan2),
			"avg":         avgOp,
			"c":           wrapConstant("c", "speed of light in m/s", 299792458),
			"cbrt":        wrapUnaryOp("cbrt", "cube root", math.Cbrt),
			"ceil":        wrapUnaryOp("ceil", "least integer value greater than or equal to stack.Top()", math.Ceil),
			"cf":          cfOp,
			"cos":         wrapUnaryOp("cos", "cosine", math.Cos),
			"cosh":        wrapUnaryOp("cosh", "hyperbolic cosine", math.Cosh),
			"dim":         wrapBinaryOp("dim", "maximum of x-y or 0", math.Dim),
			"e":           wrapConstant("e", "euler's constant", math.E),
			"erf":         wrapUnaryOp("erf", "error function", math.Erf),
			"erfc":        wrapUnaryOp("erfc", "complementary error function", math.Erfc),
			"erfcinv":     wrapUnaryOp("erfcinv", "inverse of erfc", math.Erfcinv),
			"erfinv":      wrapUnaryOp("erfinv", "inverse error function", math.Erfinv),
			"exp":         wrapUnaryOp("exp", "e^x, the base-e exponential", math.Exp),
			"exp2":        wrapUnaryOp("exp2", "2^x, the base-2 exponential", math.Exp2),
			"expm1":       wrapUnaryOp("expm1", "e^x - 1, the base-e exponential of x minus 1. It is more accurate than exp - 1 when x is near zero", math.Expm1),
			"fc":          fcOp,
			"fj":          fjOp,
			"floor":       wrapUnaryOp("floor", "greatest integer value less than or equal to stack.Top()", math.Floor),
			"fm":          fmOp,
			"fma":         wrapTernaryOp("fma", "fused multiply-add of x, y, and z", math.FMA),
			"gamma":       wrapUnaryOp("gamma", "gamma function ", math.Gamma),
			"gl":          glOp,
			"hw":          hwOp,
			"hypot":       wrapBinaryOp("hypot", "sqrt(p*p + q*q), taking care to avoid unnecessary overflow and underflow", math.Hypot),
			"ilogb":       ilogbOp,
			"inf":         wrapConstant("inf", "positive infinity", math.Inf(1)),
			"isinf":       isInfOp,
			"isnan":       isNanOp,
			"isninf":      isNInfOp,
			"j0":          wrapUnaryOp("j0", "order-zero Bessel function of the first kind", math.J0),
			"j1":          wrapUnaryOp("j1", "order-one Bessel function of the first kind", math.J1),
			"jf":          jfOp,
			"jn":          jnOp,
			"kp":          kpOp,
			"lg":          lgOp,
			"ln2":         wrapConstant("ln2", "natural log of 2", math.Ln2),
			"ln10":        wrapConstant("ln10", "natural log of 10", math.Ln10),
			"log2e":       wrapConstant("log2e", "1 / ln2", math.Log2E),
			"log10e":      wrapConstant("log10e", "1 / ln10", math.Log10E),
			"log":         wrapUnaryOp("log", "natural logarithm", math.Log),
			"log10":       wrapUnaryOp("log10", "decimal logarithm", math.Log10),
			"log1p":       wrapUnaryOp("log1p", "natural logarithm of 1 plus its argument x. It is more accurate than log(1 + x) when x is near zero", math.Log1p),
			"log2":        wrapUnaryOp("log2", "binary logarithm", math.Log2),
			"logb":        wrapUnaryOp("logb", "binary exponent", math.Logb),
			"lor":         lorOp,
			"max":         maxOp,
			"mf":          mfOp,
			"mil":         milOp,
			"min":         minOp,
			"mod":         wrapBinaryOp("mod", "floating-point remainder of x/y", math.Mod),
			"mph":         mphOp,
			"nan":         wrapConstant("nan", "not a number", math.NaN()),
			"neg":         negOp,
			"nextafter":   wrapBinaryOp("nextafter", "next representable float64 value after x towards y", math.Nextafter),
			"ninf":        wrapConstant("ninf", "negative infinity", math.Inf(-1)),
			"noop":        noOp,
			"p":           pOp,
			"pas":         pasOp,
			"phi":         wrapConstant("phi", "golden ratio", math.Phi),
			"pi":          wrapConstant("pi", "ratio of a circle's circumference to its diameter", math.Pi),
			"pk":          pkOp,
			"pop":         pOp,
			"pow":         wrapBinaryOp("pow", "x^y, the base-x exponential of y", math.Pow),
			"pow10":       pow10Op,
			"pr":          prOp,
			"r":           randOp,
			"remainder":   wrapBinaryOp("remainder", "IEEE 754 floating-point remainder of x/y", math.Remainder),
			"rn":          randNOp,
			"round":       wrapUnaryOp("round", "returns the nearest integer, rounding half away from zero", math.Round),
			"roundtoeven": wrapUnaryOp("roundtoeven", "returns the nearest integer, rounding ties to even", math.RoundToEven),
			"sd":          sdOp,
			"signbit":     signbitOp,
			"sin":         wrapUnaryOp("sin", "sine", math.Sin),
			"sinh":        wrapUnaryOp("sinh", "hyperbolic sine", math.Sinh),
			"sort":        sortOp,
			"sqrt":        wrapUnaryOp("sqrt", "square root", math.Sqrt),
			"sqrt2":       wrapConstant("sqrt2", "square root of 2", math.Sqrt2),
			"sqrte":       wrapConstant("sqrte", "square root of e", math.SqrtE),
			"sqrtphi":     wrapConstant("sqrtphi", "square root of the golden ratio", math.SqrtPhi),
			"sqrtpi":      wrapConstant("sqrtpi", "square root of pi", math.SqrtPi),
			"sum":         sumOp,
			"sw":          swapOp,
			"swa":         swapOp,
			"swap":        swapOp,
			"tan":         wrapUnaryOp("tan", "tangent", math.Tan),
			"tanh":        wrapUnaryOp("tanh", "hyperbolic tangent", math.Tanh),
			"trunc":       wrapUnaryOp("trunc", "integer value of stack.Top()", math.Trunc),
			"var":         varOp,
			"wh":          whOp,
			"y0":          wrapUnaryOp("y0", "order-zero Bessel function of the second kind", math.Y0),
			"y1":          wrapUnaryOp("y1", "order-one Bessel function of the second kind", math.Y1),
			"yn":          ynOp,
		},
	}
	return &ops
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

	avgOp = Op{
		"avg",
		"average (mean) of the entire stack",
		func(stack *Stack) (Floats, error) {
			stats := welford.New()
			arr := stack.Copy()
			for _, n := range arr {
				stats.Add(n)
			}
			return Floats{stats.Mean()}, nil
		},
	}

	cfOp = Op{
		"cf",
		"celcius to fahrenheit conversion",
		func(stack *Stack) (Floats, error) {
			top, err := stack.Pop()
			if err != nil {
				return nil, err
			}
			return Floats{top*9/5 + 32}, nil
		},
	}

	fcOp = Op{
		"fc",
		"fahrenheit to celcius conversion",
		func(stack *Stack) (Floats, error) {
			top, err := stack.Pop()
			if err != nil {
				return nil, err
			}
			return Floats{(top - 32) * 5 / 9}, nil
		},
	}

	fjOp = Op{
		"fj",
		"foot-lbs to joules conversion",
		func(stack *Stack) (Floats, error) {
			top, err := stack.Pop()
			if err != nil {
				return nil, err
			}
			return Floats{top * _jPerFtLb}, nil
		},
	}

	fmOp = Op{
		"fm",
		"feet to meters conversion",
		func(stack *Stack) (Floats, error) {
			top, err := stack.Pop()
			if err != nil {
				return nil, err
			}
			return Floats{top / _ftPerM}, nil
		},
	}

	glOp = Op{
		"gl",
		"gallons to liters conversion",
		func(stack *Stack) (Floats, error) {
			top, err := stack.Pop()
			if err != nil {
				return nil, err
			}
			return Floats{top * _lPerGal}, nil
		},
	}

	hwOp = Op{
		"hw",
		"horsepower to watts conversion",
		func(stack *Stack) (Floats, error) {
			top, err := stack.Pop()
			if err != nil {
				return nil, err
			}
			return Floats{top * _wPerHp}, nil
		},
	}

	ilogbOp = Op{
		"ilogb",
		"binary exponent of stack.Top()",
		func(stack *Stack) (Floats, error) {
			top, err := stack.Pop()
			if err != nil {
				return nil, err
			}
			result := float64(math.Ilogb(top))
			return Floats{result}, nil
		},
	}

	isInfOp = Op{
		"isinf",
		"1 if stack.Top() is +Inf",
		func(stack *Stack) (Floats, error) {
			top := stack.Top()
			if math.IsInf(top, 1) {
				return Floats{1}, nil
			}
			return Floats{0}, nil
		},
	}

	isNanOp = Op{
		"isnan",
		"1 if stack.Top() is NaN",
		func(stack *Stack) (Floats, error) {
			top := stack.Top()
			if math.IsNaN(top) {
				return Floats{1}, nil
			}
			return Floats{0}, nil
		},
	}

	isNInfOp = Op{
		"isninf",
		"1 if stack.Top() is -Inf",
		func(stack *Stack) (Floats, error) {
			top := stack.Top()
			if math.IsInf(top, -1) {
				return Floats{1}, nil
			}
			return Floats{0}, nil
		},
	}

	jfOp = Op{
		"jf",
		"joules to foot-lbs conversion",
		func(stack *Stack) (Floats, error) {
			top, err := stack.Pop()
			if err != nil {
				return nil, err
			}
			return Floats{top / _jPerFtLb}, nil
		},
	}

	jnOp = Op{
		"jn",
		"order-n Bessel function of the first kind",
		func(stack *Stack) (Floats, error) {
			elems, err := stack.PopN(2)
			if err != nil {
				return nil, err
			}
			result := math.Jn(int(elems[0]), elems[1])
			return Floats{result}, nil
		},
	}

	kpOp = Op{
		"kp",
		"kilograms to pounds conversion",
		func(stack *Stack) (Floats, error) {
			top, err := stack.Pop()
			if err != nil {
				return nil, err
			}
			return Floats{top / _pPerKg}, nil
		},
	}

	lgOp = Op{
		"lg",
		"liters to gallons conversion",
		func(stack *Stack) (Floats, error) {
			top, err := stack.Pop()
			if err != nil {
				return nil, err
			}
			return Floats{top / _lPerGal}, nil
		},
	}

	lorOp = Op{
		"lor",
		"lorentz factor",
		func(stack *Stack) (Floats, error) {
			top, err := stack.Pop()
			if err != nil {
				return nil, err
			}
			topSq := top * top
			cSq := 89875517873681764.0
			return Floats{1.0 / math.Sqrt(1-topSq/cSq)}, nil
		},
	}

	maxOp = Op{
		"max",
		"find the maximum value of the entire stack",
		func(stack *Stack) (Floats, error) {
			stats := welford.New()
			arr := stack.Copy()
			for _, n := range arr {
				stats.Add(n)
			}
			return Floats{stats.Max()}, nil
		},
	}

	mfOp = Op{
		"mf",
		"meters to feet conversion",
		func(stack *Stack) (Floats, error) {
			top, err := stack.Pop()
			if err != nil {
				return nil, err
			}
			return Floats{top * _ftPerM}, nil
		},
	}

	milOp = Op{
		"mil",
		"given yards to target and target speed in mph, return target speed in millradians per second (multiply result by time of flight for lead)",
		func(stack *Stack) (Floats, error) {
			elems, err := stack.PopR(2)
			if err != nil {
				return nil, err
			}
			distance_yds := elems[0]
			speed_mph := elems[1]
			speed_yps := speed_mph * 1760.0 / 3600.0
			mrads_per_s := 1000.0 * math.Atan(speed_yps/distance_yds)
			return Floats{mrads_per_s}, nil
		},
	}

	minOp = Op{
		"min",
		"find the minimum value of the entire stack",
		func(stack *Stack) (Floats, error) {
			stats := welford.New()
			arr := stack.Copy()
			for _, n := range arr {
				stats.Add(n)
			}
			return Floats{stats.Min()}, nil
		},
	}

	mphOp = Op{
		"mph",
		"given yards to target and target speed in millradians per second, return target speed in mph",
		func(stack *Stack) (Floats, error) {
			elems, err := stack.PopR(2)
			if err != nil {
				return nil, err
			}
			distance_yds := elems[0]
			mrads_per_s := elems[1]
			rads_per_s := mrads_per_s / 1000.0
			displacement_per_s := math.Tan(rads_per_s)
			speed_yps := distance_yds * displacement_per_s
			speed_mph := speed_yps * 3600.0 / 1760.0

			return Floats{speed_mph}, nil
		},
	}

	noOp = Op{
		"noop",
		"no op",
		func(stack *Stack) (Floats, error) {
			return nil, nil
		},
	}

	pOp = Op{
		"p",
		"pop an item from the stack",
		func(stack *Stack) (Floats, error) {
			_, _ = stack.Pop()
			return nil, nil
		},
	}

	pasOp = Op{
		"pas",
		"pasteurization time in seconds for a given temperature in fahrenheit",
		func(stack *Stack) (Floats, error) {
			top, err := stack.Pop()
			if err != nil {
				return nil, err
			}
			return Floats{math.Exp(top*-0.231) * 1.23e15 * 60}, nil
		},
	}

	pkOp = Op{
		"pk",
		"pounds to kilograms conversion",
		func(stack *Stack) (Floats, error) {
			top, err := stack.Pop()
			if err != nil {
				return nil, err
			}
			return Floats{top * _pPerKg}, nil
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

	prOp = Op{
		"pr",
		"pressure in inHg for a given altitude in feet",
		func(stack *Stack) (Floats, error) {
			top, err := stack.Pop()
			if err != nil {
				return nil, err
			}
			return Floats{29.9212524 * math.Pow(1-math.Pow(10, -5)*2.25577*(top/3.280839895), 5.25588)}, nil
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

	sdOp = Op{
		"sd",
		"standard deviation of the entire stack",
		func(stack *Stack) (Floats, error) {
			if stack.Len() <= 1 {
				return Floats{0}, nil
			}
			stats := welford.New()
			arr := stack.Copy()
			for _, n := range arr {
				stats.Add(n)
			}
			return Floats{stats.Stddev()}, nil
		},
	}

	sortOp = Op{
		"sort",
		"sort the entire stack",
		func(stack *Stack) (Floats, error) {
			stack.Sort()
			return nil, nil
		},
	}

	sumOp = Op{
		"sum",
		"sum the entire stack",
		func(stack *Stack) (Floats, error) {
			arr := stack.Copy()
			result := 0.0
			for _, n := range arr {
				result += n
			}
			return Floats{result}, nil
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

	negOp = Op{
		"neg",
		"negate stack.Top()",
		func(stack *Stack) (Floats, error) {
			top, err := stack.Pop()
			if err != nil {
				return nil, err
			}
			return Floats{-top}, nil
		},
	}

	varOp = Op{
		"var",
		"variance of the entire stack",
		func(stack *Stack) (Floats, error) {
			stats := welford.New()
			arr := stack.Copy()
			for _, n := range arr {
				stats.Add(n)
			}
			return Floats{stats.Variance()}, nil
		},
	}

	whOp = Op{
		"wh",
		"watts to horsepower conversion",
		func(stack *Stack) (Floats, error) {
			top, err := stack.Pop()
			if err != nil {
				return nil, err
			}
			return Floats{top / _wPerHp}, nil
		},
	}

	ynOp = Op{
		"yn",
		"order-n Bessel function of the second kind",
		func(stack *Stack) (Floats, error) {
			elems, err := stack.PopN(2)
			if err != nil {
				return nil, err
			}
			return Floats{math.Yn(int(elems[0]), elems[1])}, nil
		},
	}
)

func wrapConstant(name, doc string, value float64) Op {
	return Op{
		name,
		doc,
		func(stack *Stack) (Floats, error) {
			return Floats{value}, nil
		},
	}
}

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

func wrapBinaryOp(name, doc string, f func(float64, float64) float64) Op {
	return Op{
		name,
		doc,
		func(stack *Stack) (Floats, error) {
			elems, err := stack.PopN(2)
			if err != nil {
				return nil, err
			}
			result := f(elems[0], elems[1])
			return Floats{result}, nil
		},
	}
}

func wrapTernaryOp(name, doc string, f func(float64, float64, float64) float64) Op {
	return Op{
		name,
		doc,
		func(stack *Stack) (Floats, error) {
			elems, err := stack.PopN(3)
			if err != nil {
				return nil, err
			}
			result := f(elems[0], elems[1], elems[2])
			return Floats{result}, nil
		},
	}
}

func (o *Ops) OpNames() ([]string, int) {
	names := make([]string, 0, len(o.opmap))
	longest := math.MinInt
	for name := range o.opmap {
		names = append(names, name)
		longest = max(longest, len(name))
	}
	slices.Sort(names)
	return names, longest
}

func (o *Ops) Help() string {
	var b strings.Builder
	names, longest := o.OpNames()
	for _, name := range names {
		verbs := fmt.Sprintf("%%-%ds - %%s\n", longest)
		fmt.Fprintf(&b, verbs, name, o.opmap[name].doc)
	}
	return b.String()
}
