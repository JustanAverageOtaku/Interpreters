package ast

import "GoInterpreter/src/token"

type (
	Node interface {
		Literal() string
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

	Let struct {
		Token token.Token
		Name  *Identifier
		Value Expression
	}
)

func (p *Program) Literal() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].Literal()
	} else {
		return "NoLiterals"
	}
}

func (ls *Let) statementNode()  {}
func (ls *Let) Literal() string { return ls.Token.Literal }

func (i *Identifier) expressionNode() {}
func (i *Identifier) Literal() string { return i.Token.Literal }
