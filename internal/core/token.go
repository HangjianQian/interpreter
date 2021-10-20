package core

type Token struct {
	kind    TokenKind
	lexeme  string
	literal interface{}
	line    int
}

func NewToken(kind TokenKind, lexeme string, literal interface{}, line int) Token {
	return Token{
		kind:    kind,
		lexeme:  lexeme,
		literal: literal,
		line:    line,
	}
}
