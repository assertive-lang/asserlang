package parser

import (
	"fmt"
	"os"

	"github.com/assertive-lang/asserlang/Asserlang_go/ast"
	"github.com/assertive-lang/asserlang/Asserlang_go/lexer"
	"github.com/assertive-lang/asserlang/Asserlang_go/object"
	"github.com/assertive-lang/asserlang/Asserlang_go/token"
)

type (
	prefixParseFunc  func() ast.Expression
	infixParseFunc   func(ast.Expression) ast.Expression
	postfixParseFunc func() ast.Expression
)

type Parser struct {
	l *lexer.Lexer

	errors    []string
	curToken  token.Token
	peekToken token.Token

	prefixParseFuncs map[token.TokenType]prefixParseFunc

	infixParseFuncs   map[token.TokenType]infixParseFunc
	postfixParseFuncs map[token.TokenType]postfixParseFunc
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:                 l,
		errors:            []string{},
		prefixParseFuncs:  make(map[token.TokenType]prefixParseFunc),
		infixParseFuncs:   make(map[token.TokenType]infixParseFunc),
		postfixParseFuncs: make(map[token.TokenType]postfixParseFunc),
	}

	p.registerPrefix(token.KI, p.parseIntegerLiteral)
	p.registerPrefix(token.HU, p.parseIntegerLiteral)
	p.registerPrefix(token.IDENT, p.parseIdentitifier)
	p.registerPrefix(token.EOF, func() ast.Expression { os.Exit(0); return nil })
	p.registerPrefix(token.ANGUNG, p.parseCallExpression)
	p.registerPrefix(token.ANMUL, p.parseFunctionLiteral)
	p.registerPrefix(token.IF, p.parseIFExpression)

	p.registerInfix(token.KI, p.parseInfixIntegerExpression)
	p.registerInfix(token.HU, p.parseInfixIntegerExpression)
	p.registerInfix(token.TU, p.parseTUExpression)

	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()

	//fmt.Printf("'%v': %v\n", strings.Replace(p.curToken.Literal, "\n", "\\n", 1), p.curToken.Type)
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for p.peekToken.Type != token.EOF {
		stmt := p.parseStatement()

		if stmt != nil {
			program.Statements = append(program.Statements, stmt)

		}

		p.nextToken()
	}

	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.NEWLINE:
		return nil
	case token.BOF:
		return nil
	default:
		return p.parseExprStatement()
	}
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.curToken}

	if !p.expectPeek(token.IDENT) {
		p.errors = append(p.errors, fmt.Sprintf("어쩔변수 at line %d", p.peekToken.Line))
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(token.WAVE) {
		return nil
	}

	p.nextToken()

	stmt.Value = p.parseExpr()

	if fl, ok := stmt.Value.(*ast.FunctionLiteral); ok {
		fl.Name = stmt.Name.Value
	}

	return stmt
}

func (p *Parser) parseIFExpression() ast.Expression {
	expr := &ast.IfExpression{
		Token: p.curToken,
	}
	if !p.expectPeek(token.QUESTION) {
		return nil
	}
	p.nextToken()
	cond := p.parseExpr()

	if cond == nil {
		return nil
	}
	expr.Condition = cond

	if !p.expectPeek(token.THEN) {
		return nil
	}

	if !p.expectPeek(token.QUESTION) {
		return nil
	}
	expr.Consequence = p.parseBlockStatement(token.NEWLINE)

	return expr

}

func (p *Parser) parseFunctionLiteral() ast.Expression {
	lit := &ast.FunctionLiteral{Token: p.peekToken, Name: p.peekToken.Literal}

	if !p.peekTokenIs(token.IDENT) {
		return nil
	}
	p.nextToken()
	p.nextToken()

	lit.Parameters = p.parseFunctionParams()
	p.nextToken()
	lit.Body = p.parseBlockStatement(token.ANMUL)

	return lit
}

func (p *Parser) parseFunctionParams() []*ast.Identifier {
	identifiers := []*ast.Identifier{}

	if p.peekTokenIs(token.NEWLINE) {
		p.nextToken()
		return identifiers
	}
	p.nextToken()

	ident := &ast.Identifier{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}
	identifiers = append(identifiers, ident)

	for p.peekTokenIs(token.WAVE) {
		p.nextToken()
		p.nextToken()

		ident := &ast.Identifier{
			Token: p.curToken,
			Value: p.curToken.Literal,
		}
		identifiers = append(identifiers, ident)
	}

	return identifiers
}

func (p *Parser) parseBlockStatement(end token.TokenType) *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.curToken}
	block.Statements = []ast.Statement{}

	p.nextToken()

	for !p.curTokenIs(end) && !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}

		p.nextToken()

	}

	return block
}

func (p *Parser) parseExprStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}
	stmt.Expression = p.parseExpr()

	if p.peekTokenIs(token.NEWLINE) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpr() ast.Expression {
	prefix := p.prefixParseFuncs[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFuncError(p.curToken.Type)
		return nil
	}
	leftExpr := prefix()

	for !p.peekTokenIs(token.NEWLINE) {
		infix := p.infixParseFuncs[p.peekToken.Type]
		if infix == nil {
			return leftExpr
		}

		p.nextToken()
		leftExpr = infix(leftExpr)
	}
	return leftExpr
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.curToken}

	value := int64(0)

	if p.curTokenIs(token.KI) {
		value++
	} else if p.curTokenIs(token.HU) {
		value--
	}

	for p.peekTokenIs(token.KI) || p.peekTokenIs(token.HU) {
		switch p.peekToken.Type {
		case token.KI:
			value++
		case token.HU:
			value--
		}

		p.nextToken()
	}

	lit.Value = value

	return lit
}

func (p *Parser) parseInfixIntegerExpression(left ast.Expression) ast.Expression {

	expr := &ast.InfixIntegerExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}

	value := int64(0)

	if p.curTokenIs(token.KI) {
		value++
	} else if p.curTokenIs(token.HU) {
		value--
	}

	for p.peekTokenIs(token.KI) || p.peekTokenIs(token.HU) {
		switch p.peekToken.Type {
		case token.KI:
			value++
		case token.HU:
			value--
		}
		p.nextToken()
	}

	expr.Right = value

	return expr

}

func (p *Parser) parseTUExpression(left ast.Expression) ast.Expression {
	exp := &ast.TUExpression{Token: p.curToken, Left: left}

	p.nextToken()
	exp.Right = p.parseExpr()

	return exp
}

func (p *Parser) parseIdentitifier() ast.Expression {
	if object.GetBuiltinByName(p.curToken.Literal) != nil {

		expr := &ast.CallExpression{
			Token: p.curToken,
			Function: &ast.Identifier{
				Token: p.curToken,
				Value: p.curToken.Literal,
			},
		}
		if p.peekTokenIs(token.NEWLINE) {
			p.nextToken()
			expr.Arguments = []ast.Expression{}

			return expr
		}
		p.nextToken()
		val := p.parseExpr()

		expr.Arguments = []ast.Expression{val}

		return expr

	}
	exp := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	return exp
}

func (p *Parser) parseCallExpression() ast.Expression {
	expr := &ast.CallExpression{
		Token:    p.curToken,
		Function: &ast.Identifier{Token: p.peekToken, Value: p.peekToken.Literal},
	}
	expr.Arguments = p.parseExprList(token.NEWLINE)

	return expr
}

// TODO: fix functions stuff
func (p *Parser) parseExprList(end token.TokenType) []ast.Expression {
	list := []ast.Expression{}

	if p.peekTokenIs(end) {
		p.nextToken()
		return list
	}
	p.nextToken()

	for p.peekTokenIs(token.WAVE) {
		p.nextToken()
		p.nextToken()
		val := p.parseExpr()
		list = append(list, val)
	}
	p.nextToken()

	return list
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFunc) {
	p.prefixParseFuncs[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFunc) {
	p.infixParseFuncs[tokenType] = fn
}

func (p *Parser) registerPostfix(tokenType token.TokenType, fn postfixParseFunc) {
	p.postfixParseFuncs[tokenType] = fn
}

func (p *Parser) noPrefixParseFuncError(t token.TokenType) {
	msg := fmt.Sprintf("line %d: No prefix parse function for %s found", p.curToken.Line, t)
	p.errors = append(p.errors, msg)
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("line %d: expected next token to be %s, got %s instead", p.peekToken.Line, t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}
