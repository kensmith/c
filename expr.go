package main

import (
	"fmt"
	"slices"
	"strconv"

	"github.com/expr-lang/expr"
)

func tryExpr(line string, stack *Stack) error {
	localStack := stack.Copy()
	slices.Reverse(localStack)
	env := map[string]any{
		"s": localStack,
	}

	output, err := expr.Eval(line, env)
	if err != nil {
		return err
	}

	outputStr := fmt.Sprintf("%v", output)
	value, err := strconv.ParseFloat(outputStr, 64)
	if err != nil {
		return err
	}
	stack.Push(value)

	return nil
}
