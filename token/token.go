package token

type TokenType string

const (
	NotQuoted = iota
	SingleQuoted
	DoubleQuoted
)

type Token struct {
	Type       TokenType
	Literal    string
	QuotesType int
}

const (
	EOF = "EOF"

	DIRECTORY_PATH = "DIRECTORY"
	RELATIVE_PATH  = "RELATIVE"
	GO_BACK_PATH   = "GO_BACK"
	USER_PATH      = "USER_PATH"

	SPACE  = "SPACE"
	ESCAPE = "ESCAPE"
	STDOUT = "STDOUT"

	IDENTIFIER = "IDENT"
)

func NewToken(tokenType TokenType, literal string, quotesType int) Token {
	return Token{
		Type:       tokenType,
		Literal:    literal,
		QuotesType: quotesType,
	}
}
