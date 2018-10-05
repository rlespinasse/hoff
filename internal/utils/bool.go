/*
Package utils expose utility functions.
*/
package utils

import "github.com/google/go-cmp/cmp"

var (
	truePointer  = newBool(true)
	falsePointer = newBool(false)
)

// BoolPointer give a fixed pointer to a corresponding bool value.
// e.g. true as value will always have the same pointer
func BoolPointer(value bool) *bool {
	if value {
		return truePointer
	}
	return falsePointer
}

var (
	// ErrorComparator is a google/go-cmp comparator of errors
	ErrorComparator = cmp.Comparer(func(x, y error) bool {
		return (x == nil && y == nil) || (x != nil && y != nil && x.Error() == y.Error())
	})
)

func newBool(value bool) *bool {
	return &value
}
