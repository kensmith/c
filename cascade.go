package main

import "fmt"

func cascade(line string, stack *Stack, ops *Ops) error {
	err := tryExpr(line, stack)
	if err == nil {
		return nil
	}

	err = ops.Run(line, stack)
	if err == nil {
		return nil
	}

	fmt.Println(ops.Help())

	return nil
}
