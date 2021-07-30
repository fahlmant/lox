package lox

import (
	"fmt"
	"time"
)

type LoxCallable interface {
	arity() int
	call(interpreter *Interpreter, arguments []Expr) (Literal, error)
}

type Clock struct{}

func (c Clock) arity() int {
	return 0
}

func (c Clock) call(interpreter *Interpreter, arguments []Expr) (Literal, error) {
	return Literal{time.Now().Unix()}, nil
}

type Print struct{}

func (p Print) arity() int {
	return 1
}

func (p Print) call(interpreter *Interpreter, arguments []Expr) (Literal, error) {
	fmt.Println(arguments[0])
	return Literal{}, nil
}

func (f FuncStmt) arity() int {
	return len(f.params)
}

func (f FuncStmt) call(interpreter *Interpreter, arguments []Expr) (Literal, error) {

	// Create new scope for the function
	interpreter.environment = NewEnvironment(interpreter.environment)

	// Place all arguments into the scope of the function as variables
	for i, arg := range arguments {
		expr, err := interpreter.evaluate(arg)
		if err != nil {
			return Literal{}, err
		}

		if err := interpreter.environment.Define(Variable{f.params[i]}, expr); err != nil {
			return Literal{}, err
		}
	}

	// Execute stmts in the body
	for _, stmt := range f.body {
		if err := stmt.Accept(interpreter); err != nil {
			// If the statement is a return (as an error), escape the scope of the func and return the value
			if r, ok := err.(ReturnValue); ok {
				interpreter.environment = interpreter.environment.enclosing
				return r.Literal, nil

			}
			return Literal{}, err
		}
	}

	interpreter.environment = interpreter.environment.enclosing

	return Literal{}, nil

}
