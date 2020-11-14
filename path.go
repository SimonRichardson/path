package path

import (
	"fmt"

	"github.com/pkg/errors"
)

// Box represents a boxed datatype.
type Box interface {
	// Value defines the shadow type value of the Box.
	Value() interface{}
}

// Scope is used to identify a given expression of a global mutated object.
type Scope interface {
	// GetIdentValue returns the value of the identifier in a given scope.
	GetIdentValue(string) (Box, error)
}

// Path holds all the arguments for a given query.
type Path struct {
	ast *QueryExpression
}

// Parse attempts to parse a given query into a argument query.
// Returns an error if it's not in the correct layout.
func Parse(src string) (Path, error) {
	lex := NewLexer(src)
	parser := NewParser(lex)
	ast, err := parser.Run()
	if err != nil {
		return Path{}, errors.WithStack(err)
	}

	return Path{
		ast: ast,
	}, nil
}

// Run the query over a given scope.
func (q Path) Run(scope Scope) (interface{}, error) {
	box, err := q.run(q.ast, scope)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return box.Value(), nil
}

func (q Path) run(e Expression, scope Scope) (Box, error) {
	// Useful for debugging.
	// fmt.Printf("%[1]T %[1]v\n", e)

	switch node := e.(type) {
	case *QueryExpression:
		for _, exp := range node.Expressions {
			result, err := q.run(exp, scope)
			if err != nil {
				return nil, errors.WithStack(err)
			}
			if result != nil {
				return result, nil
			}
		}
		return nil, nil

	case *ExpressionStatement:
		return q.run(node.Expression, scope)

	case *Identifier:
		return scope.GetIdentValue(node.Token.Literal)

	case *AccessorExpression:
		parent, err := q.getName(node.Left, scope)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		child, err := q.getName(node.Right, scope)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		return scope.GetIdentValue(fmt.Sprintf("%s.%s", parent, child))
	}
	return nil, RuntimeErrorf("Syntax Error: Unexpected expression %T", e)
}

func (q *Path) getName(node Expression, scope Scope) (string, error) {
	parent, ok := node.(*Identifier)
	if ok {
		return parent.Token.Literal, nil
	}

	box, err := q.run(node, scope)
	if err != nil {
		return "", errors.WithStack(err)
	}
	b, ok := box.(Box)
	if !ok {
		return "", RuntimeErrorf("%T %v unexpected identifier", node, node.Pos())
	}
	raw, ok := b.Value().(string)
	if !ok {
		return "", RuntimeErrorf("%T %v unexpected name type", node, node.Pos())
	}
	return raw, nil
}
