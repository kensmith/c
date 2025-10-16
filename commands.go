package main

import (
	"fmt"
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
	}

	return nil
}
