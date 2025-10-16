package main

import (
	"fmt"
	"os"

	"github.com/chzyer/readline"
)

func NewShell() *readline.Instance {
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
		if err != nil {
			panic(err)
		}
	}

	return shell
}
