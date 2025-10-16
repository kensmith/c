package main

import (
	"fmt"
)

func main() {
	operators := NewOperatorMap()

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
		err := tryCascade(line, stack, operators)
		if err != nil {
			if len(lastLine) > 0 {
				err = tryCascade(lastLine, stack, operators)
				if err != nil {
					fmt.Println(err)
				}
				continue
			}
		}
		lastLine = line
	}
}
