package ast

import (
	"bytes"

	"github.com/assertive-lang/asserlang/Asserlang_go/token"
)

type IfExpression struct {
	Token       token.Token // The If token
	Condition   Expression
	Consequence *BlockStatement
}

func (ie *IfExpression) expressionNode() {}

func (ie *IfExpression) TokenLiteral() string { return ie.Token.Literal }

func (ie *IfExpression) String() string {
	var out bytes.Buffer

	out.WriteString("if ")
	out.WriteString(ie.Condition.String())
out.WriteString(" ")
	out.WriteString(ie.Consequence.String())


	return out.String()
}
