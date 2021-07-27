package lox

import (
	"fmt"
)

// Token contains the Type, Lexeme, Literal and Line Number
type Token struct {
	tType   TokenType
	lexeme  string
	line    int
	literal string
}

func (t Token) String() string {
	return fmt.Sprintf("%v %d %s %s", t.tType, t.line, t.lexeme, t.literal)
}
