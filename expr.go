package main

import (
	"fmt"
	"strconv"

	"github.com/expr-lang/expr"
)

func tryExpr(line string, stack *Stack) error {
	output, err := expr.Eval(line, nil)
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
