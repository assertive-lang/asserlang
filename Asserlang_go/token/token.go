package token

const (
	ILLEGAL = "ILLEGAL"
	BOF     = "BOF"
	EOF     = "EOF"
	NEWLINE = "NEWLINE"
	
	

	IDENT = "IDENT"

	WAVE = "WAVE"
	KI     = "KI"
	HU     = "HU"
	TU     = "TU"

	SEMICOLON = ";"

	FUNCTION = "FUNCTION"
	LET      = "LET"
)

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
	Line    int
}

var keywords = map[string]TokenType{
	"안물": FUNCTION,
	"어쩔": LET,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
