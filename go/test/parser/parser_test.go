package parser

import (
	"GoInterpreter/src/ast"
	"GoInterpreter/src/lexer"
	"GoInterpreter/src/parser"
	"GoInterpreter/src/token"
	"fmt"
	"testing"
)

func TestParsingPrefixFunc(t *testing.T) {
	es := 1
	prefixTests := []struct {
		input    string
		operator string
		value    int64
	}{
		{"!5", "!", 5},
		{"-15", "-", 15},
	}

	for _, tt := range prefixTests {
		lexer := lexer.New(tt.input)
		parser := parser.New(lexer)

		tree := parser.Parse()
		checkParseErrors(t, parser)

		if len(tree.Statements) != es {
			t.Fatalf("Incorrect number of statements parsed. Expected:%d Got:%d", es, len(tree.Statements))
		}

		stmt, ok := tree.Statements[0].(*ast.ExpressionSt)
		if !ok {
			t.Fatalf("The parsed statement is not an epression statement. Got:%T", tree.Statements[0])
		}

		expr, ok := stmt.Expression.(*ast.PrefixEpression)
		if !ok {
			t.Fatalf("The parsed expression is not a prefix expression. Got:%T", stmt.Expression)
		}

		if expr.Operator != tt.operator {
			t.Fatalf("The parsed operator does not match the expected one. Expected:%s Got:%s", tt.operator, expr.Operator)
		}

		if !testIntegerLiteral(t, expr.Right, tt.value) {
			return
		}
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := "expression"

	lexer := lexer.New(input)
	parser := parser.New(lexer)

	tree := parser.Parse()
	checkParseErrors(t, parser)

	if len(tree.Statements) != 1 {
		t.Fatalf("Incorrect number of statements parsed. Got:%d", len(tree.Statements))
	}

	stmt, ok := tree.Statements[0].(*ast.ExpressionSt)
	if !ok {
		t.Fatalf("The parsed statement is not an expression statemet. Got:%T", tree.Statements[0])
	}

	identifier, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("The parsed statemet is not a valid identifier. Got:%T", stmt.Expression)
	}

	if identifier.Value != "expression" {
		t.Errorf("The parsed value for the identifier is incorrect. Got:%s", identifier.Value)
	}

	if identifier.Literal() != "expression" {
		t.Errorf("The parsed literal for the identifier is incorrect. Got:%s", identifier.Literal())
	}
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "5;"

	lexer := lexer.New(input)
	parser := parser.New(lexer)

	tree := parser.Parse()
	checkParseErrors(t, parser)

	if len(tree.Statements) != 1 {
		t.Fatalf("Incorrect number of statements parsed. Got:%d", len(tree.Statements))
	}

	stmt, ok := tree.Statements[0].(*ast.ExpressionSt)
	if !ok {
		t.Fatalf("The parsed statement is not an expression statemet. Got:%T", tree.Statements[0])
	}

	identifier, ok := stmt.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("The parsed statemet is not a valid integer. Got:%T", stmt.Expression)
	}

	if identifier.Value != 5 {
		t.Errorf("The parsed value for the identifier is incorrect. Got:%d", identifier.Value)
	}

	if identifier.Literal() != "5" {
		t.Errorf("The parsed literal for the identifier is incorrect. Got:%s", identifier.Literal())
	}
}

func TestString(t *testing.T) {
	p := &ast.Program{
		Statements: []ast.Statement{
			&ast.Let{
				Token: token.Token{Type: token.LET, Literal: "let"},
				Name: &ast.Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "foo"},
					Value: "foo",
				},
				Value: &ast.Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "bar"},
					Value: "bar",
				},
			},
		},
	}

	if p.String() != "let foo = bar;" {
		t.Errorf("String Mistmatch. Got: %s", p.String())
	}
}

func TestReturnStatements(t *testing.T) {
	input := `
   	return 5;
   	return 10;
   	return 993322;
   	`

	l := lexer.New(input)
	p := parser.New(l)

	program := p.Parse()
	checkParseErrors(t, p)

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain the required statements. Expected: %d, Got: %d", 3, len(program.Statements))
	}

	for _, s := range program.Statements {
		rs, ok := s.(*ast.Return)
		if !ok {
			t.Errorf("Starement is not an ast return statement. Got: %T", s)
			continue
		}

		if rs.Token.Literal != "return" {
			t.Errorf("The literal for return statement is not 'return'. Got: %s", rs.Token.Literal)
		}
	}
}

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

func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	lit, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("The parsed expression is not a valid integer literal. Got:%T", il)
		return false
	}

	if lit.Value != value {
		t.Errorf("The parsed value does not equal the expected one. Expected:%d Got:%d", value, lit.Value)
		return false
	}

	if lit.Literal() != fmt.Sprintf("%d", value) {
		t.Errorf("Token mismatch. Expected:%d Got:%s", value, lit.Literal())
		return false
	}

	return true
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
