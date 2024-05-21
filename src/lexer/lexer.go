package lexer

import "GoInterpreter/src/token"

type Lexer struct {
	input    string
	position int
	next     int
	ch       byte
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.read()
	return l
}

func (l *Lexer) read() {
	if l.next >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.next]
	}

	l.position = l.next
	l.next += 1
}

func (l *Lexer) NextToken() token.Token {
	var t token.Token

	l.skip()

	switch l.ch {
	case '=':
		t = token.New(token.ASSIGN, l.ch)
	case ';':
		t = token.New(token.SEMICOLON, l.ch)
	case '(':
		t = token.New(token.LPAREN, l.ch)
	case ')':
		t = token.New(token.RPAREN, l.ch)
	case ',':
		t = token.New(token.COMMA, l.ch)
	case '+':
		t = token.New(token.PLUS, l.ch)
	case '{':
		t = token.New(token.LBRACE, l.ch)
	case '}':
		t = token.New(token.RBRACE, l.ch)
	case 0:
		t.Literal = ""
		t.Type = token.EOF
	default:
		if isLetter(l.ch) {
			t.Literal = l.identifier()
			t.Type = token.IdentifierLookUp(t.Literal)
			return t
		} else if isDigit(l.ch) {
			t.Type = token.INT
			t.Literal = l.number()
			return t
		} else {
			t = token.New(token.ILLEGAL, l.ch)
		}
	}

	l.read()
	return t
}

func (l *Lexer) identifier() string {
	p := l.position
	for isLetter(l.ch) {
		l.read()
	}

	return l.input[p:l.position]
}

func (l *Lexer) number() string {
	p := l.position
	for isDigit(l.ch) {
		l.read()
	}

	return l.input[p:l.position]
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func (l *Lexer) skip() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.read()
	}
}
