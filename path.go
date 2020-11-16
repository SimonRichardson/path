package path

import (
	"github.com/pkg/errors"
)

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
func (q Path) Run(scope Scope) (Scope, error) {
	result, err := q.run(q.ast, scope)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return result, nil
}

func (q Path) run(e Expression, scope Scope) (Scope, error) {
	// Useful for debugging.
	//fmt.Printf("%[1]T %[1]v\n", e)

	switch node := e.(type) {
	case *QueryExpression:
		var scopes []Scope
		for _, exp := range node.Expressions {
			result, err := q.run(exp, scope)
			if err != nil {
				return nil, errors.WithStack(err)
			}
			scopes = append(scopes, result)
		}

		return NewScopes(scopes), nil

	case *ExpressionStatement:
		return q.run(node.Expression, scope)

	case *Identifier:
		return scope.GetIdentValue(node.Token.Literal)

	case *String:
		return scope.GetIdentValue(node.Token.Literal)

	case *AccessorExpression:
		parent, err := q.run(node.Left, scope)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		return q.run(node.Right, parent)

	case *IndexExpression:
		left, err := q.run(node.Left, scope)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		return q.run(node.Index, left)
	}

	return nil, RuntimeErrorf("Syntax Error: Unexpected expression %T", e)
}
