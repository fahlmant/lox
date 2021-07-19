package lox

type Expr interface {
	// Visitor pattern
	Accept(ExprVisitor) error
}

type ExprVisitor interface {
	visitBinary(Binary) error
	visitGrouping(Grouping) error
	visitLiteral(Literal) error
	visitUnary(Unary) error
}

type Binary struct {
	Left     Expr
	Operator Token
	Right    Expr
}

func (b Binary) Accept(visitor ExprVisitor) error {
	return visitor.visitBinary(b)
}

type Grouping struct {
	Expression Expr
}

func (g Grouping) Accept(visitor ExprVisitor) error {
	return visitor.visitGrouping(g)
}

type Literal struct {
	Value interface{}
}

func (l Literal) Accept(visitor ExprVisitor) error {
	return visitor.visitLiteral(l)
}

type Unary struct {
	Operator Token
	Right    Expr
}

func (u Unary) Accept(visitor ExprVisitor) error {
	return visitor.visitUnary(u)
}
