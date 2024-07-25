package ast

import (
	"GoInterpreter/src/token"
	"bytes"
)

type (
	Node interface {
		Literal() string
		String() string
	}

	Statement interface {
		Node
		statementNode()
	}

	Expression interface {
		Node
		expressionNode()
	}

	Program struct {
		Statements []Statement
	}

	Identifier struct {
		Token token.Token
		Value string
	}

	// Statements
	Let struct {
		Token token.Token
		Name  *Identifier
		Value Expression
	}

	Return struct {
		Token token.Token
		Value Expression
	}

	ExpressionSt struct {
		Token      token.Token
		Expression Expression
	}

	// Expressions
	IntegerLiteral struct {
		Token token.Token
		Value int64
	}

	PrefixEpression struct {
		Token    token.Token
		Operator string
		Right    Expression
	}
)

func (p *Program) Literal() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].Literal()
	} else {
		return "NoLiterals"
	}
}

func (p *Program) String() string {
	var out bytes.Buffer

	for _, s := range p.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

func (ls *Let) statementNode()  {}
func (ls *Let) Literal() string { return ls.Token.Literal }
func (ls *Let) String() string {
	var out bytes.Buffer

	out.WriteString(ls.Literal() + " ")
	out.WriteString(ls.Name.String())
	out.WriteString(" = ")

	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}

	out.WriteString(";")

	return out.String()
}

func (i *Identifier) expressionNode() {}
func (i *Identifier) Literal() string { return i.Token.Literal }
func (i *Identifier) String() string  { return i.Value }

func (rs *Return) statementNode()  {}
func (rs *Return) Literal() string { return rs.Token.Literal }
func (rs *Return) String() string {
	var out bytes.Buffer

	out.WriteString(rs.Literal() + " ")

	if rs.Value != nil {
		out.WriteString(rs.Value.String())
	}

	out.WriteString(";")

	return out.String()
}

func (es *ExpressionSt) statementNode()  {}
func (es *ExpressionSt) Literal() string { return es.Token.Literal }
func (es *ExpressionSt) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}

	return ""
}

func (il *IntegerLiteral) expressionNode() {}
func (il *IntegerLiteral) Literal() string { return il.Token.Literal }
func (il *IntegerLiteral) String() string  { return il.Token.Literal }

func (pe *PrefixEpression) expressionNode() {}
func (pe *PrefixEpression) Literal() string { return pe.Token.Literal }
func (pe *PrefixEpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")

	return out.String()
}
