package path

import (
	"bytes"
	"strings"
	"unicode/utf8"
)

// Expression defines a type of AST node for outlining an expression.
type Expression interface {
	Pos() Position
	End() Position

	String() string
}

// QueryExpression represents a query full of expressions
type QueryExpression struct {
	Expressions []Expression
}

// Pos returns the first position of the query expression.
func (e *QueryExpression) Pos() Position {
	if len(e.Expressions) > 0 {
		return e.Expressions[0].Pos()
	}
	return Position{}
}

// End returns the last position of the query expression.
func (e *QueryExpression) End() Position {
	if num := len(e.Expressions); num > 0 {
		return e.Expressions[num-1].End()
	}
	return Position{}
}

func (e *QueryExpression) String() string {
	var out bytes.Buffer

	for _, s := range e.Expressions {
		out.WriteString(s.String())
	}

	return out.String()
}

// ExpressionStatement is a group of expressions that allows us to group a
// subset of expressions.
type ExpressionStatement struct {
	Token      Token
	Expression Expression
}

// Pos returns the first position of the expression statement.
func (es *ExpressionStatement) Pos() Position {
	return es.Token.Pos
}

// End returns the last position of the expression statement.
func (es *ExpressionStatement) End() Position {
	return es.Expression.End()
}

func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		str := es.Expression.String()
		if str == "" {
			return ";"
		}
		if str[len(str)-1:] != ";" {
			str += ";"
		}
		return strings.TrimSpace(str)
	}
	return ""
}

// InfixExpression represents an expression that is associated with an operator.
type InfixExpression struct {
	Token    Token
	Operator string
	Right    Expression
	Left     Expression
}

// Pos returns the first position of the identifier.
func (ie *InfixExpression) Pos() Position {
	return ie.Token.Pos
}

// End returns the last position of the identifier.
func (ie *InfixExpression) End() Position {
	return ie.Right.End()
}

func (ie *InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" " + ie.Operator + " ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")

	return out.String()
}

// AccessorExpression represents an expression that is associated with an operator.
type AccessorExpression struct {
	Token Token
	Right Expression
	Left  Expression
}

// Pos returns the first position of the identifier.
func (ie *AccessorExpression) Pos() Position {
	return ie.Token.Pos
}

// End returns the last position of the identifier.
func (ie *AccessorExpression) End() Position {
	return ie.Right.End()
}

func (ie *AccessorExpression) String() string {
	var out bytes.Buffer

	out.WriteString(ie.Left.String())
	out.WriteString(".")
	out.WriteString(ie.Right.String())

	return out.String()
}

// Identifier represents an identifier for a given AST block
type Identifier struct {
	Token Token
}

// Pos returns the first position of the identifier.
func (i *Identifier) Pos() Position {
	return i.Token.Pos
}

// End returns the last position of the identifier.
func (i *Identifier) End() Position {
	length := utf8.RuneCountInString(i.Token.Literal)
	return Position{
		Line:   i.Token.Pos.Line,
		Column: i.Token.Pos.Column + length,
	}
}

func (i *Identifier) String() string { return i.Token.Literal }
