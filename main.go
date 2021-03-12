package main

import (
	"fmt"
	"os"

	"github.com/fahlmant/lox/pkg/lox"
)

func main() {

	argLength := len(os.Args)

	if argLength > 2 {
		fmt.Println("Usage: golox [script]")
	} else if argLength == 2 {
		lox.RunFile(os.Args[1])
	} else {
		lox.RunPrompt()
	}
}
