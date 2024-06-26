package parser

import (
	"GoInterpreter/src/ast"
	"GoInterpreter/src/lexer"
	"GoInterpreter/src/parser"
	"testing"
)

func TestLetStatements(t *testing.T) {
	input := `
	let x = 5;
	let y = 10;
	let foobar = 878787;
	`

	l := lexer.New(input)
	p := parser.New(l)

	program := p.Parse()
	checkParseErrors(t, p)
	if program == nil {
		t.Fatal("Parse() returned nil")
	}
	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain the required statements. Expected: %d, Got: %d", 3, len(program.Statements))
	}

	tests := []struct {
		expected string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	for i, tt := range tests {
		s := program.Statements[i]
		if testLetStatement(t, s, tt.expected) {
			return
		}
	}
}

func testLetStatement(t *testing.T, s ast.Statement, name string) bool {
	if s.Literal() != "let" {
		t.Errorf("Expected: let, Got: %q", s.Literal())
	}

	ls, ok := s.(*ast.Let)
	if !ok {
		t.Errorf("s not ast.LetStatement, Got: %d", s)
	}

	if ls.Name.Value != name {
		t.Errorf("Identifier mismatch, Expected: %s, Got: %s", name, ls.Name.Value)
	}

	if ls.Name.Literal() != name {
		t.Errorf("Name Literal mismatch, Expected: %s, Got: %s", name, ls.Name.Literal())
	}

	return true
}

func checkParseErrors(t *testing.T, p *parser.Parser) {
	errors := p.Errors()

	if len(errors) == 0 {
		return
	}

	t.Errorf("Parser Errors: %d", len(errors))
	for _, err := range errors {
		t.Errorf("Parsing Errors: %q", err)
	}

	t.FailNow()
}
