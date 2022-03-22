package ast

import "github.com/assertive-lang/asserlang/Asserlang_go/token"

// Identifier - holds IDENTIFIER token and it's value (add, foobar, x, y, ...)
type Identifier struct {
	Token token.Token // The token.Identifier token
	Value string
}

func (i *Identifier) expressionNode() {}

func (i *Identifier) TokenLiteral() string { return i.Token.Literal }

// String - returns a string representation of the Identifier and satisfies our Node interface
func (i *Identifier) String() string { return i.Value }
