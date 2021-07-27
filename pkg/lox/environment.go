package lox

import "fmt"

// Tracks all variable assignments with a map
type Environment struct {
	values    map[string]interface{}
	enclosing *Environment
}

// Returns a new environment with an initalized map
func NewEnvironment(enclosing *Environment) *Environment {
	return &Environment{values: make(map[string]interface{}), enclosing: enclosing}
}

// Assigns a value to an existing variable.
// Does NOT allow for first time variable declaration. See Define()
func (e *Environment) Assign(v Variable, expr Expr) error {
	// Checks if the var name exists in the values map
	if _, ok := e.values[v.token.lexeme]; ok {
		// If the var name exists, update the value
		e.values[v.token.lexeme] = expr
		return nil
	}

	// Check the higher scope for the variable if not in the current scope
	if e.enclosing != nil {
		return e.enclosing.Assign(v, expr)
	}

	// If the var is not in the map, return an error
	return fmt.Errorf("error at line %d: undefined variable %v", v.token.line, v.token.lexeme)

}

// Retrieves a value from the environment map
func (e *Environment) Get(v Variable) (interface{}, error) {

	// Check to ensure the var exists in the map
	if value, ok := e.values[v.token.lexeme]; ok {
		// If the var exists, return the value
		return value, nil
	}

	// If there's a higher scope than the current one, check it for the variable
	if e.enclosing != nil {
		return e.enclosing.Get(v)
	}

	// If the var is not in the map, return an error
	return nil, fmt.Errorf("error at line %d: undefined variable %v", v.token.line, v.token.lexeme)
}

// Defines a new variable in the map
func (e *Environment) Define(v Variable, expr Expr) error {

	// Creates an entry in the map for a new variable and its definition
	e.values[v.token.lexeme] = expr

	return nil
}
