package lexer

import (
	"github.com/assertive-lang/asserlang/Asserlang_go/token"
)

type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           string
	line         int
}

func New(input string, console bool) *Lexer {
	l := &Lexer{input: input, line: 1}

	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = ""

		l.position = l.readPosition
		l.readPosition += 1
		return
	}
	Byte1 := l.input[l.readPosition]

	if Byte1 >= 128 {
		l.ch = string([]byte{Byte1, l.input[l.readPosition+1], l.input[l.readPosition+2]})

		l.position = l.readPosition
		l.readPosition += 3
	} else {
		l.ch = string(Byte1)
		l.position = l.readPosition
		l.readPosition += 1
	}
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	switch l.ch {
	case "어":
		if l.peek() == "쩔" {
			l.readChar()
			tok = newToken(token.LET, "어쩔", l.line)
		}
	case "저":
		if l.peek() == "쩔" {
			l.readChar()
			tok = newToken(token.LET, "저쩔", l.line)
		}
	case "ㅋ":
		tok = newToken(token.KI, "ㅋ", l.line)
	case "~":
		tok = newToken(token.WAVE, "~", l.line)
	case "ㅎ":
		tok = newToken(token.HU, "ㅎ", l.line)
	case "ㅌ":
		tok = newToken(token.TU, "ㅌ", l.line)
	case "슉":
		if l.IsLast() {
			tok = newToken(token.EOF, "슉슈슉슉", l.line)

		}
	case "쿠":
		l.readChar()
		if l.ch == "쿠" {
			l.readChar()
			if l.ch == "루" {
				l.readChar()
				if l.ch == "삥" {
					l.readChar()
					if l.ch == "뽕" {
						l.readChar()
						tok = newToken(token.BOF, "쿠쿠루삥뽕", l.line)

					}
				}
			}
		}
	case "\n":
		tok = newToken(token.NEWLINE, "\n", l.line)
		l.line++
	case "":
		tok.Literal = ""
		tok.Type = token.EOF
		tok.Line = l.line

	default:
		tok.Line = l.line
		tok.Literal = l.readIdentifier()
		tok.Type = token.LookupIdent(tok.Literal)
		return tok

	}
	l.readChar()
	return tok
}

func (l *Lexer) readIdentifier() string {
	position := l.position

	for !isIdentifier(l.ch, l.peek()) {
		l.readChar()
	}

	return string(l.input[position:l.position])
}

func isIdentifier(tok string, peek string) bool {
	switch tok {
	case "어":
		if peek == "쩔" {
			return true
		}
	case "ㅋ":
		return true
	case "~":
		return true
	case "ㅎ":
		return true
	case "ㅌ":
		return true
	case "\n":
		return true
	case " ":
		return true
	}
	return false
}

func (l *Lexer) peek() string {
	if l.readPosition >= len(l.input) {
		return ""
	}
	if l.input[l.readPosition] >= 128 {
		return string([]byte{l.input[l.readPosition], l.input[l.readPosition+1], l.input[l.readPosition+2]})
	} else {
		return string(l.input[l.readPosition])
	}
}

func (l *Lexer) skipWhitespace() {
	for l.ch == " " || l.ch == "\t" || l.ch == "\v" || l.ch == "\r" {
		l.readChar()
	}
}

func (l *Lexer) IsLast() bool {
	return l.input[l.position:l.position+12] == "슉슈슉슉"
}

func newToken(tokenType token.TokenType, ch string, line int) token.Token {
	return token.Token{Type: tokenType, Literal: ch, Line: line}
}
