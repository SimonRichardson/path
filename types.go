package path

import "github.com/pkg/errors"

type StringScope struct {
	v string
}

func MakeStringScope(v string) StringScope {
	return StringScope{
		v: v,
	}
}

// GetAllIdents returns all the identifiers for a given scope.
func (s StringScope) GetAllIdents() []string {
	return make([]string, 0)
}

// GetIdentValue returns the value of the identifier in a given scope.
func (s StringScope) GetIdentValue(v string) (Scope, error) {
	return s, nil
}

// RunOperation attempts to run an operation on a given scope
func (s StringScope) RunOperation(op Operation, scope Scope) (Scope, error) {
	o, ok := scope.(StringScope)
	if !ok {
		return nil, errors.Errorf("invalid scope comparision")
	}

	switch op {
	case OpEQ:
		if s.v == o.v {
			return s, nil
		}
	case OpNEQ:
		if s.v != o.v {
			return s, nil
		}
	case OpLT:
		if s.v < o.v {
			return s, nil
		}
	case OpLE:
		if s.v <= o.v {
			return s, nil
		}
	case OpGT:
		if s.v > o.v {
			return s, nil
		}
	case OpGE:
		if s.v >= o.v {
			return s, nil
		}
	}

	return nil, errors.Errorf("no match")
}
