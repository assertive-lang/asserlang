package ast

import "github.com/assertive-lang/asserlang/Asserlang_go/token"

type IntegerExpression struct {
	Token       token.Token
	Experession string
	Value       int
}

func (ie *IntegerExpression) expressionNode() {}

func (ie *IntegerExpression) TokenLiteral() string { return ie.Token.Literal }

func (ie *IntegerExpression) String() string { return ie.Token.Literal }
