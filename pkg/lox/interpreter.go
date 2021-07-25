package lox

import (
	"fmt"
)

// Represents an interpreter and associated logic
type Interpreter struct {
	literal     Literal
	environment *Environment
}

// Main interpretation loop
func (i *Interpreter) Interpret(stmts []Stmt) error {

	// Create a new environment to initalize the empty map
	i.environment = NewEnvironment(nil)

	// Loop through all statements
	for _, stmt := range stmts {
		// Exectue the logic for each statement with the vistiro pattern
		if err := stmt.Accept(i); err != nil {
			return err
		}
	}
	return nil
}

func (i *Interpreter) visitAssign(a Assign) error {

	l, err := i.evaluate(a.Value)
	if err != nil {
		return err
	}

	if err := i.environment.Assign(a.Var, l); err != nil {
		return err
	}

	if l, ok := a.Value.(Literal); ok {
		i.literal = l
	}

	return nil
}

// Implementations of required functions for visitor pattern
func (i *Interpreter) visitLiteral(l Literal) error {

	// With a literal, it can be returned as is
	i.literal = l
	return nil
}

func (i *Interpreter) visitBinary(b Binary) error {
	// For binary, both sides need to be evaluated

	// Evalaute the right side expression, expeting a single value
	right, err := i.evaluate(b.Right)
	if err != nil {
		return err
	}

	// Evalaute the left side expression, expeting a single value
	left, err := i.evaluate(b.Left)
	if err != nil {
		return err
	}

	operationError := fmt.Errorf("error at line %d: bad operand for binary %s: %T, %T", b.Operator.Line, b.Operator.Lexeme, left.Value, right.Value)

	switch b.Operator.TType {
	// Minus, Slash and Star attempt to convert both operands to float64 and then calculate the result
	case MINUS:
		// Check both sides to make sure they can be converted to float64s
		if l, ok := left.Value.(float64); ok {
			if r, ok := right.Value.(float64); ok {
				i.literal = Literal{l - r}
				return nil
			}
		}
		return operationError

	case SLASH:
		// Check both sides to make sure they can be converted to float64s
		if l, ok := left.Value.(float64); ok {
			if r, ok := right.Value.(float64); ok {
				i.literal = Literal{l / r}
				return nil
			}
		}
		return operationError
	case STAR:
		// Check both sides to make sure they can be converted to float64s
		if l, ok := left.Value.(float64); ok {
			if r, ok := right.Value.(float64); ok {
				i.literal = Literal{l * r}
				return nil
			}
		}
		return operationError
	// Plus does the same operation on float64. If the values do not convert, string concatenation is attempted
	case PLUS:

		// Check both sides to make sure they can be converted to float64s
		if l, ok := left.Value.(float64); ok {
			if r, ok := right.Value.(float64); ok {
				i.literal = Literal{l + r}
				return nil
			}
			// Plus can also work on strings , so check that as well
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
		// Check both sides to make sure they can be converted to float64s
		if l, ok := left.Value.(float64); ok {
			if r, ok := right.Value.(float64); ok {
				i.literal = Literal{l > r}
				return nil
			}
		}
		return operationError

	case GREATER_EQUAL:
		// Check both sides to make sure they can be converted to float64s
		if l, ok := left.Value.(float64); ok {
			if r, ok := right.Value.(float64); ok {
				i.literal = Literal{l >= r}
			}
		} else {
			return operationError
		}
	case LESS:
		// Check both sides to make sure they can be converted to float64s
		if l, ok := left.Value.(float64); ok {
			if r, ok := right.Value.(float64); ok {
				i.literal = Literal{l < r}
				return nil
			}
		}
		return operationError

	case LESS_EQUAL:
		// Check both sides to make sure they can be converted to float64s
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

// Visitor pattern for block statements.
// Creates a new environment, environment b, and sets it parent to the current environment, environment a
// Then sets the current environment to the new one (environment b) and evaluates all statements in it
// Then sets the current environment back to the original
func (i *Interpreter) visitBlockStmt(b BlockStmt) error {
	// Create a new environment as the child of the current environment
	// as assign it as our current environment
	i.environment = NewEnvironment(i.environment)

	// Range through statements and evaluate them
	for _, stmt := range b.statements {
		if err := stmt.Accept(i); err != nil {
			return err
		}
	}

	// Reset the environment back to the original environment
	i.environment = i.environment.enclosing

	return nil

}

// Visitor pattern for expressions statements. Evaluates the expression with the vistior patern
func (i *Interpreter) visitExprStmt(e ExprStmt) error {
	return e.expression.Accept(i)
}

// Visitor pattern of if statements.
// Evaluates the condition, and if its true, evaluates the branch
// IF not true, checks the else, and if it exists, evaluates the else
func (i *Interpreter) visitIfStmt(ifStmt IfStmt) error {

	expr, err := i.evaluate(ifStmt.condition)
	if err != nil {
		return err
	}

	if isTruthy(expr) {
		err = ifStmt.branch.Accept(i)
		if err != nil {
			return err
		}
	} else {
		if ifStmt.elseStmt != nil {
			err = ifStmt.elseStmt.Accept(i)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Visitor pattern for print statements. Evaluates and then prints the result
func (i *Interpreter) visitPrintStmt(p PrintStmt) error {

	// Evaluate the expression
	expr, err := i.evaluate(p.expression)
	if err != nil {
		return err
	}
	// Print the result
	fmt.Println(expr)

	return nil
}

// Visitor pattern for Var statements
func (i *Interpreter) visitVarStmt(v VarStmt) error {

	// Reset the interpreter literal
	i.literal = Literal{nil}

	// If the right hand size exists, evaluate it
	// Evaluate here sets i.literal to the result of the expression
	if v.initializer != nil {
		if _, err := i.evaluate(v.initializer); err != nil {
			return err
		}
	}

	// Define the variable and map it to the current value of i.literal
	// If there's no value to assign, it will map to a Literal{nil}
	// Otherwise it will map to the result of the expression in the initializer
	if err := i.environment.Define(Variable{v.name}, i.literal); err != nil {
		return err
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

func (i *Interpreter) visitVariable(v Variable) error {

	// Attempt to lookup the varible in the environment map
	value, err := i.environment.Get(v)
	if err != nil {
		return err
	}

	// Place the value of variable in the interpreter literal
	if lit, ok := value.(Literal); ok {
		i.literal = lit
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
