package parser

import (
	"GoInterpreter/src/ast"
	"GoInterpreter/src/lexer"
	"GoInterpreter/src/token"
	"fmt"
	"strconv"
)

type (
	Parser struct {
		lexer *lexer.Lexer

		current token.Token
		peek    token.Token

		prefix map[token.TokenType]prefixParseFn
		infix  map[token.TokenType]infixParseFn
		errors []string
	}

	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

const (
	_ = iota
	LOWEST
	EQUALS
	LESSGREATER
	SUM
	PRODUCT
	PREFIX
	CALL
)

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

func (p *Parser) next() {
	p.current = p.peek
	p.peek = p.lexer.NextToken()
}

func (p *Parser) parse() ast.Statement {
	switch p.current.Type {
	case token.LET:
		return p.letSt()
	case token.RETURN:
		return p.returnSt()
	default:
		return p.expressionSt()
	}
}

func (p *Parser) letSt() *ast.Let {
	st := &ast.Let{Token: p.current}
	if !p.expectPeek(token.IDENT) {
		return nil
	}

	st.Name = &ast.Identifier{Token: p.current, Value: p.current.Literal}
	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	for !p.currentTokenIs(token.SEMICOLON) {
		p.next()
	}

	return st
}

func (p *Parser) returnSt() *ast.Return {
	st := &ast.Return{Token: p.current}

	p.next()

	for !p.currentTokenIs(token.SEMICOLON) {
		p.next()
	}

	return st
}

func (p *Parser) expressionSt() *ast.ExpressionSt {
	st := &ast.ExpressionSt{
		Token:      p.current,
		Expression: p.expression(LOWEST),
	}

	if p.peekTokenIs(token.SEMICOLON) {
		p.next()
	}

	return st
}

func (p *Parser) integer() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.current}

	value, err := strconv.ParseInt(p.current.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("Error while parsing %q as an integer", p.current.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	lit.Value = value
	return lit
}

func (p *Parser) expression(precedence int) ast.Expression {
	prefix := p.prefix[p.current.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.current.Type)
		return nil
	}

	left := prefix()
	return left
}

func (p *Parser) prefixfn() ast.Expression {
	expression := &ast.PrefixEpression{
		Token:    p.current,
		Operator: p.current.Literal,
	}

	p.next()

	expression.Right = p.expression(PREFIX)
	return expression
}

func (p *Parser) identifier() ast.Expression {
	return &ast.Identifier{Token: p.current, Value: p.current.Literal}
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

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("No prefix func found for token type %s", t)
	p.errors = append(p.errors, msg)
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(t token.TokenType) {
	err := "Expected next token to be " + t + " , got " + p.peek.Type
	p.errors = append(p.errors, string(err))
}

func (p *Parser) registerInfix(t token.TokenType, fn infixParseFn) {
	p.infix[t] = fn
}

func (p *Parser) registerPrefix(t token.TokenType, fn prefixParseFn) {
	p.prefix[t] = fn
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		lexer:  l,
		errors: []string{},
	}

	p.prefix = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.IDENT, p.identifier)
	p.registerPrefix(token.INT, p.integer)
	p.registerPrefix(token.BANG, p.prefixfn)
	p.registerPrefix(token.MINUS, p.prefixfn)

	p.next()
	p.next()

	return p
}
