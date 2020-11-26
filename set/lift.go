package set

import "github.com/spoke-d/path"

func Lift(v interface{}) path.Scope {
	switch t := v.(type) {
	case map[string]interface{}:
		return MakeSet(t)
	case string:
		return path.MakeStringScope(t)
	}
	panic("missing type")
}
