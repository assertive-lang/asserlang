package parser

import (
	"testing"

	"github.com/assertive-lang/asserlang/Asserlang_go/ast"
	"github.com/assertive-lang/asserlang/Asserlang_go/lexer"
)

func TestLetStatements(t *testing.T) {
	input := `쿠쿠루삥뽕
	어쩔냉장고~ㅋㅋ
	어쩔인트~ㅋㅋㅋㅋ
	어쩔바보~ㅋㅋㅋ
	슉슈슉슉`
	l := lexer.New(input)
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
	}{{"냉장고"}, {"인트"}, {"바보"}}
	for i, tt := range tests {
		stmt := program.Statements[i]
		if testLetStatement(t, stmt, tt.expectedIdentifier) {
			return
		}
	}
}

func testLetStatement(t *testing.T, s ast.Statement, name string) bool {
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
