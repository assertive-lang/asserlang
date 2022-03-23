package lexer

import (
	"testing"

	"github.com/assertive-lang/asserlang/Asserlang_go/token"
)

func TestNextToken(t *testing.T) {
	input := `쿠쿠루삥뽕
	어쩔냉장고~ㅋㅋ
	슉슈슉슉`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.LET, "어쩔"},
		{token.IDENT, "냉장고"},
		{token.WAVE, "~"},
		{token.KI, "ㅋ"},
		{token.KI, "ㅋ"},
		{token.NEWLINE, "\n"},
		{token.EOF, "슉슈슉슉"},
	}

	l := New(input, false)

	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokenType wrong. expected=%q, got=%q", i, tt.expectedType, tok.Type)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
		}
	}
}
