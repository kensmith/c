package main

import (
	"fmt"
	"os"
)

func tryCommands(line string, stack *Stack, ops *Ops) error {
	switch line {
	case "?":
		fallthrough
	case "h":
		fallthrough
	case "he":
		fallthrough
	case "hel":
		fallthrough
	case "help":
		fmt.Println(ops.Help())
	case "sort":
		stack.Sort()
	case "f":
		fmt.Println(stack.StringF())
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
