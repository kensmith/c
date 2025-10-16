package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/chzyer/readline"
	"github.com/expr-lang/expr"
)

func main() {
	operators := NewOperatorMap()
	installUnaryFunctions(operators)
	installBinaryFunctions(operators)
	installTernaryFunctions(operators)
	installConstants(operators)

	var shell *readline.Instance
	err := os.MkdirAll(_histDirname, 0o750)
	if err != nil {
		fmt.Printf("history disabled due to inability to create directory: %s", _histDirname)
		shell, err = readline.New("[  ]> ")
		if err != nil {
			panic(err)
		}
	} else {
		shell, err = readline.NewEx(&readline.Config{
			Prompt:      "[  ]> ",
			HistoryFile: _histFilename,
		})
	}
	if err != nil {
		panic(err)
	}
	defer func() {
		err := shell.Close()
		if err != nil {
			fmt.Println(err)
		}
	}()

	stack := Stack{}
	env := map[string]any{
		"stack": &stack,
		"s":     &stack,
	}

	for {
		line, err := shell.Readline()
		if err != nil {
			// normal exit due to ctrl-c, ctrl-d
			return
		}
		commasRemoved := strings.ReplaceAll(line, ",", "")
		lineTrimmed := strings.TrimSpace(commasRemoved)
		rawOperator, ok := operators[lineTrimmed]
		if ok {
			result, err := rawOperator(&stack)
			if err != nil {
				fmt.Println(err)
			} else {
				stack.Push(result)
			}
		} else {
			switch lineTrimmed {
			case "sw":
				fallthrough
			case "swap":
				err := stack.Swap()
				if err != nil {
					fmt.Println(err)
					continue
				}
			case "p":
				fallthrough
			case "pop":
				_, err := stack.Pop()
				if err != nil {
					fmt.Println(err)
					continue
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
				return
			default:
				output, err := expr.Eval(lineTrimmed, env)
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
				stack.Push(value)
			}
		}
		shell.SetPrompt(stack.String() + "> ")
	}
}
