package lox

// Expressions are combinations of values and operators
// that create a new value
type Expr interface {
	// Visitor pattern
	Accept(ExprVisitor) error
}

// Interface to implement to interact with expressions
// using the visitor pattern
type ExprVisitor interface {
	visitAssign(Assign) error
	visitBinary(Binary) error
	visitGrouping(Grouping) error
	visitLiteral(Literal) error
	visitLogical(Logical) error
	visitUnary(Unary) error
	visitVariable(Variable) error
}

// Represents an assignment expression
// example: var a = 1
type Assign struct {
	variable Variable
	name     Token
	value    Expr
}

// Boilerplate visitor pattern for Assign
func (a Assign) Accept(visitor ExprVisitor) error {
	return visitor.visitAssign(a)
}

// Represents a binary expressions
// example 1 + 2
// example (a+2) / (b-2) {Nested binary expressions}
type Binary struct {
	left     Expr
	operator Token
	right    Expr
}

// Boilerplate visitor pattern for Binary
func (b Binary) Accept(visitor ExprVisitor) error {
	return visitor.visitBinary(b)
}

// Represents a grouping of expressions
type Grouping struct {
	expression Expr
}

// Boilerplate visitor pattern for Grouping
func (g Grouping) Accept(visitor ExprVisitor) error {
	return visitor.visitGrouping(g)
}

// Represents a singular value, such as a number or a string
type Literal struct {
	value interface{}
}

// Boilerplate visitor pattern for Literal
func (l Literal) Accept(visitor ExprVisitor) error {
	return visitor.visitLiteral(l)
}

// Represents a logical "and" or "or"
type Logical struct {
	left     Expr
	operator Token
	right    Expr
}

func (l Logical) Accept(visitor ExprVisitor) error {
	return visitor.visitLogical(l)
}

// Represetns unary operations
// example: -1 or !true
type Unary struct {
	operator Token
	right    Expr
}

// Boilerplate visitor pattern for Unary
func (u Unary) Accept(visitor ExprVisitor) error {
	return visitor.visitUnary(u)
}

// Represents a variable name for access
// Example: print foo
type Variable struct {
	token Token
}

// Boilerplate visitor pattern for Variable
func (v Variable) Accept(visitor ExprVisitor) error {
	return visitor.visitVariable(v)
}
