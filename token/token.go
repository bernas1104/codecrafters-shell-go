package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

var escapeCharacters = make(map[string]string)

const (
	EOF = "EOF"

	DIRECTORY_PATH = "DIRECTORY"
	RELATIVE_PATH  = "RELATIVE"
	GO_BACK_PATH   = "GO_BACK"
	USER_PATH      = "USER_PATH"

	SPACE  = "SPACE"
	ESCAPE = "ESCAPE"

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
			current := tk.Literal

			if tk.Type == ESCAPE {
				current = escapeCharacter(tk.Literal)
			}

			arg += current
		} else {
			args = append(args, arg)
			arg = ""
		}
	}

	if arg != "" {
		args = append(args, arg)
	}

	for i := 0; i < len(args); i++ {
		if args[i][0] == '\'' || args[i][0] == '"' {
			args[i] = removeQuotes(args[i])
		}
	}

	return args
}

func escapeCharacter(arg string) string {
	ch, escaped := escapeCharacters[arg]

	if escaped {
		return ch
	}

	return ""
}

func removeQuotes(arg string) string {
	return arg[1 : len(arg)-1]
}

func InitializeEscapeCharacters() {
	escapeCharacters["\\ "] = " "
	escapeCharacters["\\\\"] = "\\"
	escapeCharacters["\\'"] = "'"
	escapeCharacters["\\\""] = "\""
	escapeCharacters["\\n"] = "n"
	escapeCharacters["\\t"] = "t"
	escapeCharacters["\\r"] = "r"
	escapeCharacters["\\b"] = "b"
	escapeCharacters["\\f"] = "f"
	escapeCharacters["\\v"] = "v"
}
