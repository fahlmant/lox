package lox

type Scanner struct {
	source  []rune
	tokens  []Token
	start   int
	current int
	line    int
}

func (s *Scanner) scanTokens() {
	// Scan until the end of the file
	for ok := true; ok; ok = !s.isAtEnd() {
		// Set start of a new lexeme
		s.start = s.current
		s.scanToken()
	}

	// Add EOF to the end of token list
	s.tokens = append(s.tokens, Token{tokenType: EOF, lexeme: "", literal: "", line: s.line})
}

// Check if all runes have been checked
func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

// Scan a single token
func (s *Scanner) scanToken() {

	// Get the next rune
	r := s.advance()
	switch r {
	case '(':
		s.addToken(LEFT_PAREN, "")
	case ')':
		s.addToken(RIGHT_PAREN, "")
	case '{':
		s.addToken(LEFT_BRACE, "")
	case '}':
		s.addToken(RIGHT_BRACE, "")
	case ',':
		s.addToken(COMMA, "")
	case '.':
		s.addToken(DOT, "")
	case '-':
		s.addToken(MINUS, "")
	case '+':
		s.addToken(PLUS, "")
	case ';':
		s.addToken(SEMICOLON, "")
	case '*':
		s.addToken(STAR, "")
	case '!':
		if s.match('=') {
			s.addToken(BANG_EQUAL, "")
		} else {
			s.addToken(BANG, "")
		}
	case '=':
		if s.match('=') {
			s.addToken(EQUAL_EQUAL, "")
		} else {
			s.addToken(EQUAL, "")
		}
	case '<':
		if s.match('=') {
			s.addToken(LESS_EQUAL, "")
		} else {
			s.addToken(LESS, "")
		}
	case '>':
		if s.match('=') {
			s.addToken(GREATER_EQUAL, "")
		} else {
			s.addToken(GREATER, "")
		}
	case '/':
		if s.match('/') {
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
		} else {
			s.addToken(SLASH, "")
		}
	// Ignore whitespace
	case ' ':
	case '\r':
	case '\t':
		break
	// Increase line count for each newline
	case '\n':
		s.line++
	// Handle strings encased in ""
	case '"':
		s.string()

	default:
		if isDigit(r) {
			s.number()
		} else if isAlpha(r) {
			s.identifier()
		} else {
			errorReport(s.line, "Unexpected error\n")
		}
	}

}

// Add a token to the Scanner's list
func (s *Scanner) addToken(tokenType TokenType, literal string) {

	// Get the textual representation of the token
	text := s.source[s.start:s.current]
	// Create token with tokentype, string, string literal provided and line number
	s.tokens = append(s.tokens, Token{tokenType: tokenType, lexeme: string(text), literal: literal, line: s.line})
}

// Consume the next rune
func (s *Scanner) advance() rune {
	// Increment to the next rune, but return the current rune.
	// Not sure why Go won't allow s.source[s.current++] instead
	s.current++
	return s.source[s.current-1]
}

// Check the following rune to see if its expected
// Used for multi-length tokens like != == <= >=
func (s *Scanner) match(expected rune) bool {

	if s.isAtEnd() {
		return false
	}

	if s.source[s.current] != expected {
		return false
	}

	s.current++
	return true
}

// Return the next rune without advancing and consuming
func (s *Scanner) peek() rune {

	// return null if at the end
	if s.isAtEnd() {
		return '\000'
	}

	return s.source[s.current]
}

// Returns two runes ahead without advancing and consuming
func (s *Scanner) peekNext() rune {
	// Returns null if at the end
	if s.current+1 >= len(s.source) {
		return '\000'
	}

	return s.source[s.current+1]
}

// Handle strings encased by ""
func (s *Scanner) string() {

	// Look for the closing "
	for s.peek() != '"' && !s.isAtEnd() {
		s.advance()
	}

	// If the end is reached, the string is not properly terminated
	if s.isAtEnd() {
		errorReport(s.line, "Unterminated string\n")
		return
	}

	// Handle the ending "
	s.advance()

	// Add the token with the literal, sans quotes
	s.addToken(STRING, string(s.source[s.start+1:s.current-1]))
}

// Handle number literals
func (s *Scanner) number() {

	// Keep going until there are no more numbers
	for isDigit(s.peek()) {
		s.advance()
	}

	// If there's a dot with numbers following, treat it as a decimal
	if s.peek() == '.' && isDigit(s.peekNext()) {
		// Keep going until there are no more numbers
		s.advance()
		for isDigit(s.peek()) {
			s.advance()
		}
	}

	// Add the token with the literal
	// Not sure when this will get converted to an actual number
	s.addToken(NUMBER, string(s.source[s.start:s.current]))
}

// Handle identifiers
func (s *Scanner) identifier() {

	// Keep going until there are no more characters
	for isAlphaNumeric(s.peek()) {
		s.advance()
	}
	text := s.source[s.start:s.current]

	if tokenType, ok := keywords[string(text)]; ok {
		s.addToken(tokenType, "")
	} else {
		s.addToken(IDENTIFIER, "")
	}
}

func isDigit(r rune) bool {
	// Check if rune is between values for 0 and 9
	return r >= '0' && r <= '9'
}

func isAlpha(r rune) bool {
	return (r >= 'a' && r <= 'z') ||
		r >= 'A' && r <= 'Z' ||
		r == '_'
}

func isAlphaNumeric(r rune) bool {
	return isAlpha(r) || isDigit(r)
}
