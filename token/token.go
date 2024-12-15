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

	SPACE = "SPACE"

	IDENTIFIER = "IDENT"
)

func NewToken(tokenType TokenType, literal string) Token {
	return Token{
		Type:    tokenType,
		Literal: literal,
	}
}

func GetCommand(tks []Token) string {
	if tks[0].Type != IDENTIFIER {
		return "cd"
	}

	return tks[0].Literal
}

func GetArguments(tks []Token) []string {
	var arg string
	var args []string

	for _, tk := range tks {
		if tk.Type != SPACE {
			arg += tk.Literal
		} else {
			args = append(args, arg)
			arg = ""
		}
	}

	if arg != "" {
		args = append(args, arg)
	}

	return args
}
