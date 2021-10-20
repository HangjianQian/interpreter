package core

import (
	"os"
	"strconv"

	"github.com/sirupsen/logrus"
)

type scanner struct {
	tokens []Token
	source string
	length int

	start   int
	current int
	line    int
}

func NewScanner(source string) *scanner {
	return &scanner{
		source: source,
		line:   1,
		length: len(source),
	}
}

func (s *scanner) scanTokens() {
	for !s.isAtEnd() {
		s.start = s.current
		s.scanToken()
	}

	s.tokens = append(s.tokens, NewToken(EOF, "", nil, s.line))
	return
}

func (s *scanner) scanToken() {
	c := s.advance()
	switch c {
	case '(':
		s.addToken(LEFT_PAREN)
	case ')':
		s.addToken(RIGHT_PAREN)
	case '{':
		s.addToken(LEFT_BRACE)
	case '}':
		s.addToken(RIGHT_BRACE)
	case ',':
		s.addToken(COMMA)
	case '.':
		s.addToken(DOT)
	case '-':
		s.addToken(MINUS)
	case '+':
		s.addToken(PLUS)
	case ';':
		s.addToken(SEMICOLON)
	case '*':
		s.addToken(STAR)
	case '!':
		if s.match("=") {
			s.addToken(BANG_EQUAL)
		} else {
			s.addToken(BANG)
		}
	case '=':
		if s.match("=") {
			s.addToken(EQUAL_EQUAL)
		} else {
			s.addToken(EQUAL)
		}
	case '<':
		if s.match("=") {
			s.addToken(LESS_EQUAL)
		} else {
			s.addToken(LESS)
		}
	case '>':
		if s.match("=") {
			s.addToken(GREATER_EQUAL)
		} else {
			s.addToken(GREATER)
		}
	case '/':
		if s.match("/") {
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
		} else {
			s.addToken(SLASH)
		}
	case ' ', '\r', '\t':
	case '\n':
		s.line += 1
	case '"':
		str := s.string()
		s.addTokenValue(STRING, str)
	default:
		if s.isDigit(c) {
			v := s.number()
			s.addTokenValue(NUMBER, v)
		} else if s.isAlpha(c) {
			v := s.identifier()
			s.tokens = append(s.tokens, v)
		} else {
			logrus.Errorf("invalid char: %s at line %d", c, s.line)
		}
	}
}

func (s *scanner) isAtEnd() bool {
	return s.current >= s.length
}

func (s *scanner) advance() rune {
	if s.isAtEnd() {
		return '0'
	}
	emit := rune(s.source[s.current])
	s.current += 1
	return emit
}

func (s *scanner) match(c string) bool {
	if s.isAtEnd() {
		return false
	}

	if s.source[s.current:s.current+1] != c {
		return false
	}
	s.current += 1
	return true
}

func (s *scanner) peek() rune {
	if s.isAtEnd() {
		return 0
	}
	return rune(s.source[s.current])
}

func (s *scanner) peekNext() rune {
	if s.current+1 >= s.length {
		return 0
	}
	return rune(s.source[s.current])
}

func (s *scanner) string() string {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line += 1
		}
		s.advance()
	}

	if s.isAtEnd() {
		logrus.Errorf("unterminate string at line %d\n", s.line)
		os.Exit(-1)
	}
	s.advance()

	return s.source[s.start+1 : s.current-1]
}

func (s *scanner) isDigit(c rune) bool {
	return c >= '0' && c <= '9'

}

func (s *scanner) isAlpha(c rune) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_'
}

func (s *scanner) number() interface{} {
	for s.isDigit(s.peek()) {
		s.advance()
	}

	if s.peek() == '.' && s.isDigit(s.peekNext()) {
		s.advance()
		for s.isDigit(s.peek()) {
			s.advance()
		}
	}

	v, err := strconv.ParseFloat(s.source[s.start:s.current], 64)
	if err != nil {
		logrus.Errorf("invalid number %s at line %d\n", s.source[s.start:s.current], s.line)
		os.Exit(-1)
	}
	return v
}

func (s *scanner) identifier() Token {
	for s.isAlpha(s.peek()) {
		s.advance()
	}

	str := s.source[s.start : s.current+1]

	// keywords ?
	if v, ok := KEYWORDS[str]; ok {
		return NewToken(v, "", nil, s.line)
	}
	return NewToken(IDENTIFIER, s.source[s.start:s.current], nil, s.line)
}

func (s *scanner) addToken(kind TokenKind) {
	s.tokens = append(s.tokens, NewToken(kind, s.source[s.start:s.current], nil, s.line))
}

func (s *scanner) addTokenValue(kind TokenKind, v interface{}) {
	s.tokens = append(s.tokens, NewToken(kind, "", v, s.line))
}
