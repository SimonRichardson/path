package set

import (
	"github.com/pkg/errors"
	"github.com/spoke-d/path"
)

// Set defines the type for querying from a path.
type Set struct {
	m map[string]interface{}
}

func MakeSet(m map[string]interface{}) Set {
	return Set{
		m: m,
	}
}

// GetAllIdents returns all the identifiers for a given scope.
func (s Set) GetAllIdents() []string {
	result := make([]string, 0, len(s.m))
	for k := range s.m {
		result = append(result, k)
	}
	return result
}

// GetIdentValue returns the value of the identifier in a given scope.
func (s Set) GetIdentValue(v string) (path.Scope, error) {
	if i, ok := s.m[v]; ok {
		switch t := i.(type) {
		case map[string]interface{}:
			return MakeSet(t), nil
		case string:
			return path.MakeStringScope(t), nil
		}
	}
	return nil, errors.Errorf("no ident value %q found in scope", v)
}

// RunOperation attempts to run an operation on a given scope
func (s Set) RunOperation(op path.Operation, scope path.Scope) (path.Scope, error) {
	result := make(map[string]interface{})
	for k, v := range s.m {
		if _, err := scope.RunOperation(op, Lift(v)); err == nil {
			result[k] = v
		}
	}
	return MakeSet(result), nil
}
