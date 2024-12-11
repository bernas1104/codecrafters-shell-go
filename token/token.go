package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	EOF = "EOF"

	DIRECTORY_PATH = "DIRECTORY"
	RELATIVE_PATH  = "RELATIVE"
	GO_BACK_PATH   = "GO_BACK"
	USER_PATH      = "USER_PATH"

	FILE_NAME = "FILE"
)

func NewToken(tokenType TokenType, literal string) Token {
	return Token{
		Type:    tokenType,
		Literal: literal,
	}
}
