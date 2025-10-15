package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/chzyer/readline"
	"github.com/expr-lang/expr"
)

func pop2(stack *[]float64) (float64, float64, error) {
	size := len(*stack)
	if size < 2 {
		return 0.0, 0.0, fmt.Errorf("can't pop 2 for stack of length %d", size)
	}
	rhs := (*stack)[size-1]
	lhs := (*stack)[size-2]
	*stack = (*stack)[:size-2]
	return lhs, rhs, nil
}

type Operator func(*[]float64) float64

var operators = map[string]Operator{
	"+": func(stack *[]float64) float64 {
		fmt.Println("hi")
		lhs, rhs, err := pop2(stack)
		if err != nil {
			fmt.Println("couldn't add")
			return 0.0
		}
		return lhs + rhs
	},
	"-": func(stack *[]float64) float64 {
		lhs, rhs, err := pop2(stack)
		if err != nil {
			fmt.Println("couldn't add")
			return 0.0
		}
		return lhs - rhs
	},
}

func main() {
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
			result := rawOperator(&stack)
			stack = append(stack, result)
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
			fmt.Fprintf(&b, "%v", n)
			if i < stackSize-1 {
				fmt.Fprintf(&b, "  ")
			}
		}
		fmt.Fprintf(&b, " ]> ")
		shell.SetPrompt(b.String())
	}
}
