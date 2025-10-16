package main

import (
	"fmt"
)

func main() {
	ops := NewOps()

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
		err := cascade(line, stack, ops)
		if err != nil {
			if len(lastLine) > 0 {
				err = cascade(lastLine, stack, ops)
				if err != nil {
					fmt.Println(err)
				}
				continue
			}
		}
		lastLine = line
	}
}
