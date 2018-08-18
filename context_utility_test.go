package flow

import "github.com/google/go-cmp/cmp"

var emptyContext = NewContext()

var contextEqualOpts = cmp.Comparer(func(x, y contextData) bool {
	return (x == nil && y == nil) || (x != nil && y != nil && cmp.Equal(x, y))
})
