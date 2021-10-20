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
		literal: line,
		line:    line,
	}
}

func (t Token) String() string {
	return ""
}
