package ast

import (
	"bytes"

	"github.com/assertive-lang/asserlang/Asserlang_go/token"
)

type TUExpression struct {
	Token token.Token
	Left  Expression
	Right Expression
}

func (se *TUExpression) expressionNode()      {}
func (se *TUExpression) TokenLiteral() string { return se.Token.Literal }
func (se *TUExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(se.Left.String())
	out.WriteString(" * " + se.Right.String())
	out.WriteString(")")

	return out.String()
}
