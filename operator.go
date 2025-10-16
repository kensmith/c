package main

import (
	"crypto/rand"
	"fmt"
	"math"
	"math/big"
	"os"
	"slices"
	"strings"

	"github.com/eclesh/welford"
)

// TODO add remaining math pkg funcs

type (
	Floats []float64

	OperatorFunc func(*Stack) (Floats, error)

	Op struct {
		doc string
		f   OperatorFunc
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
			"%":           wrapBinaryOp("floating-point remainder of x/y", math.Mod),
			"*":           mulOp,
			"**":          wrapBinaryOp("x^y, the base-x exponential of y", math.Pow),
			"+":           plusOp,
			"++":          incrOp,
			"-":           minusOp,
			"--":          decrOp,
			"/":           divOp,
			"<<":          leftShiftOp,
			">>":          rightShiftOp,
			"^":           wrapBinaryOp("x^y, the base-x exponential of y", math.Pow),
			"abs":         wrapUnaryOp("absolute value", math.Abs),
			"acos":        wrapUnaryOp("arccosine, in radians", math.Acos),
			"acosh":       wrapUnaryOp("inverse hyperbolic cosine", math.Acosh),
			"asin":        wrapUnaryOp("arcsine", math.Asin),
			"asinh":       wrapUnaryOp("inverse hyperbolic sine ", math.Asinh),
			"atan":        wrapUnaryOp("arctangent", math.Atan),
			"atan2":       wrapBinaryOp("tangent of y/x", math.Atan2),
			"avg":         avgOp,
			"c":           wrapConstant("speed of light in m/s", 299792458),
			"cl":          clearOp,
			"clr":         clearOp,
			"clear":       clearOp,
			"cbrt":        wrapUnaryOp("cube root", math.Cbrt),
			"ceil":        wrapUnaryOp("least integer value greater than or equal to stack.Top()", math.Ceil),
			"cf":          cfOp,
			"cos":         wrapUnaryOp("cosine", math.Cos),
			"cosh":        wrapUnaryOp("hyperbolic cosine", math.Cosh),
			"dim":         wrapBinaryOp("maximum of x-y or 0", math.Dim),
			"e":           wrapConstant("euler's constant", math.E),
			"erf":         wrapUnaryOp("error function", math.Erf),
			"erfc":        wrapUnaryOp("complementary error function", math.Erfc),
			"erfcinv":     wrapUnaryOp("inverse of erfc", math.Erfcinv),
			"erfinv":      wrapUnaryOp("inverse error function", math.Erfinv),
			"exit":        qOp,
			"exp":         wrapUnaryOp("e^x, the base-e exponential", math.Exp),
			"exp2":        wrapUnaryOp("2^x, the base-2 exponential", math.Exp2),
			"expm1":       wrapUnaryOp("e^x - 1, the base-e exponential of x minus 1. It is more accurate than exp - 1 when x is near zero", math.Expm1),
			"f":           fOp,
			"fc":          fcOp,
			"fj":          fjOp,
			"floor":       wrapUnaryOp("greatest integer value less than or equal to stack.Top()", math.Floor),
			"fm":          fmOp,
			"fma":         wrapTernaryOp("fused multiply-add of x, y, and z", math.FMA),
			"frexp":       frexpOp,
			"gamma":       wrapUnaryOp("gamma function ", math.Gamma),
			"gl":          glOp,
			"hw":          hwOp,
			"hypot":       wrapBinaryOp("sqrt(p*p + q*q), taking care to avoid unnecessary overflow and underflow", math.Hypot),
			"ilogb":       ilogbOp,
			"inf":         wrapConstant("positive infinity", math.Inf(1)),
			"isinf":       isInfOp,
			"isnan":       isNanOp,
			"isninf":      isNInfOp,
			"j0":          wrapUnaryOp("order-zero Bessel function of the first kind", math.J0),
			"j1":          wrapUnaryOp("order-one Bessel function of the first kind", math.J1),
			"jf":          jfOp,
			"jn":          jnOp,
			"kp":          kpOp,
			"lg":          lgOp,
			"ln2":         wrapConstant("natural log of 2", math.Ln2),
			"ln10":        wrapConstant("natural log of 10", math.Ln10),
			"log2e":       wrapConstant("1 / ln2", math.Log2E),
			"log10e":      wrapConstant("1 / ln10", math.Log10E),
			"log":         wrapUnaryOp("natural logarithm", math.Log),
			"log10":       wrapUnaryOp("decimal logarithm", math.Log10),
			"log1p":       wrapUnaryOp("natural logarithm of 1 plus its argument x. It is more accurate than log(1 + x) when x is near zero", math.Log1p),
			"log2":        wrapUnaryOp("binary logarithm", math.Log2),
			"logb":        wrapUnaryOp("binary exponent", math.Logb),
			"lor":         lorOp,
			"max":         maxOp,
			"mf":          mfOp,
			"mil":         milOp,
			"min":         minOp,
			"mod":         wrapBinaryOp("floating-point remainder of x/y", math.Mod),
			"mph":         mphOp,
			"nan":         wrapConstant("not a number", math.NaN()),
			"neg":         negOp,
			"nextafter":   wrapBinaryOp("next representable float64 value after x towards y", math.Nextafter),
			"ninf":        wrapConstant("negative infinity", math.Inf(-1)),
			"noop":        noOp,
			"p":           pOp,
			"pas":         pasOp,
			"phi":         wrapConstant("golden ratio", math.Phi),
			"pi":          wrapConstant("ratio of a circle's circumference to its diameter", math.Pi),
			"pk":          pkOp,
			"pop":         pOp,
			"pow":         wrapBinaryOp("x^y, the base-x exponential of y", math.Pow),
			"pow10":       pow10Op,
			"pr":          prOp,
			"q":           qOp,
			"r":           randOp,
			"remainder":   wrapBinaryOp("IEEE 754 floating-point remainder of x/y", math.Remainder),
			"rn":          randNOp,
			"round":       wrapUnaryOp("returns the nearest integer, rounding half away from zero", math.Round),
			"roundtoeven": wrapUnaryOp("returns the nearest integer, rounding ties to even", math.RoundToEven),
			"sd":          sdOp,
			"signbit":     signbitOp,
			"sin":         wrapUnaryOp("sine", math.Sin),
			"sinh":        wrapUnaryOp("hyperbolic sine", math.Sinh),
			"sort":        sortOp,
			"sqrt":        wrapUnaryOp("square root", math.Sqrt),
			"sqrt2":       wrapConstant("square root of 2", math.Sqrt2),
			"sqrte":       wrapConstant("square root of e", math.SqrtE),
			"sqrtphi":     wrapConstant("square root of the golden ratio", math.SqrtPhi),
			"sqrtpi":      wrapConstant("square root of pi", math.SqrtPi),
			"sum":         sumOp,
			"sw":          swapOp,
			"swa":         swapOp,
			"swap":        swapOp,
			"tan":         wrapUnaryOp("tangent", math.Tan),
			"tanh":        wrapUnaryOp("hyperbolic tangent", math.Tanh),
			"trunc":       wrapUnaryOp("integer value of stack.Top()", math.Trunc),
			"var":         varOp,
			"wh":          whOp,
			"y0":          wrapUnaryOp("order-zero Bessel function of the second kind", math.Y0),
			"y1":          wrapUnaryOp("order-one Bessel function of the second kind", math.Y1),
			"yn":          ynOp,
		},
	}
	return &ops
}

var (
	factorialOp = Op{
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
		"celcius to fahrenheit conversion",
		func(stack *Stack) (Floats, error) {
			top, err := stack.Pop()
			if err != nil {
				return nil, err
			}
			return Floats{top*9/5 + 32}, nil
		},
	}

	clearOp = Op{
		"remove everything from the stack",
		func(stack *Stack) (Floats, error) {
			stack.Clear()
			return nil, nil
		},
	}

	fOp = Op{
		"print stack using %f",
		func(stack *Stack) (Floats, error) {
			fmt.Println(stack.StringF())
			return nil, nil
		},
	}

	fcOp = Op{
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
		"feet to meters conversion",
		func(stack *Stack) (Floats, error) {
			top, err := stack.Pop()
			if err != nil {
				return nil, err
			}
			return Floats{top / _ftPerM}, nil
		},
	}

	frexpOp = Op{
		"breaks stack.Top() into a normalized fraction and an integral power of two",
		func(stack *Stack) (Floats, error) {
			top, err := stack.Pop()
			if err != nil {
				return nil, err
			}
			frac, exp := math.Frexp(top)
			return Floats{frac, float64(exp)}, nil
		},
	}

	glOp = Op{
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
		"no op",
		func(stack *Stack) (Floats, error) {
			return nil, nil
		},
	}

	pOp = Op{
		"pop an item from the stack",
		func(stack *Stack) (Floats, error) {
			_, _ = stack.Pop()
			return nil, nil
		},
	}

	pasOp = Op{
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
		"pressure in inHg for a given altitude in feet",
		func(stack *Stack) (Floats, error) {
			top, err := stack.Pop()
			if err != nil {
				return nil, err
			}
			return Floats{29.9212524 * math.Pow(1-math.Pow(10, -5)*2.25577*(top/3.280839895), 5.25588)}, nil
		},
	}

	qOp = Op{
		"exit the program",
		func(stack *Stack) (Floats, error) {
			os.Exit(0)
			return nil, nil
		},
	}

	randOp = Op{
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
		"sort the entire stack",
		func(stack *Stack) (Floats, error) {
			stack.Sort()
			return nil, nil
		},
	}

	sumOp = Op{
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

func wrapConstant(doc string, value float64) Op {
	return Op{
		doc,
		func(stack *Stack) (Floats, error) {
			return Floats{value}, nil
		},
	}
}

func wrapUnaryOp(doc string, f func(float64) float64) Op {
	return Op{
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

func wrapBinaryOp(doc string, f func(float64, float64) float64) Op {
	return Op{
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

func wrapTernaryOp(doc string, f func(float64, float64, float64) float64) Op {
	return Op{
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
