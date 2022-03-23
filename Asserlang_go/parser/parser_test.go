package parser

import (
	"fmt"
	"testing"

	"github.com/assertive-lang/asserlang/Asserlang_go/ast"
	"github.com/assertive-lang/asserlang/Asserlang_go/lexer"
)

func TestLetStatements(t *testing.T) {
	input := `쿠쿠루삥뽕
	어쩔냉장고~ㅋㅋㅌㅋㅋ
	어쩔인트~ㅋㅋㅋㅋ
	어쩔바보~ㅋㅋㅋ
	슉슈슉슉`
	l := lexer.New(input, false)
	p := New(l)
	program := p.ParseProgram()

	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got= %d", len(program.Statements))
	}

	tests := []struct {
		expectedIdentifier string
		expectedValue      string
	}{{"냉장고", "ㅋㅋㅌㅋㅋ"}, {"인트", "ㅋㅋㅋㅋ"}, {"바보", "ㅋㅋㅋ"}}
	for i, tt := range tests {
		stmt := program.Statements[i]
		if testLetStatement(t, stmt, tt.expectedIdentifier) {
			return
		}

		val := stmt.(*ast.LetStatement).Value

		if !testLiteralExpression(t, val, tt.expectedValue) {
			return
		}
	}
}

func testLetStatement(t *testing.T, s ast.Statement, name string) bool {
	t.Helper()

	if s.TokenLiteral() != "어쩔" {
		t.Errorf("stmt.TokenLiteral not '어쩔'. got=%q", s.TokenLiteral())
		return false
	}

	letStmt, ok := s.(*ast.LetStatement)
	if !ok {
		t.Errorf("s is not an *ast.LetStatement. got=%T", s)
		return false
	}

	if letStmt.Name.Value != name {
		t.Errorf("letStmt.Name.Value not '%s'.  got=%s", name, letStmt.Name.Value)
		return false
	}

	if letStmt.Name.TokenLiteral() != name {
		t.Errorf("letStmt.Name.TokenLiteral not '%s'. got=%s", name, letStmt.Name.TokenLiteral())
		return false
	}

	return true
}

func testLiteralExpression(t *testing.T, expr ast.Expression, expected interface{}) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, expr, int64(v))
	case int64:
		return testIntegerLiteral(t, expr, v)
	}
	t.Errorf("Type of expr not handled. Got: %T", expr)
	return false
}

func testIntegerLiteral(t *testing.T, intLiteral ast.Expression, value int64) bool {
	integer, ok := intLiteral.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("intLiteral not an *ast.IntegerLiteral. Got: %T", intLiteral)
		return false
	}
	if integer.Value != value {
		t.Errorf("integer.Value not %d. Got: %d", value, integer.Value)
		return false
	}
	if integer.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("integer.TokenLiteral() not %d. Got: %s", value, integer.TokenLiteral())
		return false
	}

	return true
}
