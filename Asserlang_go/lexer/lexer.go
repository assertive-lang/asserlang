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

func New(input string) *Lexer {
	l := &Lexer{input: input, line: 2}

	l.readChar()
	return l
}

func (l *Lexer) readChar() {
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
	case "ㅇ":
		if l.peek() == "ㅉ" {
			l.readChar()

			tok = newToken(token.IDENT, "ㅇㅉ", l.line)
			l.readChar()
			return tok
		}

	case "우":
		if l.peek() == "짤" {
			l.readChar()
			if l.peek() == "래" {
				l.readChar()
				if l.peek() == "미" {
					l.readChar()
					tok = newToken(token.LET, "우짤래미", l.line)

				}
			}
		}

	case "저":
		if l.peek() == "쩔" {
			l.readChar()
			tok = newToken(token.LET, "저쩔", l.line)
		} else if l.peek() == "짤" {
			l.readChar()
			if l.peek() == "래" {
				l.readChar()
				if l.peek() == "미" {
					l.readChar()
					tok = newToken(token.LET, "저짤래미", l.line)
				}
			}
		}

	case "화":
		if l.peek() == "났" {
			l.readChar()
			if l.peek() == "쥬" {
				l.readChar()
				tok = newToken(token.IF, "화났쥬", l.line)
			}
		}

	case "안":
		if l.peek() == "물" {
			l.readChar()
			tok = newToken(token.ANMUL, "안물", l.line)
		} else if l.peek() == "궁" {
			l.readChar()
			tok = newToken(token.ANGUNG, "안궁", l.line)
		} else {
			tok.Line = l.line
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
		}

	case "ㅋ":
		tok = newToken(token.KI, "ㅋ", l.line)

	case "ㅎ":
		tok = newToken(token.HU, "ㅎ", l.line)

	case "ㅌ":
		if l.peek() == "ㅂ" {
			l.readChar()
			tok = newToken(token.IDENT, "ㅌㅂ", l.line)
		} else {
			tok = newToken(token.TU, "ㅌ", l.line)
		}

	case "?":
		tok = newToken(token.QUESTION, "?", l.line)

	case "~":
		tok = newToken(token.WAVE, "~", l.line)

	case "슉":
		if l.IsLast() {
			tok = newToken(token.EOF, "슉슈슉슉", l.line)
			return tok
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
		tok = newToken(token.NEWLINE, "\\n", l.line)
		l.line++

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

	for !l.isIdentifier() {
		l.readChar()
	}

	return string(l.input[position:l.position])
}

func (l *Lexer) isIdentifier() bool {
	internalReadPos := l.readPosition

	internalReadChar := func() {
		Byte1 := l.input[internalReadPos]

		if Byte1 >= 128 {
			l.ch = string([]byte{Byte1, l.input[internalReadPos+1], l.input[internalReadPos+2]})

			internalReadPos += 3
		} else {
			l.ch = string(Byte1)
			internalReadPos += 1
		}
	}

	internalPeek := func() string {
		if l.input[internalReadPos] >= 128 {
			return string([]byte{l.input[internalReadPos], l.input[internalReadPos+1], l.input[internalReadPos+2]})
		} else {
			return string(l.input[internalReadPos])
		}
	}

	switch l.ch {
	case "어":
		if internalPeek() == "쩔" {
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
	case "킹":
		if internalPeek() == "받" {
			internalReadChar()
			if internalPeek() == "쥬" {
				return true
			}

		}
	case "?":
		return true

	}
	return false
}

func (l *Lexer) peek() string {
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
