package flow

import (
	"github.com/google/go-cmp/cmp"
)

var computationEqualOpts = cmp.Comparer(func(x, y *computation) bool {
	return (x == nil && y == nil) || (x != nil && y != nil && *x == *y)
})
