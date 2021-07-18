package lox

import (
	"fmt"
)

var hadErr bool = false

func errorReport(line int, message string) {

	report(line, "", message)
}

func report(line int, where, message string) {

	fmt.Printf("[line %d] Error %s: %s", line, where, message)
	hadErr = true
}
