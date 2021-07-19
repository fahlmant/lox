package lox

import (
	"fmt"
)

// Token contains the Type, Lexeme, Literal and Line Number
type Token struct {
	TType   TokenType
	Lexeme  string
	Line    int
	Literal string
}

func (t Token) String() string {
	return fmt.Sprintf("%v %d %s %s", t.TType, t.Line, t.Lexeme, t.Literal)
}
