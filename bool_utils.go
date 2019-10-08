package hoff

import "github.com/google/go-cmp/cmp"

var (
	truePointer  = newBool(true)
	falsePointer = newBool(false)
)

// BoolPointer give a fixed pointer to a corresponding bool value.
// e.g. true as value will always have the same pointer
func boolPointer(value bool) *bool {
	if value {
		return truePointer
	}
	return falsePointer
}

var (
	// errorComparator is a google/go-cmp comparator of errors
	errorComparator = cmp.Comparer(func(x, y error) bool {
		return (x == nil && y == nil) || (x != nil && y != nil && x.Error() == y.Error())
	})
)

func newBool(value bool) *bool {
	return &value
}
