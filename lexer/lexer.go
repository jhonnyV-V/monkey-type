package lexer

import (
	"mokey-type/token"
)

type Lexer struct {
	input        string
	position     int // current position in input (current char)
	readPosition int // current reading position in input (after current char)
	ch           byte
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.ReadChar()
	return l
}

func (l *Lexer) ReadChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.ReadChar()
			tok = token.Token{Type: token.EQ, Literal: string(ch) + string(l.ch)}
		} else {
			tok = token.NewToken(token.ASSIGN, l.ch)
		}
	case ';':
		tok = token.NewToken(token.SEMICOLON, l.ch)
	case '(':
		tok = token.NewToken(token.LPAREN, l.ch)
	case ')':
		tok = token.NewToken(token.RPAREN, l.ch)
	case ',':
		tok = token.NewToken(token.COMMA, l.ch)
	case '+':
		if l.peekChar() == '+' {
			ch := l.ch
			l.ReadChar()
			tok = token.Token{Type: token.INCREMENT, Literal: string(ch) + string(l.ch)}
		} else {
			tok = token.NewToken(token.PLUS, l.ch)
		}
	case '-':
		if l.peekChar() == '-' {
			ch := l.ch
			l.ReadChar()
			tok = token.Token{Type: token.DECREMENT, Literal: string(ch) + string(l.ch)}
		} else {
			tok = token.NewToken(token.MINUS, l.ch)
		}
	case '{':
		tok = token.NewToken(token.LBRACE, l.ch)
	case '}':
		tok = token.NewToken(token.RBRACE, l.ch)
	case '[':
		tok = token.NewToken(token.LBRACKET, l.ch)
	case ']':
		tok = token.NewToken(token.RBRACKET, l.ch)
	case '/':
		tok = token.NewToken(token.SLASH, l.ch)
	case '*':
		tok = token.NewToken(token.ASTERISK, l.ch)
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.ReadChar()
			tok = token.Token{Type: token.NOT_EQ, Literal: string(ch) + string(l.ch)}
		} else {
			tok = token.NewToken(token.BANG, l.ch)
		}
	case '<':
		if l.peekChar() == '=' {
			ch := l.ch
			l.ReadChar()
			tok = token.Token{Type: token.LT_EQ, Literal: string(ch) + string(l.ch)}
		} else {
			tok = token.NewToken(token.LT, l.ch)
		}
	case '>':
		if l.peekChar() == '=' {
			ch := l.ch
			l.ReadChar()
			tok = token.Token{Type: token.GT_EQ, Literal: string(ch) + string(l.ch)}
		} else {
			tok = token.NewToken(token.GT, l.ch)
		}
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	case '"':
		tok.Type = token.STRING
		tok.Literal = l.readString()
	case ':':
		tok = token.NewToken(token.COLON, l.ch)
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			tok.Type = token.INT
			tok.Literal = l.readNumber()
			return tok
		} else {
			tok = token.NewToken(token.ILLEGAL, l.ch)
		}
	}

	l.ReadChar()
	return tok
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) {
		l.ReadChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.ReadChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.ReadChar()
	}
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func (l *Lexer) readString() string {
	position := l.position + 1
	for {
		l.ReadChar()
		if l.ch == '"' || l.ch == 0 {
			break
		}
	}
	return l.input[position:l.position]
}
