package lox

import (
	"fmt"
	"strconv"
)

type Parser struct {
	tokens      []Token
	current     int
	expressions []Expr
	hadError    bool
}

func (p *Parser) parse() (Expr, error) {
	var expr Expr
	var err error
	for !p.isAtEnd() {
		expr, err = p.expression()
		if err != nil {
			p.hadError = true
			fmt.Printf("%v\n", err)
			return nil, err
		}
		p.expressions = append(p.expressions, expr)
	}

	return expr, nil
}

func (p *Parser) expression() (Expr, error) {
	return p.equality()
}

func (p *Parser) equality() (Expr, error) {
	expr, err := p.comparison()
	if err != nil {
		return nil, err
	}

	for p.match(EQUAL_EQUAL, BANG_EQUAL) {
		if operator, ok := p.previous(); ok {
			right, err := p.comparison()
			if err != nil {
				return nil, err
			}
			expr = Binary{Left: expr, Operator: operator, Right: right}
		}

	}

	return expr, nil
}

func (p *Parser) comparison() (Expr, error) {

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

	return nil, fmt.Errorf("error at line %d: unknown token '%v'", p.peek().Line, p.peek().Lexeme)
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
		return Token{}, fmt.Errorf("error at line %d: expected '%v'", token.Line, token.String())
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
