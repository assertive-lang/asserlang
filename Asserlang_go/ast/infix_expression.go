package ast

import (
	"bytes"
	"fmt"

	"github.com/assertive-lang/asserlang/Asserlang_go/token"
)

type InfixIntegerExpression struct {
	Token    token.Token // The operator token (+, -, *, etc)
	Left     Expression
	Operator string // string (examples: "+", "-", "*", etc)
	Right    int64
}

func (ie *InfixIntegerExpression) expressionNode() {}

func (ie *InfixIntegerExpression) TokenLiteral() string { return ie.Token.Literal }

func (ie *InfixIntegerExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" " + ie.Operator + " ")
	out.WriteString(fmt.Sprintf("%d", ie.Right))
	out.WriteString(")")

	return out.String()
}
