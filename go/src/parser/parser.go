package parser

import (
	"GoInterpreter/src/ast"
	"GoInterpreter/src/lexer"
	"GoInterpreter/src/token"
)

type (
	Parser struct {
		l *lexer.Lexer

		current token.Token
		peek    token.Token
		errors  []string
	}
)

func (p *Parser) next() {
	p.current = p.peek
	p.peek = p.l.NextToken()
}

func (p *Parser) Parse() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for p.current.Type != token.EOF {
		s := p.parse()
		if s != nil {
			program.Statements = append(program.Statements, s)
		}

		p.next()
	}

	return program
}

func (p *Parser) parse() ast.Statement {
	switch p.current.Type {
	case token.LET:
		return p.let()
	default:
		return nil
	}
}

func (p *Parser) let() *ast.Let {
	s := &ast.Let{Token: p.current}
	if !p.expectPeek(token.IDENT) {
		return nil
	}

	s.Name = &ast.Identifier{Token: p.current, Value: p.current.Literal}
	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	for !p.currentTokenIs(token.SEMICOLON) {
		p.next()
	}

	return s
}

func (p *Parser) currentTokenIs(t token.TokenType) bool {
	return p.current.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peek.Type == t
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.next()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(t token.TokenType) {
	err := "Expected next token to be " + t + " , got " + p.peek.Type
	p.errors = append(p.errors, string(err))
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	p.next()
	p.next()

	return p
}
