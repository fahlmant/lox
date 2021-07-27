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
	if t.tType == EOF {
		report(t.line, " at end", message)
	} else {
		report(t.line, "at '"+t.lexeme+"'", message)
	}
}

// Test AST Printer to test the Visitor pattern
type AstPrinter struct{}

func (a AstPrinter) VisitBinary(b Binary) error {
	fmt.Printf("(%s %v %v)", b.operator.lexeme, b.left, b.right)

	return nil
}

func (a AstPrinter) VisitGrouping(g Grouping) error {
	fmt.Printf("(group %v)", g.expression)

	return nil
}

func (a AstPrinter) VisitLiteral(l Literal) error {
	str := fmt.Sprintf("%v", l.value)
	fmt.Printf("(%v\n)", str)

	return nil
}

func (a AstPrinter) VisitUnary(u Unary) error {
	fmt.Printf("(%s %v)", u.operator.lexeme, u.right)

	return nil
}
