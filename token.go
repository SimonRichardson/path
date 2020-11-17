package path

import "fmt"

// TokenType represents a way to identify an individual token.
type TokenType int

const (
	UNKNOWN TokenType = (iota - 1)
	EOF

	IDENT
	STRING

	BITAND  // &
	BITOR   // |
	CONDAND // &&
	CONDOR  // ||

	LPAREN   // (
	RPAREN   // )
	LBRACKET // [
	RBRACKET // ]

	PERIOD    // .
	SEMICOLON // ;
)

func (t TokenType) String() string {
	switch t {
	case EOF:
		return "<EOF>"
	case IDENT:
		return "<IDENT>"
	case STRING:
		return "<STRING>"
	case LPAREN:
		return "("
	case RPAREN:
		return ")"
	case LBRACKET:
		return "["
	case RBRACKET:
		return "]"
	case BITAND:
		return "&"
	case BITOR:
		return "|"
	case CONDAND:
		return "&&"
	case CONDOR:
		return "||"
	case PERIOD:
		return "."
	case SEMICOLON:
		return ";"
	default:
		return "<UNKNOWN>"
	}
}

// Position holds the location of the token within the query.
type Position struct {
	Offset int
	Line   int
	Column int
}

func (p Position) String() string {
	return fmt.Sprintf("<:%d:%d>", p.Line, p.Column)
}

// Token defines a token found with in a query, along with the position and what
// type it is.
type Token struct {
	Pos     Position
	Type    TokenType
	Literal string
}

// MakeToken creates a new token value.
func MakeToken(tokenType TokenType, char string) Token {
	return Token{
		Type:    tokenType,
		Literal: char,
	}
}

var (
	// UnknownToken can be used as a sentinel token for a unknown state.
	UnknownToken = Token{
		Type: UNKNOWN,
	}
)

var tokenMap = map[string]TokenType{
	";": SEMICOLON,
	".": PERIOD,
	"&": BITAND,
	"|": BITOR,
	"(": LPAREN,
	")": RPAREN,
	"[": LBRACKET,
	"]": RBRACKET,
}
