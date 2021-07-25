package lox

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

// RunFile runs a supplied file
func RunFile(path string) {
	fmt.Printf("Running file %s\n", path)
	out, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Printf("%n", err)
		return
	}
	run(string(out))
}

// RunPrompt begins an interactive session
func RunPrompt() {

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		text, _ := reader.ReadString('\n')
		if text == "\n" || text == "" {
			fmt.Println("Recieved blank line, qutting")
			break
		}
		// convert CRLF to LF
		text = strings.Replace(text, "\n", "", -1)
		//fmt.Printf("Confiming message: %s\n", text)
		run(text)
	}
}

func run(source string) {
	// Create Scanner with the input as a []rune
	s := Scanner{source: []rune(source)}
	// Generate token list
	s.scanTokens()

	// Iterate through all tokens
	/*for _, token := range s.tokens {
		fmt.Printf("TOKEN: %+v\n", token)
	}*/

	p := Parser{tokens: s.tokens}
	stms, err := p.parse()
	if err != nil {
		fmt.Println(err)
	}

	if !p.hadError {
		var i Interpreter
		err = i.Interpret(stms)
		if err != nil {
			fmt.Println(err)
		}
	}
}
