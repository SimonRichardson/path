package path

import (
	"github.com/pkg/errors"
)

// Scope is used to identify a given expression of a global mutated object.
type Scope interface {
	// GetAllIdents returns all the identifiers for a given scope.
	GetAllIdents() []string
	// GetIdentValue returns the value of the identifier in a given scope.
	GetIdentValue(string) (Scope, error)
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
