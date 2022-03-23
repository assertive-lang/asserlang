package ast

import "github.com/assertive-lang/asserlang/Asserlang_go/token"

type IntegerLiteral struct {
	Token       token.Token
	Experession string
	Value       int64
}

func (ie *IntegerLiteral) expressionNode() {}

func (ke *IntegerLiteral) TokenLiteral() string { return ke.Token.Literal }

func (ke *IntegerLiteral) String() string { return ke.Token.Literal }
