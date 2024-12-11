package lexer

import "github.com/codecrafters-io/shell-starter-go/token"

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

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	switch l.ch {
	case '/':
		tok = token.NewToken(token.DIRECTORY_PATH, string(l.ch))
	case '.':
		if l.peekChar() == '/' {
			ch := l.ch
			l.readChar()

			literal := string(ch) + string(l.ch)
			tok = token.NewToken(token.RELATIVE_PATH, literal)
		} else if l.peekChar() == '.' {
			var chs []rune
			chs = append(chs, l.ch)
			l.readChar()

			if l.peekChar() == '/' {
				chs = append(chs, l.ch)
				l.readChar()

				literal := string(chs) + string(l.ch)
				tok = token.NewToken(token.GO_BACK_PATH, literal)
			} else {
				literal := l.readLiteral()
				literal = string(chs) + literal
				tok = token.NewToken(token.FILE_NAME, literal)
			}
		} else {
			literal := l.readLiteral()
			tok = token.NewToken(token.FILE_NAME, literal)
		}
	case 0:
		tok = token.NewToken(token.EOF, "")
	default:
		literal := l.readLiteral()
		tok = token.NewToken(token.FILE_NAME, literal)
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

func (l *Lexer) readLiteral() string {
	var chs []rune

	for l.peekChar() != '/' && l.peekChar() != 0 {
		ch := l.ch
		l.readChar()

		chs = append(chs, ch)
	}

	literal := string(chs) + string(l.ch)

	return literal
}
