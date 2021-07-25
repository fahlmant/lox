package lox

import (
	"fmt"
	"strconv"
)

type Parser struct {
	tokens     []Token
	current    int
	statements []Stmt
	hadError   bool
}

// Main parsing loop
func (p *Parser) parse() ([]Stmt, error) {
	for !p.isAtEnd() {
		// Start at the top of the statement recursive tree and get a single statement
		// for each section of code
		stmt, err := p.declaration()
		if err != nil {
			// Set hadError so the interpreter doesn't run
			p.hadError = true
			return nil, err
		}
		p.statements = append(p.statements, stmt)
	}

	return p.statements, nil
}

func (p *Parser) declaration() (Stmt, error) {
	// If there's a variable declaration, handle it
	if p.match(VAR) {
		return p.varDeclaration()
	}
	// Otherwise, handle the section as a statement
	return p.statement()
}

func (p *Parser) varDeclaration() (Stmt, error) {

	// Consume and load the identifier name into a token
	token, err := p.consume(IDENTIFIER, "Expect variable name")
	if err != nil {
		return nil, err
	}

	var initalizer Expr
	// If there's an assignment, then evaluate the expression for use in the statement
	if p.match(EQUAL) {
		if initalizer, err = p.expression(); err != nil {
			return nil, err
		}
	}

	// Ensure there is an ending semicolon
	if _, err := p.consume(SEMICOLON, "Expect ; after variable declaration"); err != nil {
		return nil, err
	}

	// Return as a Var Statement so the interpreter knows to
	// track the variable with the value
	return VarStmt{name: token, initializer: initalizer}, nil
}

func (p *Parser) statement() (Stmt, error) {

	// If there's an if statement, handle it
	if p.match(IF) {
		return p.ifStatement()
	}

	// If there's a print statement, handle it
	if p.match(PRINT) {
		return p.printStatement()
	}
	// If theres a {, build a BlockStmt with all the statements before }
	if p.match(LEFT_BRACE) {
		// Build the statement list by calling block()
		stmts, err := p.block()
		if err != nil {
			return nil, err
		}
		// Build the BlockStmt statment and return
		return BlockStmt{statements: stmts}, nil
	}

	// Otherwise, handle the generic expression case
	return p.expressionStatement()
}

func (p *Parser) block() ([]Stmt, error) {
	var statements []Stmt

	// Create a list of statements between { and }
	for !(p.peek().TType == RIGHT_BRACE) && !p.isAtEnd() {
		stmt, err := p.declaration()
		if err != nil {
			return nil, err
		}
		statements = append(statements, stmt)
	}

	// Consume the right brace to finish the scope
	if _, err := p.consume(RIGHT_BRACE, "Expeted } to close block"); err != nil {
		return nil, err
	}

	// Return the list of statments
	// Note this does NOT return a singluar statement (wrapped in a BlockStmt)
	// This is becuase its a convient pice of code to also use for getting all statements
	// within a lox function. If this returned a BlockStmt{}, then we could not reuse this code
	return statements, nil
}

func (p *Parser) ifStatement() (Stmt, error) {

	// Ensure there is a left paren after an "if"
	if _, err := p.consume(LEFT_PAREN, "Expect ( after an if"); err != nil {
		return nil, err
	}

	// Get the condition expression of the if
	condition, err := p.expression()
	if err != nil {
		return nil, err
	}

	// Ensure there is a right paren after the condition
	if _, err := p.consume(RIGHT_PAREN, "Expect ) after an if condition"); err != nil {
		return nil, err
	}

	// Get the statement to do with the conditional
	thenStmt, err := p.statement()
	if err != nil {
		return nil, err
	}

	// If there's an "else", also get that statement
	var elseStmt Stmt
	if p.match(ELSE) {
		elseStmt, err = p.statement()
		if err != nil {
			return nil, err
		}
	}

	// Return an IfStatement with the condition, the expression and the else if there is one
	return IfStmt{condition: condition, branch: thenStmt, elseStmt: elseStmt}, nil
}

func (p *Parser) printStatement() (Stmt, error) {
	// Expand the following espression to print out
	value, err := p.expression()
	if err != nil {
		return nil, err
	}

	// Ensure there is an ending semicolon
	_, err = p.consume(SEMICOLON, "Expect ; after value")
	if err != nil {
		return nil, err
	}

	// Return as a print statement so the interpreter knows to print
	return PrintStmt{expression: value}, nil
}

func (p *Parser) expressionStatement() (Stmt, error) {
	// Expand the expression
	value, err := p.expression()
	if err != nil {
		return nil, err
	}

	// Ensure there is an ending semicolon
	_, err = p.consume(SEMICOLON, "Expect ; after value")
	if err != nil {
		return nil, err
	}

	// Return as a generic Expression
	return ExprStmt{expression: value}, nil
}

func (p *Parser) expression() (Expr, error) {
	// As of now, a passthrough in the recursive expansion
	return p.assignment()
}

func (p *Parser) assignment() (Expr, error) {

	// Goes down the recursive tree to get the expression
	expr, err := p.equality()
	if err != nil {
		return nil, err
	}

	// If there's an assignment happening, handle it
	if p.match(EQUAL) {
		if equals, ok := p.previous(); ok {
			// Evaluate the value to assign to the variable
			value, err := p.assignment()
			if err != nil {
				return nil, err
			}

			// Ensure the left side is a variable that can be assigned
			if v, ok := expr.(Variable); ok {
				// Build the variable assignment expression
				return Assign{Var: v, Name: equals, Value: value}, nil
			}

			return nil, fmt.Errorf("error at line %d: invalid assiment target", equals.Line)
		}
	}

	// Returns the expression by itself if there's no assignment happening
	return expr, nil
}

func (p *Parser) equality() (Expr, error) {
	// Build the first expression
	expr, err := p.comparison()
	if err != nil {
		return nil, err
	}

	// If there's a comparison happening (== or !=) handle it
	for p.match(EQUAL_EQUAL, BANG_EQUAL) {
		if operator, ok := p.previous(); ok {
			// Build the right side expression of the comparison
			right, err := p.comparison()
			if err != nil {
				return nil, err
			}
			return Binary{Left: expr, Operator: operator, Right: right}, nil
		}
	}

	// Return the expression by itself if there's no comparison
	return expr, nil
}

func (p *Parser) comparison() (Expr, error) {

	//
	expr, err := p.term()
	if err != nil {
		return nil, err
	}

	for p.match(GREATER, GREATER_EQUAL, LESS, LESS_EQUAL) {
		if operator, ok := p.previous(); ok {
			right, err := p.term()
			if err != nil {
				return nil, err
			}
			expr = Binary{Left: expr, Operator: operator, Right: right}
		}
	}

	return expr, nil
}

func (p *Parser) term() (Expr, error) {

	expr, err := p.factor()
	if err != nil {
		return nil, err
	}

	for p.match(MINUS, PLUS) {
		if operator, ok := p.previous(); ok {
			right, err := p.factor()
			if err != nil {
				return nil, err
			}
			expr = Binary{Left: expr, Operator: operator, Right: right}

		}
	}

	return expr, nil
}

func (p *Parser) factor() (Expr, error) {
	expr, err := p.unary()
	if err != nil {
		return nil, err
	}

	for p.match(SLASH, STAR) {
		if operator, ok := p.previous(); ok {
			right, err := p.unary()
			if err != nil {
				return nil, err
			}
			expr = Binary{Left: expr, Operator: operator, Right: right}
		}
	}

	return expr, nil
}

func (p *Parser) unary() (Expr, error) {

	if p.match(BANG, MINUS) {
		if operator, ok := p.previous(); ok {
			right, err := p.unary()
			if err != nil {
				return nil, err
			}
			return Unary{Operator: operator, Right: right}, nil
		}
	}

	return p.primary()
}

func (p *Parser) primary() (Expr, error) {

	if p.match(FALSE) {
		return Literal{Value: false}, nil
	}
	if p.match(TRUE) {
		return Literal{Value: true}, nil
	}
	if p.match(NIL) {
		return Literal{Value: nil}, nil
	}
	if p.match(NUMBER) {
		if token, ok := p.previous(); ok {
			value, err := strconv.ParseFloat(token.Literal, 64)
			if err != nil {
				return nil, err
			}
			return Literal{value}, nil
		}
	}
	if p.match(STRING) {
		if token, ok := p.previous(); ok {
			return Literal{token.Literal}, nil
		}
	}
	if p.match(IDENTIFIER) {
		if token, ok := p.previous(); ok {
			return Variable{token}, nil
		}
	}
	if p.match(LEFT_PAREN) {
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}
		if _, err := p.consume(RIGHT_PAREN, "Expect ')' after expression."); err != nil {
			return nil, err
		}

		return Grouping{expr}, nil
	}

	return nil, fmt.Errorf("error at line %d: unexpected token '%v'", p.peek().Line, p.peek().Lexeme)
}

func (p *Parser) match(tokenType ...TokenType) bool {
	for _, t := range tokenType {
		if p.check(t) {
			p.advance()
			return true
		}
	}

	return false
}

func (p *Parser) consume(tokenType TokenType, message string) (Token, error) {
	token := p.peek()

	if token.TType != tokenType {
		return Token{}, fmt.Errorf(message)
	}

	p.advance()

	return token, nil
}

func (p *Parser) previous() (Token, bool) {
	if p.current-1 < 0 {
		return Token{}, false
	}

	return p.tokens[p.current-1], true
}

func (p *Parser) check(t TokenType) bool {
	if p.isAtEnd() {
		return false
	}

	return p.peek().TType == t
}

func (p *Parser) peek() Token {
	return p.tokens[p.current]
}

func (p *Parser) advance() (Token, bool) {

	if p.isAtEnd() {
		return Token{}, false
	}
	p.current++

	return p.previous()
}

func (p *Parser) isAtEnd() bool {

	return p.peek().TType == EOF
}
