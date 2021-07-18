package lox

import (
	"fmt"
)

// Token contains the Type, Lexeme, Literal and Line Number
type Token struct {
	tokenType TokenType
	lexeme    string
	line      int
	literal   string
}

func (t Token) String() string {
	return fmt.Sprintf("%v %d %s %s", t.tokenType, t.line, t.lexeme, t.literal)
}
