package lexer

import (
	"github.com/codecrafters-io/shell-starter-go/token"
	"github.com/codecrafters-io/shell-starter-go/utils"
)

var breakers = []string{"\\", " ", "/"}
var escapers = []string{"'", "\"", "\\", " ", "n"}
var exceptions = []string{"\\ ", "\\n", "\\'"}
var escapeCharactersMap = make(map[string]string)

type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           rune
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()

	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = rune(l.input[l.readPosition])
	}

	l.position = l.readPosition
	l.readPosition++
}

func GetTokens(arg string) []token.Token {
	l := New(arg)
	var tokens []token.Token

	for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
		tokens = append(tokens, tok)
	}

	return tokens
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	switch l.ch {
	case '/':
		tok = token.NewToken(token.DIRECTORY_PATH, string(l.ch), token.NotQuoted)
	case '.':
		if l.peekChar() == '/' {
			ch := l.ch
			l.readChar()

			literal := string(ch) + string(l.ch)
			tok = token.NewToken(token.RELATIVE_PATH, literal, token.NotQuoted)
		} else if l.peekChar() == '.' {
			ch := l.ch
			l.readChar()

			literal := string(ch) + string(l.ch)
			tok = token.NewToken(token.GO_BACK_PATH, literal, token.NotQuoted)
		} else {
			literal, quotesType := l.readLiteral()
			tok = token.NewToken(token.IDENTIFIER, literal, quotesType)
		}
	case '~':
		tok = token.NewToken(token.USER_PATH, string(l.ch), token.NotQuoted)
	case ' ':
		spaces := string(l.ch)

		for l.peekChar() == ' ' {
			l.readChar()
			spaces += string(l.ch)
		}

		tok = token.NewToken(token.SPACE, spaces, token.NotQuoted)
	case '\\':
		ch := l.ch

		if !utils.Contains(escapers, string(l.peekChar())) {
			tok = token.NewToken(token.ESCAPE, string(ch), token.NotQuoted)
		} else {
			l.readChar()
			tok = token.NewToken(token.ESCAPE, string(ch)+string(l.ch), token.NotQuoted)
		}
	case '>':
		tok = token.NewToken(token.STDOUT, ">", token.NotQuoted)
	case '1':
		ch := l.ch
		l.readChar()

		tok = token.NewToken(token.STDOUT, string(ch)+string(l.ch), token.NotQuoted)
	case 0:
		tok = token.NewToken(token.EOF, "", token.NotQuoted)
	default:
		literal, quotesType := l.readLiteral()
		tok = token.NewToken(token.IDENTIFIER, literal, quotesType)
	}

	l.readChar()

	return tok
}

func (l *Lexer) peekChar() rune {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return rune(l.input[l.readPosition])
	}
}

func (l *Lexer) previousNthChar(nth int) rune {
	return rune(l.input[l.readPosition-nth])
}

func (l *Lexer) readLiteral() (string, int) {
	var chs []rune
	var literal string
	var quotesType int

	if l.ch != '\'' && l.ch != '"' {
		for !utils.Contains(breakers, string(l.peekChar())) && l.peekChar() != 0 {
			ch := l.ch
			l.readChar()

			chs = append(chs, ch)
		}

		literal = string(chs) + string(l.ch)
		quotesType = token.NotQuoted
	} else {
		delimiter := '\''
		quotesType = token.SingleQuoted

		if l.ch != '\'' {
			delimiter = '"'
			quotesType = token.DoubleQuoted
		}

		ch := l.ch
		l.readChar()

		chs = append(chs, ch)

		for (l.previousNthChar(1) != delimiter || l.previousNthChar(2) == '\\') && l.peekChar() != 0 {
			ch := l.ch
			l.readChar()

			chs = append(chs, ch)
		}

		literal = string(chs) + string(l.ch)
	}

	return literal, quotesType
}

func GetCommand(tks []token.Token) string {
	if tks[0].Type != token.IDENTIFIER {
		return "cd"
	}

	return tks[0].Literal
}

func GetArguments(tks []token.Token) []string {
	var arg string
	var args []string
	var quotesType int

	if tks[0].Type != token.IDENTIFIER {
		args = append(args, "cd")
	}

	for i := 0; i < len(tks); i++ {
		if tks[i].Literal[0] == '\'' || tks[i].Literal[0] == '"' {
			tks[i].Literal = removeQuotes(tks[i].Literal)
		}
	}

	for _, tk := range tks {
		quotesType = tk.QuotesType

		if tk.Type != token.SPACE {
			var current string

			if tk.Type == token.ESCAPE {
				current = escapeCharacter(tk.Literal)
			} else {
				current = escapeCharacters(tk.Literal, quotesType)
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

	return args
}

func removeQuotes(arg string) string {
	return arg[1 : len(arg)-1]
}

func escapeCharacter(arg string) string {
	ch, escaped := escapeCharactersMap[arg]

	if escaped {
		return ch
	}

	return arg
}

func escapeCharacters(arg string, quotesType int) string {
	escaped := ""
	tks := GetTokens(arg)

	for _, tk := range tks {
		if tk.Type == token.ESCAPE && quotesType == token.DoubleQuoted {
			if utils.Contains(exceptions, tk.Literal) {
				escaped += tk.Literal
			} else {
				escaped += escapeCharacter(tk.Literal)
			}
		} else {
			escaped += tk.Literal
		}
	}

	return escaped
}

func InitializeEscapeCharacters() {
	escapeCharactersMap["\\\\"] = "\\"
	escapeCharactersMap["\\ "] = " "
	escapeCharactersMap["\\'"] = "'"
	escapeCharactersMap["\\\""] = "\""
	escapeCharactersMap["\\n"] = "n"
	escapeCharactersMap["\\t"] = "t"
	escapeCharactersMap["\\r"] = "r"
	escapeCharactersMap["\\b"] = "b"
	escapeCharactersMap["\\f"] = "f"
	escapeCharactersMap["\\v"] = "v"
}
