package main

import (
	"os"
)

func tryCommands(line string, stack *Stack) error {
	switch line {
	case "sw":
		fallthrough
	case "swap":
		err := stack.Swap()
		if err != nil {
			return err
		}
	case "p":
		fallthrough
	case "pop":
		_, err := stack.Pop()
		if err != nil {
			return err
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
		fallthrough
	case "exit":
		os.Exit(0)
	}

	return nil
}
