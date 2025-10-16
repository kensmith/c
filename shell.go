package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/chzyer/readline"
)

type Shell struct {
	internal *readline.Instance
}

func NewShell() *Shell {
	shell := Shell{}
	err := os.MkdirAll(_histDirname, 0o750)
	if err != nil {
		fmt.Printf("history disabled due to inability to create directory: %s", _histDirname)
		shell.internal, err = readline.New("[  ]> ")
		if err != nil {
			panic(err)
		}
	} else {
		shell.internal, err = readline.NewEx(&readline.Config{
			Prompt:      "> ",
			HistoryFile: _histFilename,
		})
		if err != nil {
			panic(err)
		}
	}

	return &shell
}

func (s *Shell) SetPrompt(prompt string) {
	s.internal.SetPrompt(prompt)
}

func (s *Shell) ReadLine() string {
	line, err := s.internal.Readline()
	if err != nil {
		// normal exit due to ctrl-c, ctrl-d
		return "exit"
	}
	commasRemoved := strings.ReplaceAll(line, ",", "")
	lineTrimmed := strings.TrimSpace(commasRemoved)

	return lineTrimmed
}

func (s *Shell) Close() {
	err := s.internal.Close()
	if err != nil {
		fmt.Println(err)
	}
}
