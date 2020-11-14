package path

import "github.com/pkg/errors"

// RuntimeError creates an invalid error.
type RuntimeError struct {
	err error
}

func (e *RuntimeError) Error() string {
	return e.err.Error()
}

// RuntimeErrorf defines a sentinel error for invalid index.
func RuntimeErrorf(msg string, args ...interface{}) error {
	return &RuntimeError{
		err: errors.Errorf("Runtime Error: "+msg, args...),
	}
}

// IsRuntimeError returns if the error is an ErrInvalidIndex error
func IsRuntimeError(err error) bool {
	err = errors.Cause(err)
	_, ok := err.(*RuntimeError)
	return ok
}
