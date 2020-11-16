package path

import (
	"strings"
	"text/scanner"
	"unicode"
	"unicode/utf8"
)

// Lexer takes a query and breaks it down into tokens that can be consumed at
// at later date.
// The lexer in question is lazy and requires the calling of next to move it
// forward.
type Lexer struct {
	scanner scanner.Scanner
	text    string
	isEOF   bool
}

// NewLexer creates a new Lexer from a given input.
func NewLexer(input string) *Lexer {
	var scanner scanner.Scanner
	scanner.Init(strings.NewReader(input))
	lex := &Lexer{
		scanner: scanner,
	}
	lex.ReadNext()
	return lex
}

// ReadNext will attempt to read the next character and correctly setup the
// positional values for the input.
func (l *Lexer) ReadNext() {
	if l.scanner.Scan() == scanner.EOF {
		l.isEOF = true
		l.text = ""
		return
	}
	l.text = l.scanner.TokenText()
}

// NextToken attempts to grab the next token available.
func (l *Lexer) NextToken() Token {
	defer l.ReadNext()

	var tok Token
	pos := l.getPosition()
	pos.Column--

	if t, ok := tokenMap[l.text]; ok {
		switch t {
		default:
			tok = MakeToken(t, l.text)
		}
		tok.Pos = pos
		return tok
	}

	newToken := l.readRunesToken()
	newToken.Pos = pos
	return newToken
}

func (l *Lexer) readRunesToken() Token {
	var tok Token
	switch {
	case l.isEOF:
		tok.Type = EOF
		return tok
	case len(l.text) > 0 && isLetter(l.text[0]):
		tok.Type = IDENT
		tok.Literal = l.text
		return tok
	case len(l.text) > 0 && isQuote(l.text[0]):
		tok.Type = STRING
		tok.Literal = l.text[1 : len(l.text)-1]
		return tok
	}

	return MakeToken(UNKNOWN, l.text)
}

func (l *Lexer) getPosition() Position {
	pos := l.scanner.Pos()
	return Position{
		Offset: pos.Offset,
		Line:   pos.Line,
		Column: pos.Column,
	}
}

func isLetter(char byte) bool {
	return 'a' <= char && char <= 'z' || 'A' <= char && char <= 'Z' || char == '_' || char >= utf8.RuneSelf && unicode.IsLetter(rune(char))
}

func isQuote(char byte) bool {
	return char == 34
}
