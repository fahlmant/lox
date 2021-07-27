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

	// If there's a for statement, handle it
	if p.match(FOR) {
		return p.forStatement()
	}

	// If there's an if statement, handle it
	if p.match(IF) {
		return p.ifStatement()
	}

	// If there's a print statement, handle it
	if p.match(PRINT) {
		return p.printStatement()
	}
	// If there's a while statement, handle it
	if p.match(WHILE) {
		return p.whileStatement()
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
	for !(p.peek().tType == RIGHT_BRACE) && !p.isAtEnd() {
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

func (p *Parser) forStatement() (Stmt, error) {

	var err error
	// Ensure there is a left parent after a "for"
	if _, err = p.consume(LEFT_PAREN, "Expect ( after a for"); err != nil {
		return nil, err
	}

	// Example
	// for (var i = 0; i < 5; i = i + 1)

	// Get the initalizer, i.e. the "var i = 0" part of th example
	// this can be a var declaration, an expression, or just empty
	var initializer Stmt
	if p.match(SEMICOLON) {
		initializer = nil
	} else if p.match(VAR) {
		initializer, err = p.varDeclaration()
		if err != nil {
			return nil, err
		}
	} else {
		initializer, err = p.expressionStatement()
		if err != nil {
			return nil, err
		}
	}

	// Get the conndition, i.e. "i < 5" from the example
	var condition Expr
	if p.peek().tType != SEMICOLON {
		condition, err = p.expression()
		if err != nil {
			return nil, err
		}
	}

	// Consume the semicolon after the condition
	if _, err := p.consume(SEMICOLON, "Expect ; after loop condition"); err != nil {
		return nil, err
	}

	// Get the increment, i.e "i = i + 1" from the example
	var increment Expr
	if p.peek().tType != RIGHT_PAREN {
		increment, err = p.expression()
		if err != nil {
			return nil, err
		}
	}

	// Consume the right paren after the condition
	if _, err := p.consume(RIGHT_PAREN, "Expect ) after for clauses"); err != nil {
		return nil, err
	}

	var body Stmt
	body, err = p.statement()
	if err != nil {
		return nil, err
	}

	// Create a new block with the body statement and the increment if there is one
	if increment != nil {
		body = BlockStmt{[]Stmt{body, ExprStmt{expression: increment}}}
	}

	// If there's no condition, default to true
	if condition == nil {
		condition = Literal{true}
	}

	// Create a while statement with the condition and the body
	// This is part of desugaring the for loop into a while loop
	body = WhileStmt{condition: condition, body: body}

	// If there's an initializer, build a block where that statement is executed before the
	// while loop
	if initializer != nil {
		body = BlockStmt{[]Stmt{initializer, body}}
	}

	return body, nil
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

func (p *Parser) whileStatement() (Stmt, error) {

	// Ensure there is a left paren after a "while"
	_, err := p.consume(LEFT_PAREN, "Expect ( after while")
	if err != nil {
		return nil, err
	}

	// Get the condition expression of the "while"
	condition, err := p.expression()
	if err != nil {
		return nil, err
	}

	// Ensure there is a left paren after the "while" condition
	if _, err := p.consume(RIGHT_PAREN, "Expect ) after a while"); err != nil {
		return nil, err
	}

	// Get the statement to do with the conditional
	body, err := p.statement()
	if err != nil {
		return nil, err
	}

	// Return a While statement with the condition and stmt body
	return WhileStmt{condition: condition, body: body}, nil

}

func (p *Parser) expression() (Expr, error) {
	// As of now, a passthrough in the recursive expansion
	return p.assignment()
}

func (p *Parser) assignment() (Expr, error) {

	// Goes down the recursive tree to get the expression
	expr, err := p.or()
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
				return Assign{variable: v, name: equals, value: value}, nil
			}

			return nil, fmt.Errorf("error at line %d: invalid assiment target", equals.line)
		}
	}

	// Returns the expression by itself if there's no assignment happening
	return expr, nil
}

func (p *Parser) or() (Expr, error) {

	// Get the left expression
	expr, err := p.and()
	if err != nil {
		return nil, err
	}

	//If there is an "or", handle it
	for p.match(OR) {
		// Tryo t get the operator and the right size
		if operator, ok := p.previous(); ok {
			right, err := p.and()
			if err != nil {
				return nil, err
			}

			// Return a Logical "left or right"
			return Logical{left: expr, operator: operator, right: right}, nil
		}

	}

	// Return the signular expression if there's no "or"
	return expr, nil
}

func (p *Parser) and() (Expr, error) {

	// Get the left expression
	expr, err := p.equality()
	if err != nil {
		return nil, err
	}

	// If there's an "and" handle it
	for p.match(AND) {
		// Try to get the operator and the right side
		if operator, ok := p.previous(); ok {
			right, err := p.equality()
			if err != nil {
				return nil, err
			}
			// Return a Logical "left and right"
			return Logical{left: expr, operator: operator, right: right}, nil
		}
	}

	// Return the singular expression if there's no "and"
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
			return Binary{left: expr, operator: operator, right: right}, nil
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
			expr = Binary{left: expr, operator: operator, right: right}
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
			expr = Binary{left: expr, operator: operator, right: right}

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
			expr = Binary{left: expr, operator: operator, right: right}
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
			return Unary{operator: operator, right: right}, nil
		}
	}

	return p.primary()
}

func (p *Parser) primary() (Expr, error) {

	if p.match(FALSE) {
		return Literal{value: false}, nil
	}
	if p.match(TRUE) {
		return Literal{value: true}, nil
	}
	if p.match(NIL) {
		return Literal{value: nil}, nil
	}
	if p.match(NUMBER) {
		if token, ok := p.previous(); ok {
			value, err := strconv.ParseFloat(token.literal, 64)
			if err != nil {
				return nil, err
			}
			return Literal{value}, nil
		}
	}
	if p.match(STRING) {
		if token, ok := p.previous(); ok {
			return Literal{token.literal}, nil
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

	return nil, fmt.Errorf("error at line %d: unexpected token '%v'", p.peek().line, p.peek().lexeme)
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

	if token.tType != tokenType {
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

	return p.peek().tType == t
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

	return p.peek().tType == EOF
}
