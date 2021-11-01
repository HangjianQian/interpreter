package core

import "fmt"

type Parser struct {
	tokens  []Token
	current int
}

func NewParser(t []Token) *Parser {
	return &Parser{
		tokens:  t,
		current: 0,
	}
}

// 返回全部的声明和语句，不进行 evaluate 操作
func (p *Parser) doParse() []Stmt {
	var stmts []Stmt
	for !p.isAtEnd() {
		stmts = append(stmts, p.declaration())
	}
	return stmts
}

func (p *Parser) declaration() Stmt {
	if p.match(VAR) {
		return p.varDeclaration()
	}
	return p.statement()
}

func (p *Parser) varDeclaration() Stmt {
	name := p.consume(IDENTIFIER, "expected identifier name")

	var initializer Expr
	if p.match(EQUAL) {
		initializer = p.expression()
	}
	p.consume(SEMICOLON, "Expect ';' after expression.")
	return VarStmt{name: name, initializer: initializer}
}

func (p *Parser) statement() Stmt {
	return p.expressionStmt()
}

func (p *Parser) expressionStmt() Stmt {
	ex := p.expression()
	p.consume(SEMICOLON, "Expect ';' after expression.")
	return ExprStmt{expr: ex}
}

func (p *Parser) expression() Expr {
	return p.assignment()
}

func (p *Parser) assignment() Expr {
	ex := p.equality()

	if p.match(EQUAL) {
		equals := p.previous()
		value := p.assignment()

		if _, ok := ex.(VarExpr); ok {
			name := ex.(VarExpr).name
			return AssignExpr{
				name:  name,
				value: value,
			}
		}
		panic(fmt.Sprintf("invalid assign target, %+v", equals))
	}
	return ex
}

func (p *Parser) equality() Expr {
	ex := p.comparison()
	for p.match(BANG_EQUAL, EQUAL_EQUAL) {
		op := p.previous()
		right := p.comparison()
		ex = BinaryExpr{
			left:     ex,
			right:    right,
			operator: op,
		}
	}
	return ex
}

func (p *Parser) comparison() Expr {
	ex := p.term()
	for p.match(GREATER, GREATER_EQUAL, LESS, LESS_EQUAL) {
		op := p.previous()
		right := p.term()
		ex = BinaryExpr{
			left:     ex,
			right:    right,
			operator: op,
		}
	}
	return ex
}

func (p *Parser) term() Expr {
	ex := p.factor()
	for p.match(PLUS, MINUS) {
		op := p.previous()
		right := p.factor()
		ex = BinaryExpr{
			left:     ex,
			right:    right,
			operator: op,
		}
	}
	return ex
}

func (p *Parser) factor() Expr {
	ex := p.unary()
	for p.match(SLASH, STAR) {
		op := p.previous()
		right := p.unary()
		ex = BinaryExpr{
			left:     ex,
			right:    right,
			operator: op,
		}
	}
	return ex
}

func (p *Parser) unary() Expr {
	if p.match(BANG, MINUS) {
		op := p.previous()
		ex := p.unary()
		return UnaryExpr{
			operator: op,
			right:    ex,
		}
	}
	return p.primary()
}

func (p *Parser) primary() Expr {
	if p.match(TRUE) {
		return LiteralExpr{obj: true}
	}
	if p.match(FALSE) {
		return LiteralExpr{obj: false}
	}
	if p.match(NIL) {
		return LiteralExpr{obj: nil}
	}
	if p.match(NUMBER, STRING) {
		return LiteralExpr{obj: p.previous().literal}
	}
	if p.match(IDENTIFIER) {
		return VarExpr{
			name: p.previous(),
		}
	}
	if p.match(LEFT_PAREN) {
		ex := p.expression()
		p.consume(RIGHT_PAREN, "EXPECT ')'")
		return GroupExpr{ex}
	}

	return nil
}

func (p *Parser) isAtEnd() bool {
	return p.peek().kind == EOF
}

func (p *Parser) advance() Token {
	if !p.isAtEnd() {
		p.current += 1
	}
	return p.previous()
}

func (p *Parser) peek() Token {
	return p.tokens[p.current]
}

func (p *Parser) previous() Token {
	return p.tokens[p.current-1]
}

func (p *Parser) match(tks ...TokenKind) bool {
	for _, tk := range tks {
		if p.check(tk) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) check(tk TokenKind) bool {
	return p.peek().kind == tk
}

func (p *Parser) consume(tk TokenKind, msg string) Token {
	if p.check(tk) {
		return p.advance()
	}
	panic(msg)
}
