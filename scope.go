package path

import (
	"github.com/pkg/errors"
)

// Operation describes an operation that should be run.
type Operation int

const (
	OpEQ Operation = iota
	OpNEQ
	OpLT
	OpLE
	OpGT
	OpGE
)

func liftOperation(token TokenType) (Operation, error) {
	switch token {
	case EQ:
		return OpEQ, nil
	case NEQ:
		return OpNEQ, nil
	case LT:
		return OpLT, nil
	case LE:
		return OpLE, nil
	case GT:
		return OpGT, nil
	case GE:
		return OpGE, nil
	}
	return -1, errors.Errorf("unexpected token type %q", token)
}

// Scope is used to identify a given expression of a global mutated object.
type Scope interface {
	// GetAllIdents returns all the identifiers for a given scope.
	GetAllIdents() []string
	// GetIdentValue returns the value of the identifier in a given scope.
	GetIdentValue(string) (Scope, error)
	// RunOperation attempts to run an operation on a given scope
	RunOperation(Operation, Scope) (Scope, error)
}

// Scopes holds a list of scopes to walk over.
type Scopes struct {
	scopes []Scope
}

// NewScopes creates a new set of scopes.
func NewScopes(scopes []Scope) *Scopes {
	return &Scopes{
		scopes: scopes,
	}
}

// GetIdentValue returns the value of the identifier in a given scope.
func (s Scopes) GetIdentValue(v string) (Scope, error) {
	for _, x := range s.scopes {
		r, err := x.GetIdentValue(v)
		if err != nil {
			continue
		}
		return r, nil
	}
	return nil, errors.Errorf("No ident value %q found in scope", v)
}

// GetAllIdents returns all the identifiers for a given scope.
func (s Scopes) GetAllIdents() []string {
	var res []string
	for _, x := range s.scopes {
		res = append(res, x.GetAllIdents()...)
	}
	return res
}

// RunOperation runs a given operation on a scope.
func (s Scopes) RunOperation(op Operation, other Scope) (Scope, error) {
	var lastErr error
	for _, scope := range s.scopes {
		res, err := scope.RunOperation(op, scope)
		if err != nil {
			lastErr = err
			continue
		}
		if res != nil {
			return res, nil
		}
	}
	return nil, errors.WithStack(lastErr)
}
