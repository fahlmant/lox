package lox

import (
	"fmt"
)

type Interpreter struct {
	literal Literal
}

func (i *Interpreter) Interpret(e Expr) (Literal, error) {
	value, err := i.evaluate(e)
	if err != nil {
		return value, err
	}
	return value, nil
}

// Implementations of required functions for visitor pattern
func (i *Interpreter) visitLiteral(l Literal) error {

	// With a literal, we can just return it as is
	i.literal = l
	return nil
}

func (i *Interpreter) visitBinary(b Binary) error {
	right, err := i.evaluate(b.Right)
	if err != nil {
		return err
	}
	left, err := i.evaluate(b.Left)
	if err != nil {
		return err
	}

	operationError := fmt.Errorf("error at line %d: bad operand for binary %s: %T, %T", b.Operator.Line, b.Operator.Lexeme, left.Value, right.Value)

	switch b.Operator.TType {
	// Minus, Slash and Star attempt to convert both operands to float64 and then calculate the result
	case MINUS:
		if l, ok := left.Value.(float64); ok {
			if r, ok := right.Value.(float64); ok {
				i.literal = Literal{l - r}
				return nil
			}
		}
		return operationError

	case SLASH:
		if l, ok := left.Value.(float64); ok {
			if r, ok := right.Value.(float64); ok {
				i.literal = Literal{l / r}
				return nil
			}
		}
		return operationError
	case STAR:
		if l, ok := left.Value.(float64); ok {
			if r, ok := right.Value.(float64); ok {
				i.literal = Literal{l * r}
				return nil
			}
		}
		return operationError
	// Plus does the same operation on float64. If the values do not convert, string concatenation is attempted
	case PLUS:

		if l, ok := left.Value.(float64); ok {
			if r, ok := right.Value.(float64); ok {
				i.literal = Literal{l + r}
				return nil
			}
		} else if l, ok := left.Value.(string); ok {
			if r, ok := right.Value.(string); ok {
				i.literal = Literal{l + r}
				return nil
			}
		}
		return operationError
	// The rest are simple truthy-checks
	// Greater, Greater Equal, Less and Less Equal only operator on float64
	// Bang Equal and Equal Equal operate on any values
	case GREATER:
		if l, ok := left.Value.(float64); ok {
			if r, ok := right.Value.(float64); ok {
				i.literal = Literal{l > r}
				return nil
			}
		}
		return operationError

	case GREATER_EQUAL:
		if l, ok := left.Value.(float64); ok {
			if r, ok := right.Value.(float64); ok {
				i.literal = Literal{l >= r}
			}
		} else {
			return operationError
		}
	case LESS:
		if l, ok := left.Value.(float64); ok {
			if r, ok := right.Value.(float64); ok {
				i.literal = Literal{l < r}
				return nil
			}
		}
		return operationError

	case LESS_EQUAL:
		if l, ok := left.Value.(float64); ok {
			if r, ok := right.Value.(float64); ok {
				i.literal = Literal{l <= r}
				return nil
			}
		}
		return operationError
	case BANG_EQUAL:
		i.literal = Literal{left != right}
	case EQUAL_EQUAL:
		i.literal = Literal{left == right}
	}

	return nil
}

func (i *Interpreter) visitGrouping(g Grouping) error {

	// Send the expression back into the visitor
	_, err := i.evaluate(g.Expression)

	return err
}

func (i *Interpreter) visitUnary(u Unary) error {

	// Evaluate the expression that is being operated on
	_, err := i.evaluate(u.Right)
	if err != nil {
		return err
	}

	switch u.Operator.TType {
	case MINUS:
		// Ensure the value can be converted to a float64 and return the negation if it
		if value, ok := i.literal.Value.(float64); ok {
			i.literal = Literal{-value}
		} else {
			// Indicates the value cannot be converted into a number and cannot be negated
			return fmt.Errorf("error at line %d: bad operand for unary %s: %T", u.Operator.Line, u.Operator.Lexeme, i.literal.Value)
		}
	case BANG:
		// Invert the truthiness i.e. var a = true; !a;
		i.literal = Literal{!isTruthy(i.literal)}
	}

	return nil
}

func (i *Interpreter) evaluate(expr Expr) (Literal, error) {
	// Use the visitor to continue to evaluate the expression
	err := expr.Accept(i)
	return i.literal, err
}

func isTruthy(l Literal) bool {
	// If the value is already a bool, just return in
	if b, ok := l.Value.(bool); ok {
		return b
	} else {
		// If the value is nil, return false
		if l.Value == nil {
			return false
		} else {
			//Every other value is true
			return true
		}
	}
}
