package main

import (
	"fmt"
)

func tryCascade(line string, stack *Stack, operators OperatorMap) error {
	err := tryExpr(line, stack)
	if err == nil {
		return nil
	}

	err = tryOperator(line, stack, operators)
	if err == nil {
		return nil
	}

	err = tryCommands(line, stack)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	operators := NewOperatorMap()
	installUnaryFunctions(operators)
	installBinaryFunctions(operators)
	installTernaryFunctions(operators)
	installConstants(operators)

	shell := NewShell()
	defer shell.Close()

	stack := NewStack()

	lastLine := ""
	for {
		shell.SetPrompt(stack.String() + "> ")
		line := shell.ReadLine()
		if len(line) <= 0 {
			line = lastLine
		}
		err := tryCascade(line, &stack, operators)
		if err != nil {
			if len(lastLine) > 0 {
				err = tryCascade(lastLine, &stack, operators)
				if err != nil {
					fmt.Println(err)
				}
				continue
			}
		}
		lastLine = line
	}
}
