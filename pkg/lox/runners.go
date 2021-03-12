package lox

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// RunFile runs a supplied file
func RunFile(path string) {
	fmt.Printf("Running file %s\n", path)
}

// RunPrompt begins an interactive session
func RunPrompt() {

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		text, _ := reader.ReadString('\n')
		if text == "\n" {
			fmt.Println("Recieved blank line, qutting")
			break
		}
		// convert CRLF to LF
		text = strings.Replace(text, "\n", "", -1)
		fmt.Printf("Confiming message: %s\n", text)
	}
}

func run() {

}
