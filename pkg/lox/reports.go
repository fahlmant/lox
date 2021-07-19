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

func errorToken(t Token, message string) {
	if t.TType == EOF {
		report(t.Line, " at end", message)
	} else {
		report(t.Line, "at '"+t.Lexeme+"'", message)
	}
}

// Test AST Printer to test the Visitor pattern
type AstPrinter struct{}

func (a AstPrinter) VisitBinary(b Binary) error {
	fmt.Printf("(%v %v %v)", b.Operator.Lexeme, b.Left, b.Right)

	return nil
}

func (a AstPrinter) VisitGrouping(g Grouping) error {
	fmt.Printf("(group %v)", g.Expression)

	return nil
}

func (a AstPrinter) VisitLiteral(l Literal) error {
	str := fmt.Sprintf("%v", l.Value)
	fmt.Printf("(%v\n)", str)

	return nil
}

func (a AstPrinter) VisitUnary(u Unary) error {
	fmt.Printf("(%v %v)", u.Operator.Lexeme, u.Right)

	return nil
}
