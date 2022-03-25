package token

const (
	ILLEGAL = "ILLEGAL"
	BOF     = "BOF"
	EOF     = "EOF"
	NEWLINE = "NEWLINE"

	ANMUL  = "ANMUL"
	ANGUNG = "ANGUNG"

	IDENT = "IDENT"

	WAVE = "WAVE"
	KI   = "KI"
	HU   = "HU"
	TU   = "TU"

	JUMP     = ";;"
	QUESTION = "?"

	FUNCTION = "FUNCTION"
	LET      = "LET"
	IF       = "IF"
	THEN     = "THEN"
)

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
	Line    int
}

var keywords = map[string]TokenType{
	"안물":  FUNCTION,
	"어쩔":  LET,
	"화나쥬": IF,
	"킹받쥬": THEN,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
